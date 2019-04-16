// Released under MIT license.
// https://raw.githubusercontent.com/jung-kurt/caddy-pubsub/master/LICENSE

(function() {

  var sub, pub, subscribe, status, append, ids, mode;

  subscribe = null;
  ids = util.domElements('id');
  append = util.append(ids.list, 'li', 10);
  mode = "run";

  pub = function(str) {
    var body, authStr;
    if (ids.nm.value === '') ids.nm.value = 'user_' + util.randomInt(10, 100);
    body = 'Event ' + str + ' from publisher ' + ids.nm.value;
    authStr = util.basicAuthStr(ids.pub_user.value, ids.pub_pw.value);
    ps.publish(ids.pub_category.value, ids.pub_url.value, body, authStr);
  };

  sub = function() {
    var authStr, subFnc;

    subFnc = function(data, time) {
      append(data);
    };

    if (subscribe) {
      subscribe.stop();
      subscribe = null;
      ids.sub_action.textContent = "Start";
    } else {
      authStr = util.basicAuthStr(ids.sub_user.value, ids.sub_pw.value);
      subscribe = ps.Subscriber(ids.sub_category.value, ids.sub_url.value, subFnc, authStr);
      subscribe.start();
      ids.sub_action.textContent = "Stop";
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
      case "sub-action":
        sub();
        break;
      case "mode":
        if (mode === "run") {
          util.bankShow(ids.configure);
          ids.mode.textContent = "Run";
          mode = "configure";
        } else {
          util.bankShow(ids.run);
          ids.mode.textContent = "Configure";
          mode = "run";
        }
        break;
    }
  });

})();
