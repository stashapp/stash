import videojs, { VideoJsPlayer } from "video.js";
import CryptoJS from "crypto-js";

export interface IMarker {
  title: string;
  seconds: number;
  end_seconds?: number | null;
  primaryTag: { name: string };
}

interface IMarkersOptions {
  markers?: IMarker[];
}

class MarkersPlugin extends videojs.getPlugin("plugin") {
  private markers: IMarker[] = [];
  private markerDivs: {
    dot?: HTMLDivElement;
    range?: HTMLDivElement;
    containedRanges?: HTMLDivElement[];
  }[] = [];
  private markerTooltip: HTMLElement | null = null;
  private defaultTooltip: HTMLElement | null = null;

  private layerHeight: number = 9;

  private tagColors: { [tag: string]: string } = {};

  constructor(player: VideoJsPlayer) {
    super(player);
    player.ready(() => {
      const tooltip = videojs.dom.createEl("div") as HTMLElement;
      tooltip.className = "vjs-marker-tooltip";
      tooltip.style.visibility = "hidden";

      const parent = player
        .el()
        .querySelector(".vjs-progress-holder .vjs-mouse-display");
      if (parent) parent.appendChild(tooltip);
      this.markerTooltip = tooltip;

      this.defaultTooltip = player
        .el()
        .querySelector<HTMLElement>(
          ".vjs-progress-holder .vjs-mouse-display .vjs-time-tooltip"
        );
    });
  }

  private showMarkerTooltip(title: string, layer: number = 0) {
    if (!this.markerTooltip) return;
    this.markerTooltip.innerText = title;
    this.markerTooltip.style.right = `${-this.markerTooltip.clientWidth / 2}px`;
    this.markerTooltip.style.top = `-${this.layerHeight * layer + 50}px`;
    this.markerTooltip.style.visibility = "visible";
    if (this.defaultTooltip) this.defaultTooltip.style.visibility = "hidden";
  }

  private hideMarkerTooltip() {
    if (this.markerTooltip) this.markerTooltip.style.visibility = "hidden";
    if (this.defaultTooltip) this.defaultTooltip.style.visibility = "visible";
  }

  addDotMarker(marker: IMarker) {
    const duration = this.player.duration();
    const markerSet: {
      dot?: HTMLDivElement;
      range?: HTMLDivElement;
    } = {};
    const seekBar = this.player.el().querySelector(".vjs-progress-holder");

    markerSet.dot = videojs.dom.createEl("div") as HTMLDivElement;
    markerSet.dot.className = "vjs-marker";
    if (duration) {
      // marker is 6px wide - adjust by 3px to align to center not left side
      markerSet.dot.style.left = `calc(${
        (marker.seconds / duration) * 100
      }% - 3px)`;
      markerSet.dot.style.visibility = "visible";
    }

    // Add event listeners to dot
    markerSet.dot.addEventListener("click", () =>
      this.player.currentTime(marker.seconds)
    );
    markerSet.dot.toggleAttribute("marker-tooltip-shown", true);

    // Set background color based on tag (if available)
    if (
      marker.primaryTag &&
      marker.primaryTag.name &&
      this.tagColors[marker.primaryTag.name]
    ) {
      markerSet.dot.style.backgroundColor =
        this.tagColors[marker.primaryTag.name];
    }
    markerSet.dot.addEventListener("mouseenter", () => {
      this.showMarkerTooltip(marker.title);
      markerSet.dot?.toggleAttribute("marker-tooltip-shown", true);
    });

    markerSet.dot.addEventListener("mouseout", () => {
      this.hideMarkerTooltip();
      markerSet.dot?.toggleAttribute("marker-tooltip-shown", false);
    });

    if (seekBar) {
      seekBar.appendChild(markerSet.dot);
    }
    this.markers.push(marker);
    this.markerDivs.push(markerSet);
  }

  addDotMarkers(markers: IMarker[]) {
    markers.forEach(this.addDotMarker, this);
  }

