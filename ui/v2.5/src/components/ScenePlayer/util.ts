import videojs from "video.js";

export const VIDEO_PLAYER_ID = "VideoJsPlayer";

export const getPlayer = () => videojs.getPlayer(VIDEO_PLAYER_ID);

export const getPlayerPosition = () => getPlayer()?.currentTime();
