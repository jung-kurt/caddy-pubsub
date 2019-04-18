// Released under MIT license. This file is part of the project at
// https://github.com/jung-kurt/caddy-pubsub

var ps = ps || {};

(function() {

  var ajax, jsonDecode, mergeInto;

  // Copy (shallow) properties and values from objB into objA and return objA.
  mergeInto = function(objA, objB)  {
    var key;
    if (objB) {
      for (key in objB) {
        if (objB.hasOwnProperty(key)) {
          objA[key] = objB[key];
        }
      }
    }
    return objA;
  };

  // Parse the specified JSON record. If an error occurs, null is returned.
  jsonDecode = function(str) {
    try {
      return JSON.parse(str);
    } catch (error) {
      return null;
    }
  };

  // Return an AJAX request object. The specified optional success and error
  // callback are registered but the calls to 'open', 'send' and other methods
  // are left to the caller.
  ajax = function(success, error) {
    var httpRequest;
    httpRequest = new XMLHttpRequest();
    if (httpRequest) {
      httpRequest.onreadystatechange = function() {
        if (httpRequest.readyState === XMLHttpRequest.DONE) {
          if (httpRequest.status === 200) {
            if (success) success(httpRequest.responseText);
          } else if (error) error(httpRequest.status);
        }
      };
    }
    return httpRequest;
  };

  // Subscribe to events in the specified category.
  //
  // url corresponds to the subscription path specified in the server Caddyfile.
  //
  // fnc is called when an event is received from the server. Its first parameter
  // is the body that was sent to the server by the event publisher. The second
  // parameter is the event's Unix timestamp.
  //
  // authStr is the string value associated with the 'Authorization' header, for
  // example,
  //   'Bearer eyJhbGciOiJIU...otqvAerKI'
  // for JSON web token authentication, and
  //   'Basic QWxhZGRpbjpPcGVuU2VzYW1l'
  // for basic access authentication. No authentication header will be used if
  // authStr is missing or empty.
  //
  // options is an object that configures the subscription. Its fields are:
  // * 'timeout': number of seconds after which the longpoll will recycle,
  //   default 45
  // * 'successDelay': number of milliseconds to delay after an event is
  //   successfully received, default 3
  // * 'errorDelay': number of milliseconds to delay after an error occurs,
  //   default 3000
  // * 'json': boolean indicating whether event bodies are JSON-encoded and
  //   should be decoded automatically, default false
  //
  // The return object has the methods start(), isActive(), and stop().
  //
  // This function was adapted from github.com/jcuga/golongpoll, file
  // examples/basic/basic.go, MIT license
  ps.Subscriber = function(category, url, fnc, authStr, options) {

    if (!(this instanceof ps.Subscriber)) {
      return new ps.Subscriber(category, url, fnc, authStr, options);
    }

    var cycle, decode, active, timeoutId, opt, sinceTime, success, err, pollUrl, poll;

    opt = mergeInto({
      'timeout': 45,  // in seconds
      'successDelay': 10,  // milliseconds
      'errorDelay': 3000,  // milliseconds
      'json': false // do not automatically decode event bodies as JSON
    }, options);

    active = false;
    pollUrl = url + '?timeout=' + opt.timeout + '&category=' + category + '&since_time=';
    timeoutId = null;
    decode = opt.json;

    cycle = function(ok) {
      if (active) {
        timeoutId = window.setTimeout(poll, ok ? opt.successDelay : opt.errorDelay);
      }
    };

    success = function(data) {
      var ok;

      console.log('ajax success');

      ok = false;
      if (data) {
        data = jsonDecode(data);
        if (data) {
          if (data.events && data.events.length > 0) {
            data.events.forEach(function(evt) {
              var body;
              sinceTime = evt.timestamp;
              if (active) {
                if (decode) {
                  body = jsonDecode(evt.data);
                  if (body !== null) evt.data = body;
                }
                fnc(evt.data, sinceTime);
              }
            });
            ok = true;
          }
        }
      }
      cycle(ok);
    };

    err = function(code) {
      console.log('ajax error', code);
      cycle(false);
    };

    poll = function() {
      var req;

      req = ajax(success, err);
      req.open('GET', pollUrl + sinceTime);
      if (authStr && authStr.length > 0) {
        req.setRequestHeader('Authorization', authStr);
      }
      req.send();
    };

    this.start = function() {
      this.stop();
      sinceTime = (new Date(Date.now())).getTime();
      active = true;
      poll();
    };

    this.stop = function() {
      if (typeof timeoutId === "number") {
        window.clearTimeout(timeoutId);
        timeoutId = null;
      }
      active = false;
    };

    this.isActive = function() {
      return active;
    };

  };

  // Publish an event in the category specified by categoryStr. The event will be
  // dispatched to all subscribers of the specified category. urlStr corresponds
  // to the publish path in the Caddyfile. bodyStr is the event body that will be
  // sent to subscribers. authStr is the string value associated with the
  // 'Authorization' header, for example,
  //   'Bearer eyJhbGciOiJIU...otqvAerKI'
  // for JSON web token authentication, and 
  //   'Basic QWxhZGRpbjpPcGVuU2VzYW1l' 
  // for basic access authentication. No authentication header will be used if
  // authStr is missing or empty.
  ps.publish = function(categoryStr, urlStr, bodyStr, authStr) {
    var req;

    req = ajax();
    req.open('POST', urlStr);
    req.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    if (authStr && authStr.length > 0) {
      req.setRequestHeader('Authorization', authStr);
    }
    req.send('category=' + encodeURIComponent(categoryStr) + '&body=' + 
      encodeURIComponent(bodyStr));
  };

  // This function is similar to ps.publish except that the third argument will
  // be JSON-encoded.
  ps.publishObj = function(categoryStr, urlStr, bodyObj, authStr) {
    ps.publish(categoryStr, urlStr, JSON.stringify(bodyObj), authStr);
  };

})();
