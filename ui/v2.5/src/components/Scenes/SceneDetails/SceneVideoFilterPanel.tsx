import React, { useState } from "react";
import { Button, Form } from "react-bootstrap";
import { JWUtils } from "../../../utils";
import * as GQL from "../../../core/generated-graphql";

interface ISceneVideoFilterPanelProps {
  scene: GQL.SceneDataFragment;
}

// References
// https://yoksel.github.io/svg-filters/#/
// https://codepen.io/chriscoyier/pen/zbakI
// http://xahlee.info/js/js_scritping_svg_basics.html#:~:text=Just%20use%20JavaScript%20to%20script,%2C%20path%2C%20%E2%80%A6.).

type SliderRange = {
  min: number;
  default: number;
  max: number;
  divider: number;
};

export const SceneVideoFilterPanel: React.FC<ISceneVideoFilterPanelProps> = (
  props: ISceneVideoFilterPanelProps
) => {
  const contrastRange: SliderRange = {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  };
  const brightnessRange: SliderRange = {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  };
  const gammaRange: SliderRange = {
    min: 0,
    default: 100,
    max: 200,
    divider: 200,
  };
  const saturateRange: SliderRange = {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  };
  const hueRotateRange: SliderRange = {
    min: 0,
    default: 0,
    max: 360,
    divider: 1,
  };
  const whiteBalanceRange: SliderRange = {
    min: 0,
    default: 100,
    max: 200,
    divider: 200,
  };
  const colourRange: SliderRange = {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  };
  const blurRange: SliderRange = { min: 0, default: 0, max: 250, divider: 10 };
  const rotateRange: SliderRange = {
    min: 0,
    default: 2,
    max: 4,
    divider: 1 / 90,
  };
  const scaleRange: SliderRange = {
    min: 0,
    default: 100,
    max: 200,
    divider: 1,
  };
  const aspectRatioRange: SliderRange = {
    min: 0,
    default: 150,
    max: 300,
    divider: 100,
  };

  const [contrastValue, setContrastValue] = useState(contrastRange.default);
  const [brightnessValue, setBrightnessValue] = useState(
    brightnessRange.default
  );
  const [gammaValue, setGammaValue] = useState(gammaRange.default);
  const [saturateValue, setSaturateValue] = useState(saturateRange.default);
  const [hueRotateValue, setHueRotateValue] = useState(hueRotateRange.default);
  const [whiteBalanceValue, setWhiteBalanceValue] = useState(
    whiteBalanceRange.default
  );
  const [redValue, setRedValue] = useState(colourRange.default);
  const [greenValue, setGreenValue] = useState(colourRange.default);
  const [blueValue, setBlueValue] = useState(colourRange.default);
  const [blurValue, setBlurValue] = useState(blurRange.default);
  const [rotateValue, setRotateValue] = useState(rotateRange.default);
  const [scaleValue, setScaleValue] = useState(scaleRange.default);
  const [aspectRatioValue, setAspectRatioValue] = useState(
    aspectRatioRange.default
  );

  function updateVideoStyle() {
    const playerId = JWUtils.playerID;
    const playerVideoElement = document
      .getElementById(playerId)
      ?.getElementsByClassName("jw-video")[0];

    if (playerVideoElement != null) {
      let styleString = "filter:";
      let style = playerVideoElement.attributes.getNamedItem("style");

      if (style == null) {
        style = document.createAttribute("style");
        playerVideoElement.attributes.setNamedItem(style);
      }

      if (
        whiteBalanceValue !== whiteBalanceRange.default ||
        redValue !== colourRange.default ||
        greenValue !== colourRange.default ||
        blueValue !== colourRange.default ||
        gammaValue !== gammaRange.default
      ) {
        styleString += " url(#videoFilter)";
      }

      if (contrastValue !== contrastRange.default) {
        styleString += ` contrast(${contrastValue}%)`;
      }

      if (brightnessValue !== brightnessRange.default) {
        styleString += ` brightness(${brightnessValue}%)`;
      }

      if (saturateValue !== saturateRange.default) {
        styleString += ` saturate(${saturateValue}%)`;
      }

      if (hueRotateValue !== hueRotateRange.default) {
        styleString += ` hue-rotate(${hueRotateValue}deg)`;
      }

      if (blurValue > blurRange.default) {
        styleString += ` blur(${blurValue / blurRange.divider}px)`;
      }

      styleString += "; transform:";

      if (rotateValue !== rotateRange.default) {
        styleString += ` rotate(${
          (rotateValue - rotateRange.default) / rotateRange.divider
        }deg)`;
      }

      if (
        scaleValue !== scaleRange.default ||
        aspectRatioValue !== aspectRatioRange.default
      ) {
        let xScale = scaleValue / scaleRange.divider / 100.0;
        let yScale = scaleValue / scaleRange.divider / 100.0;

        if (aspectRatioValue > aspectRatioRange.default) {
          xScale *=
            (aspectRatioRange.divider +
              aspectRatioValue -
              aspectRatioRange.default) /
            aspectRatioRange.divider;
        } else if (aspectRatioValue < aspectRatioRange.default) {
          yScale *=
            (aspectRatioRange.divider +
              aspectRatioRange.default -
              aspectRatioValue) /
            aspectRatioRange.divider;
        }

        styleString += ` scale(${xScale},${yScale})`;
      }

      style.value = `${styleString};`;
    }
  }

  function updateVideoFilters() {
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
      whiteBalanceValue !== whiteBalanceRange.default ||
      redValue !== colourRange.default ||
      greenValue !== colourRange.default ||
      blueValue !== colourRange.default
    ) {
      const feColorMatrix = document.createElementNS(
        "http://www.w3.org/2000/svg",
        "feColorMatrix"
      );
      feColorMatrix.setAttribute(
        "values",
        `${
          1 +
          (whiteBalanceValue - whiteBalanceRange.default) /
            whiteBalanceRange.divider +
          (redValue - colourRange.default) / colourRange.divider
        } 0 0 0 0   0 ${
          1.0 + (greenValue - colourRange.default) / colourRange.divider
        } 0 0 0   0 0 ${
          1 -
          (whiteBalanceValue - whiteBalanceRange.default) /
            whiteBalanceRange.divider +
          (blueValue - colourRange.default) / colourRange.divider
        } 0 0   0 0 0 1.0 0`
      );
      videoFilter.appendChild(feColorMatrix);
    }

    if (gammaValue !== gammaRange.default) {
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
        `${1 + (gammaRange.default - gammaValue) / gammaRange.divider}`
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
        `${1 + (gammaRange.default - gammaValue) / gammaRange.divider}`
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
        `${1 + (gammaRange.default - gammaValue) / gammaRange.divider}`
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

  function renderBlur() {
    return (
      <div className="row form-group">
        <span className="col-3">Blur</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={blurRange.min}
            max={blurRange.max}
            value={blurValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setBlurValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setBlurValue(blurRange.default)}
          onKeyPress={() => setBlurValue(blurRange.default)}
        >
          {blurValue / blurRange.divider}px
        </span>
      </div>
    );
  }

  function renderContrast() {
    return (
      <div className="row form-group">
        <span className="col-3">Contrast</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider contrast-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={contrastRange.min}
            max={contrastRange.max}
            value={contrastValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setContrastValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setContrastValue(contrastRange.default)}
          onKeyPress={() => setContrastValue(contrastRange.default)}
        >
          {contrastValue / brightnessRange.divider}%
        </span>
      </div>
    );
  }

  function renderBrightness() {
    return (
      <div className="row form-group">
        <span className="col-3">Brightness</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider brightness-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={brightnessRange.min}
            max={brightnessRange.max}
            value={brightnessValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setBrightnessValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setBrightnessValue(brightnessRange.default)}
          onKeyPress={() => setBrightnessValue(brightnessRange.default)}
        >
          {brightnessValue / brightnessRange.divider}%
        </span>
      </div>
    );
  }

  function renderGammaSlider() {
    return (
      <div className="row form-group">
        <span className="col-3">Gamma</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider gamma-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={gammaRange.min}
            max={gammaRange.max}
            value={gammaValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setGammaValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setGammaValue(gammaRange.default)}
          onKeyPress={() => setGammaValue(gammaRange.default)}
        >
          {(gammaValue - gammaRange.default) / gammaRange.divider}
        </span>
      </div>
    );
  }

  function renderSaturate() {
    return (
      <div className="row form-group">
        <span className="col-3">Saturation</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider saturation-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={saturateRange.min}
            max={saturateRange.max}
            value={saturateValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setSaturateValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setSaturateValue(saturateRange.default)}
          onKeyPress={() => setSaturateValue(saturateRange.default)}
        >
          {saturateValue / saturateRange.divider}%
        </span>
      </div>
    );
  }

  function renderHueRotateSlider() {
    return (
      <div className="row form-group">
        <span className="col-3">Hue</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider hue-rotate-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={hueRotateRange.min}
            max={hueRotateRange.max}
            value={hueRotateValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setHueRotateValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setHueRotateValue(hueRotateRange.default)}
          onKeyPress={() => setHueRotateValue(hueRotateRange.default)}
        >
          {hueRotateValue / hueRotateRange.divider}&deg;
        </span>
      </div>
    );
  }

  function renderWhiteBalance() {
    return (
      <div className="row form-group">
        <span className="col-3">Warmth</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider white-balance-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={whiteBalanceRange.min}
            max={whiteBalanceRange.max}
            value={whiteBalanceValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setWhiteBalanceValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setWhiteBalanceValue(whiteBalanceRange.default)}
          onKeyPress={() => setWhiteBalanceValue(whiteBalanceRange.default)}
        >
          {(whiteBalanceValue - whiteBalanceRange.default) /
            whiteBalanceRange.divider}
        </span>
      </div>
    );
  }

  function renderRedSlider() {
    return (
      <div className="row form-group">
        <span className="col-3">Red</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider red-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={colourRange.min}
            max={colourRange.max}
            value={redValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setRedValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setRedValue(colourRange.default)}
          onKeyPress={() => setRedValue(colourRange.default)}
        >
          {(redValue - colourRange.default) / colourRange.divider}%
        </span>
      </div>
    );
  }

  function renderGreenSlider() {
    return (
      <div className="row form-group">
        <span className="col-3">Green</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider green-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={colourRange.min}
            max={colourRange.max}
            value={greenValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setGreenValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setGreenValue(colourRange.default)}
          onKeyPress={() => setGreenValue(colourRange.default)}
        >
          {(greenValue - colourRange.default) / colourRange.divider}%
        </span>
      </div>
    );
  }

  function renderBlueSlider() {
    return (
      <div className="row form-group">
        <span className="col-3">Blue</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider blue-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={colourRange.min}
            max={colourRange.max}
            value={blueValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setBlueValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setBlueValue(colourRange.default)}
          onKeyPress={() => setBlueValue(colourRange.default)}
        >
          {(blueValue - colourRange.default) / colourRange.divider}%
        </span>
      </div>
    );
  }

  function renderRotate() {
    return (
      <div className="row form-group">
        <span className="col-3">Rotate</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={rotateRange.min}
            max={rotateRange.max}
            value={rotateValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setRotateValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setRotateValue(rotateRange.default)}
          onKeyPress={() => setRotateValue(rotateRange.default)}
        >
          {(rotateValue - rotateRange.default) / rotateRange.divider}&deg;
        </span>
      </div>
    );
  }

  function renderScale() {
    return (
      <div className="row form-group">
        <span className="col-3">Scale</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={scaleRange.min}
            max={scaleRange.max}
            value={scaleValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setScaleValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setScaleValue(scaleRange.default)}
          onKeyPress={() => setScaleValue(scaleRange.default)}
        >
          {scaleValue / scaleRange.divider}%
        </span>
      </div>
    );
  }

  function renderAspectRatio() {
    return (
      <div className="row form-group">
        <span className="col-3">Aspect</span>
        <span className="col-7">
          <Form.Control
            className="filter-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={aspectRatioRange.min}
            max={aspectRatioRange.max}
            value={aspectRatioValue}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setAspectRatioValue(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </span>
        <span
          className="col-2 text-truncate"
          role="presentation"
          onClick={() => setAspectRatioValue(aspectRatioRange.default)}
          onKeyPress={() => setAspectRatioValue(aspectRatioRange.default)}
        >
          {(aspectRatioValue - aspectRatioRange.default) /
            aspectRatioRange.divider}
        </span>
      </div>
    );
  }

  function onRotateAndScale(direction: number) {
    if (direction === 0) {
      // Left -90
      setRotateValue(1);
    } else {
      // Right +90
      setRotateValue(3);
    }

    // Calculate Required Scaling.
    const sceneWidth = props.scene.file.width ?? 1;
    const sceneHeight = props.scene.file.height ?? 1;
    const sceneAspectRatio = sceneWidth / sceneHeight;
    const sceneNewAspectRatio = sceneHeight / sceneWidth;

    const playerId = JWUtils.playerID;
    const playerVideoElement = document
      .getElementById(playerId)
      ?.getElementsByClassName("jw-video")[0];

    const playerWidth = playerVideoElement?.clientWidth ?? 1;
    const playerHeight = playerVideoElement?.clientHeight ?? 1;
    const playerAspectRation = playerWidth / playerHeight;

    // rs > ri ? (wi * hs/hi, hs) : (ws, hi * ws/wi)
    // Determine if video is currently constrained by player height or width.
    let scaledVideoHeight = 0;
    let scaledVideoWidth = 0;
    if (playerAspectRation > sceneAspectRatio) {
      // Video has it's width scaled
      // Video is constrained by it's height
      scaledVideoHeight = playerHeight;
      scaledVideoWidth = (playerHeight / sceneHeight) * sceneWidth;
    } else {
      // Video has it's height scaled
      // Video is constrained by it's width
      scaledVideoWidth = playerWidth;
      scaledVideoHeight = (playerWidth / sceneWidth) * sceneHeight;
    }

    // but now the video is rotated
    let scaleFactor = 1;
    if (playerAspectRation > sceneNewAspectRatio) {
      // Rotated video will be constrained by it's height
      // so we need to scaledVideoWidth to match the player height
      scaleFactor = playerHeight / scaledVideoWidth;
    } else {
      // Rotated video will be constrained by it's width
      // so we need to scaledVideoHeight to match the player width
      scaleFactor = playerWidth / scaledVideoHeight;
    }

    setScaleValue(scaleFactor * 100);
  }

  function renderRotateAndScale() {
    return (
      <div className="row form-group">
        <span className="col-6">
          <Button
            id="rotateAndScaleLeft"
            variant="primary"
            type="button"
            onClick={() => onRotateAndScale(0)}
          >
            Rotate Left & Scale
          </Button>
        </span>
        <span className="col-6">
          <Button
            id="rotateAndScaleRight"
            variant="primary"
            type="button"
            onClick={() => onRotateAndScale(1)}
          >
            Rotate Right & Scale
          </Button>
        </span>
      </div>
    );
  }

  function onResetFilters() {
    setContrastValue(contrastRange.default);
    setBrightnessValue(brightnessRange.default);
    setGammaValue(gammaRange.default);
    setSaturateValue(saturateRange.default);
    setHueRotateValue(hueRotateRange.default);
    setWhiteBalanceValue(whiteBalanceRange.default);
    setRedValue(colourRange.default);
    setGreenValue(colourRange.default);
    setBlueValue(colourRange.default);
    setBlurValue(blurRange.default);
  }

  function onResetTransforms() {
    setScaleValue(scaleRange.default);
    setRotateValue(rotateRange.default);
    setAspectRatioValue(aspectRatioRange.default);
  }

  function renderResetButton() {
    return (
      <div className="row form-group">
        <span className="col-6">
          <Button
            id="resetFilters"
            variant="primary"
            type="button"
            onClick={() => onResetFilters()}
          >
            Reset Filters
          </Button>
        </span>
        <span className="col-6">
          <Button
            id="resetTransforms"
            variant="secondary"
            type="button"
            onClick={() => onResetTransforms()}
          >
            Reset Transforms
          </Button>
        </span>
      </div>
    );
  }

  function renderFilterContainer() {
    return <div id="video-filter-container" />;
  }

  // On render update video style.
  updateVideoFilters();
  updateVideoStyle();

  return (
    <div className="container scene-video-filter">
      <div className="row form-group">
        <span className="col-12">
          <h5>Filters</h5>
        </span>
      </div>
      {renderBrightness()}
      {renderContrast()}
      {renderGammaSlider()}
      {renderSaturate()}
      {renderHueRotateSlider()}
      {renderWhiteBalance()}
      {renderRedSlider()}
      {renderGreenSlider()}
      {renderBlueSlider()}
      {renderBlur()}
      <div className="row form-group">
        <span className="col-12">
          <h5>Transforms</h5>
        </span>
      </div>
      {renderRotate()}
      {renderScale()}
      {renderAspectRatio()}
      <div className="row form-group">
        <span className="col-12">
          <h5>Actions</h5>
        </span>
      </div>
      {renderRotateAndScale()}
      {renderResetButton()}
      <div className="row form-group">
        <span className="col-12">
          <i>
            Note: Filters is a beta level feature and all functions may not work
            in all browsers equally. Chrome seems best but isn&apos;t without
            issues.
          </i>
        </span>
      </div>
      {renderFilterContainer()}
    </div>
  );
};
