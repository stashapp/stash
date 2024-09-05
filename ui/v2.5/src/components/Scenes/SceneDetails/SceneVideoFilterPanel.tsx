import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { VIDEO_PLAYER_ID } from "src/components/ScenePlayer/util";
import * as GQL from "src/core/generated-graphql";

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

function getMatrixValue(value: number, range: SliderRange) {
  return (value - range.default) / range.divider;
}

interface ISliderProps {
  title: string;
  className?: string;
  range: SliderRange;
  value: number;
  setValue: (value: React.SetStateAction<number>) => void;
  displayValue: string;
}

const Slider: React.FC<ISliderProps> = (sliderProps: ISliderProps) => {
  return (
    <div className="row form-group">
      <span className="col-sm-3">{sliderProps.title}</span>
      <span className="col-sm-7">
        <Form.Control
          className={`filter-slider d-inline-flex ml-sm-3 ${sliderProps.className}`}
          type="range"
          min={sliderProps.range.min}
          max={sliderProps.range.max}
          value={sliderProps.value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            sliderProps.setValue(Number.parseInt(e.currentTarget.value, 10))
          }
        />
      </span>
      <span
        className="col-sm-2 filter-slider-value"
        role="presentation"
        onClick={() => sliderProps.setValue(sliderProps.range.default)}
        onKeyPress={() => sliderProps.setValue(sliderProps.range.default)}
      >
        <TruncatedText text={sliderProps.displayValue} />
      </span>
    </div>
  );
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
    divider: 100,
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

  const intl = useIntl();

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

  function updateVideoStyle() {
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

      if (playerVideoElement.tagName == "CANVAS") {
        styleString += "; width: 100%; height: 100%; position: absolute; top:0";
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

      const wbMatrixValue = getMatrixValue(
        whiteBalanceValue,
        whiteBalanceRange
      );

      feColorMatrix.setAttribute(
        "values",
        `${
          1 + wbMatrixValue + getMatrixValue(redValue, colourRange)
        } 0 0 0 0   0 ${
          1.0 + getMatrixValue(greenValue, colourRange)
        } 0 0 0   0 0 ${
          1 - wbMatrixValue + getMatrixValue(blueValue, colourRange)
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

  function onRotateAndScale(direction: number) {
    if (direction === 0) {
      // Left -90
      setRotateValue(1);
    } else {
      // Right +90
      setRotateValue(3);
    }

    const file =
      props.scene.files.length > 0 ? props.scene.files[0] : undefined;

    // Calculate Required Scaling.
    const sceneWidth = file?.width ?? 1;
    const sceneHeight = file?.height ?? 1;
    const sceneAspectRatio = sceneWidth / sceneHeight;
    const sceneNewAspectRatio = sceneHeight / sceneWidth;

    const playerVideoElement = document.getElementById(VIDEO_PLAYER_ID);
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
            <FormattedMessage id="effect_filters.rotate_left_and_scale" />
          </Button>
        </span>
        <span className="col-6">
          <Button
            id="rotateAndScaleRight"
            variant="primary"
            type="button"
            onClick={() => onRotateAndScale(1)}
          >
            <FormattedMessage id="effect_filters.rotate_right_and_scale" />
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
            <FormattedMessage id="effect_filters.reset_filters" />
          </Button>
        </span>
        <span className="col-6">
          <Button
            id="resetTransforms"
            variant="secondary"
            type="button"
            onClick={() => onResetTransforms()}
          >
            <FormattedMessage id="effect_filters.reset_transforms" />
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
          <h5>
            <FormattedMessage id="effect_filters.name" />
          </h5>
        </span>
      </div>
      <Slider
        title={intl.formatMessage({ id: "effect_filters.brightness" })}
        className="brightness-slider"
        range={brightnessRange}
        value={brightnessValue}
        setValue={setBrightnessValue}
        displayValue={`${brightnessValue / brightnessRange.divider}%`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.contrast" })}
        className="contrast-slider"
        range={contrastRange}
        value={contrastValue}
        setValue={setContrastValue}
        displayValue={`${contrastValue / brightnessRange.divider}%`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.gamma" })}
        className="gamma-slider"
        range={gammaRange}
        value={gammaValue}
        setValue={setGammaValue}
        displayValue={`${
          (gammaValue - gammaRange.default) / gammaRange.divider
        }`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.saturation" })}
        className="saturation-slider"
        range={saturateRange}
        value={saturateValue}
        setValue={setSaturateValue}
        displayValue={`${saturateValue / saturateRange.divider}%`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.hue" })}
        className="hue-rotate-slider"
        range={hueRotateRange}
        value={hueRotateValue}
        setValue={setHueRotateValue}
        displayValue={`${hueRotateValue / hueRotateRange.divider}\xB0`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.warmth" })}
        className="white-balance-slider"
        range={whiteBalanceRange}
        value={whiteBalanceValue}
        setValue={setWhiteBalanceValue}
        displayValue={`${
          (whiteBalanceValue - whiteBalanceRange.default) /
          whiteBalanceRange.divider
        }`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.red" })}
        className="red-slider"
        range={colourRange}
        value={redValue}
        setValue={setRedValue}
        displayValue={`${redValue}%`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.green" })}
        className="green-slider"
        range={colourRange}
        value={greenValue}
        setValue={setGreenValue}
        displayValue={`${greenValue}%`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.blue" })}
        className="blue-slider"
        range={colourRange}
        value={blueValue}
        setValue={setBlueValue}
        displayValue={`${blueValue}%`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.blur" })}
        range={blurRange}
        value={blurValue}
        setValue={setBlurValue}
        displayValue={`${blurValue / blurRange.divider}px`}
      />

      <div className="row form-group">
        <span className="col-12">
          <h5>
            <FormattedMessage id="effect_filters.name_transforms" />
          </h5>
        </span>
      </div>
      <Slider
        title={intl.formatMessage({ id: "effect_filters.rotate" })}
        range={rotateRange}
        value={rotateValue}
        setValue={setRotateValue}
        displayValue={`${
          (rotateValue - rotateRange.default) / rotateRange.divider
        }\xB0`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.scale" })}
        range={scaleRange}
        value={scaleValue}
        setValue={setScaleValue}
        displayValue={`${scaleValue / scaleRange.divider}%`}
      />
      <Slider
        title={intl.formatMessage({ id: "effect_filters.aspect" })}
        range={aspectRatioRange}
        value={aspectRatioValue}
        setValue={setAspectRatioValue}
        displayValue={`${
          (aspectRatioValue - aspectRatioRange.default) /
          aspectRatioRange.divider
        }`}
      />
      <div className="row form-group">
        <span className="col-12">
          <h5>
            <FormattedMessage id="actions_name" />
          </h5>
        </span>
      </div>
      {renderRotateAndScale()}
      {renderResetButton()}
      {renderFilterContainer()}
    </div>
  );
};

export default SceneVideoFilterPanel;
