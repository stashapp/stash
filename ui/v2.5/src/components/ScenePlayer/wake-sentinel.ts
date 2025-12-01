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
        await this.acquireWakeLock();
      }
    });

    // acquire wake lock on ready and play
    player.ready(async () => {
      player.addClass("vjs-wake-sentinel");
      await this.acquireWakeLock(true);
    });
    player.on("play", () => this.acquireWakeLock());

    // release wake lock on pause, dispose and end
    player.on("pause", () => this.releaseWakeLock());
    player.on("dispose", () => this.releaseWakeLock());
    player.on("ended", () => this.releaseWakeLock());
  }

  private async releaseWakeLock(): Promise<void> {
    this.wakeLock?.release().then(() => (this.wakeLock = null));
  }

  private async acquireWakeLock(log = false): Promise<void> {
    // if wake lock failed, don't even try
    if (this.wakeLockFail) return;
    // check for wake lock on startup
    if ("wakeLock" in navigator) {
      try {
        this.wakeLock = await navigator.wakeLock.request("screen");
      } catch (err) {
        if (log) console.error("Failed to obtain Screen Wake Lock:", err);
        this.wakeLockFail = true;
      }
    } else {
      if (log) {
        console.warn(
          "Screen Wake Lock API not supported. Secure context (https or localhost) and modern browser required."
        );
      }
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
