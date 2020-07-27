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
(window.webpackJsonpjwplayer=window.webpackJsonpjwplayer||[]).push([[10],{97:function(t,e,r){"use strict";r.r(e);var n=r(42),i=r(67),s=/^(\d+):(\d{2})(:\d{2})?\.(\d{3})/,a=/^-?\d+$/,u=/\r\n|\n/,o=/^NOTE($|[ \t])/,c=/^[^\sa-zA-Z-]+/,l=/:/,f=/\s/,h=/^\s+/,g=/-->/,d=/^WEBVTT([ \t].*)?$/,p=function(t,e){this.window=t,this.state="INITIAL",this.buffer="",this.decoder=e||new b,this.regionList=[],this.maxCueBatch=1e3};function b(){return{decode:function(t){if(!t)return"";if("string"!=typeof t)throw new Error("Error - expected string data.");return decodeURIComponent(encodeURIComponent(t))}}}function v(){this.values=Object.create(null)}v.prototype={set:function(t,e){this.get(t)||""===e||(this.values[t]=e)},get:function(t,e,r){return r?this.has(t)?this.values[t]:e[r]:this.has(t)?this.values[t]:e},has:function(t){return t in this.values},alt:function(t,e,r){for(var n=0;n<r.length;++n)if(e===r[n]){this.set(t,e);break}},integer:function(t,e){a.test(e)&&this.set(t,parseInt(e,10))},percent:function(t,e){return(e=parseFloat(e))>=0&&e<=100&&(this.set(t,e),!0)}};var E=new i.a(0,0,0),w="middle"===E.align?"middle":"center";function T(t,e,r){var n=t;function i(){var e=function(t){function e(t,e,r,n){return 3600*(0|t)+60*(0|e)+(0|r)+(0|n)/1e3}var r=t.match(s);return r?r[3]?e(r[1],r[2],r[3].replace(":",""),r[4]):r[1]>59?e(r[1],r[2],0,r[4]):e(0,r[1],r[2],r[4]):null}(t);if(null===e)throw new Error("Malformed timestamp: "+n);return t=t.replace(c,""),e}function a(){t=t.replace(h,"")}if(a(),e.startTime=i(),a(),"--\x3e"!==t.substr(0,3))throw new Error("Malformed time stamp (time stamps must be separated by '--\x3e'): "+n);t=t.substr(3),a(),e.endTime=i(),a(),function(t,e){var n=new v;!function(t,e,r,n){for(var i=n?t.split(n):[t],s=0;s<=i.length;s+=1)if("string"==typeof i[s]){var a=i[s].split(r);if(2===a.length)e(a[0],a[1])}}(t,(function(t,e){switch(t){case"region":for(var i=r.length-1;i>=0;i--)if(r[i].id===e){n.set(t,r[i].region);break}break;case"vertical":n.alt(t,e,["rl","lr"]);break;case"line":var s=e.split(","),a=s[0];n.integer(t,a),n.percent(t,a)&&n.set("snapToLines",!1),n.alt(t,a,["auto"]),2===s.length&&n.alt("lineAlign",s[1],["start",w,"end"]);break;case"position":var u=e.split(",");n.percent(t,u[0]),2===u.length&&n.alt("positionAlign",u[1],["start",w,"end","line-left","line-right","auto"]);break;case"size":n.percent(t,e);break;case"align":n.alt(t,e,["start",w,"end","left","right"])}}),l,f),e.region=n.get("region",null),e.vertical=n.get("vertical","");var i=n.get("line","auto");"auto"===i&&-1===E.line&&(i=-1),e.line=i,e.lineAlign=n.get("lineAlign","start"),e.snapToLines=n.get("snapToLines",!0),e.size=n.get("size",100),e.align=n.get("align",w);var s=n.get("position","auto");"auto"===s&&50===E.position&&(s="start"===e.align||"left"===e.align?0:"end"===e.align||"right"===e.align?100:50),e.position=s}(t,e)}p.prototype={parse:function(t,e){var r,s=this;function a(){for(var t=s.buffer,e=0;e<t.length&&"\r"!==t[e]&&"\n"!==t[e];)++e;var r=t.substr(0,e);return"\r"===t[e]&&++e,"\n"===t[e]&&++e,s.buffer=t.substr(e),r}function c(){"CUETEXT"===s.state&&s.cue&&s.oncue&&s.oncue(s.cue),s.cue=null,s.state="INITIAL"===s.state?"BADWEBVTT":"BADCUE"}t&&(s.buffer+=s.decoder.decode(t,{stream:!0}));try{if("INITIAL"===s.state){if(!u.test(s.buffer))return this;var f=(r=a()).match(d);if(!f||!f[0])throw new Error("Malformed WebVTT signature.");s.state="HEADER"}}catch(t){return c(),this}var h=!1,p=0;!function t(){try{for(;s.buffer&&p<=s.maxCueBatch;){if(!u.test(s.buffer))return s.flush(),this;switch(h?h=!1:r=a(),s.state){case"HEADER":l.test(r)||r||(s.state="ID");break;case"NOTE":r||(s.state="ID");break;case"ID":if(o.test(r)){s.state="NOTE";break}if(!r)break;if(s.cue=new i.a(0,0,""),s.state="CUE",!g.test(r)){s.cue.id=r;break}case"CUE":try{T(r,s.cue,s.regionList)}catch(t){s.cue=null,s.state="BADCUE";break}s.state="CUETEXT";break;case"CUETEXT":var f=g.test(r);if(!r||f&&(h=!0)){s.oncue&&(p+=1,s.oncue(s.cue)),s.cue=null,s.state="ID";break}s.cue.text&&(s.cue.text+="\n"),s.cue.text+=r;break;case"BADCUE":r||(s.state="ID")}}if(p=0,s.buffer)Object(n.b)(t);else if(!e)return s.flush(),this}catch(t){return c(),this}}()},flush:function(){try{if(this.buffer+=this.decoder.decode(),(this.cue||"HEADER"===this.state)&&(this.buffer+="\n\n",this.parse(void 0,!0)),"INITIAL"===this.state)throw new Error("Malformed WebVTT signature.")}catch(t){throw t}return this.onflush&&this.onflush(),this}},e.default=p}}]);