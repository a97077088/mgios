(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
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

},{}],2:[function(require,module,exports){
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

},{"process/browser.js":1,"timers":2}],3:[function(require,module,exports){
(function (setImmediate){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
}); // @ts-ignore

var RequestData = ObjC.classes.RequestData;
var NSString = ObjC.classes.NSString;

try {
  // ios.fast.setExceptionHandle(null)
  setImmediate(function () {
    main();
  });
} catch (e) {
  console.log(e);
}

function main() {
  // //aquireToken
  Interceptor.attach(RequestData["+ requestVerifyWithURL:withMethod:withData:withResult:"].implementation, {
    onEnter: function (args) {
      var a2 = new ObjC.Object(args[2]);
      var a4 = NSString.alloc().initWithData_encoding_(new ObjC.Object(args[4]), 1);
      console.log("aquireToken");
      console.log();
    }
  }); //dataCollectionService

  Interceptor.attach(RequestData["+ requestQulityWithURL:withMethod:withToken:withData:withResult:"].implementation, {
    onEnter: function (args) {
      var a2 = new ObjC.Object(args[2]);
      var a5 = NSString.alloc().initWithData_encoding_(new ObjC.Object(args[5]), 1);
      console.log("dataCollectionService\u539F\u59CB:" + a5);
      var evs = events_with_s(a5);

      for (var _i = 0, evs_1 = evs; _i < evs_1.length; _i++) {
        var ev = evs_1[_i];
        console.log(ev);
      }

      console.log();
    }
  }); //dataExposureService

  Interceptor.attach(RequestData["+ requestConfigWithURL:withMethod:withToken:withData:withResult:"].implementation, {
    onEnter: function (args) {
      var a2 = new ObjC.Object(args[2]);
      var a5 = NSString.alloc().initWithData_encoding_(new ObjC.Object(args[5]), 1);
      console.log("dataExposureService\u539F\u59CB:" + a5);
      var evs = events_with_s(a5);

      for (var _i = 0, evs_2 = evs; _i < evs_2.length; _i++) {
        var ev = evs_2[_i];
        console.log(ev);
      }

      console.log();
    }
  }); //dataEventService

  Interceptor.attach(RequestData["+ requestDefineWithURL:withMethod:withToken:withData:withResult:"].implementation, {
    onEnter: function (args) {
      var a2 = new ObjC.Object(args[2]);
      var a5 = NSString.alloc().initWithData_encoding_(new ObjC.Object(args[5]), 1);
      console.log("dataEventService\u539F\u59CB:" + a5);
      var evs = events_with_s(a5);

      for (var _i = 0, evs_3 = evs; _i < evs_3.length; _i++) {
        var ev = evs_3[_i];
        console.log(ev);
      }

      console.log();
    }
  });
}

var tss = {
  "4": "logIdentify",
  "5": "disConnectToServer",
  "8": "programDownloadFailed",
  "9": "playFailed",
  "10": "firstPageDelay",
  "11": "videoPageLoadDelay",
  "12": "searchPageLoadDelay",
  "14": "videoDownloadSpeed",
  "15": "heartBeatMessage",
  "16": "logIdentify",
  "17": "videoPlayStatus",
  "18": "purchaseStatistics",
  "19": "picloadDelay",
  "20": "pushMessageStatistics",
  "21": "exitApp",
  "22": "getPlayUrlDuration",
  "25": "downloadFlow",
  "26": "playFlow",
  "27": "crashexitdetail",
  "28": "errorLog",
  "37": "subSessionDownloadFlow",
  "26001000": "MGDownloadBytesEvent",
  "56100002": "MGPlayerGSLBReqEvent",
  "56100000": "MGPlayerCreateEvent",
  "56100001": "MGPlayerSetURLEvent",
  "56100009": "MGPlayerStartCmdEvent",
  "56100005": "MGPlayerM3U8ReqEvent",
  "56100006": "MGPlayerM3U8ParseEvent",
  "56100003": "MGPlayerDNSParseEvent",
  "56100004": "MGPlayerTCPOpenEvent",
  "56100010": "MGPlayerTSReqEvent",
  "56100014": "MGPlayerHttpConnectEvent",
  "56100007": "MGPlayerMediaOpenEvent",
  "56100008": "MGPlayerMediaInfoEvent",
  "56100011": "MGPlayerDecoderInitEvent",
  "56100012": "MGPlayerContiSeekEvent",
  "56100013": "MGPlayerFirstDecodeEvent",
  "56000015": "MGStuckEvent",
  "56000004": "MGFirstVideoRenderEvent",
  "57000000": "MGTrafficStatusticsEvent",
  "58000000": "MGErrorEvent",
  "56000016": "MGPlayerDecoderSwitchEvent",
  "59000000": "MGUserOperationEvent ",
  "70000000": "MGPlayerShutdownEvent",
  "60000000": "MGPlayerDurationEvent",
  "60100000": "MGPlayerAVDiffEvent",
  "80000000": "MGAuthenticationEvent"
};

function name_with_type(_s) {
  var r = tss[_s];

  if (r == null) {
    return "不支持解析的名字";
  }

  return r;
}

function events_with_s(_s) {
  var r = new Array();

  try {
    var jso = JSON.parse(_s);

    if (jso.sessionStart != null && Object.keys(jso.sessionStart).length > 0) {
      r.push("<sessionStart>->" + JSON.stringify(jso.sessionStart));
    } else if (jso.sessionEnd != null && Object.keys(jso.sessionEnd).length > 0) {
      r.push("<sessionEnd>->" + JSON.stringify(jso.sessionEnd));
    } else if (jso.customInfo != null && Object.keys(jso.customInfo).length > 0) {
      for (var i = 0; i < jso.customInfo.length; i++) {
        r.push("<" + name_with_type(jso.customInfo[i].type) + ":" + jso.customInfo[i].type + ">->" + JSON.stringify(jso.customInfo[i]));
      }
    } else if (jso.customEvent != null && Object.keys(jso.customEvent).length) {
      for (var i = 0; i < jso.customEvent.length; i++) {
        r.push("<" + jso.customEvent[i].eventName + ">->" + JSON.stringify(jso.customEvent[i]));
      }
    } else {
      console.log("不支持的解析格式");
    }
  } catch (e) {
    console.log(e);
  }

  return r;
}

}).call(this,require("timers").setImmediate)

},{"timers":2}]},{},[3])
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIkM6L1VzZXJzL2ZiaS9BcHBEYXRhL1JvYW1pbmcvbnBtL25vZGVfbW9kdWxlcy9mcmlkYS1jb21waWxlL25vZGVfbW9kdWxlcy9icm93c2VyLXBhY2svX3ByZWx1ZGUuanMiLCJDOi9Vc2Vycy9mYmkvQXBwRGF0YS9Sb2FtaW5nL25wbS9ub2RlX21vZHVsZXMvZnJpZGEtY29tcGlsZS9ub2RlX21vZHVsZXMvcHJvY2Vzcy9icm93c2VyLmpzIiwiQzovVXNlcnMvZmJpL0FwcERhdGEvUm9hbWluZy9ucG0vbm9kZV9tb2R1bGVzL2ZyaWRhLWNvbXBpbGUvbm9kZV9tb2R1bGVzL3RpbWVycy1icm93c2VyaWZ5L21haW4uanMiLCJpbmRleC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQTtBQ0FBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7OztBQ3hMQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTs7Ozs7Ozs7O0lDdkVBOztBQUNBLElBQU8sV0FBVyxHQUFDLElBQUksQ0FBQyxPQUFMLENBQWEsV0FBaEM7QUFDQSxJQUFPLFFBQVEsR0FBRyxJQUFJLENBQUMsT0FBTCxDQUFhLFFBQS9COztBQUVBLElBQUc7QUFDQztBQUNBLEVBQUEs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QTlCLEdBQXFDLENBQWhFLEVBQWtFO0FBQzlELE1BQUEsQ0FBQyxDQUFDLElBQUYsQ0FBTyxxQkFBbUIsSUFBSSxDQUFDLFNBQUwsQ0FBZSxHQUFHLENBQUMsWUFBbkIsQ0FBMUI7QUFDSCxLQUZELE1BRU0sSUFBRyxHQUFHLENBQUMsVUFBSixJQUFnQixJQUFoQixJQUFzQixNQUFNLENBQUMsSUFBUCxDQUFZLEdBQUcsQ0FBQyxVQUFoQixFQUE0QixNQUE1QixHQUFtQyxDQUE1RCxFQUE4RDtBQUNoRSxNQUFBLENBQUMsQ0FBQyxJQUFGLENBQU8sbUJBQWlCLElBQUksQ0FBQyxTQUFMLENBQWUsR0FBRyxDQUFDLFVBQW5CLENBQXhCO0FBQ0gsS0FGSyxNQUVBLElBQUcsR0FBRyxDQUFDLFVBQUosSUFBZ0IsSUFBaEIsSUFBc0IsTUFBTSxDQUFDLElBQVAsQ0FBWSxHQUFHLENBQUMsVUFBaEIsRUFBNEIsTUFBNUIsR0FBbUMsQ0FBNUQsRUFBOEQ7QUFDaEUsV0FBSSxJQUFJLENBQUMsR0FBQyxDQUFWLEVBQVksQ0FBQyxHQUFDLEdBQUcsQ0FBQyxVQUFKLENBQWUsTUFBN0IsRUFBb0MsQ0FBQyxFQUFyQyxFQUF3QztBQUNwQyxRQUFBLENBQUMsQ0FBQyxJQUFGLENBQU8sTUFBSSxjQUFjLENBQUMsR0FBRyxDQUFDLFVBQUosQ0FBZSxDQUFmLEVBQWtCLElBQW5CLENBQWxCLEdBQTBDLEdBQTFDLEdBQThDLEdBQUcsQ0FBQyxVQUFKLENBQWUsQ0FBZixFQUFrQixJQUFoRSxHQUFvRSxLQUFwRSxHQUEwRSxJQUFJLENBQUMsU0FBTCxDQUFlLEdBQUcsQ0FBQyxVQUFKLENBQWUsQ0FBZixDQUFmLENBQWpGO0FBQ0g7QUFDSixLQUpLLE1BSUEsSUFBRyxHQUFHLENBQUMsV0FBSixJQUFpQixJQUFqQixJQUF1QixNQUFNLENBQUMsSUFBUCxDQUFZLEdBQUcsQ0FBQyxXQUFoQixFQUE2QixNQUF2RCxFQUE4RDtBQUNoRSxXQUFJLElBQUksQ0FBQyxHQUFDLENBQVYsRUFBWSxDQUFDLEdBQUMsR0FBRyxDQUFDLFdBQUosQ0FBZ0IsTUFBOUIsRUFBcUMsQ0FBQyxFQUF0QyxFQUF5QztBQUNyQyxRQUFBLENBQUMsQ0FBQyxJQUFGLENBQU8sTUFBSSxHQUFHLENBQUMsV0FBSixDQUFnQixDQUFoQixFQUFtQixTQUF2QixHQUFnQyxLQUFoQyxHQUFzQyxJQUFJLENBQUMsU0FBTCxDQUFlLEdBQUcsQ0FBQyxXQUFKLENBQWdCLENBQWhCLENBQWYsQ0FBN0M7QUFDSDtBQUNKLEtBSkssTUFJRDtBQUNELE1BQUEsT0FBTyxDQUFDLEdBQVIsQ0FBWSxVQUFaO0FBQ0g7QUFDSixHQWpCRCxDQWlCQyxPQUFNLENBQU4sRUFBUTtBQUNMLElBQUEsT0FBTyxDQUFDLEdBQVIsQ0FBWSxDQUFaO0FBQ0g7O0FBQ0QsU0FBTyxDQUFQO0FBQ0giLCJmaWxlIjoiZ2VuZXJhdGVkLmpzIiwic291cmNlUm9vdCI6IiJ9
