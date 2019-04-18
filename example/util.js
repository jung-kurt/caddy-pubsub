// Released under MIT license. This file is part of the project at
// https://github.com/jung-kurt/caddy-pubsub

var util = util || {};

// Return object that associates, for each DOM element that has at least one of
// the attributes passed as an argument, the value of the attribute and the
// element itself. For example, els = util.DomElements('id', 'data-click').
util.domElements = function() {
  var attr, obj, el, j, k, list;

  obj = {};
  for (j = 0; j < arguments.length; j++) {
    attr = arguments[j];
    list = document.querySelectorAll('[' + attr + ']');
    for (k = 0; k < list.length; k++) {
      el = list[k];
      obj[el.getAttribute(attr)] = el;
    }
  }
  return obj;
};

// Set specified element as visible and each of its siblings as hidden.
util.bankShow = function(el) {
  var list, sib, j;
  list = el.parentNode.children;
  if (list) {
    for (j = 0; j < list.length; j++) {
      sib = list[j];
      if (sib.isEqualNode(el)) {
        el.style.display = '';
      } else {
        sib.style.display = 'none';
      }
    }
  }
};

// Call fnc when any element that has the data-click attribute is clicked. The
// function is called with the value associated with data-click.
util.click = function(fnc) {
  document.defaultView.addEventListener('click', function(evt) {
    var val, el;
    el = evt.target;
    val = el.getAttribute('data-click');
    if (val) {
      fnc(val);
    }
  });
};

// Return pseudo-random integer in [min,max).
util.randomInt = function(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min)) + min;
};

// Return the value to be associated with the 'Authorization' header for 
// basic access authentication.
util.basicAuthStr = function(user, password) {
  return 'Basic ' + window.btoa(user + ':' + password);
};

// Return a function that can be called to append a child element of type
// childTag (for example, 'li') to the specified parent element. After count
// child elements have been appended, the collection is constrained and each
// further append shifts off the oldest child to make room for the new child.
// This function efficiently reuses DOM elements once the maximum is reached.
// The function that is returned is called with a single string.
util.append = function(parentElement, childTag, count) {
  var list, append;
  list = [];
  append = function(str) {
    var child;
    child = (list.length < count) ? document.createElement(childTag) : list.shift();
    list.push(child);
    parentElement.appendChild(child);
    child.textContent = str;
  };
  return append;
};
