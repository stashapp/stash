import videojs from "video.js";
import type VideoJsPlayer from "video.js/dist/types/player";
import type Plugin from "video.js/dist/types/plugin";

class MediaSessionPlugin extends (videojs.getPlugin(
  "plugin"
) as typeof Plugin) {
  constructor(player: VideoJsPlayer) {
    super(player);

    player.ready(() => {
      player.addClass("vjs-media-session");
      this.setActionHandlers();
    });

    player.on("play", () => {
      this.updatePlaybackState();
    });

    player.on("pause", () => {
      this.updatePlaybackState();
    });
    this.updatePlaybackState();
  }

  // manually set poster since it's only set on useEffect
  public setMetadata(title: string, artist: string, poster: string): void {
    if ("mediaSession" in navigator) {
      navigator.mediaSession.metadata = new MediaMetadata({
        title,
        artist,
        artwork: [
          {
            src: poster || this.player.poster() || "",
            type: "image/jpeg",
          },
        ],
      });
    }
  }

  private updatePlaybackState(): void {
    if ("mediaSession" in navigator) {
      const playbackState = this.player.paused() ? "paused" : "playing";
      navigator.mediaSession.playbackState = playbackState;
    }
  }

  private setActionHandlers(): void {
    // method initialization
    navigator.mediaSession.setActionHandler("play", () => {
      this.player.play();
    });
    navigator.mediaSession.setActionHandler("pause", () => {
      this.player.pause();
    });
    navigator.mediaSession.setActionHandler("nexttrack", () => {
      this.player.skipButtons()?.handleForward();
    });
    navigator.mediaSession.setActionHandler("previoustrack", () => {
      this.player.skipButtons()?.handleBackward();
    });
  }
}

videojs.registerPlugin("mediaSession", MediaSessionPlugin);

declare module "video.js/dist/types/player" {
  export default interface Player {
    mediaSession: () => MediaSessionPlugin;
  }
}

export default MediaSessionPlugin;
