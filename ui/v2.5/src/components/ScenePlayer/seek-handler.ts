import videojs, { VideoJsPlayer } from "video.js";

// Most of this was copied directly from the video.js source code.
// This prevents the player from requesting video data until the mouse is released.
const seekHandler = function (this: VideoJsPlayer) {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  const SeekBar = videojs.getComponent("SeekBar") as any;

  if (!SeekBar.__super__ || !SeekBar.__super__.__seekHandlerInit) {
    SeekBar.__super__ = {
      __seekHandlerInit: true,
      userSeek_: SeekBar.prototype.userSeek_,
      handleMouseUp: SeekBar.prototype.handleMouseUp,
    };

    SeekBar.prototype.userSeek_ = function userSeek_(ct: number) {
      if (this.player_.liveTracker && this.player_.liveTracker.isLive()) {
        this.player_.liveTracker.nextSeekedFromUser();
      }

      if (!this.player_.scrubbing()) {
        this.player_.currentTime(ct);
      } else {
        // hijack the player's time setting
        this.player_.cache_.scrubTime = ct; // for thumbnails
        // on timeupdate, the player revises the currentTime cache to the actual time
        this.player_.cache_.currentTime = ct;
        this.player_.trigger({
          type: "timeupdate",
          target: this,
          manuallyTriggered: true,
        });
      }
    };

    SeekBar.prototype.handleMouseUp = function (
      event: videojs.EventTarget.Event
    ) {
      this.player_.scrubbing(false);
      this.handleMouseMove(event); // edit
      SeekBar.__super__.handleMouseUp.apply(this, arguments);
    };
  }
};

// Register the plugin with video.js.
videojs.registerPlugin("seekHandler", seekHandler);

declare module "video.js" {
  /* eslint-disable-next-line @typescript-eslint/naming-convention */
  export interface VideoJsPlayer {
    seekHandler: () => typeof seekHandler;
  }
}

export default seekHandler;
