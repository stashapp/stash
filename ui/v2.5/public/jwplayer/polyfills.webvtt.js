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
(window.webpackJsonpjwplayer=window.webpackJsonpjwplayer||[]).push([[8],{68:function(t,e,i){"use strict";function r(t,e){this.name="ParsingError",this.code=t.code,this.message=e||t.message}function n(t){function e(t,e,i,r){return 3600*(0|t)+60*(0|e)+(0|i)+(0|r)/1e3}var i=t.match(/^(\d+):(\d{2})(:\d{2})?\.(\d{3})/);return i?i[3]?e(i[1],i[2],i[3].replace(":",""),i[4]):i[1]>59?e(i[1],i[2],0,i[4]):e(0,i[1],i[2],i[4]):null}i.r(e),r.prototype=Object.create(Error.prototype),r.prototype.constructor=r,r.Errors={BadSignature:{code:0,message:"Malformed WebVTT signature."},BadTimeStamp:{code:1,message:"Malformed time stamp."}};var o={"&amp;":"&","&lt;":"<","&gt;":">","&lrm;":"‎","&rlm;":"‏","&nbsp;":" "},a={c:"span",i:"i",b:"b",u:"u",ruby:"ruby",rt:"rt",v:"span",lang:"span"},s={v:"title",lang:"lang"},l={rt:"ruby"};function h(t,e){function i(){if(!e)return null;var t,i=e.match(/^([^<]*)(<[^>]+>?)?/);return t=i[1]?i[1]:i[2],e=e.substr(t.length),t}function r(t){return o[t]}function h(t){for(var e;e=t.match(/&(amp|lt|gt|lrm|rlm|nbsp);/);)t=t.replace(e[0],r);return t}function c(t,e){return!l[e.localName]||l[e.localName]===t.localName}function p(e,i){var r=a[e];if(!r)return null;var n=t.document.createElement(r),o=s[e];return o&&i&&(n[o]=i.trim()),n}for(var f,u=t.document.createElement("div"),d=u,g=[];null!==(f=i());)if("<"!==f[0])d.appendChild(t.document.createTextNode(h(f)));else{if("/"===f[1]){g.length&&g[g.length-1]===f.substr(2).replace(">","")&&(g.pop(),d=d.parentNode);continue}var m=n(f.substr(1,f.length-2)),v=void 0;if(m){v=t.document.createProcessingInstruction("timestamp",m),d.appendChild(v);continue}var y=f.match(/^<([^.\s/0-9>]+)(\.[^\s\\>]+)?([^>\\]+)?(\\?)>?$/);if(!y)continue;if(!(v=p(y[1],y[3])))continue;if(!c(d,v))continue;y[2]&&(v.className=y[2].substr(1).replace("."," ")),g.push(y[1]),d.appendChild(v),d=v}return u}var c=[[1470,1470],[1472,1472],[1475,1475],[1478,1478],[1488,1514],[1520,1524],[1544,1544],[1547,1547],[1549,1549],[1563,1563],[1566,1610],[1645,1647],[1649,1749],[1765,1766],[1774,1775],[1786,1805],[1807,1808],[1810,1839],[1869,1957],[1969,1969],[1984,2026],[2036,2037],[2042,2042],[2048,2069],[2074,2074],[2084,2084],[2088,2088],[2096,2110],[2112,2136],[2142,2142],[2208,2208],[2210,2220],[8207,8207],[64285,64285],[64287,64296],[64298,64310],[64312,64316],[64318,64318],[64320,64321],[64323,64324],[64326,64449],[64467,64829],[64848,64911],[64914,64967],[65008,65020],[65136,65140],[65142,65276],[67584,67589],[67592,67592],[67594,67637],[67639,67640],[67644,67644],[67647,67669],[67671,67679],[67840,67867],[67872,67897],[67903,67903],[67968,68023],[68030,68031],[68096,68096],[68112,68115],[68117,68119],[68121,68147],[68160,68167],[68176,68184],[68192,68223],[68352,68405],[68416,68437],[68440,68466],[68472,68479],[68608,68680],[126464,126467],[126469,126495],[126497,126498],[126500,126500],[126503,126503],[126505,126514],[126516,126519],[126521,126521],[126523,126523],[126530,126530],[126535,126535],[126537,126537],[126539,126539],[126541,126543],[126545,126546],[126548,126548],[126551,126551],[126553,126553],[126555,126555],[126557,126557],[126559,126559],[126561,126562],[126564,126564],[126567,126570],[126572,126578],[126580,126583],[126585,126588],[126590,126590],[126592,126601],[126603,126619],[126625,126627],[126629,126633],[126635,126651],[1114109,1114109]];function p(t){for(var e=0;e<c.length;e++){var i=c[e];if(t>=i[0]&&t<=i[1])return!0}return!1}function f(t,e){for(var i=e.childNodes.length-1;i>=0;i--)t.push(e.childNodes[i])}function u(t){if(!t||!t.length)return null;var e=t.pop(),i=e.textContent||e.innerText;if(i){var r=i.match(/^.*(\n|\r)/);return r?(t.length=0,r[0]):i}return"ruby"===e.tagName?u(t):e.childNodes?(f(t,e),u(t)):void 0}function d(t){if(!t||!t.childNodes)return"ltr";var e,i=[];for(f(i,t);e=u(i);)for(var r=0;r<e.length;r++)if(p(e.charCodeAt(r)))return"rtl";return"ltr"}function g(){}function m(t,e){g.call(this),this.cue=e,this.cueDiv=h(t,e.text),this.cueDiv.className="jw-text-track-cue jw-reset";var i="horizontal-tb";/^(lr|rl)$/.test(e.vertical)&&(i="vertical-"+e.vertical);var r={textShadow:"",position:"relative",paddingLeft:0,paddingRight:0,left:0,top:0,bottom:0,display:"inline","white-space":"pre",writingMode:i,unicodeBidi:"plaintext"};this.applyStyles(r,this.cueDiv),this.div=t.document.createElement("div"),r={textAlign:"middle"===e.align?"center":e.align,whiteSpace:"pre-line",position:"absolute",direction:d(this.cueDiv),writingMode:i,unicodeBidi:"plaintext"},this.applyStyles(r),this.div.appendChild(this.cueDiv);var n=0,o="";switch(e.align){case"start":case"left":n=e.position;break;case"middle":case"center":n="auto"===e.position?50:e.position,o=e.vertical?"translateY(-50%)":"translateX(-50%)";break;case"end":case"right":n="auto"===e.position?100:e.position,o=e.vertical?"translateY(-100%)":"translateX(-100%)"}n=Math.max(Math.min(100,n),0),e.vertical?this.applyStyles({top:this.formatStyle(n,"%"),height:this.formatStyle(e.size,"%"),transform:o}):this.applyStyles({left:this.formatStyle(n,"%"),transform:o}),this.move=function(t){this.applyStyles({top:this.formatStyle(t.top,"px"),bottom:this.formatStyle(t.bottom,"px"),left:this.formatStyle(t.left,"px"),paddingRight:this.formatStyle(t.right,"px"),height:this.formatStyle(t.height,"px"),transform:""})}}function v(t){var e,i,r,n,o=t.div;if(o){i=o.offsetHeight,r=o.offsetWidth,n=o.offsetTop;var a=o.firstChild,s=a&&a.getClientRects&&a.getClientRects();t=o.getBoundingClientRect(),e=s?Math.max(s[0]&&s[0].height||0,t.height/s.length):0}this.left=t.left,this.right=t.right,this.top=t.top||n,this.height=t.height||i,this.bottom=t.bottom||n+(t.height||i),this.width=t.width||r,this.lineHeight=void 0!==e?e:t.lineHeight,this.width=Math.ceil(this.width+1)}function y(t,e,i,r,n){var o=new v(e),a=e.cue,s=function(t){if("number"==typeof t.line&&(t.snapToLines||t.line>=0&&t.line<=100))return t.line;if(!t.track||!t.track.textTrackList||!t.track.textTrackList.mediaElement)return-1;for(var e=t.track,i=e.textTrackList,r=0,n=0;n<i.length&&i[n]!==e;n++)"showing"===i[n].mode&&r++;return-1*++r}(a),l=[];if(a.snapToLines){var h;switch(a.vertical){case"":l=["+y","-y"],h="height";break;case"rl":l=["+x","-x"],h="width";break;case"lr":l=["-x","+x"],h="width"}var c=o.lineHeight,p=Math.floor(i[h]/c);s=Math.min(s,p-n);var f=c*Math.round(s),u=i[h]+c,d=l[0];if(Math.abs(f)>u&&(f=f<0?-1:1,f*=Math.ceil(u/c)*c),s<0)f+=a.vertical?i.width:i.height,f-=n*c,l=l.slice().reverse();f-=n,o.move(d,f)}else{var g=o.lineHeight/i.height*100;switch(a.lineAlign){case"middle":s-=g/2;break;case"end":s-=g}switch(a.vertical){case"":e.applyStyles({top:e.formatStyle(s,"%")});break;case"rl":e.applyStyles({left:e.formatStyle(s,"%")});break;case"lr":e.applyStyles({paddingRight:e.formatStyle(s,"%")})}l=["+y","-x","+x","-y"],o=new v(e)}var m=function t(e,n){for(var o,a,s=arguments.length>2&&void 0!==arguments[2]?arguments[2]:0,l=new v(e),h=0,c=0;c<n.length;c++){for(;e.overlapsOppositeAxis(i,n[c])||e.within(i)&&e.overlapsAny(r);)e.move(n[c]);if(e.within(i))return e;var p=e.intersectPercentage(i);h<=p&&(o=new v(e),h=p,a=n[c]),e=new v(l)}var f=o||l;return a&&0===s?t(f,-1===a.indexOf("y")?["-y","+y"]:["-x","+x"],s+1):f}(o,l);e.move(m.toCSSCompatValues(i))}function b(){}g.prototype.applyStyles=function(t,e){for(var i in e=e||this.div,t)t.hasOwnProperty(i)&&(e.style[i]=t[i])},g.prototype.formatStyle=function(t,e){return 0===t?0:t+e},m.prototype=Object.create(g.prototype),m.prototype.constructor=m,v.prototype.move=function(t,e){switch(e=void 0!==e?e:this.lineHeight,t){case"+x":this.left+=e,this.right+=e;break;case"-x":this.left-=e,this.right-=e;break;case"+y":this.top+=e,this.bottom+=e;break;case"-y":this.top-=e,this.bottom-=e}},v.prototype.overlaps=function(t){return this.left<t.right&&this.right>t.left&&this.top<t.bottom&&this.bottom>t.top},v.prototype.overlapsAny=function(t){for(var e=0;e<t.length;e++)if(this.overlaps(t[e]))return!0;return!1},v.prototype.within=function(t){return this.top>=t.top&&this.bottom<=t.bottom&&this.left>=t.left&&this.right<=t.right},v.prototype.overlapsOppositeAxis=function(t,e){switch(e){case"+x":return this.left<t.left;case"-x":return this.right>t.right;case"+y":return this.top<t.top;case"-y":return this.bottom>t.bottom}},v.prototype.intersectPercentage=function(t){return Math.max(0,Math.min(this.right,t.right)-Math.max(this.left,t.left))*Math.max(0,Math.min(this.bottom,t.bottom)-Math.max(this.top,t.top))/(this.height*this.width)},v.prototype.toCSSCompatValues=function(t){return{top:this.top-t.top,bottom:t.bottom-this.bottom,left:this.left-t.left,paddingRight:t.right-this.right,height:this.height,width:this.width}},v.getSimpleBoxPosition=function(t){var e=t.div?t.div.offsetHeight:t.tagName?t.offsetHeight:0,i=t.div?t.div.offsetWidth:t.tagName?t.offsetWidth:0,r=t.div?t.div.offsetTop:t.tagName?t.offsetTop:0,n=(t=t.div?t.div.getBoundingClientRect():t.tagName?t.getBoundingClientRect():t).height||e;return{left:t.left,right:t.right,top:t.top||r,height:n,bottom:t.bottom||r+n,width:t.width||i}},b.StringDecoder=function(){return{decode:function(t){if(!t)return"";if("string"!=typeof t)throw new Error("Error - expected string data.");return decodeURIComponent(encodeURIComponent(t))}}},b.convertCueToDOMTree=function(t,e){return t&&e?h(t,e):null};b.processCues=function(t,e,i,r){if(!t||!e||!i)return null;for(;i.firstChild;)i.removeChild(i.firstChild);if(!e.length)return null;var n=t.document.createElement("div");if(n.className="jw-text-track-container jw-reset",n.style.position="absolute",n.style.left="0",n.style.right="0",n.style.top="0",n.style.bottom="0",n.style.margin="1.5%",i.appendChild(n),function(t){for(var e=0;e<t.length;e++)if(t[e].hasBeenReset||!t[e].displayState)return!0;return!1}(e)||r){var o=[],a=v.getSimpleBoxPosition(n),s=e.reduce((function(t,e){return t+e.text.split("\n").length}),0);!function(){for(var i=0;i<e.length;i++){var r=e[i],l=new m(t,r);l.div.className="jw-text-track-display jw-reset",n.appendChild(l.div),y(0,l,a,o,s),s-=r.text.split("\n").length,r.displayState=l.div,o.push(v.getSimpleBoxPosition(l))}}()}else for(var l=0;l<e.length;l++)n.appendChild(e[l].displayState)};var w=window.WebVTT;w||(window.WebVTT=w=b),e.default=w}}]);