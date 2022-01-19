/*!
JW Player version 8.11.5
Copyright (c) 2020, JW Player, All Rights Reserved 
https://github.com/jwplayer/jwplayer/blob/v8.11.5/README.md

This source code and its use and distribution is subject to the terms and conditions of the applicable license agreement. 
https://www.jwplayer.com/tos/

This product includes portions of other software. For the full text of licenses, see below:

JW Player Third Party Software Notices and/or Additional Terms and Conditions

**************************************************************************************************
The following software is used under Apache License 2.0
**************************************************************************************************

vtt.js v0.13.0
Copyright (c) 2020 Mozilla (http://mozilla.org)
https://github.com/mozilla/vtt.js/blob/v0.13.0/LICENSE

* * *

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and
limitations under the License.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

**************************************************************************************************
The following software is used under MIT license
**************************************************************************************************

Underscore.js v1.6.0
Copyright (c) 2009-2014 Jeremy Ashkenas, DocumentCloud and Investigative
https://github.com/jashkenas/underscore/blob/1.6.0/LICENSE

Backbone backbone.events.js v1.1.2
Copyright (c) 2010-2014 Jeremy Ashkenas, DocumentCloud
https://github.com/jashkenas/backbone/blob/1.1.2/LICENSE

Promise Polyfill v7.1.1
Copyright (c) 2014 Taylor Hakes and Forbes Lindesay
https://github.com/taylorhakes/promise-polyfill/blob/v7.1.1/LICENSE

can-autoplay.js v3.0.0
Copyright (c) 2017 video-dev
https://github.com/video-dev/can-autoplay/blob/v3.0.0/LICENSE

* * *

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

**************************************************************************************************
The following software is used under W3C license
**************************************************************************************************

Intersection Observer v0.5.0
Copyright (c) 2016 Google Inc. (http://google.com)
https://github.com/w3c/IntersectionObserver/blob/v0.5.0/LICENSE.md

* * *

W3C SOFTWARE AND DOCUMENT NOTICE AND LICENSE
Status: This license takes effect 13 May, 2015.

This work is being provided by the copyright holders under the following license.

License
By obtaining and/or copying this work, you (the licensee) agree that you have read, understood, and will comply with the following terms and conditions.

Permission to copy, modify, and distribute this work, with or without modification, for any purpose and without fee or royalty is hereby granted, provided that you include the following on ALL copies of the work or portions thereof, including modifications:

The full text of this NOTICE in a location viewable to users of the redistributed or derivative work.

Any pre-existing intellectual property disclaimers, notices, or terms and conditions. If none exist, the W3C Software and Document Short Notice should be included.

Notice of any changes or modifications, through a copyright statement on the new code or document such as "This software or document includes material copied from or derived from [title and URI of the W3C document]. Copyright © [YEAR] W3C® (MIT, ERCIM, Keio, Beihang)."

Disclaimers
THIS WORK IS PROVIDED "AS IS," AND COPYRIGHT HOLDERS MAKE NO REPRESENTATIONS OR WARRANTIES, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO, WARRANTIES OF MERCHANTABILITY OR FITNESS FOR ANY PARTICULAR PURPOSE OR THAT THE USE OF THE SOFTWARE OR DOCUMENT WILL NOT INFRINGE ANY THIRD PARTY PATENTS, COPYRIGHTS, TRADEMARKS OR OTHER RIGHTS.

COPYRIGHT HOLDERS WILL NOT BE LIABLE FOR ANY DIRECT, INDIRECT, SPECIAL OR CONSEQUENTIAL DAMAGES ARISING OUT OF ANY USE OF THE SOFTWARE OR DOCUMENT.

The name and trademarks of copyright holders may NOT be used in advertising or publicity pertaining to the work without specific, written prior permission. Title to copyright in this work will at all times remain with copyright holders.
*/
(window.webpackJsonpjwplayer = window.webpackJsonpjwplayer || []).push([
  [2],
  {
    18: function (e, t, n) {
      "use strict";
      n.r(t);
      var i = n(0),
        o = n(12),
        r = n(50),
        a = n(36);
      var s = n(44),
        l = n(51),
        c = n(26),
        u = n(25),
        d = n(3),
        f = n(46),
        g = n(2),
        h = n(7),
        p = n(34);
      function b(e) {
        var t = !1;
        return {
          async: function () {
            var n = this,
              i = arguments;
            return Promise.resolve().then(function () {
              if (!t) return e.apply(n, i);
            });
          },
          cancel: function () {
            t = !0;
          },
          cancelled: function () {
            return t;
          },
        };
      }
      var m = n(1);
      function w(e) {
        return function (t, n) {
          var o = e.mediaModel,
            r = Object(i.g)({}, n, { type: t });
          switch (t) {
            case d.T:
              if (o.get(d.T) === n.mediaType) return;
              o.set(d.T, n.mediaType);
              break;
            case d.U:
              return void o.set(d.U, Object(i.g)({}, n));
            case d.M:
              if (n[t] === e.model.getMute()) return;
              break;
            case d.bb:
              n.newstate === d.mb && (e.thenPlayPromise.cancel(), o.srcReset());
              var a = o.attributes.mediaState;
              (o.attributes.mediaState = n.newstate),
                o.trigger("change:mediaState", o, n.newstate, a);
              break;
            case d.F:
              return (
                (e.beforeComplete = !0),
                e.trigger(d.B, r),
                void (e.attached && !e.background && e._playbackComplete())
              );
            case d.G:
              o.get("setup")
                ? (e.thenPlayPromise.cancel(), o.srcReset())
                : ((t = d.tb), (r.code += 1e5));
              break;
            case d.K:
              r.metadataType || (r.metadataType = "unknown");
              var s = n.duration;
              Object(i.u)(s) &&
                (o.set("seekRange", n.seekRange), o.set("duration", s));
              break;
            case d.D:
              o.set("buffer", n.bufferPercent);
            case d.S:
              o.set("seekRange", n.seekRange),
                o.set("position", n.position),
                o.set("currentTime", n.currentTime);
              var l = n.duration;
              Object(i.u)(l) && o.set("duration", l),
                t === d.S &&
                  Object(i.r)(e.item.starttime) &&
                  delete e.item.starttime;
              break;
            case d.R:
              var c = e.mediaElement;
              c && c.paused && o.set("mediaState", "paused");
              break;
            case d.I:
              o.set(d.I, n.levels);
            case d.J:
              var u = n.currentQuality,
                f = n.levels;
              u > -1 && f.length > 1 && o.set("currentLevel", parseInt(u));
              break;
            case d.f:
              o.set(d.f, n.tracks);
            case d.g:
              var g = n.currentTrack,
                h = n.tracks;
              g > -1 &&
                h.length > 0 &&
                g < h.length &&
                o.set("currentAudioTrack", parseInt(g));
          }
          e.trigger(t, r);
        };
      }
      var v = n(8),
        y = n(45),
        j = n(41);
      function k(e) {
        return (k =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (e) {
                return typeof e;
              }
            : function (e) {
                return e &&
                  "function" == typeof Symbol &&
                  e.constructor === Symbol &&
                  e !== Symbol.prototype
                  ? "symbol"
                  : typeof e;
              })(e);
      }
      function O(e, t) {
        if (!(e instanceof t))
          throw new TypeError("Cannot call a class as a function");
      }
      function x(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function C(e, t, n) {
        return t && x(e.prototype, t), n && x(e, n), e;
      }
      function M(e, t) {
        return !t || ("object" !== k(t) && "function" != typeof t)
          ? (function (e) {
              if (void 0 === e)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return e;
            })(e)
          : t;
      }
      function _(e) {
        return (_ = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function S(e, t) {
        if ("function" != typeof t && null !== t)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (e.prototype = Object.create(t && t.prototype, {
          constructor: { value: e, writable: !0, configurable: !0 },
        })),
          t && P(e, t);
      }
      function P(e, t) {
        return (P =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var T = (function (e) {
          function t() {
            var e;
            return (
              O(this, t),
              ((e = M(this, _(t).call(this))).providerController = null),
              (e._provider = null),
              e.addAttributes({ mediaModel: new E() }),
              e
            );
          }
          return (
            S(t, e),
            C(t, [
              {
                key: "setup",
                value: function (e) {
                  return (
                    (e = e || {}),
                    this._normalizeConfig(e),
                    Object(i.g)(this.attributes, e, j.b),
                    (this.providerController = new p.a(
                      this.getConfiguration()
                    )),
                    this.setAutoStart(),
                    this
                  );
                },
              },
              {
                key: "getConfiguration",
                value: function () {
                  var e = this.clone(),
                    t = e.mediaModel.attributes;
                  return (
                    Object.keys(j.a).forEach(function (n) {
                      e[n] = t[n];
                    }),
                    (e.instreamMode = !!e.instream),
                    delete e.instream,
                    delete e.mediaModel,
                    e
                  );
                },
              },
              {
                key: "persistQualityLevel",
                value: function (e, t) {
                  var n = t[e] || {},
                    o = n.label,
                    r = Object(i.u)(n.bitrate) ? n.bitrate : null;
                  this.set("bitrateSelection", r), this.set("qualityLabel", o);
                },
              },
              {
                key: "setActiveItem",
                value: function (e) {
                  var t = this.get("playlist")[e];
                  this.resetItem(t),
                    (this.attributes.playlistItem = null),
                    this.set("item", e),
                    this.set("minDvrWindow", t.minDvrWindow),
                    this.set("dvrSeekLimit", t.dvrSeekLimit),
                    this.set("playlistItem", t);
                },
              },
              {
                key: "setMediaModel",
                value: function (e) {
                  this.mediaModel &&
                    this.mediaModel !== e &&
                    this.mediaModel.off(),
                    (e = e || new E()),
                    this.set("mediaModel", e),
                    (function (e) {
                      var t = e.get("mediaState");
                      e.trigger("change:mediaState", e, t, t);
                    })(e);
                },
              },
              {
                key: "destroy",
                value: function () {
                  (this.attributes._destroyed = !0),
                    this.off(),
                    this._provider &&
                      (this._provider.off(null, null, this),
                      this._provider.destroy());
                },
              },
              {
                key: "getVideo",
                value: function () {
                  return this._provider;
                },
              },
              {
                key: "setFullscreen",
                value: function (e) {
                  (e = !!e) !== this.get("fullscreen") &&
                    this.set("fullscreen", e);
                },
              },
              {
                key: "getProviders",
                value: function () {
                  return this.providerController;
                },
              },
              {
                key: "setVolume",
                value: function (e) {
                  if (Object(i.u)(e)) {
                    var t = Math.min(Math.max(0, e), 100);
                    this.set("volume", t);
                    var n = 0 === t;
                    n !== this.getMute() && this.setMute(n);
                  }
                },
              },
              {
                key: "getMute",
                value: function () {
                  return this.get("autostartMuted") || this.get("mute");
                },
              },
              {
                key: "setMute",
                value: function (e) {
                  if (
                    (void 0 === e && (e = !this.getMute()),
                    this.set("mute", !!e),
                    !e)
                  ) {
                    var t = Math.max(10, this.get("volume"));
                    this.set("autostartMuted", !1), this.setVolume(t);
                  }
                },
              },
              {
                key: "setStreamType",
                value: function (e) {
                  this.set("streamType", e),
                    "LIVE" === e && this.setPlaybackRate(1);
                },
              },
              {
                key: "setProvider",
                value: function (e) {
                  (this._provider = e), A(this, e);
                },
              },
              {
                key: "resetProvider",
                value: function () {
                  (this._provider = null), this.set("provider", void 0);
                },
              },
              {
                key: "setPlaybackRate",
                value: function (e) {
                  Object(i.r)(e) &&
                    ((e = Math.max(Math.min(e, 4), 0.25)),
                    "LIVE" === this.get("streamType") && (e = 1),
                    this.set("defaultPlaybackRate", e),
                    this._provider &&
                      this._provider.setPlaybackRate &&
                      this._provider.setPlaybackRate(e));
                },
              },
              {
                key: "persistCaptionsTrack",
                value: function () {
                  var e = this.get("captionsTrack");
                  e
                    ? this.set("captionLabel", e.name)
                    : this.set("captionLabel", "Off");
                },
              },
              {
                key: "setVideoSubtitleTrack",
                value: function (e, t) {
                  this.set("captionsIndex", e),
                    e &&
                      t &&
                      e <= t.length &&
                      t[e - 1].data &&
                      this.set("captionsTrack", t[e - 1]);
                },
              },
              {
                key: "persistVideoSubtitleTrack",
                value: function (e, t) {
                  this.setVideoSubtitleTrack(e, t), this.persistCaptionsTrack();
                },
              },
              {
                key: "setAutoStart",
                value: function (e) {
                  void 0 !== e && this.set("autostart", e);
                  var t = v.OS.mobile && this.get("autostart");
                  this.set(
                    "playOnViewable",
                    t || "viewable" === this.get("autostart")
                  );
                },
              },
              {
                key: "resetItem",
                value: function (e) {
                  var t = e ? Object(g.g)(e.starttime) : 0,
                    n = e ? Object(g.g)(e.duration) : 0,
                    i = this.mediaModel;
                  this.set("playRejected", !1),
                    (this.attributes.itemMeta = {}),
                    i.set("position", t),
                    i.set("currentTime", 0),
                    i.set("duration", n);
                },
              },
              {
                key: "persistBandwidthEstimate",
                value: function (e) {
                  Object(i.u)(e) && this.set("bandwidthEstimate", e);
                },
              },
              {
                key: "_normalizeConfig",
                value: function (e) {
                  var t = e.floating;
                  t && t.disabled && delete e.floating;
                },
              },
            ]),
            t
          );
        })(y.a),
        A = function (e, t) {
          e.set("provider", t.getName()),
            !0 === e.get("instreamMode") && (t.instreamMode = !0),
            -1 === t.getName().name.indexOf("flash") &&
              (e.set("flashThrottle", void 0), e.set("flashBlocked", !1)),
            e.setPlaybackRate(e.get("defaultPlaybackRate")),
            e.set("supportsPlaybackRate", t.supportsPlaybackRate),
            e.set("playbackRate", t.getPlaybackRate()),
            e.set("renderCaptionsNatively", t.renderNatively);
        };
      var E = (function (e) {
          function t() {
            var e;
            return (
              O(this, t),
              (e = M(this, _(t).call(this))).addAttributes({
                mediaState: d.mb,
              }),
              e
            );
          }
          return (
            S(t, e),
            C(t, [
              {
                key: "srcReset",
                value: function () {
                  Object(i.g)(this.attributes, {
                    setup: !1,
                    started: !1,
                    preloaded: !1,
                    visualQuality: null,
                    buffer: 0,
                    currentTime: 0,
                  });
                },
              },
            ]),
            t
          );
        })(y.a),
        R = T;
      function I(e) {
        return (I =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (e) {
                return typeof e;
              }
            : function (e) {
                return e &&
                  "function" == typeof Symbol &&
                  e.constructor === Symbol &&
                  e !== Symbol.prototype
                  ? "symbol"
                  : typeof e;
              })(e);
      }
      function L(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function V(e) {
        return (V = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function F(e, t) {
        return (F =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      function z(e) {
        if (void 0 === e)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return e;
      }
      var N = (function (e) {
        function t(e, n) {
          var i, o, r, a;
          return (
            (function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
            (o = this),
            (r = V(t).call(this)),
            ((i =
              !r || ("object" !== I(r) && "function" != typeof r)
                ? z(o)
                : r).attached = !0),
            (i.beforeComplete = !1),
            (i.item = null),
            (i.mediaModel = new E()),
            (i.model = n),
            (i.provider = e),
            (i.providerListener = new w(z(z(i)))),
            (i.thenPlayPromise = b(function () {})),
            (a = z(z(i))).provider.on("all", a.providerListener, a),
            (i.eventQueue = new s.a(z(z(i)), ["trigger"], function () {
              return !i.attached || i.background;
            })),
            i
          );
        }
        var n, o, r;
        return (
          (function (e, t) {
            if ("function" != typeof t && null !== t)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (e.prototype = Object.create(t && t.prototype, {
              constructor: { value: e, writable: !0, configurable: !0 },
            })),
              t && F(e, t);
          })(t, e),
          (n = t),
          (o = [
            {
              key: "play",
              value: function (e) {
                var t = this.item,
                  n = this.model,
                  i = this.mediaModel,
                  o = this.provider;
                if (
                  (e || (e = n.get("playReason")),
                  n.set("playRejected", !1),
                  i.get("setup"))
                )
                  return o.play() || Promise.resolve();
                i.set("setup", !0);
                var r = this._loadAndPlay(t, o);
                return i.get("started") ? r : this._playAttempt(r, e);
              },
            },
            {
              key: "stop",
              value: function () {
                var e = this.provider;
                (this.beforeComplete = !1), e.stop();
              },
            },
            {
              key: "pause",
              value: function () {
                this.provider.pause();
              },
            },
            {
              key: "preload",
              value: function () {
                var e = this.item,
                  t = this.mediaModel,
                  n = this.provider;
                !e ||
                  (e && "none" === e.preload) ||
                  !this.attached ||
                  this.setup ||
                  this.preloaded ||
                  (t.set("preloaded", !0), n.preload(e));
              },
            },
            {
              key: "destroy",
              value: function () {
                var e = this.provider,
                  t = this.mediaModel;
                this.off(),
                  t.off(),
                  e.off(),
                  this.eventQueue.destroy(),
                  this.detach(),
                  e.getContainer() && e.remove(),
                  delete e.instreamMode,
                  (this.provider = null),
                  (this.item = null);
              },
            },
            {
              key: "attach",
              value: function () {
                var e = this.model,
                  t = this.provider;
                e.setPlaybackRate(e.get("defaultPlaybackRate")),
                  t.attachMedia(),
                  (this.attached = !0),
                  this.eventQueue.flush(),
                  this.beforeComplete && this._playbackComplete();
              },
            },
            {
              key: "detach",
              value: function () {
                var e = this.provider;
                this.thenPlayPromise.cancel();
                var t = e.detachMedia();
                return (this.attached = !1), t;
              },
            },
            {
              key: "_playAttempt",
              value: function (e, t) {
                var n = this,
                  o = this.item,
                  r = this.mediaModel,
                  a = this.model,
                  s = this.provider,
                  l = s ? s.video : null;
                return (
                  this.trigger(d.N, { item: o, playReason: t }),
                  (l ? l.paused : a.get(d.bb) !== d.pb) || a.set(d.bb, d.jb),
                  e
                    .then(function () {
                      r.get("setup") &&
                        (r.set("started", !0),
                        r === a.mediaModel &&
                          (function (e) {
                            var t = e.get("mediaState");
                            e.trigger("change:mediaState", e, t, t);
                          })(r));
                    })
                    .catch(function (e) {
                      if (n.item && r === a.mediaModel) {
                        if ((a.set("playRejected", !0), l && l.paused)) {
                          if (l.src === location.href)
                            return n._loadAndPlay(o, s);
                          r.set("mediaState", d.ob);
                        }
                        var c = Object(i.g)(new m.n(null, Object(m.w)(e), e), {
                          error: e,
                          item: o,
                          playReason: t,
                        });
                        throw (delete c.key, n.trigger(d.O, c), e);
                      }
                    })
                );
              },
            },
            {
              key: "_playbackComplete",
              value: function () {
                var e = this.item,
                  t = this.provider;
                e && delete e.starttime,
                  (this.beforeComplete = !1),
                  t.setState(d.kb),
                  this.trigger(d.F, {});
              },
            },
            {
              key: "_loadAndPlay",
              value: function () {
                var e = this.item,
                  t = this.provider,
                  n = t.load(e);
                if (n) {
                  var i = b(function () {
                    return t.play() || Promise.resolve();
                  });
                  return (this.thenPlayPromise = i), n.then(i.async);
                }
                return t.play() || Promise.resolve();
              },
            },
            {
              key: "audioTrack",
              get: function () {
                return this.provider.getCurrentAudioTrack();
              },
              set: function (e) {
                this.provider.setCurrentAudioTrack(e);
              },
            },
            {
              key: "quality",
              get: function () {
                return this.provider.getCurrentQuality();
              },
              set: function (e) {
                this.provider.setCurrentQuality(e);
              },
            },
            {
              key: "audioTracks",
              get: function () {
                return this.provider.getAudioTracks();
              },
            },
            {
              key: "background",
              get: function () {
                var e = this.container,
                  t = this.provider;
                return (
                  !!this.attached &&
                  !!t.video &&
                  (!e || (e && !e.contains(t.video)))
                );
              },
              set: function (e) {
                var t = this.container,
                  n = this.provider;
                n.video
                  ? t &&
                    (e
                      ? this.background ||
                        (this.thenPlayPromise.cancel(),
                        this.pause(),
                        t.removeChild(n.video),
                        (this.container = null))
                      : (this.eventQueue.flush(),
                        this.beforeComplete && this._playbackComplete()))
                  : e
                  ? this.detach()
                  : this.attach();
              },
            },
            {
              key: "container",
              get: function () {
                return this.provider.getContainer();
              },
              set: function (e) {
                this.provider.setContainer(e);
              },
            },
            {
              key: "mediaElement",
              get: function () {
                return this.provider.video;
              },
            },
            {
              key: "preloaded",
              get: function () {
                return this.mediaModel.get("preloaded");
              },
            },
            {
              key: "qualities",
              get: function () {
                return this.provider.getQualityLevels();
              },
            },
            {
              key: "setup",
              get: function () {
                return this.mediaModel.get("setup");
              },
            },
            {
              key: "started",
              get: function () {
                return this.mediaModel.get("started");
              },
            },
            {
              key: "activeItem",
              set: function (e) {
                var t = (this.mediaModel = new E()),
                  n = e ? Object(g.g)(e.starttime) : 0,
                  i = e ? Object(g.g)(e.duration) : 0,
                  o = t.attributes;
                t.srcReset(),
                  (o.position = n),
                  (o.duration = i),
                  (this.item = e),
                  this.provider.init(e);
              },
            },
            {
              key: "controls",
              set: function (e) {
                this.provider.setControls(e);
              },
            },
            {
              key: "mute",
              set: function (e) {
                this.provider.mute(e);
              },
            },
            {
              key: "position",
              set: function (e) {
                var t = this.provider;
                this.model.get("scrubbing") && t.fastSeek
                  ? t.fastSeek(e)
                  : t.seek(e);
              },
            },
            {
              key: "subtitles",
              set: function (e) {
                this.provider.setSubtitlesTrack &&
                  this.provider.setSubtitlesTrack(e);
              },
            },
            {
              key: "volume",
              set: function (e) {
                this.provider.volume(e);
              },
            },
          ]) && L(n.prototype, o),
          r && L(n, r),
          t
        );
      })(h.a);
      function H(e) {
        return (H =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (e) {
                return typeof e;
              }
            : function (e) {
                return e &&
                  "function" == typeof Symbol &&
                  e.constructor === Symbol &&
                  e !== Symbol.prototype
                  ? "symbol"
                  : typeof e;
              })(e);
      }
      function B(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function q(e) {
        return (q = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function D(e, t) {
        return (D =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      function W(e) {
        if (void 0 === e)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return e;
      }
      function U(e, t) {
        var n = t.mediaControllerListener;
        e.off().on("all", n, t);
      }
      function Q(e) {
        return e && e.sources && e.sources[0];
      }
      var Y = (function (e) {
        function t(e, n) {
          var o, r, a, s, l;
          return (
            (function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
            (r = this),
            ((o =
              !(a = q(t).call(this)) ||
              ("object" !== H(a) && "function" != typeof a)
                ? W(r)
                : a).adPlaying = !1),
            (o.background =
              ((s = null),
              (l = null),
              Object.defineProperties(
                {
                  setNext: function (e, t) {
                    l = { item: e, loadPromise: t };
                  },
                  isNext: function (e) {
                    return !(
                      !l ||
                      JSON.stringify(l.item.sources[0]) !==
                        JSON.stringify(e.sources[0])
                    );
                  },
                  clearNext: function () {
                    l = null;
                  },
                },
                {
                  nextLoadPromise: {
                    get: function () {
                      return l ? l.loadPromise : null;
                    },
                  },
                  currentMedia: {
                    get: function () {
                      return s;
                    },
                    set: function (e) {
                      s = e;
                    },
                  },
                }
              ))),
            (o.mediaPool = n),
            (o.mediaController = null),
            (o.mediaControllerListener = (function (e, t) {
              return function (n, o) {
                switch (n) {
                  case d.bb:
                    return;
                  case "flashThrottle":
                  case "flashBlocked":
                    return void e.set(n, o.value);
                  case d.V:
                  case d.M:
                    return void e.set(n, o[n]);
                  case d.P:
                    return void e.set("playbackRate", o.playbackRate);
                  case d.K:
                    Object(i.g)(e.get("itemMeta"), o.metadata);
                    break;
                  case d.J:
                    e.persistQualityLevel(o.currentQuality, o.levels);
                    break;
                  case "subtitlesTrackChanged":
                    e.persistVideoSubtitleTrack(o.currentTrack, o.tracks);
                    break;
                  case d.S:
                  case d.Q:
                  case d.R:
                  case d.X:
                  case "subtitlesTracks":
                  case "subtitlesTracksData":
                    e.trigger(n, o);
                    break;
                  case d.i:
                    return void e.persistBandwidthEstimate(o.bandwidthEstimate);
                }
                t.trigger(n, o);
              };
            })(e, W(W(o)))),
            (o.model = e),
            (o.providers = new p.a(e.getConfiguration())),
            (o.loadPromise = Promise.resolve()),
            (o.backgroundLoading = e.get("backgroundLoading")),
            o.backgroundLoading ||
              e.set("mediaElement", o.mediaPool.getPrimedElement()),
            o
          );
        }
        var n, o, r;
        return (
          (function (e, t) {
            if ("function" != typeof t && null !== t)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (e.prototype = Object.create(t && t.prototype, {
              constructor: { value: e, writable: !0, configurable: !0 },
            })),
              t && D(e, t);
          })(t, e),
          (n = t),
          (o = [
            {
              key: "setActiveItem",
              value: function (e) {
                var t = this,
                  n = this.model,
                  i = n.get("playlist")[e];
                (n.attributes.itemReady = !1), n.setActiveItem(e);
                var o = Q(i);
                if (!o) return Promise.reject(new m.n(m.k, m.h));
                var r = this.background,
                  a = this.mediaController;
                if (r.isNext(i))
                  return (
                    this._destroyActiveMedia(),
                    (this.loadPromise = this._activateBackgroundMedia()),
                    this.loadPromise
                  );
                if ((this._destroyBackgroundMedia(), a)) {
                  if (
                    n.get("castActive") ||
                    this._providerCanPlay(a.provider, o)
                  )
                    return (
                      (this.loadPromise = Promise.resolve(a)),
                      (a.activeItem = i),
                      this._setActiveMedia(a),
                      this.loadPromise
                    );
                  this._destroyActiveMedia();
                }
                var s = n.mediaModel;
                return (
                  (this.loadPromise = this._setupMediaController(o)
                    .then(function (e) {
                      if (s === n.mediaModel)
                        return (e.activeItem = i), t._setActiveMedia(e), e;
                    })
                    .catch(function (e) {
                      throw (t._destroyActiveMedia(), e);
                    })),
                  this.loadPromise
                );
              },
            },
            {
              key: "setAttached",
              value: function (e) {
                var t = this.mediaController;
                if (((this.attached = e), t)) {
                  if (!e) {
                    var n = t.detach(),
                      i = t.item,
                      o = t.mediaModel.get("position");
                    return o && (i.starttime = o), n;
                  }
                  t.attach();
                }
              },
            },
            {
              key: "playVideo",
              value: function (e) {
                var t,
                  n = this,
                  i = this.mediaController,
                  o = this.model;
                if (!o.get("playlistItem"))
                  return Promise.reject(new Error("No media"));
                if ((e || (e = o.get("playReason")), i)) t = i.play(e);
                else {
                  o.set(d.bb, d.jb);
                  var r = b(function (t) {
                    if (
                      n.mediaController &&
                      n.mediaController.mediaModel === t.mediaModel
                    )
                      return t.play(e);
                    throw new Error("Playback cancelled.");
                  });
                  t = this.loadPromise
                    .catch(function (e) {
                      throw (r.cancel(), e);
                    })
                    .then(r.async);
                }
                return t;
              },
            },
            {
              key: "stopVideo",
              value: function () {
                var e = this.mediaController,
                  t = this.model,
                  n = t.get("playlist")[t.get("item")];
                (t.attributes.playlistItem = n), t.resetItem(n), e && e.stop();
              },
            },
            {
              key: "preloadVideo",
              value: function () {
                var e = this.background,
                  t = this.mediaController || e.currentMedia;
                t && t.preload();
              },
            },
            {
              key: "pause",
              value: function () {
                var e = this.mediaController;
                e && e.pause();
              },
            },
            {
              key: "castVideo",
              value: function (e, t) {
                var n = this.model;
                n.attributes.itemReady = !1;
                var o = Object(i.g)({}, t),
                  r = (o.starttime = n.mediaModel.get("currentTime"));
                this._destroyActiveMedia();
                var a = new N(e, n);
                (a.activeItem = o),
                  this._setActiveMedia(a),
                  n.mediaModel.set("currentTime", r);
              },
            },
            {
              key: "stopCast",
              value: function () {
                var e = this.model,
                  t = e.get("item");
                (e.get("playlist")[t].starttime = e.mediaModel.get(
                  "currentTime"
                )),
                  this.stopVideo(),
                  this.setActiveItem(t);
              },
            },
            {
              key: "backgroundActiveMedia",
              value: function () {
                this.adPlaying = !0;
                var e = this.background,
                  t = this.mediaController;
                t &&
                  (e.currentMedia &&
                    this._destroyMediaController(e.currentMedia),
                  (t.background = !0),
                  (e.currentMedia = t),
                  (this.mediaController = null));
              },
            },
            {
              key: "restoreBackgroundMedia",
              value: function () {
                this.adPlaying = !1;
                var e = this.background,
                  t = this.mediaController,
                  n = e.currentMedia;
                if (n) {
                  if (t)
                    return (
                      this._destroyMediaController(n),
                      void (e.currentMedia = null)
                    );
                  var i = n.mediaModel.attributes;
                  i.mediaState === d.mb
                    ? (i.mediaState = d.ob)
                    : i.mediaState !== d.ob && (i.mediaState = d.jb),
                    this._setActiveMedia(n),
                    (n.background = !1),
                    (e.currentMedia = null);
                }
              },
            },
            {
              key: "backgroundLoad",
              value: function (e) {
                var t = this.background,
                  n = Q(e);
                t.setNext(
                  e,
                  this._setupMediaController(n)
                    .then(function (t) {
                      return (t.activeItem = e), t.preload(), t;
                    })
                    .catch(function () {
                      t.clearNext();
                    })
                );
              },
            },
            {
              key: "forwardEvents",
              value: function () {
                var e = this.mediaController;
                e && U(e, this);
              },
            },
            {
              key: "routeEvents",
              value: function (e) {
                var t = this.mediaController;
                t && (t.off(), e && U(t, e));
              },
            },
            {
              key: "destroy",
              value: function () {
                this.off(),
                  this._destroyBackgroundMedia(),
                  this._destroyActiveMedia();
              },
            },
            {
              key: "_setActiveMedia",
              value: function (e) {
                var t = this.model,
                  n = e.mediaModel,
                  i = e.provider;
                !(function (e, t) {
                  var n = e.get("mediaContainer");
                  n
                    ? (t.container = n)
                    : e.once("change:mediaContainer", function (e, n) {
                        t.container = n;
                      });
                })(t, e),
                  (this.mediaController = e),
                  t.set("mediaElement", e.mediaElement),
                  t.setMediaModel(n),
                  t.setProvider(i),
                  U(e, this),
                  t.set("itemReady", !0);
              },
            },
            {
              key: "_destroyActiveMedia",
              value: function () {
                var e = this.mediaController,
                  t = this.model;
                e &&
                  (e.detach(),
                  this._destroyMediaController(e),
                  t.resetProvider(),
                  (this.mediaController = null));
              },
            },
            {
              key: "_destroyBackgroundMedia",
              value: function () {
                var e = this.background;
                this._destroyMediaController(e.currentMedia),
                  (e.currentMedia = null),
                  this._destroyBackgroundLoadingMedia();
              },
            },
            {
              key: "_destroyMediaController",
              value: function (e) {
                var t = this.mediaPool;
                e && (t.recycle(e.mediaElement), e.destroy());
              },
            },
            {
              key: "_setupMediaController",
              value: function (e) {
                var t = this,
                  n = this.model,
                  i = this.providers,
                  o = function (e) {
                    return new N(
                      new e(n.get("id"), n.getConfiguration(), t.primedElement),
                      n
                    );
                  },
                  r = i.choose(e),
                  a = r.provider,
                  s = r.name;
                return a
                  ? Promise.resolve(o(a))
                  : i.load(s).then(function (e) {
                      return o(e);
                    });
              },
            },
            {
              key: "_activateBackgroundMedia",
              value: function () {
                var e = this,
                  t = this.background,
                  n = this.background.nextLoadPromise,
                  i = this.model;
                return (
                  this._destroyMediaController(t.currentMedia),
                  (t.currentMedia = null),
                  n.then(function (n) {
                    if (n)
                      return (
                        t.clearNext(),
                        e.adPlaying
                          ? ((i.attributes.itemReady = !0),
                            (t.currentMedia = n))
                          : (e._setActiveMedia(n), (n.background = !1)),
                        n
                      );
                  })
                );
              },
            },
            {
              key: "_destroyBackgroundLoadingMedia",
              value: function () {
                var e = this,
                  t = this.background,
                  n = this.background.nextLoadPromise;
                n &&
                  n.then(function (n) {
                    e._destroyMediaController(n), t.clearNext();
                  });
              },
            },
            {
              key: "_providerCanPlay",
              value: function (e, t) {
                var n = this.providers.choose(t).provider;
                return n && e && e instanceof n;
              },
            },
            {
              key: "audioTrack",
              get: function () {
                var e = this.mediaController;
                return e ? e.audioTrack : -1;
              },
              set: function (e) {
                var t = this.mediaController;
                t && (t.audioTrack = parseInt(e, 10) || 0);
              },
            },
            {
              key: "audioTracks",
              get: function () {
                var e = this.mediaController;
                if (e) return e.audioTracks;
              },
            },
            {
              key: "beforeComplete",
              get: function () {
                var e = this.mediaController,
                  t = this.background.currentMedia;
                return !(!e && !t) && (e ? e.beforeComplete : t.beforeComplete);
              },
            },
            {
              key: "primedElement",
              get: function () {
                return this.backgroundLoading
                  ? this.mediaPool.getPrimedElement()
                  : this.model.get("mediaElement");
              },
            },
            {
              key: "quality",
              get: function () {
                return this.mediaController ? this.mediaController.quality : -1;
              },
              set: function (e) {
                var t = this.mediaController;
                t && (t.quality = parseInt(e, 10) || 0);
              },
            },
            {
              key: "qualities",
              get: function () {
                var e = this.mediaController;
                return e ? e.qualities : null;
              },
            },
            {
              key: "controls",
              set: function (e) {
                var t = this.mediaController;
                t && (t.controls = e);
              },
            },
            {
              key: "mute",
              set: function (e) {
                var t = this.background,
                  n = this.mediaController,
                  i = this.mediaPool;
                n && (n.mute = e),
                  t.currentMedia && (t.currentMedia.mute = e),
                  i.syncMute(e);
              },
            },
            {
              key: "position",
              set: function (e) {
                var t = this.mediaController;
                t && ((t.item.starttime = e), t.attached && (t.position = e));
              },
            },
            {
              key: "subtitles",
              set: function (e) {
                var t = this.mediaController;
                t && (t.subtitles = e);
              },
            },
            {
              key: "volume",
              set: function (e) {
                var t = this.background,
                  n = this.mediaController,
                  i = this.mediaPool;
                n && (n.volume = e),
                  t.currentMedia && (t.currentMedia.volume = e),
                  i.syncVolume(e);
              },
            },
          ]) && B(n.prototype, o),
          r && B(n, r),
          t
        );
      })(h.a);
      function X(e) {
        return e === d.kb || e === d.lb ? d.mb : e;
      }
      function J(e, t, n) {
        if ((t = X(t)) !== (n = X(n))) {
          var i = t.replace(/(?:ing|d)$/, ""),
            o = {
              type: i,
              newstate: t,
              oldstate: n,
              reason: (function (e, t) {
                return e === d.jb ? (t === d.qb ? t : d.nb) : t;
              })(t, e.mediaModel.get("mediaState")),
            };
          "play" === i
            ? (o.playReason = e.get("playReason"))
            : "pause" === i && (o.pauseReason = e.get("pauseReason")),
            this.trigger(i, o);
        }
      }
      var G = n(48);
      function K(e) {
        return (K =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (e) {
                return typeof e;
              }
            : function (e) {
                return e &&
                  "function" == typeof Symbol &&
                  e.constructor === Symbol &&
                  e !== Symbol.prototype
                  ? "symbol"
                  : typeof e;
              })(e);
      }
      function $(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function Z(e, t) {
        return !t || ("object" !== K(t) && "function" != typeof t)
          ? (function (e) {
              if (void 0 === e)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return e;
            })(e)
          : t;
      }
      function ee(e, t, n, i) {
        return (ee =
          "undefined" != typeof Reflect && Reflect.set
            ? Reflect.set
            : function (e, t, n, i) {
                var o,
                  r = ie(e, t);
                if (r) {
                  if ((o = Object.getOwnPropertyDescriptor(r, t)).set)
                    return o.set.call(i, n), !0;
                  if (!o.writable) return !1;
                }
                if ((o = Object.getOwnPropertyDescriptor(i, t))) {
                  if (!o.writable) return !1;
                  (o.value = n), Object.defineProperty(i, t, o);
                } else
                  !(function (e, t, n) {
                    t in e
                      ? Object.defineProperty(e, t, {
                          value: n,
                          enumerable: !0,
                          configurable: !0,
                          writable: !0,
                        })
                      : (e[t] = n);
                  })(i, t, n);
                return !0;
              })(e, t, n, i);
      }
      function te(e, t, n, i, o) {
        if (!ee(e, t, n, i || e) && o)
          throw new Error("failed to set property");
        return n;
      }
      function ne(e, t, n) {
        return (ne =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (e, t, n) {
                var i = ie(e, t);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, t);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(e, t, n || e);
      }
      function ie(e, t) {
        for (
          ;
          !Object.prototype.hasOwnProperty.call(e, t) && null !== (e = oe(e));

        );
        return e;
      }
      function oe(e) {
        return (oe = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function re(e, t) {
        return (re =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var ae = (function (e) {
          function t(e, n) {
            var i;
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, t);
            var o,
              r = ((i = Z(this, oe(t).call(this, e, n))).model = new R());
            if (
              ((i.playerModel = e),
              (i.provider = null),
              (i.backgroundLoading = e.get("backgroundLoading")),
              (r.mediaModel.attributes.mediaType = "video"),
              i.backgroundLoading)
            )
              o = n.getAdElement();
            else {
              (o = e.get("mediaElement")),
                (r.attributes.mediaElement = o),
                (r.attributes.mediaSrc = o.src);
              var a = (i.srcResetListener = function () {
                i.srcReset();
              });
              o.addEventListener("emptied", a),
                (o.playbackRate = o.defaultPlaybackRate = 1);
            }
            return (i.mediaPool = Object(G.a)(o, n)), i;
          }
          var n, o, r;
          return (
            (function (e, t) {
              if ("function" != typeof t && null !== t)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (e.prototype = Object.create(t && t.prototype, {
                constructor: { value: e, writable: !0, configurable: !0 },
              })),
                t && re(e, t);
            })(t, e),
            (n = t),
            (o = [
              {
                key: "setup",
                value: function () {
                  var e = this.model,
                    t = this.playerModel,
                    n = this.primedElement,
                    i = t.attributes,
                    o = t.mediaModel;
                  e.setup({
                    id: i.id,
                    volume: i.volume,
                    instreamMode: !0,
                    edition: i.edition,
                    mediaContext: o,
                    mute: i.mute,
                    streamType: "VOD",
                    autostartMuted: i.autostartMuted,
                    autostart: i.autostart,
                    advertising: i.advertising,
                    sdkplatform: i.sdkplatform,
                    skipButton: !1,
                  }),
                    e.on("change:state", J, this),
                    e.on(
                      d.w,
                      function (e) {
                        this.trigger(d.w, e);
                      },
                      this
                    ),
                    n.paused || n.pause();
                },
              },
              {
                key: "setActiveItem",
                value: function (e) {
                  var n = this;
                  return (
                    this.stopVideo(),
                    (this.provider = null),
                    ne(oe(t.prototype), "setActiveItem", this)
                      .call(this, e)
                      .then(function (e) {
                        n._setProvider(e.provider);
                      }),
                    this.playVideo()
                  );
                },
              },
              {
                key: "usePsuedoProvider",
                value: function (e) {
                  (this.provider = e),
                    e &&
                      (this._setProvider(e),
                      e.off(d.w),
                      e.on(
                        d.w,
                        function (e) {
                          this.trigger(d.w, e);
                        },
                        this
                      ));
                },
              },
              {
                key: "_setProvider",
                value: function (e) {
                  var t = this;
                  if (e && this.mediaPool) {
                    var n = this.model,
                      o = this.playerModel,
                      r = "vpaid" === e.type;
                    e.off(),
                      e.on(
                        "all",
                        function (e, t) {
                          (r && e === d.F) ||
                            this.trigger(e, Object(i.g)({}, t, { type: e }));
                        },
                        this
                      );
                    var a = n.mediaModel;
                    e.on(d.bb, function (e) {
                      (e.oldstate = e.oldstate || n.get(d.bb)),
                        a.set("mediaState", e.newstate);
                    }),
                      e.on(d.X, this._nativeFullscreenHandler, this),
                      a.on("change:mediaState", function (e, n) {
                        t._stateHandler(n);
                      }),
                      e.attachMedia(),
                      e.volume(o.get("volume")),
                      e.mute(o.getMute()),
                      e.setPlaybackRate && e.setPlaybackRate(1),
                      o.on(
                        "change:volume",
                        function (e, t) {
                          this.volume = t;
                        },
                        this
                      ),
                      o.on(
                        "change:mute",
                        function (e, t) {
                          (this.mute = t), t || (this.volume = o.get("volume"));
                        },
                        this
                      ),
                      o.on(
                        "change:autostartMuted",
                        function (e, t) {
                          t ||
                            (n.set("autostartMuted", t),
                            (this.mute = o.get("mute")));
                        },
                        this
                      );
                  }
                },
              },
              {
                key: "destroy",
                value: function () {
                  var e = this.model,
                    t = this.mediaPool,
                    n = this.playerModel;
                  e.off();
                  var i = t.getPrimedElement();
                  if (this.backgroundLoading) {
                    t.clean();
                    var o = n.get("mediaContainer");
                    i.parentNode === o && o.removeChild(i);
                  } else
                    i &&
                      (i.removeEventListener("emptied", this.srcResetListener),
                      i.src !== e.get("mediaSrc") && this.srcReset());
                },
              },
              {
                key: "srcReset",
                value: function () {
                  var e = this.playerModel,
                    t = e.get("mediaModel"),
                    n = e.getVideo();
                  t.srcReset(), n && (n.src = null);
                },
              },
              {
                key: "_nativeFullscreenHandler",
                value: function (e) {
                  this.model.trigger(d.X, e),
                    this.trigger(d.y, { fullscreen: e.jwstate });
                },
              },
              {
                key: "_stateHandler",
                value: function (e) {
                  var t = this.model;
                  switch (e) {
                    case d.pb:
                    case d.ob:
                      t.set(d.bb, e);
                  }
                },
              },
              {
                key: "mute",
                set: function (e) {
                  var n = this.mediaController,
                    i = this.model,
                    o = this.provider;
                  i.set("mute", e),
                    te(oe(t.prototype), "mute", e, this, !0),
                    n || o.mute(e);
                },
              },
              {
                key: "volume",
                set: function (e) {
                  var n = this.mediaController,
                    i = this.model,
                    o = this.provider;
                  i.set("volume", e),
                    te(oe(t.prototype), "volume", e, this, !0),
                    n || o.volume(e);
                },
              },
            ]) && $(n.prototype, o),
            r && $(n, r),
            t
          );
        })(Y),
        se = { skipoffset: null, tag: null },
        le = function (e, t, n, o) {
          var r,
            a,
            s,
            l,
            c = this,
            u = this,
            h = new ae(t, o),
            p = 0,
            b = {},
            m = null,
            w = {},
            v = R,
            y = !1,
            j = !1,
            k = !1,
            O = !1,
            x = function (e) {
              j ||
                (((e = e || {}).hasControls = !!t.get("controls")),
                c.trigger(d.z, e),
                h.model.get("state") === d.ob
                  ? e.hasControls && h.playVideo().catch(function () {})
                  : h.pause());
            },
            C = function () {
              j ||
                (h.model.get("state") === d.ob &&
                  t.get("controls") &&
                  (e.setFullscreen(), e.play()));
            };
          function M() {
            h.model.set("playRejected", !0);
          }
          function _() {
            p++, u.loadItem(r).catch(function () {});
          }
          function S(e, t) {
            "complete" !== e &&
              ((t = t || {}),
              w.tag && !t.tag && (t.tag = w.tag),
              this.trigger(e, t),
              ("mediaError" !== e && "error" !== e) ||
                (r && p + 1 < r.length && _()));
          }
          function P(e) {
            var t = e.newstate,
              n = e.oldstate || h.model.get("state");
            n !== t && T(Object(i.g)({ oldstate: n }, b, e));
          }
          function T(t) {
            var n = t.newstate;
            n === d.pb ? e.trigger(d.c, t) : n === d.ob && e.trigger(d.b, t);
          }
          function A(t) {
            var n = t.duration,
              i = t.position,
              o = h.model.mediaModel || h.model;
            o.set("duration", n),
              o.set("position", i),
              l || (l = (Object(g.d)(s, n) || n) - f.b),
              !y && i >= Math.max(l, f.a) && (e.preloadNextItem(), (y = !0));
          }
          function E(e) {
            var t = {};
            w.tag && (t.tag = w.tag), this.trigger(d.F, t), R.call(this, e);
          }
          function R(e) {
            (b = {}),
              r && p + 1 < r.length
                ? _()
                : (e.type === d.F && this.trigger(d.cb, {}), this.destroy());
          }
          function I() {
            j ||
              (n.clickHandler() &&
                n.clickHandler().setAlternateClickHandlers(x, C));
          }
          function L(e) {
            e.width && e.height && n.resizeMedia();
          }
          (this.init = function () {
            if (!k && !j) {
              (k = !0),
                (b = {}),
                h.setup(),
                h.on("all", S, this),
                h.on(d.O, M, this),
                h.on(d.S, A, this),
                h.on(d.F, E, this),
                h.on(d.K, L, this),
                h.on(d.bb, P, this),
                (m = e.detachMedia());
              var i = h.primedElement;
              t.get("mediaContainer").appendChild(i),
                t.set("instream", h),
                h.model.set("state", d.jb);
              var o = n.clickHandler();
              return (
                o && o.setAlternateClickHandlers(function () {}, null),
                this.setText(t.get("localization").loadingAd),
                (O = e.isBeforeComplete() || t.get("state") === d.kb),
                this
              );
            }
          }),
            (this.enableAdsMode = function (i) {
              var o = this;
              if (!k && !j)
                return (
                  e.routeEvents({
                    mediaControllerListener: function (e, t) {
                      o.trigger(e, t);
                    },
                  }),
                  t.set("instream", h),
                  h.model.set("state", d.pb),
                  (function (i) {
                    var o = n.clickHandler();
                    o &&
                      o.setAlternateClickHandlers(function (n) {
                        j ||
                          (((n = n || {}).hasControls = !!t.get("controls")),
                          u.trigger(d.z, n),
                          i &&
                            (t.get("state") === d.ob
                              ? e.playVideo()
                              : (e.pause(),
                                i &&
                                  (e.trigger(d.a, { clickThroughUrl: i }),
                                  window.open(i)))));
                      }, null);
                  })(i),
                  this
                );
            }),
            (this.setEventData = function (e) {
              b = e;
            }),
            (this.setState = function (e) {
              var t = e.newstate,
                n = h.model;
              (e.oldstate = n.get("state")), n.set("state", t), T(e);
            }),
            (this.setTime = function (t) {
              A(t), e.trigger(d.e, t);
            }),
            (this.loadItem = function (e, n) {
              if (j || !k)
                return Promise.reject(new Error("Instream not setup"));
              b = {};
              var o = e;
              Array.isArray(e)
                ? ((a = n || a), (e = (r = e)[p]), a && (n = a[p]))
                : (o = [e]);
              var l = h.model;
              l.set("playlist", o),
                t.set("hideAdsControls", !1),
                (e.starttime = 0),
                u.trigger(d.db, { index: p, item: e }),
                (w = Object(i.g)({}, se, n)),
                I(),
                l.set("skipButton", !1);
              var c =
                !t.get("backgroundLoading") && m
                  ? m.then(function () {
                      return h.setActiveItem(p);
                    })
                  : h.setActiveItem(p);
              return (
                (y = !1),
                void 0 !== (s = e.skipoffset || w.skipoffset) &&
                  u.setupSkipButton(s, w),
                c
              );
            }),
            (this.setupSkipButton = function (e, t, n) {
              var i = h.model;
              (v = n || R),
                i.set("skipMessage", t.skipMessage),
                i.set("skipText", t.skipText),
                i.set("skipOffset", e),
                (i.attributes.skipButton = !1),
                i.set("skipButton", !0);
            }),
            (this.applyProviderListeners = function (e) {
              h.usePsuedoProvider(e), I();
            }),
            (this.play = function () {
              (b = {}), h.playVideo();
            }),
            (this.pause = function () {
              (b = {}), h.pause();
            }),
            (this.skipAd = function (e) {
              var n = t.get("autoPause").pauseAds,
                i = "autostart" === t.get("playReason"),
                o = t.get("viewable");
              !n || i || o || (this.noResume = !0);
              var r = d.d;
              this.trigger(r, e), v.call(this, { type: r });
            }),
            (this.replacePlaylistItem = function (e) {
              j || (t.set("playlistItem", e), h.srcReset());
            }),
            (this.destroy = function () {
              j ||
                ((j = !0),
                this.trigger("destroyed"),
                this.off(),
                n.clickHandler() &&
                  n.clickHandler().revertAlternateClickHandlers(),
                t.off(null, null, h),
                h.off(null, null, u),
                h.destroy(),
                k && h.model && (t.attributes.state = d.ob),
                e.forwardEvents(),
                t.set("instream", null),
                (h = null),
                (b = {}),
                (m = null),
                k &&
                  !t.attributes._destroyed &&
                  (e.attachMedia(),
                  this.noResume || (O ? e.stopVideo() : e.playVideo())));
            }),
            (this.getState = function () {
              return !j && h.model.get("state");
            }),
            (this.setText = function (e) {
              return j ? this : (n.setAltText(e || ""), this);
            }),
            (this.hide = function () {
              j || t.set("hideAdsControls", !0);
            }),
            (this.getMediaElement = function () {
              return j ? null : h.primedElement;
            }),
            (this.setSkipOffset = function (e) {
              (s = e > 0 ? e : null), h && h.model.set("skipOffset", s);
            });
        };
      Object(i.g)(le.prototype, h.a);
      var ce = le,
        ue = n(66),
        de = n(63),
        fe = function (e) {
          var t = this,
            n = [],
            i = {},
            o = 0,
            r = 0;
          function a(e) {
            if (
              ((e.data = e.data || []),
              (e.name = e.label || e.name || e.language),
              (e._id = Object(de.a)(e, n.length)),
              !e.name)
            ) {
              var t = Object(de.b)(e, o);
              (e.name = t.label), (o = t.unknownCount);
            }
            (i[e._id] = e), n.push(e);
          }
          function s() {
            for (
              var e = [{ id: "off", label: "Off" }], t = 0;
              t < n.length;
              t++
            )
              e.push({
                id: n[t]._id,
                label: n[t].name || "Unknown CC",
                language: n[t].language,
              });
            return e;
          }
          function l(t) {
            var i = (r = t),
              o = e.get("captionLabel");
            if ("Off" !== o) {
              for (var a = 0; a < n.length; a++) {
                var s = n[a];
                if (o && o === s.name) {
                  i = a + 1;
                  break;
                }
                s.default || s.defaulttrack || "default" === s._id
                  ? (i = a + 1)
                  : s.autoselect;
              }
              var l;
              (l = i),
                n.length
                  ? e.setVideoSubtitleTrack(l, n)
                  : e.set("captionsIndex", l);
            } else e.set("captionsIndex", 0);
          }
          function c() {
            var t = s();
            u(t) !== u(e.get("captionsList")) &&
              (l(r), e.set("captionsList", t));
          }
          function u(e) {
            return e
              .map(function (e) {
                return "".concat(e.id, "-").concat(e.label);
              })
              .join(",");
          }
          e.on(
            "change:playlistItem",
            function (e) {
              (n = []), (i = {}), (o = 0);
              var t = e.attributes;
              (t.captionsIndex = 0),
                (t.captionsList = s()),
                e.set("captionsTrack", null);
            },
            this
          ),
            e.on(
              "change:itemReady",
              function () {
                var n = e.get("playlistItem").tracks,
                  o = n && n.length;
                if (o && !e.get("renderCaptionsNatively"))
                  for (
                    var r = function (e) {
                        var o,
                          r = n[e];
                        ("subtitles" !== (o = r.kind) && "captions" !== o) ||
                          i[r._id] ||
                          (a(r),
                          Object(ue.c)(
                            r,
                            function (e) {
                              !(function (e, t) {
                                e.data = t;
                              })(r, e);
                            },
                            function (e) {
                              t.trigger(d.tb, e);
                            }
                          ));
                      },
                      s = 0;
                    s < o;
                    s++
                  )
                    r(s);
                c();
              },
              this
            ),
            e.on(
              "change:captionsIndex",
              function (e, t) {
                var i = null;
                0 !== t && (i = n[t - 1]), e.set("captionsTrack", i);
              },
              this
            ),
            (this.setSubtitlesTracks = function (e) {
              if (Array.isArray(e)) {
                if (e.length) {
                  for (var t = 0; t < e.length; t++) a(e[t]);
                  n = Object.keys(i).map(function (e) {
                    return i[e];
                  });
                } else (n = []), (i = {}), (o = 0);
                c();
              }
            }),
            (this.selectDefaultIndex = l),
            (this.getCurrentIndex = function () {
              return e.get("captionsIndex");
            }),
            (this.getCaptionsList = function () {
              return e.get("captionsList");
            }),
            (this.destroy = function () {
              this.off(null, null, this);
            });
        };
      Object(i.g)(fe.prototype, h.a);
      var ge = fe,
        he = function (e, t) {
          return (
            '<div id="'
              .concat(
                e,
                '" class="jwplayer jw-reset jw-state-setup" tabindex="0" aria-label="'
              )
              .concat(t || "", '" role="application">') +
            '<div class="jw-aspect jw-reset"></div><div class="jw-wrapper jw-reset"><div class="jw-top jw-reset"></div><div class="jw-aspect jw-reset"></div><div class="jw-media jw-reset"></div><div class="jw-preview jw-reset"></div><div class="jw-title jw-reset-text" dir="auto"><div class="jw-title-primary jw-reset-text"></div><div class="jw-title-secondary jw-reset-text"></div></div><div class="jw-overlays jw-reset"></div><div class="jw-hidden-accessibility"><span class="jw-time-update" aria-live="assertive"></span><span class="jw-volume-update" aria-live="assertive"></span></div></div></div>'
          );
        },
        pe = n(35),
        be = 44,
        me = function (e) {
          var t = e.get("height");
          if (e.get("aspectratio")) return !1;
          if ("string" == typeof t && t.indexOf("%") > -1) return !1;
          var n = 1 * t || NaN;
          return (
            !!(n = isNaN(n) ? e.get("containerHeight") : n) && n && n <= be
          );
        },
        we = n(54);
      function ve(e, t) {
        if (e.get("fullscreen")) return 1;
        if (!e.get("activeTab")) return 0;
        if (e.get("isFloating")) return 1;
        var n = e.get("intersectionRatio");
        return void 0 === n &&
          ((n = (function (e) {
            var t = document.documentElement,
              n = document.body,
              i = {
                top: 0,
                left: 0,
                right: t.clientWidth || n.clientWidth,
                width: t.clientWidth || n.clientWidth,
                bottom: t.clientHeight || n.clientHeight,
                height: t.clientHeight || n.clientHeight,
              };
            if (!n.contains(e)) return 0;
            if ("none" === window.getComputedStyle(e).display) return 0;
            var o = ye(e);
            if (!o) return 0;
            var r = o,
              a = e.parentNode,
              s = !1;
            for (; !s; ) {
              var l = null;
              if (
                (a === n || a === t || 1 !== a.nodeType
                  ? ((s = !0), (l = i))
                  : "visible" !== window.getComputedStyle(a).overflow &&
                    (l = ye(a)),
                l &&
                  ((c = l),
                  (u = r),
                  (d = void 0),
                  (f = void 0),
                  (g = void 0),
                  (h = void 0),
                  (p = void 0),
                  (b = void 0),
                  (d = Math.max(c.top, u.top)),
                  (f = Math.min(c.bottom, u.bottom)),
                  (g = Math.max(c.left, u.left)),
                  (h = Math.min(c.right, u.right)),
                  (b = f - d),
                  !(r = (p = h - g) >= 0 &&
                    b >= 0 && {
                      top: d,
                      bottom: f,
                      left: g,
                      right: h,
                      width: p,
                      height: b,
                    })))
              )
                return 0;
              a = a.parentNode;
            }
            var c, u, d, f, g, h, p, b;
            var m = o.width * o.height,
              w = r.width * r.height;
            return m ? w / m : 0;
          })(t)),
          window.top !== window.self && n)
          ? 0
          : n;
      }
      function ye(e) {
        try {
          return e.getBoundingClientRect();
        } catch (e) {}
      }
      var je = n(49),
        ke = n(42),
        Oe = n(58),
        xe = n(10);
      var Ce = n(32),
        Me = n(5),
        _e = n(6),
        Se = [
          "fullscreenchange",
          "webkitfullscreenchange",
          "mozfullscreenchange",
          "MSFullscreenChange",
        ],
        Pe = function (e, t, n) {
          for (
            var i =
                e.requestFullscreen ||
                e.webkitRequestFullscreen ||
                e.webkitRequestFullScreen ||
                e.mozRequestFullScreen ||
                e.msRequestFullscreen,
              o =
                t.exitFullscreen ||
                t.webkitExitFullscreen ||
                t.webkitCancelFullScreen ||
                t.mozCancelFullScreen ||
                t.msExitFullscreen,
              r = !(!i || !o),
              a = Se.length;
            a--;

          )
            t.addEventListener(Se[a], n);
          return {
            events: Se,
            supportsDomFullscreen: function () {
              return r;
            },
            requestFullscreen: function () {
              i.call(e, { navigationUI: "hide" });
            },
            exitFullscreen: function () {
              null !== this.fullscreenElement() && o.apply(t);
            },
            fullscreenElement: function () {
              var e = t.fullscreenElement,
                n = t.webkitCurrentFullScreenElement,
                i = t.mozFullScreenElement,
                o = t.msFullscreenElement;
              return null === e ? e : e || n || i || o;
            },
            destroy: function () {
              for (var e = Se.length; e--; ) t.removeEventListener(Se[e], n);
            },
          };
        },
        Te = n(40);
      function Ae(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var Ee = (function () {
          function e(t, n) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              Object(i.g)(this, h.a),
              this.revertAlternateClickHandlers(),
              (this.domElement = n),
              (this.model = t),
              (this.ui = new Te.a(n)
                .on("click tap", this.clickHandler, this)
                .on(
                  "doubleClick doubleTap",
                  function () {
                    this.alternateDoubleClickHandler
                      ? this.alternateDoubleClickHandler()
                      : this.trigger("doubleClick");
                  },
                  this
                ));
          }
          var t, n, o;
          return (
            (t = e),
            (n = [
              {
                key: "destroy",
                value: function () {
                  this.ui &&
                    (this.ui.destroy(),
                    (this.ui = this.domElement = this.model = null),
                    this.revertAlternateClickHandlers());
                },
              },
              {
                key: "clickHandler",
                value: function (e) {
                  this.model.get("flashBlocked") ||
                    (this.alternateClickHandler
                      ? this.alternateClickHandler(e)
                      : this.trigger(e.type === d.n ? "click" : "tap"));
                },
              },
              {
                key: "element",
                value: function () {
                  return this.domElement;
                },
              },
              {
                key: "setAlternateClickHandlers",
                value: function (e, t) {
                  (this.alternateClickHandler = e),
                    (this.alternateDoubleClickHandler = t || null);
                },
              },
              {
                key: "revertAlternateClickHandlers",
                value: function () {
                  (this.alternateClickHandler = null),
                    (this.alternateDoubleClickHandler = null);
                },
              },
            ]) && Ae(t.prototype, n),
            o && Ae(t, o),
            e
          );
        })(),
        Re = n(59),
        Ie = function (e, t) {
          var n = t ? " jw-hide" : "";
          return '<div class="jw-logo jw-logo-'
            .concat(e)
            .concat(n, ' jw-reset"></div>');
        },
        Le = {
          linktarget: "_blank",
          margin: 8,
          hide: !1,
          position: "top-right",
        };
      function Ve(e) {
        var t, n;
        Object(i.g)(this, h.a);
        var o = new Image();
        (this.setup = function () {
          ((n = Object(i.g)({}, Le, e.get("logo"))).position =
            n.position || Le.position),
            (n.hide = "true" === n.hide.toString()),
            n.file &&
              "control-bar" !== n.position &&
              (t || (t = Object(Me.e)(Ie(n.position, n.hide))),
              e.set("logo", n),
              (o.onload = function () {
                var i = this.height,
                  o = this.width,
                  r = { backgroundImage: 'url("' + this.src + '")' };
                if (n.margin !== Le.margin) {
                  var a = /(\w+)-(\w+)/.exec(n.position);
                  3 === a.length &&
                    ((r["margin-" + a[1]] = n.margin),
                    (r["margin-" + a[2]] = n.margin));
                }
                var s = 0.15 * e.get("containerHeight"),
                  l = 0.15 * e.get("containerWidth");
                if (i > s || o > l) {
                  var c = o / i;
                  l / s > c ? ((i = s), (o = s * c)) : ((o = l), (i = l / c));
                }
                (r.width = Math.round(o)),
                  (r.height = Math.round(i)),
                  Object(xe.d)(t, r),
                  e.set("logoWidth", r.width);
              }),
              (o.src = n.file),
              n.link &&
                (t.setAttribute("tabindex", "0"),
                t.setAttribute("aria-label", e.get("localization").logo)),
              (this.ui = new Te.a(t).on(
                "click tap enter",
                function (e) {
                  e && e.stopPropagation && e.stopPropagation(),
                    this.trigger(d.A, {
                      link: n.link,
                      linktarget: n.linktarget,
                    });
                },
                this
              )));
        }),
          (this.setContainer = function (e) {
            t && e.appendChild(t);
          }),
          (this.element = function () {
            return t;
          }),
          (this.position = function () {
            return n.position;
          }),
          (this.destroy = function () {
            (o.onload = null), this.ui && this.ui.destroy();
          });
      }
      var Fe = function (e) {
        (this.model = e), (this.image = null);
      };
      Object(i.g)(Fe.prototype, {
        setup: function (e) {
          this.el = e;
        },
        setImage: function (e) {
          var t = this.image;
          t && (t.onload = null), (this.image = null);
          var n = "";
          "string" == typeof e &&
            ((n = 'url("' + e + '")'),
            ((t = this.image = new Image()).src = e)),
            Object(xe.d)(this.el, { backgroundImage: n });
        },
        resize: function (e, t, n) {
          if ("uniform" === n) {
            if (
              (e && (this.playerAspectRatio = e / t),
              !this.playerAspectRatio ||
                !this.image ||
                ("complete" !== (s = this.model.get("state")) &&
                  "idle" !== s &&
                  "error" !== s &&
                  "buffering" !== s))
            )
              return;
            var i = this.image,
              o = null;
            if (i) {
              if (0 === i.width) {
                var r = this;
                return void (i.onload = function () {
                  r.resize(e, t, n);
                });
              }
              var a = i.width / i.height;
              Math.abs(this.playerAspectRatio - a) < 0.09 && (o = "cover");
            }
            Object(xe.d)(this.el, { backgroundSize: o });
          }
          var s;
        },
        element: function () {
          return this.el;
        },
      });
      var ze = Fe,
        Ne = function (e) {
          this.model = e.player;
        };
      Object(i.g)(Ne.prototype, {
        hide: function () {
          Object(xe.d)(this.el, { display: "none" });
        },
        show: function () {
          Object(xe.d)(this.el, { display: "" });
        },
        setup: function (e) {
          this.el = e;
          var t = this.el.getElementsByTagName("div");
          (this.title = t[0]),
            (this.description = t[1]),
            this.model.on("change:logoWidth", this.update, this),
            this.model.change("playlistItem", this.playlistItem, this);
        },
        update: function (e) {
          var t = {},
            n = e.get("logo");
          if (n) {
            var i = 1 * ("" + n.margin).replace("px", ""),
              o = e.get("logoWidth") + (isNaN(i) ? 0 : i + 10);
            "top-left" === n.position
              ? (t.paddingLeft = o)
              : "top-right" === n.position && (t.paddingRight = o);
          }
          Object(xe.d)(this.el, t);
        },
        playlistItem: function (e, t) {
          if (t)
            if (e.get("displaytitle") || e.get("displaydescription")) {
              var n = "",
                i = "";
              t.title && e.get("displaytitle") && (n = t.title),
                t.description &&
                  e.get("displaydescription") &&
                  (i = t.description),
                this.updateText(n, i);
            } else this.hide();
        },
        updateText: function (e, t) {
          Object(Me.q)(this.title, e),
            Object(Me.q)(this.description, t),
            this.title.firstChild || this.description.firstChild
              ? this.show()
              : this.hide();
        },
        element: function () {
          return this.el;
        },
      });
      var He = Ne;
      function Be(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var qe,
        De = (function () {
          function e(t) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              (this.container = t),
              (this.input = t.querySelector(".jw-media"));
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "disable",
                value: function () {
                  this.ui && (this.ui.destroy(), (this.ui = null));
                },
              },
              {
                key: "enable",
                value: function () {
                  var e,
                    t,
                    n,
                    i,
                    o = this.container,
                    r = this.input,
                    a = (this.ui = new Te.a(r, { preventScrolling: !0 })
                      .on("dragStart", function () {
                        (e = o.offsetLeft),
                          (t = o.offsetTop),
                          (n = window.innerHeight),
                          (i = window.innerWidth);
                      })
                      .on("drag", function (r) {
                        var s = Math.max(e + r.pageX - a.startX, 0),
                          l = Math.max(t + r.pageY - a.startY, 0),
                          c = Math.max(i - (s + o.clientWidth), 0),
                          u = Math.max(n - (l + o.clientHeight), 0);
                        0 === c ? (s = "auto") : (c = "auto"),
                          0 === l ? (u = "auto") : (l = "auto"),
                          Object(xe.d)(o, {
                            left: s,
                            right: c,
                            top: l,
                            bottom: u,
                            margin: 0,
                          });
                      })
                      .on("dragEnd", function () {
                        e = t = i = n = null;
                      }));
                },
              },
            ]) && Be(t.prototype, n),
            i && Be(t, i),
            e
          );
        })(),
        We = n(55);
      n(69);
      var Ue = v.OS.mobile,
        Qe = v.Browser.ie,
        Ye = null;
      var Xe = function (e, t) {
        var n,
          o,
          r,
          a,
          s = this,
          l = Object(i.g)(this, h.a, { isSetup: !1, api: e, model: t }),
          c = t.get("localization"),
          u = Object(Me.e)(he(t.get("id"), c.player)),
          f = u.querySelector(".jw-wrapper"),
          p = u.querySelector(".jw-media"),
          b = new De(f),
          m = new ze(t, e),
          w = new He(t),
          y = new Re.b(t);
        y.on("all", l.trigger, l);
        var j = -1,
          k = -1,
          O = -1,
          x = t.get("floating");
        this.dismissible = x && x.dismissible;
        var C,
          M,
          _,
          S = !1,
          P = {},
          T = null,
          A = null;
        function E() {
          return Ue && !Object(Me.f)();
        }
        function R() {
          Object(ke.a)(k), (k = Object(ke.b)(I));
        }
        function I() {
          l.isSetup && (l.updateBounds(), l.updateStyles(), l.checkResized());
        }
        function L(e, n) {
          if (Object(i.r)(e) && Object(i.r)(n)) {
            var o = Object(Oe.a)(e);
            Object(Oe.b)(u, o);
            var r = o < 2;
            Object(Me.v)(u, "jw-flag-small-player", r),
              Object(Me.v)(u, "jw-orientation-portrait", n > e);
          }
          if (t.get("controls")) {
            var a = me(t);
            Object(Me.v)(u, "jw-flag-audio-player", a), t.set("audioMode", a);
          }
        }
        function V() {
          t.set("visibility", ve(t, u));
        }
        (this.updateBounds = function () {
          Object(ke.a)(k);
          var e = t.get("isFloating") ? f : u,
            n = document.body.contains(e),
            i = Object(Me.c)(e),
            a = Math.round(i.width),
            s = Math.round(i.height);
          if (((P = Object(Me.c)(u)), a === o && s === r))
            return (o && r) || R(), void t.set("inDom", n);
          (a && s) || (o && r) || R(),
            (a || s || n) &&
              (t.set("containerWidth", a), t.set("containerHeight", s)),
            t.set("inDom", n),
            n && we.a.observe(u);
        }),
          (this.updateStyles = function () {
            var e = t.get("containerWidth"),
              n = t.get("containerHeight");
            L(e, n), A && A.resize(e, n), Z(e, n), y.resize(), x && B();
          }),
          (this.checkResized = function () {
            var e = t.get("containerWidth"),
              n = t.get("containerHeight"),
              i = t.get("isFloating");
            if (e !== o || n !== r) {
              this.resizeListener ||
                (this.resizeListener = new We.a(f, this, t)),
                (o = e),
                (r = n),
                l.trigger(d.hb, { width: e, height: n });
              var s = Object(Oe.a)(e);
              T !== s && ((T = s), l.trigger(d.j, { breakpoint: T }));
            }
            i !== a && ((a = i), l.trigger(d.x, { floating: i }), V());
          }),
          (this.responsiveListener = R),
          (this.setup = function () {
            m.setup(u.querySelector(".jw-preview")),
              w.setup(u.querySelector(".jw-title")),
              (n = new Ve(t)).setup(),
              n.setContainer(f),
              n.on(d.A, G),
              y.setup(u.id, t.get("captions")),
              w.element().parentNode.insertBefore(y.element(), w.element()),
              (C = (function (e, t, n) {
                var i = new Ee(t, n),
                  o = t.get("controls");
                i.on({
                  click: function () {
                    l.trigger(d.p),
                      A &&
                        (ce()
                          ? A.settingsMenu.close()
                          : ue()
                          ? A.infoOverlay.close()
                          : e.playToggle({ reason: "interaction" }));
                  },
                  tap: function () {
                    l.trigger(d.p),
                      ce() && A.settingsMenu.close(),
                      ue() && A.infoOverlay.close();
                    var n = t.get("state");
                    if (
                      (o &&
                        (n === d.mb ||
                          n === d.kb ||
                          (t.get("instream") && n === d.ob)) &&
                        e.playToggle({ reason: "interaction" }),
                      o && n === d.ob)
                    ) {
                      if (
                        t.get("instream") ||
                        t.get("castActive") ||
                        "audio" === t.get("mediaType")
                      )
                        return;
                      Object(Me.v)(u, "jw-flag-controls-hidden"),
                        l.dismissible &&
                          Object(Me.v)(
                            u,
                            "jw-floating-dismissible",
                            Object(Me.i)(u, "jw-flag-controls-hidden")
                          ),
                        y.renderCues(!0);
                    } else A && (A.showing ? A.userInactive() : A.userActive());
                  },
                  doubleClick: function () {
                    return A && e.setFullscreen();
                  },
                }),
                  Ue ||
                    (u.addEventListener("mousemove", U),
                    u.addEventListener("mouseover", Q),
                    u.addEventListener("mouseout", Y));
                return i;
              })(e, t, p)),
              (_ = new Te.a(u).on("click", function () {})),
              (M = Pe(u, document, te)),
              t.on("change:hideAdsControls", function (e, t) {
                Object(Me.v)(u, "jw-flag-ads-hide-controls", t);
              }),
              t.on("change:scrubbing", function (e, t) {
                Object(Me.v)(u, "jw-flag-dragging", t);
              }),
              t.on("change:playRejected", function (e, t) {
                Object(Me.v)(u, "jw-flag-play-rejected", t);
              }),
              t.on(d.X, ee),
              t.on("change:".concat(d.U), function () {
                Z(), y.resize();
              }),
              t.player.on("change:errorEvent", re),
              t.change("stretching", X);
            var i = t.get("width"),
              o = t.get("height"),
              r = $(i, o);
            Object(xe.d)(u, r),
              t.change("aspectratio", J),
              L(i, o),
              t.get("controls") ||
                (Object(Me.a)(u, "jw-flag-controls-hidden"),
                Object(Me.o)(u, "jw-floating-dismissible")),
              Qe && Object(Me.a)(u, "jw-ie");
            var a = t.get("skin") || {};
            a.name && Object(Me.p)(u, /jw-skin-\S+/, "jw-skin-" + a.name);
            var s = (function (e) {
              e || (e = {});
              var t = e.active,
                n = e.inactive,
                i = e.background,
                o = {};
              return (
                (o.controlbar = (function (e) {
                  if (e || t || n || i) {
                    var o = {};
                    return (
                      (e = e || {}),
                      (o.iconsActive = e.iconsActive || t),
                      (o.icons = e.icons || n),
                      (o.text = e.text || n),
                      (o.background = e.background || i),
                      o
                    );
                  }
                })(e.controlbar)),
                (o.timeslider = (function (e) {
                  if (e || t) {
                    var n = {};
                    return (
                      (e = e || {}),
                      (n.progress = e.progress || t),
                      (n.rail = e.rail),
                      n
                    );
                  }
                })(e.timeslider)),
                (o.menus = (function (e) {
                  if (e || t || n || i) {
                    var o = {};
                    return (
                      (e = e || {}),
                      (o.text = e.text || n),
                      (o.textActive = e.textActive || t),
                      (o.background = e.background || i),
                      o
                    );
                  }
                })(e.menus)),
                (o.tooltips = (function (e) {
                  if (e || n || i) {
                    var t = {};
                    return (
                      (e = e || {}),
                      (t.text = e.text || n),
                      (t.background = e.background || i),
                      t
                    );
                  }
                })(e.tooltips)),
                o
              );
            })(a);
            !(function (e, t) {
              var n;
              function i(t, n, i, o) {
                if (i) {
                  t = Object(g.f)(t, "#" + e + (o ? "" : " "));
                  var r = {};
                  (r[n] = i), Object(xe.b)(t.join(", "), r, e);
                }
              }
              t &&
                (t.controlbar &&
                  (function (t) {
                    i(
                      [
                        ".jw-controlbar .jw-icon-inline.jw-text",
                        ".jw-title-primary",
                        ".jw-title-secondary",
                      ],
                      "color",
                      t.text
                    ),
                      t.icons &&
                        (i(
                          [
                            ".jw-button-color:not(.jw-icon-cast)",
                            ".jw-button-color.jw-toggle.jw-off:not(.jw-icon-cast)",
                          ],
                          "color",
                          t.icons
                        ),
                        i(
                          [".jw-display-icon-container .jw-button-color"],
                          "color",
                          t.icons
                        ),
                        Object(xe.b)(
                          "#".concat(
                            e,
                            " .jw-icon-cast google-cast-launcher.jw-off"
                          ),
                          "{--disconnected-color: ".concat(t.icons, "}"),
                          e
                        ));
                    t.iconsActive &&
                      (i(
                        [
                          ".jw-display-icon-container .jw-button-color:hover",
                          ".jw-display-icon-container .jw-button-color:focus",
                        ],
                        "color",
                        t.iconsActive
                      ),
                      i(
                        [
                          ".jw-button-color.jw-toggle:not(.jw-icon-cast)",
                          ".jw-button-color:hover:not(.jw-icon-cast)",
                          ".jw-button-color:focus:not(.jw-icon-cast)",
                          ".jw-button-color.jw-toggle.jw-off:hover:not(.jw-icon-cast)",
                        ],
                        "color",
                        t.iconsActive
                      ),
                      i([".jw-svg-icon-buffer"], "fill", t.icons),
                      Object(xe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast:hover google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(xe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast:focus google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(xe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast google-cast-launcher.jw-off:focus"
                        ),
                        "{--disconnected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(xe.b)(
                        "#".concat(e, " .jw-icon-cast google-cast-launcher"),
                        "{--connected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(xe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast google-cast-launcher:focus"
                        ),
                        "{--connected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(xe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast:hover google-cast-launcher"
                        ),
                        "{--connected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(xe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast:focus google-cast-launcher"
                        ),
                        "{--connected-color: ".concat(t.iconsActive, "}"),
                        e
                      ));
                    i(
                      [
                        " .jw-settings-topbar",
                        ":not(.jw-state-idle) .jw-controlbar",
                        ".jw-flag-audio-player .jw-controlbar",
                      ],
                      "background",
                      t.background,
                      !0
                    );
                  })(t.controlbar),
                t.timeslider &&
                  (function (e) {
                    var t = e.progress;
                    "none" !== t &&
                      (i([".jw-progress", ".jw-knob"], "background-color", t),
                      i(
                        [".jw-buffer"],
                        "background-color",
                        Object(xe.c)(t, 50)
                      ));
                    i([".jw-rail"], "background-color", e.rail),
                      i(
                        [
                          ".jw-background-color.jw-slider-time",
                          ".jw-slider-time .jw-cue",
                        ],
                        "background-color",
                        e.background
                      );
                  })(t.timeslider),
                t.menus &&
                  (i(
                    [
                      ".jw-option",
                      ".jw-toggle.jw-off",
                      ".jw-skip .jw-skip-icon",
                      ".jw-nextup-tooltip",
                      ".jw-nextup-close",
                      ".jw-settings-content-item",
                      ".jw-related-title",
                    ],
                    "color",
                    (n = t.menus).text
                  ),
                  i(
                    [
                      ".jw-option.jw-active-option",
                      ".jw-option:not(.jw-active-option):hover",
                      ".jw-option:not(.jw-active-option):focus",
                      ".jw-settings-content-item:hover",
                      ".jw-nextup-tooltip:hover",
                      ".jw-nextup-tooltip:focus",
                      ".jw-nextup-close:hover",
                    ],
                    "color",
                    n.textActive
                  ),
                  i(
                    [".jw-nextup", ".jw-settings-menu"],
                    "background",
                    n.background
                  )),
                t.tooltips &&
                  (function (e) {
                    i(
                      [
                        ".jw-skip",
                        ".jw-tooltip .jw-text",
                        ".jw-time-tip .jw-text",
                      ],
                      "background-color",
                      e.background
                    ),
                      i([".jw-time-tip", ".jw-tooltip"], "color", e.background),
                      i([".jw-skip"], "border", "none"),
                      i(
                        [
                          ".jw-skip .jw-text",
                          ".jw-skip .jw-icon",
                          ".jw-time-tip .jw-text",
                          ".jw-tooltip .jw-text",
                        ],
                        "color",
                        e.text
                      );
                  })(t.tooltips),
                t.menus &&
                  (function (t) {
                    if (t.textActive) {
                      var n = {
                        color: t.textActive,
                        borderColor: t.textActive,
                        stroke: t.textActive,
                      };
                      Object(xe.b)("#".concat(e, " .jw-color-active"), n, e),
                        Object(xe.b)(
                          "#".concat(e, " .jw-color-active-hover:hover"),
                          n,
                          e
                        );
                    }
                    if (t.text) {
                      var i = {
                        color: t.text,
                        borderColor: t.text,
                        stroke: t.text,
                      };
                      Object(xe.b)("#".concat(e, " .jw-color-inactive"), i, e),
                        Object(xe.b)(
                          "#".concat(e, " .jw-color-inactive-hover:hover"),
                          i,
                          e
                        );
                    }
                  })(t.menus));
            })(t.get("id"), s),
              t.set("mediaContainer", p),
              t.set("iFrame", v.Features.iframe),
              t.set("activeTab", Object(je.a)()),
              t.set("touchMode", Ue && ("string" == typeof o || o >= be)),
              we.a.add(this),
              t.get("enableGradient") &&
                !Qe &&
                Object(Me.a)(u, "jw-ab-drop-shadow"),
              (this.isSetup = !0),
              t.trigger("viewSetup", u);
            var c = document.body.contains(u);
            c && we.a.observe(u), t.set("inDom", c);
          }),
          (this.init = function () {
            this.updateBounds(),
              t.on("change:fullscreen", K),
              t.on("change:activeTab", V),
              t.on("change:fullscreen", V),
              t.on("change:intersectionRatio", V),
              t.on("change:visibility", W),
              t.on("instreamMode", function (e) {
                e ? de() : fe();
              }),
              V(),
              1 !== we.a.size() || t.get("visibility") || W(t, 1, 0);
            var e = t.player;
            t.change("state", ae),
              e.change("controls", q),
              t.change("streamType", ie),
              t.change("mediaType", oe),
              e.change("playlistItem", function (e, t) {
                le(e, t);
              }),
              (o = r = null),
              x && Ue && we.a.addScrollHandler(B),
              this.checkResized();
          });
        var F,
          z = 62,
          N = !0;
        function H() {
          var e = t.get("isFloating"),
            n = P.top < z,
            i = n ? P.top <= window.scrollY : P.top <= window.scrollY + z;
          !e && i ? ge(0, n) : e && !i && ge(1, n);
        }
        function B() {
          E() &&
            t.get("inDom") &&
            (clearTimeout(F),
            (F = setTimeout(H, 150)),
            N &&
              ((N = !1),
              H(),
              setTimeout(function () {
                N = !0;
              }, 50)));
        }
        function q(e, t) {
          var n = { controls: t };
          t
            ? (qe = Ce.a.controls)
              ? D()
              : ((n.loadPromise = Object(Ce.b)().then(function (t) {
                  qe = t;
                  var n = e.get("controls");
                  return n && D(), n;
                })),
                n.loadPromise.catch(function (e) {
                  l.trigger(d.tb, e);
                }))
            : l.removeControls(),
            o && r && l.trigger(d.o, n);
        }
        function D() {
          var e = new qe(document, l.element());
          l.addControls(e);
        }
        function W(e, t, n) {
          t && !n && (ae(e, e.get("state")), l.updateStyles());
        }
        function U(e) {
          A && A.mouseMove(e);
        }
        function Q(e) {
          A && !A.showing && "IFRAME" === e.target.nodeName && A.userActive();
        }
        function Y(e) {
          A &&
            A.showing &&
            ((e.relatedTarget && !u.contains(e.relatedTarget)) ||
              (!e.relatedTarget && v.Features.iframe)) &&
            A.userActive();
        }
        function X(e, t) {
          Object(Me.p)(u, /jw-stretch-\S+/, "jw-stretch-" + t);
        }
        function J(e, n) {
          Object(Me.v)(u, "jw-flag-aspect-mode", !!n);
          var i = u.querySelectorAll(".jw-aspect");
          Object(xe.d)(i, { paddingTop: n || null }),
            l.isSetup &&
              n &&
              !t.get("isFloating") &&
              (Object(xe.d)(u, $(e.get("width"))), I());
        }
        function G(n) {
          n.link
            ? (e.pause({ reason: "interaction" }),
              e.setFullscreen(!1),
              Object(Me.l)(n.link, n.linktarget, { rel: "noreferrer" }))
            : t.get("controls") && e.playToggle({ reason: "interaction" });
        }
        (this.addControls = function (n) {
          var i = this;
          (A = n),
            Object(Me.o)(u, "jw-flag-controls-hidden"),
            Object(Me.v)(u, "jw-floating-dismissible", this.dismissible),
            n.enable(e, t),
            r && (L(o, r), n.resize(o, r), y.renderCues(!0)),
            n.on("userActive userInactive", function () {
              var e = t.get("state");
              (e !== d.pb && e !== d.jb) || y.renderCues(!0);
            }),
            n.on("dismissFloating", function () {
              i.stopFloating(!0), e.pause({ reason: "interaction" });
            }),
            n.on("all", l.trigger, l),
            t.get("instream") && A.setupInstream();
        }),
          (this.removeControls = function () {
            A && (A.disable(t), (A = null)),
              Object(Me.a)(u, "jw-flag-controls-hidden"),
              Object(Me.o)(u, "jw-floating-dismissible");
          });
        var K = function (t, n) {
          if (
            (n && A && t.get("autostartMuted") && A.unmuteAutoplay(e, t),
            M.supportsDomFullscreen())
          )
            n ? M.requestFullscreen() : M.exitFullscreen(), ne(u, n);
          else if (Qe) ne(u, n);
          else {
            var i = t.get("instream"),
              o = i ? i.provider : null,
              r = t.getVideo() || o;
            r && r.setFullscreen && r.setFullscreen(n);
          }
        };
        function $(e, n, o) {
          var r = { width: e };
          if (
            (o && void 0 !== n && t.set("aspectratio", null),
            !t.get("aspectratio"))
          ) {
            var a = n;
            Object(i.r)(a) && 0 !== a && (a = Math.max(a, be)), (r.height = a);
          }
          return r;
        }
        function Z(e, n) {
          if (
            ((e && !isNaN(1 * e)) || (e = t.get("containerWidth"))) &&
            ((n && !isNaN(1 * n)) || (n = t.get("containerHeight")))
          ) {
            m && m.resize(e, n, t.get("stretching"));
            var i = t.getVideo();
            i && i.resize(e, n, t.get("stretching"));
          }
        }
        function ee(e) {
          Object(Me.v)(u, "jw-flag-ios-fullscreen", e.jwstate), te(e);
        }
        function te(e) {
          var n = t.get("fullscreen"),
            i =
              void 0 !== e.jwstate
                ? e.jwstate
                : (function () {
                    if (M.supportsDomFullscreen()) {
                      var e = M.fullscreenElement();
                      return !(!e || e !== u);
                    }
                    return t.getVideo().getFullScreen();
                  })();
          n !== i && t.set("fullscreen", i),
            R(),
            clearTimeout(j),
            (j = setTimeout(Z, 200));
        }
        function ne(e, t) {
          Object(Me.v)(e, "jw-flag-fullscreen", t),
            Object(xe.d)(document.body, { overflowY: t ? "hidden" : "" }),
            t && A && A.userActive(),
            Z(),
            R();
        }
        function ie(e, t) {
          var n = "LIVE" === t;
          Object(Me.v)(u, "jw-flag-live", n);
        }
        function oe(e, t) {
          var n = "audio" === t,
            i = e.get("provider");
          Object(Me.v)(u, "jw-flag-media-audio", n);
          var o = i && 0 === i.name.indexOf("flash"),
            r = n && !o ? p : p.nextSibling;
          m.el.parentNode.insertBefore(m.el, r);
        }
        function re(e, t) {
          if (t) {
            var n = Object(pe.a)(e, t);
            pe.a.cloneIcon &&
              n.querySelector(".jw-icon").appendChild(pe.a.cloneIcon("error")),
              w.hide(),
              u.appendChild(n.firstChild),
              Object(Me.v)(u, "jw-flag-audio-player", !!e.get("audioMode"));
          } else w.playlistItem(e, e.get("playlistItem"));
        }
        function ae(e, t, n) {
          if (l.isSetup) {
            if (n === d.lb) {
              var i = u.querySelector(".jw-error-msg");
              i && i.parentNode.removeChild(i);
            }
            Object(ke.a)(O),
              t === d.pb
                ? se(t)
                : (O = Object(ke.b)(function () {
                    return se(t);
                  }));
          }
        }
        function se(e) {
          switch (
            (t.get("controls") &&
              e !== d.ob &&
              Object(Me.i)(u, "jw-flag-controls-hidden") &&
              (Object(Me.o)(u, "jw-flag-controls-hidden"),
              Object(Me.v)(u, "jw-floating-dismissible", l.dismissible)),
            Object(Me.p)(u, /jw-state-\S+/, "jw-state-" + e),
            e)
          ) {
            case d.lb:
              l.stopFloating();
            case d.mb:
            case d.kb:
              y && y.hide();
              break;
            default:
              y &&
                (y.show(), e === d.ob && A && !A.showing && y.renderCues(!0));
          }
        }
        (this.resize = function (e, n) {
          var i = $(e, n, !0);
          void 0 !== e &&
            void 0 !== n &&
            (t.set("width", e), t.set("height", n)),
            Object(xe.d)(u, i),
            t.get("isFloating") && ye(),
            I();
        }),
          (this.resizeMedia = Z),
          (this.setPosterImage = function (e, t) {
            t.setImage(e && e.image);
          });
        var le = function (e, t) {
            s.setPosterImage(t, m),
              Ue &&
                (function (e, t) {
                  var n = e.get("mediaElement");
                  if (n) {
                    var i = Object(Me.j)(t.title || "");
                    n.setAttribute("title", i.textContent);
                  }
                })(e, t);
          },
          ce = function () {
            var e = A && A.settingsMenu;
            return !(!e || !e.visible);
          },
          ue = function () {
            var e = A && A.infoOverlay;
            return !(!e || !e.visible);
          },
          de = function () {
            Object(Me.a)(u, "jw-flag-ads"), A && A.setupInstream(), b.disable();
          },
          fe = function () {
            if (C) {
              A && A.destroyInstream(t),
                Ye !== u || Object(_e.m)() || b.enable(),
                l.setAltText(""),
                Object(Me.o)(u, ["jw-flag-ads", "jw-flag-ads-hide-controls"]),
                t.set("hideAdsControls", !1);
              var e = t.getVideo();
              e && e.setContainer(p), C.revertAlternateClickHandlers();
            }
          };
        function ge(e, n) {
          if (e < 0.5 && !Object(_e.m)()) {
            var i = t.get("state");
            i !== d.mb &&
              i !== d.lb &&
              i !== d.kb &&
              null === Ye &&
              ((Ye = u),
              t.set("isFloating", !0),
              Object(Me.a)(u, "jw-flag-floating"),
              n &&
                (Object(xe.d)(f, {
                  transform: "translateY(-".concat(z - P.top, "px)"),
                }),
                setTimeout(function () {
                  Object(xe.d)(f, {
                    transform: "translateY(0)",
                    transition:
                      "transform 150ms cubic-bezier(0, 0.25, 0.25, 1)",
                  });
                })),
              Object(xe.d)(u, {
                backgroundImage: m.el.style.backgroundImage || t.get("image"),
              }),
              ye(),
              t.get("instreamMode") || b.enable(),
              R());
          } else l.stopFloating(!1, n);
        }
        function ye() {
          var e = t.get("width"),
            n = t.get("height"),
            o = $(e);
          if (((o.maxWidth = Math.min(400, P.width)), !t.get("aspectratio"))) {
            var r = P.width,
              a = P.height / r || 0.5625;
            Object(i.r)(e) && Object(i.r)(n) && (a = n / e),
              J(t, 100 * a + "%");
          }
          Object(xe.d)(f, o);
        }
        (this.setAltText = function (e) {
          t.set("altText", e);
        }),
          (this.clickHandler = function () {
            return C;
          }),
          (this.getContainer = this.element = function () {
            return u;
          }),
          (this.getWrapper = function () {
            return f;
          }),
          (this.controlsContainer = function () {
            return A ? A.element() : null;
          }),
          (this.getSafeRegion = function () {
            var e =
                !(arguments.length > 0 && void 0 !== arguments[0]) ||
                arguments[0],
              t = { x: 0, y: 0, width: o || 0, height: r || 0 };
            return A && e && (t.height -= A.controlbarHeight()), t;
          }),
          (this.setCaptions = function (e) {
            y.clear(), y.setup(t.get("id"), e), y.resize();
          }),
          (this.setIntersection = function (e) {
            var n = Math.round(100 * e.intersectionRatio) / 100;
            t.set("intersectionRatio", n),
              x && !E() && (S = S || n >= 0.5) && ge(n);
          }),
          (this.stopFloating = function (e, n) {
            if ((e && ((x = null), we.a.removeScrollHandler(B)), Ye === u)) {
              (Ye = null), t.set("isFloating", !1);
              var i = function () {
                Object(Me.o)(u, "jw-flag-floating"),
                  J(t, t.get("aspectratio")),
                  Object(xe.d)(u, { backgroundImage: null }),
                  Object(xe.d)(f, {
                    maxWidth: null,
                    width: null,
                    height: null,
                    left: null,
                    right: null,
                    top: null,
                    bottom: null,
                    margin: null,
                    transform: null,
                    transition: null,
                    "transition-timing-function": null,
                  });
              };
              n
                ? (Object(xe.d)(f, {
                    transform: "translateY(-".concat(z - P.top, "px)"),
                    "transition-timing-function": "ease-out",
                  }),
                  setTimeout(i, 150))
                : i(),
                b.disable(),
                R();
            }
          }),
          (this.destroy = function () {
            t.destroy(),
              we.a.unobserve(u),
              we.a.remove(this),
              (this.isSetup = !1),
              this.off(),
              Object(ke.a)(k),
              clearTimeout(j),
              Ye === u && (Ye = null),
              _ && (_.destroy(), (_ = null)),
              M && (M.destroy(), (M = null)),
              A && A.disable(t),
              C &&
                (C.destroy(),
                u.removeEventListener("mousemove", U),
                u.removeEventListener("mouseout", Y),
                u.removeEventListener("mouseover", Q),
                (C = null)),
              y.destroy(),
              n && (n.destroy(), (n = null)),
              Object(xe.a)(t.get("id")),
              this.resizeListener &&
                (this.resizeListener.destroy(), delete this.resizeListener),
              x && Ue && we.a.removeScrollHandler(B);
          });
      };
      function Je(e, t, n) {
        return (Je =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (e, t, n) {
                var i = (function (e, t) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(e, t) &&
                    null !== (e = tt(e));

                  );
                  return e;
                })(e, t);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, t);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(e, t, n || e);
      }
      function Ge(e) {
        return (Ge =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (e) {
                return typeof e;
              }
            : function (e) {
                return e &&
                  "function" == typeof Symbol &&
                  e.constructor === Symbol &&
                  e !== Symbol.prototype
                  ? "symbol"
                  : typeof e;
              })(e);
      }
      function Ke(e, t) {
        if (!(e instanceof t))
          throw new TypeError("Cannot call a class as a function");
      }
      function $e(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function Ze(e, t, n) {
        return t && $e(e.prototype, t), n && $e(e, n), e;
      }
      function et(e, t) {
        return !t || ("object" !== Ge(t) && "function" != typeof t) ? ot(e) : t;
      }
      function tt(e) {
        return (tt = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function nt(e, t) {
        if ("function" != typeof t && null !== t)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (e.prototype = Object.create(t && t.prototype, {
          constructor: { value: e, writable: !0, configurable: !0 },
        })),
          t && it(e, t);
      }
      function it(e, t) {
        return (it =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      function ot(e) {
        if (void 0 === e)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return e;
      }
      var rt = /^change:(.+)$/;
      function at(e, t, n) {
        Object.keys(t).forEach(function (i) {
          i in t &&
            t[i] !== n[i] &&
            e.trigger("change:".concat(i), e, t[i], n[i]);
        });
      }
      function st(e, t) {
        e && e.off(null, null, t);
      }
      var lt = (function (e) {
          function t(e, n) {
            var o;
            return (
              Ke(this, t),
              ((o = et(this, tt(t).call(this)))._model = e),
              (o._mediaModel = null),
              Object(i.g)(e.attributes, {
                altText: "",
                fullscreen: !1,
                logoWidth: 0,
                scrubbing: !1,
              }),
              e.on(
                "all",
                function (t, i, r, a) {
                  i === e && (i = ot(ot(o))),
                    (n && !n(t, i, r, a)) || o.trigger(t, i, r, a);
                },
                ot(ot(o))
              ),
              e.on(
                "change:mediaModel",
                function (e, t) {
                  o.mediaModel = t;
                },
                ot(ot(o))
              ),
              o
            );
          }
          return (
            nt(t, e),
            Ze(t, [
              {
                key: "get",
                value: function (e) {
                  var t = this._mediaModel;
                  return t && e in t.attributes ? t.get(e) : this._model.get(e);
                },
              },
              {
                key: "set",
                value: function (e, t) {
                  return this._model.set(e, t);
                },
              },
              {
                key: "getVideo",
                value: function () {
                  return this._model.getVideo();
                },
              },
              {
                key: "destroy",
                value: function () {
                  st(this._model, this), st(this._mediaModel, this), this.off();
                },
              },
              {
                key: "mediaModel",
                set: function (e) {
                  var t = this,
                    n = this._mediaModel;
                  st(n, this),
                    (this._mediaModel = e),
                    e.on(
                      "all",
                      function (n, i, o, r) {
                        i === e && (i = t), t.trigger(n, i, o, r);
                      },
                      this
                    ),
                    n && at(this, e.attributes, n.attributes);
                },
              },
            ]),
            t
          );
        })(y.a),
        ct = (function (e) {
          function t(e) {
            var n;
            return (
              Ke(this, t),
              ((n = et(
                this,
                tt(t).call(this, e, function (e) {
                  var t = n._instreamModel;
                  if (t) {
                    var i = rt.exec(e);
                    if (i) if (i[1] in t.attributes) return !1;
                  }
                  return !0;
                })
              ))._instreamModel = null),
              (n._playerViewModel = new lt(n._model)),
              e.on(
                "change:instream",
                function (e, t) {
                  n.instreamModel = t ? t.model : null;
                },
                ot(ot(n))
              ),
              n
            );
          }
          return (
            nt(t, e),
            Ze(t, [
              {
                key: "get",
                value: function (e) {
                  var t = this._mediaModel;
                  if (t && e in t.attributes) return t.get(e);
                  var n = this._instreamModel;
                  return n && e in n.attributes ? n.get(e) : this._model.get(e);
                },
              },
              {
                key: "getVideo",
                value: function () {
                  var e = this._instreamModel;
                  return e && e.getVideo()
                    ? e.getVideo()
                    : Je(tt(t.prototype), "getVideo", this).call(this);
                },
              },
              {
                key: "destroy",
                value: function () {
                  Je(tt(t.prototype), "destroy", this).call(this),
                    st(this._instreamModel, this);
                },
              },
              {
                key: "player",
                get: function () {
                  return this._playerViewModel;
                },
              },
              {
                key: "instreamModel",
                set: function (e) {
                  var t = this,
                    n = this._instreamModel;
                  if (
                    (st(n, this),
                    this._model.off("change:mediaModel", null, this),
                    (this._instreamModel = e),
                    this.trigger("instreamMode", !!e),
                    e)
                  )
                    e.on(
                      "all",
                      function (n, i, o, r) {
                        i === e && (i = t), t.trigger(n, i, o, r);
                      },
                      this
                    ),
                      e.change(
                        "mediaModel",
                        function (e, n) {
                          t.mediaModel = n;
                        },
                        this
                      ),
                      at(this, e.attributes, this._model.attributes);
                  else if (n) {
                    this._model.change(
                      "mediaModel",
                      function (e, n) {
                        t.mediaModel = n;
                      },
                      this
                    );
                    var o = Object(i.g)(
                      {},
                      this._model.attributes,
                      n.attributes
                    );
                    at(this, this._model.attributes, o);
                  }
                },
              },
            ]),
            t
          );
        })(lt);
      var ut,
        dt,
        ft = n(64),
        gt =
          (ut = window).URL && ut.URL.createObjectURL
            ? ut.URL
            : ut.webkitURL || ut.mozURL;
      function ht(e, t) {
        var n = t.muted;
        return (
          dt ||
            (dt = new Blob(
              [
                new Uint8Array([
                  0,
                  0,
                  0,
                  28,
                  102,
                  116,
                  121,
                  112,
                  105,
                  115,
                  111,
                  109,
                  0,
                  0,
                  2,
                  0,
                  105,
                  115,
                  111,
                  109,
                  105,
                  115,
                  111,
                  50,
                  109,
                  112,
                  52,
                  49,
                  0,
                  0,
                  0,
                  8,
                  102,
                  114,
                  101,
                  101,
                  0,
                  0,
                  2,
                  239,
                  109,
                  100,
                  97,
                  116,
                  33,
                  16,
                  5,
                  32,
                  164,
                  27,
                  255,
                  192,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  55,
                  167,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  112,
                  33,
                  16,
                  5,
                  32,
                  164,
                  27,
                  255,
                  192,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  55,
                  167,
                  128,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  112,
                  0,
                  0,
                  2,
                  194,
                  109,
                  111,
                  111,
                  118,
                  0,
                  0,
                  0,
                  108,
                  109,
                  118,
                  104,
                  100,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  3,
                  232,
                  0,
                  0,
                  0,
                  47,
                  0,
                  1,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  64,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  3,
                  0,
                  0,
                  1,
                  236,
                  116,
                  114,
                  97,
                  107,
                  0,
                  0,
                  0,
                  92,
                  116,
                  107,
                  104,
                  100,
                  0,
                  0,
                  0,
                  3,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  2,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  47,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  1,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  64,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  36,
                  101,
                  100,
                  116,
                  115,
                  0,
                  0,
                  0,
                  28,
                  101,
                  108,
                  115,
                  116,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  47,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  1,
                  100,
                  109,
                  100,
                  105,
                  97,
                  0,
                  0,
                  0,
                  32,
                  109,
                  100,
                  104,
                  100,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  172,
                  68,
                  0,
                  0,
                  8,
                  0,
                  85,
                  196,
                  0,
                  0,
                  0,
                  0,
                  0,
                  45,
                  104,
                  100,
                  108,
                  114,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  115,
                  111,
                  117,
                  110,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  83,
                  111,
                  117,
                  110,
                  100,
                  72,
                  97,
                  110,
                  100,
                  108,
                  101,
                  114,
                  0,
                  0,
                  0,
                  1,
                  15,
                  109,
                  105,
                  110,
                  102,
                  0,
                  0,
                  0,
                  16,
                  115,
                  109,
                  104,
                  100,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  36,
                  100,
                  105,
                  110,
                  102,
                  0,
                  0,
                  0,
                  28,
                  100,
                  114,
                  101,
                  102,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  12,
                  117,
                  114,
                  108,
                  32,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  211,
                  115,
                  116,
                  98,
                  108,
                  0,
                  0,
                  0,
                  103,
                  115,
                  116,
                  115,
                  100,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  87,
                  109,
                  112,
                  52,
                  97,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  2,
                  0,
                  16,
                  0,
                  0,
                  0,
                  0,
                  172,
                  68,
                  0,
                  0,
                  0,
                  0,
                  0,
                  51,
                  101,
                  115,
                  100,
                  115,
                  0,
                  0,
                  0,
                  0,
                  3,
                  128,
                  128,
                  128,
                  34,
                  0,
                  2,
                  0,
                  4,
                  128,
                  128,
                  128,
                  20,
                  64,
                  21,
                  0,
                  0,
                  0,
                  0,
                  1,
                  244,
                  0,
                  0,
                  1,
                  243,
                  249,
                  5,
                  128,
                  128,
                  128,
                  2,
                  18,
                  16,
                  6,
                  128,
                  128,
                  128,
                  1,
                  2,
                  0,
                  0,
                  0,
                  24,
                  115,
                  116,
                  116,
                  115,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  2,
                  0,
                  0,
                  4,
                  0,
                  0,
                  0,
                  0,
                  28,
                  115,
                  116,
                  115,
                  99,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  2,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  28,
                  115,
                  116,
                  115,
                  122,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  2,
                  0,
                  0,
                  1,
                  115,
                  0,
                  0,
                  1,
                  116,
                  0,
                  0,
                  0,
                  20,
                  115,
                  116,
                  99,
                  111,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  44,
                  0,
                  0,
                  0,
                  98,
                  117,
                  100,
                  116,
                  97,
                  0,
                  0,
                  0,
                  90,
                  109,
                  101,
                  116,
                  97,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  33,
                  104,
                  100,
                  108,
                  114,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  109,
                  100,
                  105,
                  114,
                  97,
                  112,
                  112,
                  108,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  0,
                  45,
                  105,
                  108,
                  115,
                  116,
                  0,
                  0,
                  0,
                  37,
                  169,
                  116,
                  111,
                  111,
                  0,
                  0,
                  0,
                  29,
                  100,
                  97,
                  116,
                  97,
                  0,
                  0,
                  0,
                  1,
                  0,
                  0,
                  0,
                  0,
                  76,
                  97,
                  118,
                  102,
                  53,
                  54,
                  46,
                  52,
                  48,
                  46,
                  49,
                  48,
                  49,
                ]),
              ],
              { type: "video/mp4" }
            )),
          (e.muted = n),
          (e.src = gt.createObjectURL(dt)),
          e.play() || Object(ft.a)(e)
        );
      }
      var pt = "autoplayEnabled",
        bt = "autoplayMuted",
        mt = "autoplayDisabled",
        wt = {};
      var vt = n(65);
      function yt(e) {
        return (
          (e = e || window.event) &&
          /^(?:mouse|pointer|touch|gesture|click|key)/.test(e.type)
        );
      }
      var jt = n(24),
        kt = "tabHidden",
        Ot = "tabVisible",
        xt = function (e) {
          var t = 0;
          return function (n) {
            var i = n.position;
            i > t && e(), (t = i);
          };
        };
      function Ct(e, t) {
        t.off(d.N, e._onPlayAttempt),
          t.off(d.fb, e._triggerFirstFrame),
          t.off(d.S, e._onTime),
          e.off("change:activeTab", e._onTabVisible);
      }
      var Mt = function (e, t) {
        e.change("mediaModel", function (e, n, i) {
          e._qoeItem && i && e._qoeItem.end(i.get("mediaState")),
            (e._qoeItem = new jt.a()),
            (e._qoeItem.getFirstFrame = function () {
              var e = this.between(d.N, d.H),
                t = this.between(Ot, d.H);
              return t > 0 && t < e ? t : e;
            }),
            e._qoeItem.tick(d.db),
            e._qoeItem.start(n.get("mediaState")),
            (function (e, t) {
              e._onTabVisible && Ct(e, t);
              var n = !1;
              (e._triggerFirstFrame = function () {
                if (!n) {
                  n = !0;
                  var i = e._qoeItem;
                  i.tick(d.H);
                  var o = i.getFirstFrame();
                  if ((t.trigger(d.H, { loadTime: o }), t.mediaController)) {
                    var r = t.mediaController.mediaModel;
                    r.off("change:".concat(d.U), null, r),
                      r.change(
                        d.U,
                        function (e, n) {
                          n && t.trigger(d.U, n);
                        },
                        r
                      );
                  }
                  Ct(e, t);
                }
              }),
                (e._onTime = xt(e._triggerFirstFrame)),
                (e._onPlayAttempt = function () {
                  e._qoeItem.tick(d.N);
                }),
                (e._onTabVisible = function (t, n) {
                  n ? e._qoeItem.tick(Ot) : e._qoeItem.tick(kt);
                }),
                e.on("change:activeTab", e._onTabVisible),
                t.on(d.N, e._onPlayAttempt),
                t.once(d.fb, e._triggerFirstFrame),
                t.on(d.S, e._onTime);
            })(e, t),
            n.on("change:mediaState", function (t, n, i) {
              n !== i && (e._qoeItem.end(i), e._qoeItem.start(n));
            });
        });
      };
      function _t(e) {
        return (_t =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (e) {
                return typeof e;
              }
            : function (e) {
                return e &&
                  "function" == typeof Symbol &&
                  e.constructor === Symbol &&
                  e !== Symbol.prototype
                  ? "symbol"
                  : typeof e;
              })(e);
      }
      var St = function () {},
        Pt = function () {};
      Object(i.g)(St.prototype, {
        setup: function (e, t, n, g, p, w) {
          var y,
            j,
            k,
            O,
            x = this,
            C = this,
            M = (C._model = new R()),
            _ = !1,
            S = !1,
            P = null,
            T = b(H),
            A = b(Pt);
          (C.originalContainer = C.currentContainer = n),
            (C._events = g),
            (C.trigger = function (e, t) {
              var n = (function (e, t, n) {
                var o = n;
                switch (t) {
                  case "time":
                  case "beforePlay":
                  case "pause":
                  case "play":
                  case "ready":
                    var r = e.get("viewable");
                    void 0 !== r && (o = Object(i.g)({}, n, { viewable: r }));
                }
                return o;
              })(M, e, t);
              return h.a.trigger.call(this, e, n);
            });
          var E = new s.a(C, ["trigger"], function () {
              return !0;
            }),
            I = function (e, t) {
              C.trigger(e, t);
            };
          M.setup(e);
          var L = M.get("backgroundLoading"),
            V = new ct(M);
          (y = this._view = new Xe(t, V)).on(
            "all",
            function (e, t) {
              (t && t.doNotForward) || I(e, t);
            },
            C
          );
          var F = (this._programController = new Y(M, w));
          ue(),
            F.on("all", I, C)
              .on(
                "subtitlesTracks",
                function (e) {
                  j.setSubtitlesTracks(e.tracks);
                  var t = j.getCurrentIndex();
                  t > 0 && ae(t, e.tracks);
                },
                C
              )
              .on(
                d.F,
                function () {
                  Promise.resolve().then(re);
                },
                C
              )
              .on(d.G, C.triggerError, C),
            Mt(M, F),
            M.on(d.w, C.triggerError, C),
            M.on(
              "change:state",
              function (e, t, n) {
                X() || J.call(x, e, t, n);
              },
              this
            ),
            M.on("change:castState", function (e, t) {
              C.trigger(d.m, t);
            }),
            M.on("change:fullscreen", function (e, t) {
              C.trigger(d.y, { fullscreen: t }),
                t && e.set("playOnViewable", !1);
            }),
            M.on("change:volume", function (e, t) {
              C.trigger(d.V, { volume: t });
            }),
            M.on("change:mute", function (e) {
              C.trigger(d.M, { mute: e.getMute() });
            }),
            M.on("change:playbackRate", function (e, t) {
              C.trigger(d.ab, { playbackRate: t, position: e.get("position") });
            });
          var z = function e(t, n) {
            ("clickthrough" !== n && "interaction" !== n && "external" !== n) ||
              (M.set("playOnViewable", !1),
              M.off("change:playReason change:pauseReason", e));
          };
          function N(e, t) {
            Object(i.t)(t) || M.set("viewable", Math.round(t));
          }
          function H() {
            de &&
              (!0 !== M.get("autostart") ||
                M.get("playOnViewable") ||
                Z("autostart"),
              de.flush());
          }
          function B(e, t) {
            C.trigger("viewable", { viewable: t }), q();
          }
          function q() {
            if (
              (o.a[0] === t || 1 === M.get("viewable")) &&
              "idle" === M.get("state") &&
              !1 === M.get("autostart")
            )
              if (!w.primed() && v.OS.android) {
                var e = w.getTestElement(),
                  n = C.getMute();
                Promise.resolve()
                  .then(function () {
                    return ht(e, { muted: n });
                  })
                  .then(function () {
                    "idle" === M.get("state") && F.preloadVideo();
                  })
                  .catch(Pt);
              } else F.preloadVideo();
          }
          function D(e) {
            (C._instreamAdapter.noResume = !e), e || te({ reason: "viewable" });
          }
          function W(e) {
            e || (C.pause({ reason: "viewable" }), M.set("playOnViewable", !e));
          }
          function U(e, t) {
            var n = X();
            if (e.get("playOnViewable")) {
              if (t) {
                var i = e.get("autoPause").pauseAds,
                  o = e.get("pauseReason");
                G() === d.mb
                  ? Z("viewable")
                  : (n && !i) ||
                    "interaction" === o ||
                    K({ reason: "viewable" });
              } else
                v.OS.mobile &&
                  !n &&
                  (C.pause({ reason: "autostart" }),
                  M.set("playOnViewable", !0));
              v.OS.mobile && n && D(t);
            }
          }
          function Q(e, t) {
            var n = e.get("state"),
              i = X(),
              o = e.get("playReason");
            i
              ? e.get("autoPause").pauseAds
                ? W(t)
                : D(t)
              : n === d.pb || n === d.jb
              ? W(t)
              : n === d.mb &&
                "playlist" === o &&
                e.once("change:state", function () {
                  W(t);
                });
          }
          function X() {
            var e = C._instreamAdapter;
            return !!e && e.getState();
          }
          function G() {
            var e = X();
            return e || M.get("state");
          }
          function K(e) {
            if ((T.cancel(), (S = !1), M.get("state") === d.lb))
              return Promise.resolve();
            var n = $(e);
            return (
              M.set("playReason", n),
              X()
                ? (t.pauseAd(!1, e), Promise.resolve())
                : (M.get("state") === d.kb && (ee(!0), C.setItemIndex(0)),
                  !_ &&
                  ((_ = !0),
                  C.trigger(d.C, {
                    playReason: n,
                    startTime:
                      e && e.startTime
                        ? e.startTime
                        : M.get("playlistItem").starttime,
                  }),
                  (_ = !1),
                  yt() && !w.primed() && w.prime(),
                  "playlist" === n &&
                    M.get("autoPause").viewability &&
                    Q(M, M.get("viewable")),
                  O)
                    ? (yt() && !L && M.get("mediaElement").load(),
                      (O = !1),
                      (k = null),
                      Promise.resolve())
                    : F.playVideo(n).then(w.played))
            );
          }
          function $(e) {
            return e && e.reason ? e.reason : "unknown";
          }
          function Z(e) {
            if (G() === d.mb) {
              T = b(H);
              var t = M.get("advertising");
              (function (e, t) {
                var n = t.cancelable,
                  i = t.muted,
                  o = void 0 !== i && i,
                  r = t.allowMuted,
                  a = void 0 !== r && r,
                  s = t.timeout,
                  l = void 0 === s ? 1e4 : s,
                  c = e.getTestElement(),
                  u = o ? "muted" : "".concat(a);
                wt[u] ||
                  (wt[u] = ht(c, { muted: o })
                    .catch(function (e) {
                      if (!n.cancelled() && !1 === o && a)
                        return ht(c, { muted: (o = !0) });
                      throw e;
                    })
                    .then(function () {
                      return o ? ((wt[u] = null), bt) : pt;
                    })
                    .catch(function (e) {
                      throw (
                        (clearTimeout(d), (wt[u] = null), (e.reason = mt), e)
                      );
                    }));
                var d,
                  f = wt[u].then(function (e) {
                    if ((clearTimeout(d), n.cancelled())) {
                      var t = new Error("Autoplay test was cancelled");
                      throw ((t.reason = "cancelled"), t);
                    }
                    return e;
                  }),
                  g = new Promise(function (e, t) {
                    d = setTimeout(function () {
                      wt[u] = null;
                      var e = new Error("Autoplay test timed out");
                      (e.reason = "timeout"), t(e);
                    }, l);
                  });
                return Promise.race([f, g]);
              })(w, {
                cancelable: T,
                muted: C.getMute(),
                allowMuted: !t || t.autoplayadsmuted,
              })
                .then(function (t) {
                  return (
                    M.set("canAutoplay", t),
                    t !== bt ||
                      C.getMute() ||
                      (M.set("autostartMuted", !0),
                      ue(),
                      M.once("change:autostartMuted", function (e) {
                        e.off("change:viewable", U),
                          C.trigger(d.M, { mute: M.getMute() });
                      })),
                    C.getMute() &&
                      M.get("enableDefaultCaptions") &&
                      j.selectDefaultIndex(1),
                    K({ reason: e }).catch(function () {
                      C._instreamAdapter || M.set("autostartFailed", !0),
                        (k = null);
                    })
                  );
                })
                .catch(function (e) {
                  if (
                    (M.set("canAutoplay", mt),
                    M.set("autostart", !1),
                    !T.cancelled())
                  ) {
                    var t = Object(m.w)(e);
                    C.trigger(d.h, { reason: e.reason, code: t, error: e });
                  }
                });
            }
          }
          function ee(e) {
            if ((T.cancel(), de.empty(), X())) {
              var t = C._instreamAdapter;
              return (
                t && (t.noResume = !0),
                void (k = function () {
                  return F.stopVideo();
                })
              );
            }
            (k = null),
              !e && (S = !0),
              _ && (O = !0),
              M.set("errorEvent", void 0),
              F.stopVideo();
          }
          function te(e) {
            var t = $(e);
            M.set("pauseReason", t), M.set("playOnViewable", "viewable" === t);
          }
          function ne(e) {
            (k = null), T.cancel();
            var n = X();
            if (n && n !== d.ob) return te(e), void t.pauseAd(!0, e);
            switch (M.get("state")) {
              case d.lb:
                return;
              case d.pb:
              case d.jb:
                te(e), F.pause();
                break;
              default:
                _ && (O = !0);
            }
          }
          function ie(e, t) {
            ee(!0), C.setItemIndex(e), C.play(t);
          }
          function oe(e) {
            ie(M.get("item") + 1, e);
          }
          function re() {
            C.completeCancelled() ||
              ((k = C.completeHandler),
              C.shouldAutoAdvance()
                ? C.nextItem()
                : M.get("repeat")
                ? oe({ reason: "repeat" })
                : (v.OS.iOS && le(!1),
                  M.set("playOnViewable", !1),
                  M.set("state", d.kb),
                  C.trigger(d.cb, {})));
          }
          function ae(e, t) {
            (e = parseInt(e, 10) || 0),
              M.persistVideoSubtitleTrack(e, t),
              (F.subtitles = e),
              C.trigger(d.k, { tracks: se(), track: e });
          }
          function se() {
            return j.getCaptionsList();
          }
          function le(e) {
            Object(i.n)(e) || (e = !M.get("fullscreen")),
              M.set("fullscreen", e),
              C._instreamAdapter &&
                C._instreamAdapter._adModel &&
                C._instreamAdapter._adModel.set("fullscreen", e);
          }
          function ue() {
            (F.mute = M.getMute()), (F.volume = M.get("volume"));
          }
          M.on("change:playReason change:pauseReason", z),
            C.on(d.c, function (e) {
              return z(0, e.playReason);
            }),
            C.on(d.b, function (e) {
              return z(0, e.pauseReason);
            }),
            M.on("change:scrubbing", function (e, t) {
              t
                ? ((P = M.get("state") !== d.ob), ne())
                : P && K({ reason: "interaction" });
            }),
            M.on("change:captionsList", function (e, t) {
              C.trigger(d.l, { tracks: t, track: M.get("captionsIndex") || 0 });
            }),
            M.on("change:mediaModel", function (e, t) {
              var n = this;
              e.set("errorEvent", void 0),
                t.change(
                  "mediaState",
                  function (t, n) {
                    var i;
                    e.get("errorEvent") ||
                      e.set(d.bb, (i = n) === d.nb || i === d.qb ? d.jb : i);
                  },
                  this
                ),
                t.change(
                  "duration",
                  function (t, n) {
                    if (0 !== n) {
                      var i = e.get("minDvrWindow"),
                        o = Object(vt.b)(n, i);
                      e.setStreamType(o);
                    }
                  },
                  this
                );
              var i = e.get("item") + 1,
                o = "autoplay" === (e.get("related") || {}).oncomplete,
                r = e.get("playlist")[i];
              if ((r || o) && L) {
                t.on(
                  "change:position",
                  function e(i, a) {
                    var s = r && !r.daiSetting,
                      l = t.get("duration");
                    s && a && l > 0 && a >= l - f.b
                      ? (t.off("change:position", e, n), F.backgroundLoad(r))
                      : o && (r = M.get("nextUp"));
                  },
                  this
                );
              }
            }),
            (j = new ge(M)).on("all", I, C),
            V.on("viewSetup", function (e) {
              Object(r.b)(x, e);
            }),
            (this.playerReady = function () {
              y.once(d.hb, function () {
                try {
                  !(function () {
                    M.change("visibility", N),
                      E.off(),
                      C.trigger(d.gb, { setupTime: 0 }),
                      M.change("playlist", function (e, t) {
                        if (t.length) {
                          var n = { playlist: t },
                            o = M.get("feedData");
                          o && (n.feedData = Object(i.g)({}, o)),
                            C.trigger(d.eb, n);
                        }
                      }),
                      M.change("playlistItem", function (e, t) {
                        if (t) {
                          var n = t.title,
                            i = t.image;
                          if (
                            "mediaSession" in navigator &&
                            window.MediaMetadata &&
                            (n || i)
                          )
                            try {
                              navigator.mediaSession.metadata = new window.MediaMetadata(
                                {
                                  title: n,
                                  artist: window.location.hostname,
                                  artwork: [{ src: i || "" }],
                                }
                              );
                            } catch (e) {}
                          e.set("cues", []),
                            C.trigger(d.db, { index: M.get("item"), item: t });
                        }
                      }),
                      E.flush(),
                      E.destroy(),
                      (E = null),
                      M.change("viewable", B),
                      M.change("viewable", U),
                      M.get("autoPause").viewability
                        ? M.change("viewable", Q)
                        : M.once(
                            "change:autostartFailed change:mute",
                            function (e) {
                              e.off("change:viewable", U);
                            }
                          );
                    H(),
                      M.on("change:itemReady", function (e, t) {
                        t && de.flush();
                      });
                  })();
                } catch (e) {
                  C.triggerError(Object(m.v)(m.m, m.a, e));
                }
              }),
                y.init();
            }),
            (this.preload = q),
            (this.load = function (e, t) {
              var n,
                i = C._instreamAdapter;
              switch (
                (i && (i.noResume = !0),
                C.trigger("destroyPlugin", {}),
                ee(!0),
                T.cancel(),
                (T = b(H)),
                A.cancel(),
                yt() && w.prime(),
                _t(e))
              ) {
                case "string":
                  (M.attributes.item = 0),
                    (M.attributes.itemReady = !1),
                    (A = b(function (e) {
                      if (e)
                        return C.updatePlaylist(Object(c.a)(e.playlist), e);
                    })),
                    (n = (function (e) {
                      var t = this;
                      return new Promise(function (n, i) {
                        var o = new l.a();
                        o.on(d.eb, function (e) {
                          n(e);
                        }),
                          o.on(d.w, i, t),
                          o.load(e);
                      });
                    })(e).then(A.async));
                  break;
                case "object":
                  (M.attributes.item = 0),
                    (n = C.updatePlaylist(Object(c.a)(e), t || {}));
                  break;
                case "number":
                  n = C.setItemIndex(e);
                  break;
                default:
                  return;
              }
              n.catch(function (e) {
                C.triggerError(Object(m.u)(e, m.c));
              }),
                n.then(T.async).catch(Pt);
            }),
            (this.play = function (e) {
              return K(e).catch(Pt);
            }),
            (this.pause = ne),
            (this.seek = function (e, t) {
              var n = M.get("state");
              if (n !== d.lb) {
                F.position = e;
                var i = n === d.mb;
                M.get("scrubbing") ||
                  (!i && n !== d.kb) ||
                  (i && ((t = t || {}).startTime = e), this.play(t));
              }
            }),
            (this.stop = ee),
            (this.playlistItem = ie),
            (this.playlistNext = oe),
            (this.playlistPrev = function (e) {
              ie(M.get("item") - 1, e);
            }),
            (this.setCurrentCaptions = ae),
            (this.setCurrentQuality = function (e) {
              F.quality = e;
            }),
            (this.setFullscreen = le),
            (this.getCurrentQuality = function () {
              return F.quality;
            }),
            (this.getQualityLevels = function () {
              return F.qualities;
            }),
            (this.setCurrentAudioTrack = function (e) {
              F.audioTrack = e;
            }),
            (this.getCurrentAudioTrack = function () {
              return F.audioTrack;
            }),
            (this.getAudioTracks = function () {
              return F.audioTracks;
            }),
            (this.getCurrentCaptions = function () {
              return j.getCurrentIndex();
            }),
            (this.getCaptionsList = se),
            (this.getVisualQuality = function () {
              var e = this._model.get("mediaModel");
              return e ? e.get(d.U) : null;
            }),
            (this.getConfig = function () {
              return this._model ? this._model.getConfiguration() : void 0;
            }),
            (this.getState = G),
            (this.next = Pt),
            (this.completeHandler = re),
            (this.completeCancelled = function () {
              return (
                ((e = M.get("state")) !== d.mb && e !== d.kb && e !== d.lb) ||
                (!!S && ((S = !1), !0))
              );
              var e;
            }),
            (this.shouldAutoAdvance = function () {
              return M.get("item") !== M.get("playlist").length - 1;
            }),
            (this.nextItem = function () {
              oe({ reason: "playlist" });
            }),
            (this.setConfig = function (e) {
              !(function (e, t) {
                var n = e._model,
                  i = n.attributes;
                t.height &&
                  ((t.height = Object(a.b)(t.height)),
                  (t.width = t.width || i.width)),
                  t.width &&
                    ((t.width = Object(a.b)(t.width)),
                    t.aspectratio
                      ? ((i.width = t.width), delete t.width)
                      : (t.height = i.height)),
                  t.width &&
                    t.height &&
                    !t.aspectratio &&
                    e._view.resize(t.width, t.height),
                  Object.keys(t).forEach(function (o) {
                    var r = t[o];
                    if (void 0 !== r)
                      switch (o) {
                        case "aspectratio":
                          n.set(o, Object(a.a)(r, i.width));
                          break;
                        case "autostart":
                          !(function (e, t, n) {
                            e.setAutoStart(n),
                              "idle" === e.get("state") &&
                                !0 === n &&
                                t.play({ reason: "autostart" });
                          })(n, e, r);
                          break;
                        case "mute":
                          e.setMute(r);
                          break;
                        case "volume":
                          e.setVolume(r);
                          break;
                        case "playbackRateControls":
                        case "playbackRates":
                        case "repeat":
                        case "stretching":
                          n.set(o, r);
                      }
                  });
              })(C, e);
            }),
            (this.setItemIndex = function (e) {
              F.stopVideo();
              var t = M.get("playlist").length;
              return (
                (e = (parseInt(e, 10) || 0) % t) < 0 && (e += t),
                F.setActiveItem(e).catch(function (e) {
                  e.code >= 151 && e.code <= 162 && (e = Object(m.u)(e, m.e)),
                    x.triggerError(Object(m.v)(m.k, m.d, e));
                })
              );
            }),
            (this.detachMedia = function () {
              if (
                (_ && (O = !0),
                M.get("autoPause").viewability && Q(M, M.get("viewable")),
                !L)
              )
                return F.setAttached(!1);
              F.backgroundActiveMedia();
            }),
            (this.attachMedia = function () {
              L ? F.restoreBackgroundMedia() : F.setAttached(!0),
                "function" == typeof k && k();
            }),
            (this.routeEvents = function (e) {
              return F.routeEvents(e);
            }),
            (this.forwardEvents = function () {
              return F.forwardEvents();
            }),
            (this.playVideo = function (e) {
              return F.playVideo(e);
            }),
            (this.stopVideo = function () {
              return F.stopVideo();
            }),
            (this.castVideo = function (e, t) {
              return F.castVideo(e, t);
            }),
            (this.stopCast = function () {
              return F.stopCast();
            }),
            (this.backgroundActiveMedia = function () {
              return F.backgroundActiveMedia();
            }),
            (this.restoreBackgroundMedia = function () {
              return F.restoreBackgroundMedia();
            }),
            (this.preloadNextItem = function () {
              F.background.currentMedia && F.preloadVideo();
            }),
            (this.isBeforeComplete = function () {
              return F.beforeComplete;
            }),
            (this.setVolume = function (e) {
              M.setVolume(e), ue();
            }),
            (this.setMute = function (e) {
              M.setMute(e), ue();
            }),
            (this.setPlaybackRate = function (e) {
              M.setPlaybackRate(e);
            }),
            (this.getProvider = function () {
              return M.get("provider");
            }),
            (this.getWidth = function () {
              return M.get("containerWidth");
            }),
            (this.getHeight = function () {
              return M.get("containerHeight");
            }),
            (this.getItemQoe = function () {
              return M._qoeItem;
            }),
            (this.addButton = function (e, t, n, i, o) {
              var r = M.get("customButtons") || [],
                a = !1,
                s = { img: e, tooltip: t, callback: n, id: i, btnClass: o };
              (r = r.reduce(function (e, t) {
                return t.id === i ? ((a = !0), e.push(s)) : e.push(t), e;
              }, [])),
                a || r.unshift(s),
                M.set("customButtons", r);
            }),
            (this.removeButton = function (e) {
              var t = M.get("customButtons") || [];
              (t = t.filter(function (t) {
                return t.id !== e;
              })),
                M.set("customButtons", t);
            }),
            (this.resize = y.resize),
            (this.getSafeRegion = y.getSafeRegion),
            (this.setCaptions = y.setCaptions),
            (this.checkBeforePlay = function () {
              return _;
            }),
            (this.setControls = function (e) {
              Object(i.n)(e) || (e = !M.get("controls")),
                M.set("controls", e),
                (F.controls = e);
            }),
            (this.addCues = function (e) {
              this.setCues(M.get("cues").concat(e));
            }),
            (this.setCues = function (e) {
              M.set("cues", e);
            }),
            (this.updatePlaylist = function (e, t) {
              try {
                var n = Object(c.b)(e, M, t);
                Object(c.e)(n);
                var o = Object(i.g)({}, t);
                delete o.playlist, M.set("feedData", o), M.set("playlist", n);
              } catch (e) {
                return Promise.reject(e);
              }
              return this.setItemIndex(M.get("item"));
            }),
            (this.setPlaylistItem = function (e, t) {
              (t = Object(c.d)(M, new u.a(t), t.feedData || {})) &&
                ((M.get("playlist")[e] = t),
                e === M.get("item") &&
                  "idle" === M.get("state") &&
                  this.setItemIndex(e));
            }),
            (this.playerDestroy = function () {
              this.off(),
                this.stop(),
                Object(r.b)(this, this.originalContainer),
                y && y.destroy(),
                M && M.destroy(),
                de && de.destroy(),
                j && j.destroy(),
                F && F.destroy(),
                this.instreamDestroy();
            }),
            (this.isBeforePlay = this.checkBeforePlay),
            (this.createInstream = function () {
              return (
                this.instreamDestroy(),
                (this._instreamAdapter = new ce(this, M, y, w)),
                this._instreamAdapter
              );
            }),
            (this.instreamDestroy = function () {
              C._instreamAdapter &&
                (C._instreamAdapter.destroy(), (C._instreamAdapter = null));
            });
          var de = new s.a(
            this,
            [
              "play",
              "pause",
              "setCurrentAudioTrack",
              "setCurrentCaptions",
              "setCurrentQuality",
              "setFullscreen",
            ],
            function () {
              return !x._model.get("itemReady") || E;
            }
          );
          de.queue.push.apply(de.queue, p), y.setup();
        },
        get: function (e) {
          if (e in j.a) {
            var t = this._model.get("mediaModel");
            return t ? t.get(e) : j.a[e];
          }
          return this._model.get(e);
        },
        getContainer: function () {
          return this.currentContainer || this.originalContainer;
        },
        getMute: function () {
          return this._model.getMute();
        },
        triggerError: function (e) {
          var t = this._model;
          (e.message = t.get("localization").errors[e.key]),
            delete e.key,
            t.set("errorEvent", e),
            t.set("state", d.lb),
            t.once(
              "change:state",
              function () {
                this.set("errorEvent", void 0);
              },
              t
            ),
            this.trigger(d.w, e);
        },
      });
      t.default = St;
    },
    57: function (e, t, n) {
      "use strict";
      n.d(t, "a", function () {
        return o;
      });
      var i = n(2);
      function o(e) {
        var t = [],
          n = (e = Object(i.i)(e)).split("\r\n\r\n");
        1 === n.length && (n = e.split("\n\n"));
        for (var o = 0; o < n.length; o++)
          if ("WEBVTT" !== n[o]) {
            var a = r(n[o]);
            a.text && t.push(a);
          }
        return t;
      }
      function r(e) {
        var t = {},
          n = e.split("\r\n");
        1 === n.length && (n = e.split("\n"));
        var o = 1;
        if (
          (n[0].indexOf(" --\x3e ") > 0 && (o = 0),
          n.length > o + 1 && n[o + 1])
        ) {
          var r = n[o],
            a = r.indexOf(" --\x3e ");
          a > 0 &&
            ((t.begin = Object(i.g)(r.substr(0, a))),
            (t.end = Object(i.g)(r.substr(a + 5))),
            (t.text = n.slice(o + 1).join("\r\n")));
        }
        return t;
      }
    },
    58: function (e, t, n) {
      "use strict";
      n.d(t, "a", function () {
        return o;
      }),
        n.d(t, "b", function () {
          return r;
        });
      var i = n(5);
      function o(e) {
        var t = -1;
        return (
          e >= 1280
            ? (t = 7)
            : e >= 960
            ? (t = 6)
            : e >= 800
            ? (t = 5)
            : e >= 640
            ? (t = 4)
            : e >= 540
            ? (t = 3)
            : e >= 420
            ? (t = 2)
            : e >= 320
            ? (t = 1)
            : e >= 250 && (t = 0),
          t
        );
      }
      function r(e, t) {
        var n = "jw-breakpoint-" + t;
        Object(i.p)(e, /jw-breakpoint--?\d+/, n);
      }
    },
    59: function (e, t, n) {
      "use strict";
      n.d(t, "a", function () {
        return d;
      });
      var i,
        o = n(0),
        r = n(8),
        a = n(16),
        s = n(7),
        l = n(3),
        c = n(10),
        u = n(5),
        d = {
          back: !0,
          backgroundOpacity: 50,
          edgeStyle: null,
          fontSize: 14,
          fontOpacity: 100,
          fontScale: 0.05,
          preprocessor: o.k,
          windowOpacity: 0,
        },
        f = function (e) {
          var t,
            s,
            f,
            g,
            h,
            p,
            b,
            m,
            w,
            v = this,
            y = e.player;
          function j() {
            Object(o.o)(t.fontSize) &&
              (y.get("containerHeight")
                ? (m =
                    (d.fontScale * (t.userFontScale || 1) * t.fontSize) /
                    d.fontSize)
                : y.once("change:containerHeight", j, this));
          }
          function k() {
            var e = y.get("containerHeight");
            if (e) {
              var t;
              if (y.get("fullscreen") && r.OS.iOS) t = null;
              else {
                var n = e * m;
                t =
                  Math.round(
                    10 *
                      (function (e) {
                        var t = y.get("mediaElement");
                        if (t && t.videoHeight) {
                          var n = t.videoWidth,
                            i = t.videoHeight,
                            o = n / i,
                            a = y.get("containerHeight"),
                            s = y.get("containerWidth");
                          if (y.get("fullscreen") && r.OS.mobile) {
                            var l = window.screen;
                            l.orientation &&
                              ((a = l.availHeight), (s = l.availWidth));
                          }
                          if (s && a && n && i)
                            return (s / a > o ? a : (i * s) / n) * m;
                        }
                        return e;
                      })(n)
                  ) / 10;
              }
              y.get("renderCaptionsNatively")
                ? (function (e, t) {
                    var n = "#".concat(
                      e,
                      " .jw-video::-webkit-media-text-track-display"
                    );
                    t &&
                      ((t += "px"),
                      r.OS.iOS &&
                        Object(c.b)(n, { fontSize: "inherit" }, e, !0));
                    (w.fontSize = t), Object(c.b)(n, w, e, !0);
                  })(y.get("id"), t)
                : Object(c.d)(h, { fontSize: t });
            }
          }
          function O(e, t, n) {
            var i = Object(c.c)("#000000", n);
            "dropshadow" === e
              ? (t.textShadow = "0 2px 1px " + i)
              : "raised" === e
              ? (t.textShadow =
                  "0 0 5px " + i + ", 0 1px 5px " + i + ", 0 2px 5px " + i)
              : "depressed" === e
              ? (t.textShadow = "0 -2px 1px " + i)
              : "uniform" === e &&
                (t.textShadow =
                  "-2px 0 1px " +
                  i +
                  ",2px 0 1px " +
                  i +
                  ",0 -2px 1px " +
                  i +
                  ",0 2px 1px " +
                  i +
                  ",-1px 1px 1px " +
                  i +
                  ",1px 1px 1px " +
                  i +
                  ",1px -1px 1px " +
                  i +
                  ",1px 1px 1px " +
                  i);
          }
          ((h = document.createElement("div")).className =
            "jw-captions jw-reset"),
            (this.show = function () {
              Object(u.a)(h, "jw-captions-enabled");
            }),
            (this.hide = function () {
              Object(u.o)(h, "jw-captions-enabled");
            }),
            (this.populate = function (e) {
              y.get("renderCaptionsNatively") ||
                ((f = []),
                (s = e),
                e ? this.selectCues(e, g) : this.renderCues());
            }),
            (this.resize = function () {
              k(), this.renderCues(!0);
            }),
            (this.renderCues = function (e) {
              (e = !!e), i && i.processCues(window, f, h, e);
            }),
            (this.selectCues = function (e, t) {
              if (e && e.data && t && !y.get("renderCaptionsNatively")) {
                var n = this.getAlignmentPosition(e, t);
                !1 !== n &&
                  ((f = this.getCurrentCues(e.data, n)), this.renderCues(!0));
              }
            }),
            (this.getCurrentCues = function (e, t) {
              return Object(o.h)(e, function (e) {
                return t >= e.startTime && (!e.endTime || t <= e.endTime);
              });
            }),
            (this.getAlignmentPosition = function (e, t) {
              var n = e.source,
                i = t.metadata,
                r = t.currentTime;
              return n && i && Object(o.r)(i[n]) && (r = i[n]), r;
            }),
            (this.clear = function () {
              Object(u.g)(h);
            }),
            (this.setup = function (e, n) {
              (p = document.createElement("div")),
                (b = document.createElement("span")),
                (p.className = "jw-captions-window jw-reset"),
                (b.className = "jw-captions-text jw-reset"),
                (t = Object(o.g)({}, d, n)),
                (m = d.fontScale);
              var i = function () {
                if (!y.get("renderCaptionsNatively")) {
                  j(t.fontSize);
                  var n = t.windowColor,
                    i = t.windowOpacity,
                    o = t.edgeStyle;
                  w = {};
                  var a = {};
                  !(function (e, t) {
                    var n = t.color,
                      i = t.fontOpacity;
                    (n || i !== d.fontOpacity) &&
                      (e.color = Object(c.c)(n || "#ffffff", i));
                    if (t.back) {
                      var o = t.backgroundColor,
                        r = t.backgroundOpacity;
                      (o === d.backgroundColor && r === d.backgroundOpacity) ||
                        (e.backgroundColor = Object(c.c)(o, r));
                    } else e.background = "transparent";
                    t.fontFamily && (e.fontFamily = t.fontFamily);
                    t.fontStyle && (e.fontStyle = t.fontStyle);
                    t.fontWeight && (e.fontWeight = t.fontWeight);
                    t.textDecoration && (e.textDecoration = t.textDecoration);
                  })(a, t),
                    (n || i !== d.windowOpacity) &&
                      (w.backgroundColor = Object(c.c)(n || "#000000", i)),
                    O(o, a, t.fontOpacity),
                    t.back || null !== o || O("uniform", a),
                    Object(c.d)(p, w),
                    Object(c.d)(b, a),
                    (function (e, t) {
                      k(),
                        (function (e, t) {
                          r.Browser.safari &&
                            Object(c.b)(
                              "#" +
                                e +
                                " .jw-video::-webkit-media-text-track-display-backdrop",
                              { backgroundColor: t.backgroundColor },
                              e,
                              !0
                            );
                          Object(c.b)(
                            "#" +
                              e +
                              " .jw-video::-webkit-media-text-track-display",
                            w,
                            e,
                            !0
                          ),
                            Object(c.b)("#" + e + " .jw-video::cue", t, e, !0);
                        })(e, t),
                        (function (e, t) {
                          Object(c.b)(
                            "#" + e + " .jw-text-track-display",
                            w,
                            e
                          ),
                            Object(c.b)("#" + e + " .jw-text-track-cue", t, e);
                        })(e, t);
                    })(e, a);
                }
              };
              i(),
                p.appendChild(b),
                h.appendChild(p),
                y.change(
                  "captionsTrack",
                  function (e, t) {
                    this.populate(t);
                  },
                  this
                ),
                y.set("captions", t),
                y.on("change:captions", function (e, n) {
                  (t = n), i();
                });
            }),
            (this.element = function () {
              return h;
            }),
            (this.destroy = function () {
              y.off(null, null, this), this.off();
            });
          var x = function (e) {
            (g = e), v.selectCues(s, g);
          };
          y.on(
            "change:playlistItem",
            function () {
              (g = null), (f = []);
            },
            this
          ),
            y.on(
              l.Q,
              function (e) {
                (f = []), x(e);
              },
              this
            ),
            y.on(l.S, x, this),
            y.on(
              "subtitlesTrackData",
              function () {
                this.selectCues(s, g);
              },
              this
            ),
            y.on(
              "change:captionsList",
              function e(t, o) {
                var r = this;
                1 !== o.length &&
                  (t.get("renderCaptionsNatively") ||
                    i ||
                    (n
                      .e(8)
                      .then(
                        function (e) {
                          i = n(68).default;
                        }.bind(null, n)
                      )
                      .catch(Object(a.c)(301121))
                      .catch(function (e) {
                        r.trigger(l.tb, e);
                      }),
                    t.off("change:captionsList", e, this)));
              },
              this
            );
        };
      Object(o.g)(f.prototype, s.a), (t.b = f);
    },
    60: function (e, t, n) {
      "use strict";
      e.exports = function (e) {
        var t = [];
        return (
          (t.toString = function () {
            return this.map(function (t) {
              var n = (function (e, t) {
                var n = e[1] || "",
                  i = e[3];
                if (!i) return n;
                if (t && "function" == typeof btoa) {
                  var o =
                      ((a = i),
                      "/*# sourceMappingURL=data:application/json;charset=utf-8;base64," +
                        btoa(unescape(encodeURIComponent(JSON.stringify(a)))) +
                        " */"),
                    r = i.sources.map(function (e) {
                      return "/*# sourceURL=" + i.sourceRoot + e + " */";
                    });
                  return [n].concat(r).concat([o]).join("\n");
                }
                var a;
                return [n].join("\n");
              })(t, e);
              return t[2] ? "@media " + t[2] + "{" + n + "}" : n;
            }).join("");
          }),
          (t.i = function (e, n) {
            "string" == typeof e && (e = [[null, e, ""]]);
            for (var i = {}, o = 0; o < this.length; o++) {
              var r = this[o][0];
              null != r && (i[r] = !0);
            }
            for (o = 0; o < e.length; o++) {
              var a = e[o];
              (null != a[0] && i[a[0]]) ||
                (n && !a[2]
                  ? (a[2] = n)
                  : n && (a[2] = "(" + a[2] + ") and (" + n + ")"),
                t.push(a));
            }
          }),
          t
        );
      };
    },
    61: function (e, t) {
      var n,
        i,
        o = {},
        r = {},
        a =
          ((n = function () {
            return document.head || document.getElementsByTagName("head")[0];
          }),
          function () {
            return void 0 === i && (i = n.apply(this, arguments)), i;
          });
      function s(e) {
        var t = document.createElement("style");
        return (
          (t.type = "text/css"),
          t.setAttribute("data-jwplayer-id", e),
          (function (e) {
            a().appendChild(e);
          })(t),
          t
        );
      }
      function l(e, t) {
        var n,
          i,
          o,
          a = r[e];
        a || (a = r[e] = { element: s(e), counter: 0 });
        var l = a.counter++;
        return (
          (n = a.element),
          (o = function () {
            d(n, l, "");
          }),
          (i = function (e) {
            d(n, l, e);
          })(t.css),
          function (e) {
            if (e) {
              if (e.css === t.css && e.media === t.media) return;
              i((t = e).css);
            } else o();
          }
        );
      }
      e.exports = {
        style: function (e, t) {
          !(function (e, t) {
            for (var n = 0; n < t.length; n++) {
              var i = t[n],
                r = (o[e] || {})[i.id];
              if (r) {
                for (var a = 0; a < r.parts.length; a++) r.parts[a](i.parts[a]);
                for (; a < i.parts.length; a++) r.parts.push(l(e, i.parts[a]));
              } else {
                var s = [];
                for (a = 0; a < i.parts.length; a++) s.push(l(e, i.parts[a]));
                (o[e] = o[e] || {}), (o[e][i.id] = { id: i.id, parts: s });
              }
            }
          })(
            t,
            (function (e) {
              for (var t = [], n = {}, i = 0; i < e.length; i++) {
                var o = e[i],
                  r = o[0],
                  a = o[1],
                  s = o[2],
                  l = { css: a, media: s };
                n[r]
                  ? n[r].parts.push(l)
                  : t.push((n[r] = { id: r, parts: [l] }));
              }
              return t;
            })(e)
          );
        },
        clear: function (e, t) {
          var n = o[e];
          if (!n) return;
          if (t) {
            var i = n[t];
            if (i) for (var r = 0; r < i.parts.length; r += 1) i.parts[r]();
            return;
          }
          for (var a = Object.keys(n), s = 0; s < a.length; s += 1)
            for (var l = n[a[s]], c = 0; c < l.parts.length; c += 1)
              l.parts[c]();
          delete o[e];
        },
      };
      var c,
        u =
          ((c = []),
          function (e, t) {
            return (c[e] = t), c.filter(Boolean).join("\n");
          });
      function d(e, t, n) {
        if (e.styleSheet) e.styleSheet.cssText = u(t, n);
        else {
          var i = document.createTextNode(n),
            o = e.childNodes[t];
          o ? e.replaceChild(i, o) : e.appendChild(i);
        }
      }
    },
    63: function (e, t, n) {
      "use strict";
      function i(e, t) {
        var n = e.kind || "cc";
        return e.default || e.defaulttrack
          ? "default"
          : e._id || e.file || n + t;
      }
      function o(e, t) {
        var n = e.label || e.name || e.language;
        return (
          n || ((n = "Unknown CC"), (t += 1) > 1 && (n += " [" + t + "]")),
          { label: n, unknownCount: t }
        );
      }
      n.d(t, "a", function () {
        return i;
      }),
        n.d(t, "b", function () {
          return o;
        });
    },
    64: function (e, t, n) {
      "use strict";
      function i(e) {
        return new Promise(function (t, n) {
          if (e.paused) return n(o("NotAllowedError", 0, "play() failed."));
          var i = function () {
              e.removeEventListener("play", r),
                e.removeEventListener("playing", a),
                e.removeEventListener("pause", a),
                e.removeEventListener("abort", a),
                e.removeEventListener("error", a);
            },
            r = function () {
              e.addEventListener("playing", a),
                e.addEventListener("abort", a),
                e.addEventListener("error", a),
                e.addEventListener("pause", a);
            },
            a = function (e) {
              if ((i(), "playing" === e.type)) t();
              else {
                var r = 'The play() request was interrupted by a "'.concat(
                  e.type,
                  '" event.'
                );
                "error" === e.type
                  ? n(o("NotSupportedError", 9, r))
                  : n(o("AbortError", 20, r));
              }
            };
          e.addEventListener("play", r);
        });
      }
      function o(e, t, n) {
        var i = new Error(n);
        return (i.name = e), (i.code = t), i;
      }
      n.d(t, "a", function () {
        return i;
      });
    },
    65: function (e, t, n) {
      "use strict";
      function i(e, t) {
        return e !== 1 / 0 && Math.abs(e) >= Math.max(r(t), 0);
      }
      function o(e, t) {
        var n = "VOD";
        return (
          e === 1 / 0
            ? (n = "LIVE")
            : e < 0 && (n = i(e, r(t)) ? "DVR" : "LIVE"),
          n
        );
      }
      function r(e) {
        return void 0 === e ? 120 : Math.max(e, 0);
      }
      n.d(t, "a", function () {
        return i;
      }),
        n.d(t, "b", function () {
          return o;
        });
    },
    66: function (e, t, n) {
      "use strict";
      var i = n(67),
        o = n(16),
        r = n(22),
        a = n(4),
        s = n(57),
        l = n(2),
        c = n(1);
      function u(e) {
        throw new c.n(null, e);
      }
      function d(e, t, i) {
        e.xhr = Object(r.a)(
          e.file,
          function (r) {
            !(function (e, t, i, r) {
              var d,
                f,
                h = e.responseXML ? e.responseXML.firstChild : null;
              if (h)
                for (
                  "xml" === Object(a.b)(h) && (h = h.nextSibling);
                  h.nodeType === h.COMMENT_NODE;

                )
                  h = h.nextSibling;
              try {
                if (h && "tt" === Object(a.b)(h))
                  (d = (function (e) {
                    e || u(306007);
                    var t = [],
                      n = e.getElementsByTagName("p"),
                      i = 30,
                      o = e.getElementsByTagName("tt");
                    if (o && o[0]) {
                      var r = parseFloat(o[0].getAttribute("ttp:frameRate"));
                      isNaN(r) || (i = r);
                    }
                    n || u(306005),
                      n.length ||
                        (n = e.getElementsByTagName("tt:p")).length ||
                        (n = e.getElementsByTagName("tts:p"));
                    for (var a = 0; a < n.length; a++) {
                      for (
                        var s = n[a], c = s.getElementsByTagName("br"), d = 0;
                        d < c.length;
                        d++
                      ) {
                        var f = c[d];
                        f.parentNode.replaceChild(e.createTextNode("\r\n"), f);
                      }
                      var g = s.innerHTML || s.textContent || s.text || "",
                        h = Object(l.i)(g)
                          .replace(/>\s+</g, "><")
                          .replace(/(<\/?)tts?:/g, "$1")
                          .replace(/<br.*?\/>/g, "\r\n");
                      if (h) {
                        var p = s.getAttribute("begin"),
                          b = s.getAttribute("dur"),
                          m = s.getAttribute("end"),
                          w = { begin: Object(l.g)(p, i), text: h };
                        m
                          ? (w.end = Object(l.g)(m, i))
                          : b && (w.end = w.begin + Object(l.g)(b, i)),
                          t.push(w);
                      }
                    }
                    return t.length || u(306005), t;
                  })(e.responseXML)),
                    (f = g(d)),
                    delete t.xhr,
                    i(f);
                else {
                  var p = e.responseText;
                  p.indexOf("WEBVTT") >= 0
                    ? n
                        .e(10)
                        .then(
                          function (e) {
                            return n(97).default;
                          }.bind(null, n)
                        )
                        .catch(Object(o.c)(301131))
                        .then(function (e) {
                          var n = new e(window);
                          (f = []),
                            (n.oncue = function (e) {
                              f.push(e);
                            }),
                            (n.onflush = function () {
                              delete t.xhr, i(f);
                            }),
                            n.parse(p);
                        })
                        .catch(function (e) {
                          delete t.xhr, r(Object(c.v)(null, c.b, e));
                        })
                    : ((d = Object(s.a)(p)), (f = g(d)), delete t.xhr, i(f));
                }
              } catch (e) {
                delete t.xhr, r(Object(c.v)(null, c.b, e));
              }
            })(r, e, t, i);
          },
          function (e, t, n, o) {
            i(Object(c.u)(o, c.b));
          }
        );
      }
      function f(e) {
        e &&
          e.forEach(function (e) {
            var t = e.xhr;
            t &&
              ((t.onload = null),
              (t.onreadystatechange = null),
              (t.onerror = null),
              "abort" in t && t.abort()),
              delete e.xhr;
          });
      }
      function g(e) {
        return e.map(function (e) {
          return new i.a(e.begin, e.end, e.text);
        });
      }
      n.d(t, "c", function () {
        return d;
      }),
        n.d(t, "a", function () {
          return f;
        }),
        n.d(t, "b", function () {
          return g;
        });
    },
    67: function (e, t, n) {
      "use strict";
      var i = window.VTTCue;
      function o(e) {
        if ("string" != typeof e) return !1;
        return (
          !!{ start: !0, middle: !0, end: !0, left: !0, right: !0 }[
            e.toLowerCase()
          ] && e.toLowerCase()
        );
      }
      if (!i) {
        (i = function (e, t, n) {
          var i = this;
          i.hasBeenReset = !1;
          var r = "",
            a = !1,
            s = e,
            l = t,
            c = n,
            u = null,
            d = "",
            f = !0,
            g = "auto",
            h = "start",
            p = "auto",
            b = 100,
            m = "middle";
          Object.defineProperty(i, "id", {
            enumerable: !0,
            get: function () {
              return r;
            },
            set: function (e) {
              r = "" + e;
            },
          }),
            Object.defineProperty(i, "pauseOnExit", {
              enumerable: !0,
              get: function () {
                return a;
              },
              set: function (e) {
                a = !!e;
              },
            }),
            Object.defineProperty(i, "startTime", {
              enumerable: !0,
              get: function () {
                return s;
              },
              set: function (e) {
                if ("number" != typeof e)
                  throw new TypeError("Start time must be set to a number.");
                (s = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "endTime", {
              enumerable: !0,
              get: function () {
                return l;
              },
              set: function (e) {
                if ("number" != typeof e)
                  throw new TypeError("End time must be set to a number.");
                (l = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "text", {
              enumerable: !0,
              get: function () {
                return c;
              },
              set: function (e) {
                (c = "" + e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "region", {
              enumerable: !0,
              get: function () {
                return u;
              },
              set: function (e) {
                (u = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "vertical", {
              enumerable: !0,
              get: function () {
                return d;
              },
              set: function (e) {
                var t = (function (e) {
                  return (
                    "string" == typeof e &&
                    !!{ "": !0, lr: !0, rl: !0 }[e.toLowerCase()] &&
                    e.toLowerCase()
                  );
                })(e);
                if (!1 === t)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (d = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "snapToLines", {
              enumerable: !0,
              get: function () {
                return f;
              },
              set: function (e) {
                (f = !!e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "line", {
              enumerable: !0,
              get: function () {
                return g;
              },
              set: function (e) {
                if ("number" != typeof e && "auto" !== e)
                  throw new SyntaxError(
                    "An invalid number or illegal string was specified."
                  );
                (g = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "lineAlign", {
              enumerable: !0,
              get: function () {
                return h;
              },
              set: function (e) {
                var t = o(e);
                if (!t)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (h = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "position", {
              enumerable: !0,
              get: function () {
                return p;
              },
              set: function (e) {
                if (e < 0 || e > 100)
                  throw new Error("Position must be between 0 and 100.");
                (p = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "size", {
              enumerable: !0,
              get: function () {
                return b;
              },
              set: function (e) {
                if (e < 0 || e > 100)
                  throw new Error("Size must be between 0 and 100.");
                (b = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "align", {
              enumerable: !0,
              get: function () {
                return m;
              },
              set: function (e) {
                var t = o(e);
                if (!t)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (m = t), (this.hasBeenReset = !0);
              },
            }),
            (i.displayState = void 0);
        }).prototype.getCueAsHTML = function () {
          return window.WebVTT.convertCueToDOMTree(window, this.text);
        };
      }
      t.a = i;
    },
    69: function (e, t, n) {
      var i = n(70);
      "string" == typeof i && (i = [["all-players", i, ""]]),
        n(61).style(i, "all-players"),
        i.locals && (e.exports = i.locals);
    },
    70: function (e, t, n) {
      (e.exports = n(60)(!1)).push([
        e.i,
        '.jw-reset{text-align:left;direction:ltr}.jw-reset-text,.jw-reset{color:inherit;background-color:transparent;padding:0;margin:0;float:none;font-family:Arial,Helvetica,sans-serif;font-size:1em;line-height:1em;list-style:none;text-transform:none;vertical-align:baseline;border:0;font-variant:inherit;font-stretch:inherit;-webkit-tap-highlight-color:rgba(255,255,255,0)}body .jw-error,body .jwplayer.jw-state-error{height:100%;width:100%}.jw-title{position:absolute;top:0}.jw-background-color{background:rgba(0,0,0,0.4)}.jw-text{color:rgba(255,255,255,0.8)}.jw-knob{color:rgba(255,255,255,0.8);background-color:#fff}.jw-button-color{color:rgba(255,255,255,0.8)}:not(.jw-flag-touch) .jw-button-color:not(.jw-logo-button):focus,:not(.jw-flag-touch) .jw-button-color:not(.jw-logo-button):hover{color:#fff}.jw-toggle{color:#fff}.jw-toggle.jw-off{color:rgba(255,255,255,0.8)}.jw-toggle.jw-off:focus{color:#fff}.jw-toggle:focus{outline:none}:not(.jw-flag-touch) .jw-toggle.jw-off:hover{color:#fff}.jw-rail{background:rgba(255,255,255,0.3)}.jw-buffer{background:rgba(255,255,255,0.3)}.jw-progress{background:#f2f2f2}.jw-time-tip,.jw-volume-tip{border:0}.jw-slider-volume.jw-volume-tip.jw-background-color.jw-slider-vertical{background:none}.jw-skip{padding:.5em;outline:none}.jw-skip .jw-skiptext,.jw-skip .jw-skip-icon{color:rgba(255,255,255,0.8)}.jw-skip.jw-skippable:hover .jw-skip-icon,.jw-skip.jw-skippable:focus .jw-skip-icon{color:#fff}.jw-icon-cast google-cast-launcher{--connected-color:#fff;--disconnected-color:rgba(255,255,255,0.8)}.jw-icon-cast google-cast-launcher:focus{outline:none}.jw-icon-cast google-cast-launcher.jw-off{--connected-color:rgba(255,255,255,0.8)}.jw-icon-cast:focus google-cast-launcher{--connected-color:#fff;--disconnected-color:#fff}.jw-icon-cast:hover google-cast-launcher{--connected-color:#fff;--disconnected-color:#fff}.jw-nextup-container{bottom:2.5em;padding:5px .5em}.jw-nextup{border-radius:0}.jw-color-active{color:#fff;stroke:#fff;border-color:#fff}:not(.jw-flag-touch) .jw-color-active-hover:hover,:not(.jw-flag-touch) .jw-color-active-hover:focus{color:#fff;stroke:#fff;border-color:#fff}.jw-color-inactive{color:rgba(255,255,255,0.8);stroke:rgba(255,255,255,0.8);border-color:rgba(255,255,255,0.8)}:not(.jw-flag-touch) .jw-color-inactive-hover:hover{color:rgba(255,255,255,0.8);stroke:rgba(255,255,255,0.8);border-color:rgba(255,255,255,0.8)}.jw-option{color:rgba(255,255,255,0.8)}.jw-option.jw-active-option{color:#fff;background-color:rgba(255,255,255,0.1)}:not(.jw-flag-touch) .jw-option:hover{color:#fff}.jwplayer{width:100%;font-size:16px;position:relative;display:block;min-height:0;overflow:hidden;box-sizing:border-box;font-family:Arial,Helvetica,sans-serif;-webkit-touch-callout:none;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none;outline:none}.jwplayer *{box-sizing:inherit}.jwplayer.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jwplayer.jw-flag-aspect-mode{height:auto !important}.jwplayer.jw-flag-aspect-mode .jw-aspect{display:block}.jwplayer .jw-aspect{display:none}.jwplayer .jw-swf{outline:none}.jw-media,.jw-preview{position:absolute;width:100%;height:100%;top:0;left:0;bottom:0;right:0}.jw-media{overflow:hidden;cursor:pointer}.jw-plugin{position:absolute;bottom:66px}.jw-breakpoint-7 .jw-plugin{bottom:132px}.jw-plugin .jw-banner{max-width:100%;opacity:0;cursor:pointer;position:absolute;margin:auto auto 0;left:0;right:0;bottom:0;display:block}.jw-preview,.jw-captions,.jw-title{pointer-events:none}.jw-media,.jw-logo{pointer-events:all}.jw-wrapper{background-color:#000;position:absolute;top:0;left:0;right:0;bottom:0}.jw-hidden-accessibility{border:0;clip:rect(0 0 0 0);height:1px;margin:-1px;overflow:hidden;padding:0;position:absolute;width:1px}.jw-contract-trigger::before{content:"";overflow:hidden;width:200%;height:200%;display:block;position:absolute;top:0;left:0}.jwplayer .jw-media video{position:absolute;top:0;right:0;bottom:0;left:0;width:100%;height:100%;margin:auto;background:transparent}.jwplayer .jw-media video::-webkit-media-controls-start-playback-button{display:none}.jwplayer.jw-stretch-uniform .jw-media video{object-fit:contain}.jwplayer.jw-stretch-none .jw-media video{object-fit:none}.jwplayer.jw-stretch-fill .jw-media video{object-fit:cover}.jwplayer.jw-stretch-exactfit .jw-media video{object-fit:fill}.jw-preview{position:absolute;display:none;opacity:1;visibility:visible;width:100%;height:100%;background:#000 no-repeat 50% 50%}.jwplayer .jw-preview,.jw-error .jw-preview{background-size:contain}.jw-stretch-none .jw-preview{background-size:auto auto}.jw-stretch-fill .jw-preview{background-size:cover}.jw-stretch-exactfit .jw-preview{background-size:100% 100%}.jw-title{display:none;padding-top:20px;width:100%;z-index:1}.jw-title-primary,.jw-title-secondary{color:#fff;padding-left:20px;padding-right:20px;padding-bottom:.5em;overflow:hidden;text-overflow:ellipsis;direction:unset;white-space:nowrap;width:100%}.jw-title-primary{font-size:1.625em}.jw-breakpoint-2 .jw-title-primary,.jw-breakpoint-3 .jw-title-primary{font-size:1.5em}.jw-flag-small-player .jw-title-primary{font-size:1.25em}.jw-flag-small-player .jw-title-secondary,.jw-title-secondary:empty{display:none}.jw-captions{position:absolute;width:100%;height:100%;text-align:center;display:none;letter-spacing:normal;word-spacing:normal;text-transform:none;text-indent:0;text-decoration:none;pointer-events:none;overflow:hidden;top:0}.jw-captions.jw-captions-enabled{display:block}.jw-captions-window{display:none;padding:.25em;border-radius:.25em}.jw-captions-window.jw-captions-window-active{display:inline-block}.jw-captions-text{display:inline-block;color:#fff;background-color:#000;word-wrap:normal;word-break:normal;white-space:pre-line;font-style:normal;font-weight:normal;text-align:center;text-decoration:none}.jw-text-track-display{font-size:inherit;line-height:1.5}.jw-text-track-cue{background-color:rgba(0,0,0,0.5);color:#fff;padding:.1em .3em}.jwplayer video::-webkit-media-controls{display:none;justify-content:flex-start}.jwplayer video::-webkit-media-text-track-display{min-width:-webkit-min-content}.jwplayer video::cue{background-color:rgba(0,0,0,0.5)}.jwplayer video::-webkit-media-controls-panel-container{display:none}.jwplayer:not(.jw-flag-controls-hidden):not(.jw-state-playing) .jw-captions,.jwplayer.jw-flag-media-audio.jw-state-playing .jw-captions,.jwplayer.jw-state-playing:not(.jw-flag-user-inactive):not(.jw-flag-controls-hidden) .jw-captions{max-height:calc(100% - 60px)}.jwplayer:not(.jw-flag-controls-hidden):not(.jw-state-playing):not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container,.jwplayer.jw-flag-media-audio.jw-state-playing:not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container,.jwplayer.jw-state-playing:not(.jw-flag-user-inactive):not(.jw-flag-controls-hidden):not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container{max-height:calc(100% - 60px)}.jw-logo{position:absolute;margin:20px;cursor:pointer;pointer-events:all;background-repeat:no-repeat;background-size:contain;top:auto;right:auto;left:auto;bottom:auto;outline:none}.jw-logo.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jw-flag-audio-player .jw-logo{display:none}.jw-logo-top-right{top:0;right:0}.jw-logo-top-left{top:0;left:0}.jw-logo-bottom-left{left:0}.jw-logo-bottom-right{right:0}.jw-logo-bottom-left,.jw-logo-bottom-right{bottom:44px;transition:bottom 150ms cubic-bezier(0, .25, .25, 1)}.jw-state-idle .jw-logo{z-index:1}.jw-state-setup .jw-wrapper{background-color:inherit}.jw-state-setup .jw-logo,.jw-state-setup .jw-controls,.jw-state-setup .jw-controls-backdrop{visibility:hidden}span.jw-break{display:block}body .jw-error,body .jwplayer.jw-state-error{background-color:#333;color:#fff;font-size:16px;display:table;opacity:1;position:relative}body .jw-error .jw-display,body .jwplayer.jw-state-error .jw-display{display:none}body .jw-error .jw-media,body .jwplayer.jw-state-error .jw-media{cursor:default}body .jw-error .jw-preview,body .jwplayer.jw-state-error .jw-preview{background-color:#333}body .jw-error .jw-error-msg,body .jwplayer.jw-state-error .jw-error-msg{background-color:#000;border-radius:2px;display:flex;flex-direction:row;align-items:stretch;padding:20px}body .jw-error .jw-error-msg .jw-icon,body .jwplayer.jw-state-error .jw-error-msg .jw-icon{height:30px;width:30px;margin-right:20px;flex:0 0 auto;align-self:center}body .jw-error .jw-error-msg .jw-icon:empty,body .jwplayer.jw-state-error .jw-error-msg .jw-icon:empty{display:none}body .jw-error .jw-error-msg .jw-info-container,body .jwplayer.jw-state-error .jw-error-msg .jw-info-container{margin:0;padding:0}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg{flex-direction:column}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-error-text,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-error-text,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-error-text,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-error-text{text-align:center}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-icon,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-icon,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-icon,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-icon{flex:.5 0 auto;margin-right:0;margin-bottom:20px}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg .jw-break,.jwplayer.jw-state-error.jw-flag-small-player .jw-error-msg .jw-break,.jwplayer.jw-state-error.jw-breakpoint-2 .jw-error-msg .jw-break{display:inline}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg .jw-break:before,.jwplayer.jw-state-error.jw-flag-small-player .jw-error-msg .jw-break:before,.jwplayer.jw-state-error.jw-breakpoint-2 .jw-error-msg .jw-break:before{content:" "}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg{height:100%;width:100%;top:0;position:absolute;left:0;background:#000;-webkit-transform:none;transform:none;padding:4px 16px;z-index:1}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg.jw-info-overlay{max-width:none;max-height:none}body .jwplayer.jw-state-error .jw-title,.jw-state-idle .jw-title,.jwplayer.jw-state-complete:not(.jw-flag-casting):not(.jw-flag-audio-player):not(.jw-flag-overlay-open-related) .jw-title{display:block}body .jwplayer.jw-state-error .jw-preview,.jw-state-idle .jw-preview,.jwplayer.jw-state-complete:not(.jw-flag-casting):not(.jw-flag-audio-player):not(.jw-flag-overlay-open-related) .jw-preview{display:block}.jw-state-idle .jw-captions,.jwplayer.jw-state-complete .jw-captions,body .jwplayer.jw-state-error .jw-captions{display:none}.jw-state-idle video::-webkit-media-text-track-container,.jwplayer.jw-state-complete video::-webkit-media-text-track-container,body .jwplayer.jw-state-error video::-webkit-media-text-track-container{display:none}.jwplayer.jw-flag-fullscreen{width:100% !important;height:100% !important;top:0;right:0;bottom:0;left:0;z-index:1000;margin:0;position:fixed}body .jwplayer.jw-flag-flash-blocked .jw-title{display:block}.jwplayer.jw-flag-controls-hidden .jw-media{cursor:default}.jw-flag-audio-player:not(.jw-flag-flash-blocked) .jw-media{visibility:hidden}.jw-flag-audio-player .jw-title{background:none}.jw-flag-audio-player object{min-height:45px}.jw-flag-floating{background-size:cover;background-color:#000}.jw-flag-floating .jw-wrapper{position:fixed;z-index:2147483647;-webkit-animation:jw-float-to-bottom 150ms cubic-bezier(0, .25, .25, 1) forwards 1;animation:jw-float-to-bottom 150ms cubic-bezier(0, .25, .25, 1) forwards 1;top:auto;bottom:1rem;left:auto;right:1rem;max-width:400px;max-height:400px;margin:0 auto}@media screen and (max-width:480px){.jw-flag-floating .jw-wrapper{width:100%;left:0;right:0}}.jw-flag-floating .jw-wrapper .jw-media{touch-action:none}@media screen and (max-device-width:480px) and (orientation:portrait){.jw-flag-touch.jw-flag-floating .jw-wrapper{-webkit-animation:none;animation:none;top:62px;bottom:auto;left:0;right:0;max-width:none;max-height:none}}.jw-flag-floating .jw-float-icon{pointer-events:all;cursor:pointer;display:none}.jw-flag-floating .jw-float-icon .jw-svg-icon{-webkit-filter:drop-shadow(0 0 1px #000);filter:drop-shadow(0 0 1px #000)}.jw-flag-floating.jw-floating-dismissible .jw-dismiss-icon{display:none}.jw-flag-floating.jw-floating-dismissible.jw-flag-ads .jw-float-icon{display:flex}.jw-flag-floating.jw-floating-dismissible.jw-state-paused .jw-logo,.jw-flag-floating.jw-floating-dismissible:not(.jw-flag-user-inactive) .jw-logo{display:none}.jw-flag-floating.jw-floating-dismissible.jw-state-paused .jw-float-icon,.jw-flag-floating.jw-floating-dismissible:not(.jw-flag-user-inactive) .jw-float-icon{display:flex}.jw-float-icon{display:none;position:absolute;top:3px;right:5px;align-items:center;justify-content:center}@-webkit-keyframes jw-float-to-bottom{from{-webkit-transform:translateY(100%);transform:translateY(100%)}to{-webkit-transform:translateY(0);transform:translateY(0)}}@keyframes jw-float-to-bottom{from{-webkit-transform:translateY(100%);transform:translateY(100%)}to{-webkit-transform:translateY(0);transform:translateY(0)}}.jw-flag-top{margin-top:2em;overflow:visible}.jw-top{height:2em;line-height:2;pointer-events:none;text-align:center;opacity:.8;position:absolute;top:-2em;width:100%}.jw-top .jw-icon{cursor:pointer;pointer-events:all;height:auto;width:auto}.jw-top .jw-text{color:#555}',
        "",
      ]);
    },
  },
]);