  private renderRangeMarkers(markers: IMarker[], layer: number) {
    const duration = this.player.duration();
    const parent = this.player.el().querySelector(".vjs-progress-control");
    const seekBar = this.player.el().querySelector(".vjs-progress-holder");
    if (!seekBar || !parent || !duration) return;

    markers.forEach((marker) => {
      this.renderRangeMarker(marker, layer, duration, seekBar, parent);
    });
  }

  private renderRangeMarker(
    marker: IMarker,
    layer: number,
    duration: number,
    seekBar: Element,
    parent: Element
  ) {
    if (!marker.end_seconds) return;

    const markerSet: {
      dot?: HTMLDivElement;
      range?: HTMLDivElement;
    } = {};
    const rangeDiv = videojs.dom.createEl("div") as HTMLDivElement;
    rangeDiv.className = "vjs-marker-range";

    // Use percentage-based positioning for proper scaling in fullscreen mode
    // The range marker is inside vjs-progress-control, but needs to align with
    // vjs-progress-holder which has 15px margins on each side.
    // We use calc() to combine percentage positioning with the fixed margin offset.
    const startPercent = (marker.seconds / duration) * 100;
    const widthPercent =
      ((marker.end_seconds - marker.seconds) / duration) * 100;

    // left: 15px margin + percentage of the progress holder width
    // Since progress-holder has margin: 0 15px, we need calc(15px + X% of remaining width)
    // The progress-holder width is (100% - 30px), so the actual left position is:
    // 15px + startPercent% * (100% - 30px) = 15px + startPercent% * 100% - startPercent% * 30px
    rangeDiv.style.left = `calc(15px + ${startPercent}% - ${
      startPercent * 0.3
    }px)`;

    rangeDiv.style.width = `calc(${widthPercent}% - ${widthPercent * 0.3}px)`;
    rangeDiv.style.bottom = `${layer * this.layerHeight}px`; // Adjust height based on layer
    rangeDiv.style.display = "none"; // Initially hidden

    // Set background color based on tag (if available)
    if (
      marker.primaryTag &&
      marker.primaryTag.name &&
      this.tagColors[marker.primaryTag.name]
    ) {
      rangeDiv.style.backgroundColor = this.tagColors[marker.primaryTag.name];
    }

    markerSet.range = rangeDiv;
    markerSet.range.style.display = "block";
    markerSet.range.addEventListener("pointermove", (e) => {
      e.stopPropagation();
    });
    markerSet.range.addEventListener("pointerover", (e) => {
      e.stopPropagation();
    });
    markerSet.range.addEventListener("pointerout", (e) => {
      e.stopPropagation();
    });
    markerSet.range.addEventListener("mouseenter", () => {
      this.showMarkerTooltip(marker.title, layer);
      markerSet.range?.toggleAttribute("marker-tooltip-shown", true);
    });

    markerSet.range.addEventListener("mouseout", () => {
      this.hideMarkerTooltip();
      markerSet.range?.toggleAttribute("marker-tooltip-shown", false);
    });
    parent.appendChild(rangeDiv);
    this.markers.push(marker);
    this.markerDivs.push(markerSet);
  }

  addRangeMarkers(markers: IMarker[]) {
    let remainingMarkers = [...markers];
    let layerNum = 0;

    while (remainingMarkers.length > 0) {
      // Get the set of markers that currently have the highest total duration that don't overlap. We do this layer by layer to prioritize filling
      // the lower layers when possible
      const mwis = this.findMWIS(remainingMarkers);
      if (!mwis.length) break;

      this.renderRangeMarkers(mwis, layerNum);
      remainingMarkers = remainingMarkers.filter(
        (marker) => !mwis.includes(marker)
      );
      layerNum++;
    }
  }

