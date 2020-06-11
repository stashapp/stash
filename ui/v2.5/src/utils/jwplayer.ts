const playerID = "main-jwplayer";
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const getPlayer = () => (window as any).jwplayer(playerID);

// eslint-disable-next-line @typescript-eslint/no-explicit-any
declare const Modernizr: any;

Modernizr.addTest("video.mkv", () => {
  // the following code does not return true even in browsers that support matroska (Chrome, Edge)

  // let elem = document.createElement("video");
  // return elem.canPlayType("video/x-matroska").replace(/^no$/,'');

  // HACK - detect mkv support based on browser
  // known mkv support is in Chrome and Edge (Chromium)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const isChromium = !!(window as any).chrome;
  return isChromium;
});

const getSupportedFormats = () => {
  const { h264, h265, webm, vp9, mkv } = Modernizr.video;
  const ret: string[] = [];

  if (h264) ret.push("h264");
  if (h265) ret.push("hevc");
  if (webm) ret.push("vp8");
  if (vp9) ret.push("vp9");
  if (mkv) ret.push("mkv");

  // not supported on the backend
  // if (video.ogg) ret.push("ogg");
  // if (video.hls) ret.push("hls");
  // if (video.av1) ret.push("av1");

  return ret;
};

export default {
  playerID,
  getPlayer,
  getSupportedFormats,
};
