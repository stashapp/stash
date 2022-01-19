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
  [6, 1, 2, 3, 4, 5, 7, 9],
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
    function (t, e, i) {
      "use strict";
      i.r(e);
      var n,
        o = i(8),
        a = i(3),
        r = i(7),
        s = i(43),
        l = i(5),
        c = i(15),
        u = i(40);
      function d(t) {
        return (
          n || (n = new DOMParser()),
          Object(l.r)(
            Object(l.s)(n.parseFromString(t, "image/svg+xml").documentElement)
          )
        );
      }
      var p = function (t, e, i, n) {
          var o = document.createElement("div");
          (o.className =
            "jw-icon jw-icon-inline jw-button-color jw-reset " + t),
            o.setAttribute("role", "button"),
            o.setAttribute("tabindex", "0"),
            i && o.setAttribute("aria-label", i),
            (o.style.display = "none");
          var a = new u.a(o).on("click tap enter", e || function () {});
          return (
            n &&
              Array.prototype.forEach.call(n, function (t) {
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
        h = i(0),
        f = i(71),
        w = i.n(f),
        g = i(72),
        j = i.n(g),
        b = i(73),
        m = i.n(b),
        v = i(74),
        y = i.n(v),
        k = i(75),
        x = i.n(k),
        T = i(76),
        O = i.n(T),
        C = i(77),
        _ = i.n(C),
        M = i(78),
        S = i.n(M),
        E = i(79),
        I = i.n(E),
        L = i(80),
        A = i.n(L),
        P = i(81),
        R = i.n(P),
        z = i(82),
        B = i.n(z),
        V = i(83),
        N = i.n(V),
        H = i(84),
        F = i.n(H),
        D = i(85),
        q = i.n(D),
        U = i(86),
        W = i.n(U),
        Q = i(62),
        Y = i.n(Q),
        X = i(87),
        K = i.n(X),
        J = i(88),
        Z = i.n(J),
        G = i(89),
        $ = i.n(G),
        tt = i(90),
        et = i.n(tt),
        it = i(91),
        nt = i.n(it),
        ot = i(92),
        at = i.n(ot),
        rt = i(93),
        st = i.n(rt),
        lt = i(94),
        ct = i.n(lt),
        ut = null;
      function dt(t) {
        var e = wt().querySelector(ht(t));
        if (e) return ft(e);
        throw new Error("Icon not found " + t);
      }
      function pt(t) {
        var e = wt().querySelectorAll(t.split(",").map(ht).join(","));
        if (!e.length) throw new Error("Icons not found " + t);
        return Array.prototype.map.call(e, function (t) {
          return ft(t);
        });
      }
      function ht(t) {
        return ".jw-svg-icon-".concat(t);
      }
      function ft(t) {
        return t.cloneNode(!0);
      }
      function wt() {
        return (
          ut ||
            (ut = d(
              "<xml>" +
                w.a +
                j.a +
                m.a +
                y.a +
                x.a +
                O.a +
                _.a +
                S.a +
                I.a +
                A.a +
                R.a +
                B.a +
                N.a +
                F.a +
                q.a +
                W.a +
                Y.a +
                K.a +
                Z.a +
                $.a +
                et.a +
                nt.a +
                at.a +
                st.a +
                ct.a +
                "</xml>"
            )),
          ut
        );
      }
      var gt = i(10);
      function jt(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var bt = {};
      var mt = (function () {
          function t(e, i, n, o, a) {
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
              i && s.setAttribute("aria-label", i),
              e && "<svg" === e.substring(0, 4)
                ? (r = (function (t) {
                    if (!bt[t]) {
                      var e = Object.keys(bt);
                      e.length > 10 && delete bt[e[0]];
                      var i = d(t);
                      bt[t] = i;
                    }
                    return bt[t].cloneNode(!0);
                  })(e))
                : (((r = document.createElement("div")).className =
                    "jw-icon jw-button-image jw-button-color jw-reset"),
                  e &&
                    Object(gt.d)(r, {
                      backgroundImage: "url(".concat(e, ")"),
                    })),
              s.appendChild(r),
              new u.a(s).on("click tap enter", n, this),
              s.addEventListener("mousedown", function (t) {
                t.preventDefault();
              }),
              (this.id = o),
              (this.buttonElement = s);
          }
          var e, i, n;
          return (
            (e = t),
            (i = [
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
            ]) && jt(e.prototype, i),
            n && jt(e, n),
            t
          );
        })(),
        vt = i(11);
      function yt(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var kt = function (t) {
          var e = Object(l.c)(t),
            i = window.pageXOffset;
          return (
            i &&
              o.OS.android &&
              document.body.parentElement.getBoundingClientRect().left >= 0 &&
              ((e.left -= i), (e.right -= i)),
            e
          );
        },
        xt = (function () {
          function t(e, i) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(h.g)(this, r.a),
              (this.className = e + " jw-background-color jw-reset"),
              (this.orientation = i);
          }
          var e, i, n;
          return (
            (e = t),
            (i = [
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
                    i,
                    n = (this.railBounds = this.railBounds
                      ? this.railBounds
                      : kt(this.elementRail));
                  return (
                    (i =
                      "horizontal" === this.orientation
                        ? (e = t.pageX) < n.left
                          ? 0
                          : e > n.right
                          ? 100
                          : 100 * Object(s.a)((e - n.left) / n.width, 0, 1)
                        : (e = t.pageY) >= n.bottom
                        ? 0
                        : e <= n.top
                        ? 100
                        : 100 *
                          Object(s.a)(
                            (n.height - (e - n.top)) / n.height,
                            0,
                            1
                          )),
                    this.render(i),
                    this.update(i),
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
            ]) && yt(e.prototype, i),
            n && yt(e, n),
            t
          );
        })(),
        Tt = function (t, e) {
          t &&
            e &&
            (t.setAttribute("aria-label", e),
            t.setAttribute("role", "button"),
            t.setAttribute("tabindex", "0"));
        };
      function Ot(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var Ct = (function () {
          function t(e, i, n, o) {
            var a = this;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(h.g)(this, r.a),
              (this.el = document.createElement("div"));
            var s =
              "jw-icon jw-icon-tooltip " + e + " jw-button-color jw-reset";
            n || (s += " jw-hidden"),
              Tt(this.el, i),
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
          var e, i, n;
          return (
            (e = t),
            (i = [
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
            ]) && Ot(e.prototype, i),
            n && Ot(e, n),
            t
          );
        })(),
        _t = i(22),
        Mt = i(57);
      function St(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var Et = (function () {
          function t(e, i) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              (this.time = e),
              (this.text = i),
              (this.el = document.createElement("div")),
              (this.el.className = "jw-cue jw-reset");
          }
          var e, i, n;
          return (
            (e = t),
            (i = [
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
            ]) && St(e.prototype, i),
            n && St(e, n),
            t
          );
        })(),
        It = {
          loadChapters: function (t) {
            Object(_t.a)(
              t,
              this.chaptersLoaded.bind(this),
              this.chaptersFailed,
              { plainText: !0 }
            );
          },
          chaptersLoaded: function (t) {
            var e = Object(Mt.a)(t.responseText);
            if (Array.isArray(e)) {
              var i = this._model.get("cues").concat(e);
              this._model.set("cues", i);
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
              this.cues.forEach(function (i) {
                i.align(e),
                  i.el.addEventListener("mouseover", function () {
                    t.activeCue = i;
                  }),
                  i.el.addEventListener("mouseout", function () {
                    t.activeCue = null;
                  }),
                  t.elementRail.appendChild(i.el);
              });
          },
          resetCues: function () {
            this.cues.forEach(function (t) {
              t.el.parentNode && t.el.parentNode.removeChild(t.el);
            }),
              (this.cues = []);
          },
        };
      function Lt(t) {
        (this.begin = t.begin), (this.end = t.end), (this.img = t.text);
      }
      var At = {
        loadThumbnails: function (t) {
          t &&
            ((this.vttPath = t.split("?")[0].split("/").slice(0, -1).join("/")),
            (this.individualImage = null),
            Object(_t.a)(
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
              this.thumbnails.push(new Lt(t));
            }, this),
            this.drawCues());
        },
        thumbnailsFailed: function () {},
        chooseThumbnail: function (t) {
          var e = Object(h.A)(this.thumbnails, { end: t }, Object(h.z)("end"));
          e >= this.thumbnails.length && (e = this.thumbnails.length - 1);
          var i = this.thumbnails[e].img;
          return (
            i.indexOf("://") < 0 &&
              (i = this.vttPath ? this.vttPath + "/" + i : i),
            i
          );
        },
        loadThumbnail: function (t) {
          var e = this.chooseThumbnail(t),
            i = { margin: "0 auto", backgroundPosition: "0 0" };
          if (e.indexOf("#xywh") > 0)
            try {
              var n = /(.+)#xywh=(\d+),(\d+),(\d+),(\d+)/.exec(e);
              (e = n[1]),
                (i.backgroundPosition = -1 * n[2] + "px " + -1 * n[3] + "px"),
                (i.width = n[4]),
                this.timeTip.setWidth(+i.width),
                (i.height = n[5]);
            } catch (t) {
              return;
            }
          else
            this.individualImage ||
              ((this.individualImage = new Image()),
              (this.individualImage.onload = Object(h.a)(function () {
                (this.individualImage.onload = null),
                  this.timeTip.image({
                    width: this.individualImage.width,
                    height: this.individualImage.height,
                  }),
                  this.timeTip.setWidth(this.individualImage.width);
              }, this)),
              (this.individualImage.src = e));
          return (i.backgroundImage = 'url("' + e + '")'), i;
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
      function Pt(t, e, i) {
        return (Pt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, i) {
                var n = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = Ht(t));

                  );
                  return t;
                })(t, e);
                if (n) {
                  var o = Object.getOwnPropertyDescriptor(n, e);
                  return o.get ? o.get.call(i) : o.value;
                }
              })(t, e, i || t);
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
      function zt(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Bt(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function Vt(t, e, i) {
        return e && Bt(t.prototype, e), i && Bt(t, i), t;
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
          e && Dt(t, e);
      }
      function Dt(t, e) {
        return (Dt =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var qt = (function (t) {
        function e() {
          return zt(this, e), Nt(this, Ht(e).apply(this, arguments));
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
                Object(gt.d)(this.img, t);
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
      })(Ct);
      var Ut = (function (t) {
        function e(t, i, n) {
          var o;
          return (
            zt(this, e),
            ((o = Nt(
              this,
              Ht(e).call(this, "jw-slider-time", "horizontal")
            ))._model = t),
            (o._api = i),
            (o.timeUpdateKeeper = n),
            (o.timeTip = new qt("jw-tooltip-time", null, !0)),
            o.timeTip.setup(),
            (o.cues = []),
            (o.seekThrottled = Object(h.B)(o.performSeek, 400)),
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
                Pt(Ht(e.prototype), "setup", this).apply(this, arguments),
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
                var i = this.el;
                Object(l.t)(i, "tabindex", "0"),
                  Object(l.t)(i, "role", "slider"),
                  Object(l.t)(
                    i,
                    "aria-label",
                    this._model.get("localization").slider
                  ),
                  i.removeAttribute("aria-hidden"),
                  this.elementRail.appendChild(this.timeTip.element()),
                  (this.ui = (this.ui || new u.a(i))
                    .on("move drag", this.showTimeTooltip, this)
                    .on("dragEnd out", this.hideTimeTooltip, this)
                    .on("click", function () {
                      return i.focus();
                    })
                    .on("focus", this.updateAriaText, this));
              },
            },
            {
              key: "update",
              value: function (t) {
                (this.seekTo = t),
                  this.seekThrottled(),
                  Pt(Ht(e.prototype), "update", this).apply(this, arguments);
              },
            },
            {
              key: "dragStart",
              value: function () {
                this._model.set("scrubbing", !0),
                  Pt(Ht(e.prototype), "dragStart", this).apply(this, arguments);
              },
            },
            {
              key: "dragEnd",
              value: function () {
                Pt(Ht(e.prototype), "dragEnd", this).apply(this, arguments),
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
                var i = 0;
                if (e)
                  if ("DVR" === this.streamType) {
                    var n = this._model.get("dvrSeekLimit"),
                      o = e + n;
                    i = ((o - (t + n)) / o) * 100;
                  } else
                    ("VOD" !== this.streamType && this.streamType) ||
                      (i = (t / e) * 100);
                this.render(i);
              },
            },
            {
              key: "onPlaylistItem",
              value: function (t, e) {
                this.reset();
                var i = t.get("cues");
                !this.cues.length && i.length && this.updateCues(null, i);
                var n = e.tracks;
                Object(h.f)(
                  n,
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
                  i = this._model.get("duration");
                if (0 === i) this._api.play({ reason: "interaction" });
                else if ("DVR" === this.streamType) {
                  var n = this._model.get("seekRange") || { start: 0 },
                    o = this._model.get("dvrSeekLimit");
                  (t = n.start + ((-i - o) * e) / 100),
                    this._api.seek(t, { reason: "interaction" });
                } else
                  (t = (e / 100) * i),
                    this._api.seek(Math.min(t, i - 0.25), {
                      reason: "interaction",
                    });
              },
            },
            {
              key: "showTimeTooltip",
              value: function (t) {
                var e = this,
                  i = this._model.get("duration");
                if (0 !== i) {
                  var n,
                    o = this._model.get("containerWidth"),
                    a = Object(l.c)(this.elementRail),
                    r = t.pageX ? t.pageX - a.left : t.x,
                    c = (r = Object(s.a)(r, 0, a.width)) / a.width,
                    u = i * c;
                  if (i < 0)
                    u = (i += this._model.get("dvrSeekLimit")) - (u = i * c);
                  if (
                    ("touch" === t.pointerType &&
                      (this.activeCue = this.cues.reduce(function (t, i) {
                        return Math.abs(r - (parseInt(i.pct) / 100) * a.width) <
                          e.mobileHoverDistance
                          ? i
                          : t;
                      }, void 0)),
                    this.activeCue)
                  )
                    n = this.activeCue.text;
                  else {
                    (n = Object(vt.timeFormat)(u, !0)),
                      i < 0 && u > -1 && (n = "Live");
                  }
                  var d = this.timeTip;
                  d.update(n),
                    this.textLength !== n.length &&
                      ((this.textLength = n.length), d.resetWidth()),
                    this.showThumbnail(u),
                    Object(l.a)(d.el, "jw-open");
                  var p = d.getWidth(),
                    h = a.width / 100,
                    f = o - a.width,
                    w = 0;
                  p > f && (w = (p - f) / (200 * h));
                  var g = 100 * Math.min(1 - w, Math.max(w, c)).toFixed(3);
                  Object(gt.d)(d.el, { left: g + "%" });
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
                var i = this;
                this.resetCues(),
                  e &&
                    e.length &&
                    (e.forEach(function (t) {
                      i.addCue(t);
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
                    i = t.get("duration"),
                    n = Object(vt.timeFormat)(e);
                  "DVR" !== this.streamType &&
                    (n += " of ".concat(Object(vt.timeFormat)(i)));
                  var o = this.el;
                  document.activeElement !== o &&
                    (this.timeUpdateKeeper.textContent = n),
                    Object(l.t)(o, "aria-valuenow", e),
                    Object(l.t)(o, "aria-valuetext", n);
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
      Object(h.g)(Ut.prototype, It, At);
      var Wt = Ut;
      function Qt(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function Yt(t, e, i) {
        return (Yt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, i) {
                var n = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = Gt(t));

                  );
                  return t;
                })(t, e);
                if (n) {
                  var o = Object.getOwnPropertyDescriptor(n, e);
                  return o.get ? o.get.call(i) : o.value;
                }
              })(t, e, i || t);
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
      function Kt(t, e) {
        if (!(t instanceof e))
          throw new TypeError("Cannot call a class as a function");
      }
      function Jt(t, e) {
        return !e || ("object" !== Xt(e) && "function" != typeof e) ? Zt(t) : e;
      }
      function Zt(t) {
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
          function e(t, i, n) {
            var o;
            Kt(this, e);
            var a = "jw-slider-volume";
            return (
              "vertical" === t && (a += " jw-volume-tip"),
              (o = Jt(this, Gt(e).call(this, a, t))).setup(),
              o.element().classList.remove("jw-background-color"),
              Object(l.t)(n, "tabindex", "0"),
              Object(l.t)(n, "aria-label", i),
              Object(l.t)(n, "aria-orientation", t),
              Object(l.t)(n, "aria-valuemin", 0),
              Object(l.t)(n, "aria-valuemax", 100),
              Object(l.t)(n, "role", "slider"),
              (o.uiOver = new u.a(n).on("click", function () {})),
              o
            );
          }
          return $t(e, t), e;
        })(xt),
        ie = (function (t) {
          function e(t, i, n, o, a) {
            var r;
            Kt(this, e),
              ((r = Jt(this, Gt(e).call(this, i, n, !0, o)))._model = t),
              (r.horizontalContainer = a);
            var s = t.get("localization").volumeSlider;
            return (
              (r.horizontalSlider = new ee("horizontal", s, a, Zt(Zt(r)))),
              (r.verticalSlider = new ee("vertical", s, r.tooltip, Zt(Zt(r)))),
              a.appendChild(r.horizontalSlider.element()),
              r.addContent(r.verticalSlider.element()),
              r.verticalSlider.on(
                "update",
                function (t) {
                  this.trigger("update", t);
                },
                Zt(Zt(r))
              ),
              r.horizontalSlider.on(
                "update",
                function (t) {
                  this.trigger("update", t);
                },
                Zt(Zt(r))
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
                .on("click enter", r.toggleValue, Zt(Zt(r)))
                .on("tap", r.toggleOpenState, Zt(Zt(r)))),
              r.addSliderHandlers(r.ui),
              r.addSliderHandlers(r.horizontalSlider.uiOver),
              r.addSliderHandlers(r.verticalSlider.uiOver),
              r.onAudioMode(null, t.get("audioMode")),
              r._model.on("change:audioMode", r.onAudioMode, Zt(Zt(r))),
              r._model.on("change:volume", r.onVolume, Zt(Zt(r))),
              r
            );
          }
          var i, n, o;
          return (
            $t(e, t),
            (i = e),
            (n = [
              {
                key: "onAudioMode",
                value: function (t, e) {
                  var i = e ? 0 : -1;
                  Object(l.t)(this.horizontalContainer, "tabindex", i);
                },
              },
              {
                key: "addSliderHandlers",
                value: function (t) {
                  var e = this.openSlider,
                    i = this.closeSlider;
                  t.on("over", e, this)
                    .on("out", i, this)
                    .on("focus", e, this)
                    .on("blur", i, this);
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
            ]) && Qt(i.prototype, n),
            o && Qt(i, o),
            e
          );
        })(Ct);
      function ne(t, e, i, n, o) {
        var a = document.createElement("div");
        (a.className = "jw-reset-text jw-tooltip jw-tooltip-".concat(e)),
          a.setAttribute("dir", "auto");
        var r = document.createElement("div");
        (r.className = "jw-text"), a.appendChild(r), t.appendChild(a);
        var s = {
            dirty: !!i,
            opened: !1,
            text: i,
            open: function () {
              s.touchEvent ||
                (s.suppress ? (s.suppress = !1) : (c(!0), n && n()));
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
      var oe = i(47);
      function ae(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function re(t, e) {
        var i = document.createElement("div");
        return (
          (i.className = "jw-icon jw-icon-inline jw-text jw-reset " + t),
          e && Object(l.t)(i, "role", e),
          i
        );
      }
      function se(t) {
        var e = document.createElement("div");
        return (e.className = "jw-reset ".concat(t)), e;
      }
      function le(t, e) {
        if (o.Browser.safari) {
          var i = p(
            "jw-icon-airplay jw-off",
            t,
            e.airplay,
            pt("airplay-off,airplay-on")
          );
          return ne(i.element(), "airplay", e.airplay), i;
        }
        if (o.Browser.chrome && window.chrome) {
          var n = document.createElement("google-cast-launcher");
          Object(l.t)(n, "tabindex", "-1"), (n.className += " jw-reset");
          var a = p("jw-icon-cast", null, e.cast);
          a.ui.off();
          var r = a.element();
          return (
            (r.style.cursor = "pointer"),
            r.appendChild(n),
            (a.button = n),
            ne(r, "chromecast", e.cast),
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
          function t(e, i, n) {
            var s = this;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(h.g)(this, r.a),
              (this._api = e),
              (this._model = i),
              (this._isMobile = o.OS.mobile),
              (this._volumeAnnouncer = n.querySelector(".jw-volume-update"));
            var c,
              d,
              f,
              w = i.get("localization"),
              g = new Wt(i, e, n.querySelector(".jw-time-update")),
              j = (this.menus = []);
            this.ui = [];
            var b = "",
              m = w.volume;
            if (this._isMobile) {
              if (
                !(i.get("sdkplatform") || (o.OS.iOS && o.OS.version.major < 10))
              ) {
                var v = pt("volume-0,volume-100");
                f = p(
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
              var y = (c = new ie(
                i,
                "jw-icon-volume",
                m,
                pt("volume-0,volume-50,volume-100"),
                d
              )).element();
              j.push(c),
                Object(l.t)(y, "role", "button"),
                i.change(
                  "mute",
                  function (t, e) {
                    var i = e ? w.unmute : w.mute;
                    Object(l.t)(y, "aria-label", i);
                  },
                  this
                );
            }
            var k = p(
                "jw-icon-next",
                function () {
                  e.next({ feedShownId: b, reason: "interaction" });
                },
                w.next,
                pt("next")
              ),
              x = p(
                "jw-icon-settings jw-settings-submenu-button",
                function (t) {
                  s.trigger("settingsInteraction", "quality", !0, t);
                },
                w.settings,
                pt("settings")
              );
            Object(l.t)(x.element(), "aria-haspopup", "true");
            var T = p(
              "jw-icon-cc jw-settings-submenu-button",
              function (t) {
                s.trigger("settingsInteraction", "captions", !1, t);
              },
              w.cc,
              pt("cc-off,cc-on")
            );
            Object(l.t)(T.element(), "aria-haspopup", "true");
            var O = p(
              "jw-text-live",
              function () {
                s.goToLiveEdge();
              },
              w.liveBroadcast
            );
            O.element().textContent = w.liveBroadcast;
            var C,
              _,
              M,
              S = (this.elements = {
                alt:
                  ((C = "jw-text-alt"),
                  (_ = "status"),
                  (M = document.createElement("span")),
                  (M.className = "jw-text jw-reset " + C),
                  _ && Object(l.t)(M, "role", _),
                  M),
                play: p(
                  "jw-icon-playback",
                  function () {
                    e.playToggle({ reason: "interaction" });
                  },
                  w.play,
                  pt("play,pause,stop")
                ),
                rewind: p(
                  "jw-icon-rewind",
                  function () {
                    s.rewind();
                  },
                  w.rewind,
                  pt("rewind")
                ),
                live: O,
                next: k,
                elapsed: re("jw-text-elapsed", "timer"),
                countdown: re("jw-text-countdown", "timer"),
                time: g,
                duration: re("jw-text-duration", "timer"),
                mute: f,
                volumetooltip: c,
                horizontalVolumeContainer: d,
                cast: le(function () {
                  e.castToggle();
                }, w),
                fullscreen: p(
                  "jw-icon-fullscreen",
                  function () {
                    e.setFullscreen();
                  },
                  w.fullscreen,
                  pt("fullscreen-off,fullscreen-on")
                ),
                spacer: se("jw-spacer"),
                buttonContainer: se("jw-button-container"),
                settingsButton: x,
                captionsButton: T,
              }),
              E = ne(T.element(), "captions", w.cc),
              I = function (t) {
                var e = t.get("captionsList")[t.get("captionsIndex")],
                  i = w.cc;
                e && "Off" !== e.label && (i = e.label), E.setText(i);
              },
              L = ne(S.play.element(), "play", w.play);
            this.setPlayText = function (t) {
              L.setText(t);
            };
            var A = S.next.element(),
              P = ne(
                A,
                "next",
                w.nextUp,
                function () {
                  var t = i.get("nextUp");
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
              ne(S.rewind.element(), "rewind", w.rewind),
              ne(S.settingsButton.element(), "settings", w.settings);
            var R = ne(S.fullscreen.element(), "fullscreen", w.fullscreen),
              z = [
                S.play,
                S.rewind,
                S.next,
                S.volumetooltip,
                S.mute,
                S.horizontalVolumeContainer,
                S.alt,
                S.live,
                S.elapsed,
                S.countdown,
                S.duration,
                S.spacer,
                S.cast,
                S.captionsButton,
                S.settingsButton,
                S.fullscreen,
              ].filter(function (t) {
                return t;
              }),
              B = [S.time, S.buttonContainer].filter(function (t) {
                return t;
              });
            (this.el = document.createElement("div")),
              (this.el.className = "jw-controlbar jw-reset"),
              ue(S.buttonContainer, z),
              ue(this.el, B);
            var V = i.get("logo");
            if (
              (V && "control-bar" === V.position && this.addLogo(V),
              S.play.show(),
              S.fullscreen.show(),
              S.mute && S.mute.show(),
              i.change("volume", this.onVolume, this),
              i.change(
                "mute",
                function (t, e) {
                  s.renderVolume(e, t.get("volume"));
                },
                this
              ),
              i.change("state", this.onState, this),
              i.change("duration", this.onDuration, this),
              i.change("position", this.onElapsed, this),
              i.change(
                "fullscreen",
                function (t, e) {
                  var i = s.elements.fullscreen.element();
                  Object(l.v)(i, "jw-off", e);
                  var n = t.get("fullscreen") ? w.exitFullscreen : w.fullscreen;
                  R.setText(n), Object(l.t)(i, "aria-label", n);
                },
                this
              ),
              i.change("streamType", this.onStreamTypeChange, this),
              i.change(
                "dvrLive",
                function (t, e) {
                  var i = w.liveBroadcast,
                    n = w.notLive,
                    o = s.elements.live.element(),
                    a = !1 === e;
                  Object(l.v)(o, "jw-dvr-live", a),
                    Object(l.t)(o, "aria-label", a ? n : i),
                    (o.textContent = i);
                },
                this
              ),
              i.change("altText", this.setAltText, this),
              i.change("customButtons", this.updateButtons, this),
              i.on("change:captionsIndex", I, this),
              i.on("change:captionsList", I, this),
              i.change(
                "nextUp",
                function (t, e) {
                  b = Object(oe.b)(oe.a);
                  var i = w.nextUp;
                  e && e.title && (i += ": ".concat(e.title)),
                    P.setText(i),
                    S.next.toggle(!!e);
                },
                this
              ),
              i.change("audioMode", this.onAudioMode, this),
              S.cast &&
                (i.change("castAvailable", this.onCastAvailable, this),
                i.change("castActive", this.onCastActive, this)),
              S.volumetooltip &&
                (S.volumetooltip.on(
                  "update",
                  function (t) {
                    var e = t.percentage;
                    this._api.setVolume(e);
                  },
                  this
                ),
                S.volumetooltip.on(
                  "toggleValue",
                  function () {
                    this._api.setMute();
                  },
                  this
                ),
                S.volumetooltip.on(
                  "adjustVolume",
                  function (t) {
                    this.trigger("adjustVolume", t);
                  },
                  this
                )),
              S.cast && S.cast.button)
            ) {
              var N = S.cast.ui.on(
                "click tap enter",
                function (t) {
                  "click" !== t.type && S.cast.button.click(),
                    this._model.set("castClicked", !0);
                },
                this
              );
              this.ui.push(N);
            }
            var H = new u.a(S.duration).on(
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
              j.forEach(function (t) {
                t.on("open-tooltip", s.closeMenus, s);
              });
          }
          var e, i, n;
          return (
            (e = t),
            (i = [
              {
                key: "onVolume",
                value: function (t, e) {
                  this.renderVolume(t.get("mute"), e);
                },
              },
              {
                key: "renderVolume",
                value: function (t, e) {
                  var i = this.elements.mute,
                    n = this.elements.volumetooltip;
                  if (
                    (i &&
                      (Object(l.v)(i.element(), "jw-off", t),
                      Object(l.v)(i.element(), "jw-full", !t)),
                    n)
                  ) {
                    var o = t ? 0 : e,
                      a = n.element();
                    n.verticalSlider.render(o), n.horizontalSlider.render(o);
                    var r = n.tooltip,
                      s = n.horizontalContainer;
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
                  var i,
                    n,
                    o = t.get("duration");
                  if ("DVR" === t.get("streamType")) {
                    var a = Math.ceil(e),
                      r = this._model.get("dvrSeekLimit");
                    (i = n =
                      a >= -r ? "" : "-" + Object(vt.timeFormat)(-(e + r))),
                      t.set("dvrLive", a >= -r);
                  } else
                    (i = Object(vt.timeFormat)(e)),
                      (n = Object(vt.timeFormat)(o - e));
                  (this.elements.elapsed.textContent = i),
                    (this.elements.countdown.textContent = n);
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
                  var i = this.elements.time.element();
                  e
                    ? this.elements.buttonContainer.insertBefore(
                        i,
                        this.elements.elapsed
                      )
                    : Object(l.m)(this.el, i);
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
                    i = this._model.get("currentTime");
                  i
                    ? (t = i - 10)
                    : ((t = this._model.get("position") - 10),
                      "DVR" === this._model.get("streamType") &&
                        (e = this._model.get("duration"))),
                    this._api.seek(Math.max(t, e), { reason: "interaction" });
                },
              },
              {
                key: "onState",
                value: function (t, e) {
                  var i = t.get("localization"),
                    n = i.play;
                  this.setPlayText(n),
                    e === a.pb &&
                      ("LIVE" !== t.get("streamType")
                        ? ((n = i.pause), this.setPlayText(n))
                        : ((n = i.stop), this.setPlayText(n))),
                    Object(l.t)(this.elements.play.element(), "aria-label", n);
                },
              },
              {
                key: "onStreamTypeChange",
                value: function (t, e) {
                  var i = "LIVE" === e,
                    n = "DVR" === e;
                  this.elements.rewind.toggle(!i),
                    this.elements.live.toggle(i || n),
                    Object(l.t)(
                      this.elements.live.element(),
                      "tabindex",
                      i ? "-1" : "0"
                    ),
                    (this.elements.duration.style.display = n ? "none" : ""),
                    this.onDuration(t, t.get("duration")),
                    this.onState(t, t.get("state"));
                },
              },
              {
                key: "addLogo",
                value: function (t) {
                  var e = this.elements.buttonContainer,
                    i = new mt(
                      t.file,
                      this._model.get("localization").logo,
                      function () {
                        t.link &&
                          Object(l.l)(t.link, "_blank", { rel: "noreferrer" });
                      },
                      "logo",
                      "jw-logo-button"
                    );
                  t.link || Object(l.t)(i.element(), "tabindex", "-1"),
                    e.insertBefore(
                      i.element(),
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
                value: function (t, e, i) {
                  if (e) {
                    var n,
                      o,
                      a = this.elements.buttonContainer;
                    e !== i && i
                      ? ((n = ce(e, i)),
                        (o = ce(i, e)),
                        this.removeButtons(a, o))
                      : (n = e);
                    for (var r = n.length - 1; r >= 0; r--) {
                      var s = n[r],
                        l = new mt(
                          s.img,
                          s.tooltip,
                          s.callback,
                          s.id,
                          s.btnClass
                        );
                      s.tooltip && ne(l.element(), s.id, s.tooltip);
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
                  for (var i = e.length; i--; ) {
                    var n = t.querySelector('[button="'.concat(e[i].id, '"]'));
                    n && t.removeChild(n);
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
                      var i = t.elements[e];
                      i &&
                        "function" == typeof i.destroy &&
                        t.elements[e].destroy();
                    }),
                    this.ui.forEach(function (t) {
                      t.destroy();
                    }),
                    (this.ui = []);
                },
              },
            ]) && ae(e.prototype, i),
            n && ae(e, n),
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
        he = function (t) {
          return (
            '<div class="jw-display jw-reset"><div class="jw-display-container jw-reset"><div class="jw-display-controls jw-reset">' +
            pe("rewind", t.rewind) +
            pe("display", t.playback) +
            pe("next", t.next) +
            "</div></div></div>"
          );
        };
      function fe(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var we = (function () {
        function t(e, i, n) {
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, t);
          var o = n.querySelector(".jw-icon");
          (this.el = n),
            (this.ui = new u.a(o).on("click tap enter", function () {
              var t = e.get("position"),
                n = e.get("duration"),
                o = t - 10,
                a = 0;
              "DVR" === e.get("streamType") && (a = n), i.seek(Math.max(o, a));
            }));
        }
        var e, i, n;
        return (
          (e = t),
          (i = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && fe(e.prototype, i),
          n && fe(e, n),
          t
        );
      })();
      function ge(t) {
        return (ge =
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
      function je(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function be(t, e) {
        return !e || ("object" !== ge(e) && "function" != typeof e)
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
        function e(t, i, n) {
          var o;
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, e),
            (o = be(this, me(e).call(this)));
          var a = t.get("localization"),
            r = n.querySelector(".jw-icon");
          if (
            ((o.icon = r),
            (o.el = n),
            (o.ui = new u.a(r).on("click tap enter", function (t) {
              o.trigger(t.type);
            })),
            t.on("change:state", function (t, e) {
              var i;
              switch (e) {
                case "buffering":
                  i = a.buffer;
                  break;
                case "playing":
                  i = a.pause;
                  break;
                case "idle":
                case "paused":
                  i = a.playback;
                  break;
                case "complete":
                  i = a.replay;
                  break;
                default:
                  i = "";
              }
              "" !== i
                ? r.setAttribute("aria-label", i)
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
        var i, n, o;
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
          (i = e),
          (n = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && je(i.prototype, n),
          o && je(i, o),
          e
        );
      })(r.a);
      function ke(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var xe = (function () {
        function t(e, i, n) {
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, t);
          var o = n.querySelector(".jw-icon");
          (this.ui = new u.a(o).on("click tap enter", function () {
            i.next({ reason: "interaction" });
          })),
            e.change("nextUp", function (t, e) {
              n.style.visibility = e ? "" : "hidden";
            }),
            (this.el = n);
        }
        var e, i, n;
        return (
          (e = t),
          (i = [
            {
              key: "element",
              value: function () {
                return this.el;
              },
            },
          ]) && ke(e.prototype, i),
          n && ke(e, n),
          t
        );
      })();
      function Te(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var Oe = (function () {
        function t(e, i) {
          !(function (t, e) {
            if (!(t instanceof e))
              throw new TypeError("Cannot call a class as a function");
          })(this, t),
            (this.el = Object(l.e)(he(e.get("localization"))));
          var n = this.el.querySelector(".jw-display-controls"),
            o = {};
          Ce("rewind", pt("rewind"), we, n, o, e, i),
            Ce("display", pt("play,pause,buffer,replay"), ye, n, o, e, i),
            Ce("next", pt("next"), xe, n, o, e, i),
            (this.container = n),
            (this.buttons = o);
        }
        var e, i, n;
        return (
          (e = t),
          (i = [
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
          ]) && Te(e.prototype, i),
          n && Te(e, n),
          t
        );
      })();
      function Ce(t, e, i, n, o, a, r) {
        var s = n.querySelector(".jw-display-icon-".concat(t)),
          l = n.querySelector(".jw-icon-".concat(t));
        e.forEach(function (t) {
          l.appendChild(t);
        }),
          (o[t] = new i(a, r, s));
      }
      var _e = i(2);
      function Me(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var Se = (function () {
          function t(e, i, n) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(h.g)(this, r.a),
              (this._model = e),
              (this._api = i),
              (this._playerElement = n),
              (this.localization = e.get("localization")),
              (this.state = "tooltip"),
              (this.enabled = !1),
              (this.shown = !1),
              (this.feedShownId = ""),
              (this.closeUi = null),
              (this.tooltipUi = null),
              this.reset();
          }
          var e, i, n;
          return (
            (e = t),
            (i = [
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
                        i =
                          arguments.length > 2 && void 0 !== arguments[2]
                            ? arguments[2]
                            : "",
                        n =
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
                          i,
                          "</div>"
                        ) +
                        "</div></div>" +
                        '<button type="button" class="jw-icon jw-nextup-close jw-reset" aria-label="'.concat(
                          n,
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
                  var i = this._model,
                    n = i.player;
                  (this.enabled = !1),
                    i.on("change:nextUp", this.onNextUp, this),
                    n.change("duration", this.onDuration, this),
                    n.change("position", this.onElapsed, this),
                    n.change("streamType", this.onStreamType, this),
                    n.change(
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
                    var i = this._model.get("nextUp");
                    t && i
                      ? ((this.feedShownId = Object(oe.b)(oe.a)),
                        this.trigger("nextShown", {
                          mode: i.mode,
                          ui: "nextup",
                          itemsShown: [i],
                          feedData: i.feedData,
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
                      var i = e.loadThumbnail(t.image);
                      Object(gt.d)(e.thumbnail, i);
                    }
                    (e.header = e.content.querySelector(".jw-nextup-header")),
                      (e.header.textContent = Object(l.e)(
                        e.localization.nextUp
                      ).textContent),
                      (e.title = e.content.querySelector(".jw-nextup-title"));
                    var n = t.title;
                    e.title.textContent = n ? Object(l.e)(n).textContent : "";
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
                    var i = t.get("nextupoffset"),
                      n = -10;
                    i && (n = Object(_e.d)(i, e)),
                      n < 0 && (n += e),
                      Object(_e.c)(i) && e - 5 < n && (n = e - 5),
                      (this.offset = n);
                  }
                },
              },
              {
                key: "onElapsed",
                value: function (t, e) {
                  var i = this.nextUpSticky;
                  if (this.enabled && !1 !== i) {
                    var n = e >= this.offset;
                    n && void 0 === i
                      ? ((this.nextUpSticky = n), this.toggle(n, "time"))
                      : !n && i && this.reset();
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
            ]) && Me(e.prototype, i),
            n && Me(e, n),
            t
          );
        })(),
        Ee = function (t, e) {
          var i = t.featured,
            n = t.showLogo,
            o = t.type;
          return (
            (t.logo = n
              ? '<span class="jw-rightclick-logo jw-reset"></span>'
              : ""),
            '<li class="jw-reset jw-rightclick-item '
              .concat(i ? "jw-featured" : "", '">')
              .concat(Ie[o](t, e), "</li>")
          );
        },
        Ie = {
          link: function (t) {
            var e = t.link,
              i = t.title,
              n = t.logo;
            return '<a href="'
              .concat(
                e || "",
                '" class="jw-rightclick-link jw-reset-text" target="_blank" rel="noreferrer" dir="auto">'
              )
              .concat(n)
              .concat(i || "", "</a>");
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
        Le = i(23),
        Ae = i(6),
        Pe = i(13);
      function Re(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var ze = {
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
          i = e.querySelector(".jw-rightclick-logo");
        return i && i.appendChild(dt("jwplayer-logo")), e;
      }
      var Ve = (function () {
          function t(e, i) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              (this.infoOverlay = e),
              (this.shortcutsTooltip = i);
          }
          var e, i, n;
          return (
            (e = t),
            (i = [
              {
                key: "buildArray",
                value: function () {
                  var t = Le.a.split("+")[0],
                    e = this.model,
                    i = e.get("edition"),
                    n = e.get("localization").poweredBy,
                    o = '<span class="jw-reset">JW Player '.concat(
                      t,
                      "</span>"
                    ),
                    a = {
                      items: [
                        { type: "info" },
                        {
                          title: Object(Pe.e)(n)
                            ? "".concat(o, " ").concat(n)
                            : "".concat(n, " ").concat(o),
                          type: "link",
                          featured: !0,
                          showLogo: !0,
                          link: "https://jwplayer.com/learn-more?e=".concat(
                            ze[i]
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
                    i = t.pageX - e.left,
                    n = t.pageY - e.top;
                  return (
                    this.model.get("touchMode") && (n -= 100), { x: i, y: n }
                  );
                },
              },
              {
                key: "showMenu",
                value: function (t) {
                  var e = this,
                    i = this.getOffset(t);
                  return (
                    (this.el.style.left = i.x + "px"),
                    (this.el.style.top = i.y + "px"),
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
                    i,
                    n,
                    o = this,
                    a =
                      ((t = this.buildArray()),
                      (e = this.model.get("localization")),
                      (i = t.items),
                      (n = (void 0 === i ? [] : i).map(function (t) {
                        return Ee(t, e);
                      })),
                      '<div class="jw-rightclick jw-reset">' +
                        '<ul class="jw-rightclick-list jw-reset">'.concat(
                          n.join(""),
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
                value: function (t, e, i) {
                  (this.wrapperElement = i),
                    (this.model = t),
                    (this.mouseOverContext = !1),
                    (this.playerContainer = e),
                    (this.ui = new u.a(i).on(
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
            ]) && Re(e.prototype, i),
            n && Re(e, n),
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
      function qe(t, e) {
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
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function Xe(t, e, i) {
        return e && Ye(t.prototype, e), i && Ye(t, i), t;
      }
      var Ke,
        Je = (function () {
          function t(e, i) {
            var n =
              arguments.length > 2 && void 0 !== arguments[2]
                ? arguments[2]
                : Ne;
            Qe(this, t),
              (this.el = Object(l.e)(n(e))),
              (this.ui = new u.a(this.el).on("click tap enter", i, this));
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
        Ze = (function (t) {
          function e(t, i) {
            var n =
              arguments.length > 2 && void 0 !== arguments[2]
                ? arguments[2]
                : Fe;
            return Qe(this, e), qe(this, Ue(e).call(this, t, i, n));
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
        })(Je),
        Ge = function (t, e) {
          return t
            ? '<div class="jw-reset jw-settings-submenu jw-settings-submenu-'.concat(
                e,
                '" role="menu" aria-expanded="false">'
              ) + '<div class="jw-settings-submenu-items"></div></div>'
            : '<div class="jw-reset jw-settings-menu" role="menu" aria-expanded="false"><div class="jw-reset jw-settings-topbar" role="menubar"><div class="jw-reset jw-settings-topbar-text" tabindex="0"></div><div class="jw-reset jw-settings-topbar-buttons"></div></div></div>';
        },
        $e = function (t, e) {
          var i = t.name,
            n = {
              captions: "cc-off",
              audioTracks: "audio-tracks",
              quality: "quality-100",
              playbackRates: "playback-rate",
            }[i];
          if (n || t.icon) {
            var o = p(
                "jw-settings-".concat(i, " jw-submenu-").concat(i),
                function (e) {
                  t.open(e);
                },
                i,
                [(t.icon && Object(l.e)(t.icon)) || dt(n)]
              ),
              a = o.element();
            return (
              a.setAttribute("role", "menuitemradio"),
              a.setAttribute("aria-checked", "false"),
              a.setAttribute("aria-label", e),
              "ontouchstart" in window || (o.tooltip = ne(a, i, e)),
              o
            );
          }
        };
      function ti(t) {
        return (ti =
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
      function ei(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function ii(t) {
        return (ii = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function ni(t, e) {
        return (ni =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      function oi(t) {
        if (void 0 === t)
          throw new ReferenceError(
            "this hasn't been initialised - super() hasn't been called"
          );
        return t;
      }
      var ai = (function (t) {
          function e(t, i, n) {
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
                !(r = ii(e).call(this)) ||
                ("object" !== ti(r) && "function" != typeof r)
                  ? oi(a)
                  : r).open = o.open.bind(oi(oi(o)))),
              (o.close = o.close.bind(oi(oi(o)))),
              (o.toggle = o.toggle.bind(oi(oi(o)))),
              (o.onDocumentClick = o.onDocumentClick.bind(oi(oi(o)))),
              (o.name = t),
              (o.isSubmenu = !!i),
              (o.el = Object(l.e)(s(o.isSubmenu, t))),
              (o.topbar = o.el.querySelector(".jw-".concat(o.name, "-topbar"))),
              (o.buttonContainer = o.el.querySelector(
                ".jw-".concat(o.name, "-topbar-buttons")
              )),
              (o.children = {}),
              (o.openMenus = []),
              (o.items = []),
              (o.visible = !1),
              (o.parentMenu = i),
              (o.mainMenu = o.parentMenu ? o.parentMenu.mainMenu : oi(oi(o))),
              (o.categoryButton = null),
              (o.closeButton =
                (o.parentMenu && o.parentMenu.closeButton) ||
                o.createCloseButton(n)),
              o.isSubmenu
                ? ((o.categoryButton =
                    o.parentMenu.categoryButton || o.createCategoryButton(n)),
                  o.parentMenu.parentMenu &&
                    !o.mainMenu.backButton &&
                    (o.mainMenu.backButton = o.createBackButton(n)),
                  (o.itemsContainer = o.createItemsContainer()),
                  o.parentMenu.appendMenu(oi(oi(o))))
                : (o.ui = ri(oi(oi(o)))),
              o
            );
          }
          var i, n, o;
          return (
            (function (t, e) {
              if ("function" != typeof e && null !== e)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (t.prototype = Object.create(e && e.prototype, {
                constructor: { value: t, writable: !0, configurable: !0 },
              })),
                e && ni(t, e);
            })(e, t),
            (i = e),
            (n = [
              {
                key: "createItemsContainer",
                value: function () {
                  var t,
                    e,
                    i = this,
                    n = this.el.querySelector(".jw-settings-submenu-items"),
                    o = new u.a(n),
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
                      if (o.target.parentNode === n) {
                        var r = function (t, e) {
                            t
                              ? t.focus()
                              : void 0 !== e && n.childNodes[e].focus();
                          },
                          s = o.sourceEvent,
                          c = s.target,
                          u = n.firstChild === c,
                          d = n.lastChild === c,
                          p = i.topbar,
                          h = t || Object(l.k)(a),
                          f = e || Object(l.n)(a),
                          w = Object(l.k)(s.target),
                          g = Object(l.n)(s.target),
                          j = s.key.replace(/(Arrow|ape)/, "");
                        switch (j) {
                          case "Tab":
                            r(s.shiftKey ? f : h);
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
                              : r(g, n.childNodes.length - 1);
                            break;
                          case "Right":
                            r(h);
                            break;
                          case "Down":
                            p && d ? r(p.firstChild) : r(w, 0);
                        }
                        s.preventDefault(), "Esc" !== j && s.stopPropagation();
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
                          i = e.key.replace(/(Arrow|ape)/, "");
                        ("Enter" === i ||
                          "Right" === i ||
                          ("Tab" === i && !e.shiftKey)) &&
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
                  var i = $e(this, e);
                  return i.element().setAttribute("name", this.name), i;
                },
              },
              {
                key: "createBackButton",
                value: function (t) {
                  var e = p(
                    "jw-settings-back",
                    function (t) {
                      Ke && Ke.open(t);
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
                  var i = this,
                    n =
                      arguments.length > 2 && void 0 !== arguments[2]
                        ? arguments[2]
                        : {},
                    o =
                      arguments.length > 3 && void 0 !== arguments[3]
                        ? arguments[3]
                        : Ze,
                    a = this.name,
                    r = t.map(function (t, r) {
                      var s, l;
                      switch (a) {
                        case "quality":
                          s =
                            "Auto" === t.label && 0 === r
                              ? "".concat(
                                  n.defaultText,
                                  '&nbsp;<span class="jw-reset jw-auto-label"></span>'
                                )
                              : t.label;
                          break;
                        case "captions":
                          s =
                            ("Off" !== t.label && "off" !== t.id) || 0 !== r
                              ? t.label
                              : n.defaultText;
                          break;
                        case "playbackRates":
                          (l = t),
                            (s = Object(Pe.e)(n.tooltipText)
                              ? "x" + t
                              : t + "x");
                          break;
                        case "audioTracks":
                          s = t.name;
                      }
                      s || ((s = t), "object" === ti(t) && (s.options = n));
                      var c = new o(
                        s,
                        function (t) {
                          c.active ||
                            (e(l || r),
                            c.deactivate &&
                              (i.items
                                .filter(function (t) {
                                  return !0 === t.active;
                                })
                                .forEach(function (t) {
                                  t.deactivate();
                                }),
                              Ke ? Ke.open(t) : i.mainMenu.close(t)),
                            c.activate && c.activate());
                        }.bind(i)
                      );
                      return c;
                    });
                  return r;
                },
              },
              {
                key: "setMenuItems",
                value: function (t, e) {
                  var i = this;
                  t
                    ? ((this.items = []),
                      Object(l.h)(this.itemsContainer.el),
                      t.forEach(function (t) {
                        i.items.push(t), i.itemsContainer.el.appendChild(t.el);
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
                      i = t.name,
                      n = t.categoryButton;
                    if (((this.children[i] = t), n)) {
                      var o = this.mainMenu.buttonContainer,
                        a = o.querySelector(".jw-settings-sharing"),
                        r =
                          "quality" === i
                            ? o.firstChild
                            : a || this.closeButton.element();
                      o.insertBefore(n.element(), r);
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
                    if (((Ke = null), this.isSubmenu)) {
                      var i = this.mainMenu,
                        n = this.parentMenu,
                        o = this.categoryButton;
                      if (
                        (n.openMenus.length && n.closeChildren(),
                        o && o.element().setAttribute("aria-checked", "true"),
                        n.isSubmenu)
                      ) {
                        n.el.classList.remove("jw-settings-submenu-active"),
                          i.topbar.classList.add("jw-nested-menu-open");
                        var a = i.topbar.querySelector(
                          ".jw-settings-topbar-text"
                        );
                        a.setAttribute("name", this.name),
                          (a.innerText = this.title || this.name),
                          i.backButton.show(),
                          (Ke = this.parentMenu),
                          (e = this.topbar
                            ? this.topbar.firstChild
                            : t && "enter" === t.type
                            ? this.items[0].el
                            : a);
                      } else
                        i.topbar.classList.remove("jw-nested-menu-open"),
                          i.backButton && i.backButton.hide();
                      this.el.classList.add("jw-settings-submenu-active"),
                        n.openMenus.push(this.name),
                        i.visible ||
                          (i.open(t),
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
                    var i = t.children[e];
                    i && i.close();
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
                      i = e.indexOf(this.name);
                    e.length && i > -1 && this.openMenus.splice(i, 1),
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
                    i = t.captions,
                    n = t.audioTracks,
                    o = t.sharing,
                    a = t.playbackRates;
                  return e || i || n || o || a;
                },
              },
            ]) && ei(i.prototype, n),
            o && ei(i, o),
            e
          );
        })(r.a),
        ri = function (t) {
          var e = t.closeButton,
            i = t.el;
          return new u.a(i).on("keydown", function (i) {
            var n = i.sourceEvent,
              o = i.target,
              a = Object(l.k)(o),
              r = Object(l.n)(o),
              s = n.key.replace(/(Arrow|ape)/, ""),
              c = function (e) {
                r ? e || r.focus() : t.close(i);
              };
            switch (s) {
              case "Esc":
                t.close(i);
                break;
              case "Left":
                c();
                break;
              case "Right":
                a && e.element() && o !== e.element() && a.focus();
                break;
              case "Tab":
                n.shiftKey && c(!0);
                break;
              case "Up":
              case "Down":
                !(function () {
                  var e = t.children[o.getAttribute("name")];
                  if ((!e && Ke && (e = Ke.children[Ke.openMenus]), e))
                    return (
                      e.open(i),
                      void (e.topbar
                        ? e.topbar.firstChild.focus()
                        : e.items && e.items.length && e.items[0].el.focus())
                    );
                  if (
                    i.target.parentNode.classList.contains("jw-submenu-topbar")
                  ) {
                    var n = i.target.parentNode.parentNode.querySelector(
                      ".jw-settings-submenu-items"
                    );
                    ("Down" === s
                      ? n.childNodes[0]
                      : n.childNodes[n.childNodes.length - 1]
                    ).focus();
                  }
                })();
            }
            if ((n.stopPropagation(), /13|32|37|38|39|40/.test(n.keyCode)))
              return n.preventDefault(), !1;
          });
        },
        si = i(59),
        li = function (t) {
          return hi[t];
        },
        ci = function (t) {
          for (var e, i = Object.keys(hi), n = 0; n < i.length; n++)
            if (hi[i[n]] === t) {
              e = i[n];
              break;
            }
          return e;
        },
        ui = function (t) {
          return t + "%";
        },
        di = function (t) {
          return parseInt(t);
        },
        pi = [
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
            getTypedValue: li,
            getOption: ci,
          },
          {
            name: "Font Opacity",
            propertyName: "fontOpacity",
            options: ["100%", "75%", "25%"],
            defaultVal: "100%",
            getTypedValue: di,
            getOption: ui,
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
            getTypedValue: li,
            getOption: ci,
          },
          {
            name: "Background Opacity",
            propertyName: "backgroundOpacity",
            options: ["100%", "75%", "50%", "25%", "0%"],
            defaultVal: "50%",
            getTypedValue: di,
            getOption: ui,
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
            getTypedValue: li,
            getOption: ci,
          },
          {
            name: "Window Opacity",
            propertyName: "windowOpacity",
            options: ["100%", "75%", "50%", "25%", "0%"],
            defaultVal: "0%",
            getTypedValue: di,
            getOption: ui,
          },
        ],
        hi = {
          White: "#ffffff",
          Black: "#000000",
          Red: "#ff0000",
          Green: "#00ff00",
          Blue: "#0000ff",
          Yellow: "#ffff00",
          Magenta: "ff00ff",
          Cyan: "#00ffff",
        },
        fi = function (t, e, i, n) {
          var o = new ai("settings", null, n),
            a = function (t, e, a, r, s) {
              var l = i.elements["".concat(t, "Button")];
              if (!e || e.length <= 1)
                return o.removeMenu(t), void (l && l.hide());
              var c = o.children[t];
              c || (c = new ai(t, o, n)),
                c.setMenuItems(c.createItems(e, a, s), r),
                l && l.show();
            },
            r = function (r) {
              var s = { defaultText: n.auto };
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
              i.elements.settingsButton.toggle(c);
            };
          e.change(
            "levels",
            function (t, e) {
              r(e);
            },
            o
          );
          var s = function (t, i, n) {
            var o = e.get("levels");
            if (o && "Auto" === o[0].label && i && i.items.length) {
              var a = i.items[0].el.querySelector(".jw-auto-label"),
                r = o[t.index] || { label: "" };
              a.textContent = n ? "" : r.label;
            }
          };
          e.on("change:visualQuality", function (t, i) {
            var n = o.children.quality;
            i && n && s(i.level, n, e.get("currentLevel"));
          }),
            e.on(
              "change:currentLevel",
              function (t, i) {
                var n = o.children.quality,
                  a = e.get("visualQuality");
                a && n && s(a.level, n, i);
              },
              o
            ),
            e.change("captionsList", function (i, r) {
              var s = { defaultText: n.off },
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
                var u = new ai("captionsSettings", c, n);
                u.title = "Subtitle Settings";
                var d = new Je("Settings", u.open);
                c.topbar.appendChild(d.el);
                var p = new Ze("Reset", function () {
                  e.set("captions", si.a), w();
                });
                p.el.classList.add("jw-settings-reset");
                var f = e.get("captions"),
                  w = function () {
                    var t = [];
                    pi.forEach(function (i) {
                      f &&
                        f[i.propertyName] &&
                        (i.defaultVal = i.getOption(f[i.propertyName]));
                      var o = new ai(i.name, u, n),
                        a = new Je(
                          { label: i.name, value: i.defaultVal },
                          o.open,
                          He
                        ),
                        r = o.createItems(
                          i.options,
                          function (t) {
                            var n = a.el.querySelector(
                              ".jw-settings-content-item-value"
                            );
                            !(function (t, i) {
                              var n = e.get("captions"),
                                o = t.propertyName,
                                a = t.options && t.options[i],
                                r = t.getTypedValue(a),
                                s = Object(h.g)({}, n);
                              (s[o] = r), e.set("captions", s);
                            })(i, t),
                              (n.innerText = i.options[t]);
                          },
                          null
                        );
                      o.setMenuItems(r, i.options.indexOf(i.defaultVal) || 0),
                        t.push(a);
                    }),
                      t.push(p),
                      u.setMenuItems(t);
                  };
                w();
              }
            });
          var l = function (t, e) {
            t && e > -1 && t.items[e].activate();
          };
          e.change(
            "captionsIndex",
            function (t, e) {
              var n = o.children.captions;
              n && l(n, e), i.toggleCaptionsButtonState(!!e);
            },
            o
          );
          var c = function (i) {
            if (
              e.get("supportsPlaybackRate") &&
              "LIVE" !== e.get("streamType") &&
              e.get("playbackRateControls")
            ) {
              var r = i.indexOf(e.get("playbackRate")),
                s = { tooltipText: n.playbackRates };
              a(
                "playbackRates",
                i,
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
          var u = function (i) {
            a(
              "audioTracks",
              i,
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
              function (t, i) {
                var n = e.get("playbackRates"),
                  a = -1;
                n && (a = n.indexOf(i)), l(o.children.playbackRates, a);
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
                  i.elements.captionsButton.hide(),
                  o.visible && o.close();
              },
              o
            ),
            e.on("change:playbackRateControls", function () {
              c(e.get("playbackRates"));
            }),
            e.on(
              "change:castActive",
              function (t, i, n) {
                i !== n &&
                  (i
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
        wi = i(58),
        gi = i(35),
        ji = i(12),
        bi = function (t, e, i, n) {
          var o = Object(l.e)(
              '<div class="jw-reset jw-info-overlay jw-modal"><div class="jw-reset jw-info-container"><div class="jw-reset-text jw-info-title" dir="auto"></div><div class="jw-reset-text jw-info-duration" dir="auto"></div><div class="jw-reset-text jw-info-description" dir="auto"></div></div><div class="jw-reset jw-info-clientid"></div></div>'
            ),
            r = !1,
            s = null,
            c = !1,
            u = function (t) {
              /jw-info/.test(t.target.className) || h.close();
            },
            d = function () {
              var n,
                a,
                s,
                c,
                u,
                d = p(
                  "jw-info-close",
                  function () {
                    h.close();
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
                  var i = e.description,
                    n = e.title;
                  Object(l.q)(c, i || ""), Object(l.q)(a, n || "Unknown Title");
                }),
                e.change(
                  "duration",
                  function (t, i) {
                    var n = "";
                    switch (e.get("streamType")) {
                      case "LIVE":
                        n = "Live";
                        break;
                      case "DVR":
                        n = "DVR";
                        break;
                      default:
                        i && (n = Object(vt.timeFormat)(i));
                    }
                    s.textContent = n;
                  },
                  h
                ),
                (u.textContent =
                  (n = i.getPlugin("jwpsrv")) &&
                  "function" == typeof n.doNotTrackUser &&
                  n.doNotTrackUser()
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
          var h = {
            open: function () {
              r || d(), document.addEventListener("click", u), (c = !0);
              var t = e.get("state");
              t === a.pb && i.pause("infoOverlayInteraction"), (s = t), n(!0);
            },
            close: function () {
              document.removeEventListener("click", u),
                (c = !1),
                e.get("state") === a.ob &&
                  s === a.pb &&
                  i.play("infoOverlayInteraction"),
                (s = null),
                n(!1);
            },
            destroy: function () {
              this.close(), e.off(null, null, this);
            },
          };
          return (
            Object.defineProperties(h, {
              visible: {
                enumerable: !0,
                get: function () {
                  return c;
                },
              },
            }),
            h
          );
        };
      var mi = function (t, e, i) {
          var n,
            o = !1,
            r = null,
            s = i.get("localization").shortcuts,
            c = Object(l.e)(
              (function (t, e) {
                var i = t
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
                  "".concat(i) +
                  "</div></div></div></div>"
                );
              })(
                (function (t) {
                  var e = t.playPause,
                    i = t.volumeToggle,
                    n = t.fullscreenToggle,
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
                    { key: "f", description: n },
                    { key: "m", description: i },
                    { key: "0-9", description: o },
                  ];
                })(s),
                s.keyboardShortcuts
              )
            ),
            d = { reason: "settingsInteraction" },
            h = new u.a(c.querySelector(".jw-switch")),
            f = function () {
              h.el.setAttribute("aria-checked", i.get("enableShortcuts")),
                Object(l.a)(c, "jw-open"),
                (r = i.get("state")),
                c.querySelector(".jw-shortcuts-close").focus(),
                document.addEventListener("click", g),
                (o = !0),
                e.pause(d);
            },
            w = function () {
              Object(l.o)(c, "jw-open"),
                document.removeEventListener("click", g),
                t.focus(),
                (o = !1),
                r === a.pb && e.play(d);
            },
            g = function (t) {
              /jw-shortcuts|jw-switch/.test(t.target.className) || w();
            },
            j = function (t) {
              var e = t.currentTarget,
                n = "true" !== e.getAttribute("aria-checked");
              e.setAttribute("aria-checked", n), i.set("enableShortcuts", n);
            };
          return (
            (n = p("jw-shortcuts-close", w, i.get("localization").close, [
              dt("close"),
            ])),
            Object(l.m)(c, n.element()),
            n.show(),
            t.appendChild(c),
            h.on("click tap enter", j),
            {
              el: c,
              open: f,
              close: w,
              destroy: function () {
                w(), h.destroy();
              },
              toggleVisibility: function () {
                o ? w() : f();
              },
            }
          );
        },
        vi = function (t) {
          return (
            '<div class="jw-float-icon jw-icon jw-button-color jw-reset" aria-label='.concat(
              t,
              ' tabindex="0">'
            ) + "</div>"
          );
        };
      function yi(t) {
        return (yi =
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
      function ki(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function xi(t, e) {
        return !e || ("object" !== yi(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Ti(t) {
        return (Ti = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function Oi(t, e) {
        return (Oi =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      var Ci = (function (t) {
        function e(t, i) {
          var n;
          return (
            (function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
            ((n = xi(this, Ti(e).call(this))).element = Object(l.e)(vi(i))),
            n.element.appendChild(dt("close")),
            (n.ui = new u.a(n.element, { directSelect: !0 }).on(
              "click tap enter",
              function () {
                n.trigger(a.sb);
              }
            )),
            t.appendChild(n.element),
            n
          );
        }
        var i, n, o;
        return (
          (function (t, e) {
            if ("function" != typeof e && null !== e)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (t.prototype = Object.create(e && e.prototype, {
              constructor: { value: t, writable: !0, configurable: !0 },
            })),
              e && Oi(t, e);
          })(e, t),
          (i = e),
          (n = [
            {
              key: "destroy",
              value: function () {
                this.element &&
                  (this.ui.destroy(),
                  this.element.parentNode.removeChild(this.element),
                  (this.element = null));
              },
            },
          ]) && ki(i.prototype, n),
          o && ki(i, o),
          e
        );
      })(r.a);
      function _i(t) {
        return (_i =
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
      function Mi(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function Si(t, e) {
        return !e || ("object" !== _i(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function Ei(t) {
        return (Ei = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function Ii(t, e) {
        return (Ii =
          Object.setPrototypeOf ||
          function (t, e) {
            return (t.__proto__ = e), t;
          })(t, e);
      }
      i.d(e, "default", function () {
        return Ri;
      }),
        i(95);
      var Li = o.OS.mobile ? 4e3 : 2e3,
        Ai = [27];
      (gi.a.cloneIcon = dt),
        ji.a.forEach(function (t) {
          if (t.getState() === a.lb) {
            var e = t.getContainer().querySelector(".jw-error-msg .jw-icon");
            e && !e.hasChildNodes() && e.appendChild(gi.a.cloneIcon("error"));
          }
        });
      var Pi = function () {
          return { reason: "interaction" };
        },
        Ri = (function (t) {
          function e(t, i) {
            var n;
            return (
              (function (t, e) {
                if (!(t instanceof e))
                  throw new TypeError("Cannot call a class as a function");
              })(this, e),
              ((n = Si(this, Ei(e).call(this))).activeTimeout = -1),
              (n.inactiveTime = 0),
              (n.context = t),
              (n.controlbar = null),
              (n.displayContainer = null),
              (n.backdrop = null),
              (n.enabled = !0),
              (n.instreamState = null),
              (n.keydownCallback = null),
              (n.keyupCallback = null),
              (n.blurCallback = null),
              (n.mute = null),
              (n.nextUpToolTip = null),
              (n.playerContainer = i),
              (n.wrapperElement = i.querySelector(".jw-wrapper")),
              (n.rightClickMenu = null),
              (n.settingsMenu = null),
              (n.shortcutsTooltip = null),
              (n.showing = !1),
              (n.muteChangeCallback = null),
              (n.unmuteCallback = null),
              (n.logo = null),
              (n.div = null),
              (n.dimensions = {}),
              (n.infoOverlay = null),
              (n.userInactiveTimeout = function () {
                var t = n.inactiveTime - Object(c.a)();
                n.inactiveTime && t > 16
                  ? (n.activeTimeout = setTimeout(n.userInactiveTimeout, t))
                  : n.playerContainer.querySelector(".jw-tab-focus")
                  ? n.resetActiveTimeout()
                  : n.userInactive();
              }),
              n
            );
          }
          var i, n, r;
          return (
            (function (t, e) {
              if ("function" != typeof e && null !== e)
                throw new TypeError(
                  "Super expression must either be null or a function"
                );
              (t.prototype = Object.create(e && e.prototype, {
                constructor: { value: t, writable: !0, configurable: !0 },
              })),
                e && Ii(t, e);
            })(e, t),
            (i = e),
            (n = [
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
                  var i = this,
                    n = this.context.createElement("div");
                  (n.className = "jw-controls jw-reset"), (this.div = n);
                  var r = this.context.createElement("div");
                  (r.className = "jw-controls-backdrop jw-reset"),
                    (this.backdrop = r),
                    (this.logo = this.playerContainer.querySelector(
                      ".jw-logo"
                    ));
                  var c = e.get("touchMode"),
                    u = function () {
                      (e.get("isFloating")
                        ? i.wrapperElement
                        : i.playerContainer
                      ).focus();
                    };
                  if (!this.displayContainer) {
                    var d = new Oe(e, t);
                    d.buttons.display.on("click tap enter", function () {
                      i.trigger(a.p),
                        i.userActive(1e3),
                        t.playToggle(Pi()),
                        u();
                    }),
                      this.div.appendChild(d.element()),
                      (this.displayContainer = d);
                  }
                  (this.infoOverlay = new bi(n, e, t, function (t) {
                    Object(l.v)(i.div, "jw-info-open", t),
                      t && i.div.querySelector(".jw-info-close").focus();
                  })),
                    o.OS.mobile ||
                      (this.shortcutsTooltip = new mi(
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
                              ? i.rightClickMenu.destroy()
                              : i.rightClickMenu.setup(
                                  t,
                                  i.playerContainer,
                                  i.wrapperElement
                                );
                          },
                          this
                        );
                  var h = e.get("floating");
                  if (h) {
                    var f = new Ci(n, e.get("localization").close);
                    f.on(a.sb, function () {
                      return i.trigger("dismissFloating", { doNotForward: !0 });
                    }),
                      !1 !== h.dismissible &&
                        Object(l.a)(
                          this.playerContainer,
                          "jw-floating-dismissible"
                        );
                  }
                  var w = (this.controlbar = new de(
                    t,
                    e,
                    this.playerContainer.querySelector(
                      ".jw-hidden-accessibility"
                    )
                  ));
                  if (
                    (w.on(a.sb, function () {
                      return i.userActive();
                    }),
                    w.on(
                      "nextShown",
                      function (t) {
                        this.trigger("nextShown", t);
                      },
                      this
                    ),
                    w.on("adjustVolume", k, this),
                    e.get("nextUpDisplay") && !w.nextUpToolTip)
                  ) {
                    var g = new Se(e, t, this.playerContainer);
                    g.on("all", this.trigger, this),
                      g.setup(this.context),
                      (w.nextUpToolTip = g),
                      this.div.appendChild(g.element());
                  }
                  this.div.appendChild(w.element());
                  var j = e.get("localization"),
                    b = (this.settingsMenu = fi(
                      t,
                      e.player,
                      this.controlbar,
                      j
                    )),
                    m = null;
                  this.controlbar.on("menuVisibility", function (n) {
                    var o = n.visible,
                      r = n.evt,
                      s = e.get("state"),
                      l = { reason: "settingsInteraction" },
                      c = i.controlbar.elements.settingsButton,
                      d = "keydown" === ((r && r.sourceEvent) || r || {}).type,
                      p = o || d ? 0 : Li;
                    i.userActive(p),
                      (m = s),
                      Object(wi.a)(e.get("containerWidth")) < 2 &&
                        (o && s === a.pb
                          ? t.pause(l)
                          : o || s !== a.ob || m !== a.pb || t.play(l)),
                      !o && d && c ? c.element().focus() : r && u();
                  }),
                    b.on("menuVisibility", function (t) {
                      return i.controlbar.trigger("menuVisibility", t);
                    }),
                    this.controlbar.on(
                      "settingsInteraction",
                      function (t, e, i) {
                        if (e) return b.defaultChild.toggle(i);
                        b.children[t].toggle(i);
                      }
                    ),
                    o.OS.mobile
                      ? this.div.appendChild(b.el)
                      : (this.playerContainer.setAttribute(
                          "aria-describedby",
                          "jw-shortcuts-tooltip-explanation"
                        ),
                        this.div.insertBefore(b.el, w.element()));
                  var v = function (e) {
                    if (e.get("autostartMuted")) {
                      var n = function () {
                          return i.unmuteAutoplay(t, e);
                        },
                        a = function (t, e) {
                          e || n();
                        };
                      o.OS.mobile &&
                        ((i.mute = p(
                          "jw-autostart-mute jw-off",
                          n,
                          e.get("localization").unmute,
                          [dt("volume-0")]
                        )),
                        i.mute.show(),
                        i.div.appendChild(i.mute.element())),
                        w.renderVolume(!0, e.get("volume")),
                        Object(l.a)(i.playerContainer, "jw-flag-autostart"),
                        e.on("change:autostartFailed", n, i),
                        e.on("change:autostartMuted change:mute", a, i),
                        (i.muteChangeCallback = a),
                        (i.unmuteCallback = n);
                    }
                  };
                  function y(i) {
                    var n = 0,
                      o = e.get("duration"),
                      a = e.get("position");
                    if ("DVR" === e.get("streamType")) {
                      var r = e.get("dvrSeekLimit");
                      (n = o), (o = Math.max(a, -r));
                    }
                    var l = Object(s.a)(a + i, n, o);
                    t.seek(l, Pi());
                  }
                  function k(i) {
                    var n = Object(s.a)(e.get("volume") + i, 0, 100);
                    t.setVolume(n);
                  }
                  e.once("change:autostartMuted", v), v(e);
                  var x = function (n) {
                    if (n.ctrlKey || n.metaKey) return !0;
                    var o = !i.settingsMenu.visible,
                      a = !0 === e.get("enableShortcuts"),
                      r = i.instreamState;
                    if (a || -1 !== Ai.indexOf(n.keyCode)) {
                      switch (n.keyCode) {
                        case 27:
                          if (e.get("fullscreen"))
                            t.setFullscreen(!1),
                              i.playerContainer.blur(),
                              i.userInactive();
                          else {
                            var s = t.getPlugin("related");
                            s && s.close({ type: "escape" });
                          }
                          i.rightClickMenu.el &&
                            i.rightClickMenu.hideMenuHandler(),
                            i.infoOverlay.visible && i.infoOverlay.close(),
                            i.shortcutsTooltip && i.shortcutsTooltip.close();
                          break;
                        case 13:
                        case 32:
                          if (
                            document.activeElement.classList.contains(
                              "jw-switch"
                            ) &&
                            13 === n.keyCode
                          )
                            return !0;
                          t.playToggle(Pi());
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
                          i.shortcutsTooltip &&
                            i.shortcutsTooltip.toggleVisibility();
                          break;
                        default:
                          if (n.keyCode >= 48 && n.keyCode <= 59) {
                            var u = ((n.keyCode - 48) / 10) * e.get("duration");
                            t.seek(u, Pi());
                          }
                      }
                      return /13|32|37|38|39|40/.test(n.keyCode)
                        ? (n.preventDefault(), !1)
                        : void 0;
                    }
                  };
                  this.playerContainer.addEventListener("keydown", x),
                    (this.keydownCallback = x);
                  var T = function (t) {
                    switch (t.keyCode) {
                      case 9:
                        var e = i.playerContainer.contains(t.target) ? 0 : Li;
                        i.userActive(e);
                        break;
                      case 32:
                        t.preventDefault();
                    }
                  };
                  this.playerContainer.addEventListener("keyup", T),
                    (this.keyupCallback = T);
                  var O = function (t) {
                    var e = t.relatedTarget || document.querySelector(":focus");
                    e && (i.playerContainer.contains(e) || i.userInactive());
                  };
                  this.playerContainer.addEventListener("blur", O, !0),
                    (this.blurCallback = O);
                  var C = function t() {
                    "jw-shortcuts-tooltip-explanation" ===
                      i.playerContainer.getAttribute("aria-describedby") &&
                      i.playerContainer.removeAttribute("aria-describedby"),
                      i.playerContainer.removeEventListener("blur", t, !0);
                  };
                  this.shortcutsTooltip &&
                    (this.playerContainer.addEventListener("blur", C, !0),
                    (this.onRemoveShortcutsDescription = C)),
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
                    i = this.settingsMenu,
                    n = this.infoOverlay,
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
                    i && i.destroy(),
                    n && n.destroy(),
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
                  var i = !e.get("autostartFailed"),
                    n = e.get("mute");
                  i ? (n = !1) : e.set("playOnViewable", !1),
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
                    t.setMute(n),
                    this.controlbar.renderVolume(n, e.get("volume")),
                    this.mute && this.mute.hide(),
                    Object(l.o)(this.playerContainer, "jw-flag-autostart"),
                    this.userActive();
                },
              },
              {
                key: "mouseMove",
                value: function (t) {
                  var e = this.controlbar.element().contains(t.target),
                    i =
                      this.controlbar.nextUpToolTip &&
                      this.controlbar.nextUpToolTip
                        .element()
                        .contains(t.target),
                    n = this.logo && this.logo.contains(t.target),
                    o = e || i || n ? 0 : Li;
                  this.userActive(o);
                },
              },
              {
                key: "userActive",
                value: function () {
                  var t =
                    arguments.length > 0 && void 0 !== arguments[0]
                      ? arguments[0]
                      : Li;
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
            ]) && Mi(i.prototype, n),
            r && Mi(i, r),
            e
          );
        })(r.a);
    },
    function (t, e, i) {
      "use strict";
      i.r(e);
      var n = i(0),
        o = i(12),
        a = i(50),
        r = i(36);
      var s = i(44),
        l = i(51),
        c = i(26),
        u = i(25),
        d = i(3),
        p = i(46),
        h = i(2),
        f = i(7),
        w = i(34);
      function g(t) {
        var e = !1;
        return {
          async: function () {
            var i = this,
              n = arguments;
            return Promise.resolve().then(function () {
              if (!e) return t.apply(i, n);
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
      var j = i(1);
      function b(t) {
        return function (e, i) {
          var o = t.mediaModel,
            a = Object(n.g)({}, i, { type: e });
          switch (e) {
            case d.T:
              if (o.get(d.T) === i.mediaType) return;
              o.set(d.T, i.mediaType);
              break;
            case d.U:
              return void o.set(d.U, Object(n.g)({}, i));
            case d.M:
              if (i[e] === t.model.getMute()) return;
              break;
            case d.bb:
              i.newstate === d.mb && (t.thenPlayPromise.cancel(), o.srcReset());
              var r = o.attributes.mediaState;
              (o.attributes.mediaState = i.newstate),
                o.trigger("change:mediaState", o, i.newstate, r);
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
              var s = i.duration;
              Object(n.u)(s) &&
                (o.set("seekRange", i.seekRange), o.set("duration", s));
              break;
            case d.D:
              o.set("buffer", i.bufferPercent);
            case d.S:
              o.set("seekRange", i.seekRange),
                o.set("position", i.position),
                o.set("currentTime", i.currentTime);
              var l = i.duration;
              Object(n.u)(l) && o.set("duration", l),
                e === d.S &&
                  Object(n.r)(t.item.starttime) &&
                  delete t.item.starttime;
              break;
            case d.R:
              var c = t.mediaElement;
              c && c.paused && o.set("mediaState", "paused");
              break;
            case d.I:
              o.set(d.I, i.levels);
            case d.J:
              var u = i.currentQuality,
                p = i.levels;
              u > -1 && p.length > 1 && o.set("currentLevel", parseInt(u));
              break;
            case d.f:
              o.set(d.f, i.tracks);
            case d.g:
              var h = i.currentTrack,
                f = i.tracks;
              h > -1 &&
                f.length > 0 &&
                h < f.length &&
                o.set("currentAudioTrack", parseInt(h));
          }
          t.trigger(e, a);
        };
      }
      var m = i(8),
        v = i(45),
        y = i(41);
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
      function T(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function O(t, e, i) {
        return e && T(t.prototype, e), i && T(t, i), t;
      }
      function C(t, e) {
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
      function _(t) {
        return (_ = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function M(t, e) {
        if ("function" != typeof e && null !== e)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (t.prototype = Object.create(e && e.prototype, {
          constructor: { value: t, writable: !0, configurable: !0 },
        })),
          e && S(t, e);
      }
      function S(t, e) {
        return (S =
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
              ((t = C(this, _(e).call(this))).providerController = null),
              (t._provider = null),
              t.addAttributes({ mediaModel: new L() }),
              t
            );
          }
          return (
            M(e, t),
            O(e, [
              {
                key: "setup",
                value: function (t) {
                  return (
                    (t = t || {}),
                    this._normalizeConfig(t),
                    Object(n.g)(this.attributes, t, y.b),
                    (this.providerController = new w.a(
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
                    Object.keys(y.a).forEach(function (i) {
                      t[i] = e[i];
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
                  var i = e[t] || {},
                    o = i.label,
                    a = Object(n.u)(i.bitrate) ? i.bitrate : null;
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
                    (t = t || new L()),
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
                  if (Object(n.u)(t)) {
                    var e = Math.min(Math.max(0, t), 100);
                    this.set("volume", e);
                    var i = 0 === e;
                    i !== this.getMute() && this.setMute(i);
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
                  (this._provider = t), I(this, t);
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
                  Object(n.r)(t) &&
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
                  var e = t ? Object(h.g)(t.starttime) : 0,
                    i = t ? Object(h.g)(t.duration) : 0,
                    n = this.mediaModel;
                  this.set("playRejected", !1),
                    (this.attributes.itemMeta = {}),
                    n.set("position", e),
                    n.set("currentTime", 0),
                    n.set("duration", i);
                },
              },
              {
                key: "persistBandwidthEstimate",
                value: function (t) {
                  Object(n.u)(t) && this.set("bandwidthEstimate", t);
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
        I = function (t, e) {
          t.set("provider", e.getName()),
            !0 === t.get("instreamMode") && (e.instreamMode = !0),
            -1 === e.getName().name.indexOf("flash") &&
              (t.set("flashThrottle", void 0), t.set("flashBlocked", !1)),
            t.setPlaybackRate(t.get("defaultPlaybackRate")),
            t.set("supportsPlaybackRate", e.supportsPlaybackRate),
            t.set("playbackRate", e.getPlaybackRate()),
            t.set("renderCaptionsNatively", e.renderNatively);
        };
      var L = (function (t) {
          function e() {
            var t;
            return (
              x(this, e),
              (t = C(this, _(e).call(this))).addAttributes({
                mediaState: d.mb,
              }),
              t
            );
          }
          return (
            M(e, t),
            O(e, [
              {
                key: "srcReset",
                value: function () {
                  Object(n.g)(this.attributes, {
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
      function P(t) {
        return (P =
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
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function z(t) {
        return (z = Object.setPrototypeOf
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
        function e(t, i) {
          var n, o, a, r;
          return (
            (function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
            (o = this),
            (a = z(e).call(this)),
            ((n =
              !a || ("object" !== P(a) && "function" != typeof a)
                ? V(o)
                : a).attached = !0),
            (n.beforeComplete = !1),
            (n.item = null),
            (n.mediaModel = new L()),
            (n.model = i),
            (n.provider = t),
            (n.providerListener = new b(V(V(n)))),
            (n.thenPlayPromise = g(function () {})),
            (r = V(V(n))).provider.on("all", r.providerListener, r),
            (n.eventQueue = new s.a(V(V(n)), ["trigger"], function () {
              return !n.attached || n.background;
            })),
            n
          );
        }
        var i, o, a;
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
          (i = e),
          (o = [
            {
              key: "play",
              value: function (t) {
                var e = this.item,
                  i = this.model,
                  n = this.mediaModel,
                  o = this.provider;
                if (
                  (t || (t = i.get("playReason")),
                  i.set("playRejected", !1),
                  n.get("setup"))
                )
                  return o.play() || Promise.resolve();
                n.set("setup", !0);
                var a = this._loadAndPlay(e, o);
                return n.get("started") ? a : this._playAttempt(a, t);
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
                  i = this.provider;
                !t ||
                  (t && "none" === t.preload) ||
                  !this.attached ||
                  this.setup ||
                  this.preloaded ||
                  (e.set("preloaded", !0), i.preload(t));
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
                var i = this,
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
                      if (i.item && a === r.mediaModel) {
                        if ((r.set("playRejected", !0), l && l.paused)) {
                          if (l.src === location.href)
                            return i._loadAndPlay(o, s);
                          a.set("mediaState", d.ob);
                        }
                        var c = Object(n.g)(new j.n(null, Object(j.w)(t), t), {
                          error: t,
                          item: o,
                          playReason: e,
                        });
                        throw (delete c.key, i.trigger(d.O, c), t);
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
                  i = e.load(t);
                if (i) {
                  var n = g(function () {
                    return e.play() || Promise.resolve();
                  });
                  return (this.thenPlayPromise = n), i.then(n.async);
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
                  i = this.provider;
                i.video
                  ? e &&
                    (t
                      ? this.background ||
                        (this.thenPlayPromise.cancel(),
                        this.pause(),
                        e.removeChild(i.video),
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
                var e = (this.mediaModel = new L()),
                  i = t ? Object(h.g)(t.starttime) : 0,
                  n = t ? Object(h.g)(t.duration) : 0,
                  o = e.attributes;
                e.srcReset(),
                  (o.position = i),
                  (o.duration = n),
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
          ]) && R(i.prototype, o),
          a && R(i, a),
          e
        );
      })(f.a);
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
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function D(t) {
        return (D = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function q(t, e) {
        return (q =
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
        var i = e.mediaControllerListener;
        t.off().on("all", i, e);
      }
      function Q(t) {
        return t && t.sources && t.sources[0];
      }
      var Y = (function (t) {
        function e(t, i) {
          var o, a, r, s, l;
          return (
            (function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e),
            (a = this),
            ((o =
              !(r = D(e).call(this)) ||
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
            (o.mediaPool = i),
            (o.mediaController = null),
            (o.mediaControllerListener = (function (t, e) {
              return function (i, o) {
                switch (i) {
                  case d.bb:
                    return;
                  case "flashThrottle":
                  case "flashBlocked":
                    return void t.set(i, o.value);
                  case d.V:
                  case d.M:
                    return void t.set(i, o[i]);
                  case d.P:
                    return void t.set("playbackRate", o.playbackRate);
                  case d.K:
                    Object(n.g)(t.get("itemMeta"), o.metadata);
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
                    t.trigger(i, o);
                    break;
                  case d.i:
                    return void t.persistBandwidthEstimate(o.bandwidthEstimate);
                }
                e.trigger(i, o);
              };
            })(t, U(U(o)))),
            (o.model = t),
            (o.providers = new w.a(t.getConfiguration())),
            (o.loadPromise = Promise.resolve()),
            (o.backgroundLoading = t.get("backgroundLoading")),
            o.backgroundLoading ||
              t.set("mediaElement", o.mediaPool.getPrimedElement()),
            o
          );
        }
        var i, o, a;
        return (
          (function (t, e) {
            if ("function" != typeof e && null !== e)
              throw new TypeError(
                "Super expression must either be null or a function"
              );
            (t.prototype = Object.create(e && e.prototype, {
              constructor: { value: t, writable: !0, configurable: !0 },
            })),
              e && q(t, e);
          })(e, t),
          (i = e),
          (o = [
            {
              key: "setActiveItem",
              value: function (t) {
                var e = this,
                  i = this.model,
                  n = i.get("playlist")[t];
                (i.attributes.itemReady = !1), i.setActiveItem(t);
                var o = Q(n);
                if (!o) return Promise.reject(new j.n(j.k, j.h));
                var a = this.background,
                  r = this.mediaController;
                if (a.isNext(n))
                  return (
                    this._destroyActiveMedia(),
                    (this.loadPromise = this._activateBackgroundMedia()),
                    this.loadPromise
                  );
                if ((this._destroyBackgroundMedia(), r)) {
                  if (
                    i.get("castActive") ||
                    this._providerCanPlay(r.provider, o)
                  )
                    return (
                      (this.loadPromise = Promise.resolve(r)),
                      (r.activeItem = n),
                      this._setActiveMedia(r),
                      this.loadPromise
                    );
                  this._destroyActiveMedia();
                }
                var s = i.mediaModel;
                return (
                  (this.loadPromise = this._setupMediaController(o)
                    .then(function (t) {
                      if (s === i.mediaModel)
                        return (t.activeItem = n), e._setActiveMedia(t), t;
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
                    var i = e.detach(),
                      n = e.item,
                      o = e.mediaModel.get("position");
                    return o && (n.starttime = o), i;
                  }
                  e.attach();
                }
              },
            },
            {
              key: "playVideo",
              value: function (t) {
                var e,
                  i = this,
                  n = this.mediaController,
                  o = this.model;
                if (!o.get("playlistItem"))
                  return Promise.reject(new Error("No media"));
                if ((t || (t = o.get("playReason")), n)) e = n.play(t);
                else {
                  o.set(d.bb, d.jb);
                  var a = g(function (e) {
                    if (
                      i.mediaController &&
                      i.mediaController.mediaModel === e.mediaModel
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
                  i = e.get("playlist")[e.get("item")];
                (e.attributes.playlistItem = i), e.resetItem(i), t && t.stop();
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
                var i = this.model;
                i.attributes.itemReady = !1;
                var o = Object(n.g)({}, e),
                  a = (o.starttime = i.mediaModel.get("currentTime"));
                this._destroyActiveMedia();
                var r = new N(t, i);
                (r.activeItem = o),
                  this._setActiveMedia(r),
                  i.mediaModel.set("currentTime", a);
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
                  i = t.currentMedia;
                if (i) {
                  if (e)
                    return (
                      this._destroyMediaController(i),
                      void (t.currentMedia = null)
                    );
                  var n = i.mediaModel.attributes;
                  n.mediaState === d.mb
                    ? (n.mediaState = d.ob)
                    : n.mediaState !== d.ob && (n.mediaState = d.jb),
                    this._setActiveMedia(i),
                    (i.background = !1),
                    (t.currentMedia = null);
                }
              },
            },
            {
              key: "backgroundLoad",
              value: function (t) {
                var e = this.background,
                  i = Q(t);
                e.setNext(
                  t,
                  this._setupMediaController(i)
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
                  i = t.mediaModel,
                  n = t.provider;
                !(function (t, e) {
                  var i = t.get("mediaContainer");
                  i
                    ? (e.container = i)
                    : t.once("change:mediaContainer", function (t, i) {
                        e.container = i;
                      });
                })(e, t),
                  (this.mediaController = t),
                  e.set("mediaElement", t.mediaElement),
                  e.setMediaModel(i),
                  e.setProvider(n),
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
                  i = this.model,
                  n = this.providers,
                  o = function (t) {
                    return new N(
                      new t(i.get("id"), i.getConfiguration(), e.primedElement),
                      i
                    );
                  },
                  a = n.choose(t),
                  r = a.provider,
                  s = a.name;
                return r
                  ? Promise.resolve(o(r))
                  : n.load(s).then(function (t) {
                      return o(t);
                    });
              },
            },
            {
              key: "_activateBackgroundMedia",
              value: function () {
                var t = this,
                  e = this.background,
                  i = this.background.nextLoadPromise,
                  n = this.model;
                return (
                  this._destroyMediaController(e.currentMedia),
                  (e.currentMedia = null),
                  i.then(function (i) {
                    if (i)
                      return (
                        e.clearNext(),
                        t.adPlaying
                          ? ((n.attributes.itemReady = !0),
                            (e.currentMedia = i))
                          : (t._setActiveMedia(i), (i.background = !1)),
                        i
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
                  i = this.background.nextLoadPromise;
                i &&
                  i.then(function (i) {
                    t._destroyMediaController(i), e.clearNext();
                  });
              },
            },
            {
              key: "_providerCanPlay",
              value: function (t, e) {
                var i = this.providers.choose(e).provider;
                return i && t && t instanceof i;
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
                  i = this.mediaController,
                  n = this.mediaPool;
                i && (i.mute = t),
                  e.currentMedia && (e.currentMedia.mute = t),
                  n.syncMute(t);
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
                  i = this.mediaController,
                  n = this.mediaPool;
                i && (i.volume = t),
                  e.currentMedia && (e.currentMedia.volume = t),
                  n.syncVolume(t);
              },
            },
          ]) && F(i.prototype, o),
          a && F(i, a),
          e
        );
      })(f.a);
      function X(t) {
        return t === d.kb || t === d.lb ? d.mb : t;
      }
      function K(t, e, i) {
        if ((e = X(e)) !== (i = X(i))) {
          var n = e.replace(/(?:ing|d)$/, ""),
            o = {
              type: n,
              newstate: e,
              oldstate: i,
              reason: (function (t, e) {
                return t === d.jb ? (e === d.qb ? e : d.nb) : e;
              })(e, t.mediaModel.get("mediaState")),
            };
          "play" === n
            ? (o.playReason = t.get("playReason"))
            : "pause" === n && (o.pauseReason = t.get("pauseReason")),
            this.trigger(n, o);
        }
      }
      var J = i(48);
      function Z(t) {
        return (Z =
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
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function $(t, e) {
        return !e || ("object" !== Z(e) && "function" != typeof e)
          ? (function (t) {
              if (void 0 === t)
                throw new ReferenceError(
                  "this hasn't been initialised - super() hasn't been called"
                );
              return t;
            })(t)
          : e;
      }
      function tt(t, e, i, n) {
        return (tt =
          "undefined" != typeof Reflect && Reflect.set
            ? Reflect.set
            : function (t, e, i, n) {
                var o,
                  a = nt(t, e);
                if (a) {
                  if ((o = Object.getOwnPropertyDescriptor(a, e)).set)
                    return o.set.call(n, i), !0;
                  if (!o.writable) return !1;
                }
                if ((o = Object.getOwnPropertyDescriptor(n, e))) {
                  if (!o.writable) return !1;
                  (o.value = i), Object.defineProperty(n, e, o);
                } else
                  !(function (t, e, i) {
                    e in t
                      ? Object.defineProperty(t, e, {
                          value: i,
                          enumerable: !0,
                          configurable: !0,
                          writable: !0,
                        })
                      : (t[e] = i);
                  })(n, e, i);
                return !0;
              })(t, e, i, n);
      }
      function et(t, e, i, n, o) {
        if (!tt(t, e, i, n || t) && o)
          throw new Error("failed to set property");
        return i;
      }
      function it(t, e, i) {
        return (it =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, i) {
                var n = nt(t, e);
                if (n) {
                  var o = Object.getOwnPropertyDescriptor(n, e);
                  return o.get ? o.get.call(i) : o.value;
                }
              })(t, e, i || t);
      }
      function nt(t, e) {
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
          function e(t, i) {
            var n;
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, e);
            var o,
              a = ((n = $(this, ot(e).call(this, t, i))).model = new A());
            if (
              ((n.playerModel = t),
              (n.provider = null),
              (n.backgroundLoading = t.get("backgroundLoading")),
              (a.mediaModel.attributes.mediaType = "video"),
              n.backgroundLoading)
            )
              o = i.getAdElement();
            else {
              (o = t.get("mediaElement")),
                (a.attributes.mediaElement = o),
                (a.attributes.mediaSrc = o.src);
              var r = (n.srcResetListener = function () {
                n.srcReset();
              });
              o.addEventListener("emptied", r),
                (o.playbackRate = o.defaultPlaybackRate = 1);
            }
            return (n.mediaPool = Object(J.a)(o, i)), n;
          }
          var i, o, a;
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
            (i = e),
            (o = [
              {
                key: "setup",
                value: function () {
                  var t = this.model,
                    e = this.playerModel,
                    i = this.primedElement,
                    n = e.attributes,
                    o = e.mediaModel;
                  t.setup({
                    id: n.id,
                    volume: n.volume,
                    instreamMode: !0,
                    edition: n.edition,
                    mediaContext: o,
                    mute: n.mute,
                    streamType: "VOD",
                    autostartMuted: n.autostartMuted,
                    autostart: n.autostart,
                    advertising: n.advertising,
                    sdkplatform: n.sdkplatform,
                    skipButton: !1,
                  }),
                    t.on("change:state", K, this),
                    t.on(
                      d.w,
                      function (t) {
                        this.trigger(d.w, t);
                      },
                      this
                    ),
                    i.paused || i.pause();
                },
              },
              {
                key: "setActiveItem",
                value: function (t) {
                  var i = this;
                  return (
                    this.stopVideo(),
                    (this.provider = null),
                    it(ot(e.prototype), "setActiveItem", this)
                      .call(this, t)
                      .then(function (t) {
                        i._setProvider(t.provider);
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
                    var i = this.model,
                      o = this.playerModel,
                      a = "vpaid" === t.type;
                    t.off(),
                      t.on(
                        "all",
                        function (t, e) {
                          (a && t === d.F) ||
                            this.trigger(t, Object(n.g)({}, e, { type: t }));
                        },
                        this
                      );
                    var r = i.mediaModel;
                    t.on(d.bb, function (t) {
                      (t.oldstate = t.oldstate || i.get(d.bb)),
                        r.set("mediaState", t.newstate);
                    }),
                      t.on(d.X, this._nativeFullscreenHandler, this),
                      r.on("change:mediaState", function (t, i) {
                        e._stateHandler(i);
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
                            (i.set("autostartMuted", e),
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
                    i = this.playerModel;
                  t.off();
                  var n = e.getPrimedElement();
                  if (this.backgroundLoading) {
                    e.clean();
                    var o = i.get("mediaContainer");
                    n.parentNode === o && o.removeChild(n);
                  } else
                    n &&
                      (n.removeEventListener("emptied", this.srcResetListener),
                      n.src !== t.get("mediaSrc") && this.srcReset());
                },
              },
              {
                key: "srcReset",
                value: function () {
                  var t = this.playerModel,
                    e = t.get("mediaModel"),
                    i = t.getVideo();
                  e.srcReset(), i && (i.src = null);
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
                  var i = this.mediaController,
                    n = this.model,
                    o = this.provider;
                  n.set("mute", t),
                    et(ot(e.prototype), "mute", t, this, !0),
                    i || o.mute(t);
                },
              },
              {
                key: "volume",
                set: function (t) {
                  var i = this.mediaController,
                    n = this.model,
                    o = this.provider;
                  n.set("volume", t),
                    et(ot(e.prototype), "volume", t, this, !0),
                    i || o.volume(t);
                },
              },
            ]) && G(i.prototype, o),
            a && G(i, a),
            e
          );
        })(Y),
        st = { skipoffset: null, tag: null },
        lt = function (t, e, i, o) {
          var a,
            r,
            s,
            l,
            c = this,
            u = this,
            f = new rt(e, o),
            w = 0,
            g = {},
            j = null,
            b = {},
            m = A,
            v = !1,
            y = !1,
            k = !1,
            x = !1,
            T = function (t) {
              y ||
                (((t = t || {}).hasControls = !!e.get("controls")),
                c.trigger(d.z, t),
                f.model.get("state") === d.ob
                  ? t.hasControls && f.playVideo().catch(function () {})
                  : f.pause());
            },
            O = function () {
              y ||
                (f.model.get("state") === d.ob &&
                  e.get("controls") &&
                  (t.setFullscreen(), t.play()));
            };
          function C() {
            f.model.set("playRejected", !0);
          }
          function _() {
            w++, u.loadItem(a).catch(function () {});
          }
          function M(t, e) {
            "complete" !== t &&
              ((e = e || {}),
              b.tag && !e.tag && (e.tag = b.tag),
              this.trigger(t, e),
              ("mediaError" !== t && "error" !== t) ||
                (a && w + 1 < a.length && _()));
          }
          function S(t) {
            var e = t.newstate,
              i = t.oldstate || f.model.get("state");
            i !== e && E(Object(n.g)({ oldstate: i }, g, t));
          }
          function E(e) {
            var i = e.newstate;
            i === d.pb ? t.trigger(d.c, e) : i === d.ob && t.trigger(d.b, e);
          }
          function I(e) {
            var i = e.duration,
              n = e.position,
              o = f.model.mediaModel || f.model;
            o.set("duration", i),
              o.set("position", n),
              l || (l = (Object(h.d)(s, i) || i) - p.b),
              !v && n >= Math.max(l, p.a) && (t.preloadNextItem(), (v = !0));
          }
          function L(t) {
            var e = {};
            b.tag && (e.tag = b.tag), this.trigger(d.F, e), A.call(this, t);
          }
          function A(t) {
            (g = {}),
              a && w + 1 < a.length
                ? _()
                : (t.type === d.F && this.trigger(d.cb, {}), this.destroy());
          }
          function P() {
            y ||
              (i.clickHandler() &&
                i.clickHandler().setAlternateClickHandlers(T, O));
          }
          function R(t) {
            t.width && t.height && i.resizeMedia();
          }
          (this.init = function () {
            if (!k && !y) {
              (k = !0),
                (g = {}),
                f.setup(),
                f.on("all", M, this),
                f.on(d.O, C, this),
                f.on(d.S, I, this),
                f.on(d.F, L, this),
                f.on(d.K, R, this),
                f.on(d.bb, S, this),
                (j = t.detachMedia());
              var n = f.primedElement;
              e.get("mediaContainer").appendChild(n),
                e.set("instream", f),
                f.model.set("state", d.jb);
              var o = i.clickHandler();
              return (
                o && o.setAlternateClickHandlers(function () {}, null),
                this.setText(e.get("localization").loadingAd),
                (x = t.isBeforeComplete() || e.get("state") === d.kb),
                this
              );
            }
          }),
            (this.enableAdsMode = function (n) {
              var o = this;
              if (!k && !y)
                return (
                  t.routeEvents({
                    mediaControllerListener: function (t, e) {
                      o.trigger(t, e);
                    },
                  }),
                  e.set("instream", f),
                  f.model.set("state", d.pb),
                  (function (n) {
                    var o = i.clickHandler();
                    o &&
                      o.setAlternateClickHandlers(function (i) {
                        y ||
                          (((i = i || {}).hasControls = !!e.get("controls")),
                          u.trigger(d.z, i),
                          n &&
                            (e.get("state") === d.ob
                              ? t.playVideo()
                              : (t.pause(),
                                n &&
                                  (t.trigger(d.a, { clickThroughUrl: n }),
                                  window.open(n)))));
                      }, null);
                  })(n),
                  this
                );
            }),
            (this.setEventData = function (t) {
              g = t;
            }),
            (this.setState = function (t) {
              var e = t.newstate,
                i = f.model;
              (t.oldstate = i.get("state")), i.set("state", e), E(t);
            }),
            (this.setTime = function (e) {
              I(e), t.trigger(d.e, e);
            }),
            (this.loadItem = function (t, i) {
              if (y || !k)
                return Promise.reject(new Error("Instream not setup"));
              g = {};
              var o = t;
              Array.isArray(t)
                ? ((r = i || r), (t = (a = t)[w]), r && (i = r[w]))
                : (o = [t]);
              var l = f.model;
              l.set("playlist", o),
                e.set("hideAdsControls", !1),
                (t.starttime = 0),
                u.trigger(d.db, { index: w, item: t }),
                (b = Object(n.g)({}, st, i)),
                P(),
                l.set("skipButton", !1);
              var c =
                !e.get("backgroundLoading") && j
                  ? j.then(function () {
                      return f.setActiveItem(w);
                    })
                  : f.setActiveItem(w);
              return (
                (v = !1),
                void 0 !== (s = t.skipoffset || b.skipoffset) &&
                  u.setupSkipButton(s, b),
                c
              );
            }),
            (this.setupSkipButton = function (t, e, i) {
              var n = f.model;
              (m = i || A),
                n.set("skipMessage", e.skipMessage),
                n.set("skipText", e.skipText),
                n.set("skipOffset", t),
                (n.attributes.skipButton = !1),
                n.set("skipButton", !0);
            }),
            (this.applyProviderListeners = function (t) {
              f.usePsuedoProvider(t), P();
            }),
            (this.play = function () {
              (g = {}), f.playVideo();
            }),
            (this.pause = function () {
              (g = {}), f.pause();
            }),
            (this.skipAd = function (t) {
              var i = e.get("autoPause").pauseAds,
                n = "autostart" === e.get("playReason"),
                o = e.get("viewable");
              !i || n || o || (this.noResume = !0);
              var a = d.d;
              this.trigger(a, t), m.call(this, { type: a });
            }),
            (this.replacePlaylistItem = function (t) {
              y || (e.set("playlistItem", t), f.srcReset());
            }),
            (this.destroy = function () {
              y ||
                ((y = !0),
                this.trigger("destroyed"),
                this.off(),
                i.clickHandler() &&
                  i.clickHandler().revertAlternateClickHandlers(),
                e.off(null, null, f),
                f.off(null, null, u),
                f.destroy(),
                k && f.model && (e.attributes.state = d.ob),
                t.forwardEvents(),
                e.set("instream", null),
                (f = null),
                (g = {}),
                (j = null),
                k &&
                  !e.attributes._destroyed &&
                  (t.attachMedia(),
                  this.noResume || (x ? t.stopVideo() : t.playVideo())));
            }),
            (this.getState = function () {
              return !y && f.model.get("state");
            }),
            (this.setText = function (t) {
              return y ? this : (i.setAltText(t || ""), this);
            }),
            (this.hide = function () {
              y || e.set("hideAdsControls", !0);
            }),
            (this.getMediaElement = function () {
              return y ? null : f.primedElement;
            }),
            (this.setSkipOffset = function (t) {
              (s = t > 0 ? t : null), f && f.model.set("skipOffset", s);
            });
        };
      Object(n.g)(lt.prototype, f.a);
      var ct = lt,
        ut = i(66),
        dt = i(63),
        pt = function (t) {
          var e = this,
            i = [],
            n = {},
            o = 0,
            a = 0;
          function r(t) {
            if (
              ((t.data = t.data || []),
              (t.name = t.label || t.name || t.language),
              (t._id = Object(dt.a)(t, i.length)),
              !t.name)
            ) {
              var e = Object(dt.b)(t, o);
              (t.name = e.label), (o = e.unknownCount);
            }
            (n[t._id] = t), i.push(t);
          }
          function s() {
            for (
              var t = [{ id: "off", label: "Off" }], e = 0;
              e < i.length;
              e++
            )
              t.push({
                id: i[e]._id,
                label: i[e].name || "Unknown CC",
                language: i[e].language,
              });
            return t;
          }
          function l(e) {
            var n = (a = e),
              o = t.get("captionLabel");
            if ("Off" !== o) {
              for (var r = 0; r < i.length; r++) {
                var s = i[r];
                if (o && o === s.name) {
                  n = r + 1;
                  break;
                }
                s.default || s.defaulttrack || "default" === s._id
                  ? (n = r + 1)
                  : s.autoselect;
              }
              var l;
              (l = n),
                i.length
                  ? t.setVideoSubtitleTrack(l, i)
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
              (i = []), (n = {}), (o = 0);
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
                var i = t.get("playlistItem").tracks,
                  o = i && i.length;
                if (o && !t.get("renderCaptionsNatively"))
                  for (
                    var a = function (t) {
                        var o,
                          a = i[t];
                        ("subtitles" !== (o = a.kind) && "captions" !== o) ||
                          n[a._id] ||
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
                var n = null;
                0 !== e && (n = i[e - 1]), t.set("captionsTrack", n);
              },
              this
            ),
            (this.setSubtitlesTracks = function (t) {
              if (Array.isArray(t)) {
                if (t.length) {
                  for (var e = 0; e < t.length; e++) r(t[e]);
                  i = Object.keys(n).map(function (t) {
                    return n[t];
                  });
                } else (i = []), (n = {}), (o = 0);
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
      Object(n.g)(pt.prototype, f.a);
      var ht = pt,
        ft = function (t, e) {
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
        wt = i(35),
        gt = 44,
        jt = function (t) {
          var e = t.get("height");
          if (t.get("aspectratio")) return !1;
          if ("string" == typeof e && e.indexOf("%") > -1) return !1;
          var i = 1 * e || NaN;
          return (
            !!(i = isNaN(i) ? t.get("containerHeight") : i) && i && i <= gt
          );
        },
        bt = i(54);
      function mt(t, e) {
        if (t.get("fullscreen")) return 1;
        if (!t.get("activeTab")) return 0;
        if (t.get("isFloating")) return 1;
        var i = t.get("intersectionRatio");
        return void 0 === i &&
          ((i = (function (t) {
            var e = document.documentElement,
              i = document.body,
              n = {
                top: 0,
                left: 0,
                right: e.clientWidth || i.clientWidth,
                width: e.clientWidth || i.clientWidth,
                bottom: e.clientHeight || i.clientHeight,
                height: e.clientHeight || i.clientHeight,
              };
            if (!i.contains(t)) return 0;
            if ("none" === window.getComputedStyle(t).display) return 0;
            var o = vt(t);
            if (!o) return 0;
            var a = o,
              r = t.parentNode,
              s = !1;
            for (; !s; ) {
              var l = null;
              if (
                (r === i || r === e || 1 !== r.nodeType
                  ? ((s = !0), (l = n))
                  : "visible" !== window.getComputedStyle(r).overflow &&
                    (l = vt(r)),
                l &&
                  ((c = l),
                  (u = a),
                  (d = void 0),
                  (p = void 0),
                  (h = void 0),
                  (f = void 0),
                  (w = void 0),
                  (g = void 0),
                  (d = Math.max(c.top, u.top)),
                  (p = Math.min(c.bottom, u.bottom)),
                  (h = Math.max(c.left, u.left)),
                  (f = Math.min(c.right, u.right)),
                  (g = p - d),
                  !(a = (w = f - h) >= 0 &&
                    g >= 0 && {
                      top: d,
                      bottom: p,
                      left: h,
                      right: f,
                      width: w,
                      height: g,
                    })))
              )
                return 0;
              r = r.parentNode;
            }
            var c, u, d, p, h, f, w, g;
            var j = o.width * o.height,
              b = a.width * a.height;
            return j ? b / j : 0;
          })(e)),
          window.top !== window.self && i)
          ? 0
          : i;
      }
      function vt(t) {
        try {
          return t.getBoundingClientRect();
        } catch (t) {}
      }
      var yt = i(49),
        kt = i(42),
        xt = i(58),
        Tt = i(10);
      var Ot = i(32),
        Ct = i(5),
        _t = i(6),
        Mt = [
          "fullscreenchange",
          "webkitfullscreenchange",
          "mozfullscreenchange",
          "MSFullscreenChange",
        ],
        St = function (t, e, i) {
          for (
            var n =
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
              a = !(!n || !o),
              r = Mt.length;
            r--;

          )
            e.addEventListener(Mt[r], i);
          return {
            events: Mt,
            supportsDomFullscreen: function () {
              return a;
            },
            requestFullscreen: function () {
              n.call(t, { navigationUI: "hide" });
            },
            exitFullscreen: function () {
              null !== this.fullscreenElement() && o.apply(e);
            },
            fullscreenElement: function () {
              var t = e.fullscreenElement,
                i = e.webkitCurrentFullScreenElement,
                n = e.mozFullScreenElement,
                o = e.msFullscreenElement;
              return null === t ? t : t || i || n || o;
            },
            destroy: function () {
              for (var t = Mt.length; t--; ) e.removeEventListener(Mt[t], i);
            },
          };
        },
        Et = i(40);
      function It(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var Lt = (function () {
          function t(e, i) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              Object(n.g)(this, f.a),
              this.revertAlternateClickHandlers(),
              (this.domElement = i),
              (this.model = e),
              (this.ui = new Et.a(i)
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
          var e, i, o;
          return (
            (e = t),
            (i = [
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
            ]) && It(e.prototype, i),
            o && It(e, o),
            t
          );
        })(),
        At = i(59),
        Pt = function (t, e) {
          var i = e ? " jw-hide" : "";
          return '<div class="jw-logo jw-logo-'
            .concat(t)
            .concat(i, ' jw-reset"></div>');
        },
        Rt = {
          linktarget: "_blank",
          margin: 8,
          hide: !1,
          position: "top-right",
        };
      function zt(t) {
        var e, i;
        Object(n.g)(this, f.a);
        var o = new Image();
        (this.setup = function () {
          ((i = Object(n.g)({}, Rt, t.get("logo"))).position =
            i.position || Rt.position),
            (i.hide = "true" === i.hide.toString()),
            i.file &&
              "control-bar" !== i.position &&
              (e || (e = Object(Ct.e)(Pt(i.position, i.hide))),
              t.set("logo", i),
              (o.onload = function () {
                var n = this.height,
                  o = this.width,
                  a = { backgroundImage: 'url("' + this.src + '")' };
                if (i.margin !== Rt.margin) {
                  var r = /(\w+)-(\w+)/.exec(i.position);
                  3 === r.length &&
                    ((a["margin-" + r[1]] = i.margin),
                    (a["margin-" + r[2]] = i.margin));
                }
                var s = 0.15 * t.get("containerHeight"),
                  l = 0.15 * t.get("containerWidth");
                if (n > s || o > l) {
                  var c = o / n;
                  l / s > c ? ((n = s), (o = s * c)) : ((o = l), (n = l / c));
                }
                (a.width = Math.round(o)),
                  (a.height = Math.round(n)),
                  Object(Tt.d)(e, a),
                  t.set("logoWidth", a.width);
              }),
              (o.src = i.file),
              i.link &&
                (e.setAttribute("tabindex", "0"),
                e.setAttribute("aria-label", t.get("localization").logo)),
              (this.ui = new Et.a(e).on(
                "click tap enter",
                function (t) {
                  t && t.stopPropagation && t.stopPropagation(),
                    this.trigger(d.A, {
                      link: i.link,
                      linktarget: i.linktarget,
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
            return i.position;
          }),
          (this.destroy = function () {
            (o.onload = null), this.ui && this.ui.destroy();
          });
      }
      var Bt = function (t) {
        (this.model = t), (this.image = null);
      };
      Object(n.g)(Bt.prototype, {
        setup: function (t) {
          this.el = t;
        },
        setImage: function (t) {
          var e = this.image;
          e && (e.onload = null), (this.image = null);
          var i = "";
          "string" == typeof t &&
            ((i = 'url("' + t + '")'),
            ((e = this.image = new Image()).src = t)),
            Object(Tt.d)(this.el, { backgroundImage: i });
        },
        resize: function (t, e, i) {
          if ("uniform" === i) {
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
            var n = this.image,
              o = null;
            if (n) {
              if (0 === n.width) {
                var a = this;
                return void (n.onload = function () {
                  a.resize(t, e, i);
                });
              }
              var r = n.width / n.height;
              Math.abs(this.playerAspectRatio - r) < 0.09 && (o = "cover");
            }
            Object(Tt.d)(this.el, { backgroundSize: o });
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
      Object(n.g)(Nt.prototype, {
        hide: function () {
          Object(Tt.d)(this.el, { display: "none" });
        },
        show: function () {
          Object(Tt.d)(this.el, { display: "" });
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
            i = t.get("logo");
          if (i) {
            var n = 1 * ("" + i.margin).replace("px", ""),
              o = t.get("logoWidth") + (isNaN(n) ? 0 : n + 10);
            "top-left" === i.position
              ? (e.paddingLeft = o)
              : "top-right" === i.position && (e.paddingRight = o);
          }
          Object(Tt.d)(this.el, e);
        },
        playlistItem: function (t, e) {
          if (e)
            if (t.get("displaytitle") || t.get("displaydescription")) {
              var i = "",
                n = "";
              e.title && t.get("displaytitle") && (i = e.title),
                e.description &&
                  t.get("displaydescription") &&
                  (n = e.description),
                this.updateText(i, n);
            } else this.hide();
        },
        updateText: function (t, e) {
          Object(Ct.q)(this.title, t),
            Object(Ct.q)(this.description, e),
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
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      var Dt,
        qt = (function () {
          function t(e) {
            !(function (t, e) {
              if (!(t instanceof e))
                throw new TypeError("Cannot call a class as a function");
            })(this, t),
              (this.container = e),
              (this.input = e.querySelector(".jw-media"));
          }
          var e, i, n;
          return (
            (e = t),
            (i = [
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
                    i,
                    n,
                    o = this.container,
                    a = this.input,
                    r = (this.ui = new Et.a(a, { preventScrolling: !0 })
                      .on("dragStart", function () {
                        (t = o.offsetLeft),
                          (e = o.offsetTop),
                          (i = window.innerHeight),
                          (n = window.innerWidth);
                      })
                      .on("drag", function (a) {
                        var s = Math.max(t + a.pageX - r.startX, 0),
                          l = Math.max(e + a.pageY - r.startY, 0),
                          c = Math.max(n - (s + o.clientWidth), 0),
                          u = Math.max(i - (l + o.clientHeight), 0);
                        0 === c ? (s = "auto") : (c = "auto"),
                          0 === l ? (u = "auto") : (l = "auto"),
                          Object(Tt.d)(o, {
                            left: s,
                            right: c,
                            top: l,
                            bottom: u,
                            margin: 0,
                          });
                      })
                      .on("dragEnd", function () {
                        t = e = n = i = null;
                      }));
                },
              },
            ]) && Ft(e.prototype, i),
            n && Ft(e, n),
            t
          );
        })(),
        Ut = i(55);
      i(69);
      var Wt = m.OS.mobile,
        Qt = m.Browser.ie,
        Yt = null;
      var Xt = function (t, e) {
        var i,
          o,
          a,
          r,
          s = this,
          l = Object(n.g)(this, f.a, { isSetup: !1, api: t, model: e }),
          c = e.get("localization"),
          u = Object(Ct.e)(ft(e.get("id"), c.player)),
          p = u.querySelector(".jw-wrapper"),
          w = u.querySelector(".jw-media"),
          g = new qt(p),
          j = new Vt(e, t),
          b = new Ht(e),
          v = new At.b(e);
        v.on("all", l.trigger, l);
        var y = -1,
          k = -1,
          x = -1,
          T = e.get("floating");
        this.dismissible = T && T.dismissible;
        var O,
          C,
          _,
          M = !1,
          S = {},
          E = null,
          I = null;
        function L() {
          return Wt && !Object(Ct.f)();
        }
        function A() {
          Object(kt.a)(k), (k = Object(kt.b)(P));
        }
        function P() {
          l.isSetup && (l.updateBounds(), l.updateStyles(), l.checkResized());
        }
        function R(t, i) {
          if (Object(n.r)(t) && Object(n.r)(i)) {
            var o = Object(xt.a)(t);
            Object(xt.b)(u, o);
            var a = o < 2;
            Object(Ct.v)(u, "jw-flag-small-player", a),
              Object(Ct.v)(u, "jw-orientation-portrait", i > t);
          }
          if (e.get("controls")) {
            var r = jt(e);
            Object(Ct.v)(u, "jw-flag-audio-player", r), e.set("audioMode", r);
          }
        }
        function z() {
          e.set("visibility", mt(e, u));
        }
        (this.updateBounds = function () {
          Object(kt.a)(k);
          var t = e.get("isFloating") ? p : u,
            i = document.body.contains(t),
            n = Object(Ct.c)(t),
            r = Math.round(n.width),
            s = Math.round(n.height);
          if (((S = Object(Ct.c)(u)), r === o && s === a))
            return (o && a) || A(), void e.set("inDom", i);
          (r && s) || (o && a) || A(),
            (r || s || i) &&
              (e.set("containerWidth", r), e.set("containerHeight", s)),
            e.set("inDom", i),
            i && bt.a.observe(u);
        }),
          (this.updateStyles = function () {
            var t = e.get("containerWidth"),
              i = e.get("containerHeight");
            R(t, i), I && I.resize(t, i), $(t, i), v.resize(), T && F();
          }),
          (this.checkResized = function () {
            var t = e.get("containerWidth"),
              i = e.get("containerHeight"),
              n = e.get("isFloating");
            if (t !== o || i !== a) {
              this.resizeListener ||
                (this.resizeListener = new Ut.a(p, this, e)),
                (o = t),
                (a = i),
                l.trigger(d.hb, { width: t, height: i });
              var s = Object(xt.a)(t);
              E !== s && ((E = s), l.trigger(d.j, { breakpoint: E }));
            }
            n !== r && ((r = n), l.trigger(d.x, { floating: n }), z());
          }),
          (this.responsiveListener = A),
          (this.setup = function () {
            j.setup(u.querySelector(".jw-preview")),
              b.setup(u.querySelector(".jw-title")),
              (i = new zt(e)).setup(),
              i.setContainer(p),
              i.on(d.A, J),
              v.setup(u.id, e.get("captions")),
              b.element().parentNode.insertBefore(v.element(), b.element()),
              (O = (function (t, e, i) {
                var n = new Lt(e, i),
                  o = e.get("controls");
                n.on({
                  click: function () {
                    l.trigger(d.p),
                      I &&
                        (ct()
                          ? I.settingsMenu.close()
                          : ut()
                          ? I.infoOverlay.close()
                          : t.playToggle({ reason: "interaction" }));
                  },
                  tap: function () {
                    l.trigger(d.p),
                      ct() && I.settingsMenu.close(),
                      ut() && I.infoOverlay.close();
                    var i = e.get("state");
                    if (
                      (o &&
                        (i === d.mb ||
                          i === d.kb ||
                          (e.get("instream") && i === d.ob)) &&
                        t.playToggle({ reason: "interaction" }),
                      o && i === d.ob)
                    ) {
                      if (
                        e.get("instream") ||
                        e.get("castActive") ||
                        "audio" === e.get("mediaType")
                      )
                        return;
                      Object(Ct.v)(u, "jw-flag-controls-hidden"),
                        l.dismissible &&
                          Object(Ct.v)(
                            u,
                            "jw-floating-dismissible",
                            Object(Ct.i)(u, "jw-flag-controls-hidden")
                          ),
                        v.renderCues(!0);
                    } else I && (I.showing ? I.userInactive() : I.userActive());
                  },
                  doubleClick: function () {
                    return I && t.setFullscreen();
                  },
                }),
                  Wt ||
                    (u.addEventListener("mousemove", W),
                    u.addEventListener("mouseover", Q),
                    u.addEventListener("mouseout", Y));
                return n;
              })(t, e, w)),
              (_ = new Et.a(u).on("click", function () {})),
              (C = St(u, document, et)),
              e.on("change:hideAdsControls", function (t, e) {
                Object(Ct.v)(u, "jw-flag-ads-hide-controls", e);
              }),
              e.on("change:scrubbing", function (t, e) {
                Object(Ct.v)(u, "jw-flag-dragging", e);
              }),
              e.on("change:playRejected", function (t, e) {
                Object(Ct.v)(u, "jw-flag-play-rejected", e);
              }),
              e.on(d.X, tt),
              e.on("change:".concat(d.U), function () {
                $(), v.resize();
              }),
              e.player.on("change:errorEvent", at),
              e.change("stretching", X);
            var n = e.get("width"),
              o = e.get("height"),
              a = G(n, o);
            Object(Tt.d)(u, a),
              e.change("aspectratio", K),
              R(n, o),
              e.get("controls") ||
                (Object(Ct.a)(u, "jw-flag-controls-hidden"),
                Object(Ct.o)(u, "jw-floating-dismissible")),
              Qt && Object(Ct.a)(u, "jw-ie");
            var r = e.get("skin") || {};
            r.name && Object(Ct.p)(u, /jw-skin-\S+/, "jw-skin-" + r.name);
            var s = (function (t) {
              t || (t = {});
              var e = t.active,
                i = t.inactive,
                n = t.background,
                o = {};
              return (
                (o.controlbar = (function (t) {
                  if (t || e || i || n) {
                    var o = {};
                    return (
                      (t = t || {}),
                      (o.iconsActive = t.iconsActive || e),
                      (o.icons = t.icons || i),
                      (o.text = t.text || i),
                      (o.background = t.background || n),
                      o
                    );
                  }
                })(t.controlbar)),
                (o.timeslider = (function (t) {
                  if (t || e) {
                    var i = {};
                    return (
                      (t = t || {}),
                      (i.progress = t.progress || e),
                      (i.rail = t.rail),
                      i
                    );
                  }
                })(t.timeslider)),
                (o.menus = (function (t) {
                  if (t || e || i || n) {
                    var o = {};
                    return (
                      (t = t || {}),
                      (o.text = t.text || i),
                      (o.textActive = t.textActive || e),
                      (o.background = t.background || n),
                      o
                    );
                  }
                })(t.menus)),
                (o.tooltips = (function (t) {
                  if (t || i || n) {
                    var e = {};
                    return (
                      (t = t || {}),
                      (e.text = t.text || i),
                      (e.background = t.background || n),
                      e
                    );
                  }
                })(t.tooltips)),
                o
              );
            })(r);
            !(function (t, e) {
              var i;
              function n(e, i, n, o) {
                if (n) {
                  e = Object(h.f)(e, "#" + t + (o ? "" : " "));
                  var a = {};
                  (a[i] = n), Object(Tt.b)(e.join(", "), a, t);
                }
              }
              e &&
                (e.controlbar &&
                  (function (e) {
                    n(
                      [
                        ".jw-controlbar .jw-icon-inline.jw-text",
                        ".jw-title-primary",
                        ".jw-title-secondary",
                      ],
                      "color",
                      e.text
                    ),
                      e.icons &&
                        (n(
                          [
                            ".jw-button-color:not(.jw-icon-cast)",
                            ".jw-button-color.jw-toggle.jw-off:not(.jw-icon-cast)",
                          ],
                          "color",
                          e.icons
                        ),
                        n(
                          [".jw-display-icon-container .jw-button-color"],
                          "color",
                          e.icons
                        ),
                        Object(Tt.b)(
                          "#".concat(
                            t,
                            " .jw-icon-cast google-cast-launcher.jw-off"
                          ),
                          "{--disconnected-color: ".concat(e.icons, "}"),
                          t
                        ));
                    e.iconsActive &&
                      (n(
                        [
                          ".jw-display-icon-container .jw-button-color:hover",
                          ".jw-display-icon-container .jw-button-color:focus",
                        ],
                        "color",
                        e.iconsActive
                      ),
                      n(
                        [
                          ".jw-button-color.jw-toggle:not(.jw-icon-cast)",
                          ".jw-button-color:hover:not(.jw-icon-cast)",
                          ".jw-button-color:focus:not(.jw-icon-cast)",
                          ".jw-button-color.jw-toggle.jw-off:hover:not(.jw-icon-cast)",
                        ],
                        "color",
                        e.iconsActive
                      ),
                      n([".jw-svg-icon-buffer"], "fill", e.icons),
                      Object(Tt.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:hover google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Tt.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:focus google-cast-launcher.jw-off"
                        ),
                        "{--disconnected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Tt.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast google-cast-launcher.jw-off:focus"
                        ),
                        "{--disconnected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Tt.b)(
                        "#".concat(t, " .jw-icon-cast google-cast-launcher"),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Tt.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast google-cast-launcher:focus"
                        ),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Tt.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:hover google-cast-launcher"
                        ),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ),
                      Object(Tt.b)(
                        "#".concat(
                          t,
                          " .jw-icon-cast:focus google-cast-launcher"
                        ),
                        "{--connected-color: ".concat(e.iconsActive, "}"),
                        t
                      ));
                    n(
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
                      (n([".jw-progress", ".jw-knob"], "background-color", e),
                      n(
                        [".jw-buffer"],
                        "background-color",
                        Object(Tt.c)(e, 50)
                      ));
                    n([".jw-rail"], "background-color", t.rail),
                      n(
                        [
                          ".jw-background-color.jw-slider-time",
                          ".jw-slider-time .jw-cue",
                        ],
                        "background-color",
                        t.background
                      );
                  })(e.timeslider),
                e.menus &&
                  (n(
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
                    (i = e.menus).text
                  ),
                  n(
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
                    i.textActive
                  ),
                  n(
                    [".jw-nextup", ".jw-settings-menu"],
                    "background",
                    i.background
                  )),
                e.tooltips &&
                  (function (t) {
                    n(
                      [
                        ".jw-skip",
                        ".jw-tooltip .jw-text",
                        ".jw-time-tip .jw-text",
                      ],
                      "background-color",
                      t.background
                    ),
                      n([".jw-time-tip", ".jw-tooltip"], "color", t.background),
                      n([".jw-skip"], "border", "none"),
                      n(
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
                      var i = {
                        color: e.textActive,
                        borderColor: e.textActive,
                        stroke: e.textActive,
                      };
                      Object(Tt.b)("#".concat(t, " .jw-color-active"), i, t),
                        Object(Tt.b)(
                          "#".concat(t, " .jw-color-active-hover:hover"),
                          i,
                          t
                        );
                    }
                    if (e.text) {
                      var n = {
                        color: e.text,
                        borderColor: e.text,
                        stroke: e.text,
                      };
                      Object(Tt.b)("#".concat(t, " .jw-color-inactive"), n, t),
                        Object(Tt.b)(
                          "#".concat(t, " .jw-color-inactive-hover:hover"),
                          n,
                          t
                        );
                    }
                  })(e.menus));
            })(e.get("id"), s),
              e.set("mediaContainer", w),
              e.set("iFrame", m.Features.iframe),
              e.set("activeTab", Object(yt.a)()),
              e.set("touchMode", Wt && ("string" == typeof o || o >= gt)),
              bt.a.add(this),
              e.get("enableGradient") &&
                !Qt &&
                Object(Ct.a)(u, "jw-ab-drop-shadow"),
              (this.isSetup = !0),
              e.trigger("viewSetup", u);
            var c = document.body.contains(u);
            c && bt.a.observe(u), e.set("inDom", c);
          }),
          (this.init = function () {
            this.updateBounds(),
              e.on("change:fullscreen", Z),
              e.on("change:activeTab", z),
              e.on("change:fullscreen", z),
              e.on("change:intersectionRatio", z),
              e.on("change:visibility", U),
              e.on("instreamMode", function (t) {
                t ? dt() : pt();
              }),
              z(),
              1 !== bt.a.size() || e.get("visibility") || U(e, 1, 0);
            var t = e.player;
            e.change("state", rt),
              t.change("controls", D),
              e.change("streamType", nt),
              e.change("mediaType", ot),
              t.change("playlistItem", function (t, e) {
                lt(t, e);
              }),
              (o = a = null),
              T && Wt && bt.a.addScrollHandler(F),
              this.checkResized();
          });
        var B,
          V = 62,
          N = !0;
        function H() {
          var t = e.get("isFloating"),
            i = S.top < V,
            n = i ? S.top <= window.scrollY : S.top <= window.scrollY + V;
          !t && n ? ht(0, i) : t && !n && ht(1, i);
        }
        function F() {
          L() &&
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
        function D(t, e) {
          var i = { controls: e };
          e
            ? (Dt = Ot.a.controls)
              ? q()
              : ((i.loadPromise = Object(Ot.b)().then(function (e) {
                  Dt = e;
                  var i = t.get("controls");
                  return i && q(), i;
                })),
                i.loadPromise.catch(function (t) {
                  l.trigger(d.tb, t);
                }))
            : l.removeControls(),
            o && a && l.trigger(d.o, i);
        }
        function q() {
          var t = new Dt(document, l.element());
          l.addControls(t);
        }
        function U(t, e, i) {
          e && !i && (rt(t, t.get("state")), l.updateStyles());
        }
        function W(t) {
          I && I.mouseMove(t);
        }
        function Q(t) {
          I && !I.showing && "IFRAME" === t.target.nodeName && I.userActive();
        }
        function Y(t) {
          I &&
            I.showing &&
            ((t.relatedTarget && !u.contains(t.relatedTarget)) ||
              (!t.relatedTarget && m.Features.iframe)) &&
            I.userActive();
        }
        function X(t, e) {
          Object(Ct.p)(u, /jw-stretch-\S+/, "jw-stretch-" + e);
        }
        function K(t, i) {
          Object(Ct.v)(u, "jw-flag-aspect-mode", !!i);
          var n = u.querySelectorAll(".jw-aspect");
          Object(Tt.d)(n, { paddingTop: i || null }),
            l.isSetup &&
              i &&
              !e.get("isFloating") &&
              (Object(Tt.d)(u, G(t.get("width"))), P());
        }
        function J(i) {
          i.link
            ? (t.pause({ reason: "interaction" }),
              t.setFullscreen(!1),
              Object(Ct.l)(i.link, i.linktarget, { rel: "noreferrer" }))
            : e.get("controls") && t.playToggle({ reason: "interaction" });
        }
        (this.addControls = function (i) {
          var n = this;
          (I = i),
            Object(Ct.o)(u, "jw-flag-controls-hidden"),
            Object(Ct.v)(u, "jw-floating-dismissible", this.dismissible),
            i.enable(t, e),
            a && (R(o, a), i.resize(o, a), v.renderCues(!0)),
            i.on("userActive userInactive", function () {
              var t = e.get("state");
              (t !== d.pb && t !== d.jb) || v.renderCues(!0);
            }),
            i.on("dismissFloating", function () {
              n.stopFloating(!0), t.pause({ reason: "interaction" });
            }),
            i.on("all", l.trigger, l),
            e.get("instream") && I.setupInstream();
        }),
          (this.removeControls = function () {
            I && (I.disable(e), (I = null)),
              Object(Ct.a)(u, "jw-flag-controls-hidden"),
              Object(Ct.o)(u, "jw-floating-dismissible");
          });
        var Z = function (e, i) {
          if (
            (i && I && e.get("autostartMuted") && I.unmuteAutoplay(t, e),
            C.supportsDomFullscreen())
          )
            i ? C.requestFullscreen() : C.exitFullscreen(), it(u, i);
          else if (Qt) it(u, i);
          else {
            var n = e.get("instream"),
              o = n ? n.provider : null,
              a = e.getVideo() || o;
            a && a.setFullscreen && a.setFullscreen(i);
          }
        };
        function G(t, i, o) {
          var a = { width: t };
          if (
            (o && void 0 !== i && e.set("aspectratio", null),
            !e.get("aspectratio"))
          ) {
            var r = i;
            Object(n.r)(r) && 0 !== r && (r = Math.max(r, gt)), (a.height = r);
          }
          return a;
        }
        function $(t, i) {
          if (
            ((t && !isNaN(1 * t)) || (t = e.get("containerWidth"))) &&
            ((i && !isNaN(1 * i)) || (i = e.get("containerHeight")))
          ) {
            j && j.resize(t, i, e.get("stretching"));
            var n = e.getVideo();
            n && n.resize(t, i, e.get("stretching"));
          }
        }
        function tt(t) {
          Object(Ct.v)(u, "jw-flag-ios-fullscreen", t.jwstate), et(t);
        }
        function et(t) {
          var i = e.get("fullscreen"),
            n =
              void 0 !== t.jwstate
                ? t.jwstate
                : (function () {
                    if (C.supportsDomFullscreen()) {
                      var t = C.fullscreenElement();
                      return !(!t || t !== u);
                    }
                    return e.getVideo().getFullScreen();
                  })();
          i !== n && e.set("fullscreen", n),
            A(),
            clearTimeout(y),
            (y = setTimeout($, 200));
        }
        function it(t, e) {
          Object(Ct.v)(t, "jw-flag-fullscreen", e),
            Object(Tt.d)(document.body, { overflowY: e ? "hidden" : "" }),
            e && I && I.userActive(),
            $(),
            A();
        }
        function nt(t, e) {
          var i = "LIVE" === e;
          Object(Ct.v)(u, "jw-flag-live", i);
        }
        function ot(t, e) {
          var i = "audio" === e,
            n = t.get("provider");
          Object(Ct.v)(u, "jw-flag-media-audio", i);
          var o = n && 0 === n.name.indexOf("flash"),
            a = i && !o ? w : w.nextSibling;
          j.el.parentNode.insertBefore(j.el, a);
        }
        function at(t, e) {
          if (e) {
            var i = Object(wt.a)(t, e);
            wt.a.cloneIcon &&
              i.querySelector(".jw-icon").appendChild(wt.a.cloneIcon("error")),
              b.hide(),
              u.appendChild(i.firstChild),
              Object(Ct.v)(u, "jw-flag-audio-player", !!t.get("audioMode"));
          } else b.playlistItem(t, t.get("playlistItem"));
        }
        function rt(t, e, i) {
          if (l.isSetup) {
            if (i === d.lb) {
              var n = u.querySelector(".jw-error-msg");
              n && n.parentNode.removeChild(n);
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
              Object(Ct.i)(u, "jw-flag-controls-hidden") &&
              (Object(Ct.o)(u, "jw-flag-controls-hidden"),
              Object(Ct.v)(u, "jw-floating-dismissible", l.dismissible)),
            Object(Ct.p)(u, /jw-state-\S+/, "jw-state-" + t),
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
                (v.show(), t === d.ob && I && !I.showing && v.renderCues(!0));
          }
        }
        (this.resize = function (t, i) {
          var n = G(t, i, !0);
          void 0 !== t &&
            void 0 !== i &&
            (e.set("width", t), e.set("height", i)),
            Object(Tt.d)(u, n),
            e.get("isFloating") && vt(),
            P();
        }),
          (this.resizeMedia = $),
          (this.setPosterImage = function (t, e) {
            e.setImage(t && t.image);
          });
        var lt = function (t, e) {
            s.setPosterImage(e, j),
              Wt &&
                (function (t, e) {
                  var i = t.get("mediaElement");
                  if (i) {
                    var n = Object(Ct.j)(e.title || "");
                    i.setAttribute("title", n.textContent);
                  }
                })(t, e);
          },
          ct = function () {
            var t = I && I.settingsMenu;
            return !(!t || !t.visible);
          },
          ut = function () {
            var t = I && I.infoOverlay;
            return !(!t || !t.visible);
          },
          dt = function () {
            Object(Ct.a)(u, "jw-flag-ads"), I && I.setupInstream(), g.disable();
          },
          pt = function () {
            if (O) {
              I && I.destroyInstream(e),
                Yt !== u || Object(_t.m)() || g.enable(),
                l.setAltText(""),
                Object(Ct.o)(u, ["jw-flag-ads", "jw-flag-ads-hide-controls"]),
                e.set("hideAdsControls", !1);
              var t = e.getVideo();
              t && t.setContainer(w), O.revertAlternateClickHandlers();
            }
          };
        function ht(t, i) {
          if (t < 0.5 && !Object(_t.m)()) {
            var n = e.get("state");
            n !== d.mb &&
              n !== d.lb &&
              n !== d.kb &&
              null === Yt &&
              ((Yt = u),
              e.set("isFloating", !0),
              Object(Ct.a)(u, "jw-flag-floating"),
              i &&
                (Object(Tt.d)(p, {
                  transform: "translateY(-".concat(V - S.top, "px)"),
                }),
                setTimeout(function () {
                  Object(Tt.d)(p, {
                    transform: "translateY(0)",
                    transition:
                      "transform 150ms cubic-bezier(0, 0.25, 0.25, 1)",
                  });
                })),
              Object(Tt.d)(u, {
                backgroundImage: j.el.style.backgroundImage || e.get("image"),
              }),
              vt(),
              e.get("instreamMode") || g.enable(),
              A());
          } else l.stopFloating(!1, i);
        }
        function vt() {
          var t = e.get("width"),
            i = e.get("height"),
            o = G(t);
          if (((o.maxWidth = Math.min(400, S.width)), !e.get("aspectratio"))) {
            var a = S.width,
              r = S.height / a || 0.5625;
            Object(n.r)(t) && Object(n.r)(i) && (r = i / t),
              K(e, 100 * r + "%");
          }
          Object(Tt.d)(p, o);
        }
        (this.setAltText = function (t) {
          e.set("altText", t);
        }),
          (this.clickHandler = function () {
            return O;
          }),
          (this.getContainer = this.element = function () {
            return u;
          }),
          (this.getWrapper = function () {
            return p;
          }),
          (this.controlsContainer = function () {
            return I ? I.element() : null;
          }),
          (this.getSafeRegion = function () {
            var t =
                !(arguments.length > 0 && void 0 !== arguments[0]) ||
                arguments[0],
              e = { x: 0, y: 0, width: o || 0, height: a || 0 };
            return I && t && (e.height -= I.controlbarHeight()), e;
          }),
          (this.setCaptions = function (t) {
            v.clear(), v.setup(e.get("id"), t), v.resize();
          }),
          (this.setIntersection = function (t) {
            var i = Math.round(100 * t.intersectionRatio) / 100;
            e.set("intersectionRatio", i),
              T && !L() && (M = M || i >= 0.5) && ht(i);
          }),
          (this.stopFloating = function (t, i) {
            if ((t && ((T = null), bt.a.removeScrollHandler(F)), Yt === u)) {
              (Yt = null), e.set("isFloating", !1);
              var n = function () {
                Object(Ct.o)(u, "jw-flag-floating"),
                  K(e, e.get("aspectratio")),
                  Object(Tt.d)(u, { backgroundImage: null }),
                  Object(Tt.d)(p, {
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
              i
                ? (Object(Tt.d)(p, {
                    transform: "translateY(-".concat(V - S.top, "px)"),
                    "transition-timing-function": "ease-out",
                  }),
                  setTimeout(n, 150))
                : n(),
                g.disable(),
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
              _ && (_.destroy(), (_ = null)),
              C && (C.destroy(), (C = null)),
              I && I.disable(e),
              O &&
                (O.destroy(),
                u.removeEventListener("mousemove", W),
                u.removeEventListener("mouseout", Y),
                u.removeEventListener("mouseover", Q),
                (O = null)),
              v.destroy(),
              i && (i.destroy(), (i = null)),
              Object(Tt.a)(e.get("id")),
              this.resizeListener &&
                (this.resizeListener.destroy(), delete this.resizeListener),
              T && Wt && bt.a.removeScrollHandler(F);
          });
      };
      function Kt(t, e, i) {
        return (Kt =
          "undefined" != typeof Reflect && Reflect.get
            ? Reflect.get
            : function (t, e, i) {
                var n = (function (t, e) {
                  for (
                    ;
                    !Object.prototype.hasOwnProperty.call(t, e) &&
                    null !== (t = ee(t));

                  );
                  return t;
                })(t, e);
                if (n) {
                  var o = Object.getOwnPropertyDescriptor(n, e);
                  return o.get ? o.get.call(i) : o.value;
                }
              })(t, e, i || t);
      }
      function Jt(t) {
        return (Jt =
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
      function Gt(t, e) {
        for (var i = 0; i < e.length; i++) {
          var n = e[i];
          (n.enumerable = n.enumerable || !1),
            (n.configurable = !0),
            "value" in n && (n.writable = !0),
            Object.defineProperty(t, n.key, n);
        }
      }
      function $t(t, e, i) {
        return e && Gt(t.prototype, e), i && Gt(t, i), t;
      }
      function te(t, e) {
        return !e || ("object" !== Jt(e) && "function" != typeof e) ? oe(t) : e;
      }
      function ee(t) {
        return (ee = Object.setPrototypeOf
          ? Object.getPrototypeOf
          : function (t) {
              return t.__proto__ || Object.getPrototypeOf(t);
            })(t);
      }
      function ie(t, e) {
        if ("function" != typeof e && null !== e)
          throw new TypeError(
            "Super expression must either be null or a function"
          );
        (t.prototype = Object.create(e && e.prototype, {
          constructor: { value: t, writable: !0, configurable: !0 },
        })),
          e && ne(t, e);
      }
      function ne(t, e) {
        return (ne =
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
      function re(t, e, i) {
        Object.keys(e).forEach(function (n) {
          n in e &&
            e[n] !== i[n] &&
            t.trigger("change:".concat(n), t, e[n], i[n]);
        });
      }
      function se(t, e) {
        t && t.off(null, null, e);
      }
      var le = (function (t) {
          function e(t, i) {
            var o;
            return (
              Zt(this, e),
              ((o = te(this, ee(e).call(this)))._model = t),
              (o._mediaModel = null),
              Object(n.g)(t.attributes, {
                altText: "",
                fullscreen: !1,
                logoWidth: 0,
                scrubbing: !1,
              }),
              t.on(
                "all",
                function (e, n, a, r) {
                  n === t && (n = oe(oe(o))),
                    (i && !i(e, n, a, r)) || o.trigger(e, n, a, r);
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
            ie(e, t),
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
                    i = this._mediaModel;
                  se(i, this),
                    (this._mediaModel = t),
                    t.on(
                      "all",
                      function (i, n, o, a) {
                        n === t && (n = e), e.trigger(i, n, o, a);
                      },
                      this
                    ),
                    i && re(this, t.attributes, i.attributes);
                },
              },
            ]),
            e
          );
        })(v.a),
        ce = (function (t) {
          function e(t) {
            var i;
            return (
              Zt(this, e),
              ((i = te(
                this,
                ee(e).call(this, t, function (t) {
                  var e = i._instreamModel;
                  if (e) {
                    var n = ae.exec(t);
                    if (n) if (n[1] in e.attributes) return !1;
                  }
                  return !0;
                })
              ))._instreamModel = null),
              (i._playerViewModel = new le(i._model)),
              t.on(
                "change:instream",
                function (t, e) {
                  i.instreamModel = e ? e.model : null;
                },
                oe(oe(i))
              ),
              i
            );
          }
          return (
            ie(e, t),
            $t(e, [
              {
                key: "get",
                value: function (t) {
                  var e = this._mediaModel;
                  if (e && t in e.attributes) return e.get(t);
                  var i = this._instreamModel;
                  return i && t in i.attributes ? i.get(t) : this._model.get(t);
                },
              },
              {
                key: "getVideo",
                value: function () {
                  var t = this._instreamModel;
                  return t && t.getVideo()
                    ? t.getVideo()
                    : Kt(ee(e.prototype), "getVideo", this).call(this);
                },
              },
              {
                key: "destroy",
                value: function () {
                  Kt(ee(e.prototype), "destroy", this).call(this),
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
                    i = this._instreamModel;
                  if (
                    (se(i, this),
                    this._model.off("change:mediaModel", null, this),
                    (this._instreamModel = t),
                    this.trigger("instreamMode", !!t),
                    t)
                  )
                    t.on(
                      "all",
                      function (i, n, o, a) {
                        n === t && (n = e), e.trigger(i, n, o, a);
                      },
                      this
                    ),
                      t.change(
                        "mediaModel",
                        function (t, i) {
                          e.mediaModel = i;
                        },
                        this
                      ),
                      re(this, t.attributes, this._model.attributes);
                  else if (i) {
                    this._model.change(
                      "mediaModel",
                      function (t, i) {
                        e.mediaModel = i;
                      },
                      this
                    );
                    var o = Object(n.g)(
                      {},
                      this._model.attributes,
                      i.attributes
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
        pe = i(64),
        he =
          (ue = window).URL && ue.URL.createObjectURL
            ? ue.URL
            : ue.webkitURL || ue.mozURL;
      function fe(t, e) {
        var i = e.muted;
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
          (t.muted = i),
          (t.src = he.createObjectURL(de)),
          t.play() || Object(pe.a)(t)
        );
      }
      var we = "autoplayEnabled",
        ge = "autoplayMuted",
        je = "autoplayDisabled",
        be = {};
      var me = i(65);
      function ve(t) {
        return (
          (t = t || window.event) &&
          /^(?:mouse|pointer|touch|gesture|click|key)/.test(t.type)
        );
      }
      var ye = i(24),
        ke = "tabHidden",
        xe = "tabVisible",
        Te = function (t) {
          var e = 0;
          return function (i) {
            var n = i.position;
            n > e && t(), (e = n);
          };
        };
      function Oe(t, e) {
        e.off(d.N, t._onPlayAttempt),
          e.off(d.fb, t._triggerFirstFrame),
          e.off(d.S, t._onTime),
          t.off("change:activeTab", t._onTabVisible);
      }
      var Ce = function (t, e) {
        t.change("mediaModel", function (t, i, n) {
          t._qoeItem && n && t._qoeItem.end(n.get("mediaState")),
            (t._qoeItem = new ye.a()),
            (t._qoeItem.getFirstFrame = function () {
              var t = this.between(d.N, d.H),
                e = this.between(xe, d.H);
              return e > 0 && e < t ? e : t;
            }),
            t._qoeItem.tick(d.db),
            t._qoeItem.start(i.get("mediaState")),
            (function (t, e) {
              t._onTabVisible && Oe(t, e);
              var i = !1;
              (t._triggerFirstFrame = function () {
                if (!i) {
                  i = !0;
                  var n = t._qoeItem;
                  n.tick(d.H);
                  var o = n.getFirstFrame();
                  if ((e.trigger(d.H, { loadTime: o }), e.mediaController)) {
                    var a = e.mediaController.mediaModel;
                    a.off("change:".concat(d.U), null, a),
                      a.change(
                        d.U,
                        function (t, i) {
                          i && e.trigger(d.U, i);
                        },
                        a
                      );
                  }
                  Oe(t, e);
                }
              }),
                (t._onTime = Te(t._triggerFirstFrame)),
                (t._onPlayAttempt = function () {
                  t._qoeItem.tick(d.N);
                }),
                (t._onTabVisible = function (e, i) {
                  i ? t._qoeItem.tick(xe) : t._qoeItem.tick(ke);
                }),
                t.on("change:activeTab", t._onTabVisible),
                e.on(d.N, t._onPlayAttempt),
                e.once(d.fb, t._triggerFirstFrame),
                e.on(d.S, t._onTime);
            })(t, e),
            i.on("change:mediaState", function (e, i, n) {
              i !== n && (t._qoeItem.end(n), t._qoeItem.start(i));
            });
        });
      };
      function _e(t) {
        return (_e =
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
      var Me = function () {},
        Se = function () {};
      Object(n.g)(Me.prototype, {
        setup: function (t, e, i, h, w, b) {
          var v,
            y,
            k,
            x,
            T = this,
            O = this,
            C = (O._model = new A()),
            _ = !1,
            M = !1,
            S = null,
            E = g(H),
            I = g(Se);
          (O.originalContainer = O.currentContainer = i),
            (O._events = h),
            (O.trigger = function (t, e) {
              var i = (function (t, e, i) {
                var o = i;
                switch (e) {
                  case "time":
                  case "beforePlay":
                  case "pause":
                  case "play":
                  case "ready":
                    var a = t.get("viewable");
                    void 0 !== a && (o = Object(n.g)({}, i, { viewable: a }));
                }
                return o;
              })(C, t, e);
              return f.a.trigger.call(this, t, i);
            });
          var L = new s.a(O, ["trigger"], function () {
              return !0;
            }),
            P = function (t, e) {
              O.trigger(t, e);
            };
          C.setup(t);
          var R = C.get("backgroundLoading"),
            z = new ce(C);
          (v = this._view = new Xt(e, z)).on(
            "all",
            function (t, e) {
              (e && e.doNotForward) || P(t, e);
            },
            O
          );
          var B = (this._programController = new Y(C, b));
          ut(),
            B.on("all", P, O)
              .on(
                "subtitlesTracks",
                function (t) {
                  y.setSubtitlesTracks(t.tracks);
                  var e = y.getCurrentIndex();
                  e > 0 && rt(e, t.tracks);
                },
                O
              )
              .on(
                d.F,
                function () {
                  Promise.resolve().then(at);
                },
                O
              )
              .on(d.G, O.triggerError, O),
            Ce(C, B),
            C.on(d.w, O.triggerError, O),
            C.on(
              "change:state",
              function (t, e, i) {
                X() || K.call(T, t, e, i);
              },
              this
            ),
            C.on("change:castState", function (t, e) {
              O.trigger(d.m, e);
            }),
            C.on("change:fullscreen", function (t, e) {
              O.trigger(d.y, { fullscreen: e }),
                e && t.set("playOnViewable", !1);
            }),
            C.on("change:volume", function (t, e) {
              O.trigger(d.V, { volume: e });
            }),
            C.on("change:mute", function (t) {
              O.trigger(d.M, { mute: t.getMute() });
            }),
            C.on("change:playbackRate", function (t, e) {
              O.trigger(d.ab, { playbackRate: e, position: t.get("position") });
            });
          var V = function t(e, i) {
            ("clickthrough" !== i && "interaction" !== i && "external" !== i) ||
              (C.set("playOnViewable", !1),
              C.off("change:playReason change:pauseReason", t));
          };
          function N(t, e) {
            Object(n.t)(e) || C.set("viewable", Math.round(e));
          }
          function H() {
            dt &&
              (!0 !== C.get("autostart") ||
                C.get("playOnViewable") ||
                $("autostart"),
              dt.flush());
          }
          function F(t, e) {
            O.trigger("viewable", { viewable: e }), D();
          }
          function D() {
            if (
              (o.a[0] === e || 1 === C.get("viewable")) &&
              "idle" === C.get("state") &&
              !1 === C.get("autostart")
            )
              if (!b.primed() && m.OS.android) {
                var t = b.getTestElement(),
                  i = O.getMute();
                Promise.resolve()
                  .then(function () {
                    return fe(t, { muted: i });
                  })
                  .then(function () {
                    "idle" === C.get("state") && B.preloadVideo();
                  })
                  .catch(Se);
              } else B.preloadVideo();
          }
          function q(t) {
            (O._instreamAdapter.noResume = !t), t || et({ reason: "viewable" });
          }
          function U(t) {
            t || (O.pause({ reason: "viewable" }), C.set("playOnViewable", !t));
          }
          function W(t, e) {
            var i = X();
            if (t.get("playOnViewable")) {
              if (e) {
                var n = t.get("autoPause").pauseAds,
                  o = t.get("pauseReason");
                J() === d.mb
                  ? $("viewable")
                  : (i && !n) ||
                    "interaction" === o ||
                    Z({ reason: "viewable" });
              } else
                m.OS.mobile &&
                  !i &&
                  (O.pause({ reason: "autostart" }),
                  C.set("playOnViewable", !0));
              m.OS.mobile && i && q(e);
            }
          }
          function Q(t, e) {
            var i = t.get("state"),
              n = X(),
              o = t.get("playReason");
            n
              ? t.get("autoPause").pauseAds
                ? U(e)
                : q(e)
              : i === d.pb || i === d.jb
              ? U(e)
              : i === d.mb &&
                "playlist" === o &&
                t.once("change:state", function () {
                  U(e);
                });
          }
          function X() {
            var t = O._instreamAdapter;
            return !!t && t.getState();
          }
          function J() {
            var t = X();
            return t || C.get("state");
          }
          function Z(t) {
            if ((E.cancel(), (M = !1), C.get("state") === d.lb))
              return Promise.resolve();
            var i = G(t);
            return (
              C.set("playReason", i),
              X()
                ? (e.pauseAd(!1, t), Promise.resolve())
                : (C.get("state") === d.kb && (tt(!0), O.setItemIndex(0)),
                  !_ &&
                  ((_ = !0),
                  O.trigger(d.C, {
                    playReason: i,
                    startTime:
                      t && t.startTime
                        ? t.startTime
                        : C.get("playlistItem").starttime,
                  }),
                  (_ = !1),
                  ve() && !b.primed() && b.prime(),
                  "playlist" === i &&
                    C.get("autoPause").viewability &&
                    Q(C, C.get("viewable")),
                  x)
                    ? (ve() && !R && C.get("mediaElement").load(),
                      (x = !1),
                      (k = null),
                      Promise.resolve())
                    : B.playVideo(i).then(b.played))
            );
          }
          function G(t) {
            return t && t.reason ? t.reason : "unknown";
          }
          function $(t) {
            if (J() === d.mb) {
              E = g(H);
              var e = C.get("advertising");
              (function (t, e) {
                var i = e.cancelable,
                  n = e.muted,
                  o = void 0 !== n && n,
                  a = e.allowMuted,
                  r = void 0 !== a && a,
                  s = e.timeout,
                  l = void 0 === s ? 1e4 : s,
                  c = t.getTestElement(),
                  u = o ? "muted" : "".concat(r);
                be[u] ||
                  (be[u] = fe(c, { muted: o })
                    .catch(function (t) {
                      if (!i.cancelled() && !1 === o && r)
                        return fe(c, { muted: (o = !0) });
                      throw t;
                    })
                    .then(function () {
                      return o ? ((be[u] = null), ge) : we;
                    })
                    .catch(function (t) {
                      throw (
                        (clearTimeout(d), (be[u] = null), (t.reason = je), t)
                      );
                    }));
                var d,
                  p = be[u].then(function (t) {
                    if ((clearTimeout(d), i.cancelled())) {
                      var e = new Error("Autoplay test was cancelled");
                      throw ((e.reason = "cancelled"), e);
                    }
                    return t;
                  }),
                  h = new Promise(function (t, e) {
                    d = setTimeout(function () {
                      be[u] = null;
                      var t = new Error("Autoplay test timed out");
                      (t.reason = "timeout"), e(t);
                    }, l);
                  });
                return Promise.race([p, h]);
              })(b, {
                cancelable: E,
                muted: O.getMute(),
                allowMuted: !e || e.autoplayadsmuted,
              })
                .then(function (e) {
                  return (
                    C.set("canAutoplay", e),
                    e !== ge ||
                      O.getMute() ||
                      (C.set("autostartMuted", !0),
                      ut(),
                      C.once("change:autostartMuted", function (t) {
                        t.off("change:viewable", W),
                          O.trigger(d.M, { mute: C.getMute() });
                      })),
                    O.getMute() &&
                      C.get("enableDefaultCaptions") &&
                      y.selectDefaultIndex(1),
                    Z({ reason: t }).catch(function () {
                      O._instreamAdapter || C.set("autostartFailed", !0),
                        (k = null);
                    })
                  );
                })
                .catch(function (t) {
                  if (
                    (C.set("canAutoplay", je),
                    C.set("autostart", !1),
                    !E.cancelled())
                  ) {
                    var e = Object(j.w)(t);
                    O.trigger(d.h, { reason: t.reason, code: e, error: t });
                  }
                });
            }
          }
          function tt(t) {
            if ((E.cancel(), dt.empty(), X())) {
              var e = O._instreamAdapter;
              return (
                e && (e.noResume = !0),
                void (k = function () {
                  return B.stopVideo();
                })
              );
            }
            (k = null),
              !t && (M = !0),
              _ && (x = !0),
              C.set("errorEvent", void 0),
              B.stopVideo();
          }
          function et(t) {
            var e = G(t);
            C.set("pauseReason", e), C.set("playOnViewable", "viewable" === e);
          }
          function it(t) {
            (k = null), E.cancel();
            var i = X();
            if (i && i !== d.ob) return et(t), void e.pauseAd(!0, t);
            switch (C.get("state")) {
              case d.lb:
                return;
              case d.pb:
              case d.jb:
                et(t), B.pause();
                break;
              default:
                _ && (x = !0);
            }
          }
          function nt(t, e) {
            tt(!0), O.setItemIndex(t), O.play(e);
          }
          function ot(t) {
            nt(C.get("item") + 1, t);
          }
          function at() {
            O.completeCancelled() ||
              ((k = O.completeHandler),
              O.shouldAutoAdvance()
                ? O.nextItem()
                : C.get("repeat")
                ? ot({ reason: "repeat" })
                : (m.OS.iOS && lt(!1),
                  C.set("playOnViewable", !1),
                  C.set("state", d.kb),
                  O.trigger(d.cb, {})));
          }
          function rt(t, e) {
            (t = parseInt(t, 10) || 0),
              C.persistVideoSubtitleTrack(t, e),
              (B.subtitles = t),
              O.trigger(d.k, { tracks: st(), track: t });
          }
          function st() {
            return y.getCaptionsList();
          }
          function lt(t) {
            Object(n.n)(t) || (t = !C.get("fullscreen")),
              C.set("fullscreen", t),
              O._instreamAdapter &&
                O._instreamAdapter._adModel &&
                O._instreamAdapter._adModel.set("fullscreen", t);
          }
          function ut() {
            (B.mute = C.getMute()), (B.volume = C.get("volume"));
          }
          C.on("change:playReason change:pauseReason", V),
            O.on(d.c, function (t) {
              return V(0, t.playReason);
            }),
            O.on(d.b, function (t) {
              return V(0, t.pauseReason);
            }),
            C.on("change:scrubbing", function (t, e) {
              e
                ? ((S = C.get("state") !== d.ob), it())
                : S && Z({ reason: "interaction" });
            }),
            C.on("change:captionsList", function (t, e) {
              O.trigger(d.l, { tracks: e, track: C.get("captionsIndex") || 0 });
            }),
            C.on("change:mediaModel", function (t, e) {
              var i = this;
              t.set("errorEvent", void 0),
                e.change(
                  "mediaState",
                  function (e, i) {
                    var n;
                    t.get("errorEvent") ||
                      t.set(d.bb, (n = i) === d.nb || n === d.qb ? d.jb : n);
                  },
                  this
                ),
                e.change(
                  "duration",
                  function (e, i) {
                    if (0 !== i) {
                      var n = t.get("minDvrWindow"),
                        o = Object(me.b)(i, n);
                      t.setStreamType(o);
                    }
                  },
                  this
                );
              var n = t.get("item") + 1,
                o = "autoplay" === (t.get("related") || {}).oncomplete,
                a = t.get("playlist")[n];
              if ((a || o) && R) {
                e.on(
                  "change:position",
                  function t(n, r) {
                    var s = a && !a.daiSetting,
                      l = e.get("duration");
                    s && r && l > 0 && r >= l - p.b
                      ? (e.off("change:position", t, i), B.backgroundLoad(a))
                      : o && (a = C.get("nextUp"));
                  },
                  this
                );
              }
            }),
            (y = new ht(C)).on("all", P, O),
            z.on("viewSetup", function (t) {
              Object(a.b)(T, t);
            }),
            (this.playerReady = function () {
              v.once(d.hb, function () {
                try {
                  !(function () {
                    C.change("visibility", N),
                      L.off(),
                      O.trigger(d.gb, { setupTime: 0 }),
                      C.change("playlist", function (t, e) {
                        if (e.length) {
                          var i = { playlist: e },
                            o = C.get("feedData");
                          o && (i.feedData = Object(n.g)({}, o)),
                            O.trigger(d.eb, i);
                        }
                      }),
                      C.change("playlistItem", function (t, e) {
                        if (e) {
                          var i = e.title,
                            n = e.image;
                          if (
                            "mediaSession" in navigator &&
                            window.MediaMetadata &&
                            (i || n)
                          )
                            try {
                              navigator.mediaSession.metadata = new window.MediaMetadata(
                                {
                                  title: i,
                                  artist: window.location.hostname,
                                  artwork: [{ src: n || "" }],
                                }
                              );
                            } catch (t) {}
                          t.set("cues", []),
                            O.trigger(d.db, { index: C.get("item"), item: e });
                        }
                      }),
                      L.flush(),
                      L.destroy(),
                      (L = null),
                      C.change("viewable", F),
                      C.change("viewable", W),
                      C.get("autoPause").viewability
                        ? C.change("viewable", Q)
                        : C.once(
                            "change:autostartFailed change:mute",
                            function (t) {
                              t.off("change:viewable", W);
                            }
                          );
                    H(),
                      C.on("change:itemReady", function (t, e) {
                        e && dt.flush();
                      });
                  })();
                } catch (t) {
                  O.triggerError(Object(j.v)(j.m, j.a, t));
                }
              }),
                v.init();
            }),
            (this.preload = D),
            (this.load = function (t, e) {
              var i,
                n = O._instreamAdapter;
              switch (
                (n && (n.noResume = !0),
                O.trigger("destroyPlugin", {}),
                tt(!0),
                E.cancel(),
                (E = g(H)),
                I.cancel(),
                ve() && b.prime(),
                _e(t))
              ) {
                case "string":
                  (C.attributes.item = 0),
                    (C.attributes.itemReady = !1),
                    (I = g(function (t) {
                      if (t)
                        return O.updatePlaylist(Object(c.a)(t.playlist), t);
                    })),
                    (i = (function (t) {
                      var e = this;
                      return new Promise(function (i, n) {
                        var o = new l.a();
                        o.on(d.eb, function (t) {
                          i(t);
                        }),
                          o.on(d.w, n, e),
                          o.load(t);
                      });
                    })(t).then(I.async));
                  break;
                case "object":
                  (C.attributes.item = 0),
                    (i = O.updatePlaylist(Object(c.a)(t), e || {}));
                  break;
                case "number":
                  i = O.setItemIndex(t);
                  break;
                default:
                  return;
              }
              i.catch(function (t) {
                O.triggerError(Object(j.u)(t, j.c));
              }),
                i.then(E.async).catch(Se);
            }),
            (this.play = function (t) {
              return Z(t).catch(Se);
            }),
            (this.pause = it),
            (this.seek = function (t, e) {
              var i = C.get("state");
              if (i !== d.lb) {
                B.position = t;
                var n = i === d.mb;
                C.get("scrubbing") ||
                  (!n && i !== d.kb) ||
                  (n && ((e = e || {}).startTime = t), this.play(e));
              }
            }),
            (this.stop = tt),
            (this.playlistItem = nt),
            (this.playlistNext = ot),
            (this.playlistPrev = function (t) {
              nt(C.get("item") - 1, t);
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
            (this.getState = J),
            (this.next = Se),
            (this.completeHandler = at),
            (this.completeCancelled = function () {
              return (
                ((t = C.get("state")) !== d.mb && t !== d.kb && t !== d.lb) ||
                (!!M && ((M = !1), !0))
              );
              var t;
            }),
            (this.shouldAutoAdvance = function () {
              return C.get("item") !== C.get("playlist").length - 1;
            }),
            (this.nextItem = function () {
              ot({ reason: "playlist" });
            }),
            (this.setConfig = function (t) {
              !(function (t, e) {
                var i = t._model,
                  n = i.attributes;
                e.height &&
                  ((e.height = Object(r.b)(e.height)),
                  (e.width = e.width || n.width)),
                  e.width &&
                    ((e.width = Object(r.b)(e.width)),
                    e.aspectratio
                      ? ((n.width = e.width), delete e.width)
                      : (e.height = n.height)),
                  e.width &&
                    e.height &&
                    !e.aspectratio &&
                    t._view.resize(e.width, e.height),
                  Object.keys(e).forEach(function (o) {
                    var a = e[o];
                    if (void 0 !== a)
                      switch (o) {
                        case "aspectratio":
                          i.set(o, Object(r.a)(a, n.width));
                          break;
                        case "autostart":
                          !(function (t, e, i) {
                            t.setAutoStart(i),
                              "idle" === t.get("state") &&
                                !0 === i &&
                                e.play({ reason: "autostart" });
                          })(i, t, a);
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
                          i.set(o, a);
                      }
                  });
              })(O, t);
            }),
            (this.setItemIndex = function (t) {
              B.stopVideo();
              var e = C.get("playlist").length;
              return (
                (t = (parseInt(t, 10) || 0) % e) < 0 && (t += e),
                B.setActiveItem(t).catch(function (t) {
                  t.code >= 151 && t.code <= 162 && (t = Object(j.u)(t, j.e)),
                    T.triggerError(Object(j.v)(j.k, j.d, t));
                })
              );
            }),
            (this.detachMedia = function () {
              if (
                (_ && (x = !0),
                C.get("autoPause").viewability && Q(C, C.get("viewable")),
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
              C.setVolume(t), ut();
            }),
            (this.setMute = function (t) {
              C.setMute(t), ut();
            }),
            (this.setPlaybackRate = function (t) {
              C.setPlaybackRate(t);
            }),
            (this.getProvider = function () {
              return C.get("provider");
            }),
            (this.getWidth = function () {
              return C.get("containerWidth");
            }),
            (this.getHeight = function () {
              return C.get("containerHeight");
            }),
            (this.getItemQoe = function () {
              return C._qoeItem;
            }),
            (this.addButton = function (t, e, i, n, o) {
              var a = C.get("customButtons") || [],
                r = !1,
                s = { img: t, tooltip: e, callback: i, id: n, btnClass: o };
              (a = a.reduce(function (t, e) {
                return e.id === n ? ((r = !0), t.push(s)) : t.push(e), t;
              }, [])),
                r || a.unshift(s),
                C.set("customButtons", a);
            }),
            (this.removeButton = function (t) {
              var e = C.get("customButtons") || [];
              (e = e.filter(function (e) {
                return e.id !== t;
              })),
                C.set("customButtons", e);
            }),
            (this.resize = v.resize),
            (this.getSafeRegion = v.getSafeRegion),
            (this.setCaptions = v.setCaptions),
            (this.checkBeforePlay = function () {
              return _;
            }),
            (this.setControls = function (t) {
              Object(n.n)(t) || (t = !C.get("controls")),
                C.set("controls", t),
                (B.controls = t);
            }),
            (this.addCues = function (t) {
              this.setCues(C.get("cues").concat(t));
            }),
            (this.setCues = function (t) {
              C.set("cues", t);
            }),
            (this.updatePlaylist = function (t, e) {
              try {
                var i = Object(c.b)(t, C, e);
                Object(c.e)(i);
                var o = Object(n.g)({}, e);
                delete o.playlist, C.set("feedData", o), C.set("playlist", i);
              } catch (t) {
                return Promise.reject(t);
              }
              return this.setItemIndex(C.get("item"));
            }),
            (this.setPlaylistItem = function (t, e) {
              (e = Object(c.d)(C, new u.a(e), e.feedData || {})) &&
                ((C.get("playlist")[t] = e),
                t === C.get("item") &&
                  "idle" === C.get("state") &&
                  this.setItemIndex(t));
            }),
            (this.playerDestroy = function () {
              this.off(),
                this.stop(),
                Object(a.b)(this, this.originalContainer),
                v && v.destroy(),
                C && C.destroy(),
                dt && dt.destroy(),
                y && y.destroy(),
                B && B.destroy(),
                this.instreamDestroy();
            }),
            (this.isBeforePlay = this.checkBeforePlay),
            (this.createInstream = function () {
              return (
                this.instreamDestroy(),
                (this._instreamAdapter = new ct(this, C, v, b)),
                this._instreamAdapter
              );
            }),
            (this.instreamDestroy = function () {
              O._instreamAdapter &&
                (O._instreamAdapter.destroy(), (O._instreamAdapter = null));
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
              return !T._model.get("itemReady") || L;
            }
          );
          dt.queue.push.apply(dt.queue, w), v.setup();
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
      e.default = Me;
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
          var i = [];
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
                e.sort().filter(function (t, e, i) {
                  if ("number" != typeof t || isNaN(t) || t < 0 || t > 1)
                    throw new Error(
                      "threshold must be a number between 0 and 1 inclusively"
                    );
                  return t !== i[e - 1];
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
                i = e
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
                  u = e && l && this._computeTargetAndRootIntersection(a, i),
                  d = (o.entry = new n({
                    time: t.performance && performance.now && performance.now(),
                    target: a,
                    boundingClientRect: r,
                    rootBounds: i,
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
            (o.prototype._computeTargetAndRootIntersection = function (i, n) {
              if ("none" != t.getComputedStyle(i).display) {
                for (
                  var o, a, r, l, u, d, p, h, f = s(i), w = c(i), g = !1;
                  !g;

                ) {
                  var j = null,
                    b = 1 == w.nodeType ? t.getComputedStyle(w) : {};
                  if ("none" == b.display) return;
                  if (
                    (w == this.root || w == e
                      ? ((g = !0), (j = n))
                      : w != e.body &&
                        w != e.documentElement &&
                        "visible" != b.overflow &&
                        (j = s(w)),
                    j &&
                      ((o = j),
                      (a = f),
                      (r = void 0),
                      (l = void 0),
                      (u = void 0),
                      (d = void 0),
                      (p = void 0),
                      (h = void 0),
                      (r = Math.max(o.top, a.top)),
                      (l = Math.min(o.bottom, a.bottom)),
                      (u = Math.max(o.left, a.left)),
                      (d = Math.min(o.right, a.right)),
                      (h = l - r),
                      !(f = (p = d - u) >= 0 &&
                        h >= 0 && {
                          top: r,
                          bottom: l,
                          left: u,
                          right: d,
                          width: p,
                          height: h,
                        })))
                  )
                    break;
                  w = c(w);
                }
                return f;
              }
            }),
            (o.prototype._getRootRect = function () {
              var t;
              if (this.root) t = s(this.root);
              else {
                var i = e.documentElement,
                  n = e.body;
                t = {
                  top: 0,
                  left: 0,
                  right: i.clientWidth || n.clientWidth,
                  width: i.clientWidth || n.clientWidth,
                  bottom: i.clientHeight || n.clientHeight,
                  height: i.clientHeight || n.clientHeight,
                };
              }
              return this._expandRectByRootMargin(t);
            }),
            (o.prototype._expandRectByRootMargin = function (t) {
              var e = this._rootMarginValues.map(function (e, i) {
                  return "px" == e.unit
                    ? e.value
                    : (e.value * (i % 2 ? t.width : t.height)) / 100;
                }),
                i = {
                  top: t.top - e[0],
                  right: t.right + e[1],
                  bottom: t.bottom + e[2],
                  left: t.left - e[3],
                };
              return (
                (i.width = i.right - i.left), (i.height = i.bottom - i.top), i
              );
            }),
            (o.prototype._hasCrossedThreshold = function (t, e) {
              var i = t && t.isIntersecting ? t.intersectionRatio || 0 : -1,
                n = e.isIntersecting ? e.intersectionRatio || 0 : -1;
              if (i !== n)
                for (var o = 0; o < this.thresholds.length; o++) {
                  var a = this.thresholds[o];
                  if (a == i || a == n || a < i != a < n) return !0;
                }
            }),
            (o.prototype._rootIsInDom = function () {
              return !this.root || l(e, this.root);
            }),
            (o.prototype._rootContainsTarget = function (t) {
              return l(this.root || e, t);
            }),
            (o.prototype._registerInstance = function () {
              i.indexOf(this) < 0 && i.push(this);
            }),
            (o.prototype._unregisterInstance = function () {
              var t = i.indexOf(this);
              -1 != t && i.splice(t, 1);
            }),
            (t.IntersectionObserver = o),
            (t.IntersectionObserverEntry = n);
        }
        function n(t) {
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
            i = e.width * e.height,
            n = this.intersectionRect,
            o = n.width * n.height;
          this.intersectionRatio = i ? o / i : this.isIntersecting ? 1 : 0;
        }
        function o(t, e) {
          var i,
            n,
            o,
            a = e || {};
          if ("function" != typeof t)
            throw new Error("callback must be a function");
          if (a.root && 1 != a.root.nodeType)
            throw new Error("root must be an Element");
          (this._checkForIntersections =
            ((i = this._checkForIntersections.bind(this)),
            (n = this.THROTTLE_TIMEOUT),
            (o = null),
            function () {
              o ||
                (o = setTimeout(function () {
                  i(), (o = null);
                }, n));
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
        function a(t, e, i, n) {
          "function" == typeof t.addEventListener
            ? t.addEventListener(e, i, n || !1)
            : "function" == typeof t.attachEvent && t.attachEvent("on" + e, i);
        }
        function r(t, e, i, n) {
          "function" == typeof t.removeEventListener
            ? t.removeEventListener(e, i, n || !1)
            : "function" == typeof t.detatchEvent &&
              t.detatchEvent("on" + e, i);
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
          for (var i = e; i; ) {
            if (i == t) return !0;
            i = c(i);
          }
          return !1;
        }
        function c(t) {
          var e = t.parentNode;
          return e && 11 == e.nodeType && e.host ? e.host : e;
        }
      })(window, document);
    },
    function (t, e, i) {
      "use strict";
      i.r(e);
      var n = i(0);
      var o = i(8),
        a = i(52),
        r = i(3),
        s = i(43),
        l = {
          canplay: function () {
            this.trigger(r.E);
          },
          play: function () {
            (this.stallTime = -1),
              this.video.paused || this.state === r.pb || this.setState(r.nb);
          },
          loadedmetadata: function () {
            var t = {
                metadataType: "media",
                duration: this.getDuration(),
                height: this.video.videoHeight,
                width: this.video.videoWidth,
                seekRange: this.getSeekRange(),
              },
              e = this.drmUsed;
            e && (t.drm = e), this.trigger(r.K, t);
          },
          timeupdate: function () {
            var t = this.getVideoCurrentTime(),
              e = this.getCurrentTime(),
              i = this.getDuration();
            if (!isNaN(i)) {
              this.seeking ||
                this.video.paused ||
                (this.state !== r.qb && this.state !== r.nb) ||
                this.stallTime === t ||
                ((this.stallTime = -1),
                this.setState(r.pb),
                this.trigger(r.fb));
              var n = {
                position: e,
                duration: i,
                currentTime: t,
                seekRange: this.getSeekRange(),
                metadata: { currentTime: t },
              };
              if (this.getPtsOffset) {
                var o = this.getPtsOffset();
                o >= 0 && (n.metadata.mpegts = o + e);
              }
              var a = this.getLiveLatency();
              null !== a && (n.latency = a),
                (this.state === r.pb || this.seeking) && this.trigger(r.S, n);
            }
          },
          click: function (t) {
            this.trigger(r.n, t);
          },
          volumechange: function () {
            var t = this.video;
            this.trigger(r.V, { volume: Math.round(100 * t.volume) }),
              this.trigger(r.M, { mute: t.muted });
          },
          seeked: function () {
            this.seeking && ((this.seeking = !1), this.trigger(r.R));
          },
          playing: function () {
            -1 === this.stallTime && this.setState(r.pb), this.trigger(r.fb);
          },
          pause: function () {
            this.state !== r.kb &&
              (this.video.ended ||
                this.video.error ||
                (this.getVideoCurrentTime() !== this.getDuration() &&
                  this.setState(r.ob)));
          },
          progress: function () {
            var t = this.getDuration();
            if (!(t <= 0 || t === 1 / 0)) {
              var e = this.video.buffered;
              if (e && 0 !== e.length) {
                var i = Object(s.a)(e.end(e.length - 1) / t, 0, 1);
                this.trigger(r.D, {
                  bufferPercent: 100 * i,
                  position: this.getCurrentTime(),
                  duration: t,
                  currentTime: this.getVideoCurrentTime(),
                  seekRange: this.getSeekRange(),
                });
              }
            }
          },
          ratechange: function () {
            this.trigger(r.P, { playbackRate: this.video.playbackRate });
          },
          ended: function () {
            (this.videoHeight = 0),
              (this.streamBitrate = -1),
              this.state !== r.mb && this.state !== r.kb && this.trigger(r.F);
          },
          loadeddata: function () {
            this.renderNatively && this.setTextTracks(this.video.textTracks);
          },
        },
        c = i(10);
      function u(t) {
        return t && t.length ? t.end(t.length - 1) : 0;
      }
      var d = {
          container: null,
          volume: function (t) {
            this.video.volume = Math.min(Math.max(0, t / 100), 1);
          },
          mute: function (t) {
            (this.video.muted = !!t),
              this.video.muted || this.video.removeAttribute("muted");
          },
          resize: function (t, e, i) {
            var n = this.video,
              a = n.videoWidth,
              r = n.videoHeight;
            if (t && e && a && r) {
              var s = { objectFit: "", width: "", height: "" };
              if ("uniform" === i) {
                var l = t / e,
                  u = a / r,
                  d = Math.abs(l - u);
                d < 0.09 &&
                  d > 0.0025 &&
                  ((s.objectFit = "fill"), (i = "exactfit"));
              }
              if (
                o.Browser.ie ||
                (o.OS.iOS && o.OS.version.major < 9) ||
                o.Browser.androidNative
              )
                if ("uniform" !== i) {
                  s.objectFit = "contain";
                  var p = t / e,
                    h = a / r,
                    f = 1,
                    w = 1;
                  "none" === i
                    ? (f = w =
                        p > h
                          ? Math.ceil((100 * r) / e) / 100
                          : Math.ceil((100 * a) / t) / 100)
                    : "fill" === i
                    ? (f = w = p > h ? p / h : h / p)
                    : "exactfit" === i &&
                      (p > h ? ((f = p / h), (w = 1)) : ((f = 1), (w = h / p))),
                    Object(c.e)(
                      n,
                      "matrix("
                        .concat(f.toFixed(2), ", 0, 0, ")
                        .concat(w.toFixed(2), ", 0, 0)")
                    );
                } else (s.top = s.left = s.margin = ""), Object(c.e)(n, "");
              Object(c.d)(n, s);
            }
          },
          getContainer: function () {
            return this.container;
          },
          setContainer: function (t) {
            (this.container = t),
              this.video.parentNode !== t && t.appendChild(this.video);
          },
          remove: function () {
            this.stop(), this.destroy();
            var t = this.container;
            t && t === this.video.parentNode && t.removeChild(this.video);
          },
          atEdgeOfLiveStream: function () {
            if (!this.isLive()) return !1;
            return u(this.video.buffered) - this.video.currentTime <= 2;
          },
        },
        p = {
          eventsOn_: function () {},
          eventsOff_: function () {},
          attachMedia: function () {
            this.eventsOn_();
          },
          detachMedia: function () {
            return this.eventsOff_();
          },
        },
        h = i(65),
        f = i(5),
        w = i(53),
        g = i(7),
        j = i(66),
        b = i(63),
        m = {
          TIT2: "title",
          TT2: "title",
          WXXX: "url",
          TPE1: "artist",
          TP1: "artist",
          TALB: "album",
          TAL: "album",
        };
      function v(t, e) {
        for (var i, n, o, a = t.length, r = "", s = e || 0; s < a; )
          if (0 !== (i = t[s++]) && 3 !== i)
            switch (i >> 4) {
              case 0:
              case 1:
              case 2:
              case 3:
              case 4:
              case 5:
              case 6:
              case 7:
                r += String.fromCharCode(i);
                break;
              case 12:
              case 13:
                (n = t[s++]),
                  (r += String.fromCharCode(((31 & i) << 6) | (63 & n)));
                break;
              case 14:
                (n = t[s++]),
                  (o = t[s++]),
                  (r += String.fromCharCode(
                    ((15 & i) << 12) | ((63 & n) << 6) | ((63 & o) << 0)
                  ));
            }
        return r;
      }
      function y(t) {
        var e = (function (t) {
          for (var e = "0x", i = 0; i < t.length; i++)
            t[i] < 16 && (e += "0"), (e += t[i].toString(16));
          return parseInt(e);
        })(t);
        return (
          (127 & e) |
          ((32512 & e) >> 1) |
          ((8323072 & e) >> 2) |
          ((2130706432 & e) >> 3)
        );
      }
      function k() {
        return (arguments.length > 0 && void 0 !== arguments[0]
          ? arguments[0]
          : []
        ).reduce(function (t, e) {
          if (!("value" in e) && "data" in e && e.data instanceof ArrayBuffer) {
            var i = new Uint8Array(e.data),
              n = i.length;
            e = { value: { key: "", data: "" } };
            for (var o = 10; o < 14 && o < i.length && 0 !== i[o]; )
              (e.value.key += String.fromCharCode(i[o])), o++;
            var a = 19,
              r = i[a];
            (3 !== r && 0 !== r) || ((r = i[++a]), n--);
            var s = 0;
            if (1 !== r && 2 !== r)
              for (var l = a + 1; l < n; l++)
                if (0 === i[l]) {
                  s = l - a;
                  break;
                }
            if (s > 0) {
              var c = v(i.subarray(a, (a += s)), 0);
              if ("PRIV" === e.value.key) {
                if ("com.apple.streaming.transportStreamTimestamp" === c) {
                  var u = 1 & y(i.subarray(a, (a += 4))),
                    d = y(i.subarray(a, (a += 4))) + (u ? 4294967296 : 0);
                  e.value.data = d;
                } else e.value.data = v(i, a + 1);
                e.value.info = c;
              } else (e.value.info = c), (e.value.data = v(i, a + 1));
            } else {
              var p = i[a];
              e.value.data =
                1 === p || 2 === p
                  ? (function (t, e) {
                      for (var i = t.length - 1, n = "", o = e || 0; o < i; )
                        (254 === t[o] && 255 === t[o + 1]) ||
                          (n += String.fromCharCode((t[o] << 8) + t[o + 1])),
                          (o += 2);
                      return n;
                    })(i, a + 1)
                  : v(i, a + 1);
            }
          }
          if (
            (m.hasOwnProperty(e.value.key) &&
              (t[m[e.value.key]] = e.value.data),
            e.value.info)
          ) {
            var h = t[e.value.key];
            h !== Object(h) && ((h = {}), (t[e.value.key] = h)),
              (h[e.value.info] = e.value.data);
          } else t[e.value.key] = e.value.data;
          return t;
        }, {});
      }
      function x(t, e, i) {
        t &&
          (t.removeEventListener
            ? t.removeEventListener(e, i)
            : (t["on" + e] = null));
      }
      function T() {
        var t = this.video.textTracks,
          e = Object(n.h)(t, function (t) {
            return (t.inuse || !t._id) && S(t.kind);
          });
        if (this._textTracks && !R.call(this, e)) {
          for (var i = -1, o = 0; o < this._textTracks.length; o++)
            if ("showing" === this._textTracks[o].mode) {
              i = o;
              break;
            }
          i !== this._currentTextTrackIndex && this.setSubtitlesTrack(i + 1);
        } else this.setTextTracks(t);
      }
      function O() {
        this.setTextTracks(this.video.textTracks);
      }
      function C(t) {
        var e = this;
        t &&
          (this._textTracks || this._initTextTracks(),
          t.forEach(function (t) {
            if (!t.kind || S(t.kind)) {
              var i = E.call(e, t);
              I.call(e, i),
                t.file &&
                  ((t.data = []),
                  Object(j.c)(
                    t,
                    function (t) {
                      e.addVTTCuesToTrack(i, t);
                    },
                    function (t) {
                      e.trigger(r.tb, t);
                    }
                  ));
            }
          }),
          this._textTracks &&
            this._textTracks.length &&
            this.trigger("subtitlesTracks", { tracks: this._textTracks }));
      }
      function _(t, e, i) {
        if (o.Browser.ie) {
          var n = i;
          (t || "metadata" === e.kind) &&
            (n = new window.TextTrackCue(i.startTime, i.endTime, i.text)),
            (function (t, e) {
              var i = [],
                n = t.mode;
              t.mode = "hidden";
              for (
                var o = t.cues, a = o.length - 1;
                a >= 0 && o[a].startTime > e.startTime;
                a--
              )
                i.unshift(o[a]), t.removeCue(o[a]);
              try {
                t.addCue(e),
                  i.forEach(function (e) {
                    return t.addCue(e);
                  });
              } catch (t) {
                console.error(t);
              }
              t.mode = n;
            })(e, n);
        } else
          try {
            e.addCue(i);
          } catch (t) {
            console.error(t);
          }
      }
      function M(t, e) {
        e &&
          e.length &&
          Object(n.f)(e, function (e) {
            if (!(o.Browser.ie && t && /^(native|subtitle|cc)/.test(e._id))) {
              (o.Browser.ie && "disabled" === e.mode) ||
                ((e.mode = "disabled"), (e.mode = "hidden"));
              for (var i = e.cues.length; i--; ) e.removeCue(e.cues[i]);
              e.embedded || (e.mode = "disabled"), (e.inuse = !1);
            }
          });
      }
      function S(t) {
        return "subtitles" === t || "captions" === t;
      }
      function E(t) {
        var e,
          i = Object(b.b)(t, this._unknownCount),
          o = i.label;
        if (
          ((this._unknownCount = i.unknownCount),
          this.renderNatively || "metadata" === t.kind)
        ) {
          var a = this.video.textTracks;
          (e = Object(n.j)(a, { label: o })) ||
            (e = this.video.addTextTrack(t.kind, o, t.language || "")),
            (e.default = t.default),
            (e.mode = "disabled"),
            (e.inuse = !0);
        } else (e = t).data = e.data || [];
        return e._id || (e._id = Object(b.a)(t, this._textTracks.length)), e;
      }
      function I(t) {
        this._textTracks.push(t), (this._tracksById[t._id] = t);
      }
      function L() {
        if (this._textTracks) {
          var t = this._textTracks.filter(function (t) {
            return t.embedded || "subs" === t.groupid;
          });
          this._initTextTracks(),
            t.forEach(function (t) {
              this._tracksById[t._id] = t;
            }),
            (this._textTracks = t);
        }
      }
      function A(t) {
        this.triggerActiveCues(t.currentTarget.activeCues);
      }
      function P(t, e, i) {
        var n = t.kind;
        this._cachedVTTCues[t._id] || (this._cachedVTTCues[t._id] = {});
        var o,
          a = this._cachedVTTCues[t._id];
        switch (n) {
          case "captions":
          case "subtitles":
            o = i || Math.floor(20 * e.startTime);
            var r = "_" + e.line,
              s = Math.floor(20 * e.endTime),
              l = a[o + r] || a[o + 1 + r] || a[o - 1 + r];
            return !(l && Math.abs(l - s) <= 1) && ((a[o + r] = s), !0);
          case "metadata":
            var c = e.data ? new Uint8Array(e.data).join("") : e.text;
            return !a[(o = i || e.startTime + c)] && ((a[o] = e.endTime), !0);
          default:
            return !1;
        }
      }
      function R(t) {
        if (t.length > this._textTracks.length) return !0;
        for (var e = 0; e < t.length; e++) {
          var i = t[e];
          if (!i._id || !this._tracksById[i._id]) return !0;
        }
        return !1;
      }
      var z = {
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
          addTracksListener: function (t, e, i) {
            if (!t) return;
            if ((x(t, e, i), this.instreamMode)) return;
            t.addEventListener ? t.addEventListener(e, i) : (t["on" + e] = i);
          },
          clearTracks: function () {
            Object(j.a)(this._itemTracks);
            var t = this._tracksById && this._tracksById.nativemetadata;
            (this.renderNatively || t) &&
              (M(this.renderNatively, this.video.textTracks),
              t && (t.oncuechange = null));
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
                M(this.renderNatively, this.video.textTracks));
          },
          clearMetaCues: function () {
            var t = this._tracksById && this._tracksById.nativemetadata;
            t &&
              (M(this.renderNatively, [t]),
              (t.mode = "hidden"),
              (t.inuse = !0),
              (this._cachedVTTCues[t._id] = {}));
          },
          clearCueData: function (t) {
            var e = this._cachedVTTCues;
            e &&
              e[t] &&
              ((e[t] = {}),
              this._tracksById && (this._tracksById[t].data = []));
          },
          disableTextTrack: function () {
            if (this._textTracks) {
              var t = this._textTracks[this._currentTextTrackIndex];
              if (t) {
                t.mode = "disabled";
                var e = t._id;
                e && 0 === e.indexOf("nativecaptions") && (t.mode = "hidden");
              }
            }
          },
          enableTextTrack: function () {
            if (this._textTracks) {
              var t = this._textTracks[this._currentTextTrackIndex];
              t && (t.mode = "showing");
            }
          },
          getSubtitlesTrack: function () {
            return this._currentTextTrackIndex;
          },
          removeTracksListener: x,
          addTextTracks: C,
          setTextTracks: function (t) {
            if (((this._currentTextTrackIndex = -1), !t)) return;
            this._textTracks
              ? ((this._unknownCount = 0),
                (this._textTracks = this._textTracks.filter(function (t) {
                  var e = t._id;
                  return this.renderNatively &&
                    e &&
                    0 === e.indexOf("nativecaptions")
                    ? (delete this._tracksById[e], !1)
                    : (t.name &&
                        0 === t.name.indexOf("Unknown") &&
                        this._unknownCount++,
                      !0);
                }, this)),
                delete this._tracksById.nativemetadata)
              : this._initTextTracks();
            if (t.length)
              for (var e = 0, i = t.length; e < i; e++) {
                var n = t[e];
                if (!n._id) {
                  if ("captions" === n.kind || "metadata" === n.kind) {
                    if (
                      ((n._id = "native" + n.kind + e),
                      !n.label && "captions" === n.kind)
                    ) {
                      var a = Object(b.b)(n, this._unknownCount);
                      (n.name = a.label), (this._unknownCount = a.unknownCount);
                    }
                  } else n._id = Object(b.a)(n, this._textTracks.length);
                  if (this._tracksById[n._id]) continue;
                  n.inuse = !0;
                }
                if (n.inuse && !this._tracksById[n._id])
                  if ("metadata" === n.kind)
                    (n.mode = "hidden"),
                      (n.oncuechange = A.bind(this)),
                      (this._tracksById[n._id] = n);
                  else if (S(n.kind)) {
                    var r = n.mode,
                      s = void 0;
                    if (((n.mode = "hidden"), !n.cues.length && n.embedded))
                      continue;
                    if (
                      ((n.mode = r),
                      this._cuesByTrackId[n._id] &&
                        !this._cuesByTrackId[n._id].loaded)
                    ) {
                      for (
                        var l = this._cuesByTrackId[n._id].cues;
                        (s = l.shift());

                      )
                        _(this.renderNatively, n, s);
                      (n.mode = r), (this._cuesByTrackId[n._id].loaded = !0);
                    }
                    I.call(this, n);
                  }
              }
            this.renderNatively &&
              ((this.textTrackChangeHandler =
                this.textTrackChangeHandler || T.bind(this)),
              this.addTracksListener(
                this.video.textTracks,
                "change",
                this.textTrackChangeHandler
              ),
              (o.Browser.edge || o.Browser.firefox || o.Browser.safari) &&
                ((this.addTrackHandler = this.addTrackHandler || O.bind(this)),
                this.addTracksListener(
                  this.video.textTracks,
                  "addtrack",
                  this.addTrackHandler
                )));
            this._textTracks.length &&
              this.trigger("subtitlesTracks", { tracks: this._textTracks });
          },
          setupSideloadedTracks: function (t) {
            if (!this.renderNatively) return;
            var e = t === this._itemTracks;
            e || Object(j.a)(this._itemTracks);
            if (((this._itemTracks = t), !t)) return;
            e || (this.disableTextTrack(), L.call(this), this.addTextTracks(t));
          },
          setSubtitlesTrack: function (t) {
            if (!this.renderNatively)
              return void (
                this.setCurrentSubtitleTrack &&
                this.setCurrentSubtitleTrack(t - 1)
              );
            if (!this._textTracks) return;
            0 === t &&
              this._textTracks.forEach(function (t) {
                t.mode = t.embedded ? "hidden" : "disabled";
              });
            if (this._currentTextTrackIndex === t - 1) return;
            this.disableTextTrack(),
              (this._currentTextTrackIndex = t - 1),
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
          addCuesToTrack: function (t) {
            var e = this._tracksById[t.name];
            if (!e) return;
            e.source = t.source;
            for (
              var i = t.captions || [], n = [], o = !1, a = 0;
              a < i.length;
              a++
            ) {
              var r = i[a],
                s = t.name + "_" + r.begin + "_" + r.end;
              this._metaCuesByTextTime[s] ||
                ((this._metaCuesByTextTime[s] = r), n.push(r), (o = !0));
            }
            o &&
              n.sort(function (t, e) {
                return t.begin - e.begin;
              });
            var l = Object(j.b)(n);
            Array.prototype.push.apply(e.data, l);
          },
          addCaptionsCue: function (t) {
            if (!t.text || !t.begin || !t.end) return;
            var e,
              i = t.trackid.toString(),
              n = this._tracksById && this._tracksById[i];
            n ||
              ((n = { kind: "captions", _id: i, data: [] }),
              this.addTextTracks([n]),
              this.trigger("subtitlesTracks", { tracks: this._textTracks }));
            t.useDTS && (n.source || (n.source = t.source || "mpegts"));
            e = t.begin + "_" + t.text;
            var o = this._metaCuesByTextTime[e];
            if (!o) {
              (o = { begin: t.begin, end: t.end, text: t.text }),
                (this._metaCuesByTextTime[e] = o);
              var a = Object(j.b)([o])[0];
              n.data.push(a);
            }
          },
          createCue: function (t, e, i) {
            var n = window.VTTCue || window.TextTrackCue,
              o = Math.max(e || 0, t + 0.25);
            return new n(t, o, i);
          },
          addVTTCue: function (t, e) {
            this._tracksById || this._initTextTracks();
            var i = t.track ? t.track : "native" + t.type,
              n = this._tracksById[i],
              o = "captions" === t.type ? "Unknown CC" : "ID3 Metadata",
              a = t.cue;
            if (!n) {
              var r = { kind: t.type, _id: i, label: o, embedded: !0 };
              (n = E.call(this, r)),
                this.renderNatively || "metadata" === n.kind
                  ? this.setTextTracks(this.video.textTracks)
                  : C.call(this, [n]);
            }
            if (P.call(this, n, a, e)) {
              var s = this.renderNatively || "metadata" === n.kind;
              return s ? _(s, n, a) : n.data.push(a), a;
            }
            return null;
          },
          addVTTCuesToTrack: function (t, e) {
            if (!this.renderNatively) return;
            var i,
              n = this._tracksById[t._id];
            if (!n)
              return (
                this._cuesByTrackId || (this._cuesByTrackId = {}),
                void (this._cuesByTrackId[t._id] = { cues: e, loaded: !1 })
              );
            if (this._cuesByTrackId[t._id] && this._cuesByTrackId[t._id].loaded)
              return;
            this._cuesByTrackId[t._id] = { cues: e, loaded: !0 };
            for (; (i = e.shift()); ) _(this.renderNatively, n, i);
          },
          triggerActiveCues: function (t) {
            var e = this;
            if (!t || !t.length) return void (this._activeCues = null);
            var i = this._activeCues || [],
              n = Array.prototype.filter.call(t, function (t) {
                if (
                  i.some(function (e) {
                    return (
                      (n = e),
                      (i = t).startTime === n.startTime &&
                        i.endTime === n.endTime &&
                        i.text === n.text &&
                        i.data === n.data &&
                        i.value === n.value
                    );
                    var i, n;
                  })
                )
                  return !1;
                if (t.data || t.value) return !0;
                if (t.text) {
                  var n = JSON.parse(t.text),
                    o = { metadataTime: t.startTime, metadata: n };
                  n.programDateTime && (o.programDateTime = n.programDateTime),
                    n.metadataType &&
                      ((o.metadataType = n.metadataType),
                      delete n.metadataType),
                    e.trigger(r.K, o);
                }
                return !1;
              });
            if (n.length) {
              var o = k(n),
                a = n[0].startTime;
              this.trigger(r.K, {
                metadataType: "id3",
                metadataTime: a,
                metadata: o,
              });
            }
            this._activeCues = Array.prototype.slice.call(t);
          },
          renderNatively: !1,
        },
        B = i(64),
        V = i(15),
        N = i(1),
        H = 224e3,
        F = 224005,
        D = 221e3,
        q = 324e3,
        U = window.clearTimeout,
        W = "html5",
        Q = function () {};
      function Y(t, e) {
        Object.keys(t).forEach(function (i) {
          e.removeEventListener(i, t[i]);
        });
      }
      function X(t, e, i) {
        (this.state = r.mb),
          (this.seeking = !1),
          (this.currentTime = -1),
          (this.retries = 0),
          (this.maxRetries = 3);
        var s,
          w = this,
          j = e.minDvrWindow,
          b = {
            progress: function () {
              l.progress.call(w), ft();
            },
            timeupdate: function () {
              w.currentTime >= 0 && (w.retries = 0);
              var t = w.getVideoCurrentTime();
              (w.currentTime = t),
                M && C !== t && $(t),
                l.timeupdate.call(w),
                ft(),
                o.Browser.ie && G();
            },
            resize: G,
            ended: function () {
              (_ = -1), wt(), l.ended.call(w);
            },
            loadedmetadata: function () {
              var t = w.getDuration();
              R && t === 1 / 0 && (t = 0);
              var e = {
                metadataType: "media",
                duration: t,
                height: v.videoHeight,
                width: v.videoWidth,
                seekRange: w.getSeekRange(),
              };
              w.trigger(r.K, e), G();
            },
            durationchange: function () {
              R || l.progress.call(w);
            },
            loadeddata: function () {
              var t;
              !(function () {
                if (v.getStartDate) {
                  var t = v.getStartDate(),
                    e = t.getTime ? t.getTime() : NaN;
                  if (e !== w.startDateTime && !isNaN(e)) {
                    w.startDateTime = e;
                    var i = t.toISOString(),
                      n = w.getSeekRange(),
                      o = n.start,
                      a = n.end,
                      s = {
                        metadataType: "program-date-time",
                        programDateTime: i,
                        start: o,
                        end: a,
                      },
                      l = w.createCue(o, a, JSON.stringify(s));
                    w.addVTTCue({ type: "metadata", cue: l }),
                      delete s.metadataType,
                      w.trigger(r.L, {
                        metadataType: "program-date-time",
                        metadata: s,
                      });
                  }
                }
              })(),
                l.loadeddata.call(w),
                (function (t) {
                  if (((E = null), !t)) return;
                  if (t.length) {
                    for (var e = 0; e < t.length; e++)
                      if (t[e].enabled) {
                        I = e;
                        break;
                      }
                    -1 === I && (t[(I = 0)].enabled = !0),
                      (E = Object(n.v)(t, function (t) {
                        return {
                          name: t.label || t.language,
                          language: t.language,
                        };
                      }));
                  }
                  w.addTracksListener(t, "change", ct),
                    E &&
                      w.trigger("audioTracks", { currentTrack: I, tracks: E });
                })(v.audioTracks),
                (t = w.getDuration()),
                T && -1 !== T && t && t !== 1 / 0 && w.seek(T),
                G();
            },
            canplay: function () {
              (x = !0),
                R || ht(),
                o.Browser.ie &&
                  9 === o.Browser.version.major &&
                  w.setTextTracks(w._textTracks),
                l.canplay.call(w);
            },
            seeking: function () {
              var t = null !== O ? tt(O) : w.getCurrentTime(),
                e = tt(C);
              (C = O),
                (O = null),
                (T = 0),
                (w.seeking = !0),
                w.trigger(r.Q, { position: e, offset: t });
            },
            seeked: function () {
              l.seeked.call(w);
            },
            waiting: function () {
              w.seeking
                ? w.setState(r.nb)
                : w.state === r.pb &&
                  (w.atEdgeOfLiveStream() && w.setPlaybackRate(1),
                  (w.stallTime = w.video.currentTime),
                  w.setState(r.qb));
            },
            webkitbeginfullscreen: function (t) {
              (M = !0), ut(t);
            },
            webkitendfullscreen: function (t) {
              (M = !1), ut(t);
            },
            error: function () {
              var t = w.video,
                e = t.error,
                i = (e && e.code) || -1;
              if ((3 === i || 4 === i) && w.retries < w.maxRetries)
                return (
                  w.trigger(r.tb, new N.n(null, q + i - 1, e)),
                  w.retries++,
                  v.load(),
                  void (
                    -1 !== w.currentTime &&
                    ((x = !1), w.seek(w.currentTime), (w.currentTime = -1))
                  )
                );
              var n = H,
                o = N.k;
              1 === i
                ? (n += i)
                : 2 === i
                ? ((o = N.i), (n = D))
                : 3 === i || 4 === i
                ? ((n += i - 1), 4 === i && t.src === location.href && (n = F))
                : (o = N.m),
                rt(),
                w.trigger(r.G, new N.n(o, n, e));
            },
          };
        Object.keys(l).forEach(function (t) {
          if (!b[t]) {
            var e = l[t];
            b[t] = function (t) {
              e.call(w, t);
            };
          }
        }),
          Object(n.g)(this, g.a, d, p, z, {
            renderNatively:
              ((s = e.renderCaptionsNatively),
              !(!o.OS.iOS && !o.Browser.safari) || (s && o.Browser.chrome)),
            eventsOn_: function () {
              var t, e;
              (t = b),
                (e = v),
                Object.keys(t).forEach(function (i) {
                  e.removeEventListener(i, t[i]), e.addEventListener(i, t[i]);
                });
            },
            eventsOff_: function () {
              Y(b, v);
            },
            detachMedia: function () {
              p.detachMedia.call(w),
                wt(),
                this.removeTracksListener(
                  v.textTracks,
                  "change",
                  this.textTrackChangeHandler
                ),
                this.disableTextTrack();
            },
            attachMedia: function () {
              p.attachMedia.call(w),
                (x = !1),
                (this.seeking = !1),
                (v.loop = !1),
                this.enableTextTrack(),
                this.renderNatively &&
                  this.setTextTracks(this.video.textTracks),
                this.addTracksListener(
                  v.textTracks,
                  "change",
                  this.textTrackChangeHandler
                );
            },
            isLive: function () {
              return this.getDuration() === 1 / 0;
            },
          });
        var m,
          v = i,
          y = { level: {} },
          k = null !== e.liveTimeout ? e.liveTimeout : 3e4,
          x = !1,
          T = 0,
          O = null,
          C = null,
          _ = -1,
          M = !1,
          S = Q,
          E = null,
          I = -1,
          L = -1,
          A = !1,
          P = null,
          R = !1,
          X = null,
          J = null,
          Z = 0;
        function G() {
          var t = y.level;
          if (t.width !== v.videoWidth || t.height !== v.videoHeight) {
            if ((!v.videoWidth && !pt()) || -1 === _) return;
            (t.width = v.videoWidth),
              (t.height = v.videoHeight),
              ht(),
              (y.reason = y.reason || "auto"),
              (y.mode = "hls" === m[_].type ? "auto" : "manual"),
              (y.bitrate = 0),
              (t.index = _),
              (t.label = m[_].label),
              w.trigger(r.U, y),
              (y.reason = "");
          }
        }
        function $(t) {
          C = t;
        }
        function tt(t) {
          var e = w.getSeekRange();
          return w.isLive() && Object(h.a)(e.end - e.start, j)
            ? Math.min(0, t - e.end)
            : t;
        }
        function et(t) {
          var e;
          return (
            Array.isArray(t) &&
              t.length > 0 &&
              (e = t.map(function (t, e) {
                return { label: t.label || e };
              })),
            e
          );
        }
        function it(t) {
          (w.currentTime = -1),
            (j = t.minDvrWindow),
            (m = t.sources),
            (_ = (function (t) {
              var i = Math.max(0, _),
                n = e.qualityLabel;
              if (t)
                for (var o = 0; o < t.length; o++)
                  if ((t[o].default && (i = o), n && t[o].label === n))
                    return o;
              (y.reason = "initial choice"),
                (y.level.width && y.level.height) || (y.level = {});
              return i;
            })(m));
        }
        function nt() {
          return (
            v.paused &&
              v.played &&
              v.played.length &&
              w.isLive() &&
              !Object(h.a)(lt() - st(), j) &&
              (w.clearTracks(), v.load()),
            v.play() || Object(B.a)(v)
          );
        }
        function ot(t) {
          (w.currentTime = -1), (T = 0), wt();
          var e = v.src,
            i = document.createElement("source");
          (i.src = m[_].file),
            i.src !== e
              ? (at(m[_]), e && v.load())
              : 0 === t && w.getVideoCurrentTime() > 0 && ((T = -1), w.seek(t)),
            t > 0 && w.getVideoCurrentTime() !== t && w.seek(t);
          var n = et(m);
          n && w.trigger(r.I, { levels: n, currentQuality: _ }),
            m.length && "hls" !== m[0].type && ht();
        }
        function at(t) {
          (E = null),
            (I = -1),
            y.reason || ((y.reason = "initial choice"), (y.level = {})),
            (x = !1);
          var e = document.createElement("source");
          (e.src = t.file), v.src !== e.src && (v.src = t.file);
        }
        function rt() {
          v &&
            (w.disableTextTrack(),
            v.removeAttribute("preload"),
            v.removeAttribute("src"),
            Object(f.h)(v),
            Object(c.d)(v, { objectFit: "" }),
            (_ = -1),
            !o.Browser.msie && "load" in v && v.load());
        }
        function st() {
          var t = 1 / 0;
          return (
            ["buffered", "seekable"].forEach(function (e) {
              for (var i = v[e], o = i ? i.length : 0; o--; ) {
                var a = Math.min(t, i.start(o));
                Object(n.o)(a) && (t = a);
              }
            }),
            t
          );
        }
        function lt() {
          var t = 0;
          return (
            ["buffered", "seekable"].forEach(function (e) {
              for (var i = v[e], o = i ? i.length : 0; o--; ) {
                var a = Math.max(t, i.end(o));
                Object(n.o)(a) && (t = a);
              }
            }),
            t
          );
        }
        function ct() {
          for (var t = -1, e = 0; e < v.audioTracks.length; e++)
            if (v.audioTracks[e].enabled) {
              t = e;
              break;
            }
          dt(t);
        }
        function ut(t) {
          w.trigger(r.X, { target: t.target, jwstate: M });
        }
        function dt(t) {
          v &&
            v.audioTracks &&
            E &&
            t > -1 &&
            t < v.audioTracks.length &&
            t !== I &&
            ((v.audioTracks[I].enabled = !1),
            (I = t),
            (v.audioTracks[I].enabled = !0),
            w.trigger("audioTrackChanged", { currentTrack: I, tracks: E }));
        }
        function pt() {
          if (!(v.readyState < 2)) return 0 === v.videoHeight;
        }
        function ht() {
          var t = pt();
          if (void 0 !== t) {
            var e = t ? "audio" : "video";
            w.trigger(r.T, { mediaType: e });
          }
        }
        function ft() {
          if (0 !== k) {
            var t = u(v.buffered);
            w.isLive() && t && P === t
              ? -1 === L &&
                (L = setTimeout(function () {
                  (A = !0),
                    (function () {
                      if (A && w.atEdgeOfLiveStream())
                        return w.trigger(r.G, new N.n(N.l, K)), !0;
                    })();
                }, k))
              : (wt(), (A = !1)),
              (P = t);
          }
        }
        function wt() {
          U(L), (L = -1);
        }
        (this.video = v),
          (this.supportsPlaybackRate = !0),
          (this.startDateTime = 0),
          (w.getVideoCurrentTime = function () {
            return e.getCurrentTimeHook
              ? e.getCurrentTimeHook(v)
              : v.currentTime;
          }),
          (w.getCurrentTime = function () {
            return (function (t) {
              var e = w.getSeekRange();
              if (w.isLive()) {
                if (
                  ((!J || Math.abs(X - e.end) > 1) &&
                    (function (t) {
                      (X = t.end),
                        (J = Math.min(0, w.getVideoCurrentTime() - X)),
                        (Z = Object(V.a)());
                    })(e),
                  Object(h.a)(e.end - e.start, j))
                )
                  return J;
              }
              return t;
            })(w.getVideoCurrentTime());
          }),
          (w.getDuration = function () {
            if (e.getDurationHook) return e.getDurationHook();
            var t = v.duration;
            if ((R && t === 1 / 0 && 0 === w.getVideoCurrentTime()) || isNaN(t))
              return 0;
            var i = lt();
            if (v.duration === 1 / 0 && i) {
              var n = i - st();
              Object(h.a)(n, j) && (t = -n);
            }
            return t;
          }),
          (w.getSeekRange = function () {
            var t = { start: 0, end: w.getDuration() };
            return v.seekable.length && ((t.end = lt()), (t.start = st())), t;
          }),
          (w.getLiveLatency = function () {
            var t = null,
              e = lt();
            return (
              w.isLive() &&
                e &&
                (t = e + (Object(V.a)() - Z) / 1e3 - w.getVideoCurrentTime()),
              t
            );
          }),
          (this.stop = function () {
            wt(),
              rt(),
              this.clearTracks(),
              o.Browser.ie && v.pause(),
              this.setState(r.mb);
          }),
          (this.destroy = function () {
            (S = Q),
              Y(b, v),
              this.removeTracksListener(v.audioTracks, "change", ct),
              this.removeTracksListener(
                v.textTracks,
                "change",
                w.textTrackChangeHandler
              ),
              this.off();
          }),
          (this.init = function (t) {
            (w.retries = 0), (w.maxRetries = t.adType ? 0 : 3), it(t);
            var e = m[_];
            (R = Object(a.a)(e)) &&
              ((w.supportsPlaybackRate = !1), (b.waiting = Q)),
              w.eventsOn_(),
              m.length && "hls" !== m[0].type && this.sendMediaType(m),
              (y.reason = "");
          }),
          (this.preload = function (t) {
            it(t);
            var e = m[_],
              i = e.preload || "metadata";
            "none" !== i && (v.setAttribute("preload", i), at(e));
          }),
          (this.load = function (t) {
            it(t), ot(t.starttime), this.setupSideloadedTracks(t.tracks);
          }),
          (this.play = function () {
            return S(), nt();
          }),
          (this.pause = function () {
            wt(),
              (S = function () {
                if (v.paused && w.getVideoCurrentTime() && w.isLive()) {
                  var t = lt(),
                    e = t - st(),
                    i = !Object(h.a)(e, j),
                    o = t - w.getVideoCurrentTime();
                  if (i && t && (o > 15 || o < 0)) {
                    if (((O = Math.max(t - 10, t - e)), !Object(n.o)(O)))
                      return;
                    $(w.getVideoCurrentTime()), (v.currentTime = O);
                  }
                }
              }),
              v.pause();
          }),
          (this.seek = function (t) {
            if (!e.seekHook || !e.seekHook(t, v)) {
              var i = w.getSeekRange(),
                n = t;
              if ((t < 0 && (n += i.end), x || (x = !!lt()), x)) {
                T = 0;
                try {
                  if (
                    ((w.seeking = !0),
                    w.isLive() && Object(h.a)(i.end - i.start, j))
                  )
                    if (((J = Math.min(0, n - X)), t < 0))
                      n += Math.min(12, (Object(V.a)() - Z) / 1e3);
                  (O = n), $(w.getVideoCurrentTime()), (v.currentTime = n);
                } catch (t) {
                  (w.seeking = !1), (T = n);
                }
              } else (T = n), o.Browser.firefox && v.paused && nt();
            }
          }),
          (this.setVisibility = function (t) {
            (t = !!t) || o.OS.android
              ? Object(c.d)(w.container, { visibility: "visible", opacity: 1 })
              : Object(c.d)(w.container, { visibility: "", opacity: 0 });
          }),
          (this.setFullscreen = function (t) {
            if ((t = !!t)) {
              try {
                var e = v.webkitEnterFullscreen || v.webkitEnterFullScreen;
                e && e.apply(v);
              } catch (t) {
                return !1;
              }
              return w.getFullScreen();
            }
            var i = v.webkitExitFullscreen || v.webkitExitFullScreen;
            return i && i.apply(v), t;
          }),
          (w.getFullScreen = function () {
            return M || !!v.webkitDisplayingFullscreen;
          }),
          (this.setCurrentQuality = function (t) {
            _ !== t &&
              t >= 0 &&
              m &&
              m.length > t &&
              ((_ = t),
              (y.reason = "api"),
              (y.level = {}),
              this.trigger(r.J, { currentQuality: t, levels: et(m) }),
              (e.qualityLabel = m[t].label),
              ot(w.getVideoCurrentTime() || 0),
              nt());
          }),
          (this.setPlaybackRate = function (t) {
            v.playbackRate = v.defaultPlaybackRate = t;
          }),
          (this.getPlaybackRate = function () {
            return v.playbackRate;
          }),
          (this.getCurrentQuality = function () {
            return _;
          }),
          (this.getQualityLevels = function () {
            return Array.isArray(m)
              ? m.map(function (t) {
                  return (function (t) {
                    return {
                      bitrate: t.bitrate,
                      label: t.label,
                      width: t.width,
                      height: t.height,
                    };
                  })(t);
                })
              : [];
          }),
          (this.getName = function () {
            return { name: W };
          }),
          (this.setCurrentAudioTrack = dt),
          (this.getAudioTracks = function () {
            return E || [];
          }),
          (this.getCurrentAudioTrack = function () {
            return I;
          });
      }
      Object(n.g)(X.prototype, w.a),
        (X.getName = function () {
          return { name: "html5" };
        });
      e.default = X;
      var K = 220001;
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
    function (t, e, i) {
      "use strict";
      i.d(e, "a", function () {
        return o;
      });
      var n = i(2);
      function o(t) {
        var e = [],
          i = (t = Object(n.i)(t)).split("\r\n\r\n");
        1 === i.length && (i = t.split("\n\n"));
        for (var o = 0; o < i.length; o++)
          if ("WEBVTT" !== i[o]) {
            var r = a(i[o]);
            r.text && e.push(r);
          }
        return e;
      }
      function a(t) {
        var e = {},
          i = t.split("\r\n");
        1 === i.length && (i = t.split("\n"));
        var o = 1;
        if (
          (i[0].indexOf(" --\x3e ") > 0 && (o = 0),
          i.length > o + 1 && i[o + 1])
        ) {
          var a = i[o],
            r = a.indexOf(" --\x3e ");
          r > 0 &&
            ((e.begin = Object(n.g)(a.substr(0, r))),
            (e.end = Object(n.g)(a.substr(r + 5))),
            (e.text = i.slice(o + 1).join("\r\n")));
        }
        return e;
      }
    },
    function (t, e, i) {
      "use strict";
      i.d(e, "a", function () {
        return o;
      }),
        i.d(e, "b", function () {
          return a;
        });
      var n = i(5);
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
        var i = "jw-breakpoint-" + e;
        Object(n.p)(t, /jw-breakpoint--?\d+/, i);
      }
    },
    function (t, e, i) {
      "use strict";
      i.d(e, "a", function () {
        return d;
      });
      var n,
        o = i(0),
        a = i(8),
        r = i(16),
        s = i(7),
        l = i(3),
        c = i(10),
        u = i(5),
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
            h,
            f,
            w,
            g,
            j,
            b,
            m = this,
            v = t.player;
          function y() {
            Object(o.o)(e.fontSize) &&
              (v.get("containerHeight")
                ? (j =
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
                var i = t * j;
                e =
                  Math.round(
                    10 *
                      (function (t) {
                        var e = v.get("mediaElement");
                        if (e && e.videoHeight) {
                          var i = e.videoWidth,
                            n = e.videoHeight,
                            o = i / n,
                            r = v.get("containerHeight"),
                            s = v.get("containerWidth");
                          if (v.get("fullscreen") && a.OS.mobile) {
                            var l = window.screen;
                            l.orientation &&
                              ((r = l.availHeight), (s = l.availWidth));
                          }
                          if (s && r && i && n)
                            return (s / r > o ? r : (n * s) / i) * j;
                        }
                        return t;
                      })(i)
                  ) / 10;
              }
              v.get("renderCaptionsNatively")
                ? (function (t, e) {
                    var i = "#".concat(
                      t,
                      " .jw-video::-webkit-media-text-track-display"
                    );
                    e &&
                      ((e += "px"),
                      a.OS.iOS &&
                        Object(c.b)(i, { fontSize: "inherit" }, t, !0));
                    (b.fontSize = e), Object(c.b)(i, b, t, !0);
                  })(v.get("id"), e)
                : Object(c.d)(f, { fontSize: e });
            }
          }
          function x(t, e, i) {
            var n = Object(c.c)("#000000", i);
            "dropshadow" === t
              ? (e.textShadow = "0 2px 1px " + n)
              : "raised" === t
              ? (e.textShadow =
                  "0 0 5px " + n + ", 0 1px 5px " + n + ", 0 2px 5px " + n)
              : "depressed" === t
              ? (e.textShadow = "0 -2px 1px " + n)
              : "uniform" === t &&
                (e.textShadow =
                  "-2px 0 1px " +
                  n +
                  ",2px 0 1px " +
                  n +
                  ",0 -2px 1px " +
                  n +
                  ",0 2px 1px " +
                  n +
                  ",-1px 1px 1px " +
                  n +
                  ",1px 1px 1px " +
                  n +
                  ",1px -1px 1px " +
                  n +
                  ",1px 1px 1px " +
                  n);
          }
          ((f = document.createElement("div")).className =
            "jw-captions jw-reset"),
            (this.show = function () {
              Object(u.a)(f, "jw-captions-enabled");
            }),
            (this.hide = function () {
              Object(u.o)(f, "jw-captions-enabled");
            }),
            (this.populate = function (t) {
              v.get("renderCaptionsNatively") ||
                ((p = []),
                (s = t),
                t ? this.selectCues(t, h) : this.renderCues());
            }),
            (this.resize = function () {
              k(), this.renderCues(!0);
            }),
            (this.renderCues = function (t) {
              (t = !!t), n && n.processCues(window, p, f, t);
            }),
            (this.selectCues = function (t, e) {
              if (t && t.data && e && !v.get("renderCaptionsNatively")) {
                var i = this.getAlignmentPosition(t, e);
                !1 !== i &&
                  ((p = this.getCurrentCues(t.data, i)), this.renderCues(!0));
              }
            }),
            (this.getCurrentCues = function (t, e) {
              return Object(o.h)(t, function (t) {
                return e >= t.startTime && (!t.endTime || e <= t.endTime);
              });
            }),
            (this.getAlignmentPosition = function (t, e) {
              var i = t.source,
                n = e.metadata,
                a = e.currentTime;
              return i && n && Object(o.r)(n[i]) && (a = n[i]), a;
            }),
            (this.clear = function () {
              Object(u.g)(f);
            }),
            (this.setup = function (t, i) {
              (w = document.createElement("div")),
                (g = document.createElement("span")),
                (w.className = "jw-captions-window jw-reset"),
                (g.className = "jw-captions-text jw-reset"),
                (e = Object(o.g)({}, d, i)),
                (j = d.fontScale);
              var n = function () {
                if (!v.get("renderCaptionsNatively")) {
                  y(e.fontSize);
                  var i = e.windowColor,
                    n = e.windowOpacity,
                    o = e.edgeStyle;
                  b = {};
                  var r = {};
                  !(function (t, e) {
                    var i = e.color,
                      n = e.fontOpacity;
                    (i || n !== d.fontOpacity) &&
                      (t.color = Object(c.c)(i || "#ffffff", n));
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
                    (i || n !== d.windowOpacity) &&
                      (b.backgroundColor = Object(c.c)(i || "#000000", n)),
                    x(o, r, e.fontOpacity),
                    e.back || null !== o || x("uniform", r),
                    Object(c.d)(w, b),
                    Object(c.d)(g, r),
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
              n(),
                w.appendChild(g),
                f.appendChild(w),
                v.change(
                  "captionsTrack",
                  function (t, e) {
                    this.populate(e);
                  },
                  this
                ),
                v.set("captions", e),
                v.on("change:captions", function (t, i) {
                  (e = i), n();
                });
            }),
            (this.element = function () {
              return f;
            }),
            (this.destroy = function () {
              v.off(null, null, this), this.off();
            });
          var T = function (t) {
            (h = t), m.selectCues(s, h);
          };
          v.on(
            "change:playlistItem",
            function () {
              (h = null), (p = []);
            },
            this
          ),
            v.on(
              l.Q,
              function (t) {
                (p = []), T(t);
              },
              this
            ),
            v.on(l.S, T, this),
            v.on(
              "subtitlesTrackData",
              function () {
                this.selectCues(s, h);
              },
              this
            ),
            v.on(
              "change:captionsList",
              function t(e, o) {
                var a = this;
                1 !== o.length &&
                  (e.get("renderCaptionsNatively") ||
                    n ||
                    (i
                      .e(8)
                      .then(
                        function (t) {
                          n = i(68).default;
                        }.bind(null, i)
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
    function (t, e, i) {
      "use strict";
      t.exports = function (t) {
        var e = [];
        return (
          (e.toString = function () {
            return this.map(function (e) {
              var i = (function (t, e) {
                var i = t[1] || "",
                  n = t[3];
                if (!n) return i;
                if (e && "function" == typeof btoa) {
                  var o =
                      ((r = n),
                      "/*# sourceMappingURL=data:application/json;charset=utf-8;base64," +
                        btoa(unescape(encodeURIComponent(JSON.stringify(r)))) +
                        " */"),
                    a = n.sources.map(function (t) {
                      return "/*# sourceURL=" + n.sourceRoot + t + " */";
                    });
                  return [i].concat(a).concat([o]).join("\n");
                }
                var r;
                return [i].join("\n");
              })(e, t);
              return e[2] ? "@media " + e[2] + "{" + i + "}" : i;
            }).join("");
          }),
          (e.i = function (t, i) {
            "string" == typeof t && (t = [[null, t, ""]]);
            for (var n = {}, o = 0; o < this.length; o++) {
              var a = this[o][0];
              null != a && (n[a] = !0);
            }
            for (o = 0; o < t.length; o++) {
              var r = t[o];
              (null != r[0] && n[r[0]]) ||
                (i && !r[2]
                  ? (r[2] = i)
                  : i && (r[2] = "(" + r[2] + ") and (" + i + ")"),
                e.push(r));
            }
          }),
          e
        );
      };
    },
    function (t, e) {
      var i,
        n,
        o = {},
        a = {},
        r =
          ((i = function () {
            return document.head || document.getElementsByTagName("head")[0];
          }),
          function () {
            return void 0 === n && (n = i.apply(this, arguments)), n;
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
        var i,
          n,
          o,
          r = a[t];
        r || (r = a[t] = { element: s(t), counter: 0 });
        var l = r.counter++;
        return (
          (i = r.element),
          (o = function () {
            d(i, l, "");
          }),
          (n = function (t) {
            d(i, l, t);
          })(e.css),
          function (t) {
            if (t) {
              if (t.css === e.css && t.media === e.media) return;
              n((e = t).css);
            } else o();
          }
        );
      }
      t.exports = {
        style: function (t, e) {
          !(function (t, e) {
            for (var i = 0; i < e.length; i++) {
              var n = e[i],
                a = (o[t] || {})[n.id];
              if (a) {
                for (var r = 0; r < a.parts.length; r++) a.parts[r](n.parts[r]);
                for (; r < n.parts.length; r++) a.parts.push(l(t, n.parts[r]));
              } else {
                var s = [];
                for (r = 0; r < n.parts.length; r++) s.push(l(t, n.parts[r]));
                (o[t] = o[t] || {}), (o[t][n.id] = { id: n.id, parts: s });
              }
            }
          })(
            e,
            (function (t) {
              for (var e = [], i = {}, n = 0; n < t.length; n++) {
                var o = t[n],
                  a = o[0],
                  r = o[1],
                  s = o[2],
                  l = { css: r, media: s };
                i[a]
                  ? i[a].parts.push(l)
                  : e.push((i[a] = { id: a, parts: [l] }));
              }
              return e;
            })(t)
          );
        },
        clear: function (t, e) {
          var i = o[t];
          if (!i) return;
          if (e) {
            var n = i[e];
            if (n) for (var a = 0; a < n.parts.length; a += 1) n.parts[a]();
            return;
          }
          for (var r = Object.keys(i), s = 0; s < r.length; s += 1)
            for (var l = i[r[s]], c = 0; c < l.parts.length; c += 1)
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
      function d(t, e, i) {
        if (t.styleSheet) t.styleSheet.cssText = u(e, i);
        else {
          var n = document.createTextNode(i),
            o = t.childNodes[e];
          o ? t.replaceChild(n, o) : t.appendChild(n);
        }
      }
    },
    function (t, e) {
      t.exports =
        '<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-arrow-right" viewBox="0 0 240 240" focusable="false"><path d="M183.6,104.4L81.8,0L45.4,36.3l84.9,84.9l-84.9,84.9L79.3,240l101.9-101.7c9.9-6.9,12.4-20.4,5.5-30.4C185.8,106.7,184.8,105.4,183.6,104.4L183.6,104.4z"></path></svg>';
    },
    function (t, e, i) {
      "use strict";
      function n(t, e) {
        var i = t.kind || "cc";
        return t.default || t.defaulttrack
          ? "default"
          : t._id || t.file || i + e;
      }
      function o(t, e) {
        var i = t.label || t.name || t.language;
        return (
          i || ((i = "Unknown CC"), (e += 1) > 1 && (i += " [" + e + "]")),
          { label: i, unknownCount: e }
        );
      }
      i.d(e, "a", function () {
        return n;
      }),
        i.d(e, "b", function () {
          return o;
        });
    },
    function (t, e, i) {
      "use strict";
      function n(t) {
        return new Promise(function (e, i) {
          if (t.paused) return i(o("NotAllowedError", 0, "play() failed."));
          var n = function () {
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
              if ((n(), "playing" === t.type)) e();
              else {
                var a = 'The play() request was interrupted by a "'.concat(
                  t.type,
                  '" event.'
                );
                "error" === t.type
                  ? i(o("NotSupportedError", 9, a))
                  : i(o("AbortError", 20, a));
              }
            };
          t.addEventListener("play", a);
        });
      }
      function o(t, e, i) {
        var n = new Error(i);
        return (n.name = t), (n.code = e), n;
      }
      i.d(e, "a", function () {
        return n;
      });
    },
    function (t, e, i) {
      "use strict";
      function n(t, e) {
        return t !== 1 / 0 && Math.abs(t) >= Math.max(a(e), 0);
      }
      function o(t, e) {
        var i = "VOD";
        return (
          t === 1 / 0
            ? (i = "LIVE")
            : t < 0 && (i = n(t, a(e)) ? "DVR" : "LIVE"),
          i
        );
      }
      function a(t) {
        return void 0 === t ? 120 : Math.max(t, 0);
      }
      i.d(e, "a", function () {
        return n;
      }),
        i.d(e, "b", function () {
          return o;
        });
    },
    function (t, e, i) {
      "use strict";
      var n = i(67),
        o = i(16),
        a = i(22),
        r = i(4),
        s = i(57),
        l = i(2),
        c = i(1);
      function u(t) {
        throw new c.n(null, t);
      }
      function d(t, e, n) {
        t.xhr = Object(a.a)(
          t.file,
          function (a) {
            !(function (t, e, n, a) {
              var d,
                p,
                f = t.responseXML ? t.responseXML.firstChild : null;
              if (f)
                for (
                  "xml" === Object(r.b)(f) && (f = f.nextSibling);
                  f.nodeType === f.COMMENT_NODE;

                )
                  f = f.nextSibling;
              try {
                if (f && "tt" === Object(r.b)(f))
                  (d = (function (t) {
                    t || u(306007);
                    var e = [],
                      i = t.getElementsByTagName("p"),
                      n = 30,
                      o = t.getElementsByTagName("tt");
                    if (o && o[0]) {
                      var a = parseFloat(o[0].getAttribute("ttp:frameRate"));
                      isNaN(a) || (n = a);
                    }
                    i || u(306005),
                      i.length ||
                        (i = t.getElementsByTagName("tt:p")).length ||
                        (i = t.getElementsByTagName("tts:p"));
                    for (var r = 0; r < i.length; r++) {
                      for (
                        var s = i[r], c = s.getElementsByTagName("br"), d = 0;
                        d < c.length;
                        d++
                      ) {
                        var p = c[d];
                        p.parentNode.replaceChild(t.createTextNode("\r\n"), p);
                      }
                      var h = s.innerHTML || s.textContent || s.text || "",
                        f = Object(l.i)(h)
                          .replace(/>\s+</g, "><")
                          .replace(/(<\/?)tts?:/g, "$1")
                          .replace(/<br.*?\/>/g, "\r\n");
                      if (f) {
                        var w = s.getAttribute("begin"),
                          g = s.getAttribute("dur"),
                          j = s.getAttribute("end"),
                          b = { begin: Object(l.g)(w, n), text: f };
                        j
                          ? (b.end = Object(l.g)(j, n))
                          : g && (b.end = b.begin + Object(l.g)(g, n)),
                          e.push(b);
                      }
                    }
                    return e.length || u(306005), e;
                  })(t.responseXML)),
                    (p = h(d)),
                    delete e.xhr,
                    n(p);
                else {
                  var w = t.responseText;
                  w.indexOf("WEBVTT") >= 0
                    ? i
                        .e(10)
                        .then(
                          function (t) {
                            return i(97).default;
                          }.bind(null, i)
                        )
                        .catch(Object(o.c)(301131))
                        .then(function (t) {
                          var i = new t(window);
                          (p = []),
                            (i.oncue = function (t) {
                              p.push(t);
                            }),
                            (i.onflush = function () {
                              delete e.xhr, n(p);
                            }),
                            i.parse(w);
                        })
                        .catch(function (t) {
                          delete e.xhr, a(Object(c.v)(null, c.b, t));
                        })
                    : ((d = Object(s.a)(w)), (p = h(d)), delete e.xhr, n(p));
                }
              } catch (t) {
                delete e.xhr, a(Object(c.v)(null, c.b, t));
              }
            })(a, t, e, n);
          },
          function (t, e, i, o) {
            n(Object(c.u)(o, c.b));
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
      function h(t) {
        return t.map(function (t) {
          return new n.a(t.begin, t.end, t.text);
        });
      }
      i.d(e, "c", function () {
        return d;
      }),
        i.d(e, "a", function () {
          return p;
        }),
        i.d(e, "b", function () {
          return h;
        });
    },
    function (t, e, i) {
      "use strict";
      var n = window.VTTCue;
      function o(t) {
        if ("string" != typeof t) return !1;
        return (
          !!{ start: !0, middle: !0, end: !0, left: !0, right: !0 }[
            t.toLowerCase()
          ] && t.toLowerCase()
        );
      }
      if (!n) {
        (n = function (t, e, i) {
          var n = this;
          n.hasBeenReset = !1;
          var a = "",
            r = !1,
            s = t,
            l = e,
            c = i,
            u = null,
            d = "",
            p = !0,
            h = "auto",
            f = "start",
            w = "auto",
            g = 100,
            j = "middle";
          Object.defineProperty(n, "id", {
            enumerable: !0,
            get: function () {
              return a;
            },
            set: function (t) {
              a = "" + t;
            },
          }),
            Object.defineProperty(n, "pauseOnExit", {
              enumerable: !0,
              get: function () {
                return r;
              },
              set: function (t) {
                r = !!t;
              },
            }),
            Object.defineProperty(n, "startTime", {
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
            Object.defineProperty(n, "endTime", {
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
            Object.defineProperty(n, "text", {
              enumerable: !0,
              get: function () {
                return c;
              },
              set: function (t) {
                (c = "" + t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(n, "region", {
              enumerable: !0,
              get: function () {
                return u;
              },
              set: function (t) {
                (u = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(n, "vertical", {
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
            Object.defineProperty(n, "snapToLines", {
              enumerable: !0,
              get: function () {
                return p;
              },
              set: function (t) {
                (p = !!t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(n, "line", {
              enumerable: !0,
              get: function () {
                return h;
              },
              set: function (t) {
                if ("number" != typeof t && "auto" !== t)
                  throw new SyntaxError(
                    "An invalid number or illegal string was specified."
                  );
                (h = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(n, "lineAlign", {
              enumerable: !0,
              get: function () {
                return f;
              },
              set: function (t) {
                var e = o(t);
                if (!e)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (f = e), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(n, "position", {
              enumerable: !0,
              get: function () {
                return w;
              },
              set: function (t) {
                if (t < 0 || t > 100)
                  throw new Error("Position must be between 0 and 100.");
                (w = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(n, "size", {
              enumerable: !0,
              get: function () {
                return g;
              },
              set: function (t) {
                if (t < 0 || t > 100)
                  throw new Error("Size must be between 0 and 100.");
                (g = t), (this.hasBeenReset = !0);
              },
            }),
            Object.defineProperty(n, "align", {
              enumerable: !0,
              get: function () {
                return j;
              },
              set: function (t) {
                var e = o(t);
                if (!e)
                  throw new SyntaxError(
                    "An invalid or illegal string was specified."
                  );
                (j = e), (this.hasBeenReset = !0);
              },
            }),
            (n.displayState = void 0);
        }).prototype.getCueAsHTML = function () {
          return window.WebVTT.convertCueToDOMTree(window, this.text);
        };
      }
      e.a = n;
    },
    ,
    function (t, e, i) {
      var n = i(70);
      "string" == typeof n && (n = [["all-players", n, ""]]),
        i(61).style(n, "all-players"),
        n.locals && (t.exports = n.locals);
    },
    function (t, e, i) {
      (t.exports = i(60)(!1)).push([
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
    function (t, e, i) {
      var n = i(96);
      "string" == typeof n && (n = [["all-players", n, ""]]),
        i(61).style(n, "all-players"),
        n.locals && (t.exports = n.locals);
    },
    function (t, e, i) {
      (t.exports = i(60)(!1)).push([
        t.i,
        '.jw-overlays,.jw-controls,.jw-controls-backdrop,.jw-flag-small-player .jw-settings-menu,.jw-settings-submenu{height:100%;width:100%}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-settings-menu .jw-icon.jw-button-color::after{position:absolute;right:0}.jw-overlays,.jw-controls,.jw-controls-backdrop,.jw-settings-item-active::before{top:0;position:absolute;left:0}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-settings-menu .jw-icon.jw-button-color::after{position:absolute;bottom:0;left:0}.jw-nextup-close{position:absolute;top:0;right:0}.jw-overlays,.jw-controls,.jw-flag-small-player .jw-settings-menu{position:absolute;bottom:0;right:0}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after,.jw-time-tip::after,.jw-settings-menu .jw-icon.jw-button-color::after,.jw-text-live::before,.jw-controlbar .jw-tooltip::after,.jw-settings-menu .jw-tooltip::after{content:"";display:block}.jw-svg-icon{height:24px;width:24px;fill:currentColor;pointer-events:none}.jw-icon{height:44px;width:44px;background-color:transparent;outline:none}.jw-icon.jw-tab-focus:focus{border:solid 2px #4d90fe}.jw-icon-airplay .jw-svg-icon-airplay-off{display:none}.jw-off.jw-icon-airplay .jw-svg-icon-airplay-off{display:block}.jw-icon-airplay .jw-svg-icon-airplay-on{display:block}.jw-off.jw-icon-airplay .jw-svg-icon-airplay-on{display:none}.jw-icon-cc .jw-svg-icon-cc-off{display:none}.jw-off.jw-icon-cc .jw-svg-icon-cc-off{display:block}.jw-icon-cc .jw-svg-icon-cc-on{display:block}.jw-off.jw-icon-cc .jw-svg-icon-cc-on{display:none}.jw-icon-fullscreen .jw-svg-icon-fullscreen-off{display:none}.jw-off.jw-icon-fullscreen .jw-svg-icon-fullscreen-off{display:block}.jw-icon-fullscreen .jw-svg-icon-fullscreen-on{display:block}.jw-off.jw-icon-fullscreen .jw-svg-icon-fullscreen-on{display:none}.jw-icon-volume .jw-svg-icon-volume-0{display:none}.jw-off.jw-icon-volume .jw-svg-icon-volume-0{display:block}.jw-icon-volume .jw-svg-icon-volume-100{display:none}.jw-full.jw-icon-volume .jw-svg-icon-volume-100{display:block}.jw-icon-volume .jw-svg-icon-volume-50{display:block}.jw-off.jw-icon-volume .jw-svg-icon-volume-50,.jw-full.jw-icon-volume .jw-svg-icon-volume-50{display:none}.jw-settings-menu .jw-icon::after,.jw-icon-settings::after,.jw-icon-volume::after{height:100%;width:24px;box-shadow:inset 0 -3px 0 -1px currentColor;margin:auto;opacity:0;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-settings-menu .jw-icon[aria-checked="true"]::after,.jw-settings-open .jw-icon-settings::after,.jw-icon-volume.jw-open::after{opacity:1}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-cc,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-settings,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-audio-tracks,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-hd,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-settings-sharing,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-fullscreen,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player).jw-flag-cast-available .jw-icon-airplay,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player).jw-flag-cast-available .jw-icon-cast{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-volume,.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-text-live{bottom:6px}.jwplayer.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-icon-volume::after{display:none}.jw-overlays,.jw-controls{pointer-events:none}.jw-controls-backdrop{display:block;background:linear-gradient(to bottom, transparent, rgba(0,0,0,0.4) 77%, rgba(0,0,0,0.4) 100%) 100% 100% / 100% 240px no-repeat transparent;transition:opacity 250ms cubic-bezier(0, .25, .25, 1),background-size 250ms cubic-bezier(0, .25, .25, 1);pointer-events:none}.jw-overlays{cursor:auto}.jw-controls{overflow:hidden}.jw-flag-small-player .jw-controls{text-align:center}.jw-text{height:1em;font-family:Arial,Helvetica,sans-serif;font-size:.75em;font-style:normal;font-weight:normal;color:#fff;text-align:center;font-variant:normal;font-stretch:normal}.jw-controlbar,.jw-skip,.jw-display-icon-container .jw-icon,.jw-nextup-container,.jw-autostart-mute,.jw-overlays .jw-plugin{pointer-events:all}.jwplayer .jw-display-icon-container,.jw-error .jw-display-icon-container{width:auto;height:auto;box-sizing:content-box}.jw-display{display:table;height:100%;padding:57px 0;position:relative;width:100%}.jw-flag-dragging .jw-display{display:none}.jw-state-idle:not(.jw-flag-cast-available) .jw-display{padding:0}.jw-display-container{display:table-cell;height:100%;text-align:center;vertical-align:middle}.jw-display-controls{display:inline-block}.jwplayer .jw-display-icon-container{float:left}.jw-display-icon-container{display:inline-block;padding:5.5px;margin:0 22px}.jw-display-icon-container .jw-icon{height:75px;width:75px;cursor:pointer;display:flex;justify-content:center;align-items:center}.jw-display-icon-container .jw-icon .jw-svg-icon{height:33px;width:33px;padding:0;position:relative}.jw-display-icon-container .jw-icon .jw-svg-icon-rewind{padding:.2em .05em}.jw-breakpoint--1 .jw-nextup-container{display:none}.jw-breakpoint-0 .jw-display-icon-next,.jw-breakpoint--1 .jw-display-icon-next,.jw-breakpoint-0 .jw-display-icon-rewind,.jw-breakpoint--1 .jw-display-icon-rewind{display:none}.jw-breakpoint-0 .jw-display .jw-icon,.jw-breakpoint--1 .jw-display .jw-icon,.jw-breakpoint-0 .jw-display .jw-svg-icon,.jw-breakpoint--1 .jw-display .jw-svg-icon{width:44px;height:44px;line-height:44px}.jw-breakpoint-0 .jw-display .jw-icon:before,.jw-breakpoint--1 .jw-display .jw-icon:before,.jw-breakpoint-0 .jw-display .jw-svg-icon:before,.jw-breakpoint--1 .jw-display .jw-svg-icon:before{width:22px;height:22px}.jw-breakpoint-1 .jw-display .jw-icon,.jw-breakpoint-1 .jw-display .jw-svg-icon{width:44px;height:44px;line-height:44px}.jw-breakpoint-1 .jw-display .jw-icon:before,.jw-breakpoint-1 .jw-display .jw-svg-icon:before{width:22px;height:22px}.jw-breakpoint-1 .jw-display .jw-icon.jw-icon-rewind:before{width:33px;height:33px}.jw-breakpoint-2 .jw-display .jw-icon,.jw-breakpoint-3 .jw-display .jw-icon,.jw-breakpoint-2 .jw-display .jw-svg-icon,.jw-breakpoint-3 .jw-display .jw-svg-icon{width:77px;height:77px;line-height:77px}.jw-breakpoint-2 .jw-display .jw-icon:before,.jw-breakpoint-3 .jw-display .jw-icon:before,.jw-breakpoint-2 .jw-display .jw-svg-icon:before,.jw-breakpoint-3 .jw-display .jw-svg-icon:before{width:38.5px;height:38.5px}.jw-breakpoint-4 .jw-display .jw-icon,.jw-breakpoint-5 .jw-display .jw-icon,.jw-breakpoint-6 .jw-display .jw-icon,.jw-breakpoint-7 .jw-display .jw-icon,.jw-breakpoint-4 .jw-display .jw-svg-icon,.jw-breakpoint-5 .jw-display .jw-svg-icon,.jw-breakpoint-6 .jw-display .jw-svg-icon,.jw-breakpoint-7 .jw-display .jw-svg-icon{width:88px;height:88px;line-height:88px}.jw-breakpoint-4 .jw-display .jw-icon:before,.jw-breakpoint-5 .jw-display .jw-icon:before,.jw-breakpoint-6 .jw-display .jw-icon:before,.jw-breakpoint-7 .jw-display .jw-icon:before,.jw-breakpoint-4 .jw-display .jw-svg-icon:before,.jw-breakpoint-5 .jw-display .jw-svg-icon:before,.jw-breakpoint-6 .jw-display .jw-svg-icon:before,.jw-breakpoint-7 .jw-display .jw-svg-icon:before{width:44px;height:44px}.jw-controlbar{display:flex;flex-flow:row wrap;align-items:center;justify-content:center;position:absolute;left:0;bottom:0;width:100%;border:none;border-radius:0;background-size:auto;box-shadow:none;max-height:72px;transition:250ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility;transition-delay:0s}.jw-breakpoint-7 .jw-controlbar{max-height:140px}.jw-breakpoint-7 .jw-controlbar .jw-button-container{padding:0 48px 20px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-tooltip{margin-bottom:-7px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume .jw-overlay{padding-bottom:40%}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-text{font-size:1em}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-text.jw-text-elapsed{justify-content:flex-end}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-inline,.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume{height:60px;width:60px}.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-inline .jw-svg-icon,.jw-breakpoint-7 .jw-controlbar .jw-button-container .jw-icon-volume .jw-svg-icon{height:30px;width:30px}.jw-breakpoint-7 .jw-controlbar .jw-slider-time{padding:0 60px;height:34px}.jw-breakpoint-7 .jw-controlbar .jw-slider-time .jw-slider-container{height:10px}.jw-controlbar .jw-button-image{background:no-repeat 50% 50%;background-size:contain;max-height:24px}.jw-controlbar .jw-spacer{flex:1 1 auto;align-self:stretch}.jw-controlbar .jw-icon.jw-button-color:hover{color:#fff}.jw-button-container{display:flex;flex-flow:row nowrap;flex:1 1 auto;align-items:center;justify-content:center;width:100%;padding:0 12px}.jw-slider-horizontal{background-color:transparent}.jw-icon-inline{position:relative}.jw-icon-inline,.jw-icon-tooltip{height:44px;width:44px;align-items:center;display:flex;justify-content:center}.jw-icon-inline:not(.jw-text),.jw-icon-tooltip,.jw-slider-horizontal{cursor:pointer}.jw-text-elapsed,.jw-text-duration{justify-content:flex-start;width:-webkit-fit-content;width:-moz-fit-content;width:fit-content}.jw-icon-tooltip{position:relative}.jw-knob:hover,.jw-icon-inline:hover,.jw-icon-tooltip:hover,.jw-icon-display:hover,.jw-option:before:hover{color:#fff}.jw-time-tip,.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{pointer-events:none}.jw-icon-cast{display:none;margin:0;padding:0}.jw-icon-cast google-cast-launcher{background-color:transparent;border:none;padding:0;width:24px;height:24px;cursor:pointer}.jw-icon-inline.jw-icon-volume{display:none}.jwplayer .jw-text-countdown{display:none}.jw-flag-small-player .jw-display{padding-top:0;padding-bottom:0}.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-rewind,.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-next,.jw-flag-small-player:not(.jw-flag-audio-player):not(.jw-flag-ads) .jw-controlbar .jw-button-container>.jw-icon-playback{display:none}.jw-flag-ads-vpaid:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar,.jw-flag-user-inactive.jw-state-playing:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar,.jw-flag-user-inactive.jw-state-buffering:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controlbar{visibility:hidden;pointer-events:none;opacity:0;transition-delay:0s, 250ms}.jw-flag-ads-vpaid:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop,.jw-flag-user-inactive.jw-state-playing:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop,.jw-flag-user-inactive.jw-state-buffering:not(.jw-flag-media-audio):not(.jw-flag-audio-player):not(.jw-flag-ads-vpaid-controls):not(.jw-flag-casting) .jw-controls-backdrop{opacity:0}.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-countdown{display:flex}.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint--1 .jw-text-elapsed,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-elapsed,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint--1 .jw-text-duration,.jwplayer:not(.jw-flag-ads):not(.jw-flag-live).jw-breakpoint-0 .jw-text-duration{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-text-countdown,.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-related-btn,.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-slider-volume{display:none}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-controlbar{flex-direction:column-reverse}.jwplayer.jw-breakpoint--1:not(.jw-flag-ads):not(.jw-flag-audio-player) .jw-button-container{height:30px}.jw-breakpoint--1.jw-flag-ads:not(.jw-flag-audio-player) .jw-icon-volume,.jw-breakpoint--1.jw-flag-ads:not(.jw-flag-audio-player) .jw-icon-fullscreen{display:none}.jwplayer:not(.jw-breakpoint-0) .jw-text-duration:before,.jwplayer:not(.jw-breakpoint--1) .jw-text-duration:before{content:"/";padding-right:1ch;padding-left:1ch}.jwplayer:not(.jw-flag-user-inactive) .jw-controlbar{will-change:transform}.jwplayer:not(.jw-flag-user-inactive) .jw-controlbar .jw-text{-webkit-transform-style:preserve-3d;transform-style:preserve-3d}.jw-slider-container{display:flex;align-items:center;position:relative;touch-action:none}.jw-rail,.jw-buffer,.jw-progress{position:absolute;cursor:pointer}.jw-progress{background-color:#f2f2f2}.jw-rail{background-color:rgba(255,255,255,0.3)}.jw-buffer{background-color:rgba(255,255,255,0.3)}.jw-knob{height:13px;width:13px;background-color:#fff;border-radius:50%;box-shadow:0 0 10px rgba(0,0,0,0.4);opacity:1;pointer-events:none;position:absolute;-webkit-transform:translate(-50%, -50%) scale(0);transform:translate(-50%, -50%) scale(0);transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, -webkit-transform;transition-property:opacity, transform;transition-property:opacity, transform, -webkit-transform}.jw-flag-dragging .jw-slider-time .jw-knob,.jw-icon-volume:active .jw-slider-volume .jw-knob{box-shadow:0 0 26px rgba(0,0,0,0.2),0 0 10px rgba(0,0,0,0.4),0 0 0 6px rgba(255,255,255,0.2)}.jw-slider-horizontal,.jw-slider-vertical{display:flex}.jw-slider-horizontal .jw-slider-container{height:5px;width:100%}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress,.jw-slider-horizontal .jw-cue,.jw-slider-horizontal .jw-knob{top:50%}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress,.jw-slider-horizontal .jw-cue{-webkit-transform:translate(0, -50%);transform:translate(0, -50%)}.jw-slider-horizontal .jw-rail,.jw-slider-horizontal .jw-buffer,.jw-slider-horizontal .jw-progress{height:5px}.jw-slider-horizontal .jw-rail{width:100%}.jw-slider-vertical{align-items:center;flex-direction:column}.jw-slider-vertical .jw-slider-container{height:88px;width:5px}.jw-slider-vertical .jw-rail,.jw-slider-vertical .jw-buffer,.jw-slider-vertical .jw-progress,.jw-slider-vertical .jw-knob{left:50%}.jw-slider-vertical .jw-rail,.jw-slider-vertical .jw-buffer,.jw-slider-vertical .jw-progress{height:100%;width:5px;-webkit-backface-visibility:hidden;backface-visibility:hidden;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);transition:-webkit-transform 150ms ease-in-out;transition:transform 150ms ease-in-out;transition:transform 150ms ease-in-out, -webkit-transform 150ms ease-in-out;bottom:0}.jw-slider-vertical .jw-knob{-webkit-transform:translate(-50%, 50%);transform:translate(-50%, 50%)}.jw-slider-time.jw-tab-focus:focus .jw-rail{outline:solid 2px #4d90fe}.jw-slider-time,.jw-flag-audio-player .jw-slider-volume{height:17px;width:100%;align-items:center;background:transparent none;padding:0 12px}.jw-slider-time .jw-cue{background-color:rgba(33,33,33,0.8);cursor:pointer;position:absolute;width:6px}.jw-slider-time,.jw-horizontal-volume-container{z-index:1;outline:none}.jw-slider-time .jw-rail,.jw-horizontal-volume-container .jw-rail,.jw-slider-time .jw-buffer,.jw-horizontal-volume-container .jw-buffer,.jw-slider-time .jw-progress,.jw-horizontal-volume-container .jw-progress,.jw-slider-time .jw-cue,.jw-horizontal-volume-container .jw-cue{-webkit-backface-visibility:hidden;backface-visibility:hidden;height:100%;-webkit-transform:translate(0, -50%) scale(1, .6);transform:translate(0, -50%) scale(1, .6);transition:-webkit-transform 150ms ease-in-out;transition:transform 150ms ease-in-out;transition:transform 150ms ease-in-out, -webkit-transform 150ms ease-in-out}.jw-slider-time:hover .jw-rail,.jw-horizontal-volume-container:hover .jw-rail,.jw-slider-time:focus .jw-rail,.jw-horizontal-volume-container:focus .jw-rail,.jw-flag-dragging .jw-slider-time .jw-rail,.jw-flag-dragging .jw-horizontal-volume-container .jw-rail,.jw-flag-touch .jw-slider-time .jw-rail,.jw-flag-touch .jw-horizontal-volume-container .jw-rail,.jw-slider-time:hover .jw-buffer,.jw-horizontal-volume-container:hover .jw-buffer,.jw-slider-time:focus .jw-buffer,.jw-horizontal-volume-container:focus .jw-buffer,.jw-flag-dragging .jw-slider-time .jw-buffer,.jw-flag-dragging .jw-horizontal-volume-container .jw-buffer,.jw-flag-touch .jw-slider-time .jw-buffer,.jw-flag-touch .jw-horizontal-volume-container .jw-buffer,.jw-slider-time:hover .jw-progress,.jw-horizontal-volume-container:hover .jw-progress,.jw-slider-time:focus .jw-progress,.jw-horizontal-volume-container:focus .jw-progress,.jw-flag-dragging .jw-slider-time .jw-progress,.jw-flag-dragging .jw-horizontal-volume-container .jw-progress,.jw-flag-touch .jw-slider-time .jw-progress,.jw-flag-touch .jw-horizontal-volume-container .jw-progress,.jw-slider-time:hover .jw-cue,.jw-horizontal-volume-container:hover .jw-cue,.jw-slider-time:focus .jw-cue,.jw-horizontal-volume-container:focus .jw-cue,.jw-flag-dragging .jw-slider-time .jw-cue,.jw-flag-dragging .jw-horizontal-volume-container .jw-cue,.jw-flag-touch .jw-slider-time .jw-cue,.jw-flag-touch .jw-horizontal-volume-container .jw-cue{-webkit-transform:translate(0, -50%) scale(1, 1);transform:translate(0, -50%) scale(1, 1)}.jw-slider-time:hover .jw-knob,.jw-horizontal-volume-container:hover .jw-knob,.jw-slider-time:focus .jw-knob,.jw-horizontal-volume-container:focus .jw-knob{-webkit-transform:translate(-50%, -50%) scale(1);transform:translate(-50%, -50%) scale(1)}.jw-slider-time .jw-rail,.jw-horizontal-volume-container .jw-rail{background-color:rgba(255,255,255,0.2)}.jw-slider-time .jw-buffer,.jw-horizontal-volume-container .jw-buffer{background-color:rgba(255,255,255,0.4)}.jw-flag-touch .jw-slider-time::before,.jw-flag-touch .jw-horizontal-volume-container::before{height:44px;width:100%;content:"";position:absolute;display:block;bottom:calc(100% - 17px);left:0}.jw-slider-time.jw-tab-focus:focus .jw-rail,.jw-horizontal-volume-container.jw-tab-focus:focus .jw-rail{outline:solid 2px #4d90fe}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time{height:17px;padding:0}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time .jw-slider-container{height:10px}.jw-breakpoint--1:not(.jw-flag-audio-player) .jw-slider-time .jw-knob{border-radius:0;border:1px solid rgba(0,0,0,0.75);height:12px;width:10px}.jw-modal{width:284px}.jw-breakpoint-7 .jw-modal,.jw-breakpoint-6 .jw-modal,.jw-breakpoint-5 .jw-modal{height:232px}.jw-breakpoint-4 .jw-modal,.jw-breakpoint-3 .jw-modal{height:192px}.jw-breakpoint-2 .jw-modal,.jw-flag-small-player .jw-modal{bottom:0;right:0;height:100%;width:100%;max-height:none;max-width:none;z-index:2}.jwplayer .jw-rightclick{display:none;position:absolute;white-space:nowrap}.jwplayer .jw-rightclick.jw-open{display:block}.jwplayer .jw-rightclick .jw-rightclick-list{border-radius:1px;list-style:none;margin:0;padding:0}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item{background-color:rgba(0,0,0,0.8);border-bottom:1px solid #444;margin:0}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-logo{color:#fff;display:inline-flex;padding:0 10px 0 0;vertical-align:middle}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-logo .jw-svg-icon{height:20px;width:20px}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item .jw-rightclick-link{border:none;color:#fff;display:block;font-size:11px;line-height:1em;padding:15px 23px;text-align:start;text-decoration:none;width:100%}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item:last-child{border-bottom:none}.jwplayer .jw-rightclick .jw-rightclick-list .jw-rightclick-item:hover{cursor:pointer}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured{vertical-align:middle}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured .jw-rightclick-link{color:#fff}.jwplayer .jw-rightclick .jw-rightclick-list .jw-featured .jw-rightclick-link span{color:#fff}.jwplayer .jw-rightclick .jw-info-overlay-item,.jwplayer .jw-rightclick .jw-share-item,.jwplayer .jw-rightclick .jw-shortcuts-item{border:none;background-color:transparent;outline:none;cursor:pointer}.jw-icon-tooltip.jw-open .jw-overlay{opacity:1;pointer-events:auto;transition-delay:0s}.jw-icon-tooltip.jw-open .jw-overlay:focus{outline:none}.jw-icon-tooltip.jw-open .jw-overlay:focus.jw-tab-focus{outline:solid 2px #4d90fe}.jw-slider-time .jw-overlay:before{height:1em;top:auto}.jw-slider-time .jw-icon-tooltip.jw-open .jw-overlay{pointer-events:none}.jw-volume-tip{padding:13px 0 26px}.jw-time-tip,.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{height:auto;width:100%;box-shadow:0 0 10px rgba(0,0,0,0.4);color:#fff;display:block;margin:0 0 14px;pointer-events:none;position:relative;z-index:0}.jw-time-tip::after,.jw-controlbar .jw-tooltip::after,.jw-settings-menu .jw-tooltip::after{top:100%;position:absolute;left:50%;height:14px;width:14px;border-radius:1px;background-color:currentColor;-webkit-transform-origin:75% 50%;transform-origin:75% 50%;-webkit-transform:translate(-50%, -50%) rotate(45deg);transform:translate(-50%, -50%) rotate(45deg);z-index:-1}.jw-time-tip .jw-text,.jw-controlbar .jw-tooltip .jw-text,.jw-settings-menu .jw-tooltip .jw-text{background-color:#fff;border-radius:1px;color:#000;font-size:10px;height:auto;line-height:1;padding:7px 10px;display:inline-block;min-width:100%;vertical-align:middle}.jw-controlbar .jw-overlay{position:absolute;bottom:100%;left:50%;margin:0;min-height:44px;min-width:44px;opacity:0;pointer-events:none;transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility;transition-delay:0s, 150ms;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);width:100%;z-index:1}.jw-controlbar .jw-overlay .jw-contents{position:relative}.jw-controlbar .jw-option{position:relative;white-space:nowrap;cursor:pointer;list-style:none;height:1.5em;font-family:inherit;line-height:1.5em;padding:0 .5em;font-size:.8em;margin:0}.jw-controlbar .jw-option::before{padding-right:.125em}.jw-controlbar .jw-tooltip,.jw-settings-menu .jw-tooltip{position:absolute;bottom:100%;left:50%;opacity:0;-webkit-transform:translate(-50%, 0);transform:translate(-50%, 0);transition:100ms 0s cubic-bezier(0, .25, .25, 1);transition-property:opacity, visibility, -webkit-transform;transition-property:opacity, transform, visibility;transition-property:opacity, transform, visibility, -webkit-transform;visibility:hidden;white-space:nowrap;width:auto;z-index:1}.jw-controlbar .jw-tooltip.jw-open,.jw-settings-menu .jw-tooltip.jw-open{opacity:1;-webkit-transform:translate(-50%, -10px);transform:translate(-50%, -10px);transition-duration:150ms;transition-delay:500ms,0s,500ms;visibility:visible}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen{left:auto;right:0;-webkit-transform:translate(0, 0);transform:translate(0, 0)}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen.jw-open,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen.jw-open{-webkit-transform:translate(0, -10px);transform:translate(0, -10px)}.jw-controlbar .jw-tooltip.jw-tooltip-fullscreen::after,.jw-settings-menu .jw-tooltip.jw-tooltip-fullscreen::after{left:auto;right:9px}.jw-tooltip-time{height:auto;width:0;bottom:100%;line-height:normal;padding:0;pointer-events:none;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none}.jw-tooltip-time .jw-overlay{bottom:0;min-height:0;width:auto}.jw-tooltip{bottom:57px;display:none;position:absolute}.jw-tooltip .jw-text{height:100%;white-space:nowrap;text-overflow:ellipsis;direction:unset;max-width:246px;overflow:hidden}.jw-flag-audio-player .jw-tooltip{display:none}.jw-flag-small-player .jw-time-thumb{display:none}.jwplayer .jw-shortcuts-tooltip{top:50%;position:absolute;left:50%;background:#333;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);display:none;color:#fff;pointer-events:all;-webkit-user-select:text;-moz-user-select:text;-ms-user-select:text;user-select:text;overflow:hidden;flex-direction:column;z-index:1}.jwplayer .jw-shortcuts-tooltip.jw-open{display:flex}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-close{flex:0 0 auto;margin:5px 5px 5px auto}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container{display:flex;flex:1 1 auto;flex-flow:column;font-size:12px;margin:0 20px 20px;overflow-y:auto;padding:5px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container::-webkit-scrollbar{background-color:transparent;width:6px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-title{font-weight:bold}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-header{align-items:center;display:flex;justify-content:space-between;margin-bottom:10px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list{display:flex;max-width:340px;margin:0 10px}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-tooltip-descriptions{width:100%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row{display:flex;align-items:center;justify-content:space-between;margin:10px 0;width:100%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row .jw-shortcuts-description{margin-right:10px;max-width:70%}.jwplayer .jw-shortcuts-tooltip .jw-shortcuts-container .jw-shortcuts-tooltip-list .jw-shortcuts-row .jw-shortcuts-key{background:#fefefe;color:#333;overflow:hidden;padding:7px 10px;text-overflow:ellipsis;white-space:nowrap}.jw-skip{color:rgba(255,255,255,0.8);cursor:default;position:absolute;display:flex;right:.75em;bottom:56px;padding:.5em;border:1px solid #333;background-color:#000;align-items:center;height:2em}.jw-skip.jw-tab-focus:focus{outline:solid 2px #4d90fe}.jw-skip.jw-skippable{cursor:pointer;padding:.25em .75em}.jw-skip.jw-skippable:hover{cursor:pointer;color:#fff}.jw-skip.jw-skippable .jw-skip-icon{display:inline;height:24px;width:24px;margin:0}.jw-breakpoint-7 .jw-skip{padding:1.35em 1em;bottom:130px}.jw-breakpoint-7 .jw-skip .jw-text{font-size:1em;font-weight:normal}.jw-breakpoint-7 .jw-skip .jw-icon-inline{height:30px;width:30px}.jw-breakpoint-7 .jw-skip .jw-icon-inline .jw-svg-icon{height:30px;width:30px}.jw-skip .jw-skip-icon{display:none;margin-left:-0.75em;padding:0 .5em;pointer-events:none}.jw-skip .jw-skip-icon .jw-svg-icon-next{display:block;padding:0}.jw-skip .jw-text,.jw-skip .jw-skip-icon{vertical-align:middle;font-size:.7em}.jw-skip .jw-text{font-weight:bold}.jw-cast{background-size:cover;display:none;height:100%;position:relative;width:100%}.jw-cast-container{background:linear-gradient(180deg, rgba(25,25,25,0.75), rgba(25,25,25,0.25), rgba(25,25,25,0));left:0;padding:20px 20px 80px;position:absolute;top:0;width:100%}.jw-cast-text{color:#fff;font-size:1.6em}.jw-breakpoint--1 .jw-cast-text,.jw-breakpoint-0 .jw-cast-text{font-size:1.15em}.jw-breakpoint-1 .jw-cast-text,.jw-breakpoint-2 .jw-cast-text,.jw-breakpoint-3 .jw-cast-text{font-size:1.3em}.jw-nextup-container{position:absolute;bottom:66px;left:0;background-color:transparent;cursor:pointer;margin:0 auto;padding:12px;pointer-events:none;right:0;text-align:right;visibility:hidden;width:100%}.jw-settings-open .jw-nextup-container,.jw-info-open .jw-nextup-container{display:none}.jw-breakpoint-7 .jw-nextup-container{padding:60px}.jw-flag-small-player .jw-nextup-container{padding:0 12px 0 0}.jw-flag-small-player .jw-nextup-container .jw-nextup-title,.jw-flag-small-player .jw-nextup-container .jw-nextup-duration,.jw-flag-small-player .jw-nextup-container .jw-nextup-close{display:none}.jw-flag-small-player .jw-nextup-container .jw-nextup-tooltip{height:30px}.jw-flag-small-player .jw-nextup-container .jw-nextup-header{font-size:12px}.jw-flag-small-player .jw-nextup-container .jw-nextup-body{justify-content:center;align-items:center;padding:.75em .3em}.jw-flag-small-player .jw-nextup-container .jw-nextup-thumbnail{width:50%}.jw-flag-small-player .jw-nextup-container .jw-nextup{max-width:65px}.jw-flag-small-player .jw-nextup-container .jw-nextup.jw-nextup-thumbnail-visible{max-width:120px}.jw-nextup{background:#333;border-radius:0;box-shadow:0 0 10px rgba(0,0,0,0.5);color:rgba(255,255,255,0.8);display:inline-block;max-width:280px;overflow:hidden;opacity:0;position:relative;width:64%;pointer-events:all;-webkit-transform:translate(0, -5px);transform:translate(0, -5px);transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:opacity, -webkit-transform;transition-property:opacity, transform;transition-property:opacity, transform, -webkit-transform;transition-delay:0s}.jw-nextup:hover .jw-nextup-tooltip{color:#fff}.jw-nextup.jw-nextup-thumbnail-visible{max-width:400px}.jw-nextup.jw-nextup-thumbnail-visible .jw-nextup-thumbnail{display:block}.jw-nextup-container-visible{visibility:visible}.jw-nextup-container-visible .jw-nextup{opacity:1;-webkit-transform:translate(0, 0);transform:translate(0, 0);transition-delay:0s, 0s, 150ms}.jw-nextup-tooltip{display:flex;height:80px}.jw-nextup-thumbnail{width:120px;background-position:center;background-size:cover;flex:0 0 auto;display:none}.jw-nextup-body{flex:1 1 auto;overflow:hidden;padding:.75em .875em;display:flex;flex-flow:column wrap;justify-content:space-between}.jw-nextup-header,.jw-nextup-title{font-size:14px;line-height:1.35}.jw-nextup-header{font-weight:bold}.jw-nextup-title{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;width:100%}.jw-nextup-duration{align-self:flex-end;text-align:right;font-size:12px}.jw-nextup-close{height:24px;width:24px;border:none;color:rgba(255,255,255,0.8);cursor:pointer;margin:6px;visibility:hidden}.jw-nextup-close:hover{color:#fff}.jw-nextup-sticky .jw-nextup-close{visibility:visible}.jw-autostart-mute{position:absolute;bottom:0;right:12px;height:44px;width:44px;background-color:rgba(33,33,33,0.4);padding:5px 4px 5px 6px;display:none}.jwplayer.jw-flag-autostart:not(.jw-flag-media-audio) .jw-nextup{display:none}.jw-settings-menu{position:absolute;bottom:57px;right:12px;align-items:flex-start;background-color:#333;display:none;flex-flow:column nowrap;max-width:284px;pointer-events:auto}.jw-settings-open .jw-settings-menu{display:flex}.jw-breakpoint-7 .jw-settings-menu{bottom:130px;right:60px;max-height:none;max-width:none;height:35%;width:25%}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline{height:60px;width:60px}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline .jw-svg-icon{height:30px;width:30px}.jw-breakpoint-7 .jw-settings-menu .jw-settings-topbar:not(.jw-nested-menu-open) .jw-icon-inline .jw-tooltip .jw-text{font-size:1em}.jw-breakpoint-7 .jw-settings-menu .jw-settings-back{min-width:60px}.jw-breakpoint-6 .jw-settings-menu,.jw-breakpoint-5 .jw-settings-menu{height:232px;width:284px;max-height:232px}.jw-breakpoint-4 .jw-settings-menu,.jw-breakpoint-3 .jw-settings-menu{height:192px;width:284px;max-height:192px}.jw-breakpoint-2 .jw-settings-menu{height:179px;width:284px;max-height:179px}.jw-flag-small-player .jw-settings-menu{max-width:none}.jw-settings-menu .jw-icon.jw-button-color::after{height:100%;width:24px;box-shadow:inset 0 -3px 0 -1px currentColor;margin:auto;opacity:0;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-settings-menu .jw-icon.jw-button-color[aria-checked="true"]::after{opacity:1}.jw-settings-menu .jw-settings-reset{text-decoration:underline}.jw-settings-topbar{align-items:center;background-color:rgba(0,0,0,0.4);display:flex;flex:0 0 auto;padding:3px 5px 0;width:100%}.jw-settings-topbar.jw-nested-menu-open{padding:0}.jw-settings-topbar.jw-nested-menu-open .jw-icon:not(.jw-settings-close):not(.jw-settings-back){display:none}.jw-settings-topbar.jw-nested-menu-open .jw-svg-icon-close{width:20px}.jw-settings-topbar.jw-nested-menu-open .jw-svg-icon-arrow-left{height:12px}.jw-settings-topbar.jw-nested-menu-open .jw-settings-topbar-text{display:block;outline:none}.jw-settings-topbar .jw-settings-back{min-width:44px}.jw-settings-topbar .jw-settings-topbar-buttons{display:inherit;width:100%;height:100%}.jw-settings-topbar .jw-settings-topbar-text{display:none;color:#fff;font-size:13px;width:100%}.jw-settings-topbar .jw-settings-close{margin-left:auto}.jw-settings-submenu{display:none;flex:1 1 auto;overflow-y:auto;padding:8px 20px 0 5px}.jw-settings-submenu::-webkit-scrollbar{background-color:transparent;width:6px}.jw-settings-submenu::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jw-settings-submenu.jw-settings-submenu-active{display:block}.jw-settings-submenu .jw-submenu-topbar{box-shadow:0 2px 9px 0 #1d1d1d;background-color:#2f2d2d;margin:-8px -20px 0 -5px}.jw-settings-submenu .jw-submenu-topbar .jw-settings-content-item{cursor:pointer;text-align:right;padding-right:15px;text-decoration:underline}.jw-settings-submenu .jw-settings-value-wrapper{float:right;display:flex;align-items:center}.jw-settings-submenu .jw-settings-value-wrapper .jw-settings-content-item-arrow{display:flex}.jw-settings-submenu .jw-settings-value-wrapper .jw-svg-icon-arrow-right{width:8px;margin-left:5px;height:12px}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item{font-size:1em;padding:11px 15px 11px 30px}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item .jw-settings-item-active::before{justify-content:flex-end}.jw-breakpoint-7 .jw-settings-submenu .jw-settings-content-item .jw-auto-label{font-size:.85em;padding-left:10px}.jw-flag-touch .jw-settings-submenu{overflow-y:scroll;-webkit-overflow-scrolling:touch}.jw-auto-label{font-size:10px;font-weight:initial;opacity:.75;padding-left:5px}.jw-settings-content-item{position:relative;color:rgba(255,255,255,0.8);cursor:pointer;font-size:12px;line-height:1;padding:7px 0 7px 15px;width:100%;text-align:left;outline:none}.jw-settings-content-item:hover{color:#fff}.jw-settings-content-item:focus{font-weight:bold}.jw-flag-small-player .jw-settings-content-item{line-height:1.75}.jw-settings-content-item.jw-tab-focus:focus{border:solid 2px #4d90fe}.jw-settings-item-active{font-weight:bold;position:relative}.jw-settings-item-active::before{height:100%;width:1em;align-items:center;content:"\\2022";display:inline-flex;justify-content:center}.jw-breakpoint-2 .jw-settings-open .jw-display-container,.jw-flag-small-player .jw-settings-open .jw-display-container,.jw-flag-touch .jw-settings-open .jw-display-container{display:none}.jw-breakpoint-2 .jw-settings-open.jw-controls,.jw-flag-small-player .jw-settings-open.jw-controls,.jw-flag-touch .jw-settings-open.jw-controls{z-index:1}.jw-flag-small-player .jw-settings-open .jw-controlbar{display:none}.jw-settings-open .jw-icon-settings::after{opacity:1}.jw-settings-open .jw-tooltip-settings{display:none}.jw-sharing-link{cursor:pointer}.jw-shortcuts-container .jw-switch{position:relative;display:inline-block;transition:ease-out .15s;transition-property:opacity, background;border-radius:18px;width:80px;height:20px;padding:10px;background:rgba(80,80,80,0.8);cursor:pointer;font-size:inherit;vertical-align:middle}.jw-shortcuts-container .jw-switch.jw-tab-focus{outline:solid 2px #4d90fe}.jw-shortcuts-container .jw-switch .jw-switch-knob{position:absolute;top:2px;left:1px;transition:ease-out .15s;box-shadow:0 0 10px rgba(0,0,0,0.4);border-radius:13px;width:15px;height:15px;background:#fefefe}.jw-shortcuts-container .jw-switch:before,.jw-shortcuts-container .jw-switch:after{position:absolute;top:3px;transition:inherit;color:#fefefe}.jw-shortcuts-container .jw-switch:before{content:attr(data-jw-switch-disabled);right:8px}.jw-shortcuts-container .jw-switch:after{content:attr(data-jw-switch-enabled);left:8px;opacity:0}.jw-shortcuts-container .jw-switch[aria-checked="true"]{background:#475470}.jw-shortcuts-container .jw-switch[aria-checked="true"]:before{opacity:0}.jw-shortcuts-container .jw-switch[aria-checked="true"]:after{opacity:1}.jw-shortcuts-container .jw-switch[aria-checked="true"] .jw-switch-knob{left:60px}.jw-idle-icon-text{display:none;line-height:1;position:absolute;text-align:center;text-indent:.35em;top:100%;white-space:nowrap;left:50%;-webkit-transform:translateX(-50%);transform:translateX(-50%)}.jw-idle-label{border-radius:50%;color:#fff;-webkit-filter:drop-shadow(1px 1px 5px rgba(12,26,71,0.25));filter:drop-shadow(1px 1px 5px rgba(12,26,71,0.25));font:normal 16px/1 Arial,Helvetica,sans-serif;position:relative;transition:background-color 150ms cubic-bezier(0, .25, .25, 1);transition-property:background-color,-webkit-filter;transition-property:background-color,filter;transition-property:background-color,filter,-webkit-filter;-webkit-font-smoothing:antialiased}.jw-state-idle .jw-icon-display.jw-idle-label .jw-idle-icon-text{display:block}.jw-state-idle .jw-icon-display.jw-idle-label .jw-svg-icon-play{-webkit-transform:scale(.7, .7);transform:scale(.7, .7)}.jw-breakpoint-0.jw-state-idle .jw-icon-display.jw-idle-label,.jw-breakpoint--1.jw-state-idle .jw-icon-display.jw-idle-label{font-size:12px}.jw-info-overlay{top:50%;position:absolute;left:50%;background:#333;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);display:none;color:#fff;pointer-events:all;-webkit-user-select:text;-moz-user-select:text;-ms-user-select:text;user-select:text;overflow:hidden;flex-direction:column}.jw-info-overlay .jw-info-close{flex:0 0 auto;margin:5px 5px 5px auto}.jw-info-open .jw-info-overlay{display:flex}.jw-info-container{display:flex;flex:1 1 auto;flex-flow:column;margin:0 20px 20px;overflow-y:auto;padding:5px}.jw-info-container [class*="jw-info"]:not(:first-of-type){color:rgba(255,255,255,0.8);padding-top:10px;font-size:12px}.jw-info-container .jw-info-description{margin-bottom:30px;text-align:start}.jw-info-container .jw-info-description:empty{display:none}.jw-info-container .jw-info-duration{text-align:start}.jw-info-container .jw-info-title{text-align:start;font-size:12px;font-weight:bold}.jw-info-container::-webkit-scrollbar{background-color:transparent;width:6px}.jw-info-container::-webkit-scrollbar-thumb{background-color:#fff;border:1px solid #333;border-radius:6px}.jw-info-clientid{align-self:flex-end;font-size:12px;color:rgba(255,255,255,0.8);margin:0 20px 20px 44px;text-align:right}.jw-flag-touch .jw-info-open .jw-display-container{display:none}@supports ((-webkit-filter: drop-shadow(0 0 3px #000)) or (filter: drop-shadow(0 0 3px #000))){.jwplayer.jw-ab-drop-shadow .jw-controls .jw-svg-icon,.jwplayer.jw-ab-drop-shadow .jw-controls .jw-icon.jw-text,.jwplayer.jw-ab-drop-shadow .jw-slider-container .jw-rail,.jwplayer.jw-ab-drop-shadow .jw-title{text-shadow:none;box-shadow:none;-webkit-filter:drop-shadow(0 2px 3px rgba(0,0,0,0.3));filter:drop-shadow(0 2px 3px rgba(0,0,0,0.3))}.jwplayer.jw-ab-drop-shadow .jw-button-color{opacity:.8;transition-property:color, opacity}.jwplayer.jw-ab-drop-shadow .jw-button-color:not(:hover){color:#fff;opacity:.8}.jwplayer.jw-ab-drop-shadow .jw-button-color:hover{opacity:1}.jwplayer.jw-ab-drop-shadow .jw-controls-backdrop{background-image:linear-gradient(to bottom, hsla(0, 0%, 0%, 0), hsla(0, 0%, 0%, 0.00787) 10.79%, hsla(0, 0%, 0%, 0.02963) 21.99%, hsla(0, 0%, 0%, 0.0625) 33.34%, hsla(0, 0%, 0%, 0.1037) 44.59%, hsla(0, 0%, 0%, 0.15046) 55.48%, hsla(0, 0%, 0%, 0.2) 65.75%, hsla(0, 0%, 0%, 0.24954) 75.14%, hsla(0, 0%, 0%, 0.2963) 83.41%, hsla(0, 0%, 0%, 0.3375) 90.28%, hsla(0, 0%, 0%, 0.37037) 95.51%, hsla(0, 0%, 0%, 0.39213) 98.83%, hsla(0, 0%, 0%, 0.4));mix-blend-mode:multiply;transition-property:opacity}.jw-state-idle.jwplayer.jw-ab-drop-shadow .jw-controls-backdrop{background-image:linear-gradient(to bottom, hsla(0, 0%, 0%, 0.2), hsla(0, 0%, 0%, 0.19606) 1.17%, hsla(0, 0%, 0%, 0.18519) 4.49%, hsla(0, 0%, 0%, 0.16875) 9.72%, hsla(0, 0%, 0%, 0.14815) 16.59%, hsla(0, 0%, 0%, 0.12477) 24.86%, hsla(0, 0%, 0%, 0.1) 34.25%, hsla(0, 0%, 0%, 0.07523) 44.52%, hsla(0, 0%, 0%, 0.05185) 55.41%, hsla(0, 0%, 0%, 0.03125) 66.66%, hsla(0, 0%, 0%, 0.01481) 78.01%, hsla(0, 0%, 0%, 0.00394) 89.21%, hsla(0, 0%, 0%, 0));background-size:100% 7rem;background-position:50% 0}.jwplayer.jw-ab-drop-shadow.jw-state-idle .jw-controls{background-color:transparent}}.jw-video-thumbnail-container{position:relative;overflow:hidden}.jw-video-thumbnail-container:not(.jw-related-shelf-item-image){height:100%;width:100%}.jw-video-thumbnail-container.jw-video-thumbnail-generated{position:absolute;top:0;left:0}.jw-video-thumbnail-container:hover,.jw-related-item-content:hover .jw-video-thumbnail-container,.jw-related-shelf-item:hover .jw-video-thumbnail-container{cursor:pointer}.jw-video-thumbnail-container:hover .jw-video-thumbnail:not(.jw-video-thumbnail-completed),.jw-related-item-content:hover .jw-video-thumbnail-container .jw-video-thumbnail:not(.jw-video-thumbnail-completed),.jw-related-shelf-item:hover .jw-video-thumbnail-container .jw-video-thumbnail:not(.jw-video-thumbnail-completed){opacity:1}.jw-video-thumbnail-container .jw-video-thumbnail{position:absolute;top:50%;left:50%;bottom:unset;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%);width:100%;height:auto;min-width:100%;min-height:100%;opacity:0;transition:opacity .3s ease;object-fit:cover;background:#000}.jw-related-item-next-up .jw-video-thumbnail-container .jw-video-thumbnail{height:100%;width:auto}.jw-video-thumbnail-container .jw-video-thumbnail.jw-video-thumbnail-visible:not(.jw-video-thumbnail-completed){opacity:1}.jw-video-thumbnail-container .jw-video-thumbnail.jw-video-thumbnail-completed{opacity:0}.jw-video-thumbnail-container .jw-video-thumbnail~.jw-svg-icon-play{display:none}.jw-video-thumbnail-container .jw-video-thumbnail+.jw-related-shelf-item-aspect{pointer-events:none}.jw-video-thumbnail-container .jw-video-thumbnail+.jw-related-item-poster-content{pointer-events:none}.jw-state-idle:not(.jw-flag-cast-available) .jw-display{padding:0}.jw-state-idle .jw-controls{background:rgba(0,0,0,0.4)}.jw-state-idle.jw-flag-cast-available:not(.jw-flag-audio-player) .jw-controlbar .jw-slider-time,.jw-state-idle.jw-flag-cardboard-available .jw-controlbar .jw-slider-time,.jw-state-idle.jw-flag-cast-available:not(.jw-flag-audio-player) .jw-controlbar .jw-icon:not(.jw-icon-cardboard):not(.jw-icon-cast):not(.jw-icon-airplay),.jw-state-idle.jw-flag-cardboard-available .jw-controlbar .jw-icon:not(.jw-icon-cardboard):not(.jw-icon-cast):not(.jw-icon-airplay){display:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-icon:focus{border:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-icon .jw-svg-icon-buffer{-webkit-animation:jw-spin 2s linear infinite;animation:jw-spin 2s linear infinite;display:block}@-webkit-keyframes jw-spin{100%{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}@keyframes jw-spin{100%{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}.jwplayer.jw-state-buffering .jw-icon-playback .jw-svg-icon-play{display:none}.jwplayer.jw-state-buffering .jw-icon-display .jw-svg-icon-pause{display:none}.jwplayer.jw-state-playing .jw-display .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-playing .jw-icon-playback .jw-svg-icon-play{display:none}.jwplayer.jw-state-playing .jw-display .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-playing .jw-icon-playback .jw-svg-icon-pause{display:block}.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio) .jw-controls-backdrop{opacity:0}.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio) .jw-logo-bottom-left,.jwplayer.jw-state-playing.jw-flag-user-inactive:not(.jw-flag-audio-player):not(.jw-flag-casting):not(.jw-flag-media-audio):not(.jw-flag-autostart) .jw-logo-bottom-right{bottom:0}.jwplayer .jw-icon-playback .jw-svg-icon-stop{display:none}.jwplayer.jw-state-paused .jw-svg-icon-pause,.jwplayer.jw-state-idle .jw-svg-icon-pause,.jwplayer.jw-state-error .jw-svg-icon-pause,.jwplayer.jw-state-complete .jw-svg-icon-pause{display:none}.jwplayer.jw-state-error .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-complete .jw-icon-display .jw-svg-icon-play,.jwplayer.jw-state-buffering .jw-icon-display .jw-svg-icon-play{display:none}.jwplayer:not(.jw-state-buffering) .jw-svg-icon-buffer{display:none}.jwplayer:not(.jw-state-complete) .jw-svg-icon-replay{display:none}.jwplayer:not(.jw-state-error) .jw-svg-icon-error{display:none}.jwplayer.jw-state-complete .jw-display .jw-icon-display .jw-svg-icon-replay{display:block}.jwplayer.jw-state-complete .jw-display .jw-text{display:none}.jwplayer.jw-state-complete .jw-controls{background:rgba(0,0,0,0.4);height:100%}.jw-state-idle .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-paused .jw-icon-playback .jw-svg-icon-pause,.jwplayer.jw-state-paused .jw-icon-display .jw-svg-icon-pause,.jwplayer.jw-state-complete .jw-icon-playback .jw-svg-icon-pause{display:none}.jw-state-idle .jw-display-icon-rewind,.jwplayer.jw-state-buffering .jw-display-icon-rewind,.jwplayer.jw-state-complete .jw-display-icon-rewind,body .jw-error .jw-display-icon-rewind,body .jwplayer.jw-state-error .jw-display-icon-rewind,.jw-state-idle .jw-display-icon-next,.jwplayer.jw-state-buffering .jw-display-icon-next,.jwplayer.jw-state-complete .jw-display-icon-next,body .jw-error .jw-display-icon-next,body .jwplayer.jw-state-error .jw-display-icon-next{display:none}body .jw-error .jw-icon-display,body .jwplayer.jw-state-error .jw-icon-display{cursor:default}body .jw-error .jw-icon-display .jw-svg-icon-error,body .jwplayer.jw-state-error .jw-icon-display .jw-svg-icon-error{display:block}body .jw-error .jw-icon-container{position:absolute;width:100%;height:100%;top:0;left:0;bottom:0;right:0}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-preview{display:none}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title{padding-top:4px}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title-primary{width:auto;display:inline-block;padding-right:.5ch}body .jwplayer.jw-state-error.jw-flag-audio-player .jw-title-secondary{width:auto;display:inline-block;padding-left:0}body .jwplayer.jw-state-error .jw-controlbar,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-controlbar{display:none}body .jwplayer.jw-state-error .jw-settings-menu,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-settings-menu{height:100%;top:50%;left:50%;-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%)}body .jwplayer.jw-state-error .jw-display,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-display{padding:0}body .jwplayer.jw-state-error .jw-logo-bottom-left,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-logo-bottom-left,body .jwplayer.jw-state-error .jw-logo-bottom-right,.jwplayer.jw-state-idle:not(.jw-flag-audio-player):not(.jw-flag-cast-available):not(.jw-flag-cardboard-available) .jw-logo-bottom-right{bottom:0}.jwplayer.jw-state-playing.jw-flag-user-inactive .jw-display{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-state-playing:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display,.jwplayer.jw-state-paused:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting):not(.jw-flag-play-rejected) .jw-display{display:none}.jwplayer.jw-state-paused.jw-flag-play-rejected:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display-icon-rewind,.jwplayer.jw-state-paused.jw-flag-play-rejected:not(.jw-flag-touch):not(.jw-flag-small-player):not(.jw-flag-casting) .jw-display-icon-next{display:none}.jwplayer.jw-state-buffering .jw-display-icon-display .jw-text,.jwplayer.jw-state-complete .jw-display .jw-text{display:none}.jwplayer.jw-flag-casting:not(.jw-flag-audio-player) .jw-cast{display:block}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-display-icon-container{display:none}.jwplayer.jw-flag-casting .jw-icon-hd,.jwplayer.jw-flag-casting .jw-captions,.jwplayer.jw-flag-casting .jw-icon-fullscreen,.jwplayer.jw-flag-casting .jw-icon-audio-tracks{display:none}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-icon-volume{display:none}.jwplayer.jw-flag-casting.jw-flag-airplay-casting .jw-icon-airplay{color:#fff}.jw-state-playing.jw-flag-casting:not(.jw-flag-audio-player) .jw-display,.jw-state-paused.jw-flag-casting:not(.jw-flag-audio-player) .jw-display{display:table}.jwplayer.jw-flag-cast-available .jw-icon-cast,.jwplayer.jw-flag-cast-available .jw-icon-airplay{display:flex}.jwplayer.jw-flag-cardboard-available .jw-icon-cardboard{display:flex}.jwplayer.jw-flag-live .jw-display-icon-rewind{visibility:hidden}.jwplayer.jw-flag-live .jw-controlbar .jw-text-elapsed,.jwplayer.jw-flag-live .jw-controlbar .jw-text-duration,.jwplayer.jw-flag-live .jw-controlbar .jw-text-countdown,.jwplayer.jw-flag-live .jw-controlbar .jw-slider-time{display:none}.jwplayer.jw-flag-live .jw-controlbar .jw-text-alt{display:flex}.jwplayer.jw-flag-live .jw-controlbar .jw-overlay:after{display:none}.jwplayer.jw-flag-live .jw-nextup-container{bottom:44px}.jwplayer.jw-flag-live .jw-text-elapsed,.jwplayer.jw-flag-live .jw-text-duration{display:none}.jwplayer.jw-flag-live .jw-text-live{cursor:default}.jwplayer.jw-flag-live .jw-text-live:hover{color:rgba(255,255,255,0.8)}.jwplayer.jw-flag-live.jw-state-playing .jw-icon-playback .jw-svg-icon-stop,.jwplayer.jw-flag-live.jw-state-buffering .jw-icon-playback .jw-svg-icon-stop{display:block}.jwplayer.jw-flag-live.jw-state-playing .jw-icon-playback .jw-svg-icon-pause,.jwplayer.jw-flag-live.jw-state-buffering .jw-icon-playback .jw-svg-icon-pause{display:none}.jw-text-live{height:24px;width:auto;align-items:center;border-radius:1px;color:rgba(255,255,255,0.8);display:flex;font-size:12px;font-weight:bold;margin-right:10px;padding:0 1ch;text-rendering:geometricPrecision;text-transform:uppercase;transition:150ms cubic-bezier(0, .25, .25, 1);transition-property:box-shadow,color}.jw-text-live::before{height:8px;width:8px;background-color:currentColor;border-radius:50%;margin-right:6px;opacity:1;transition:opacity 150ms cubic-bezier(0, .25, .25, 1)}.jw-text-live.jw-dvr-live{box-shadow:inset 0 0 0 2px currentColor}.jw-text-live.jw-dvr-live::before{opacity:.5}.jw-text-live.jw-dvr-live:hover{color:#fff}.jwplayer.jw-flag-controls-hidden .jw-logo.jw-hide{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-flag-controls-hidden:not(.jw-flag-casting) .jw-logo-top-right{top:0}.jwplayer.jw-flag-controls-hidden .jw-plugin{bottom:.5em}.jwplayer.jw-flag-controls-hidden .jw-nextup-container{bottom:0}.jw-flag-controls-hidden .jw-controlbar,.jw-flag-controls-hidden .jw-display{visibility:hidden;pointer-events:none;opacity:0;transition-delay:0s, 250ms}.jw-flag-controls-hidden .jw-controls-backdrop{opacity:0}.jw-flag-controls-hidden .jw-logo{visibility:visible}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-logo.jw-hide{visibility:hidden;pointer-events:none;opacity:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-casting) .jw-logo-top-right{top:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-plugin{bottom:.5em}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing .jw-nextup-container{bottom:0}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-controls-hidden) .jw-media{cursor:none;-webkit-cursor-visibility:auto-hide}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing.jw-flag-casting .jw-display{display:table}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-state-playing:not(.jw-flag-ads) .jw-autostart-mute{display:flex}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-flag-casting .jw-nextup-container{bottom:66px}.jwplayer.jw-flag-user-inactive:not(.jw-flag-media-audio).jw-flag-casting.jw-state-idle .jw-nextup-container{display:none}.jw-flag-media-audio .jw-preview{display:block}.jwplayer.jw-flag-ads .jw-preview,.jwplayer.jw-flag-ads .jw-logo,.jwplayer.jw-flag-ads .jw-captions.jw-captions-enabled,.jwplayer.jw-flag-ads .jw-nextup-container,.jwplayer.jw-flag-ads .jw-text-duration,.jwplayer.jw-flag-ads .jw-text-elapsed{display:none}.jwplayer.jw-flag-ads video::-webkit-media-text-track-container{display:none}.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-rewind,.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-next,.jwplayer.jw-flag-ads.jw-flag-small-player .jw-display-icon-display{display:none}.jwplayer.jw-flag-ads.jw-flag-small-player.jw-state-buffering .jw-display-icon-display{display:inline-block}.jwplayer.jw-flag-ads .jw-controlbar{flex-wrap:wrap-reverse}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time{height:auto;padding:0;pointer-events:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-slider-container{height:5px}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-rail,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-knob,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-buffer,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-cue,.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-icon-settings{display:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-slider-time .jw-progress{-webkit-transform:none;transform:none;top:auto}.jwplayer.jw-flag-ads .jw-controlbar .jw-tooltip,.jwplayer.jw-flag-ads .jw-controlbar .jw-icon-tooltip:not(.jw-icon-volume),.jwplayer.jw-flag-ads .jw-controlbar .jw-icon-inline:not(.jw-icon-playback):not(.jw-icon-fullscreen):not(.jw-icon-volume){display:none}.jwplayer.jw-flag-ads .jw-controlbar .jw-volume-tip{padding:13px 0}.jwplayer.jw-flag-ads .jw-controlbar .jw-text-alt{display:flex}.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid) .jw-controls .jw-controlbar,.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-autostart .jw-controls .jw-controlbar{display:flex;pointer-events:all;visibility:visible;opacity:1}.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-user-inactive .jw-controls-backdrop,.jwplayer.jw-flag-ads.jw-flag-ads.jw-state-playing.jw-flag-touch:not(.jw-flag-ads-vpaid).jw-flag-autostart.jw-flag-user-inactive .jw-controls-backdrop{opacity:1;background-size:100% 60px}.jwplayer.jw-flag-ads-vpaid .jw-display-container,.jwplayer.jw-flag-touch.jw-flag-ads-vpaid .jw-display-container,.jwplayer.jw-flag-ads-vpaid .jw-skip,.jwplayer.jw-flag-touch.jw-flag-ads-vpaid .jw-skip{display:none}.jwplayer.jw-flag-ads-vpaid.jw-flag-small-player .jw-controls{background:none}.jwplayer.jw-flag-ads-vpaid.jw-flag-small-player .jw-controls::after{content:none}.jwplayer.jw-flag-ads-hide-controls .jw-controls-backdrop,.jwplayer.jw-flag-ads-hide-controls .jw-controls{display:none !important}.jw-flag-overlay-open-related .jw-controls,.jw-flag-overlay-open-related .jw-title,.jw-flag-overlay-open-related .jw-logo{display:none}.jwplayer.jw-flag-rightclick-open{overflow:visible}.jwplayer.jw-flag-rightclick-open .jw-rightclick{z-index:16777215}body .jwplayer.jw-flag-flash-blocked .jw-controls,body .jwplayer.jw-flag-flash-blocked .jw-overlays,body .jwplayer.jw-flag-flash-blocked .jw-controls-backdrop,body .jwplayer.jw-flag-flash-blocked .jw-preview{display:none}body .jwplayer.jw-flag-flash-blocked .jw-error-msg{top:25%}.jw-flag-touch.jw-breakpoint-7 .jw-captions,.jw-flag-touch.jw-breakpoint-6 .jw-captions,.jw-flag-touch.jw-breakpoint-5 .jw-captions,.jw-flag-touch.jw-breakpoint-4 .jw-captions,.jw-flag-touch.jw-breakpoint-7 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-6 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-5 .jw-nextup-container,.jw-flag-touch.jw-breakpoint-4 .jw-nextup-container{bottom:4.25em}.jw-flag-touch .jw-controlbar .jw-icon-volume{display:flex}.jw-flag-touch .jw-display,.jw-flag-touch .jw-display-container,.jw-flag-touch .jw-display-controls{pointer-events:none}.jw-flag-touch.jw-state-paused:not(.jw-breakpoint-1) .jw-display-icon-next,.jw-flag-touch.jw-state-playing:not(.jw-breakpoint-1) .jw-display-icon-next,.jw-flag-touch.jw-state-paused:not(.jw-breakpoint-1) .jw-display-icon-rewind,.jw-flag-touch.jw-state-playing:not(.jw-breakpoint-1) .jw-display-icon-rewind{display:none}.jw-flag-touch.jw-state-paused.jw-flag-dragging .jw-display{display:none}.jw-flag-audio-player{background-color:#000}.jw-flag-audio-player:not(.jw-flag-flash-blocked) .jw-media{visibility:hidden}.jw-flag-audio-player .jw-title{background:none}.jw-flag-audio-player object{min-height:44px}.jw-flag-audio-player:not(.jw-flag-live) .jw-spacer{display:none}.jw-flag-audio-player .jw-preview,.jw-flag-audio-player .jw-display,.jw-flag-audio-player .jw-title,.jw-flag-audio-player .jw-nextup-container{display:none}.jw-flag-audio-player .jw-controlbar{position:relative}.jw-flag-audio-player .jw-controlbar .jw-button-container{padding-right:3px;padding-left:0}.jw-flag-audio-player .jw-controlbar .jw-icon-tooltip,.jw-flag-audio-player .jw-controlbar .jw-icon-inline{display:none}.jw-flag-audio-player .jw-controlbar .jw-icon-volume,.jw-flag-audio-player .jw-controlbar .jw-icon-playback,.jw-flag-audio-player .jw-controlbar .jw-icon-next,.jw-flag-audio-player .jw-controlbar .jw-icon-rewind,.jw-flag-audio-player .jw-controlbar .jw-icon-cast,.jw-flag-audio-player .jw-controlbar .jw-text-live,.jw-flag-audio-player .jw-controlbar .jw-icon-airplay,.jw-flag-audio-player .jw-controlbar .jw-logo-button,.jw-flag-audio-player .jw-controlbar .jw-text-elapsed,.jw-flag-audio-player .jw-controlbar .jw-text-duration{display:flex;flex:0 0 auto}.jw-flag-audio-player .jw-controlbar .jw-text-duration,.jw-flag-audio-player .jw-controlbar .jw-text-countdown{padding-right:10px}.jw-flag-audio-player .jw-controlbar .jw-slider-time{flex:0 1 auto;align-items:center;display:flex;order:1}.jw-flag-audio-player .jw-controlbar .jw-icon-volume{margin-right:0;transition:margin-right 150ms cubic-bezier(0, .25, .25, 1)}.jw-flag-audio-player .jw-controlbar .jw-icon-volume .jw-overlay{display:none}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container{transition:width 300ms cubic-bezier(0, .25, .25, 1);width:0}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open{width:140px}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open .jw-slider-volume{padding-right:24px;transition:opacity 300ms;opacity:1}.jw-flag-audio-player .jw-controlbar .jw-horizontal-volume-container.jw-open~.jw-slider-time{flex:1 1 auto;width:auto;transition:opacity 300ms, width 300ms}.jw-flag-audio-player .jw-controlbar .jw-slider-volume{opacity:0}.jw-flag-audio-player .jw-controlbar .jw-slider-volume .jw-knob{-webkit-transform:translate(-50%, -50%);transform:translate(-50%, -50%)}.jw-flag-audio-player .jw-controlbar .jw-slider-volume~.jw-icon-volume{margin-right:140px}.jw-flag-audio-player.jw-breakpoint-1 .jw-horizontal-volume-container.jw-open~.jw-slider-time,.jw-flag-audio-player.jw-breakpoint-2 .jw-horizontal-volume-container.jw-open~.jw-slider-time{opacity:0}.jw-flag-audio-player.jw-flag-small-player .jw-text-elapsed,.jw-flag-audio-player.jw-flag-small-player .jw-text-duration{display:none}.jw-flag-audio-player.jw-flag-ads .jw-slider-time{display:none}.jw-hidden{display:none}',
        "",
      ]);
    },
  ],
]);