  // Use dynamic programming to find maximum weight independent set (ie the set of markers that have the highest total duration that don't overlap)
  private findMWIS(markers: IMarker[]): IMarker[] {
    if (!markers.length) return [];

    // Sort markers by end time
    markers = markers
      .slice()
      .sort((a, b) => (a.end_seconds || 0) - (b.end_seconds || 0));
    const n = markers.length;

    // Compute p(j) for each marker. This is the index of the marker that has the highest end time that doesn't overlap with marker j
    const p: number[] = new Array(n).fill(-1);
    for (let j = 0; j < n; j++) {
      for (let i = j - 1; i >= 0; i--) {
        if ((markers[i].end_seconds || 0) <= markers[j].seconds) {
          p[j] = i;
          break;
        }
      }
    }

    // Initialize M[j]
    // Compute M[j] for each marker. This is the maximum total duration of markers that don't overlap with marker j
    const M: number[] = new Array(n).fill(0);
    for (let j = 0; j < n; j++) {
      const include =
        (markers[j].end_seconds || 0) - markers[j].seconds + (M[p[j]] || 0);
      const exclude = j > 0 ? M[j - 1] : 0;
      M[j] = Math.max(include, exclude);
    }

    // Reconstruct optimal solution
    const findSolution = (j: number): IMarker[] => {
      if (j < 0) return [];
      const include =
        (markers[j].end_seconds || 0) - markers[j].seconds + (M[p[j]] || 0);
      const exclude = j > 0 ? M[j - 1] : 0;
      if (include >= exclude) {
        return [...findSolution(p[j]), markers[j]];
      } else {
        return findSolution(j - 1);
      }
    };

    return findSolution(n - 1);
  }

  removeMarker(marker: IMarker) {
    const i = this.markers.indexOf(marker);
    if (i === -1) return;

    this.markers.splice(i, 1);
    const markerSet = this.markerDivs.splice(i, 1)[0];

    if (markerSet.dot?.hasAttribute("marker-tooltip-shown")) {
      this.hideMarkerTooltip();
    }

    markerSet.dot?.remove();
    if (markerSet.range) markerSet.range.remove();
  }

  removeMarkers(markers: IMarker[]) {
    markers.forEach(this.removeMarker, this);
  }

  clearMarkers() {
    for (const markerSet of this.markerDivs) {
      if (markerSet.dot?.hasAttribute("marker-tooltip-shown")) {
        this.hideMarkerTooltip();
      }

      markerSet.dot?.remove();
      if (markerSet.range) markerSet.range.remove();
    }
    this.markers = [];
    this.markerDivs = [];
  }

  // Implementing the findColors method
  findColors(tagNames: string[]) {
    // Compute base hues for each tag
    const baseHues: { [tag: string]: number } = {};
    for (const tag of tagNames) {
      baseHues[tag] = this.computeBaseHue(tag);
    }

    // Adjust hues to avoid similar colors
    const adjustedHues = this.adjustHues(baseHues);

    // Convert adjusted hues to colors and store in tagColors dictionary
    for (const tag of tagNames) {
      this.tagColors[tag] = this.hueToColor(adjustedHues[tag]);
    }
  }

  // Helper methods translated from Python

  // Compute base hue from tag name
  private computeBaseHue(tag: string): number {
    const hash = CryptoJS.SHA256(tag);
    const hashHex = hash.toString(CryptoJS.enc.Hex);
    const hashInt = BigInt(`0x${hashHex}`);
    const baseHue = Number(hashInt % BigInt(360)); // Map to [0, 360)
    return baseHue;
  }

  // Calculate minimum acceptable hue difference based on number of tags
  private calculateDeltaMin(N: number): number {
    const maxDeltaNeeded = 35;
    let scalingFactor: number;

    if (N <= 4) {
      scalingFactor = 0.8;
    } else if (N <= 10) {
      scalingFactor = 0.6;
    } else {
      scalingFactor = 0.4;
    }

    const deltaMin = Math.min((360 / N) * scalingFactor, maxDeltaNeeded);
    return deltaMin;
  }

