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
  [3, 1, 2],
  [
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    function (e, t, n) {
      "use strict";
      n.r(t);
      var i,
        o = n(8),
        a = n(3),
        r = n(7),
        s = n(43),
        l = n(5),
        c = n(15),
        u = n(40);
      function d(e) {
        return (
          i || (i = new DOMParser()),
          Object(l.r)(
            Object(l.s)(i.parseFromString(e, "image/svg+xml").documentElement)
          )
        );
      }
      var p = function (e, t, n, i) {
          var o = document.createElement("div");
          (o.className =
            "jw-icon jw-icon-inline jw-button-color jw-reset " + e),
            o.setAttribute("role", "button"),
            o.setAttribute("tabindex", "0"),
            n && o.setAttribute("aria-label", n),
            (o.style.display = "none");
          var a = new u.a(o).on("click tap enter", t || function () {});
          return (
            i &&
              Array.prototype.forEach.call(i, function (e) {
                "string" == typeof e ? o.appendChild(d(e)) : o.appendChild(e);
              }),
            {
              ui: a,
              element: function () {
                return o;
              },
              toggle: function (e) {
                e ? this.show() : this.hide();
              },
              show: function () {
                o.style.display = "";
              },
              hide: function () {
                o.style.display = "none";
              },
            }
          );
        },
        w = n(0),
        f = n(71),
        h = n.n(f),
        j = n(72),
        g = n.n(j),
        b = n(73),
        m = n.n(b),
        v = n(74),
        y = n.n(v),
        k = n(75),
        x = n.n(k),
        O = n(76),
        C = n.n(O),
        M = n(77),
        S = n.n(M),
        T = n(78),
        E = n.n(T),
        _ = n(79),
        z = n.n(_),
        P = n(80),
        A = n.n(P),
        L = n(81),
        R = n.n(L),
        I = n(82),
        B = n.n(I),
        V = n(83),
        N = n.n(V),
        H = n(84),
        F = n.n(H),
        q = n(85),
        D = n.n(q),
        U = n(86),
        W = n.n(U),
        Q = n(62),
        Y = n.n(Q),
        X = n(87),
        Z = n.n(X),
        K = n(88),
        J = n.n(K),
        G = n(89),
        $ = n.n(G),
        ee = n(90),
        te = n.n(ee),
        ne = n(91),
        ie = n.n(ne),
        oe = n(92),
        ae = n.n(oe),
        re = n(93),
        se = n.n(re),
        le = n(94),
        ce = n.n(le),
        ue = null;
      function de(e) {
        var t = he().querySelector(we(e));
        if (t) return fe(t);
        throw new Error("Icon not found " + e);
      }
      function pe(e) {
        var t = he().querySelectorAll(e.split(",").map(we).join(","));
        if (!t.length) throw new Error("Icons not found " + e);
        return Array.prototype.map.call(t, function (e) {
          return fe(e);
        });
      }
      function we(e) {
        return ".jw-svg-icon-".concat(e);
      }
      function fe(e) {
        return e.cloneNode(!0);
      }
      function he() {
        return (
          ue ||
            (ue = d(
              "<xml>" +
                h.a +
                g.a +
                m.a +
                y.a +
                x.a +
                C.a +
                S.a +
                E.a +
                z.a +
                A.a +
                R.a +
                B.a +
                N.a +
                F.a +
                D.a +
                W.a +
                Y.a +
                Z.a +
                J.a +
                $.a +
                te.a +
                ie.a +
                ae.a +
                se.a +
                ce.a +
                "</xml>"
            )),
          ue
        );
      }
      var je = n(10);
      function ge(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var be = {};
      var me = (function () {
          function e(t, n, i, o, a) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e);
            var r,
              s = document.createElement("div");
            (s.className = "jw-icon jw-icon-inline jw-button-color jw-reset ".concat(
              a || ""
            )),
              s.setAttribute("button", o),
              s.setAttribute("role", "button"),
              s.setAttribute("tabindex", "0"),
              n && s.setAttribute("aria-label", n),
              t && "<svg" === t.substring(0, 4)
                ? (r = (function (e) {
                    if (!be[e]) {
                      var t = Object.keys(be);
                      t.length > 10 && delete be[t[0]];
                      var n = d(e);
                      be[e] = n;
                    }
                    return be[e].cloneNode(!0);
                  })(t))
                : (((r = document.createElement("div")).className =
                    "jw-icon jw-button-image jw-button-color jw-reset"),
                  t &&
                    Object(je.d)(r, {
                      backgroundImage: "url(".concat(t, ")"),
                    })),
              s.appendChild(r),
              new u.a(s).on("click tap enter", i, this),
              s.addEventListener("mousedown", function (e) {
                e.preventDefault();
              }),
              (this.id = o),
              (this.buttonElement = s);
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "element",
                value: function () {
                  return this.buttonElement;
                },
              },
              {
                key: "toggle",
                value: function (e) {
                  e ? this.show() : this.hide();
                },
              },
              {
                key: "show",
                value: function () {
                  this.buttonElement.style.display = "";
                },
              },
              {
                key: "hide",
                value: function () {
                  this.buttonElement.style.display = "none";
                },
              },
            ]) && ge(t.prototype, n),
            i && ge(t, i),
            e
          );
        })(),
        ve = n(11);
      function ye(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var ke = function (e) {
          var t = Object(l.c)(e),
            n = window.pageXOffset;
          return (
            n &&
              o.OS.android &&
              document.body.parentElement.getBoundingClientRect().left >= 0 &&
              ((t.left -= n), (t.right -= n)),
            t
          );
        },
        xe = (function () {
          function e(t, n) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              Object(w.g)(this, r.a),
              (this.className = t + " jw-background-color jw-reset"),
              (this.orientation = n);
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "setup",
                value: function () {
                  (this.el = Object(l.e)(
                    (function () {
                      var e =
                          arguments.length > 0 && void 0 !== arguments[0]
                            ? arguments[0]
                            : "",
                        t =
                          arguments.length > 1 && void 0 !== arguments[1]
                            ? arguments[1]
                            : "";
                      return (
                        '<div class="'
                          .concat(e, " ")
                          .concat(t, ' jw-reset" aria-hidden="true">') +
                        '<div class="jw-slider-container jw-reset"><div class="jw-rail jw-reset"></div><div class="jw-buffer jw-reset"></div><div class="jw-progress jw-reset"></div><div class="jw-knob jw-reset"></div></div></div>'
                      );
                    })(this.className, "jw-slider-" + this.orientation)
                  )),
                    (this.elementRail = this.el.getElementsByClassName(
                      "jw-slider-container"
                    )[0]),
                    (this.elementBuffer = this.el.getElementsByClassName(
                      "jw-buffer"
                    )[0]),
                    (this.elementProgress = this.el.getElementsByClassName(
                      "jw-progress"
                    )[0]),
                    (this.elementThumb = this.el.getElementsByClassName(
                      "jw-knob"
                    )[0]),
                    (this.ui = new u.a(this.element(), { preventScrolling: !0 })
                      .on("dragStart", this.dragStart, this)
                      .on("drag", this.dragMove, this)
                      .on("dragEnd", this.dragEnd, this)
                      .on("click tap", this.tap, this));
                },
              },
              {
                key: "dragStart",
                value: function () {
                  this.trigger("dragStart"),
                    (this.railBounds = ke(this.elementRail));
                },
              },
              {
                key: "dragEnd",
                value: function (e) {
                  this.dragMove(e), this.trigger("dragEnd");
                },
              },
              {
                key: "dragMove",
                value: function (e) {
                  var t,
                    n,
                    i = (this.railBounds = this.railBounds
                      ? this.railBounds
                      : ke(this.elementRail));
                  return (
                    (n =
                      "horizontal" === this.orientation
                        ? (t = e.pageX) < i.left
                          ? 0
                          : t > i.right
                          ? 100
                          : 100 * Object(s.a)((t - i.left) / i.width, 0, 1)
                        : (t = e.pageY) >= i.bottom
                        ? 0
                        : t <= i.top
                        ? 100
                        : 100 *
                          Object(s.a)(
                            (i.height - (t - i.top)) / i.height,
                            0,
                            1
                          )),
                    this.render(n),
                    this.update(n),
                    !1
                  );
                },
              },
              {
                key: "tap",
                value: function (e) {
                  (this.railBounds = ke(this.elementRail)), this.dragMove(e);
                },
              },
              {
                key: "limit",
                value: function (e) {
                  return e;
                },
              },
              {
                key: "update",
                value: function (e) {
                  this.trigger("update", { percentage: e });
                },
              },
              {
                key: "render",
                value: function (e) {
                  (e = Math.max(0, Math.min(e, 100))),
                    "horizontal" === this.orientation
                      ? ((this.elementThumb.style.left = e + "%"),
                        (this.elementProgress.style.width = e + "%"))
                      : ((this.elementThumb.style.bottom = e + "%"),
                        (this.elementProgress.style.height = e + "%"));
                },
              },
              {
                key: "updateBuffer",
                value: function (e) {
                  this.elementBuffer.style.width = e + "%";
                },
              },
              {
                key: "element",
                value: function () {
                  return this.el;
                },
              },
            ]) && ye(t.prototype, n),
            i && ye(t, i),
            e
          );
        })(),
        Oe = function (e, t) {
          e &&
            t &&
            (e.setAttribute("aria-label", t),
            e.setAttribute("role", "button"),
            e.setAttribute("tabindex", "0"));
        };
      function Ce(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var Me = (function () {
          function e(t, n, i, o) {
            var a = this;
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              Object(w.g)(this, r.a),
              (this.el = document.createElement("div"));
            var s =
              "jw-icon jw-icon-tooltip " + t + " jw-button-color jw-reset";
            i || (s += " jw-hidden"),
              Oe(this.el, n),
              (this.el.className = s),
              (this.tooltip = document.createElement("div")),
              (this.tooltip.className = "jw-overlay jw-reset"),
              (this.openClass = "jw-open"),
              (this.componentType = "tooltip"),
              this.el.appendChild(this.tooltip),
              o &&
                o.length > 0 &&
                Array.prototype.forEach.call(o, function (e) {
                  "string" == typeof e
                    ? a.el.appendChild(d(e))
                    : a.el.appendChild(e);
                });
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "addContent",
                value: function (e) {
                  this.content && this.removeContent(),
                    (this.content = e),
                    this.tooltip.appendChild(e);
                },
              },
              {
                key: "removeContent",
                value: function () {
                  this.content &&
                    (this.tooltip.removeChild(this.content),
                    (this.content = null));
                },
              },
              {
                key: "hasContent",
                value: function () {
                  return !!this.content;
                },
              },
              {
                key: "element",
                value: function () {
                  return this.el;
                },
              },
              {
                key: "openTooltip",
                value: function (e) {
                  this.isOpen ||
                    (this.trigger("open-" + this.componentType, e, {
                      isOpen: !0,
                    }),
                    (this.isOpen = !0),
                    Object(l.v)(this.el, this.openClass, this.isOpen));
                },
              },
              {
                key: "closeTooltip",
                value: function (e) {
                  this.isOpen &&
                    (this.trigger("close-" + this.componentType, e, {
                      isOpen: !1,
                    }),
                    (this.isOpen = !1),
                    Object(l.v)(this.el, this.openClass, this.isOpen));
                },
              },
              {
                key: "toggleOpenState",
                value: function (e) {
                  this.isOpen ? this.closeTooltip(e) : this.openTooltip(e);
                },
              },
            ]) && Ce(t.prototype, n),
            i && Ce(t, i),
            e
          );
        })(),
        Se = n(22),
        Te = n(57);
      function Ee(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var _e = (function () {
          function e(t, n) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              (this.time = t),
              (this.text = n),
              (this.el = document.createElement("div")),
              (this.el.className = "jw-cue jw-reset");
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "align",
                value: function (e) {
                  if ("%" === this.time.toString().slice(-1))
                    this.pct = this.time;
                  else {
                    var t = (this.time / e) * 100;
                    this.pct = t + "%";
                  }
                  this.el.style.left = this.pct;
                },
              },
            ]) && Ee(t.prototype, n),
            i && Ee(t, i),
            e
          );
        })(),
        ze = {
          loadChapters: function (e) {
            Object(Se.a)(
              e,
              this.chaptersLoaded.bind(this),
              this.chaptersFailed,
              { plainText: !0 }
            );
          },
          chaptersLoaded: function (e) {
            var t = Object(Te.a)(e.responseText);
            if (Array.isArray(t)) {
              var n = this._model.get("cues").concat(t);
              this._model.set("cues", n);
            }
          },
          chaptersFailed: function () {},
          addCue: function (e) {
            this.cues.push(new _e(e.begin, e.text));
          },
          drawCues: function () {
            var e = this,
              t = this._model.get("duration");
            !t ||
              t <= 0 ||
              this.cues.forEach(function (n) {
                n.align(t),
                  n.el.addEventListener("mouseover", function () {
                    e.activeCue = n;
                  }),
                  n.el.addEventListener("mouseout", function () {
                    e.activeCue = null;
                  }),
                  e.elementRail.appendChild(n.el);
              });
          },
          resetCues: function () {
            this.cues.forEach(function (e) {
              e.el.parentNode && e.el.parentNode.removeChild(e.el);
            }),
              (this.cues = []);
          },
        };
      function Pe(e) {
        (this.begin = e.begin), (this.end = e.end), (this.img = e.text);
      }
      var Ae = {
        loadThumbnails: function (e) {
          e &&
            ((this.vttPath = e.split("?")[0].split("/").slice(0, -1).join("/")),
            (this.individualImage = null),
            Object(Se.a)(
              e,
              this.thumbnailsLoaded.bind(this),
              this.thumbnailsFailed.bind(this),
              { plainText: !0 }
            ));
        },
        thumbnailsLoaded: function (e) {
          var t = Object(Te.a)(e.responseText);
          Array.isArray(t) &&
            (t.forEach(function (e) {
              this.thumbnails.push(new Pe(e));
            }, this),
            this.drawCues());
        },
        thumbnailsFailed: function () {},
        chooseThumbnail: function (e) {
          var t = Object(w.A)(this.thumbnails, { end: e }, Object(w.z)("end"));
          t >= this.thumbnails.length && (t = this.thumbnails.length - 1);
          var n = this.thumbnails[t].img;
          return (
            n.indexOf("://") < 0 &&
              (n = this.vttPath ? this.vttPath + "/" + n : n),
            n
          );
        },
        loadThumbnail: function (e) {
          var t = this.chooseThumbnail(e),
            n = { margin: "0 auto", backgroundPosition: "0 0" };
          if (t.indexOf("#xywh") > 0)
            try {
              var i = /(.+)#xywh=(\d+),(\d+),(\d+),(\d+)/.exec(t);
              (t = i[1]),
                (n.backgroundPosition = -1 * i[2] + "px " + -1 * i[3] + "px"),
                (n.width = i[4]),
                this.timeTip.setWidth(+n.width),
                (n.height = i[5]);
            } catch (e) {
              return;
            }
          else
            this.individualImage ||
              ((this.individualImage = new Image()),
              (this.individualImage.onload = Object(w.a)(function () {
                (this.individualImage.onload = null),
                  this.timeTip.image({
                    width: this.individualImage.width,
                    height: this.individualImage.height,
                  }),
                  this.timeTip.setWidth(this.individualImage.width);
              }, this)),
              (this.individualImage.src = t));
          return (n.backgroundImage = 'url("' + t + '")'), n;
        },
        showThumbnail: function (e) {
          this._model.get("containerWidth") <= 420 ||
            this.thumbnails.length < 1 ||
            this.timeTip.image(this.loadThumbnail(e));
        },
        resetThumbnails: function () {
          this.timeTip.image({ backgroundImage: "", width: 0, height: 0 }),
            (this.thumbnails = []);
        },
      };
      function Le(e, t, n) {
        return (Le =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (e, t, n) {
                var i = (function (e, t) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(e, t) &&
                    null !== (e = He(e));

                  );
                  return e;
                })(e, t);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, t);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(e, t, n || e);
      }
      function Re(e) {
        return (Re =
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
      function Ie(e, t) {
        if (!(e instanceof t))
          throw new TypeError("Cannot call a class as a function");
      }
      function Be(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function Ve(e, t, n) {
        return t && Be(e.prototype, t), n && Be(e, n), e;
      }
      function Ne(e, t) {
        return !t || ("object" !== Re(t) && "function" != typeof t)
          ? (function (e) {
              if (void 0 === e)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return e;
            })(e)
          : t;
      }
      function He(e) {
        return (He = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function Fe(e, t) {
        if ("function" != typeof t && null !== t)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (e.prototype = Object.create(t && t.prototype, {
          constructor: { value: e, writable: !0, configurable: !0 },
        })),
          t && qe(e, t);
      }
      function qe(e, t) {
        return (qe =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var De = (function (e) {
        function t() {
          return Ie(this, t), Ne(this, He(t).apply(this, arguments));
        }
        return (
          Fe(t, e),
          Ve(t, [
            {
              key: "setup",
              value: function () {
                (this.text = document.createElement("span")),
                  (this.text.className = "jw-text jw-reset"),
                  (this.img = document.createElement("div")),
                  (this.img.className = "jw-time-thumb jw-reset"),
                  (this.containerWidth = 0),
                  (this.textLength = 0),
                  (this.dragJustReleased = !1);
                var e = document.createElement("div");
                (e.className = "jw-time-tip jw-reset"),
                  e.appendChild(this.img),
                  e.appendChild(this.text),
                  this.addContent(e);
              },
            },
            {
              key: "image",
              value: function (e) {
                Object(je.d)(this.img, e);
              },
            },
            {
              key: "update",
              value: function (e) {
                this.text.textContent = e;
              },
            },
            {
              key: "getWidth",
              value: function () {
                return (
                  this.containerWidth || this.setWidth(), this.containerWidth
                );
              },
            },
            {
              key: "setWidth",
              value: function (e) {
                e
                  ? (this.containerWidth = e + 16)
                  : this.tooltip &&
                    (this.containerWidth =
                      Object(l.c)(this.container).width + 16);
              },
            },
            {
              key: "resetWidth",
              value: function () {
                this.containerWidth = 0;
              },
            },
          ]),
          t
        );
      })(Me);
      var Ue = (function (e) {
        function t(e, n, i) {
          var o;
          return (
            Ie(this, t),
            ((o = Ne(
              this,
              He(t).call(this, "jw-slider-time", "horizontal")
            ))._model = e),
            (o._api = n),
            (o.timeUpdateKeeper = i),
            (o.timeTip = new De("jw-tooltip-time", null, !0)),
            o.timeTip.setup(),
            (o.cues = []),
            (o.seekThrottled = Object(w.B)(o.performSeek, 400)),
            (o.mobileHoverDistance = 5),
            o.setup(),
            o
          );
        }
        return (
          Fe(t, e),
          Ve(t, [
            {
              key: "setup",
              value: function () {
                var e = this;
                Le(He(t.prototype), "setup", this).apply(this, arguments),
                  this._model
                    .on("change:duration", this.onDuration, this)
                    .on("change:cues", this.updateCues, this)
                    .on("seeked", function () {
                      e._model.get("scrubbing") || e.updateAriaText();
                    })
                    .change("position", this.onPosition, this)
                    .change("buffer", this.onBuffer, this)
                    .change("streamType", this.onStreamType, this),
                  this._model.player.change(
                    "playlistItem",
                    this.onPlaylistItem,
                    this
                  );
                var n = this.el;
                Object(l.t)(n, "tabindex", "0"),
                  Object(l.t)(n, "role", "slider"),
                  Object(l.t)(
                    n,
                    "aria-label",
                    this._model.get("localization").slider
                  ),
                  n.removeAttribute("aria-hidden"),
                  this.elementRail.appendChild(this.timeTip.element()),
                  (this.ui = (this.ui || new u.a(n))
                    .on("move drag", this.showTimeTooltip, this)
                    .on("dragEnd out", this.hideTimeTooltip, this)
                    .on("click", function () {
                      return n.focus();
                    })
                    .on("focus", this.updateAriaText, this));
              },
            },
            {
              key: "update",
              value: function (e) {
                (this.seekTo = e),
                  this.seekThrottled(),
                  Le(He(t.prototype), "update", this).apply(this, arguments);
              },
            },
            {
              key: "dragStart",
              value: function () {
                this._model.set("scrubbing", !0),
                  Le(He(t.prototype), "dragStart", this).apply(this, arguments);
              },
            },
            {
              key: "dragEnd",
              value: function () {
                Le(He(t.prototype), "dragEnd", this).apply(this, arguments),
                  this._model.set("scrubbing", !1);
              },
            },
            {
              key: "onBuffer",
              value: function (e, t) {
                this.updateBuffer(t);
              },
            },
            {
              key: "onPosition",
              value: function (e, t) {
                this.updateTime(t, e.get("duration"));
              },
            },
            {
              key: "onDuration",
              value: function (e, t) {
                this.updateTime(e.get("position"), t),
                  Object(l.t)(this.el, "aria-valuemin", 0),
                  Object(l.t)(this.el, "aria-valuemax", t),
                  this.drawCues();
              },
            },
            {
              key: "onStreamType",
              value: function (e, t) {
                this.streamType = t;
              },
            },
            {
              key: "updateTime",
              value: function (e, t) {
                var n = 0;
                if (t)
                  if ("DVR" === this.streamType) {
                    var i = this._model.get("dvrSeekLimit"),
                      o = t + i;
                    n = ((o - (e + i)) / o) * 100;
                  } else
                    ("VOD" !== this.streamType && this.streamType) ||
                      (n = (e / t) * 100);
                this.render(n);
              },
            },
            {
              key: "onPlaylistItem",
              value: function (e, t) {
                this.reset();
                var n = e.get("cues");
                !this.cues.length && n.length && this.updateCues(null, n);
                var i = t.tracks;
                Object(w.f)(
                  i,
                  function (e) {
                    e && e.kind && "thumbnails" === e.kind.toLowerCase()
                      ? this.loadThumbnails(e.file)
                      : e &&
                        e.kind &&
                        "chapters" === e.kind.toLowerCase() &&
                        this.loadChapters(e.file);
                  },
                  this
                );
              },
            },
            {
              key: "performSeek",
              value: function () {
                var e,
                  t = this.seekTo,
                  n = this._model.get("duration");
                if (0 === n) this._api.play({ reason: "interaction" });
                else if ("DVR" === this.streamType) {
                  var i = this._model.get("seekRange") || { start: 0 },
                    o = this._model.get("dvrSeekLimit");
                  (e = i.start + ((-n - o) * t) / 100),
                    this._api.seek(e, { reason: "interaction" });
                } else
                  (e = (t / 100) * n),
                    this._api.seek(Math.min(e, n - 0.25), {
                      reason: "interaction",
                    });
              },
            },
            {
              key: "showTimeTooltip",
              value: function (e) {
                var t = this,
                  n = this._model.get("duration");
                if (0 !== n) {
                  var i,
                    o = this._model.get("containerWidth"),
                    a = Object(l.c)(this.elementRail),
                    r = e.pageX ? e.pageX - a.left : e.x,
                    c = (r = Object(s.a)(r, 0, a.width)) / a.width,
                    u = n * c;
                  if (n < 0)
                    u = (n += this._model.get("dvrSeekLimit")) - (u = n * c);
                  if (
                    ("touch" === e.pointerType &&
                      (this.activeCue = this.cues.reduce(function (e, n) {
                        return Math.abs(r - (parseInt(n.pct) / 100) * a.width) <
                          t.mobileHoverDistance
                          ? n
                          : e;
                      }, void 0)),
                    this.activeCue)
                  )
                    i = this.activeCue.text;
                  else {
                    (i = Object(ve.timeFormat)(u, !0)),
                      n < 0 && u > -1 && (i = "Live");
                  }
                  var d = this.timeTip;
                  d.update(i),
                    this.textLength !== i.length &&
                      ((this.textLength = i.length), d.resetWidth()),
                    this.showThumbnail(u),
                    Object(l.a)(d.el, "jw-open");
                  var p = d.getWidth(),
                    w = a.width / 100,
                    f = o - a.width,
                    h = 0;
                  p > f && (h = (p - f) / (200 * w));
                  var j = 100 * Math.min(1 - h, Math.max(h, c)).toFixed(3);
                  Object(je.d)(d.el, { left: j + "%" });
                }
              },
            },
            {
              key: "hideTimeTooltip",
              value: function () {
                Object(l.o)(this.timeTip.el, "jw-open");
              },
            },
            {
              key: "updateCues",
              value: function (e, t) {
                var n = this;
                this.resetCues(),
                  t &&
                    t.length &&
                    (t.forEach(function (e) {
                      n.addCue(e);
                    }),
                    this.drawCues());
              },
            },
            {
              key: "updateAriaText",
              value: function () {
                var e = this._model;
                if (!e.get("seeking")) {
                  var t = e.get("position"),
                    n = e.get("duration"),
                    i = Object(ve.timeFormat)(t);
                  "DVR" !== this.streamType &&
                    (i += " of ".concat(Object(ve.timeFormat)(n)));
                  var o = this.el;
                  document.activeElement !== o &&
                    (this.timeUpdateKeeper.textContent = i),
                    Object(l.t)(o, "aria-valuenow", t),
                    Object(l.t)(o, "aria-valuetext", i);
                }
              },
            },
            {
              key: "reset",
              value: function () {
                this.resetThumbnails(),
                  this.timeTip.resetWidth(),
                  (this.textLength = 0);
              },
            },
          ]),
          t
        );
      })(xe);
      Object(w.g)(Ue.prototype, ze, Ae);
      var We = Ue;
      function Qe(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function Ye(e, t, n) {
        return (Ye =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (e, t, n) {
                var i = (function (e, t) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(e, t) &&
                    null !== (e = Ge(e));

                  );
                  return e;
                })(e, t);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, t);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(e, t, n || e);
      }
      function Xe(e) {
        return (Xe =
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
      function Ze(e, t) {
        if (!(e instanceof t))
          throw new TypeError("Cannot call a class as a function");
      }
      function Ke(e, t) {
        return !t || ("object" !== Xe(t) && "function" != typeof t) ? Je(e) : t;
      }
      function Je(e) {
        if (void 0 === e)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return e;
      }
      function Ge(e) {
        return (Ge = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function $e(e, t) {
        if ("function" != typeof t && null !== t)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (e.prototype = Object.create(t && t.prototype, {
          constructor: { value: e, writable: !0, configurable: !0 },
        })),
          t && et(e, t);
      }
      function et(e, t) {
        return (et =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var tt = (function (e) {
          function t(e, n, i) {
            var o;
            Ze(this, t);
            var a = "jw-slider-volume";
            return (
              "vertical" === e && (a += " jw-volume-tip"),
              (o = Ke(this, Ge(t).call(this, a, e))).setup(),
              o.element().classList.remove("jw-background-color"),
              Object(l.t)(i, "tabindex", "0"),
              Object(l.t)(i, "aria-label", n),
              Object(l.t)(i, "aria-orientation", e),
              Object(l.t)(i, "aria-valuemin", 0),
              Object(l.t)(i, "aria-valuemax", 100),
              Object(l.t)(i, "role", "slider"),
              (o.uiOver = new u.a(i).on("click", function () {})),
              o
            );
          }
          return $e(t, e), t;
        })(xe),
        nt = (function (e) {
          function t(e, n, i, o, a) {
            var r;
            Ze(this, t),
              ((r = Ke(this, Ge(t).call(this, n, i, !0, o)))._model = e),
              (r.horizontalContainer = a);
            var s = e.get("localization").volumeSlider;
            return (
              (r.horizontalSlider = new tt("horizontal", s, a, Je(Je(r)))),
              (r.verticalSlider = new tt("vertical", s, r.tooltip, Je(Je(r)))),
              a.appendChild(r.horizontalSlider.element()),
              r.addContent(r.verticalSlider.element()),
              r.verticalSlider.on(
                "update",
                function (e) {
                  this.trigger("update", e);
                },
                Je(Je(r))
              ),
              r.horizontalSlider.on(
                "update",
                function (e) {
                  this.trigger("update", e);
                },
                Je(Je(r))
              ),
              r.horizontalSlider.uiOver.on("keydown", function (e) {
                var t = e.sourceEvent;
                switch (t.keyCode) {
                  case 37:
                    t.stopPropagation(), r.trigger("adjustVolume", -10);
                    break;
                  case 39:
                    t.stopPropagation(), r.trigger("adjustVolume", 10);
                }
              }),
              (r.ui = new u.a(r.el, { directSelect: !0 })
                .on("click enter", r.toggleValue, Je(Je(r)))
                .on("tap", r.toggleOpenState, Je(Je(r)))),
              r.addSliderHandlers(r.ui),
              r.addSliderHandlers(r.horizontalSlider.uiOver),
              r.addSliderHandlers(r.verticalSlider.uiOver),
              r.onAudioMode(null, e.get("audioMode")),
              r._model.on("change:audioMode", r.onAudioMode, Je(Je(r))),
              r._model.on("change:volume", r.onVolume, Je(Je(r))),
              r
            );
          }
          var n, i, o;
          return (
            $e(t, e),
            (n = t),
            (i = [
              {
                key: "onAudioMode",
                value: function (e, t) {
                  var n = t ? 0 : -1;
                  Object(l.t)(this.horizontalContainer, "tabindex", n);
                },
              },
              {
                key: "addSliderHandlers",
                value: function (e) {
                  var t = this.openSlider,
                    n = this.closeSlider;
                  e.on("over", t, this)
                    .on("out", n, this)
                    .on("focus", t, this)
                    .on("blur", n, this);
                },
              },
              {
                key: "openSlider",
                value: function (e) {
                  Ye(Ge(t.prototype), "openTooltip", this).call(this, e),
                    Object(l.v)(this.horizontalContainer, this.openClass, !0);
                },
              },
              {
                key: "closeSlider",
                value: function (e) {
                  Ye(Ge(t.prototype), "closeTooltip", this).call(this, e),
                    Object(l.v)(this.horizontalContainer, this.openClass, !1),
                    this.horizontalContainer.blur();
                },
              },
              {
                key: "toggleValue",
                value: function () {
                  this.trigger("toggleValue");
                },
              },
              {
                key: "destroy",
                value: function () {
                  this.horizontalSlider.uiOver.destroy(),
                    this.verticalSlider.uiOver.destroy(),
                    this.ui.destroy();
                },
              },
            ]) && Qe(n.prototype, i),
            o && Qe(n, o),
            t
          );
        })(Me);
      function it(e, t, n, i, o) {
        var a = document.createElement("div");
        (a.className = "jw-reset-text jw-tooltip jw-tooltip-".concat(t)),
          a.setAttribute("dir", "auto");
        var r = document.createElement("div");
        (r.className = "jw-text"), a.appendChild(r), e.appendChild(a);
        var s = {
            dirty: !!n,
            opened: !1,
            text: n,
            open: function () {
              s.touchEvent ||
                (s.suppress ? (s.suppress = !1) : (c(!0), i && i()));
            },
            close: function () {
              s.touchEvent || (c(!1), o && o());
            },
            setText: function (e) {
              e !== s.text && ((s.text = e), (s.dirty = !0)), s.opened && c(!0);
            },
          },
          c = function (e) {
            e && s.dirty && (Object(l.q)(r, s.text), (s.dirty = !1)),
              (s.opened = e),
              Object(l.v)(a, "jw-open", e);
          };
        return (
          e.addEventListener("mouseover", s.open),
          e.addEventListener("focus", s.open),
          e.addEventListener("blur", s.close),
          e.addEventListener("mouseout", s.close),
          e.addEventListener(
            "touchstart",
            function () {
              s.touchEvent = !0;
            },
            { passive: !0 }
          ),
          s
        );
      }
      var ot = n(47);
      function at(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function rt(e, t) {
        var n = document.createElement("div");
        return (
          (n.className = "jw-icon jw-icon-inline jw-text jw-reset " + e),
          t && Object(l.t)(n, "role", t),
          n
        );
      }
      function st(e) {
        var t = document.createElement("div");
        return (t.className = "jw-reset ".concat(e)), t;
      }
      function lt(e, t) {
        if (o.Browser.safari) {
          var n = p(
            "jw-icon-airplay jw-off",
            e,
            t.airplay,
            pe("airplay-off,airplay-on")
          );
          return it(n.element(), "airplay", t.airplay), n;
        }
        if (o.Browser.chrome && window.chrome) {
          var i = document.createElement("google-cast-launcher");
          Object(l.t)(i, "tabindex", "-1"), (i.className += " jw-reset");
          var a = p("jw-icon-cast", null, t.cast);
          a.ui.off();
          var r = a.element();
          return (
            (r.style.cursor = "pointer"),
            r.appendChild(i),
            (a.button = i),
            it(r, "chromecast", t.cast),
            a
          );
        }
      }
      function ct(e, t) {
        return e.filter(function (e) {
          return !t.some(function (t) {
            return (
              t.id + t.btnClass === e.id + e.btnClass &&
              e.callback === t.callback
            );
          });
        });
      }
      var ut = function (e, t) {
          t.forEach(function (t) {
            t.element && (t = t.element()), e.appendChild(t);
          });
        },
        dt = (function () {
          function e(t, n, i) {
            var s = this;
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              Object(w.g)(this, r.a),
              (this._api = t),
              (this._model = n),
              (this._isMobile = o.OS.mobile),
              (this._volumeAnnouncer = i.querySelector(".jw-volume-update"));
            var c,
              d,
              f,
              h = n.get("localization"),
              j = new We(n, t, i.querySelector(".jw-time-update")),
              g = (this.menus = []);
            this.ui = [];
            var b = "",
              m = h.volume;
            if (this._isMobile) {
              if (
                !(n.get("sdkplatform") || (o.OS.iOS && o.OS.version.major < 10))
              ) {
                var v = pe("volume-0,volume-100");
                f = p(
                  "jw-icon-volume",
                  function () {
                    t.setMute();
                  },
                  m,
                  v
                );
              }
            } else {
              (d = document.createElement("div")).className =
                "jw-horizontal-volume-container";
              var y = (c = new nt(
                n,
                "jw-icon-volume",
                m,
                pe("volume-0,volume-50,volume-100"),
                d
              )).element();
              g.push(c),
                Object(l.t)(y, "role", "button"),
                n.change(
                  "mute",
                  function (e, t) {
                    var n = t ? h.unmute : h.mute;
                    Object(l.t)(y, "aria-label", n);
                  },
                  this
                );
            }
            var k = p(
                "jw-icon-next",
                function () {
                  t.next({ feedShownId: b, reason: "interaction" });
                },
                h.next,
                pe("next")
              ),
              x = p(
                "jw-icon-settings jw-settings-submenu-button",
                function (e) {
                  s.trigger("settingsInteraction", "quality", !0, e);
                },
                h.settings,
                pe("settings")
              );
            Object(l.t)(x.element(), "aria-haspopup", "true");
            var O = p(
              "jw-icon-cc jw-settings-submenu-button",
              function (e) {
                s.trigger("settingsInteraction", "captions", !1, e);
              },
              h.cc,
              pe("cc-off,cc-on")
            );
            Object(l.t)(O.element(), "aria-haspopup", "true");
            var C = p(
              "jw-text-live",
              function () {
                s.goToLiveEdge();
              },
              h.liveBroadcast
            );
            C.element().textContent = h.liveBroadcast;
            var M,
              S,
              T,
              E = (this.elements = {
                alt:
                  ((M = "jw-text-alt"),
                  (S = "status"),
                  (T = document.createElement("span")),
                  (T.className = "jw-text jw-reset " + M),
                  S && Object(l.t)(T, "role", S),
                  T),
                play: p(
                  "jw-icon-playback",
                  function () {
                    t.playToggle({ reason: "interaction" });
                  },
                  h.play,
                  pe("play,pause,stop")
                ),
                rewind: p(
                  "jw-icon-rewind",
                  function () {
                    s.rewind();
                  },
                  h.rewind,
                  pe("rewind")
                ),
                live: C,
                next: k,
                elapsed: rt("jw-text-elapsed", "timer"),
                countdown: rt("jw-text-countdown", "timer"),
                time: j,
                duration: rt("jw-text-duration", "timer"),
                mute: f,
                volumetooltip: c,
                horizontalVolumeContainer: d,
                cast: lt(function () {
                  t.castToggle();
                }, h),
                fullscreen: p(
                  "jw-icon-fullscreen",
                  function () {
                    t.setFullscreen();
                  },
                  h.fullscreen,
                  pe("fullscreen-off,fullscreen-on")
                ),
                spacer: st("jw-spacer"),
                buttonContainer: st("jw-button-container"),
                settingsButton: x,
                captionsButton: O,
              }),
              _ = it(O.element(), "captions", h.cc),
              z = function (e) {
                var t = e.get("captionsList")[e.get("captionsIndex")],
                  n = h.cc;
                t && "Off" !== t.label && (n = t.label), _.setText(n);
              },
              P = it(E.play.element(), "play", h.play);
            this.setPlayText = function (e) {
              P.setText(e);
            };
            var A = E.next.element(),
              L = it(
                A,
                "next",
                h.nextUp,
                function () {
                  var e = n.get("nextUp");
                  (b = Object(ot.b)(ot.a)),
                    s.trigger("nextShown", {
                      mode: e.mode,
                      ui: "nextup",
                      itemsShown: [e],
                      feedData: e.feedData,
                      reason: "hover",
                      feedShownId: b,
                    });
                },
                function () {
                  b = "";
                }
              );
            Object(l.t)(A, "dir", "auto"),
              it(E.rewind.element(), "rewind", h.rewind),
              it(E.settingsButton.element(), "settings", h.settings);
            var R = it(E.fullscreen.element(), "fullscreen", h.fullscreen),
              I = [
                E.play,
                E.rewind,
                E.next,
                E.volumetooltip,
                E.mute,
                E.horizontalVolumeContainer,
                E.alt,
                E.live,
                E.elapsed,
                E.countdown,
                E.duration,
                E.spacer,
                E.cast,
                E.captionsButton,
                E.settingsButton,
                E.fullscreen,
              ].filter(function (e) {
                return e;
              }),
              B = [E.time, E.buttonContainer].filter(function (e) {
                return e;
              });
            (this.el = document.createElement("div")),
              (this.el.className = "jw-controlbar jw-reset"),
              ut(E.buttonContainer, I),
              ut(this.el, B);
            var V = n.get("logo");
            if (
              (V && "control-bar" === V.position && this.addLogo(V),
              E.play.show(),
              E.fullscreen.show(),
              E.mute && E.mute.show(),
              n.change("volume", this.onVolume, this),
              n.change(
                "mute",
                function (e, t) {
                  s.renderVolume(t, e.get("volume"));
                },
                this
              ),
              n.change("state", this.onState, this),
              n.change("duration", this.onDuration, this),
              n.change("position", this.onElapsed, this),
              n.change(
                "fullscreen",
                function (e, t) {
                  var n = s.elements.fullscreen.element();
                  Object(l.v)(n, "jw-off", t);
                  var i = e.get("fullscreen") ? h.exitFullscreen : h.fullscreen;
                  R.setText(i), Object(l.t)(n, "aria-label", i);
                },
                this
              ),
              n.change("streamType", this.onStreamTypeChange, this),
              n.change(
                "dvrLive",
                function (e, t) {
                  var n = h.liveBroadcast,
                    i = h.notLive,
                    o = s.elements.live.element(),
                    a = !1 === t;
                  Object(l.v)(o, "jw-dvr-live", a),
                    Object(l.t)(o, "aria-label", a ? i : n),
                    (o.textContent = n);
                },
                this
              ),
              n.change("altText", this.setAltText, this),
              n.change("customButtons", this.updateButtons, this),
              n.on("change:captionsIndex", z, this),
              n.on("change:captionsList", z, this),
              n.change(
                "nextUp",
                function (e, t) {
                  b = Object(ot.b)(ot.a);
                  var n = h.nextUp;
                  t && t.title && (n += ": ".concat(t.title)),
                    L.setText(n),
                    E.next.toggle(!!t);
                },
                this
              ),
              n.change("audioMode", this.onAudioMode, this),
              E.cast &&
                (n.change("castAvailable", this.onCastAvailable, this),
                n.change("castActive", this.onCastActive, this)),
              E.volumetooltip &&
                (E.volumetooltip.on(
                  "update",
                  function (e) {
                    var t = e.percentage;
                    this._api.setVolume(t);
                  },
                  this
                ),
                E.volumetooltip.on(
                  "toggleValue",
                  function () {
                    this._api.setMute();
                  },
                  this
                ),
                E.volumetooltip.on(
                  "adjustVolume",
                  function (e) {
                    this.trigger("adjustVolume", e);
                  },
                  this
                )),
              E.cast && E.cast.button)
            ) {
              var N = E.cast.ui.on(
                "click tap enter",
                function (e) {
                  "click" !== e.type && E.cast.button.click(),
                    this._model.set("castClicked", !0);
                },
                this
              );
              this.ui.push(N);
            }
            var H = new u.a(E.duration).on(
              "click tap enter",
              function () {
                if ("DVR" === this._model.get("streamType")) {
                  var e = this._model.get("position"),
                    t = this._model.get("dvrSeekLimit");
                  this._api.seek(Math.max(-t, e), { reason: "interaction" });
                }
              },
              this
            );
            this.ui.push(H);
            var F = new u.a(this.el).on(
              "click tap drag",
              function () {
                this.trigger(a.sb);
              },
              this
            );
            this.ui.push(F),
              g.forEach(function (e) {
                e.on("open-tooltip", s.closeMenus, s);
              });
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "onVolume",
                value: function (e, t) {
                  this.renderVolume(e.get("mute"), t);
                },
              },
              {
                key: "renderVolume",
                value: function (e, t) {
                  var n = this.elements.mute,
                    i = this.elements.volumetooltip;
                  if (
                    (n &&
                      (Object(l.v)(n.element(), "jw-off", e),
                      Object(l.v)(n.element(), "jw-full", !e)),
                    i)
                  ) {
                    var o = e ? 0 : t,
                      a = i.element();
                    i.verticalSlider.render(o), i.horizontalSlider.render(o);
                    var r = i.tooltip,
                      s = i.horizontalContainer;
                    Object(l.v)(a, "jw-off", e),
                      Object(l.v)(a, "jw-full", t >= 75 && !e),
                      Object(l.t)(r, "aria-valuenow", o),
                      Object(l.t)(s, "aria-valuenow", o);
                    var c = "Volume ".concat(o, "%");
                    Object(l.t)(r, "aria-valuetext", c),
                      Object(l.t)(s, "aria-valuetext", c),
                      document.activeElement !== r &&
                        document.activeElement !== s &&
                        (this._volumeAnnouncer.textContent = c);
                  }
                },
              },
              {
                key: "onCastAvailable",
                value: function (e, t) {
                  this.elements.cast.toggle(t);
                },
              },
              {
                key: "onCastActive",
                value: function (e, t) {
                  this.elements.fullscreen.toggle(!t),
                    this.elements.cast.button &&
                      Object(l.v)(this.elements.cast.button, "jw-off", !t);
                },
              },
              {
                key: "onElapsed",
                value: function (e, t) {
                  var n,
                    i,
                    o = e.get("duration");
                  if ("DVR" === e.get("streamType")) {
                    var a = Math.ceil(t),
                      r = this._model.get("dvrSeekLimit");
                    (n = i =
                      a >= -r ? "" : "-" + Object(ve.timeFormat)(-(t + r))),
                      e.set("dvrLive", a >= -r);
                  } else
                    (n = Object(ve.timeFormat)(t)),
                      (i = Object(ve.timeFormat)(o - t));
                  (this.elements.elapsed.textContent = n),
                    (this.elements.countdown.textContent = i);
                },
              },
              {
                key: "onDuration",
                value: function (e, t) {
                  this.elements.duration.textContent = Object(ve.timeFormat)(
                    Math.abs(t)
                  );
                },
              },
              {
                key: "onAudioMode",
                value: function (e, t) {
                  var n = this.elements.time.element();
                  t
                    ? this.elements.buttonContainer.insertBefore(
                        n,
                        this.elements.elapsed
                      )
                    : Object(l.m)(this.el, n);
                },
              },
              {
                key: "element",
                value: function () {
                  return this.el;
                },
              },
              {
                key: "setAltText",
                value: function (e, t) {
                  this.elements.alt.textContent = t;
                },
              },
              {
                key: "closeMenus",
                value: function (e) {
                  this.menus.forEach(function (t) {
                    (e && e.target === t.el) || t.closeTooltip(e);
                  });
                },
              },
              {
                key: "rewind",
                value: function () {
                  var e,
                    t = 0,
                    n = this._model.get("currentTime");
                  n
                    ? (e = n - 10)
                    : ((e = this._model.get("position") - 10),
                      "DVR" === this._model.get("streamType") &&
                        (t = this._model.get("duration"))),
                    this._api.seek(Math.max(e, t), { reason: "interaction" });
                },
              },
              {
                key: "onState",
                value: function (e, t) {
                  var n = e.get("localization"),
                    i = n.play;
                  this.setPlayText(i),
                    t === a.pb &&
                      ("LIVE" !== e.get("streamType")
                        ? ((i = n.pause), this.setPlayText(i))
                        : ((i = n.stop), this.setPlayText(i))),
                    Object(l.t)(this.elements.play.element(), "aria-label", i);
                },
              },
              {
                key: "onStreamTypeChange",
                value: function (e, t) {
                  var n = "LIVE" === t,
                    i = "DVR" === t;
                  this.elements.rewind.toggle(!n),
                    this.elements.live.toggle(n || i),
                    Object(l.t)(
                      this.elements.live.element(),
                      "tabindex",
                      n ? "-1" : "0"
                    ),
                    (this.elements.duration.style.display = i ? "none" : ""),
                    this.onDuration(e, e.get("duration")),
                    this.onState(e, e.get("state"));
                },
              },
              {
                key: "addLogo",
                value: function (e) {
                  var t = this.elements.buttonContainer,
                    n = new me(
                      e.file,
                      this._model.get("localization").logo,
                      function () {
                        e.link &&
                          Object(l.l)(e.link, "_blank", { rel: "noreferrer" });
                      },
                      "logo",
                      "jw-logo-button"
                    );
                  e.link || Object(l.t)(n.element(), "tabindex", "-1"),
                    t.insertBefore(
                      n.element(),
                      t.querySelector(".jw-spacer").nextSibling
                    );
                },
              },
              {
                key: "goToLiveEdge",
                value: function () {
                  if ("DVR" === this._model.get("streamType")) {
                    var e = Math.min(this._model.get("position"), -1),
                      t = this._model.get("dvrSeekLimit");
                    this._api.seek(Math.max(-t, e), { reason: "interaction" }),
                      this._api.play({ reason: "interaction" });
                  }
                },
              },
              {
                key: "updateButtons",
                value: function (e, t, n) {
                  if (t) {
                    var i,
                      o,
                      a = this.elements.buttonContainer;
                    t !== n && n
                      ? ((i = ct(t, n)),
                        (o = ct(n, t)),
                        this.removeButtons(a, o))
                      : (i = t);
                    for (var r = i.length - 1; r >= 0; r--) {
                      var s = i[r],
                        l = new me(
                          s.img,
                          s.tooltip,
                          s.callback,
                          s.id,
                          s.btnClass
                        );
                      s.tooltip && it(l.element(), s.id, s.tooltip);
                      var c = void 0;
                      "related" === l.id
                        ? (c = this.elements.settingsButton.element())
                        : "share" === l.id
                        ? (c =
                            a.querySelector('[button="related"]') ||
                            this.elements.settingsButton.element())
                        : (c = this.elements.spacer.nextSibling) &&
                          "logo" === c.getAttribute("button") &&
                          (c = c.nextSibling),
                        a.insertBefore(l.element(), c);
                    }
                  }
                },
              },
              {
                key: "removeButtons",
                value: function (e, t) {
                  for (var n = t.length; n--; ) {
                    var i = e.querySelector('[button="'.concat(t[n].id, '"]'));
                    i && e.removeChild(i);
                  }
                },
              },
              {
                key: "toggleCaptionsButtonState",
                value: function (e) {
                  var t = this.elements.captionsButton;
                  t && Object(l.v)(t.element(), "jw-off", !e);
                },
              },
              {
                key: "destroy",
                value: function () {
                  var e = this;
                  this._model.off(null, null, this),
                    Object.keys(this.elements).forEach(function (t) {
                      var n = e.elements[t];
                      n &&
                        "function" == typeof n.destroy &&
                        e.elements[t].destroy();
                    }),
                    this.ui.forEach(function (e) {
                      e.destroy();
                    }),
                    (this.ui = []);
                },
              },
            ]) && at(t.prototype, n),
            i && at(t, i),
            e
          );
        })(),
        pt = function () {
          var e =
              arguments.length > 0 && void 0 !== arguments[0]
                ? arguments[0]
                : "",
            t =
              arguments.length > 1 && void 0 !== arguments[1]
                ? arguments[1]
                : "";
          return (
            '<div class="jw-display-icon-container jw-display-icon-'.concat(
              e,
              ' jw-reset">'
            ) +
            '<div class="jw-icon jw-icon-'
              .concat(
                e,
                ' jw-button-color jw-reset" role="button" tabindex="0" aria-label="'
              )
              .concat(t, '"></div>') +
            "</div>"
          );
        },
        wt = function (e) {
          return (
            '<div class="jw-display jw-reset"><div class="jw-display-container jw-reset"><div class="jw-display-controls jw-reset">' +
            pt("rewind", e.rewind) +
            pt("display", e.playback) +
            pt("next", e.next) +
            "</div></div></div>"
          );
        };
      function ft(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var ht = (function () {
        function e(t, n, i) {
          !(function (e, t) {
            if (!(e instanceof t))
              throw new TypeError("Cannot call a class as a function");
          })(this, e);
          var o = i.querySelector(".jw-icon");
          (this.el = i),
            (this.ui = new u.a(o).on("click tap enter", function () {
              var e = t.get("position"),
                i = t.get("duration"),
                o = e - 10,
                a = 0;
              "DVR" === t.get("streamType") && (a = i), n.seek(Math.max(o, a));
            }));
        }
        var t, n, i;
        return (
          (t = e),
          (n = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && ft(t.prototype, n),
          i && ft(t, i),
          e
        );
      })();
      function jt(e) {
        return (jt =
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
      function gt(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function bt(e, t) {
        return !t || ("object" !== jt(t) && "function" != typeof t)
          ? (function (e) {
              if (void 0 === e)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return e;
            })(e)
          : t;
      }
      function mt(e) {
        return (mt = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function vt(e, t) {
        return (vt =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var yt = (function (e) {
        function t(e, n, i) {
          var o;
          !(function (e, t) {
            if (!(e instanceof t))
              throw new TypeError("Cannot call a class as a function");
          })(this, t),
            (o = bt(this, mt(t).call(this)));
          var a = e.get("localization"),
            r = i.querySelector(".jw-icon");
          if (
            ((o.icon = r),
            (o.el = i),
            (o.ui = new u.a(r).on("click tap enter", function (e) {
              o.trigger(e.type);
            })),
            e.on("change:state", function (e, t) {
              var n;
              switch (t) {
                case "buffering":
                  n = a.buffer;
                  break;
                case "playing":
                  n = a.pause;
                  break;
                case "idle":
                case "paused":
                  n = a.playback;
                  break;
                case "complete":
                  n = a.replay;
                  break;
                default:
                  n = "";
              }
              "" !== n
                ? r.setAttribute("aria-label", n)
                : r.removeAttribute("aria-label");
            }),
            e.get("displayPlaybackLabel"))
          ) {
            var s = o.icon.getElementsByClassName("jw-idle-icon-text")[0];
            s ||
              ((s = Object(l.e)(
                '<div class="jw-idle-icon-text">'.concat(a.playback, "</div>")
              )),
              Object(l.a)(o.icon, "jw-idle-label"),
              o.icon.appendChild(s));
          }
          return o;
        }
        var n, i, o;
        return (
          (function (e, t) {
            if ("function" != typeof t && null !== t)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (e.prototype = Object.create(t && t.prototype, {
              constructor: { value: e, writable: !0, configurable: !0 },
            })),
              t && vt(e, t);
          })(t, e),
          (n = t),
          (i = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && gt(n.prototype, i),
          o && gt(n, o),
          t
        );
      })(r.a);
      function kt(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var xt = (function () {
        function e(t, n, i) {
          !(function (e, t) {
            if (!(e instanceof t))
              throw new TypeError("Cannot call a class as a function");
          })(this, e);
          var o = i.querySelector(".jw-icon");
          (this.ui = new u.a(o).on("click tap enter", function () {
            n.next({ reason: "interaction" });
          })),
            t.change("nextUp", function (e, t) {
              i.style.visibility = t ? "" : "hidden";
            }),
            (this.el = i);
        }
        var t, n, i;
        return (
          (t = e),
          (n = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && kt(t.prototype, n),
          i && kt(t, i),
          e
        );
      })();
      function Ot(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var Ct = (function () {
        function e(t, n) {
          !(function (e, t) {
            if (!(e instanceof t))
              throw new TypeError("Cannot call a class as a function");
          })(this, e),
            (this.el = Object(l.e)(wt(t.get("localization"))));
          var i = this.el.querySelector(".jw-display-controls"),
            o = {};
          Mt("rewind", pe("rewind"), ht, i, o, t, n),
            Mt("display", pe("play,pause,buffer,replay"), yt, i, o, t, n),
            Mt("next", pe("next"), xt, i, o, t, n),
            (this.container = i),
            (this.buttons = o);
        }
        var t, n, i;
        return (
          (t = e),
          (n = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
            {
              key: "destroy",
              value: function () {
                var e = this.buttons;
                Object.keys(e).forEach(function (t) {
                  e[t].ui && e[t].ui.destroy();
                });
              },
            },
          ]) && Ot(t.prototype, n),
          i && Ot(t, i),
          e
        );
      })();
      function Mt(e, t, n, i, o, a, r) {
        var s = i.querySelector(".jw-display-icon-".concat(e)),
          l = i.querySelector(".jw-icon-".concat(e));
        t.forEach(function (e) {
          l.appendChild(e);
        }),
          (o[e] = new n(a, r, s));
      }
      var St = n(2);
      function Tt(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var Et = (function () {
          function e(t, n, i) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              Object(w.g)(this, r.a),
              (this._model = t),
              (this._api = n),
              (this._playerElement = i),
              (this.localization = t.get("localization")),
              (this.state = "tooltip"),
              (this.enabled = !1),
              (this.shown = !1),
              (this.feedShownId = ""),
              (this.closeUi = null),
              (this.tooltipUi = null),
              this.reset();
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "setup",
                value: function (e) {
                  (this.container = e.createElement("div")),
                    (this.container.className = "jw-nextup-container jw-reset");
                  var t = Object(l.e)(
                    (function () {
                      var e =
                          arguments.length > 0 && void 0 !== arguments[0]
                            ? arguments[0]
                            : "",
                        t =
                          arguments.length > 1 && void 0 !== arguments[1]
                            ? arguments[1]
                            : "",
                        n =
                          arguments.length > 2 && void 0 !== arguments[2]
                            ? arguments[2]
                            : "",
                        i =
                          arguments.length > 3 && void 0 !== arguments[3]
                            ? arguments[3]
                            : "";
                      return (
                        '<div class="jw-nextup jw-background-color jw-reset"><div class="jw-nextup-tooltip jw-reset"><div class="jw-nextup-thumbnail jw-reset"></div><div class="jw-nextup-body jw-reset">' +
                        '<div class="jw-nextup-header jw-reset">'.concat(
                          e,
                          "</div>"
                        ) +
                        '<div class="jw-nextup-title jw-reset-text" dir="auto">'.concat(
                          t,
                          "</div>"
                        ) +
                        '<div class="jw-nextup-duration jw-reset">'.concat(
                          n,
                          "</div>"
                        ) +
                        "</div></div>" +
                        '<button type="button" class="jw-icon jw-nextup-close jw-reset" aria-label="'.concat(
                          i,
                          '"></button>'
                        ) +
                        "</div>"
                      );
                    })()
                  );
                  t.querySelector(".jw-nextup-close").appendChild(de("close")),
                    this.addContent(t),
                    (this.closeButton = this.content.querySelector(
                      ".jw-nextup-close"
                    )),
                    this.closeButton.setAttribute(
                      "aria-label",
                      this.localization.close
                    ),
                    (this.tooltip = this.content.querySelector(
                      ".jw-nextup-tooltip"
                    ));
                  var n = this._model,
                    i = n.player;
                  (this.enabled = !1),
                    n.on("change:nextUp", this.onNextUp, this),
                    i.change("duration", this.onDuration, this),
                    i.change("position", this.onElapsed, this),
                    i.change("streamType", this.onStreamType, this),
                    i.change(
                      "state",
                      function (e, t) {
                        "complete" === t && this.toggle(!1);
                      },
                      this
                    ),
                    (this.closeUi = new u.a(this.closeButton, {
                      directSelect: !0,
                    }).on(
                      "click tap enter",
                      function () {
                        (this.nextUpSticky = !1), this.toggle(!1);
                      },
                      this
                    )),
                    (this.tooltipUi = new u.a(this.tooltip).on(
                      "click tap",
                      this.click,
                      this
                    ));
                },
              },
              {
                key: "loadThumbnail",
                value: function (e) {
                  return (
                    (this.nextUpImage = new Image()),
                    (this.nextUpImage.onload = function () {
                      this.nextUpImage.onload = null;
                    }.bind(this)),
                    (this.nextUpImage.src = e),
                    { backgroundImage: 'url("' + e + '")' }
                  );
                },
              },
              {
                key: "click",
                value: function () {
                  var e = this.feedShownId;
                  this.reset(),
                    this._api.next({ feedShownId: e, reason: "interaction" });
                },
              },
              {
                key: "toggle",
                value: function (e, t) {
                  if (
                    this.enabled &&
                    (Object(l.v)(
                      this.container,
                      "jw-nextup-sticky",
                      !!this.nextUpSticky
                    ),
                    this.shown !== e)
                  ) {
                    (this.shown = e),
                      Object(l.v)(
                        this.container,
                        "jw-nextup-container-visible",
                        e
                      ),
                      Object(l.v)(this._playerElement, "jw-flag-nextup", e);
                    var n = this._model.get("nextUp");
                    e && n
                      ? ((this.feedShownId = Object(ot.b)(ot.a)),
                        this.trigger("nextShown", {
                          mode: n.mode,
                          ui: "nextup",
                          itemsShown: [n],
                          feedData: n.feedData,
                          reason: t,
                          feedShownId: this.feedShownId,
                        }))
                      : (this.feedShownId = "");
                  }
                },
              },
              {
                key: "setNextUpItem",
                value: function (e) {
                  var t = this;
                  setTimeout(function () {
                    if (
                      ((t.thumbnail = t.content.querySelector(
                        ".jw-nextup-thumbnail"
                      )),
                      Object(l.v)(
                        t.content,
                        "jw-nextup-thumbnail-visible",
                        !!e.image
                      ),
                      e.image)
                    ) {
                      var n = t.loadThumbnail(e.image);
                      Object(je.d)(t.thumbnail, n);
                    }
                    (t.header = t.content.querySelector(".jw-nextup-header")),
                      (t.header.textContent = Object(l.e)(
                        t.localization.nextUp
                      ).textContent),
                      (t.title = t.content.querySelector(".jw-nextup-title"));
                    var i = e.title;
                    t.title.textContent = i ? Object(l.e)(i).textContent : "";
                    var o = e.duration;
                    o &&
                      ((t.duration = t.content.querySelector(
                        ".jw-nextup-duration"
                      )),
                      (t.duration.textContent =
                        "number" == typeof o ? Object(ve.timeFormat)(o) : o));
                  }, 500);
                },
              },
              {
                key: "onNextUp",
                value: function (e, t) {
                  this.reset(),
                    t || (t = { showNextUp: !1 }),
                    (this.enabled = !(!t.title && !t.image)),
                    this.enabled &&
                      (t.showNextUp ||
                        ((this.nextUpSticky = !1), this.toggle(!1)),
                      this.setNextUpItem(t));
                },
              },
              {
                key: "onDuration",
                value: function (e, t) {
                  if (t) {
                    var n = e.get("nextupoffset"),
                      i = -10;
                    n && (i = Object(St.d)(n, t)),
                      i < 0 && (i += t),
                      Object(St.c)(n) && t - 5 < i && (i = t - 5),
                      (this.offset = i);
                  }
                },
              },
              {
                key: "onElapsed",
                value: function (e, t) {
                  var n = this.nextUpSticky;
                  if (this.enabled && !1 !== n) {
                    var i = t >= this.offset;
                    i && void 0 === n
                      ? ((this.nextUpSticky = i), this.toggle(i, "time"))
                      : !i && n && this.reset();
                  }
                },
              },
              {
                key: "onStreamType",
                value: function (e, t) {
                  "VOD" !== t && ((this.nextUpSticky = !1), this.toggle(!1));
                },
              },
              {
                key: "element",
                value: function () {
                  return this.container;
                },
              },
              {
                key: "addContent",
                value: function (e) {
                  this.content && this.removeContent(),
                    (this.content = e),
                    this.container.appendChild(e);
                },
              },
              {
                key: "removeContent",
                value: function () {
                  this.content &&
                    (this.container.removeChild(this.content),
                    (this.content = null));
                },
              },
              {
                key: "reset",
                value: function () {
                  (this.nextUpSticky = void 0), this.toggle(!1);
                },
              },
              {
                key: "destroy",
                value: function () {
                  this.off(),
                    this._model.off(null, null, this),
                    this.closeUi && this.closeUi.destroy(),
                    this.tooltipUi && this.tooltipUi.destroy();
                },
              },
            ]) && Tt(t.prototype, n),
            i && Tt(t, i),
            e
          );
        })(),
        _t = function (e, t) {
          var n = e.featured,
            i = e.showLogo,
            o = e.type;
          return (
            (e.logo = i
              ? '<span class="jw-rightclick-logo jw-reset"></span>'
              : ""),
            '<li class="jw-reset jw-rightclick-item '
              .concat(n ? "jw-featured" : "", '">')
              .concat(zt[o](e, t), "</li>")
          );
        },
        zt = {
          link: function (e) {
            var t = e.link,
              n = e.title,
              i = e.logo;
            return '<a href="'
              .concat(
                t || "",
                '" class="jw-rightclick-link jw-reset-text" target="_blank" rel="noreferrer" dir="auto">'
              )
              .concat(i)
              .concat(n || "", "</a>");
          },
          info: function (e, t) {
            return '<button type="button" class="jw-reset-text jw-rightclick-link jw-info-overlay-item" dir="auto">'.concat(
              t.videoInfo,
              "</button>"
            );
          },
          share: function (e, t) {
            return '<button type="button" class="jw-reset-text jw-rightclick-link jw-share-item" dir="auto">'.concat(
              t.sharing.heading,
              "</button>"
            );
          },
          keyboardShortcuts: function (e, t) {
            return '<button type="button" class="jw-reset-text jw-rightclick-link jw-shortcuts-item" dir="auto">'.concat(
              t.shortcuts.keyboardShortcuts,
              "</button>"
            );
          },
        },
        Pt = n(23),
        At = n(6),
        Lt = n(13);
      function Rt(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var It = {
        free: 0,
        pro: 1,
        premium: 2,
        ads: 3,
        invalid: 4,
        enterprise: 6,
        trial: 7,
        platinum: 8,
        starter: 9,
        business: 10,
        developer: 11,
      };
      function Bt(e) {
        var t = Object(l.e)(e),
          n = t.querySelector(".jw-rightclick-logo");
        return n && n.appendChild(de("jwplayer-logo")), t;
      }
      var Vt = (function () {
          function e(t, n) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              (this.infoOverlay = t),
              (this.shortcutsTooltip = n);
          }
          var t, n, i;
          return (
            (t = e),
            (n = [
              {
                key: "buildArray",
                value: function () {
                  var e = Pt.a.split("+")[0],
                    t = this.model,
                    n = t.get("edition"),
                    i = t.get("localization").poweredBy,
                    o = '<span class="jw-reset">JW Player '.concat(
                      e,
                      "</span>"
                    ),
                    a = {
                      items: [
                        { type: "info" },
                        {
                          title: Object(Lt.e)(i)
                            ? "".concat(o, " ").concat(i)
                            : "".concat(i, " ").concat(o),
                          type: "link",
                          featured: !0,
                          showLogo: !0,
                          link: "https://jwplayer.com/learn-more?e=".concat(
                            It[n]
                          ),
                        },
                      ],
                    },
                    r = t.get("provider"),
                    s = a.items;
                  if (r && r.name.indexOf("flash") >= 0) {
                    var l = "Flash Version " + Object(At.a)();
                    s.push({
                      title: l,
                      type: "link",
                      link: "http://www.adobe.com/software/flash/about/",
                    });
                  }
                  return (
                    this.shortcutsTooltip &&
                      s.splice(s.length - 1, 0, { type: "keyboardShortcuts" }),
                    a
                  );
                },
              },
              {
                key: "rightClick",
                value: function (e) {
                  if ((this.lazySetup(), this.mouseOverContext)) return !1;
                  this.hideMenu(), this.showMenu(e), this.addHideMenuHandlers();
                },
              },
              {
                key: "getOffset",
                value: function (e) {
                  var t = Object(l.c)(this.wrapperElement),
                    n = e.pageX - t.left,
                    i = e.pageY - t.top;
                  return (
                    this.model.get("touchMode") && (i -= 100), { x: n, y: i }
                  );
                },
              },
              {
                key: "showMenu",
                value: function (e) {
                  var t = this,
                    n = this.getOffset(e);
                  return (
                    (this.el.style.left = n.x + "px"),
                    (this.el.style.top = n.y + "px"),
                    (this.outCount = 0),
                    Object(l.a)(
                      this.playerContainer,
                      "jw-flag-rightclick-open"
                    ),
                    Object(l.a)(this.el, "jw-open"),
                    clearTimeout(this._menuTimeout),
                    (this._menuTimeout = setTimeout(function () {
                      return t.hideMenu();
                    }, 3e3)),
                    !1
                  );
                },
              },
              {
                key: "hideMenu",
                value: function (e) {
                  (e && this.el && this.el.contains(e.target)) ||
                    (Object(l.o)(
                      this.playerContainer,
                      "jw-flag-rightclick-open"
                    ),
                    Object(l.o)(this.el, "jw-open"));
                },
              },
              {
                key: "lazySetup",
                value: function () {
                  var e,
                    t,
                    n,
                    i,
                    o = this,
                    a =
                      ((e = this.buildArray()),
                      (t = this.model.get("localization")),
                      (n = e.items),
                      (i = (void 0 === n ? [] : n).map(function (e) {
                        return _t(e, t);
                      })),
                      '<div class="jw-rightclick jw-reset">' +
                        '<ul class="jw-rightclick-list jw-reset">'.concat(
                          i.join(""),
                          "</ul>"
                        ) +
                        "</div>");
                  if (this.el) {
                    if (this.html !== a) {
                      this.html = a;
                      var r = Bt(a);
                      Object(l.h)(this.el);
                      for (var s = r.childNodes.length; s--; )
                        this.el.appendChild(r.firstChild);
                    }
                  } else
                    (this.html = a),
                      (this.el = Bt(this.html)),
                      this.wrapperElement.appendChild(this.el),
                      (this.hideMenuHandler = function (e) {
                        return o.hideMenu(e);
                      }),
                      (this.overHandler = function () {
                        o.mouseOverContext = !0;
                      }),
                      (this.outHandler = function (e) {
                        (o.mouseOverContext = !1),
                          e.relatedTarget &&
                            !o.el.contains(e.relatedTarget) &&
                            ++o.outCount > 1 &&
                            o.hideMenu();
                      }),
                      (this.infoOverlayHandler = function () {
                        (o.mouseOverContext = !1),
                          o.hideMenu(),
                          o.infoOverlay.open();
                      }),
                      (this.shortcutsTooltipHandler = function () {
                        (o.mouseOverContext = !1),
                          o.hideMenu(),
                          o.shortcutsTooltip.open();
                      });
                },
              },
              {
                key: "setup",
                value: function (e, t, n) {
                  (this.wrapperElement = n),
                    (this.model = e),
                    (this.mouseOverContext = !1),
                    (this.playerContainer = t),
                    (this.ui = new u.a(n).on(
                      "longPress",
                      this.rightClick,
                      this
                    ));
                },
              },
              {
                key: "addHideMenuHandlers",
                value: function () {
                  this.removeHideMenuHandlers(),
                    this.wrapperElement.addEventListener(
                      "touchstart",
                      this.hideMenuHandler
                    ),
                    document.addEventListener(
                      "touchstart",
                      this.hideMenuHandler
                    ),
                    o.OS.mobile ||
                      (this.wrapperElement.addEventListener(
                        "click",
                        this.hideMenuHandler
                      ),
                      document.addEventListener("click", this.hideMenuHandler),
                      this.el.addEventListener("mouseover", this.overHandler),
                      this.el.addEventListener("mouseout", this.outHandler)),
                    this.el
                      .querySelector(".jw-info-overlay-item")
                      .addEventListener("click", this.infoOverlayHandler),
                    this.shortcutsTooltip &&
                      this.el
                        .querySelector(".jw-shortcuts-item")
                        .addEventListener(
                          "click",
                          this.shortcutsTooltipHandler
                        );
                },
              },
              {
                key: "removeHideMenuHandlers",
                value: function () {
                  this.wrapperElement &&
                    (this.wrapperElement.removeEventListener(
                      "click",
                      this.hideMenuHandler
                    ),
                    this.wrapperElement.removeEventListener(
                      "touchstart",
                      this.hideMenuHandler
                    )),
                    this.el &&
                      (this.el
                        .querySelector(".jw-info-overlay-item")
                        .removeEventListener("click", this.infoOverlayHandler),
                      this.el.removeEventListener(
                        "mouseover",
                        this.overHandler
                      ),
                      this.el.removeEventListener("mouseout", this.outHandler),
                      this.shortcutsTooltip &&
                        this.el
                          .querySelector(".jw-shortcuts-item")
                          .removeEventListener(
                            "click",
                            this.shortcutsTooltipHandler
                          )),
                    document.removeEventListener("click", this.hideMenuHandler),
                    document.removeEventListener(
                      "touchstart",
                      this.hideMenuHandler
                    );
                },
              },
              {
                key: "destroy",
                value: function () {
                  clearTimeout(this._menuTimeout),
                    this.removeHideMenuHandlers(),
                    this.el &&
                      (this.hideMenu(),
                      (this.hideMenuHandler = null),
                      (this.el = null)),
                    this.wrapperElement &&
                      ((this.wrapperElement.oncontextmenu = null),
                      (this.wrapperElement = null)),
                    this.model && (this.model = null),
                    this.ui && (this.ui.destroy(), (this.ui = null));
                },
              },
            ]) && Rt(t.prototype, n),
            i && Rt(t, i),
            e
          );
        })(),
        Nt = function (e) {
          return '<button type="button" class="jw-reset-text jw-settings-content-item" dir="auto">'.concat(
            e,
            "</button>"
          );
        },
        Ht = function (e) {
          return (
            '<button type="button" class="jw-reset-text jw-settings-content-item" dir="auto">' +
            "".concat(e.label) +
            "<div class='jw-reset jw-settings-value-wrapper'>" +
            '<div class="jw-reset-text jw-settings-content-item-value">'.concat(
              e.value,
              "</div>"
            ) +
            '<div class="jw-reset-text jw-settings-content-item-arrow">'.concat(
              Y.a,
              "</div>"
            ) +
            "</div></button>"
          );
        },
        Ft = function (e) {
          return (
            '<button type="button" class="jw-reset-text jw-settings-content-item" role="menuitemradio" aria-checked="false" dir="auto">' +
            "".concat(e) +
            "</button>"
          );
        };
      function qt(e) {
        return (qt =
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
      function Dt(e, t) {
        return !t || ("object" !== qt(t) && "function" != typeof t)
          ? (function (e) {
              if (void 0 === e)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return e;
            })(e)
          : t;
      }
      function Ut(e) {
        return (Ut = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function Wt(e, t) {
        return (Wt =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      function Qt(e, t) {
        if (!(e instanceof t))
          throw new TypeError("Cannot call a class as a function");
      }
      function Yt(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function Xt(e, t, n) {
        return t && Yt(e.prototype, t), n && Yt(e, n), e;
      }
      var Zt,
        Kt = (function () {
          function e(t, n) {
            var i =
              arguments.length > 2 && void 0 !== arguments[2]
                ? arguments[2]
                : Nt;
            Qt(this, e),
              (this.el = Object(l.e)(i(t))),
              (this.ui = new u.a(this.el).on("click tap enter", n, this));
          }
          return (
            Xt(e, [
              {
                key: "destroy",
                value: function () {
                  this.ui.destroy();
                },
              },
            ]),
            e
          );
        })(),
        Jt = (function (e) {
          function t(e, n) {
            var i =
              arguments.length > 2 && void 0 !== arguments[2]
                ? arguments[2]
                : Ft;
            return Qt(this, t), Dt(this, Ut(t).call(this, e, n, i));
          }
          return (
            (function (e, t) {
              if ("function" != typeof t && null !== t)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (e.prototype = Object.create(t && t.prototype, {
                constructor: { value: e, writable: !0, configurable: !0 },
              })),
                t && Wt(e, t);
            })(t, e),
            Xt(t, [
              {
                key: "activate",
                value: function () {
                  Object(l.v)(this.el, "jw-settings-item-active", !0),
                    this.el.setAttribute("aria-checked", "true"),
                    (this.active = !0);
                },
              },
              {
                key: "deactivate",
                value: function () {
                  Object(l.v)(this.el, "jw-settings-item-active", !1),
                    this.el.setAttribute("aria-checked", "false"),
                    (this.active = !1);
                },
              },
            ]),
            t
          );
        })(Kt),
        Gt = function (e, t) {
          return e
            ? '<div class="jw-reset jw-settings-submenu jw-settings-submenu-'.concat(
                t,
                '" role="menu" aria-expanded="false">'
              ) + '<div class="jw-settings-submenu-items"></div></div>'
            : '<div class="jw-reset jw-settings-menu" role="menu" aria-expanded="false"><div class="jw-reset jw-settings-topbar" role="menubar"><div class="jw-reset jw-settings-topbar-text" tabindex="0"></div><div class="jw-reset jw-settings-topbar-buttons"></div></div></div>';
        },
        $t = function (e, t) {
          var n = e.name,
            i = {
              captions: "cc-off",
              audioTracks: "audio-tracks",
              quality: "quality-100",
              playbackRates: "playback-rate",
            }[n];
          if (i || e.icon) {
            var o = p(
                "jw-settings-".concat(n, " jw-submenu-").concat(n),
                function (t) {
                  e.open(t);
                },
                n,
                [(e.icon && Object(l.e)(e.icon)) || de(i)]
              ),
              a = o.element();
            return (
              a.setAttribute("role", "menuitemradio"),
              a.setAttribute("aria-checked", "false"),
              a.setAttribute("aria-label", t),
              "ontouchstart" in window || (o.tooltip = it(a, n, t)),
              o
            );
          }
        };
      function en(e) {
        return (en =
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
      function tn(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function nn(e) {
        return (nn = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function on(e, t) {
        return (on =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      function an(e) {
        if (void 0 === e)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return e;
      }
      var rn = (function (e) {
          function t(e, n, i) {
            var o,
              a,
              r,
              s =
                arguments.length > 3 && void 0 !== arguments[3]
                  ? arguments[3]
                  : Gt;
            return (
              (function (e, t) {
                if (!(e instanceof t))
                  throw new TypeError("Cannot call a class as a function");
              })(this, t),
              (a = this),
              ((o =
                !(r = nn(t).call(this)) ||
                ("object" !== en(r) && "function" != typeof r)
                  ? an(a)
                  : r).open = o.open.bind(an(an(o)))),
              (o.close = o.close.bind(an(an(o)))),
              (o.toggle = o.toggle.bind(an(an(o)))),
              (o.onDocumentClick = o.onDocumentClick.bind(an(an(o)))),
              (o.name = e),
              (o.isSubmenu = !!n),
              (o.el = Object(l.e)(s(o.isSubmenu, e))),
              (o.topbar = o.el.querySelector(".jw-".concat(o.name, "-topbar"))),
              (o.buttonContainer = o.el.querySelector(
                ".jw-".concat(o.name, "-topbar-buttons")
              )),
              (o.children = {}),
              (o.openMenus = []),
              (o.items = []),
              (o.visible = !1),
              (o.parentMenu = n),
              (o.mainMenu = o.parentMenu ? o.parentMenu.mainMenu : an(an(o))),
              (o.categoryButton = null),
              (o.closeButton =
                (o.parentMenu && o.parentMenu.closeButton) ||
                o.createCloseButton(i)),
              o.isSubmenu
                ? ((o.categoryButton =
                    o.parentMenu.categoryButton || o.createCategoryButton(i)),
                  o.parentMenu.parentMenu &&
                    !o.mainMenu.backButton &&
                    (o.mainMenu.backButton = o.createBackButton(i)),
                  (o.itemsContainer = o.createItemsContainer()),
                  o.parentMenu.appendMenu(an(an(o))))
                : (o.ui = sn(an(an(o)))),
              o
            );
          }
          var n, i, o;
          return (
            (function (e, t) {
              if ("function" != typeof t && null !== t)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (e.prototype = Object.create(t && t.prototype, {
                constructor: { value: e, writable: !0, configurable: !0 },
              })),
                t && on(e, t);
            })(t, e),
            (n = t),
            (i = [
              {
                key: "createItemsContainer",
                value: function () {
                  var e,
                    t,
                    n = this,
                    i = this.el.querySelector(".jw-settings-submenu-items"),
                    o = new u.a(i),
                    a =
                      (this.categoryButton && this.categoryButton.element()) ||
                      (this.parentMenu.categoryButton &&
                        this.parentMenu.categoryButton.element()) ||
                      this.mainMenu.buttonContainer.firstChild;
                  return (
                    this.parentMenu.isSubmenu &&
                      ((e = this.mainMenu.closeButton.element()),
                      (t = this.mainMenu.backButton.element())),
                    o.on("keydown", function (o) {
                      if (o.target.parentNode === i) {
                        var r = function (e, t) {
                            e
                              ? e.focus()
                              : void 0 !== t && i.childNodes[t].focus();
                          },
                          s = o.sourceEvent,
                          c = s.target,
                          u = i.firstChild === c,
                          d = i.lastChild === c,
                          p = n.topbar,
                          w = e || Object(l.k)(a),
                          f = t || Object(l.n)(a),
                          h = Object(l.k)(s.target),
                          j = Object(l.n)(s.target),
                          g = s.key.replace(/(Arrow|ape)/, "");
                        switch (g) {
                          case "Tab":
                            r(s.shiftKey ? f : w);
                            break;
                          case "Left":
                            r(
                              f ||
                                Object(l.n)(
                                  document.getElementsByClassName(
                                    "jw-icon-settings"
                                  )[0]
                                )
                            );
                            break;
                          case "Up":
                            p && u
                              ? r(p.firstChild)
                              : r(j, i.childNodes.length - 1);
                            break;
                          case "Right":
                            r(w);
                            break;
                          case "Down":
                            p && d ? r(p.firstChild) : r(h, 0);
                        }
                        s.preventDefault(), "Esc" !== g && s.stopPropagation();
                      }
                    }),
                    o
                  );
                },
              },
              {
                key: "createCloseButton",
                value: function (e) {
                  var t = p("jw-settings-close", this.close, e.close, [
                    de("close"),
                  ]);
                  return (
                    this.topbar.appendChild(t.element()),
                    t.show(),
                    t.ui.on(
                      "keydown",
                      function (e) {
                        var t = e.sourceEvent,
                          n = t.key.replace(/(Arrow|ape)/, "");
                        ("Enter" === n ||
                          "Right" === n ||
                          ("Tab" === n && !t.shiftKey)) &&
                          this.close(e);
                      },
                      this
                    ),
                    this.buttonContainer.appendChild(t.element()),
                    t
                  );
                },
              },
              {
                key: "createCategoryButton",
                value: function (e) {
                  var t =
                    e[
                      {
                        captions: "cc",
                        audioTracks: "audioTracks",
                        quality: "hd",
                        playbackRates: "playbackRates",
                      }[this.name]
                    ];
                  "sharing" === this.name && (t = e.sharing.heading);
                  var n = $t(this, t);
                  return n.element().setAttribute("name", this.name), n;
                },
              },
              {
                key: "createBackButton",
                value: function (e) {
                  var t = p(
                    "jw-settings-back",
                    function (e) {
                      Zt && Zt.open(e);
                    },
                    e.close,
                    [de("arrow-left")]
                  );
                  return Object(l.m)(this.mainMenu.topbar, t.element()), t;
                },
              },
              {
                key: "createTopbar",
                value: function () {
                  var e = Object(l.e)('<div class="jw-submenu-topbar"></div>');
                  return Object(l.m)(this.el, e), e;
                },
              },
              {
                key: "createItems",
                value: function (e, t) {
                  var n = this,
                    i =
                      arguments.length > 2 && void 0 !== arguments[2]
                        ? arguments[2]
                        : {},
                    o =
                      arguments.length > 3 && void 0 !== arguments[3]
                        ? arguments[3]
                        : Jt,
                    a = this.name,
                    r = e.map(function (e, r) {
                      var s, l;
                      switch (a) {
                        case "quality":
                          s =
                            "Auto" === e.label && 0 === r
                              ? "".concat(
                                  i.defaultText,
                                  '&nbsp;<span class="jw-reset jw-auto-label"></span>'
                                )
                              : e.label;
                          break;
                        case "captions":
                          s =
                            ("Off" !== e.label && "off" !== e.id) || 0 !== r
                              ? e.label
                              : i.defaultText;
                          break;
                        case "playbackRates":
                          (l = e),
                            (s = Object(Lt.e)(i.tooltipText)
                              ? "x" + e
                              : e + "x");
                          break;
                        case "audioTracks":
                          s = e.name;
                      }
                      s || ((s = e), "object" === en(e) && (s.options = i));
                      var c = new o(
                        s,
                        function (e) {
                          c.active ||
                            (t(l || r),
                            c.deactivate &&
                              (n.items
                                .filter(function (e) {
                                  return !0 === e.active;
                                })
                                .forEach(function (e) {
                                  e.deactivate();
                                }),
                              Zt ? Zt.open(e) : n.mainMenu.close(e)),
                            c.activate && c.activate());
                        }.bind(n)
                      );
                      return c;
                    });
                  return r;
                },
              },
              {
                key: "setMenuItems",
                value: function (e, t) {
                  var n = this;
                  e
                    ? ((this.items = []),
                      Object(l.h)(this.itemsContainer.el),
                      e.forEach(function (e) {
                        n.items.push(e), n.itemsContainer.el.appendChild(e.el);
                      }),
                      t > -1 && e[t].activate(),
                      this.categoryButton.show())
                    : this.removeMenu();
                },
              },
              {
                key: "appendMenu",
                value: function (e) {
                  if (e) {
                    var t = e.el,
                      n = e.name,
                      i = e.categoryButton;
                    if (((this.children[n] = e), i)) {
                      var o = this.mainMenu.buttonContainer,
                        a = o.querySelector(".jw-settings-sharing"),
                        r =
                          "quality" === n
                            ? o.firstChild
                            : a || this.closeButton.element();
                      o.insertBefore(i.element(), r);
                    }
                    this.mainMenu.el.appendChild(t);
                  }
                },
              },
              {
                key: "removeMenu",
                value: function (e) {
                  if (!e) return this.parentMenu.removeMenu(this.name);
                  var t = this.children[e];
                  t && (delete this.children[e], t.destroy());
                },
              },
              {
                key: "open",
                value: function (e) {
                  if (!this.visible || this.openMenus) {
                    var t;
                    if (((Zt = null), this.isSubmenu)) {
                      var n = this.mainMenu,
                        i = this.parentMenu,
                        o = this.categoryButton;
                      if (
                        (i.openMenus.length && i.closeChildren(),
                        o && o.element().setAttribute("aria-checked", "true"),
                        i.isSubmenu)
                      ) {
                        i.el.classList.remove("jw-settings-submenu-active"),
                          n.topbar.classList.add("jw-nested-menu-open");
                        var a = n.topbar.querySelector(
                          ".jw-settings-topbar-text"
                        );
                        a.setAttribute("name", this.name),
                          (a.innerText = this.title || this.name),
                          n.backButton.show(),
                          (Zt = this.parentMenu),
                          (t = this.topbar
                            ? this.topbar.firstChild
                            : e && "enter" === e.type
                            ? this.items[0].el
                            : a);
                      } else
                        n.topbar.classList.remove("jw-nested-menu-open"),
                          n.backButton && n.backButton.hide();
                      this.el.classList.add("jw-settings-submenu-active"),
                        i.openMenus.push(this.name),
                        n.visible ||
                          (n.open(e),
                          this.items && e && "enter" === e.type
                            ? (t = this.topbar
                                ? this.topbar.firstChild.focus()
                                : this.items[0].el)
                            : o.tooltip &&
                              ((o.tooltip.suppress = !0), (t = o.element()))),
                        this.openMenus.length && this.closeChildren(),
                        t && t.focus(),
                        (this.el.scrollTop = 0);
                    } else
                      this.el.parentNode.classList.add("jw-settings-open"),
                        this.trigger("menuVisibility", { visible: !0, evt: e }),
                        document.addEventListener(
                          "click",
                          this.onDocumentClick
                        );
                    (this.visible = !0),
                      this.el.setAttribute("aria-expanded", "true");
                  }
                },
              },
              {
                key: "close",
                value: function (e) {
                  var t = this;
                  this.visible &&
                    ((this.visible = !1),
                    this.el.setAttribute("aria-expanded", "false"),
                    this.isSubmenu
                      ? (this.el.classList.remove("jw-settings-submenu-active"),
                        this.categoryButton
                          .element()
                          .setAttribute("aria-checked", "false"),
                        (this.parentMenu.openMenus = this.parentMenu.openMenus.filter(
                          function (e) {
                            return e !== t.name;
                          }
                        )),
                        !this.mainMenu.openMenus.length &&
                          this.mainMenu.visible &&
                          this.mainMenu.close(e))
                      : (this.el.parentNode.classList.remove(
                          "jw-settings-open"
                        ),
                        this.trigger("menuVisibility", { visible: !1, evt: e }),
                        document.removeEventListener(
                          "click",
                          this.onDocumentClick
                        )),
                    this.openMenus.length && this.closeChildren());
                },
              },
              {
                key: "closeChildren",
                value: function () {
                  var e = this;
                  this.openMenus.forEach(function (t) {
                    var n = e.children[t];
                    n && n.close();
                  });
                },
              },
              {
                key: "toggle",
                value: function (e) {
                  this.visible ? this.close(e) : this.open(e);
                },
              },
              {
                key: "onDocumentClick",
                value: function (e) {
                  /jw-(settings|video|nextup-close|sharing-link|share-item)/.test(
                    e.target.className
                  ) || this.close();
                },
              },
              {
                key: "destroy",
                value: function () {
                  var e = this;
                  if (
                    (document.removeEventListener(
                      "click",
                      this.onDocumentClick
                    ),
                    Object.keys(this.children).map(function (t) {
                      e.children[t].destroy();
                    }),
                    this.isSubmenu)
                  ) {
                    this.parentMenu.name === this.mainMenu.name &&
                      this.categoryButton &&
                      (this.parentMenu.buttonContainer.removeChild(
                        this.categoryButton.element()
                      ),
                      this.categoryButton.ui.destroy()),
                      this.itemsContainer && this.itemsContainer.destroy();
                    var t = this.parentMenu.openMenus,
                      n = t.indexOf(this.name);
                    t.length && n > -1 && this.openMenus.splice(n, 1),
                      delete this.parentMenu;
                  } else this.ui.destroy();
                  (this.visible = !1),
                    this.el.parentNode &&
                      this.el.parentNode.removeChild(this.el);
                },
              },
              {
                key: "defaultChild",
                get: function () {
                  var e = this.children,
                    t = e.quality,
                    n = e.captions,
                    i = e.audioTracks,
                    o = e.sharing,
                    a = e.playbackRates;
                  return t || n || i || o || a;
                },
              },
            ]) && tn(n.prototype, i),
            o && tn(n, o),
            t
          );
        })(r.a),
        sn = function (e) {
          var t = e.closeButton,
            n = e.el;
          return new u.a(n).on("keydown", function (n) {
            var i = n.sourceEvent,
              o = n.target,
              a = Object(l.k)(o),
              r = Object(l.n)(o),
              s = i.key.replace(/(Arrow|ape)/, ""),
              c = function (t) {
                r ? t || r.focus() : e.close(n);
              };
            switch (s) {
              case "Esc":
                e.close(n);
                break;
              case "Left":
                c();
                break;
              case "Right":
                a && t.element() && o !== t.element() && a.focus();
                break;
              case "Tab":
                i.shiftKey && c(!0);
                break;
              case "Up":
              case "Down":
                !(function () {
                  var t = e.children[o.getAttribute("name")];
                  if ((!t && Zt && (t = Zt.children[Zt.openMenus]), t))
                    return (
                      t.open(n),
                      void (t.topbar
                        ? t.topbar.firstChild.focus()
                        : t.items && t.items.length && t.items[0].el.focus())
                    );
                  if (
                    n.target.parentNode.classList.contains("jw-submenu-topbar")
                  ) {
                    var i = n.target.parentNode.parentNode.querySelector(
                      ".jw-settings-submenu-items"
                    );
                    ("Down" === s
                      ? i.childNodes[0]
                      : i.childNodes[i.childNodes.length - 1]
                    ).focus();
                  }
                })();
            }
            if ((i.stopPropagation(), /13|32|37|38|39|40/.test(i.keyCode)))
              return i.preventDefault(), !1;
          });
        },
        ln = n(59),
        cn = function (e) {
          return fn[e];
        },
        un = function (e) {
          for (var t, n = Object.keys(fn), i = 0; i < n.length; i++)
            if (fn[n[i]] === e) {
              t = n[i];
              break;
            }
          return t;
        },
        dn = function (e) {
          return e + "%";
        },
        pn = function (e) {
          return parseInt(e);
        },
        wn = [
          {
            name: "Font Color",
            propertyName: "color",
            options: [
              "White",
              "Black",
              "Red",
              "Green",
              "Blue",
              "Yellow",
              "Magenta",
              "Cyan",
            ],
            defaultVal: "White",
            getTypedValue: cn,
            getOption: un,
          },
          {
            name: "Font Opacity",
            propertyName: "fontOpacity",
            options: ["100%", "75%", "25%"],
            defaultVal: "100%",
            getTypedValue: pn,
            getOption: dn,
          },
          {
            name: "Font Size",
            propertyName: "userFontScale",
            options: ["200%", "175%", "150%", "125%", "100%", "75%", "50%"],
            defaultVal: "100%",
            getTypedValue: function (e) {
              return parseInt(e) / 100;
            },
            getOption: function (e) {
              return 100 * e + "%";
            },
          },
          {
            name: "Font Family",
            propertyName: "fontFamily",
            options: [
              "Arial",
              "Courier",
              "Georgia",
              "Impact",
              "Lucida Console",
              "Tahoma",
              "Times New Roman",
              "Trebuchet MS",
              "Verdana",
            ],
            defaultVal: "Arial",
            getTypedValue: function (e) {
              return e;
            },
            getOption: function (e) {
              return e;
            },
          },
          {
            name: "Character Edge",
            propertyName: "edgeStyle",
            options: ["None", "Raised", "Depressed", "Uniform", "Drop Shadow"],
            defaultVal: "None",
            getTypedValue: function (e) {
              return e.toLowerCase().replace(/ /g, "");
            },
            getOption: function (e) {
              if ("dropshadow" === e) return "Drop Shadow";
              var t = e.replace(/([A-Z])/g, " $1");
              return t.charAt(0).toUpperCase() + t.slice(1);
            },
          },
          {
            name: "Background Color",
            propertyName: "backgroundColor",
            options: [
              "White",
              "Black",
              "Red",
              "Green",
              "Blue",
              "Yellow",
              "Magenta",
              "Cyan",
            ],
            defaultVal: "Black",
            getTypedValue: cn,
            getOption: un,
          },
          {
            name: "Background Opacity",
            propertyName: "backgroundOpacity",
            options: ["100%", "75%", "50%", "25%", "0%"],
            defaultVal: "50%",
            getTypedValue: pn,
            getOption: dn,
          },
          {
            name: "Window Color",
            propertyName: "windowColor",
            options: [
              "White",
              "Black",
              "Red",
              "Green",
              "Blue",
              "Yellow",
              "Magenta",
              "Cyan",
            ],
            defaultVal: "Black",
            getTypedValue: cn,
            getOption: un,
          },
          {
            name: "Window Opacity",
            propertyName: "windowOpacity",
            options: ["100%", "75%", "50%", "25%", "0%"],
            defaultVal: "0%",
            getTypedValue: pn,
            getOption: dn,
          },
        ],
        fn = {
          White: "#ffffff",
          Black: "#000000",
          Red: "#ff0000",
          Green: "#00ff00",
          Blue: "#0000ff",
          Yellow: "#ffff00",
          Magenta: "ff00ff",
          Cyan: "#00ffff",
        },
        hn = function (e, t, n, i) {
          var o = new rn("settings", null, i),
            a = function (e, t, a, r, s) {
              var l = n.elements["".concat(e, "Button")];
              if (!t || t.length <= 1)
                return o.removeMenu(e), void (l && l.hide());
              var c = o.children[e];
              c || (c = new rn(e, o, i)),
                c.setMenuItems(c.createItems(t, a, s), r),
                l && l.show();
            },
            r = function (r) {
              var s = { defaultText: i.auto };
              a(
                "quality",
                r,
                function (t) {
                  return e.setCurrentQuality(t);
                },
                t.get("currentLevel") || 0,
                s
              );
              var l = o.children,
                c = !!l.quality || l.playbackRates || Object.keys(l).length > 1;
              n.elements.settingsButton.toggle(c);
            };
          t.change(
            "levels",
            function (e, t) {
              r(t);
            },
            o
          );
          var s = function (e, n, i) {
            var o = t.get("levels");
            if (o && "Auto" === o[0].label && n && n.items.length) {
              var a = n.items[0].el.querySelector(".jw-auto-label"),
                r = o[e.index] || { label: "" };
              a.textContent = i ? "" : r.label;
            }
          };
          t.on("change:visualQuality", function (e, n) {
            var i = o.children.quality;
            n && i && s(n.level, i, t.get("currentLevel"));
          }),
            t.on(
              "change:currentLevel",
              function (e, n) {
                var i = o.children.quality,
                  a = t.get("visualQuality");
                a && i && s(a.level, i, n);
              },
              o
            ),
            t.change("captionsList", function (n, r) {
              var s = { defaultText: i.off },
                l = t.get("captionsIndex");
              a(
                "captions",
                r,
                function (t) {
                  return e.setCurrentCaptions(t);
                },
                l,
                s
              );
              var c = o.children.captions;
              if (c && !c.children.captionsSettings) {
                c.topbar = c.topbar || c.createTopbar();
                var u = new rn("captionsSettings", c, i);
                u.title = "Subtitle Settings";
                var d = new Kt("Settings", u.open);
                c.topbar.appendChild(d.el);
                var p = new Jt("Reset", function () {
                  t.set("captions", ln.a), h();
                });
                p.el.classList.add("jw-settings-reset");
                var f = t.get("captions"),
                  h = function () {
                    var e = [];
                    wn.forEach(function (n) {
                      f &&
                        f[n.propertyName] &&
                        (n.defaultVal = n.getOption(f[n.propertyName]));
                      var o = new rn(n.name, u, i),
                        a = new Kt(
                          { label: n.name, value: n.defaultVal },
                          o.open,
                          Ht
                        ),
                        r = o.createItems(
                          n.options,
                          function (e) {
                            var i = a.el.querySelector(
                              ".jw-settings-content-item-value"
                            );
                            !(function (e, n) {
                              var i = t.get("captions"),
                                o = e.propertyName,
                                a = e.options && e.options[n],
                                r = e.getTypedValue(a),
                                s = Object(w.g)({}, i);
                              (s[o] = r), t.set("captions", s);
                            })(n, e),
                              (i.innerText = n.options[e]);
                          },
                          null
                        );
                      o.setMenuItems(r, n.options.indexOf(n.defaultVal) || 0),
                        e.push(a);
                    }),
                      e.push(p),
                      u.setMenuItems(e);
                  };
                h();
              }
            });
          var l = function (e, t) {
            e && t > -1 && e.items[t].activate();
          };
          t.change(
            "captionsIndex",
            function (e, t) {
              var i = o.children.captions;
              i && l(i, t), n.toggleCaptionsButtonState(!!t);
            },
            o
          );
          var c = function (n) {
            if (
              t.get("supportsPlaybackRate") &&
              "LIVE" !== t.get("streamType") &&
              t.get("playbackRateControls")
            ) {
              var r = n.indexOf(t.get("playbackRate")),
                s = { tooltipText: i.playbackRates };
              a(
                "playbackRates",
                n,
                function (t) {
                  return e.setPlaybackRate(t);
                },
                r,
                s
              );
            } else o.children.playbackRates && o.removeMenu("playbackRates");
          };
          t.on(
            "change:playbackRates",
            function (e, t) {
              c(t);
            },
            o
          );
          var u = function (n) {
            a(
              "audioTracks",
              n,
              function (t) {
                return e.setCurrentAudioTrack(t);
              },
              t.get("currentAudioTrack")
            );
          };
          return (
            t.on(
              "change:audioTracks",
              function (e, t) {
                u(t);
              },
              o
            ),
            t.on(
              "change:playbackRate",
              function (e, n) {
                var i = t.get("playbackRates"),
                  a = -1;
                i && (a = i.indexOf(n)), l(o.children.playbackRates, a);
              },
              o
            ),
            t.on(
              "change:currentAudioTrack",
              function (e, t) {
                o.children.audioTracks.items[t].activate();
              },
              o
            ),
            t.on(
              "change:playlistItem",
              function () {
                o.removeMenu("captions"),
                  n.elements.captionsButton.hide(),
                  o.visible && o.close();
              },
              o
            ),
            t.on("change:playbackRateControls", function () {
              c(t.get("playbackRates"));
            }),
            t.on(
              "change:castActive",
              function (e, n, i) {
                n !== i &&
                  (n
                    ? (o.removeMenu("audioTracks"),
                      o.removeMenu("quality"),
                      o.removeMenu("playbackRates"))
                    : (u(t.get("audioTracks")),
                      r(t.get("levels")),
                      c(t.get("playbackRates"))));
              },
              o
            ),
            t.on(
              "change:streamType",
              function () {
                c(t.get("playbackRates"));
              },
              o
            ),
            o
          );
        },
        jn = n(58),
        gn = n(35),
        bn = n(12),
        mn = function (e, t, n, i) {
          var o = Object(l.e)(
              '<div class="jw-reset jw-info-overlay jw-modal"><div class="jw-reset jw-info-container"><div class="jw-reset-text jw-info-title" dir="auto"></div><div class="jw-reset-text jw-info-duration" dir="auto"></div><div class="jw-reset-text jw-info-description" dir="auto"></div></div><div class="jw-reset jw-info-clientid"></div></div>'
            ),
            r = !1,
            s = null,
            c = !1,
            u = function (e) {
              /jw-info/.test(e.target.className) || w.close();
            },
            d = function () {
              var i,
                a,
                s,
                c,
                u,
                d = p(
                  "jw-info-close",
                  function () {
                    w.close();
                  },
                  t.get("localization").close,
                  [de("close")]
                );
              d.show(),
                Object(l.m)(o, d.element()),
                (a = o.querySelector(".jw-info-title")),
                (s = o.querySelector(".jw-info-duration")),
                (c = o.querySelector(".jw-info-description")),
                (u = o.querySelector(".jw-info-clientid")),
                t.change("playlistItem", function (e, t) {
                  var n = t.description,
                    i = t.title;
                  Object(l.q)(c, n || ""), Object(l.q)(a, i || "Unknown Title");
                }),
                t.change(
                  "duration",
                  function (e, n) {
                    var i = "";
                    switch (t.get("streamType")) {
                      case "LIVE":
                        i = "Live";
                        break;
                      case "DVR":
                        i = "DVR";
                        break;
                      default:
                        n && (i = Object(ve.timeFormat)(n));
                    }
                    s.textContent = i;
                  },
                  w
                ),
                (u.textContent =
                  (i = n.getPlugin("jwpsrv")) &&
                  "function" == typeof i.doNotTrackUser &&
                  i.doNotTrackUser()
                    ? ""
                    : "Client ID: ".concat(
                        (function () {
                          try {
                            return window.localStorage.jwplayerLocalId;
                          } catch (e) {
                            return "none";
                          }
                        })()
                      )),
                e.appendChild(o),
                (r = !0);
            };
          var w = {
            open: function () {
              r || d(), document.addEventListener("click", u), (c = !0);
              var e = t.get("state");
              e === a.pb && n.pause("infoOverlayInteraction"), (s = e), i(!0);
            },
            close: function () {
              document.removeEventListener("click", u),
                (c = !1),
                t.get("state") === a.ob &&
                  s === a.pb &&
                  n.play("infoOverlayInteraction"),
                (s = null),
                i(!1);
            },
            destroy: function () {
              this.close(), t.off(null, null, this);
            },
          };
          return (
            Object.defineProperties(w, {
              visible: {
                enumerable: !0,
                get: function () {
                  return c;
                },
              },
            }),
            w
          );
        };
      var vn = function (e, t, n) {
          var i,
            o = !1,
            r = null,
            s = n.get("localization").shortcuts,
            c = Object(l.e)(
              (function (e, t) {
                var n = e
                  .map(function (e) {
                    return (
                      '<div class="jw-shortcuts-row jw-reset">' +
                      '<span class="jw-shortcuts-description jw-reset">'.concat(
                        e.description,
                        "</span>"
                      ) +
                      '<span class="jw-shortcuts-key jw-reset">'.concat(
                        e.key,
                        "</span>"
                      ) +
                      "</div>"
                    );
                  })
                  .join("");
                return (
                  '<div class="jw-shortcuts-tooltip jw-modal jw-reset" title="'.concat(
                    t,
                    '">'
                  ) +
                  '<span class="jw-hidden" id="jw-shortcuts-tooltip-explanation">Press shift question mark to access a list of keyboard shortcuts</span><div class="jw-reset jw-shortcuts-container"><div class="jw-reset jw-shortcuts-header">' +
                  '<span class="jw-reset jw-shortcuts-title">'.concat(
                    t,
                    "</span>"
                  ) +
                  '<button role="switch" class="jw-reset jw-switch" data-jw-switch-enabled="Enabled" data-jw-switch-disabled="Disabled"><span class="jw-reset jw-switch-knob"></span></button></div><div class="jw-reset jw-shortcuts-tooltip-list"><div class="jw-shortcuts-tooltip-descriptions jw-reset">' +
                  "".concat(n) +
                  "</div></div></div></div>"
                );
              })(
                (function (e) {
                  var t = e.playPause,
                    n = e.volumeToggle,
                    i = e.fullscreenToggle,
                    o = e.seekPercent,
                    a = e.increaseVolume,
                    r = e.decreaseVolume,
                    s = e.seekForward,
                    l = e.seekBackward;
                  return [
                    { key: e.spacebar, description: t },
                    { key: "↑", description: a },
                    { key: "↓", description: r },
                    { key: "→", description: s },
                    { key: "←", description: l },
                    { key: "c", description: e.captionsToggle },
                    { key: "f", description: i },
                    { key: "m", description: n },
                    { key: "0-9", description: o },
                  ];
                })(s),
                s.keyboardShortcuts
              )
            ),
            d = { reason: "settingsInteraction" },
            w = new u.a(c.querySelector(".jw-switch")),
            f = function () {
              w.el.setAttribute("aria-checked", n.get("enableShortcuts")),
                Object(l.a)(c, "jw-open"),
                (r = n.get("state")),
                c.querySelector(".jw-shortcuts-close").focus(),
                document.addEventListener("click", j),
                (o = !0),
                t.pause(d);
            },
            h = function () {
              Object(l.o)(c, "jw-open"),
                document.removeEventListener("click", j),
                e.focus(),
                (o = !1),
                r === a.pb && t.play(d);
            },
            j = function (e) {
              /jw-shortcuts|jw-switch/.test(e.target.className) || h();
            },
            g = function (e) {
              var t = e.currentTarget,
                i = "true" !== t.getAttribute("aria-checked");
              t.setAttribute("aria-checked", i), n.set("enableShortcuts", i);
            };
          return (
            (i = p("jw-shortcuts-close", h, n.get("localization").close, [
              de("close"),
            ])),
            Object(l.m)(c, i.element()),
            i.show(),
            e.appendChild(c),
            w.on("click tap enter", g),
            {
              el: c,
              open: f,
              close: h,
              destroy: function () {
                h(), w.destroy();
              },
              toggleVisibility: function () {
                o ? h() : f();
              },
            }
          );
        },
        yn = function (e) {
          return (
            '<div class="jw-float-icon jw-icon jw-button-color jw-reset" aria-label='.concat(
              e,
              ' tabindex="0">'
            ) + "</div>"
          );
        };
      function kn(e) {
        return (kn =
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
      function xn(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function On(e, t) {
        return !t || ("object" !== kn(t) && "function" != typeof t)
          ? (function (e) {
              if (void 0 === e)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return e;
            })(e)
          : t;
      }
      function Cn(e) {
        return (Cn = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function Mn(e, t) {
        return (Mn =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var Sn = (function (e) {
        function t(e, n) {
          var i;
          return (
            (function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
            ((i = On(this, Cn(t).call(this))).element = Object(l.e)(yn(n))),
            i.element.appendChild(de("close")),
            (i.ui = new u.a(i.element, { directSelect: !0 }).on(
              "click tap enter",
              function () {
                i.trigger(a.sb);
              }
            )),
            e.appendChild(i.element),
            i
          );
        }
        var n, i, o;
        return (
          (function (e, t) {
            if ("function" != typeof t && null !== t)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (e.prototype = Object.create(t && t.prototype, {
              constructor: { value: e, writable: !0, configurable: !0 },
            })),
              t && Mn(e, t);
          })(t, e),
          (n = t),
          (i = [
            {
              key: "destroy",
              value: function () {
                this.element &&
                  (this.ui.destroy(),
                  this.element.parentNode.removeChild(this.element),
                  (this.element = null));
              },
            },
          ]) && xn(n.prototype, i),
          o && xn(n, o),
          t
        );
      })(r.a);
      function Tn(e) {
        return (Tn =
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
      function En(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function _n(e, t) {
        return !t || ("object" !== Tn(t) && "function" != typeof t)
          ? (function (e) {
              if (void 0 === e)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return e;
            })(e)
          : t;
      }
      function zn(e) {
        return (zn = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function Pn(e, t) {
        return (Pn =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      n.d(t, "default", function () {
        return In;
      }),
        n(95);
      var An = o.OS.mobile ? 4e3 : 2e3,
        Ln = [27];
      (gn.a.cloneIcon = de),
        bn.a.forEach(function (e) {
          if (e.getState() === a.lb) {
            var t = e.getContainer().querySelector(".jw-error-msg .jw-icon");
            t && !t.hasChildNodes() && t.appendChild(gn.a.cloneIcon("error"));
          }
        });
      var Rn = function () {
          return { reason: "interaction" };
        },
        In = (function (e) {
          function t(e, n) {
            var i;
            return (
              (function (e, t) {
                if (!(e instanceof t))
                  throw new TypeError("Cannot call a class as a function");
              })(this, t),
              ((i = _n(this, zn(t).call(this))).activeTimeout = -1),
              (i.inactiveTime = 0),
              (i.context = e),
              (i.controlbar = null),
              (i.displayContainer = null),
              (i.backdrop = null),
              (i.enabled = !0),
              (i.instreamState = null),
              (i.keydownCallback = null),
              (i.keyupCallback = null),
              (i.blurCallback = null),
              (i.mute = null),
              (i.nextUpToolTip = null),
              (i.playerContainer = n),
              (i.wrapperElement = n.querySelector(".jw-wrapper")),
              (i.rightClickMenu = null),
              (i.settingsMenu = null),
              (i.shortcutsTooltip = null),
              (i.showing = !1),
              (i.muteChangeCallback = null),
              (i.unmuteCallback = null),
              (i.logo = null),
              (i.div = null),
              (i.dimensions = {}),
              (i.infoOverlay = null),
              (i.userInactiveTimeout = function () {
                var e = i.inactiveTime - Object(c.a)();
                i.inactiveTime && e > 16
                  ? (i.activeTimeout = setTimeout(i.userInactiveTimeout, e))
                  : i.playerContainer.querySelector(".jw-tab-focus")
                  ? i.resetActiveTimeout()
                  : i.userInactive();
              }),
              i
            );
          }
          var n, i, r;
          return (
            (function (e, t) {
              if ("function" != typeof t && null !== t)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (e.prototype = Object.create(t && t.prototype, {
                constructor: { value: e, writable: !0, configurable: !0 },
              })),
                t && Pn(e, t);
            })(t, e),
            (n = t),
            (i = [
              {
                key: "resetActiveTimeout",
                value: function () {
                  clearTimeout(this.activeTimeout),
                    (this.activeTimeout = -1),
                    (this.inactiveTime = 0);
                },
              },
              {
                key: "enable",
                value: function (e, t) {
                  var n = this,
                    i = this.context.createElement("div");
                  (i.className = "jw-controls jw-reset"), (this.div = i);
                  var r = this.context.createElement("div");
                  (r.className = "jw-controls-backdrop jw-reset"),
                    (this.backdrop = r),
                    (this.logo = this.playerContainer.querySelector(
                      ".jw-logo"
                    ));
                  var c = t.get("touchMode"),
                    u = function () {
                      (t.get("isFloating")
                        ? n.wrapperElement
                        : n.playerContainer
                      ).focus();
                    };
                  if (!this.displayContainer) {
                    var d = new Ct(t, e);
                    d.buttons.display.on("click tap enter", function () {
                      n.trigger(a.p),
                        n.userActive(1e3),
                        e.playToggle(Rn()),
                        u();
                    }),
                      this.div.appendChild(d.element()),
                      (this.displayContainer = d);
                  }
                  (this.infoOverlay = new mn(i, t, e, function (e) {
                    Object(l.v)(n.div, "jw-info-open", e),
                      e && n.div.querySelector(".jw-info-close").focus();
                  })),
                    o.OS.mobile ||
                      (this.shortcutsTooltip = new vn(
                        this.wrapperElement,
                        e,
                        t
                      )),
                    (this.rightClickMenu = new Vt(
                      this.infoOverlay,
                      this.shortcutsTooltip
                    )),
                    c
                      ? (Object(l.a)(this.playerContainer, "jw-flag-touch"),
                        this.rightClickMenu.setup(
                          t,
                          this.playerContainer,
                          this.wrapperElement
                        ))
                      : t.change(
                          "flashBlocked",
                          function (e, t) {
                            t
                              ? n.rightClickMenu.destroy()
                              : n.rightClickMenu.setup(
                                  e,
                                  n.playerContainer,
                                  n.wrapperElement
                                );
                          },
                          this
                        );
                  var w = t.get("floating");
                  if (w) {
                    var f = new Sn(i, t.get("localization").close);
                    f.on(a.sb, function () {
                      return n.trigger("dismissFloating", { doNotForward: !0 });
                    }),
                      !1 !== w.dismissible &&
                        Object(l.a)(
                          this.playerContainer,
                          "jw-floating-dismissible"
                        );
                  }
                  var h = (this.controlbar = new dt(
                    e,
                    t,
                    this.playerContainer.querySelector(
                      ".jw-hidden-accessibility"
                    )
                  ));
                  if (
                    (h.on(a.sb, function () {
                      return n.userActive();
                    }),
                    h.on(
                      "nextShown",
                      function (e) {
                        this.trigger("nextShown", e);
                      },
                      this
                    ),
                    h.on("adjustVolume", k, this),
                    t.get("nextUpDisplay") && !h.nextUpToolTip)
                  ) {
                    var j = new Et(t, e, this.playerContainer);
                    j.on("all", this.trigger, this),
                      j.setup(this.context),
                      (h.nextUpToolTip = j),
                      this.div.appendChild(j.element());
                  }
                  this.div.appendChild(h.element());
                  var g = t.get("localization"),
                    b = (this.settingsMenu = hn(
                      e,
                      t.player,
                      this.controlbar,
                      g
                    )),
                    m = null;
                  this.controlbar.on("menuVisibility", function (i) {
                    var o = i.visible,
                      r = i.evt,
                      s = t.get("state"),
                      l = { reason: "settingsInteraction" },
                      c = n.controlbar.elements.settingsButton,
                      d = "keydown" === ((r && r.sourceEvent) || r || {}).type,
                      p = o || d ? 0 : An;
                    n.userActive(p),
                      (m = s),
                      Object(jn.a)(t.get("containerWidth")) < 2 &&
                        (o && s === a.pb
                          ? e.pause(l)
                          : o || s !== a.ob || m !== a.pb || e.play(l)),
                      !o && d && c ? c.element().focus() : r && u();
                  }),
                    b.on("menuVisibility", function (e) {
                      return n.controlbar.trigger("menuVisibility", e);
                    }),
                    this.controlbar.on(
                      "settingsInteraction",
                      function (e, t, n) {
                        if (t) return b.defaultChild.toggle(n);
                        b.children[e].toggle(n);
                      }
                    ),
                    o.OS.mobile
                      ? this.div.appendChild(b.el)
                      : (this.playerContainer.setAttribute(
                          "aria-describedby",
                          "jw-shortcuts-tooltip-explanation"
                        ),
                        this.div.insertBefore(b.el, h.element()));
                  var v = function (t) {
                    if (t.get("autostartMuted")) {
                      var i = function () {
                          return n.unmuteAutoplay(e, t);
                        },
                        a = function (e, t) {
                          t || i();
                        };
                      o.OS.mobile &&
                        ((n.mute = p(
                          "jw-autostart-mute jw-off",
                          i,
                          t.get("localization").unmute,
                          [de("volume-0")]
                        )),
                        n.mute.show(),
                        n.div.appendChild(n.mute.element())),
                        h.renderVolume(!0, t.get("volume")),
                        Object(l.a)(n.playerContainer, "jw-flag-autostart"),
                        t.on("change:autostartFailed", i, n),
                        t.on("change:autostartMuted change:mute", a, n),
                        (n.muteChangeCallback = a),
                        (n.unmuteCallback = i);
                    }
                  };
                  function y(n) {
                    var i = 0,
                      o = t.get("duration"),
                      a = t.get("position");
                    if ("DVR" === t.get("streamType")) {
                      var r = t.get("dvrSeekLimit");
                      (i = o), (o = Math.max(a, -r));
                    }
                    var l = Object(s.a)(a + n, i, o);
                    e.seek(l, Rn());
                  }
                  function k(n) {
                    var i = Object(s.a)(t.get("volume") + n, 0, 100);
                    e.setVolume(i);
                  }
                  t.once("change:autostartMuted", v), v(t);
                  var x = function (i) {
                    if (i.ctrlKey || i.metaKey) return !0;
                    var o = !n.settingsMenu.visible,
                      a = !0 === t.get("enableShortcuts"),
                      r = n.instreamState;
                    if (a || -1 !== Ln.indexOf(i.keyCode)) {
                      switch (i.keyCode) {
                        case 27:
                          if (t.get("fullscreen"))
                            e.setFullscreen(!1),
                              n.playerContainer.blur(),
                              n.userInactive();
                          else {
                            var s = e.getPlugin("related");
                            s && s.close({ type: "escape" });
                          }
                          n.rightClickMenu.el &&
                            n.rightClickMenu.hideMenuHandler(),
                            n.infoOverlay.visible && n.infoOverlay.close(),
                            n.shortcutsTooltip && n.shortcutsTooltip.close();
                          break;
                        case 13:
                        case 32:
                          if (
                            document.activeElement.classList.contains(
                              "jw-switch"
                            ) &&
                            13 === i.keyCode
                          )
                            return !0;
                          e.playToggle(Rn());
                          break;
                        case 37:
                          !r && o && y(-5);
                          break;
                        case 39:
                          !r && o && y(5);
                          break;
                        case 38:
                          o && k(10);
                          break;
                        case 40:
                          o && k(-10);
                          break;
                        case 67:
                          var l = e.getCaptionsList().length;
                          if (l) {
                            var c = (e.getCurrentCaptions() + 1) % l;
                            e.setCurrentCaptions(c);
                          }
                          break;
                        case 77:
                          e.setMute();
                          break;
                        case 70:
                          e.setFullscreen();
                          break;
                        case 191:
                          n.shortcutsTooltip &&
                            n.shortcutsTooltip.toggleVisibility();
                          break;
                        default:
                          if (i.keyCode >= 48 && i.keyCode <= 59) {
                            var u = ((i.keyCode - 48) / 10) * t.get("duration");
                            e.seek(u, Rn());
                          }
                      }
                      return /13|32|37|38|39|40/.test(i.keyCode)
                        ? (i.preventDefault(), !1)
                        : void 0;
                    }
                  };
                  this.playerContainer.addEventListener("keydown", x),
                    (this.keydownCallback = x);
                  var O = function (e) {
                    switch (e.keyCode) {
                      case 9:
                        var t = n.playerContainer.contains(e.target) ? 0 : An;
                        n.userActive(t);
                        break;
                      case 32:
                        e.preventDefault();
                    }
                  };
                  this.playerContainer.addEventListener("keyup", O),
                    (this.keyupCallback = O);
                  var C = function (e) {
                    var t = e.relatedTarget || document.querySelector(":focus");
                    t && (n.playerContainer.contains(t) || n.userInactive());
                  };
                  this.playerContainer.addEventListener("blur", C, !0),
                    (this.blurCallback = C);
                  var M = function e() {
                    "jw-shortcuts-tooltip-explanation" ===
                      n.playerContainer.getAttribute("aria-describedby") &&
                      n.playerContainer.removeAttribute("aria-describedby"),
                      n.playerContainer.removeEventListener("blur", e, !0);
                  };
                  this.shortcutsTooltip &&
                    (this.playerContainer.addEventListener("blur", M, !0),
                    (this.onRemoveShortcutsDescription = M)),
                    this.userActive(),
                    this.addControls(),
                    this.addBackdrop(),
                    t.set("controlsEnabled", !0);
                },
              },
              {
                key: "addControls",
                value: function () {
                  this.wrapperElement.appendChild(this.div);
                },
              },
              {
                key: "disable",
                value: function (e) {
                  var t = this.nextUpToolTip,
                    n = this.settingsMenu,
                    i = this.infoOverlay,
                    o = this.controlbar,
                    a = this.rightClickMenu,
                    r = this.shortcutsTooltip,
                    s = this.playerContainer,
                    c = this.div;
                  clearTimeout(this.activeTimeout),
                    (this.activeTimeout = -1),
                    this.off(),
                    e.off(null, null, this),
                    e.set("controlsEnabled", !1),
                    c.parentNode &&
                      (Object(l.o)(s, "jw-flag-touch"),
                      c.parentNode.removeChild(c)),
                    o && o.destroy(),
                    a && a.destroy(),
                    this.keydownCallback &&
                      s.removeEventListener("keydown", this.keydownCallback),
                    this.keyupCallback &&
                      s.removeEventListener("keyup", this.keyupCallback),
                    this.blurCallback &&
                      s.removeEventListener("blur", this.blurCallback),
                    this.onRemoveShortcutsDescription &&
                      s.removeEventListener(
                        "blur",
                        this.onRemoveShortcutsDescription
                      ),
                    this.displayContainer && this.displayContainer.destroy(),
                    t && t.destroy(),
                    n && n.destroy(),
                    i && i.destroy(),
                    r && r.destroy(),
                    this.removeBackdrop();
                },
              },
              {
                key: "controlbarHeight",
                value: function () {
                  return (
                    this.dimensions.cbHeight ||
                      (this.dimensions.cbHeight = this.controlbar.element().clientHeight),
                    this.dimensions.cbHeight
                  );
                },
              },
              {
                key: "element",
                value: function () {
                  return this.div;
                },
              },
              {
                key: "resize",
                value: function () {
                  this.dimensions = {};
                },
              },
              {
                key: "unmuteAutoplay",
                value: function (e, t) {
                  var n = !t.get("autostartFailed"),
                    i = t.get("mute");
                  n ? (i = !1) : t.set("playOnViewable", !1),
                    this.muteChangeCallback &&
                      (t.off(
                        "change:autostartMuted change:mute",
                        this.muteChangeCallback
                      ),
                      (this.muteChangeCallback = null)),
                    this.unmuteCallback &&
                      (t.off("change:autostartFailed", this.unmuteCallback),
                      (this.unmuteCallback = null)),
                    t.set("autostartFailed", void 0),
                    t.set("autostartMuted", void 0),
                    e.setMute(i),
                    this.controlbar.renderVolume(i, t.get("volume")),
                    this.mute && this.mute.hide(),
                    Object(l.o)(this.playerContainer, "jw-flag-autostart"),
                    this.userActive();
                },
              },
              {
                key: "mouseMove",
                value: function (e) {
                  var t = this.controlbar.element().contains(e.target),
                    n =
                      this.controlbar.nextUpToolTip &&
                      this.controlbar.nextUpToolTip
                        .element()
                        .contains(e.target),
                    i = this.logo && this.logo.contains(e.target),
                    o = t || n || i ? 0 : An;
                  this.userActive(o);
                },
              },
              {
                key: "userActive",
                value: function () {
                  var e =
                    arguments.length > 0 && void 0 !== arguments[0]
                      ? arguments[0]
                      : An;
                  e > 0
                    ? ((this.inactiveTime = Object(c.a)() + e),
                      -1 === this.activeTimeout &&
                        (this.activeTimeout = setTimeout(
                          this.userInactiveTimeout,
                          e
                        )))
                    : this.resetActiveTimeout(),
                    this.showing ||
                      (Object(l.o)(
                        this.playerContainer,
                        "jw-flag-user-inactive"
                      ),
                      (this.showing = !0),
                      this.trigger("userActive"));
                },
              },
              {
                key: "userInactive",
                value: function () {
                  clearTimeout(this.activeTimeout),
                    (this.activeTimeout = -1),
                    this.settingsMenu.visible ||
                      ((this.inactiveTime = 0),
                      (this.showing = !1),
                      Object(l.a)(
                        this.playerContainer,
                        "jw-flag-user-inactive"
                      ),
                      this.trigger("userInactive"));
                },
              },
              {
                key: "addBackdrop",
                value: function () {
                  var e = this.instreamState
                    ? this.div
                    : this.wrapperElement.querySelector(".jw-captions");
                  this.wrapperElement.insertBefore(this.backdrop, e);
                },
              },
              {
                key: "removeBackdrop",
                value: function () {
                  var e = this.backdrop.parentNode;
                  e && e.removeChild(this.backdrop);
                },
              },
              {
                key: "setupInstream",
                value: function () {
                  (this.instreamState = !0),
                    this.userActive(),
                    this.addBackdrop(),
                    this.settingsMenu && this.settingsMenu.close(),
                    Object(l.o)(this.playerContainer, "jw-flag-autostart"),
                    this.controlbar.elements.time
                      .element()
                      .setAttribute("tabindex", "-1");
                },
              },
              {
                key: "destroyInstream",
                value: function (e) {
                  (this.instreamState = null),
                    this.addBackdrop(),
                    e.get("autostartMuted") &&
                      Object(l.a)(this.playerContainer, "jw-flag-autostart"),
                    this.controlbar.elements.time
                      .element()
                      .setAttribute("tabindex", "0");
                },
              },
            ]) && En(n.prototype, i),
            r && En(n, r),
            t
          );
        })(r.a);
    },
    function (e, t, n) {
      "use strict";
      n.r(t);
      var i = n(0),
        o = n(12),
        a = n(50),
        r = n(36);
      var s = n(44),
        l = n(51),
        c = n(26),
        u = n(25),
        d = n(3),
        p = n(46),
        w = n(2),
        f = n(7),
        h = n(34);
      function j(e) {
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
      var g = n(1);
      function b(e) {
        return function (t, n) {
          var o = e.mediaModel,
            a = Object(i.g)({}, n, { type: t });
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
              var r = o.attributes.mediaState;
              (o.attributes.mediaState = n.newstate),
                o.trigger("change:mediaState", o, n.newstate, r);
              break;
            case d.F:
              return (
                (e.beforeComplete = !0),
                e.trigger(d.B, a),
                void (e.attached && !e.background && e._playbackComplete())
              );
            case d.G:
              o.get("setup")
                ? (e.thenPlayPromise.cancel(), o.srcReset())
                : ((t = d.tb), (a.code += 1e5));
              break;
            case d.K:
              a.metadataType || (a.metadataType = "unknown");
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
                p = n.levels;
              u > -1 && p.length > 1 && o.set("currentLevel", parseInt(u));
              break;
            case d.f:
              o.set(d.f, n.tracks);
            case d.g:
              var w = n.currentTrack,
                f = n.tracks;
              w > -1 &&
                f.length > 0 &&
                w < f.length &&
                o.set("currentAudioTrack", parseInt(w));
          }
          e.trigger(t, a);
        };
      }
      var m = n(8),
        v = n(45),
        y = n(41);
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
      function x(e, t) {
        if (!(e instanceof t))
          throw new TypeError("Cannot call a class as a function");
      }
      function O(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function C(e, t, n) {
        return t && O(e.prototype, t), n && O(e, n), e;
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
      function S(e) {
        return (S = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function T(e, t) {
        if ("function" != typeof t && null !== t)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (e.prototype = Object.create(t && t.prototype, {
          constructor: { value: e, writable: !0, configurable: !0 },
        })),
          t && E(e, t);
      }
      function E(e, t) {
        return (E =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var _ = (function (e) {
          function t() {
            var e;
            return (
              x(this, t),
              ((e = M(this, S(t).call(this))).providerController = null),
              (e._provider = null),
              e.addAttributes({ mediaModel: new P() }),
              e
            );
          }
          return (
            T(t, e),
            C(t, [
              {
                key: "setup",
                value: function (e) {
                  return (
                    (e = e || {}),
                    this._normalizeConfig(e),
                    Object(i.g)(this.attributes, e, y.b),
                    (this.providerController = new h.a(
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
                    Object.keys(y.a).forEach(function (n) {
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
                    a = Object(i.u)(n.bitrate) ? n.bitrate : null;
                  this.set("bitrateSelection", a), this.set("qualityLabel", o);
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
                    (e = e || new P()),
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
                  (this._provider = e), z(this, e);
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
                  var t = m.OS.mobile && this.get("autostart");
                  this.set(
                    "playOnViewable",
                    t || "viewable" === this.get("autostart")
                  );
                },
              },
              {
                key: "resetItem",
                value: function (e) {
                  var t = e ? Object(w.g)(e.starttime) : 0,
                    n = e ? Object(w.g)(e.duration) : 0,
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
        })(v.a),
        z = function (e, t) {
          e.set("provider", t.getName()),
            !0 === e.get("instreamMode") && (t.instreamMode = !0),
            -1 === t.getName().name.indexOf("flash") &&
              (e.set("flashThrottle", void 0), e.set("flashBlocked", !1)),
            e.setPlaybackRate(e.get("defaultPlaybackRate")),
            e.set("supportsPlaybackRate", t.supportsPlaybackRate),
            e.set("playbackRate", t.getPlaybackRate()),
            e.set("renderCaptionsNatively", t.renderNatively);
        };
      var P = (function (e) {
          function t() {
            var e;
            return (
              x(this, t),
              (e = M(this, S(t).call(this))).addAttributes({
                mediaState: d.mb,
              }),
              e
            );
          }
          return (
            T(t, e),
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
        })(v.a),
        A = _;
      function L(e) {
        return (L =
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
      function R(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function I(e) {
        return (I = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (e) {
              return e.__proto__ || Object.getPrototypeOf(e);
            })(e);
      }
      function B(e, t) {
        return (B =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      function V(e) {
        if (void 0 === e)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return e;
      }
      var N = (function (e) {
        function t(e, n) {
          var i, o, a, r;
          return (
            (function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
            (o = this),
            (a = I(t).call(this)),
            ((i =
              !a || ("object" !== L(a) && "function" != typeof a)
                ? V(o)
                : a).attached = !0),
            (i.beforeComplete = !1),
            (i.item = null),
            (i.mediaModel = new P()),
            (i.model = n),
            (i.provider = e),
            (i.providerListener = new b(V(V(i)))),
            (i.thenPlayPromise = j(function () {})),
            (r = V(V(i))).provider.on("all", r.providerListener, r),
            (i.eventQueue = new s.a(V(V(i)), ["trigger"], function () {
              return !i.attached || i.background;
            })),
            i
          );
        }
        var n, o, a;
        return (
          (function (e, t) {
            if ("function" != typeof t && null !== t)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (e.prototype = Object.create(t && t.prototype, {
              constructor: { value: e, writable: !0, configurable: !0 },
            })),
              t && B(e, t);
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
                var a = this._loadAndPlay(t, o);
                return i.get("started") ? a : this._playAttempt(a, e);
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
                  a = this.mediaModel,
                  r = this.model,
                  s = this.provider,
                  l = s ? s.video : null;
                return (
                  this.trigger(d.N, { item: o, playReason: t }),
                  (l ? l.paused : r.get(d.bb) !== d.pb) || r.set(d.bb, d.jb),
                  e
                    .then(function () {
                      a.get("setup") &&
                        (a.set("started", !0),
                        a === r.mediaModel &&
                          (function (e) {
                            var t = e.get("mediaState");
                            e.trigger("change:mediaState", e, t, t);
                          })(a));
                    })
                    .catch(function (e) {
                      if (n.item && a === r.mediaModel) {
                        if ((r.set("playRejected", !0), l && l.paused)) {
                          if (l.src === location.href)
                            return n._loadAndPlay(o, s);
                          a.set("mediaState", d.ob);
                        }
                        var c = Object(i.g)(new g.n(null, Object(g.w)(e), e), {
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
                  var i = j(function () {
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
                var t = (this.mediaModel = new P()),
                  n = e ? Object(w.g)(e.starttime) : 0,
                  i = e ? Object(w.g)(e.duration) : 0,
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
          ]) && R(n.prototype, o),
          a && R(n, a),
          t
        );
      })(f.a);
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
      function F(e, t) {
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
      function U(e) {
        if (void 0 === e)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return e;
      }
      function W(e, t) {
        var n = t.mediaControllerListener;
        e.off().on("all", n, t);
      }
      function Q(e) {
        return e && e.sources && e.sources[0];
      }
      var Y = (function (e) {
        function t(e, n) {
          var o, a, r, s, l;
          return (
            (function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
            (a = this),
            ((o =
              !(r = q(t).call(this)) ||
              ("object" !== H(r) && "function" != typeof r)
                ? U(a)
                : r).adPlaying = !1),
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
            })(e, U(U(o)))),
            (o.model = e),
            (o.providers = new h.a(e.getConfiguration())),
            (o.loadPromise = Promise.resolve()),
            (o.backgroundLoading = e.get("backgroundLoading")),
            o.backgroundLoading ||
              e.set("mediaElement", o.mediaPool.getPrimedElement()),
            o
          );
        }
        var n, o, a;
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
                if (!o) return Promise.reject(new g.n(g.k, g.h));
                var a = this.background,
                  r = this.mediaController;
                if (a.isNext(i))
                  return (
                    this._destroyActiveMedia(),
                    (this.loadPromise = this._activateBackgroundMedia()),
                    this.loadPromise
                  );
                if ((this._destroyBackgroundMedia(), r)) {
                  if (
                    n.get("castActive") ||
                    this._providerCanPlay(r.provider, o)
                  )
                    return (
                      (this.loadPromise = Promise.resolve(r)),
                      (r.activeItem = i),
                      this._setActiveMedia(r),
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
                  var a = j(function (t) {
                    if (
                      n.mediaController &&
                      n.mediaController.mediaModel === t.mediaModel
                    )
                      return t.play(e);
                    throw new Error("Playback cancelled.");
                  });
                  t = this.loadPromise
                    .catch(function (e) {
                      throw (a.cancel(), e);
                    })
                    .then(a.async);
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
                  a = (o.starttime = n.mediaModel.get("currentTime"));
                this._destroyActiveMedia();
                var r = new N(e, n);
                (r.activeItem = o),
                  this._setActiveMedia(r),
                  n.mediaModel.set("currentTime", a);
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
                e && W(e, this);
              },
            },
            {
              key: "routeEvents",
              value: function (e) {
                var t = this.mediaController;
                t && (t.off(), e && W(t, e));
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
                  W(e, this),
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
                  a = i.choose(e),
                  r = a.provider,
                  s = a.name;
                return r
                  ? Promise.resolve(o(r))
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
          ]) && F(n.prototype, o),
          a && F(n, a),
          t
        );
      })(f.a);
      function X(e) {
        return e === d.kb || e === d.lb ? d.mb : e;
      }
      function Z(e, t, n) {
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
      var K = n(48);
      function J(e) {
        return (J =
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
      function G(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function $(e, t) {
        return !t || ("object" !== J(t) && "function" != typeof t)
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
                  a = ie(e, t);
                if (a) {
                  if ((o = Object.getOwnPropertyDescriptor(a, t)).set)
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
      function ae(e, t) {
        return (ae =
          Object.setPrototypeOf ||
          function (e, t) {
            return (e.__proto__ = t), e;
          })(e, t);
      }
      var re = (function (e) {
          function t(e, n) {
            var i;
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, t);
            var o,
              a = ((i = $(this, oe(t).call(this, e, n))).model = new A());
            if (
              ((i.playerModel = e),
              (i.provider = null),
              (i.backgroundLoading = e.get("backgroundLoading")),
              (a.mediaModel.attributes.mediaType = "video"),
              i.backgroundLoading)
            )
              o = n.getAdElement();
            else {
              (o = e.get("mediaElement")),
                (a.attributes.mediaElement = o),
                (a.attributes.mediaSrc = o.src);
              var r = (i.srcResetListener = function () {
                i.srcReset();
              });
              o.addEventListener("emptied", r),
                (o.playbackRate = o.defaultPlaybackRate = 1);
            }
            return (i.mediaPool = Object(K.a)(o, n)), i;
          }
          var n, o, a;
          return (
            (function (e, t) {
              if ("function" != typeof t && null !== t)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (e.prototype = Object.create(t && t.prototype, {
                constructor: { value: e, writable: !0, configurable: !0 },
              })),
                t && ae(e, t);
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
                    e.on("change:state", Z, this),
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
                      a = "vpaid" === e.type;
                    e.off(),
                      e.on(
                        "all",
                        function (e, t) {
                          (a && e === d.F) ||
                            this.trigger(e, Object(i.g)({}, t, { type: e }));
                        },
                        this
                      );
                    var r = n.mediaModel;
                    e.on(d.bb, function (e) {
                      (e.oldstate = e.oldstate || n.get(d.bb)),
                        r.set("mediaState", e.newstate);
                    }),
                      e.on(d.X, this._nativeFullscreenHandler, this),
                      r.on("change:mediaState", function (e, n) {
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
            ]) && G(n.prototype, o),
            a && G(n, a),
            t
          );
        })(Y),
        se = { skipoffset: null, tag: null },
        le = function (e, t, n, o) {
          var a,
            r,
            s,
            l,
            c = this,
            u = this,
            f = new re(t, o),
            h = 0,
            j = {},
            g = null,
            b = {},
            m = A,
            v = !1,
            y = !1,
            k = !1,
            x = !1,
            O = function (e) {
              y ||
                (((e = e || {}).hasControls = !!t.get("controls")),
                c.trigger(d.z, e),
                f.model.get("state") === d.ob
                  ? e.hasControls && f.playVideo().catch(function () {})
                  : f.pause());
            },
            C = function () {
              y ||
                (f.model.get("state") === d.ob &&
                  t.get("controls") &&
                  (e.setFullscreen(), e.play()));
            };
          function M() {
            f.model.set("playRejected", !0);
          }
          function S() {
            h++, u.loadItem(a).catch(function () {});
          }
          function T(e, t) {
            "complete" !== e &&
              ((t = t || {}),
              b.tag && !t.tag && (t.tag = b.tag),
              this.trigger(e, t),
              ("mediaError" !== e && "error" !== e) ||
                (a && h + 1 < a.length && S()));
          }
          function E(e) {
            var t = e.newstate,
              n = e.oldstate || f.model.get("state");
            n !== t && _(Object(i.g)({ oldstate: n }, j, e));
          }
          function _(t) {
            var n = t.newstate;
            n === d.pb ? e.trigger(d.c, t) : n === d.ob && e.trigger(d.b, t);
          }
          function z(t) {
            var n = t.duration,
              i = t.position,
              o = f.model.mediaModel || f.model;
            o.set("duration", n),
              o.set("position", i),
              l || (l = (Object(w.d)(s, n) || n) - p.b),
              !v && i >= Math.max(l, p.a) && (e.preloadNextItem(), (v = !0));
          }
          function P(e) {
            var t = {};
            b.tag && (t.tag = b.tag), this.trigger(d.F, t), A.call(this, e);
          }
          function A(e) {
            (j = {}),
              a && h + 1 < a.length
                ? S()
                : (e.type === d.F && this.trigger(d.cb, {}), this.destroy());
          }
          function L() {
            y ||
              (n.clickHandler() &&
                n.clickHandler().setAlternateClickHandlers(O, C));
          }
          function R(e) {
            e.width && e.height && n.resizeMedia();
          }
          (this.init = function () {
            if (!k && !y) {
              (k = !0),
                (j = {}),
                f.setup(),
                f.on("all", T, this),
                f.on(d.O, M, this),
                f.on(d.S, z, this),
                f.on(d.F, P, this),
                f.on(d.K, R, this),
                f.on(d.bb, E, this),
                (g = e.detachMedia());
              var i = f.primedElement;
              t.get("mediaContainer").appendChild(i),
                t.set("instream", f),
                f.model.set("state", d.jb);
              var o = n.clickHandler();
              return (
                o && o.setAlternateClickHandlers(function () {}, null),
                this.setText(t.get("localization").loadingAd),
                (x = e.isBeforeComplete() || t.get("state") === d.kb),
                this
              );
            }
          }),
            (this.enableAdsMode = function (i) {
              var o = this;
              if (!k && !y)
                return (
                  e.routeEvents({
                    mediaControllerListener: function (e, t) {
                      o.trigger(e, t);
                    },
                  }),
                  t.set("instream", f),
                  f.model.set("state", d.pb),
                  (function (i) {
                    var o = n.clickHandler();
                    o &&
                      o.setAlternateClickHandlers(function (n) {
                        y ||
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
              j = e;
            }),
            (this.setState = function (e) {
              var t = e.newstate,
                n = f.model;
              (e.oldstate = n.get("state")), n.set("state", t), _(e);
            }),
            (this.setTime = function (t) {
              z(t), e.trigger(d.e, t);
            }),
            (this.loadItem = function (e, n) {
              if (y || !k)
                return Promise.reject(new Error("Instream not setup"));
              j = {};
              var o = e;
              Array.isArray(e)
                ? ((r = n || r), (e = (a = e)[h]), r && (n = r[h]))
                : (o = [e]);
              var l = f.model;
              l.set("playlist", o),
                t.set("hideAdsControls", !1),
                (e.starttime = 0),
                u.trigger(d.db, { index: h, item: e }),
                (b = Object(i.g)({}, se, n)),
                L(),
                l.set("skipButton", !1);
              var c =
                !t.get("backgroundLoading") && g
                  ? g.then(function () {
                      return f.setActiveItem(h);
                    })
                  : f.setActiveItem(h);
              return (
                (v = !1),
                void 0 !== (s = e.skipoffset || b.skipoffset) &&
                  u.setupSkipButton(s, b),
                c
              );
            }),
            (this.setupSkipButton = function (e, t, n) {
              var i = f.model;
              (m = n || A),
                i.set("skipMessage", t.skipMessage),
                i.set("skipText", t.skipText),
                i.set("skipOffset", e),
                (i.attributes.skipButton = !1),
                i.set("skipButton", !0);
            }),
            (this.applyProviderListeners = function (e) {
              f.usePsuedoProvider(e), L();
            }),
            (this.play = function () {
              (j = {}), f.playVideo();
            }),
            (this.pause = function () {
              (j = {}), f.pause();
            }),
            (this.skipAd = function (e) {
              var n = t.get("autoPause").pauseAds,
                i = "autostart" === t.get("playReason"),
                o = t.get("viewable");
              !n || i || o || (this.noResume = !0);
              var a = d.d;
              this.trigger(a, e), m.call(this, { type: a });
            }),
            (this.replacePlaylistItem = function (e) {
              y || (t.set("playlistItem", e), f.srcReset());
            }),
            (this.destroy = function () {
              y ||
                ((y = !0),
                this.trigger("destroyed"),
                this.off(),
                n.clickHandler() &&
                  n.clickHandler().revertAlternateClickHandlers(),
                t.off(null, null, f),
                f.off(null, null, u),
                f.destroy(),
                k && f.model && (t.attributes.state = d.ob),
                e.forwardEvents(),
                t.set("instream", null),
                (f = null),
                (j = {}),
                (g = null),
                k &&
                  !t.attributes._destroyed &&
                  (e.attachMedia(),
                  this.noResume || (x ? e.stopVideo() : e.playVideo())));
            }),
            (this.getState = function () {
              return !y && f.model.get("state");
            }),
            (this.setText = function (e) {
              return y ? this : (n.setAltText(e || ""), this);
            }),
            (this.hide = function () {
              y || t.set("hideAdsControls", !0);
            }),
            (this.getMediaElement = function () {
              return y ? null : f.primedElement;
            }),
            (this.setSkipOffset = function (e) {
              (s = e > 0 ? e : null), f && f.model.set("skipOffset", s);
            });
        };
      Object(i.g)(le.prototype, f.a);
      var ce = le,
        ue = n(66),
        de = n(63),
        pe = function (e) {
          var t = this,
            n = [],
            i = {},
            o = 0,
            a = 0;
          function r(e) {
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
            var i = (a = t),
              o = e.get("captionLabel");
            if ("Off" !== o) {
              for (var r = 0; r < n.length; r++) {
                var s = n[r];
                if (o && o === s.name) {
                  i = r + 1;
                  break;
                }
                s.default || s.defaulttrack || "default" === s._id
                  ? (i = r + 1)
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
              (l(a), e.set("captionsList", t));
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
                    var a = function (e) {
                        var o,
                          a = n[e];
                        ("subtitles" !== (o = a.kind) && "captions" !== o) ||
                          i[a._id] ||
                          (r(a),
                          Object(ue.c)(
                            a,
                            function (e) {
                              !(function (e, t) {
                                e.data = t;
                              })(a, e);
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
                    a(s);
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
                  for (var t = 0; t < e.length; t++) r(e[t]);
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
      Object(i.g)(pe.prototype, f.a);
      var we = pe,
        fe = function (e, t) {
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
        he = n(35),
        je = 44,
        ge = function (e) {
          var t = e.get("height");
          if (e.get("aspectratio")) return !1;
          if ("string" == typeof t && t.indexOf("%") > -1) return !1;
          var n = 1 * t || NaN;
          return (
            !!(n = isNaN(n) ? e.get("containerHeight") : n) && n && n <= je
          );
        },
        be = n(54);
      function me(e, t) {
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
            var o = ve(e);
            if (!o) return 0;
            var a = o,
              r = e.parentNode,
              s = !1;
            for (; !s; ) {
              var l = null;
              if (
                (r === n || r === t || 1 !== r.nodeType
                  ? ((s = !0), (l = i))
                  : "visible" !== window.getComputedStyle(r).overflow &&
                    (l = ve(r)),
                l &&
                  ((c = l),
                  (u = a),
                  (d = void 0),
                  (p = void 0),
                  (w = void 0),
                  (f = void 0),
                  (h = void 0),
                  (j = void 0),
                  (d = Math.max(c.top, u.top)),
                  (p = Math.min(c.bottom, u.bottom)),
                  (w = Math.max(c.left, u.left)),
                  (f = Math.min(c.right, u.right)),
                  (j = p - d),
                  !(a = (h = f - w) >= 0 &&
                    j >= 0 && {
                      top: d,
                      bottom: p,
                      left: w,
                      right: f,
                      width: h,
                      height: j,
                    })))
              )
                return 0;
              r = r.parentNode;
            }
            var c, u, d, p, w, f, h, j;
            var g = o.width * o.height,
              b = a.width * a.height;
            return g ? b / g : 0;
          })(t)),
          window.top !== window.self && n)
          ? 0
          : n;
      }
      function ve(e) {
        try {
          return e.getBoundingClientRect();
        } catch (e) {}
      }
      var ye = n(49),
        ke = n(42),
        xe = n(58),
        Oe = n(10);
      var Ce = n(32),
        Me = n(5),
        Se = n(6),
        Te = [
          "fullscreenchange",
          "webkitfullscreenchange",
          "mozfullscreenchange",
          "MSFullscreenChange",
        ],
        Ee = function (e, t, n) {
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
              a = !(!i || !o),
              r = Te.length;
            r--;

          )
            t.addEventListener(Te[r], n);
          return {
            events: Te,
            supportsDomFullscreen: function () {
              return a;
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
              for (var e = Te.length; e--; ) t.removeEventListener(Te[e], n);
            },
          };
        },
        _e = n(40);
      function ze(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      var Pe = (function () {
          function e(t, n) {
            !(function (e, t) {
              if (!(e instanceof t))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
              Object(i.g)(this, f.a),
              this.revertAlternateClickHandlers(),
              (this.domElement = n),
              (this.model = t),
              (this.ui = new _e.a(n)
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
            ]) && ze(t.prototype, n),
            o && ze(t, o),
            e
          );
        })(),
        Ae = n(59),
        Le = function (e, t) {
          var n = t ? " jw-hide" : "";
          return '<div class="jw-logo jw-logo-'
            .concat(e)
            .concat(n, ' jw-reset"></div>');
        },
        Re = {
          linktarget: "_blank",
          margin: 8,
          hide: !1,
          position: "top-right",
        };
      function Ie(e) {
        var t, n;
        Object(i.g)(this, f.a);
        var o = new Image();
        (this.setup = function () {
          ((n = Object(i.g)({}, Re, e.get("logo"))).position =
            n.position || Re.position),
            (n.hide = "true" === n.hide.toString()),
            n.file &&
              "control-bar" !== n.position &&
              (t || (t = Object(Me.e)(Le(n.position, n.hide))),
              e.set("logo", n),
              (o.onload = function () {
                var i = this.height,
                  o = this.width,
                  a = { backgroundImage: 'url("' + this.src + '")' };
                if (n.margin !== Re.margin) {
                  var r = /(\w+)-(\w+)/.exec(n.position);
                  3 === r.length &&
                    ((a["margin-" + r[1]] = n.margin),
                    (a["margin-" + r[2]] = n.margin));
                }
                var s = 0.15 * e.get("containerHeight"),
                  l = 0.15 * e.get("containerWidth");
                if (i > s || o > l) {
                  var c = o / i;
                  l / s > c ? ((i = s), (o = s * c)) : ((o = l), (i = l / c));
                }
                (a.width = Math.round(o)),
                  (a.height = Math.round(i)),
                  Object(Oe.d)(t, a),
                  e.set("logoWidth", a.width);
              }),
              (o.src = n.file),
              n.link &&
                (t.setAttribute("tabindex", "0"),
                t.setAttribute("aria-label", e.get("localization").logo)),
              (this.ui = new _e.a(t).on(
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
      var Be = function (e) {
        (this.model = e), (this.image = null);
      };
      Object(i.g)(Be.prototype, {
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
            Object(Oe.d)(this.el, { backgroundImage: n });
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
                var a = this;
                return void (i.onload = function () {
                  a.resize(e, t, n);
                });
              }
              var r = i.width / i.height;
              Math.abs(this.playerAspectRatio - r) < 0.09 && (o = "cover");
            }
            Object(Oe.d)(this.el, { backgroundSize: o });
          }
          var s;
        },
        element: function () {
          return this.el;
        },
      });
      var Ve = Be,
        Ne = function (e) {
          this.model = e.player;
        };
      Object(i.g)(Ne.prototype, {
        hide: function () {
          Object(Oe.d)(this.el, { display: "none" });
        },
        show: function () {
          Object(Oe.d)(this.el, { display: "" });
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
          Object(Oe.d)(this.el, t);
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
      function Fe(e, t) {
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
                    a = this.input,
                    r = (this.ui = new _e.a(a, { preventScrolling: !0 })
                      .on("dragStart", function () {
                        (e = o.offsetLeft),
                          (t = o.offsetTop),
                          (n = window.innerHeight),
                          (i = window.innerWidth);
                      })
                      .on("drag", function (a) {
                        var s = Math.max(e + a.pageX - r.startX, 0),
                          l = Math.max(t + a.pageY - r.startY, 0),
                          c = Math.max(i - (s + o.clientWidth), 0),
                          u = Math.max(n - (l + o.clientHeight), 0);
                        0 === c ? (s = "auto") : (c = "auto"),
                          0 === l ? (u = "auto") : (l = "auto"),
                          Object(Oe.d)(o, {
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
            ]) && Fe(t.prototype, n),
            i && Fe(t, i),
            e
          );
        })(),
        Ue = n(55);
      n(69);
      var We = m.OS.mobile,
        Qe = m.Browser.ie,
        Ye = null;
      var Xe = function (e, t) {
        var n,
          o,
          a,
          r,
          s = this,
          l = Object(i.g)(this, f.a, { isSetup: !1, api: e, model: t }),
          c = t.get("localization"),
          u = Object(Me.e)(fe(t.get("id"), c.player)),
          p = u.querySelector(".jw-wrapper"),
          h = u.querySelector(".jw-media"),
          j = new De(p),
          g = new Ve(t, e),
          b = new He(t),
          v = new Ae.b(t);
        v.on("all", l.trigger, l);
        var y = -1,
          k = -1,
          x = -1,
          O = t.get("floating");
        this.dismissible = O && O.dismissible;
        var C,
          M,
          S,
          T = !1,
          E = {},
          _ = null,
          z = null;
        function P() {
          return We && !Object(Me.f)();
        }
        function A() {
          Object(ke.a)(k), (k = Object(ke.b)(L));
        }
        function L() {
          l.isSetup && (l.updateBounds(), l.updateStyles(), l.checkResized());
        }
        function R(e, n) {
          if (Object(i.r)(e) && Object(i.r)(n)) {
            var o = Object(xe.a)(e);
            Object(xe.b)(u, o);
            var a = o < 2;
            Object(Me.v)(u, "jw-flag-small-player", a),
              Object(Me.v)(u, "jw-orientation-portrait", n > e);
          }
          if (t.get("controls")) {
            var r = ge(t);
            Object(Me.v)(u, "jw-flag-audio-player", r), t.set("audioMode", r);
          }
        }
        function I() {
          t.set("visibility", me(t, u));
        }
        (this.updateBounds = function () {
          Object(ke.a)(k);
          var e = t.get("isFloating") ? p : u,
            n = document.body.contains(e),
            i = Object(Me.c)(e),
            r = Math.round(i.width),
            s = Math.round(i.height);
          if (((E = Object(Me.c)(u)), r === o && s === a))
            return (o && a) || A(), void t.set("inDom", n);
          (r && s) || (o && a) || A(),
            (r || s || n) &&
              (t.set("containerWidth", r), t.set("containerHeight", s)),
            t.set("inDom", n),
            n && be.a.observe(u);
        }),
          (this.updateStyles = function () {
            var e = t.get("containerWidth"),
              n = t.get("containerHeight");
            R(e, n), z && z.resize(e, n), $(e, n), v.resize(), O && F();
          }),
          (this.checkResized = function () {
            var e = t.get("containerWidth"),
              n = t.get("containerHeight"),
              i = t.get("isFloating");
            if (e !== o || n !== a) {
              this.resizeListener ||
                (this.resizeListener = new Ue.a(p, this, t)),
                (o = e),
                (a = n),
                l.trigger(d.hb, { width: e, height: n });
              var s = Object(xe.a)(e);
              _ !== s && ((_ = s), l.trigger(d.j, { breakpoint: _ }));
            }
            i !== r && ((r = i), l.trigger(d.x, { floating: i }), I());
          }),
          (this.responsiveListener = A),
          (this.setup = function () {
            g.setup(u.querySelector(".jw-preview")),
              b.setup(u.querySelector(".jw-title")),
              (n = new Ie(t)).setup(),
              n.setContainer(p),
              n.on(d.A, K),
              v.setup(u.id, t.get("captions")),
              b.element().parentNode.insertBefore(v.element(), b.element()),
              (C = (function (e, t, n) {
                var i = new Pe(t, n),
                  o = t.get("controls");
                i.on({
                  click: function () {
                    l.trigger(d.p),
                      z &&
                        (ce()
                          ? z.settingsMenu.close()
                          : ue()
                          ? z.infoOverlay.close()
                          : e.playToggle({ reason: "interaction" }));
                  },
                  tap: function () {
                    l.trigger(d.p),
                      ce() && z.settingsMenu.close(),
                      ue() && z.infoOverlay.close();
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
                        v.renderCues(!0);
                    } else z && (z.showing ? z.userInactive() : z.userActive());
                  },
                  doubleClick: function () {
                    return z && e.setFullscreen();
                  },
                }),
                  We ||
                    (u.addEventListener("mousemove", W),
                    u.addEventListener("mouseover", Q),
                    u.addEventListener("mouseout", Y));
                return i;
              })(e, t, h)),
              (S = new _e.a(u).on("click", function () {})),
              (M = Ee(u, document, te)),
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
                $(), v.resize();
              }),
              t.player.on("change:errorEvent", ae),
              t.change("stretching", X);
            var i = t.get("width"),
              o = t.get("height"),
              a = G(i, o);
            Object(Oe.d)(u, a),
              t.change("aspectratio", Z),
              R(i, o),
              t.get("controls") ||
                (Object(Me.a)(u, "jw-flag-controls-hidden"),
                Object(Me.o)(u, "jw-floating-dismissible")),
              Qe && Object(Me.a)(u, "jw-ie");
            var r = t.get("skin") || {};
            r.name && Object(Me.p)(u, /jw-skin-\S+/, "jw-skin-" + r.name);
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
            })(r);
            !(function (e, t) {
              var n;
              function i(t, n, i, o) {
                if (i) {
                  t = Object(w.f)(t, "#" + e + (o ? "" : " "));
                  var a = {};
                  (a[n] = i), Object(Oe.b)(t.join(", "), a, e);
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
                        Object(Oe.b)(
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
                      Object(Oe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast:hover google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(Oe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast:focus google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(Oe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast google-cast-launcher.jw-off:focus"
                        ),
                        "{--disconnected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(Oe.b)(
                        "#".concat(e, " .jw-icon-cast google-cast-launcher"),
                        "{--connected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(Oe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast google-cast-launcher:focus"
                        ),
                        "{--connected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(Oe.b)(
                        "#".concat(
                          e,
                          " .jw-icon-cast:hover google-cast-launcher"
                        ),
                        "{--connected-color: ".concat(t.iconsActive, "}"),
                        e
                      ),
                      Object(Oe.b)(
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
                        Object(Oe.c)(t, 50)
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
                      Object(Oe.b)("#".concat(e, " .jw-color-active"), n, e),
                        Object(Oe.b)(
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
                      Object(Oe.b)("#".concat(e, " .jw-color-inactive"), i, e),
                        Object(Oe.b)(
                          "#".concat(e, " .jw-color-inactive-hover:hover"),
                          i,
                          e
                        );
                    }
                  })(t.menus));
            })(t.get("id"), s),
              t.set("mediaContainer", h),
              t.set("iFrame", m.Features.iframe),
              t.set("activeTab", Object(ye.a)()),
              t.set("touchMode", We && ("string" == typeof o || o >= je)),
              be.a.add(this),
              t.get("enableGradient") &&
                !Qe &&
                Object(Me.a)(u, "jw-ab-drop-shadow"),
              (this.isSetup = !0),
              t.trigger("viewSetup", u);
            var c = document.body.contains(u);
            c && be.a.observe(u), t.set("inDom", c);
          }),
          (this.init = function () {
            this.updateBounds(),
              t.on("change:fullscreen", J),
              t.on("change:activeTab", I),
              t.on("change:fullscreen", I),
              t.on("change:intersectionRatio", I),
              t.on("change:visibility", U),
              t.on("instreamMode", function (e) {
                e ? de() : pe();
              }),
              I(),
              1 !== be.a.size() || t.get("visibility") || U(t, 1, 0);
            var e = t.player;
            t.change("state", re),
              e.change("controls", q),
              t.change("streamType", ie),
              t.change("mediaType", oe),
              e.change("playlistItem", function (e, t) {
                le(e, t);
              }),
              (o = a = null),
              O && We && be.a.addScrollHandler(F),
              this.checkResized();
          });
        var B,
          V = 62,
          N = !0;
        function H() {
          var e = t.get("isFloating"),
            n = E.top < V,
            i = n ? E.top <= window.scrollY : E.top <= window.scrollY + V;
          !e && i ? we(0, n) : e && !i && we(1, n);
        }
        function F() {
          P() &&
            t.get("inDom") &&
            (clearTimeout(B),
            (B = setTimeout(H, 150)),
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
            o && a && l.trigger(d.o, n);
        }
        function D() {
          var e = new qe(document, l.element());
          l.addControls(e);
        }
        function U(e, t, n) {
          t && !n && (re(e, e.get("state")), l.updateStyles());
        }
        function W(e) {
          z && z.mouseMove(e);
        }
        function Q(e) {
          z && !z.showing && "IFRAME" === e.target.nodeName && z.userActive();
        }
        function Y(e) {
          z &&
            z.showing &&
            ((e.relatedTarget && !u.contains(e.relatedTarget)) ||
              (!e.relatedTarget && m.Features.iframe)) &&
            z.userActive();
        }
        function X(e, t) {
          Object(Me.p)(u, /jw-stretch-\S+/, "jw-stretch-" + t);
        }
        function Z(e, n) {
          Object(Me.v)(u, "jw-flag-aspect-mode", !!n);
          var i = u.querySelectorAll(".jw-aspect");
          Object(Oe.d)(i, { paddingTop: n || null }),
            l.isSetup &&
              n &&
              !t.get("isFloating") &&
              (Object(Oe.d)(u, G(e.get("width"))), L());
        }
        function K(n) {
          n.link
            ? (e.pause({ reason: "interaction" }),
              e.setFullscreen(!1),
              Object(Me.l)(n.link, n.linktarget, { rel: "noreferrer" }))
            : t.get("controls") && e.playToggle({ reason: "interaction" });
        }
        (this.addControls = function (n) {
          var i = this;
          (z = n),
            Object(Me.o)(u, "jw-flag-controls-hidden"),
            Object(Me.v)(u, "jw-floating-dismissible", this.dismissible),
            n.enable(e, t),
            a && (R(o, a), n.resize(o, a), v.renderCues(!0)),
            n.on("userActive userInactive", function () {
              var e = t.get("state");
              (e !== d.pb && e !== d.jb) || v.renderCues(!0);
            }),
            n.on("dismissFloating", function () {
              i.stopFloating(!0), e.pause({ reason: "interaction" });
            }),
            n.on("all", l.trigger, l),
            t.get("instream") && z.setupInstream();
        }),
          (this.removeControls = function () {
            z && (z.disable(t), (z = null)),
              Object(Me.a)(u, "jw-flag-controls-hidden"),
              Object(Me.o)(u, "jw-floating-dismissible");
          });
        var J = function (t, n) {
          if (
            (n && z && t.get("autostartMuted") && z.unmuteAutoplay(e, t),
            M.supportsDomFullscreen())
          )
            n ? M.requestFullscreen() : M.exitFullscreen(), ne(u, n);
          else if (Qe) ne(u, n);
          else {
            var i = t.get("instream"),
              o = i ? i.provider : null,
              a = t.getVideo() || o;
            a && a.setFullscreen && a.setFullscreen(n);
          }
        };
        function G(e, n, o) {
          var a = { width: e };
          if (
            (o && void 0 !== n && t.set("aspectratio", null),
            !t.get("aspectratio"))
          ) {
            var r = n;
            Object(i.r)(r) && 0 !== r && (r = Math.max(r, je)), (a.height = r);
          }
          return a;
        }
        function $(e, n) {
          if (
            ((e && !isNaN(1 * e)) || (e = t.get("containerWidth"))) &&
            ((n && !isNaN(1 * n)) || (n = t.get("containerHeight")))
          ) {
            g && g.resize(e, n, t.get("stretching"));
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
            A(),
            clearTimeout(y),
            (y = setTimeout($, 200));
        }
        function ne(e, t) {
          Object(Me.v)(e, "jw-flag-fullscreen", t),
            Object(Oe.d)(document.body, { overflowY: t ? "hidden" : "" }),
            t && z && z.userActive(),
            $(),
            A();
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
            a = n && !o ? h : h.nextSibling;
          g.el.parentNode.insertBefore(g.el, a);
        }
        function ae(e, t) {
          if (t) {
            var n = Object(he.a)(e, t);
            he.a.cloneIcon &&
              n.querySelector(".jw-icon").appendChild(he.a.cloneIcon("error")),
              b.hide(),
              u.appendChild(n.firstChild),
              Object(Me.v)(u, "jw-flag-audio-player", !!e.get("audioMode"));
          } else b.playlistItem(e, e.get("playlistItem"));
        }
        function re(e, t, n) {
          if (l.isSetup) {
            if (n === d.lb) {
              var i = u.querySelector(".jw-error-msg");
              i && i.parentNode.removeChild(i);
            }
            Object(ke.a)(x),
              t === d.pb
                ? se(t)
                : (x = Object(ke.b)(function () {
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
              v && v.hide();
              break;
            default:
              v &&
                (v.show(), e === d.ob && z && !z.showing && v.renderCues(!0));
          }
        }
        (this.resize = function (e, n) {
          var i = G(e, n, !0);
          void 0 !== e &&
            void 0 !== n &&
            (t.set("width", e), t.set("height", n)),
            Object(Oe.d)(u, i),
            t.get("isFloating") && ve(),
            L();
        }),
          (this.resizeMedia = $),
          (this.setPosterImage = function (e, t) {
            t.setImage(e && e.image);
          });
        var le = function (e, t) {
            s.setPosterImage(t, g),
              We &&
                (function (e, t) {
                  var n = e.get("mediaElement");
                  if (n) {
                    var i = Object(Me.j)(t.title || "");
                    n.setAttribute("title", i.textContent);
                  }
                })(e, t);
          },
          ce = function () {
            var e = z && z.settingsMenu;
            return !(!e || !e.visible);
          },
          ue = function () {
            var e = z && z.infoOverlay;
            return !(!e || !e.visible);
          },
          de = function () {
            Object(Me.a)(u, "jw-flag-ads"), z && z.setupInstream(), j.disable();
          },
          pe = function () {
            if (C) {
              z && z.destroyInstream(t),
                Ye !== u || Object(Se.m)() || j.enable(),
                l.setAltText(""),
                Object(Me.o)(u, ["jw-flag-ads", "jw-flag-ads-hide-controls"]),
                t.set("hideAdsControls", !1);
              var e = t.getVideo();
              e && e.setContainer(h), C.revertAlternateClickHandlers();
            }
          };
        function we(e, n) {
          if (e < 0.5 && !Object(Se.m)()) {
            var i = t.get("state");
            i !== d.mb &&
              i !== d.lb &&
              i !== d.kb &&
              null === Ye &&
              ((Ye = u),
              t.set("isFloating", !0),
              Object(Me.a)(u, "jw-flag-floating"),
              n &&
                (Object(Oe.d)(p, {
                  transform: "translateY(-".concat(V - E.top, "px)"),
                }),
                setTimeout(function () {
                  Object(Oe.d)(p, {
                    transform: "translateY(0)",
                    transition:
                      "transform 150ms cubic-bezier(0, 0.25, 0.25, 1)",
                  });
                })),
              Object(Oe.d)(u, {
                backgroundImage: g.el.style.backgroundImage || t.get("image"),
              }),
              ve(),
              t.get("instreamMode") || j.enable(),
              A());
          } else l.stopFloating(!1, n);
        }
        function ve() {
          var e = t.get("width"),
            n = t.get("height"),
            o = G(e);
          if (((o.maxWidth = Math.min(400, E.width)), !t.get("aspectratio"))) {
            var a = E.width,
              r = E.height / a || 0.5625;
            Object(i.r)(e) && Object(i.r)(n) && (r = n / e),
              Z(t, 100 * r + "%");
          }
          Object(Oe.d)(p, o);
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
            return p;
          }),
          (this.controlsContainer = function () {
            return z ? z.element() : null;
          }),
          (this.getSafeRegion = function () {
            var e =
                !(arguments.length > 0 && void 0 !== arguments[0]) ||
                arguments[0],
              t = { x: 0, y: 0, width: o || 0, height: a || 0 };
            return z && e && (t.height -= z.controlbarHeight()), t;
          }),
          (this.setCaptions = function (e) {
            v.clear(), v.setup(t.get("id"), e), v.resize();
          }),
          (this.setIntersection = function (e) {
            var n = Math.round(100 * e.intersectionRatio) / 100;
            t.set("intersectionRatio", n),
              O && !P() && (T = T || n >= 0.5) && we(n);
          }),
          (this.stopFloating = function (e, n) {
            if ((e && ((O = null), be.a.removeScrollHandler(F)), Ye === u)) {
              (Ye = null), t.set("isFloating", !1);
              var i = function () {
                Object(Me.o)(u, "jw-flag-floating"),
                  Z(t, t.get("aspectratio")),
                  Object(Oe.d)(u, { backgroundImage: null }),
                  Object(Oe.d)(p, {
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
                ? (Object(Oe.d)(p, {
                    transform: "translateY(-".concat(V - E.top, "px)"),
                    "transition-timing-function": "ease-out",
                  }),
                  setTimeout(i, 150))
                : i(),
                j.disable(),
                A();
            }
          }),
          (this.destroy = function () {
            t.destroy(),
              be.a.unobserve(u),
              be.a.remove(this),
              (this.isSetup = !1),
              this.off(),
              Object(ke.a)(k),
              clearTimeout(y),
              Ye === u && (Ye = null),
              S && (S.destroy(), (S = null)),
              M && (M.destroy(), (M = null)),
              z && z.disable(t),
              C &&
                (C.destroy(),
                u.removeEventListener("mousemove", W),
                u.removeEventListener("mouseout", Y),
                u.removeEventListener("mouseover", Q),
                (C = null)),
              v.destroy(),
              n && (n.destroy(), (n = null)),
              Object(Oe.a)(t.get("id")),
              this.resizeListener &&
                (this.resizeListener.destroy(), delete this.resizeListener),
              O && We && be.a.removeScrollHandler(F);
          });
      };
      function Ze(e, t, n) {
        return (Ze =
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
      function Ke(e) {
        return (Ke =
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
      function Je(e, t) {
        if (!(e instanceof t))
          throw new TypeError("Cannot call a class as a function");
      }
      function Ge(e, t) {
        for (var n = 0; n < t.length; n++) {
          var i = t[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(e, i.key, i);
        }
      }
      function $e(e, t, n) {
        return t && Ge(e.prototype, t), n && Ge(e, n), e;
      }
      function et(e, t) {
        return !t || ("object" !== Ke(t) && "function" != typeof t) ? ot(e) : t;
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
      var at = /^change:(.+)$/;
      function rt(e, t, n) {
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
              Je(this, t),
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
                function (t, i, a, r) {
                  i === e && (i = ot(ot(o))),
                    (n && !n(t, i, a, r)) || o.trigger(t, i, a, r);
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
            $e(t, [
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
                      function (n, i, o, a) {
                        i === e && (i = t), t.trigger(n, i, o, a);
                      },
                      this
                    ),
                    n && rt(this, e.attributes, n.attributes);
                },
              },
            ]),
            t
          );
        })(v.a),
        ct = (function (e) {
          function t(e) {
            var n;
            return (
              Je(this, t),
              ((n = et(
                this,
                tt(t).call(this, e, function (e) {
                  var t = n._instreamModel;
                  if (t) {
                    var i = at.exec(e);
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
            $e(t, [
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
                    : Ze(tt(t.prototype), "getVideo", this).call(this);
                },
              },
              {
                key: "destroy",
                value: function () {
                  Ze(tt(t.prototype), "destroy", this).call(this),
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
                      function (n, i, o, a) {
                        i === e && (i = t), t.trigger(n, i, o, a);
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
                      rt(this, e.attributes, this._model.attributes);
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
                    rt(this, this._model.attributes, o);
                  }
                },
              },
            ]),
            t
          );
        })(lt);
      var ut,
        dt,
        pt = n(64),
        wt =
          (ut = window).URL && ut.URL.createObjectURL
            ? ut.URL
            : ut.webkitURL || ut.mozURL;
      function ft(e, t) {
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
          (e.src = wt.createObjectURL(dt)),
          e.play() || Object(pt.a)(e)
        );
      }
      var ht = "autoplayEnabled",
        jt = "autoplayMuted",
        gt = "autoplayDisabled",
        bt = {};
      var mt = n(65);
      function vt(e) {
        return (
          (e = e || window.event) &&
          /^(?:mouse|pointer|touch|gesture|click|key)/.test(e.type)
        );
      }
      var yt = n(24),
        kt = "tabHidden",
        xt = "tabVisible",
        Ot = function (e) {
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
            (e._qoeItem = new yt.a()),
            (e._qoeItem.getFirstFrame = function () {
              var e = this.between(d.N, d.H),
                t = this.between(xt, d.H);
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
                    var a = t.mediaController.mediaModel;
                    a.off("change:".concat(d.U), null, a),
                      a.change(
                        d.U,
                        function (e, n) {
                          n && t.trigger(d.U, n);
                        },
                        a
                      );
                  }
                  Ct(e, t);
                }
              }),
                (e._onTime = Ot(e._triggerFirstFrame)),
                (e._onPlayAttempt = function () {
                  e._qoeItem.tick(d.N);
                }),
                (e._onTabVisible = function (t, n) {
                  n ? e._qoeItem.tick(xt) : e._qoeItem.tick(kt);
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
      function St(e) {
        return (St =
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
      var Tt = function () {},
        Et = function () {};
      Object(i.g)(Tt.prototype, {
        setup: function (e, t, n, w, h, b) {
          var v,
            y,
            k,
            x,
            O = this,
            C = this,
            M = (C._model = new A()),
            S = !1,
            T = !1,
            E = null,
            _ = j(H),
            z = j(Et);
          (C.originalContainer = C.currentContainer = n),
            (C._events = w),
            (C.trigger = function (e, t) {
              var n = (function (e, t, n) {
                var o = n;
                switch (t) {
                  case "time":
                  case "beforePlay":
                  case "pause":
                  case "play":
                  case "ready":
                    var a = e.get("viewable");
                    void 0 !== a && (o = Object(i.g)({}, n, { viewable: a }));
                }
                return o;
              })(M, e, t);
              return f.a.trigger.call(this, e, n);
            });
          var P = new s.a(C, ["trigger"], function () {
              return !0;
            }),
            L = function (e, t) {
              C.trigger(e, t);
            };
          M.setup(e);
          var R = M.get("backgroundLoading"),
            I = new ct(M);
          (v = this._view = new Xe(t, I)).on(
            "all",
            function (e, t) {
              (t && t.doNotForward) || L(e, t);
            },
            C
          );
          var B = (this._programController = new Y(M, b));
          ue(),
            B.on("all", L, C)
              .on(
                "subtitlesTracks",
                function (e) {
                  y.setSubtitlesTracks(e.tracks);
                  var t = y.getCurrentIndex();
                  t > 0 && re(t, e.tracks);
                },
                C
              )
              .on(
                d.F,
                function () {
                  Promise.resolve().then(ae);
                },
                C
              )
              .on(d.G, C.triggerError, C),
            Mt(M, B),
            M.on(d.w, C.triggerError, C),
            M.on(
              "change:state",
              function (e, t, n) {
                X() || Z.call(O, e, t, n);
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
          var V = function e(t, n) {
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
                $("autostart"),
              de.flush());
          }
          function F(e, t) {
            C.trigger("viewable", { viewable: t }), q();
          }
          function q() {
            if (
              (o.a[0] === t || 1 === M.get("viewable")) &&
              "idle" === M.get("state") &&
              !1 === M.get("autostart")
            )
              if (!b.primed() && m.OS.android) {
                var e = b.getTestElement(),
                  n = C.getMute();
                Promise.resolve()
                  .then(function () {
                    return ft(e, { muted: n });
                  })
                  .then(function () {
                    "idle" === M.get("state") && B.preloadVideo();
                  })
                  .catch(Et);
              } else B.preloadVideo();
          }
          function D(e) {
            (C._instreamAdapter.noResume = !e), e || te({ reason: "viewable" });
          }
          function U(e) {
            e || (C.pause({ reason: "viewable" }), M.set("playOnViewable", !e));
          }
          function W(e, t) {
            var n = X();
            if (e.get("playOnViewable")) {
              if (t) {
                var i = e.get("autoPause").pauseAds,
                  o = e.get("pauseReason");
                K() === d.mb
                  ? $("viewable")
                  : (n && !i) ||
                    "interaction" === o ||
                    J({ reason: "viewable" });
              } else
                m.OS.mobile &&
                  !n &&
                  (C.pause({ reason: "autostart" }),
                  M.set("playOnViewable", !0));
              m.OS.mobile && n && D(t);
            }
          }
          function Q(e, t) {
            var n = e.get("state"),
              i = X(),
              o = e.get("playReason");
            i
              ? e.get("autoPause").pauseAds
                ? U(t)
                : D(t)
              : n === d.pb || n === d.jb
              ? U(t)
              : n === d.mb &&
                "playlist" === o &&
                e.once("change:state", function () {
                  U(t);
                });
          }
          function X() {
            var e = C._instreamAdapter;
            return !!e && e.getState();
          }
          function K() {
            var e = X();
            return e || M.get("state");
          }
          function J(e) {
            if ((_.cancel(), (T = !1), M.get("state") === d.lb))
              return Promise.resolve();
            var n = G(e);
            return (
              M.set("playReason", n),
              X()
                ? (t.pauseAd(!1, e), Promise.resolve())
                : (M.get("state") === d.kb && (ee(!0), C.setItemIndex(0)),
                  !S &&
                  ((S = !0),
                  C.trigger(d.C, {
                    playReason: n,
                    startTime:
                      e && e.startTime
                        ? e.startTime
                        : M.get("playlistItem").starttime,
                  }),
                  (S = !1),
                  vt() && !b.primed() && b.prime(),
                  "playlist" === n &&
                    M.get("autoPause").viewability &&
                    Q(M, M.get("viewable")),
                  x)
                    ? (vt() && !R && M.get("mediaElement").load(),
                      (x = !1),
                      (k = null),
                      Promise.resolve())
                    : B.playVideo(n).then(b.played))
            );
          }
          function G(e) {
            return e && e.reason ? e.reason : "unknown";
          }
          function $(e) {
            if (K() === d.mb) {
              _ = j(H);
              var t = M.get("advertising");
              (function (e, t) {
                var n = t.cancelable,
                  i = t.muted,
                  o = void 0 !== i && i,
                  a = t.allowMuted,
                  r = void 0 !== a && a,
                  s = t.timeout,
                  l = void 0 === s ? 1e4 : s,
                  c = e.getTestElement(),
                  u = o ? "muted" : "".concat(r);
                bt[u] ||
                  (bt[u] = ft(c, { muted: o })
                    .catch(function (e) {
                      if (!n.cancelled() && !1 === o && r)
                        return ft(c, { muted: (o = !0) });
                      throw e;
                    })
                    .then(function () {
                      return o ? ((bt[u] = null), jt) : ht;
                    })
                    .catch(function (e) {
                      throw (
                        (clearTimeout(d), (bt[u] = null), (e.reason = gt), e)
                      );
                    }));
                var d,
                  p = bt[u].then(function (e) {
                    if ((clearTimeout(d), n.cancelled())) {
                      var t = new Error("Autoplay test was cancelled");
                      throw ((t.reason = "cancelled"), t);
                    }
                    return e;
                  }),
                  w = new Promise(function (e, t) {
                    d = setTimeout(function () {
                      bt[u] = null;
                      var e = new Error("Autoplay test timed out");
                      (e.reason = "timeout"), t(e);
                    }, l);
                  });
                return Promise.race([p, w]);
              })(b, {
                cancelable: _,
                muted: C.getMute(),
                allowMuted: !t || t.autoplayadsmuted,
              })
                .then(function (t) {
                  return (
                    M.set("canAutoplay", t),
                    t !== jt ||
                      C.getMute() ||
                      (M.set("autostartMuted", !0),
                      ue(),
                      M.once("change:autostartMuted", function (e) {
                        e.off("change:viewable", W),
                          C.trigger(d.M, { mute: M.getMute() });
                      })),
                    C.getMute() &&
                      M.get("enableDefaultCaptions") &&
                      y.selectDefaultIndex(1),
                    J({ reason: e }).catch(function () {
                      C._instreamAdapter || M.set("autostartFailed", !0),
                        (k = null);
                    })
                  );
                })
                .catch(function (e) {
                  if (
                    (M.set("canAutoplay", gt),
                    M.set("autostart", !1),
                    !_.cancelled())
                  ) {
                    var t = Object(g.w)(e);
                    C.trigger(d.h, { reason: e.reason, code: t, error: e });
                  }
                });
            }
          }
          function ee(e) {
            if ((_.cancel(), de.empty(), X())) {
              var t = C._instreamAdapter;
              return (
                t && (t.noResume = !0),
                void (k = function () {
                  return B.stopVideo();
                })
              );
            }
            (k = null),
              !e && (T = !0),
              S && (x = !0),
              M.set("errorEvent", void 0),
              B.stopVideo();
          }
          function te(e) {
            var t = G(e);
            M.set("pauseReason", t), M.set("playOnViewable", "viewable" === t);
          }
          function ne(e) {
            (k = null), _.cancel();
            var n = X();
            if (n && n !== d.ob) return te(e), void t.pauseAd(!0, e);
            switch (M.get("state")) {
              case d.lb:
                return;
              case d.pb:
              case d.jb:
                te(e), B.pause();
                break;
              default:
                S && (x = !0);
            }
          }
          function ie(e, t) {
            ee(!0), C.setItemIndex(e), C.play(t);
          }
          function oe(e) {
            ie(M.get("item") + 1, e);
          }
          function ae() {
            C.completeCancelled() ||
              ((k = C.completeHandler),
              C.shouldAutoAdvance()
                ? C.nextItem()
                : M.get("repeat")
                ? oe({ reason: "repeat" })
                : (m.OS.iOS && le(!1),
                  M.set("playOnViewable", !1),
                  M.set("state", d.kb),
                  C.trigger(d.cb, {})));
          }
          function re(e, t) {
            (e = parseInt(e, 10) || 0),
              M.persistVideoSubtitleTrack(e, t),
              (B.subtitles = e),
              C.trigger(d.k, { tracks: se(), track: e });
          }
          function se() {
            return y.getCaptionsList();
          }
          function le(e) {
            Object(i.n)(e) || (e = !M.get("fullscreen")),
              M.set("fullscreen", e),
              C._instreamAdapter &&
                C._instreamAdapter._adModel &&
                C._instreamAdapter._adModel.set("fullscreen", e);
          }
          function ue() {
            (B.mute = M.getMute()), (B.volume = M.get("volume"));
          }
          M.on("change:playReason change:pauseReason", V),
            C.on(d.c, function (e) {
              return V(0, e.playReason);
            }),
            C.on(d.b, function (e) {
              return V(0, e.pauseReason);
            }),
            M.on("change:scrubbing", function (e, t) {
              t
                ? ((E = M.get("state") !== d.ob), ne())
                : E && J({ reason: "interaction" });
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
                        o = Object(mt.b)(n, i);
                      e.setStreamType(o);
                    }
                  },
                  this
                );
              var i = e.get("item") + 1,
                o = "autoplay" === (e.get("related") || {}).oncomplete,
                a = e.get("playlist")[i];
              if ((a || o) && R) {
                t.on(
                  "change:position",
                  function e(i, r) {
                    var s = a && !a.daiSetting,
                      l = t.get("duration");
                    s && r && l > 0 && r >= l - p.b
                      ? (t.off("change:position", e, n), B.backgroundLoad(a))
                      : o && (a = M.get("nextUp"));
                  },
                  this
                );
              }
            }),
            (y = new we(M)).on("all", L, C),
            I.on("viewSetup", function (e) {
              Object(a.b)(O, e);
            }),
            (this.playerReady = function () {
              v.once(d.hb, function () {
                try {
                  !(function () {
                    M.change("visibility", N),
                      P.off(),
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
                      P.flush(),
                      P.destroy(),
                      (P = null),
                      M.change("viewable", F),
                      M.change("viewable", W),
                      M.get("autoPause").viewability
                        ? M.change("viewable", Q)
                        : M.once(
                            "change:autostartFailed change:mute",
                            function (e) {
                              e.off("change:viewable", W);
                            }
                          );
                    H(),
                      M.on("change:itemReady", function (e, t) {
                        t && de.flush();
                      });
                  })();
                } catch (e) {
                  C.triggerError(Object(g.v)(g.m, g.a, e));
                }
              }),
                v.init();
            }),
            (this.preload = q),
            (this.load = function (e, t) {
              var n,
                i = C._instreamAdapter;
              switch (
                (i && (i.noResume = !0),
                C.trigger("destroyPlugin", {}),
                ee(!0),
                _.cancel(),
                (_ = j(H)),
                z.cancel(),
                vt() && b.prime(),
                St(e))
              ) {
                case "string":
                  (M.attributes.item = 0),
                    (M.attributes.itemReady = !1),
                    (z = j(function (e) {
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
                    })(e).then(z.async));
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
                C.triggerError(Object(g.u)(e, g.c));
              }),
                n.then(_.async).catch(Et);
            }),
            (this.play = function (e) {
              return J(e).catch(Et);
            }),
            (this.pause = ne),
            (this.seek = function (e, t) {
              var n = M.get("state");
              if (n !== d.lb) {
                B.position = e;
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
            (this.setCurrentCaptions = re),
            (this.setCurrentQuality = function (e) {
              B.quality = e;
            }),
            (this.setFullscreen = le),
            (this.getCurrentQuality = function () {
              return B.quality;
            }),
            (this.getQualityLevels = function () {
              return B.qualities;
            }),
            (this.setCurrentAudioTrack = function (e) {
              B.audioTrack = e;
            }),
            (this.getCurrentAudioTrack = function () {
              return B.audioTrack;
            }),
            (this.getAudioTracks = function () {
              return B.audioTracks;
            }),
            (this.getCurrentCaptions = function () {
              return y.getCurrentIndex();
            }),
            (this.getCaptionsList = se),
            (this.getVisualQuality = function () {
              var e = this._model.get("mediaModel");
              return e ? e.get(d.U) : null;
            }),
            (this.getConfig = function () {
              return this._model ? this._model.getConfiguration() : void 0;
            }),
            (this.getState = K),
            (this.next = Et),
            (this.completeHandler = ae),
            (this.completeCancelled = function () {
              return (
                ((e = M.get("state")) !== d.mb && e !== d.kb && e !== d.lb) ||
                (!!T && ((T = !1), !0))
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
                  ((t.height = Object(r.b)(t.height)),
                  (t.width = t.width || i.width)),
                  t.width &&
                    ((t.width = Object(r.b)(t.width)),
                    t.aspectratio
                      ? ((i.width = t.width), delete t.width)
                      : (t.height = i.height)),
                  t.width &&
                    t.height &&
                    !t.aspectratio &&
                    e._view.resize(t.width, t.height),
                  Object.keys(t).forEach(function (o) {
                    var a = t[o];
                    if (void 0 !== a)
                      switch (o) {
                        case "aspectratio":
                          n.set(o, Object(r.a)(a, i.width));
                          break;
                        case "autostart":
                          !(function (e, t, n) {
                            e.setAutoStart(n),
                              "idle" === e.get("state") &&
                                !0 === n &&
                                t.play({ reason: "autostart" });
                          })(n, e, a);
                          break;
                        case "mute":
                          e.setMute(a);
                          break;
                        case "volume":
                          e.setVolume(a);
                          break;
                        case "playbackRateControls":
                        case "playbackRates":
                        case "repeat":
                        case "stretching":
                          n.set(o, a);
                      }
                  });
              })(C, e);
            }),
            (this.setItemIndex = function (e) {
              B.stopVideo();
              var t = M.get("playlist").length;
              return (
                (e = (parseInt(e, 10) || 0) % t) < 0 && (e += t),
                B.setActiveItem(e).catch(function (e) {
                  e.code >= 151 && e.code <= 162 && (e = Object(g.u)(e, g.e)),
                    O.triggerError(Object(g.v)(g.k, g.d, e));
                })
              );
            }),
            (this.detachMedia = function () {
              if (
                (S && (x = !0),
                M.get("autoPause").viewability && Q(M, M.get("viewable")),
                !R)
              )
                return B.setAttached(!1);
              B.backgroundActiveMedia();
            }),
            (this.attachMedia = function () {
              R ? B.restoreBackgroundMedia() : B.setAttached(!0),
                "function" == typeof k && k();
            }),
            (this.routeEvents = function (e) {
              return B.routeEvents(e);
            }),
            (this.forwardEvents = function () {
              return B.forwardEvents();
            }),
            (this.playVideo = function (e) {
              return B.playVideo(e);
            }),
            (this.stopVideo = function () {
              return B.stopVideo();
            }),
            (this.castVideo = function (e, t) {
              return B.castVideo(e, t);
            }),
            (this.stopCast = function () {
              return B.stopCast();
            }),
            (this.backgroundActiveMedia = function () {
              return B.backgroundActiveMedia();
            }),
            (this.restoreBackgroundMedia = function () {
              return B.restoreBackgroundMedia();
            }),
            (this.preloadNextItem = function () {
              B.background.currentMedia && B.preloadVideo();
            }),
            (this.isBeforeComplete = function () {
              return B.beforeComplete;
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
              var a = M.get("customButtons") || [],
                r = !1,
                s = { img: e, tooltip: t, callback: n, id: i, btnClass: o };
              (a = a.reduce(function (e, t) {
                return t.id === i ? ((r = !0), e.push(s)) : e.push(t), e;
              }, [])),
                r || a.unshift(s),
                M.set("customButtons", a);
            }),
            (this.removeButton = function (e) {
              var t = M.get("customButtons") || [];
              (t = t.filter(function (t) {
                return t.id !== e;
              })),
                M.set("customButtons", t);
            }),
            (this.resize = v.resize),
            (this.getSafeRegion = v.getSafeRegion),
            (this.setCaptions = v.setCaptions),
            (this.checkBeforePlay = function () {
              return S;
            }),
            (this.setControls = function (e) {
              Object(i.n)(e) || (e = !M.get("controls")),
                M.set("controls", e),
                (B.controls = e);
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
                Object(a.b)(this, this.originalContainer),
                v && v.destroy(),
                M && M.destroy(),
                de && de.destroy(),
                y && y.destroy(),
                B && B.destroy(),
                this.instreamDestroy();
            }),
            (this.isBeforePlay = this.checkBeforePlay),
            (this.createInstream = function () {
              return (
                this.instreamDestroy(),
                (this._instreamAdapter = new ce(this, M, v, b)),
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
              return !O._model.get("itemReady") || P;
            }
          );
          de.queue.push.apply(de.queue, h), v.setup();
        },
        get: function (e) {
          if (e in y.a) {
            var t = this._model.get("mediaModel");
            return t ? t.get(e) : y.a[e];
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
      t.default = Tt;
    },
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    function (e, t, n) {
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
            var r = a(n[o]);
            r.text && t.push(r);
          }
        return t;
      }
      function a(e) {
        var t = {},
          n = e.split("\r\n");
        1 === n.length && (n = e.split("\n"));
        var o = 1;
        if (
          (n[0].indexOf(" --\x3e ") > 0 && (o = 0),
          n.length > o + 1 && n[o + 1])
        ) {
          var a = n[o],
            r = a.indexOf(" --\x3e ");
          r > 0 &&
            ((t.begin = Object(i.g)(a.substr(0, r))),
            (t.end = Object(i.g)(a.substr(r + 5))),
            (t.text = n.slice(o + 1).join("\r\n")));
        }
        return t;
      }
    },
    function (e, t, n) {
      "use strict";
      n.d(t, "a", function () {
        return o;
      }),
        n.d(t, "b", function () {
          return a;
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
      function a(e, t) {
        var n = "jw-breakpoint-" + t;
        Object(i.p)(e, /jw-breakpoint--?\d+/, n);
      }
    },
    function (e, t, n) {
      "use strict";
      n.d(t, "a", function () {
        return d;
      });
      var i,
        o = n(0),
        a = n(8),
        r = n(16),
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
        p = function (e) {
          var t,
            s,
            p,
            w,
            f,
            h,
            j,
            g,
            b,
            m = this,
            v = e.player;
          function y() {
            Object(o.o)(t.fontSize) &&
              (v.get("containerHeight")
                ? (g =
                    (d.fontScale * (t.userFontScale || 1) * t.fontSize) /
                    d.fontSize)
                : v.once("change:containerHeight", y, this));
          }
          function k() {
            var e = v.get("containerHeight");
            if (e) {
              var t;
              if (v.get("fullscreen") && a.OS.iOS) t = null;
              else {
                var n = e * g;
                t =
                  Math.round(
                    10 *
                      (function (e) {
                        var t = v.get("mediaElement");
                        if (t && t.videoHeight) {
                          var n = t.videoWidth,
                            i = t.videoHeight,
                            o = n / i,
                            r = v.get("containerHeight"),
                            s = v.get("containerWidth");
                          if (v.get("fullscreen") && a.OS.mobile) {
                            var l = window.screen;
                            l.orientation &&
                              ((r = l.availHeight), (s = l.availWidth));
                          }
                          if (s && r && n && i)
                            return (s / r > o ? r : (i * s) / n) * g;
                        }
                        return e;
                      })(n)
                  ) / 10;
              }
              v.get("renderCaptionsNatively")
                ? (function (e, t) {
                    var n = "#".concat(
                      e,
                      " .jw-video::-webkit-media-text-track-display"
                    );
                    t &&
                      ((t += "px"),
                      a.OS.iOS &&
                        Object(c.b)(n, { fontSize: "inherit" }, e, !0));
                    (b.fontSize = t), Object(c.b)(n, b, e, !0);
                  })(v.get("id"), t)
                : Object(c.d)(f, { fontSize: t });
            }
          }
          function x(e, t, n) {
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
          ((f = document.createElement("div")).className =
            "jw-captions jw-reset"),
            (this.show = function () {
              Object(u.a)(f, "jw-captions-enabled");
            }),
            (this.hide = function () {
              Object(u.o)(f, "jw-captions-enabled");
            }),
            (this.populate = function (e) {
              v.get("renderCaptionsNatively") ||
                ((p = []),
                (s = e),
                e ? this.selectCues(e, w) : this.renderCues());
            }),
            (this.resize = function () {
              k(), this.renderCues(!0);
            }),
            (this.renderCues = function (e) {
              (e = !!e), i && i.processCues(window, p, f, e);
            }),
            (this.selectCues = function (e, t) {
              if (e && e.data && t && !v.get("renderCaptionsNatively")) {
                var n = this.getAlignmentPosition(e, t);
                !1 !== n &&
                  ((p = this.getCurrentCues(e.data, n)), this.renderCues(!0));
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
                a = t.currentTime;
              return n && i && Object(o.r)(i[n]) && (a = i[n]), a;
            }),
            (this.clear = function () {
              Object(u.g)(f);
            }),
            (this.setup = function (e, n) {
              (h = document.createElement("div")),
                (j = document.createElement("span")),
                (h.className = "jw-captions-window jw-reset"),
                (j.className = "jw-captions-text jw-reset"),
                (t = Object(o.g)({}, d, n)),
                (g = d.fontScale);
              var i = function () {
                if (!v.get("renderCaptionsNatively")) {
                  y(t.fontSize);
                  var n = t.windowColor,
                    i = t.windowOpacity,
                    o = t.edgeStyle;
                  b = {};
                  var r = {};
                  !(function (e, t) {
                    var n = t.color,
                      i = t.fontOpacity;
                    (n || i !== d.fontOpacity) &&
                      (e.color = Object(c.c)(n || "#ffffff", i));
                    if (t.back) {
                      var o = t.backgroundColor,
                        a = t.backgroundOpacity;
                      (o === d.backgroundColor && a === d.backgroundOpacity) ||
                        (e.backgroundColor = Object(c.c)(o, a));
                    } else e.background = "transparent";
                    t.fontFamily && (e.fontFamily = t.fontFamily);
                    t.fontStyle && (e.fontStyle = t.fontStyle);
                    t.fontWeight && (e.fontWeight = t.fontWeight);
                    t.textDecoration && (e.textDecoration = t.textDecoration);
                  })(r, t),
                    (n || i !== d.windowOpacity) &&
                      (b.backgroundColor = Object(c.c)(n || "#000000", i)),
                    x(o, r, t.fontOpacity),
                    t.back || null !== o || x("uniform", r),
                    Object(c.d)(h, b),
                    Object(c.d)(j, r),
                    (function (e, t) {
                      k(),
                        (function (e, t) {
                          a.Browser.safari &&
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
                            b,
                            e,
                            !0
                          ),
                            Object(c.b)("#" + e + " .jw-video::cue", t, e, !0);
                        })(e, t),
                        (function (e, t) {
                          Object(c.b)(
                            "#" + e + " .jw-text-track-display",
                            b,
                            e
                          ),
                            Object(c.b)("#" + e + " .jw-text-track-cue", t, e);
                        })(e, t);
                    })(e, r);
                }
              };
              i(),
                h.appendChild(j),
                f.appendChild(h),
                v.change(
                  "captionsTrack",
                  function (e, t) {
                    this.populate(t);
                  },
                  this
                ),
                v.set("captions", t),
                v.on("change:captions", function (e, n) {
                  (t = n), i();
                });
            }),
            (this.element = function () {
              return f;
            }),
            (this.destroy = function () {
              v.off(null, null, this), this.off();
            });
          var O = function (e) {
            (w = e), m.selectCues(s, w);
          };
          v.on(
            "change:playlistItem",
            function () {
              (w = null), (p = []);
            },
            this
          ),
            v.on(
              l.Q,
              function (e) {
                (p = []), O(e);
              },
              this
            ),
            v.on(l.S, O, this),
            v.on(
              "subtitlesTrackData",
              function () {
                this.selectCues(s, w);
              },
              this
            ),
            v.on(
              "change:captionsList",
              function e(t, o) {
                var a = this;
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
                      .catch(Object(r.c)(301121))
                      .catch(function (e) {
                        a.trigger(l.tb, e);
                      }),
                    t.off("change:captionsList", e, this)));
              },
              this
            );
        };
      Object(o.g)(p.prototype, s.a), (t.b = p);
    },
    function (e, t, n) {
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
                      ((r = i),
                      "/*# sourceMappingURL=data:application/json;charset=utf-8;base64," +
                        btoa(unescape(encodeURIComponent(JSON.stringify(r)))) +
                        " */"),
                    a = i.sources.map(function (e) {
                      return "/*# sourceURL=" + i.sourceRoot + e + " */";
                    });
                  return [n].concat(a).concat([o]).join("\n");
                }
                var r;
                return [n].join("\n");
              })(t, e);
              return t[2] ? "@media " + t[2] + "{" + n + "}" : n;
            }).join("");
          }),
          (t.i = function (e, n) {
            "string" == typeof e && (e = [[null, e, ""]]);
            for (var i = {}, o = 0; o < this.length; o++) {
              var a = this[o][0];
              null != a && (i[a] = !0);
            }
            for (o = 0; o < e.length; o++) {
              var r = e[o];
              (null != r[0] && i[r[0]]) ||
                (n && !r[2]
                  ? (r[2] = n)
                  : n && (r[2] = "(" + r[2] + ") and (" + n + ")"),
                t.push(r));
            }
          }),
          t
        );
      };
    },
    function (e, t) {
      var n,
        i,
        o = {},
        a = {},
        r =
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
            r().appendChild(e);
          })(t),
          t
        );
      }
      function l(e, t) {
        var n,
          i,
          o,
          r = a[e];
        r || (r = a[e] = { element: s(e), counter: 0 });
        var l = r.counter++;
        return (
          (n = r.element),
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
                a = (o[e] || {})[i.id];
              if (a) {
                for (var r = 0; r < a.parts.length; r++) a.parts[r](i.parts[r]);
                for (; r < i.parts.length; r++) a.parts.push(l(e, i.parts[r]));
              } else {
                var s = [];
                for (r = 0; r < i.parts.length; r++) s.push(l(e, i.parts[r]));
                (o[e] = o[e] || {}), (o[e][i.id] = { id: i.id, parts: s });
              }
            }
          })(
            t,
            (function (e) {
              for (var t = [], n = {}, i = 0; i < e.length; i++) {
                var o = e[i],
                  a = o[0],
                  r = o[1],
                  s = o[2],
                  l = { css: r, media: s };
                n[a]
                  ? n[a].parts.push(l)
                  : t.push((n[a] = { id: a, parts: [l] }));
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
            if (i) for (var a = 0; a < i.parts.length; a += 1) i.parts[a]();
            return;
          }
          for (var r = Object.keys(n), s = 0; s < r.length; s += 1)
            for (var l = n[r[s]], c = 0; c < l.parts.length; c += 1)
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
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-arrow-right" viewBox="0 0 240 240" focusable="false"><path d="M183.6,104.4L81.8,0L45.4,36.3l84.9,84.9l-84.9,84.9L79.3,240l101.9-101.7c9.9-6.9,12.4-20.4,5.5-30.4C185.8,106.7,184.8,105.4,183.6,104.4L183.6,104.4z"></path></svg>';
    },
    function (e, t, n) {
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
    function (e, t, n) {
      "use strict";
      function i(e) {
        return new Promise(function (t, n) {
          if (e.paused) return n(o("NotAllowedError", 0, "play() failed."));
          var i = function () {
              e.removeEventListener("play", a),
                e.removeEventListener("playing", r),
                e.removeEventListener("pause", r),
                e.removeEventListener("abort", r),
                e.removeEventListener("error", r);
            },
            a = function () {
              e.addEventListener("playing", r),
                e.addEventListener("abort", r),
                e.addEventListener("error", r),
                e.addEventListener("pause", r);
            },
            r = function (e) {
              if ((i(), "playing" === e.type)) t();
              else {
                var a = 'The play() request was interrupted by a "'.concat(
                  e.type,
                  '" event.'
                );
                "error" === e.type
                  ? n(o("NotSupportedError", 9, a))
                  : n(o("AbortError", 20, a));
              }
            };
          e.addEventListener("play", a);
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
    function (e, t, n) {
      "use strict";
      function i(e, t) {
        return e !== 1 / 0 && Math.abs(e) >= Math.max(a(t), 0);
      }
      function o(e, t) {
        var n = "VOD";
        return (
          e === 1 / 0
            ? (n = "LIVE")
            : e < 0 && (n = i(e, a(t)) ? "DVR" : "LIVE"),
          n
        );
      }
      function a(e) {
        return void 0 === e ? 120 : Math.max(e, 0);
      }
      n.d(t, "a", function () {
        return i;
      }),
        n.d(t, "b", function () {
          return o;
        });
    },
    function (e, t, n) {
      "use strict";
      var i = n(67),
        o = n(16),
        a = n(22),
        r = n(4),
        s = n(57),
        l = n(2),
        c = n(1);
      function u(e) {
        throw new c.n(null, e);
      }
      function d(e, t, i) {
        e.xhr = Object(a.a)(
          e.file,
          function (a) {
            !(function (e, t, i, a) {
              var d,
                p,
                f = e.responseXML ? e.responseXML.firstChild : null;
              if (f)
                for (
                  "xml" === Object(r.b)(f) && (f = f.nextSibling);
                  f.nodeType === f.COMMENT_NODE;

                )
                  f = f.nextSibling;
              try {
                if (f && "tt" === Object(r.b)(f))
                  (d = (function (e) {
                    e || u(306007);
                    var t = [],
                      n = e.getElementsByTagName("p"),
                      i = 30,
                      o = e.getElementsByTagName("tt");
                    if (o && o[0]) {
                      var a = parseFloat(o[0].getAttribute("ttp:frameRate"));
                      isNaN(a) || (i = a);
                    }
                    n || u(306005),
                      n.length ||
                        (n = e.getElementsByTagName("tt:p")).length ||
                        (n = e.getElementsByTagName("tts:p"));
                    for (var r = 0; r < n.length; r++) {
                      for (
                        var s = n[r], c = s.getElementsByTagName("br"), d = 0;
                        d < c.length;
                        d++
                      ) {
                        var p = c[d];
                        p.parentNode.replaceChild(e.createTextNode("\r\n"), p);
                      }
                      var w = s.innerHTML || s.textContent || s.text || "",
                        f = Object(l.i)(w)
                          .replace(/>\s+</g, "><")
                          .replace(/(<\/?)tts?:/g, "$1")
                          .replace(/<br.*?\/>/g, "\r\n");
                      if (f) {
                        var h = s.getAttribute("begin"),
                          j = s.getAttribute("dur"),
                          g = s.getAttribute("end"),
                          b = { begin: Object(l.g)(h, i), text: f };
                        g
                          ? (b.end = Object(l.g)(g, i))
                          : j && (b.end = b.begin + Object(l.g)(j, i)),
                          t.push(b);
                      }
                    }
                    return t.length || u(306005), t;
                  })(e.responseXML)),
                    (p = w(d)),
                    delete t.xhr,
                    i(p);
                else {
                  var h = e.responseText;
                  h.indexOf("WEBVTT") >= 0
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
                          (p = []),
                            (n.oncue = function (e) {
                              p.push(e);
                            }),
                            (n.onflush = function () {
                              delete t.xhr, i(p);
                            }),
                            n.parse(h);
                        })
                        .catch(function (e) {
                          delete t.xhr, a(Object(c.v)(null, c.b, e));
                        })
                    : ((d = Object(s.a)(h)), (p = w(d)), delete t.xhr, i(p));
                }
              } catch (e) {
                delete t.xhr, a(Object(c.v)(null, c.b, e));
              }
            })(a, e, t, i);
          },
          function (e, t, n, o) {
            i(Object(c.u)(o, c.b));
          }
        );
      }
      function p(e) {
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
      function w(e) {
        return e.map(function (e) {
          return new i.a(e.begin, e.end, e.text);
        });
      }
      n.d(t, "c", function () {
        return d;
      }),
        n.d(t, "a", function () {
          return p;
        }),
        n.d(t, "b", function () {
          return w;
        });
    },
    function (e, t, n) {
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
          var a = "",
            r = !1,
            s = e,
            l = t,
            c = n,
            u = null,
            d = "",
            p = !0,
            w = "auto",
            f = "start",
            h = "auto",
            j = 100,
            g = "middle";
          Object.defineProperty(i, "id", {
            enumerable: !0,
            get: function () {
              return a;
            },
            set: function (e) {
              a = "" + e;
            },
          }),
            Object.defineProperty(i, "pauseOnExit", {
              enumerable: !0,
              get: function () {
                return r;
              },
              set: function (e) {
                r = !!e;
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
                return p;
              },
              set: function (e) {
                (p = !!e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "line", {
              enumerable: !0,
              get: function () {
                return w;
              },
              set: function (e) {
                if ("number" != typeof e && "auto" !== e)
                  throw new SyntaxError(
                    "An invalid number or illegal string was specified."
                  );
                (w = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "lineAlign", {
              enumerable: !0,
              get: function () {
                return f;
              },
              set: function (e) {
                var t = o(e);
                if (!t)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (f = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "position", {
              enumerable: !0,
              get: function () {
                return h;
              },
              set: function (e) {
                if (e < 0 || e > 100)
                  throw new Error("Position must be between 0 and 100.");
                (h = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "size", {
              enumerable: !0,
              get: function () {
                return j;
              },
              set: function (e) {
                if (e < 0 || e > 100)
                  throw new Error("Size must be between 0 and 100.");
                (j = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "align", {
              enumerable: !0,
              get: function () {
                return g;
              },
              set: function (e) {
                var t = o(e);
                if (!t)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (g = t), (this.hasBeenReset = !0);
              },
            }),
            (i.displayState = void 0);
        }).prototype.getCueAsHTML = function () {
          return window.WebVTT.convertCueToDOMTree(window, this.text);
        };
      }
      t.a = i;
    },
    ,
    function (e, t, n) {
      var i = n(70);
      "string" == typeof i && (i = [["all-players", i, ""]]),
        n(61).style(i, "all-players"),
        i.locals && (e.exports = i.locals);
    },
    function (e, t, n) {
      (e.exports = n(60)(!1)).push([
        e.i,
        '.jw-reset{text-align:left;direction:ltr}.jw-reset-text,.jw-reset{color:inherit;background-color:transparent;padding:0;margin:0;float:none;font-family:Arial,Helvetica,sans-serif;font-size:1em;line-height:1em;list-style:none;text-transform:none;vertical-align:baseline;border:0;font-variant:inherit;font-stretch:inherit;-webkit-tap-highlight-color:rgba(255,255,255,0)}body .jw-error,body .jwplayer.jw-state-error{height:100%;width:100%}.jw-title{position:absolute;top:0}.jw-background-color{background:rgba(0,0,0,0.4)}.jw-text{color:rgba(255,255,255,0.8)}.jw-knob{color:rgba(255,255,255,0.8);background-color:#fff}.jw-button-color{color:rgba(255,255,255,0.8)}:not(.jw-flag-touch) .jw-button-color:not(.jw-logo-button):focus,:not(.jw-flag-touch) .jw-button-color:not(.jw-logo-button):hover{color:#fff}.jw-toggle{color:#fff}.jw-toggle.jw-off{color:rgba(255,255,255,0.8)}.jw-toggle.jw-off:focus{color:#fff}.jw-toggle:focus{outline:none}:not(.jw-flag-touch) .jw-toggle.jw-off:hover{color:#fff}.jw-rail{background:rgba(255,255,255,0.3)}.jw-buffer{background:rgba(255,255,255,0.3)}.jw-progress{background:#f2f2f2}.jw-time-tip,.jw-volume-tip{border:0}.jw-slider-volume.jw-volume-tip.jw-background-color.jw-slider-vertical{background:none}.jw-skip{padding:.5em;outline:none}.jw-skip .jw-skiptext,.jw-skip .jw-skip-icon{color:rgba(255,255,255,0.8)}.jw-skip.jw-skippable:hover .jw-skip-icon,.jw-skip.jw-skippable:focus .jw-skip-icon{color:#fff}.jw-icon-cast google-cast-launcher{--connected-color:#fff;--disconnected-color:rgba(255,255,255,0.8)}.jw-icon-cast google-cast-launcher:focus{outline:none}.jw-icon-cast google-cast-launcher.jw-off{--connected-color:rgba(255,255,255,0.8)}.jw-icon-cast:focus google-cast-launcher{--connected-color:#fff;--disconnected-color:#fff}.jw-icon-cast:hover google-cast-launcher{--connected-color:#fff;--disconnected-color:#fff}.jw-nextup-container{bottom:2.5em;padding:5px .5em}.jw-nextup{border-radius:0}.jw-color-active{color:#fff;stroke:#fff;border-color:#fff}:not(.jw-flag-touch) .jw-color-active-hover:hover,:not(.jw-flag-touch) .jw-color-active-hover:focus{color:#fff;stroke:#fff;border-color:#fff}.jw-color-inactive{color:rgba(255,255,255,0.8);stroke:rgba(255,255,255,0.8);border-color:rgba(255,255,255,0.8)}:not(.jw-flag-touch) .jw-color-inactive-hover:hover{color:rgba(255,255,255,0.8);stroke:rgba(255,255,255,0.8);border-color:rgba(255,255,255,0.8)}.jw-option{color:rgba(255,255,255,0.8)}.jw-option.jw-active-option{color:#fff;background-color:rgba(255,255,255,0.1)}:not(.jw-flag-touch) .jw-option:hover{color:#fff}.jwplayer{width:100%;font-size:16px;position:relative;display:block;min-height:0;overflow:hidden;box-sizing:border-box;font-family:Arial,Helvetica,sans-serif;-webkit-touch-callout:none;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none;outline:none}.jwplayer *{box-sizing:inherit}.jwplayer.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jwplayer.jw-flag-aspect-mode{height:auto !important}.jwplayer.jw-flag-aspect-mode .jw-aspect{display:block}.jwplayer .jw-aspect{display:none}.jwplayer .jw-swf{outline:none}.jw-media,.jw-preview{position:absolute;width:100%;height:100%;top:0;left:0;bottom:0;right:0}.jw-media{overflow:hidden;cursor:pointer}.jw-plugin{position:absolute;bottom:66px}.jw-breakpoint-7 .jw-plugin{bottom:132px}.jw-plugin .jw-banner{max-width:100%;opacity:0;cursor:pointer;position:absolute;margin:auto auto 0;left:0;right:0;bottom:0;display:block}.jw-preview,.jw-captions,.jw-title{pointer-events:none}.jw-media,.jw-logo{pointer-events:all}.jw-wrapper{background-color:#000;position:absolute;top:0;left:0;right:0;bottom:0}.jw-hidden-accessibility{border:0;clip:rect(0 0 0 0);height:1px;margin:-1px;overflow:hidden;padding:0;position:absolute;width:1px}.jw-contract-trigger::before{content:"";overflow:hidden;width:200%;height:200%;display:block;position:absolute;top:0;left:0}.jwplayer .jw-media video{position:absolute;top:0;right:0;bottom:0;left:0;width:100%;height:100%;margin:auto;background:transparent}.jwplayer .jw-media video::-webkit-media-controls-start-playback-button{display:none}.jwplayer.jw-stretch-uniform .jw-media video{object-fit:contain}.jwplayer.jw-stretch-none .jw-media video{object-fit:none}.jwplayer.jw-stretch-fill .jw-media video{object-fit:cover}.jwplayer.jw-stretch-exactfit .jw-media video{object-fit:fill}.jw-preview{position:absolute;display:none;opacity:1;visibility:visible;width:100%;height:100%;background:#000 no-repeat 50% 50%}.jwplayer .jw-preview,.jw-error .jw-preview{background-size:contain}.jw-stretch-none .jw-preview{background-size:auto auto}.jw-stretch-fill .jw-preview{background-size:cover}.jw-stretch-exactfit .jw-preview{background-size:100% 100%}.jw-title{display:none;padding-top:20px;width:100%;z-index:1}.jw-title-primary,.jw-title-secondary{color:#fff;padding-left:20px;padding-right:20px;padding-bottom:.5em;overflow:hidden;text-overflow:ellipsis;direction:unset;white-space:nowrap;width:100%}.jw-title-primary{font-size:1.625em}.jw-breakpoint-2 .jw-title-primary,.jw-breakpoint-3 .jw-title-primary{font-size:1.5em}.jw-flag-small-player .jw-title-primary{font-size:1.25em}.jw-flag-small-player .jw-title-secondary,.jw-title-secondary:empty{display:none}.jw-captions{position:absolute;width:100%;height:100%;text-align:center;display:none;letter-spacing:normal;word-spacing:normal;text-transform:none;text-indent:0;text-decoration:none;pointer-events:none;overflow:hidden;top:0}.jw-captions.jw-captions-enabled{display:block}.jw-captions-window{display:none;padding:.25em;border-radius:.25em}.jw-captions-window.jw-captions-window-active{display:inline-block}.jw-captions-text{display:inline-block;color:#fff;background-color:#000;word-wrap:normal;word-break:normal;white-space:pre-line;font-style:normal;font-weight:normal;text-align:center;text-decoration:none}.jw-text-track-display{font-size:inherit;line-height:1.5}.jw-text-track-cue{background-color:rgba(0,0,0,0.5);color:#fff;padding:.1em .3em}.jwplayer video::-webkit-media-controls{display:none;justify-content:flex-start}.jwplayer video::-webkit-media-text-track-display{min-width:-webkit-min-content}.jwplayer video::cue{background-color:rgba(0,0,0,0.5)}.jwplayer video::-webkit-media-controls-panel-container{display:none}.jwplayer:not(.jw-flag-controls-hidden):not(.jw-state-playing) .jw-captions,.jwplayer.jw-flag-media-audio.jw-state-playing .jw-captions,.jwplayer.jw-state-playing:not(.jw-flag-user-inactive):not(.jw-flag-controls-hidden) .jw-captions{max-height:calc(100% - 60px)}.jwplayer:not(.jw-flag-controls-hidden):not(.jw-state-playing):not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container,.jwplayer.jw-flag-media-audio.jw-state-playing:not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container,.jwplayer.jw-state-playing:not(.jw-flag-user-inactive):not(.jw-flag-controls-hidden):not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container{max-height:calc(100% - 60px)}.jw-logo{position:absolute;margin:20px;cursor:pointer;pointer-events:all;background-repeat:no-repeat;background-size:contain;top:auto;right:auto;left:auto;bottom:auto;outline:none}.jw-logo.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jw-flag-audio-player .jw-logo{display:none}.jw-logo-top-right{top:0;right:0}.jw-logo-top-left{top:0;left:0}.jw-logo-bottom-left{left:0}.jw-logo-bottom-right{right:0}.jw-logo-bottom-left,.jw-logo-bottom-right{bottom:44px;transition:bottom 150ms cubic-bezier(0, .25, .25, 1)}.jw-state-idle .jw-logo{z-index:1}.jw-state-setup .jw-wrapper{background-color:inherit}.jw-state-setup .jw-logo,.jw-state-setup .jw-controls,.jw-state-setup .jw-controls-backdrop{visibility:hidden}span.jw-break{display:block}body .jw-error,body .jwplayer.jw-state-error{background-color:#333;color:#fff;font-size:16px;display:table;opacity:1;position:relative}body .jw-error .jw-display,body .jwplayer.jw-state-error .jw-display{display:none}body .jw-error .jw-media,body .jwplayer.jw-state-error .jw-media{cursor:default}body .jw-error .jw-preview,body .jwplayer.jw-state-error .jw-preview{background-color:#333}body .jw-error .jw-error-msg,body .jwplayer.jw-state-error .jw-error-msg{background-color:#000;border-radius:2px;display:flex;flex-direction:row;align-items:stretch;padding:20px}body .jw-error .jw-error-msg .jw-icon,body .jwplayer.jw-state-error .jw-error-msg .jw-icon{height:30px;width:30px;margin-right:20px;flex:0 0 auto;align-self:center}body .jw-error .jw-error-msg .jw-icon:empty,body .jwplayer.jw-state-error .jw-error-msg .jw-icon:empty{display:none}body .jw-error .jw-error-msg .jw-info-container,body .jwplayer.jw-state-error .jw-error-msg .jw-info-container{margin:0;padding:0}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg{flex-direction:column}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-error-text,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-error-text,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-error-text,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-error-text{text-align:center}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-icon,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-icon,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-icon,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-icon{flex:.5 0 auto;margin-right:0;margin-bottom:20px}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg .jw-break,.jwplayer.jw-state-error.jw-flag-small-player .jw-error-msg .jw-break,.jwplayer.jw-state-error.jw-breakpoint-2 .jw-error-msg .jw-break{display:inline}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg .jw-break:before,.jwplayer.jw-state-error.jw-flag-small-player .jw-error-msg .jw-break:before,.jwplayer.jw-state-error.jw-breakpoint-2 .jw-error-msg .jw-break:before{content:" "}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg{height:100%;width:100%;top:0;position:absolute;left:0;background:#000;-webkit-transform:none;transform:none;padding:4px 16px;z-index:1}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg.jw-info-overlay{max-width:none;max-height:none}body .jwplayer.jw-state-error .jw-title,.jw-state-idle .jw-title,.jwplayer.jw-state-complete:not(.jw-flag-casting):not(.jw-flag-audio-player):not(.jw-flag-overlay-open-related) .jw-title{display:block}body .jwplayer.jw-state-error .jw-preview,.jw-state-idle .jw-preview,.jwplayer.jw-state-complete:not(.jw-flag-casting):not(.jw-flag-audio-player):not(.jw-flag-overlay-open-related) .jw-preview{display:block}.jw-state-idle .jw-captions,.jwplayer.jw-state-complete .jw-captions,body .jwplayer.jw-state-error .jw-captions{display:none}.jw-state-idle video::-webkit-media-text-track-container,.jwplayer.jw-state-complete video::-webkit-media-text-track-container,body .jwplayer.jw-state-error video::-webkit-media-text-track-container{display:none}.jwplayer.jw-flag-fullscreen{width:100% !important;height:100% !important;top:0;right:0;bottom:0;left:0;z-index:1000;margin:0;position:fixed}body .jwplayer.jw-flag-flash-blocked .jw-title{display:block}.jwplayer.jw-flag-controls-hidden .jw-media{cursor:default}.jw-flag-audio-player:not(.jw-flag-flash-blocked) .jw-media{visibility:hidden}.jw-flag-audio-player .jw-title{background:none}.jw-flag-audio-player object{min-height:45px}.jw-flag-floating{background-size:cover;background-color:#000}.jw-flag-floating .jw-wrapper{position:fixed;z-index:2147483647;-webkit-animation:jw-float-to-bottom 150ms cubic-bezier(0, .25, .25, 1) forwards 1;animation:jw-float-to-bottom 150ms cubic-bezier(0, .25, .25, 1) forwards 1;top:auto;bottom:1rem;left:auto;right:1rem;max-width:400px;max-height:400px;margin:0 auto}@media screen and (max-width:480px){.jw-flag-floating .jw-wrapper{width:100%;left:0;right:0}}.jw-flag-floating .jw-wrapper .jw-media{touch-action:none}@media screen and (max-device-width:480px) and (orientation:portrait){.jw-flag-touch.jw-flag-floating .jw-wrapper{-webkit-animation:none;animation:none;top:62px;bottom:auto;left:0;right:0;max-width:none;max-height:none}}.jw-flag-floating .jw-float-icon{pointer-events:all;cursor:pointer;display:none}.jw-flag-floating .jw-float-icon .jw-svg-icon{-webkit-filter:drop-shadow(0 0 1px #000);filter:drop-shadow(0 0 1px #000)}.jw-flag-floating.jw-floating-dismissible .jw-dismiss-icon{display:none}.jw-flag-floating.jw-floating-dismissible.jw-flag-ads .jw-float-icon{display:flex}.jw-flag-floating.jw-floating-dismissible.jw-state-paused .jw-logo,.jw-flag-floating.jw-floating-dismissible:not(.jw-flag-user-inactive) .jw-logo{display:none}.jw-flag-floating.jw-floating-dismissible.jw-state-paused .jw-float-icon,.jw-flag-floating.jw-floating-dismissible:not(.jw-flag-user-inactive) .jw-float-icon{display:flex}.jw-float-icon{display:none;position:absolute;top:3px;right:5px;align-items:center;justify-content:center}@-webkit-keyframes jw-float-to-bottom{from{-webkit-transform:translateY(100%);transform:translateY(100%)}to{-webkit-transform:translateY(0);transform:translateY(0)}}@keyframes jw-float-to-bottom{from{-webkit-transform:translateY(100%);transform:translateY(100%)}to{-webkit-transform:translateY(0);transform:translateY(0)}}.jw-flag-top{margin-top:2em;overflow:visible}.jw-top{height:2em;line-height:2;pointer-events:none;text-align:center;opacity:.8;position:absolute;top:-2em;width:100%}.jw-top .jw-icon{cursor:pointer;pointer-events:all;height:auto;width:auto}.jw-top .jw-text{color:#555}',
        "",
      ]);
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-buffer" viewBox="0 0 240 240" focusable="false"><path d="M120,186.667a66.667,66.667,0,0,1,0-133.333V40a80,80,0,1,0,80,80H186.667A66.846,66.846,0,0,1,120,186.667Z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-replay" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M120,41.9v-20c0-5-4-8-8-4l-44,28a5.865,5.865,0,0,0-3.3,7.6A5.943,5.943,0,0,0,68,56.8l43,29c5,4,9,1,9-4v-20a60,60,0,1,1-60,60H40a80,80,0,1,0,80-79.9Z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-error" viewBox="0 0 36 36" style="width:100%;height:100%;" focusable="false"><path d="M34.6 20.2L10 33.2 27.6 16l7 3.7a.4.4 0 0 1 .2.5.4.4 0 0 1-.2.2zM33.3 0L21 12.2 9 6c-.2-.3-.6 0-.6.5V25L0 33.6 2.5 36 36 2.7z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-play" viewBox="0 0 240 240" focusable="false"><path d="M62.8,199.5c-1,0.8-2.4,0.6-3.3-0.4c-0.4-0.5-0.6-1.1-0.5-1.8V42.6c-0.2-1.3,0.7-2.4,1.9-2.6c0.7-0.1,1.3,0.1,1.9,0.4l154.7,77.7c2.1,1.1,2.1,2.8,0,3.8L62.8,199.5z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-pause" viewBox="0 0 240 240" focusable="false"><path d="M100,194.9c0.2,2.6-1.8,4.8-4.4,5c-0.2,0-0.4,0-0.6,0H65c-2.6,0.2-4.8-1.8-5-4.4c0-0.2,0-0.4,0-0.6V45c-0.2-2.6,1.8-4.8,4.4-5c0.2,0,0.4,0,0.6,0h30c2.6-0.2,4.8,1.8,5,4.4c0,0.2,0,0.4,0,0.6V194.9z M180,45.1c0.2-2.6-1.8-4.8-4.4-5c-0.2,0-0.4,0-0.6,0h-30c-2.6-0.2-4.8,1.8-5,4.4c0,0.2,0,0.4,0,0.6V195c-0.2,2.6,1.8,4.8,4.4,5c0.2,0,0.4,0,0.6,0h30c2.6,0.2,4.8-1.8,5-4.4c0-0.2,0-0.4,0-0.6V45.1z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-rewind" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M113.2,131.078a21.589,21.589,0,0,0-17.7-10.6,21.589,21.589,0,0,0-17.7,10.6,44.769,44.769,0,0,0,0,46.3,21.589,21.589,0,0,0,17.7,10.6,21.589,21.589,0,0,0,17.7-10.6,44.769,44.769,0,0,0,0-46.3Zm-17.7,47.2c-7.8,0-14.4-11-14.4-24.1s6.6-24.1,14.4-24.1,14.4,11,14.4,24.1S103.4,178.278,95.5,178.278Zm-43.4,9.7v-51l-4.8,4.8-6.8-6.8,13-13a4.8,4.8,0,0,1,8.2,3.4v62.7l-9.6-.1Zm162-130.2v125.3a4.867,4.867,0,0,1-4.8,4.8H146.6v-19.3h48.2v-96.4H79.1v19.3c0,5.3-3.6,7.2-8,4.3l-41.8-27.9a6.013,6.013,0,0,1-2.7-8,5.887,5.887,0,0,1,2.7-2.7l41.8-27.9c4.4-2.9,8-1,8,4.3v19.3H209.2A4.974,4.974,0,0,1,214.1,57.778Z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-next" viewBox="0 0 240 240" focusable="false"><path d="M165,60v53.3L59.2,42.8C56.9,41.3,55,42.3,55,45v150c0,2.7,1.9,3.8,4.2,2.2L165,126.6v53.3h20v-120L165,60L165,60z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-stop" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M190,185c0.2,2.6-1.8,4.8-4.4,5c-0.2,0-0.4,0-0.6,0H55c-2.6,0.2-4.8-1.8-5-4.4c0-0.2,0-0.4,0-0.6V55c-0.2-2.6,1.8-4.8,4.4-5c0.2,0,0.4,0,0.6,0h130c2.6-0.2,4.8,1.8,5,4.4c0,0.2,0,0.4,0,0.6V185z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-volume-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M116.4,42.8v154.5c0,2.8-1.7,3.6-3.8,1.7l-54.1-48.1H28.9c-2.8,0-5.2-2.3-5.2-5.2V94.2c0-2.8,2.3-5.2,5.2-5.2h29.6l54.1-48.1C114.6,39.1,116.4,39.9,116.4,42.8z M212.3,96.4l-14.6-14.6l-23.6,23.6l-23.6-23.6l-14.6,14.6l23.6,23.6l-23.6,23.6l14.6,14.6l23.6-23.6l23.6,23.6l14.6-14.6L188.7,120L212.3,96.4z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-volume-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M116.4,42.8v154.5c0,2.8-1.7,3.6-3.8,1.7l-54.1-48.1H28.9c-2.8,0-5.2-2.3-5.2-5.2V94.2c0-2.8,2.3-5.2,5.2-5.2h29.6l54.1-48.1C114.7,39.1,116.4,39.9,116.4,42.8z M178.2,120c0-22.7-18.5-41.2-41.2-41.2v20.6c11.4,0,20.6,9.2,20.6,20.6c0,11.4-9.2,20.6-20.6,20.6v20.6C159.8,161.2,178.2,142.7,178.2,120z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-volume-100" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M116.5,42.8v154.4c0,2.8-1.7,3.6-3.8,1.7l-54.1-48H29c-2.8,0-5.2-2.3-5.2-5.2V94.3c0-2.8,2.3-5.2,5.2-5.2h29.6l54.1-48C114.8,39.2,116.5,39.9,116.5,42.8z"></path><path d="M136.2,160v-20c11.1,0,20-8.9,20-20s-8.9-20-20-20V80c22.1,0,40,17.9,40,40S158.3,160,136.2,160z"></path><path d="M216.2,120c0-44.2-35.8-80-80-80v20c33.1,0,60,26.9,60,60s-26.9,60-60,60v20C180.4,199.9,216.1,164.1,216.2,120z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-cc-on" viewBox="0 0 240 240" focusable="false"><path d="M215,40H25c-2.7,0-5,2.2-5,5v150c0,2.7,2.2,5,5,5h190c2.7,0,5-2.2,5-5V45C220,42.2,217.8,40,215,40z M108.1,137.7c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9c-2.4-3.7-6.5-5.9-10.9-5.9c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6C90.4,141.7,102,143.5,108.1,137.7z M152.9,137.7c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9c-2.4-3.7-6.5-5.9-10.9-5.9c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6C135.2,141.7,146.8,143.5,152.9,137.7z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-cc-off" viewBox="0 0 240 240" focusable="false"><path d="M99.4,97.8c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6c0,9.6,11.6,11.4,17.7,5.5c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9C107.9,100,103.8,97.8,99.4,97.8z M144.1,97.8c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6c0,9.6,11.6,11.4,17.7,5.5c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9C152.6,100,148.5,97.8,144.1,97.8L144.1,97.8z M200,60v120H40V60H200 M215,40H25c-2.7,0-5,2.2-5,5v150c0,2.7,2.2,5,5,5h190c2.7,0,5-2.2,5-5V45C220,42.2,217.8,40,215,40z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-airplay-on" viewBox="0 0 240 240" focusable="false"><path d="M229.9,40v130c0.2,2.6-1.8,4.8-4.4,5c-0.2,0-0.4,0-0.6,0h-44l-17-20h46V55H30v100h47l-17,20h-45c-2.6,0.2-4.8-1.8-5-4.4c0-0.2,0-0.4,0-0.6V40c-0.2-2.6,1.8-4.8,4.4-5c0.2,0,0.4,0,0.6,0h209.8c2.6-0.2,4.8,1.8,5,4.4C229.9,39.7,229.9,39.9,229.9,40z M104.9,122l15-18l15,18l11,13h44V75H50v60h44L104.9,122z M179.9,205l-60-70l-60,70H179.9z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-airplay-off" viewBox="0 0 240 240" focusable="false"><path d="M210,55v100h-50l20,20h45c2.6,0.2,4.8-1.8,5-4.4c0-0.2,0-0.4,0-0.6V40c0.2-2.6-1.8-4.8-4.4-5c-0.2,0-0.4,0-0.6,0H15c-2.6-0.2-4.8,1.8-5,4.4c0,0.2,0,0.4,0,0.6v130c-0.2,2.6,1.8,4.8,4.4,5c0.2,0,0.4,0,0.6,0h45l20-20H30V55H210 M60,205l60-70l60,70H60L60,205z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-arrow-left" viewBox="0 0 240 240" focusable="false"><path d="M55.4,104.4c-1.1,1.1-2.2,2.3-3.1,3.6c-6.9,9.9-4.4,23.5,5.5,30.4L159.7,240l33.9-33.9l-84.9-84.9l84.9-84.9L157.3,0L55.4,104.4L55.4,104.4z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-playback-rate" viewBox="0 0 240 240" focusable="false"><path d="M158.83,48.83A71.17,71.17,0,1,0,230,120,71.163,71.163,0,0,0,158.83,48.83Zm45.293,77.632H152.34V74.708h12.952v38.83h38.83ZM35.878,74.708h38.83V87.66H35.878ZM10,113.538H61.755V126.49H10Zm25.878,38.83h38.83V165.32H35.878Z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-settings" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M204,145l-25-14c0.8-3.6,1.2-7.3,1-11c0.2-3.7-0.2-7.4-1-11l25-14c2.2-1.6,3.1-4.5,2-7l-16-26c-1.2-2.1-3.8-2.9-6-2l-25,14c-6-4.2-12.3-7.9-19-11V35c0.2-2.6-1.8-4.8-4.4-5c-0.2,0-0.4,0-0.6,0h-30c-2.6-0.2-4.8,1.8-5,4.4c0,0.2,0,0.4,0,0.6v28c-6.7,3.1-13,6.7-19,11L56,60c-2.2-0.9-4.8-0.1-6,2L35,88c-1.6,2.2-1.3,5.3,0.9,6.9c0,0,0.1,0,0.1,0.1l25,14c-0.8,3.6-1.2,7.3-1,11c-0.2,3.7,0.2,7.4,1,11l-25,14c-2.2,1.6-3.1,4.5-2,7l16,26c1.2,2.1,3.8,2.9,6,2l25-14c5.7,4.6,12.2,8.3,19,11v28c-0.2,2.6,1.8,4.8,4.4,5c0.2,0,0.4,0,0.6,0h30c2.6,0.2,4.8-1.8,5-4.4c0-0.2,0-0.4,0-0.6v-28c7-2.3,13.5-6,19-11l25,14c2.5,1.3,5.6,0.4,7-2l15-26C206.7,149.4,206,146.7,204,145z M120,149.9c-16.5,0-30-13.4-30-30s13.4-30,30-30s30,13.4,30,30c0.3,16.3-12.6,29.7-28.9,30C120.7,149.9,120.4,149.9,120,149.9z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-audio-tracks" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M35,34h160v20H35V34z M35,94h160V74H35V94z M35,134h60v-20H35V134z M160,114c-23.4-1.3-43.6,16.5-45,40v50h20c5.2,0.3,9.7-3.6,10-8.9c0-0.4,0-0.7,0-1.1v-20c0.3-5.2-3.6-9.7-8.9-10c-0.4,0-0.7,0-1.1,0h-10v-10c1.5-17.9,17.1-31.3,35-30c17.9-1.3,33.6,12.1,35,30v10H185c-5.2-0.3-9.7,3.6-10,8.9c0,0.4,0,0.7,0,1.1v20c-0.3,5.2,3.6,9.7,8.9,10c0.4,0,0.7,0,1.1,0h20v-50C203.5,130.6,183.4,112.7,160,114z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg class="jw-svg-icon jw-svg-icon-quality-100" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M55,200H35c-3,0-5-2-5-4c0,0,0,0,0-1v-30c0-3,2-5,4-5c0,0,0,0,1,0h20c3,0,5,2,5,4c0,0,0,0,0,1v30C60,198,58,200,55,200L55,200z M110,195v-70c0-3-2-5-4-5c0,0,0,0-1,0H85c-3,0-5,2-5,4c0,0,0,0,0,1v70c0,3,2,5,4,5c0,0,0,0,1,0h20C108,200,110,198,110,195L110,195z M160,195V85c0-3-2-5-4-5c0,0,0,0-1,0h-20c-3,0-5,2-5,4c0,0,0,0,0,1v110c0,3,2,5,4,5c0,0,0,0,1,0h20C158,200,160,198,160,195L160,195z M210,195V45c0-3-2-5-4-5c0,0,0,0-1,0h-20c-3,0-5,2-5,4c0,0,0,0,0,1v150c0,3,2,5,4,5c0,0,0,0,1,0h20C208,200,210,198,210,195L210,195z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-fullscreen-off" viewBox="0 0 240 240" focusable="false"><path d="M109.2,134.9l-8.4,50.1c-0.4,2.7-2.4,3.3-4.4,1.4L82,172l-27.9,27.9l-14.2-14.2l27.9-27.9l-14.4-14.4c-1.9-1.9-1.3-3.9,1.4-4.4l50.1-8.4c1.8-0.5,3.6,0.6,4.1,2.4C109.4,133.7,109.4,134.3,109.2,134.9L109.2,134.9z M172.1,82.1L200,54.2L185.8,40l-27.9,27.9l-14.4-14.4c-1.9-1.9-3.9-1.3-4.4,1.4l-8.4,50.1c-0.5,1.8,0.6,3.6,2.4,4.1c0.5,0.2,1.2,0.2,1.7,0l50.1-8.4c2.7-0.4,3.3-2.4,1.4-4.4L172.1,82.1z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-fullscreen-on" viewBox="0 0 240 240" focusable="false"><path d="M96.3,186.1c1.9,1.9,1.3,4-1.4,4.4l-50.6,8.4c-1.8,0.5-3.7-0.6-4.2-2.4c-0.2-0.6-0.2-1.2,0-1.7l8.4-50.6c0.4-2.7,2.4-3.4,4.4-1.4l14.5,14.5l28.2-28.2l14.3,14.3l-28.2,28.2L96.3,186.1z M195.8,39.1l-50.6,8.4c-2.7,0.4-3.4,2.4-1.4,4.4l14.5,14.5l-28.2,28.2l14.3,14.3l28.2-28.2l14.5,14.5c1.9,1.9,4,1.3,4.4-1.4l8.4-50.6c0.5-1.8-0.6-3.6-2.4-4.2C197,39,196.4,39,195.8,39.1L195.8,39.1z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-close" viewBox="0 0 240 240" focusable="false"><path d="M134.8,120l48.6-48.6c2-1.9,2.1-5.2,0.2-7.2c0,0-0.1-0.1-0.2-0.2l-7.4-7.4c-1.9-2-5.2-2.1-7.2-0.2c0,0-0.1,0.1-0.2,0.2L120,105.2L71.4,56.6c-1.9-2-5.2-2.1-7.2-0.2c0,0-0.1,0.1-0.2,0.2L56.6,64c-2,1.9-2.1,5.2-0.2,7.2c0,0,0.1,0.1,0.2,0.2l48.6,48.7l-48.6,48.6c-2,1.9-2.1,5.2-0.2,7.2c0,0,0.1,0.1,0.2,0.2l7.4,7.4c1.9,2,5.2,2.1,7.2,0.2c0,0,0.1-0.1,0.2-0.2l48.7-48.6l48.6,48.6c1.9,2,5.2,2.1,7.2,0.2c0,0,0.1-0.1,0.2-0.2l7.4-7.4c2-1.9,2.1-5.2,0.2-7.2c0,0-0.1-0.1-0.2-0.2L134.8,120z"></path></svg>';
    },
    function (e, t) {
      e.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-jwplayer-logo" viewBox="0 0 992 1024" focusable="false"><path d="M144 518.4c0 6.4-6.4 6.4-6.4 0l-3.2-12.8c0 0-6.4-19.2-12.8-38.4 0-6.4-6.4-12.8-9.6-22.4-6.4-6.4-16-9.6-28.8-6.4-9.6 3.2-16 12.8-16 22.4s0 16 0 25.6c3.2 25.6 22.4 121.6 32 140.8 9.6 22.4 35.2 32 54.4 22.4 22.4-9.6 28.8-35.2 38.4-54.4 9.6-25.6 60.8-166.4 60.8-166.4 6.4-12.8 9.6-12.8 9.6 0 0 0 0 140.8-3.2 204.8 0 25.6 0 67.2 9.6 89.6 6.4 16 12.8 28.8 25.6 38.4s28.8 12.8 44.8 12.8c6.4 0 16-3.2 22.4-6.4 9.6-6.4 16-12.8 25.6-22.4 16-19.2 28.8-44.8 38.4-64 25.6-51.2 89.6-201.6 89.6-201.6 6.4-12.8 9.6-12.8 9.6 0 0 0-9.6 256-9.6 355.2 0 25.6 6.4 48 12.8 70.4 9.6 22.4 22.4 38.4 44.8 48s48 9.6 70.4-3.2c16-9.6 28.8-25.6 38.4-38.4 12.8-22.4 25.6-48 32-70.4 19.2-51.2 35.2-102.4 51.2-153.6s153.6-540.8 163.2-582.4c0-6.4 0-9.6 0-12.8 0-9.6-6.4-19.2-16-22.4-16-6.4-32 0-38.4 12.8-6.4 16-195.2 470.4-195.2 470.4-6.4 12.8-9.6 12.8-9.6 0 0 0 0-156.8 0-288 0-70.4-35.2-108.8-83.2-118.4-22.4-3.2-44.8 0-67.2 12.8s-35.2 32-48 54.4c-16 28.8-105.6 297.6-105.6 297.6-6.4 12.8-9.6 12.8-9.6 0 0 0-3.2-115.2-6.4-144-3.2-41.6-12.8-108.8-67.2-115.2-51.2-3.2-73.6 57.6-86.4 99.2-9.6 25.6-51.2 163.2-51.2 163.2v3.2z"></path></svg>';
    },
    function (e, t, n) {
      var i = n(96);
      "string" == typeof i && (i = [["all-players", i, ""]]),
        n(61).style(i, "all-players"),
        i.locals && (e.exports = i.locals);
    },
    function (e, t, n) {
      (e.exports = n(60)(!1)).push([
        e.i,
        '.jw-overlays,.jw-controls,.jw-controls-backdrop,.jw-flag-small-player .jw-settings-menu,.jw-settings-submenu{height:100%;width:100%}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-settings-menu .jw-icon.jw-button-color::after{position:absolute;right:0}.jw-overlays,.jw-controls,.jw-controls-backdrop,.jw-settings-item-active::before{top:0;position:absolute;left:0}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-settings-menu .jw-icon.jw-button-color::after{position:absolute;bottom:0;left:0}.jw-nextup-close{position:absolute;top:0;right:0}.jw-overlays,.jw-controls,.jw-flag-small-player .jw-settings-menu{position:absolute;bottom:0;right:0}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-time-tip::after,.jw-settings-menu .jw-icon.jw-button-color::after,.jw-text-live::before,.jw-controlbar .jw-tooltip::after,.jw-settings-menu .jw-tooltip::after{content:"";display:block}.jw-svg-icon{height:24px;width:24px;fill:currentColor;pointer-events:none}.jw-icon{height:44px;width:44px;background-color:transparent;outline:none}.jw-icon.jw-tab-focus:focus{border:solid 2px #4d90fe}.jw-icon-airplay .jw-svg-icon-airplay-off{display:none}.jw-off.jw-icon-airplay .jw-svg-icon-airplay-off{display:block}.jw-icon-airplay .jw-svg-icon-airplay-on{display:block}.jw-off.jw-icon-airplay .jw-svg-icon-airplay-on{display:none}.jw-icon-cc .jw-svg-icon-cc-off{display:none}.jw-off.jw-icon-cc .jw-svg-icon-cc-off{display:block}.jw-icon-cc .jw-svg-icon-cc-on{display:block}.jw-off.jw-icon-cc .jw-svg-icon-cc-on{display:none}.jw-icon-fullscreen .jw-svg-icon-fullscreen-off{display:none}.jw-off.jw-icon-fullscreen .jw-svg-icon-fullscreen-off{display:block}.jw-icon-fullscreen .jw-svg-icon-fullscreen-on{display:block}.jw-off.jw-icon-fullscreen .jw-svg-icon-fullscreen-on{display:none}.jw-icon-volume .jw-svg-icon-volume-0{display:none}.jw-off.jw-icon-volume .jw-svg-icon-volume-0{display:block}.jw-icon-volume .jw-svg-icon-volume-100{display:none}.jw-full.jw-icon-volume .jw-svg-icon-volume-100{display:block}.jw-icon-volume .jw-svg-icon-volume-50{display:block}.jw-off.jw-icon-volume .jw-svg-icon-volume-50,.jw-full.jw-icon-volume .jw-svg-icon-volume-50{display:none}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after{height:100%;width:24px;box-shadow:inset 0 -3px 0 -1px currentColor;margin:auto;opacity:0;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-settings-menu .jw-icon[aria-checked="true"]::after,.jw-settings-open .jw-icon-settings::after,.jw-icon-volume.jw-open::after{opacity:1}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-cc,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-settings,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-audio-tracks,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-hd,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-settings-sharing,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-fullscreen,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player).jw-flag-cast-available .jw-icon-airplay,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player).jw-flag-cast-available .jw-icon-cast{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-volume,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-text-live{bottom:6px}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-volume::after{display:none}.jw-overlays,.jw-controls{pointer-events:none}.jw-controls-backdrop{display:block;background:linear-gradient(to bottom, transparent, rgba(0,0,0,0.4) 77%, rgba(0,0,0,0.4) 100%) 100% 100% / 100% 240px no-repeat transparent;transition:opacity 250ms cubic-bezier(0, .25, .25, 1),background-size 250ms cubic-bezier(0, .25, .25, 1);pointer-events:none}.jw-overlays{cursor:auto}.jw-controls{overflow:hidden}.jw-flag-small-player .jw-controls{text-align:center}.jw-text{height:1em;font-family:Arial,Helvetica,sans-serif;font-size:.75em;font-style:normal;font-weight:normal;color:#fff;text-align:center;font-variant:normal;font-stretch:normal}.jw-controlbar,.jw-skip,.jw-display-icon-container .jw-icon,.jw-nextup-container,.jw-autostart-mute,.jw-overlays .jw-plugin{pointer-events:all}.jwplayer .jw-display-icon-container,.jw-error .jw-display-icon-container{width:auto;height:auto;box-sizing:content-box}.jw-display{display:table;height:100%;padding:57px 0;position:relative;width:100%}.jw-flag-dragging .jw-display{display:none}.jw-state-idle:not(.jw-flag-cast-available) .jw-display{padding:0}.jw-display-container{display:table-cell;height:100%;text-align:center;vertical-align:middle}.jw-display-controls{display:inline-block}.jwplayer .jw-display-icon-container{float:left}.jw-display-icon-container{display:inline-block;padding:5.5px;margin:0 22px}.jw-display-icon-container .jw-icon{height:75px;width:75px;cursor:pointer;display:flex;justify-content:center;align-items:center}.jw-display-icon-container .jw-icon .jw-svg-icon{height:33px;width:33px;padding:0;position:relative}.jw-display-icon-container .jw-icon .jw-svg-icon-rewind{padding:.2em .05em}.jw-breakpoint--1 .jw-nextup-container{display:none}.jw-breakpoint-0 .jw-display-icon-next,.jw-breakpoint--1 .jw-display-icon-next,.jw-breakpoint-0 .jw-display-icon-rewind,.jw-breakpoint--1 .jw-display-icon-rewind{display:none}.jw-breakpoint-0 .jw-display .jw-icon,.jw-breakpoint--1 .jw-display .jw-icon,.jw-breakpoint-0 .jw-display .jw-svg-icon,.jw-breakpoint--1 .jw-display .jw-svg-icon{width:44px;height:44px;line-height:44px}.jw-breakpoint-0 .jw-display .jw-icon:before,.jw-breakpoint--1 .jw-display .jw-icon:before,.jw-breakpoint-0 .jw-display .jw-svg-icon:before,.jw-breakpoint--1 .jw-display .jw-svg-icon:before{width:22px;height:22px}.jw-breakpoint-1 .jw-display .jw-icon,.jw-breakpoint-1 .jw-display .jw-svg-icon{width:44px;height:44px;line-height:44px}.jw-breakpoint-1 .jw-display .jw-icon:before,.jw-breakpoint-1 .jw-display .jw-svg-icon:before{width:22px;height:22px}.jw-breakpoint-1 .jw-display .jw-icon.jw-icon-rewind:before{width:33px;height:33px}.jw-breakpoint-2 .jw-display .jw-icon,.jw-breakpoint-3 .jw-display .jw-icon,.jw-breakpoint-2 .jw-display .jw-svg-icon,.jw-breakpoint-3 .jw-display .jw-svg-icon{width:77px;height:77px;line-height:77px}.jw-breakpoint-2 .jw-display .jw-icon:before,.jw-breakpoint-3 .jw-display .jw-icon:before,.jw-breakpoint-2 .jw-display .jw-svg-icon:before,.jw-breakpoint-3 .jw-display .jw-svg-icon:before{width:38.5px;height:38.5px}.jw-breakpoint-4 .jw-display .jw-icon,.jw-breakpoint-5 .jw-display .jw-icon,.jw-breakpoint-6 .jw-display .jw-icon,.jw-breakpoint-7 .jw-display .jw-icon,.jw-breakpoint-4 .jw-display .jw-svg-icon,.jw-breakpoint-5 .jw-display .jw-svg-icon,.jw-breakpoint-6 .jw-display .jw-svg-icon,.jw-breakpoint-7 .jw-display .jw-svg-icon{width:88px;height:88px;line-height:88px}.jw-breakpoint-4 .jw-display .jw-icon:before,.jw-breakpoint-5 .jw-display .jw-icon:before,.jw-breakpoint-6 .jw-display .jw-icon:before,.jw-breakpoint-7 .jw-display .jw-icon:before,.jw-breakpoint-4 .jw-display .jw-svg-icon:before,.jw-breakpoint-5 .jw-display .jw-svg-icon:before,.jw-breakpoint-6 .jw-display .jw-svg-icon:before,.jw-breakpoint-7 .jw-display .jw-svg-icon:before{width:44px;height:44px}.jw-controlbar{display:flex;flex-flow:row wrap;align-items:center;justify-content:center;position:absolute;left:0;bottom:0;width:100%;border:none;border-radius:0;background-size:auto;box-shadow:none;max-height:72px;transition:250ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility;transition-delay:0s}.jw-breakpoint-7 .jw-controlbar{max-height:140px}.jw-breakpoint-7 .jw-controlbar .jw-button-container{padding:0 48px 20px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-tooltip{margin-bottom:-7px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume .jw-overlay{padding-bottom:40%}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-text{font-size:1em}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-text.jw-text-elapsed{justify-content:flex-end}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-inline,.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume{height:60px;width:60px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-inline .jw-svg-icon,.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume .jw-svg-icon{height:30px;width:30px}.jw-breakpoint-7 .jw-controlbar .jw-slider-time{padding:0 60px;height:34px}.jw-breakpoint-7 .jw-controlbar .jw-slider-time .jw-slider-container{height:10px}.jw-controlbar .jw-button-image{background:no-repeat 50% 50%;background-size:contain;max-height:24px}.jw-controlbar .jw-spacer{flex:1 1 auto;align-self:stretch}.jw-controlbar .jw-icon.jw-button-color:hover{color:#fff}.jw-button-container{display:flex;flex-flow:row nowrap;flex:1 1 auto;align-items:center;justify-content:center;width:100%;padding:0 12px}.jw-slider-horizontal{background-color:transparent}.jw-icon-inline{position:relative}.jw-icon-inline,.jw-icon-tooltip{height:44px;width:44px;align-items:center;display:flex;justify-content:center}.jw-icon-inline:not(.jw-text),.jw-icon-tooltip,.jw-slider-horizontal{cursor:pointer}.jw-text-elapsed,.jw-text-duration{justify-content:flex-start;width:-webkit-fit-content;width:-moz-fit-content;width:fit-content}.jw-icon-tooltip{position:relative}.jw-knob:hover,.jw-icon-inline:hover,.jw-icon-tooltip:hover,.jw-icon-display:hover,.jw-option:before:hover{color:#fff}.jw-time-tip,.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{pointer-events:none}.jw-icon-cast{display:none;margin:0;padding:0}.jw-icon-cast google-cast-launcher{background-color:transparent;border:none;padding:0;width:24px;height:24px;cursor:pointer}.jw-icon-inline.jw-icon-volume{display:none}.jwplayer .jw-text-countdown{display:none}.jw-flag-small-player .jw-display{padding-top:0;padding-bottom:0}.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-rewind,.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-next,.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-playback{display:none}.jw-flag-ads-vpaid:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar,.jw-flag-user-inactive.jw-state-playing:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar,.jw-flag-user-inactive.jw-state-buffering:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar{visibility:hidden;pointer-events:none;opacity:0;transition-delay:0s, 250ms}.jw-flag-ads-vpaid:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop,.jw-flag-user-inactive.jw-state-playing:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop,.jw-flag-user-inactive.jw-state-buffering:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop{opacity:0}.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-countdown{display:flex}.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint--1 .jw-text-elapsed,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-elapsed,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint--1 .jw-text-duration,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-duration{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-text-countdown,.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-related-btn,.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-slider-volume{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-controlbar{flex-direction:column-reverse}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-button-container{height:30px}.jw-breakpoint--1.jw-flag-ads:not(.jw-flag-audio-player) .jw-icon-volume,.jw-breakpoint--1.jw-flag-ads:not(.jw-flag-audio-player) .jw-icon-fullscreen{display:none}.jwplayer:not(.jw-breakpoint-0) .jw-text-duration:before,.jwplayer:not(.jw-breakpoint--1) .jw-text-duration:before{content:"/";padding-right:1ch;padding-left:1ch}.jwplayer:not(.jw-flag-user-inactive) .jw-controlbar{will-change:transform}.jwplayer:not(.jw-flag-user-inactive) .jw-controlbar .jw-text{-webkit-transform-style:preserve-3d;transform-style:preserve-3d}.jw-slider-container{display:flex;align-items:center;position:relative;touch-action:none}.jw-rail,.jw-buffer,.jw-progress{position:absolute;cursor:pointer}.jw-progress{background-color:#f2f2f2}.jw-rail{background-color:rgba(255,255,255,0.3)}.jw-buffer{background-color:rgba(255,255,255,0.3)}.jw-knob{height:13px;width:13px;background-color:#fff;border-radius:50%;box-shadow:0 0 10px rgba(0,0,0,0.4);opacity:1;pointer-events:none;position:absolute;-webkit-transform:translate(-50%, -50%) scale(0);transform:translate(-50%, -50%) scale(0);transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, -webkit-transform;transition-property:opacity, transform;transition-property:opacity, transform, -webkit-transform}.jw-flag-dragging .jw-slider-time .jw-knob,.jw-icon-volume:active .jw-slider-volume .jw-knob{box-shadow:0 0 26px rgba(0,0,0,0.2),0 0 10px rgba(0,0,0,0.4),0 0 0 6px rgba(255,255,255,0.2)}.jw-slider-horizontal,.jw-slider-vertical{display:flex}.jw-slider-horizontal .jw-slider-container{height:5px;width:100%}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress,.jw-slider-horizontal .jw-cue,.jw-slider-horizontal .jw-knob{top:50%}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress,.jw-slider-horizontal .jw-cue{-webkit-transform:translate(0, -50%);transform:translate(0, -50%)}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress{height:5px}.jw-slider-horizontal .jw-rail{width:100%}.jw-slider-vertical{align-items:center;flex-direction:column}.jw-slider-vertical .jw-slider-container{height:88px;width:5px}.jw-slider-vertical .jw-rail,.jw-slider-vertical .jw-buffer,.jw-slider-vertical .jw-progress,.jw-slider-vertical .jw-knob{left:50%}.jw-slider-vertical .jw-rail,.jw-slider-vertical .jw-buffer,.jw-slider-vertical .jw-progress{height:100%;width:5px;-webkit-backface-visibility:hidden;backface-visibility:hidden;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);transition:-webkit-transform 150ms ease-in-out;transition:transform 150ms ease-in-out;transition:transform 150ms ease-in-out, -webkit-transform 150ms ease-in-out;bottom:0}.jw-slider-vertical .jw-knob{-webkit-transform:translate(-50%, 50%);transform:translate(-50%, 50%)}.jw-slider-time.jw-tab-focus:focus .jw-rail{outline:solid 2px #4d90fe}.jw-slider-time,.jw-flag-audio-player .jw-slider-volume{height:17px;width:100%;align-items:center;background:transparent none;padding:0 12px}.jw-slider-time .jw-cue{background-color:rgba(33,33,33,0.8);cursor:pointer;position:absolute;width:6px}.jw-slider-time,.jw-horizontal-volume-container{z-index:1;outline:none}.jw-slider-time .jw-rail,.jw-horizontal-volume-container .jw-rail,.jw-slider-time .jw-buffer,.jw-horizontal-volume-container .jw-buffer,.jw-slider-time .jw-progress,.jw-horizontal-volume-container .jw-progress,.jw-slider-time .jw-cue,.jw-horizontal-volume-container .jw-cue{-webkit-backface-visibility:hidden;backface-visibility:hidden;height:100%;-webkit-transform:translate(0, -50%) scale(1, .6);transform:translate(0, -50%) scale(1, .6);transition:-webkit-transform 150ms ease-in-out;transition:transform 150ms ease-in-out;transition:transform 150ms ease-in-out, -webkit-transform 150ms ease-in-out}.jw-slider-time:hover .jw-rail,.jw-horizontal-volume-container:hover .jw-rail,.jw-slider-time:focus .jw-rail,.jw-horizontal-volume-container:focus .jw-rail,.jw-flag-dragging .jw-slider-time .jw-rail,.jw-flag-dragging .jw-horizontal-volume-container .jw-rail,.jw-flag-touch .jw-slider-time .jw-rail,.jw-flag-touch .jw-horizontal-volume-container .jw-rail,.jw-slider-time:hover .jw-buffer,.jw-horizontal-volume-container:hover .jw-buffer,.jw-slider-time:focus .jw-buffer,.jw-horizontal-volume-container:focus .jw-buffer,.jw-flag-dragging .jw-slider-time .jw-buffer,.jw-flag-dragging .jw-horizontal-volume-container .jw-buffer,.jw-flag-touch .jw-slider-time .jw-buffer,.jw-flag-touch .jw-horizontal-volume-container .jw-buffer,.jw-slider-time:hover .jw-progress,.jw-horizontal-volume-container:hover .jw-progress,.jw-slider-time:focus .jw-progress,.jw-horizontal-volume-container:focus .jw-progress,.jw-flag-dragging .jw-slider-time .jw-progress,.jw-flag-dragging .jw-horizontal-volume-container .jw-progress,.jw-flag-touch .jw-slider-time .jw-progress,.jw-flag-touch .jw-horizontal-volume-container .jw-progress,.jw-slider-time:hover .jw-cue,.jw-horizontal-volume-container:hover .jw-cue,.jw-slider-time:focus .jw-cue,.jw-horizontal-volume-container:focus .jw-cue,.jw-flag-dragging .jw-slider-time .jw-cue,.jw-flag-dragging .jw-horizontal-volume-container .jw-cue,.jw-flag-touch .jw-slider-time .jw-cue,.jw-flag-touch .jw-horizontal-volume-container .jw-cue{-webkit-transform:translate(0, -50%) scale(1, 1);transform:translate(0, -50%) scale(1, 1)}.jw-slider-time:hover .jw-knob,.jw-horizontal-volume-container:hover .jw-knob,.jw-slider-time:focus .jw-knob,.jw-horizontal-volume-container:focus .jw-knob{-webkit-transform:translate(-50%, -50%) scale(1);transform:translate(-50%, -50%) scale(1)}.jw-slider-time .jw-rail,.jw-horizontal-volume-container .jw-rail{background-color:rgba(255,255,255,0.2)}.jw-slider-time .jw-buffer,.jw-horizontal-volume-container .jw-buffer{background-color:rgba(255,255,255,0.4)}.jw-flag-touch .jw-slider-time::before,.jw-flag-touch .jw-horizontal-volume-container::before{height:44px;width:100%;content:"";position:absolute;display:block;bottom:calc(100% - 17px);left:0}.jw-slider-time.jw-tab-focus:focus .jw-rail,.jw-horizontal-volume-container.jw-tab-focus:focus .jw-rail{outline:solid 2px #4d90fe}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time{height:17px;padding:0}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time .jw-slider-container{height:10px}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time .jw-knob{border-radius:0;border:1px solid rgba(0,0,0,0.75);height:12px;width:10px}.jw-modal{width:284px}.jw-breakpoint-7 .jw-modal,.jw-breakpoint-6 .jw-modal,.jw-breakpoint-5 .jw-modal{height:232px}.jw-breakpoint-4 .jw-modal,.jw-breakpoint-3 .jw-modal{height:192px}.jw-breakpoint-2 .jw-modal,.jw-flag-small-player .jw-modal{bottom:0;right:0;height:100%;width:100%;max-height:none;max-width:none;z-index:2}.jwplayer .jw-rightclick{display:none;position:absolute;white-space:nowrap}.jwplayer .jw-rightclick.jw-open{display:block}.jwplayer .jw-rightclick .jw-rightclick-list{border-radius:1px;list-style:none;margin:0;padding:0}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item{background-color:rgba(0,0,0,0.8);border-bottom:1px solid #444;margin:0}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-logo{color:#fff;display:inline-flex;padding:0 10px 0 0;vertical-align:middle}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-logo .jw-svg-icon{height:20px;width:20px}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-link{border:none;color:#fff;display:block;font-size:11px;line-height:1em;padding:15px 23px;text-align:start;text-decoration:none;width:100%}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item:last-child{border-bottom:none}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item:hover{cursor:pointer}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured{vertical-align:middle}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured .jw-rightclick-link{color:#fff}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured .jw-rightclick-link span{color:#fff}.jwplayer .jw-rightclick .jw-info-overlay-item,.jwplayer .jw-rightclick .jw-share-item,.jwplayer .jw-rightclick .jw-shortcuts-item{border:none;background-color:transparent;outline:none;cursor:pointer}.jw-icon-tooltip.jw-open .jw-overlay{opacity:1;pointer-events:auto;transition-delay:0s}.jw-icon-tooltip.jw-open .jw-overlay:focus{outline:none}.jw-icon-tooltip.jw-open .jw-overlay:focus.jw-tab-focus{outline:solid 2px #4d90fe}.jw-slider-time .jw-overlay:before{height:1em;top:auto}.jw-slider-time .jw-icon-tooltip.jw-open .jw-overlay{pointer-events:none}.jw-volume-tip{padding:13px 0 26px}.jw-time-tip,.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{height:auto;width:100%;box-shadow:0 0 10px rgba(0,0,0,0.4);color:#fff;display:block;margin:0 0 14px;pointer-events:none;position:relative;z-index:0}.jw-time-tip::after,.jw-controlbar .jw-tooltip::after,.jw-settings-menu .jw-tooltip::after{top:100%;position:absolute;left:50%;height:14px;width:14px;border-radius:1px;background-color:currentColor;-webkit-transform-origin:75% 50%;transform-origin:75% 50%;-webkit-transform:translate(-50%, -50%) rotate(45deg);transform:translate(-50%, -50%) rotate(45deg);z-index:-1}.jw-time-tip .jw-text,.jw-controlbar .jw-tooltip .jw-text,.jw-settings-menu .jw-tooltip .jw-text{background-color:#fff;border-radius:1px;color:#000;font-size:10px;height:auto;line-height:1;padding:7px 10px;display:inline-block;min-width:100%;vertical-align:middle}.jw-controlbar .jw-overlay{position:absolute;bottom:100%;left:50%;margin:0;min-height:44px;min-width:44px;opacity:0;pointer-events:none;transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility;transition-delay:0s, 150ms;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);width:100%;z-index:1}.jw-controlbar .jw-overlay .jw-contents{position:relative}.jw-controlbar .jw-option{position:relative;white-space:nowrap;cursor:pointer;list-style:none;height:1.5em;font-family:inherit;line-height:1.5em;padding:0 .5em;font-size:.8em;margin:0}.jw-controlbar .jw-option::before{padding-right:.125em}.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{position:absolute;bottom:100%;left:50%;opacity:0;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);transition:100ms 0s cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility, -webkit-transform;transition-property:opacity, transform, visibility;transition-property:opacity, transform, visibility, -webkit-transform;visibility:hidden;white-space:nowrap;width:auto;z-index:1}.jw-controlbar .jw-tooltip.jw-open,.jw-settings-menu .jw-tooltip.jw-open{opacity:1;-webkit-transform:translate(-50%, -10px);transform:translate(-50%, -10px);transition-duration:150ms;transition-delay:500ms,0s,500ms;visibility:visible}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen{left:auto;right:0;-webkit-transform:translate(0, 0);transform:translate(0, 0)}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen.jw-open,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen.jw-open{-webkit-transform:translate(0, -10px);transform:translate(0, -10px)}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen::after,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen::after{left:auto;right:9px}.jw-tooltip-time{height:auto;width:0;bottom:100%;line-height:normal;padding:0;pointer-events:none;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none}.jw-tooltip-time .jw-overlay{bottom:0;min-height:0;width:auto}.jw-tooltip{bottom:57px;display:none;position:absolute}.jw-tooltip .jw-text{height:100%;white-space:nowrap;text-overflow:ellipsis;direction:unset;max-width:246px;overflow:hidden}.jw-flag-audio-player .jw-tooltip{display:none}.jw-flag-small-player .jw-time-thumb{display:none}.jwplayer .jw-shortcuts-tooltip{top:50%;position:absolute;left:50%;background:#333;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);display:none;color:#fff;pointer-events:all;-webkit-user-select:text;-moz-user-select:text;-ms-user-select:text;user-select:text;overflow:hidden;flex-direction:column;z-index:1}.jwplayer .jw-shortcuts-tooltip.jw-open{display:flex}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-close{flex:0 0 auto;margin:5px 5px 5px auto}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container{display:flex;flex:1 1 auto;flex-flow:column;font-size:12px;margin:0 20px 20px;overflow-y:auto;padding:5px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container::-webkit-scrollbar{background-color:transparent;width:6px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-title{font-weight:bold}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-header{align-items:center;display:flex;justify-content:space-between;margin-bottom:10px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list{display:flex;max-width:340px;margin:0 10px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-tooltip-descriptions{width:100%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row{display:flex;align-items:center;justify-content:space-between;margin:10px 0;width:100%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row .jw-shortcuts-description{margin-right:10px;max-width:70%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row .jw-shortcuts-key{background:#fefefe;color:#333;overflow:hidden;padding:7px 10px;text-overflow:ellipsis;white-space:nowrap}.jw-skip{color:rgba(255,255,255,0.8);cursor:default;position:absolute;display:flex;right:.75em;bottom:56px;padding:.5em;border:1px solid #333;background-color:#000;align-items:center;height:2em}.jw-skip.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jw-skip.jw-skippable{cursor:pointer;padding:.25em .75em}.jw-skip.jw-skippable:hover{cursor:pointer;color:#fff}.jw-skip.jw-skippable .jw-skip-icon{display:inline;height:24px;width:24px;margin:0}.jw-breakpoint-7 .jw-skip{padding:1.35em 1em;bottom:130px}.jw-breakpoint-7 .jw-skip .jw-text{font-size:1em;font-weight:normal}.jw-breakpoint-7 .jw-skip .jw-icon-inline{height:30px;width:30px}.jw-breakpoint-7 .jw-skip .jw-icon-inline .jw-svg-icon{height:30px;width:30px}.jw-skip .jw-skip-icon{display:none;margin-left:-0.75em;padding:0 .5em;pointer-events:none}.jw-skip .jw-skip-icon .jw-svg-icon-next{display:block;padding:0}.jw-skip .jw-text,.jw-skip .jw-skip-icon{vertical-align:middle;font-size:.7em}.jw-skip .jw-text{font-weight:bold}.jw-cast{background-size:cover;display:none;height:100%;position:relative;width:100%}.jw-cast-container{background:linear-gradient(180deg, rgba(25,25,25,0.75), rgba(25,25,25,0.25), rgba(25,25,25,0));left:0;padding:20px 20px 80px;position:absolute;top:0;width:100%}.jw-cast-text{color:#fff;font-size:1.6em}.jw-breakpoint--1 .jw-cast-text,.jw-breakpoint-0 .jw-cast-text{font-size:1.15em}.jw-breakpoint-1 .jw-cast-text,.jw-breakpoint-2 .jw-cast-text,.jw-breakpoint-3 .jw-cast-text{font-size:1.3em}.jw-nextup-container{position:absolute;bottom:66px;left:0;background-color:transparent;cursor:pointer;margin:0 auto;padding:12px;pointer-events:none;right:0;text-align:right;visibility:hidden;width:100%}.jw-settings-open .jw-nextup-container,.jw-info-open .jw-nextup-container{display:none}.jw-breakpoint-7 .jw-nextup-container{padding:60px}.jw-flag-small-player .jw-nextup-container{padding:0 12px 0 0}.jw-flag-small-player .jw-nextup-container .jw-nextup-title,.jw-flag-small-player .jw-nextup-container .jw-nextup-duration,.jw-flag-small-player .jw-nextup-container .jw-nextup-close{display:none}.jw-flag-small-player .jw-nextup-container .jw-nextup-tooltip{height:30px}.jw-flag-small-player .jw-nextup-container .jw-nextup-header{font-size:12px}.jw-flag-small-player .jw-nextup-container .jw-nextup-body{justify-content:center;align-items:center;padding:.75em .3em}.jw-flag-small-player .jw-nextup-container .jw-nextup-thumbnail{width:50%}.jw-flag-small-player .jw-nextup-container .jw-nextup{max-width:65px}.jw-flag-small-player .jw-nextup-container .jw-nextup.jw-nextup-thumbnail-visible{max-width:120px}.jw-nextup{background:#333;border-radius:0;box-shadow:0 0 10px rgba(0,0,0,0.5);color:rgba(255,255,255,0.8);display:inline-block;max-width:280px;overflow:hidden;opacity:0;position:relative;width:64%;pointer-events:all;-webkit-transform:translate(0, -5px);transform:translate(0, -5px);transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, -webkit-transform;transition-property:opacity, transform;transition-property:opacity, transform, -webkit-transform;transition-delay:0s}.jw-nextup:hover .jw-nextup-tooltip{color:#fff}.jw-nextup.jw-nextup-thumbnail-visible{max-width:400px}.jw-nextup.jw-nextup-thumbnail-visible .jw-nextup-thumbnail{display:block}.jw-nextup-container-visible{visibility:visible}.jw-nextup-container-visible .jw-nextup{opacity:1;-webkit-transform:translate(0, 0);transform:translate(0, 0);transition-delay:0s, 0s, 150ms}.jw-nextup-tooltip{display:flex;height:80px}.jw-nextup-thumbnail{width:120px;background-position:center;background-size:cover;flex:0 0 auto;display:none}.jw-nextup-body{flex:1 1 auto;overflow:hidden;padding:.75em .875em;display:flex;flex-flow:column wrap;justify-content:space-between}.jw-nextup-header,.jw-nextup-title{font-size:14px;line-height:1.35}.jw-nextup-header{font-weight:bold}.jw-nextup-title{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;width:100%}.jw-nextup-duration{align-self:flex-end;text-align:right;font-size:12px}.jw-nextup-close{height:24px;width:24px;border:none;color:rgba(255,255,255,0.8);cursor:pointer;margin:6px;visibility:hidden}.jw-nextup-close:hover{color:#fff}.jw-nextup-sticky .jw-nextup-close{visibility:visible}.jw-autostart-mute{position:absolute;bottom:0;right:12px;height:44px;width:44px;background-color:rgba(33,33,33,0.4);padding:5px 4px 5px 6px;display:none}.jwplayer.jw-flag-autostart:not(.jw-flag-media-audio) .jw-nextup{display:none}.jw-settings-menu{position:absolute;bottom:57px;right:12px;align-items:flex-start;background-color:#333;display:none;flex-flow:column nowrap;max-width:284px;pointer-events:auto}.jw-settings-open .jw-settings-menu{display:flex}.jw-breakpoint-7 .jw-settings-menu{bottom:130px;right:60px;max-height:none;max-width:none;height:35%;width:25%}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline{height:60px;width:60px}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline .jw-svg-icon{height:30px;width:30px}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline .jw-tooltip .jw-text{font-size:1em}.jw-breakpoint-7 .jw-settings-menu .jw-settings-back{min-width:60px}.jw-breakpoint-6 .jw-settings-menu,.jw-breakpoint-5 .jw-settings-menu{height:232px;width:284px;max-height:232px}.jw-breakpoint-4 .jw-settings-menu,.jw-breakpoint-3 .jw-settings-menu{height:192px;width:284px;max-height:192px}.jw-breakpoint-2 .jw-settings-menu{height:179px;width:284px;max-height:179px}.jw-flag-small-player .jw-settings-menu{max-width:none}.jw-settings-menu .jw-icon.jw-button-color::after{height:100%;width:24px;box-shadow:inset 0 -3px 0 -1px currentColor;margin:auto;opacity:0;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-settings-menu .jw-icon.jw-button-color[aria-checked="true"]::after{opacity:1}.jw-settings-menu .jw-settings-reset{text-decoration:underline}.jw-settings-topbar{align-items:center;background-color:rgba(0,0,0,0.4);display:flex;flex:0 0 auto;padding:3px 5px 0;width:100%}.jw-settings-topbar.jw-nested-menu-open{padding:0}.jw-settings-topbar.jw-nested-menu-open .jw-icon:not(.jw-settings-close):not(.jw-settings-back){display:none}.jw-settings-topbar.jw-nested-menu-open .jw-svg-icon-close{width:20px}.jw-settings-topbar.jw-nested-menu-open .jw-svg-icon-arrow-left{height:12px}.jw-settings-topbar.jw-nested-menu-open .jw-settings-topbar-text{display:block;outline:none}.jw-settings-topbar .jw-settings-back{min-width:44px}.jw-settings-topbar .jw-settings-topbar-buttons{display:inherit;width:100%;height:100%}.jw-settings-topbar .jw-settings-topbar-text{display:none;color:#fff;font-size:13px;width:100%}.jw-settings-topbar .jw-settings-close{margin-left:auto}.jw-settings-submenu{display:none;flex:1 1 auto;overflow-y:auto;padding:8px 20px 0 5px}.jw-settings-submenu::-webkit-scrollbar{background-color:transparent;width:6px}.jw-settings-submenu::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jw-settings-submenu.jw-settings-submenu-active{display:block}.jw-settings-submenu .jw-submenu-topbar{box-shadow:0 2px 9px 0 #1d1d1d;background-color:#2f2d2d;margin:-8px -20px 0 -5px}.jw-settings-submenu .jw-submenu-topbar .jw-settings-content-item{cursor:pointer;text-align:right;padding-right:15px;text-decoration:underline}.jw-settings-submenu .jw-settings-value-wrapper{float:right;display:flex;align-items:center}.jw-settings-submenu .jw-settings-value-wrapper .jw-settings-content-item-arrow{display:flex}.jw-settings-submenu .jw-settings-value-wrapper .jw-svg-icon-arrow-right{width:8px;margin-left:5px;height:12px}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item{font-size:1em;padding:11px 15px 11px 30px}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item .jw-settings-item-active::before{justify-content:flex-end}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item .jw-auto-label{font-size:.85em;padding-left:10px}.jw-flag-touch .jw-settings-submenu{overflow-y:scroll;-webkit-overflow-scrolling:touch}.jw-auto-label{font-size:10px;font-weight:initial;opacity:.75;padding-left:5px}.jw-settings-content-item{position:relative;color:rgba(255,255,255,0.8);cursor:pointer;font-size:12px;line-height:1;padding:7px 0 7px 15px;width:100%;text-align:left;outline:none}.jw-settings-content-item:hover{color:#fff}.jw-settings-content-item:focus{font-weight:bold}.jw-flag-small-player .jw-settings-content-item{line-height:1.75}.jw-settings-content-item.jw-tab-focus:focus{border:solid 2px #4d90fe}.jw-settings-item-active{font-weight:bold;position:relative}.jw-settings-item-active::before{height:100%;width:1em;align-items:center;content:"\\2022";display:inline-flex;justify-content:center}.jw-breakpoint-2 .jw-settings-open .jw-display-container,.jw-flag-small-player .jw-settings-open .jw-display-container,.jw-flag-touch .jw-settings-open .jw-display-container{display:none}.jw-breakpoint-2 .jw-settings-open.jw-controls,.jw-flag-small-player .jw-settings-open.jw-controls,.jw-flag-touch .jw-settings-open.jw-controls{z-index:1}.jw-flag-small-player .jw-settings-open .jw-controlbar{display:none}.jw-settings-open .jw-icon-settings::after{opacity:1}.jw-settings-open .jw-tooltip-settings{display:none}.jw-sharing-link{cursor:pointer}.jw-shortcuts-container .jw-switch{position:relative;display:inline-block;transition:ease-out .15s;transition-property:opacity, background;border-radius:18px;width:80px;height:20px;padding:10px;background:rgba(80,80,80,0.8);cursor:pointer;font-size:inherit;vertical-align:middle}.jw-shortcuts-container .jw-switch.jw-tab-focus{outline:solid 2px #4d90fe}.jw-shortcuts-container .jw-switch .jw-switch-knob{position:absolute;top:2px;left:1px;transition:ease-out .15s;box-shadow:0 0 10px rgba(0,0,0,0.4);border-radius:13px;width:15px;height:15px;background:#fefefe}.jw-shortcuts-container .jw-switch:before,.jw-shortcuts-container .jw-switch:after{position:absolute;top:3px;transition:inherit;color:#fefefe}.jw-shortcuts-container .jw-switch:before{content:attr(data-jw-switch-disabled);right:8px}.jw-shortcuts-container .jw-switch:after{content:attr(data-jw-switch-enabled);left:8px;opacity:0}.jw-shortcuts-container .jw-switch[aria-checked="true"]{background:#475470}.jw-shortcuts-container .jw-switch[aria-checked="true"]:before{opacity:0}.jw-shortcuts-container .jw-switch[aria-checked="true"]:after{opacity:1}.jw-shortcuts-container .jw-switch[aria-checked="true"] .jw-switch-knob{left:60px}.jw-idle-icon-text{display:none;line-height:1;position:absolute;text-align:center;text-indent:.35em;top:100%;white-space:nowrap;left:50%;-webkit-transform:translateX(-50%);transform:translateX(-50%)}.jw-idle-label{border-radius:50%;color:#fff;-webkit-filter:drop-shadow(1px 1px 5px rgba(12,26,71,0.25));filter:drop-shadow(1px 1px 5px rgba(12,26,71,0.25));font:normal 16px/1 Arial,Helvetica,sans-serif;position:relative;transition:background-color 150ms cubic-bezier(0, .25, .25, 1);transition-property:background-color,-webkit-filter;transition-property:background-color,filter;transition-property:background-color,filter,-webkit-filter;-webkit-font-smoothing:antialiased}.jw-state-idle .jw-icon-display.jw-idle-label .jw-idle-icon-text{display:block}.jw-state-idle .jw-icon-display.jw-idle-label .jw-svg-icon-play{-webkit-transform:scale(.7, .7);transform:scale(.7, .7)}.jw-breakpoint-0.jw-state-idle .jw-icon-display.jw-idle-label,.jw-breakpoint--1.jw-state-idle .jw-icon-display.jw-idle-label{font-size:12px}.jw-info-overlay{top:50%;position:absolute;left:50%;background:#333;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);display:none;color:#fff;pointer-events:all;-webkit-user-select:text;-moz-user-select:text;-ms-user-select:text;user-select:text;overflow:hidden;flex-direction:column}.jw-info-overlay .jw-info-close{flex:0 0 auto;margin:5px 5px 5px auto}.jw-info-open .jw-info-overlay{display:flex}.jw-info-container{display:flex;flex:1 1 auto;flex-flow:column;margin:0 20px 20px;overflow-y:auto;padding:5px}.jw-info-container [class*="jw-info"]:not(:first-of-type){color:rgba(255,255,255,0.8);padding-top:10px;font-size:12px}.jw-info-container .jw-info-description{margin-bottom:30px;text-align:start}.jw-info-container .jw-info-description:empty{display:none}.jw-info-container .jw-info-duration{text-align:start}.jw-info-container .jw-info-title{text-align:start;font-size:12px;font-weight:bold}.jw-info-container::-webkit-scrollbar{background-color:transparent;width:6px}.jw-info-container::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jw-info-clientid{align-self:flex-end;font-size:12px;color:rgba(255,255,255,0.8);margin:0 20px 20px 44px;text-align:right}.jw-flag-touch .jw-info-open .jw-display-container{display:none}@supports ((-webkit-filter: drop-shadow(0 0 3px #000)) or (filter: drop-shadow(0 0 3px #000))){.jwplayer.jw-ab-drop-shadow .jw-controls .jw-svg-icon,.jwplayer.jw-ab-drop-shadow .jw-controls .jw-icon.jw-text,.jwplayer.jw-ab-drop-shadow .jw-slider-container .jw-rail,.jwplayer.jw-ab-drop-shadow .jw-title{text-shadow:none;box-shadow:none;-webkit-filter:drop-shadow(0 2px 3px rgba(0,0,0,0.3));filter:drop-shadow(0 2px 3px rgba(0,0,0,0.3))}.jwplayer.jw-ab-drop-shadow .jw-button-color{opacity:.8;transition-property:color, opacity}.jwplayer.jw-ab-drop-shadow .jw-button-color:not(:hover){color:#fff;opacity:.8}.jwplayer.jw-ab-drop-shadow .jw-button-color:hover{opacity:1}.jwplayer.jw-ab-drop-shadow .jw-controls-backdrop{background-image:linear-gradient(to bottom, hsla(0, 0%, 0%, 0), hsla(0, 0%, 0%, 0.00787) 10.79%, hsla(0, 0%, 0%, 0.02963) 21.99%, hsla(0, 0%, 0%, 0.0625) 33.34%, hsla(0, 0%, 0%, 0.1037) 44.59%, hsla(0, 0%, 0%, 0.15046) 55.48%, hsla(0, 0%, 0%, 0.2) 65.75%, hsla(0, 0%, 0%, 0.24954) 75.14%, hsla(0, 0%, 0%, 0.2963) 83.41%, hsla(0, 0%, 0%, 0.3375) 90.28%, hsla(0, 0%, 0%, 0.37037) 95.51%, hsla(0, 0%, 0%, 0.39213) 98.83%, hsla(0, 0%, 0%, 0.4));mix-blend-mode:multiply;transition-property:opacity}.jw-state-idle.jwplayer.jw-ab-drop-shadow .jw-controls-backdrop{background-image:linear-gradient(to bottom, hsla(0, 0%, 0%, 0.2), hsla(0, 0%, 0%, 0.19606) 1.17%, hsla(0, 0%, 0%, 0.18519) 4.49%, hsla(0, 0%, 0%, 0.16875) 9.72%, hsla(0, 0%, 0%, 0.14815) 16.59%, hsla(0, 0%, 0%, 0.12477) 24.86%, hsla(0, 0%, 0%, 0.1) 34.25%, hsla(0, 0%, 0%, 0.07523) 44.52%, hsla(0, 0%, 0%, 0.05185) 55.41%, hsla(0, 0%, 0%, 0.03125) 66.66%, hsla(0, 0%, 0%, 0.01481) 78.01%, hsla(0, 0%, 0%, 0.00394) 89.21%, hsla(0, 0%, 0%, 0));background-size:100% 7rem;background-position:50% 0}.jwplayer.jw-ab-drop-shadow.jw-state-idle .jw-controls{background-color:transparent}}.jw-video-thumbnail-container{position:relative;overflow:hidden}.jw-video-thumbnail-container:not(.jw-related-shelf-item-image){height:100%;width:100%}.jw-video-thumbnail-container.jw-video-thumbnail-generated{position:absolute;top:0;left:0}.jw-video-thumbnail-container:hover,.jw-related-item-content:hover .jw-video-thumbnail-container,.jw-related-shelf-item:hover .jw-video-thumbnail-container{cursor:pointer}.jw-video-thumbnail-container:hover .jw-video-thumbnail:not(.jw-video-thumbnail-completed),.jw-related-item-content:hover .jw-video-thumbnail-container .jw-video-thumbnail:not(.jw-video-thumbnail-completed),.jw-related-shelf-item:hover .jw-video-thumbnail-container .jw-video-thumbnail:not(.jw-video-thumbnail-completed){opacity:1}.jw-video-thumbnail-container .jw-video-thumbnail{position:absolute;top:50%;left:50%;bottom:unset;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);width:100%;height:auto;min-width:100%;min-height:100%;opacity:0;transition:opacity .3s ease;object-fit:cover;background:#000}.jw-related-item-next-up .jw-video-thumbnail-container .jw-video-thumbnail{height:100%;width:auto}.jw-video-thumbnail-container .jw-video-thumbnail.jw-video-thumbnail-visible:not(.jw-video-thumbnail-completed){opacity:1}.jw-video-thumbnail-container .jw-video-thumbnail.jw-video-thumbnail-completed{opacity:0}.jw-video-thumbnail-container .jw-video-thumbnail~.jw-svg-icon-play{display:none}.jw-video-thumbnail-container .jw-video-thumbnail+.jw-related-shelf-item-aspect{pointer-events:none}.jw-video-thumbnail-container .jw-video-thumbnail+.jw-related-item-poster-content{pointer-events:none}.jw-state-idle:not(.jw-flag-cast-available) .jw-display{padding:0}.jw-state-idle .jw-controls{background:rgba(0,0,0,0.4)}.jw-state-idle.jw-flag-cast-available:not(.jw-flag-audio-player) .jw-controlbar .jw-slider-time,.jw-state-idle.jw-flag-cardboard-available .jw-controlbar .jw-slider-time,.jw-state-idle.jw-flag-cast-available:not(.jw-flag-audio-player) .jw-controlbar .jw-icon:not(.jw-icon-cardboard):not(.jw-icon-cast):not(.jw-icon-airplay),.jw-state-idle.jw-flag-cardboard-available .jw-controlbar .jw-icon:not(.jw-icon-cardboard):not(.jw-icon-cast):not(.jw-icon-airplay){display:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-icon:focus{border:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-icon .jw-svg-icon-buffer{-webkit-animation:jw-spin 2s linear infinite;animation:jw-spin 2s linear infinite;display:block}@-webkit-keyframes jw-spin{100%{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}@keyframes jw-spin{100%{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}.jwplayer.jw-state-buffering .jw-icon-playback .jw-svg-icon-play{display:none}.jwplayer.jw-state-buffering .jw-icon-display .jw-svg-icon-pause{display:none}.jwplayer.jw-state-playing .jw-display .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-playing .jw-icon-playback .jw-svg-icon-play{display:none}.jwplayer.jw-state-playing .jw-display .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-playing .jw-icon-playback .jw-svg-icon-pause{display:block}.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio) .jw-controls-backdrop{opacity:0}.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio) .jw-logo-bottom-left,.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio):not(.jw-flag-autostart) .jw-logo-bottom-right{bottom:0}.jwplayer .jw-icon-playback .jw-svg-icon-stop{display:none}.jwplayer.jw-state-paused .jw-svg-icon-pause,.jwplayer.jw-state-idle .jw-svg-icon-pause,.jwplayer.jw-state-error .jw-svg-icon-pause,.jwplayer.jw-state-complete .jw-svg-icon-pause{display:none}.jwplayer.jw-state-error .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-complete .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-buffering .jw-icon-display .jw-svg-icon-play{display:none}.jwplayer:not(.jw-state-buffering) .jw-svg-icon-buffer{display:none}.jwplayer:not(.jw-state-complete) .jw-svg-icon-replay{display:none}.jwplayer:not(.jw-state-error) .jw-svg-icon-error{display:none}.jwplayer.jw-state-complete .jw-display .jw-icon-display .jw-svg-icon-replay{display:block}.jwplayer.jw-state-complete .jw-display .jw-text{display:none}.jwplayer.jw-state-complete .jw-controls{background:rgba(0,0,0,0.4);height:100%}.jw-state-idle .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-paused .jw-icon-playback .jw-svg-icon-pause,.jwplayer.jw-state-paused .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-complete .jw-icon-playback .jw-svg-icon-pause{display:none}.jw-state-idle .jw-display-icon-rewind,.jwplayer.jw-state-buffering .jw-display-icon-rewind,.jwplayer.jw-state-complete .jw-display-icon-rewind,body .jw-error .jw-display-icon-rewind,body .jwplayer.jw-state-error .jw-display-icon-rewind,.jw-state-idle .jw-display-icon-next,.jwplayer.jw-state-buffering .jw-display-icon-next,.jwplayer.jw-state-complete .jw-display-icon-next,body .jw-error .jw-display-icon-next,body .jwplayer.jw-state-error .jw-display-icon-next{display:none}body .jw-error .jw-icon-display,body .jwplayer.jw-state-error .jw-icon-display{cursor:default}body .jw-error .jw-icon-display .jw-svg-icon-error,body .jwplayer.jw-state-error .jw-icon-display .jw-svg-icon-error{display:block}body .jw-error .jw-icon-container{position:absolute;width:100%;height:100%;top:0;left:0;bottom:0;right:0}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-preview{display:none}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title{padding-top:4px}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title-primary{width:auto;display:inline-block;padding-right:.5ch}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title-secondary{width:auto;display:inline-block;padding-left:0}body .jwplayer.jw-state-error .jw-controlbar,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-controlbar{display:none}body .jwplayer.jw-state-error .jw-settings-menu,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-settings-menu{height:100%;top:50%;left:50%;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%)}body .jwplayer.jw-state-error .jw-display,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-display{padding:0}body .jwplayer.jw-state-error .jw-logo-bottom-left,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-logo-bottom-left,body .jwplayer.jw-state-error .jw-logo-bottom-right,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-logo-bottom-right{bottom:0}.jwplayer.jw-state-playing.jw-flag-user-inactive .jw-display{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-state-playing:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display,.jwplayer.jw-state-paused:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting):not(.jw-flag-play-rejected) .jw-display{display:none}.jwplayer.jw-state-paused.jw-flag-play-rejected:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display-icon-rewind,.jwplayer.jw-state-paused.jw-flag-play-rejected:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display-icon-next{display:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-text,.jwplayer.jw-state-complete .jw-display .jw-text{display:none}.jwplayer.jw-flag-casting:not(.jw-flag-audio-player) .jw-cast{display:block}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-display-icon-container{display:none}.jwplayer.jw-flag-casting .jw-icon-hd,.jwplayer.jw-flag-casting .jw-captions,.jwplayer.jw-flag-casting .jw-icon-fullscreen,.jwplayer.jw-flag-casting .jw-icon-audio-tracks{display:none}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-icon-volume{display:none}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-icon-airplay{color:#fff}.jw-state-playing.jw-flag-casting:not(.jw-flag-audio-player) .jw-display,.jw-state-paused.jw-flag-casting:not(.jw-flag-audio-player) .jw-display{display:table}.jwplayer.jw-flag-cast-available .jw-icon-cast,.jwplayer.jw-flag-cast-available .jw-icon-airplay{display:flex}.jwplayer.jw-flag-cardboard-available .jw-icon-cardboard{display:flex}.jwplayer.jw-flag-live .jw-display-icon-rewind{visibility:hidden}.jwplayer.jw-flag-live .jw-controlbar .jw-text-elapsed,.jwplayer.jw-flag-live .jw-controlbar .jw-text-duration,.jwplayer.jw-flag-live .jw-controlbar .jw-text-countdown,.jwplayer.jw-flag-live .jw-controlbar .jw-slider-time{display:none}.jwplayer.jw-flag-live .jw-controlbar .jw-text-alt{display:flex}.jwplayer.jw-flag-live .jw-controlbar .jw-overlay:after{display:none}.jwplayer.jw-flag-live .jw-nextup-container{bottom:44px}.jwplayer.jw-flag-live .jw-text-elapsed,.jwplayer.jw-flag-live .jw-text-duration{display:none}.jwplayer.jw-flag-live .jw-text-live{cursor:default}.jwplayer.jw-flag-live .jw-text-live:hover{color:rgba(255,255,255,0.8)}.jwplayer.jw-flag-live.jw-state-playing .jw-icon-playback .jw-svg-icon-stop,.jwplayer.jw-flag-live.jw-state-buffering .jw-icon-playback .jw-svg-icon-stop{display:block}.jwplayer.jw-flag-live.jw-state-playing .jw-icon-playback .jw-svg-icon-pause,.jwplayer.jw-flag-live.jw-state-buffering .jw-icon-playback .jw-svg-icon-pause{display:none}.jw-text-live{height:24px;width:auto;align-items:center;border-radius:1px;color:rgba(255,255,255,0.8);display:flex;font-size:12px;font-weight:bold;margin-right:10px;padding:0 1ch;text-rendering:geometricPrecision;text-transform:uppercase;transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:box-shadow,color}.jw-text-live::before{height:8px;width:8px;background-color:currentColor;border-radius:50%;margin-right:6px;opacity:1;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-text-live.jw-dvr-live{box-shadow:inset 0 0 0 2px currentColor}.jw-text-live.jw-dvr-live::before{opacity:.5}.jw-text-live.jw-dvr-live:hover{color:#fff}.jwplayer.jw-flag-controls-hidden .jw-logo.jw-hide{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-flag-controls-hidden:not(.jw-flag-casting) .jw-logo-top-right{top:0}.jwplayer.jw-flag-controls-hidden .jw-plugin{bottom:.5em}.jwplayer.jw-flag-controls-hidden .jw-nextup-container{bottom:0}.jw-flag-controls-hidden .jw-controlbar,.jw-flag-controls-hidden .jw-display{visibility:hidden;pointer-events:none;opacity:0;transition-delay:0s, 250ms}.jw-flag-controls-hidden .jw-controls-backdrop{opacity:0}.jw-flag-controls-hidden .jw-logo{visibility:visible}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-logo.jw-hide{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-casting) .jw-logo-top-right{top:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-plugin{bottom:.5em}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-nextup-container{bottom:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-controls-hidden) .jw-media{cursor:none;-webkit-cursor-visibility:auto-hide}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing.jw-flag-casting .jw-display{display:table}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-ads) .jw-autostart-mute{display:flex}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-flag-casting .jw-nextup-container{bottom:66px}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-flag-casting.jw-state-idle .jw-nextup-container{display:none}.jw-flag-media-audio .jw-preview{display:block}.jwplayer.jw-flag-ads .jw-preview,.jwplayer.jw-flag-ads .jw-logo,.jwplayer.jw-flag-ads .jw-captions.jw-captions-enabled,.jwplayer.jw-flag-ads .jw-nextup-container,.jwplayer.jw-flag-ads .jw-text-duration,.jwplayer.jw-flag-ads .jw-text-elapsed{display:none}.jwplayer.jw-flag-ads video::-webkit-media-text-track-container{display:none}.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-rewind,.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-next,.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-display{display:none}.jwplayer.jw-flag-ads.jw-flag-small-player.jw-state-buffering .jw-display-icon-display{display:inline-block}.jwplayer.jw-flag-ads .jw-controlbar{flex-wrap:wrap-reverse}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time{height:auto;padding:0;pointer-events:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-slider-container{height:5px}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-rail,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-knob,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-buffer,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-cue,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-icon-settings{display:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-progress{-webkit-transform:none;transform:none;top:auto}.jwplayer.jw-flag-ads .jw-controlbar .jw-tooltip,.jwplayer.jw-flag-ads .jw-controlbar .jw-icon-tooltip:not(.jw-icon-volume),.jwplayer.jw-flag-ads .jw-controlbar .jw-icon-inline:not(.jw-icon-playback):not(.jw-icon-fullscreen):not(.jw-icon-volume){display:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-volume-tip{padding:13px 0}.jwplayer.jw-flag-ads .jw-controlbar .jw-text-alt{display:flex}.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid) .jw-controls .jw-controlbar,.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-autostart .jw-controls .jw-controlbar{display:flex;pointer-events:all;visibility:visible;opacity:1}.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-user-inactive .jw-controls-backdrop,.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-autostart.jw-flag-user-inactive .jw-controls-backdrop{opacity:1;background-size:100% 60px}.jwplayer.jw-flag-ads-vpaid .jw-display-container,.jwplayer.jw-flag-touch.jw-flag-ads-vpaid .jw-display-container,.jwplayer.jw-flag-ads-vpaid .jw-skip,.jwplayer.jw-flag-touch.jw-flag-ads-vpaid .jw-skip{display:none}.jwplayer.jw-flag-ads-vpaid.jw-flag-small-player .jw-controls{background:none}.jwplayer.jw-flag-ads-vpaid.jw-flag-small-player .jw-controls::after{content:none}.jwplayer.jw-flag-ads-hide-controls .jw-controls-backdrop,.jwplayer.jw-flag-ads-hide-controls .jw-controls{display:none !important}.jw-flag-overlay-open-related .jw-controls,.jw-flag-overlay-open-related .jw-title,.jw-flag-overlay-open-related .jw-logo{display:none}.jwplayer.jw-flag-rightclick-open{overflow:visible}.jwplayer.jw-flag-rightclick-open .jw-rightclick{z-index:16777215}body .jwplayer.jw-flag-flash-blocked .jw-controls,body .jwplayer.jw-flag-flash-blocked .jw-overlays,body .jwplayer.jw-flag-flash-blocked .jw-controls-backdrop,body .jwplayer.jw-flag-flash-blocked .jw-preview{display:none}body .jwplayer.jw-flag-flash-blocked .jw-error-msg{top:25%}.jw-flag-touch.jw-breakpoint-7 .jw-captions,.jw-flag-touch.jw-breakpoint-6 .jw-captions,.jw-flag-touch.jw-breakpoint-5 .jw-captions,.jw-flag-touch.jw-breakpoint-4 .jw-captions,.jw-flag-touch.jw-breakpoint-7 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-6 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-5 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-4 .jw-nextup-container{bottom:4.25em}.jw-flag-touch .jw-controlbar .jw-icon-volume{display:flex}.jw-flag-touch .jw-display,.jw-flag-touch .jw-display-container,.jw-flag-touch .jw-display-controls{pointer-events:none}.jw-flag-touch.jw-state-paused:not(.jw-breakpoint-1) .jw-display-icon-next,.jw-flag-touch.jw-state-playing:not(.jw-breakpoint-1) .jw-display-icon-next,.jw-flag-touch.jw-state-paused:not(.jw-breakpoint-1) .jw-display-icon-rewind,.jw-flag-touch.jw-state-playing:not(.jw-breakpoint-1) .jw-display-icon-rewind{display:none}.jw-flag-touch.jw-state-paused.jw-flag-dragging .jw-display{display:none}.jw-flag-audio-player{background-color:#000}.jw-flag-audio-player:not(.jw-flag-flash-blocked) .jw-media{visibility:hidden}.jw-flag-audio-player .jw-title{background:none}.jw-flag-audio-player object{min-height:44px}.jw-flag-audio-player:not(.jw-flag-live) .jw-spacer{display:none}.jw-flag-audio-player .jw-preview,.jw-flag-audio-player .jw-display,.jw-flag-audio-player .jw-title,.jw-flag-audio-player .jw-nextup-container{display:none}.jw-flag-audio-player .jw-controlbar{position:relative}.jw-flag-audio-player .jw-controlbar .jw-button-container{padding-right:3px;padding-left:0}.jw-flag-audio-player .jw-controlbar .jw-icon-tooltip,.jw-flag-audio-player .jw-controlbar .jw-icon-inline{display:none}.jw-flag-audio-player .jw-controlbar .jw-icon-volume,.jw-flag-audio-player .jw-controlbar .jw-icon-playback,.jw-flag-audio-player .jw-controlbar .jw-icon-next,.jw-flag-audio-player .jw-controlbar .jw-icon-rewind,.jw-flag-audio-player .jw-controlbar .jw-icon-cast,.jw-flag-audio-player .jw-controlbar .jw-text-live,.jw-flag-audio-player .jw-controlbar .jw-icon-airplay,.jw-flag-audio-player .jw-controlbar .jw-logo-button,.jw-flag-audio-player .jw-controlbar .jw-text-elapsed,.jw-flag-audio-player .jw-controlbar .jw-text-duration{display:flex;flex:0 0 auto}.jw-flag-audio-player .jw-controlbar .jw-text-duration,.jw-flag-audio-player .jw-controlbar .jw-text-countdown{padding-right:10px}.jw-flag-audio-player .jw-controlbar .jw-slider-time{flex:0 1 auto;align-items:center;display:flex;order:1}.jw-flag-audio-player .jw-controlbar .jw-icon-volume{margin-right:0;transition:margin-right 150ms cubic-bezier(0, .25, .25, 1)}.jw-flag-audio-player .jw-controlbar .jw-icon-volume .jw-overlay{display:none}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container{transition:width 300ms cubic-bezier(0, .25, .25, 1);width:0}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open{width:140px}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open .jw-slider-volume{padding-right:24px;transition:opacity 300ms;opacity:1}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open~.jw-slider-time{flex:1 1 auto;width:auto;transition:opacity 300ms, width 300ms}.jw-flag-audio-player .jw-controlbar .jw-slider-volume{opacity:0}.jw-flag-audio-player .jw-controlbar .jw-slider-volume .jw-knob{-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%)}.jw-flag-audio-player .jw-controlbar .jw-slider-volume~.jw-icon-volume{margin-right:140px}.jw-flag-audio-player.jw-breakpoint-1 .jw-horizontal-volume-container.jw-open~.jw-slider-time,.jw-flag-audio-player.jw-breakpoint-2 .jw-horizontal-volume-container.jw-open~.jw-slider-time{opacity:0}.jw-flag-audio-player.jw-flag-small-player .jw-text-elapsed,.jw-flag-audio-player.jw-flag-small-player .jw-text-duration{display:none}.jw-flag-audio-player.jw-flag-ads .jw-slider-time{display:none}.jw-hidden{display:none}',
        "",
      ]);
    },
  ],
]);
