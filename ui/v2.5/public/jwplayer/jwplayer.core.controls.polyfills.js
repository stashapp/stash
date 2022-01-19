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
  [5, 1, 2, 3, 7],
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
        s = n(43),
        l = n(5),
        c = n(15),
        u = n(40);
      function d(t) {
        return (
          i || (i = new DOMParser()),
          Object(l.r)(
            Object(l.s)(i.parseFromString(t, "image/svg+xml").documentElement)
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
                "string" == typeof t ? o.appendChild(d(t)) : o.appendChild(t);
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
        w = n(0),
        h = n(71),
        f = n.n(h),
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
        T = n.n(M),
        S = n(78),
        _ = n.n(S),
        E = n(79),
        z = n.n(E),
        P = n(80),
        A = n.n(P),
        I = n(81),
        R = n.n(I),
        L = n(82),
        B = n.n(L),
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
        tt = n(90),
        et = n.n(tt),
        nt = n(91),
        it = n.n(nt),
        ot = n(92),
        at = n.n(ot),
        rt = n(93),
        st = n.n(rt),
        lt = n(94),
        ct = n.n(lt),
        ut = null;
      function dt(t) {
        var e = ft().querySelector(wt(t));
        if (e) return ht(e);
        throw new Error("Icon not found " + t);
      }
      function pt(t) {
        var e = ft().querySelectorAll(t.split(",").map(wt).join(","));
        if (!e.length) throw new Error("Icons not found " + t);
        return Array.prototype.map.call(e, function (t) {
          return ht(t);
        });
      }
      function wt(t) {
        return ".jw-svg-icon-".concat(t);
      }
      function ht(t) {
        return t.cloneNode(!0);
      }
      function ft() {
        return (
          ut ||
            (ut = d(
              "<xml>" +
                f.a +
                g.a +
                m.a +
                y.a +
                x.a +
                C.a +
                T.a +
                _.a +
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
                et.a +
                it.a +
                at.a +
                st.a +
                ct.a +
                "</xml>"
            )),
          ut
        );
      }
      var jt = n(10);
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
      var mt = (function () {
          function t(e, n, i, o, a) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t);
            var r,
              s = document.createElement("div");
            (s.className = "jw-icon jw-icon-inline jw-button-color jw-reset ".concat(
              a || ""
            )),
              s.setAttribute("button", o),
              s.setAttribute("role", "button"),
              s.setAttribute("tabindex", "0"),
              n && s.setAttribute("aria-label", n),
              e && "<svg" === e.substring(0, 4)
                ? (r = (function (t) {
                    if (!bt[t]) {
                      var e = Object.keys(bt);
                      e.length > 10 && delete bt[e[0]];
                      var n = d(t);
                      bt[t] = n;
                    }
                    return bt[t].cloneNode(!0);
                  })(e))
                : (((r = document.createElement("div")).className =
                    "jw-icon jw-button-image jw-button-color jw-reset"),
                  e &&
                    Object(jt.d)(r, {
                      backgroundImage: "url(".concat(e, ")"),
                    })),
              s.appendChild(r),
              new u.a(s).on("click tap enter", i, this),
              s.addEventListener("mousedown", function (t) {
                t.preventDefault();
              }),
              (this.id = o),
              (this.buttonElement = s);
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
      function yt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var kt = function (t) {
          var e = Object(l.c)(t),
            n = window.pageXOffset;
          return (
            n &&
              o.OS.android &&
              document.body.parentElement.getBoundingClientRect().left >= 0 &&
              ((e.left -= n), (e.right -= n)),
            e
          );
        },
        xt = (function () {
          function t(e, n) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(w.g)(this, r.a),
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
                  (this.el = Object(l.e)(
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
                    (this.railBounds = kt(this.elementRail));
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
                      : kt(this.elementRail));
                  return (
                    (n =
                      "horizontal" === this.orientation
                        ? (e = t.pageX) < i.left
                          ? 0
                          : e > i.right
                          ? 100
                          : 100 * Object(s.a)((e - i.left) / i.width, 0, 1)
                        : (e = t.pageY) >= i.bottom
                        ? 0
                        : e <= i.top
                        ? 100
                        : 100 *
                          Object(s.a)(
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
                  (this.railBounds = kt(this.elementRail)), this.dragMove(t);
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
            ]) && yt(e.prototype, n),
            i && yt(e, i),
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
      var Mt = (function () {
          function t(e, n, i, o) {
            var a = this;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(w.g)(this, r.a),
              (this.el = document.createElement("div"));
            var s =
              "jw-icon jw-icon-tooltip " + e + " jw-button-color jw-reset";
            i || (s += " jw-hidden"),
              Ot(this.el, n),
              (this.el.className = s),
              (this.tooltip = document.createElement("div")),
              (this.tooltip.className = "jw-overlay jw-reset"),
              (this.openClass = "jw-open"),
              (this.componentType = "tooltip"),
              this.el.appendChild(this.tooltip),
              o &&
                o.length > 0 &&
                Array.prototype.forEach.call(o, function (t) {
                  "string" == typeof t
                    ? a.el.appendChild(d(t))
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
                    Object(l.v)(this.el, this.openClass, this.isOpen));
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
                    Object(l.v)(this.el, this.openClass, this.isOpen));
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
        St = n(57);
      function _t(t, e) {
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
            ]) && _t(e.prototype, n),
            i && _t(e, i),
            t
          );
        })(),
        zt = {
          loadChapters: function (t) {
            Object(Tt.a)(
              t,
              this.chaptersLoaded.bind(this),
              this.chaptersFailed,
              { plainText: !0 }
            );
          },
          chaptersLoaded: function (t) {
            var e = Object(St.a)(t.responseText);
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
      function Pt(t) {
        (this.begin = t.begin), (this.end = t.end), (this.img = t.text);
      }
      var At = {
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
          var e = Object(St.a)(t.responseText);
          Array.isArray(e) &&
            (e.forEach(function (t) {
              this.thumbnails.push(new Pt(t));
            }, this),
            this.drawCues());
        },
        thumbnailsFailed: function () {},
        chooseThumbnail: function (t) {
          var e = Object(w.A)(this.thumbnails, { end: t }, Object(w.z)("end"));
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
              (this.individualImage.onload = Object(w.a)(function () {
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
      function It(t, e, n) {
        return (It =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, n) {
                var i = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = Ht(t));

                  );
                  return t;
                })(t, e);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, e);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(t, e, n || t);
      }
      function Rt(t) {
        return (Rt =
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
      function Lt(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Bt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function Vt(t, e, n) {
        return e && Bt(t.prototype, e), n && Bt(t, n), t;
      }
      function Nt(t, e) {
        return !e || ("object" !== Rt(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Ht(t) {
        return (Ht = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function Ft(t, e) {
        if ("function" != typeof e && null !== e)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (t.prototype = Object.create(e && e.prototype, {
          constructor: { value: t, writable: !0, configurable: !0 },
        })),
          e && qt(t, e);
      }
      function qt(t, e) {
        return (qt =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var Dt = (function (t) {
        function e() {
          return Lt(this, e), Nt(this, Ht(e).apply(this, arguments));
        }
        return (
          Ft(e, t),
          Vt(e, [
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
                Object(jt.d)(this.img, t);
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
          e
        );
      })(Mt);
      var Ut = (function (t) {
        function e(t, n, i) {
          var o;
          return (
            Lt(this, e),
            ((o = Nt(
              this,
              Ht(e).call(this, "jw-slider-time", "horizontal")
            ))._model = t),
            (o._api = n),
            (o.timeUpdateKeeper = i),
            (o.timeTip = new Dt("jw-tooltip-time", null, !0)),
            o.timeTip.setup(),
            (o.cues = []),
            (o.seekThrottled = Object(w.B)(o.performSeek, 400)),
            (o.mobileHoverDistance = 5),
            o.setup(),
            o
          );
        }
        return (
          Ft(e, t),
          Vt(e, [
            {
              key: "setup",
              value: function () {
                var t = this;
                It(Ht(e.prototype), "setup", this).apply(this, arguments),
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
              value: function (t) {
                (this.seekTo = t),
                  this.seekThrottled(),
                  It(Ht(e.prototype), "update", this).apply(this, arguments);
              },
            },
            {
              key: "dragStart",
              value: function () {
                this._model.set("scrubbing", !0),
                  It(Ht(e.prototype), "dragStart", this).apply(this, arguments);
              },
            },
            {
              key: "dragEnd",
              value: function () {
                It(Ht(e.prototype), "dragEnd", this).apply(this, arguments),
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
                  Object(l.t)(this.el, "aria-valuemin", 0),
                  Object(l.t)(this.el, "aria-valuemax", e),
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
                Object(w.f)(
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
                    a = Object(l.c)(this.elementRail),
                    r = t.pageX ? t.pageX - a.left : t.x,
                    c = (r = Object(s.a)(r, 0, a.width)) / a.width,
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
                  var d = this.timeTip;
                  d.update(i),
                    this.textLength !== i.length &&
                      ((this.textLength = i.length), d.resetWidth()),
                    this.showThumbnail(u),
                    Object(l.a)(d.el, "jw-open");
                  var p = d.getWidth(),
                    w = a.width / 100,
                    h = o - a.width,
                    f = 0;
                  p > h && (f = (p - h) / (200 * w));
                  var j = 100 * Math.min(1 - f, Math.max(f, c)).toFixed(3);
                  Object(jt.d)(d.el, { left: j + "%" });
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
                    Object(l.t)(o, "aria-valuenow", e),
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
          e
        );
      })(xt);
      Object(w.g)(Ut.prototype, zt, At);
      var Wt = Ut;
      function Qt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function Yt(t, e, n) {
        return (Yt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, n) {
                var i = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = Gt(t));

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
      function Zt(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Kt(t, e) {
        return !e || ("object" !== Xt(e) && "function" != typeof e) ? Jt(t) : e;
      }
      function Jt(t) {
        if (void 0 === t)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return t;
      }
      function Gt(t) {
        return (Gt = Object.setPrototypeOf
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
            Zt(this, e);
            var a = "jw-slider-volume";
            return (
              "vertical" === t && (a += " jw-volume-tip"),
              (o = Kt(this, Gt(e).call(this, a, t))).setup(),
              o.element().classList.remove("jw-background-color"),
              Object(l.t)(i, "tabindex", "0"),
              Object(l.t)(i, "aria-label", n),
              Object(l.t)(i, "aria-orientation", t),
              Object(l.t)(i, "aria-valuemin", 0),
              Object(l.t)(i, "aria-valuemax", 100),
              Object(l.t)(i, "role", "slider"),
              (o.uiOver = new u.a(i).on("click", function () {})),
              o
            );
          }
          return $t(e, t), e;
        })(xt),
        ne = (function (t) {
          function e(t, n, i, o, a) {
            var r;
            Zt(this, e),
              ((r = Kt(this, Gt(e).call(this, n, i, !0, o)))._model = t),
              (r.horizontalContainer = a);
            var s = t.get("localization").volumeSlider;
            return (
              (r.horizontalSlider = new ee("horizontal", s, a, Jt(Jt(r)))),
              (r.verticalSlider = new ee("vertical", s, r.tooltip, Jt(Jt(r)))),
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
                  Object(l.t)(this.horizontalContainer, "tabindex", n);
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
                  Yt(Gt(e.prototype), "openTooltip", this).call(this, t),
                    Object(l.v)(this.horizontalContainer, this.openClass, !0);
                },
              },
              {
                key: "closeSlider",
                value: function (t) {
                  Yt(Gt(e.prototype), "closeTooltip", this).call(this, t),
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
            ]) && Qt(n.prototype, i),
            o && Qt(n, o),
            e
          );
        })(Mt);
      function ie(t, e, n, i, o) {
        var a = document.createElement("div");
        (a.className = "jw-reset-text jw-tooltip jw-tooltip-".concat(e)),
          a.setAttribute("dir", "auto");
        var r = document.createElement("div");
        (r.className = "jw-text"), a.appendChild(r), t.appendChild(a);
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
            setText: function (t) {
              t !== s.text && ((s.text = t), (s.dirty = !0)), s.opened && c(!0);
            },
          },
          c = function (t) {
            t && s.dirty && (Object(l.q)(r, s.text), (s.dirty = !1)),
              (s.opened = t),
              Object(l.v)(a, "jw-open", t);
          };
        return (
          t.addEventListener("mouseover", s.open),
          t.addEventListener("focus", s.open),
          t.addEventListener("blur", s.close),
          t.addEventListener("mouseout", s.close),
          t.addEventListener(
            "touchstart",
            function () {
              s.touchEvent = !0;
            },
            { passive: !0 }
          ),
          s
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
          e && Object(l.t)(n, "role", e),
          n
        );
      }
      function se(t) {
        var e = document.createElement("div");
        return (e.className = "jw-reset ".concat(t)), e;
      }
      function le(t, e) {
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
          Object(l.t)(i, "tabindex", "-1"), (i.className += " jw-reset");
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
        de = (function () {
          function t(e, n, i) {
            var s = this;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(w.g)(this, r.a),
              (this._api = e),
              (this._model = n),
              (this._isMobile = o.OS.mobile),
              (this._volumeAnnouncer = i.querySelector(".jw-volume-update"));
            var c,
              d,
              h,
              f = n.get("localization"),
              j = new Wt(n, e, i.querySelector(".jw-time-update")),
              g = (this.menus = []);
            this.ui = [];
            var b = "",
              m = f.volume;
            if (this._isMobile) {
              if (
                !(n.get("sdkplatform") || (o.OS.iOS && o.OS.version.major < 10))
              ) {
                var v = pt("volume-0,volume-100");
                h = p(
                  "jw-icon-volume",
                  function () {
                    e.setMute();
                  },
                  m,
                  v
                );
              }
            } else {
              (d = document.createElement("div")).className =
                "jw-horizontal-volume-container";
              var y = (c = new ne(
                n,
                "jw-icon-volume",
                m,
                pt("volume-0,volume-50,volume-100"),
                d
              )).element();
              g.push(c),
                Object(l.t)(y, "role", "button"),
                n.change(
                  "mute",
                  function (t, e) {
                    var n = e ? f.unmute : f.mute;
                    Object(l.t)(y, "aria-label", n);
                  },
                  this
                );
            }
            var k = p(
                "jw-icon-next",
                function () {
                  e.next({ feedShownId: b, reason: "interaction" });
                },
                f.next,
                pt("next")
              ),
              x = p(
                "jw-icon-settings jw-settings-submenu-button",
                function (t) {
                  s.trigger("settingsInteraction", "quality", !0, t);
                },
                f.settings,
                pt("settings")
              );
            Object(l.t)(x.element(), "aria-haspopup", "true");
            var O = p(
              "jw-icon-cc jw-settings-submenu-button",
              function (t) {
                s.trigger("settingsInteraction", "captions", !1, t);
              },
              f.cc,
              pt("cc-off,cc-on")
            );
            Object(l.t)(O.element(), "aria-haspopup", "true");
            var C = p(
              "jw-text-live",
              function () {
                s.goToLiveEdge();
              },
              f.liveBroadcast
            );
            C.element().textContent = f.liveBroadcast;
            var M,
              T,
              S,
              _ = (this.elements = {
                alt:
                  ((M = "jw-text-alt"),
                  (T = "status"),
                  (S = document.createElement("span")),
                  (S.className = "jw-text jw-reset " + M),
                  T && Object(l.t)(S, "role", T),
                  S),
                play: p(
                  "jw-icon-playback",
                  function () {
                    e.playToggle({ reason: "interaction" });
                  },
                  f.play,
                  pt("play,pause,stop")
                ),
                rewind: p(
                  "jw-icon-rewind",
                  function () {
                    s.rewind();
                  },
                  f.rewind,
                  pt("rewind")
                ),
                live: C,
                next: k,
                elapsed: re("jw-text-elapsed", "timer"),
                countdown: re("jw-text-countdown", "timer"),
                time: j,
                duration: re("jw-text-duration", "timer"),
                mute: h,
                volumetooltip: c,
                horizontalVolumeContainer: d,
                cast: le(function () {
                  e.castToggle();
                }, f),
                fullscreen: p(
                  "jw-icon-fullscreen",
                  function () {
                    e.setFullscreen();
                  },
                  f.fullscreen,
                  pt("fullscreen-off,fullscreen-on")
                ),
                spacer: se("jw-spacer"),
                buttonContainer: se("jw-button-container"),
                settingsButton: x,
                captionsButton: O,
              }),
              E = ie(O.element(), "captions", f.cc),
              z = function (t) {
                var e = t.get("captionsList")[t.get("captionsIndex")],
                  n = f.cc;
                e && "Off" !== e.label && (n = e.label), E.setText(n);
              },
              P = ie(_.play.element(), "play", f.play);
            this.setPlayText = function (t) {
              P.setText(t);
            };
            var A = _.next.element(),
              I = ie(
                A,
                "next",
                f.nextUp,
                function () {
                  var t = n.get("nextUp");
                  (b = Object(oe.b)(oe.a)),
                    s.trigger("nextShown", {
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
            Object(l.t)(A, "dir", "auto"),
              ie(_.rewind.element(), "rewind", f.rewind),
              ie(_.settingsButton.element(), "settings", f.settings);
            var R = ie(_.fullscreen.element(), "fullscreen", f.fullscreen),
              L = [
                _.play,
                _.rewind,
                _.next,
                _.volumetooltip,
                _.mute,
                _.horizontalVolumeContainer,
                _.alt,
                _.live,
                _.elapsed,
                _.countdown,
                _.duration,
                _.spacer,
                _.cast,
                _.captionsButton,
                _.settingsButton,
                _.fullscreen,
              ].filter(function (t) {
                return t;
              }),
              B = [_.time, _.buttonContainer].filter(function (t) {
                return t;
              });
            (this.el = document.createElement("div")),
              (this.el.className = "jw-controlbar jw-reset"),
              ue(_.buttonContainer, L),
              ue(this.el, B);
            var V = n.get("logo");
            if (
              (V && "control-bar" === V.position && this.addLogo(V),
              _.play.show(),
              _.fullscreen.show(),
              _.mute && _.mute.show(),
              n.change("volume", this.onVolume, this),
              n.change(
                "mute",
                function (t, e) {
                  s.renderVolume(e, t.get("volume"));
                },
                this
              ),
              n.change("state", this.onState, this),
              n.change("duration", this.onDuration, this),
              n.change("position", this.onElapsed, this),
              n.change(
                "fullscreen",
                function (t, e) {
                  var n = s.elements.fullscreen.element();
                  Object(l.v)(n, "jw-off", e);
                  var i = t.get("fullscreen") ? f.exitFullscreen : f.fullscreen;
                  R.setText(i), Object(l.t)(n, "aria-label", i);
                },
                this
              ),
              n.change("streamType", this.onStreamTypeChange, this),
              n.change(
                "dvrLive",
                function (t, e) {
                  var n = f.liveBroadcast,
                    i = f.notLive,
                    o = s.elements.live.element(),
                    a = !1 === e;
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
                function (t, e) {
                  b = Object(oe.b)(oe.a);
                  var n = f.nextUp;
                  e && e.title && (n += ": ".concat(e.title)),
                    I.setText(n),
                    _.next.toggle(!!e);
                },
                this
              ),
              n.change("audioMode", this.onAudioMode, this),
              _.cast &&
                (n.change("castAvailable", this.onCastAvailable, this),
                n.change("castActive", this.onCastActive, this)),
              _.volumetooltip &&
                (_.volumetooltip.on(
                  "update",
                  function (t) {
                    var e = t.percentage;
                    this._api.setVolume(e);
                  },
                  this
                ),
                _.volumetooltip.on(
                  "toggleValue",
                  function () {
                    this._api.setMute();
                  },
                  this
                ),
                _.volumetooltip.on(
                  "adjustVolume",
                  function (t) {
                    this.trigger("adjustVolume", t);
                  },
                  this
                )),
              _.cast && _.cast.button)
            ) {
              var N = _.cast.ui.on(
                "click tap enter",
                function (t) {
                  "click" !== t.type && _.cast.button.click(),
                    this._model.set("castClicked", !0);
                },
                this
              );
              this.ui.push(N);
            }
            var H = new u.a(_.duration).on(
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
            this.ui.push(H);
            var F = new u.a(this.el).on(
              "click tap drag",
              function () {
                this.trigger(a.sb);
              },
              this
            );
            this.ui.push(F),
              g.forEach(function (t) {
                t.on("open-tooltip", s.closeMenus, s);
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
                      (Object(l.v)(n.element(), "jw-off", t),
                      Object(l.v)(n.element(), "jw-full", !t)),
                    i)
                  ) {
                    var o = t ? 0 : e,
                      a = i.element();
                    i.verticalSlider.render(o), i.horizontalSlider.render(o);
                    var r = i.tooltip,
                      s = i.horizontalContainer;
                    Object(l.v)(a, "jw-off", t),
                      Object(l.v)(a, "jw-full", e >= 75 && !t),
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
                value: function (t, e) {
                  this.elements.cast.toggle(e);
                },
              },
              {
                key: "onCastActive",
                value: function (t, e) {
                  this.elements.fullscreen.toggle(!e),
                    this.elements.cast.button &&
                      Object(l.v)(this.elements.cast.button, "jw-off", !e);
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
                    Object(l.t)(this.elements.play.element(), "aria-label", i);
                },
              },
              {
                key: "onStreamTypeChange",
                value: function (t, e) {
                  var n = "LIVE" === e,
                    i = "DVR" === e;
                  this.elements.rewind.toggle(!n),
                    this.elements.live.toggle(n || i),
                    Object(l.t)(
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
                    n = new mt(
                      t.file,
                      this._model.get("localization").logo,
                      function () {
                        t.link &&
                          Object(l.l)(t.link, "_blank", { rel: "noreferrer" });
                      },
                      "logo",
                      "jw-logo-button"
                    );
                  t.link || Object(l.t)(n.element(), "tabindex", "-1"),
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
                      var s = i[r],
                        l = new mt(
                          s.img,
                          s.tooltip,
                          s.callback,
                          s.id,
                          s.btnClass
                        );
                      s.tooltip && ie(l.element(), s.id, s.tooltip);
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
                  e && Object(l.v)(e.element(), "jw-off", !t);
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
        we = function (t) {
          return (
            '<div class="jw-display jw-reset"><div class="jw-display-container jw-reset"><div class="jw-display-controls jw-reset">' +
            pe("rewind", t.rewind) +
            pe("display", t.playback) +
            pe("next", t.next) +
            "</div></div></div>"
          );
        };
      function he(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var fe = (function () {
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
          ]) && he(e.prototype, n),
          i && he(e, i),
          t
        );
      })();
      function je(t) {
        return (je =
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
        return !e || ("object" !== je(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function me(t) {
        return (me = Object.setPrototypeOf
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
      var ye = (function (t) {
        function e(t, n, i) {
          var o;
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, e),
            (o = be(this, me(e).call(this)));
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
      function ke(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var xe = (function () {
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
          ]) && ke(e.prototype, n),
          i && ke(e, i),
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
            (this.el = Object(l.e)(we(e.get("localization"))));
          var i = this.el.querySelector(".jw-display-controls"),
            o = {};
          Me("rewind", pt("rewind"), fe, i, o, e, n),
            Me("display", pt("play,pause,buffer,replay"), ye, i, o, e, n),
            Me("next", pt("next"), xe, i, o, e, n),
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
      function Me(t, e, n, i, o, a, r) {
        var s = i.querySelector(".jw-display-icon-".concat(t)),
          l = i.querySelector(".jw-icon-".concat(t));
        e.forEach(function (t) {
          l.appendChild(t);
        }),
          (o[t] = new n(a, r, s));
      }
      var Te = n(2);
      function Se(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var _e = (function () {
          function t(e, n, i) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(w.g)(this, r.a),
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
                  var e = Object(l.e)(
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
                  e.querySelector(".jw-nextup-close").appendChild(dt("close")),
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
                    (Object(l.v)(
                      this.container,
                      "jw-nextup-sticky",
                      !!this.nextUpSticky
                    ),
                    this.shown !== t)
                  ) {
                    (this.shown = t),
                      Object(l.v)(
                        this.container,
                        "jw-nextup-container-visible",
                        t
                      ),
                      Object(l.v)(this._playerElement, "jw-flag-nextup", t);
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
                      Object(l.v)(
                        e.content,
                        "jw-nextup-thumbnail-visible",
                        !!t.image
                      ),
                      t.image)
                    ) {
                      var n = e.loadThumbnail(t.image);
                      Object(jt.d)(e.thumbnail, n);
                    }
                    (e.header = e.content.querySelector(".jw-nextup-header")),
                      (e.header.textContent = Object(l.e)(
                        e.localization.nextUp
                      ).textContent),
                      (e.title = e.content.querySelector(".jw-nextup-title"));
                    var i = t.title;
                    e.title.textContent = i ? Object(l.e)(i).textContent : "";
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
            ]) && Se(e.prototype, n),
            i && Se(e, i),
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
              .concat(ze[o](t, e), "</li>")
          );
        },
        ze = {
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
        Pe = n(23),
        Ae = n(6),
        Ie = n(13);
      function Re(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var Le = {
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
      function Be(t) {
        var e = Object(l.e)(t),
          n = e.querySelector(".jw-rightclick-logo");
        return n && n.appendChild(dt("jwplayer-logo")), e;
      }
      var Ve = (function () {
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
                  var t = Pe.a.split("+")[0],
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
                          title: Object(Ie.e)(i)
                            ? "".concat(o, " ").concat(i)
                            : "".concat(i, " ").concat(o),
                          type: "link",
                          featured: !0,
                          showLogo: !0,
                          link: "https://jwplayer.com/learn-more?e=".concat(
                            Le[n]
                          ),
                        },
                      ],
                    },
                    r = e.get("provider"),
                    s = a.items;
                  if (r && r.name.indexOf("flash") >= 0) {
                    var l = "Flash Version " + Object(Ae.a)();
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
                value: function (t) {
                  if ((this.lazySetup(), this.mouseOverContext)) return !1;
                  this.hideMenu(), this.showMenu(t), this.addHideMenuHandlers();
                },
              },
              {
                key: "getOffset",
                value: function (t) {
                  var e = Object(l.c)(this.wrapperElement),
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
                    Object(l.a)(
                      this.playerContainer,
                      "jw-flag-rightclick-open"
                    ),
                    Object(l.a)(this.el, "jw-open"),
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
                      var r = Be(a);
                      Object(l.h)(this.el);
                      for (var s = r.childNodes.length; s--; )
                        this.el.appendChild(r.firstChild);
                    }
                  } else
                    (this.html = a),
                      (this.el = Be(this.html)),
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
            ]) && Re(e.prototype, n),
            i && Re(e, i),
            t
          );
        })(),
        Ne = function (t) {
          return '<button type="button" class="jw-reset-text jw-settings-content-item" dir="auto">'.concat(
            t,
            "</button>"
          );
        },
        He = function (t) {
          return (
            '<button type="button" class="jw-reset-text jw-settings-content-item" dir="auto">' +
            "".concat(t.label) +
            "<div class='jw-reset jw-settings-value-wrapper'>" +
            '<div class="jw-reset-text jw-settings-content-item-value">'.concat(
              t.value,
              "</div>"
            ) +
            '<div class="jw-reset-text jw-settings-content-item-arrow">'.concat(
              Y.a,
              "</div>"
            ) +
            "</div></button>"
          );
        },
        Fe = function (t) {
          return (
            '<button type="button" class="jw-reset-text jw-settings-content-item" role="menuitemradio" aria-checked="false" dir="auto">' +
            "".concat(t) +
            "</button>"
          );
        };
      function qe(t) {
        return (qe =
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
      function De(t, e) {
        return !e || ("object" !== qe(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Ue(t) {
        return (Ue = Object.setPrototypeOf
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
      function Qe(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Ye(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function Xe(t, e, n) {
        return e && Ye(t.prototype, e), n && Ye(t, n), t;
      }
      var Ze,
        Ke = (function () {
          function t(e, n) {
            var i =
              arguments.length > 2 && void 0 !== arguments[2]
                ? arguments[2]
                : Ne;
            Qe(this, t),
              (this.el = Object(l.e)(i(e))),
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
                : Fe;
            return Qe(this, e), De(this, Ue(e).call(this, t, n, i));
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
            e
          );
        })(Ke),
        Ge = function (t, e) {
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
                [(t.icon && Object(l.e)(t.icon)) || dt(i)]
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
              s =
                arguments.length > 3 && void 0 !== arguments[3]
                  ? arguments[3]
                  : Ge;
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
              (o.el = Object(l.e)(s(o.isSubmenu, t))),
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
                          s = o.sourceEvent,
                          c = s.target,
                          u = i.firstChild === c,
                          d = i.lastChild === c,
                          p = n.topbar,
                          w = t || Object(l.k)(a),
                          h = e || Object(l.n)(a),
                          f = Object(l.k)(s.target),
                          j = Object(l.n)(s.target),
                          g = s.key.replace(/(Arrow|ape)/, "");
                        switch (g) {
                          case "Tab":
                            r(s.shiftKey ? h : w);
                            break;
                          case "Left":
                            r(
                              h ||
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
                            p && d ? r(p.firstChild) : r(f, 0);
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
                value: function (t) {
                  var e = p("jw-settings-close", this.close, t.close, [
                    dt("close"),
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
                      Ze && Ze.open(t);
                    },
                    t.close,
                    [dt("arrow-left")]
                  );
                  return Object(l.m)(this.mainMenu.topbar, e.element()), e;
                },
              },
              {
                key: "createTopbar",
                value: function () {
                  var t = Object(l.e)('<div class="jw-submenu-topbar"></div>');
                  return Object(l.m)(this.el, t), t;
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
                      var s, l;
                      switch (a) {
                        case "quality":
                          s =
                            "Auto" === t.label && 0 === r
                              ? "".concat(
                                  i.defaultText,
                                  '&nbsp;<span class="jw-reset jw-auto-label"></span>'
                                )
                              : t.label;
                          break;
                        case "captions":
                          s =
                            ("Off" !== t.label && "off" !== t.id) || 0 !== r
                              ? t.label
                              : i.defaultText;
                          break;
                        case "playbackRates":
                          (l = t),
                            (s = Object(Ie.e)(i.tooltipText)
                              ? "x" + t
                              : t + "x");
                          break;
                        case "audioTracks":
                          s = t.name;
                      }
                      s || ((s = t), "object" === tn(t) && (s.options = i));
                      var c = new o(
                        s,
                        function (t) {
                          c.active ||
                            (e(l || r),
                            c.deactivate &&
                              (n.items
                                .filter(function (t) {
                                  return !0 === t.active;
                                })
                                .forEach(function (t) {
                                  t.deactivate();
                                }),
                              Ze ? Ze.open(t) : n.mainMenu.close(t)),
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
                      Object(l.h)(this.itemsContainer.el),
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
                    if (((Ze = null), this.isSubmenu)) {
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
                          (Ze = this.parentMenu),
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
        sn = function (t) {
          var e = t.closeButton,
            n = t.el;
          return new u.a(n).on("keydown", function (n) {
            var i = n.sourceEvent,
              o = n.target,
              a = Object(l.k)(o),
              r = Object(l.n)(o),
              s = i.key.replace(/(Arrow|ape)/, ""),
              c = function (e) {
                r ? e || r.focus() : t.close(n);
              };
            switch (s) {
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
                  if ((!e && Ze && (e = Ze.children[Ze.openMenus]), e))
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
        cn = function (t) {
          return hn[t];
        },
        un = function (t) {
          for (var e, n = Object.keys(hn), i = 0; i < n.length; i++)
            if (hn[n[i]] === t) {
              e = n[i];
              break;
            }
          return e;
        },
        dn = function (t) {
          return t + "%";
        },
        pn = function (t) {
          return parseInt(t);
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
        hn = {
          White: "#ffffff",
          Black: "#000000",
          Red: "#ff0000",
          Green: "#00ff00",
          Blue: "#0000ff",
          Yellow: "#ffff00",
          Magenta: "ff00ff",
          Cyan: "#00ffff",
        },
        fn = function (t, e, n, i) {
          var o = new rn("settings", null, i),
            a = function (t, e, a, r, s) {
              var l = n.elements["".concat(t, "Button")];
              if (!e || e.length <= 1)
                return o.removeMenu(t), void (l && l.hide());
              var c = o.children[t];
              c || (c = new rn(t, o, i)),
                c.setMenuItems(c.createItems(e, a, s), r),
                l && l.show();
            },
            r = function (r) {
              var s = { defaultText: i.auto };
              a(
                "quality",
                r,
                function (e) {
                  return t.setCurrentQuality(e);
                },
                e.get("currentLevel") || 0,
                s
              );
              var l = o.children,
                c = !!l.quality || l.playbackRates || Object.keys(l).length > 1;
              n.elements.settingsButton.toggle(c);
            };
          e.change(
            "levels",
            function (t, e) {
              r(e);
            },
            o
          );
          var s = function (t, n, i) {
            var o = e.get("levels");
            if (o && "Auto" === o[0].label && n && n.items.length) {
              var a = n.items[0].el.querySelector(".jw-auto-label"),
                r = o[t.index] || { label: "" };
              a.textContent = i ? "" : r.label;
            }
          };
          e.on("change:visualQuality", function (t, n) {
            var i = o.children.quality;
            n && i && s(n.level, i, e.get("currentLevel"));
          }),
            e.on(
              "change:currentLevel",
              function (t, n) {
                var i = o.children.quality,
                  a = e.get("visualQuality");
                a && i && s(a.level, i, n);
              },
              o
            ),
            e.change("captionsList", function (n, r) {
              var s = { defaultText: i.off },
                l = e.get("captionsIndex");
              a(
                "captions",
                r,
                function (e) {
                  return t.setCurrentCaptions(e);
                },
                l,
                s
              );
              var c = o.children.captions;
              if (c && !c.children.captionsSettings) {
                c.topbar = c.topbar || c.createTopbar();
                var u = new rn("captionsSettings", c, i);
                u.title = "Subtitle Settings";
                var d = new Ke("Settings", u.open);
                c.topbar.appendChild(d.el);
                var p = new Je("Reset", function () {
                  e.set("captions", ln.a), f();
                });
                p.el.classList.add("jw-settings-reset");
                var h = e.get("captions"),
                  f = function () {
                    var t = [];
                    wn.forEach(function (n) {
                      h &&
                        h[n.propertyName] &&
                        (n.defaultVal = n.getOption(h[n.propertyName]));
                      var o = new rn(n.name, u, i),
                        a = new Ke(
                          { label: n.name, value: n.defaultVal },
                          o.open,
                          He
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
                                s = Object(w.g)({}, i);
                              (s[o] = r), e.set("captions", s);
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
                f();
              }
            });
          var l = function (t, e) {
            t && e > -1 && t.items[e].activate();
          };
          e.change(
            "captionsIndex",
            function (t, e) {
              var i = o.children.captions;
              i && l(i, e), n.toggleCaptionsButtonState(!!e);
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
                s = { tooltipText: i.playbackRates };
              a(
                "playbackRates",
                n,
                function (e) {
                  return t.setPlaybackRate(e);
                },
                r,
                s
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
                i && (a = i.indexOf(n)), l(o.children.playbackRates, a);
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
        jn = n(58),
        gn = n(35),
        bn = n(12),
        mn = function (t, e, n, i) {
          var o = Object(l.e)(
              '<div class="jw-reset jw-info-overlay jw-modal"><div class="jw-reset jw-info-container"><div class="jw-reset-text jw-info-title" dir="auto"></div><div class="jw-reset-text jw-info-duration" dir="auto"></div><div class="jw-reset-text jw-info-description" dir="auto"></div></div><div class="jw-reset jw-info-clientid"></div></div>'
            ),
            r = !1,
            s = null,
            c = !1,
            u = function (t) {
              /jw-info/.test(t.target.className) || w.close();
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
                  e.get("localization").close,
                  [dt("close")]
                );
              d.show(),
                Object(l.m)(o, d.element()),
                (a = o.querySelector(".jw-info-title")),
                (s = o.querySelector(".jw-info-duration")),
                (c = o.querySelector(".jw-info-description")),
                (u = o.querySelector(".jw-info-clientid")),
                e.change("playlistItem", function (t, e) {
                  var n = e.description,
                    i = e.title;
                  Object(l.q)(c, n || ""), Object(l.q)(a, i || "Unknown Title");
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
                          } catch (t) {
                            return "none";
                          }
                        })()
                      )),
                t.appendChild(o),
                (r = !0);
            };
          var w = {
            open: function () {
              r || d(), document.addEventListener("click", u), (c = !0);
              var t = e.get("state");
              t === a.pb && n.pause("infoOverlayInteraction"), (s = t), i(!0);
            },
            close: function () {
              document.removeEventListener("click", u),
                (c = !1),
                e.get("state") === a.ob &&
                  s === a.pb &&
                  n.play("infoOverlayInteraction"),
                (s = null),
                i(!1);
            },
            destroy: function () {
              this.close(), e.off(null, null, this);
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
      var vn = function (t, e, n) {
          var i,
            o = !1,
            r = null,
            s = n.get("localization").shortcuts,
            c = Object(l.e)(
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
                    s = t.seekForward,
                    l = t.seekBackward;
                  return [
                    { key: t.spacebar, description: e },
                    { key: "↑", description: a },
                    { key: "↓", description: r },
                    { key: "→", description: s },
                    { key: "←", description: l },
                    { key: "c", description: t.captionsToggle },
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
            h = function () {
              w.el.setAttribute("aria-checked", n.get("enableShortcuts")),
                Object(l.a)(c, "jw-open"),
                (r = n.get("state")),
                c.querySelector(".jw-shortcuts-close").focus(),
                document.addEventListener("click", j),
                (o = !0),
                e.pause(d);
            },
            f = function () {
              Object(l.o)(c, "jw-open"),
                document.removeEventListener("click", j),
                t.focus(),
                (o = !1),
                r === a.pb && e.play(d);
            },
            j = function (t) {
              /jw-shortcuts|jw-switch/.test(t.target.className) || f();
            },
            g = function (t) {
              var e = t.currentTarget,
                i = "true" !== e.getAttribute("aria-checked");
              e.setAttribute("aria-checked", i), n.set("enableShortcuts", i);
            };
          return (
            (i = p("jw-shortcuts-close", f, n.get("localization").close, [
              dt("close"),
            ])),
            Object(l.m)(c, i.element()),
            i.show(),
            t.appendChild(c),
            w.on("click tap enter", g),
            {
              el: c,
              open: h,
              close: f,
              destroy: function () {
                f(), w.destroy();
              },
              toggleVisibility: function () {
                o ? f() : h();
              },
            }
          );
        },
        yn = function (t) {
          return (
            '<div class="jw-float-icon jw-icon jw-button-color jw-reset" aria-label='.concat(
              t,
              ' tabindex="0">'
            ) + "</div>"
          );
        };
      function kn(t) {
        return (kn =
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
      function xn(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function On(t, e) {
        return !e || ("object" !== kn(e) && "function" != typeof e)
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
      function Mn(t, e) {
        return (Mn =
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
            ((i = On(this, Cn(e).call(this))).element = Object(l.e)(yn(n))),
            i.element.appendChild(dt("close")),
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
              e && Mn(t, e);
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
          ]) && xn(n.prototype, i),
          o && xn(n, o),
          e
        );
      })(r.a);
      function Sn(t) {
        return (Sn =
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
      function _n(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function En(t, e) {
        return !e || ("object" !== Sn(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function zn(t) {
        return (zn = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function Pn(t, e) {
        return (Pn =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      n.d(e, "default", function () {
        return Ln;
      }),
        n(95);
      var An = o.OS.mobile ? 4e3 : 2e3,
        In = [27];
      (gn.a.cloneIcon = dt),
        bn.a.forEach(function (t) {
          if (t.getState() === a.lb) {
            var e = t.getContainer().querySelector(".jw-error-msg .jw-icon");
            e && !e.hasChildNodes() && e.appendChild(gn.a.cloneIcon("error"));
          }
        });
      var Rn = function () {
          return { reason: "interaction" };
        },
        Ln = (function (t) {
          function e(t, n) {
            var i;
            return (
              (function (t, e) {
                if (!(t instanceof e))
                  throw new TypeError("Cannot call a class as a function");
              })(this, e),
              ((i = En(this, zn(e).call(this))).activeTimeout = -1),
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
                e && Pn(t, e);
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
                    var d = new Ce(e, t);
                    d.buttons.display.on("click tap enter", function () {
                      n.trigger(a.p),
                        n.userActive(1e3),
                        t.playToggle(Rn()),
                        u();
                    }),
                      this.div.appendChild(d.element()),
                      (this.displayContainer = d);
                  }
                  (this.infoOverlay = new mn(i, e, t, function (t) {
                    Object(l.v)(n.div, "jw-info-open", t),
                      t && n.div.querySelector(".jw-info-close").focus();
                  })),
                    o.OS.mobile ||
                      (this.shortcutsTooltip = new vn(
                        this.wrapperElement,
                        t,
                        e
                      )),
                    (this.rightClickMenu = new Ve(
                      this.infoOverlay,
                      this.shortcutsTooltip
                    )),
                    c
                      ? (Object(l.a)(this.playerContainer, "jw-flag-touch"),
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
                  var w = e.get("floating");
                  if (w) {
                    var h = new Tn(i, e.get("localization").close);
                    h.on(a.sb, function () {
                      return n.trigger("dismissFloating", { doNotForward: !0 });
                    }),
                      !1 !== w.dismissible &&
                        Object(l.a)(
                          this.playerContainer,
                          "jw-floating-dismissible"
                        );
                  }
                  var f = (this.controlbar = new de(
                    t,
                    e,
                    this.playerContainer.querySelector(
                      ".jw-hidden-accessibility"
                    )
                  ));
                  if (
                    (f.on(a.sb, function () {
                      return n.userActive();
                    }),
                    f.on(
                      "nextShown",
                      function (t) {
                        this.trigger("nextShown", t);
                      },
                      this
                    ),
                    f.on("adjustVolume", k, this),
                    e.get("nextUpDisplay") && !f.nextUpToolTip)
                  ) {
                    var j = new _e(e, t, this.playerContainer);
                    j.on("all", this.trigger, this),
                      j.setup(this.context),
                      (f.nextUpToolTip = j),
                      this.div.appendChild(j.element());
                  }
                  this.div.appendChild(f.element());
                  var g = e.get("localization"),
                    b = (this.settingsMenu = fn(
                      t,
                      e.player,
                      this.controlbar,
                      g
                    )),
                    m = null;
                  this.controlbar.on("menuVisibility", function (i) {
                    var o = i.visible,
                      r = i.evt,
                      s = e.get("state"),
                      l = { reason: "settingsInteraction" },
                      c = n.controlbar.elements.settingsButton,
                      d = "keydown" === ((r && r.sourceEvent) || r || {}).type,
                      p = o || d ? 0 : An;
                    n.userActive(p),
                      (m = s),
                      Object(jn.a)(e.get("containerWidth")) < 2 &&
                        (o && s === a.pb
                          ? t.pause(l)
                          : o || s !== a.ob || m !== a.pb || t.play(l)),
                      !o && d && c ? c.element().focus() : r && u();
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
                        this.div.insertBefore(b.el, f.element()));
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
                          [dt("volume-0")]
                        )),
                        n.mute.show(),
                        n.div.appendChild(n.mute.element())),
                        f.renderVolume(!0, e.get("volume")),
                        Object(l.a)(n.playerContainer, "jw-flag-autostart"),
                        e.on("change:autostartFailed", i, n),
                        e.on("change:autostartMuted change:mute", a, n),
                        (n.muteChangeCallback = a),
                        (n.unmuteCallback = i);
                    }
                  };
                  function y(n) {
                    var i = 0,
                      o = e.get("duration"),
                      a = e.get("position");
                    if ("DVR" === e.get("streamType")) {
                      var r = e.get("dvrSeekLimit");
                      (i = o), (o = Math.max(a, -r));
                    }
                    var l = Object(s.a)(a + n, i, o);
                    t.seek(l, Rn());
                  }
                  function k(n) {
                    var i = Object(s.a)(e.get("volume") + n, 0, 100);
                    t.setVolume(i);
                  }
                  e.once("change:autostartMuted", v), v(e);
                  var x = function (i) {
                    if (i.ctrlKey || i.metaKey) return !0;
                    var o = !n.settingsMenu.visible,
                      a = !0 === e.get("enableShortcuts"),
                      r = n.instreamState;
                    if (a || -1 !== In.indexOf(i.keyCode)) {
                      switch (i.keyCode) {
                        case 27:
                          if (e.get("fullscreen"))
                            t.setFullscreen(!1),
                              n.playerContainer.blur(),
                              n.userInactive();
                          else {
                            var s = t.getPlugin("related");
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
                          t.playToggle(Rn());
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
                          var l = t.getCaptionsList().length;
                          if (l) {
                            var c = (t.getCurrentCaptions() + 1) % l;
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
                            t.seek(u, Rn());
                          }
                      }
                      return /13|32|37|38|39|40/.test(i.keyCode)
                        ? (i.preventDefault(), !1)
                        : void 0;
                    }
                  };
                  this.playerContainer.addEventListener("keydown", x),
                    (this.keydownCallback = x);
                  var O = function (t) {
                    switch (t.keyCode) {
                      case 9:
                        var e = n.playerContainer.contains(t.target) ? 0 : An;
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
                  var M = function t() {
                    "jw-shortcuts-tooltip-explanation" ===
                      n.playerContainer.getAttribute("aria-describedby") &&
                      n.playerContainer.removeAttribute("aria-describedby"),
                      n.playerContainer.removeEventListener("blur", t, !0);
                  };
                  this.shortcutsTooltip &&
                    (this.playerContainer.addEventListener("blur", M, !0),
                    (this.onRemoveShortcutsDescription = M)),
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
                    s = this.playerContainer,
                    c = this.div;
                  clearTimeout(this.activeTimeout),
                    (this.activeTimeout = -1),
                    this.off(),
                    t.off(null, null, this),
                    t.set("controlsEnabled", !1),
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
                    Object(l.o)(this.playerContainer, "jw-flag-autostart"),
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
                    o = e || n || i ? 0 : An;
                  this.userActive(o);
                },
              },
              {
                key: "userActive",
                value: function () {
                  var t =
                    arguments.length > 0 && void 0 !== arguments[0]
                      ? arguments[0]
                      : An;
                  t > 0
                    ? ((this.inactiveTime = Object(c.a)() + t),
                      -1 === this.activeTimeout &&
                        (this.activeTimeout = setTimeout(
                          this.userInactiveTimeout,
                          t
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
                    Object(l.o)(this.playerContainer, "jw-flag-autostart"),
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
                      Object(l.a)(this.playerContainer, "jw-flag-autostart"),
                    this.controlbar.elements.time
                      .element()
                      .setAttribute("tabindex", "0");
                },
              },
            ]) && _n(n.prototype, i),
            r && _n(n, r),
            e
          );
        })(r.a);
    },
    function (t, e, n) {
      "use strict";
      n.r(e);
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
        h = n(7),
        f = n(34);
      function j(t) {
        var e = !1;
        return {
          async: function () {
            var n = this,
              i = arguments;
            return Promise.resolve().then(function () {
              if (!e) return t.apply(n, i);
            });
          },
          cancel: function () {
            e = !0;
          },
          cancelled: function () {
            return e;
          },
        };
      }
      var g = n(1);
      function b(t) {
        return function (e, n) {
          var o = t.mediaModel,
            a = Object(i.g)({}, n, { type: e });
          switch (e) {
            case d.T:
              if (o.get(d.T) === n.mediaType) return;
              o.set(d.T, n.mediaType);
              break;
            case d.U:
              return void o.set(d.U, Object(i.g)({}, n));
            case d.M:
              if (n[e] === t.model.getMute()) return;
              break;
            case d.bb:
              n.newstate === d.mb && (t.thenPlayPromise.cancel(), o.srcReset());
              var r = o.attributes.mediaState;
              (o.attributes.mediaState = n.newstate),
                o.trigger("change:mediaState", o, n.newstate, r);
              break;
            case d.F:
              return (
                (t.beforeComplete = !0),
                t.trigger(d.B, a),
                void (t.attached && !t.background && t._playbackComplete())
              );
            case d.G:
              o.get("setup")
                ? (t.thenPlayPromise.cancel(), o.srcReset())
                : ((e = d.tb), (a.code += 1e5));
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
                e === d.S &&
                  Object(i.r)(t.item.starttime) &&
                  delete t.item.starttime;
              break;
            case d.R:
              var c = t.mediaElement;
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
                h = n.tracks;
              w > -1 &&
                h.length > 0 &&
                w < h.length &&
                o.set("currentAudioTrack", parseInt(w));
          }
          t.trigger(e, a);
        };
      }
      var m = n(8),
        v = n(45),
        y = n(41);
      function k(t) {
        return (k =
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
      function x(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function O(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function C(t, e, n) {
        return e && O(t.prototype, e), n && O(t, n), t;
      }
      function M(t, e) {
        return !e || ("object" !== k(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function T(t) {
        return (T = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function S(t, e) {
        if ("function" != typeof e && null !== e)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (t.prototype = Object.create(e && e.prototype, {
          constructor: { value: t, writable: !0, configurable: !0 },
        })),
          e && _(t, e);
      }
      function _(t, e) {
        return (_ =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var E = (function (t) {
          function e() {
            var t;
            return (
              x(this, e),
              ((t = M(this, T(e).call(this))).providerController = null),
              (t._provider = null),
              t.addAttributes({ mediaModel: new P() }),
              t
            );
          }
          return (
            S(e, t),
            C(e, [
              {
                key: "setup",
                value: function (t) {
                  return (
                    (t = t || {}),
                    this._normalizeConfig(t),
                    Object(i.g)(this.attributes, t, y.b),
                    (this.providerController = new f.a(
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
                  var t = this.clone(),
                    e = t.mediaModel.attributes;
                  return (
                    Object.keys(y.a).forEach(function (n) {
                      t[n] = e[n];
                    }),
                    (t.instreamMode = !!t.instream),
                    delete t.instream,
                    delete t.mediaModel,
                    t
                  );
                },
              },
              {
                key: "persistQualityLevel",
                value: function (t, e) {
                  var n = e[t] || {},
                    o = n.label,
                    a = Object(i.u)(n.bitrate) ? n.bitrate : null;
                  this.set("bitrateSelection", a), this.set("qualityLabel", o);
                },
              },
              {
                key: "setActiveItem",
                value: function (t) {
                  var e = this.get("playlist")[t];
                  this.resetItem(e),
                    (this.attributes.playlistItem = null),
                    this.set("item", t),
                    this.set("minDvrWindow", e.minDvrWindow),
                    this.set("dvrSeekLimit", e.dvrSeekLimit),
                    this.set("playlistItem", e);
                },
              },
              {
                key: "setMediaModel",
                value: function (t) {
                  this.mediaModel &&
                    this.mediaModel !== t &&
                    this.mediaModel.off(),
                    (t = t || new P()),
                    this.set("mediaModel", t),
                    (function (t) {
                      var e = t.get("mediaState");
                      t.trigger("change:mediaState", t, e, e);
                    })(t);
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
                value: function (t) {
                  (t = !!t) !== this.get("fullscreen") &&
                    this.set("fullscreen", t);
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
                value: function (t) {
                  if (Object(i.u)(t)) {
                    var e = Math.min(Math.max(0, t), 100);
                    this.set("volume", e);
                    var n = 0 === e;
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
                value: function (t) {
                  if (
                    (void 0 === t && (t = !this.getMute()),
                    this.set("mute", !!t),
                    !t)
                  ) {
                    var e = Math.max(10, this.get("volume"));
                    this.set("autostartMuted", !1), this.setVolume(e);
                  }
                },
              },
              {
                key: "setStreamType",
                value: function (t) {
                  this.set("streamType", t),
                    "LIVE" === t && this.setPlaybackRate(1);
                },
              },
              {
                key: "setProvider",
                value: function (t) {
                  (this._provider = t), z(this, t);
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
                value: function (t) {
                  Object(i.r)(t) &&
                    ((t = Math.max(Math.min(t, 4), 0.25)),
                    "LIVE" === this.get("streamType") && (t = 1),
                    this.set("defaultPlaybackRate", t),
                    this._provider &&
                      this._provider.setPlaybackRate &&
                      this._provider.setPlaybackRate(t));
                },
              },
              {
                key: "persistCaptionsTrack",
                value: function () {
                  var t = this.get("captionsTrack");
                  t
                    ? this.set("captionLabel", t.name)
                    : this.set("captionLabel", "Off");
                },
              },
              {
                key: "setVideoSubtitleTrack",
                value: function (t, e) {
                  this.set("captionsIndex", t),
                    t &&
                      e &&
                      t <= e.length &&
                      e[t - 1].data &&
                      this.set("captionsTrack", e[t - 1]);
                },
              },
              {
                key: "persistVideoSubtitleTrack",
                value: function (t, e) {
                  this.setVideoSubtitleTrack(t, e), this.persistCaptionsTrack();
                },
              },
              {
                key: "setAutoStart",
                value: function (t) {
                  void 0 !== t && this.set("autostart", t);
                  var e = m.OS.mobile && this.get("autostart");
                  this.set(
                    "playOnViewable",
                    e || "viewable" === this.get("autostart")
                  );
                },
              },
              {
                key: "resetItem",
                value: function (t) {
                  var e = t ? Object(w.g)(t.starttime) : 0,
                    n = t ? Object(w.g)(t.duration) : 0,
                    i = this.mediaModel;
                  this.set("playRejected", !1),
                    (this.attributes.itemMeta = {}),
                    i.set("position", e),
                    i.set("currentTime", 0),
                    i.set("duration", n);
                },
              },
              {
                key: "persistBandwidthEstimate",
                value: function (t) {
                  Object(i.u)(t) && this.set("bandwidthEstimate", t);
                },
              },
              {
                key: "_normalizeConfig",
                value: function (t) {
                  var e = t.floating;
                  e && e.disabled && delete t.floating;
                },
              },
            ]),
            e
          );
        })(v.a),
        z = function (t, e) {
          t.set("provider", e.getName()),
            !0 === t.get("instreamMode") && (e.instreamMode = !0),
            -1 === e.getName().name.indexOf("flash") &&
              (t.set("flashThrottle", void 0), t.set("flashBlocked", !1)),
            t.setPlaybackRate(t.get("defaultPlaybackRate")),
            t.set("supportsPlaybackRate", e.supportsPlaybackRate),
            t.set("playbackRate", e.getPlaybackRate()),
            t.set("renderCaptionsNatively", e.renderNatively);
        };
      var P = (function (t) {
          function e() {
            var t;
            return (
              x(this, e),
              (t = M(this, T(e).call(this))).addAttributes({
                mediaState: d.mb,
              }),
              t
            );
          }
          return (
            S(e, t),
            C(e, [
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
            e
          );
        })(v.a),
        A = E;
      function I(t) {
        return (I =
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
      function R(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function L(t) {
        return (L = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function B(t, e) {
        return (B =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      function V(t) {
        if (void 0 === t)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return t;
      }
      var N = (function (t) {
        function e(t, n) {
          var i, o, a, r;
          return (
            (function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
            (o = this),
            (a = L(e).call(this)),
            ((i =
              !a || ("object" !== I(a) && "function" != typeof a)
                ? V(o)
                : a).attached = !0),
            (i.beforeComplete = !1),
            (i.item = null),
            (i.mediaModel = new P()),
            (i.model = n),
            (i.provider = t),
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
          (function (t, e) {
            if ("function" != typeof e && null !== e)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (t.prototype = Object.create(e && e.prototype, {
              constructor: { value: t, writable: !0, configurable: !0 },
            })),
              e && B(t, e);
          })(e, t),
          (n = e),
          (o = [
            {
              key: "play",
              value: function (t) {
                var e = this.item,
                  n = this.model,
                  i = this.mediaModel,
                  o = this.provider;
                if (
                  (t || (t = n.get("playReason")),
                  n.set("playRejected", !1),
                  i.get("setup"))
                )
                  return o.play() || Promise.resolve();
                i.set("setup", !0);
                var a = this._loadAndPlay(e, o);
                return i.get("started") ? a : this._playAttempt(a, t);
              },
            },
            {
              key: "stop",
              value: function () {
                var t = this.provider;
                (this.beforeComplete = !1), t.stop();
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
                var t = this.item,
                  e = this.mediaModel,
                  n = this.provider;
                !t ||
                  (t && "none" === t.preload) ||
                  !this.attached ||
                  this.setup ||
                  this.preloaded ||
                  (e.set("preloaded", !0), n.preload(t));
              },
            },
            {
              key: "destroy",
              value: function () {
                var t = this.provider,
                  e = this.mediaModel;
                this.off(),
                  e.off(),
                  t.off(),
                  this.eventQueue.destroy(),
                  this.detach(),
                  t.getContainer() && t.remove(),
                  delete t.instreamMode,
                  (this.provider = null),
                  (this.item = null);
              },
            },
            {
              key: "attach",
              value: function () {
                var t = this.model,
                  e = this.provider;
                t.setPlaybackRate(t.get("defaultPlaybackRate")),
                  e.attachMedia(),
                  (this.attached = !0),
                  this.eventQueue.flush(),
                  this.beforeComplete && this._playbackComplete();
              },
            },
            {
              key: "detach",
              value: function () {
                var t = this.provider;
                this.thenPlayPromise.cancel();
                var e = t.detachMedia();
                return (this.attached = !1), e;
              },
            },
            {
              key: "_playAttempt",
              value: function (t, e) {
                var n = this,
                  o = this.item,
                  a = this.mediaModel,
                  r = this.model,
                  s = this.provider,
                  l = s ? s.video : null;
                return (
                  this.trigger(d.N, { item: o, playReason: e }),
                  (l ? l.paused : r.get(d.bb) !== d.pb) || r.set(d.bb, d.jb),
                  t
                    .then(function () {
                      a.get("setup") &&
                        (a.set("started", !0),
                        a === r.mediaModel &&
                          (function (t) {
                            var e = t.get("mediaState");
                            t.trigger("change:mediaState", t, e, e);
                          })(a));
                    })
                    .catch(function (t) {
                      if (n.item && a === r.mediaModel) {
                        if ((r.set("playRejected", !0), l && l.paused)) {
                          if (l.src === location.href)
                            return n._loadAndPlay(o, s);
                          a.set("mediaState", d.ob);
                        }
                        var c = Object(i.g)(new g.n(null, Object(g.w)(t), t), {
                          error: t,
                          item: o,
                          playReason: e,
                        });
                        throw (delete c.key, n.trigger(d.O, c), t);
                      }
                    })
                );
              },
            },
            {
              key: "_playbackComplete",
              value: function () {
                var t = this.item,
                  e = this.provider;
                t && delete t.starttime,
                  (this.beforeComplete = !1),
                  e.setState(d.kb),
                  this.trigger(d.F, {});
              },
            },
            {
              key: "_loadAndPlay",
              value: function () {
                var t = this.item,
                  e = this.provider,
                  n = e.load(t);
                if (n) {
                  var i = j(function () {
                    return e.play() || Promise.resolve();
                  });
                  return (this.thenPlayPromise = i), n.then(i.async);
                }
                return e.play() || Promise.resolve();
              },
            },
            {
              key: "audioTrack",
              get: function () {
                return this.provider.getCurrentAudioTrack();
              },
              set: function (t) {
                this.provider.setCurrentAudioTrack(t);
              },
            },
            {
              key: "quality",
              get: function () {
                return this.provider.getCurrentQuality();
              },
              set: function (t) {
                this.provider.setCurrentQuality(t);
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
                var t = this.container,
                  e = this.provider;
                return (
                  !!this.attached &&
                  !!e.video &&
                  (!t || (t && !t.contains(e.video)))
                );
              },
              set: function (t) {
                var e = this.container,
                  n = this.provider;
                n.video
                  ? e &&
                    (t
                      ? this.background ||
                        (this.thenPlayPromise.cancel(),
                        this.pause(),
                        e.removeChild(n.video),
                        (this.container = null))
                      : (this.eventQueue.flush(),
                        this.beforeComplete && this._playbackComplete()))
                  : t
                  ? this.detach()
                  : this.attach();
              },
            },
            {
              key: "container",
              get: function () {
                return this.provider.getContainer();
              },
              set: function (t) {
                this.provider.setContainer(t);
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
              set: function (t) {
                var e = (this.mediaModel = new P()),
                  n = t ? Object(w.g)(t.starttime) : 0,
                  i = t ? Object(w.g)(t.duration) : 0,
                  o = e.attributes;
                e.srcReset(),
                  (o.position = n),
                  (o.duration = i),
                  (this.item = t),
                  this.provider.init(t);
              },
            },
            {
              key: "controls",
              set: function (t) {
                this.provider.setControls(t);
              },
            },
            {
              key: "mute",
              set: function (t) {
                this.provider.mute(t);
              },
            },
            {
              key: "position",
              set: function (t) {
                var e = this.provider;
                this.model.get("scrubbing") && e.fastSeek
                  ? e.fastSeek(t)
                  : e.seek(t);
              },
            },
            {
              key: "subtitles",
              set: function (t) {
                this.provider.setSubtitlesTrack &&
                  this.provider.setSubtitlesTrack(t);
              },
            },
            {
              key: "volume",
              set: function (t) {
                this.provider.volume(t);
              },
            },
          ]) && R(n.prototype, o),
          a && R(n, a),
          e
        );
      })(h.a);
      function H(t) {
        return (H =
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
      function F(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function q(t) {
        return (q = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function D(t, e) {
        return (D =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      function U(t) {
        if (void 0 === t)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return t;
      }
      function W(t, e) {
        var n = e.mediaControllerListener;
        t.off().on("all", n, e);
      }
      function Q(t) {
        return t && t.sources && t.sources[0];
      }
      var Y = (function (t) {
        function e(t, n) {
          var o, a, r, s, l;
          return (
            (function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
            (a = this),
            ((o =
              !(r = q(e).call(this)) ||
              ("object" !== H(r) && "function" != typeof r)
                ? U(a)
                : r).adPlaying = !1),
            (o.background =
              ((s = null),
              (l = null),
              Object.defineProperties(
                {
                  setNext: function (t, e) {
                    l = { item: t, loadPromise: e };
                  },
                  isNext: function (t) {
                    return !(
                      !l ||
                      JSON.stringify(l.item.sources[0]) !==
                        JSON.stringify(t.sources[0])
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
                    set: function (t) {
                      s = t;
                    },
                  },
                }
              ))),
            (o.mediaPool = n),
            (o.mediaController = null),
            (o.mediaControllerListener = (function (t, e) {
              return function (n, o) {
                switch (n) {
                  case d.bb:
                    return;
                  case "flashThrottle":
                  case "flashBlocked":
                    return void t.set(n, o.value);
                  case d.V:
                  case d.M:
                    return void t.set(n, o[n]);
                  case d.P:
                    return void t.set("playbackRate", o.playbackRate);
                  case d.K:
                    Object(i.g)(t.get("itemMeta"), o.metadata);
                    break;
                  case d.J:
                    t.persistQualityLevel(o.currentQuality, o.levels);
                    break;
                  case "subtitlesTrackChanged":
                    t.persistVideoSubtitleTrack(o.currentTrack, o.tracks);
                    break;
                  case d.S:
                  case d.Q:
                  case d.R:
                  case d.X:
                  case "subtitlesTracks":
                  case "subtitlesTracksData":
                    t.trigger(n, o);
                    break;
                  case d.i:
                    return void t.persistBandwidthEstimate(o.bandwidthEstimate);
                }
                e.trigger(n, o);
              };
            })(t, U(U(o)))),
            (o.model = t),
            (o.providers = new f.a(t.getConfiguration())),
            (o.loadPromise = Promise.resolve()),
            (o.backgroundLoading = t.get("backgroundLoading")),
            o.backgroundLoading ||
              t.set("mediaElement", o.mediaPool.getPrimedElement()),
            o
          );
        }
        var n, o, a;
        return (
          (function (t, e) {
            if ("function" != typeof e && null !== e)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (t.prototype = Object.create(e && e.prototype, {
              constructor: { value: t, writable: !0, configurable: !0 },
            })),
              e && D(t, e);
          })(e, t),
          (n = e),
          (o = [
            {
              key: "setActiveItem",
              value: function (t) {
                var e = this,
                  n = this.model,
                  i = n.get("playlist")[t];
                (n.attributes.itemReady = !1), n.setActiveItem(t);
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
                    .then(function (t) {
                      if (s === n.mediaModel)
                        return (t.activeItem = i), e._setActiveMedia(t), t;
                    })
                    .catch(function (t) {
                      throw (e._destroyActiveMedia(), t);
                    })),
                  this.loadPromise
                );
              },
            },
            {
              key: "setAttached",
              value: function (t) {
                var e = this.mediaController;
                if (((this.attached = t), e)) {
                  if (!t) {
                    var n = e.detach(),
                      i = e.item,
                      o = e.mediaModel.get("position");
                    return o && (i.starttime = o), n;
                  }
                  e.attach();
                }
              },
            },
            {
              key: "playVideo",
              value: function (t) {
                var e,
                  n = this,
                  i = this.mediaController,
                  o = this.model;
                if (!o.get("playlistItem"))
                  return Promise.reject(new Error("No media"));
                if ((t || (t = o.get("playReason")), i)) e = i.play(t);
                else {
                  o.set(d.bb, d.jb);
                  var a = j(function (e) {
                    if (
                      n.mediaController &&
                      n.mediaController.mediaModel === e.mediaModel
                    )
                      return e.play(t);
                    throw new Error("Playback cancelled.");
                  });
                  e = this.loadPromise
                    .catch(function (t) {
                      throw (a.cancel(), t);
                    })
                    .then(a.async);
                }
                return e;
              },
            },
            {
              key: "stopVideo",
              value: function () {
                var t = this.mediaController,
                  e = this.model,
                  n = e.get("playlist")[e.get("item")];
                (e.attributes.playlistItem = n), e.resetItem(n), t && t.stop();
              },
            },
            {
              key: "preloadVideo",
              value: function () {
                var t = this.background,
                  e = this.mediaController || t.currentMedia;
                e && e.preload();
              },
            },
            {
              key: "pause",
              value: function () {
                var t = this.mediaController;
                t && t.pause();
              },
            },
            {
              key: "castVideo",
              value: function (t, e) {
                var n = this.model;
                n.attributes.itemReady = !1;
                var o = Object(i.g)({}, e),
                  a = (o.starttime = n.mediaModel.get("currentTime"));
                this._destroyActiveMedia();
                var r = new N(t, n);
                (r.activeItem = o),
                  this._setActiveMedia(r),
                  n.mediaModel.set("currentTime", a);
              },
            },
            {
              key: "stopCast",
              value: function () {
                var t = this.model,
                  e = t.get("item");
                (t.get("playlist")[e].starttime = t.mediaModel.get(
                  "currentTime"
                )),
                  this.stopVideo(),
                  this.setActiveItem(e);
              },
            },
            {
              key: "backgroundActiveMedia",
              value: function () {
                this.adPlaying = !0;
                var t = this.background,
                  e = this.mediaController;
                e &&
                  (t.currentMedia &&
                    this._destroyMediaController(t.currentMedia),
                  (e.background = !0),
                  (t.currentMedia = e),
                  (this.mediaController = null));
              },
            },
            {
              key: "restoreBackgroundMedia",
              value: function () {
                this.adPlaying = !1;
                var t = this.background,
                  e = this.mediaController,
                  n = t.currentMedia;
                if (n) {
                  if (e)
                    return (
                      this._destroyMediaController(n),
                      void (t.currentMedia = null)
                    );
                  var i = n.mediaModel.attributes;
                  i.mediaState === d.mb
                    ? (i.mediaState = d.ob)
                    : i.mediaState !== d.ob && (i.mediaState = d.jb),
                    this._setActiveMedia(n),
                    (n.background = !1),
                    (t.currentMedia = null);
                }
              },
            },
            {
              key: "backgroundLoad",
              value: function (t) {
                var e = this.background,
                  n = Q(t);
                e.setNext(
                  t,
                  this._setupMediaController(n)
                    .then(function (e) {
                      return (e.activeItem = t), e.preload(), e;
                    })
                    .catch(function () {
                      e.clearNext();
                    })
                );
              },
            },
            {
              key: "forwardEvents",
              value: function () {
                var t = this.mediaController;
                t && W(t, this);
              },
            },
            {
              key: "routeEvents",
              value: function (t) {
                var e = this.mediaController;
                e && (e.off(), t && W(e, t));
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
              value: function (t) {
                var e = this.model,
                  n = t.mediaModel,
                  i = t.provider;
                !(function (t, e) {
                  var n = t.get("mediaContainer");
                  n
                    ? (e.container = n)
                    : t.once("change:mediaContainer", function (t, n) {
                        e.container = n;
                      });
                })(e, t),
                  (this.mediaController = t),
                  e.set("mediaElement", t.mediaElement),
                  e.setMediaModel(n),
                  e.setProvider(i),
                  W(t, this),
                  e.set("itemReady", !0);
              },
            },
            {
              key: "_destroyActiveMedia",
              value: function () {
                var t = this.mediaController,
                  e = this.model;
                t &&
                  (t.detach(),
                  this._destroyMediaController(t),
                  e.resetProvider(),
                  (this.mediaController = null));
              },
            },
            {
              key: "_destroyBackgroundMedia",
              value: function () {
                var t = this.background;
                this._destroyMediaController(t.currentMedia),
                  (t.currentMedia = null),
                  this._destroyBackgroundLoadingMedia();
              },
            },
            {
              key: "_destroyMediaController",
              value: function (t) {
                var e = this.mediaPool;
                t && (e.recycle(t.mediaElement), t.destroy());
              },
            },
            {
              key: "_setupMediaController",
              value: function (t) {
                var e = this,
                  n = this.model,
                  i = this.providers,
                  o = function (t) {
                    return new N(
                      new t(n.get("id"), n.getConfiguration(), e.primedElement),
                      n
                    );
                  },
                  a = i.choose(t),
                  r = a.provider,
                  s = a.name;
                return r
                  ? Promise.resolve(o(r))
                  : i.load(s).then(function (t) {
                      return o(t);
                    });
              },
            },
            {
              key: "_activateBackgroundMedia",
              value: function () {
                var t = this,
                  e = this.background,
                  n = this.background.nextLoadPromise,
                  i = this.model;
                return (
                  this._destroyMediaController(e.currentMedia),
                  (e.currentMedia = null),
                  n.then(function (n) {
                    if (n)
                      return (
                        e.clearNext(),
                        t.adPlaying
                          ? ((i.attributes.itemReady = !0),
                            (e.currentMedia = n))
                          : (t._setActiveMedia(n), (n.background = !1)),
                        n
                      );
                  })
                );
              },
            },
            {
              key: "_destroyBackgroundLoadingMedia",
              value: function () {
                var t = this,
                  e = this.background,
                  n = this.background.nextLoadPromise;
                n &&
                  n.then(function (n) {
                    t._destroyMediaController(n), e.clearNext();
                  });
              },
            },
            {
              key: "_providerCanPlay",
              value: function (t, e) {
                var n = this.providers.choose(e).provider;
                return n && t && t instanceof n;
              },
            },
            {
              key: "audioTrack",
              get: function () {
                var t = this.mediaController;
                return t ? t.audioTrack : -1;
              },
              set: function (t) {
                var e = this.mediaController;
                e && (e.audioTrack = parseInt(t, 10) || 0);
              },
            },
            {
              key: "audioTracks",
              get: function () {
                var t = this.mediaController;
                if (t) return t.audioTracks;
              },
            },
            {
              key: "beforeComplete",
              get: function () {
                var t = this.mediaController,
                  e = this.background.currentMedia;
                return !(!t && !e) && (t ? t.beforeComplete : e.beforeComplete);
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
              set: function (t) {
                var e = this.mediaController;
                e && (e.quality = parseInt(t, 10) || 0);
              },
            },
            {
              key: "qualities",
              get: function () {
                var t = this.mediaController;
                return t ? t.qualities : null;
              },
            },
            {
              key: "controls",
              set: function (t) {
                var e = this.mediaController;
                e && (e.controls = t);
              },
            },
            {
              key: "mute",
              set: function (t) {
                var e = this.background,
                  n = this.mediaController,
                  i = this.mediaPool;
                n && (n.mute = t),
                  e.currentMedia && (e.currentMedia.mute = t),
                  i.syncMute(t);
              },
            },
            {
              key: "position",
              set: function (t) {
                var e = this.mediaController;
                e && ((e.item.starttime = t), e.attached && (e.position = t));
              },
            },
            {
              key: "subtitles",
              set: function (t) {
                var e = this.mediaController;
                e && (e.subtitles = t);
              },
            },
            {
              key: "volume",
              set: function (t) {
                var e = this.background,
                  n = this.mediaController,
                  i = this.mediaPool;
                n && (n.volume = t),
                  e.currentMedia && (e.currentMedia.volume = t),
                  i.syncVolume(t);
              },
            },
          ]) && F(n.prototype, o),
          a && F(n, a),
          e
        );
      })(h.a);
      function X(t) {
        return t === d.kb || t === d.lb ? d.mb : t;
      }
      function Z(t, e, n) {
        if ((e = X(e)) !== (n = X(n))) {
          var i = e.replace(/(?:ing|d)$/, ""),
            o = {
              type: i,
              newstate: e,
              oldstate: n,
              reason: (function (t, e) {
                return t === d.jb ? (e === d.qb ? e : d.nb) : e;
              })(e, t.mediaModel.get("mediaState")),
            };
          "play" === i
            ? (o.playReason = t.get("playReason"))
            : "pause" === i && (o.pauseReason = t.get("pauseReason")),
            this.trigger(i, o);
        }
      }
      var K = n(48);
      function J(t) {
        return (J =
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
      function G(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function $(t, e) {
        return !e || ("object" !== J(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function tt(t, e, n, i) {
        return (tt =
          "undefined" != typeof Reflect && Reflect.set
            ? Reflect.set
            : function (t, e, n, i) {
                var o,
                  a = it(t, e);
                if (a) {
                  if ((o = Object.getOwnPropertyDescriptor(a, e)).set)
                    return o.set.call(i, n), !0;
                  if (!o.writable) return !1;
                }
                if ((o = Object.getOwnPropertyDescriptor(i, e))) {
                  if (!o.writable) return !1;
                  (o.value = n), Object.defineProperty(i, e, o);
                } else
                  !(function (t, e, n) {
                    e in t
                      ? Object.defineProperty(t, e, {
                          value: n,
                          enumerable: !0,
                          configurable: !0,
                          writable: !0,
                        })
                      : (t[e] = n);
                  })(i, e, n);
                return !0;
              })(t, e, n, i);
      }
      function et(t, e, n, i, o) {
        if (!tt(t, e, n, i || t) && o)
          throw new Error("failed to set property");
        return n;
      }
      function nt(t, e, n) {
        return (nt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, n) {
                var i = it(t, e);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, e);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(t, e, n || t);
      }
      function it(t, e) {
        for (
          ;
          !Object.prototype.hasOwnProperty.call(t, e) && null !== (t = ot(t));

        );
        return t;
      }
      function ot(t) {
        return (ot = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function at(t, e) {
        return (at =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var rt = (function (t) {
          function e(t, n) {
            var i;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e);
            var o,
              a = ((i = $(this, ot(e).call(this, t, n))).model = new A());
            if (
              ((i.playerModel = t),
              (i.provider = null),
              (i.backgroundLoading = t.get("backgroundLoading")),
              (a.mediaModel.attributes.mediaType = "video"),
              i.backgroundLoading)
            )
              o = n.getAdElement();
            else {
              (o = t.get("mediaElement")),
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
            (function (t, e) {
              if ("function" != typeof e && null !== e)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (t.prototype = Object.create(e && e.prototype, {
                constructor: { value: t, writable: !0, configurable: !0 },
              })),
                e && at(t, e);
            })(e, t),
            (n = e),
            (o = [
              {
                key: "setup",
                value: function () {
                  var t = this.model,
                    e = this.playerModel,
                    n = this.primedElement,
                    i = e.attributes,
                    o = e.mediaModel;
                  t.setup({
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
                    t.on("change:state", Z, this),
                    t.on(
                      d.w,
                      function (t) {
                        this.trigger(d.w, t);
                      },
                      this
                    ),
                    n.paused || n.pause();
                },
              },
              {
                key: "setActiveItem",
                value: function (t) {
                  var n = this;
                  return (
                    this.stopVideo(),
                    (this.provider = null),
                    nt(ot(e.prototype), "setActiveItem", this)
                      .call(this, t)
                      .then(function (t) {
                        n._setProvider(t.provider);
                      }),
                    this.playVideo()
                  );
                },
              },
              {
                key: "usePsuedoProvider",
                value: function (t) {
                  (this.provider = t),
                    t &&
                      (this._setProvider(t),
                      t.off(d.w),
                      t.on(
                        d.w,
                        function (t) {
                          this.trigger(d.w, t);
                        },
                        this
                      ));
                },
              },
              {
                key: "_setProvider",
                value: function (t) {
                  var e = this;
                  if (t && this.mediaPool) {
                    var n = this.model,
                      o = this.playerModel,
                      a = "vpaid" === t.type;
                    t.off(),
                      t.on(
                        "all",
                        function (t, e) {
                          (a && t === d.F) ||
                            this.trigger(t, Object(i.g)({}, e, { type: t }));
                        },
                        this
                      );
                    var r = n.mediaModel;
                    t.on(d.bb, function (t) {
                      (t.oldstate = t.oldstate || n.get(d.bb)),
                        r.set("mediaState", t.newstate);
                    }),
                      t.on(d.X, this._nativeFullscreenHandler, this),
                      r.on("change:mediaState", function (t, n) {
                        e._stateHandler(n);
                      }),
                      t.attachMedia(),
                      t.volume(o.get("volume")),
                      t.mute(o.getMute()),
                      t.setPlaybackRate && t.setPlaybackRate(1),
                      o.on(
                        "change:volume",
                        function (t, e) {
                          this.volume = e;
                        },
                        this
                      ),
                      o.on(
                        "change:mute",
                        function (t, e) {
                          (this.mute = e), e || (this.volume = o.get("volume"));
                        },
                        this
                      ),
                      o.on(
                        "change:autostartMuted",
                        function (t, e) {
                          e ||
                            (n.set("autostartMuted", e),
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
                  var t = this.model,
                    e = this.mediaPool,
                    n = this.playerModel;
                  t.off();
                  var i = e.getPrimedElement();
                  if (this.backgroundLoading) {
                    e.clean();
                    var o = n.get("mediaContainer");
                    i.parentNode === o && o.removeChild(i);
                  } else
                    i &&
                      (i.removeEventListener("emptied", this.srcResetListener),
                      i.src !== t.get("mediaSrc") && this.srcReset());
                },
              },
              {
                key: "srcReset",
                value: function () {
                  var t = this.playerModel,
                    e = t.get("mediaModel"),
                    n = t.getVideo();
                  e.srcReset(), n && (n.src = null);
                },
              },
              {
                key: "_nativeFullscreenHandler",
                value: function (t) {
                  this.model.trigger(d.X, t),
                    this.trigger(d.y, { fullscreen: t.jwstate });
                },
              },
              {
                key: "_stateHandler",
                value: function (t) {
                  var e = this.model;
                  switch (t) {
                    case d.pb:
                    case d.ob:
                      e.set(d.bb, t);
                  }
                },
              },
              {
                key: "mute",
                set: function (t) {
                  var n = this.mediaController,
                    i = this.model,
                    o = this.provider;
                  i.set("mute", t),
                    et(ot(e.prototype), "mute", t, this, !0),
                    n || o.mute(t);
                },
              },
              {
                key: "volume",
                set: function (t) {
                  var n = this.mediaController,
                    i = this.model,
                    o = this.provider;
                  i.set("volume", t),
                    et(ot(e.prototype), "volume", t, this, !0),
                    n || o.volume(t);
                },
              },
            ]) && G(n.prototype, o),
            a && G(n, a),
            e
          );
        })(Y),
        st = { skipoffset: null, tag: null },
        lt = function (t, e, n, o) {
          var a,
            r,
            s,
            l,
            c = this,
            u = this,
            h = new rt(e, o),
            f = 0,
            j = {},
            g = null,
            b = {},
            m = A,
            v = !1,
            y = !1,
            k = !1,
            x = !1,
            O = function (t) {
              y ||
                (((t = t || {}).hasControls = !!e.get("controls")),
                c.trigger(d.z, t),
                h.model.get("state") === d.ob
                  ? t.hasControls && h.playVideo().catch(function () {})
                  : h.pause());
            },
            C = function () {
              y ||
                (h.model.get("state") === d.ob &&
                  e.get("controls") &&
                  (t.setFullscreen(), t.play()));
            };
          function M() {
            h.model.set("playRejected", !0);
          }
          function T() {
            f++, u.loadItem(a).catch(function () {});
          }
          function S(t, e) {
            "complete" !== t &&
              ((e = e || {}),
              b.tag && !e.tag && (e.tag = b.tag),
              this.trigger(t, e),
              ("mediaError" !== t && "error" !== t) ||
                (a && f + 1 < a.length && T()));
          }
          function _(t) {
            var e = t.newstate,
              n = t.oldstate || h.model.get("state");
            n !== e && E(Object(i.g)({ oldstate: n }, j, t));
          }
          function E(e) {
            var n = e.newstate;
            n === d.pb ? t.trigger(d.c, e) : n === d.ob && t.trigger(d.b, e);
          }
          function z(e) {
            var n = e.duration,
              i = e.position,
              o = h.model.mediaModel || h.model;
            o.set("duration", n),
              o.set("position", i),
              l || (l = (Object(w.d)(s, n) || n) - p.b),
              !v && i >= Math.max(l, p.a) && (t.preloadNextItem(), (v = !0));
          }
          function P(t) {
            var e = {};
            b.tag && (e.tag = b.tag), this.trigger(d.F, e), A.call(this, t);
          }
          function A(t) {
            (j = {}),
              a && f + 1 < a.length
                ? T()
                : (t.type === d.F && this.trigger(d.cb, {}), this.destroy());
          }
          function I() {
            y ||
              (n.clickHandler() &&
                n.clickHandler().setAlternateClickHandlers(O, C));
          }
          function R(t) {
            t.width && t.height && n.resizeMedia();
          }
          (this.init = function () {
            if (!k && !y) {
              (k = !0),
                (j = {}),
                h.setup(),
                h.on("all", S, this),
                h.on(d.O, M, this),
                h.on(d.S, z, this),
                h.on(d.F, P, this),
                h.on(d.K, R, this),
                h.on(d.bb, _, this),
                (g = t.detachMedia());
              var i = h.primedElement;
              e.get("mediaContainer").appendChild(i),
                e.set("instream", h),
                h.model.set("state", d.jb);
              var o = n.clickHandler();
              return (
                o && o.setAlternateClickHandlers(function () {}, null),
                this.setText(e.get("localization").loadingAd),
                (x = t.isBeforeComplete() || e.get("state") === d.kb),
                this
              );
            }
          }),
            (this.enableAdsMode = function (i) {
              var o = this;
              if (!k && !y)
                return (
                  t.routeEvents({
                    mediaControllerListener: function (t, e) {
                      o.trigger(t, e);
                    },
                  }),
                  e.set("instream", h),
                  h.model.set("state", d.pb),
                  (function (i) {
                    var o = n.clickHandler();
                    o &&
                      o.setAlternateClickHandlers(function (n) {
                        y ||
                          (((n = n || {}).hasControls = !!e.get("controls")),
                          u.trigger(d.z, n),
                          i &&
                            (e.get("state") === d.ob
                              ? t.playVideo()
                              : (t.pause(),
                                i &&
                                  (t.trigger(d.a, { clickThroughUrl: i }),
                                  window.open(i)))));
                      }, null);
                  })(i),
                  this
                );
            }),
            (this.setEventData = function (t) {
              j = t;
            }),
            (this.setState = function (t) {
              var e = t.newstate,
                n = h.model;
              (t.oldstate = n.get("state")), n.set("state", e), E(t);
            }),
            (this.setTime = function (e) {
              z(e), t.trigger(d.e, e);
            }),
            (this.loadItem = function (t, n) {
              if (y || !k)
                return Promise.reject(new Error("Instream not setup"));
              j = {};
              var o = t;
              Array.isArray(t)
                ? ((r = n || r), (t = (a = t)[f]), r && (n = r[f]))
                : (o = [t]);
              var l = h.model;
              l.set("playlist", o),
                e.set("hideAdsControls", !1),
                (t.starttime = 0),
                u.trigger(d.db, { index: f, item: t }),
                (b = Object(i.g)({}, st, n)),
                I(),
                l.set("skipButton", !1);
              var c =
                !e.get("backgroundLoading") && g
                  ? g.then(function () {
                      return h.setActiveItem(f);
                    })
                  : h.setActiveItem(f);
              return (
                (v = !1),
                void 0 !== (s = t.skipoffset || b.skipoffset) &&
                  u.setupSkipButton(s, b),
                c
              );
            }),
            (this.setupSkipButton = function (t, e, n) {
              var i = h.model;
              (m = n || A),
                i.set("skipMessage", e.skipMessage),
                i.set("skipText", e.skipText),
                i.set("skipOffset", t),
                (i.attributes.skipButton = !1),
                i.set("skipButton", !0);
            }),
            (this.applyProviderListeners = function (t) {
              h.usePsuedoProvider(t), I();
            }),
            (this.play = function () {
              (j = {}), h.playVideo();
            }),
            (this.pause = function () {
              (j = {}), h.pause();
            }),
            (this.skipAd = function (t) {
              var n = e.get("autoPause").pauseAds,
                i = "autostart" === e.get("playReason"),
                o = e.get("viewable");
              !n || i || o || (this.noResume = !0);
              var a = d.d;
              this.trigger(a, t), m.call(this, { type: a });
            }),
            (this.replacePlaylistItem = function (t) {
              y || (e.set("playlistItem", t), h.srcReset());
            }),
            (this.destroy = function () {
              y ||
                ((y = !0),
                this.trigger("destroyed"),
                this.off(),
                n.clickHandler() &&
                  n.clickHandler().revertAlternateClickHandlers(),
                e.off(null, null, h),
                h.off(null, null, u),
                h.destroy(),
                k && h.model && (e.attributes.state = d.ob),
                t.forwardEvents(),
                e.set("instream", null),
                (h = null),
                (j = {}),
                (g = null),
                k &&
                  !e.attributes._destroyed &&
                  (t.attachMedia(),
                  this.noResume || (x ? t.stopVideo() : t.playVideo())));
            }),
            (this.getState = function () {
              return !y && h.model.get("state");
            }),
            (this.setText = function (t) {
              return y ? this : (n.setAltText(t || ""), this);
            }),
            (this.hide = function () {
              y || e.set("hideAdsControls", !0);
            }),
            (this.getMediaElement = function () {
              return y ? null : h.primedElement;
            }),
            (this.setSkipOffset = function (t) {
              (s = t > 0 ? t : null), h && h.model.set("skipOffset", s);
            });
        };
      Object(i.g)(lt.prototype, h.a);
      var ct = lt,
        ut = n(66),
        dt = n(63),
        pt = function (t) {
          var e = this,
            n = [],
            i = {},
            o = 0,
            a = 0;
          function r(t) {
            if (
              ((t.data = t.data || []),
              (t.name = t.label || t.name || t.language),
              (t._id = Object(dt.a)(t, n.length)),
              !t.name)
            ) {
              var e = Object(dt.b)(t, o);
              (t.name = e.label), (o = e.unknownCount);
            }
            (i[t._id] = t), n.push(t);
          }
          function s() {
            for (
              var t = [{ id: "off", label: "Off" }], e = 0;
              e < n.length;
              e++
            )
              t.push({
                id: n[e]._id,
                label: n[e].name || "Unknown CC",
                language: n[e].language,
              });
            return t;
          }
          function l(e) {
            var i = (a = e),
              o = t.get("captionLabel");
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
                  ? t.setVideoSubtitleTrack(l, n)
                  : t.set("captionsIndex", l);
            } else t.set("captionsIndex", 0);
          }
          function c() {
            var e = s();
            u(e) !== u(t.get("captionsList")) &&
              (l(a), t.set("captionsList", e));
          }
          function u(t) {
            return t
              .map(function (t) {
                return "".concat(t.id, "-").concat(t.label);
              })
              .join(",");
          }
          t.on(
            "change:playlistItem",
            function (t) {
              (n = []), (i = {}), (o = 0);
              var e = t.attributes;
              (e.captionsIndex = 0),
                (e.captionsList = s()),
                t.set("captionsTrack", null);
            },
            this
          ),
            t.on(
              "change:itemReady",
              function () {
                var n = t.get("playlistItem").tracks,
                  o = n && n.length;
                if (o && !t.get("renderCaptionsNatively"))
                  for (
                    var a = function (t) {
                        var o,
                          a = n[t];
                        ("subtitles" !== (o = a.kind) && "captions" !== o) ||
                          i[a._id] ||
                          (r(a),
                          Object(ut.c)(
                            a,
                            function (t) {
                              !(function (t, e) {
                                t.data = e;
                              })(a, t);
                            },
                            function (t) {
                              e.trigger(d.tb, t);
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
            t.on(
              "change:captionsIndex",
              function (t, e) {
                var i = null;
                0 !== e && (i = n[e - 1]), t.set("captionsTrack", i);
              },
              this
            ),
            (this.setSubtitlesTracks = function (t) {
              if (Array.isArray(t)) {
                if (t.length) {
                  for (var e = 0; e < t.length; e++) r(t[e]);
                  n = Object.keys(i).map(function (t) {
                    return i[t];
                  });
                } else (n = []), (i = {}), (o = 0);
                c();
              }
            }),
            (this.selectDefaultIndex = l),
            (this.getCurrentIndex = function () {
              return t.get("captionsIndex");
            }),
            (this.getCaptionsList = function () {
              return t.get("captionsList");
            }),
            (this.destroy = function () {
              this.off(null, null, this);
            });
        };
      Object(i.g)(pt.prototype, h.a);
      var wt = pt,
        ht = function (t, e) {
          return (
            '<div id="'
              .concat(
                t,
                '" class="jwplayer jw-reset jw-state-setup" tabindex="0" aria-label="'
              )
              .concat(e || "", '" role="application">') +
            '<div class="jw-aspect jw-reset"></div><div class="jw-wrapper jw-reset"><div class="jw-top jw-reset"></div><div class="jw-aspect jw-reset"></div><div class="jw-media jw-reset"></div><div class="jw-preview jw-reset"></div><div class="jw-title jw-reset-text" dir="auto"><div class="jw-title-primary jw-reset-text"></div><div class="jw-title-secondary jw-reset-text"></div></div><div class="jw-overlays jw-reset"></div><div class="jw-hidden-accessibility"><span class="jw-time-update" aria-live="assertive"></span><span class="jw-volume-update" aria-live="assertive"></span></div></div></div>'
          );
        },
        ft = n(35),
        jt = 44,
        gt = function (t) {
          var e = t.get("height");
          if (t.get("aspectratio")) return !1;
          if ("string" == typeof e && e.indexOf("%") > -1) return !1;
          var n = 1 * e || NaN;
          return (
            !!(n = isNaN(n) ? t.get("containerHeight") : n) && n && n <= jt
          );
        },
        bt = n(54);
      function mt(t, e) {
        if (t.get("fullscreen")) return 1;
        if (!t.get("activeTab")) return 0;
        if (t.get("isFloating")) return 1;
        var n = t.get("intersectionRatio");
        return void 0 === n &&
          ((n = (function (t) {
            var e = document.documentElement,
              n = document.body,
              i = {
                top: 0,
                left: 0,
                right: e.clientWidth || n.clientWidth,
                width: e.clientWidth || n.clientWidth,
                bottom: e.clientHeight || n.clientHeight,
                height: e.clientHeight || n.clientHeight,
              };
            if (!n.contains(t)) return 0;
            if ("none" === window.getComputedStyle(t).display) return 0;
            var o = vt(t);
            if (!o) return 0;
            var a = o,
              r = t.parentNode,
              s = !1;
            for (; !s; ) {
              var l = null;
              if (
                (r === n || r === e || 1 !== r.nodeType
                  ? ((s = !0), (l = i))
                  : "visible" !== window.getComputedStyle(r).overflow &&
                    (l = vt(r)),
                l &&
                  ((c = l),
                  (u = a),
                  (d = void 0),
                  (p = void 0),
                  (w = void 0),
                  (h = void 0),
                  (f = void 0),
                  (j = void 0),
                  (d = Math.max(c.top, u.top)),
                  (p = Math.min(c.bottom, u.bottom)),
                  (w = Math.max(c.left, u.left)),
                  (h = Math.min(c.right, u.right)),
                  (j = p - d),
                  !(a = (f = h - w) >= 0 &&
                    j >= 0 && {
                      top: d,
                      bottom: p,
                      left: w,
                      right: h,
                      width: f,
                      height: j,
                    })))
              )
                return 0;
              r = r.parentNode;
            }
            var c, u, d, p, w, h, f, j;
            var g = o.width * o.height,
              b = a.width * a.height;
            return g ? b / g : 0;
          })(e)),
          window.top !== window.self && n)
          ? 0
          : n;
      }
      function vt(t) {
        try {
          return t.getBoundingClientRect();
        } catch (t) {}
      }
      var yt = n(49),
        kt = n(42),
        xt = n(58),
        Ot = n(10);
      var Ct = n(32),
        Mt = n(5),
        Tt = n(6),
        St = [
          "fullscreenchange",
          "webkitfullscreenchange",
          "mozfullscreenchange",
          "MSFullscreenChange",
        ],
        _t = function (t, e, n) {
          for (
            var i =
                t.requestFullscreen ||
                t.webkitRequestFullscreen ||
                t.webkitRequestFullScreen ||
                t.mozRequestFullScreen ||
                t.msRequestFullscreen,
              o =
                e.exitFullscreen ||
                e.webkitExitFullscreen ||
                e.webkitCancelFullScreen ||
                e.mozCancelFullScreen ||
                e.msExitFullscreen,
              a = !(!i || !o),
              r = St.length;
            r--;

          )
            e.addEventListener(St[r], n);
          return {
            events: St,
            supportsDomFullscreen: function () {
              return a;
            },
            requestFullscreen: function () {
              i.call(t, { navigationUI: "hide" });
            },
            exitFullscreen: function () {
              null !== this.fullscreenElement() && o.apply(e);
            },
            fullscreenElement: function () {
              var t = e.fullscreenElement,
                n = e.webkitCurrentFullScreenElement,
                i = e.mozFullScreenElement,
                o = e.msFullscreenElement;
              return null === t ? t : t || n || i || o;
            },
            destroy: function () {
              for (var t = St.length; t--; ) e.removeEventListener(St[t], n);
            },
          };
        },
        Et = n(40);
      function zt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var Pt = (function () {
          function t(e, n) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(i.g)(this, h.a),
              this.revertAlternateClickHandlers(),
              (this.domElement = n),
              (this.model = e),
              (this.ui = new Et.a(n)
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
          var e, n, o;
          return (
            (e = t),
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
                value: function (t) {
                  this.model.get("flashBlocked") ||
                    (this.alternateClickHandler
                      ? this.alternateClickHandler(t)
                      : this.trigger(t.type === d.n ? "click" : "tap"));
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
                value: function (t, e) {
                  (this.alternateClickHandler = t),
                    (this.alternateDoubleClickHandler = e || null);
                },
              },
              {
                key: "revertAlternateClickHandlers",
                value: function () {
                  (this.alternateClickHandler = null),
                    (this.alternateDoubleClickHandler = null);
                },
              },
            ]) && zt(e.prototype, n),
            o && zt(e, o),
            t
          );
        })(),
        At = n(59),
        It = function (t, e) {
          var n = e ? " jw-hide" : "";
          return '<div class="jw-logo jw-logo-'
            .concat(t)
            .concat(n, ' jw-reset"></div>');
        },
        Rt = {
          linktarget: "_blank",
          margin: 8,
          hide: !1,
          position: "top-right",
        };
      function Lt(t) {
        var e, n;
        Object(i.g)(this, h.a);
        var o = new Image();
        (this.setup = function () {
          ((n = Object(i.g)({}, Rt, t.get("logo"))).position =
            n.position || Rt.position),
            (n.hide = "true" === n.hide.toString()),
            n.file &&
              "control-bar" !== n.position &&
              (e || (e = Object(Mt.e)(It(n.position, n.hide))),
              t.set("logo", n),
              (o.onload = function () {
                var i = this.height,
                  o = this.width,
                  a = { backgroundImage: 'url("' + this.src + '")' };
                if (n.margin !== Rt.margin) {
                  var r = /(\w+)-(\w+)/.exec(n.position);
                  3 === r.length &&
                    ((a["margin-" + r[1]] = n.margin),
                    (a["margin-" + r[2]] = n.margin));
                }
                var s = 0.15 * t.get("containerHeight"),
                  l = 0.15 * t.get("containerWidth");
                if (i > s || o > l) {
                  var c = o / i;
                  l / s > c ? ((i = s), (o = s * c)) : ((o = l), (i = l / c));
                }
                (a.width = Math.round(o)),
                  (a.height = Math.round(i)),
                  Object(Ot.d)(e, a),
                  t.set("logoWidth", a.width);
              }),
              (o.src = n.file),
              n.link &&
                (e.setAttribute("tabindex", "0"),
                e.setAttribute("aria-label", t.get("localization").logo)),
              (this.ui = new Et.a(e).on(
                "click tap enter",
                function (t) {
                  t && t.stopPropagation && t.stopPropagation(),
                    this.trigger(d.A, {
                      link: n.link,
                      linktarget: n.linktarget,
                    });
                },
                this
              )));
        }),
          (this.setContainer = function (t) {
            e && t.appendChild(e);
          }),
          (this.element = function () {
            return e;
          }),
          (this.position = function () {
            return n.position;
          }),
          (this.destroy = function () {
            (o.onload = null), this.ui && this.ui.destroy();
          });
      }
      var Bt = function (t) {
        (this.model = t), (this.image = null);
      };
      Object(i.g)(Bt.prototype, {
        setup: function (t) {
          this.el = t;
        },
        setImage: function (t) {
          var e = this.image;
          e && (e.onload = null), (this.image = null);
          var n = "";
          "string" == typeof t &&
            ((n = 'url("' + t + '")'),
            ((e = this.image = new Image()).src = t)),
            Object(Ot.d)(this.el, { backgroundImage: n });
        },
        resize: function (t, e, n) {
          if ("uniform" === n) {
            if (
              (t && (this.playerAspectRatio = t / e),
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
                  a.resize(t, e, n);
                });
              }
              var r = i.width / i.height;
              Math.abs(this.playerAspectRatio - r) < 0.09 && (o = "cover");
            }
            Object(Ot.d)(this.el, { backgroundSize: o });
          }
          var s;
        },
        element: function () {
          return this.el;
        },
      });
      var Vt = Bt,
        Nt = function (t) {
          this.model = t.player;
        };
      Object(i.g)(Nt.prototype, {
        hide: function () {
          Object(Ot.d)(this.el, { display: "none" });
        },
        show: function () {
          Object(Ot.d)(this.el, { display: "" });
        },
        setup: function (t) {
          this.el = t;
          var e = this.el.getElementsByTagName("div");
          (this.title = e[0]),
            (this.description = e[1]),
            this.model.on("change:logoWidth", this.update, this),
            this.model.change("playlistItem", this.playlistItem, this);
        },
        update: function (t) {
          var e = {},
            n = t.get("logo");
          if (n) {
            var i = 1 * ("" + n.margin).replace("px", ""),
              o = t.get("logoWidth") + (isNaN(i) ? 0 : i + 10);
            "top-left" === n.position
              ? (e.paddingLeft = o)
              : "top-right" === n.position && (e.paddingRight = o);
          }
          Object(Ot.d)(this.el, e);
        },
        playlistItem: function (t, e) {
          if (e)
            if (t.get("displaytitle") || t.get("displaydescription")) {
              var n = "",
                i = "";
              e.title && t.get("displaytitle") && (n = e.title),
                e.description &&
                  t.get("displaydescription") &&
                  (i = e.description),
                this.updateText(n, i);
            } else this.hide();
        },
        updateText: function (t, e) {
          Object(Mt.q)(this.title, t),
            Object(Mt.q)(this.description, e),
            this.title.firstChild || this.description.firstChild
              ? this.show()
              : this.hide();
        },
        element: function () {
          return this.el;
        },
      });
      var Ht = Nt;
      function Ft(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      var qt,
        Dt = (function () {
          function t(e) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              (this.container = e),
              (this.input = e.querySelector(".jw-media"));
          }
          var e, n, i;
          return (
            (e = t),
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
                  var t,
                    e,
                    n,
                    i,
                    o = this.container,
                    a = this.input,
                    r = (this.ui = new Et.a(a, { preventScrolling: !0 })
                      .on("dragStart", function () {
                        (t = o.offsetLeft),
                          (e = o.offsetTop),
                          (n = window.innerHeight),
                          (i = window.innerWidth);
                      })
                      .on("drag", function (a) {
                        var s = Math.max(t + a.pageX - r.startX, 0),
                          l = Math.max(e + a.pageY - r.startY, 0),
                          c = Math.max(i - (s + o.clientWidth), 0),
                          u = Math.max(n - (l + o.clientHeight), 0);
                        0 === c ? (s = "auto") : (c = "auto"),
                          0 === l ? (u = "auto") : (l = "auto"),
                          Object(Ot.d)(o, {
                            left: s,
                            right: c,
                            top: l,
                            bottom: u,
                            margin: 0,
                          });
                      })
                      .on("dragEnd", function () {
                        t = e = i = n = null;
                      }));
                },
              },
            ]) && Ft(e.prototype, n),
            i && Ft(e, i),
            t
          );
        })(),
        Ut = n(55);
      n(69);
      var Wt = m.OS.mobile,
        Qt = m.Browser.ie,
        Yt = null;
      var Xt = function (t, e) {
        var n,
          o,
          a,
          r,
          s = this,
          l = Object(i.g)(this, h.a, { isSetup: !1, api: t, model: e }),
          c = e.get("localization"),
          u = Object(Mt.e)(ht(e.get("id"), c.player)),
          p = u.querySelector(".jw-wrapper"),
          f = u.querySelector(".jw-media"),
          j = new Dt(p),
          g = new Vt(e, t),
          b = new Ht(e),
          v = new At.b(e);
        v.on("all", l.trigger, l);
        var y = -1,
          k = -1,
          x = -1,
          O = e.get("floating");
        this.dismissible = O && O.dismissible;
        var C,
          M,
          T,
          S = !1,
          _ = {},
          E = null,
          z = null;
        function P() {
          return Wt && !Object(Mt.f)();
        }
        function A() {
          Object(kt.a)(k), (k = Object(kt.b)(I));
        }
        function I() {
          l.isSetup && (l.updateBounds(), l.updateStyles(), l.checkResized());
        }
        function R(t, n) {
          if (Object(i.r)(t) && Object(i.r)(n)) {
            var o = Object(xt.a)(t);
            Object(xt.b)(u, o);
            var a = o < 2;
            Object(Mt.v)(u, "jw-flag-small-player", a),
              Object(Mt.v)(u, "jw-orientation-portrait", n > t);
          }
          if (e.get("controls")) {
            var r = gt(e);
            Object(Mt.v)(u, "jw-flag-audio-player", r), e.set("audioMode", r);
          }
        }
        function L() {
          e.set("visibility", mt(e, u));
        }
        (this.updateBounds = function () {
          Object(kt.a)(k);
          var t = e.get("isFloating") ? p : u,
            n = document.body.contains(t),
            i = Object(Mt.c)(t),
            r = Math.round(i.width),
            s = Math.round(i.height);
          if (((_ = Object(Mt.c)(u)), r === o && s === a))
            return (o && a) || A(), void e.set("inDom", n);
          (r && s) || (o && a) || A(),
            (r || s || n) &&
              (e.set("containerWidth", r), e.set("containerHeight", s)),
            e.set("inDom", n),
            n && bt.a.observe(u);
        }),
          (this.updateStyles = function () {
            var t = e.get("containerWidth"),
              n = e.get("containerHeight");
            R(t, n), z && z.resize(t, n), $(t, n), v.resize(), O && F();
          }),
          (this.checkResized = function () {
            var t = e.get("containerWidth"),
              n = e.get("containerHeight"),
              i = e.get("isFloating");
            if (t !== o || n !== a) {
              this.resizeListener ||
                (this.resizeListener = new Ut.a(p, this, e)),
                (o = t),
                (a = n),
                l.trigger(d.hb, { width: t, height: n });
              var s = Object(xt.a)(t);
              E !== s && ((E = s), l.trigger(d.j, { breakpoint: E }));
            }
            i !== r && ((r = i), l.trigger(d.x, { floating: i }), L());
          }),
          (this.responsiveListener = A),
          (this.setup = function () {
            g.setup(u.querySelector(".jw-preview")),
              b.setup(u.querySelector(".jw-title")),
              (n = new Lt(e)).setup(),
              n.setContainer(p),
              n.on(d.A, K),
              v.setup(u.id, e.get("captions")),
              b.element().parentNode.insertBefore(v.element(), b.element()),
              (C = (function (t, e, n) {
                var i = new Pt(e, n),
                  o = e.get("controls");
                i.on({
                  click: function () {
                    l.trigger(d.p),
                      z &&
                        (ct()
                          ? z.settingsMenu.close()
                          : ut()
                          ? z.infoOverlay.close()
                          : t.playToggle({ reason: "interaction" }));
                  },
                  tap: function () {
                    l.trigger(d.p),
                      ct() && z.settingsMenu.close(),
                      ut() && z.infoOverlay.close();
                    var n = e.get("state");
                    if (
                      (o &&
                        (n === d.mb ||
                          n === d.kb ||
                          (e.get("instream") && n === d.ob)) &&
                        t.playToggle({ reason: "interaction" }),
                      o && n === d.ob)
                    ) {
                      if (
                        e.get("instream") ||
                        e.get("castActive") ||
                        "audio" === e.get("mediaType")
                      )
                        return;
                      Object(Mt.v)(u, "jw-flag-controls-hidden"),
                        l.dismissible &&
                          Object(Mt.v)(
                            u,
                            "jw-floating-dismissible",
                            Object(Mt.i)(u, "jw-flag-controls-hidden")
                          ),
                        v.renderCues(!0);
                    } else z && (z.showing ? z.userInactive() : z.userActive());
                  },
                  doubleClick: function () {
                    return z && t.setFullscreen();
                  },
                }),
                  Wt ||
                    (u.addEventListener("mousemove", W),
                    u.addEventListener("mouseover", Q),
                    u.addEventListener("mouseout", Y));
                return i;
              })(t, e, f)),
              (T = new Et.a(u).on("click", function () {})),
              (M = _t(u, document, et)),
              e.on("change:hideAdsControls", function (t, e) {
                Object(Mt.v)(u, "jw-flag-ads-hide-controls", e);
              }),
              e.on("change:scrubbing", function (t, e) {
                Object(Mt.v)(u, "jw-flag-dragging", e);
              }),
              e.on("change:playRejected", function (t, e) {
                Object(Mt.v)(u, "jw-flag-play-rejected", e);
              }),
              e.on(d.X, tt),
              e.on("change:".concat(d.U), function () {
                $(), v.resize();
              }),
              e.player.on("change:errorEvent", at),
              e.change("stretching", X);
            var i = e.get("width"),
              o = e.get("height"),
              a = G(i, o);
            Object(Ot.d)(u, a),
              e.change("aspectratio", Z),
              R(i, o),
              e.get("controls") ||
                (Object(Mt.a)(u, "jw-flag-controls-hidden"),
                Object(Mt.o)(u, "jw-floating-dismissible")),
              Qt && Object(Mt.a)(u, "jw-ie");
            var r = e.get("skin") || {};
            r.name && Object(Mt.p)(u, /jw-skin-\S+/, "jw-skin-" + r.name);
            var s = (function (t) {
              t || (t = {});
              var e = t.active,
                n = t.inactive,
                i = t.background,
                o = {};
              return (
                (o.controlbar = (function (t) {
                  if (t || e || n || i) {
                    var o = {};
                    return (
                      (t = t || {}),
                      (o.iconsActive = t.iconsActive || e),
                      (o.icons = t.icons || n),
                      (o.text = t.text || n),
                      (o.background = t.background || i),
                      o
                    );
                  }
                })(t.controlbar)),
                (o.timeslider = (function (t) {
                  if (t || e) {
                    var n = {};
                    return (
                      (t = t || {}),
                      (n.progress = t.progress || e),
                      (n.rail = t.rail),
                      n
                    );
                  }
                })(t.timeslider)),
                (o.menus = (function (t) {
                  if (t || e || n || i) {
                    var o = {};
                    return (
                      (t = t || {}),
                      (o.text = t.text || n),
                      (o.textActive = t.textActive || e),
                      (o.background = t.background || i),
                      o
                    );
                  }
                })(t.menus)),
                (o.tooltips = (function (t) {
                  if (t || n || i) {
                    var e = {};
                    return (
                      (t = t || {}),
                      (e.text = t.text || n),
                      (e.background = t.background || i),
                      e
                    );
                  }
                })(t.tooltips)),
                o
              );
            })(r);
            !(function (t, e) {
              var n;
              function i(e, n, i, o) {
                if (i) {
                  e = Object(w.f)(e, "#" + t + (o ? "" : " "));
                  var a = {};
                  (a[n] = i), Object(Ot.b)(e.join(", "), a, t);
                }
              }
              e &&
                (e.controlbar &&
                  (function (e) {
                    i(
                      [
                        ".jw-controlbar .jw-icon-inline.jw-text",
                        ".jw-title-primary",
                        ".jw-title-secondary",
                      ],
                      "color",
                      e.text
                    ),
                      e.icons &&
                        (i(
                          [
                            ".jw-button-color:not(.jw-icon-cast)",
                            ".jw-button-color.jw-toggle.jw-off:not(.jw-icon-cast)",
                          ],
                          "color",
                          e.icons
                        ),
                        i(
                          [".jw-display-icon-container .jw-button-color"],
                          "color",
                          e.icons
                        ),
                        Object(Ot.b)(
                          "#".concat(
                            t,
                            " .jw-icon-cast google-cast-launcher.jw-off"
                          ),
                          "{--disconnected-color: ".concat(e.icons, "}"),
                          t
                        ));
                    e.iconsActive &&
                      (i(
                        [
                          ".jw-display-icon-container .jw-button-color:hover",
                          ".jw-display-icon-container .jw-button-color:focus",
                        ],
                        "color",
                        e.iconsActive
                      ),
                      i(
                        [
                          ".jw-button-color.jw-toggle:not(.jw-icon-cast)",
                          ".jw-button-color:hover:not(.jw-icon-cast)",
                          ".jw-button-color:focus:not(.jw-icon-cast)",
                          ".jw-button-color.jw-toggle.jw-off:hover:not(.jw-icon-cast)",
                        ],
                        "color",
                        e.iconsActive
                      ),
                      i([".jw-svg-icon-buffer"], "fill", e.icons),
                      Object(Ot.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:hover google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Ot.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:focus google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Ot.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast google-cast-launcher.jw-off:focus"
                        ),
                        "{--disconnected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Ot.b)(
                        "#".concat(t, " .jw-icon-cast google-cast-launcher"),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Ot.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast google-cast-launcher:focus"
                        ),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Ot.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:hover google-cast-launcher"
                        ),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Ot.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:focus google-cast-launcher"
                        ),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ));
                    i(
                      [
                        " .jw-settings-topbar",
                        ":not(.jw-state-idle) .jw-controlbar",
                        ".jw-flag-audio-player .jw-controlbar",
                      ],
                      "background",
                      e.background,
                      !0
                    );
                  })(e.controlbar),
                e.timeslider &&
                  (function (t) {
                    var e = t.progress;
                    "none" !== e &&
                      (i([".jw-progress", ".jw-knob"], "background-color", e),
                      i(
                        [".jw-buffer"],
                        "background-color",
                        Object(Ot.c)(e, 50)
                      ));
                    i([".jw-rail"], "background-color", t.rail),
                      i(
                        [
                          ".jw-background-color.jw-slider-time",
                          ".jw-slider-time .jw-cue",
                        ],
                        "background-color",
                        t.background
                      );
                  })(e.timeslider),
                e.menus &&
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
                    (n = e.menus).text
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
                e.tooltips &&
                  (function (t) {
                    i(
                      [
                        ".jw-skip",
                        ".jw-tooltip .jw-text",
                        ".jw-time-tip .jw-text",
                      ],
                      "background-color",
                      t.background
                    ),
                      i([".jw-time-tip", ".jw-tooltip"], "color", t.background),
                      i([".jw-skip"], "border", "none"),
                      i(
                        [
                          ".jw-skip .jw-text",
                          ".jw-skip .jw-icon",
                          ".jw-time-tip .jw-text",
                          ".jw-tooltip .jw-text",
                        ],
                        "color",
                        t.text
                      );
                  })(e.tooltips),
                e.menus &&
                  (function (e) {
                    if (e.textActive) {
                      var n = {
                        color: e.textActive,
                        borderColor: e.textActive,
                        stroke: e.textActive,
                      };
                      Object(Ot.b)("#".concat(t, " .jw-color-active"), n, t),
                        Object(Ot.b)(
                          "#".concat(t, " .jw-color-active-hover:hover"),
                          n,
                          t
                        );
                    }
                    if (e.text) {
                      var i = {
                        color: e.text,
                        borderColor: e.text,
                        stroke: e.text,
                      };
                      Object(Ot.b)("#".concat(t, " .jw-color-inactive"), i, t),
                        Object(Ot.b)(
                          "#".concat(t, " .jw-color-inactive-hover:hover"),
                          i,
                          t
                        );
                    }
                  })(e.menus));
            })(e.get("id"), s),
              e.set("mediaContainer", f),
              e.set("iFrame", m.Features.iframe),
              e.set("activeTab", Object(yt.a)()),
              e.set("touchMode", Wt && ("string" == typeof o || o >= jt)),
              bt.a.add(this),
              e.get("enableGradient") &&
                !Qt &&
                Object(Mt.a)(u, "jw-ab-drop-shadow"),
              (this.isSetup = !0),
              e.trigger("viewSetup", u);
            var c = document.body.contains(u);
            c && bt.a.observe(u), e.set("inDom", c);
          }),
          (this.init = function () {
            this.updateBounds(),
              e.on("change:fullscreen", J),
              e.on("change:activeTab", L),
              e.on("change:fullscreen", L),
              e.on("change:intersectionRatio", L),
              e.on("change:visibility", U),
              e.on("instreamMode", function (t) {
                t ? dt() : pt();
              }),
              L(),
              1 !== bt.a.size() || e.get("visibility") || U(e, 1, 0);
            var t = e.player;
            e.change("state", rt),
              t.change("controls", q),
              e.change("streamType", it),
              e.change("mediaType", ot),
              t.change("playlistItem", function (t, e) {
                lt(t, e);
              }),
              (o = a = null),
              O && Wt && bt.a.addScrollHandler(F),
              this.checkResized();
          });
        var B,
          V = 62,
          N = !0;
        function H() {
          var t = e.get("isFloating"),
            n = _.top < V,
            i = n ? _.top <= window.scrollY : _.top <= window.scrollY + V;
          !t && i ? wt(0, n) : t && !i && wt(1, n);
        }
        function F() {
          P() &&
            e.get("inDom") &&
            (clearTimeout(B),
            (B = setTimeout(H, 150)),
            N &&
              ((N = !1),
              H(),
              setTimeout(function () {
                N = !0;
              }, 50)));
        }
        function q(t, e) {
          var n = { controls: e };
          e
            ? (qt = Ct.a.controls)
              ? D()
              : ((n.loadPromise = Object(Ct.b)().then(function (e) {
                  qt = e;
                  var n = t.get("controls");
                  return n && D(), n;
                })),
                n.loadPromise.catch(function (t) {
                  l.trigger(d.tb, t);
                }))
            : l.removeControls(),
            o && a && l.trigger(d.o, n);
        }
        function D() {
          var t = new qt(document, l.element());
          l.addControls(t);
        }
        function U(t, e, n) {
          e && !n && (rt(t, t.get("state")), l.updateStyles());
        }
        function W(t) {
          z && z.mouseMove(t);
        }
        function Q(t) {
          z && !z.showing && "IFRAME" === t.target.nodeName && z.userActive();
        }
        function Y(t) {
          z &&
            z.showing &&
            ((t.relatedTarget && !u.contains(t.relatedTarget)) ||
              (!t.relatedTarget && m.Features.iframe)) &&
            z.userActive();
        }
        function X(t, e) {
          Object(Mt.p)(u, /jw-stretch-\S+/, "jw-stretch-" + e);
        }
        function Z(t, n) {
          Object(Mt.v)(u, "jw-flag-aspect-mode", !!n);
          var i = u.querySelectorAll(".jw-aspect");
          Object(Ot.d)(i, { paddingTop: n || null }),
            l.isSetup &&
              n &&
              !e.get("isFloating") &&
              (Object(Ot.d)(u, G(t.get("width"))), I());
        }
        function K(n) {
          n.link
            ? (t.pause({ reason: "interaction" }),
              t.setFullscreen(!1),
              Object(Mt.l)(n.link, n.linktarget, { rel: "noreferrer" }))
            : e.get("controls") && t.playToggle({ reason: "interaction" });
        }
        (this.addControls = function (n) {
          var i = this;
          (z = n),
            Object(Mt.o)(u, "jw-flag-controls-hidden"),
            Object(Mt.v)(u, "jw-floating-dismissible", this.dismissible),
            n.enable(t, e),
            a && (R(o, a), n.resize(o, a), v.renderCues(!0)),
            n.on("userActive userInactive", function () {
              var t = e.get("state");
              (t !== d.pb && t !== d.jb) || v.renderCues(!0);
            }),
            n.on("dismissFloating", function () {
              i.stopFloating(!0), t.pause({ reason: "interaction" });
            }),
            n.on("all", l.trigger, l),
            e.get("instream") && z.setupInstream();
        }),
          (this.removeControls = function () {
            z && (z.disable(e), (z = null)),
              Object(Mt.a)(u, "jw-flag-controls-hidden"),
              Object(Mt.o)(u, "jw-floating-dismissible");
          });
        var J = function (e, n) {
          if (
            (n && z && e.get("autostartMuted") && z.unmuteAutoplay(t, e),
            M.supportsDomFullscreen())
          )
            n ? M.requestFullscreen() : M.exitFullscreen(), nt(u, n);
          else if (Qt) nt(u, n);
          else {
            var i = e.get("instream"),
              o = i ? i.provider : null,
              a = e.getVideo() || o;
            a && a.setFullscreen && a.setFullscreen(n);
          }
        };
        function G(t, n, o) {
          var a = { width: t };
          if (
            (o && void 0 !== n && e.set("aspectratio", null),
            !e.get("aspectratio"))
          ) {
            var r = n;
            Object(i.r)(r) && 0 !== r && (r = Math.max(r, jt)), (a.height = r);
          }
          return a;
        }
        function $(t, n) {
          if (
            ((t && !isNaN(1 * t)) || (t = e.get("containerWidth"))) &&
            ((n && !isNaN(1 * n)) || (n = e.get("containerHeight")))
          ) {
            g && g.resize(t, n, e.get("stretching"));
            var i = e.getVideo();
            i && i.resize(t, n, e.get("stretching"));
          }
        }
        function tt(t) {
          Object(Mt.v)(u, "jw-flag-ios-fullscreen", t.jwstate), et(t);
        }
        function et(t) {
          var n = e.get("fullscreen"),
            i =
              void 0 !== t.jwstate
                ? t.jwstate
                : (function () {
                    if (M.supportsDomFullscreen()) {
                      var t = M.fullscreenElement();
                      return !(!t || t !== u);
                    }
                    return e.getVideo().getFullScreen();
                  })();
          n !== i && e.set("fullscreen", i),
            A(),
            clearTimeout(y),
            (y = setTimeout($, 200));
        }
        function nt(t, e) {
          Object(Mt.v)(t, "jw-flag-fullscreen", e),
            Object(Ot.d)(document.body, { overflowY: e ? "hidden" : "" }),
            e && z && z.userActive(),
            $(),
            A();
        }
        function it(t, e) {
          var n = "LIVE" === e;
          Object(Mt.v)(u, "jw-flag-live", n);
        }
        function ot(t, e) {
          var n = "audio" === e,
            i = t.get("provider");
          Object(Mt.v)(u, "jw-flag-media-audio", n);
          var o = i && 0 === i.name.indexOf("flash"),
            a = n && !o ? f : f.nextSibling;
          g.el.parentNode.insertBefore(g.el, a);
        }
        function at(t, e) {
          if (e) {
            var n = Object(ft.a)(t, e);
            ft.a.cloneIcon &&
              n.querySelector(".jw-icon").appendChild(ft.a.cloneIcon("error")),
              b.hide(),
              u.appendChild(n.firstChild),
              Object(Mt.v)(u, "jw-flag-audio-player", !!t.get("audioMode"));
          } else b.playlistItem(t, t.get("playlistItem"));
        }
        function rt(t, e, n) {
          if (l.isSetup) {
            if (n === d.lb) {
              var i = u.querySelector(".jw-error-msg");
              i && i.parentNode.removeChild(i);
            }
            Object(kt.a)(x),
              e === d.pb
                ? st(e)
                : (x = Object(kt.b)(function () {
                    return st(e);
                  }));
          }
        }
        function st(t) {
          switch (
            (e.get("controls") &&
              t !== d.ob &&
              Object(Mt.i)(u, "jw-flag-controls-hidden") &&
              (Object(Mt.o)(u, "jw-flag-controls-hidden"),
              Object(Mt.v)(u, "jw-floating-dismissible", l.dismissible)),
            Object(Mt.p)(u, /jw-state-\S+/, "jw-state-" + t),
            t)
          ) {
            case d.lb:
              l.stopFloating();
            case d.mb:
            case d.kb:
              v && v.hide();
              break;
            default:
              v &&
                (v.show(), t === d.ob && z && !z.showing && v.renderCues(!0));
          }
        }
        (this.resize = function (t, n) {
          var i = G(t, n, !0);
          void 0 !== t &&
            void 0 !== n &&
            (e.set("width", t), e.set("height", n)),
            Object(Ot.d)(u, i),
            e.get("isFloating") && vt(),
            I();
        }),
          (this.resizeMedia = $),
          (this.setPosterImage = function (t, e) {
            e.setImage(t && t.image);
          });
        var lt = function (t, e) {
            s.setPosterImage(e, g),
              Wt &&
                (function (t, e) {
                  var n = t.get("mediaElement");
                  if (n) {
                    var i = Object(Mt.j)(e.title || "");
                    n.setAttribute("title", i.textContent);
                  }
                })(t, e);
          },
          ct = function () {
            var t = z && z.settingsMenu;
            return !(!t || !t.visible);
          },
          ut = function () {
            var t = z && z.infoOverlay;
            return !(!t || !t.visible);
          },
          dt = function () {
            Object(Mt.a)(u, "jw-flag-ads"), z && z.setupInstream(), j.disable();
          },
          pt = function () {
            if (C) {
              z && z.destroyInstream(e),
                Yt !== u || Object(Tt.m)() || j.enable(),
                l.setAltText(""),
                Object(Mt.o)(u, ["jw-flag-ads", "jw-flag-ads-hide-controls"]),
                e.set("hideAdsControls", !1);
              var t = e.getVideo();
              t && t.setContainer(f), C.revertAlternateClickHandlers();
            }
          };
        function wt(t, n) {
          if (t < 0.5 && !Object(Tt.m)()) {
            var i = e.get("state");
            i !== d.mb &&
              i !== d.lb &&
              i !== d.kb &&
              null === Yt &&
              ((Yt = u),
              e.set("isFloating", !0),
              Object(Mt.a)(u, "jw-flag-floating"),
              n &&
                (Object(Ot.d)(p, {
                  transform: "translateY(-".concat(V - _.top, "px)"),
                }),
                setTimeout(function () {
                  Object(Ot.d)(p, {
                    transform: "translateY(0)",
                    transition:
                      "transform 150ms cubic-bezier(0, 0.25, 0.25, 1)",
                  });
                })),
              Object(Ot.d)(u, {
                backgroundImage: g.el.style.backgroundImage || e.get("image"),
              }),
              vt(),
              e.get("instreamMode") || j.enable(),
              A());
          } else l.stopFloating(!1, n);
        }
        function vt() {
          var t = e.get("width"),
            n = e.get("height"),
            o = G(t);
          if (((o.maxWidth = Math.min(400, _.width)), !e.get("aspectratio"))) {
            var a = _.width,
              r = _.height / a || 0.5625;
            Object(i.r)(t) && Object(i.r)(n) && (r = n / t),
              Z(e, 100 * r + "%");
          }
          Object(Ot.d)(p, o);
        }
        (this.setAltText = function (t) {
          e.set("altText", t);
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
            var t =
                !(arguments.length > 0 && void 0 !== arguments[0]) ||
                arguments[0],
              e = { x: 0, y: 0, width: o || 0, height: a || 0 };
            return z && t && (e.height -= z.controlbarHeight()), e;
          }),
          (this.setCaptions = function (t) {
            v.clear(), v.setup(e.get("id"), t), v.resize();
          }),
          (this.setIntersection = function (t) {
            var n = Math.round(100 * t.intersectionRatio) / 100;
            e.set("intersectionRatio", n),
              O && !P() && (S = S || n >= 0.5) && wt(n);
          }),
          (this.stopFloating = function (t, n) {
            if ((t && ((O = null), bt.a.removeScrollHandler(F)), Yt === u)) {
              (Yt = null), e.set("isFloating", !1);
              var i = function () {
                Object(Mt.o)(u, "jw-flag-floating"),
                  Z(e, e.get("aspectratio")),
                  Object(Ot.d)(u, { backgroundImage: null }),
                  Object(Ot.d)(p, {
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
                ? (Object(Ot.d)(p, {
                    transform: "translateY(-".concat(V - _.top, "px)"),
                    "transition-timing-function": "ease-out",
                  }),
                  setTimeout(i, 150))
                : i(),
                j.disable(),
                A();
            }
          }),
          (this.destroy = function () {
            e.destroy(),
              bt.a.unobserve(u),
              bt.a.remove(this),
              (this.isSetup = !1),
              this.off(),
              Object(kt.a)(k),
              clearTimeout(y),
              Yt === u && (Yt = null),
              T && (T.destroy(), (T = null)),
              M && (M.destroy(), (M = null)),
              z && z.disable(e),
              C &&
                (C.destroy(),
                u.removeEventListener("mousemove", W),
                u.removeEventListener("mouseout", Y),
                u.removeEventListener("mouseover", Q),
                (C = null)),
              v.destroy(),
              n && (n.destroy(), (n = null)),
              Object(Ot.a)(e.get("id")),
              this.resizeListener &&
                (this.resizeListener.destroy(), delete this.resizeListener),
              O && Wt && bt.a.removeScrollHandler(F);
          });
      };
      function Zt(t, e, n) {
        return (Zt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, n) {
                var i = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = ee(t));

                  );
                  return t;
                })(t, e);
                if (i) {
                  var o = Object.getOwnPropertyDescriptor(i, e);
                  return o.get ? o.get.call(n) : o.value;
                }
              })(t, e, n || t);
      }
      function Kt(t) {
        return (Kt =
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
      function Jt(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Gt(t, e) {
        for (var n = 0; n < e.length; n++) {
          var i = e[n];
          (i.enumerable = i.enumerable || !1),
            (i.configurable = !0),
            "value" in i && (i.writable = !0),
            Object.defineProperty(t, i.key, i);
        }
      }
      function $t(t, e, n) {
        return e && Gt(t.prototype, e), n && Gt(t, n), t;
      }
      function te(t, e) {
        return !e || ("object" !== Kt(e) && "function" != typeof e) ? oe(t) : e;
      }
      function ee(t) {
        return (ee = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function ne(t, e) {
        if ("function" != typeof e && null !== e)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (t.prototype = Object.create(e && e.prototype, {
          constructor: { value: t, writable: !0, configurable: !0 },
        })),
          e && ie(t, e);
      }
      function ie(t, e) {
        return (ie =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      function oe(t) {
        if (void 0 === t)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return t;
      }
      var ae = /^change:(.+)$/;
      function re(t, e, n) {
        Object.keys(e).forEach(function (i) {
          i in e &&
            e[i] !== n[i] &&
            t.trigger("change:".concat(i), t, e[i], n[i]);
        });
      }
      function se(t, e) {
        t && t.off(null, null, e);
      }
      var le = (function (t) {
          function e(t, n) {
            var o;
            return (
              Jt(this, e),
              ((o = te(this, ee(e).call(this)))._model = t),
              (o._mediaModel = null),
              Object(i.g)(t.attributes, {
                altText: "",
                fullscreen: !1,
                logoWidth: 0,
                scrubbing: !1,
              }),
              t.on(
                "all",
                function (e, i, a, r) {
                  i === t && (i = oe(oe(o))),
                    (n && !n(e, i, a, r)) || o.trigger(e, i, a, r);
                },
                oe(oe(o))
              ),
              t.on(
                "change:mediaModel",
                function (t, e) {
                  o.mediaModel = e;
                },
                oe(oe(o))
              ),
              o
            );
          }
          return (
            ne(e, t),
            $t(e, [
              {
                key: "get",
                value: function (t) {
                  var e = this._mediaModel;
                  return e && t in e.attributes ? e.get(t) : this._model.get(t);
                },
              },
              {
                key: "set",
                value: function (t, e) {
                  return this._model.set(t, e);
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
                  se(this._model, this), se(this._mediaModel, this), this.off();
                },
              },
              {
                key: "mediaModel",
                set: function (t) {
                  var e = this,
                    n = this._mediaModel;
                  se(n, this),
                    (this._mediaModel = t),
                    t.on(
                      "all",
                      function (n, i, o, a) {
                        i === t && (i = e), e.trigger(n, i, o, a);
                      },
                      this
                    ),
                    n && re(this, t.attributes, n.attributes);
                },
              },
            ]),
            e
          );
        })(v.a),
        ce = (function (t) {
          function e(t) {
            var n;
            return (
              Jt(this, e),
              ((n = te(
                this,
                ee(e).call(this, t, function (t) {
                  var e = n._instreamModel;
                  if (e) {
                    var i = ae.exec(t);
                    if (i) if (i[1] in e.attributes) return !1;
                  }
                  return !0;
                })
              ))._instreamModel = null),
              (n._playerViewModel = new le(n._model)),
              t.on(
                "change:instream",
                function (t, e) {
                  n.instreamModel = e ? e.model : null;
                },
                oe(oe(n))
              ),
              n
            );
          }
          return (
            ne(e, t),
            $t(e, [
              {
                key: "get",
                value: function (t) {
                  var e = this._mediaModel;
                  if (e && t in e.attributes) return e.get(t);
                  var n = this._instreamModel;
                  return n && t in n.attributes ? n.get(t) : this._model.get(t);
                },
              },
              {
                key: "getVideo",
                value: function () {
                  var t = this._instreamModel;
                  return t && t.getVideo()
                    ? t.getVideo()
                    : Zt(ee(e.prototype), "getVideo", this).call(this);
                },
              },
              {
                key: "destroy",
                value: function () {
                  Zt(ee(e.prototype), "destroy", this).call(this),
                    se(this._instreamModel, this);
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
                set: function (t) {
                  var e = this,
                    n = this._instreamModel;
                  if (
                    (se(n, this),
                    this._model.off("change:mediaModel", null, this),
                    (this._instreamModel = t),
                    this.trigger("instreamMode", !!t),
                    t)
                  )
                    t.on(
                      "all",
                      function (n, i, o, a) {
                        i === t && (i = e), e.trigger(n, i, o, a);
                      },
                      this
                    ),
                      t.change(
                        "mediaModel",
                        function (t, n) {
                          e.mediaModel = n;
                        },
                        this
                      ),
                      re(this, t.attributes, this._model.attributes);
                  else if (n) {
                    this._model.change(
                      "mediaModel",
                      function (t, n) {
                        e.mediaModel = n;
                      },
                      this
                    );
                    var o = Object(i.g)(
                      {},
                      this._model.attributes,
                      n.attributes
                    );
                    re(this, this._model.attributes, o);
                  }
                },
              },
            ]),
            e
          );
        })(le);
      var ue,
        de,
        pe = n(64),
        we =
          (ue = window).URL && ue.URL.createObjectURL
            ? ue.URL
            : ue.webkitURL || ue.mozURL;
      function he(t, e) {
        var n = e.muted;
        return (
          de ||
            (de = new Blob(
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
          (t.muted = n),
          (t.src = we.createObjectURL(de)),
          t.play() || Object(pe.a)(t)
        );
      }
      var fe = "autoplayEnabled",
        je = "autoplayMuted",
        ge = "autoplayDisabled",
        be = {};
      var me = n(65);
      function ve(t) {
        return (
          (t = t || window.event) &&
          /^(?:mouse|pointer|touch|gesture|click|key)/.test(t.type)
        );
      }
      var ye = n(24),
        ke = "tabHidden",
        xe = "tabVisible",
        Oe = function (t) {
          var e = 0;
          return function (n) {
            var i = n.position;
            i > e && t(), (e = i);
          };
        };
      function Ce(t, e) {
        e.off(d.N, t._onPlayAttempt),
          e.off(d.fb, t._triggerFirstFrame),
          e.off(d.S, t._onTime),
          t.off("change:activeTab", t._onTabVisible);
      }
      var Me = function (t, e) {
        t.change("mediaModel", function (t, n, i) {
          t._qoeItem && i && t._qoeItem.end(i.get("mediaState")),
            (t._qoeItem = new ye.a()),
            (t._qoeItem.getFirstFrame = function () {
              var t = this.between(d.N, d.H),
                e = this.between(xe, d.H);
              return e > 0 && e < t ? e : t;
            }),
            t._qoeItem.tick(d.db),
            t._qoeItem.start(n.get("mediaState")),
            (function (t, e) {
              t._onTabVisible && Ce(t, e);
              var n = !1;
              (t._triggerFirstFrame = function () {
                if (!n) {
                  n = !0;
                  var i = t._qoeItem;
                  i.tick(d.H);
                  var o = i.getFirstFrame();
                  if ((e.trigger(d.H, { loadTime: o }), e.mediaController)) {
                    var a = e.mediaController.mediaModel;
                    a.off("change:".concat(d.U), null, a),
                      a.change(
                        d.U,
                        function (t, n) {
                          n && e.trigger(d.U, n);
                        },
                        a
                      );
                  }
                  Ce(t, e);
                }
              }),
                (t._onTime = Oe(t._triggerFirstFrame)),
                (t._onPlayAttempt = function () {
                  t._qoeItem.tick(d.N);
                }),
                (t._onTabVisible = function (e, n) {
                  n ? t._qoeItem.tick(xe) : t._qoeItem.tick(ke);
                }),
                t.on("change:activeTab", t._onTabVisible),
                e.on(d.N, t._onPlayAttempt),
                e.once(d.fb, t._triggerFirstFrame),
                e.on(d.S, t._onTime);
            })(t, e),
            n.on("change:mediaState", function (e, n, i) {
              n !== i && (t._qoeItem.end(i), t._qoeItem.start(n));
            });
        });
      };
      function Te(t) {
        return (Te =
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
      var Se = function () {},
        _e = function () {};
      Object(i.g)(Se.prototype, {
        setup: function (t, e, n, w, f, b) {
          var v,
            y,
            k,
            x,
            O = this,
            C = this,
            M = (C._model = new A()),
            T = !1,
            S = !1,
            _ = null,
            E = j(H),
            z = j(_e);
          (C.originalContainer = C.currentContainer = n),
            (C._events = w),
            (C.trigger = function (t, e) {
              var n = (function (t, e, n) {
                var o = n;
                switch (e) {
                  case "time":
                  case "beforePlay":
                  case "pause":
                  case "play":
                  case "ready":
                    var a = t.get("viewable");
                    void 0 !== a && (o = Object(i.g)({}, n, { viewable: a }));
                }
                return o;
              })(M, t, e);
              return h.a.trigger.call(this, t, n);
            });
          var P = new s.a(C, ["trigger"], function () {
              return !0;
            }),
            I = function (t, e) {
              C.trigger(t, e);
            };
          M.setup(t);
          var R = M.get("backgroundLoading"),
            L = new ce(M);
          (v = this._view = new Xt(e, L)).on(
            "all",
            function (t, e) {
              (e && e.doNotForward) || I(t, e);
            },
            C
          );
          var B = (this._programController = new Y(M, b));
          ut(),
            B.on("all", I, C)
              .on(
                "subtitlesTracks",
                function (t) {
                  y.setSubtitlesTracks(t.tracks);
                  var e = y.getCurrentIndex();
                  e > 0 && rt(e, t.tracks);
                },
                C
              )
              .on(
                d.F,
                function () {
                  Promise.resolve().then(at);
                },
                C
              )
              .on(d.G, C.triggerError, C),
            Me(M, B),
            M.on(d.w, C.triggerError, C),
            M.on(
              "change:state",
              function (t, e, n) {
                X() || Z.call(O, t, e, n);
              },
              this
            ),
            M.on("change:castState", function (t, e) {
              C.trigger(d.m, e);
            }),
            M.on("change:fullscreen", function (t, e) {
              C.trigger(d.y, { fullscreen: e }),
                e && t.set("playOnViewable", !1);
            }),
            M.on("change:volume", function (t, e) {
              C.trigger(d.V, { volume: e });
            }),
            M.on("change:mute", function (t) {
              C.trigger(d.M, { mute: t.getMute() });
            }),
            M.on("change:playbackRate", function (t, e) {
              C.trigger(d.ab, { playbackRate: e, position: t.get("position") });
            });
          var V = function t(e, n) {
            ("clickthrough" !== n && "interaction" !== n && "external" !== n) ||
              (M.set("playOnViewable", !1),
              M.off("change:playReason change:pauseReason", t));
          };
          function N(t, e) {
            Object(i.t)(e) || M.set("viewable", Math.round(e));
          }
          function H() {
            dt &&
              (!0 !== M.get("autostart") ||
                M.get("playOnViewable") ||
                $("autostart"),
              dt.flush());
          }
          function F(t, e) {
            C.trigger("viewable", { viewable: e }), q();
          }
          function q() {
            if (
              (o.a[0] === e || 1 === M.get("viewable")) &&
              "idle" === M.get("state") &&
              !1 === M.get("autostart")
            )
              if (!b.primed() && m.OS.android) {
                var t = b.getTestElement(),
                  n = C.getMute();
                Promise.resolve()
                  .then(function () {
                    return he(t, { muted: n });
                  })
                  .then(function () {
                    "idle" === M.get("state") && B.preloadVideo();
                  })
                  .catch(_e);
              } else B.preloadVideo();
          }
          function D(t) {
            (C._instreamAdapter.noResume = !t), t || et({ reason: "viewable" });
          }
          function U(t) {
            t || (C.pause({ reason: "viewable" }), M.set("playOnViewable", !t));
          }
          function W(t, e) {
            var n = X();
            if (t.get("playOnViewable")) {
              if (e) {
                var i = t.get("autoPause").pauseAds,
                  o = t.get("pauseReason");
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
              m.OS.mobile && n && D(e);
            }
          }
          function Q(t, e) {
            var n = t.get("state"),
              i = X(),
              o = t.get("playReason");
            i
              ? t.get("autoPause").pauseAds
                ? U(e)
                : D(e)
              : n === d.pb || n === d.jb
              ? U(e)
              : n === d.mb &&
                "playlist" === o &&
                t.once("change:state", function () {
                  U(e);
                });
          }
          function X() {
            var t = C._instreamAdapter;
            return !!t && t.getState();
          }
          function K() {
            var t = X();
            return t || M.get("state");
          }
          function J(t) {
            if ((E.cancel(), (S = !1), M.get("state") === d.lb))
              return Promise.resolve();
            var n = G(t);
            return (
              M.set("playReason", n),
              X()
                ? (e.pauseAd(!1, t), Promise.resolve())
                : (M.get("state") === d.kb && (tt(!0), C.setItemIndex(0)),
                  !T &&
                  ((T = !0),
                  C.trigger(d.C, {
                    playReason: n,
                    startTime:
                      t && t.startTime
                        ? t.startTime
                        : M.get("playlistItem").starttime,
                  }),
                  (T = !1),
                  ve() && !b.primed() && b.prime(),
                  "playlist" === n &&
                    M.get("autoPause").viewability &&
                    Q(M, M.get("viewable")),
                  x)
                    ? (ve() && !R && M.get("mediaElement").load(),
                      (x = !1),
                      (k = null),
                      Promise.resolve())
                    : B.playVideo(n).then(b.played))
            );
          }
          function G(t) {
            return t && t.reason ? t.reason : "unknown";
          }
          function $(t) {
            if (K() === d.mb) {
              E = j(H);
              var e = M.get("advertising");
              (function (t, e) {
                var n = e.cancelable,
                  i = e.muted,
                  o = void 0 !== i && i,
                  a = e.allowMuted,
                  r = void 0 !== a && a,
                  s = e.timeout,
                  l = void 0 === s ? 1e4 : s,
                  c = t.getTestElement(),
                  u = o ? "muted" : "".concat(r);
                be[u] ||
                  (be[u] = he(c, { muted: o })
                    .catch(function (t) {
                      if (!n.cancelled() && !1 === o && r)
                        return he(c, { muted: (o = !0) });
                      throw t;
                    })
                    .then(function () {
                      return o ? ((be[u] = null), je) : fe;
                    })
                    .catch(function (t) {
                      throw (
                        (clearTimeout(d), (be[u] = null), (t.reason = ge), t)
                      );
                    }));
                var d,
                  p = be[u].then(function (t) {
                    if ((clearTimeout(d), n.cancelled())) {
                      var e = new Error("Autoplay test was cancelled");
                      throw ((e.reason = "cancelled"), e);
                    }
                    return t;
                  }),
                  w = new Promise(function (t, e) {
                    d = setTimeout(function () {
                      be[u] = null;
                      var t = new Error("Autoplay test timed out");
                      (t.reason = "timeout"), e(t);
                    }, l);
                  });
                return Promise.race([p, w]);
              })(b, {
                cancelable: E,
                muted: C.getMute(),
                allowMuted: !e || e.autoplayadsmuted,
              })
                .then(function (e) {
                  return (
                    M.set("canAutoplay", e),
                    e !== je ||
                      C.getMute() ||
                      (M.set("autostartMuted", !0),
                      ut(),
                      M.once("change:autostartMuted", function (t) {
                        t.off("change:viewable", W),
                          C.trigger(d.M, { mute: M.getMute() });
                      })),
                    C.getMute() &&
                      M.get("enableDefaultCaptions") &&
                      y.selectDefaultIndex(1),
                    J({ reason: t }).catch(function () {
                      C._instreamAdapter || M.set("autostartFailed", !0),
                        (k = null);
                    })
                  );
                })
                .catch(function (t) {
                  if (
                    (M.set("canAutoplay", ge),
                    M.set("autostart", !1),
                    !E.cancelled())
                  ) {
                    var e = Object(g.w)(t);
                    C.trigger(d.h, { reason: t.reason, code: e, error: t });
                  }
                });
            }
          }
          function tt(t) {
            if ((E.cancel(), dt.empty(), X())) {
              var e = C._instreamAdapter;
              return (
                e && (e.noResume = !0),
                void (k = function () {
                  return B.stopVideo();
                })
              );
            }
            (k = null),
              !t && (S = !0),
              T && (x = !0),
              M.set("errorEvent", void 0),
              B.stopVideo();
          }
          function et(t) {
            var e = G(t);
            M.set("pauseReason", e), M.set("playOnViewable", "viewable" === e);
          }
          function nt(t) {
            (k = null), E.cancel();
            var n = X();
            if (n && n !== d.ob) return et(t), void e.pauseAd(!0, t);
            switch (M.get("state")) {
              case d.lb:
                return;
              case d.pb:
              case d.jb:
                et(t), B.pause();
                break;
              default:
                T && (x = !0);
            }
          }
          function it(t, e) {
            tt(!0), C.setItemIndex(t), C.play(e);
          }
          function ot(t) {
            it(M.get("item") + 1, t);
          }
          function at() {
            C.completeCancelled() ||
              ((k = C.completeHandler),
              C.shouldAutoAdvance()
                ? C.nextItem()
                : M.get("repeat")
                ? ot({ reason: "repeat" })
                : (m.OS.iOS && lt(!1),
                  M.set("playOnViewable", !1),
                  M.set("state", d.kb),
                  C.trigger(d.cb, {})));
          }
          function rt(t, e) {
            (t = parseInt(t, 10) || 0),
              M.persistVideoSubtitleTrack(t, e),
              (B.subtitles = t),
              C.trigger(d.k, { tracks: st(), track: t });
          }
          function st() {
            return y.getCaptionsList();
          }
          function lt(t) {
            Object(i.n)(t) || (t = !M.get("fullscreen")),
              M.set("fullscreen", t),
              C._instreamAdapter &&
                C._instreamAdapter._adModel &&
                C._instreamAdapter._adModel.set("fullscreen", t);
          }
          function ut() {
            (B.mute = M.getMute()), (B.volume = M.get("volume"));
          }
          M.on("change:playReason change:pauseReason", V),
            C.on(d.c, function (t) {
              return V(0, t.playReason);
            }),
            C.on(d.b, function (t) {
              return V(0, t.pauseReason);
            }),
            M.on("change:scrubbing", function (t, e) {
              e
                ? ((_ = M.get("state") !== d.ob), nt())
                : _ && J({ reason: "interaction" });
            }),
            M.on("change:captionsList", function (t, e) {
              C.trigger(d.l, { tracks: e, track: M.get("captionsIndex") || 0 });
            }),
            M.on("change:mediaModel", function (t, e) {
              var n = this;
              t.set("errorEvent", void 0),
                e.change(
                  "mediaState",
                  function (e, n) {
                    var i;
                    t.get("errorEvent") ||
                      t.set(d.bb, (i = n) === d.nb || i === d.qb ? d.jb : i);
                  },
                  this
                ),
                e.change(
                  "duration",
                  function (e, n) {
                    if (0 !== n) {
                      var i = t.get("minDvrWindow"),
                        o = Object(me.b)(n, i);
                      t.setStreamType(o);
                    }
                  },
                  this
                );
              var i = t.get("item") + 1,
                o = "autoplay" === (t.get("related") || {}).oncomplete,
                a = t.get("playlist")[i];
              if ((a || o) && R) {
                e.on(
                  "change:position",
                  function t(i, r) {
                    var s = a && !a.daiSetting,
                      l = e.get("duration");
                    s && r && l > 0 && r >= l - p.b
                      ? (e.off("change:position", t, n), B.backgroundLoad(a))
                      : o && (a = M.get("nextUp"));
                  },
                  this
                );
              }
            }),
            (y = new wt(M)).on("all", I, C),
            L.on("viewSetup", function (t) {
              Object(a.b)(O, t);
            }),
            (this.playerReady = function () {
              v.once(d.hb, function () {
                try {
                  !(function () {
                    M.change("visibility", N),
                      P.off(),
                      C.trigger(d.gb, { setupTime: 0 }),
                      M.change("playlist", function (t, e) {
                        if (e.length) {
                          var n = { playlist: e },
                            o = M.get("feedData");
                          o && (n.feedData = Object(i.g)({}, o)),
                            C.trigger(d.eb, n);
                        }
                      }),
                      M.change("playlistItem", function (t, e) {
                        if (e) {
                          var n = e.title,
                            i = e.image;
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
                            } catch (t) {}
                          t.set("cues", []),
                            C.trigger(d.db, { index: M.get("item"), item: e });
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
                            function (t) {
                              t.off("change:viewable", W);
                            }
                          );
                    H(),
                      M.on("change:itemReady", function (t, e) {
                        e && dt.flush();
                      });
                  })();
                } catch (t) {
                  C.triggerError(Object(g.v)(g.m, g.a, t));
                }
              }),
                v.init();
            }),
            (this.preload = q),
            (this.load = function (t, e) {
              var n,
                i = C._instreamAdapter;
              switch (
                (i && (i.noResume = !0),
                C.trigger("destroyPlugin", {}),
                tt(!0),
                E.cancel(),
                (E = j(H)),
                z.cancel(),
                ve() && b.prime(),
                Te(t))
              ) {
                case "string":
                  (M.attributes.item = 0),
                    (M.attributes.itemReady = !1),
                    (z = j(function (t) {
                      if (t)
                        return C.updatePlaylist(Object(c.a)(t.playlist), t);
                    })),
                    (n = (function (t) {
                      var e = this;
                      return new Promise(function (n, i) {
                        var o = new l.a();
                        o.on(d.eb, function (t) {
                          n(t);
                        }),
                          o.on(d.w, i, e),
                          o.load(t);
                      });
                    })(t).then(z.async));
                  break;
                case "object":
                  (M.attributes.item = 0),
                    (n = C.updatePlaylist(Object(c.a)(t), e || {}));
                  break;
                case "number":
                  n = C.setItemIndex(t);
                  break;
                default:
                  return;
              }
              n.catch(function (t) {
                C.triggerError(Object(g.u)(t, g.c));
              }),
                n.then(E.async).catch(_e);
            }),
            (this.play = function (t) {
              return J(t).catch(_e);
            }),
            (this.pause = nt),
            (this.seek = function (t, e) {
              var n = M.get("state");
              if (n !== d.lb) {
                B.position = t;
                var i = n === d.mb;
                M.get("scrubbing") ||
                  (!i && n !== d.kb) ||
                  (i && ((e = e || {}).startTime = t), this.play(e));
              }
            }),
            (this.stop = tt),
            (this.playlistItem = it),
            (this.playlistNext = ot),
            (this.playlistPrev = function (t) {
              it(M.get("item") - 1, t);
            }),
            (this.setCurrentCaptions = rt),
            (this.setCurrentQuality = function (t) {
              B.quality = t;
            }),
            (this.setFullscreen = lt),
            (this.getCurrentQuality = function () {
              return B.quality;
            }),
            (this.getQualityLevels = function () {
              return B.qualities;
            }),
            (this.setCurrentAudioTrack = function (t) {
              B.audioTrack = t;
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
            (this.getCaptionsList = st),
            (this.getVisualQuality = function () {
              var t = this._model.get("mediaModel");
              return t ? t.get(d.U) : null;
            }),
            (this.getConfig = function () {
              return this._model ? this._model.getConfiguration() : void 0;
            }),
            (this.getState = K),
            (this.next = _e),
            (this.completeHandler = at),
            (this.completeCancelled = function () {
              return (
                ((t = M.get("state")) !== d.mb && t !== d.kb && t !== d.lb) ||
                (!!S && ((S = !1), !0))
              );
              var t;
            }),
            (this.shouldAutoAdvance = function () {
              return M.get("item") !== M.get("playlist").length - 1;
            }),
            (this.nextItem = function () {
              ot({ reason: "playlist" });
            }),
            (this.setConfig = function (t) {
              !(function (t, e) {
                var n = t._model,
                  i = n.attributes;
                e.height &&
                  ((e.height = Object(r.b)(e.height)),
                  (e.width = e.width || i.width)),
                  e.width &&
                    ((e.width = Object(r.b)(e.width)),
                    e.aspectratio
                      ? ((i.width = e.width), delete e.width)
                      : (e.height = i.height)),
                  e.width &&
                    e.height &&
                    !e.aspectratio &&
                    t._view.resize(e.width, e.height),
                  Object.keys(e).forEach(function (o) {
                    var a = e[o];
                    if (void 0 !== a)
                      switch (o) {
                        case "aspectratio":
                          n.set(o, Object(r.a)(a, i.width));
                          break;
                        case "autostart":
                          !(function (t, e, n) {
                            t.setAutoStart(n),
                              "idle" === t.get("state") &&
                                !0 === n &&
                                e.play({ reason: "autostart" });
                          })(n, t, a);
                          break;
                        case "mute":
                          t.setMute(a);
                          break;
                        case "volume":
                          t.setVolume(a);
                          break;
                        case "playbackRateControls":
                        case "playbackRates":
                        case "repeat":
                        case "stretching":
                          n.set(o, a);
                      }
                  });
              })(C, t);
            }),
            (this.setItemIndex = function (t) {
              B.stopVideo();
              var e = M.get("playlist").length;
              return (
                (t = (parseInt(t, 10) || 0) % e) < 0 && (t += e),
                B.setActiveItem(t).catch(function (t) {
                  t.code >= 151 && t.code <= 162 && (t = Object(g.u)(t, g.e)),
                    O.triggerError(Object(g.v)(g.k, g.d, t));
                })
              );
            }),
            (this.detachMedia = function () {
              if (
                (T && (x = !0),
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
            (this.routeEvents = function (t) {
              return B.routeEvents(t);
            }),
            (this.forwardEvents = function () {
              return B.forwardEvents();
            }),
            (this.playVideo = function (t) {
              return B.playVideo(t);
            }),
            (this.stopVideo = function () {
              return B.stopVideo();
            }),
            (this.castVideo = function (t, e) {
              return B.castVideo(t, e);
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
            (this.setVolume = function (t) {
              M.setVolume(t), ut();
            }),
            (this.setMute = function (t) {
              M.setMute(t), ut();
            }),
            (this.setPlaybackRate = function (t) {
              M.setPlaybackRate(t);
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
            (this.addButton = function (t, e, n, i, o) {
              var a = M.get("customButtons") || [],
                r = !1,
                s = { img: t, tooltip: e, callback: n, id: i, btnClass: o };
              (a = a.reduce(function (t, e) {
                return e.id === i ? ((r = !0), t.push(s)) : t.push(e), t;
              }, [])),
                r || a.unshift(s),
                M.set("customButtons", a);
            }),
            (this.removeButton = function (t) {
              var e = M.get("customButtons") || [];
              (e = e.filter(function (e) {
                return e.id !== t;
              })),
                M.set("customButtons", e);
            }),
            (this.resize = v.resize),
            (this.getSafeRegion = v.getSafeRegion),
            (this.setCaptions = v.setCaptions),
            (this.checkBeforePlay = function () {
              return T;
            }),
            (this.setControls = function (t) {
              Object(i.n)(t) || (t = !M.get("controls")),
                M.set("controls", t),
                (B.controls = t);
            }),
            (this.addCues = function (t) {
              this.setCues(M.get("cues").concat(t));
            }),
            (this.setCues = function (t) {
              M.set("cues", t);
            }),
            (this.updatePlaylist = function (t, e) {
              try {
                var n = Object(c.b)(t, M, e);
                Object(c.e)(n);
                var o = Object(i.g)({}, e);
                delete o.playlist, M.set("feedData", o), M.set("playlist", n);
              } catch (t) {
                return Promise.reject(t);
              }
              return this.setItemIndex(M.get("item"));
            }),
            (this.setPlaylistItem = function (t, e) {
              (e = Object(c.d)(M, new u.a(e), e.feedData || {})) &&
                ((M.get("playlist")[t] = e),
                t === M.get("item") &&
                  "idle" === M.get("state") &&
                  this.setItemIndex(t));
            }),
            (this.playerDestroy = function () {
              this.off(),
                this.stop(),
                Object(a.b)(this, this.originalContainer),
                v && v.destroy(),
                M && M.destroy(),
                dt && dt.destroy(),
                y && y.destroy(),
                B && B.destroy(),
                this.instreamDestroy();
            }),
            (this.isBeforePlay = this.checkBeforePlay),
            (this.createInstream = function () {
              return (
                this.instreamDestroy(),
                (this._instreamAdapter = new ct(this, M, v, b)),
                this._instreamAdapter
              );
            }),
            (this.instreamDestroy = function () {
              C._instreamAdapter &&
                (C._instreamAdapter.destroy(), (C._instreamAdapter = null));
            });
          var dt = new s.a(
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
          dt.queue.push.apply(dt.queue, f), v.setup();
        },
        get: function (t) {
          if (t in y.a) {
            var e = this._model.get("mediaModel");
            return e ? e.get(t) : y.a[t];
          }
          return this._model.get(t);
        },
        getContainer: function () {
          return this.currentContainer || this.originalContainer;
        },
        getMute: function () {
          return this._model.getMute();
        },
        triggerError: function (t) {
          var e = this._model;
          (t.message = e.get("localization").errors[t.key]),
            delete t.key,
            e.set("errorEvent", t),
            e.set("state", d.lb),
            e.once(
              "change:state",
              function () {
                this.set("errorEvent", void 0);
              },
              e
            ),
            this.trigger(d.w, t);
        },
      });
      e.default = Se;
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
    function (t, e) {
      !(function (t, e) {
        "use strict";
        if (
          "IntersectionObserver" in t &&
          "IntersectionObserverEntry" in t &&
          "intersectionRatio" in t.IntersectionObserverEntry.prototype
        )
          "isIntersecting" in t.IntersectionObserverEntry.prototype ||
            Object.defineProperty(
              t.IntersectionObserverEntry.prototype,
              "isIntersecting",
              {
                get: function () {
                  return this.intersectionRatio > 0;
                },
              }
            );
        else {
          var n = [];
          (o.prototype.THROTTLE_TIMEOUT = 100),
            (o.prototype.POLL_INTERVAL = null),
            (o.prototype.USE_MUTATION_OBSERVER = !0),
            (o.prototype.observe = function (t) {
              if (
                !this._observationTargets.some(function (e) {
                  return e.element == t;
                })
              ) {
                if (!t || 1 != t.nodeType)
                  throw new Error("target must be an Element");
                this._registerInstance(),
                  this._observationTargets.push({ element: t, entry: null }),
                  this._monitorIntersections(),
                  this._checkForIntersections();
              }
            }),
            (o.prototype.unobserve = function (t) {
              (this._observationTargets = this._observationTargets.filter(
                function (e) {
                  return e.element != t;
                }
              )),
                this._observationTargets.length ||
                  (this._unmonitorIntersections(), this._unregisterInstance());
            }),
            (o.prototype.disconnect = function () {
              (this._observationTargets = []),
                this._unmonitorIntersections(),
                this._unregisterInstance();
            }),
            (o.prototype.takeRecords = function () {
              var t = this._queuedEntries.slice();
              return (this._queuedEntries = []), t;
            }),
            (o.prototype._initThresholds = function (t) {
              var e = t || [0];
              return (
                Array.isArray(e) || (e = [e]),
                e.sort().filter(function (t, e, n) {
                  if ("number" != typeof t || isNaN(t) || t < 0 || t > 1)
                    throw new Error(
                      "threshold must be a number between 0 and 1 inclusively"
                    );
                  return t !== n[e - 1];
                })
              );
            }),
            (o.prototype._parseRootMargin = function (t) {
              var e = (t || "0px").split(/\s+/).map(function (t) {
                var e = /^(-?\d*\.?\d+)(px|%)$/.exec(t);
                if (!e)
                  throw new Error(
                    "rootMargin must be specified in pixels or percent"
                  );
                return { value: parseFloat(e[1]), unit: e[2] };
              });
              return (
                (e[1] = e[1] || e[0]),
                (e[2] = e[2] || e[0]),
                (e[3] = e[3] || e[1]),
                e
              );
            }),
            (o.prototype._monitorIntersections = function () {
              this._monitoringIntersections ||
                ((this._monitoringIntersections = !0),
                this.POLL_INTERVAL
                  ? (this._monitoringInterval = setInterval(
                      this._checkForIntersections,
                      this.POLL_INTERVAL
                    ))
                  : (a(t, "resize", this._checkForIntersections, !0),
                    a(e, "scroll", this._checkForIntersections, !0),
                    this.USE_MUTATION_OBSERVER &&
                      "MutationObserver" in t &&
                      ((this._domObserver = new MutationObserver(
                        this._checkForIntersections
                      )),
                      this._domObserver.observe(e, {
                        attributes: !0,
                        childList: !0,
                        characterData: !0,
                        subtree: !0,
                      }))));
            }),
            (o.prototype._unmonitorIntersections = function () {
              this._monitoringIntersections &&
                ((this._monitoringIntersections = !1),
                clearInterval(this._monitoringInterval),
                (this._monitoringInterval = null),
                r(t, "resize", this._checkForIntersections, !0),
                r(e, "scroll", this._checkForIntersections, !0),
                this._domObserver &&
                  (this._domObserver.disconnect(), (this._domObserver = null)));
            }),
            (o.prototype._checkForIntersections = function () {
              var e = this._rootIsInDom(),
                n = e
                  ? this._getRootRect()
                  : {
                      top: 0,
                      bottom: 0,
                      left: 0,
                      right: 0,
                      width: 0,
                      height: 0,
                    };
              this._observationTargets.forEach(function (o) {
                var a = o.element,
                  r = s(a),
                  l = this._rootContainsTarget(a),
                  c = o.entry,
                  u = e && l && this._computeTargetAndRootIntersection(a, n),
                  d = (o.entry = new i({
                    time: t.performance && performance.now && performance.now(),
                    target: a,
                    boundingClientRect: r,
                    rootBounds: n,
                    intersectionRect: u,
                  }));
                c
                  ? e && l
                    ? this._hasCrossedThreshold(c, d) &&
                      this._queuedEntries.push(d)
                    : c && c.isIntersecting && this._queuedEntries.push(d)
                  : this._queuedEntries.push(d);
              }, this),
                this._queuedEntries.length &&
                  this._callback(this.takeRecords(), this);
            }),
            (o.prototype._computeTargetAndRootIntersection = function (n, i) {
              if ("none" != t.getComputedStyle(n).display) {
                for (
                  var o, a, r, l, u, d, p, w, h = s(n), f = c(n), j = !1;
                  !j;

                ) {
                  var g = null,
                    b = 1 == f.nodeType ? t.getComputedStyle(f) : {};
                  if ("none" == b.display) return;
                  if (
                    (f == this.root || f == e
                      ? ((j = !0), (g = i))
                      : f != e.body &&
                        f != e.documentElement &&
                        "visible" != b.overflow &&
                        (g = s(f)),
                    g &&
                      ((o = g),
                      (a = h),
                      (r = void 0),
                      (l = void 0),
                      (u = void 0),
                      (d = void 0),
                      (p = void 0),
                      (w = void 0),
                      (r = Math.max(o.top, a.top)),
                      (l = Math.min(o.bottom, a.bottom)),
                      (u = Math.max(o.left, a.left)),
                      (d = Math.min(o.right, a.right)),
                      (w = l - r),
                      !(h = (p = d - u) >= 0 &&
                        w >= 0 && {
                          top: r,
                          bottom: l,
                          left: u,
                          right: d,
                          width: p,
                          height: w,
                        })))
                  )
                    break;
                  f = c(f);
                }
                return h;
              }
            }),
            (o.prototype._getRootRect = function () {
              var t;
              if (this.root) t = s(this.root);
              else {
                var n = e.documentElement,
                  i = e.body;
                t = {
                  top: 0,
                  left: 0,
                  right: n.clientWidth || i.clientWidth,
                  width: n.clientWidth || i.clientWidth,
                  bottom: n.clientHeight || i.clientHeight,
                  height: n.clientHeight || i.clientHeight,
                };
              }
              return this._expandRectByRootMargin(t);
            }),
            (o.prototype._expandRectByRootMargin = function (t) {
              var e = this._rootMarginValues.map(function (e, n) {
                  return "px" == e.unit
                    ? e.value
                    : (e.value * (n % 2 ? t.width : t.height)) / 100;
                }),
                n = {
                  top: t.top - e[0],
                  right: t.right + e[1],
                  bottom: t.bottom + e[2],
                  left: t.left - e[3],
                };
              return (
                (n.width = n.right - n.left), (n.height = n.bottom - n.top), n
              );
            }),
            (o.prototype._hasCrossedThreshold = function (t, e) {
              var n = t && t.isIntersecting ? t.intersectionRatio || 0 : -1,
                i = e.isIntersecting ? e.intersectionRatio || 0 : -1;
              if (n !== i)
                for (var o = 0; o < this.thresholds.length; o++) {
                  var a = this.thresholds[o];
                  if (a == n || a == i || a < n != a < i) return !0;
                }
            }),
            (o.prototype._rootIsInDom = function () {
              return !this.root || l(e, this.root);
            }),
            (o.prototype._rootContainsTarget = function (t) {
              return l(this.root || e, t);
            }),
            (o.prototype._registerInstance = function () {
              n.indexOf(this) < 0 && n.push(this);
            }),
            (o.prototype._unregisterInstance = function () {
              var t = n.indexOf(this);
              -1 != t && n.splice(t, 1);
            }),
            (t.IntersectionObserver = o),
            (t.IntersectionObserverEntry = i);
        }
        function i(t) {
          (this.time = t.time),
            (this.target = t.target),
            (this.rootBounds = t.rootBounds),
            (this.boundingClientRect = t.boundingClientRect),
            (this.intersectionRect = t.intersectionRect || {
              top: 0,
              bottom: 0,
              left: 0,
              right: 0,
              width: 0,
              height: 0,
            }),
            (this.isIntersecting = !!t.intersectionRect);
          var e = this.boundingClientRect,
            n = e.width * e.height,
            i = this.intersectionRect,
            o = i.width * i.height;
          this.intersectionRatio = n ? o / n : this.isIntersecting ? 1 : 0;
        }
        function o(t, e) {
          var n,
            i,
            o,
            a = e || {};
          if ("function" != typeof t)
            throw new Error("callback must be a function");
          if (a.root && 1 != a.root.nodeType)
            throw new Error("root must be an Element");
          (this._checkForIntersections =
            ((n = this._checkForIntersections.bind(this)),
            (i = this.THROTTLE_TIMEOUT),
            (o = null),
            function () {
              o ||
                (o = setTimeout(function () {
                  n(), (o = null);
                }, i));
            })),
            (this._callback = t),
            (this._observationTargets = []),
            (this._queuedEntries = []),
            (this._rootMarginValues = this._parseRootMargin(a.rootMargin)),
            (this.thresholds = this._initThresholds(a.threshold)),
            (this.root = a.root || null),
            (this.rootMargin = this._rootMarginValues
              .map(function (t) {
                return t.value + t.unit;
              })
              .join(" "));
        }
        function a(t, e, n, i) {
          "function" == typeof t.addEventListener
            ? t.addEventListener(e, n, i || !1)
            : "function" == typeof t.attachEvent && t.attachEvent("on" + e, n);
        }
        function r(t, e, n, i) {
          "function" == typeof t.removeEventListener
            ? t.removeEventListener(e, n, i || !1)
            : "function" == typeof t.detatchEvent &&
              t.detatchEvent("on" + e, n);
        }
        function s(t) {
          var e;
          try {
            e = t.getBoundingClientRect();
          } catch (t) {}
          return e
            ? ((e.width && e.height) ||
                (e = {
                  top: e.top,
                  right: e.right,
                  bottom: e.bottom,
                  left: e.left,
                  width: e.right - e.left,
                  height: e.bottom - e.top,
                }),
              e)
            : { top: 0, bottom: 0, left: 0, right: 0, width: 0, height: 0 };
        }
        function l(t, e) {
          for (var n = e; n; ) {
            if (n == t) return !0;
            n = c(n);
          }
          return !1;
        }
        function c(t) {
          var e = t.parentNode;
          return e && 11 == e.nodeType && e.host ? e.host : e;
        }
      })(window, document);
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
        p = function (t) {
          var e,
            s,
            p,
            w,
            h,
            f,
            j,
            g,
            b,
            m = this,
            v = t.player;
          function y() {
            Object(o.o)(e.fontSize) &&
              (v.get("containerHeight")
                ? (g =
                    (d.fontScale * (e.userFontScale || 1) * e.fontSize) /
                    d.fontSize)
                : v.once("change:containerHeight", y, this));
          }
          function k() {
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
                            s = v.get("containerWidth");
                          if (v.get("fullscreen") && a.OS.mobile) {
                            var l = window.screen;
                            l.orientation &&
                              ((r = l.availHeight), (s = l.availWidth));
                          }
                          if (s && r && n && i)
                            return (s / r > o ? r : (i * s) / n) * g;
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
                : Object(c.d)(h, { fontSize: e });
            }
          }
          function x(t, e, n) {
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
          ((h = document.createElement("div")).className =
            "jw-captions jw-reset"),
            (this.show = function () {
              Object(u.a)(h, "jw-captions-enabled");
            }),
            (this.hide = function () {
              Object(u.o)(h, "jw-captions-enabled");
            }),
            (this.populate = function (t) {
              v.get("renderCaptionsNatively") ||
                ((p = []),
                (s = t),
                t ? this.selectCues(t, w) : this.renderCues());
            }),
            (this.resize = function () {
              k(), this.renderCues(!0);
            }),
            (this.renderCues = function (t) {
              (t = !!t), i && i.processCues(window, p, h, t);
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
              Object(u.g)(h);
            }),
            (this.setup = function (t, n) {
              (f = document.createElement("div")),
                (j = document.createElement("span")),
                (f.className = "jw-captions-window jw-reset"),
                (j.className = "jw-captions-text jw-reset"),
                (e = Object(o.g)({}, d, n)),
                (g = d.fontScale);
              var i = function () {
                if (!v.get("renderCaptionsNatively")) {
                  y(e.fontSize);
                  var n = e.windowColor,
                    i = e.windowOpacity,
                    o = e.edgeStyle;
                  b = {};
                  var r = {};
                  !(function (t, e) {
                    var n = e.color,
                      i = e.fontOpacity;
                    (n || i !== d.fontOpacity) &&
                      (t.color = Object(c.c)(n || "#ffffff", i));
                    if (e.back) {
                      var o = e.backgroundColor,
                        a = e.backgroundOpacity;
                      (o === d.backgroundColor && a === d.backgroundOpacity) ||
                        (t.backgroundColor = Object(c.c)(o, a));
                    } else t.background = "transparent";
                    e.fontFamily && (t.fontFamily = e.fontFamily);
                    e.fontStyle && (t.fontStyle = e.fontStyle);
                    e.fontWeight && (t.fontWeight = e.fontWeight);
                    e.textDecoration && (t.textDecoration = e.textDecoration);
                  })(r, e),
                    (n || i !== d.windowOpacity) &&
                      (b.backgroundColor = Object(c.c)(n || "#000000", i)),
                    x(o, r, e.fontOpacity),
                    e.back || null !== o || x("uniform", r),
                    Object(c.d)(f, b),
                    Object(c.d)(j, r),
                    (function (t, e) {
                      k(),
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
                f.appendChild(j),
                h.appendChild(f),
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
              return h;
            }),
            (this.destroy = function () {
              v.off(null, null, this), this.off();
            });
          var O = function (t) {
            (w = t), m.selectCues(s, w);
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
              function (t) {
                (p = []), O(t);
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
                        a.trigger(l.tb, t);
                      }),
                    e.off("change:captionsList", t, this)));
              },
              this
            );
        };
      Object(o.g)(p.prototype, s.a), (e.b = p);
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
      function s(t) {
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
      function l(t, e) {
        var n,
          i,
          o,
          r = a[t];
        r || (r = a[t] = { element: s(t), counter: 0 });
        var l = r.counter++;
        return (
          (n = r.element),
          (o = function () {
            d(n, l, "");
          }),
          (i = function (t) {
            d(n, l, t);
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
                for (; r < i.parts.length; r++) a.parts.push(l(t, i.parts[r]));
              } else {
                var s = [];
                for (r = 0; r < i.parts.length; r++) s.push(l(t, i.parts[r]));
                (o[t] = o[t] || {}), (o[t][i.id] = { id: i.id, parts: s });
              }
            }
          })(
            e,
            (function (t) {
              for (var e = [], n = {}, i = 0; i < t.length; i++) {
                var o = t[i],
                  a = o[0],
                  r = o[1],
                  s = o[2],
                  l = { css: r, media: s };
                n[a]
                  ? n[a].parts.push(l)
                  : e.push((n[a] = { id: a, parts: [l] }));
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
          for (var r = Object.keys(n), s = 0; s < r.length; s += 1)
            for (var l = n[r[s]], c = 0; c < l.parts.length; c += 1)
              l.parts[c]();
          delete o[t];
        },
      };
      var c,
        u =
          ((c = []),
          function (t, e) {
            return (c[t] = e), c.filter(Boolean).join("\n");
          });
      function d(t, e, n) {
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
    function (t, e, n) {
      "use strict";
      function i(t, e) {
        var n = t.kind || "cc";
        return t.default || t.defaulttrack
          ? "default"
          : t._id || t.file || n + e;
      }
      function o(t, e) {
        var n = t.label || t.name || t.language;
        return (
          n || ((n = "Unknown CC"), (e += 1) > 1 && (n += " [" + e + "]")),
          { label: n, unknownCount: e }
        );
      }
      n.d(e, "a", function () {
        return i;
      }),
        n.d(e, "b", function () {
          return o;
        });
    },
    function (t, e, n) {
      "use strict";
      function i(t) {
        return new Promise(function (e, n) {
          if (t.paused) return n(o("NotAllowedError", 0, "play() failed."));
          var i = function () {
              t.removeEventListener("play", a),
                t.removeEventListener("playing", r),
                t.removeEventListener("pause", r),
                t.removeEventListener("abort", r),
                t.removeEventListener("error", r);
            },
            a = function () {
              t.addEventListener("playing", r),
                t.addEventListener("abort", r),
                t.addEventListener("error", r),
                t.addEventListener("pause", r);
            },
            r = function (t) {
              if ((i(), "playing" === t.type)) e();
              else {
                var a = 'The play() request was interrupted by a "'.concat(
                  t.type,
                  '" event.'
                );
                "error" === t.type
                  ? n(o("NotSupportedError", 9, a))
                  : n(o("AbortError", 20, a));
              }
            };
          t.addEventListener("play", a);
        });
      }
      function o(t, e, n) {
        var i = new Error(n);
        return (i.name = t), (i.code = e), i;
      }
      n.d(e, "a", function () {
        return i;
      });
    },
    function (t, e, n) {
      "use strict";
      function i(t, e) {
        return t !== 1 / 0 && Math.abs(t) >= Math.max(a(e), 0);
      }
      function o(t, e) {
        var n = "VOD";
        return (
          t === 1 / 0
            ? (n = "LIVE")
            : t < 0 && (n = i(t, a(e)) ? "DVR" : "LIVE"),
          n
        );
      }
      function a(t) {
        return void 0 === t ? 120 : Math.max(t, 0);
      }
      n.d(e, "a", function () {
        return i;
      }),
        n.d(e, "b", function () {
          return o;
        });
    },
    function (t, e, n) {
      "use strict";
      var i = n(67),
        o = n(16),
        a = n(22),
        r = n(4),
        s = n(57),
        l = n(2),
        c = n(1);
      function u(t) {
        throw new c.n(null, t);
      }
      function d(t, e, i) {
        t.xhr = Object(a.a)(
          t.file,
          function (a) {
            !(function (t, e, i, a) {
              var d,
                p,
                h = t.responseXML ? t.responseXML.firstChild : null;
              if (h)
                for (
                  "xml" === Object(r.b)(h) && (h = h.nextSibling);
                  h.nodeType === h.COMMENT_NODE;

                )
                  h = h.nextSibling;
              try {
                if (h && "tt" === Object(r.b)(h))
                  (d = (function (t) {
                    t || u(306007);
                    var e = [],
                      n = t.getElementsByTagName("p"),
                      i = 30,
                      o = t.getElementsByTagName("tt");
                    if (o && o[0]) {
                      var a = parseFloat(o[0].getAttribute("ttp:frameRate"));
                      isNaN(a) || (i = a);
                    }
                    n || u(306005),
                      n.length ||
                        (n = t.getElementsByTagName("tt:p")).length ||
                        (n = t.getElementsByTagName("tts:p"));
                    for (var r = 0; r < n.length; r++) {
                      for (
                        var s = n[r], c = s.getElementsByTagName("br"), d = 0;
                        d < c.length;
                        d++
                      ) {
                        var p = c[d];
                        p.parentNode.replaceChild(t.createTextNode("\r\n"), p);
                      }
                      var w = s.innerHTML || s.textContent || s.text || "",
                        h = Object(l.i)(w)
                          .replace(/>\s+</g, "><")
                          .replace(/(<\/?)tts?:/g, "$1")
                          .replace(/<br.*?\/>/g, "\r\n");
                      if (h) {
                        var f = s.getAttribute("begin"),
                          j = s.getAttribute("dur"),
                          g = s.getAttribute("end"),
                          b = { begin: Object(l.g)(f, i), text: h };
                        g
                          ? (b.end = Object(l.g)(g, i))
                          : j && (b.end = b.begin + Object(l.g)(j, i)),
                          e.push(b);
                      }
                    }
                    return e.length || u(306005), e;
                  })(t.responseXML)),
                    (p = w(d)),
                    delete e.xhr,
                    i(p);
                else {
                  var f = t.responseText;
                  f.indexOf("WEBVTT") >= 0
                    ? n
                        .e(10)
                        .then(
                          function (t) {
                            return n(97).default;
                          }.bind(null, n)
                        )
                        .catch(Object(o.c)(301131))
                        .then(function (t) {
                          var n = new t(window);
                          (p = []),
                            (n.oncue = function (t) {
                              p.push(t);
                            }),
                            (n.onflush = function () {
                              delete e.xhr, i(p);
                            }),
                            n.parse(f);
                        })
                        .catch(function (t) {
                          delete e.xhr, a(Object(c.v)(null, c.b, t));
                        })
                    : ((d = Object(s.a)(f)), (p = w(d)), delete e.xhr, i(p));
                }
              } catch (t) {
                delete e.xhr, a(Object(c.v)(null, c.b, t));
              }
            })(a, t, e, i);
          },
          function (t, e, n, o) {
            i(Object(c.u)(o, c.b));
          }
        );
      }
      function p(t) {
        t &&
          t.forEach(function (t) {
            var e = t.xhr;
            e &&
              ((e.onload = null),
              (e.onreadystatechange = null),
              (e.onerror = null),
              "abort" in e && e.abort()),
              delete t.xhr;
          });
      }
      function w(t) {
        return t.map(function (t) {
          return new i.a(t.begin, t.end, t.text);
        });
      }
      n.d(e, "c", function () {
        return d;
      }),
        n.d(e, "a", function () {
          return p;
        }),
        n.d(e, "b", function () {
          return w;
        });
    },
    function (t, e, n) {
      "use strict";
      var i = window.VTTCue;
      function o(t) {
        if ("string" != typeof t) return !1;
        return (
          !!{ start: !0, middle: !0, end: !0, left: !0, right: !0 }[
            t.toLowerCase()
          ] && t.toLowerCase()
        );
      }
      if (!i) {
        (i = function (t, e, n) {
          var i = this;
          i.hasBeenReset = !1;
          var a = "",
            r = !1,
            s = t,
            l = e,
            c = n,
            u = null,
            d = "",
            p = !0,
            w = "auto",
            h = "start",
            f = "auto",
            j = 100,
            g = "middle";
          Object.defineProperty(i, "id", {
            enumerable: !0,
            get: function () {
              return a;
            },
            set: function (t) {
              a = "" + t;
            },
          }),
            Object.defineProperty(i, "pauseOnExit", {
              enumerable: !0,
              get: function () {
                return r;
              },
              set: function (t) {
                r = !!t;
              },
            }),
            Object.defineProperty(i, "startTime", {
              enumerable: !0,
              get: function () {
                return s;
              },
              set: function (t) {
                if ("number" != typeof t)
                  throw new TypeError("Start time must be set to a number.");
                (s = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "endTime", {
              enumerable: !0,
              get: function () {
                return l;
              },
              set: function (t) {
                if ("number" != typeof t)
                  throw new TypeError("End time must be set to a number.");
                (l = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "text", {
              enumerable: !0,
              get: function () {
                return c;
              },
              set: function (t) {
                (c = "" + t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "region", {
              enumerable: !0,
              get: function () {
                return u;
              },
              set: function (t) {
                (u = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "vertical", {
              enumerable: !0,
              get: function () {
                return d;
              },
              set: function (t) {
                var e = (function (t) {
                  return (
                    "string" == typeof t &&
                    !!{ "": !0, lr: !0, rl: !0 }[t.toLowerCase()] &&
                    t.toLowerCase()
                  );
                })(t);
                if (!1 === e)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (d = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "snapToLines", {
              enumerable: !0,
              get: function () {
                return p;
              },
              set: function (t) {
                (p = !!t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "line", {
              enumerable: !0,
              get: function () {
                return w;
              },
              set: function (t) {
                if ("number" != typeof t && "auto" !== t)
                  throw new SyntaxError(
                    "An invalid number or illegal string was specified."
                  );
                (w = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "lineAlign", {
              enumerable: !0,
              get: function () {
                return h;
              },
              set: function (t) {
                var e = o(t);
                if (!e)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (h = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "position", {
              enumerable: !0,
              get: function () {
                return f;
              },
              set: function (t) {
                if (t < 0 || t > 100)
                  throw new Error("Position must be between 0 and 100.");
                (f = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "size", {
              enumerable: !0,
              get: function () {
                return j;
              },
              set: function (t) {
                if (t < 0 || t > 100)
                  throw new Error("Size must be between 0 and 100.");
                (j = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(i, "align", {
              enumerable: !0,
              get: function () {
                return g;
              },
              set: function (t) {
                var e = o(t);
                if (!e)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (g = e), (this.hasBeenReset = !0);
              },
            }),
            (i.displayState = void 0);
        }).prototype.getCueAsHTML = function () {
          return window.WebVTT.convertCueToDOMTree(window, this.text);
        };
      }
      e.a = i;
    },
    ,
    function (t, e, n) {
      var i = n(70);
      "string" == typeof i && (i = [["all-players", i, ""]]),
        n(61).style(i, "all-players"),
        i.locals && (t.exports = i.locals);
    },
    function (t, e, n) {
      (t.exports = n(60)(!1)).push([
        t.i,
        '.jw-reset{text-align:left;direction:ltr}.jw-reset-text,.jw-reset{color:inherit;background-color:transparent;padding:0;margin:0;float:none;font-family:Arial,Helvetica,sans-serif;font-size:1em;line-height:1em;list-style:none;text-transform:none;vertical-align:baseline;border:0;font-variant:inherit;font-stretch:inherit;-webkit-tap-highlight-color:rgba(255,255,255,0)}body .jw-error,body .jwplayer.jw-state-error{height:100%;width:100%}.jw-title{position:absolute;top:0}.jw-background-color{background:rgba(0,0,0,0.4)}.jw-text{color:rgba(255,255,255,0.8)}.jw-knob{color:rgba(255,255,255,0.8);background-color:#fff}.jw-button-color{color:rgba(255,255,255,0.8)}:not(.jw-flag-touch) .jw-button-color:not(.jw-logo-button):focus,:not(.jw-flag-touch) .jw-button-color:not(.jw-logo-button):hover{color:#fff}.jw-toggle{color:#fff}.jw-toggle.jw-off{color:rgba(255,255,255,0.8)}.jw-toggle.jw-off:focus{color:#fff}.jw-toggle:focus{outline:none}:not(.jw-flag-touch) .jw-toggle.jw-off:hover{color:#fff}.jw-rail{background:rgba(255,255,255,0.3)}.jw-buffer{background:rgba(255,255,255,0.3)}.jw-progress{background:#f2f2f2}.jw-time-tip,.jw-volume-tip{border:0}.jw-slider-volume.jw-volume-tip.jw-background-color.jw-slider-vertical{background:none}.jw-skip{padding:.5em;outline:none}.jw-skip .jw-skiptext,.jw-skip .jw-skip-icon{color:rgba(255,255,255,0.8)}.jw-skip.jw-skippable:hover .jw-skip-icon,.jw-skip.jw-skippable:focus .jw-skip-icon{color:#fff}.jw-icon-cast google-cast-launcher{--connected-color:#fff;--disconnected-color:rgba(255,255,255,0.8)}.jw-icon-cast google-cast-launcher:focus{outline:none}.jw-icon-cast google-cast-launcher.jw-off{--connected-color:rgba(255,255,255,0.8)}.jw-icon-cast:focus google-cast-launcher{--connected-color:#fff;--disconnected-color:#fff}.jw-icon-cast:hover google-cast-launcher{--connected-color:#fff;--disconnected-color:#fff}.jw-nextup-container{bottom:2.5em;padding:5px .5em}.jw-nextup{border-radius:0}.jw-color-active{color:#fff;stroke:#fff;border-color:#fff}:not(.jw-flag-touch) .jw-color-active-hover:hover,:not(.jw-flag-touch) .jw-color-active-hover:focus{color:#fff;stroke:#fff;border-color:#fff}.jw-color-inactive{color:rgba(255,255,255,0.8);stroke:rgba(255,255,255,0.8);border-color:rgba(255,255,255,0.8)}:not(.jw-flag-touch) .jw-color-inactive-hover:hover{color:rgba(255,255,255,0.8);stroke:rgba(255,255,255,0.8);border-color:rgba(255,255,255,0.8)}.jw-option{color:rgba(255,255,255,0.8)}.jw-option.jw-active-option{color:#fff;background-color:rgba(255,255,255,0.1)}:not(.jw-flag-touch) .jw-option:hover{color:#fff}.jwplayer{width:100%;font-size:16px;position:relative;display:block;min-height:0;overflow:hidden;box-sizing:border-box;font-family:Arial,Helvetica,sans-serif;-webkit-touch-callout:none;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none;outline:none}.jwplayer *{box-sizing:inherit}.jwplayer.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jwplayer.jw-flag-aspect-mode{height:auto !important}.jwplayer.jw-flag-aspect-mode .jw-aspect{display:block}.jwplayer .jw-aspect{display:none}.jwplayer .jw-swf{outline:none}.jw-media,.jw-preview{position:absolute;width:100%;height:100%;top:0;left:0;bottom:0;right:0}.jw-media{overflow:hidden;cursor:pointer}.jw-plugin{position:absolute;bottom:66px}.jw-breakpoint-7 .jw-plugin{bottom:132px}.jw-plugin .jw-banner{max-width:100%;opacity:0;cursor:pointer;position:absolute;margin:auto auto 0;left:0;right:0;bottom:0;display:block}.jw-preview,.jw-captions,.jw-title{pointer-events:none}.jw-media,.jw-logo{pointer-events:all}.jw-wrapper{background-color:#000;position:absolute;top:0;left:0;right:0;bottom:0}.jw-hidden-accessibility{border:0;clip:rect(0 0 0 0);height:1px;margin:-1px;overflow:hidden;padding:0;position:absolute;width:1px}.jw-contract-trigger::before{content:"";overflow:hidden;width:200%;height:200%;display:block;position:absolute;top:0;left:0}.jwplayer .jw-media video{position:absolute;top:0;right:0;bottom:0;left:0;width:100%;height:100%;margin:auto;background:transparent}.jwplayer .jw-media video::-webkit-media-controls-start-playback-button{display:none}.jwplayer.jw-stretch-uniform .jw-media video{object-fit:contain}.jwplayer.jw-stretch-none .jw-media video{object-fit:none}.jwplayer.jw-stretch-fill .jw-media video{object-fit:cover}.jwplayer.jw-stretch-exactfit .jw-media video{object-fit:fill}.jw-preview{position:absolute;display:none;opacity:1;visibility:visible;width:100%;height:100%;background:#000 no-repeat 50% 50%}.jwplayer .jw-preview,.jw-error .jw-preview{background-size:contain}.jw-stretch-none .jw-preview{background-size:auto auto}.jw-stretch-fill .jw-preview{background-size:cover}.jw-stretch-exactfit .jw-preview{background-size:100% 100%}.jw-title{display:none;padding-top:20px;width:100%;z-index:1}.jw-title-primary,.jw-title-secondary{color:#fff;padding-left:20px;padding-right:20px;padding-bottom:.5em;overflow:hidden;text-overflow:ellipsis;direction:unset;white-space:nowrap;width:100%}.jw-title-primary{font-size:1.625em}.jw-breakpoint-2 .jw-title-primary,.jw-breakpoint-3 .jw-title-primary{font-size:1.5em}.jw-flag-small-player .jw-title-primary{font-size:1.25em}.jw-flag-small-player .jw-title-secondary,.jw-title-secondary:empty{display:none}.jw-captions{position:absolute;width:100%;height:100%;text-align:center;display:none;letter-spacing:normal;word-spacing:normal;text-transform:none;text-indent:0;text-decoration:none;pointer-events:none;overflow:hidden;top:0}.jw-captions.jw-captions-enabled{display:block}.jw-captions-window{display:none;padding:.25em;border-radius:.25em}.jw-captions-window.jw-captions-window-active{display:inline-block}.jw-captions-text{display:inline-block;color:#fff;background-color:#000;word-wrap:normal;word-break:normal;white-space:pre-line;font-style:normal;font-weight:normal;text-align:center;text-decoration:none}.jw-text-track-display{font-size:inherit;line-height:1.5}.jw-text-track-cue{background-color:rgba(0,0,0,0.5);color:#fff;padding:.1em .3em}.jwplayer video::-webkit-media-controls{display:none;justify-content:flex-start}.jwplayer video::-webkit-media-text-track-display{min-width:-webkit-min-content}.jwplayer video::cue{background-color:rgba(0,0,0,0.5)}.jwplayer video::-webkit-media-controls-panel-container{display:none}.jwplayer:not(.jw-flag-controls-hidden):not(.jw-state-playing) .jw-captions,.jwplayer.jw-flag-media-audio.jw-state-playing .jw-captions,.jwplayer.jw-state-playing:not(.jw-flag-user-inactive):not(.jw-flag-controls-hidden) .jw-captions{max-height:calc(100% - 60px)}.jwplayer:not(.jw-flag-controls-hidden):not(.jw-state-playing):not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container,.jwplayer.jw-flag-media-audio.jw-state-playing:not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container,.jwplayer.jw-state-playing:not(.jw-flag-user-inactive):not(.jw-flag-controls-hidden):not(.jw-flag-ios-fullscreen) video::-webkit-media-text-track-container{max-height:calc(100% - 60px)}.jw-logo{position:absolute;margin:20px;cursor:pointer;pointer-events:all;background-repeat:no-repeat;background-size:contain;top:auto;right:auto;left:auto;bottom:auto;outline:none}.jw-logo.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jw-flag-audio-player .jw-logo{display:none}.jw-logo-top-right{top:0;right:0}.jw-logo-top-left{top:0;left:0}.jw-logo-bottom-left{left:0}.jw-logo-bottom-right{right:0}.jw-logo-bottom-left,.jw-logo-bottom-right{bottom:44px;transition:bottom 150ms cubic-bezier(0, .25, .25, 1)}.jw-state-idle .jw-logo{z-index:1}.jw-state-setup .jw-wrapper{background-color:inherit}.jw-state-setup .jw-logo,.jw-state-setup .jw-controls,.jw-state-setup .jw-controls-backdrop{visibility:hidden}span.jw-break{display:block}body .jw-error,body .jwplayer.jw-state-error{background-color:#333;color:#fff;font-size:16px;display:table;opacity:1;position:relative}body .jw-error .jw-display,body .jwplayer.jw-state-error .jw-display{display:none}body .jw-error .jw-media,body .jwplayer.jw-state-error .jw-media{cursor:default}body .jw-error .jw-preview,body .jwplayer.jw-state-error .jw-preview{background-color:#333}body .jw-error .jw-error-msg,body .jwplayer.jw-state-error .jw-error-msg{background-color:#000;border-radius:2px;display:flex;flex-direction:row;align-items:stretch;padding:20px}body .jw-error .jw-error-msg .jw-icon,body .jwplayer.jw-state-error .jw-error-msg .jw-icon{height:30px;width:30px;margin-right:20px;flex:0 0 auto;align-self:center}body .jw-error .jw-error-msg .jw-icon:empty,body .jwplayer.jw-state-error .jw-error-msg .jw-icon:empty{display:none}body .jw-error .jw-error-msg .jw-info-container,body .jwplayer.jw-state-error .jw-error-msg .jw-info-container{margin:0;padding:0}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg{flex-direction:column}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-error-text,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-error-text,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-error-text,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-error-text{text-align:center}body .jw-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-icon,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-flag-small-player .jw-error-msg .jw-icon,body .jw-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-icon,body .jwplayer.jw-state-error:not(.jw-flag-audio-player).jw-breakpoint-2 .jw-error-msg .jw-icon{flex:.5 0 auto;margin-right:0;margin-bottom:20px}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg .jw-break,.jwplayer.jw-state-error.jw-flag-small-player .jw-error-msg .jw-break,.jwplayer.jw-state-error.jw-breakpoint-2 .jw-error-msg .jw-break{display:inline}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg .jw-break:before,.jwplayer.jw-state-error.jw-flag-small-player .jw-error-msg .jw-break:before,.jwplayer.jw-state-error.jw-breakpoint-2 .jw-error-msg .jw-break:before{content:" "}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg{height:100%;width:100%;top:0;position:absolute;left:0;background:#000;-webkit-transform:none;transform:none;padding:4px 16px;z-index:1}.jwplayer.jw-state-error.jw-flag-audio-player .jw-error-msg.jw-info-overlay{max-width:none;max-height:none}body .jwplayer.jw-state-error .jw-title,.jw-state-idle .jw-title,.jwplayer.jw-state-complete:not(.jw-flag-casting):not(.jw-flag-audio-player):not(.jw-flag-overlay-open-related) .jw-title{display:block}body .jwplayer.jw-state-error .jw-preview,.jw-state-idle .jw-preview,.jwplayer.jw-state-complete:not(.jw-flag-casting):not(.jw-flag-audio-player):not(.jw-flag-overlay-open-related) .jw-preview{display:block}.jw-state-idle .jw-captions,.jwplayer.jw-state-complete .jw-captions,body .jwplayer.jw-state-error .jw-captions{display:none}.jw-state-idle video::-webkit-media-text-track-container,.jwplayer.jw-state-complete video::-webkit-media-text-track-container,body .jwplayer.jw-state-error video::-webkit-media-text-track-container{display:none}.jwplayer.jw-flag-fullscreen{width:100% !important;height:100% !important;top:0;right:0;bottom:0;left:0;z-index:1000;margin:0;position:fixed}body .jwplayer.jw-flag-flash-blocked .jw-title{display:block}.jwplayer.jw-flag-controls-hidden .jw-media{cursor:default}.jw-flag-audio-player:not(.jw-flag-flash-blocked) .jw-media{visibility:hidden}.jw-flag-audio-player .jw-title{background:none}.jw-flag-audio-player object{min-height:45px}.jw-flag-floating{background-size:cover;background-color:#000}.jw-flag-floating .jw-wrapper{position:fixed;z-index:2147483647;-webkit-animation:jw-float-to-bottom 150ms cubic-bezier(0, .25, .25, 1) forwards 1;animation:jw-float-to-bottom 150ms cubic-bezier(0, .25, .25, 1) forwards 1;top:auto;bottom:1rem;left:auto;right:1rem;max-width:400px;max-height:400px;margin:0 auto}@media screen and (max-width:480px){.jw-flag-floating .jw-wrapper{width:100%;left:0;right:0}}.jw-flag-floating .jw-wrapper .jw-media{touch-action:none}@media screen and (max-device-width:480px) and (orientation:portrait){.jw-flag-touch.jw-flag-floating .jw-wrapper{-webkit-animation:none;animation:none;top:62px;bottom:auto;left:0;right:0;max-width:none;max-height:none}}.jw-flag-floating .jw-float-icon{pointer-events:all;cursor:pointer;display:none}.jw-flag-floating .jw-float-icon .jw-svg-icon{-webkit-filter:drop-shadow(0 0 1px #000);filter:drop-shadow(0 0 1px #000)}.jw-flag-floating.jw-floating-dismissible .jw-dismiss-icon{display:none}.jw-flag-floating.jw-floating-dismissible.jw-flag-ads .jw-float-icon{display:flex}.jw-flag-floating.jw-floating-dismissible.jw-state-paused .jw-logo,.jw-flag-floating.jw-floating-dismissible:not(.jw-flag-user-inactive) .jw-logo{display:none}.jw-flag-floating.jw-floating-dismissible.jw-state-paused .jw-float-icon,.jw-flag-floating.jw-floating-dismissible:not(.jw-flag-user-inactive) .jw-float-icon{display:flex}.jw-float-icon{display:none;position:absolute;top:3px;right:5px;align-items:center;justify-content:center}@-webkit-keyframes jw-float-to-bottom{from{-webkit-transform:translateY(100%);transform:translateY(100%)}to{-webkit-transform:translateY(0);transform:translateY(0)}}@keyframes jw-float-to-bottom{from{-webkit-transform:translateY(100%);transform:translateY(100%)}to{-webkit-transform:translateY(0);transform:translateY(0)}}.jw-flag-top{margin-top:2em;overflow:visible}.jw-top{height:2em;line-height:2;pointer-events:none;text-align:center;opacity:.8;position:absolute;top:-2em;width:100%}.jw-top .jw-icon{cursor:pointer;pointer-events:all;height:auto;width:auto}.jw-top .jw-text{color:#555}',
        "",
      ]);
    },
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
