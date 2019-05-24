(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
// Copyright Joyent, Inc. and other Node contributors.
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to permit
// persons to whom the Software is furnished to do so, subject to the
// following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN
// NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
// USE OR OTHER DEALINGS IN THE SOFTWARE.

var objectCreate = Object.create || objectCreatePolyfill
var objectKeys = Object.keys || objectKeysPolyfill
var bind = Function.prototype.bind || functionBindPolyfill

function EventEmitter() {
  if (!this._events || !Object.prototype.hasOwnProperty.call(this, '_events')) {
    this._events = objectCreate(null);
    this._eventsCount = 0;
  }

  this._maxListeners = this._maxListeners || undefined;
}
module.exports = EventEmitter;

// Backwards-compat with node 0.10.x
EventEmitter.EventEmitter = EventEmitter;

EventEmitter.prototype._events = undefined;
EventEmitter.prototype._maxListeners = undefined;

// By default EventEmitters will print a warning if more than 10 listeners are
// added to it. This is a useful default which helps finding memory leaks.
var defaultMaxListeners = 10;

var hasDefineProperty;
try {
  var o = {};
  if (Object.defineProperty) Object.defineProperty(o, 'x', { value: 0 });
  hasDefineProperty = o.x === 0;
} catch (err) { hasDefineProperty = false }
if (hasDefineProperty) {
  Object.defineProperty(EventEmitter, 'defaultMaxListeners', {
    enumerable: true,
    get: function() {
      return defaultMaxListeners;
    },
    set: function(arg) {
      // check whether the input is a positive number (whose value is zero or
      // greater and not a NaN).
      if (typeof arg !== 'number' || arg < 0 || arg !== arg)
        throw new TypeError('"defaultMaxListeners" must be a positive number');
      defaultMaxListeners = arg;
    }
  });
} else {
  EventEmitter.defaultMaxListeners = defaultMaxListeners;
}

// Obviously not all Emitters should be limited to 10. This function allows
// that to be increased. Set to zero for unlimited.
EventEmitter.prototype.setMaxListeners = function setMaxListeners(n) {
  if (typeof n !== 'number' || n < 0 || isNaN(n))
    throw new TypeError('"n" argument must be a positive number');
  this._maxListeners = n;
  return this;
};

function $getMaxListeners(that) {
  if (that._maxListeners === undefined)
    return EventEmitter.defaultMaxListeners;
  return that._maxListeners;
}

EventEmitter.prototype.getMaxListeners = function getMaxListeners() {
  return $getMaxListeners(this);
};

// These standalone emit* functions are used to optimize calling of event
// handlers for fast cases because emit() itself often has a variable number of
// arguments and can be deoptimized because of that. These functions always have
// the same number of arguments and thus do not get deoptimized, so the code
// inside them can execute faster.
function emitNone(handler, isFn, self) {
  if (isFn)
    handler.call(self);
  else {
    var len = handler.length;
    var listeners = arrayClone(handler, len);
    for (var i = 0; i < len; ++i)
      listeners[i].call(self);
  }
}
function emitOne(handler, isFn, self, arg1) {
  if (isFn)
    handler.call(self, arg1);
  else {
    var len = handler.length;
    var listeners = arrayClone(handler, len);
    for (var i = 0; i < len; ++i)
      listeners[i].call(self, arg1);
  }
}
function emitTwo(handler, isFn, self, arg1, arg2) {
  if (isFn)
    handler.call(self, arg1, arg2);
  else {
    var len = handler.length;
    var listeners = arrayClone(handler, len);
    for (var i = 0; i < len; ++i)
      listeners[i].call(self, arg1, arg2);
  }
}
function emitThree(handler, isFn, self, arg1, arg2, arg3) {
  if (isFn)
    handler.call(self, arg1, arg2, arg3);
  else {
    var len = handler.length;
    var listeners = arrayClone(handler, len);
    for (var i = 0; i < len; ++i)
      listeners[i].call(self, arg1, arg2, arg3);
  }
}

function emitMany(handler, isFn, self, args) {
  if (isFn)
    handler.apply(self, args);
  else {
    var len = handler.length;
    var listeners = arrayClone(handler, len);
    for (var i = 0; i < len; ++i)
      listeners[i].apply(self, args);
  }
}

EventEmitter.prototype.emit = function emit(type) {
  var er, handler, len, args, i, events;
  var doError = (type === 'error');

  events = this._events;
  if (events)
    doError = (doError && events.error == null);
  else if (!doError)
    return false;

  // If there is no 'error' event listener then throw.
  if (doError) {
    if (arguments.length > 1)
      er = arguments[1];
    if (er instanceof Error) {
      throw er; // Unhandled 'error' event
    } else {
      // At least give some kind of context to the user
      var err = new Error('Unhandled "error" event. (' + er + ')');
      err.context = er;
      throw err;
    }
    return false;
  }

  handler = events[type];

  if (!handler)
    return false;

  var isFn = typeof handler === 'function';
  len = arguments.length;
  switch (len) {
      // fast cases
    case 1:
      emitNone(handler, isFn, this);
      break;
    case 2:
      emitOne(handler, isFn, this, arguments[1]);
      break;
    case 3:
      emitTwo(handler, isFn, this, arguments[1], arguments[2]);
      break;
    case 4:
      emitThree(handler, isFn, this, arguments[1], arguments[2], arguments[3]);
      break;
      // slower
    default:
      args = new Array(len - 1);
      for (i = 1; i < len; i++)
        args[i - 1] = arguments[i];
      emitMany(handler, isFn, this, args);
  }

  return true;
};

function _addListener(target, type, listener, prepend) {
  var m;
  var events;
  var existing;

  if (typeof listener !== 'function')
    throw new TypeError('"listener" argument must be a function');

  events = target._events;
  if (!events) {
    events = target._events = objectCreate(null);
    target._eventsCount = 0;
  } else {
    // To avoid recursion in the case that type === "newListener"! Before
    // adding it to the listeners, first emit "newListener".
    if (events.newListener) {
      target.emit('newListener', type,
          listener.listener ? listener.listener : listener);

      // Re-assign `events` because a newListener handler could have caused the
      // this._events to be assigned to a new object
      events = target._events;
    }
    existing = events[type];
  }

  if (!existing) {
    // Optimize the case of one listener. Don't need the extra array object.
    existing = events[type] = listener;
    ++target._eventsCount;
  } else {
    if (typeof existing === 'function') {
      // Adding the second element, need to change to array.
      existing = events[type] =
          prepend ? [listener, existing] : [existing, listener];
    } else {
      // If we've already got an array, just append.
      if (prepend) {
        existing.unshift(listener);
      } else {
        existing.push(listener);
      }
    }

    // Check for listener leak
    if (!existing.warned) {
      m = $getMaxListeners(target);
      if (m && m > 0 && existing.length > m) {
        existing.warned = true;
        var w = new Error('Possible EventEmitter memory leak detected. ' +
            existing.length + ' "' + String(type) + '" listeners ' +
            'added. Use emitter.setMaxListeners() to ' +
            'increase limit.');
        w.name = 'MaxListenersExceededWarning';
        w.emitter = target;
        w.type = type;
        w.count = existing.length;
        if (typeof console === 'object' && console.warn) {
          console.warn('%s: %s', w.name, w.message);
        }
      }
    }
  }

  return target;
}

EventEmitter.prototype.addListener = function addListener(type, listener) {
  return _addListener(this, type, listener, false);
};

EventEmitter.prototype.on = EventEmitter.prototype.addListener;

EventEmitter.prototype.prependListener =
    function prependListener(type, listener) {
      return _addListener(this, type, listener, true);
    };

function onceWrapper() {
  if (!this.fired) {
    this.target.removeListener(this.type, this.wrapFn);
    this.fired = true;
    switch (arguments.length) {
      case 0:
        return this.listener.call(this.target);
      case 1:
        return this.listener.call(this.target, arguments[0]);
      case 2:
        return this.listener.call(this.target, arguments[0], arguments[1]);
      case 3:
        return this.listener.call(this.target, arguments[0], arguments[1],
            arguments[2]);
      default:
        var args = new Array(arguments.length);
        for (var i = 0; i < args.length; ++i)
          args[i] = arguments[i];
        this.listener.apply(this.target, args);
    }
  }
}

function _onceWrap(target, type, listener) {
  var state = { fired: false, wrapFn: undefined, target: target, type: type, listener: listener };
  var wrapped = bind.call(onceWrapper, state);
  wrapped.listener = listener;
  state.wrapFn = wrapped;
  return wrapped;
}

EventEmitter.prototype.once = function once(type, listener) {
  if (typeof listener !== 'function')
    throw new TypeError('"listener" argument must be a function');
  this.on(type, _onceWrap(this, type, listener));
  return this;
};

EventEmitter.prototype.prependOnceListener =
    function prependOnceListener(type, listener) {
      if (typeof listener !== 'function')
        throw new TypeError('"listener" argument must be a function');
      this.prependListener(type, _onceWrap(this, type, listener));
      return this;
    };

// Emits a 'removeListener' event if and only if the listener was removed.
EventEmitter.prototype.removeListener =
    function removeListener(type, listener) {
      var list, events, position, i, originalListener;

      if (typeof listener !== 'function')
        throw new TypeError('"listener" argument must be a function');

      events = this._events;
      if (!events)
        return this;

      list = events[type];
      if (!list)
        return this;

      if (list === listener || list.listener === listener) {
        if (--this._eventsCount === 0)
          this._events = objectCreate(null);
        else {
          delete events[type];
          if (events.removeListener)
            this.emit('removeListener', type, list.listener || listener);
        }
      } else if (typeof list !== 'function') {
        position = -1;

        for (i = list.length - 1; i >= 0; i--) {
          if (list[i] === listener || list[i].listener === listener) {
            originalListener = list[i].listener;
            position = i;
            break;
          }
        }

        if (position < 0)
          return this;

        if (position === 0)
          list.shift();
        else
          spliceOne(list, position);

        if (list.length === 1)
          events[type] = list[0];

        if (events.removeListener)
          this.emit('removeListener', type, originalListener || listener);
      }

      return this;
    };

EventEmitter.prototype.removeAllListeners =
    function removeAllListeners(type) {
      var listeners, events, i;

      events = this._events;
      if (!events)
        return this;

      // not listening for removeListener, no need to emit
      if (!events.removeListener) {
        if (arguments.length === 0) {
          this._events = objectCreate(null);
          this._eventsCount = 0;
        } else if (events[type]) {
          if (--this._eventsCount === 0)
            this._events = objectCreate(null);
          else
            delete events[type];
        }
        return this;
      }

      // emit removeListener for all listeners on all events
      if (arguments.length === 0) {
        var keys = objectKeys(events);
        var key;
        for (i = 0; i < keys.length; ++i) {
          key = keys[i];
          if (key === 'removeListener') continue;
          this.removeAllListeners(key);
        }
        this.removeAllListeners('removeListener');
        this._events = objectCreate(null);
        this._eventsCount = 0;
        return this;
      }

      listeners = events[type];

      if (typeof listeners === 'function') {
        this.removeListener(type, listeners);
      } else if (listeners) {
        // LIFO order
        for (i = listeners.length - 1; i >= 0; i--) {
          this.removeListener(type, listeners[i]);
        }
      }

      return this;
    };

function _listeners(target, type, unwrap) {
  var events = target._events;

  if (!events)
    return [];

  var evlistener = events[type];
  if (!evlistener)
    return [];

  if (typeof evlistener === 'function')
    return unwrap ? [evlistener.listener || evlistener] : [evlistener];

  return unwrap ? unwrapListeners(evlistener) : arrayClone(evlistener, evlistener.length);
}

EventEmitter.prototype.listeners = function listeners(type) {
  return _listeners(this, type, true);
};

EventEmitter.prototype.rawListeners = function rawListeners(type) {
  return _listeners(this, type, false);
};

EventEmitter.listenerCount = function(emitter, type) {
  if (typeof emitter.listenerCount === 'function') {
    return emitter.listenerCount(type);
  } else {
    return listenerCount.call(emitter, type);
  }
};

EventEmitter.prototype.listenerCount = listenerCount;
function listenerCount(type) {
  var events = this._events;

  if (events) {
    var evlistener = events[type];

    if (typeof evlistener === 'function') {
      return 1;
    } else if (evlistener) {
      return evlistener.length;
    }
  }

  return 0;
}

EventEmitter.prototype.eventNames = function eventNames() {
  return this._eventsCount > 0 ? Reflect.ownKeys(this._events) : [];
};

// About 1.5x faster than the two-arg version of Array#splice().
function spliceOne(list, index) {
  for (var i = index, k = i + 1, n = list.length; k < n; i += 1, k += 1)
    list[i] = list[k];
  list.pop();
}

function arrayClone(arr, n) {
  var copy = new Array(n);
  for (var i = 0; i < n; ++i)
    copy[i] = arr[i];
  return copy;
}

function unwrapListeners(arr) {
  var ret = new Array(arr.length);
  for (var i = 0; i < ret.length; ++i) {
    ret[i] = arr[i].listener || arr[i];
  }
  return ret;
}

function objectCreatePolyfill(proto) {
  var F = function() {};
  F.prototype = proto;
  return new F;
}
function objectKeysPolyfill(obj) {
  var keys = [];
  for (var k in obj) if (Object.prototype.hasOwnProperty.call(obj, k)) {
    keys.push(k);
  }
  return k;
}
function functionBindPolyfill(context) {
  var fn = this;
  return function () {
    return fn.apply(context, arguments);
  };
}

},{}],2:[function(require,module,exports){
"use strict";

// Based on https://github.com/shtylman/node-process
var EventEmitter = require('events');

var process = module.exports = {};
process.nextTick = Script.nextTick;
process.title = 'Frida';
process.browser = true;
process.env = {};
process.argv = [];
process.version = ''; // empty string to avoid regexp issues

process.versions = {};
process.EventEmitter = EventEmitter;
process.on = noop;
process.addListener = noop;
process.once = noop;
process.off = noop;
process.removeListener = noop;
process.removeAllListeners = noop;
process.emit = noop;

process.binding = function (name) {
  throw new Error('process.binding is not supported');
};

process.cwd = function () {
  return '/';
};

process.chdir = function (dir) {
  throw new Error('process.chdir is not supported');
};

process.umask = function () {
  return 0;
};

function noop() {}

},{"events":1}],3:[function(require,module,exports){
(function (process){
// .dirname, .basename, and .extname methods are extracted from Node.js v8.11.1,
// backported and transplited with Babel, with backwards-compat fixes

// Copyright Joyent, Inc. and other Node contributors.
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to permit
// persons to whom the Software is furnished to do so, subject to the
// following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN
// NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
// USE OR OTHER DEALINGS IN THE SOFTWARE.

// resolves . and .. elements in a path array with directory names there
// must be no slashes, empty elements, or device names (c:\) in the array
// (so also no leading and trailing slashes - it does not distinguish
// relative and absolute paths)
function normalizeArray(parts, allowAboveRoot) {
  // if the path tries to go above the root, `up` ends up > 0
  var up = 0;
  for (var i = parts.length - 1; i >= 0; i--) {
    var last = parts[i];
    if (last === '.') {
      parts.splice(i, 1);
    } else if (last === '..') {
      parts.splice(i, 1);
      up++;
    } else if (up) {
      parts.splice(i, 1);
      up--;
    }
  }

  // if the path is allowed to go above the root, restore leading ..s
  if (allowAboveRoot) {
    for (; up--; up) {
      parts.unshift('..');
    }
  }

  return parts;
}

// path.resolve([from ...], to)
// posix version
exports.resolve = function() {
  var resolvedPath = '',
      resolvedAbsolute = false;

  for (var i = arguments.length - 1; i >= -1 && !resolvedAbsolute; i--) {
    var path = (i >= 0) ? arguments[i] : process.cwd();

    // Skip empty and invalid entries
    if (typeof path !== 'string') {
      throw new TypeError('Arguments to path.resolve must be strings');
    } else if (!path) {
      continue;
    }

    resolvedPath = path + '/' + resolvedPath;
    resolvedAbsolute = path.charAt(0) === '/';
  }

  // At this point the path should be resolved to a full absolute path, but
  // handle relative paths to be safe (might happen when process.cwd() fails)

  // Normalize the path
  resolvedPath = normalizeArray(filter(resolvedPath.split('/'), function(p) {
    return !!p;
  }), !resolvedAbsolute).join('/');

  return ((resolvedAbsolute ? '/' : '') + resolvedPath) || '.';
};

// path.normalize(path)
// posix version
exports.normalize = function(path) {
  var isAbsolute = exports.isAbsolute(path),
      trailingSlash = substr(path, -1) === '/';

  // Normalize the path
  path = normalizeArray(filter(path.split('/'), function(p) {
    return !!p;
  }), !isAbsolute).join('/');

  if (!path && !isAbsolute) {
    path = '.';
  }
  if (path && trailingSlash) {
    path += '/';
  }

  return (isAbsolute ? '/' : '') + path;
};

// posix version
exports.isAbsolute = function(path) {
  return path.charAt(0) === '/';
};

// posix version
exports.join = function() {
  var paths = Array.prototype.slice.call(arguments, 0);
  return exports.normalize(filter(paths, function(p, index) {
    if (typeof p !== 'string') {
      throw new TypeError('Arguments to path.join must be strings');
    }
    return p;
  }).join('/'));
};


// path.relative(from, to)
// posix version
exports.relative = function(from, to) {
  from = exports.resolve(from).substr(1);
  to = exports.resolve(to).substr(1);

  function trim(arr) {
    var start = 0;
    for (; start < arr.length; start++) {
      if (arr[start] !== '') break;
    }

    var end = arr.length - 1;
    for (; end >= 0; end--) {
      if (arr[end] !== '') break;
    }

    if (start > end) return [];
    return arr.slice(start, end - start + 1);
  }

  var fromParts = trim(from.split('/'));
  var toParts = trim(to.split('/'));

  var length = Math.min(fromParts.length, toParts.length);
  var samePartsLength = length;
  for (var i = 0; i < length; i++) {
    if (fromParts[i] !== toParts[i]) {
      samePartsLength = i;
      break;
    }
  }

  var outputParts = [];
  for (var i = samePartsLength; i < fromParts.length; i++) {
    outputParts.push('..');
  }

  outputParts = outputParts.concat(toParts.slice(samePartsLength));

  return outputParts.join('/');
};

exports.sep = '/';
exports.delimiter = ':';

exports.dirname = function (path) {
  if (typeof path !== 'string') path = path + '';
  if (path.length === 0) return '.';
  var code = path.charCodeAt(0);
  var hasRoot = code === 47 /*/*/;
  var end = -1;
  var matchedSlash = true;
  for (var i = path.length - 1; i >= 1; --i) {
    code = path.charCodeAt(i);
    if (code === 47 /*/*/) {
        if (!matchedSlash) {
          end = i;
          break;
        }
      } else {
      // We saw the first non-path separator
      matchedSlash = false;
    }
  }

  if (end === -1) return hasRoot ? '/' : '.';
  if (hasRoot && end === 1) {
    // return '//';
    // Backwards-compat fix:
    return '/';
  }
  return path.slice(0, end);
};

function basename(path) {
  if (typeof path !== 'string') path = path + '';

  var start = 0;
  var end = -1;
  var matchedSlash = true;
  var i;

  for (i = path.length - 1; i >= 0; --i) {
    if (path.charCodeAt(i) === 47 /*/*/) {
        // If we reached a path separator that was not part of a set of path
        // separators at the end of the string, stop now
        if (!matchedSlash) {
          start = i + 1;
          break;
        }
      } else if (end === -1) {
      // We saw the first non-path separator, mark this as the end of our
      // path component
      matchedSlash = false;
      end = i + 1;
    }
  }

  if (end === -1) return '';
  return path.slice(start, end);
}

// Uses a mixed approach for backwards-compatibility, as ext behavior changed
// in new Node.js versions, so only basename() above is backported here
exports.basename = function (path, ext) {
  var f = basename(path);
  if (ext && f.substr(-1 * ext.length) === ext) {
    f = f.substr(0, f.length - ext.length);
  }
  return f;
};

exports.extname = function (path) {
  if (typeof path !== 'string') path = path + '';
  var startDot = -1;
  var startPart = 0;
  var end = -1;
  var matchedSlash = true;
  // Track the state of characters (if any) we see before our first dot and
  // after any path separator we find
  var preDotState = 0;
  for (var i = path.length - 1; i >= 0; --i) {
    var code = path.charCodeAt(i);
    if (code === 47 /*/*/) {
        // If we reached a path separator that was not part of a set of path
        // separators at the end of the string, stop now
        if (!matchedSlash) {
          startPart = i + 1;
          break;
        }
        continue;
      }
    if (end === -1) {
      // We saw the first non-path separator, mark this as the end of our
      // extension
      matchedSlash = false;
      end = i + 1;
    }
    if (code === 46 /*.*/) {
        // If this is our first dot, mark it as the start of our extension
        if (startDot === -1)
          startDot = i;
        else if (preDotState !== 1)
          preDotState = 1;
    } else if (startDot !== -1) {
      // We saw a non-dot and non-path separator before our dot, so we should
      // have a good chance at having a non-empty extension
      preDotState = -1;
    }
  }

  if (startDot === -1 || end === -1 ||
      // We saw a non-dot character immediately before the dot
      preDotState === 0 ||
      // The (right-most) trimmed path component is exactly '..'
      preDotState === 1 && startDot === end - 1 && startDot === startPart + 1) {
    return '';
  }
  return path.slice(startDot, end);
};

function filter (xs, f) {
    if (xs.filter) return xs.filter(f);
    var res = [];
    for (var i = 0; i < xs.length; i++) {
        if (f(xs[i], i, xs)) res.push(xs[i]);
    }
    return res;
}

// String.prototype.substr - negative index don't work in IE8
var substr = 'ab'.substr(-1) === 'b'
    ? function (str, start, len) { return str.substr(start, len) }
    : function (str, start, len) {
        if (start < 0) start = str.length + start;
        return str.substr(start, len);
    }
;

}).call(this,require('_process'))

},{"_process":2}],4:[function(require,module,exports){
// shim for using process in browser
var process = module.exports = {};

// cached from whatever global is present so that test runners that stub it
// don't break things.  But we need to wrap it in a try catch in case it is
// wrapped in strict mode code which doesn't define any globals.  It's inside a
// function because try/catches deoptimize in certain engines.

var cachedSetTimeout;
var cachedClearTimeout;

function defaultSetTimout() {
    throw new Error('setTimeout has not been defined');
}
function defaultClearTimeout () {
    throw new Error('clearTimeout has not been defined');
}
(function () {
    try {
        if (typeof setTimeout === 'function') {
            cachedSetTimeout = setTimeout;
        } else {
            cachedSetTimeout = defaultSetTimout;
        }
    } catch (e) {
        cachedSetTimeout = defaultSetTimout;
    }
    try {
        if (typeof clearTimeout === 'function') {
            cachedClearTimeout = clearTimeout;
        } else {
            cachedClearTimeout = defaultClearTimeout;
        }
    } catch (e) {
        cachedClearTimeout = defaultClearTimeout;
    }
} ())
function runTimeout(fun) {
    if (cachedSetTimeout === setTimeout) {
        //normal enviroments in sane situations
        return setTimeout(fun, 0);
    }
    // if setTimeout wasn't available but was latter defined
    if ((cachedSetTimeout === defaultSetTimout || !cachedSetTimeout) && setTimeout) {
        cachedSetTimeout = setTimeout;
        return setTimeout(fun, 0);
    }
    try {
        // when when somebody has screwed with setTimeout but no I.E. maddness
        return cachedSetTimeout(fun, 0);
    } catch(e){
        try {
            // When we are in I.E. but the script has been evaled so I.E. doesn't trust the global object when called normally
            return cachedSetTimeout.call(null, fun, 0);
        } catch(e){
            // same as above but when it's a version of I.E. that must have the global object for 'this', hopfully our context correct otherwise it will throw a global error
            return cachedSetTimeout.call(this, fun, 0);
        }
    }


}
function runClearTimeout(marker) {
    if (cachedClearTimeout === clearTimeout) {
        //normal enviroments in sane situations
        return clearTimeout(marker);
    }
    // if clearTimeout wasn't available but was latter defined
    if ((cachedClearTimeout === defaultClearTimeout || !cachedClearTimeout) && clearTimeout) {
        cachedClearTimeout = clearTimeout;
        return clearTimeout(marker);
    }
    try {
        // when when somebody has screwed with setTimeout but no I.E. maddness
        return cachedClearTimeout(marker);
    } catch (e){
        try {
            // When we are in I.E. but the script has been evaled so I.E. doesn't  trust the global object when called normally
            return cachedClearTimeout.call(null, marker);
        } catch (e){
            // same as above but when it's a version of I.E. that must have the global object for 'this', hopfully our context correct otherwise it will throw a global error.
            // Some versions of I.E. have different rules for clearTimeout vs setTimeout
            return cachedClearTimeout.call(this, marker);
        }
    }



}
var queue = [];
var draining = false;
var currentQueue;
var queueIndex = -1;

function cleanUpNextTick() {
    if (!draining || !currentQueue) {
        return;
    }
    draining = false;
    if (currentQueue.length) {
        queue = currentQueue.concat(queue);
    } else {
        queueIndex = -1;
    }
    if (queue.length) {
        drainQueue();
    }
}

function drainQueue() {
    if (draining) {
        return;
    }
    var timeout = runTimeout(cleanUpNextTick);
    draining = true;

    var len = queue.length;
    while(len) {
        currentQueue = queue;
        queue = [];
        while (++queueIndex < len) {
            if (currentQueue) {
                currentQueue[queueIndex].run();
            }
        }
        queueIndex = -1;
        len = queue.length;
    }
    currentQueue = null;
    draining = false;
    runClearTimeout(timeout);
}

process.nextTick = function (fun) {
    var args = new Array(arguments.length - 1);
    if (arguments.length > 1) {
        for (var i = 1; i < arguments.length; i++) {
            args[i - 1] = arguments[i];
        }
    }
    queue.push(new Item(fun, args));
    if (queue.length === 1 && !draining) {
        runTimeout(drainQueue);
    }
};

// v8 likes predictible objects
function Item(fun, array) {
    this.fun = fun;
    this.array = array;
}
Item.prototype.run = function () {
    this.fun.apply(null, this.array);
};
process.title = 'browser';
process.browser = true;
process.env = {};
process.argv = [];
process.version = ''; // empty string to avoid regexp issues
process.versions = {};

function noop() {}

process.on = noop;
process.addListener = noop;
process.once = noop;
process.off = noop;
process.removeListener = noop;
process.removeAllListeners = noop;
process.emit = noop;
process.prependListener = noop;
process.prependOnceListener = noop;

process.listeners = function (name) { return [] }

process.binding = function (name) {
    throw new Error('process.binding is not supported');
};

process.cwd = function () { return '/' };
process.chdir = function (dir) {
    throw new Error('process.chdir is not supported');
};
process.umask = function() { return 0; };

},{}],5:[function(require,module,exports){
(function (setImmediate,clearImmediate){
var nextTick = require('process/browser.js').nextTick;
var apply = Function.prototype.apply;
var slice = Array.prototype.slice;
var immediateIds = {};
var nextImmediateId = 0;

// DOM APIs, for completeness

exports.setTimeout = function() {
  return new Timeout(apply.call(setTimeout, window, arguments), clearTimeout);
};
exports.setInterval = function() {
  return new Timeout(apply.call(setInterval, window, arguments), clearInterval);
};
exports.clearTimeout =
exports.clearInterval = function(timeout) { timeout.close(); };

function Timeout(id, clearFn) {
  this._id = id;
  this._clearFn = clearFn;
}
Timeout.prototype.unref = Timeout.prototype.ref = function() {};
Timeout.prototype.close = function() {
  this._clearFn.call(window, this._id);
};

// Does not start the time, just sets up the members needed.
exports.enroll = function(item, msecs) {
  clearTimeout(item._idleTimeoutId);
  item._idleTimeout = msecs;
};

exports.unenroll = function(item) {
  clearTimeout(item._idleTimeoutId);
  item._idleTimeout = -1;
};

exports._unrefActive = exports.active = function(item) {
  clearTimeout(item._idleTimeoutId);

  var msecs = item._idleTimeout;
  if (msecs >= 0) {
    item._idleTimeoutId = setTimeout(function onTimeout() {
      if (item._onTimeout)
        item._onTimeout();
    }, msecs);
  }
};

// That's not how node.js implements it but the exposed api is the same.
exports.setImmediate = typeof setImmediate === "function" ? setImmediate : function(fn) {
  var id = nextImmediateId++;
  var args = arguments.length < 2 ? false : slice.call(arguments, 1);

  immediateIds[id] = true;

  nextTick(function onNextTick() {
    if (immediateIds[id]) {
      // fn.call() is faster so we optimize for the common use-case
      // @see http://jsperf.com/call-apply-segu
      if (args) {
        fn.apply(null, args);
      } else {
        fn.call(null);
      }
      // Prevent ids from leaking
      exports.clearImmediate(id);
    }
  });

  return id;
};

exports.clearImmediate = typeof clearImmediate === "function" ? clearImmediate : function(id) {
  delete immediateIds[id];
};
}).call(this,require("timers").setImmediate,require("timers").clearImmediate)

},{"process/browser.js":4,"timers":5}],6:[function(require,module,exports){
(function (setImmediate){
"use strict";

var __importStar = this && this.__importStar || function (mod) {
  if (mod && mod.__esModule) return mod;
  var result = {};
  if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
  result["default"] = mod;
  return result;
};

Object.defineProperty(exports, "__esModule", {
  value: true
});

var ios = __importStar(require("frida-lib/ios")); // @ts-ignore


var MGStatic = ObjC.classes.MGStatic; // @ts-ignore

var SessionSDK = ObjC.classes.SessionSDK;

try {
  // ios.fast.setExceptionHandle(null)
  setImmediate(function () {
    main();
  });
} catch (e) {
  console.log(e);
}

function main() {
  ios.fast.showcallmethod_with_hookoccls(SessionSDK);
  ios.fast.showcallmethod_with_hookoccls(MGStatic); // ios.fast.showcallmethod_with_hookoccls(SessionAppInfo)
}

}).call(this,require("timers").setImmediate)

},{"frida-lib/ios":11,"timers":5}],7:[function(require,module,exports){
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var libc_1 = require("./libc");
var dyld;
(function (dyld) {
    dyld.FAT_MAGIC = 0xcafebabe;
    dyld.FAT_CIGAM = 0xbebafeca;
    dyld.MH_MAGIC = 0xfeedface;
    dyld.MH_CIGAM = 0xcefaedfe;
    dyld.MH_MAGIC_64 = 0xfeedfacf;
    dyld.MH_CIGAM_64 = 0xcffaedfe;
    dyld.LC_SEGMENT = 0x1;
    dyld.LC_SEGMENT_64 = 0x19;
    dyld.LC_LOAD_DYLIB = 0xc;
    dyld.LC_ENCRYPTION_INFO = 0x21;
    dyld.LC_ENCRYPTION_INFO_64 = 0x2C;
    dyld.LC_SYMTAB = 0x2;
    dyld.LC_DYSYMTAB = 0xb;
    dyld.SECTION_TYPE = 0x000000ff;
    dyld.S_LAZY_SYMBOL_POINTERS = 0x7;
    dyld.S_NON_LAZY_SYMBOL_POINTERS = 0x6;
    dyld.S_SYMBOL_STUBS = 0x8;
    dyld.S_MOD_INIT_FUNC_POINTERS = 0x9;
    dyld.S_MOD_TERM_FUNC_POINTERS = 0xa;
    dyld.S_COALESCED = 0xb;
    dyld.S_GB_ZEROFILL = 0xc;
    dyld.S_INTERPOSING = 0xd;
    dyld.INDIRECT_SYMBOL_ABS = 0x40000000;
    dyld.INDIRECT_SYMBOL_LOCAL = 0x80000000;
    dyld.RTLD_LAZY = 0x1;
    dyld.RTLD_NOW = 0x2;
    dyld.RTLD_LOCAL = 0x4;
    dyld.RTLD_GLOBAL = 0x8;
    dyld.dlopen = libc_1.libc.getExportFunction("dlopen", "pointer", ['pointer', 'int']);
    dyld.dlsym = libc_1.libc.getExportFunction("dlsym", "pointer", ['pointer', 'pointer']);
    dyld.dlclose = libc_1.libc.getExportFunction("dlclose", "int", ['pointer']);
    dyld.dlerror = libc_1.libc.getExportFunction("dlerror", "pointer", []);
    dyld.exit = libc_1.libc.getExportFunction("exit", "void", ["int"]);
    dyld._exit = libc_1.libc.getExportFunction("_exit", "void", ["int"]);
    dyld.dladdr = libc_1.libc.getExportFunction("dladdr", "int", ["pointer", "pointer"]);
    dyld._dyld_register_func_for_add_image = libc_1.libc.getExportFunction("_dyld_register_func_for_add_image", "void", ["pointer"]);
    dyld._dyld_register_func_for_remove_image = libc_1.libc.getExportFunction("_dyld_register_func_for_remove_image", "void", ["pointer"]);
    dyld._dyld_image_count = libc_1.libc.getExportFunction("_dyld_image_count", "uint", []);
    dyld._dyld_get_image_header = libc_1.libc.getExportFunction("_dyld_get_image_header", "pointer", ["uint"]);
    dyld._dyld_get_image_vmaddr_slide = libc_1.libc.getExportFunction("_dyld_get_image_vmaddr_slide", "long", ["uint"]);
    dyld._dyld_get_image_name = libc_1.libc.getExportFunction("_dyld_get_image_name", "pointer", ["uint"]);
})(dyld = exports.dyld || (exports.dyld = {}));

},{"./libc":12}],8:[function(require,module,exports){
"use strict";
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
var struct = __importStar(require("./struct"));
var dyld_1 = require("./dyld");
var ios = __importStar(require("frida-lib/ios/index"));
var fast;
(function (fast) {
    var UIAlertView = ObjC.classes.UIAlertView;
    var NSString = ObjC.classes.NSString;
    var UIApplication = ObjC.classes.UIApplication;
    var UITabBarController = ObjC.classes.UITabBarController;
    var UINavigationController = ObjC.classes.UINavigationController;
    var UIStoryboard = ObjC.classes.UIStoryboard;
    var NSThread = ObjC.classes.NSThread;
    function show_backtrace() {
        var traces = new ObjC.Object(NSThread.callStackSymbols());
        for (var i = 0; i < traces.count(); i++) {
            console.log(traces.objectAtIndex_(i));
        }
        // // @ts-ignore
        // var traces=Thread.backtrace()
        // for(var addrtrace of traces){
        //     var trace=DebugSymbol.fromAddress(<NativePointer><any>addrtrace)
        //     console.log(trace)
        // }
    }
    fast.show_backtrace = show_backtrace;
    function curViewController() {
        var vc;
        vc = UIApplication.sharedApplication().keyWindow().rootViewController();
        while (1) {
            if (vc.isKindOfClass_(UITabBarController.class())) {
                vc = vc.selectedViewController();
            }
            if (vc.isKindOfClass_(UINavigationController.class())) {
                vc = vc.visibleViewController();
            }
            if (vc.presentedViewController()) {
                vc = vc.presentedViewController();
            }
            else {
                break;
            }
        }
        return vc;
    }
    fast.curViewController = curViewController;
    function simple_alert(_text, _title) {
        var alertview = UIAlertView.alloc();
        alertview.initWithTitle_message_delegate_cancelButtonTitle_otherButtonTitles_(jsocstr(_title), jsocstr(_text), NULL, jsocstr("OK"), NULL);
        console.log(alertview);
        alertview.show(); //////////
    }
    fast.simple_alert = simple_alert;
    function hookoc(method, fn) {
        var oldfunc = method.implementation;
        method.implementation = ObjC.implement(method, function (self, selector) {
            var args = [];
            for (var _i = 2; _i < arguments.length; _i++) {
                args[_i - 2] = arguments[_i];
            }
            return fn.apply(void 0, [oldfunc, self, selector].concat(args));
        });
    }
    fast.hookoc = hookoc;
    function hookoccls(cls, fn) {
        var metnames = cls.$methods;
        // @ts-ignore
        var objnames = ObjC.classes.NSObject.$methods;
        for (var _i = 0, metnames_1 = metnames; _i < metnames_1.length; _i++) {
            var name = metnames_1[_i];
            if (objnames.indexOf(name) == -1) {
                // @ts-ignore
                var method = cls[name];
                hookoc(method, fn);
                // Interceptor.attach(method.implementation,{
                //     onEnter:fn,
                // })
            }
        }
    }
    fast.hookoccls = hookoccls;
    //直接显示所有调用，指定一个类
    function showcallmethod_with_hookoccls(cls) {
        ios.fast.hookoccls(cls, function (oldfn, _self, sel) {
            var args = [];
            for (var _i = 3; _i < arguments.length; _i++) {
                args[_i - 3] = arguments[_i];
            }
            var self = new ObjC.Object(_self);
            var selname = ObjC.selectorAsString(sel);
            var r = oldfn.apply(void 0, [self, sel].concat(args));
            // @ts-ignore
            var methodsign = self.methodSignatureForSelector_(sel);
            var sargs = [];
            var argscount = methodsign.numberOfArguments() - 2; //解析参数个数
            for (var i = 0; i < argscount; i++) {
                var argtype = String(methodsign.getArgumentTypeAtIndex_(i + 2));
                var arg;
                if (argtype == "@" || argtype == "@?") {
                    arg = new ObjC.Object(args[i]);
                }
                else {
                    arg = args[i];
                }
                // @ts-ignore
                sargs.push(arg);
            } //解析参数类型
            var sret;
            var srettype = String(methodsign.methodReturnType());
            if (srettype == "@") {
                sret = new ObjC.Object(r);
            }
            else {
                sret = r;
            }
            var spselname = selname.split(":");
            var newselname = "";
            if (sargs.length == 0) {
                newselname = spselname[0];
            }
            else {
                var toargscount = argscount;
                for (var i = 0; i < toargscount; i++) {
                    newselname = newselname + " " + spselname[i] + ":" + sargs[i] + " ";
                }
            }
            var xsret = "";
            if (sret != null && sret != undefined && sret != "") {
                xsret = sret + "<- ";
            }
            console.log("" + xsret + self.$className + " " + newselname);
            console.log();
            return r;
        });
    }
    fast.showcallmethod_with_hookoccls = showcallmethod_with_hookoccls;
    // settingsStoryboard
    // settingsView
    function storytoviewcontroller() {
        var story = UIStoryboard.storyboardWithName_bundle_(jsocstr("settingsStoryboard"), NULL);
        var vc = story.instantiateViewControllerWithIdentifier_(jsocstr("settingsView"));
        return vc;
    }
    fast.storytoviewcontroller = storytoviewcontroller;
    function hook_objc_msgSend(func) {
        Interceptor.attach(Module.findExportByName('/usr/lib/libobjc.A.dylib', 'objc_msgSend'), {
            onEnter: function (args) {
                var m = cjsstr(args[1]);
                if (m != 'length' && m != 'self' && m != 'retain' && m != "_mutate" && m != 'count' && m != 'hash' && !m.startsWith("dealloc") && !m.startsWith("__new:::") && !m.startsWith('release') && !m.startsWith('alloc') && !m.startsWith('_fastC')) {
                    func(args);
                }
            }
        });
    }
    fast.hook_objc_msgSend = hook_objc_msgSend;
    function mainfunc(fn) {
        ObjC.schedule(ObjC.mainQueue, fn);
    }
    fast.mainfunc = mainfunc;
    function jsocstr(str) {
        return NSString.stringWithUTF8String_(jscstr(str));
    }
    fast.jsocstr = jsocstr;
    function ocjsstr(str) {
        return str.UTF8String().readUtf8String();
    }
    fast.ocjsstr = ocjsstr;
    function jscstr(str) {
        return Memory.allocUtf8String(str);
    }
    fast.jscstr = jscstr;
    function cjsstr(str) {
        if (str == null) {
            return "";
        }
        if (str == NULL) {
            return "";
        }
        var r = str.readUtf8String();
        if (r == null) {
            return "";
        }
        else {
            return r;
        }
    }
    fast.cjsstr = cjsstr;
    function showinfo(obj) {
        Object.keys(obj.$ivars).forEach(function (v) {
            console.log('\t', v, '=', obj.$ivars[v]);
        });
    }
    fast.showinfo = showinfo;
    function observeClass(obj) {
        var k = obj;
        var name = obj.$className;
        k.$ownMethods.forEach(function (m) {
            var impl = k[m].implementation;
            console.log('Observing ' + name + ' ' + m);
            Interceptor.attach(impl, {
                onEnter: function (a) {
                    this.log = [];
                    this.log.push('(' + a[0] + ',' + a[1].readUtf8String() + ') ' + name + ' ' + m);
                    if (m.indexOf(':') !== -1) {
                        var params = m.split(':');
                        params[0] = params[0].split(' ')[1];
                        for (var i = 0; i < params.length - 1; i++) {
                            try {
                                this.log.push(params[i] + ': ' + new ObjC.Object(a[2 + i]).toString());
                            }
                            catch (e) {
                                this.log.push(params[i] + ': ' + a[2 + i].toString());
                            }
                        }
                    }
                    this.log.push(Thread.backtrace(this.context, Backtracer.ACCURATE)
                        .map(DebugSymbol.fromAddress)
                        .join('\n'));
                },
                onLeave: function (r) {
                    try {
                        this.log.push('RET: ' + new ObjC.Object(r).toString());
                    }
                    catch (e) {
                        this.log.push('RET: ' + r.toString());
                    }
                    console.log(this.log.join('\n') + '\n');
                }
            });
        });
    }
    fast.observeClass = observeClass;
    function setExceptionHandle(fn) {
        if (fn != null) {
            Process.setExceptionHandler(fn);
        }
        else {
            Process.setExceptionHandler(function (exception) {
                var Dlinfo = new struct.Dl_info(null);
                dyld_1.dyld.dladdr(exception.address, Dlinfo.ptr);
                console.log("\u6700\u8FD1\u7684\u7B26\u53F7:" + Dlinfo.dli_sname());
                console.log('异常堆栈:\n' + Thread.backtrace(exception.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
                console.log("\u6355\u83B7\u5230\u5F02\u5E38" + exception);
                send({ "done": true });
                Dlinfo.free();
                return false;
            });
        }
    }
    fast.setExceptionHandle = setExceptionHandle;
})(fast = exports.fast || (exports.fast = {}));

},{"./dyld":7,"./struct":13,"frida-lib/ios/index":11}],9:[function(require,module,exports){
"use strict";
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
var struct = __importStar(require("./struct"));
var dyld_1 = require("./dyld");
var path_1 = __importDefault(require("path"));
var fast_1 = require("./fast");
var natcall = new NativeCallback(function () {
    console.log("nav");
}, "void", ["pointer", "pointer", "pointer"]);
// let tzset=libc.getExportFunction("test","void",[])   //随便导出一个函数来屏蔽entry
var fishhook;
(function (fishhook) {
    var bindclass = /** @class */ (function () {
        function bindclass(_modename, _funcname, _newcalladdr) {
            this.name = "";
            this.oldaddr = NULL;
            this.newcall = NULL;
            this.modname = "";
            this.modname = _modename;
            this.name = _funcname;
            this.newcall = _newcalladdr;
        }
        return bindclass;
    }());
    fishhook.bindclass = bindclass;
    fishhook.bindclsvar = [];
    var my_bind = new NativeCallback(function (mh, intptr_t) {
        bind(mh, intptr_t);
    }, "void", ["pointer", "pointer"]);
    function isinbind_with_modname(_modfname) {
        for (var _i = 0, bindclsvar_1 = fishhook.bindclsvar; _i < bindclsvar_1.length; _i++) {
            var it = bindclsvar_1[_i];
            var itname = it.modname.toLowerCase();
            var itfname = it.name.toLocaleLowerCase();
            var modfname = _modfname.toLowerCase();
            if (itname == modfname || itname == "") {
                return true;
                // let callf = Module.findExportByName(modfname, itfname)
                // if (callf != null) {
                //     return true
                // }
            }
        }
        return false;
    }
    function perform_rebinding_with_section(_modpath, _section, _intptr_t, _symtab, _strtab, _indirect_symtab) {
        var modfpath = _modpath;
        var modfname = path_1.default.basename(modfpath);
        if (modfname == "") {
            return;
        }
        var section = new struct.section(_section);
        var intptr_t = _intptr_t;
        var symtab = _symtab;
        var strtab = _strtab;
        var indirect_symtab = _indirect_symtab;
        var indirect_symbol_indices = indirect_symtab.add(section.reserved1().readU32() * 4);
        var indirect_symbol_bindings = intptr_t.add(section.addr().readU64());
        if (isinbind_with_modname(modfname) == false) {
            return;
        }
        var ton = section.psize().readU64().toNumber() / Process.pointerSize;
        var _loop_1 = function () {
            var symtab_index = indirect_symbol_indices.add(i * 4).readU32();
            if (symtab_index == dyld_1.dyld.INDIRECT_SYMBOL_ABS || symtab_index == dyld_1.dyld.INDIRECT_SYMBOL_LOCAL || symtab_index == ((dyld_1.dyld.INDIRECT_SYMBOL_LOCAL | dyld_1.dyld.INDIRECT_SYMBOL_ABS) >>> 0)) {
                return "continue";
            }
            strtab_offset = symtab.add(symtab_index * 16).readU32();
            var symbol_name = strtab.add(strtab_offset).readCString();
            if (symbol_name.length < 1) {
                return { value: void 0 };
            }
            symbol_name = symbol_name.slice(1);
            fishhook.bindclsvar.forEach(function (it) {
                if (it.name == symbol_name) {
                    //console.log(`fname:${modfname} symbol_name: ${symbol_name}`)
                    var fnptr = indirect_symbol_bindings.add(i * Process.pointerSize);
                    if (it.name == symbol_name) {
                        if (it.oldaddr == NULL) {
                            it.oldaddr = fnptr.readPointer(); //保存第一个模块的符号，因为符号都是相同的
                        }
                        console.log("\u6210\u529FHook " + modfname + " " + symbol_name);
                        fnptr.writePointer(it.newcall);
                    }
                }
            });
        };
        var strtab_offset;
        for (var i = 0; i < ton; i++) {
            var state_1 = _loop_1();
            if (typeof state_1 === "object")
                return state_1.value;
        }
        // bindclsvar.forEach((it) => {
        //     if ((it.modname != "" && it.modname != modfname) || it.modname == "#initfunc") {  //如果不是指定模块
        //         return
        //     } else {
        //         //console.log(`name:${modfname} reserved1:${section.reserved1().readU32()} secsize:${section.psize().readU64()}  sectname: ${fast.cjsstr(section.sectname())}`)
        //         return;
        //         for (var i = 0; i < (section.psize().readU64().toNumber() / Process.pointerSize); i++) {
        //             let symtab_index = indirect_symbol_indices.add(i * 4).readU32()
        //             if (symtab_index == dyld.INDIRECT_SYMBOL_ABS || symtab_index == dyld.INDIRECT_SYMBOL_LOCAL || symtab_index == ((dyld.INDIRECT_SYMBOL_LOCAL | dyld.INDIRECT_SYMBOL_ABS) >>> 0)) {
        //                 continue
        //             }
        //             var strtab_offset = symtab.add(symtab_index * 16).readU32();
        //             let symbol_name = <string>strtab.add(strtab_offset).readCString()
        //             if (symbol_name.length < 1) {
        //                 return
        //             }
        //             symbol_name = symbol_name.slice(1)
        //             let fnptr = indirect_symbol_bindings.add(i * Process.pointerSize)
        //             bindclsvar.forEach((it) => {
        //                 //console.log(`${section.sectname().readUtf8String()}:${symbol_name}`)
        //                 if (it.name == symbol_name) {
        //                     if (it.oldaddr == NULL) {
        //                         it.oldaddr = fnptr.readPointer()  //保存第一个模块的符号，因为符号都是相同的
        //                     }
        //                     //console.log(`成功Hook ${modfname} ${symbol_name}`)
        //                     fnptr.writePointer(it.newcall)
        //                 }
        //             })
        //         }
        //     }
        // })
    }
    // function perform_rebindinginit_with_section(_modpath:string,_section: NativePointer, _intptr_t: NativePointer, _symtab: NativePointer, _strtab: NativePointer, _indirect_symtab: NativePointer) {
    //     let modfpath = _modpath
    //     let modfname = path.basename(modfpath)
    //     if (modfname == "") {
    //         return
    //     }
    //     let section = new struct.section(_section)
    //     let intptr_t = _intptr_t
    //     let symtab = _symtab
    //     let strtab = _strtab
    //     let indirect_symtab = _indirect_symtab
    //     let indirect_symbol_bindings = intptr_t.add(section.addr().readU64())
    //     bindclsvar.forEach((it) => {
    //         if (it.modname != "" && it.modname != modfname) {  //如果不是指定模块
    //             return
    //         } else {
    //             console.log(`name:${modfname} reserved1:${section.reserved1().readU32()} secsize:${section.psize().readU64()}  sectname: ${fast.cjsstr(section.sectname())}`)
    //             for (var i = 0; i < (section.psize().readU64().toNumber() / Process.pointerSize); i++) {
    //                 let fnptr = indirect_symbol_bindings.add(i * 8)
    //                 bindclsvar.forEach((it) => {       
    //                     fnptr.writePointer(tzset)
    //                     console.log(`成功Hook ${modfname} init`)
    //                 })
    //             }
    //         }
    //     })
    // }
    function bind(mh, intptr_t) {
        var dlinfo = new struct.Dl_info(null);
        if (dyld_1.dyld.dladdr(mh, dlinfo.ptr) == 0) {
            return;
        }
        var modpath = fast_1.fast.cjsstr(dlinfo.dli_fname());
        dlinfo.free();
        var linkedit_segment = new struct.segment_command(NULL);
        var symtab_cmd = new struct.symtab_command(NULL);
        var dysymtab_cmd = new struct.dysymtab_command(NULL);
        var macho = new struct.mach_header(mh);
        var pmh = mh.add(macho.size());
        for (var i = 0; i < macho.ncmds().readU32(); i++) {
            var lccommand = new struct.segment_command(pmh);
            if (lccommand.cmd().readU32() == dyld_1.dyld.LC_SEGMENT_64) {
                var segname = lccommand.segname().readUtf8String();
                if (segname == "__LINKEDIT") {
                    linkedit_segment.ptr = lccommand.ptr;
                }
            }
            else if (lccommand.cmd().readU32() == dyld_1.dyld.LC_SYMTAB) {
                symtab_cmd.ptr = lccommand.ptr;
            }
            else if (lccommand.cmd().readU32() == dyld_1.dyld.LC_DYSYMTAB) {
                dysymtab_cmd.ptr = lccommand.ptr;
            }
            pmh = pmh.add(lccommand.cmdsize().readU32());
        }
        if (symtab_cmd.isnull() || dysymtab_cmd.isnull() || linkedit_segment.isnull()) {
            console.log("遍历seg失败");
            return;
        }
        var linkedit_base = intptr_t.add(linkedit_segment.vmaddr().readU64()).sub(linkedit_segment.fileoff().readU64());
        var symtab = linkedit_base.add(symtab_cmd.symoff().readU32());
        var strtab = linkedit_base.add(symtab_cmd.stroff().readU32());
        var indirect_symtab = linkedit_base.add(dysymtab_cmd.indrectsymoff().readU32());
        var cur = mh.add(macho.size());
        for (var i = 0; i < macho.ncmds().readU32(); i++) {
            var cur_seg_cmd = new struct.segment_command(cur);
            if (cur_seg_cmd.cmd().readU32() == dyld_1.dyld.LC_SEGMENT_64) {
                var segname = cur_seg_cmd.segname().readUtf8String();
                if (segname != "__DATA" && segname != "__DATA_CONST") {
                    cur = cur.add(cur_seg_cmd.cmdsize().readU32());
                    continue;
                }
                for (var j = 0; j < cur_seg_cmd.nsects().readU64().toNumber(); j++) {
                    var tmpsec = new struct.section(NULL);
                    var sect = new struct.section(cur.add(cur_seg_cmd.size()).add(tmpsec.size() * j));
                    if ((sect.flags().readU32() & dyld_1.dyld.SECTION_TYPE) == dyld_1.dyld.S_LAZY_SYMBOL_POINTERS) {
                        perform_rebinding_with_section(modpath, sect.ptr, intptr_t, symtab, strtab, indirect_symtab);
                    }
                    if ((sect.flags().readU32() & dyld_1.dyld.SECTION_TYPE) == dyld_1.dyld.S_NON_LAZY_SYMBOL_POINTERS) {
                        perform_rebinding_with_section(modpath, sect.ptr, intptr_t, symtab, strtab, indirect_symtab);
                    }
                    // if ((sect.flags().readU32() & dyld.SECTION_TYPE) == dyld.S_MOD_INIT_FUNC_POINTERS) {
                    //     perform_rebindinginit_with_section(modpath,sect.ptr, intptr_t, symtab, strtab, indirect_symtab);
                    // }
                }
            }
            cur = cur.add(cur_seg_cmd.cmdsize().readU32());
        }
    }
    function startbind(bindcls) {
        if (bindcls.name == "") {
            return;
        }
        if (bindcls.newcall == NULL) {
            return;
        }
        fishhook.bindclsvar.push(bindcls);
        var ncount = dyld_1.dyld._dyld_image_count();
        for (var i = 0; i < ncount; i++) {
            var dyname = fast_1.fast.cjsstr(dyld_1.dyld._dyld_get_image_name(i));
            var pheader = dyld_1.dyld._dyld_get_image_header(i);
            var pslide = dyld_1.dyld._dyld_get_image_vmaddr_slide(i);
            bind(pheader, new NativePointer(pslide));
        }
        dyld_1.dyld._dyld_register_func_for_add_image(my_bind);
    }
    fishhook.startbind = startbind;
})(fishhook = exports.fishhook || (exports.fishhook = {}));

},{"./dyld":7,"./fast":8,"./struct":13,"path":3}],10:[function(require,module,exports){
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var libc_1 = require("./libc");
function NSSearchPathForDirectoriesInDomains() {
    return libc_1.libc.getExportFunction("NSSearchPathForDirectoriesInDomains", "pointer", ["int", "int", "int"]);
}
exports.NSSearchPathForDirectoriesInDomains = NSSearchPathForDirectoriesInDomains;
function NSStringFromSelector(_sel) {
    return new ObjC.Object((libc_1.libc.getExportFunction("NSStringFromSelector", "pointer", ["pointer",])(_sel)));
}
exports.NSStringFromSelector = NSStringFromSelector;
function exitdebug() {
    send({ "done": true });
}
exports.exitdebug = exitdebug;

},{"./libc":12}],11:[function(require,module,exports){
"use strict";
function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
__export(require("./dyld"));
__export(require("./struct"));
__export(require("./fishhook"));
__export(require("./libc"));
__export(require("./fast"));
__export(require("./func"));

},{"./dyld":7,"./fast":8,"./fishhook":9,"./func":10,"./libc":12,"./struct":13}],12:[function(require,module,exports){
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var libc;
(function (libc) {
    libc.O_RDONLY = 0;
    libc.O_WRONLY = 1;
    libc.O_RDWR = 2;
    libc.O_CREAT = 512;
    libc.SEEK_SET = 0;
    libc.SEEK_CUR = 1;
    libc.SEEK_END = 2;
    libc.read = getExportFunction("read", "int", ["int", "pointer", "int"]);
    libc.write = getExportFunction("write", "int", ["int", "pointer", "int"]);
    libc.lseek = getExportFunction("lseek", "int64", ["int", "int64", "int"]);
    libc.close = getExportFunction("close", "int", ["int"]);
    libc.remove = getExportFunction("remove", "int", ["pointer"]);
    libc.access = getExportFunction("access", "int", ["pointer", "int"]);
    libc.open = getExportFunction("open", "int", ["pointer", "int", "int"]);
    libc.free = getExportFunction("free", "void", ["pointer"]);
    libc.malloc = getExportFunction("malloc", "pointer", ["int"]);
    libc.sleep = getExportFunction("sleep", "int", ["int"]);
    libc.getenv = getExportFunction("getenv", "pointer", ['pointer']);
    function getExportFunction(name, ret, args) {
        var nptr;
        nptr = Module.findExportByName(null, name);
        if (nptr === null) {
            console.log("cannot find " + name);
            return null;
        }
        else {
            var funclet = new NativeFunction(nptr, ret, args);
            if (typeof funclet === "undefined") {
                console.log("parse error " + name);
                return null;
            }
            return funclet;
        }
    }
    libc.getExportFunction = getExportFunction;
})(libc = exports.libc || (exports.libc = {}));

},{}],13:[function(require,module,exports){
"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
var dyld_1 = require("./dyld");
var struct = /** @class */ (function () {
    function struct(ptr) {
        if (ptr == null) {
            this.ptr = Memory.alloc(this.size());
        }
        else {
            this.ptr = ptr;
        }
    }
    struct.prototype.isnull = function () {
        if (this.ptr == NULL) {
            return true;
        }
        else {
            return false;
        }
    };
    struct.prototype.free = function () {
        // if(this.ptr!=NULL){
        //     libc.free(this.ptr)
        //     this.ptr=NULL
        // }
    };
    return struct;
}());
exports.struct = struct;
var section = /** @class */ (function (_super) {
    __extends(section, _super);
    function section() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    section.prototype.sectname = function () {
        return this.ptr;
    };
    section.prototype.segname = function () {
        return this.ptr.add(16);
    };
    section.prototype.addr = function () {
        return this.ptr.add(32);
    };
    section.prototype.psize = function () {
        return this.ptr.add(40);
    };
    section.prototype.offset = function () {
        return this.ptr.add(48);
    };
    section.prototype.align = function () {
        return this.ptr.add(52);
    };
    section.prototype.reloff = function () {
        return this.ptr.add(56);
    };
    section.prototype.nreloc = function () {
        return this.ptr.add(60);
    };
    section.prototype.flags = function () {
        return this.ptr.add(64);
    };
    section.prototype.reserved1 = function () {
        return this.ptr.add(68);
    };
    section.prototype.reserved2 = function () {
        return this.ptr.add(72);
    };
    section.prototype.reserved3 = function () {
        return this.ptr.add(76);
    };
    section.prototype.size = function () {
        return 80;
    };
    return section;
}(struct));
exports.section = section;
var symtab_command = /** @class */ (function (_super) {
    __extends(symtab_command, _super);
    function symtab_command() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    symtab_command.prototype.cmd = function () {
        return this.ptr;
    };
    symtab_command.prototype.cmdsize = function () {
        return this.ptr.add(4);
    };
    symtab_command.prototype.symoff = function () {
        return this.ptr.add(8);
    };
    symtab_command.prototype.nsyms = function () {
        return this.ptr.add(12);
    };
    symtab_command.prototype.stroff = function () {
        return this.ptr.add(16);
    };
    symtab_command.prototype.strsize = function () {
        return this.ptr.add(20);
    };
    symtab_command.prototype.size = function () {
        return 24;
    };
    return symtab_command;
}(struct));
exports.symtab_command = symtab_command;
var dysymtab_command = /** @class */ (function (_super) {
    __extends(dysymtab_command, _super);
    function dysymtab_command() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    dysymtab_command.prototype.cmd = function () {
        return this.ptr;
    };
    dysymtab_command.prototype.cmdsize = function () {
        return this.ptr.add(4);
    };
    dysymtab_command.prototype.ilocalsym = function () {
        return this.ptr.add(8);
    };
    dysymtab_command.prototype.nlocalsym = function () {
        return this.ptr.add(12);
    };
    dysymtab_command.prototype.iextdefsym = function () {
        return this.ptr.add(16);
    };
    dysymtab_command.prototype.nextdefsym = function () {
        return this.ptr.add(20);
    };
    dysymtab_command.prototype.iundefsym = function () {
        return this.ptr.add(24);
    };
    dysymtab_command.prototype.nundefsym = function () {
        return this.ptr.add(28);
    };
    dysymtab_command.prototype.tocoff = function () {
        return this.ptr.add(32);
    };
    dysymtab_command.prototype.ntoc = function () {
        return this.ptr.add(36);
    };
    dysymtab_command.prototype.modtaboff = function () {
        return this.ptr.add(40);
    };
    dysymtab_command.prototype.nmodtab = function () {
        return this.ptr.add(44);
    };
    dysymtab_command.prototype.extrefsymoff = function () {
        return this.ptr.add(48);
    };
    dysymtab_command.prototype.nextrefsyms = function () {
        return this.ptr.add(52);
    };
    dysymtab_command.prototype.indrectsymoff = function () {
        return this.ptr.add(56);
    };
    dysymtab_command.prototype.nindirectsyms = function () {
        return this.ptr.add(60);
    };
    dysymtab_command.prototype.extreloff = function () {
        return this.ptr.add(64);
    };
    dysymtab_command.prototype.nextrel = function () {
        return this.ptr.add(68);
    };
    dysymtab_command.prototype.locreloff = function () {
        return this.ptr.add(72);
    };
    dysymtab_command.prototype.nlocrel = function () {
        return this.ptr.add(76);
    };
    dysymtab_command.prototype.size = function () {
        return 80;
    };
    return dysymtab_command;
}(struct));
exports.dysymtab_command = dysymtab_command;
var segment_command = /** @class */ (function (_super) {
    __extends(segment_command, _super);
    function segment_command() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    segment_command.prototype.cmd = function () {
        return this.ptr;
    };
    segment_command.prototype.cmdsize = function () {
        return this.ptr.add(4);
    };
    segment_command.prototype.segname = function () {
        return this.ptr.add(8);
    };
    segment_command.prototype.vmaddr = function () {
        return this.ptr.add(24);
    };
    segment_command.prototype.vmsize = function () {
        return this.ptr.add(32);
    };
    segment_command.prototype.fileoff = function () {
        return this.ptr.add(40);
    };
    segment_command.prototype.filesize = function () {
        return this.ptr.add(48);
    };
    segment_command.prototype.maxport = function () {
        return this.ptr.add(56);
    };
    segment_command.prototype.initport = function () {
        return this.ptr.add(60);
    };
    segment_command.prototype.nsects = function () {
        return this.ptr.add(64);
    };
    segment_command.prototype.flags = function () {
        return this.ptr.add(68);
    };
    segment_command.prototype.size = function () {
        return 72;
    };
    return segment_command;
}(struct));
exports.segment_command = segment_command;
//
var encryption_info_command = /** @class */ (function (_super) {
    __extends(encryption_info_command, _super);
    function encryption_info_command() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    encryption_info_command.prototype.size = function () {
        return 20;
    };
    encryption_info_command.prototype.cmd = function () {
        return this.ptr;
    };
    encryption_info_command.prototype.cmdsize = function () {
        return this.ptr.add(4);
    };
    encryption_info_command.prototype.cryptoff = function () {
        return this.ptr.add(8);
    };
    encryption_info_command.prototype.cryptsize = function () {
        return this.ptr.add(12);
    };
    encryption_info_command.prototype.cryptid = function () {
        return this.ptr.add(16);
    };
    return encryption_info_command;
}(struct));
exports.encryption_info_command = encryption_info_command;
var mach_header = /** @class */ (function (_super) {
    __extends(mach_header, _super);
    function mach_header() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    mach_header.prototype.magic = function () {
        return this.ptr;
    };
    mach_header.prototype.cputype = function () {
        return this.ptr.add(4);
    };
    mach_header.prototype.cpusubtype = function () {
        return this.ptr.add(8);
    };
    mach_header.prototype.filetype = function () {
        return this.ptr.add(12);
    };
    mach_header.prototype.ncmds = function () {
        return this.ptr.add(16);
    };
    mach_header.prototype.sizeofcmds = function () {
        return this.ptr.add(20);
    };
    mach_header.prototype.flags = function () {
        return this.ptr.add(24);
    };
    mach_header.prototype.size = function () {
        if (this.magic().readU32() == dyld_1.dyld.MH_MAGIC || this.magic().readU32() == dyld_1.dyld.MH_CIGAM) {
            return 28;
        }
        else if (this.magic().readU32() == dyld_1.dyld.MH_MAGIC_64 || this.magic().readU32() == dyld_1.dyld.MH_CIGAM_64) {
            return 32;
        }
        else {
            return 0;
        }
    };
    return mach_header;
}(struct));
exports.mach_header = mach_header;
var Dl_info = /** @class */ (function (_super) {
    __extends(Dl_info, _super);
    function Dl_info() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    Dl_info.prototype.dli_fname = function () {
        return this.ptr.readPointer();
    };
    Dl_info.prototype.dli_fbase = function () {
        return this.ptr.readPointer().add(8);
    };
    Dl_info.prototype.dli_sname = function () {
        return this.ptr.readPointer().add(16);
    };
    Dl_info.prototype.dli_saddr = function () {
        return this.ptr.readPointer().add(24);
    };
    Dl_info.prototype.size = function () {
        return 32;
    };
    return Dl_info;
}(struct));
exports.Dl_info = Dl_info;

},{"./dyld":7}]},{},[6])
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIkM6L1VzZXJzL2ZiaS9BcHBEYXRhL1JvYW1pbmcvbnBtL25vZGVfbW9kdWxlcy9mcmlkYS1jb21waWxlL25vZGVfbW9kdWxlcy9icm93c2VyLXBhY2svX3ByZWx1ZGUuanMiLCJDOi9Vc2Vycy9mYmkvQXBwRGF0YS9Sb2FtaW5nL25wbS9ub2RlX21vZHVsZXMvZnJpZGEtY29tcGlsZS9ub2RlX21vZHVsZXMvZXZlbnRzL2V2ZW50cy5qcyIsIkM6L1VzZXJzL2ZiaS9BcHBEYXRhL1JvYW1pbmcvbnBtL25vZGVfbW9kdWxlcy9mcmlkYS1jb21waWxlL25vZGVfbW9kdWxlcy9mcmlkYS1wcm9jZXNzL2luZGV4LmpzIiwiQzovVXNlcnMvZmJpL0FwcERhdGEvUm9hbWluZy9ucG0vbm9kZV9tb2R1bGVzL2ZyaWRhLWNvbXBpbGUvbm9kZV9tb2R1bGVzL3BhdGgtYnJvd3NlcmlmeS9pbmRleC5qcyIsIkM6L1VzZXJzL2ZiaS9BcHBEYXRhL1JvYW1pbmcvbnBtL25vZGVfbW9kdWxlcy9mcmlkYS1jb21waWxlL25vZGVfbW9kdWxlcy9wcm9jZXNzL2Jyb3dzZXIuanMiLCJDOi9Vc2Vycy9mYmkvQXBwRGF0YS9Sb2FtaW5nL25wbS9ub2RlX21vZHVsZXMvZnJpZGEtY29tcGlsZS9ub2RlX21vZHVsZXMvdGltZXJzLWJyb3dzZXJpZnkvbWFpbi5qcyIsImluZGV4LnRzIiwibm9kZV9tb2R1bGVzL2ZyaWRhLWxpYi9pb3MvZHlsZC5qcyIsIm5vZGVfbW9kdWxlcy9mcmlkYS1saWIvaW9zL2Zhc3QuanMiLCJub2RlX21vZHVsZXMvZnJpZGEtbGliL2lvcy9maXNoaG9vay5qcyIsIm5vZGVfbW9kdWxlcy9mcmlkYS1saWIvaW9zL2Z1bmMuanMiLCJub2RlX21vZHVsZXMvZnJpZGEtbGliL2lvcy9pbmRleC5qcyIsIm5vZGVfbW9kdWxlcy9mcmlkYS1saWIvaW9zL2xpYmMuanMiLCJub2RlX21vZHVsZXMvZnJpZGEtbGliL2lvcy9zdHJ1Y3QuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7QUNBQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBOzs7O0FDM2dCQTtBQUVBLElBQU0sWUFBWSxHQUFHLE9BQU8sQ0FBQyxRQUFELENBQTVCOztBQUVBLElBQU0sT0FBTyxHQUFHLE1BQU0sQ0FBQyxPQUFQLEdBQWlCLEVBQWpDO0FBRUEsT0FBTyxDQUFDLFFBQVIsR0FBbUIsTUFBTSxDQUFDLFFBQTFCO0FBRUEsT0FBTyxDQUFDLEtBQVIsR0FBZ0IsT0FBaEI7QUFDQSxPQUFPLENBQUMsT0FBUixHQUFrQixJQUFsQjtBQUNBLE9BQU8sQ0FBQyxHQUFSLEdBQWMsRUFBZDtBQUNBLE9BQU8sQ0FBQyxJQUFSLEdBQWUsRUFBZjtBQUNBLE9BQU8sQ0FBQyxPQUFSLEdBQWtCLEVBQWxCLEMsQ0FBc0I7O0FBQ3RCLE9BQU8sQ0FBQyxRQUFSLEdBQW1CLEVBQW5CO0FBRUEsT0FBTyxDQUFDLFlBQVIsR0FBdUIsWUFBdkI7QUFDQSxPQUFPLENBQUMsRUFBUixHQUFhLElBQWI7QUFDQSxPQUFPLENBQUMsV0FBUixHQUFzQixJQUF0QjtBQUNBLE9BQU8sQ0FBQyxJQUFSLEdBQWUsSUFBZjtBQUNBLE9BQU8sQ0FBQyxHQUFSLEdBQWMsSUFBZDtBQUNBLE9BQU8sQ0FBQyxjQUFSLEdBQXlCLElBQXpCO0FBQ0EsT0FBTyxDQUFDLGtCQUFSLEdBQTZCLElBQTdCO0FBQ0EsT0FBTyxDQUFDLElBQVIsR0FBZSxJQUFmOztBQUVBLE9BQU8sQ0FBQyxPQUFSLEdBQWtCLFVBQVUsSUFBVixFQUFnQjtBQUNoQyxRQUFNLElBQUksS0FBSixDQUFVLGtDQUFWLENBQU47QUFDRCxDQUZEOztBQUlBLE9BQU8sQ0FBQyxHQUFSLEdBQWMsWUFBWTtBQUN4QixTQUFPLEdBQVA7QUFDRCxDQUZEOztBQUdBLE9BQU8sQ0FBQyxLQUFSLEdBQWdCLFVBQVUsR0FBVixFQUFlO0FBQzdCLFFBQU0sSUFBSSxLQUFKLENBQVUsZ0NBQVYsQ0FBTjtBQUNELENBRkQ7O0FBR0EsT0FBTyxDQUFDLEtBQVIsR0FBZ0IsWUFBWTtBQUMxQixTQUFPLENBQVA7QUFDRCxDQUZEOztBQUlBLFNBQVMsSUFBVCxHQUFpQixDQUFFOzs7O0FDdENuQjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7Ozs7QUM5U0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTs7O0FDeExBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBOzs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FDM0VBLElBQUEsR0FBQSxHQUFBLFlBQUEsQ0FBQSxPQUFBLENBQUEsZUFBQSxDQUFBLENBQUEsQyxDQUNBOzs7QUFDQSxJQUFPLFFBQVEsR0FBRyxJQUFJLENBQUMsT0FBTCxDQUFhLFFBQS9CLEMsQ0FNQTs7QUFDQSxJQUFPLFVBQVUsR0FBQyxJQUFJLENBQUMsT0FBTCxDQUFhLFVBQS9COztBQUVBLElBQUc7QUFDQztBQUNBLEVBQUEsWUFBWSxDQUFDLFlBQUE7QUFDVCxJQUFBLElBQUk7QUFDUCxHQUZXLENBQVo7QUFJSCxDQU5ELENBTUMsT0FBTSxDQUFOLEVBQVE7QUFDTCxFQUFBLE9BQU8sQ0FBQyxHQUFSLENBQVksQ0FBWjtBQUNIOztBQUVELFNBQVMsSUFBVCxHQUFhO0FBQ1QsRUFBQSxHQUFHLENBQUMsSUFBSixDQUFTLDZCQUFULENBQXVDLFVBQXZDO0FBQ0EsRUFBQSxHQUFHLENBQUMsSUFBSixDQUFTLDZCQUFULENBQXVDLFFBQXZDLEVBRlMsQ0FHVDtBQUNIOzs7OztBQ3pCRDtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7O0FDL0NBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7O0FDblFBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBOztBQ25QQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTs7QUNmQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7O0FDWEE7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTs7QUN4Q0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBIiwiZmlsZSI6ImdlbmVyYXRlZC5qcyIsInNvdXJjZVJvb3QiOiIifQ==
