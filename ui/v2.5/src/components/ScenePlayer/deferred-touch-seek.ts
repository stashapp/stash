import videojs from "video.js";

const SeekBar = videojs.getComponent(
  "SeekBar"
) as unknown as typeof videojs.SeekBar;

// The stock SeekBar seeks live - it calls player.currentTime() (a real
// seek) on every touchmove while dragging, the same code path used for
// mouse dragging. That's essentially free for a local file, but Stash
// scenes are typically streamed (HLS/DASH, or range-requested direct
// play), so sweeping a finger across the bar issues a real segment/range
// load at every intermediate point swept over on the way to wherever the
// finger lands, not just the one it's released on - most noticeable as a
// long stall/lots of loading when scrubbing on mobile.
//
// Defer the actual seek to touchend/release. The fill bar is still moved
// on every touchmove - computed directly from pointer position rather
// than from the player's actual current time - so dragging still looks
// and feels responsive; vtt-thumbnails.ts's hover preview tracks pointer
// position independently of any of this, so it's unaffected either way.
// Mouse dragging (desktop) is untouched, since only touch* events take
// this path.
class DeferredTouchSeekBar extends SeekBar {
  private pendingSeekPercent?: number;

  handleMouseMove(event: videojs.EventTarget.Event, mouseDown?: boolean) {
    if (event.type !== "touchstart" && event.type !== "touchmove") {
      // the public type only declares handleMouseMove(event), but the
      // real implementation also takes the mouseDown flag Slider.
      // handleMouseDown() passes on the initial down event - forward it
      // through so mouse-drag behavior is unchanged
      (
        super.handleMouseMove as (
          event: videojs.EventTarget.Event,
          mouseDown?: boolean
        ) => void
      )(event, mouseDown);
      return;
    }

    const percent = this.calculateDistance(event);
    this.pendingSeekPercent = percent;

    /* biome-ignore lint/suspicious/noExplicitAny: `bar` isn't part of video.js's public Slider/SeekBar types, but is set by the base Slider constructor */
    const bar = (this as any).bar;
    if (bar) {
      (bar.el() as HTMLElement).style.width = `${(percent * 100).toFixed(2)}%`;
    }
    this.el().setAttribute("aria-valuenow", (percent * 100).toFixed(2));
  }

  handleMouseUp(event: videojs.EventTarget.Event) {
    const percent = this.pendingSeekPercent;
    this.pendingSeekPercent = undefined;

    // apply the deferred seek before delegating to the base class, so that
    // by the time it decides whether to resume playback, currentTime
    // already reflects where the drag ended - same end state as the stock
    // continuous-seek behavior, just without the intermediate loads
    if (percent !== undefined) {
      const player = this.player();
      player.currentTime(percent * player.duration());
    }

    super.handleMouseUp(event);
  }
}

// Replace the built-in SeekBar - this is video.js's own supported way to
// customize a default component, and every consumer (ProgressControl,
// keyboard handling, etc.) looks it up by this same registered name.
videojs.registerComponent("SeekBar", DeferredTouchSeekBar);

export default DeferredTouchSeekBar;
