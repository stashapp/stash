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
  [1],
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
    function (t, e, n) {
      "use strict";
      n.r(e);
      var i,
        o = n(8),
        a = n(3),
        r = n(7),
        l = n(43),
        s = n(5),
        c = n(15),
        u = n(40);
      function w(t) {
        return (
          i || (i = new DOMParser()),
          Object(s.r)(
            Object(s.s)(i.parseFromString(t, "image/svg+xml").documentElement)
          )
        );
      }
      var p = function (t, e, n, i) {
          var o = document.createElement("div");
          (o.className =
            "jw-icon jw-icon-inline jw-button-color jw-reset " + t),
            o.setAttribute("role", "button"),
            o.setAttribute("tabindex", "0"),
            n && o.setAttribute("aria-label", n),
            (o.style.display = "none");
          var a = new u.a(o).on("click tap enter", e || function () {});
          return (
            i &&
              Array.prototype.forEach.call(i, function (t) {
                "string" == typeof t ? o.appendChild(w(t)) : o.appendChild(t);
              }),
            {
              ui: a,
              element: function () {
                return o;
              },
              toggle: function (t) {
                t ? this.show() : this.hide();
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
        d = n(0),
        j = n(71),
        h = n.n(j),
        f = n(72),
        g = n.n(f),
        b = n(73),
        y = n.n(b),
        v = n(74),
        m = n.n(v),
        x = n(75),
        k = n.n(x),
        O = n(76),
        C = n.n(O),
        S = n(77),
        T = n.n(S),
        M = n(78),
        z = n.n(M),
        E = n(79),
        L = n.n(E),
        B = n(80),
        _ = n.n(B),
        V = n(81),
        A = n.n(V),
        N = n(82),
        H = n.n(N),
        P = n(83),
        I = n.n(P),
        R = n(84),
        q = n.n(R),
        D = n(85),
        U = n.n(D),
        F = n(86),
        W = n.n(F),
        Z = n(62),
        K = n.n(Z),
        X = n(87),
        Y = n.n(X),
        G = n(88),
        J = n.n(G),
        Q = n(89),
        $ = n.n(Q),
        tt = n(90),
        et = n.n(tt),
        nt = n(91),
        it = n.n(nt),
        ot = n(92),
        at = n.n(ot),
        rt = n(93),
        lt = n.n(rt),
        st = n(94),
        ct = n.n(st),
        ut = null;
      function wt(t) {
        var e = ht().querySelector(dt(t));
        if (e) return jt(e);
        throw new Error("Icon not found " + t);
      }
      function pt(t) {
        var e = ht().querySelectorAll(t.split(",").map(dt).join(","));
        if (!e.length) throw new Error("Icons not found " + t);
        return Array.prototype.map.call(e, function (t) {
          return jt(t);
        });
      }
      function dt(t) {
        return ".jw-svg-icon-".concat(t);
      }
      function jt(t) {
        return t.cloneNode(!0);
      }
      function ht() {
        return (
          ut ||
            (ut = w(
              "<xml>" +
                h.a +
                g.a +
                y.a +
                m.a +
                k.a +
                C.a +
                T.a +
                z.a +
                L.a +
                _.a +
                A.a +
                H.a +
                I.a +
                q.a +
                U.a +
                W.a +
                K.a +
                Y.a +
                J.a +
                $.a +
                et.a +
                it.a +
                at.a +
                lt.a +
                ct.a +
                "</xml>"
            )),
          ut
        );
      }
      var ft = n(10);
      function gt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var bt = {};
      var yt = (function () {
          function t(e, n, i, o, a) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t);
            var r,
              l = document.createElement("div");
            (l.className = "jw-icon jw-icon-inline jw-button-color jw-reset ".concat(
              a || ""
            )),
              l.setAttribute("button", o),
              l.setAttribute("role", "button"),
              l.setAttribute("tabindex", "0"),
              n && l.setAttribute("aria-label", n),
              e && "<svg" === e.substring(0, 4)
                ? (r = (function (t) {
                    if (!bt[t]) {
                      var e = Object.keys(bt);
                      e.length > 10 && delete bt[e[0]];
                      var n = w(t);
                      bt[t] = n;
                    }
                    return bt[t].cloneNode(!0);
                  })(e))
                : (((r = document.createElement("div")).className =
                    "jw-icon jw-button-image jw-button-color jw-reset"),
                  e &&
                    Object(ft.d)(r, {
                      backgroundImage: "url(".concat(e, ")"),
                    })),
              l.appendChild(r),
              new u.a(l).on("click tap enter", i, this),
              l.addEventListener("mousedown", function (t) {
                t.preventDefault();
              }),
              (this.id = o),
              (this.buttonElement = l);
          }
          var e, n, i;
          return (
            (e = t),
            (n = [
              {
                key: "element",
                value: function () {
                  return this.buttonElement;
                },
              },
              {
                key: "toggle",
                value: function (t) {
                  t ? this.show() : this.hide();
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
            ]) && gt(e.prototype, n),
            i && gt(e, i),
            t
          );
        })(),
        vt = n(11);
      function mt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var xt = function (t) {
          var e = Object(s.c)(t),
            n = window.pageXOffset;
          return (
            n &&
              o.OS.android &&
              document.body.parentElement.getBoundingClientRect().left >= 0 &&
              ((e.left -= n), (e.right -= n)),
            e
          );
        },
        kt = (function () {
          function t(e, n) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(d.g)(this, r.a),
              (this.className = e + " jw-background-color jw-reset"),
              (this.orientation = n);
          }
          var e, n, i;
          return (
            (e = t),
            (n = [
              {
                key: "setup",
                value: function () {
                  (this.el = Object(s.e)(
                    (function () {
                      var t =
                          arguments.length > 0 && void 0 !== arguments[0]
                            ? arguments[0]
                            : "",
                        e =
                          arguments.length > 1 && void 0 !== arguments[1]
                            ? arguments[1]
                            : "";
                      return (
                        '<div class="'
                          .concat(t, " ")
                          .concat(e, ' jw-reset" aria-hidden="true">') +
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
                    (this.railBounds = xt(this.elementRail));
                },
              },
              {
                key: "dragEnd",
                value: function (t) {
                  this.dragMove(t), this.trigger("dragEnd");
                },
              },
              {
                key: "dragMove",
                value: function (t) {
                  var e,
                    n,
                    i = (this.railBounds = this.railBounds
                      ? this.railBounds
                      : xt(this.elementRail));
                  return (
                    (n =
                      "horizontal" === this.orientation
                        ? (e = t.pageX) < i.left
                          ? 0
                          : e > i.right
                          ? 100
                          : 100 * Object(l.a)((e - i.left) / i.width, 0, 1)
                        : (e = t.pageY) >= i.bottom
                        ? 0
                        : e <= i.top
                        ? 100
                        : 100 *
                          Object(l.a)(
                            (i.height - (e - i.top)) / i.height,
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
                value: function (t) {
                  (this.railBounds = xt(this.elementRail)), this.dragMove(t);
                },
              },
              {
                key: "limit",
                value: function (t) {
                  return t;
                },
              },
              {
                key: "update",
                value: function (t) {
                  this.trigger("update", { percentage: t });
                },
              },
              {
                key: "render",
                value: function (t) {
                  (t = Math.max(0, Math.min(t, 100))),
                    "horizontal" === this.orientation
                      ? ((this.elementThumb.style.left = t + "%"),
                        (this.elementProgress.style.width = t + "%"))
                      : ((this.elementThumb.style.bottom = t + "%"),
                        (this.elementProgress.style.height = t + "%"));
                },
              },
              {
                key: "updateBuffer",
                value: function (t) {
                  this.elementBuffer.style.width = t + "%";
                },
              },
              {
                key: "element",
                value: function () {
                  return this.el;
                },
              },
            ]) && mt(e.prototype, n),
            i && mt(e, i),
            t
          );
        })(),
        Ot = function (t, e) {
          t &&
            e &&
            (t.setAttribute("aria-label", e),
            t.setAttribute("role", "button"),
            t.setAttribute("tabindex", "0"));
        };
      function Ct(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var St = (function () {
          function t(e, n, i, o) {
            var a = this;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(d.g)(this, r.a),
              (this.el = document.createElement("div"));
            var l =
              "jw-icon jw-icon-tooltip " + e + " jw-button-color jw-reset";
            i || (l += " jw-hidden"),
              Ot(this.el, n),
              (this.el.className = l),
              (this.tooltip = document.createElement("div")),
              (this.tooltip.className = "jw-overlay jw-reset"),
              (this.openClass = "jw-open"),
              (this.componentType = "tooltip"),
              this.el.appendChild(this.tooltip),
              o &&
                o.length > 0 &&
                Array.prototype.forEach.call(o, function (t) {
                  "string" == typeof t
                    ? a.el.appendChild(w(t))
                    : a.el.appendChild(t);
                });
          }
          var e, n, i;
          return (
            (e = t),
            (n = [
              {
                key: "addContent",
                value: function (t) {
                  this.content && this.removeContent(),
                    (this.content = t),
                    this.tooltip.appendChild(t);
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
                value: function (t) {
                  this.isOpen ||
                    (this.trigger("open-" + this.componentType, t, {
                      isOpen: !0,
                    }),
                    (this.isOpen = !0),
                    Object(s.v)(this.el, this.openClass, this.isOpen));
                },
              },
              {
                key: "closeTooltip",
                value: function (t) {
                  this.isOpen &&
                    (this.trigger("close-" + this.componentType, t, {
                      isOpen: !1,
                    }),
                    (this.isOpen = !1),
                    Object(s.v)(this.el, this.openClass, this.isOpen));
                },
              },
              {
                key: "toggleOpenState",
                value: function (t) {
                  this.isOpen ? this.closeTooltip(t) : this.openTooltip(t);
                },
              },
            ]) && Ct(e.prototype, n),
            i && Ct(e, i),
            t
          );
        })(),
        Tt = n(22),
        Mt = n(57);
      function zt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var Et = (function () {
          function t(e, n) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              (this.time = e),
              (this.text = n),
              (this.el = document.createElement("div")),
              (this.el.className = "jw-cue jw-reset");
          }
          var e, n, i;
          return (
            (e = t),
            (n = [
              {
                key: "align",
                value: function (t) {
                  if ("%" === this.time.toString().slice(-1))
                    this.pct = this.time;
                  else {
                    var e = (this.time / t) * 100;
                    this.pct = e + "%";
                  }
                  this.el.style.left = this.pct;
                },
              },
            ]) && zt(e.prototype, n),
            i && zt(e, i),
            t
          );
        })(),
        Lt = {
          loadChapters: function (t) {
            Object(Tt.a)(
              t,
              this.chaptersLoaded.bind(this),
              this.chaptersFailed,
              { plainText: !0 }
            );
          },
          chaptersLoaded: function (t) {
            var e = Object(Mt.a)(t.responseText);
            if (Array.isArray(e)) {
              var n = this._model.get("cues").concat(e);
              this._model.set("cues", n);
            }
          },
          chaptersFailed: function () {},
          addCue: function (t) {
            this.cues.push(new Et(t.begin, t.text));
          },
          drawCues: function () {
            var t = this,
              e = this._model.get("duration");
            !e ||
              e <= 0 ||
              this.cues.forEach(function (n) {
                n.align(e),
                  n.el.addEventListener("mouseover", function () {
                    t.activeCue = n;
                  }),
                  n.el.addEventListener("mouseout", function () {
                    t.activeCue = null;
                  }),
                  t.elementRail.appendChild(n.el);
              });
          },
          resetCues: function () {
            this.cues.forEach(function (t) {
              t.el.parentNode && t.el.parentNode.removeChild(t.el);
            }),
              (this.cues = []);
          },
        };
      function Bt(t) {
        (this.begin = t.begin), (this.end = t.end), (this.img = t.text);
      }
      var _t = {
        loadThumbnails: function (t) {
          t &&
            ((this.vttPath = t.split("?")[0].split("/").slice(0, -1).join("/")),
            (this.individualImage = null),
            Object(Tt.a)(
              t,
              this.thumbnailsLoaded.bind(this),
              this.thumbnailsFailed.bind(this),
              { plainText: !0 }
            ));
        },
        thumbnailsLoaded: function (t) {
          var e = Object(Mt.a)(t.responseText);
          Array.isArray(e) &&
            (e.forEach(function (t) {
              this.thumbnails.push(new Bt(t));
            }, this),
            this.drawCues());
        },
        thumbnailsFailed: function () {},
        chooseThumbnail: function (t) {
          var e = Object(d.A)(this.thumbnails, { end: t }, Object(d.z)("end"));
          e >= this.thumbnails.length && (e = this.thumbnails.length - 1);
          var n = this.thumbnails[e].img;
          return (
            n.indexOf("://") < 0 &&
              (n = this.vttPath ? this.vttPath + "/" + n : n),
            n
          );
        },
        loadThumbnail: function (t) {
          var e = this.chooseThumbnail(t),
            n = { margin: "0 auto", backgroundPosition: "0 0" };
          if (e.indexOf("#xywh") > 0)
            try {
              var i = /(.+)#xywh=(\d+),(\d+),(\d+),(\d+)/.exec(e);
              (e = i[1]),
                (n.backgroundPosition = -1 * i[2] + "px " + -1 * i[3] + "px"),
                (n.width = i[4]),
                this.timeTip.setWidth(+n.width),
                (n.height = i[5]);
            } catch (t) {
              return;
            }
          else
            this.individualImage ||
              ((this.individualImage = new Image()),
              (this.individualImage.onload = Object(d.a)(function () {
                (this.individualImage.onload = null),
                  this.timeTip.image({
                    width: this.individualImage.width,
                    height: this.individualImage.height,
                  }),
                  this.timeTip.setWidth(this.individualImage.width);
              }, this)),
              (this.individualImage.src = e));
          return (n.backgroundImage = 'url("' + e + '")'), n;
        },
        showThumbnail: function (t) {
          this._model.get("containerWidth") <= 420 ||
            this.thumbnails.length < 1 ||
            this.timeTip.image(this.loadThumbnail(t));
        },
        resetThumbnails: function () {
          this.timeTip.image({ backgroundImage: "", width: 0, height: 0 }),
            (this.thumbnails = []);
        },
      };
      function Vt(t, e, n) {
        return (Vt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, n) {
                var i = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = Rt(t));

                  );
                  return t;
                })(t, e);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, e);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(t, e, n || t);
      }
      function At(t) {
        return (At =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (t) {
                return typeof t;
              }
            : function (t) {
                return t &&
                  "function" == typeof Symbol &&
                  t.constructor === Symbol &&
                  t !== Symbol.prototype
                  ? "symbol"
                  : typeof t;
              })(t);
      }
      function Nt(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Ht(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function Pt(t, e, n) {
        return e && Ht(t.prototype, e), n && Ht(t, n), t;
      }
      function It(t, e) {
        return !e || ("object" !== At(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Rt(t) {
        return (Rt = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function qt(t, e) {
        if ("function" != typeof e && null !== e)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (t.prototype = Object.create(e && e.prototype, {
          constructor: { value: t, writable: !0, configurable: !0 },
        })),
          e && Dt(t, e);
      }
      function Dt(t, e) {
        return (Dt =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var Ut = (function (t) {
        function e() {
          return Nt(this, e), It(this, Rt(e).apply(this, arguments));
        }
        return (
          qt(e, t),
          Pt(e, [
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
                var t = document.createElement("div");
                (t.className = "jw-time-tip jw-reset"),
                  t.appendChild(this.img),
                  t.appendChild(this.text),
                  this.addContent(t);
              },
            },
            {
              key: "image",
              value: function (t) {
                Object(ft.d)(this.img, t);
              },
            },
            {
              key: "update",
              value: function (t) {
                this.text.textContent = t;
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
              value: function (t) {
                t
                  ? (this.containerWidth = t + 16)
                  : this.tooltip &&
                    (this.containerWidth =
                      Object(s.c)(this.container).width + 16);
              },
            },
            {
              key: "resetWidth",
              value: function () {
                this.containerWidth = 0;
              },
            },
          ]),
          e
        );
      })(St);
      var Ft = (function (t) {
        function e(t, n, i) {
          var o;
          return (
            Nt(this, e),
            ((o = It(
              this,
              Rt(e).call(this, "jw-slider-time", "horizontal")
            ))._model = t),
            (o._api = n),
            (o.timeUpdateKeeper = i),
            (o.timeTip = new Ut("jw-tooltip-time", null, !0)),
            o.timeTip.setup(),
            (o.cues = []),
            (o.seekThrottled = Object(d.B)(o.performSeek, 400)),
            (o.mobileHoverDistance = 5),
            o.setup(),
            o
          );
        }
        return (
          qt(e, t),
          Pt(e, [
            {
              key: "setup",
              value: function () {
                var t = this;
                Vt(Rt(e.prototype), "setup", this).apply(this, arguments),
                  this._model
                    .on("change:duration", this.onDuration, this)
                    .on("change:cues", this.updateCues, this)
                    .on("seeked", function () {
                      t._model.get("scrubbing") || t.updateAriaText();
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
                Object(s.t)(n, "tabindex", "0"),
                  Object(s.t)(n, "role", "slider"),
                  Object(s.t)(
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
              value: function (t) {
                (this.seekTo = t),
                  this.seekThrottled(),
                  Vt(Rt(e.prototype), "update", this).apply(this, arguments);
              },
            },
            {
              key: "dragStart",
              value: function () {
                this._model.set("scrubbing", !0),
                  Vt(Rt(e.prototype), "dragStart", this).apply(this, arguments);
              },
            },
            {
              key: "dragEnd",
              value: function () {
                Vt(Rt(e.prototype), "dragEnd", this).apply(this, arguments),
                  this._model.set("scrubbing", !1);
              },
            },
            {
              key: "onBuffer",
              value: function (t, e) {
                this.updateBuffer(e);
              },
            },
            {
              key: "onPosition",
              value: function (t, e) {
                this.updateTime(e, t.get("duration"));
              },
            },
            {
              key: "onDuration",
              value: function (t, e) {
                this.updateTime(t.get("position"), e),
                  Object(s.t)(this.el, "aria-valuemin", 0),
                  Object(s.t)(this.el, "aria-valuemax", e),
                  this.drawCues();
              },
            },
            {
              key: "onStreamType",
              value: function (t, e) {
                this.streamType = e;
              },
            },
            {
              key: "updateTime",
              value: function (t, e) {
                var n = 0;
                if (e)
                  if ("DVR" === this.streamType) {
                    var i = this._model.get("dvrSeekLimit"),
                      o = e + i;
                    n = ((o - (t + i)) / o) * 100;
                  } else
                    ("VOD" !== this.streamType && this.streamType) ||
                      (n = (t / e) * 100);
                this.render(n);
              },
            },
            {
              key: "onPlaylistItem",
              value: function (t, e) {
                this.reset();
                var n = t.get("cues");
                !this.cues.length && n.length && this.updateCues(null, n);
                var i = e.tracks;
                Object(d.f)(
                  i,
                  function (t) {
                    t && t.kind && "thumbnails" === t.kind.toLowerCase()
                      ? this.loadThumbnails(t.file)
                      : t &&
                        t.kind &&
                        "chapters" === t.kind.toLowerCase() &&
                        this.loadChapters(t.file);
                  },
                  this
                );
              },
            },
            {
              key: "performSeek",
              value: function () {
                var t,
                  e = this.seekTo,
                  n = this._model.get("duration");
                if (0 === n) this._api.play({ reason: "interaction" });
                else if ("DVR" === this.streamType) {
                  var i = this._model.get("seekRange") || { start: 0 },
                    o = this._model.get("dvrSeekLimit");
                  (t = i.start + ((-n - o) * e) / 100),
                    this._api.seek(t, { reason: "interaction" });
                } else
                  (t = (e / 100) * n),
                    this._api.seek(Math.min(t, n - 0.25), {
                      reason: "interaction",
                    });
              },
            },
            {
              key: "showTimeTooltip",
              value: function (t) {
                var e = this,
                  n = this._model.get("duration");
                if (0 !== n) {
                  var i,
                    o = this._model.get("containerWidth"),
                    a = Object(s.c)(this.elementRail),
                    r = t.pageX ? t.pageX - a.left : t.x,
                    c = (r = Object(l.a)(r, 0, a.width)) / a.width,
                    u = n * c;
                  if (n < 0)
                    u = (n += this._model.get("dvrSeekLimit")) - (u = n * c);
                  if (
                    ("touch" === t.pointerType &&
                      (this.activeCue = this.cues.reduce(function (t, n) {
                        return Math.abs(r - (parseInt(n.pct) / 100) * a.width) <
                          e.mobileHoverDistance
                          ? n
                          : t;
                      }, void 0)),
                    this.activeCue)
                  )
                    i = this.activeCue.text;
                  else {
                    (i = Object(vt.timeFormat)(u, !0)),
                      n < 0 && u > -1 && (i = "Live");
                  }
                  var w = this.timeTip;
                  w.update(i),
                    this.textLength !== i.length &&
                      ((this.textLength = i.length), w.resetWidth()),
                    this.showThumbnail(u),
                    Object(s.a)(w.el, "jw-open");
                  var p = w.getWidth(),
                    d = a.width / 100,
                    j = o - a.width,
                    h = 0;
                  p > j && (h = (p - j) / (200 * d));
                  var f = 100 * Math.min(1 - h, Math.max(h, c)).toFixed(3);
                  Object(ft.d)(w.el, { left: f + "%" });
                }
              },
            },
            {
              key: "hideTimeTooltip",
              value: function () {
                Object(s.o)(this.timeTip.el, "jw-open");
              },
            },
            {
              key: "updateCues",
              value: function (t, e) {
                var n = this;
                this.resetCues(),
                  e &&
                    e.length &&
                    (e.forEach(function (t) {
                      n.addCue(t);
                    }),
                    this.drawCues());
              },
            },
            {
              key: "updateAriaText",
              value: function () {
                var t = this._model;
                if (!t.get("seeking")) {
                  var e = t.get("position"),
                    n = t.get("duration"),
                    i = Object(vt.timeFormat)(e);
                  "DVR" !== this.streamType &&
                    (i += " of ".concat(Object(vt.timeFormat)(n)));
                  var o = this.el;
                  document.activeElement !== o &&
                    (this.timeUpdateKeeper.textContent = i),
                    Object(s.t)(o, "aria-valuenow", e),
                    Object(s.t)(o, "aria-valuetext", i);
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
          e
        );
      })(kt);
      Object(d.g)(Ft.prototype, Lt, _t);
      var Wt = Ft;
      function Zt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function Kt(t, e, n) {
        return (Kt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, n) {
                var i = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = Qt(t));

                  );
                  return t;
                })(t, e);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, e);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(t, e, n || t);
      }
      function Xt(t) {
        return (Xt =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (t) {
                return typeof t;
              }
            : function (t) {
                return t &&
                  "function" == typeof Symbol &&
                  t.constructor === Symbol &&
                  t !== Symbol.prototype
                  ? "symbol"
                  : typeof t;
              })(t);
      }
      function Yt(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Gt(t, e) {
        return !e || ("object" !== Xt(e) && "function" != typeof e) ? Jt(t) : e;
      }
      function Jt(t) {
        if (void 0 === t)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return t;
      }
      function Qt(t) {
        return (Qt = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function $t(t, e) {
        if ("function" != typeof e && null !== e)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (t.prototype = Object.create(e && e.prototype, {
          constructor: { value: t, writable: !0, configurable: !0 },
        })),
          e && te(t, e);
      }
      function te(t, e) {
        return (te =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var ee = (function (t) {
          function e(t, n, i) {
            var o;
            Yt(this, e);
            var a = "jw-slider-volume";
            return (
              "vertical" === t && (a += " jw-volume-tip"),
              (o = Gt(this, Qt(e).call(this, a, t))).setup(),
              o.element().classList.remove("jw-background-color"),
              Object(s.t)(i, "tabindex", "0"),
              Object(s.t)(i, "aria-label", n),
              Object(s.t)(i, "aria-orientation", t),
              Object(s.t)(i, "aria-valuemin", 0),
              Object(s.t)(i, "aria-valuemax", 100),
              Object(s.t)(i, "role", "slider"),
              (o.uiOver = new u.a(i).on("click", function () {})),
              o
            );
          }
          return $t(e, t), e;
        })(kt),
        ne = (function (t) {
          function e(t, n, i, o, a) {
            var r;
            Yt(this, e),
              ((r = Gt(this, Qt(e).call(this, n, i, !0, o)))._model = t),
              (r.horizontalContainer = a);
            var l = t.get("localization").volumeSlider;
            return (
              (r.horizontalSlider = new ee("horizontal", l, a, Jt(Jt(r)))),
              (r.verticalSlider = new ee("vertical", l, r.tooltip, Jt(Jt(r)))),
              a.appendChild(r.horizontalSlider.element()),
              r.addContent(r.verticalSlider.element()),
              r.verticalSlider.on(
                "update",
                function (t) {
                  this.trigger("update", t);
                },
                Jt(Jt(r))
              ),
              r.horizontalSlider.on(
                "update",
                function (t) {
                  this.trigger("update", t);
                },
                Jt(Jt(r))
              ),
              r.horizontalSlider.uiOver.on("keydown", function (t) {
                var e = t.sourceEvent;
                switch (e.keyCode) {
                  case 37:
                    e.stopPropagation(), r.trigger("adjustVolume", -10);
                    break;
                  case 39:
                    e.stopPropagation(), r.trigger("adjustVolume", 10);
                }
              }),
              (r.ui = new u.a(r.el, { directSelect: !0 })
                .on("click enter", r.toggleValue, Jt(Jt(r)))
                .on("tap", r.toggleOpenState, Jt(Jt(r)))),
              r.addSliderHandlers(r.ui),
              r.addSliderHandlers(r.horizontalSlider.uiOver),
              r.addSliderHandlers(r.verticalSlider.uiOver),
              r.onAudioMode(null, t.get("audioMode")),
              r._model.on("change:audioMode", r.onAudioMode, Jt(Jt(r))),
              r._model.on("change:volume", r.onVolume, Jt(Jt(r))),
              r
            );
          }
          var n, i, o;
          return (
            $t(e, t),
            (n = e),
            (i = [
              {
                key: "onAudioMode",
                value: function (t, e) {
                  var n = e ? 0 : -1;
                  Object(s.t)(this.horizontalContainer, "tabindex", n);
                },
              },
              {
                key: "addSliderHandlers",
                value: function (t) {
                  var e = this.openSlider,
                    n = this.closeSlider;
                  t.on("over", e, this)
                    .on("out", n, this)
                    .on("focus", e, this)
                    .on("blur", n, this);
                },
              },
              {
                key: "openSlider",
                value: function (t) {
                  Kt(Qt(e.prototype), "openTooltip", this).call(this, t),
                    Object(s.v)(this.horizontalContainer, this.openClass, !0);
                },
              },
              {
                key: "closeSlider",
                value: function (t) {
                  Kt(Qt(e.prototype), "closeTooltip", this).call(this, t),
                    Object(s.v)(this.horizontalContainer, this.openClass, !1),
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
            ]) && Zt(n.prototype, i),
            o && Zt(n, o),
            e
          );
        })(St);
      function ie(t, e, n, i, o) {
        var a = document.createElement("div");
        (a.className = "jw-reset-text jw-tooltip jw-tooltip-".concat(e)),
          a.setAttribute("dir", "auto");
        var r = document.createElement("div");
        (r.className = "jw-text"), a.appendChild(r), t.appendChild(a);
        var l = {
            dirty: !!n,
            opened: !1,
            text: n,
            open: function () {
              l.touchEvent ||
                (l.suppress ? (l.suppress = !1) : (c(!0), i && i()));
            },
            close: function () {
              l.touchEvent || (c(!1), o && o());
            },
            setText: function (t) {
              t !== l.text && ((l.text = t), (l.dirty = !0)), l.opened && c(!0);
            },
          },
          c = function (t) {
            t && l.dirty && (Object(s.q)(r, l.text), (l.dirty = !1)),
              (l.opened = t),
              Object(s.v)(a, "jw-open", t);
          };
        return (
          t.addEventListener("mouseover", l.open),
          t.addEventListener("focus", l.open),
          t.addEventListener("blur", l.close),
          t.addEventListener("mouseout", l.close),
          t.addEventListener(
            "touchstart",
            function () {
              l.touchEvent = !0;
            },
            { passive: !0 }
          ),
          l
        );
      }
      var oe = n(47);
      function ae(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function re(t, e) {
        var n = document.createElement("div");
        return (
          (n.className = "jw-icon jw-icon-inline jw-text jw-reset " + t),
          e && Object(s.t)(n, "role", e),
          n
        );
      }
      function le(t) {
        var e = document.createElement("div");
        return (e.className = "jw-reset ".concat(t)), e;
      }
      function se(t, e) {
        if (o.Browser.safari) {
          var n = p(
            "jw-icon-airplay jw-off",
            t,
            e.airplay,
            pt("airplay-off,airplay-on")
          );
          return ie(n.element(), "airplay", e.airplay), n;
        }
        if (o.Browser.chrome && window.chrome) {
          var i = document.createElement("google-cast-launcher");
          Object(s.t)(i, "tabindex", "-1"), (i.className += " jw-reset");
          var a = p("jw-icon-cast", null, e.cast);
          a.ui.off();
          var r = a.element();
          return (
            (r.style.cursor = "pointer"),
            r.appendChild(i),
            (a.button = i),
            ie(r, "chromecast", e.cast),
            a
          );
        }
      }
      function ce(t, e) {
        return t.filter(function (t) {
          return !e.some(function (e) {
            return (
              e.id + e.btnClass === t.id + t.btnClass &&
              t.callback === e.callback
            );
          });
        });
      }
      var ue = function (t, e) {
          e.forEach(function (e) {
            e.element && (e = e.element()), t.appendChild(e);
          });
        },
        we = (function () {
          function t(e, n, i) {
            var l = this;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(d.g)(this, r.a),
              (this._api = e),
              (this._model = n),
              (this._isMobile = o.OS.mobile),
              (this._volumeAnnouncer = i.querySelector(".jw-volume-update"));
            var c,
              w,
              j,
              h = n.get("localization"),
              f = new Wt(n, e, i.querySelector(".jw-time-update")),
              g = (this.menus = []);
            this.ui = [];
            var b = "",
              y = h.volume;
            if (this._isMobile) {
              if (
                !(n.get("sdkplatform") || (o.OS.iOS && o.OS.version.major < 10))
              ) {
                var v = pt("volume-0,volume-100");
                j = p(
                  "jw-icon-volume",
                  function () {
                    e.setMute();
                  },
                  y,
                  v
                );
              }
            } else {
              (w = document.createElement("div")).className =
                "jw-horizontal-volume-container";
              var m = (c = new ne(
                n,
                "jw-icon-volume",
                y,
                pt("volume-0,volume-50,volume-100"),
                w
              )).element();
              g.push(c),
                Object(s.t)(m, "role", "button"),
                n.change(
                  "mute",
                  function (t, e) {
                    var n = e ? h.unmute : h.mute;
                    Object(s.t)(m, "aria-label", n);
                  },
                  this
                );
            }
            var x = p(
                "jw-icon-next",
                function () {
                  e.next({ feedShownId: b, reason: "interaction" });
                },
                h.next,
                pt("next")
              ),
              k = p(
                "jw-icon-settings jw-settings-submenu-button",
                function (t) {
                  l.trigger("settingsInteraction", "quality", !0, t);
                },
                h.settings,
                pt("settings")
              );
            Object(s.t)(k.element(), "aria-haspopup", "true");
            var O = p(
              "jw-icon-cc jw-settings-submenu-button",
              function (t) {
                l.trigger("settingsInteraction", "captions", !1, t);
              },
              h.cc,
              pt("cc-off,cc-on")
            );
            Object(s.t)(O.element(), "aria-haspopup", "true");
            var C = p(
              "jw-text-live",
              function () {
                l.goToLiveEdge();
              },
              h.liveBroadcast
            );
            C.element().textContent = h.liveBroadcast;
            var S,
              T,
              M,
              z = (this.elements = {
                alt:
                  ((S = "jw-text-alt"),
                  (T = "status"),
                  (M = document.createElement("span")),
                  (M.className = "jw-text jw-reset " + S),
                  T && Object(s.t)(M, "role", T),
                  M),
                play: p(
                  "jw-icon-playback",
                  function () {
                    e.playToggle({ reason: "interaction" });
                  },
                  h.play,
                  pt("play,pause,stop")
                ),
                rewind: p(
                  "jw-icon-rewind",
                  function () {
                    l.rewind();
                  },
                  h.rewind,
                  pt("rewind")
                ),
                live: C,
                next: x,
                elapsed: re("jw-text-elapsed", "timer"),
                countdown: re("jw-text-countdown", "timer"),
                time: f,
                duration: re("jw-text-duration", "timer"),
                mute: j,
                volumetooltip: c,
                horizontalVolumeContainer: w,
                cast: se(function () {
                  e.castToggle();
                }, h),
                fullscreen: p(
                  "jw-icon-fullscreen",
                  function () {
                    e.setFullscreen();
                  },
                  h.fullscreen,
                  pt("fullscreen-off,fullscreen-on")
                ),
                spacer: le("jw-spacer"),
                buttonContainer: le("jw-button-container"),
                settingsButton: k,
                captionsButton: O,
              }),
              E = ie(O.element(), "captions", h.cc),
              L = function (t) {
                var e = t.get("captionsList")[t.get("captionsIndex")],
                  n = h.cc;
                e && "Off" !== e.label && (n = e.label), E.setText(n);
              },
              B = ie(z.play.element(), "play", h.play);
            this.setPlayText = function (t) {
              B.setText(t);
            };
            var _ = z.next.element(),
              V = ie(
                _,
                "next",
                h.nextUp,
                function () {
                  var t = n.get("nextUp");
                  (b = Object(oe.b)(oe.a)),
                    l.trigger("nextShown", {
                      mode: t.mode,
                      ui: "nextup",
                      itemsShown: [t],
                      feedData: t.feedData,
                      reason: "hover",
                      feedShownId: b,
                    });
                },
                function () {
                  b = "";
                }
              );
            Object(s.t)(_, "dir", "auto"),
              ie(z.rewind.element(), "rewind", h.rewind),
              ie(z.settingsButton.element(), "settings", h.settings);
            var A = ie(z.fullscreen.element(), "fullscreen", h.fullscreen),
              N = [
                z.play,
                z.rewind,
                z.next,
                z.volumetooltip,
                z.mute,
                z.horizontalVolumeContainer,
                z.alt,
                z.live,
                z.elapsed,
                z.countdown,
                z.duration,
                z.spacer,
                z.cast,
                z.captionsButton,
                z.settingsButton,
                z.fullscreen,
              ].filter(function (t) {
                return t;
              }),
              H = [z.time, z.buttonContainer].filter(function (t) {
                return t;
              });
            (this.el = document.createElement("div")),
              (this.el.className = "jw-controlbar jw-reset"),
              ue(z.buttonContainer, N),
              ue(this.el, H);
            var P = n.get("logo");
            if (
              (P && "control-bar" === P.position && this.addLogo(P),
              z.play.show(),
              z.fullscreen.show(),
              z.mute && z.mute.show(),
              n.change("volume", this.onVolume, this),
              n.change(
                "mute",
                function (t, e) {
                  l.renderVolume(e, t.get("volume"));
                },
                this
              ),
              n.change("state", this.onState, this),
              n.change("duration", this.onDuration, this),
              n.change("position", this.onElapsed, this),
              n.change(
                "fullscreen",
                function (t, e) {
                  var n = l.elements.fullscreen.element();
                  Object(s.v)(n, "jw-off", e);
                  var i = t.get("fullscreen") ? h.exitFullscreen : h.fullscreen;
                  A.setText(i), Object(s.t)(n, "aria-label", i);
                },
                this
              ),
              n.change("streamType", this.onStreamTypeChange, this),
              n.change(
                "dvrLive",
                function (t, e) {
                  var n = h.liveBroadcast,
                    i = h.notLive,
                    o = l.elements.live.element(),
                    a = !1 === e;
                  Object(s.v)(o, "jw-dvr-live", a),
                    Object(s.t)(o, "aria-label", a ? i : n),
                    (o.textContent = n);
                },
                this
              ),
              n.change("altText", this.setAltText, this),
              n.change("customButtons", this.updateButtons, this),
              n.on("change:captionsIndex", L, this),
              n.on("change:captionsList", L, this),
              n.change(
                "nextUp",
                function (t, e) {
                  b = Object(oe.b)(oe.a);
                  var n = h.nextUp;
                  e && e.title && (n += ": ".concat(e.title)),
                    V.setText(n),
                    z.next.toggle(!!e);
                },
                this
              ),
              n.change("audioMode", this.onAudioMode, this),
              z.cast &&
                (n.change("castAvailable", this.onCastAvailable, this),
                n.change("castActive", this.onCastActive, this)),
              z.volumetooltip &&
                (z.volumetooltip.on(
                  "update",
                  function (t) {
                    var e = t.percentage;
                    this._api.setVolume(e);
                  },
                  this
                ),
                z.volumetooltip.on(
                  "toggleValue",
                  function () {
                    this._api.setMute();
                  },
                  this
                ),
                z.volumetooltip.on(
                  "adjustVolume",
                  function (t) {
                    this.trigger("adjustVolume", t);
                  },
                  this
                )),
              z.cast && z.cast.button)
            ) {
              var I = z.cast.ui.on(
                "click tap enter",
                function (t) {
                  "click" !== t.type && z.cast.button.click(),
                    this._model.set("castClicked", !0);
                },
                this
              );
              this.ui.push(I);
            }
            var R = new u.a(z.duration).on(
              "click tap enter",
              function () {
                if ("DVR" === this._model.get("streamType")) {
                  var t = this._model.get("position"),
                    e = this._model.get("dvrSeekLimit");
                  this._api.seek(Math.max(-e, t), { reason: "interaction" });
                }
              },
              this
            );
            this.ui.push(R);
            var q = new u.a(this.el).on(
              "click tap drag",
              function () {
                this.trigger(a.sb);
              },
              this
            );
            this.ui.push(q),
              g.forEach(function (t) {
                t.on("open-tooltip", l.closeMenus, l);
              });
          }
          var e, n, i;
          return (
            (e = t),
            (n = [
              {
                key: "onVolume",
                value: function (t, e) {
                  this.renderVolume(t.get("mute"), e);
                },
              },
              {
                key: "renderVolume",
                value: function (t, e) {
                  var n = this.elements.mute,
                    i = this.elements.volumetooltip;
                  if (
                    (n &&
                      (Object(s.v)(n.element(), "jw-off", t),
                      Object(s.v)(n.element(), "jw-full", !t)),
                    i)
                  ) {
                    var o = t ? 0 : e,
                      a = i.element();
                    i.verticalSlider.render(o), i.horizontalSlider.render(o);
                    var r = i.tooltip,
                      l = i.horizontalContainer;
                    Object(s.v)(a, "jw-off", t),
                      Object(s.v)(a, "jw-full", e >= 75 && !t),
                      Object(s.t)(r, "aria-valuenow", o),
                      Object(s.t)(l, "aria-valuenow", o);
                    var c = "Volume ".concat(o, "%");
                    Object(s.t)(r, "aria-valuetext", c),
                      Object(s.t)(l, "aria-valuetext", c),
                      document.activeElement !== r &&
                        document.activeElement !== l &&
                        (this._volumeAnnouncer.textContent = c);
                  }
                },
              },
              {
                key: "onCastAvailable",
                value: function (t, e) {
                  this.elements.cast.toggle(e);
                },
              },
              {
                key: "onCastActive",
                value: function (t, e) {
                  this.elements.fullscreen.toggle(!e),
                    this.elements.cast.button &&
                      Object(s.v)(this.elements.cast.button, "jw-off", !e);
                },
              },
              {
                key: "onElapsed",
                value: function (t, e) {
                  var n,
                    i,
                    o = t.get("duration");
                  if ("DVR" === t.get("streamType")) {
                    var a = Math.ceil(e),
                      r = this._model.get("dvrSeekLimit");
                    (n = i =
                      a >= -r ? "" : "-" + Object(vt.timeFormat)(-(e + r))),
                      t.set("dvrLive", a >= -r);
                  } else
                    (n = Object(vt.timeFormat)(e)),
                      (i = Object(vt.timeFormat)(o - e));
                  (this.elements.elapsed.textContent = n),
                    (this.elements.countdown.textContent = i);
                },
              },
              {
                key: "onDuration",
                value: function (t, e) {
                  this.elements.duration.textContent = Object(vt.timeFormat)(
                    Math.abs(e)
                  );
                },
              },
              {
                key: "onAudioMode",
                value: function (t, e) {
                  var n = this.elements.time.element();
                  e
                    ? this.elements.buttonContainer.insertBefore(
                        n,
                        this.elements.elapsed
                      )
                    : Object(s.m)(this.el, n);
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
                value: function (t, e) {
                  this.elements.alt.textContent = e;
                },
              },
              {
                key: "closeMenus",
                value: function (t) {
                  this.menus.forEach(function (e) {
                    (t && t.target === e.el) || e.closeTooltip(t);
                  });
                },
              },
              {
                key: "rewind",
                value: function () {
                  var t,
                    e = 0,
                    n = this._model.get("currentTime");
                  n
                    ? (t = n - 10)
                    : ((t = this._model.get("position") - 10),
                      "DVR" === this._model.get("streamType") &&
                        (e = this._model.get("duration"))),
                    this._api.seek(Math.max(t, e), { reason: "interaction" });
                },
              },
              {
                key: "onState",
                value: function (t, e) {
                  var n = t.get("localization"),
                    i = n.play;
                  this.setPlayText(i),
                    e === a.pb &&
                      ("LIVE" !== t.get("streamType")
                        ? ((i = n.pause), this.setPlayText(i))
                        : ((i = n.stop), this.setPlayText(i))),
                    Object(s.t)(this.elements.play.element(), "aria-label", i);
                },
              },
              {
                key: "onStreamTypeChange",
                value: function (t, e) {
                  var n = "LIVE" === e,
                    i = "DVR" === e;
                  this.elements.rewind.toggle(!n),
                    this.elements.live.toggle(n || i),
                    Object(s.t)(
                      this.elements.live.element(),
                      "tabindex",
                      n ? "-1" : "0"
                    ),
                    (this.elements.duration.style.display = i ? "none" : ""),
                    this.onDuration(t, t.get("duration")),
                    this.onState(t, t.get("state"));
                },
              },
              {
                key: "addLogo",
                value: function (t) {
                  var e = this.elements.buttonContainer,
                    n = new yt(
                      t.file,
                      this._model.get("localization").logo,
                      function () {
                        t.link &&
                          Object(s.l)(t.link, "_blank", { rel: "noreferrer" });
                      },
                      "logo",
                      "jw-logo-button"
                    );
                  t.link || Object(s.t)(n.element(), "tabindex", "-1"),
                    e.insertBefore(
                      n.element(),
                      e.querySelector(".jw-spacer").nextSibling
                    );
                },
              },
              {
                key: "goToLiveEdge",
                value: function () {
                  if ("DVR" === this._model.get("streamType")) {
                    var t = Math.min(this._model.get("position"), -1),
                      e = this._model.get("dvrSeekLimit");
                    this._api.seek(Math.max(-e, t), { reason: "interaction" }),
                      this._api.play({ reason: "interaction" });
                  }
                },
              },
              {
                key: "updateButtons",
                value: function (t, e, n) {
                  if (e) {
                    var i,
                      o,
                      a = this.elements.buttonContainer;
                    e !== n && n
                      ? ((i = ce(e, n)),
                        (o = ce(n, e)),
                        this.removeButtons(a, o))
                      : (i = e);
                    for (var r = i.length - 1; r >= 0; r--) {
                      var l = i[r],
                        s = new yt(
                          l.img,
                          l.tooltip,
                          l.callback,
                          l.id,
                          l.btnClass
                        );
                      l.tooltip && ie(s.element(), l.id, l.tooltip);
                      var c = void 0;
                      "related" === s.id
                        ? (c = this.elements.settingsButton.element())
                        : "share" === s.id
                        ? (c =
                            a.querySelector('[button="related"]') ||
                            this.elements.settingsButton.element())
                        : (c = this.elements.spacer.nextSibling) &&
                          "logo" === c.getAttribute("button") &&
                          (c = c.nextSibling),
                        a.insertBefore(s.element(), c);
                    }
                  }
                },
              },
              {
                key: "removeButtons",
                value: function (t, e) {
                  for (var n = e.length; n--; ) {
                    var i = t.querySelector('[button="'.concat(e[n].id, '"]'));
                    i && t.removeChild(i);
                  }
                },
              },
              {
                key: "toggleCaptionsButtonState",
                value: function (t) {
                  var e = this.elements.captionsButton;
                  e && Object(s.v)(e.element(), "jw-off", !t);
                },
              },
              {
                key: "destroy",
                value: function () {
                  var t = this;
                  this._model.off(null, null, this),
                    Object.keys(this.elements).forEach(function (e) {
                      var n = t.elements[e];
                      n &&
                        "function" == typeof n.destroy &&
                        t.elements[e].destroy();
                    }),
                    this.ui.forEach(function (t) {
                      t.destroy();
                    }),
                    (this.ui = []);
                },
              },
            ]) && ae(e.prototype, n),
            i && ae(e, i),
            t
          );
        })(),
        pe = function () {
          var t =
              arguments.length > 0 && void 0 !== arguments[0]
                ? arguments[0]
                : "",
            e =
              arguments.length > 1 && void 0 !== arguments[1]
                ? arguments[1]
                : "";
          return (
            '<div class="jw-display-icon-container jw-display-icon-'.concat(
              t,
              ' jw-reset">'
            ) +
            '<div class="jw-icon jw-icon-'
              .concat(
                t,
                ' jw-button-color jw-reset" role="button" tabindex="0" aria-label="'
              )
              .concat(e, '"></div>') +
            "</div>"
          );
        },
        de = function (t) {
          return (
            '<div class="jw-display jw-reset"><div class="jw-display-container jw-reset"><div class="jw-display-controls jw-reset">' +
            pe("rewind", t.rewind) +
            pe("display", t.playback) +
            pe("next", t.next) +
            "</div></div></div>"
          );
        };
      function je(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var he = (function () {
        function t(e, n, i) {
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, t);
          var o = i.querySelector(".jw-icon");
          (this.el = i),
            (this.ui = new u.a(o).on("click tap enter", function () {
              var t = e.get("position"),
                i = e.get("duration"),
                o = t - 10,
                a = 0;
              "DVR" === e.get("streamType") && (a = i), n.seek(Math.max(o, a));
            }));
        }
        var e, n, i;
        return (
          (e = t),
          (n = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && je(e.prototype, n),
          i && je(e, i),
          t
        );
      })();
      function fe(t) {
        return (fe =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (t) {
                return typeof t;
              }
            : function (t) {
                return t &&
                  "function" == typeof Symbol &&
                  t.constructor === Symbol &&
                  t !== Symbol.prototype
                  ? "symbol"
                  : typeof t;
              })(t);
      }
      function ge(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function be(t, e) {
        return !e || ("object" !== fe(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function ye(t) {
        return (ye = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function ve(t, e) {
        return (ve =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var me = (function (t) {
        function e(t, n, i) {
          var o;
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, e),
            (o = be(this, ye(e).call(this)));
          var a = t.get("localization"),
            r = i.querySelector(".jw-icon");
          if (
            ((o.icon = r),
            (o.el = i),
            (o.ui = new u.a(r).on("click tap enter", function (t) {
              o.trigger(t.type);
            })),
            t.on("change:state", function (t, e) {
              var n;
              switch (e) {
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
            t.get("displayPlaybackLabel"))
          ) {
            var l = o.icon.getElementsByClassName("jw-idle-icon-text")[0];
            l ||
              ((l = Object(s.e)(
                '<div class="jw-idle-icon-text">'.concat(a.playback, "</div>")
              )),
              Object(s.a)(o.icon, "jw-idle-label"),
              o.icon.appendChild(l));
          }
          return o;
        }
        var n, i, o;
        return (
          (function (t, e) {
            if ("function" != typeof e && null !== e)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (t.prototype = Object.create(e && e.prototype, {
              constructor: { value: t, writable: !0, configurable: !0 },
            })),
              e && ve(t, e);
          })(e, t),
          (n = e),
          (i = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && ge(n.prototype, i),
          o && ge(n, o),
          e
        );
      })(r.a);
      function xe(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var ke = (function () {
        function t(e, n, i) {
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, t);
          var o = i.querySelector(".jw-icon");
          (this.ui = new u.a(o).on("click tap enter", function () {
            n.next({ reason: "interaction" });
          })),
            e.change("nextUp", function (t, e) {
              i.style.visibility = e ? "" : "hidden";
            }),
            (this.el = i);
        }
        var e, n, i;
        return (
          (e = t),
          (n = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && xe(e.prototype, n),
          i && xe(e, i),
          t
        );
      })();
      function Oe(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var Ce = (function () {
        function t(e, n) {
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, t),
            (this.el = Object(s.e)(de(e.get("localization"))));
          var i = this.el.querySelector(".jw-display-controls"),
            o = {};
          Se("rewind", pt("rewind"), he, i, o, e, n),
            Se("display", pt("play,pause,buffer,replay"), me, i, o, e, n),
            Se("next", pt("next"), ke, i, o, e, n),
            (this.container = i),
            (this.buttons = o);
        }
        var e, n, i;
        return (
          (e = t),
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
                var t = this.buttons;
                Object.keys(t).forEach(function (e) {
                  t[e].ui && t[e].ui.destroy();
                });
              },
            },
          ]) && Oe(e.prototype, n),
          i && Oe(e, i),
          t
        );
      })();
      function Se(t, e, n, i, o, a, r) {
        var l = i.querySelector(".jw-display-icon-".concat(t)),
          s = i.querySelector(".jw-icon-".concat(t));
        e.forEach(function (t) {
          s.appendChild(t);
        }),
          (o[t] = new n(a, r, l));
      }
      var Te = n(2);
      function Me(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var ze = (function () {
          function t(e, n, i) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(d.g)(this, r.a),
              (this._model = e),
              (this._api = n),
              (this._playerElement = i),
              (this.localization = e.get("localization")),
              (this.state = "tooltip"),
              (this.enabled = !1),
              (this.shown = !1),
              (this.feedShownId = ""),
              (this.closeUi = null),
              (this.tooltipUi = null),
              this.reset();
          }
          var e, n, i;
          return (
            (e = t),
            (n = [
              {
                key: "setup",
                value: function (t) {
                  (this.container = t.createElement("div")),
                    (this.container.className = "jw-nextup-container jw-reset");
                  var e = Object(s.e)(
                    (function () {
                      var t =
                          arguments.length > 0 && void 0 !== arguments[0]
                            ? arguments[0]
                            : "",
                        e =
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
                          t,
                          "</div>"
                        ) +
                        '<div class="jw-nextup-title jw-reset-text" dir="auto">'.concat(
                          e,
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
                  e.querySelector(".jw-nextup-close").appendChild(wt("close")),
                    this.addContent(e),
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
                      function (t, e) {
                        "complete" === e && this.toggle(!1);
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
                value: function (t) {
                  return (
                    (this.nextUpImage = new Image()),
                    (this.nextUpImage.onload = function () {
                      this.nextUpImage.onload = null;
                    }.bind(this)),
                    (this.nextUpImage.src = t),
                    { backgroundImage: 'url("' + t + '")' }
                  );
                },
              },
              {
                key: "click",
                value: function () {
                  var t = this.feedShownId;
                  this.reset(),
                    this._api.next({ feedShownId: t, reason: "interaction" });
                },
              },
              {
                key: "toggle",
                value: function (t, e) {
                  if (
                    this.enabled &&
                    (Object(s.v)(
                      this.container,
                      "jw-nextup-sticky",
                      !!this.nextUpSticky
                    ),
                    this.shown !== t)
                  ) {
                    (this.shown = t),
                      Object(s.v)(
                        this.container,
                        "jw-nextup-container-visible",
                        t
                      ),
                      Object(s.v)(this._playerElement, "jw-flag-nextup", t);
                    var n = this._model.get("nextUp");
                    t && n
                      ? ((this.feedShownId = Object(oe.b)(oe.a)),
                        this.trigger("nextShown", {
                          mode: n.mode,
                          ui: "nextup",
                          itemsShown: [n],
                          feedData: n.feedData,
                          reason: e,
                          feedShownId: this.feedShownId,
                        }))
                      : (this.feedShownId = "");
                  }
                },
              },
              {
                key: "setNextUpItem",
                value: function (t) {
                  var e = this;
                  setTimeout(function () {
                    if (
                      ((e.thumbnail = e.content.querySelector(
                        ".jw-nextup-thumbnail"
                      )),
                      Object(s.v)(
                        e.content,
                        "jw-nextup-thumbnail-visible",
                        !!t.image
                      ),
                      t.image)
                    ) {
                      var n = e.loadThumbnail(t.image);
                      Object(ft.d)(e.thumbnail, n);
                    }
                    (e.header = e.content.querySelector(".jw-nextup-header")),
                      (e.header.textContent = Object(s.e)(
                        e.localization.nextUp
                      ).textContent),
                      (e.title = e.content.querySelector(".jw-nextup-title"));
                    var i = t.title;
                    e.title.textContent = i ? Object(s.e)(i).textContent : "";
                    var o = t.duration;
                    o &&
                      ((e.duration = e.content.querySelector(
                        ".jw-nextup-duration"
                      )),
                      (e.duration.textContent =
                        "number" == typeof o ? Object(vt.timeFormat)(o) : o));
                  }, 500);
                },
              },
              {
                key: "onNextUp",
                value: function (t, e) {
                  this.reset(),
                    e || (e = { showNextUp: !1 }),
                    (this.enabled = !(!e.title && !e.image)),
                    this.enabled &&
                      (e.showNextUp ||
                        ((this.nextUpSticky = !1), this.toggle(!1)),
                      this.setNextUpItem(e));
                },
              },
              {
                key: "onDuration",
                value: function (t, e) {
                  if (e) {
                    var n = t.get("nextupoffset"),
                      i = -10;
                    n && (i = Object(Te.d)(n, e)),
                      i < 0 && (i += e),
                      Object(Te.c)(n) && e - 5 < i && (i = e - 5),
                      (this.offset = i);
                  }
                },
              },
              {
                key: "onElapsed",
                value: function (t, e) {
                  var n = this.nextUpSticky;
                  if (this.enabled && !1 !== n) {
                    var i = e >= this.offset;
                    i && void 0 === n
                      ? ((this.nextUpSticky = i), this.toggle(i, "time"))
                      : !i && n && this.reset();
                  }
                },
              },
              {
                key: "onStreamType",
                value: function (t, e) {
                  "VOD" !== e && ((this.nextUpSticky = !1), this.toggle(!1));
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
                value: function (t) {
                  this.content && this.removeContent(),
                    (this.content = t),
                    this.container.appendChild(t);
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
            ]) && Me(e.prototype, n),
            i && Me(e, i),
            t
          );
        })(),
        Ee = function (t, e) {
          var n = t.featured,
            i = t.showLogo,
            o = t.type;
          return (
            (t.logo = i
              ? '<span class="jw-rightclick-logo jw-reset"></span>'
              : ""),
            '<li class="jw-reset jw-rightclick-item '
              .concat(n ? "jw-featured" : "", '">')
              .concat(Le[o](t, e), "</li>")
          );
        },
        Le = {
          link: function (t) {
            var e = t.link,
              n = t.title,
              i = t.logo;
            return '<a href="'
              .concat(
                e || "",
                '" class="jw-rightclick-link jw-reset-text" target="_blank" rel="noreferrer" dir="auto">'
              )
              .concat(i)
              .concat(n || "", "</a>");
          },
          info: function (t, e) {
            return '<button type="button" class="jw-reset-text jw-rightclick-link jw-info-overlay-item" dir="auto">'.concat(
              e.videoInfo,
              "</button>"
            );
          },
          share: function (t, e) {
            return '<button type="button" class="jw-reset-text jw-rightclick-link jw-share-item" dir="auto">'.concat(
              e.sharing.heading,
              "</button>"
            );
          },
          keyboardShortcuts: function (t, e) {
            return '<button type="button" class="jw-reset-text jw-rightclick-link jw-shortcuts-item" dir="auto">'.concat(
              e.shortcuts.keyboardShortcuts,
              "</button>"
            );
          },
        },
        Be = n(23),
        _e = n(6),
        Ve = n(13);
      function Ae(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var Ne = {
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
      function He(t) {
        var e = Object(s.e)(t),
          n = e.querySelector(".jw-rightclick-logo");
        return n && n.appendChild(wt("jwplayer-logo")), e;
      }
      var Pe = (function () {
          function t(e, n) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              (this.infoOverlay = e),
              (this.shortcutsTooltip = n);
          }
          var e, n, i;
          return (
            (e = t),
            (n = [
              {
                key: "buildArray",
                value: function () {
                  var t = Be.a.split("+")[0],
                    e = this.model,
                    n = e.get("edition"),
                    i = e.get("localization").poweredBy,
                    o = '<span class="jw-reset">JW Player '.concat(
                      t,
                      "</span>"
                    ),
                    a = {
                      items: [
                        { type: "info" },
                        {
                          title: Object(Ve.e)(i)
                            ? "".concat(o, " ").concat(i)
                            : "".concat(i, " ").concat(o),
                          type: "link",
                          featured: !0,
                          showLogo: !0,
                          link: "https://jwplayer.com/learn-more?e=".concat(
                            Ne[n]
                          ),
                        },
                      ],
                    },
                    r = e.get("provider"),
                    l = a.items;
                  if (r && r.name.indexOf("flash") >= 0) {
                    var s = "Flash Version " + Object(_e.a)();
                    l.push({
                      title: s,
                      type: "link",
                      link: "http://www.adobe.com/software/flash/about/",
                    });
                  }
                  return (
                    this.shortcutsTooltip &&
                      l.splice(l.length - 1, 0, { type: "keyboardShortcuts" }),
                    a
                  );
                },
              },
              {
                key: "rightClick",
                value: function (t) {
                  if ((this.lazySetup(), this.mouseOverContext)) return !1;
                  this.hideMenu(), this.showMenu(t), this.addHideMenuHandlers();
                },
              },
              {
                key: "getOffset",
                value: function (t) {
                  var e = Object(s.c)(this.wrapperElement),
                    n = t.pageX - e.left,
                    i = t.pageY - e.top;
                  return (
                    this.model.get("touchMode") && (i -= 100), { x: n, y: i }
                  );
                },
              },
              {
                key: "showMenu",
                value: function (t) {
                  var e = this,
                    n = this.getOffset(t);
                  return (
                    (this.el.style.left = n.x + "px"),
                    (this.el.style.top = n.y + "px"),
                    (this.outCount = 0),
                    Object(s.a)(
                      this.playerContainer,
                      "jw-flag-rightclick-open"
                    ),
                    Object(s.a)(this.el, "jw-open"),
                    clearTimeout(this._menuTimeout),
                    (this._menuTimeout = setTimeout(function () {
                      return e.hideMenu();
                    }, 3e3)),
                    !1
                  );
                },
              },
              {
                key: "hideMenu",
                value: function (t) {
                  (t && this.el && this.el.contains(t.target)) ||
                    (Object(s.o)(
                      this.playerContainer,
                      "jw-flag-rightclick-open"
                    ),
                    Object(s.o)(this.el, "jw-open"));
                },
              },
              {
                key: "lazySetup",
                value: function () {
                  var t,
                    e,
                    n,
                    i,
                    o = this,
                    a =
                      ((t = this.buildArray()),
                      (e = this.model.get("localization")),
                      (n = t.items),
                      (i = (void 0 === n ? [] : n).map(function (t) {
                        return Ee(t, e);
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
                      var r = He(a);
                      Object(s.h)(this.el);
                      for (var l = r.childNodes.length; l--; )
                        this.el.appendChild(r.firstChild);
                    }
                  } else
                    (this.html = a),
                      (this.el = He(this.html)),
                      this.wrapperElement.appendChild(this.el),
                      (this.hideMenuHandler = function (t) {
                        return o.hideMenu(t);
                      }),
                      (this.overHandler = function () {
                        o.mouseOverContext = !0;
                      }),
                      (this.outHandler = function (t) {
                        (o.mouseOverContext = !1),
                          t.relatedTarget &&
                            !o.el.contains(t.relatedTarget) &&
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
                value: function (t, e, n) {
                  (this.wrapperElement = n),
                    (this.model = t),
                    (this.mouseOverContext = !1),
                    (this.playerContainer = e),
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
            ]) && Ae(e.prototype, n),
            i && Ae(e, i),
            t
          );
        })(),
        Ie = function (t) {
          return '<button type="button" class="jw-reset-text jw-settings-content-item" dir="auto">'.concat(
            t,
            "</button>"
          );
        },
        Re = function (t) {
          return (
            '<button type="button" class="jw-reset-text jw-settings-content-item" dir="auto">' +
            "".concat(t.label) +
            "<div class='jw-reset jw-settings-value-wrapper'>" +
            '<div class="jw-reset-text jw-settings-content-item-value">'.concat(
              t.value,
              "</div>"
            ) +
            '<div class="jw-reset-text jw-settings-content-item-arrow">'.concat(
              K.a,
              "</div>"
            ) +
            "</div></button>"
          );
        },
        qe = function (t) {
          return (
            '<button type="button" class="jw-reset-text jw-settings-content-item" role="menuitemradio" aria-checked="false" dir="auto">' +
            "".concat(t) +
            "</button>"
          );
        };
      function De(t) {
        return (De =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (t) {
                return typeof t;
              }
            : function (t) {
                return t &&
                  "function" == typeof Symbol &&
                  t.constructor === Symbol &&
                  t !== Symbol.prototype
                  ? "symbol"
                  : typeof t;
              })(t);
      }
      function Ue(t, e) {
        return !e || ("object" !== De(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Fe(t) {
        return (Fe = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function We(t, e) {
        return (We =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      function Ze(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Ke(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function Xe(t, e, n) {
        return e && Ke(t.prototype, e), n && Ke(t, n), t;
      }
      var Ye,
        Ge = (function () {
          function t(e, n) {
            var i =
              arguments.length > 2 && void 0 !== arguments[2]
                ? arguments[2]
                : Ie;
            Ze(this, t),
              (this.el = Object(s.e)(i(e))),
              (this.ui = new u.a(this.el).on("click tap enter", n, this));
          }
          return (
            Xe(t, [
              {
                key: "destroy",
                value: function () {
                  this.ui.destroy();
                },
              },
            ]),
            t
          );
        })(),
        Je = (function (t) {
          function e(t, n) {
            var i =
              arguments.length > 2 && void 0 !== arguments[2]
                ? arguments[2]
                : qe;
            return Ze(this, e), Ue(this, Fe(e).call(this, t, n, i));
          }
          return (
            (function (t, e) {
              if ("function" != typeof e && null !== e)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (t.prototype = Object.create(e && e.prototype, {
                constructor: { value: t, writable: !0, configurable: !0 },
              })),
                e && We(t, e);
            })(e, t),
            Xe(e, [
              {
                key: "activate",
                value: function () {
                  Object(s.v)(this.el, "jw-settings-item-active", !0),
                    this.el.setAttribute("aria-checked", "true"),
                    (this.active = !0);
                },
              },
              {
                key: "deactivate",
                value: function () {
                  Object(s.v)(this.el, "jw-settings-item-active", !1),
                    this.el.setAttribute("aria-checked", "false"),
                    (this.active = !1);
                },
              },
            ]),
            e
          );
        })(Ge),
        Qe = function (t, e) {
          return t
            ? '<div class="jw-reset jw-settings-submenu jw-settings-submenu-'.concat(
                e,
                '" role="menu" aria-expanded="false">'
              ) + '<div class="jw-settings-submenu-items"></div></div>'
            : '<div class="jw-reset jw-settings-menu" role="menu" aria-expanded="false"><div class="jw-reset jw-settings-topbar" role="menubar"><div class="jw-reset jw-settings-topbar-text" tabindex="0"></div><div class="jw-reset jw-settings-topbar-buttons"></div></div></div>';
        },
        $e = function (t, e) {
          var n = t.name,
            i = {
              captions: "cc-off",
              audioTracks: "audio-tracks",
              quality: "quality-100",
              playbackRates: "playback-rate",
            }[n];
          if (i || t.icon) {
            var o = p(
                "jw-settings-".concat(n, " jw-submenu-").concat(n),
                function (e) {
                  t.open(e);
                },
                n,
                [(t.icon && Object(s.e)(t.icon)) || wt(i)]
              ),
              a = o.element();
            return (
              a.setAttribute("role", "menuitemradio"),
              a.setAttribute("aria-checked", "false"),
              a.setAttribute("aria-label", e),
              "ontouchstart" in window || (o.tooltip = ie(a, n, e)),
              o
            );
          }
        };
      function tn(t) {
        return (tn =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (t) {
                return typeof t;
              }
            : function (t) {
                return t &&
                  "function" == typeof Symbol &&
                  t.constructor === Symbol &&
                  t !== Symbol.prototype
                  ? "symbol"
                  : typeof t;
              })(t);
      }
      function en(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function nn(t) {
        return (nn = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function on(t, e) {
        return (on =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      function an(t) {
        if (void 0 === t)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return t;
      }
      var rn = (function (t) {
          function e(t, n, i) {
            var o,
              a,
              r,
              l =
                arguments.length > 3 && void 0 !== arguments[3]
                  ? arguments[3]
                  : Qe;
            return (
              (function (t, e) {
                if (!(t instanceof e))
                  throw new TypeError("Cannot call a class as a function");
              })(this, e),
              (a = this),
              ((o =
                !(r = nn(e).call(this)) ||
                ("object" !== tn(r) && "function" != typeof r)
                  ? an(a)
                  : r).open = o.open.bind(an(an(o)))),
              (o.close = o.close.bind(an(an(o)))),
              (o.toggle = o.toggle.bind(an(an(o)))),
              (o.onDocumentClick = o.onDocumentClick.bind(an(an(o)))),
              (o.name = t),
              (o.isSubmenu = !!n),
              (o.el = Object(s.e)(l(o.isSubmenu, t))),
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
                : (o.ui = ln(an(an(o)))),
              o
            );
          }
          var n, i, o;
          return (
            (function (t, e) {
              if ("function" != typeof e && null !== e)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (t.prototype = Object.create(e && e.prototype, {
                constructor: { value: t, writable: !0, configurable: !0 },
              })),
                e && on(t, e);
            })(e, t),
            (n = e),
            (i = [
              {
                key: "createItemsContainer",
                value: function () {
                  var t,
                    e,
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
                      ((t = this.mainMenu.closeButton.element()),
                      (e = this.mainMenu.backButton.element())),
                    o.on("keydown", function (o) {
                      if (o.target.parentNode === i) {
                        var r = function (t, e) {
                            t
                              ? t.focus()
                              : void 0 !== e && i.childNodes[e].focus();
                          },
                          l = o.sourceEvent,
                          c = l.target,
                          u = i.firstChild === c,
                          w = i.lastChild === c,
                          p = n.topbar,
                          d = t || Object(s.k)(a),
                          j = e || Object(s.n)(a),
                          h = Object(s.k)(l.target),
                          f = Object(s.n)(l.target),
                          g = l.key.replace(/(Arrow|ape)/, "");
                        switch (g) {
                          case "Tab":
                            r(l.shiftKey ? j : d);
                            break;
                          case "Left":
                            r(
                              j ||
                                Object(s.n)(
                                  document.getElementsByClassName(
                                    "jw-icon-settings"
                                  )[0]
                                )
                            );
                            break;
                          case "Up":
                            p && u
                              ? r(p.firstChild)
                              : r(f, i.childNodes.length - 1);
                            break;
                          case "Right":
                            r(d);
                            break;
                          case "Down":
                            p && w ? r(p.firstChild) : r(h, 0);
                        }
                        l.preventDefault(), "Esc" !== g && l.stopPropagation();
                      }
                    }),
                    o
                  );
                },
              },
              {
                key: "createCloseButton",
                value: function (t) {
                  var e = p("jw-settings-close", this.close, t.close, [
                    wt("close"),
                  ]);
                  return (
                    this.topbar.appendChild(e.element()),
                    e.show(),
                    e.ui.on(
                      "keydown",
                      function (t) {
                        var e = t.sourceEvent,
                          n = e.key.replace(/(Arrow|ape)/, "");
                        ("Enter" === n ||
                          "Right" === n ||
                          ("Tab" === n && !e.shiftKey)) &&
                          this.close(t);
                      },
                      this
                    ),
                    this.buttonContainer.appendChild(e.element()),
                    e
                  );
                },
              },
              {
                key: "createCategoryButton",
                value: function (t) {
                  var e =
                    t[
                      {
                        captions: "cc",
                        audioTracks: "audioTracks",
                        quality: "hd",
                        playbackRates: "playbackRates",
                      }[this.name]
                    ];
                  "sharing" === this.name && (e = t.sharing.heading);
                  var n = $e(this, e);
                  return n.element().setAttribute("name", this.name), n;
                },
              },
              {
                key: "createBackButton",
                value: function (t) {
                  var e = p(
                    "jw-settings-back",
                    function (t) {
                      Ye && Ye.open(t);
                    },
                    t.close,
                    [wt("arrow-left")]
                  );
                  return Object(s.m)(this.mainMenu.topbar, e.element()), e;
                },
              },
              {
                key: "createTopbar",
                value: function () {
                  var t = Object(s.e)('<div class="jw-submenu-topbar"></div>');
                  return Object(s.m)(this.el, t), t;
                },
              },
              {
                key: "createItems",
                value: function (t, e) {
                  var n = this,
                    i =
                      arguments.length > 2 && void 0 !== arguments[2]
                        ? arguments[2]
                        : {},
                    o =
                      arguments.length > 3 && void 0 !== arguments[3]
                        ? arguments[3]
                        : Je,
                    a = this.name,
                    r = t.map(function (t, r) {
                      var l, s;
                      switch (a) {
                        case "quality":
                          l =
                            "Auto" === t.label && 0 === r
                              ? "".concat(
                                  i.defaultText,
                                  '&nbsp;<span class="jw-reset jw-auto-label"></span>'
                                )
                              : t.label;
                          break;
                        case "captions":
                          l =
                            ("Off" !== t.label && "off" !== t.id) || 0 !== r
                              ? t.label
                              : i.defaultText;
                          break;
                        case "playbackRates":
                          (s = t),
                            (l = Object(Ve.e)(i.tooltipText)
                              ? "x" + t
                              : t + "x");
                          break;
                        case "audioTracks":
                          l = t.name;
                      }
                      l || ((l = t), "object" === tn(t) && (l.options = i));
                      var c = new o(
                        l,
                        function (t) {
                          c.active ||
                            (e(s || r),
                            c.deactivate &&
                              (n.items
                                .filter(function (t) {
                                  return !0 === t.active;
                                })
                                .forEach(function (t) {
                                  t.deactivate();
                                }),
                              Ye ? Ye.open(t) : n.mainMenu.close(t)),
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
                value: function (t, e) {
                  var n = this;
                  t
                    ? ((this.items = []),
                      Object(s.h)(this.itemsContainer.el),
                      t.forEach(function (t) {
                        n.items.push(t), n.itemsContainer.el.appendChild(t.el);
                      }),
                      e > -1 && t[e].activate(),
                      this.categoryButton.show())
                    : this.removeMenu();
                },
              },
              {
                key: "appendMenu",
                value: function (t) {
                  if (t) {
                    var e = t.el,
                      n = t.name,
                      i = t.categoryButton;
                    if (((this.children[n] = t), i)) {
                      var o = this.mainMenu.buttonContainer,
                        a = o.querySelector(".jw-settings-sharing"),
                        r =
                          "quality" === n
                            ? o.firstChild
                            : a || this.closeButton.element();
                      o.insertBefore(i.element(), r);
                    }
                    this.mainMenu.el.appendChild(e);
                  }
                },
              },
              {
                key: "removeMenu",
                value: function (t) {
                  if (!t) return this.parentMenu.removeMenu(this.name);
                  var e = this.children[t];
                  e && (delete this.children[t], e.destroy());
                },
              },
              {
                key: "open",
                value: function (t) {
                  if (!this.visible || this.openMenus) {
                    var e;
                    if (((Ye = null), this.isSubmenu)) {
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
                          (Ye = this.parentMenu),
                          (e = this.topbar
                            ? this.topbar.firstChild
                            : t && "enter" === t.type
                            ? this.items[0].el
                            : a);
                      } else
                        n.topbar.classList.remove("jw-nested-menu-open"),
                          n.backButton && n.backButton.hide();
                      this.el.classList.add("jw-settings-submenu-active"),
                        i.openMenus.push(this.name),
                        n.visible ||
                          (n.open(t),
                          this.items && t && "enter" === t.type
                            ? (e = this.topbar
                                ? this.topbar.firstChild.focus()
                                : this.items[0].el)
                            : o.tooltip &&
                              ((o.tooltip.suppress = !0), (e = o.element()))),
                        this.openMenus.length && this.closeChildren(),
                        e && e.focus(),
                        (this.el.scrollTop = 0);
                    } else
                      this.el.parentNode.classList.add("jw-settings-open"),
                        this.trigger("menuVisibility", { visible: !0, evt: t }),
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
                value: function (t) {
                  var e = this;
                  this.visible &&
                    ((this.visible = !1),
                    this.el.setAttribute("aria-expanded", "false"),
                    this.isSubmenu
                      ? (this.el.classList.remove("jw-settings-submenu-active"),
                        this.categoryButton
                          .element()
                          .setAttribute("aria-checked", "false"),
                        (this.parentMenu.openMenus = this.parentMenu.openMenus.filter(
                          function (t) {
                            return t !== e.name;
                          }
                        )),
                        !this.mainMenu.openMenus.length &&
                          this.mainMenu.visible &&
                          this.mainMenu.close(t))
                      : (this.el.parentNode.classList.remove(
                          "jw-settings-open"
                        ),
                        this.trigger("menuVisibility", { visible: !1, evt: t }),
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
                  var t = this;
                  this.openMenus.forEach(function (e) {
                    var n = t.children[e];
                    n && n.close();
                  });
                },
              },
              {
                key: "toggle",
                value: function (t) {
                  this.visible ? this.close(t) : this.open(t);
                },
              },
              {
                key: "onDocumentClick",
                value: function (t) {
                  /jw-(settings|video|nextup-close|sharing-link|share-item)/.test(
                    t.target.className
                  ) || this.close();
                },
              },
              {
                key: "destroy",
                value: function () {
                  var t = this;
                  if (
                    (document.removeEventListener(
                      "click",
                      this.onDocumentClick
                    ),
                    Object.keys(this.children).map(function (e) {
                      t.children[e].destroy();
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
                    var e = this.parentMenu.openMenus,
                      n = e.indexOf(this.name);
                    e.length && n > -1 && this.openMenus.splice(n, 1),
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
                  var t = this.children,
                    e = t.quality,
                    n = t.captions,
                    i = t.audioTracks,
                    o = t.sharing,
                    a = t.playbackRates;
                  return e || n || i || o || a;
                },
              },
            ]) && en(n.prototype, i),
            o && en(n, o),
            e
          );
        })(r.a),
        ln = function (t) {
          var e = t.closeButton,
            n = t.el;
          return new u.a(n).on("keydown", function (n) {
            var i = n.sourceEvent,
              o = n.target,
              a = Object(s.k)(o),
              r = Object(s.n)(o),
              l = i.key.replace(/(Arrow|ape)/, ""),
              c = function (e) {
                r ? e || r.focus() : t.close(n);
              };
            switch (l) {
              case "Esc":
                t.close(n);
                break;
              case "Left":
                c();
                break;
              case "Right":
                a && e.element() && o !== e.element() && a.focus();
                break;
              case "Tab":
                i.shiftKey && c(!0);
                break;
              case "Up":
              case "Down":
                !(function () {
                  var e = t.children[o.getAttribute("name")];
                  if ((!e && Ye && (e = Ye.children[Ye.openMenus]), e))
                    return (
                      e.open(n),
                      void (e.topbar
                        ? e.topbar.firstChild.focus()
                        : e.items && e.items.length && e.items[0].el.focus())
                    );
                  if (
                    n.target.parentNode.classList.contains("jw-submenu-topbar")
                  ) {
                    var i = n.target.parentNode.parentNode.querySelector(
                      ".jw-settings-submenu-items"
                    );
                    ("Down" === l
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
        sn = n(59),
        cn = function (t) {
          return jn[t];
        },
        un = function (t) {
          for (var e, n = Object.keys(jn), i = 0; i < n.length; i++)
            if (jn[n[i]] === t) {
              e = n[i];
              break;
            }
          return e;
        },
        wn = function (t) {
          return t + "%";
        },
        pn = function (t) {
          return parseInt(t);
        },
        dn = [
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
            getOption: wn,
          },
          {
            name: "Font Size",
            propertyName: "userFontScale",
            options: ["200%", "175%", "150%", "125%", "100%", "75%", "50%"],
            defaultVal: "100%",
            getTypedValue: function (t) {
              return parseInt(t) / 100;
            },
            getOption: function (t) {
              return 100 * t + "%";
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
            getTypedValue: function (t) {
              return t;
            },
            getOption: function (t) {
              return t;
            },
          },
          {
            name: "Character Edge",
            propertyName: "edgeStyle",
            options: ["None", "Raised", "Depressed", "Uniform", "Drop Shadow"],
            defaultVal: "None",
            getTypedValue: function (t) {
              return t.toLowerCase().replace(/ /g, "");
            },
            getOption: function (t) {
              if ("dropshadow" === t) return "Drop Shadow";
              var e = t.replace(/([A-Z])/g, " $1");
              return e.charAt(0).toUpperCase() + e.slice(1);
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
            getOption: wn,
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
            getOption: wn,
          },
        ],
        jn = {
          White: "#ffffff",
          Black: "#000000",
          Red: "#ff0000",
          Green: "#00ff00",
          Blue: "#0000ff",
          Yellow: "#ffff00",
          Magenta: "ff00ff",
          Cyan: "#00ffff",
        },
        hn = function (t, e, n, i) {
          var o = new rn("settings", null, i),
            a = function (t, e, a, r, l) {
              var s = n.elements["".concat(t, "Button")];
              if (!e || e.length <= 1)
                return o.removeMenu(t), void (s && s.hide());
              var c = o.children[t];
              c || (c = new rn(t, o, i)),
                c.setMenuItems(c.createItems(e, a, l), r),
                s && s.show();
            },
            r = function (r) {
              var l = { defaultText: i.auto };
              a(
                "quality",
                r,
                function (e) {
                  return t.setCurrentQuality(e);
                },
                e.get("currentLevel") || 0,
                l
              );
              var s = o.children,
                c = !!s.quality || s.playbackRates || Object.keys(s).length > 1;
              n.elements.settingsButton.toggle(c);
            };
          e.change(
            "levels",
            function (t, e) {
              r(e);
            },
            o
          );
          var l = function (t, n, i) {
            var o = e.get("levels");
            if (o && "Auto" === o[0].label && n && n.items.length) {
              var a = n.items[0].el.querySelector(".jw-auto-label"),
                r = o[t.index] || { label: "" };
              a.textContent = i ? "" : r.label;
            }
          };
          e.on("change:visualQuality", function (t, n) {
            var i = o.children.quality;
            n && i && l(n.level, i, e.get("currentLevel"));
          }),
            e.on(
              "change:currentLevel",
              function (t, n) {
                var i = o.children.quality,
                  a = e.get("visualQuality");
                a && i && l(a.level, i, n);
              },
              o
            ),
            e.change("captionsList", function (n, r) {
              var l = { defaultText: i.off },
                s = e.get("captionsIndex");
              a(
                "captions",
                r,
                function (e) {
                  return t.setCurrentCaptions(e);
                },
                s,
                l
              );
              var c = o.children.captions;
              if (c && !c.children.captionsSettings) {
                c.topbar = c.topbar || c.createTopbar();
                var u = new rn("captionsSettings", c, i);
                u.title = "Subtitle Settings";
                var w = new Ge("Settings", u.open);
                c.topbar.appendChild(w.el);
                var p = new Je("Reset", function () {
                  e.set("captions", sn.a), h();
                });
                p.el.classList.add("jw-settings-reset");
                var j = e.get("captions"),
                  h = function () {
                    var t = [];
                    dn.forEach(function (n) {
                      j &&
                        j[n.propertyName] &&
                        (n.defaultVal = n.getOption(j[n.propertyName]));
                      var o = new rn(n.name, u, i),
                        a = new Ge(
                          { label: n.name, value: n.defaultVal },
                          o.open,
                          Re
                        ),
                        r = o.createItems(
                          n.options,
                          function (t) {
                            var i = a.el.querySelector(
                              ".jw-settings-content-item-value"
                            );
                            !(function (t, n) {
                              var i = e.get("captions"),
                                o = t.propertyName,
                                a = t.options && t.options[n],
                                r = t.getTypedValue(a),
                                l = Object(d.g)({}, i);
                              (l[o] = r), e.set("captions", l);
                            })(n, t),
                              (i.innerText = n.options[t]);
                          },
                          null
                        );
                      o.setMenuItems(r, n.options.indexOf(n.defaultVal) || 0),
                        t.push(a);
                    }),
                      t.push(p),
                      u.setMenuItems(t);
                  };
                h();
              }
            });
          var s = function (t, e) {
            t && e > -1 && t.items[e].activate();
          };
          e.change(
            "captionsIndex",
            function (t, e) {
              var i = o.children.captions;
              i && s(i, e), n.toggleCaptionsButtonState(!!e);
            },
            o
          );
          var c = function (n) {
            if (
              e.get("supportsPlaybackRate") &&
              "LIVE" !== e.get("streamType") &&
              e.get("playbackRateControls")
            ) {
              var r = n.indexOf(e.get("playbackRate")),
                l = { tooltipText: i.playbackRates };
              a(
                "playbackRates",
                n,
                function (e) {
                  return t.setPlaybackRate(e);
                },
                r,
                l
              );
            } else o.children.playbackRates && o.removeMenu("playbackRates");
          };
          e.on(
            "change:playbackRates",
            function (t, e) {
              c(e);
            },
            o
          );
          var u = function (n) {
            a(
              "audioTracks",
              n,
              function (e) {
                return t.setCurrentAudioTrack(e);
              },
              e.get("currentAudioTrack")
            );
          };
          return (
            e.on(
              "change:audioTracks",
              function (t, e) {
                u(e);
              },
              o
            ),
            e.on(
              "change:playbackRate",
              function (t, n) {
                var i = e.get("playbackRates"),
                  a = -1;
                i && (a = i.indexOf(n)), s(o.children.playbackRates, a);
              },
              o
            ),
            e.on(
              "change:currentAudioTrack",
              function (t, e) {
                o.children.audioTracks.items[e].activate();
              },
              o
            ),
            e.on(
              "change:playlistItem",
              function () {
                o.removeMenu("captions"),
                  n.elements.captionsButton.hide(),
                  o.visible && o.close();
              },
              o
            ),
            e.on("change:playbackRateControls", function () {
              c(e.get("playbackRates"));
            }),
            e.on(
              "change:castActive",
              function (t, n, i) {
                n !== i &&
                  (n
                    ? (o.removeMenu("audioTracks"),
                      o.removeMenu("quality"),
                      o.removeMenu("playbackRates"))
                    : (u(e.get("audioTracks")),
                      r(e.get("levels")),
                      c(e.get("playbackRates"))));
              },
              o
            ),
            e.on(
              "change:streamType",
              function () {
                c(e.get("playbackRates"));
              },
              o
            ),
            o
          );
        },
        fn = n(58),
        gn = n(35),
        bn = n(12),
        yn = function (t, e, n, i) {
          var o = Object(s.e)(
              '<div class="jw-reset jw-info-overlay jw-modal"><div class="jw-reset jw-info-container"><div class="jw-reset-text jw-info-title" dir="auto"></div><div class="jw-reset-text jw-info-duration" dir="auto"></div><div class="jw-reset-text jw-info-description" dir="auto"></div></div><div class="jw-reset jw-info-clientid"></div></div>'
            ),
            r = !1,
            l = null,
            c = !1,
            u = function (t) {
              /jw-info/.test(t.target.className) || d.close();
            },
            w = function () {
              var i,
                a,
                l,
                c,
                u,
                w = p(
                  "jw-info-close",
                  function () {
                    d.close();
                  },
                  e.get("localization").close,
                  [wt("close")]
                );
              w.show(),
                Object(s.m)(o, w.element()),
                (a = o.querySelector(".jw-info-title")),
                (l = o.querySelector(".jw-info-duration")),
                (c = o.querySelector(".jw-info-description")),
                (u = o.querySelector(".jw-info-clientid")),
                e.change("playlistItem", function (t, e) {
                  var n = e.description,
                    i = e.title;
                  Object(s.q)(c, n || ""), Object(s.q)(a, i || "Unknown Title");
                }),
                e.change(
                  "duration",
                  function (t, n) {
                    var i = "";
                    switch (e.get("streamType")) {
                      case "LIVE":
                        i = "Live";
                        break;
                      case "DVR":
                        i = "DVR";
                        break;
                      default:
                        n && (i = Object(vt.timeFormat)(n));
                    }
                    l.textContent = i;
                  },
                  d
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
                          } catch (t) {
                            return "none";
                          }
                        })()
                      )),
                t.appendChild(o),
                (r = !0);
            };
          var d = {
            open: function () {
              r || w(), document.addEventListener("click", u), (c = !0);
              var t = e.get("state");
              t === a.pb && n.pause("infoOverlayInteraction"), (l = t), i(!0);
            },
            close: function () {
              document.removeEventListener("click", u),
                (c = !1),
                e.get("state") === a.ob &&
                  l === a.pb &&
                  n.play("infoOverlayInteraction"),
                (l = null),
                i(!1);
            },
            destroy: function () {
              this.close(), e.off(null, null, this);
            },
          };
          return (
            Object.defineProperties(d, {
              visible: {
                enumerable: !0,
                get: function () {
                  return c;
                },
              },
            }),
            d
          );
        };
      var vn = function (t, e, n) {
          var i,
            o = !1,
            r = null,
            l = n.get("localization").shortcuts,
            c = Object(s.e)(
              (function (t, e) {
                var n = t
                  .map(function (t) {
                    return (
                      '<div class="jw-shortcuts-row jw-reset">' +
                      '<span class="jw-shortcuts-description jw-reset">'.concat(
                        t.description,
                        "</span>"
                      ) +
                      '<span class="jw-shortcuts-key jw-reset">'.concat(
                        t.key,
                        "</span>"
                      ) +
                      "</div>"
                    );
                  })
                  .join("");
                return (
                  '<div class="jw-shortcuts-tooltip jw-modal jw-reset" title="'.concat(
                    e,
                    '">'
                  ) +
                  '<span class="jw-hidden" id="jw-shortcuts-tooltip-explanation">Press shift question mark to access a list of keyboard shortcuts</span><div class="jw-reset jw-shortcuts-container"><div class="jw-reset jw-shortcuts-header">' +
                  '<span class="jw-reset jw-shortcuts-title">'.concat(
                    e,
                    "</span>"
                  ) +
                  '<button role="switch" class="jw-reset jw-switch" data-jw-switch-enabled="Enabled" data-jw-switch-disabled="Disabled"><span class="jw-reset jw-switch-knob"></span></button></div><div class="jw-reset jw-shortcuts-tooltip-list"><div class="jw-shortcuts-tooltip-descriptions jw-reset">' +
                  "".concat(n) +
                  "</div></div></div></div>"
                );
              })(
                (function (t) {
                  var e = t.playPause,
                    n = t.volumeToggle,
                    i = t.fullscreenToggle,
                    o = t.seekPercent,
                    a = t.increaseVolume,
                    r = t.decreaseVolume,
                    l = t.seekForward,
                    s = t.seekBackward;
                  return [
                    { key: t.spacebar, description: e },
                    { key: "↑", description: a },
                    { key: "↓", description: r },
                    { key: "→", description: l },
                    { key: "←", description: s },
                    { key: "c", description: t.captionsToggle },
                    { key: "f", description: i },
                    { key: "m", description: n },
                    { key: "0-9", description: o },
                  ];
                })(l),
                l.keyboardShortcuts
              )
            ),
            w = { reason: "settingsInteraction" },
            d = new u.a(c.querySelector(".jw-switch")),
            j = function () {
              d.el.setAttribute("aria-checked", n.get("enableShortcuts")),
                Object(s.a)(c, "jw-open"),
                (r = n.get("state")),
                c.querySelector(".jw-shortcuts-close").focus(),
                document.addEventListener("click", f),
                (o = !0),
                e.pause(w);
            },
            h = function () {
              Object(s.o)(c, "jw-open"),
                document.removeEventListener("click", f),
                t.focus(),
                (o = !1),
                r === a.pb && e.play(w);
            },
            f = function (t) {
              /jw-shortcuts|jw-switch/.test(t.target.className) || h();
            },
            g = function (t) {
              var e = t.currentTarget,
                i = "true" !== e.getAttribute("aria-checked");
              e.setAttribute("aria-checked", i), n.set("enableShortcuts", i);
            };
          return (
            (i = p("jw-shortcuts-close", h, n.get("localization").close, [
              wt("close"),
            ])),
            Object(s.m)(c, i.element()),
            i.show(),
            t.appendChild(c),
            d.on("click tap enter", g),
            {
              el: c,
              open: j,
              close: h,
              destroy: function () {
                h(), d.destroy();
              },
              toggleVisibility: function () {
                o ? h() : j();
              },
            }
          );
        },
        mn = function (t) {
          return (
            '<div class="jw-float-icon jw-icon jw-button-color jw-reset" aria-label='.concat(
              t,
              ' tabindex="0">'
            ) + "</div>"
          );
        };
      function xn(t) {
        return (xn =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (t) {
                return typeof t;
              }
            : function (t) {
                return t &&
                  "function" == typeof Symbol &&
                  t.constructor === Symbol &&
                  t !== Symbol.prototype
                  ? "symbol"
                  : typeof t;
              })(t);
      }
      function kn(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function On(t, e) {
        return !e || ("object" !== xn(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Cn(t) {
        return (Cn = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function Sn(t, e) {
        return (Sn =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var Tn = (function (t) {
        function e(t, n) {
          var i;
          return (
            (function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
            ((i = On(this, Cn(e).call(this))).element = Object(s.e)(mn(n))),
            i.element.appendChild(wt("close")),
            (i.ui = new u.a(i.element, { directSelect: !0 }).on(
              "click tap enter",
              function () {
                i.trigger(a.sb);
              }
            )),
            t.appendChild(i.element),
            i
          );
        }
        var n, i, o;
        return (
          (function (t, e) {
            if ("function" != typeof e && null !== e)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (t.prototype = Object.create(e && e.prototype, {
              constructor: { value: t, writable: !0, configurable: !0 },
            })),
              e && Sn(t, e);
          })(e, t),
          (n = e),
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
          ]) && kn(n.prototype, i),
          o && kn(n, o),
          e
        );
      })(r.a);
      function Mn(t) {
        return (Mn =
          "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
            ? function (t) {
                return typeof t;
              }
            : function (t) {
                return t &&
                  "function" == typeof Symbol &&
                  t.constructor === Symbol &&
                  t !== Symbol.prototype
                  ? "symbol"
                  : typeof t;
              })(t);
      }
      function zn(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function En(t, e) {
        return !e || ("object" !== Mn(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Ln(t) {
        return (Ln = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function Bn(t, e) {
        return (Bn =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      n.d(e, "default", function () {
        return Nn;
      }),
        n(95);
      var _n = o.OS.mobile ? 4e3 : 2e3,
        Vn = [27];
      (gn.a.cloneIcon = wt),
        bn.a.forEach(function (t) {
          if (t.getState() === a.lb) {
            var e = t.getContainer().querySelector(".jw-error-msg .jw-icon");
            e && !e.hasChildNodes() && e.appendChild(gn.a.cloneIcon("error"));
          }
        });
      var An = function () {
          return { reason: "interaction" };
        },
        Nn = (function (t) {
          function e(t, n) {
            var i;
            return (
              (function (t, e) {
                if (!(t instanceof e))
                  throw new TypeError("Cannot call a class as a function");
              })(this, e),
              ((i = En(this, Ln(e).call(this))).activeTimeout = -1),
              (i.inactiveTime = 0),
              (i.context = t),
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
                var t = i.inactiveTime - Object(c.a)();
                i.inactiveTime && t > 16
                  ? (i.activeTimeout = setTimeout(i.userInactiveTimeout, t))
                  : i.playerContainer.querySelector(".jw-tab-focus")
                  ? i.resetActiveTimeout()
                  : i.userInactive();
              }),
              i
            );
          }
          var n, i, r;
          return (
            (function (t, e) {
              if ("function" != typeof e && null !== e)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (t.prototype = Object.create(e && e.prototype, {
                constructor: { value: t, writable: !0, configurable: !0 },
              })),
                e && Bn(t, e);
            })(e, t),
            (n = e),
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
                value: function (t, e) {
                  var n = this,
                    i = this.context.createElement("div");
                  (i.className = "jw-controls jw-reset"), (this.div = i);
                  var r = this.context.createElement("div");
                  (r.className = "jw-controls-backdrop jw-reset"),
                    (this.backdrop = r),
                    (this.logo = this.playerContainer.querySelector(
                      ".jw-logo"
                    ));
                  var c = e.get("touchMode"),
                    u = function () {
                      (e.get("isFloating")
                        ? n.wrapperElement
                        : n.playerContainer
                      ).focus();
                    };
                  if (!this.displayContainer) {
                    var w = new Ce(e, t);
                    w.buttons.display.on("click tap enter", function () {
                      n.trigger(a.p),
                        n.userActive(1e3),
                        t.playToggle(An()),
                        u();
                    }),
                      this.div.appendChild(w.element()),
                      (this.displayContainer = w);
                  }
                  (this.infoOverlay = new yn(i, e, t, function (t) {
                    Object(s.v)(n.div, "jw-info-open", t),
                      t && n.div.querySelector(".jw-info-close").focus();
                  })),
                    o.OS.mobile ||
                      (this.shortcutsTooltip = new vn(
                        this.wrapperElement,
                        t,
                        e
                      )),
                    (this.rightClickMenu = new Pe(
                      this.infoOverlay,
                      this.shortcutsTooltip
                    )),
                    c
                      ? (Object(s.a)(this.playerContainer, "jw-flag-touch"),
                        this.rightClickMenu.setup(
                          e,
                          this.playerContainer,
                          this.wrapperElement
                        ))
                      : e.change(
                          "flashBlocked",
                          function (t, e) {
                            e
                              ? n.rightClickMenu.destroy()
                              : n.rightClickMenu.setup(
                                  t,
                                  n.playerContainer,
                                  n.wrapperElement
                                );
                          },
                          this
                        );
                  var d = e.get("floating");
                  if (d) {
                    var j = new Tn(i, e.get("localization").close);
                    j.on(a.sb, function () {
                      return n.trigger("dismissFloating", { doNotForward: !0 });
                    }),
                      !1 !== d.dismissible &&
                        Object(s.a)(
                          this.playerContainer,
                          "jw-floating-dismissible"
                        );
                  }
                  var h = (this.controlbar = new we(
                    t,
                    e,
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
                      function (t) {
                        this.trigger("nextShown", t);
                      },
                      this
                    ),
                    h.on("adjustVolume", x, this),
                    e.get("nextUpDisplay") && !h.nextUpToolTip)
                  ) {
                    var f = new ze(e, t, this.playerContainer);
                    f.on("all", this.trigger, this),
                      f.setup(this.context),
                      (h.nextUpToolTip = f),
                      this.div.appendChild(f.element());
                  }
                  this.div.appendChild(h.element());
                  var g = e.get("localization"),
                    b = (this.settingsMenu = hn(
                      t,
                      e.player,
                      this.controlbar,
                      g
                    )),
                    y = null;
                  this.controlbar.on("menuVisibility", function (i) {
                    var o = i.visible,
                      r = i.evt,
                      l = e.get("state"),
                      s = { reason: "settingsInteraction" },
                      c = n.controlbar.elements.settingsButton,
                      w = "keydown" === ((r && r.sourceEvent) || r || {}).type,
                      p = o || w ? 0 : _n;
                    n.userActive(p),
                      (y = l),
                      Object(fn.a)(e.get("containerWidth")) < 2 &&
                        (o && l === a.pb
                          ? t.pause(s)
                          : o || l !== a.ob || y !== a.pb || t.play(s)),
                      !o && w && c ? c.element().focus() : r && u();
                  }),
                    b.on("menuVisibility", function (t) {
                      return n.controlbar.trigger("menuVisibility", t);
                    }),
                    this.controlbar.on(
                      "settingsInteraction",
                      function (t, e, n) {
                        if (e) return b.defaultChild.toggle(n);
                        b.children[t].toggle(n);
                      }
                    ),
                    o.OS.mobile
                      ? this.div.appendChild(b.el)
                      : (this.playerContainer.setAttribute(
                          "aria-describedby",
                          "jw-shortcuts-tooltip-explanation"
                        ),
                        this.div.insertBefore(b.el, h.element()));
                  var v = function (e) {
                    if (e.get("autostartMuted")) {
                      var i = function () {
                          return n.unmuteAutoplay(t, e);
                        },
                        a = function (t, e) {
                          e || i();
                        };
                      o.OS.mobile &&
                        ((n.mute = p(
                          "jw-autostart-mute jw-off",
                          i,
                          e.get("localization").unmute,
                          [wt("volume-0")]
                        )),
                        n.mute.show(),
                        n.div.appendChild(n.mute.element())),
                        h.renderVolume(!0, e.get("volume")),
                        Object(s.a)(n.playerContainer, "jw-flag-autostart"),
                        e.on("change:autostartFailed", i, n),
                        e.on("change:autostartMuted change:mute", a, n),
                        (n.muteChangeCallback = a),
                        (n.unmuteCallback = i);
                    }
                  };
                  function m(n) {
                    var i = 0,
                      o = e.get("duration"),
                      a = e.get("position");
                    if ("DVR" === e.get("streamType")) {
                      var r = e.get("dvrSeekLimit");
                      (i = o), (o = Math.max(a, -r));
                    }
                    var s = Object(l.a)(a + n, i, o);
                    t.seek(s, An());
                  }
                  function x(n) {
                    var i = Object(l.a)(e.get("volume") + n, 0, 100);
                    t.setVolume(i);
                  }
                  e.once("change:autostartMuted", v), v(e);
                  var k = function (i) {
                    if (i.ctrlKey || i.metaKey) return !0;
                    var o = !n.settingsMenu.visible,
                      a = !0 === e.get("enableShortcuts"),
                      r = n.instreamState;
                    if (a || -1 !== Vn.indexOf(i.keyCode)) {
                      switch (i.keyCode) {
                        case 27:
                          if (e.get("fullscreen"))
                            t.setFullscreen(!1),
                              n.playerContainer.blur(),
                              n.userInactive();
                          else {
                            var l = t.getPlugin("related");
                            l && l.close({ type: "escape" });
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
                          t.playToggle(An());
                          break;
                        case 37:
                          !r && o && m(-5);
                          break;
                        case 39:
                          !r && o && m(5);
                          break;
                        case 38:
                          o && x(10);
                          break;
                        case 40:
                          o && x(-10);
                          break;
                        case 67:
                          var s = t.getCaptionsList().length;
                          if (s) {
                            var c = (t.getCurrentCaptions() + 1) % s;
                            t.setCurrentCaptions(c);
                          }
                          break;
                        case 77:
                          t.setMute();
                          break;
                        case 70:
                          t.setFullscreen();
                          break;
                        case 191:
                          n.shortcutsTooltip &&
                            n.shortcutsTooltip.toggleVisibility();
                          break;
                        default:
                          if (i.keyCode >= 48 && i.keyCode <= 59) {
                            var u = ((i.keyCode - 48) / 10) * e.get("duration");
                            t.seek(u, An());
                          }
                      }
                      return /13|32|37|38|39|40/.test(i.keyCode)
                        ? (i.preventDefault(), !1)
                        : void 0;
                    }
                  };
                  this.playerContainer.addEventListener("keydown", k),
                    (this.keydownCallback = k);
                  var O = function (t) {
                    switch (t.keyCode) {
                      case 9:
                        var e = n.playerContainer.contains(t.target) ? 0 : _n;
                        n.userActive(e);
                        break;
                      case 32:
                        t.preventDefault();
                    }
                  };
                  this.playerContainer.addEventListener("keyup", O),
                    (this.keyupCallback = O);
                  var C = function (t) {
                    var e = t.relatedTarget || document.querySelector(":focus");
                    e && (n.playerContainer.contains(e) || n.userInactive());
                  };
                  this.playerContainer.addEventListener("blur", C, !0),
                    (this.blurCallback = C);
                  var S = function t() {
                    "jw-shortcuts-tooltip-explanation" ===
                      n.playerContainer.getAttribute("aria-describedby") &&
                      n.playerContainer.removeAttribute("aria-describedby"),
                      n.playerContainer.removeEventListener("blur", t, !0);
                  };
                  this.shortcutsTooltip &&
                    (this.playerContainer.addEventListener("blur", S, !0),
                    (this.onRemoveShortcutsDescription = S)),
                    this.userActive(),
                    this.addControls(),
                    this.addBackdrop(),
                    e.set("controlsEnabled", !0);
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
                value: function (t) {
                  var e = this.nextUpToolTip,
                    n = this.settingsMenu,
                    i = this.infoOverlay,
                    o = this.controlbar,
                    a = this.rightClickMenu,
                    r = this.shortcutsTooltip,
                    l = this.playerContainer,
                    c = this.div;
                  clearTimeout(this.activeTimeout),
                    (this.activeTimeout = -1),
                    this.off(),
                    t.off(null, null, this),
                    t.set("controlsEnabled", !1),
                    c.parentNode &&
                      (Object(s.o)(l, "jw-flag-touch"),
                      c.parentNode.removeChild(c)),
                    o && o.destroy(),
                    a && a.destroy(),
                    this.keydownCallback &&
                      l.removeEventListener("keydown", this.keydownCallback),
                    this.keyupCallback &&
                      l.removeEventListener("keyup", this.keyupCallback),
                    this.blurCallback &&
                      l.removeEventListener("blur", this.blurCallback),
                    this.onRemoveShortcutsDescription &&
                      l.removeEventListener(
                        "blur",
                        this.onRemoveShortcutsDescription
                      ),
                    this.displayContainer && this.displayContainer.destroy(),
                    e && e.destroy(),
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
                value: function (t, e) {
                  var n = !e.get("autostartFailed"),
                    i = e.get("mute");
                  n ? (i = !1) : e.set("playOnViewable", !1),
                    this.muteChangeCallback &&
                      (e.off(
                        "change:autostartMuted change:mute",
                        this.muteChangeCallback
                      ),
                      (this.muteChangeCallback = null)),
                    this.unmuteCallback &&
                      (e.off("change:autostartFailed", this.unmuteCallback),
                      (this.unmuteCallback = null)),
                    e.set("autostartFailed", void 0),
                    e.set("autostartMuted", void 0),
                    t.setMute(i),
                    this.controlbar.renderVolume(i, e.get("volume")),
                    this.mute && this.mute.hide(),
                    Object(s.o)(this.playerContainer, "jw-flag-autostart"),
                    this.userActive();
                },
              },
              {
                key: "mouseMove",
                value: function (t) {
                  var e = this.controlbar.element().contains(t.target),
                    n =
                      this.controlbar.nextUpToolTip &&
                      this.controlbar.nextUpToolTip
                        .element()
                        .contains(t.target),
                    i = this.logo && this.logo.contains(t.target),
                    o = e || n || i ? 0 : _n;
                  this.userActive(o);
                },
              },
              {
                key: "userActive",
                value: function () {
                  var t =
                    arguments.length > 0 && void 0 !== arguments[0]
                      ? arguments[0]
                      : _n;
                  t > 0
                    ? ((this.inactiveTime = Object(c.a)() + t),
                      -1 === this.activeTimeout &&
                        (this.activeTimeout = setTimeout(
                          this.userInactiveTimeout,
                          t
                        )))
                    : this.resetActiveTimeout(),
                    this.showing ||
                      (Object(s.o)(
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
                      Object(s.a)(
                        this.playerContainer,
                        "jw-flag-user-inactive"
                      ),
                      this.trigger("userInactive"));
                },
              },
              {
                key: "addBackdrop",
                value: function () {
                  var t = this.instreamState
                    ? this.div
                    : this.wrapperElement.querySelector(".jw-captions");
                  this.wrapperElement.insertBefore(this.backdrop, t);
                },
              },
              {
                key: "removeBackdrop",
                value: function () {
                  var t = this.backdrop.parentNode;
                  t && t.removeChild(this.backdrop);
                },
              },
              {
                key: "setupInstream",
                value: function () {
                  (this.instreamState = !0),
                    this.userActive(),
                    this.addBackdrop(),
                    this.settingsMenu && this.settingsMenu.close(),
                    Object(s.o)(this.playerContainer, "jw-flag-autostart"),
                    this.controlbar.elements.time
                      .element()
                      .setAttribute("tabindex", "-1");
                },
              },
              {
                key: "destroyInstream",
                value: function (t) {
                  (this.instreamState = null),
                    this.addBackdrop(),
                    t.get("autostartMuted") &&
                      Object(s.a)(this.playerContainer, "jw-flag-autostart"),
                    this.controlbar.elements.time
                      .element()
                      .setAttribute("tabindex", "0");
                },
              },
            ]) && zn(n.prototype, i),
            r && zn(n, r),
            e
          );
        })(r.a);
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
    ,
    function (t, e, n) {
      "use strict";
      n.d(e, "a", function () {
        return o;
      });
      var i = n(2);
      function o(t) {
        var e = [],
          n = (t = Object(i.i)(t)).split("\r\n\r\n");
        1 === n.length && (n = t.split("\n\n"));
        for (var o = 0; o < n.length; o++)
          if ("WEBVTT" !== n[o]) {
            var r = a(n[o]);
            r.text && e.push(r);
          }
        return e;
      }
      function a(t) {
        var e = {},
          n = t.split("\r\n");
        1 === n.length && (n = t.split("\n"));
        var o = 1;
        if (
          (n[0].indexOf(" --\x3e ") > 0 && (o = 0),
          n.length > o + 1 && n[o + 1])
        ) {
          var a = n[o],
            r = a.indexOf(" --\x3e ");
          r > 0 &&
            ((e.begin = Object(i.g)(a.substr(0, r))),
            (e.end = Object(i.g)(a.substr(r + 5))),
            (e.text = n.slice(o + 1).join("\r\n")));
        }
        return e;
      }
    },
    function (t, e, n) {
      "use strict";
      n.d(e, "a", function () {
        return o;
      }),
        n.d(e, "b", function () {
          return a;
        });
      var i = n(5);
      function o(t) {
        var e = -1;
        return (
          t >= 1280
            ? (e = 7)
            : t >= 960
            ? (e = 6)
            : t >= 800
            ? (e = 5)
            : t >= 640
            ? (e = 4)
            : t >= 540
            ? (e = 3)
            : t >= 420
            ? (e = 2)
            : t >= 320
            ? (e = 1)
            : t >= 250 && (e = 0),
          e
        );
      }
      function a(t, e) {
        var n = "jw-breakpoint-" + e;
        Object(i.p)(t, /jw-breakpoint--?\d+/, n);
      }
    },
    function (t, e, n) {
      "use strict";
      n.d(e, "a", function () {
        return w;
      });
      var i,
        o = n(0),
        a = n(8),
        r = n(16),
        l = n(7),
        s = n(3),
        c = n(10),
        u = n(5),
        w = {
          back: !0,
          backgroundOpacity: 50,
          edgeStyle: null,
          fontSize: 14,
          fontOpacity: 100,
          fontScale: 0.05,
          preprocessor: o.k,
          windowOpacity: 0,
        },
        p = function (t) {
          var e,
            l,
            p,
            d,
            j,
            h,
            f,
            g,
            b,
            y = this,
            v = t.player;
          function m() {
            Object(o.o)(e.fontSize) &&
              (v.get("containerHeight")
                ? (g =
                    (w.fontScale * (e.userFontScale || 1) * e.fontSize) /
                    w.fontSize)
                : v.once("change:containerHeight", m, this));
          }
          function x() {
            var t = v.get("containerHeight");
            if (t) {
              var e;
              if (v.get("fullscreen") && a.OS.iOS) e = null;
              else {
                var n = t * g;
                e =
                  Math.round(
                    10 *
                      (function (t) {
                        var e = v.get("mediaElement");
                        if (e && e.videoHeight) {
                          var n = e.videoWidth,
                            i = e.videoHeight,
                            o = n / i,
                            r = v.get("containerHeight"),
                            l = v.get("containerWidth");
                          if (v.get("fullscreen") && a.OS.mobile) {
                            var s = window.screen;
                            s.orientation &&
                              ((r = s.availHeight), (l = s.availWidth));
                          }
                          if (l && r && n && i)
                            return (l / r > o ? r : (i * l) / n) * g;
                        }
                        return t;
                      })(n)
                  ) / 10;
              }
              v.get("renderCaptionsNatively")
                ? (function (t, e) {
                    var n = "#".concat(
                      t,
                      " .jw-video::-webkit-media-text-track-display"
                    );
                    e &&
                      ((e += "px"),
                      a.OS.iOS &&
                        Object(c.b)(n, { fontSize: "inherit" }, t, !0));
                    (b.fontSize = e), Object(c.b)(n, b, t, !0);
                  })(v.get("id"), e)
                : Object(c.d)(j, { fontSize: e });
            }
          }
          function k(t, e, n) {
            var i = Object(c.c)("#000000", n);
            "dropshadow" === t
              ? (e.textShadow = "0 2px 1px " + i)
              : "raised" === t
              ? (e.textShadow =
                  "0 0 5px " + i + ", 0 1px 5px " + i + ", 0 2px 5px " + i)
              : "depressed" === t
              ? (e.textShadow = "0 -2px 1px " + i)
              : "uniform" === t &&
                (e.textShadow =
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
          ((j = document.createElement("div")).className =
            "jw-captions jw-reset"),
            (this.show = function () {
              Object(u.a)(j, "jw-captions-enabled");
            }),
            (this.hide = function () {
              Object(u.o)(j, "jw-captions-enabled");
            }),
            (this.populate = function (t) {
              v.get("renderCaptionsNatively") ||
                ((p = []),
                (l = t),
                t ? this.selectCues(t, d) : this.renderCues());
            }),
            (this.resize = function () {
              x(), this.renderCues(!0);
            }),
            (this.renderCues = function (t) {
              (t = !!t), i && i.processCues(window, p, j, t);
            }),
            (this.selectCues = function (t, e) {
              if (t && t.data && e && !v.get("renderCaptionsNatively")) {
                var n = this.getAlignmentPosition(t, e);
                !1 !== n &&
                  ((p = this.getCurrentCues(t.data, n)), this.renderCues(!0));
              }
            }),
            (this.getCurrentCues = function (t, e) {
              return Object(o.h)(t, function (t) {
                return e >= t.startTime && (!t.endTime || e <= t.endTime);
              });
            }),
            (this.getAlignmentPosition = function (t, e) {
              var n = t.source,
                i = e.metadata,
                a = e.currentTime;
              return n && i && Object(o.r)(i[n]) && (a = i[n]), a;
            }),
            (this.clear = function () {
              Object(u.g)(j);
            }),
            (this.setup = function (t, n) {
              (h = document.createElement("div")),
                (f = document.createElement("span")),
                (h.className = "jw-captions-window jw-reset"),
                (f.className = "jw-captions-text jw-reset"),
                (e = Object(o.g)({}, w, n)),
                (g = w.fontScale);
              var i = function () {
                if (!v.get("renderCaptionsNatively")) {
                  m(e.fontSize);
                  var n = e.windowColor,
                    i = e.windowOpacity,
                    o = e.edgeStyle;
                  b = {};
                  var r = {};
                  !(function (t, e) {
                    var n = e.color,
                      i = e.fontOpacity;
                    (n || i !== w.fontOpacity) &&
                      (t.color = Object(c.c)(n || "#ffffff", i));
                    if (e.back) {
                      var o = e.backgroundColor,
                        a = e.backgroundOpacity;
                      (o === w.backgroundColor && a === w.backgroundOpacity) ||
                        (t.backgroundColor = Object(c.c)(o, a));
                    } else t.background = "transparent";
                    e.fontFamily && (t.fontFamily = e.fontFamily);
                    e.fontStyle && (t.fontStyle = e.fontStyle);
                    e.fontWeight && (t.fontWeight = e.fontWeight);
                    e.textDecoration && (t.textDecoration = e.textDecoration);
                  })(r, e),
                    (n || i !== w.windowOpacity) &&
                      (b.backgroundColor = Object(c.c)(n || "#000000", i)),
                    k(o, r, e.fontOpacity),
                    e.back || null !== o || k("uniform", r),
                    Object(c.d)(h, b),
                    Object(c.d)(f, r),
                    (function (t, e) {
                      x(),
                        (function (t, e) {
                          a.Browser.safari &&
                            Object(c.b)(
                              "#" +
                                t +
                                " .jw-video::-webkit-media-text-track-display-backdrop",
                              { backgroundColor: e.backgroundColor },
                              t,
                              !0
                            );
                          Object(c.b)(
                            "#" +
                              t +
                              " .jw-video::-webkit-media-text-track-display",
                            b,
                            t,
                            !0
                          ),
                            Object(c.b)("#" + t + " .jw-video::cue", e, t, !0);
                        })(t, e),
                        (function (t, e) {
                          Object(c.b)(
                            "#" + t + " .jw-text-track-display",
                            b,
                            t
                          ),
                            Object(c.b)("#" + t + " .jw-text-track-cue", e, t);
                        })(t, e);
                    })(t, r);
                }
              };
              i(),
                h.appendChild(f),
                j.appendChild(h),
                v.change(
                  "captionsTrack",
                  function (t, e) {
                    this.populate(e);
                  },
                  this
                ),
                v.set("captions", e),
                v.on("change:captions", function (t, n) {
                  (e = n), i();
                });
            }),
            (this.element = function () {
              return j;
            }),
            (this.destroy = function () {
              v.off(null, null, this), this.off();
            });
          var O = function (t) {
            (d = t), y.selectCues(l, d);
          };
          v.on(
            "change:playlistItem",
            function () {
              (d = null), (p = []);
            },
            this
          ),
            v.on(
              s.Q,
              function (t) {
                (p = []), O(t);
              },
              this
            ),
            v.on(s.S, O, this),
            v.on(
              "subtitlesTrackData",
              function () {
                this.selectCues(l, d);
              },
              this
            ),
            v.on(
              "change:captionsList",
              function t(e, o) {
                var a = this;
                1 !== o.length &&
                  (e.get("renderCaptionsNatively") ||
                    i ||
                    (n
                      .e(8)
                      .then(
                        function (t) {
                          i = n(68).default;
                        }.bind(null, n)
                      )
                      .catch(Object(r.c)(301121))
                      .catch(function (t) {
                        a.trigger(s.tb, t);
                      }),
                    e.off("change:captionsList", t, this)));
              },
              this
            );
        };
      Object(o.g)(p.prototype, l.a), (e.b = p);
    },
    function (t, e, n) {
      "use strict";
      t.exports = function (t) {
        var e = [];
        return (
          (e.toString = function () {
            return this.map(function (e) {
              var n = (function (t, e) {
                var n = t[1] || "",
                  i = t[3];
                if (!i) return n;
                if (e && "function" == typeof btoa) {
                  var o =
                      ((r = i),
                      "/*# sourceMappingURL=data:application/json;charset=utf-8;base64," +
                        btoa(unescape(encodeURIComponent(JSON.stringify(r)))) +
                        " */"),
                    a = i.sources.map(function (t) {
                      return "/*# sourceURL=" + i.sourceRoot + t + " */";
                    });
                  return [n].concat(a).concat([o]).join("\n");
                }
                var r;
                return [n].join("\n");
              })(e, t);
              return e[2] ? "@media " + e[2] + "{" + n + "}" : n;
            }).join("");
          }),
          (e.i = function (t, n) {
            "string" == typeof t && (t = [[null, t, ""]]);
            for (var i = {}, o = 0; o < this.length; o++) {
              var a = this[o][0];
              null != a && (i[a] = !0);
            }
            for (o = 0; o < t.length; o++) {
              var r = t[o];
              (null != r[0] && i[r[0]]) ||
                (n && !r[2]
                  ? (r[2] = n)
                  : n && (r[2] = "(" + r[2] + ") and (" + n + ")"),
                e.push(r));
            }
          }),
          e
        );
      };
    },
    function (t, e) {
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
      function l(t) {
        var e = document.createElement("style");
        return (
          (e.type = "text/css"),
          e.setAttribute("data-jwplayer-id", t),
          (function (t) {
            r().appendChild(t);
          })(e),
          e
        );
      }
      function s(t, e) {
        var n,
          i,
          o,
          r = a[t];
        r || (r = a[t] = { element: l(t), counter: 0 });
        var s = r.counter++;
        return (
          (n = r.element),
          (o = function () {
            w(n, s, "");
          }),
          (i = function (t) {
            w(n, s, t);
          })(e.css),
          function (t) {
            if (t) {
              if (t.css === e.css && t.media === e.media) return;
              i((e = t).css);
            } else o();
          }
        );
      }
      t.exports = {
        style: function (t, e) {
          !(function (t, e) {
            for (var n = 0; n < e.length; n++) {
              var i = e[n],
                a = (o[t] || {})[i.id];
              if (a) {
                for (var r = 0; r < a.parts.length; r++) a.parts[r](i.parts[r]);
                for (; r < i.parts.length; r++) a.parts.push(s(t, i.parts[r]));
              } else {
                var l = [];
                for (r = 0; r < i.parts.length; r++) l.push(s(t, i.parts[r]));
                (o[t] = o[t] || {}), (o[t][i.id] = { id: i.id, parts: l });
              }
            }
          })(
            e,
            (function (t) {
              for (var e = [], n = {}, i = 0; i < t.length; i++) {
                var o = t[i],
                  a = o[0],
                  r = o[1],
                  l = o[2],
                  s = { css: r, media: l };
                n[a]
                  ? n[a].parts.push(s)
                  : e.push((n[a] = { id: a, parts: [s] }));
              }
              return e;
            })(t)
          );
        },
        clear: function (t, e) {
          var n = o[t];
          if (!n) return;
          if (e) {
            var i = n[e];
            if (i) for (var a = 0; a < i.parts.length; a += 1) i.parts[a]();
            return;
          }
          for (var r = Object.keys(n), l = 0; l < r.length; l += 1)
            for (var s = n[r[l]], c = 0; c < s.parts.length; c += 1)
              s.parts[c]();
          delete o[t];
        },
      };
      var c,
        u =
          ((c = []),
          function (t, e) {
            return (c[t] = e), c.filter(Boolean).join("\n");
          });
      function w(t, e, n) {
        if (t.styleSheet) t.styleSheet.cssText = u(e, n);
        else {
          var i = document.createTextNode(n),
            o = t.childNodes[e];
          o ? t.replaceChild(i, o) : t.appendChild(i);
        }
      }
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-arrow-right" viewBox="0 0 240 240" focusable="false"><path d="M183.6,104.4L81.8,0L45.4,36.3l84.9,84.9l-84.9,84.9L79.3,240l101.9-101.7c9.9-6.9,12.4-20.4,5.5-30.4C185.8,106.7,184.8,105.4,183.6,104.4L183.6,104.4z"></path></svg>';
    },
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    ,
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-buffer" viewBox="0 0 240 240" focusable="false"><path d="M120,186.667a66.667,66.667,0,0,1,0-133.333V40a80,80,0,1,0,80,80H186.667A66.846,66.846,0,0,1,120,186.667Z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-replay" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M120,41.9v-20c0-5-4-8-8-4l-44,28a5.865,5.865,0,0,0-3.3,7.6A5.943,5.943,0,0,0,68,56.8l43,29c5,4,9,1,9-4v-20a60,60,0,1,1-60,60H40a80,80,0,1,0,80-79.9Z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-error" viewBox="0 0 36 36" style="width:100%;height:100%;" focusable="false"><path d="M34.6 20.2L10 33.2 27.6 16l7 3.7a.4.4 0 0 1 .2.5.4.4 0 0 1-.2.2zM33.3 0L21 12.2 9 6c-.2-.3-.6 0-.6.5V25L0 33.6 2.5 36 36 2.7z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-play" viewBox="0 0 240 240" focusable="false"><path d="M62.8,199.5c-1,0.8-2.4,0.6-3.3-0.4c-0.4-0.5-0.6-1.1-0.5-1.8V42.6c-0.2-1.3,0.7-2.4,1.9-2.6c0.7-0.1,1.3,0.1,1.9,0.4l154.7,77.7c2.1,1.1,2.1,2.8,0,3.8L62.8,199.5z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-pause" viewBox="0 0 240 240" focusable="false"><path d="M100,194.9c0.2,2.6-1.8,4.8-4.4,5c-0.2,0-0.4,0-0.6,0H65c-2.6,0.2-4.8-1.8-5-4.4c0-0.2,0-0.4,0-0.6V45c-0.2-2.6,1.8-4.8,4.4-5c0.2,0,0.4,0,0.6,0h30c2.6-0.2,4.8,1.8,5,4.4c0,0.2,0,0.4,0,0.6V194.9z M180,45.1c0.2-2.6-1.8-4.8-4.4-5c-0.2,0-0.4,0-0.6,0h-30c-2.6-0.2-4.8,1.8-5,4.4c0,0.2,0,0.4,0,0.6V195c-0.2,2.6,1.8,4.8,4.4,5c0.2,0,0.4,0,0.6,0h30c2.6,0.2,4.8-1.8,5-4.4c0-0.2,0-0.4,0-0.6V45.1z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-rewind" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M113.2,131.078a21.589,21.589,0,0,0-17.7-10.6,21.589,21.589,0,0,0-17.7,10.6,44.769,44.769,0,0,0,0,46.3,21.589,21.589,0,0,0,17.7,10.6,21.589,21.589,0,0,0,17.7-10.6,44.769,44.769,0,0,0,0-46.3Zm-17.7,47.2c-7.8,0-14.4-11-14.4-24.1s6.6-24.1,14.4-24.1,14.4,11,14.4,24.1S103.4,178.278,95.5,178.278Zm-43.4,9.7v-51l-4.8,4.8-6.8-6.8,13-13a4.8,4.8,0,0,1,8.2,3.4v62.7l-9.6-.1Zm162-130.2v125.3a4.867,4.867,0,0,1-4.8,4.8H146.6v-19.3h48.2v-96.4H79.1v19.3c0,5.3-3.6,7.2-8,4.3l-41.8-27.9a6.013,6.013,0,0,1-2.7-8,5.887,5.887,0,0,1,2.7-2.7l41.8-27.9c4.4-2.9,8-1,8,4.3v19.3H209.2A4.974,4.974,0,0,1,214.1,57.778Z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-next" viewBox="0 0 240 240" focusable="false"><path d="M165,60v53.3L59.2,42.8C56.9,41.3,55,42.3,55,45v150c0,2.7,1.9,3.8,4.2,2.2L165,126.6v53.3h20v-120L165,60L165,60z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-stop" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M190,185c0.2,2.6-1.8,4.8-4.4,5c-0.2,0-0.4,0-0.6,0H55c-2.6,0.2-4.8-1.8-5-4.4c0-0.2,0-0.4,0-0.6V55c-0.2-2.6,1.8-4.8,4.4-5c0.2,0,0.4,0,0.6,0h130c2.6-0.2,4.8,1.8,5,4.4c0,0.2,0,0.4,0,0.6V185z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-volume-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M116.4,42.8v154.5c0,2.8-1.7,3.6-3.8,1.7l-54.1-48.1H28.9c-2.8,0-5.2-2.3-5.2-5.2V94.2c0-2.8,2.3-5.2,5.2-5.2h29.6l54.1-48.1C114.6,39.1,116.4,39.9,116.4,42.8z M212.3,96.4l-14.6-14.6l-23.6,23.6l-23.6-23.6l-14.6,14.6l23.6,23.6l-23.6,23.6l14.6,14.6l23.6-23.6l23.6,23.6l14.6-14.6L188.7,120L212.3,96.4z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-volume-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M116.4,42.8v154.5c0,2.8-1.7,3.6-3.8,1.7l-54.1-48.1H28.9c-2.8,0-5.2-2.3-5.2-5.2V94.2c0-2.8,2.3-5.2,5.2-5.2h29.6l54.1-48.1C114.7,39.1,116.4,39.9,116.4,42.8z M178.2,120c0-22.7-18.5-41.2-41.2-41.2v20.6c11.4,0,20.6,9.2,20.6,20.6c0,11.4-9.2,20.6-20.6,20.6v20.6C159.8,161.2,178.2,142.7,178.2,120z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-volume-100" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M116.5,42.8v154.4c0,2.8-1.7,3.6-3.8,1.7l-54.1-48H29c-2.8,0-5.2-2.3-5.2-5.2V94.3c0-2.8,2.3-5.2,5.2-5.2h29.6l54.1-48C114.8,39.2,116.5,39.9,116.5,42.8z"></path><path d="M136.2,160v-20c11.1,0,20-8.9,20-20s-8.9-20-20-20V80c22.1,0,40,17.9,40,40S158.3,160,136.2,160z"></path><path d="M216.2,120c0-44.2-35.8-80-80-80v20c33.1,0,60,26.9,60,60s-26.9,60-60,60v20C180.4,199.9,216.1,164.1,216.2,120z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-cc-on" viewBox="0 0 240 240" focusable="false"><path d="M215,40H25c-2.7,0-5,2.2-5,5v150c0,2.7,2.2,5,5,5h190c2.7,0,5-2.2,5-5V45C220,42.2,217.8,40,215,40z M108.1,137.7c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9c-2.4-3.7-6.5-5.9-10.9-5.9c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6C90.4,141.7,102,143.5,108.1,137.7z M152.9,137.7c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9c-2.4-3.7-6.5-5.9-10.9-5.9c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6C135.2,141.7,146.8,143.5,152.9,137.7z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-cc-off" viewBox="0 0 240 240" focusable="false"><path d="M99.4,97.8c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6c0,9.6,11.6,11.4,17.7,5.5c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9C107.9,100,103.8,97.8,99.4,97.8z M144.1,97.8c-2.4-0.2-4.8,0.7-6.6,2.3c-1.7,1.7-2.5,4.1-2.4,6.5v25.6c0,9.6,11.6,11.4,17.7,5.5c0.7-0.7,1.5-1.5,2.4-2.3l6.6,7.8c-2.2,2.4-5,4.4-8,5.8c-8,3.5-17.3,2.4-24.3-2.9c-3.9-3.6-5.9-8.7-5.5-14v-25.6c0-2.7,0.5-5.3,1.5-7.8c0.9-2.2,2.4-4.3,4.2-5.9c5.7-4.5,13.2-6.2,20.3-4.6c3.3,0.5,6.3,2,8.7,4.3c1.3,1.3,2.5,2.6,3.5,4.2l-7.1,6.9C152.6,100,148.5,97.8,144.1,97.8L144.1,97.8z M200,60v120H40V60H200 M215,40H25c-2.7,0-5,2.2-5,5v150c0,2.7,2.2,5,5,5h190c2.7,0,5-2.2,5-5V45C220,42.2,217.8,40,215,40z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-airplay-on" viewBox="0 0 240 240" focusable="false"><path d="M229.9,40v130c0.2,2.6-1.8,4.8-4.4,5c-0.2,0-0.4,0-0.6,0h-44l-17-20h46V55H30v100h47l-17,20h-45c-2.6,0.2-4.8-1.8-5-4.4c0-0.2,0-0.4,0-0.6V40c-0.2-2.6,1.8-4.8,4.4-5c0.2,0,0.4,0,0.6,0h209.8c2.6-0.2,4.8,1.8,5,4.4C229.9,39.7,229.9,39.9,229.9,40z M104.9,122l15-18l15,18l11,13h44V75H50v60h44L104.9,122z M179.9,205l-60-70l-60,70H179.9z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-airplay-off" viewBox="0 0 240 240" focusable="false"><path d="M210,55v100h-50l20,20h45c2.6,0.2,4.8-1.8,5-4.4c0-0.2,0-0.4,0-0.6V40c0.2-2.6-1.8-4.8-4.4-5c-0.2,0-0.4,0-0.6,0H15c-2.6-0.2-4.8,1.8-5,4.4c0,0.2,0,0.4,0,0.6v130c-0.2,2.6,1.8,4.8,4.4,5c0.2,0,0.4,0,0.6,0h45l20-20H30V55H210 M60,205l60-70l60,70H60L60,205z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-arrow-left" viewBox="0 0 240 240" focusable="false"><path d="M55.4,104.4c-1.1,1.1-2.2,2.3-3.1,3.6c-6.9,9.9-4.4,23.5,5.5,30.4L159.7,240l33.9-33.9l-84.9-84.9l84.9-84.9L157.3,0L55.4,104.4L55.4,104.4z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-playback-rate" viewBox="0 0 240 240" focusable="false"><path d="M158.83,48.83A71.17,71.17,0,1,0,230,120,71.163,71.163,0,0,0,158.83,48.83Zm45.293,77.632H152.34V74.708h12.952v38.83h38.83ZM35.878,74.708h38.83V87.66H35.878ZM10,113.538H61.755V126.49H10Zm25.878,38.83h38.83V165.32H35.878Z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-settings" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M204,145l-25-14c0.8-3.6,1.2-7.3,1-11c0.2-3.7-0.2-7.4-1-11l25-14c2.2-1.6,3.1-4.5,2-7l-16-26c-1.2-2.1-3.8-2.9-6-2l-25,14c-6-4.2-12.3-7.9-19-11V35c0.2-2.6-1.8-4.8-4.4-5c-0.2,0-0.4,0-0.6,0h-30c-2.6-0.2-4.8,1.8-5,4.4c0,0.2,0,0.4,0,0.6v28c-6.7,3.1-13,6.7-19,11L56,60c-2.2-0.9-4.8-0.1-6,2L35,88c-1.6,2.2-1.3,5.3,0.9,6.9c0,0,0.1,0,0.1,0.1l25,14c-0.8,3.6-1.2,7.3-1,11c-0.2,3.7,0.2,7.4,1,11l-25,14c-2.2,1.6-3.1,4.5-2,7l16,26c1.2,2.1,3.8,2.9,6,2l25-14c5.7,4.6,12.2,8.3,19,11v28c-0.2,2.6,1.8,4.8,4.4,5c0.2,0,0.4,0,0.6,0h30c2.6,0.2,4.8-1.8,5-4.4c0-0.2,0-0.4,0-0.6v-28c7-2.3,13.5-6,19-11l25,14c2.5,1.3,5.6,0.4,7-2l15-26C206.7,149.4,206,146.7,204,145z M120,149.9c-16.5,0-30-13.4-30-30s13.4-30,30-30s30,13.4,30,30c0.3,16.3-12.6,29.7-28.9,30C120.7,149.9,120.4,149.9,120,149.9z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-audio-tracks" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M35,34h160v20H35V34z M35,94h160V74H35V94z M35,134h60v-20H35V134z M160,114c-23.4-1.3-43.6,16.5-45,40v50h20c5.2,0.3,9.7-3.6,10-8.9c0-0.4,0-0.7,0-1.1v-20c0.3-5.2-3.6-9.7-8.9-10c-0.4,0-0.7,0-1.1,0h-10v-10c1.5-17.9,17.1-31.3,35-30c17.9-1.3,33.6,12.1,35,30v10H185c-5.2-0.3-9.7,3.6-10,8.9c0,0.4,0,0.7,0,1.1v20c-0.3,5.2,3.6,9.7,8.9,10c0.4,0,0.7,0,1.1,0h20v-50C203.5,130.6,183.4,112.7,160,114z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg class="jw-svg-icon jw-svg-icon-quality-100" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240" focusable="false"><path d="M55,200H35c-3,0-5-2-5-4c0,0,0,0,0-1v-30c0-3,2-5,4-5c0,0,0,0,1,0h20c3,0,5,2,5,4c0,0,0,0,0,1v30C60,198,58,200,55,200L55,200z M110,195v-70c0-3-2-5-4-5c0,0,0,0-1,0H85c-3,0-5,2-5,4c0,0,0,0,0,1v70c0,3,2,5,4,5c0,0,0,0,1,0h20C108,200,110,198,110,195L110,195z M160,195V85c0-3-2-5-4-5c0,0,0,0-1,0h-20c-3,0-5,2-5,4c0,0,0,0,0,1v110c0,3,2,5,4,5c0,0,0,0,1,0h20C158,200,160,198,160,195L160,195z M210,195V45c0-3-2-5-4-5c0,0,0,0-1,0h-20c-3,0-5,2-5,4c0,0,0,0,0,1v150c0,3,2,5,4,5c0,0,0,0,1,0h20C208,200,210,198,210,195L210,195z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-fullscreen-off" viewBox="0 0 240 240" focusable="false"><path d="M109.2,134.9l-8.4,50.1c-0.4,2.7-2.4,3.3-4.4,1.4L82,172l-27.9,27.9l-14.2-14.2l27.9-27.9l-14.4-14.4c-1.9-1.9-1.3-3.9,1.4-4.4l50.1-8.4c1.8-0.5,3.6,0.6,4.1,2.4C109.4,133.7,109.4,134.3,109.2,134.9L109.2,134.9z M172.1,82.1L200,54.2L185.8,40l-27.9,27.9l-14.4-14.4c-1.9-1.9-3.9-1.3-4.4,1.4l-8.4,50.1c-0.5,1.8,0.6,3.6,2.4,4.1c0.5,0.2,1.2,0.2,1.7,0l50.1-8.4c2.7-0.4,3.3-2.4,1.4-4.4L172.1,82.1z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-fullscreen-on" viewBox="0 0 240 240" focusable="false"><path d="M96.3,186.1c1.9,1.9,1.3,4-1.4,4.4l-50.6,8.4c-1.8,0.5-3.7-0.6-4.2-2.4c-0.2-0.6-0.2-1.2,0-1.7l8.4-50.6c0.4-2.7,2.4-3.4,4.4-1.4l14.5,14.5l28.2-28.2l14.3,14.3l-28.2,28.2L96.3,186.1z M195.8,39.1l-50.6,8.4c-2.7,0.4-3.4,2.4-1.4,4.4l14.5,14.5l-28.2,28.2l14.3,14.3l28.2-28.2l14.5,14.5c1.9,1.9,4,1.3,4.4-1.4l8.4-50.6c0.5-1.8-0.6-3.6-2.4-4.2C197,39,196.4,39,195.8,39.1L195.8,39.1z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-close" viewBox="0 0 240 240" focusable="false"><path d="M134.8,120l48.6-48.6c2-1.9,2.1-5.2,0.2-7.2c0,0-0.1-0.1-0.2-0.2l-7.4-7.4c-1.9-2-5.2-2.1-7.2-0.2c0,0-0.1,0.1-0.2,0.2L120,105.2L71.4,56.6c-1.9-2-5.2-2.1-7.2-0.2c0,0-0.1,0.1-0.2,0.2L56.6,64c-2,1.9-2.1,5.2-0.2,7.2c0,0,0.1,0.1,0.2,0.2l48.6,48.7l-48.6,48.6c-2,1.9-2.1,5.2-0.2,7.2c0,0,0.1,0.1,0.2,0.2l7.4,7.4c1.9,2,5.2,2.1,7.2,0.2c0,0,0.1-0.1,0.2-0.2l48.7-48.6l48.6,48.6c1.9,2,5.2,2.1,7.2,0.2c0,0,0.1-0.1,0.2-0.2l7.4-7.4c2-1.9,2.1-5.2,0.2-7.2c0,0-0.1-0.1-0.2-0.2L134.8,120z"></path></svg>';
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-jwplayer-logo" viewBox="0 0 992 1024" focusable="false"><path d="M144 518.4c0 6.4-6.4 6.4-6.4 0l-3.2-12.8c0 0-6.4-19.2-12.8-38.4 0-6.4-6.4-12.8-9.6-22.4-6.4-6.4-16-9.6-28.8-6.4-9.6 3.2-16 12.8-16 22.4s0 16 0 25.6c3.2 25.6 22.4 121.6 32 140.8 9.6 22.4 35.2 32 54.4 22.4 22.4-9.6 28.8-35.2 38.4-54.4 9.6-25.6 60.8-166.4 60.8-166.4 6.4-12.8 9.6-12.8 9.6 0 0 0 0 140.8-3.2 204.8 0 25.6 0 67.2 9.6 89.6 6.4 16 12.8 28.8 25.6 38.4s28.8 12.8 44.8 12.8c6.4 0 16-3.2 22.4-6.4 9.6-6.4 16-12.8 25.6-22.4 16-19.2 28.8-44.8 38.4-64 25.6-51.2 89.6-201.6 89.6-201.6 6.4-12.8 9.6-12.8 9.6 0 0 0-9.6 256-9.6 355.2 0 25.6 6.4 48 12.8 70.4 9.6 22.4 22.4 38.4 44.8 48s48 9.6 70.4-3.2c16-9.6 28.8-25.6 38.4-38.4 12.8-22.4 25.6-48 32-70.4 19.2-51.2 35.2-102.4 51.2-153.6s153.6-540.8 163.2-582.4c0-6.4 0-9.6 0-12.8 0-9.6-6.4-19.2-16-22.4-16-6.4-32 0-38.4 12.8-6.4 16-195.2 470.4-195.2 470.4-6.4 12.8-9.6 12.8-9.6 0 0 0 0-156.8 0-288 0-70.4-35.2-108.8-83.2-118.4-22.4-3.2-44.8 0-67.2 12.8s-35.2 32-48 54.4c-16 28.8-105.6 297.6-105.6 297.6-6.4 12.8-9.6 12.8-9.6 0 0 0-3.2-115.2-6.4-144-3.2-41.6-12.8-108.8-67.2-115.2-51.2-3.2-73.6 57.6-86.4 99.2-9.6 25.6-51.2 163.2-51.2 163.2v3.2z"></path></svg>';
    },
    function (t, e, n) {
      var i = n(96);
      "string" == typeof i && (i = [["all-players", i, ""]]),
        n(61).style(i, "all-players"),
        i.locals && (t.exports = i.locals);
    },
    function (t, e, n) {
      (t.exports = n(60)(!1)).push([
        t.i,
        '.jw-overlays,.jw-controls,.jw-controls-backdrop,.jw-flag-small-player .jw-settings-menu,.jw-settings-submenu{height:100%;width:100%}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-settings-menu .jw-icon.jw-button-color::after{position:absolute;right:0}.jw-overlays,.jw-controls,.jw-controls-backdrop,.jw-settings-item-active::before{top:0;position:absolute;left:0}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-settings-menu .jw-icon.jw-button-color::after{position:absolute;bottom:0;left:0}.jw-nextup-close{position:absolute;top:0;right:0}.jw-overlays,.jw-controls,.jw-flag-small-player .jw-settings-menu{position:absolute;bottom:0;right:0}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-time-tip::after,.jw-settings-menu .jw-icon.jw-button-color::after,.jw-text-live::before,.jw-controlbar .jw-tooltip::after,.jw-settings-menu .jw-tooltip::after{content:"";display:block}.jw-svg-icon{height:24px;width:24px;fill:currentColor;pointer-events:none}.jw-icon{height:44px;width:44px;background-color:transparent;outline:none}.jw-icon.jw-tab-focus:focus{border:solid 2px #4d90fe}.jw-icon-airplay .jw-svg-icon-airplay-off{display:none}.jw-off.jw-icon-airplay .jw-svg-icon-airplay-off{display:block}.jw-icon-airplay .jw-svg-icon-airplay-on{display:block}.jw-off.jw-icon-airplay .jw-svg-icon-airplay-on{display:none}.jw-icon-cc .jw-svg-icon-cc-off{display:none}.jw-off.jw-icon-cc .jw-svg-icon-cc-off{display:block}.jw-icon-cc .jw-svg-icon-cc-on{display:block}.jw-off.jw-icon-cc .jw-svg-icon-cc-on{display:none}.jw-icon-fullscreen .jw-svg-icon-fullscreen-off{display:none}.jw-off.jw-icon-fullscreen .jw-svg-icon-fullscreen-off{display:block}.jw-icon-fullscreen .jw-svg-icon-fullscreen-on{display:block}.jw-off.jw-icon-fullscreen .jw-svg-icon-fullscreen-on{display:none}.jw-icon-volume .jw-svg-icon-volume-0{display:none}.jw-off.jw-icon-volume .jw-svg-icon-volume-0{display:block}.jw-icon-volume .jw-svg-icon-volume-100{display:none}.jw-full.jw-icon-volume .jw-svg-icon-volume-100{display:block}.jw-icon-volume .jw-svg-icon-volume-50{display:block}.jw-off.jw-icon-volume .jw-svg-icon-volume-50,.jw-full.jw-icon-volume .jw-svg-icon-volume-50{display:none}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after{height:100%;width:24px;box-shadow:inset 0 -3px 0 -1px currentColor;margin:auto;opacity:0;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-settings-menu .jw-icon[aria-checked="true"]::after,.jw-settings-open .jw-icon-settings::after,.jw-icon-volume.jw-open::after{opacity:1}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-cc,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-settings,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-audio-tracks,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-hd,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-settings-sharing,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-fullscreen,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player).jw-flag-cast-available .jw-icon-airplay,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player).jw-flag-cast-available .jw-icon-cast{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-volume,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-text-live{bottom:6px}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-volume::after{display:none}.jw-overlays,.jw-controls{pointer-events:none}.jw-controls-backdrop{display:block;background:linear-gradient(to bottom, transparent, rgba(0,0,0,0.4) 77%, rgba(0,0,0,0.4) 100%) 100% 100% / 100% 240px no-repeat transparent;transition:opacity 250ms cubic-bezier(0, .25, .25, 1),background-size 250ms cubic-bezier(0, .25, .25, 1);pointer-events:none}.jw-overlays{cursor:auto}.jw-controls{overflow:hidden}.jw-flag-small-player .jw-controls{text-align:center}.jw-text{height:1em;font-family:Arial,Helvetica,sans-serif;font-size:.75em;font-style:normal;font-weight:normal;color:#fff;text-align:center;font-variant:normal;font-stretch:normal}.jw-controlbar,.jw-skip,.jw-display-icon-container .jw-icon,.jw-nextup-container,.jw-autostart-mute,.jw-overlays .jw-plugin{pointer-events:all}.jwplayer .jw-display-icon-container,.jw-error .jw-display-icon-container{width:auto;height:auto;box-sizing:content-box}.jw-display{display:table;height:100%;padding:57px 0;position:relative;width:100%}.jw-flag-dragging .jw-display{display:none}.jw-state-idle:not(.jw-flag-cast-available) .jw-display{padding:0}.jw-display-container{display:table-cell;height:100%;text-align:center;vertical-align:middle}.jw-display-controls{display:inline-block}.jwplayer .jw-display-icon-container{float:left}.jw-display-icon-container{display:inline-block;padding:5.5px;margin:0 22px}.jw-display-icon-container .jw-icon{height:75px;width:75px;cursor:pointer;display:flex;justify-content:center;align-items:center}.jw-display-icon-container .jw-icon .jw-svg-icon{height:33px;width:33px;padding:0;position:relative}.jw-display-icon-container .jw-icon .jw-svg-icon-rewind{padding:.2em .05em}.jw-breakpoint--1 .jw-nextup-container{display:none}.jw-breakpoint-0 .jw-display-icon-next,.jw-breakpoint--1 .jw-display-icon-next,.jw-breakpoint-0 .jw-display-icon-rewind,.jw-breakpoint--1 .jw-display-icon-rewind{display:none}.jw-breakpoint-0 .jw-display .jw-icon,.jw-breakpoint--1 .jw-display .jw-icon,.jw-breakpoint-0 .jw-display .jw-svg-icon,.jw-breakpoint--1 .jw-display .jw-svg-icon{width:44px;height:44px;line-height:44px}.jw-breakpoint-0 .jw-display .jw-icon:before,.jw-breakpoint--1 .jw-display .jw-icon:before,.jw-breakpoint-0 .jw-display .jw-svg-icon:before,.jw-breakpoint--1 .jw-display .jw-svg-icon:before{width:22px;height:22px}.jw-breakpoint-1 .jw-display .jw-icon,.jw-breakpoint-1 .jw-display .jw-svg-icon{width:44px;height:44px;line-height:44px}.jw-breakpoint-1 .jw-display .jw-icon:before,.jw-breakpoint-1 .jw-display .jw-svg-icon:before{width:22px;height:22px}.jw-breakpoint-1 .jw-display .jw-icon.jw-icon-rewind:before{width:33px;height:33px}.jw-breakpoint-2 .jw-display .jw-icon,.jw-breakpoint-3 .jw-display .jw-icon,.jw-breakpoint-2 .jw-display .jw-svg-icon,.jw-breakpoint-3 .jw-display .jw-svg-icon{width:77px;height:77px;line-height:77px}.jw-breakpoint-2 .jw-display .jw-icon:before,.jw-breakpoint-3 .jw-display .jw-icon:before,.jw-breakpoint-2 .jw-display .jw-svg-icon:before,.jw-breakpoint-3 .jw-display .jw-svg-icon:before{width:38.5px;height:38.5px}.jw-breakpoint-4 .jw-display .jw-icon,.jw-breakpoint-5 .jw-display .jw-icon,.jw-breakpoint-6 .jw-display .jw-icon,.jw-breakpoint-7 .jw-display .jw-icon,.jw-breakpoint-4 .jw-display .jw-svg-icon,.jw-breakpoint-5 .jw-display .jw-svg-icon,.jw-breakpoint-6 .jw-display .jw-svg-icon,.jw-breakpoint-7 .jw-display .jw-svg-icon{width:88px;height:88px;line-height:88px}.jw-breakpoint-4 .jw-display .jw-icon:before,.jw-breakpoint-5 .jw-display .jw-icon:before,.jw-breakpoint-6 .jw-display .jw-icon:before,.jw-breakpoint-7 .jw-display .jw-icon:before,.jw-breakpoint-4 .jw-display .jw-svg-icon:before,.jw-breakpoint-5 .jw-display .jw-svg-icon:before,.jw-breakpoint-6 .jw-display .jw-svg-icon:before,.jw-breakpoint-7 .jw-display .jw-svg-icon:before{width:44px;height:44px}.jw-controlbar{display:flex;flex-flow:row wrap;align-items:center;justify-content:center;position:absolute;left:0;bottom:0;width:100%;border:none;border-radius:0;background-size:auto;box-shadow:none;max-height:72px;transition:250ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility;transition-delay:0s}.jw-breakpoint-7 .jw-controlbar{max-height:140px}.jw-breakpoint-7 .jw-controlbar .jw-button-container{padding:0 48px 20px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-tooltip{margin-bottom:-7px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume .jw-overlay{padding-bottom:40%}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-text{font-size:1em}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-text.jw-text-elapsed{justify-content:flex-end}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-inline,.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume{height:60px;width:60px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-inline .jw-svg-icon,.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume .jw-svg-icon{height:30px;width:30px}.jw-breakpoint-7 .jw-controlbar .jw-slider-time{padding:0 60px;height:34px}.jw-breakpoint-7 .jw-controlbar .jw-slider-time .jw-slider-container{height:10px}.jw-controlbar .jw-button-image{background:no-repeat 50% 50%;background-size:contain;max-height:24px}.jw-controlbar .jw-spacer{flex:1 1 auto;align-self:stretch}.jw-controlbar .jw-icon.jw-button-color:hover{color:#fff}.jw-button-container{display:flex;flex-flow:row nowrap;flex:1 1 auto;align-items:center;justify-content:center;width:100%;padding:0 12px}.jw-slider-horizontal{background-color:transparent}.jw-icon-inline{position:relative}.jw-icon-inline,.jw-icon-tooltip{height:44px;width:44px;align-items:center;display:flex;justify-content:center}.jw-icon-inline:not(.jw-text),.jw-icon-tooltip,.jw-slider-horizontal{cursor:pointer}.jw-text-elapsed,.jw-text-duration{justify-content:flex-start;width:-webkit-fit-content;width:-moz-fit-content;width:fit-content}.jw-icon-tooltip{position:relative}.jw-knob:hover,.jw-icon-inline:hover,.jw-icon-tooltip:hover,.jw-icon-display:hover,.jw-option:before:hover{color:#fff}.jw-time-tip,.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{pointer-events:none}.jw-icon-cast{display:none;margin:0;padding:0}.jw-icon-cast google-cast-launcher{background-color:transparent;border:none;padding:0;width:24px;height:24px;cursor:pointer}.jw-icon-inline.jw-icon-volume{display:none}.jwplayer .jw-text-countdown{display:none}.jw-flag-small-player .jw-display{padding-top:0;padding-bottom:0}.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-rewind,.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-next,.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-playback{display:none}.jw-flag-ads-vpaid:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar,.jw-flag-user-inactive.jw-state-playing:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar,.jw-flag-user-inactive.jw-state-buffering:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar{visibility:hidden;pointer-events:none;opacity:0;transition-delay:0s, 250ms}.jw-flag-ads-vpaid:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop,.jw-flag-user-inactive.jw-state-playing:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop,.jw-flag-user-inactive.jw-state-buffering:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop{opacity:0}.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-countdown{display:flex}.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint--1 .jw-text-elapsed,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-elapsed,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint--1 .jw-text-duration,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-duration{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-text-countdown,.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-related-btn,.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-slider-volume{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-controlbar{flex-direction:column-reverse}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-button-container{height:30px}.jw-breakpoint--1.jw-flag-ads:not(.jw-flag-audio-player) .jw-icon-volume,.jw-breakpoint--1.jw-flag-ads:not(.jw-flag-audio-player) .jw-icon-fullscreen{display:none}.jwplayer:not(.jw-breakpoint-0) .jw-text-duration:before,.jwplayer:not(.jw-breakpoint--1) .jw-text-duration:before{content:"/";padding-right:1ch;padding-left:1ch}.jwplayer:not(.jw-flag-user-inactive) .jw-controlbar{will-change:transform}.jwplayer:not(.jw-flag-user-inactive) .jw-controlbar .jw-text{-webkit-transform-style:preserve-3d;transform-style:preserve-3d}.jw-slider-container{display:flex;align-items:center;position:relative;touch-action:none}.jw-rail,.jw-buffer,.jw-progress{position:absolute;cursor:pointer}.jw-progress{background-color:#f2f2f2}.jw-rail{background-color:rgba(255,255,255,0.3)}.jw-buffer{background-color:rgba(255,255,255,0.3)}.jw-knob{height:13px;width:13px;background-color:#fff;border-radius:50%;box-shadow:0 0 10px rgba(0,0,0,0.4);opacity:1;pointer-events:none;position:absolute;-webkit-transform:translate(-50%, -50%) scale(0);transform:translate(-50%, -50%) scale(0);transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, -webkit-transform;transition-property:opacity, transform;transition-property:opacity, transform, -webkit-transform}.jw-flag-dragging .jw-slider-time .jw-knob,.jw-icon-volume:active .jw-slider-volume .jw-knob{box-shadow:0 0 26px rgba(0,0,0,0.2),0 0 10px rgba(0,0,0,0.4),0 0 0 6px rgba(255,255,255,0.2)}.jw-slider-horizontal,.jw-slider-vertical{display:flex}.jw-slider-horizontal .jw-slider-container{height:5px;width:100%}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress,.jw-slider-horizontal .jw-cue,.jw-slider-horizontal .jw-knob{top:50%}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress,.jw-slider-horizontal .jw-cue{-webkit-transform:translate(0, -50%);transform:translate(0, -50%)}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress{height:5px}.jw-slider-horizontal .jw-rail{width:100%}.jw-slider-vertical{align-items:center;flex-direction:column}.jw-slider-vertical .jw-slider-container{height:88px;width:5px}.jw-slider-vertical .jw-rail,.jw-slider-vertical .jw-buffer,.jw-slider-vertical .jw-progress,.jw-slider-vertical .jw-knob{left:50%}.jw-slider-vertical .jw-rail,.jw-slider-vertical .jw-buffer,.jw-slider-vertical .jw-progress{height:100%;width:5px;-webkit-backface-visibility:hidden;backface-visibility:hidden;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);transition:-webkit-transform 150ms ease-in-out;transition:transform 150ms ease-in-out;transition:transform 150ms ease-in-out, -webkit-transform 150ms ease-in-out;bottom:0}.jw-slider-vertical .jw-knob{-webkit-transform:translate(-50%, 50%);transform:translate(-50%, 50%)}.jw-slider-time.jw-tab-focus:focus .jw-rail{outline:solid 2px #4d90fe}.jw-slider-time,.jw-flag-audio-player .jw-slider-volume{height:17px;width:100%;align-items:center;background:transparent none;padding:0 12px}.jw-slider-time .jw-cue{background-color:rgba(33,33,33,0.8);cursor:pointer;position:absolute;width:6px}.jw-slider-time,.jw-horizontal-volume-container{z-index:1;outline:none}.jw-slider-time .jw-rail,.jw-horizontal-volume-container .jw-rail,.jw-slider-time .jw-buffer,.jw-horizontal-volume-container .jw-buffer,.jw-slider-time .jw-progress,.jw-horizontal-volume-container .jw-progress,.jw-slider-time .jw-cue,.jw-horizontal-volume-container .jw-cue{-webkit-backface-visibility:hidden;backface-visibility:hidden;height:100%;-webkit-transform:translate(0, -50%) scale(1, .6);transform:translate(0, -50%) scale(1, .6);transition:-webkit-transform 150ms ease-in-out;transition:transform 150ms ease-in-out;transition:transform 150ms ease-in-out, -webkit-transform 150ms ease-in-out}.jw-slider-time:hover .jw-rail,.jw-horizontal-volume-container:hover .jw-rail,.jw-slider-time:focus .jw-rail,.jw-horizontal-volume-container:focus .jw-rail,.jw-flag-dragging .jw-slider-time .jw-rail,.jw-flag-dragging .jw-horizontal-volume-container .jw-rail,.jw-flag-touch .jw-slider-time .jw-rail,.jw-flag-touch .jw-horizontal-volume-container .jw-rail,.jw-slider-time:hover .jw-buffer,.jw-horizontal-volume-container:hover .jw-buffer,.jw-slider-time:focus .jw-buffer,.jw-horizontal-volume-container:focus .jw-buffer,.jw-flag-dragging .jw-slider-time .jw-buffer,.jw-flag-dragging .jw-horizontal-volume-container .jw-buffer,.jw-flag-touch .jw-slider-time .jw-buffer,.jw-flag-touch .jw-horizontal-volume-container .jw-buffer,.jw-slider-time:hover .jw-progress,.jw-horizontal-volume-container:hover .jw-progress,.jw-slider-time:focus .jw-progress,.jw-horizontal-volume-container:focus .jw-progress,.jw-flag-dragging .jw-slider-time .jw-progress,.jw-flag-dragging .jw-horizontal-volume-container .jw-progress,.jw-flag-touch .jw-slider-time .jw-progress,.jw-flag-touch .jw-horizontal-volume-container .jw-progress,.jw-slider-time:hover .jw-cue,.jw-horizontal-volume-container:hover .jw-cue,.jw-slider-time:focus .jw-cue,.jw-horizontal-volume-container:focus .jw-cue,.jw-flag-dragging .jw-slider-time .jw-cue,.jw-flag-dragging .jw-horizontal-volume-container .jw-cue,.jw-flag-touch .jw-slider-time .jw-cue,.jw-flag-touch .jw-horizontal-volume-container .jw-cue{-webkit-transform:translate(0, -50%) scale(1, 1);transform:translate(0, -50%) scale(1, 1)}.jw-slider-time:hover .jw-knob,.jw-horizontal-volume-container:hover .jw-knob,.jw-slider-time:focus .jw-knob,.jw-horizontal-volume-container:focus .jw-knob{-webkit-transform:translate(-50%, -50%) scale(1);transform:translate(-50%, -50%) scale(1)}.jw-slider-time .jw-rail,.jw-horizontal-volume-container .jw-rail{background-color:rgba(255,255,255,0.2)}.jw-slider-time .jw-buffer,.jw-horizontal-volume-container .jw-buffer{background-color:rgba(255,255,255,0.4)}.jw-flag-touch .jw-slider-time::before,.jw-flag-touch .jw-horizontal-volume-container::before{height:44px;width:100%;content:"";position:absolute;display:block;bottom:calc(100% - 17px);left:0}.jw-slider-time.jw-tab-focus:focus .jw-rail,.jw-horizontal-volume-container.jw-tab-focus:focus .jw-rail{outline:solid 2px #4d90fe}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time{height:17px;padding:0}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time .jw-slider-container{height:10px}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time .jw-knob{border-radius:0;border:1px solid rgba(0,0,0,0.75);height:12px;width:10px}.jw-modal{width:284px}.jw-breakpoint-7 .jw-modal,.jw-breakpoint-6 .jw-modal,.jw-breakpoint-5 .jw-modal{height:232px}.jw-breakpoint-4 .jw-modal,.jw-breakpoint-3 .jw-modal{height:192px}.jw-breakpoint-2 .jw-modal,.jw-flag-small-player .jw-modal{bottom:0;right:0;height:100%;width:100%;max-height:none;max-width:none;z-index:2}.jwplayer .jw-rightclick{display:none;position:absolute;white-space:nowrap}.jwplayer .jw-rightclick.jw-open{display:block}.jwplayer .jw-rightclick .jw-rightclick-list{border-radius:1px;list-style:none;margin:0;padding:0}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item{background-color:rgba(0,0,0,0.8);border-bottom:1px solid #444;margin:0}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-logo{color:#fff;display:inline-flex;padding:0 10px 0 0;vertical-align:middle}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-logo .jw-svg-icon{height:20px;width:20px}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-link{border:none;color:#fff;display:block;font-size:11px;line-height:1em;padding:15px 23px;text-align:start;text-decoration:none;width:100%}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item:last-child{border-bottom:none}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item:hover{cursor:pointer}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured{vertical-align:middle}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured .jw-rightclick-link{color:#fff}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured .jw-rightclick-link span{color:#fff}.jwplayer .jw-rightclick .jw-info-overlay-item,.jwplayer .jw-rightclick .jw-share-item,.jwplayer .jw-rightclick .jw-shortcuts-item{border:none;background-color:transparent;outline:none;cursor:pointer}.jw-icon-tooltip.jw-open .jw-overlay{opacity:1;pointer-events:auto;transition-delay:0s}.jw-icon-tooltip.jw-open .jw-overlay:focus{outline:none}.jw-icon-tooltip.jw-open .jw-overlay:focus.jw-tab-focus{outline:solid 2px #4d90fe}.jw-slider-time .jw-overlay:before{height:1em;top:auto}.jw-slider-time .jw-icon-tooltip.jw-open .jw-overlay{pointer-events:none}.jw-volume-tip{padding:13px 0 26px}.jw-time-tip,.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{height:auto;width:100%;box-shadow:0 0 10px rgba(0,0,0,0.4);color:#fff;display:block;margin:0 0 14px;pointer-events:none;position:relative;z-index:0}.jw-time-tip::after,.jw-controlbar .jw-tooltip::after,.jw-settings-menu .jw-tooltip::after{top:100%;position:absolute;left:50%;height:14px;width:14px;border-radius:1px;background-color:currentColor;-webkit-transform-origin:75% 50%;transform-origin:75% 50%;-webkit-transform:translate(-50%, -50%) rotate(45deg);transform:translate(-50%, -50%) rotate(45deg);z-index:-1}.jw-time-tip .jw-text,.jw-controlbar .jw-tooltip .jw-text,.jw-settings-menu .jw-tooltip .jw-text{background-color:#fff;border-radius:1px;color:#000;font-size:10px;height:auto;line-height:1;padding:7px 10px;display:inline-block;min-width:100%;vertical-align:middle}.jw-controlbar .jw-overlay{position:absolute;bottom:100%;left:50%;margin:0;min-height:44px;min-width:44px;opacity:0;pointer-events:none;transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility;transition-delay:0s, 150ms;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);width:100%;z-index:1}.jw-controlbar .jw-overlay .jw-contents{position:relative}.jw-controlbar .jw-option{position:relative;white-space:nowrap;cursor:pointer;list-style:none;height:1.5em;font-family:inherit;line-height:1.5em;padding:0 .5em;font-size:.8em;margin:0}.jw-controlbar .jw-option::before{padding-right:.125em}.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{position:absolute;bottom:100%;left:50%;opacity:0;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);transition:100ms 0s cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility, -webkit-transform;transition-property:opacity, transform, visibility;transition-property:opacity, transform, visibility, -webkit-transform;visibility:hidden;white-space:nowrap;width:auto;z-index:1}.jw-controlbar .jw-tooltip.jw-open,.jw-settings-menu .jw-tooltip.jw-open{opacity:1;-webkit-transform:translate(-50%, -10px);transform:translate(-50%, -10px);transition-duration:150ms;transition-delay:500ms,0s,500ms;visibility:visible}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen{left:auto;right:0;-webkit-transform:translate(0, 0);transform:translate(0, 0)}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen.jw-open,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen.jw-open{-webkit-transform:translate(0, -10px);transform:translate(0, -10px)}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen::after,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen::after{left:auto;right:9px}.jw-tooltip-time{height:auto;width:0;bottom:100%;line-height:normal;padding:0;pointer-events:none;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none}.jw-tooltip-time .jw-overlay{bottom:0;min-height:0;width:auto}.jw-tooltip{bottom:57px;display:none;position:absolute}.jw-tooltip .jw-text{height:100%;white-space:nowrap;text-overflow:ellipsis;direction:unset;max-width:246px;overflow:hidden}.jw-flag-audio-player .jw-tooltip{display:none}.jw-flag-small-player .jw-time-thumb{display:none}.jwplayer .jw-shortcuts-tooltip{top:50%;position:absolute;left:50%;background:#333;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);display:none;color:#fff;pointer-events:all;-webkit-user-select:text;-moz-user-select:text;-ms-user-select:text;user-select:text;overflow:hidden;flex-direction:column;z-index:1}.jwplayer .jw-shortcuts-tooltip.jw-open{display:flex}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-close{flex:0 0 auto;margin:5px 5px 5px auto}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container{display:flex;flex:1 1 auto;flex-flow:column;font-size:12px;margin:0 20px 20px;overflow-y:auto;padding:5px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container::-webkit-scrollbar{background-color:transparent;width:6px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-title{font-weight:bold}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-header{align-items:center;display:flex;justify-content:space-between;margin-bottom:10px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list{display:flex;max-width:340px;margin:0 10px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-tooltip-descriptions{width:100%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row{display:flex;align-items:center;justify-content:space-between;margin:10px 0;width:100%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row .jw-shortcuts-description{margin-right:10px;max-width:70%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row .jw-shortcuts-key{background:#fefefe;color:#333;overflow:hidden;padding:7px 10px;text-overflow:ellipsis;white-space:nowrap}.jw-skip{color:rgba(255,255,255,0.8);cursor:default;position:absolute;display:flex;right:.75em;bottom:56px;padding:.5em;border:1px solid #333;background-color:#000;align-items:center;height:2em}.jw-skip.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jw-skip.jw-skippable{cursor:pointer;padding:.25em .75em}.jw-skip.jw-skippable:hover{cursor:pointer;color:#fff}.jw-skip.jw-skippable .jw-skip-icon{display:inline;height:24px;width:24px;margin:0}.jw-breakpoint-7 .jw-skip{padding:1.35em 1em;bottom:130px}.jw-breakpoint-7 .jw-skip .jw-text{font-size:1em;font-weight:normal}.jw-breakpoint-7 .jw-skip .jw-icon-inline{height:30px;width:30px}.jw-breakpoint-7 .jw-skip .jw-icon-inline .jw-svg-icon{height:30px;width:30px}.jw-skip .jw-skip-icon{display:none;margin-left:-0.75em;padding:0 .5em;pointer-events:none}.jw-skip .jw-skip-icon .jw-svg-icon-next{display:block;padding:0}.jw-skip .jw-text,.jw-skip .jw-skip-icon{vertical-align:middle;font-size:.7em}.jw-skip .jw-text{font-weight:bold}.jw-cast{background-size:cover;display:none;height:100%;position:relative;width:100%}.jw-cast-container{background:linear-gradient(180deg, rgba(25,25,25,0.75), rgba(25,25,25,0.25), rgba(25,25,25,0));left:0;padding:20px 20px 80px;position:absolute;top:0;width:100%}.jw-cast-text{color:#fff;font-size:1.6em}.jw-breakpoint--1 .jw-cast-text,.jw-breakpoint-0 .jw-cast-text{font-size:1.15em}.jw-breakpoint-1 .jw-cast-text,.jw-breakpoint-2 .jw-cast-text,.jw-breakpoint-3 .jw-cast-text{font-size:1.3em}.jw-nextup-container{position:absolute;bottom:66px;left:0;background-color:transparent;cursor:pointer;margin:0 auto;padding:12px;pointer-events:none;right:0;text-align:right;visibility:hidden;width:100%}.jw-settings-open .jw-nextup-container,.jw-info-open .jw-nextup-container{display:none}.jw-breakpoint-7 .jw-nextup-container{padding:60px}.jw-flag-small-player .jw-nextup-container{padding:0 12px 0 0}.jw-flag-small-player .jw-nextup-container .jw-nextup-title,.jw-flag-small-player .jw-nextup-container .jw-nextup-duration,.jw-flag-small-player .jw-nextup-container .jw-nextup-close{display:none}.jw-flag-small-player .jw-nextup-container .jw-nextup-tooltip{height:30px}.jw-flag-small-player .jw-nextup-container .jw-nextup-header{font-size:12px}.jw-flag-small-player .jw-nextup-container .jw-nextup-body{justify-content:center;align-items:center;padding:.75em .3em}.jw-flag-small-player .jw-nextup-container .jw-nextup-thumbnail{width:50%}.jw-flag-small-player .jw-nextup-container .jw-nextup{max-width:65px}.jw-flag-small-player .jw-nextup-container .jw-nextup.jw-nextup-thumbnail-visible{max-width:120px}.jw-nextup{background:#333;border-radius:0;box-shadow:0 0 10px rgba(0,0,0,0.5);color:rgba(255,255,255,0.8);display:inline-block;max-width:280px;overflow:hidden;opacity:0;position:relative;width:64%;pointer-events:all;-webkit-transform:translate(0, -5px);transform:translate(0, -5px);transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, -webkit-transform;transition-property:opacity, transform;transition-property:opacity, transform, -webkit-transform;transition-delay:0s}.jw-nextup:hover .jw-nextup-tooltip{color:#fff}.jw-nextup.jw-nextup-thumbnail-visible{max-width:400px}.jw-nextup.jw-nextup-thumbnail-visible .jw-nextup-thumbnail{display:block}.jw-nextup-container-visible{visibility:visible}.jw-nextup-container-visible .jw-nextup{opacity:1;-webkit-transform:translate(0, 0);transform:translate(0, 0);transition-delay:0s, 0s, 150ms}.jw-nextup-tooltip{display:flex;height:80px}.jw-nextup-thumbnail{width:120px;background-position:center;background-size:cover;flex:0 0 auto;display:none}.jw-nextup-body{flex:1 1 auto;overflow:hidden;padding:.75em .875em;display:flex;flex-flow:column wrap;justify-content:space-between}.jw-nextup-header,.jw-nextup-title{font-size:14px;line-height:1.35}.jw-nextup-header{font-weight:bold}.jw-nextup-title{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;width:100%}.jw-nextup-duration{align-self:flex-end;text-align:right;font-size:12px}.jw-nextup-close{height:24px;width:24px;border:none;color:rgba(255,255,255,0.8);cursor:pointer;margin:6px;visibility:hidden}.jw-nextup-close:hover{color:#fff}.jw-nextup-sticky .jw-nextup-close{visibility:visible}.jw-autostart-mute{position:absolute;bottom:0;right:12px;height:44px;width:44px;background-color:rgba(33,33,33,0.4);padding:5px 4px 5px 6px;display:none}.jwplayer.jw-flag-autostart:not(.jw-flag-media-audio) .jw-nextup{display:none}.jw-settings-menu{position:absolute;bottom:57px;right:12px;align-items:flex-start;background-color:#333;display:none;flex-flow:column nowrap;max-width:284px;pointer-events:auto}.jw-settings-open .jw-settings-menu{display:flex}.jw-breakpoint-7 .jw-settings-menu{bottom:130px;right:60px;max-height:none;max-width:none;height:35%;width:25%}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline{height:60px;width:60px}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline .jw-svg-icon{height:30px;width:30px}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline .jw-tooltip .jw-text{font-size:1em}.jw-breakpoint-7 .jw-settings-menu .jw-settings-back{min-width:60px}.jw-breakpoint-6 .jw-settings-menu,.jw-breakpoint-5 .jw-settings-menu{height:232px;width:284px;max-height:232px}.jw-breakpoint-4 .jw-settings-menu,.jw-breakpoint-3 .jw-settings-menu{height:192px;width:284px;max-height:192px}.jw-breakpoint-2 .jw-settings-menu{height:179px;width:284px;max-height:179px}.jw-flag-small-player .jw-settings-menu{max-width:none}.jw-settings-menu .jw-icon.jw-button-color::after{height:100%;width:24px;box-shadow:inset 0 -3px 0 -1px currentColor;margin:auto;opacity:0;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-settings-menu .jw-icon.jw-button-color[aria-checked="true"]::after{opacity:1}.jw-settings-menu .jw-settings-reset{text-decoration:underline}.jw-settings-topbar{align-items:center;background-color:rgba(0,0,0,0.4);display:flex;flex:0 0 auto;padding:3px 5px 0;width:100%}.jw-settings-topbar.jw-nested-menu-open{padding:0}.jw-settings-topbar.jw-nested-menu-open .jw-icon:not(.jw-settings-close):not(.jw-settings-back){display:none}.jw-settings-topbar.jw-nested-menu-open .jw-svg-icon-close{width:20px}.jw-settings-topbar.jw-nested-menu-open .jw-svg-icon-arrow-left{height:12px}.jw-settings-topbar.jw-nested-menu-open .jw-settings-topbar-text{display:block;outline:none}.jw-settings-topbar .jw-settings-back{min-width:44px}.jw-settings-topbar .jw-settings-topbar-buttons{display:inherit;width:100%;height:100%}.jw-settings-topbar .jw-settings-topbar-text{display:none;color:#fff;font-size:13px;width:100%}.jw-settings-topbar .jw-settings-close{margin-left:auto}.jw-settings-submenu{display:none;flex:1 1 auto;overflow-y:auto;padding:8px 20px 0 5px}.jw-settings-submenu::-webkit-scrollbar{background-color:transparent;width:6px}.jw-settings-submenu::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jw-settings-submenu.jw-settings-submenu-active{display:block}.jw-settings-submenu .jw-submenu-topbar{box-shadow:0 2px 9px 0 #1d1d1d;background-color:#2f2d2d;margin:-8px -20px 0 -5px}.jw-settings-submenu .jw-submenu-topbar .jw-settings-content-item{cursor:pointer;text-align:right;padding-right:15px;text-decoration:underline}.jw-settings-submenu .jw-settings-value-wrapper{float:right;display:flex;align-items:center}.jw-settings-submenu .jw-settings-value-wrapper .jw-settings-content-item-arrow{display:flex}.jw-settings-submenu .jw-settings-value-wrapper .jw-svg-icon-arrow-right{width:8px;margin-left:5px;height:12px}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item{font-size:1em;padding:11px 15px 11px 30px}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item .jw-settings-item-active::before{justify-content:flex-end}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item .jw-auto-label{font-size:.85em;padding-left:10px}.jw-flag-touch .jw-settings-submenu{overflow-y:scroll;-webkit-overflow-scrolling:touch}.jw-auto-label{font-size:10px;font-weight:initial;opacity:.75;padding-left:5px}.jw-settings-content-item{position:relative;color:rgba(255,255,255,0.8);cursor:pointer;font-size:12px;line-height:1;padding:7px 0 7px 15px;width:100%;text-align:left;outline:none}.jw-settings-content-item:hover{color:#fff}.jw-settings-content-item:focus{font-weight:bold}.jw-flag-small-player .jw-settings-content-item{line-height:1.75}.jw-settings-content-item.jw-tab-focus:focus{border:solid 2px #4d90fe}.jw-settings-item-active{font-weight:bold;position:relative}.jw-settings-item-active::before{height:100%;width:1em;align-items:center;content:"\\2022";display:inline-flex;justify-content:center}.jw-breakpoint-2 .jw-settings-open .jw-display-container,.jw-flag-small-player .jw-settings-open .jw-display-container,.jw-flag-touch .jw-settings-open .jw-display-container{display:none}.jw-breakpoint-2 .jw-settings-open.jw-controls,.jw-flag-small-player .jw-settings-open.jw-controls,.jw-flag-touch .jw-settings-open.jw-controls{z-index:1}.jw-flag-small-player .jw-settings-open .jw-controlbar{display:none}.jw-settings-open .jw-icon-settings::after{opacity:1}.jw-settings-open .jw-tooltip-settings{display:none}.jw-sharing-link{cursor:pointer}.jw-shortcuts-container .jw-switch{position:relative;display:inline-block;transition:ease-out .15s;transition-property:opacity, background;border-radius:18px;width:80px;height:20px;padding:10px;background:rgba(80,80,80,0.8);cursor:pointer;font-size:inherit;vertical-align:middle}.jw-shortcuts-container .jw-switch.jw-tab-focus{outline:solid 2px #4d90fe}.jw-shortcuts-container .jw-switch .jw-switch-knob{position:absolute;top:2px;left:1px;transition:ease-out .15s;box-shadow:0 0 10px rgba(0,0,0,0.4);border-radius:13px;width:15px;height:15px;background:#fefefe}.jw-shortcuts-container .jw-switch:before,.jw-shortcuts-container .jw-switch:after{position:absolute;top:3px;transition:inherit;color:#fefefe}.jw-shortcuts-container .jw-switch:before{content:attr(data-jw-switch-disabled);right:8px}.jw-shortcuts-container .jw-switch:after{content:attr(data-jw-switch-enabled);left:8px;opacity:0}.jw-shortcuts-container .jw-switch[aria-checked="true"]{background:#475470}.jw-shortcuts-container .jw-switch[aria-checked="true"]:before{opacity:0}.jw-shortcuts-container .jw-switch[aria-checked="true"]:after{opacity:1}.jw-shortcuts-container .jw-switch[aria-checked="true"] .jw-switch-knob{left:60px}.jw-idle-icon-text{display:none;line-height:1;position:absolute;text-align:center;text-indent:.35em;top:100%;white-space:nowrap;left:50%;-webkit-transform:translateX(-50%);transform:translateX(-50%)}.jw-idle-label{border-radius:50%;color:#fff;-webkit-filter:drop-shadow(1px 1px 5px rgba(12,26,71,0.25));filter:drop-shadow(1px 1px 5px rgba(12,26,71,0.25));font:normal 16px/1 Arial,Helvetica,sans-serif;position:relative;transition:background-color 150ms cubic-bezier(0, .25, .25, 1);transition-property:background-color,-webkit-filter;transition-property:background-color,filter;transition-property:background-color,filter,-webkit-filter;-webkit-font-smoothing:antialiased}.jw-state-idle .jw-icon-display.jw-idle-label .jw-idle-icon-text{display:block}.jw-state-idle .jw-icon-display.jw-idle-label .jw-svg-icon-play{-webkit-transform:scale(.7, .7);transform:scale(.7, .7)}.jw-breakpoint-0.jw-state-idle .jw-icon-display.jw-idle-label,.jw-breakpoint--1.jw-state-idle .jw-icon-display.jw-idle-label{font-size:12px}.jw-info-overlay{top:50%;position:absolute;left:50%;background:#333;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);display:none;color:#fff;pointer-events:all;-webkit-user-select:text;-moz-user-select:text;-ms-user-select:text;user-select:text;overflow:hidden;flex-direction:column}.jw-info-overlay .jw-info-close{flex:0 0 auto;margin:5px 5px 5px auto}.jw-info-open .jw-info-overlay{display:flex}.jw-info-container{display:flex;flex:1 1 auto;flex-flow:column;margin:0 20px 20px;overflow-y:auto;padding:5px}.jw-info-container [class*="jw-info"]:not(:first-of-type){color:rgba(255,255,255,0.8);padding-top:10px;font-size:12px}.jw-info-container .jw-info-description{margin-bottom:30px;text-align:start}.jw-info-container .jw-info-description:empty{display:none}.jw-info-container .jw-info-duration{text-align:start}.jw-info-container .jw-info-title{text-align:start;font-size:12px;font-weight:bold}.jw-info-container::-webkit-scrollbar{background-color:transparent;width:6px}.jw-info-container::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jw-info-clientid{align-self:flex-end;font-size:12px;color:rgba(255,255,255,0.8);margin:0 20px 20px 44px;text-align:right}.jw-flag-touch .jw-info-open .jw-display-container{display:none}@supports ((-webkit-filter: drop-shadow(0 0 3px #000)) or (filter: drop-shadow(0 0 3px #000))){.jwplayer.jw-ab-drop-shadow .jw-controls .jw-svg-icon,.jwplayer.jw-ab-drop-shadow .jw-controls .jw-icon.jw-text,.jwplayer.jw-ab-drop-shadow .jw-slider-container .jw-rail,.jwplayer.jw-ab-drop-shadow .jw-title{text-shadow:none;box-shadow:none;-webkit-filter:drop-shadow(0 2px 3px rgba(0,0,0,0.3));filter:drop-shadow(0 2px 3px rgba(0,0,0,0.3))}.jwplayer.jw-ab-drop-shadow .jw-button-color{opacity:.8;transition-property:color, opacity}.jwplayer.jw-ab-drop-shadow .jw-button-color:not(:hover){color:#fff;opacity:.8}.jwplayer.jw-ab-drop-shadow .jw-button-color:hover{opacity:1}.jwplayer.jw-ab-drop-shadow .jw-controls-backdrop{background-image:linear-gradient(to bottom, hsla(0, 0%, 0%, 0), hsla(0, 0%, 0%, 0.00787) 10.79%, hsla(0, 0%, 0%, 0.02963) 21.99%, hsla(0, 0%, 0%, 0.0625) 33.34%, hsla(0, 0%, 0%, 0.1037) 44.59%, hsla(0, 0%, 0%, 0.15046) 55.48%, hsla(0, 0%, 0%, 0.2) 65.75%, hsla(0, 0%, 0%, 0.24954) 75.14%, hsla(0, 0%, 0%, 0.2963) 83.41%, hsla(0, 0%, 0%, 0.3375) 90.28%, hsla(0, 0%, 0%, 0.37037) 95.51%, hsla(0, 0%, 0%, 0.39213) 98.83%, hsla(0, 0%, 0%, 0.4));mix-blend-mode:multiply;transition-property:opacity}.jw-state-idle.jwplayer.jw-ab-drop-shadow .jw-controls-backdrop{background-image:linear-gradient(to bottom, hsla(0, 0%, 0%, 0.2), hsla(0, 0%, 0%, 0.19606) 1.17%, hsla(0, 0%, 0%, 0.18519) 4.49%, hsla(0, 0%, 0%, 0.16875) 9.72%, hsla(0, 0%, 0%, 0.14815) 16.59%, hsla(0, 0%, 0%, 0.12477) 24.86%, hsla(0, 0%, 0%, 0.1) 34.25%, hsla(0, 0%, 0%, 0.07523) 44.52%, hsla(0, 0%, 0%, 0.05185) 55.41%, hsla(0, 0%, 0%, 0.03125) 66.66%, hsla(0, 0%, 0%, 0.01481) 78.01%, hsla(0, 0%, 0%, 0.00394) 89.21%, hsla(0, 0%, 0%, 0));background-size:100% 7rem;background-position:50% 0}.jwplayer.jw-ab-drop-shadow.jw-state-idle .jw-controls{background-color:transparent}}.jw-video-thumbnail-container{position:relative;overflow:hidden}.jw-video-thumbnail-container:not(.jw-related-shelf-item-image){height:100%;width:100%}.jw-video-thumbnail-container.jw-video-thumbnail-generated{position:absolute;top:0;left:0}.jw-video-thumbnail-container:hover,.jw-related-item-content:hover .jw-video-thumbnail-container,.jw-related-shelf-item:hover .jw-video-thumbnail-container{cursor:pointer}.jw-video-thumbnail-container:hover .jw-video-thumbnail:not(.jw-video-thumbnail-completed),.jw-related-item-content:hover .jw-video-thumbnail-container .jw-video-thumbnail:not(.jw-video-thumbnail-completed),.jw-related-shelf-item:hover .jw-video-thumbnail-container .jw-video-thumbnail:not(.jw-video-thumbnail-completed){opacity:1}.jw-video-thumbnail-container .jw-video-thumbnail{position:absolute;top:50%;left:50%;bottom:unset;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);width:100%;height:auto;min-width:100%;min-height:100%;opacity:0;transition:opacity .3s ease;object-fit:cover;background:#000}.jw-related-item-next-up .jw-video-thumbnail-container .jw-video-thumbnail{height:100%;width:auto}.jw-video-thumbnail-container .jw-video-thumbnail.jw-video-thumbnail-visible:not(.jw-video-thumbnail-completed){opacity:1}.jw-video-thumbnail-container .jw-video-thumbnail.jw-video-thumbnail-completed{opacity:0}.jw-video-thumbnail-container .jw-video-thumbnail~.jw-svg-icon-play{display:none}.jw-video-thumbnail-container .jw-video-thumbnail+.jw-related-shelf-item-aspect{pointer-events:none}.jw-video-thumbnail-container .jw-video-thumbnail+.jw-related-item-poster-content{pointer-events:none}.jw-state-idle:not(.jw-flag-cast-available) .jw-display{padding:0}.jw-state-idle .jw-controls{background:rgba(0,0,0,0.4)}.jw-state-idle.jw-flag-cast-available:not(.jw-flag-audio-player) .jw-controlbar .jw-slider-time,.jw-state-idle.jw-flag-cardboard-available .jw-controlbar .jw-slider-time,.jw-state-idle.jw-flag-cast-available:not(.jw-flag-audio-player) .jw-controlbar .jw-icon:not(.jw-icon-cardboard):not(.jw-icon-cast):not(.jw-icon-airplay),.jw-state-idle.jw-flag-cardboard-available .jw-controlbar .jw-icon:not(.jw-icon-cardboard):not(.jw-icon-cast):not(.jw-icon-airplay){display:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-icon:focus{border:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-icon .jw-svg-icon-buffer{-webkit-animation:jw-spin 2s linear infinite;animation:jw-spin 2s linear infinite;display:block}@-webkit-keyframes jw-spin{100%{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}@keyframes jw-spin{100%{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}.jwplayer.jw-state-buffering .jw-icon-playback .jw-svg-icon-play{display:none}.jwplayer.jw-state-buffering .jw-icon-display .jw-svg-icon-pause{display:none}.jwplayer.jw-state-playing .jw-display .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-playing .jw-icon-playback .jw-svg-icon-play{display:none}.jwplayer.jw-state-playing .jw-display .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-playing .jw-icon-playback .jw-svg-icon-pause{display:block}.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio) .jw-controls-backdrop{opacity:0}.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio) .jw-logo-bottom-left,.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio):not(.jw-flag-autostart) .jw-logo-bottom-right{bottom:0}.jwplayer .jw-icon-playback .jw-svg-icon-stop{display:none}.jwplayer.jw-state-paused .jw-svg-icon-pause,.jwplayer.jw-state-idle .jw-svg-icon-pause,.jwplayer.jw-state-error .jw-svg-icon-pause,.jwplayer.jw-state-complete .jw-svg-icon-pause{display:none}.jwplayer.jw-state-error .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-complete .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-buffering .jw-icon-display .jw-svg-icon-play{display:none}.jwplayer:not(.jw-state-buffering) .jw-svg-icon-buffer{display:none}.jwplayer:not(.jw-state-complete) .jw-svg-icon-replay{display:none}.jwplayer:not(.jw-state-error) .jw-svg-icon-error{display:none}.jwplayer.jw-state-complete .jw-display .jw-icon-display .jw-svg-icon-replay{display:block}.jwplayer.jw-state-complete .jw-display .jw-text{display:none}.jwplayer.jw-state-complete .jw-controls{background:rgba(0,0,0,0.4);height:100%}.jw-state-idle .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-paused .jw-icon-playback .jw-svg-icon-pause,.jwplayer.jw-state-paused .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-complete .jw-icon-playback .jw-svg-icon-pause{display:none}.jw-state-idle .jw-display-icon-rewind,.jwplayer.jw-state-buffering .jw-display-icon-rewind,.jwplayer.jw-state-complete .jw-display-icon-rewind,body .jw-error .jw-display-icon-rewind,body .jwplayer.jw-state-error .jw-display-icon-rewind,.jw-state-idle .jw-display-icon-next,.jwplayer.jw-state-buffering .jw-display-icon-next,.jwplayer.jw-state-complete .jw-display-icon-next,body .jw-error .jw-display-icon-next,body .jwplayer.jw-state-error .jw-display-icon-next{display:none}body .jw-error .jw-icon-display,body .jwplayer.jw-state-error .jw-icon-display{cursor:default}body .jw-error .jw-icon-display .jw-svg-icon-error,body .jwplayer.jw-state-error .jw-icon-display .jw-svg-icon-error{display:block}body .jw-error .jw-icon-container{position:absolute;width:100%;height:100%;top:0;left:0;bottom:0;right:0}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-preview{display:none}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title{padding-top:4px}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title-primary{width:auto;display:inline-block;padding-right:.5ch}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title-secondary{width:auto;display:inline-block;padding-left:0}body .jwplayer.jw-state-error .jw-controlbar,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-controlbar{display:none}body .jwplayer.jw-state-error .jw-settings-menu,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-settings-menu{height:100%;top:50%;left:50%;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%)}body .jwplayer.jw-state-error .jw-display,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-display{padding:0}body .jwplayer.jw-state-error .jw-logo-bottom-left,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-logo-bottom-left,body .jwplayer.jw-state-error .jw-logo-bottom-right,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-logo-bottom-right{bottom:0}.jwplayer.jw-state-playing.jw-flag-user-inactive .jw-display{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-state-playing:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display,.jwplayer.jw-state-paused:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting):not(.jw-flag-play-rejected) .jw-display{display:none}.jwplayer.jw-state-paused.jw-flag-play-rejected:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display-icon-rewind,.jwplayer.jw-state-paused.jw-flag-play-rejected:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display-icon-next{display:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-text,.jwplayer.jw-state-complete .jw-display .jw-text{display:none}.jwplayer.jw-flag-casting:not(.jw-flag-audio-player) .jw-cast{display:block}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-display-icon-container{display:none}.jwplayer.jw-flag-casting .jw-icon-hd,.jwplayer.jw-flag-casting .jw-captions,.jwplayer.jw-flag-casting .jw-icon-fullscreen,.jwplayer.jw-flag-casting .jw-icon-audio-tracks{display:none}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-icon-volume{display:none}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-icon-airplay{color:#fff}.jw-state-playing.jw-flag-casting:not(.jw-flag-audio-player) .jw-display,.jw-state-paused.jw-flag-casting:not(.jw-flag-audio-player) .jw-display{display:table}.jwplayer.jw-flag-cast-available .jw-icon-cast,.jwplayer.jw-flag-cast-available .jw-icon-airplay{display:flex}.jwplayer.jw-flag-cardboard-available .jw-icon-cardboard{display:flex}.jwplayer.jw-flag-live .jw-display-icon-rewind{visibility:hidden}.jwplayer.jw-flag-live .jw-controlbar .jw-text-elapsed,.jwplayer.jw-flag-live .jw-controlbar .jw-text-duration,.jwplayer.jw-flag-live .jw-controlbar .jw-text-countdown,.jwplayer.jw-flag-live .jw-controlbar .jw-slider-time{display:none}.jwplayer.jw-flag-live .jw-controlbar .jw-text-alt{display:flex}.jwplayer.jw-flag-live .jw-controlbar .jw-overlay:after{display:none}.jwplayer.jw-flag-live .jw-nextup-container{bottom:44px}.jwplayer.jw-flag-live .jw-text-elapsed,.jwplayer.jw-flag-live .jw-text-duration{display:none}.jwplayer.jw-flag-live .jw-text-live{cursor:default}.jwplayer.jw-flag-live .jw-text-live:hover{color:rgba(255,255,255,0.8)}.jwplayer.jw-flag-live.jw-state-playing .jw-icon-playback .jw-svg-icon-stop,.jwplayer.jw-flag-live.jw-state-buffering .jw-icon-playback .jw-svg-icon-stop{display:block}.jwplayer.jw-flag-live.jw-state-playing .jw-icon-playback .jw-svg-icon-pause,.jwplayer.jw-flag-live.jw-state-buffering .jw-icon-playback .jw-svg-icon-pause{display:none}.jw-text-live{height:24px;width:auto;align-items:center;border-radius:1px;color:rgba(255,255,255,0.8);display:flex;font-size:12px;font-weight:bold;margin-right:10px;padding:0 1ch;text-rendering:geometricPrecision;text-transform:uppercase;transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:box-shadow,color}.jw-text-live::before{height:8px;width:8px;background-color:currentColor;border-radius:50%;margin-right:6px;opacity:1;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-text-live.jw-dvr-live{box-shadow:inset 0 0 0 2px currentColor}.jw-text-live.jw-dvr-live::before{opacity:.5}.jw-text-live.jw-dvr-live:hover{color:#fff}.jwplayer.jw-flag-controls-hidden .jw-logo.jw-hide{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-flag-controls-hidden:not(.jw-flag-casting) .jw-logo-top-right{top:0}.jwplayer.jw-flag-controls-hidden .jw-plugin{bottom:.5em}.jwplayer.jw-flag-controls-hidden .jw-nextup-container{bottom:0}.jw-flag-controls-hidden .jw-controlbar,.jw-flag-controls-hidden .jw-display{visibility:hidden;pointer-events:none;opacity:0;transition-delay:0s, 250ms}.jw-flag-controls-hidden .jw-controls-backdrop{opacity:0}.jw-flag-controls-hidden .jw-logo{visibility:visible}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-logo.jw-hide{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-casting) .jw-logo-top-right{top:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-plugin{bottom:.5em}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-nextup-container{bottom:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-controls-hidden) .jw-media{cursor:none;-webkit-cursor-visibility:auto-hide}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing.jw-flag-casting .jw-display{display:table}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-ads) .jw-autostart-mute{display:flex}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-flag-casting .jw-nextup-container{bottom:66px}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-flag-casting.jw-state-idle .jw-nextup-container{display:none}.jw-flag-media-audio .jw-preview{display:block}.jwplayer.jw-flag-ads .jw-preview,.jwplayer.jw-flag-ads .jw-logo,.jwplayer.jw-flag-ads .jw-captions.jw-captions-enabled,.jwplayer.jw-flag-ads .jw-nextup-container,.jwplayer.jw-flag-ads .jw-text-duration,.jwplayer.jw-flag-ads .jw-text-elapsed{display:none}.jwplayer.jw-flag-ads video::-webkit-media-text-track-container{display:none}.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-rewind,.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-next,.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-display{display:none}.jwplayer.jw-flag-ads.jw-flag-small-player.jw-state-buffering .jw-display-icon-display{display:inline-block}.jwplayer.jw-flag-ads .jw-controlbar{flex-wrap:wrap-reverse}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time{height:auto;padding:0;pointer-events:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-slider-container{height:5px}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-rail,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-knob,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-buffer,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-cue,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-icon-settings{display:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-progress{-webkit-transform:none;transform:none;top:auto}.jwplayer.jw-flag-ads .jw-controlbar .jw-tooltip,.jwplayer.jw-flag-ads .jw-controlbar .jw-icon-tooltip:not(.jw-icon-volume),.jwplayer.jw-flag-ads .jw-controlbar .jw-icon-inline:not(.jw-icon-playback):not(.jw-icon-fullscreen):not(.jw-icon-volume){display:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-volume-tip{padding:13px 0}.jwplayer.jw-flag-ads .jw-controlbar .jw-text-alt{display:flex}.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid) .jw-controls .jw-controlbar,.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-autostart .jw-controls .jw-controlbar{display:flex;pointer-events:all;visibility:visible;opacity:1}.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-user-inactive .jw-controls-backdrop,.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-autostart.jw-flag-user-inactive .jw-controls-backdrop{opacity:1;background-size:100% 60px}.jwplayer.jw-flag-ads-vpaid .jw-display-container,.jwplayer.jw-flag-touch.jw-flag-ads-vpaid .jw-display-container,.jwplayer.jw-flag-ads-vpaid .jw-skip,.jwplayer.jw-flag-touch.jw-flag-ads-vpaid .jw-skip{display:none}.jwplayer.jw-flag-ads-vpaid.jw-flag-small-player .jw-controls{background:none}.jwplayer.jw-flag-ads-vpaid.jw-flag-small-player .jw-controls::after{content:none}.jwplayer.jw-flag-ads-hide-controls .jw-controls-backdrop,.jwplayer.jw-flag-ads-hide-controls .jw-controls{display:none !important}.jw-flag-overlay-open-related .jw-controls,.jw-flag-overlay-open-related .jw-title,.jw-flag-overlay-open-related .jw-logo{display:none}.jwplayer.jw-flag-rightclick-open{overflow:visible}.jwplayer.jw-flag-rightclick-open .jw-rightclick{z-index:16777215}body .jwplayer.jw-flag-flash-blocked .jw-controls,body .jwplayer.jw-flag-flash-blocked .jw-overlays,body .jwplayer.jw-flag-flash-blocked .jw-controls-backdrop,body .jwplayer.jw-flag-flash-blocked .jw-preview{display:none}body .jwplayer.jw-flag-flash-blocked .jw-error-msg{top:25%}.jw-flag-touch.jw-breakpoint-7 .jw-captions,.jw-flag-touch.jw-breakpoint-6 .jw-captions,.jw-flag-touch.jw-breakpoint-5 .jw-captions,.jw-flag-touch.jw-breakpoint-4 .jw-captions,.jw-flag-touch.jw-breakpoint-7 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-6 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-5 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-4 .jw-nextup-container{bottom:4.25em}.jw-flag-touch .jw-controlbar .jw-icon-volume{display:flex}.jw-flag-touch .jw-display,.jw-flag-touch .jw-display-container,.jw-flag-touch .jw-display-controls{pointer-events:none}.jw-flag-touch.jw-state-paused:not(.jw-breakpoint-1) .jw-display-icon-next,.jw-flag-touch.jw-state-playing:not(.jw-breakpoint-1) .jw-display-icon-next,.jw-flag-touch.jw-state-paused:not(.jw-breakpoint-1) .jw-display-icon-rewind,.jw-flag-touch.jw-state-playing:not(.jw-breakpoint-1) .jw-display-icon-rewind{display:none}.jw-flag-touch.jw-state-paused.jw-flag-dragging .jw-display{display:none}.jw-flag-audio-player{background-color:#000}.jw-flag-audio-player:not(.jw-flag-flash-blocked) .jw-media{visibility:hidden}.jw-flag-audio-player .jw-title{background:none}.jw-flag-audio-player object{min-height:44px}.jw-flag-audio-player:not(.jw-flag-live) .jw-spacer{display:none}.jw-flag-audio-player .jw-preview,.jw-flag-audio-player .jw-display,.jw-flag-audio-player .jw-title,.jw-flag-audio-player .jw-nextup-container{display:none}.jw-flag-audio-player .jw-controlbar{position:relative}.jw-flag-audio-player .jw-controlbar .jw-button-container{padding-right:3px;padding-left:0}.jw-flag-audio-player .jw-controlbar .jw-icon-tooltip,.jw-flag-audio-player .jw-controlbar .jw-icon-inline{display:none}.jw-flag-audio-player .jw-controlbar .jw-icon-volume,.jw-flag-audio-player .jw-controlbar .jw-icon-playback,.jw-flag-audio-player .jw-controlbar .jw-icon-next,.jw-flag-audio-player .jw-controlbar .jw-icon-rewind,.jw-flag-audio-player .jw-controlbar .jw-icon-cast,.jw-flag-audio-player .jw-controlbar .jw-text-live,.jw-flag-audio-player .jw-controlbar .jw-icon-airplay,.jw-flag-audio-player .jw-controlbar .jw-logo-button,.jw-flag-audio-player .jw-controlbar .jw-text-elapsed,.jw-flag-audio-player .jw-controlbar .jw-text-duration{display:flex;flex:0 0 auto}.jw-flag-audio-player .jw-controlbar .jw-text-duration,.jw-flag-audio-player .jw-controlbar .jw-text-countdown{padding-right:10px}.jw-flag-audio-player .jw-controlbar .jw-slider-time{flex:0 1 auto;align-items:center;display:flex;order:1}.jw-flag-audio-player .jw-controlbar .jw-icon-volume{margin-right:0;transition:margin-right 150ms cubic-bezier(0, .25, .25, 1)}.jw-flag-audio-player .jw-controlbar .jw-icon-volume .jw-overlay{display:none}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container{transition:width 300ms cubic-bezier(0, .25, .25, 1);width:0}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open{width:140px}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open .jw-slider-volume{padding-right:24px;transition:opacity 300ms;opacity:1}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open~.jw-slider-time{flex:1 1 auto;width:auto;transition:opacity 300ms, width 300ms}.jw-flag-audio-player .jw-controlbar .jw-slider-volume{opacity:0}.jw-flag-audio-player .jw-controlbar .jw-slider-volume .jw-knob{-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%)}.jw-flag-audio-player .jw-controlbar .jw-slider-volume~.jw-icon-volume{margin-right:140px}.jw-flag-audio-player.jw-breakpoint-1 .jw-horizontal-volume-container.jw-open~.jw-slider-time,.jw-flag-audio-player.jw-breakpoint-2 .jw-horizontal-volume-container.jw-open~.jw-slider-time{opacity:0}.jw-flag-audio-player.jw-flag-small-player .jw-text-elapsed,.jw-flag-audio-player.jw-flag-small-player .jw-text-duration{display:none}.jw-flag-audio-player.jw-flag-ads .jw-slider-time{display:none}.jw-hidden{display:none}',
        "",
      ]);
    },
  ],
]);
