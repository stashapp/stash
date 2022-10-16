import videojs, { VideoJsPlayer } from "video.js";

class backgroundPlaybackDisabler extends videojs.getPlugin("plugin") {
  constructor(player: VideoJsPlayer, options: any) {
    super(player, options);

    window.addEventListener("visibilitychange", this.handleVisibilityChange);
  }

  handleVisibilityChange = () => {
    if (document.hidden) {
      this.player.pause();
    }
  }

  dispose() {
    super.dispose();
    window.removeEventListener("visibilitychange", this.handleVisibilityChange);
  }
}

// Register the plugin with video.js.
videojs.registerPlugin(
  "backgroundPlaybackDisabler",
  backgroundPlaybackDisabler
);

declare module "video.js" {
  export interface VideoJsPlayer {
    backgroundPlaybackDisabler: () => backgroundPlaybackDisabler;
  }
}

export default backgroundPlaybackDisabler;
