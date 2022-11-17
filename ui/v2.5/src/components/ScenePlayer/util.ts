import videojs from "video.js";

export const VIDEO_PLAYER_ID = "VideoJsPlayer";

export const getPlayerPosition = () =>
  videojs.getPlayer(VIDEO_PLAYER_ID)?.currentTime();
