const playerID = "main-jwplayer";
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const getPlayer = () => (window as any).jwplayer(playerID);

declare const Modernizr:any;

Modernizr.addTest("video.mkv", () => {
  // the following code does not return true even in browsers that support matroska (Chrome, Edge)

  // let elem = document.createElement("video");
  // return elem.canPlayType("video/x-matroska").replace(/^no$/,'');

  // HACK - detect mkv support based on browser
  // known mkv support is in Chrome and Edge (Chromium)
  const isChromium = !!(window as any).chrome;
  return isChromium;
});

const getSupportedFormats = () => {
  const video = Modernizr.video;
  const ret : string[] = [];

  if (video.h264) ret.push("h264");
  if (video.h265) ret.push("hevc");
  if (video.webm) ret.push("vp8");
  if (video.vp9) ret.push("vp9");
  if (video.mkv) ret.push("mkv");
  
  // not supported on the backend
  // if (video.ogg) ret.push("ogg");
  // if (video.hls) ret.push("hls");
  // if (video.av1) ret.push("av1");

  return ret;
}

export default {
  playerID,
  getPlayer,
  getSupportedFormats,
};
