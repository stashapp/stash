import videojs, { VideoJsPlayer } from "video.js";
import ChromecastSessionManager from "@silvermine/videojs-chromecast/src/js/chromecast/ChromecastSessionManager";

interface ITechOptions {
  playerId: string;
  source?: videojs.Tech.SourceObject;
}

type PlayerWithManager = VideoJsPlayer & {
  chromecastSessionManager?: ChromecastSessionManager;
};

// The plugin is unmaintained, so its Chromecast tech is adjusted here instead of being patched.
// Call after the plugin has registered its tech.
export function registerStashChromecastTech() {
  // biome-ignore lint/suspicious/noExplicitAny: the plugin ships no types
  const ChromecastTech = videojs.getTech("chromecast") as any;

  class StashChromecastTech extends ChromecastTech {
    constructor(options: ITechOptions) {
      // A cast that outlives the previous player makes Video.js pick this tech while the player is built, before the plugin has created the session manager.
      const player = videojs(options.playerId) as PlayerWithManager;
      player.chromecastSessionManager ??= new ChromecastSessionManager(player);

      super(options);
    }

    // Video.js only calls setSource on a tech that defines it itself, and sets the source through src otherwise, which this tech ignores.
    setSource(source: videojs.Tech.SourceObject) {
      super.setSource(source);
    }

    // With a cast in progress Video.js picks this tech before the player has a source, and sets one once the tech is ready.
    _playSource(source: videojs.Tech.SourceObject | undefined, time: number) {
      if (source) {
        super._playSource(source, time);
      } else {
        this.triggerReady();
      }
    }

    // The device answers after the player may be gone.
    trigger(...args: unknown[]) {
      if (!this.isDisposed()) {
        super.trigger(...args);
      }
    }

    // Video.js forwards the player's loop setting to whichever tech is active, and the device has none.
    loop() {
      return false;
    }

    setLoop() {}

    // Video.js holds tech events back while the tech is seeking, until a seeked event this tech never fires.
    seeking() {
      return false;
    }
  }

  videojs.registerTech("chromecast", StashChromecastTech);
}