  // Adjust hues to ensure minimum difference
  private adjustHues(baseHues: { [tag: string]: number }): {
    [tag: string]: number;
  } {
    const adjustedHues: { [tag: string]: number } = {};
    const tags = Object.keys(baseHues);
    const N = tags.length;
    const deltaMin = this.calculateDeltaMin(N);

    // Sort the tags by base hue
    const sortedTags = tags.sort((a, b) => baseHues[a] - baseHues[b]);
    // Get sorted base hues
    const baseHuesSorted = sortedTags.map((tag) => baseHues[tag]);

    // Unwrap hues to handle circular nature
    const unwrappedHues = [...baseHuesSorted];
    for (let i = 1; i < N; i++) {
      if (unwrappedHues[i] <= unwrappedHues[i - 1]) {
        unwrappedHues[i] += 360; // Unwrap by adding 360 degrees
      }
    }

    // Adjust hues to ensure minimum difference
    for (let i = 1; i < N; i++) {
      const requiredHue = unwrappedHues[i - 1] + deltaMin;
      if (unwrappedHues[i] < requiredHue) {
        unwrappedHues[i] = requiredHue; // Adjust hue minimally
      }
    }

    // Handle wrap-around difference
    const endGap = unwrappedHues[0] + 360 - unwrappedHues[N - 1];
    if (endGap < deltaMin) {
      // Adjust first and last hues minimally to increase end gap
      const adjustmentNeeded = (deltaMin - endGap) / 2;
      // Adjust the first hue backward, ensure it doesn't go below other hues
      unwrappedHues[0] = Math.max(
        unwrappedHues[0] - adjustmentNeeded,
        unwrappedHues[1] - 360 + deltaMin
      );
      // Adjust the last hue forward
      unwrappedHues[N - 1] += adjustmentNeeded;
    }

    // Wrap adjusted hues back to [0, 360)
    const adjustedHuesList = unwrappedHues.map((hue) => hue % 360);

    // Map adjusted hues back to tags
    for (let i = 0; i < N; i++) {
      adjustedHues[sortedTags[i]] = adjustedHuesList[i];
    }

    return adjustedHues;
  }

  // Convert hue to RGB color in hex format
  private hueToColor(hue: number): string {
    // Convert hue from degrees to [0, 1)
    const hueNormalized = hue / 360.0;
    const saturation = 0.65;
    const value = 0.95;
    const rgb = this.hsvToRgb(hueNormalized, saturation, value);
    const alpha = 0.6; // Set the desired alpha value here
    const rgbColor = `#${this.toHex(rgb[0])}${this.toHex(rgb[1])}${this.toHex(
      rgb[2]
    )}${this.toHex(Math.round(alpha * 255))}`;
    return rgbColor;
  }

  // Convert HSV to RGB
  private hsvToRgb(h: number, s: number, v: number): [number, number, number] {
    const i = Math.floor(h * 6);
    const f = h * 6 - i;
    const p = v * (1 - s);
    const q = v * (1 - f * s);
    const t = v * (1 - (1 - f) * s);

    let r, g, b;
    switch (i % 6) {
      case 0:
        r = v;
        g = t;
        b = p;
        break;
      case 1:
        r = q;
        g = v;
        b = p;
        break;
      case 2:
        r = p;
        g = v;
        b = t;
        break;
      case 3:
        r = p;
        g = q;
        b = v;
        break;
      case 4:
        r = t;
        g = p;
        b = v;
        break;
      case 5:
        r = v;
        g = p;
        b = q;
        break;
      default:
        r = v;
        g = t;
        b = p;
        break;
    }

    return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)];
  }

  // Convert a number to two-digit hex string
  private toHex(value: number): string {
    return value.toString(16).padStart(2, "0");
  }
}

videojs.registerPlugin("markers", MarkersPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    markers: () => MarkersPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    markers?: IMarkersOptions;
  }
}

export default MarkersPlugin;
