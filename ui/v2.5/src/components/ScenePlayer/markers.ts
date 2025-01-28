import videojs, { VideoJsPlayer } from "video.js";
import "./markers.css";
import { layer } from "@fortawesome/fontawesome-svg-core";

export interface IMarker {
  title: string;
  seconds: number;
  end_seconds?: number | null;
  layer?: number;
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

  private layers: IMarker[][] = [];

  private layerHeight: number = 6;

  private tagColors: { [tag: string]: string } = {};

  constructor(player: VideoJsPlayer, options?: IMarkersOptions) {
    super(player);
    player.ready(() => {
      const tooltip = videojs.dom.createEl("div") as HTMLElement;
      tooltip.className = "vjs-marker-tooltip";
      tooltip.style.visibility = "hidden";

      const parent = player.el().querySelector(".vjs-progress-holder .vjs-mouse-display");
      if (parent) parent.appendChild(tooltip);
      this.markerTooltip = tooltip;

      this.defaultTooltip = player.el().querySelector<HTMLElement>(
        ".vjs-progress-holder .vjs-mouse-display .vjs-time-tooltip"
      );
    });
  }

  private showMarkerTooltip(title: string, layer:number = 0) {
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

  private formatTime(seconds: number): string {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = Math.floor(seconds % 60);
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  }

  addDotMarker(marker: IMarker) {
    const duration = this.player.duration();
    const markerSet: {
      dot?: HTMLDivElement;
      range?: HTMLDivElement;
    } = {};
    const seekBar = this.player.el().querySelector(".vjs-progress-holder");
    if (marker.end_seconds) {
      throw new Error("Cannot add range marker with addDotMarker");
    }
    markerSet.dot = videojs.dom.createEl("div") as HTMLDivElement
    markerSet.dot.className = "vjs-marker-dot";
    if (duration) {
      markerSet.dot.style.left = `calc(${(marker.seconds / duration) * 100}% - 3px)`;
    }

    // Add event listeners to dot
    markerSet.dot.addEventListener("click", () => this.player.currentTime(marker.seconds));
    markerSet.dot.toggleAttribute("marker-tooltip-shown", true);

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
  }

  private renderRangeMarkers(markers: IMarker[], layer: number) {
    const duration = this.player.duration();
    const seekBar = this.player.el().querySelector(".vjs-progress-holder");
    if (!seekBar || !duration) return;

    markers.forEach(marker => {
      this.renderRangeMarker(marker, layer, duration, seekBar);
    });
  }

  private renderRangeMarker(marker: IMarker, layer: number, duration: number, seekBar: Element) {
    if (!marker.end_seconds) return;
      
    const markerSet: {
      dot?: HTMLDivElement;
      range?: HTMLDivElement;
    } = {
    };
    const rangeDiv = videojs.dom.createEl("div") as HTMLDivElement;
    rangeDiv.className = "vjs-marker-range";
    
    const startPercent = (marker.seconds / duration) * 100;
    const endPercent = (marker.end_seconds / duration) * 100;
    const width = endPercent - startPercent;
    
    rangeDiv.style.left = `${startPercent}%`;
    rangeDiv.style.width = `${width}%`;
    rangeDiv.style.bottom = `${layer * this.layerHeight + 10}px`; // Adjust height based on layer
    rangeDiv.style.display = 'none'; // Initially hidden

    markerSet.range = rangeDiv;
    markerSet.range.style.display = 'block';
    markerSet.range.addEventListener("mouseenter", () => {
      this.showMarkerTooltip(marker.title, layer);
      markerSet.range?.toggleAttribute("marker-tooltip-shown", true);
    });

    markerSet.range.addEventListener("mouseout", () => {
      this.hideMarkerTooltip();
      markerSet.range?.toggleAttribute("marker-tooltip-shown", false);
    });
    seekBar.appendChild(rangeDiv);
    this.markers.push(marker);
    this.markerDivs.push(markerSet);
  }

  addRangeMarkersNew(markers: IMarker[]) {
    let remainingMarkers = [...markers];
    this.layers = [];
    let layerNum = 0;

    while (remainingMarkers.length > 0) {
      const mwis = this.findMWIS(remainingMarkers);
      if (!mwis.length) break;

      this.layers.push(mwis);
      console.log("Rendering layer", layerNum, mwis);
      this.renderRangeMarkers(mwis, layerNum);
      remainingMarkers = remainingMarkers.filter(marker => !mwis.includes(marker));
      layerNum++;
    }
  }

  private findMWIS(markers: IMarker[]): IMarker[] {
    if (!markers.length) return [];

    // Sort markers by end time
    markers = markers.slice().sort((a, b) => (a.end_seconds || 0) - (b.end_seconds || 0));
    const n = markers.length;

    // Compute p(j) for each marker
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
    const M: number[] = new Array(n).fill(0);
    for (let j = 0; j < n; j++) {
      const include = (markers[j].end_seconds || 0) - markers[j].seconds + (M[p[j]] || 0);
      const exclude = j > 0 ? M[j - 1] : 0;
      M[j] = Math.max(include, exclude);
    }

    // Reconstruct optimal solution
    const findSolution = (j: number): IMarker[] => {
      if (j < 0) return [];
      const include = (markers[j].end_seconds || 0) - markers[j].seconds + (M[p[j]] || 0);
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
    this.removeMarkers([...this.markers]);
  }

  findColors(tagNames: string[]){
    
  }
}

videojs.registerPlugin("markers", MarkersPlugin);

declare module "video.js" {
  interface VideoJsPlayer {
    markers: () => MarkersPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    markers?: IMarkersOptions;
  }
}

export default MarkersPlugin;