import videojs, { VideoJsPlayer } from "video.js";

// Most of this was copied directly from the video.js source code.
// This prevents the player from requesting video data until the mouse is released.
const seekHandler = function (this: VideoJsPlayer) {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  const SeekBar = videojs.getComponent("SeekBar") as any;

  if (!SeekBar.__super__ || !SeekBar.__super__.__seekHandlerInit) {
    SeekBar.__super__ = {
      __seekHandlerInit: true,
      handleMouseMove: SeekBar.prototype.handleMouseMove,
      handleMouseUp: SeekBar.prototype.handleMouseUp,
    };

    SeekBar.prototype.handleMouseMove = function (
      event: videojs.EventTarget.Event,
      seek: boolean
    ) {
      if (!videojs.dom.isSingleLeftClick(event)) {
        return;
      }
      let newTime;
      const distance = this.calculateDistance(event);
      /* eslint-disable-next-line prefer-destructuring */
      const liveTracker = this.player_.liveTracker;

      if (!liveTracker || !liveTracker.isLive()) {
        newTime = distance * this.player_.duration();

        // Don't let video end while scrubbing.
        if (newTime === this.player_.duration()) {
          newTime = newTime - 0.1;
        }
      } else {
        if (distance >= 0.99) {
          liveTracker.seekToLiveEdge();
          return;
        }
        const seekableStart = liveTracker.seekableStart();
        const seekableEnd = liveTracker.liveCurrentTime();

        newTime = seekableStart + distance * liveTracker.liveWindow();

        // Don't let video end while scrubbing.
        if (newTime >= seekableEnd) {
          newTime = seekableEnd;
        }

        // Compensate for precision differences so that currentTime is not less
        // than seekable start
        if (newTime <= seekableStart) {
          newTime = seekableStart + 0.1;
        }

        // On android seekableEnd can be Infinity sometimes,
        // this will cause newTime to be Infinity, which is
        // not a valid currentTime.
        if (newTime === Infinity) {
          return;
        }
      }

      // EDIT: Start
      if (seek) {
        // Set new time (tell player to seek to new time)
        this.userSeek_(newTime);
      } else {
        // hijack the player's time setting
        this.player_.cache_.currentTime = newTime;
        this.player_.trigger({
          type: "timeupdate",
          target: this,
          manuallyTriggered: true,
        });
      }
      // EDIT: End
    };

    SeekBar.prototype.handleMouseUp = function (
      event: videojs.EventTarget.Event
    ) {
      Object.getPrototypeOf(SeekBar).prototype.handleMouseUp.call(this, event); // Stop event propagation to prevent double fire in progress-control.js

      // Stop event propagation to prevent double fire in progress-control.js
      if (event) {
        event.stopPropagation();
      }
      // EDIT: Set new time (tell player to actually seek to new time)
      this.handleMouseMove(event, true); // edit
      this.player_.scrubbing(false);

      /**
       * Trigger timeupdate because we're done seeking and the time has changed.
       * This is particularly useful for if the player is paused to time the time displays.
       *
       * @event Tech#timeupdate
       * @type {EventTarget~Event}
       */
      this.player_.trigger({
        type: "timeupdate",
        target: this,
        manuallyTriggered: true,
      });
      if (this.videoWasPlaying) {
        // EDIT: Start
        // silencePromise(this.player_.play());
        this.player_.play().catch(() => {});
        // EDIT: End
      } else {
        // We're done seeking and the time has changed.
        // If the player is paused, make sure we display the correct time on the seek bar.
        this.update_();
      }
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
