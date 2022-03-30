import videojs, { VideoJsPlayer } from "video.js";

const markers = function (this: VideoJsPlayer) {
  const player = this;

  function getPosition(marker: VTTCue) {
    return (marker.startTime / player.duration()) * 100;
  }

  function createMarkerToolTip() {
    const tooltip = videojs.dom.createEl("div") as HTMLElement;
    tooltip.className = "vjs-marker-tooltip";

    return tooltip;
  }

  function removeMarkerToolTip() {
    const div = player
      .el()
      .querySelector(".vjs-progress-holder .vjs-marker-tooltip");
    if (div) div.remove();
  }

  function createMarkerDiv(marker: VTTCue) {
    const markerDiv = videojs.dom.createEl(
      "div",
      {},
      {
        "data-marker-time": marker.startTime,
      }
    ) as HTMLElement;

    markerDiv.className = "vjs-marker";
    markerDiv.style.left = getPosition(marker) + "%";

    // bind click event to seek to marker time
    markerDiv.addEventListener("click", function () {
      const time = this.getAttribute("data-marker-time");
      player.currentTime(Number(time));
    });

    // show tooltip on hover
    markerDiv.addEventListener("mouseenter", function () {
      // create and show tooltip
      const tooltip = createMarkerToolTip();
      tooltip.innerText = marker.text;

      const parent = player
        .el()
        .querySelector(".vjs-progress-holder .vjs-mouse-display");

      parent?.appendChild(tooltip);

      // hide default tooltip
      const defaultTooltip = parent?.querySelector(
        ".vjs-time-tooltip"
      ) as HTMLElement;
      defaultTooltip.style.visibility = "hidden";
    });

    markerDiv.addEventListener("mouseout", function () {
      removeMarkerToolTip();

      // show default tooltip
      const defaultTooltip = player
        .el()
        .querySelector(
          ".vjs-progress-holder .vjs-mouse-display .vjs-time-tooltip"
        ) as HTMLElement;
      if (defaultTooltip) defaultTooltip.style.visibility = "visible";
    });

    return markerDiv;
  }

  function removeMarkerDivs() {
    const divs = player
      .el()
      .querySelectorAll(".vjs-progress-holder .vjs-marker");
    divs.forEach((div) => {
      div.remove();
    });
  }

  this.on("loadedmetadata", function () {
    removeMarkerDivs();
    removeMarkerToolTip();

    const textTracks = player.remoteTextTracks();
    const seekBar = player.el().querySelector(".vjs-progress-holder");

    if (seekBar && textTracks.length > 0) {
      const vttTrack = textTracks[0];
      if (!vttTrack || !vttTrack.cues) return;
      for (let i = 0; i < vttTrack.cues.length; i++) {
        const cue = vttTrack.cues[i];
        const markerDiv = createMarkerDiv(cue as VTTCue);
        seekBar.appendChild(markerDiv);
      }
    }
  });
};

// Register the plugin with video.js.
videojs.registerPlugin("markers", markers);

export default markers;
