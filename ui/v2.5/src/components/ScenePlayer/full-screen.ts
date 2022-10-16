import videojs, { VideoJsPlayer } from "video.js";

class fullscreenLock extends videojs.getPlugin("plugin") {
  constructor(player: VideoJsPlayer, options: any) {
    super(player, options);

    player.on('fullscreenchange', e => {
        if (videojs.browser.IS_ANDROID || videojs.browser.IS_IOS) {
          if (player.isFullscreen()) {
            window.screen.orientation.lock('landscape');
          }
        }
      });
  }
}

// Register the plugin with video.js.
videojs.registerPlugin("fullscreenLock", fullscreenLock);

declare module "video.js" {
  export interface VideoJsPlayer {
    fullscreenLock: () => fullscreenLock;
  }
}

export default fullscreenLock;
