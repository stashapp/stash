import { VIDEO_PLAYER_ID } from "src/components/ScenePlayer/util";

export type SliderRange = {
  min: number;
  default: number;
  max: number;
  divider: number;
};

export const sliderRanges: { [key: string]: SliderRange } = {
  contrastRange: {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  },
  brightnessRange: {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  },
  gammaRange: {
    min: 0,
    default: 100,
    max: 200,
    divider: 200,
  },
  saturateRange: {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  },
  hueRotateRange: {
    min: 0,
    default: 0,
    max: 360,
    divider: 1,
  },
  warmthRange: {
    min: 0,
    default: 100,
    max: 200,
    divider: 200,
  },
  colourRange: {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  },
  blurRange: {
    min: 0,
    default: 0,
    max: 250,
    divider: 10,
  },
  rotateRange: {
    min: 0,
    default: 2,
    max: 4,
    divider: 1 / 90,
  },
  scaleRange: {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  },
  aspectRatioRange: {
    min: 0,
    default: 150,
    max: 300,
    divider: 100,
  },
};

// eslint-disable-next-line
function getVideoElement(playerVideoContainer: any) {
  let videoElements = playerVideoContainer.getElementsByTagName("canvas");

  if (videoElements.length == 0) {
    videoElements = playerVideoContainer.getElementsByTagName("video");
  }

  if (videoElements.length > 0) {
    return videoElements[0];
  }
}

export function updateVideoFilters(
  gammaValue: number,
  redValue: number,
  greenValue: number,
  blueValue: number,
  warmthValue: number
): void {
  const filterContainer = document.getElementById("video-filter-container");

  if (filterContainer == null) {
    return;
  }

  const svg1 = document.createElementNS("http://www.w3.org/2000/svg", "svg");
  const videoFilter = document.createElementNS(
    "http://www.w3.org/2000/svg",
    "filter"
  );
  videoFilter.setAttribute("id", "videoFilter");

  if (
    warmthValue !== sliderRanges.warmthRange.default ||
    redValue !== sliderRanges.colourRange.default ||
    greenValue !== sliderRanges.colourRange.default ||
    blueValue !== sliderRanges.colourRange.default
  ) {
    const feColorMatrix = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "feColorMatrix"
    );
    feColorMatrix.setAttribute(
      "values",
      `${
        1 +
        (warmthValue - sliderRanges.warmthRange.default) /
          sliderRanges.warmthRange.divider +
        (redValue - sliderRanges.colourRange.default) /
          sliderRanges.colourRange.divider
      } 0 0 0 0   0 ${
        1.0 +
        (greenValue - sliderRanges.colourRange.default) /
          sliderRanges.colourRange.divider
      } 0 0 0   0 0 ${
        1 -
        (warmthValue - sliderRanges.warmthRange.default) /
          sliderRanges.warmthRange.divider +
        (blueValue - sliderRanges.colourRange.default) /
          sliderRanges.colourRange.divider
      } 0 0   0 0 0 1.0 0`
    );
    videoFilter.appendChild(feColorMatrix);
  }

  if (gammaValue !== sliderRanges.gammaRange.default) {
    const feComponentTransfer = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "feComponentTransfer"
    );

    const feFuncR = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "feFuncR"
    );
    feFuncR.setAttribute("type", "gamma");
    feFuncR.setAttribute("amplitude", "1.0");
    feFuncR.setAttribute(
      "exponent",
      `${
        1 +
        (sliderRanges.gammaRange.default - gammaValue) /
          sliderRanges.gammaRange.divider
      }`
    );
    feFuncR.setAttribute("offset", "0.0");
    feComponentTransfer.appendChild(feFuncR);

    const feFuncG = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "feFuncG"
    );
    feFuncG.setAttribute("type", "gamma");
    feFuncG.setAttribute("amplitude", "1.0");
    feFuncG.setAttribute(
      "exponent",
      `${
        1 +
        (sliderRanges.gammaRange.default - gammaValue) /
          sliderRanges.gammaRange.divider
      }`
    );
    feFuncG.setAttribute("offset", "0.0");
    feComponentTransfer.appendChild(feFuncG);

    const feFuncB = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "feFuncB"
    );
    feFuncB.setAttribute("type", "gamma");
    feFuncB.setAttribute("amplitude", "1.0");
    feFuncB.setAttribute(
      "exponent",
      `${
        1 +
        (sliderRanges.gammaRange.default - gammaValue) /
          sliderRanges.gammaRange.divider
      }`
    );
    feFuncB.setAttribute("offset", "0.0");
    feComponentTransfer.appendChild(feFuncB);

    const feFuncA = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "feFuncA"
    );
    feFuncA.setAttribute("type", "gamma");
    feFuncA.setAttribute("amplitude", "1.0");
    feFuncA.setAttribute("exponent", "1.0");
    feFuncA.setAttribute("offset", "0.0");
    feComponentTransfer.appendChild(feFuncA);

    videoFilter.appendChild(feComponentTransfer);
  }

  svg1.appendChild(videoFilter);

  // Add or Replace existing svg
  const filterContainerSvgs = filterContainer.getElementsByTagNameNS(
    "http://www.w3.org/2000/svg",
    "svg"
  );
  if (filterContainerSvgs.length === 0) {
    // attach container to document
    filterContainer.appendChild(svg1);
  } else {
    // assume only one svg... maybe issue
    filterContainer.replaceChild(svg1, filterContainerSvgs[0]);
  }
}

