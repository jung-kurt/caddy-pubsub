// Released under MIT license.
// https://raw.githubusercontent.com/jung-kurt/caddy-pubsub/master/LICENSE

(function() {

  var sub, pub, subscribe, status, append, els, mode;

  subscribe = null;
  els = util.domElements('id', 'data-click');
  append = util.append(els.list, 'li', 10);
  mode = "run";

  pub = function(str) {
    var body, authStr;
    if (els.nm.value === '') els.nm.value = 'user_' + util.randomInt(10, 100);
    body = 'Event ' + str + ' from publisher ' + els.nm.value;
    authStr = util.basicAuthStr(els.pub_user.value, els.pub_pw.value);
    ps.publish(els.pub_category.value, els.pub_url.value, body, authStr);
  };

  sub = function() {
    var authStr, subFnc;

    subFnc = function(data, time) {
      append(data);
    };

    if (subscribe) {
      subscribe.stop();
      subscribe = null;
      els.sub_action.textContent = "Start";
      els.status.textContent = "Inactive";
    } else {
      authStr = util.basicAuthStr(els.sub_user.value, els.sub_pw.value);
      subscribe = ps.Subscriber(els.sub_category.value, els.sub_url.value, subFnc, authStr);
      subscribe.start();
      els.sub_action.textContent = "Stop";
      els.status.textContent = "Active";
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
      case "sub_action":
        sub();
        break;
      case "mode":
        if (mode === "run") {
          util.bankShow(els.configure);
          els.mode.textContent = "Run";
          mode = "configure";
        } else {
          util.bankShow(els.run);
          els.mode.textContent = "Configure";
          mode = "run";
        }
        break;
    }
  });

})();
