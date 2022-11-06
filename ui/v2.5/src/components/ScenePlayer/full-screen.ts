import videojs, { VideoJsPlayer } from "video.js";

class fullscreenLock extends videojs.getPlugin("plugin") {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  constructor(player: VideoJsPlayer, options: any) {
    super(player, options);

    player.on("fullscreenchange", () => {
      if (videojs.browser.IS_ANDROID || videojs.browser.IS_IOS) {
        if (player.isFullscreen()) {
          window.screen.orientation.lock("landscape");
        }
      }
    });
  }
}

// Register the plugin with video.js.
videojs.registerPlugin("fullscreenLock", fullscreenLock);

declare module "video.js" {
  /* eslint-disable-next-line @typescript-eslint/naming-convention */
  export interface VideoJsPlayer {
    fullscreenLock: () => fullscreenLock;
  }
}

export default fullscreenLock;
