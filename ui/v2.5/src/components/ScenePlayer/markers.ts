import videojs, { VideoJsPlayer } from "video.js";

interface IMarker {
  title: string;
  time: number;
}

interface IMarkersOptions {
  markers?: IMarker[];
}

class MarkersPlugin extends videojs.getPlugin("plugin") {
  private markers: IMarker[] = [];
  private markerDivs: HTMLDivElement[] = [];
  private markerTooltip: HTMLElement | null = null;
  private defaultTooltip: HTMLElement | null = null;

  constructor(player: VideoJsPlayer, options?: IMarkersOptions) {
    super(player);

    player.ready(() => {
      // create marker tooltip
      const tooltip = videojs.dom.createEl("div") as HTMLElement;
      tooltip.className = "vjs-marker-tooltip";
      tooltip.style.visibility = "hidden";

      const parent = player
        .el()
        .querySelector(".vjs-progress-holder .vjs-mouse-display");
      if (parent) parent.appendChild(tooltip);
      this.markerTooltip = tooltip;

      // save default tooltip
      this.defaultTooltip = player
        .el()
        .querySelector<HTMLElement>(
          ".vjs-progress-holder .vjs-mouse-display .vjs-time-tooltip"
        );

      options?.markers?.forEach(this.addMarker, this);
    });

    player.on("loadedmetadata", () => {
      const seekBar = player.el().querySelector(".vjs-progress-holder");
      const duration = this.player.duration();

      for (let i = 0; i < this.markers.length; i++) {
        const marker = this.markers[i];
        const markerDiv = this.markerDivs[i];

        if (duration) {
          // marker is 6px wide - adjust by 3px to align to center not left side
          markerDiv.style.left = `calc(${
            (marker.time / duration) * 100
          }% - 3px)`;
          markerDiv.style.visibility = "visible";
        }
        if (seekBar) seekBar.appendChild(markerDiv);
      }
    });
  }

  private showMarkerTooltip(title: string) {
    if (!this.markerTooltip) return;

    this.markerTooltip.innerText = title;
    this.markerTooltip.style.right = `${-this.markerTooltip.clientWidth / 2}px`;
    this.markerTooltip.style.visibility = "visible";

    // hide default tooltip
    if (this.defaultTooltip) this.defaultTooltip.style.visibility = "hidden";
  }

  private hideMarkerTooltip() {
    if (this.markerTooltip) this.markerTooltip.style.visibility = "hidden";

    // show default tooltip
    if (this.defaultTooltip) this.defaultTooltip.style.visibility = "visible";
  }

  addMarker(marker: IMarker) {
    const markerDiv = videojs.dom.createEl("div") as HTMLDivElement;
    markerDiv.className = "vjs-marker";

    const duration = this.player.duration();
    if (duration) {
      // marker is 6px wide - adjust by 3px to align to center not left side
      markerDiv.style.left = `calc(${(marker.time / duration) * 100}% - 3px)`;
      markerDiv.style.visibility = "visible";
    }

    // bind click event to seek to marker time
    markerDiv.addEventListener("click", () =>
      this.player.currentTime(marker.time)
    );

    // show/hide tooltip on hover
    markerDiv.addEventListener("mouseenter", () => {
      this.showMarkerTooltip(marker.title);
      markerDiv.toggleAttribute("marker-tooltip-shown", true);
    });
    markerDiv.addEventListener("mouseout", () => {
      this.hideMarkerTooltip();
      markerDiv.toggleAttribute("marker-tooltip-shown", false);
    });

    const seekBar = this.player.el().querySelector(".vjs-progress-holder");
    if (seekBar) seekBar.appendChild(markerDiv);

    this.markers.push(marker);
    this.markerDivs.push(markerDiv);
  }

  addMarkers(markers: IMarker[]) {
    markers.forEach(this.addMarker, this);
  }

  removeMarker(marker: IMarker) {
    const i = this.markers.indexOf(marker);
    if (i === -1) return;

    this.markers.splice(i, 1);

    const div = this.markerDivs.splice(i, 1)[0];
    if (div.hasAttribute("marker-tooltip-shown")) {
      this.hideMarkerTooltip();
    }
    div.remove();
  }

  removeMarkers(markers: IMarker[]) {
    markers.forEach(this.removeMarker, this);
  }

  clearMarkers() {
    this.removeMarkers([...this.markers]);
  }
}

// Register the plugin with video.js.
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
