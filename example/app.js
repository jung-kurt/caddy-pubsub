// Released under MIT license. This file is part of the project at
// https://github.com/jung-kurt/caddy-pubsub

(function() {

  var auto, autocycle, autoc, sub, pub, subscribe, append, els, mode;

  subscribe = null;
  els = util.domElements('id', 'data-click');
  append = util.append(els.list, 'li', 10);
  mode = "run";
  auto = false;
  autoc = 'A';

  pub = function(str) {
    var body, authStr;
    if (els.nm.value === '') els.nm.value = 'user_' + util.randomInt(10, 100);
    body = { 'event': str, 'publisher': els.nm.value };
    authStr = util.basicAuthStr(els.pub_user.value, els.pub_pw.value);
    ps.publishObj(els.pub_category.value, els.pub_url.value, body, authStr);
  };

  sub = function() {
    var authStr, subFnc;

    subFnc = function(data, time) {
      console.log(data);
      append('Event ' + data.event + ' from publisher ' + data.publisher + '.');
    };

    if (subscribe) {
      subscribe.stop();
      subscribe = null;
      els.sub_action.textContent = "Start";
      els.status.textContent = "Inactive";
    } else {
      authStr = util.basicAuthStr(els.sub_user.value, els.sub_pw.value);
      subscribe = ps.Subscriber(els.sub_category.value, els.sub_url.value, subFnc, authStr, { 'json': true });
      subscribe.start();
      els.sub_action.textContent = "Stop";
      els.status.textContent = "Active";
    }

  };

  (function() {
    var cycle;

    cycle = function() {
      if (auto) {
        pub(autoc);
        autoc = util.nextChar(autoc);
        window.setTimeout(cycle, util.randomInt(2, 5) * 1000);
      }
    };

    autocycle = function() {
      auto = !auto;
      if (auto) {
        els.auto.textContent = 'Auto stop';
        cycle();
      } else {
        els.auto.textContent = 'Auto start';
      }
    };

  })();

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
      case "auto":
        autocycle();
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
