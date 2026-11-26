import videojs, { VideoJsPlayer } from "video.js";
import type { AbLoopPluginApi } from "./util";

function getAbLoopApi(player: VideoJsPlayer) {
  return player.abLoopPlugin as unknown as AbLoopPluginApi | undefined;
}

interface AbLoopRangeOptions {
  /**
   * Whether to show the range highlight at all (mirrors showAbLoopControls).
   * @default false
   */
  enabled?: boolean;
}

// Highlights the active AB-loop range on the progress bar, so there's some
// visual answer to "what time range is currently selected" beyond having
// to remember what you last tapped Start/End at. Also drops a thin marker
// at the start and/or end point individually, shown as soon as that point
// is set even if the other one isn't - see refresh() below. Appended
// directly into .vjs-progress-holder (not .vjs-progress-control, like
// markers.ts's range markers) so it shares that element's own coordinate
// space - no need to separately compensate for the holder's 15px side
// margins the way markers.ts has to when positioning from the outer
// container instead.
class AbLoopRangePlugin extends videojs.getPlugin("plugin") {
  private rangeEl?: HTMLDivElement;
  private startMarkerEl?: HTMLDivElement;
  private endMarkerEl?: HTMLDivElement;

  constructor(player: VideoJsPlayer, options?: AbLoopRangeOptions) {
    super(player, options);
    if (!options?.enabled) return;

    player.on("ready", () => {
      const holder = player.el().querySelector(".vjs-progress-holder");
      if (!holder) return;

      const el = document.createElement("div");
      el.className = "vjs-ab-loop-range";
      el.style.display = "none";
      holder.appendChild(el);
      this.rangeEl = el;

      const startMarker = document.createElement("div");
      startMarker.className = "vjs-ab-loop-marker vjs-ab-loop-marker-start";
      startMarker.style.display = "none";
      holder.appendChild(startMarker);
      this.startMarkerEl = startMarker;

      const endMarker = document.createElement("div");
      endMarker.className = "vjs-ab-loop-marker vjs-ab-loop-marker-end";
      endMarker.style.display = "none";
      holder.appendChild(endMarker);
      this.endMarkerEl = endMarker;

      const refresh = () => this.refresh();
      player.on("abloopchange", refresh);
      player.on(["loadedmetadata", "durationchange"], refresh);
      refresh();
    });
  }

  private positionMarker(
    el: HTMLDivElement | undefined,
    show: boolean,
    time: number | false | undefined,
    duration: number
  ) {
    if (!el) return;
    if (!show || typeof time !== "number") {
      el.style.display = "none";
      return;
    }
    el.style.display = "block";
    el.style.left = `${(time / duration) * 100}%`;
  }

  private refresh() {
    const opts = getAbLoopApi(this.player)?.getOptions();
    const duration = this.player.duration();
    const hasDuration = Number.isFinite(duration) && duration > 0;

    // the shaded range: the region that's actually looping right now, once
    // enabled - even with neither point explicitly set, this covers the
    // whole video (the vendor plugin's own start:0/end:false defaults),
    // which is exactly what an enabled-but-untouched loop plays
    const showRange = !!opts?.enabled && hasDuration;
    if (this.rangeEl) {
      if (!showRange || !opts) {
        this.rangeEl.style.display = "none";
      } else {
        const start = typeof opts.start === "number" ? opts.start : 0;
        const end = opts.end === false ? duration : opts.end;
        this.rangeEl.style.display = "block";
        this.rangeEl.style.left = `${(start / duration) * 100}%`;
        this.rangeEl.style.width = `${
          (Math.max(0, end - start) / duration) * 100
        }%`;
      }
    }

    // the individual pins - each shown only once that specific bound has
    // been set (not the vendor plugin's 0/false default), independently of
    // whether the other one has been too
    const startIsSet = typeof opts?.start === "number" && opts.start > 0;
    const endIsSet = typeof opts?.end === "number";
    this.positionMarker(
      this.startMarkerEl,
      !!opts?.enabled && hasDuration && startIsSet,
      opts?.start,
      duration
    );
    this.positionMarker(
      this.endMarkerEl,
      !!opts?.enabled && hasDuration && endIsSet,
      opts?.end,
      duration
    );
  }
}

// Register the plugin with video.js.
videojs.registerPlugin("abLoopRange", AbLoopRangePlugin);

declare module "video.js" {
  interface VideoJsPlayer {
    abLoopRange: () => AbLoopRangePlugin;
  }
  interface VideoJsPlayerPluginOptions {
    abLoopRange?: AbLoopRangeOptions;
  }
}

export default AbLoopRangePlugin;
