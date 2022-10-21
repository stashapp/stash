import videojs, { VideoJsPlayer } from "video.js";

class backgroundPlaybackDisabler extends videojs.getPlugin("plugin") {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  constructor(player: VideoJsPlayer, options: any) {
    super(player, options);

    window.addEventListener("visibilitychange", this.handleVisibilityChange);
  }

  handleVisibilityChange = () => {
    if (document.hidden) {
      this.player.pause();
    }
  };

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
  /* eslint-disable-next-line @typescript-eslint/naming-convention */
  export interface VideoJsPlayer {
    backgroundPlaybackDisabler: () => backgroundPlaybackDisabler;
  }
}

export default backgroundPlaybackDisabler;
