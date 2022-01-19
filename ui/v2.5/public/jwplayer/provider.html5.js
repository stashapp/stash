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
  [9],
  {
    31: function (e, t, i) {
      "use strict";
      i.r(t);
      var r = i(0);
      var n = i(8),
        a = i(52),
        s = i(3),
        c = i(43),
        u = {
          canplay: function () {
            this.trigger(s.E);
          },
          play: function () {
            (this.stallTime = -1),
              this.video.paused || this.state === s.pb || this.setState(s.nb);
          },
          loadedmetadata: function () {
            var e = {
                metadataType: "media",
                duration: this.getDuration(),
                height: this.video.videoHeight,
                width: this.video.videoWidth,
                seekRange: this.getSeekRange(),
              },
              t = this.drmUsed;
            t && (e.drm = t), this.trigger(s.K, e);
          },
          timeupdate: function () {
            var e = this.getVideoCurrentTime(),
              t = this.getCurrentTime(),
              i = this.getDuration();
            if (!isNaN(i)) {
              this.seeking ||
                this.video.paused ||
                (this.state !== s.qb && this.state !== s.nb) ||
                this.stallTime === e ||
                ((this.stallTime = -1),
                this.setState(s.pb),
                this.trigger(s.fb));
              var r = {
                position: t,
                duration: i,
                currentTime: e,
                seekRange: this.getSeekRange(),
                metadata: { currentTime: e },
              };
              if (this.getPtsOffset) {
                var n = this.getPtsOffset();
                n >= 0 && (r.metadata.mpegts = n + t);
              }
              var a = this.getLiveLatency();
              null !== a && (r.latency = a),
                (this.state === s.pb || this.seeking) && this.trigger(s.S, r);
            }
          },
          click: function (e) {
            this.trigger(s.n, e);
          },
          volumechange: function () {
            var e = this.video;
            this.trigger(s.V, { volume: Math.round(100 * e.volume) }),
              this.trigger(s.M, { mute: e.muted });
          },
          seeked: function () {
            this.seeking && ((this.seeking = !1), this.trigger(s.R));
          },
          playing: function () {
            -1 === this.stallTime && this.setState(s.pb), this.trigger(s.fb);
          },
          pause: function () {
            this.state !== s.kb &&
              (this.video.ended ||
                this.video.error ||
                (this.getVideoCurrentTime() !== this.getDuration() &&
                  this.setState(s.ob)));
          },
          progress: function () {
            var e = this.getDuration();
            if (!(e <= 0 || e === 1 / 0)) {
              var t = this.video.buffered;
              if (t && 0 !== t.length) {
                var i = Object(c.a)(t.end(t.length - 1) / e, 0, 1);
                this.trigger(s.D, {
                  bufferPercent: 100 * i,
                  position: this.getCurrentTime(),
                  duration: e,
                  currentTime: this.getVideoCurrentTime(),
                  seekRange: this.getSeekRange(),
                });
              }
            }
          },
          ratechange: function () {
            this.trigger(s.P, { playbackRate: this.video.playbackRate });
          },
          ended: function () {
            (this.videoHeight = 0),
              (this.streamBitrate = -1),
              this.state !== s.mb && this.state !== s.kb && this.trigger(s.F);
          },
          loadeddata: function () {
            this.renderNatively && this.setTextTracks(this.video.textTracks);
          },
        },
        o = i(10);
      function d(e) {
        return e && e.length ? e.end(e.length - 1) : 0;
      }
      var l = {
          container: null,
          volume: function (e) {
            this.video.volume = Math.min(Math.max(0, e / 100), 1);
          },
          mute: function (e) {
            (this.video.muted = !!e),
              this.video.muted || this.video.removeAttribute("muted");
          },
          resize: function (e, t, i) {
            var r = this.video,
              a = r.videoWidth,
              s = r.videoHeight;
            if (e && t && a && s) {
              var c = { objectFit: "", width: "", height: "" };
              if ("uniform" === i) {
                var u = e / t,
                  d = a / s,
                  l = Math.abs(u - d);
                l < 0.09 &&
                  l > 0.0025 &&
                  ((c.objectFit = "fill"), (i = "exactfit"));
              }
              if (
                n.Browser.ie ||
                (n.OS.iOS && n.OS.version.major < 9) ||
                n.Browser.androidNative
              )
                if ("uniform" !== i) {
                  c.objectFit = "contain";
                  var h = e / t,
                    f = a / s,
                    v = 1,
                    T = 1;
                  "none" === i
                    ? (v = T =
                        h > f
                          ? Math.ceil((100 * s) / t) / 100
                          : Math.ceil((100 * a) / e) / 100)
                    : "fill" === i
                    ? (v = T = h > f ? h / f : f / h)
                    : "exactfit" === i &&
                      (h > f ? ((v = h / f), (T = 1)) : ((v = 1), (T = f / h))),
                    Object(o.e)(
                      r,
                      "matrix("
                        .concat(v.toFixed(2), ", 0, 0, ")
                        .concat(T.toFixed(2), ", 0, 0)")
                    );
                } else (c.top = c.left = c.margin = ""), Object(o.e)(r, "");
              Object(o.d)(r, c);
            }
          },
          getContainer: function () {
            return this.container;
          },
          setContainer: function (e) {
            (this.container = e),
              this.video.parentNode !== e && e.appendChild(this.video);
          },
          remove: function () {
            this.stop(), this.destroy();
            var e = this.container;
            e && e === this.video.parentNode && e.removeChild(this.video);
          },
          atEdgeOfLiveStream: function () {
            if (!this.isLive()) return !1;
            return d(this.video.buffered) - this.video.currentTime <= 2;
          },
        },
        h = {
          eventsOn_: function () {},
          eventsOff_: function () {},
          attachMedia: function () {
            this.eventsOn_();
          },
          detachMedia: function () {
            return this.eventsOff_();
          },
        },
        f = i(65),
        v = i(5),
        T = i(53),
        g = i(7),
        m = i(66),
        k = i(63),
        b = {
          TIT2: "title",
          TT2: "title",
          WXXX: "url",
          TPE1: "artist",
          TP1: "artist",
          TALB: "album",
          TAL: "album",
        };
      function y(e, t) {
        for (var i, r, n, a = e.length, s = "", c = t || 0; c < a; )
          if (0 !== (i = e[c++]) && 3 !== i)
            switch (i >> 4) {
              case 0:
              case 1:
              case 2:
              case 3:
              case 4:
              case 5:
              case 6:
              case 7:
                s += String.fromCharCode(i);
                break;
              case 12:
              case 13:
                (r = e[c++]),
                  (s += String.fromCharCode(((31 & i) << 6) | (63 & r)));
                break;
              case 14:
                (r = e[c++]),
                  (n = e[c++]),
                  (s += String.fromCharCode(
                    ((15 & i) << 12) | ((63 & r) << 6) | ((63 & n) << 0)
                  ));
            }
        return s;
      }
      function p(e) {
        var t = (function (e) {
          for (var t = "0x", i = 0; i < e.length; i++)
            e[i] < 16 && (t += "0"), (t += e[i].toString(16));
          return parseInt(t);
        })(e);
        return (
          (127 & t) |
          ((32512 & t) >> 1) |
          ((8323072 & t) >> 2) |
          ((2130706432 & t) >> 3)
        );
      }
      function x() {
        return (arguments.length > 0 && void 0 !== arguments[0]
          ? arguments[0]
          : []
        ).reduce(function (e, t) {
          if (!("value" in t) && "data" in t && t.data instanceof ArrayBuffer) {
            var i = new Uint8Array(t.data),
              r = i.length;
            t = { value: { key: "", data: "" } };
            for (var n = 10; n < 14 && n < i.length && 0 !== i[n]; )
              (t.value.key += String.fromCharCode(i[n])), n++;
            var a = 19,
              s = i[a];
            (3 !== s && 0 !== s) || ((s = i[++a]), r--);
            var c = 0;
            if (1 !== s && 2 !== s)
              for (var u = a + 1; u < r; u++)
                if (0 === i[u]) {
                  c = u - a;
                  break;
                }
            if (c > 0) {
              var o = y(i.subarray(a, (a += c)), 0);
              if ("PRIV" === t.value.key) {
                if ("com.apple.streaming.transportStreamTimestamp" === o) {
                  var d = 1 & p(i.subarray(a, (a += 4))),
                    l = p(i.subarray(a, (a += 4))) + (d ? 4294967296 : 0);
                  t.value.data = l;
                } else t.value.data = y(i, a + 1);
                t.value.info = o;
              } else (t.value.info = o), (t.value.data = y(i, a + 1));
            } else {
              var h = i[a];
              t.value.data =
                1 === h || 2 === h
                  ? (function (e, t) {
                      for (var i = e.length - 1, r = "", n = t || 0; n < i; )
                        (254 === e[n] && 255 === e[n + 1]) ||
                          (r += String.fromCharCode((e[n] << 8) + e[n + 1])),
                          (n += 2);
                      return r;
                    })(i, a + 1)
                  : y(i, a + 1);
            }
          }
          if (
            (b.hasOwnProperty(t.value.key) &&
              (e[b[t.value.key]] = t.value.data),
            t.value.info)
          ) {
            var f = e[t.value.key];
            f !== Object(f) && ((f = {}), (e[t.value.key] = f)),
              (f[t.value.info] = t.value.data);
          } else e[t.value.key] = t.value.data;
          return e;
        }, {});
      }
      function _(e, t, i) {
        e &&
          (e.removeEventListener
            ? e.removeEventListener(t, i)
            : (e["on" + t] = null));
      }
      function w() {
        var e = this.video.textTracks,
          t = Object(r.h)(e, function (e) {
            return (e.inuse || !e._id) && E(e.kind);
          });
        if (this._textTracks && !M.call(this, t)) {
          for (var i = -1, n = 0; n < this._textTracks.length; n++)
            if ("showing" === this._textTracks[n].mode) {
              i = n;
              break;
            }
          i !== this._currentTextTrackIndex && this.setSubtitlesTrack(i + 1);
        } else this.setTextTracks(e);
      }
      function C() {
        this.setTextTracks(this.video.textTracks);
      }
      function O(e) {
        var t = this;
        e &&
          (this._textTracks || this._initTextTracks(),
          e.forEach(function (e) {
            if (!e.kind || E(e.kind)) {
              var i = S.call(t, e);
              L.call(t, i),
                e.file &&
                  ((e.data = []),
                  Object(m.c)(
                    e,
                    function (e) {
                      t.addVTTCuesToTrack(i, e);
                    },
                    function (e) {
                      t.trigger(s.tb, e);
                    }
                  ));
            }
          }),
          this._textTracks &&
            this._textTracks.length &&
            this.trigger("subtitlesTracks", { tracks: this._textTracks }));
      }
      function j(e, t, i) {
        if (n.Browser.ie) {
          var r = i;
          (e || "metadata" === t.kind) &&
            (r = new window.TextTrackCue(i.startTime, i.endTime, i.text)),
            (function (e, t) {
              var i = [],
                r = e.mode;
              e.mode = "hidden";
              for (
                var n = e.cues, a = n.length - 1;
                a >= 0 && n[a].startTime > t.startTime;
                a--
              )
                i.unshift(n[a]), e.removeCue(n[a]);
              try {
                e.addCue(t),
                  i.forEach(function (t) {
                    return e.addCue(t);
                  });
              } catch (e) {
                console.error(e);
              }
              e.mode = r;
            })(t, r);
        } else
          try {
            t.addCue(i);
          } catch (e) {
            console.error(e);
          }
      }
      function B(e, t) {
        t &&
          t.length &&
          Object(r.f)(t, function (t) {
            if (!(n.Browser.ie && e && /^(native|subtitle|cc)/.test(t._id))) {
              (n.Browser.ie && "disabled" === t.mode) ||
                ((t.mode = "disabled"), (t.mode = "hidden"));
              for (var i = t.cues.length; i--; ) t.removeCue(t.cues[i]);
              t.embedded || (t.mode = "disabled"), (t.inuse = !1);
            }
          });
      }
      function E(e) {
        return "subtitles" === e || "captions" === e;
      }
      function S(e) {
        var t,
          i = Object(k.b)(e, this._unknownCount),
          n = i.label;
        if (
          ((this._unknownCount = i.unknownCount),
          this.renderNatively || "metadata" === e.kind)
        ) {
          var a = this.video.textTracks;
          (t = Object(r.j)(a, { label: n })) ||
            (t = this.video.addTextTrack(e.kind, n, e.language || "")),
            (t.default = e.default),
            (t.mode = "disabled"),
            (t.inuse = !0);
        } else (t = e).data = t.data || [];
        return t._id || (t._id = Object(k.a)(e, this._textTracks.length)), t;
      }
      function L(e) {
        this._textTracks.push(e), (this._tracksById[e._id] = e);
      }
      function I() {
        if (this._textTracks) {
          var e = this._textTracks.filter(function (e) {
            return e.embedded || "subs" === e.groupid;
          });
          this._initTextTracks(),
            e.forEach(function (e) {
              this._tracksById[e._id] = e;
            }),
            (this._textTracks = e);
        }
      }
      function N(e) {
        this.triggerActiveCues(e.currentTarget.activeCues);
      }
      function R(e, t, i) {
        var r = e.kind;
        this._cachedVTTCues[e._id] || (this._cachedVTTCues[e._id] = {});
        var n,
          a = this._cachedVTTCues[e._id];
        switch (r) {
          case "captions":
          case "subtitles":
            n = i || Math.floor(20 * t.startTime);
            var s = "_" + t.line,
              c = Math.floor(20 * t.endTime),
              u = a[n + s] || a[n + 1 + s] || a[n - 1 + s];
            return !(u && Math.abs(u - c) <= 1) && ((a[n + s] = c), !0);
          case "metadata":
            var o = t.data ? new Uint8Array(t.data).join("") : t.text;
            return !a[(n = i || t.startTime + o)] && ((a[n] = t.endTime), !0);
          default:
            return !1;
        }
      }
      function M(e) {
        if (e.length > this._textTracks.length) return !0;
        for (var t = 0; t < e.length; t++) {
          var i = e[t];
          if (!i._id || !this._tracksById[i._id]) return !0;
        }
        return !1;
      }
      var V = {
          _itemTracks: null,
          _textTracks: null,
          _tracksById: null,
          _cuesByTrackId: null,
          _cachedVTTCues: null,
          _metaCuesByTextTime: null,
          _currentTextTrackIndex: -1,
          _unknownCount: 0,
          _activeCues: null,
          _initTextTracks: function () {
            (this._textTracks = []),
              (this._tracksById = {}),
              (this._metaCuesByTextTime = {}),
              (this._cuesByTrackId = {}),
              (this._cachedVTTCues = {}),
              (this._unknownCount = 0);
          },
          addTracksListener: function (e, t, i) {
            if (!e) return;
            if ((_(e, t, i), this.instreamMode)) return;
            e.addEventListener ? e.addEventListener(t, i) : (e["on" + t] = i);
          },
          clearTracks: function () {
            Object(m.a)(this._itemTracks);
            var e = this._tracksById && this._tracksById.nativemetadata;
            (this.renderNatively || e) &&
              (B(this.renderNatively, this.video.textTracks),
              e && (e.oncuechange = null));
            (this._itemTracks = null),
              (this._textTracks = null),
              (this._tracksById = null),
              (this._cuesByTrackId = null),
              (this._metaCuesByTextTime = null),
              (this._unknownCount = 0),
              (this._currentTextTrackIndex = -1),
              (this._activeCues = null),
              this.renderNatively &&
                (this.removeTracksListener(
                  this.video.textTracks,
                  "change",
                  this.textTrackChangeHandler
                ),
                B(this.renderNatively, this.video.textTracks));
          },
          clearMetaCues: function () {
            var e = this._tracksById && this._tracksById.nativemetadata;
            e &&
              (B(this.renderNatively, [e]),
              (e.mode = "hidden"),
              (e.inuse = !0),
              (this._cachedVTTCues[e._id] = {}));
          },
          clearCueData: function (e) {
            var t = this._cachedVTTCues;
            t &&
              t[e] &&
              ((t[e] = {}),
              this._tracksById && (this._tracksById[e].data = []));
          },
          disableTextTrack: function () {
            if (this._textTracks) {
              var e = this._textTracks[this._currentTextTrackIndex];
              if (e) {
                e.mode = "disabled";
                var t = e._id;
                t && 0 === t.indexOf("nativecaptions") && (e.mode = "hidden");
              }
            }
          },
          enableTextTrack: function () {
            if (this._textTracks) {
              var e = this._textTracks[this._currentTextTrackIndex];
              e && (e.mode = "showing");
            }
          },
          getSubtitlesTrack: function () {
            return this._currentTextTrackIndex;
          },
          removeTracksListener: _,
          addTextTracks: O,
          setTextTracks: function (e) {
            if (((this._currentTextTrackIndex = -1), !e)) return;
            this._textTracks
              ? ((this._unknownCount = 0),
                (this._textTracks = this._textTracks.filter(function (e) {
                  var t = e._id;
                  return this.renderNatively &&
                    t &&
                    0 === t.indexOf("nativecaptions")
                    ? (delete this._tracksById[t], !1)
                    : (e.name &&
                        0 === e.name.indexOf("Unknown") &&
                        this._unknownCount++,
                      !0);
                }, this)),
                delete this._tracksById.nativemetadata)
              : this._initTextTracks();
            if (e.length)
              for (var t = 0, i = e.length; t < i; t++) {
                var r = e[t];
                if (!r._id) {
                  if ("captions" === r.kind || "metadata" === r.kind) {
                    if (
                      ((r._id = "native" + r.kind + t),
                      !r.label && "captions" === r.kind)
                    ) {
                      var a = Object(k.b)(r, this._unknownCount);
                      (r.name = a.label), (this._unknownCount = a.unknownCount);
                    }
                  } else r._id = Object(k.a)(r, this._textTracks.length);
                  if (this._tracksById[r._id]) continue;
                  r.inuse = !0;
                }
                if (r.inuse && !this._tracksById[r._id])
                  if ("metadata" === r.kind)
                    (r.mode = "hidden"),
                      (r.oncuechange = N.bind(this)),
                      (this._tracksById[r._id] = r);
                  else if (E(r.kind)) {
                    var s = r.mode,
                      c = void 0;
                    if (((r.mode = "hidden"), !r.cues.length && r.embedded))
                      continue;
                    if (
                      ((r.mode = s),
                      this._cuesByTrackId[r._id] &&
                        !this._cuesByTrackId[r._id].loaded)
                    ) {
                      for (
                        var u = this._cuesByTrackId[r._id].cues;
                        (c = u.shift());

                      )
                        j(this.renderNatively, r, c);
                      (r.mode = s), (this._cuesByTrackId[r._id].loaded = !0);
                    }
                    L.call(this, r);
                  }
              }
            this.renderNatively &&
              ((this.textTrackChangeHandler =
                this.textTrackChangeHandler || w.bind(this)),
              this.addTracksListener(
                this.video.textTracks,
                "change",
                this.textTrackChangeHandler
              ),
              (n.Browser.edge || n.Browser.firefox || n.Browser.safari) &&
                ((this.addTrackHandler = this.addTrackHandler || C.bind(this)),
                this.addTracksListener(
                  this.video.textTracks,
                  "addtrack",
                  this.addTrackHandler
                )));
            this._textTracks.length &&
              this.trigger("subtitlesTracks", { tracks: this._textTracks });
          },
          setupSideloadedTracks: function (e) {
            if (!this.renderNatively) return;
            var t = e === this._itemTracks;
            t || Object(m.a)(this._itemTracks);
            if (((this._itemTracks = e), !e)) return;
            t || (this.disableTextTrack(), I.call(this), this.addTextTracks(e));
          },
          setSubtitlesTrack: function (e) {
            if (!this.renderNatively)
              return void (
                this.setCurrentSubtitleTrack &&
                this.setCurrentSubtitleTrack(e - 1)
              );
            if (!this._textTracks) return;
            0 === e &&
              this._textTracks.forEach(function (e) {
                e.mode = e.embedded ? "hidden" : "disabled";
              });
            if (this._currentTextTrackIndex === e - 1) return;
            this.disableTextTrack(),
              (this._currentTextTrackIndex = e - 1),
              this._textTracks[this._currentTextTrackIndex] &&
                (this._textTracks[this._currentTextTrackIndex].mode =
                  "showing");
            this.trigger("subtitlesTrackChanged", {
              currentTrack: this._currentTextTrackIndex + 1,
              tracks: this._textTracks,
            });
          },
          textTrackChangeHandler: null,
          addTrackHandler: null,
          addCuesToTrack: function (e) {
            var t = this._tracksById[e.name];
            if (!t) return;
            t.source = e.source;
            for (
              var i = e.captions || [], r = [], n = !1, a = 0;
              a < i.length;
              a++
            ) {
              var s = i[a],
                c = e.name + "_" + s.begin + "_" + s.end;
              this._metaCuesByTextTime[c] ||
                ((this._metaCuesByTextTime[c] = s), r.push(s), (n = !0));
            }
            n &&
              r.sort(function (e, t) {
                return e.begin - t.begin;
              });
            var u = Object(m.b)(r);
            Array.prototype.push.apply(t.data, u);
          },
          addCaptionsCue: function (e) {
            if (!e.text || !e.begin || !e.end) return;
            var t,
              i = e.trackid.toString(),
              r = this._tracksById && this._tracksById[i];
            r ||
              ((r = { kind: "captions", _id: i, data: [] }),
              this.addTextTracks([r]),
              this.trigger("subtitlesTracks", { tracks: this._textTracks }));
            e.useDTS && (r.source || (r.source = e.source || "mpegts"));
            t = e.begin + "_" + e.text;
            var n = this._metaCuesByTextTime[t];
            if (!n) {
              (n = { begin: e.begin, end: e.end, text: e.text }),
                (this._metaCuesByTextTime[t] = n);
              var a = Object(m.b)([n])[0];
              r.data.push(a);
            }
          },
          createCue: function (e, t, i) {
            var r = window.VTTCue || window.TextTrackCue,
              n = Math.max(t || 0, e + 0.25);
            return new r(e, n, i);
          },
          addVTTCue: function (e, t) {
            this._tracksById || this._initTextTracks();
            var i = e.track ? e.track : "native" + e.type,
              r = this._tracksById[i],
              n = "captions" === e.type ? "Unknown CC" : "ID3 Metadata",
              a = e.cue;
            if (!r) {
              var s = { kind: e.type, _id: i, label: n, embedded: !0 };
              (r = S.call(this, s)),
                this.renderNatively || "metadata" === r.kind
                  ? this.setTextTracks(this.video.textTracks)
                  : O.call(this, [r]);
            }
            if (R.call(this, r, a, t)) {
              var c = this.renderNatively || "metadata" === r.kind;
              return c ? j(c, r, a) : r.data.push(a), a;
            }
            return null;
          },
          addVTTCuesToTrack: function (e, t) {
            if (!this.renderNatively) return;
            var i,
              r = this._tracksById[e._id];
            if (!r)
              return (
                this._cuesByTrackId || (this._cuesByTrackId = {}),
                void (this._cuesByTrackId[e._id] = { cues: t, loaded: !1 })
              );
            if (this._cuesByTrackId[e._id] && this._cuesByTrackId[e._id].loaded)
              return;
            this._cuesByTrackId[e._id] = { cues: t, loaded: !0 };
            for (; (i = t.shift()); ) j(this.renderNatively, r, i);
          },
          triggerActiveCues: function (e) {
            var t = this;
            if (!e || !e.length) return void (this._activeCues = null);
            var i = this._activeCues || [],
              r = Array.prototype.filter.call(e, function (e) {
                if (
                  i.some(function (t) {
                    return (
                      (r = t),
                      (i = e).startTime === r.startTime &&
                        i.endTime === r.endTime &&
                        i.text === r.text &&
                        i.data === r.data &&
                        i.value === r.value
                    );
                    var i, r;
                  })
                )
                  return !1;
                if (e.data || e.value) return !0;
                if (e.text) {
                  var r = JSON.parse(e.text),
                    n = { metadataTime: e.startTime, metadata: r };
                  r.programDateTime && (n.programDateTime = r.programDateTime),
                    r.metadataType &&
                      ((n.metadataType = r.metadataType),
                      delete r.metadataType),
                    t.trigger(s.K, n);
                }
                return !1;
              });
            if (r.length) {
              var n = x(r),
                a = r[0].startTime;
              this.trigger(s.K, {
                metadataType: "id3",
                metadataTime: a,
                metadata: n,
              });
            }
            this._activeCues = Array.prototype.slice.call(e);
          },
          renderNatively: !1,
        },
        A = i(64),
        D = i(15),
        P = i(1),
        H = 224e3,
        F = 224005,
        W = 221e3,
        U = 324e3,
        X = window.clearTimeout,
        Q = "html5",
        q = function () {};
      function J(e, t) {
        Object.keys(e).forEach(function (i) {
          t.removeEventListener(i, e[i]);
        });
      }
      function z(e, t, i) {
        (this.state = s.mb),
          (this.seeking = !1),
          (this.currentTime = -1),
          (this.retries = 0),
          (this.maxRetries = 3);
        var c,
          T = this,
          m = t.minDvrWindow,
          k = {
            progress: function () {
              u.progress.call(T), ve();
            },
            timeupdate: function () {
              T.currentTime >= 0 && (T.retries = 0);
              var e = T.getVideoCurrentTime();
              (T.currentTime = e),
                B && O !== e && Z(e),
                u.timeupdate.call(T),
                ve(),
                n.Browser.ie && Y();
            },
            resize: Y,
            ended: function () {
              (j = -1), Te(), u.ended.call(T);
            },
            loadedmetadata: function () {
              var e = T.getDuration();
              M && e === 1 / 0 && (e = 0);
              var t = {
                metadataType: "media",
                duration: e,
                height: y.videoHeight,
                width: y.videoWidth,
                seekRange: T.getSeekRange(),
              };
              T.trigger(s.K, t), Y();
            },
            durationchange: function () {
              M || u.progress.call(T);
            },
            loadeddata: function () {
              var e;
              !(function () {
                if (y.getStartDate) {
                  var e = y.getStartDate(),
                    t = e.getTime ? e.getTime() : NaN;
                  if (t !== T.startDateTime && !isNaN(t)) {
                    T.startDateTime = t;
                    var i = e.toISOString(),
                      r = T.getSeekRange(),
                      n = r.start,
                      a = r.end,
                      c = {
                        metadataType: "program-date-time",
                        programDateTime: i,
                        start: n,
                        end: a,
                      },
                      u = T.createCue(n, a, JSON.stringify(c));
                    T.addVTTCue({ type: "metadata", cue: u }),
                      delete c.metadataType,
                      T.trigger(s.L, {
                        metadataType: "program-date-time",
                        metadata: c,
                      });
                  }
                }
              })(),
                u.loadeddata.call(T),
                (function (e) {
                  if (((S = null), !e)) return;
                  if (e.length) {
                    for (var t = 0; t < e.length; t++)
                      if (e[t].enabled) {
                        L = t;
                        break;
                      }
                    -1 === L && (e[(L = 0)].enabled = !0),
                      (S = Object(r.v)(e, function (e) {
                        return {
                          name: e.label || e.language,
                          language: e.language,
                        };
                      }));
                  }
                  T.addTracksListener(e, "change", oe),
                    S &&
                      T.trigger("audioTracks", { currentTrack: L, tracks: S });
                })(y.audioTracks),
                (e = T.getDuration()),
                w && -1 !== w && e && e !== 1 / 0 && T.seek(w),
                Y();
            },
            canplay: function () {
              (_ = !0),
                M || fe(),
                n.Browser.ie &&
                  9 === n.Browser.version.major &&
                  T.setTextTracks(T._textTracks),
                u.canplay.call(T);
            },
            seeking: function () {
              var e = null !== C ? ee(C) : T.getCurrentTime(),
                t = ee(O);
              (O = C),
                (C = null),
                (w = 0),
                (T.seeking = !0),
                T.trigger(s.Q, { position: t, offset: e });
            },
            seeked: function () {
              u.seeked.call(T);
            },
            waiting: function () {
              T.seeking
                ? T.setState(s.nb)
                : T.state === s.pb &&
                  (T.atEdgeOfLiveStream() && T.setPlaybackRate(1),
                  (T.stallTime = T.video.currentTime),
                  T.setState(s.qb));
            },
            webkitbeginfullscreen: function (e) {
              (B = !0), de(e);
            },
            webkitendfullscreen: function (e) {
              (B = !1), de(e);
            },
            error: function () {
              var e = T.video,
                t = e.error,
                i = (t && t.code) || -1;
              if ((3 === i || 4 === i) && T.retries < T.maxRetries)
                return (
                  T.trigger(s.tb, new P.n(null, U + i - 1, t)),
                  T.retries++,
                  y.load(),
                  void (
                    -1 !== T.currentTime &&
                    ((_ = !1), T.seek(T.currentTime), (T.currentTime = -1))
                  )
                );
              var r = H,
                n = P.k;
              1 === i
                ? (r += i)
                : 2 === i
                ? ((n = P.i), (r = W))
                : 3 === i || 4 === i
                ? ((r += i - 1), 4 === i && e.src === location.href && (r = F))
                : (n = P.m),
                se(),
                T.trigger(s.G, new P.n(n, r, t));
            },
          };
        Object.keys(u).forEach(function (e) {
          if (!k[e]) {
            var t = u[e];
            k[e] = function (e) {
              t.call(T, e);
            };
          }
        }),
          Object(r.g)(this, g.a, l, h, V, {
            renderNatively:
              ((c = t.renderCaptionsNatively),
              !(!n.OS.iOS && !n.Browser.safari) || (c && n.Browser.chrome)),
            eventsOn_: function () {
              var e, t;
              (e = k),
                (t = y),
                Object.keys(e).forEach(function (i) {
                  t.removeEventListener(i, e[i]), t.addEventListener(i, e[i]);
                });
            },
            eventsOff_: function () {
              J(k, y);
            },
            detachMedia: function () {
              h.detachMedia.call(T),
                Te(),
                this.removeTracksListener(
                  y.textTracks,
                  "change",
                  this.textTrackChangeHandler
                ),
                this.disableTextTrack();
            },
            attachMedia: function () {
              h.attachMedia.call(T),
                (_ = !1),
                (this.seeking = !1),
                (y.loop = !1),
                this.enableTextTrack(),
                this.renderNatively &&
                  this.setTextTracks(this.video.textTracks),
                this.addTracksListener(
                  y.textTracks,
                  "change",
                  this.textTrackChangeHandler
                );
            },
            isLive: function () {
              return this.getDuration() === 1 / 0;
            },
          });
        var b,
          y = i,
          p = { level: {} },
          x = null !== t.liveTimeout ? t.liveTimeout : 3e4,
          _ = !1,
          w = 0,
          C = null,
          O = null,
          j = -1,
          B = !1,
          E = q,
          S = null,
          L = -1,
          I = -1,
          N = !1,
          R = null,
          M = !1,
          z = null,
          G = null,
          $ = 0;
        function Y() {
          var e = p.level;
          if (e.width !== y.videoWidth || e.height !== y.videoHeight) {
            if ((!y.videoWidth && !he()) || -1 === j) return;
            (e.width = y.videoWidth),
              (e.height = y.videoHeight),
              fe(),
              (p.reason = p.reason || "auto"),
              (p.mode = "hls" === b[j].type ? "auto" : "manual"),
              (p.bitrate = 0),
              (e.index = j),
              (e.label = b[j].label),
              T.trigger(s.U, p),
              (p.reason = "");
          }
        }
        function Z(e) {
          O = e;
        }
        function ee(e) {
          var t = T.getSeekRange();
          return T.isLive() && Object(f.a)(t.end - t.start, m)
            ? Math.min(0, e - t.end)
            : e;
        }
        function te(e) {
          var t;
          return (
            Array.isArray(e) &&
              e.length > 0 &&
              (t = e.map(function (e, t) {
                return { label: e.label || t };
              })),
            t
          );
        }
        function ie(e) {
          (T.currentTime = -1),
            (m = e.minDvrWindow),
            (b = e.sources),
            (j = (function (e) {
              var i = Math.max(0, j),
                r = t.qualityLabel;
              if (e)
                for (var n = 0; n < e.length; n++)
                  if ((e[n].default && (i = n), r && e[n].label === r))
                    return n;
              (p.reason = "initial choice"),
                (p.level.width && p.level.height) || (p.level = {});
              return i;
            })(b));
        }
        function re() {
          return (
            y.paused &&
              y.played &&
              y.played.length &&
              T.isLive() &&
              !Object(f.a)(ue() - ce(), m) &&
              (T.clearTracks(), y.load()),
            y.play() || Object(A.a)(y)
          );
        }
        function ne(e) {
          (T.currentTime = -1), (w = 0), Te();
          var t = y.src,
            i = document.createElement("source");
          (i.src = b[j].file),
            i.src !== t
              ? (ae(b[j]), t && y.load())
              : 0 === e && T.getVideoCurrentTime() > 0 && ((w = -1), T.seek(e)),
            e > 0 && T.getVideoCurrentTime() !== e && T.seek(e);
          var r = te(b);
          r && T.trigger(s.I, { levels: r, currentQuality: j }),
            b.length && "hls" !== b[0].type && fe();
        }
        function ae(e) {
          (S = null),
            (L = -1),
            p.reason || ((p.reason = "initial choice"), (p.level = {})),
            (_ = !1);
          var t = document.createElement("source");
          (t.src = e.file), y.src !== t.src && (y.src = e.file);
        }
        function se() {
          y &&
            (T.disableTextTrack(),
            y.removeAttribute("preload"),
            y.removeAttribute("src"),
            Object(v.h)(y),
            Object(o.d)(y, { objectFit: "" }),
            (j = -1),
            !n.Browser.msie && "load" in y && y.load());
        }
        function ce() {
          var e = 1 / 0;
          return (
            ["buffered", "seekable"].forEach(function (t) {
              for (var i = y[t], n = i ? i.length : 0; n--; ) {
                var a = Math.min(e, i.start(n));
                Object(r.o)(a) && (e = a);
              }
            }),
            e
          );
        }
        function ue() {
          var e = 0;
          return (
            ["buffered", "seekable"].forEach(function (t) {
              for (var i = y[t], n = i ? i.length : 0; n--; ) {
                var a = Math.max(e, i.end(n));
                Object(r.o)(a) && (e = a);
              }
            }),
            e
          );
        }
        function oe() {
          for (var e = -1, t = 0; t < y.audioTracks.length; t++)
            if (y.audioTracks[t].enabled) {
              e = t;
              break;
            }
          le(e);
        }
        function de(e) {
          T.trigger(s.X, { target: e.target, jwstate: B });
        }
        function le(e) {
          y &&
            y.audioTracks &&
            S &&
            e > -1 &&
            e < y.audioTracks.length &&
            e !== L &&
            ((y.audioTracks[L].enabled = !1),
            (L = e),
            (y.audioTracks[L].enabled = !0),
            T.trigger("audioTrackChanged", { currentTrack: L, tracks: S }));
        }
        function he() {
          if (!(y.readyState < 2)) return 0 === y.videoHeight;
        }
        function fe() {
          var e = he();
          if (void 0 !== e) {
            var t = e ? "audio" : "video";
            T.trigger(s.T, { mediaType: t });
          }
        }
        function ve() {
          if (0 !== x) {
            var e = d(y.buffered);
            T.isLive() && e && R === e
              ? -1 === I &&
                (I = setTimeout(function () {
                  (N = !0),
                    (function () {
                      if (N && T.atEdgeOfLiveStream())
                        return T.trigger(s.G, new P.n(P.l, K)), !0;
                    })();
                }, x))
              : (Te(), (N = !1)),
              (R = e);
          }
        }
        function Te() {
          X(I), (I = -1);
        }
        (this.video = y),
          (this.supportsPlaybackRate = !0),
          (this.startDateTime = 0),
          (T.getVideoCurrentTime = function () {
            return t.getCurrentTimeHook
              ? t.getCurrentTimeHook(y)
              : y.currentTime;
          }),
          (T.getCurrentTime = function () {
            return (function (e) {
              var t = T.getSeekRange();
              if (T.isLive()) {
                if (
                  ((!G || Math.abs(z - t.end) > 1) &&
                    (function (e) {
                      (z = e.end),
                        (G = Math.min(0, T.getVideoCurrentTime() - z)),
                        ($ = Object(D.a)());
                    })(t),
                  Object(f.a)(t.end - t.start, m))
                )
                  return G;
              }
              return e;
            })(T.getVideoCurrentTime());
          }),
          (T.getDuration = function () {
            if (t.getDurationHook) return t.getDurationHook();
            var e = y.duration;
            if ((M && e === 1 / 0 && 0 === T.getVideoCurrentTime()) || isNaN(e))
              return 0;
            var i = ue();
            if (y.duration === 1 / 0 && i) {
              var r = i - ce();
              Object(f.a)(r, m) && (e = -r);
            }
            return e;
          }),
          (T.getSeekRange = function () {
            var e = { start: 0, end: T.getDuration() };
            return y.seekable.length && ((e.end = ue()), (e.start = ce())), e;
          }),
          (T.getLiveLatency = function () {
            var e = null,
              t = ue();
            return (
              T.isLive() &&
                t &&
                (e = t + (Object(D.a)() - $) / 1e3 - T.getVideoCurrentTime()),
              e
            );
          }),
          (this.stop = function () {
            Te(),
              se(),
              this.clearTracks(),
              n.Browser.ie && y.pause(),
              this.setState(s.mb);
          }),
          (this.destroy = function () {
            (E = q),
              J(k, y),
              this.removeTracksListener(y.audioTracks, "change", oe),
              this.removeTracksListener(
                y.textTracks,
                "change",
                T.textTrackChangeHandler
              ),
              this.off();
          }),
          (this.init = function (e) {
            (T.retries = 0), (T.maxRetries = e.adType ? 0 : 3), ie(e);
            var t = b[j];
            (M = Object(a.a)(t)) &&
              ((T.supportsPlaybackRate = !1), (k.waiting = q)),
              T.eventsOn_(),
              b.length && "hls" !== b[0].type && this.sendMediaType(b),
              (p.reason = "");
          }),
          (this.preload = function (e) {
            ie(e);
            var t = b[j],
              i = t.preload || "metadata";
            "none" !== i && (y.setAttribute("preload", i), ae(t));
          }),
          (this.load = function (e) {
            ie(e), ne(e.starttime), this.setupSideloadedTracks(e.tracks);
          }),
          (this.play = function () {
            return E(), re();
          }),
          (this.pause = function () {
            Te(),
              (E = function () {
                if (y.paused && T.getVideoCurrentTime() && T.isLive()) {
                  var e = ue(),
                    t = e - ce(),
                    i = !Object(f.a)(t, m),
                    n = e - T.getVideoCurrentTime();
                  if (i && e && (n > 15 || n < 0)) {
                    if (((C = Math.max(e - 10, e - t)), !Object(r.o)(C)))
                      return;
                    Z(T.getVideoCurrentTime()), (y.currentTime = C);
                  }
                }
              }),
              y.pause();
          }),
          (this.seek = function (e) {
            if (!t.seekHook || !t.seekHook(e, y)) {
              var i = T.getSeekRange(),
                r = e;
              if ((e < 0 && (r += i.end), _ || (_ = !!ue()), _)) {
                w = 0;
                try {
                  if (
                    ((T.seeking = !0),
                    T.isLive() && Object(f.a)(i.end - i.start, m))
                  )
                    if (((G = Math.min(0, r - z)), e < 0))
                      r += Math.min(12, (Object(D.a)() - $) / 1e3);
                  (C = r), Z(T.getVideoCurrentTime()), (y.currentTime = r);
                } catch (e) {
                  (T.seeking = !1), (w = r);
                }
              } else (w = r), n.Browser.firefox && y.paused && re();
            }
          }),
          (this.setVisibility = function (e) {
            (e = !!e) || n.OS.android
              ? Object(o.d)(T.container, { visibility: "visible", opacity: 1 })
              : Object(o.d)(T.container, { visibility: "", opacity: 0 });
          }),
          (this.setFullscreen = function (e) {
            if ((e = !!e)) {
              try {
                var t = y.webkitEnterFullscreen || y.webkitEnterFullScreen;
                t && t.apply(y);
              } catch (e) {
                return !1;
              }
              return T.getFullScreen();
            }
            var i = y.webkitExitFullscreen || y.webkitExitFullScreen;
            return i && i.apply(y), e;
          }),
          (T.getFullScreen = function () {
            return B || !!y.webkitDisplayingFullscreen;
          }),
          (this.setCurrentQuality = function (e) {
            j !== e &&
              e >= 0 &&
              b &&
              b.length > e &&
              ((j = e),
              (p.reason = "api"),
              (p.level = {}),
              this.trigger(s.J, { currentQuality: e, levels: te(b) }),
              (t.qualityLabel = b[e].label),
              ne(T.getVideoCurrentTime() || 0),
              re());
          }),
          (this.setPlaybackRate = function (e) {
            y.playbackRate = y.defaultPlaybackRate = e;
          }),
          (this.getPlaybackRate = function () {
            return y.playbackRate;
          }),
          (this.getCurrentQuality = function () {
            return j;
          }),
          (this.getQualityLevels = function () {
            return Array.isArray(b)
              ? b.map(function (e) {
                  return (function (e) {
                    return {
                      bitrate: e.bitrate,
                      label: e.label,
                      width: e.width,
                      height: e.height,
                    };
                  })(e);
                })
              : [];
          }),
          (this.getName = function () {
            return { name: Q };
          }),
          (this.setCurrentAudioTrack = le),
          (this.getAudioTracks = function () {
            return S || [];
          }),
          (this.getCurrentAudioTrack = function () {
            return L;
          });
      }
      Object(r.g)(z.prototype, T.a),
        (z.getName = function () {
          return { name: "html5" };
        });
      t.default = z;
      var K = 220001;
    },
    57: function (e, t, i) {
      "use strict";
      i.d(t, "a", function () {
        return n;
      });
      var r = i(2);
      function n(e) {
        var t = [],
          i = (e = Object(r.i)(e)).split("\r\n\r\n");
        1 === i.length && (i = e.split("\n\n"));
        for (var n = 0; n < i.length; n++)
          if ("WEBVTT" !== i[n]) {
            var s = a(i[n]);
            s.text && t.push(s);
          }
        return t;
      }
      function a(e) {
        var t = {},
          i = e.split("\r\n");
        1 === i.length && (i = e.split("\n"));
        var n = 1;
        if (
          (i[0].indexOf(" --\x3e ") > 0 && (n = 0),
          i.length > n + 1 && i[n + 1])
        ) {
          var a = i[n],
            s = a.indexOf(" --\x3e ");
          s > 0 &&
            ((t.begin = Object(r.g)(a.substr(0, s))),
            (t.end = Object(r.g)(a.substr(s + 5))),
            (t.text = i.slice(n + 1).join("\r\n")));
        }
        return t;
      }
    },
    63: function (e, t, i) {
      "use strict";
      function r(e, t) {
        var i = e.kind || "cc";
        return e.default || e.defaulttrack
          ? "default"
          : e._id || e.file || i + t;
      }
      function n(e, t) {
        var i = e.label || e.name || e.language;
        return (
          i || ((i = "Unknown CC"), (t += 1) > 1 && (i += " [" + t + "]")),
          { label: i, unknownCount: t }
        );
      }
      i.d(t, "a", function () {
        return r;
      }),
        i.d(t, "b", function () {
          return n;
        });
    },
    64: function (e, t, i) {
      "use strict";
      function r(e) {
        return new Promise(function (t, i) {
          if (e.paused) return i(n("NotAllowedError", 0, "play() failed."));
          var r = function () {
              e.removeEventListener("play", a),
                e.removeEventListener("playing", s),
                e.removeEventListener("pause", s),
                e.removeEventListener("abort", s),
                e.removeEventListener("error", s);
            },
            a = function () {
              e.addEventListener("playing", s),
                e.addEventListener("abort", s),
                e.addEventListener("error", s),
                e.addEventListener("pause", s);
            },
            s = function (e) {
              if ((r(), "playing" === e.type)) t();
              else {
                var a = 'The play() request was interrupted by a "'.concat(
                  e.type,
                  '" event.'
                );
                "error" === e.type
                  ? i(n("NotSupportedError", 9, a))
                  : i(n("AbortError", 20, a));
              }
            };
          e.addEventListener("play", a);
        });
      }
      function n(e, t, i) {
        var r = new Error(i);
        return (r.name = e), (r.code = t), r;
      }
      i.d(t, "a", function () {
        return r;
      });
    },
    65: function (e, t, i) {
      "use strict";
      function r(e, t) {
        return e !== 1 / 0 && Math.abs(e) >= Math.max(a(t), 0);
      }
      function n(e, t) {
        var i = "VOD";
        return (
          e === 1 / 0
            ? (i = "LIVE")
            : e < 0 && (i = r(e, a(t)) ? "DVR" : "LIVE"),
          i
        );
      }
      function a(e) {
        return void 0 === e ? 120 : Math.max(e, 0);
      }
      i.d(t, "a", function () {
        return r;
      }),
        i.d(t, "b", function () {
          return n;
        });
    },
    66: function (e, t, i) {
      "use strict";
      var r = i(67),
        n = i(16),
        a = i(22),
        s = i(4),
        c = i(57),
        u = i(2),
        o = i(1);
      function d(e) {
        throw new o.n(null, e);
      }
      function l(e, t, r) {
        e.xhr = Object(a.a)(
          e.file,
          function (a) {
            !(function (e, t, r, a) {
              var l,
                h,
                v = e.responseXML ? e.responseXML.firstChild : null;
              if (v)
                for (
                  "xml" === Object(s.b)(v) && (v = v.nextSibling);
                  v.nodeType === v.COMMENT_NODE;

                )
                  v = v.nextSibling;
              try {
                if (v && "tt" === Object(s.b)(v))
                  (l = (function (e) {
                    e || d(306007);
                    var t = [],
                      i = e.getElementsByTagName("p"),
                      r = 30,
                      n = e.getElementsByTagName("tt");
                    if (n && n[0]) {
                      var a = parseFloat(n[0].getAttribute("ttp:frameRate"));
                      isNaN(a) || (r = a);
                    }
                    i || d(306005),
                      i.length ||
                        (i = e.getElementsByTagName("tt:p")).length ||
                        (i = e.getElementsByTagName("tts:p"));
                    for (var s = 0; s < i.length; s++) {
                      for (
                        var c = i[s], o = c.getElementsByTagName("br"), l = 0;
                        l < o.length;
                        l++
                      ) {
                        var h = o[l];
                        h.parentNode.replaceChild(e.createTextNode("\r\n"), h);
                      }
                      var f = c.innerHTML || c.textContent || c.text || "",
                        v = Object(u.i)(f)
                          .replace(/>\s+</g, "><")
                          .replace(/(<\/?)tts?:/g, "$1")
                          .replace(/<br.*?\/>/g, "\r\n");
                      if (v) {
                        var T = c.getAttribute("begin"),
                          g = c.getAttribute("dur"),
                          m = c.getAttribute("end"),
                          k = { begin: Object(u.g)(T, r), text: v };
                        m
                          ? (k.end = Object(u.g)(m, r))
                          : g && (k.end = k.begin + Object(u.g)(g, r)),
                          t.push(k);
                      }
                    }
                    return t.length || d(306005), t;
                  })(e.responseXML)),
                    (h = f(l)),
                    delete t.xhr,
                    r(h);
                else {
                  var T = e.responseText;
                  T.indexOf("WEBVTT") >= 0
                    ? i
                        .e(10)
                        .then(
                          function (e) {
                            return i(97).default;
                          }.bind(null, i)
                        )
                        .catch(Object(n.c)(301131))
                        .then(function (e) {
                          var i = new e(window);
                          (h = []),
                            (i.oncue = function (e) {
                              h.push(e);
                            }),
                            (i.onflush = function () {
                              delete t.xhr, r(h);
                            }),
                            i.parse(T);
                        })
                        .catch(function (e) {
                          delete t.xhr, a(Object(o.v)(null, o.b, e));
                        })
                    : ((l = Object(c.a)(T)), (h = f(l)), delete t.xhr, r(h));
                }
              } catch (e) {
                delete t.xhr, a(Object(o.v)(null, o.b, e));
              }
            })(a, e, t, r);
          },
          function (e, t, i, n) {
            r(Object(o.u)(n, o.b));
          }
        );
      }
      function h(e) {
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
      function f(e) {
        return e.map(function (e) {
          return new r.a(e.begin, e.end, e.text);
        });
      }
      i.d(t, "c", function () {
        return l;
      }),
        i.d(t, "a", function () {
          return h;
        }),
        i.d(t, "b", function () {
          return f;
        });
    },
    67: function (e, t, i) {
      "use strict";
      var r = window.VTTCue;
      function n(e) {
        if ("string" != typeof e) return !1;
        return (
          !!{ start: !0, middle: !0, end: !0, left: !0, right: !0 }[
            e.toLowerCase()
          ] && e.toLowerCase()
        );
      }
      if (!r) {
        (r = function (e, t, i) {
          var r = this;
          r.hasBeenReset = !1;
          var a = "",
            s = !1,
            c = e,
            u = t,
            o = i,
            d = null,
            l = "",
            h = !0,
            f = "auto",
            v = "start",
            T = "auto",
            g = 100,
            m = "middle";
          Object.defineProperty(r, "id", {
            enumerable: !0,
            get: function () {
              return a;
            },
            set: function (e) {
              a = "" + e;
            },
          }),
            Object.defineProperty(r, "pauseOnExit", {
              enumerable: !0,
              get: function () {
                return s;
              },
              set: function (e) {
                s = !!e;
              },
            }),
            Object.defineProperty(r, "startTime", {
              enumerable: !0,
              get: function () {
                return c;
              },
              set: function (e) {
                if ("number" != typeof e)
                  throw new TypeError("Start time must be set to a number.");
                (c = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "endTime", {
              enumerable: !0,
              get: function () {
                return u;
              },
              set: function (e) {
                if ("number" != typeof e)
                  throw new TypeError("End time must be set to a number.");
                (u = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "text", {
              enumerable: !0,
              get: function () {
                return o;
              },
              set: function (e) {
                (o = "" + e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "region", {
              enumerable: !0,
              get: function () {
                return d;
              },
              set: function (e) {
                (d = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "vertical", {
              enumerable: !0,
              get: function () {
                return l;
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
                (l = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "snapToLines", {
              enumerable: !0,
              get: function () {
                return h;
              },
              set: function (e) {
                (h = !!e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "line", {
              enumerable: !0,
              get: function () {
                return f;
              },
              set: function (e) {
                if ("number" != typeof e && "auto" !== e)
                  throw new SyntaxError(
                    "An invalid number or illegal string was specified."
                  );
                (f = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "lineAlign", {
              enumerable: !0,
              get: function () {
                return v;
              },
              set: function (e) {
                var t = n(e);
                if (!t)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (v = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "position", {
              enumerable: !0,
              get: function () {
                return T;
              },
              set: function (e) {
                if (e < 0 || e > 100)
                  throw new Error("Position must be between 0 and 100.");
                (T = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "size", {
              enumerable: !0,
              get: function () {
                return g;
              },
              set: function (e) {
                if (e < 0 || e > 100)
                  throw new Error("Size must be between 0 and 100.");
                (g = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(r, "align", {
              enumerable: !0,
              get: function () {
                return m;
              },
              set: function (e) {
                var t = n(e);
                if (!t)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (m = t), (this.hasBeenReset = !0);
              },
            }),
            (r.displayState = void 0);
        }).prototype.getCueAsHTML = function () {
          return window.WebVTT.convertCueToDOMTree(window, this.text);
        };
      }
      t.a = r;
    },
  },
]);
