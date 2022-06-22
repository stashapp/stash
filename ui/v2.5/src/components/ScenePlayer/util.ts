import VideoJS from "video.js";

export const VIDEO_PLAYER_ID = "VideoJsPlayer";

export const getPlayerPosition = () =>
  VideoJS.getPlayer(VIDEO_PLAYER_ID).currentTime();
