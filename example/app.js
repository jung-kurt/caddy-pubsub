(function() {

  var sub, pub, subscribe, status, append, ids;

  subscribe = null;

  ids = util.domElements('id');

  append = util.append(ids.list, 'li', 10);

  pub = function(str) {
    var body, authStr;
    if (ids.nm.value === '') ids.nm.value = 'user_' + util.randomInt(10, 100);
    body = 'Event ' + str + ' from publisher ' + ids.nm.value;
    authStr = util.basicAuthStr(ids.pub_user.value, ids.pub_pw.value);
    ps.publish(ids.pub_category.value, ids.pub_url.value, body, authStr);
  };

  status = function(active) {
    if (active) {
      ids.status.textContent = "Active";
      ids.btn_stop.removeAttribute('disabled');
      ids.btn_start.setAttribute('disabled', 'disabled');
    } else {
      ids.status.textContent = "Inactive";
      ids.btn_start.removeAttribute('disabled');
      ids.btn_stop.setAttribute('disabled', 'disabled');
    }
  };

  status(false);

  sub = function(activate) {
    var authStr, subFnc;

    subFnc = function(data, time) {
      append(data);
    };

    if (activate) {
      if (subscribe) subscribe.stop();
      authStr = util.basicAuthStr(ids.sub_user.value, ids.sub_pw.value);
      subscribe = ps.Subscriber(ids.sub_category.value, ids.sub_url.value, subFnc, authStr);
      subscribe.start();
      status(true);
    } else {
      if (subscribe) {
        subscribe.stop();
        subscribe = null;
      }
      status(false);
    }
  };

  util.click(function(val) {
    switch (val) {
      case "pub-a":
        pub("A");
        break;
      case "pub-b":
        pub("B");
        break;
      case "pub-c":
        pub("C");
        break;
      case "sub-start":
        sub(true);
        break;
      case "sub-stop":
        sub(false);
        break;
      case "configure":
        util.bankShow(ids.configure);
        break;
      case "run":
        util.bankShow(ids.run);
        break;
    }
  });

})();
