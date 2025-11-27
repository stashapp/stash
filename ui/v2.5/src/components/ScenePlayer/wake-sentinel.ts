import videojs, { VideoJsPlayer } from "video.js";

class WakeSentinelPlugin extends videojs.getPlugin("plugin") {
  public wakeLock: WakeLockSentinel | null = null;
  public wakeLockFail: boolean = false;
  constructor(player: VideoJsPlayer) {
    super(player);

    // listen for visibility change events
    document.addEventListener("visibilitychange", async () => {
      if (document.visibilityState === "visible") {
        // reacquire the wake lock when the page becomes visible
        await this.acquireWakeLock(false);
      }
    });

    player.ready(async () => {
      player.addClass("vjs-wake-sentinel");
      await this.acquireWakeLock(true);
    });

    player.on("play", () => {
      this.acquireWakeLock(false);
    });

    player.on("pause", () => {
      this.wakeLock?.release();
    });
  }

  private async acquireWakeLock(log = false): Promise<void> {
    // if wake lock failed, don't even try
    if (this.wakeLockFail) return;
    // check for wake lock on startup
    if ("wakeLock" in navigator) {
      try {
        this.wakeLock = await navigator.wakeLock.request("screen");
        if (log) console.log("Screen Wake Lock is active!");
      } catch (err) {
        if (log) console.error("Failed to obtain Screen Wake Lock:", err);
        this.wakeLockFail = true;
      }
    } else {
      if (log) console.warn("Screen Wake Lock API not supported.");
      this.wakeLockFail = true;
    }
  }
}

videojs.registerPlugin("wakeSentinel", WakeSentinelPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    wakeSentinel: () => WakeSentinelPlugin;
  }
}

export default WakeSentinelPlugin;
