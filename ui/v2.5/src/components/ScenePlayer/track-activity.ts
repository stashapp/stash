import videojs, { VideoJsPlayer } from "video.js";

const intervalSeconds = 1; // check every second
const sendInterval = 10; // send every 10 seconds

class TrackActivityPlugin extends videojs.getPlugin("plugin") {
  totalPlayDuration = 0;
  currentPlayDuration = 0;
  minimumPlayPercent = 0;
  incrementPlayCount: () => Promise<void> = () => {
    return Promise.resolve();
  };
  saveActivity: (resumeTime: number, playDuration: number) => Promise<void> =
    () => {
      return Promise.resolve();
    };

  private enabled = false;
  private playCountIncremented = false;
  private intervalID: number | undefined;

  private lastResumeTime = 0;
  private lastDuration = 0;

  constructor(player: VideoJsPlayer) {
    super(player);

    player.on("playing", () => {
      this.start();
    });

    player.on("waiting", () => {
      this.stop();
    });

    player.on("stalled", () => {
      this.stop();
    });

    player.on("pause", () => {
      this.stop();
    });

    player.on("dispose", () => {
      this.stop();
    });

    player.on("ended", () => {
      this.stop();
    });
  }

  private start() {
    if (this.enabled && !this.intervalID) {
      this.intervalID = window.setInterval(() => {
        this.intervalHandler();
      }, intervalSeconds * 1000);
      this.lastResumeTime = this.player.currentTime();
      this.lastDuration = this.player.duration();
    }
  }

  private stop() {
    if (this.intervalID) {
      window.clearInterval(this.intervalID);
      this.intervalID = undefined;
      this.sendActivity();
    }
  }

  reset() {
    this.stop();
    this.totalPlayDuration = 0;
    this.currentPlayDuration = 0;
    this.playCountIncremented = false;
  }

  setEnabled(enabled: boolean) {
    this.enabled = enabled;
    if (!enabled) {
      this.stop();
    } else if (!this.player.paused()) {
      this.start();
    }
  }

  private intervalHandler() {
    if (!this.enabled || !this.player) return;

    this.lastResumeTime = this.player.currentTime();
    this.lastDuration = this.player.duration();

    this.totalPlayDuration += intervalSeconds;
    this.currentPlayDuration += intervalSeconds;
    if (this.totalPlayDuration % sendInterval === 0) {
      this.sendActivity();
    }
  }

  private sendActivity() {
    if (!this.enabled) return;

    if (this.totalPlayDuration > 0) {
      let resumeTime = this.player?.currentTime() ?? this.lastResumeTime;
      const videoDuration = this.player?.duration() ?? this.lastDuration;
      const percentCompleted = (100 / videoDuration) * resumeTime;
      const percentPlayed = (100 / videoDuration) * this.totalPlayDuration;

      if (
        !this.playCountIncremented &&
        percentPlayed >= this.minimumPlayPercent
      ) {
        this.incrementPlayCount();
        this.playCountIncremented = true;
      }

      // if video is 98% or more complete then reset resume_time
      if (percentCompleted >= 98) {
        resumeTime = 0;
      }

      this.saveActivity(resumeTime, this.currentPlayDuration);
      this.currentPlayDuration = 0;
    }
  }
}

// Register the plugin with video.js.
videojs.registerPlugin("trackActivity", TrackActivityPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    trackActivity: () => TrackActivityPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    trackActivity?: {};
  }
}

export default TrackActivityPlugin;