export function updateVideoStyle(
  aspectRatioValue: number,
  blurValue: number,
  brightnessValue: number,
  contrastValue: number,
  gammaValue: number,
  hueRotateValue: number,
  redValue: number,
  greenValue: number,
  blueValue: number,
  rotateValue: number,
  saturateValue: number,
  scaleValue: number,
  warmthValue: number
): void {
  const playerVideoContainer = document.getElementById(VIDEO_PLAYER_ID)!;
  if (!playerVideoContainer) {
    return;
  }

  const playerVideoElement = getVideoElement(playerVideoContainer);
  if (playerVideoElement != null) {
    let styleString = "filter:";
    let style = playerVideoElement.attributes.getNamedItem("style");

    if (style == null) {
      style = document.createAttribute("style");
      playerVideoElement.attributes.setNamedItem(style);
    }

    if (
      warmthValue !== sliderRanges.warmthRange.default ||
      redValue !== sliderRanges.colourRange.default ||
      greenValue !== sliderRanges.colourRange.default ||
      blueValue !== sliderRanges.colourRange.default ||
      gammaValue !== sliderRanges.gammaRange.default
    ) {
      styleString += " url(#videoFilter)";
    }

    if (contrastValue !== sliderRanges.contrastRange.default) {
      styleString += ` contrast(${contrastValue}%)`;
    }

    if (brightnessValue !== sliderRanges.brightnessRange.default) {
      styleString += ` brightness(${brightnessValue}%)`;
    }

    if (saturateValue !== sliderRanges.saturateRange.default) {
      styleString += ` saturate(${saturateValue}%)`;
    }

    if (hueRotateValue !== sliderRanges.hueRotateRange.default) {
      styleString += ` hue-rotate(${hueRotateValue}deg)`;
    }

    if (blurValue > sliderRanges.blurRange.default) {
      styleString += ` blur(${blurValue / sliderRanges.blurRange.divider}px)`;
    }

    styleString += "; transform:";

    if (rotateValue !== sliderRanges.rotateRange.default) {
      styleString += ` rotate(${
        (rotateValue - sliderRanges.rotateRange.default) /
        sliderRanges.rotateRange.divider
      }deg)`;
    }

    if (
      scaleValue !== sliderRanges.scaleRange.default ||
      aspectRatioValue !== sliderRanges.aspectRatioRange.default
    ) {
      let xScale = scaleValue / sliderRanges.scaleRange.divider / 100.0;
      let yScale = scaleValue / sliderRanges.scaleRange.divider / 100.0;

      if (aspectRatioValue > sliderRanges.aspectRatioRange.default) {
        xScale *=
          (sliderRanges.aspectRatioRange.divider +
            aspectRatioValue -
            sliderRanges.aspectRatioRange.default) /
          sliderRanges.aspectRatioRange.divider;
      } else if (aspectRatioValue < sliderRanges.aspectRatioRange.default) {
        yScale *=
          (sliderRanges.aspectRatioRange.divider +
            sliderRanges.aspectRatioRange.default -
            aspectRatioValue) /
          sliderRanges.aspectRatioRange.divider;
      }

      styleString += ` scale(${xScale},${yScale})`;
    }

    if (playerVideoElement.tagName == "CANVAS") {
      styleString += "; width: 100%; height: 100%; position: absolute; top:0";
    }

    style.value = `${styleString};`;
  }
}
