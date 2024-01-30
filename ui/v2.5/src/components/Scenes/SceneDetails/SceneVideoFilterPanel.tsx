import React, { useEffect, useState, useMemo, useRef } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { VIDEO_PLAYER_ID } from "src/components/ScenePlayer/util";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { useFormik } from "formik";
import {
  updateSceneFilters,
  updateSceneFiltersStyle,
  sliderRanges,
} from "src/utils/sceneFilters";
import * as yup from "yup";
import { useSceneUpdate } from "src/core/StashService";
import isEqual from "lodash-es/isEqual";

interface ISceneVideoFilterPanelProps {
  scene: GQL.SceneDataFragment;
  isVisible: boolean;
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

interface IFilters {
  contrast?: number;
  brightness?: number;
  gamma?: number;
  saturate?: number;
  hue_rotate?: number;
  warmth?: number;
  red?: number;
  green?: number;
  blue?: number;
  blur?: number;
  rotate?: number;
  scale?: number;
  aspect_ratio?: number;
}

export const SceneVideoFilterPanel: React.FC<ISceneVideoFilterPanelProps> = ({
  isVisible,
  ...props
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const [sceneFilterUpdate] = useSceneUpdate();
  const hasFiltersRef = useRef<boolean>();
  const [filters, setFilters] = useState<IFilters | null>(null);

  const [contrastValue, setContrastValue] = useState(
    sliderRanges.contrastRange.default
  );
  const [brightnessValue, setBrightnessValue] = useState(
    sliderRanges.brightnessRange.default
  );
  const [gammaValue, setGammaValue] = useState(sliderRanges.gammaRange.default);
  const [saturateValue, setSaturateValue] = useState(
    filters?.saturate ?? sliderRanges.saturateRange.default
  );
  const [hueRotateValue, setHueRotateValue] = useState(
    sliderRanges.hueRotateRange.default
  );
  const [warmthValue, setWarmthValue] = useState(
    sliderRanges.warmthRange.default
  );
  const [redValue, setRedValue] = useState(sliderRanges.colourRange.default);
  const [greenValue, setGreenValue] = useState(
    sliderRanges.colourRange.default
  );
  const [blueValue, setBlueValue] = useState(sliderRanges.colourRange.default);
  const [blurValue, setBlurValue] = useState(sliderRanges.blurRange.default);
  const [rotateValue, setRotateValue] = useState(
    sliderRanges.rotateRange.default
  );
  const [scaleValue, setScaleValue] = useState(sliderRanges.scaleRange.default);
  const [aspectRatioValue, setAspectRatioValue] = useState(
    sliderRanges.aspectRatioRange.default
  );

  useEffect(() => {
    if (props.scene.filters && props.scene.filters.length > 0) {
      hasFiltersRef.current = false;
      const parsedFilters = JSON.parse(props.scene.filters);
      setFilters(parsedFilters);
      setContrastValue(parsedFilters.contrast);
      setBrightnessValue(parsedFilters.brightness);
      setGammaValue(parsedFilters.gamma);
      setHueRotateValue(parsedFilters.hue_rotate);
      setWarmthValue(parsedFilters.warmth);
      setRedValue(parsedFilters.red);
      setGreenValue(parsedFilters.green);
      setBlueValue(parsedFilters.blue);
      setBlurValue(parsedFilters.blur);
      setRotateValue(parsedFilters.rotate);
      setScaleValue(parsedFilters.scale);
      setAspectRatioValue(parsedFilters.aspect_ratio);
      setContrastValue(parsedFilters.contrast);
    } else {
      hasFiltersRef.current = true;
      setFilters(null);
      setContrastValue(sliderRanges.contrastRange.default);
      setBrightnessValue(sliderRanges.brightnessRange.default);
      setGammaValue(sliderRanges.gammaRange.default);
      setHueRotateValue(sliderRanges.hueRotateRange.default);
      setWarmthValue(sliderRanges.warmthRange.default);
      setRedValue(sliderRanges.colourRange.default);
      setGreenValue(sliderRanges.colourRange.default);
      setBlueValue(sliderRanges.colourRange.default);
      setBlurValue(sliderRanges.blurRange.default);
      setRotateValue(sliderRanges.rotateRange.default);
      setScaleValue(sliderRanges.scaleRange.default);
      setAspectRatioValue(sliderRanges.aspectRatioRange.default);
      setContrastValue(sliderRanges.contrastRange.default);
    }
  }, [props.scene]);

  const schema = yup.object({
    contrast: yup.number().required(),
    brightness: yup.number().required(),
    gamma: yup.number().required(),
    saturate: yup.number().required(),
    hue_rotate: yup.number().required(),
    warmth: yup.number().required(),
    red: yup.number().required(),
    green: yup.number().required(),
    blue: yup.number().required(),
    blur: yup.number().required(),
    rotate: yup.number().required(),
    scale: yup.number().required(),
    aspect_ratio: yup.number().required(),
  });

  const initialValues = useMemo(
    () => ({
      contrast: filters?.contrast ?? 100,
      brightness: filters?.brightness ?? 100,
      gamma: filters?.gamma ?? 100,
      saturate: filters?.saturate ?? 100,
      hue_rotate: filters?.hue_rotate ?? 0,
      warmth: filters?.warmth ?? 100,
      red: filters?.red ?? 100,
      green: filters?.green ?? 100,
      blue: filters?.blue ?? 100,
      blur: filters?.blur ?? 0,
      rotate: filters?.rotate ?? 2,
      scale: filters?.scale ?? 100,
      aspect_ratio: filters?.aspect_ratio ?? 150,
    }),
    [filters]
  );

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSceneUpdate(values),
  });

  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        if (formik.dirty) {
          formik.submitForm();
        }
      });
      Mousetrap.bind("d d", () => {
        if (isVisible) {
          onSceneUpdate(null);
        }
      });

      return () => {
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");
      };
    }
  });

  async function onSceneUpdate(input: InputValues | null) {
    setIsLoading(true);
    try {
      const result = await sceneFilterUpdate({
        variables: {
          input: {
            id: props.scene.id,
            filters: input
              ? JSON.stringify({
                  contrast: input.contrast,
                  brightness: input.brightness,
                  gamma: input.gamma,
                  saturate: input.saturate,
                  hue_rotate: input.hue_rotate,
                  warmth: input.warmth,
                  red: input.red,
                  green: input.green,
                  blue: input.blue,
                  blur: input.blur,
                  rotate: input.rotate,
                  scale: input.scale,
                  aspect_ratio: input.aspect_ratio,
                })
              : null,
          },
        },
      });

      if (result) {
        if (input) {
          Toast.success(
            intl.formatMessage({ id: "toast.scene_video_filter_saved" })
          );
        } else {
          Toast.success(
            intl.formatMessage({ id: "toast.scene_video_filter_deleted" })
          );
          onResetTransforms();
          onResetFilters();
        }
      }
      return;
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onDeleteAction() {
    await onSceneUpdate(null);
  }

  interface ISliderProps {
    title: string;
    className?: string;
    range: SliderRange;
    value: number;
    // setValue: (value: number) => void;
    setValue: (value: React.SetStateAction<number>) => void;
    displayValue: string;
    fieldName: string;
  }

  function renderSlider(sliderProps: ISliderProps) {
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
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
              const newValue = Number.parseInt(e.currentTarget.value, 10);
              sliderProps.setValue(newValue);
              formik.setFieldValue(sliderProps.fieldName, newValue);
            }}
          />
        </span>
        <span
          className="col-sm-2"
          role="presentation"
          onClick={() => sliderProps.setValue(sliderProps.range.default)}
          onKeyPress={() => sliderProps.setValue(sliderProps.range.default)}
        >
          <TruncatedText text={sliderProps.displayValue} />
        </span>
      </div>
    );
  }

  function renderBlur() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.blur" }),
      range: sliderRanges.blurRange,
      value: blurValue,
      setValue: setBlurValue,
      displayValue: `${blurValue / sliderRanges.blurRange.divider}px`,
      fieldName: "blur",
    });
  }

  function renderContrast() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.contrast" }),
      className: "contrast-slider",
      range: sliderRanges.contrastRange,
      value: contrastValue,
      setValue: setContrastValue,
      displayValue: `${contrastValue / sliderRanges.brightnessRange.divider}%`,
      fieldName: "contrast",
    });
  }

  function renderBrightness() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.brightness" }),
      className: "brightness-slider",
      range: sliderRanges.brightnessRange,
      value: brightnessValue,
      setValue: setBrightnessValue,
      displayValue: `${
        brightnessValue / sliderRanges.brightnessRange.divider
      }%`,
      fieldName: "brightness",
    });
  }

  function renderGammaSlider() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.gamma" }),
      className: "gamma-slider",
      range: sliderRanges.gammaRange,
      value: gammaValue,
      setValue: setGammaValue,
      displayValue: `${
        (gammaValue - sliderRanges.gammaRange.default) /
        sliderRanges.gammaRange.divider
      }`,
      fieldName: "gamma",
    });
  }

  function renderSaturate() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.saturation" }),
      className: "saturation-slider",
      range: sliderRanges.saturateRange,
      value: saturateValue,
      setValue: setSaturateValue,
      displayValue: `${saturateValue / sliderRanges.saturateRange.divider}%`,
      fieldName: "saturate",
    });
  }

  function renderHueRotateSlider() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.hue" }),
      className: "hue-rotate-slider",
      range: sliderRanges.hueRotateRange,
      value: hueRotateValue,
      setValue: setHueRotateValue,
      displayValue: `${
        hueRotateValue / sliderRanges.hueRotateRange.divider
      }\xB0`,
      fieldName: "hue_rotate",
    });
  }

  function renderWarmth() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.warmth" }),
      className: "white-balance-slider",
      range: sliderRanges.warmthRange,
      value: warmthValue,
      setValue: setWarmthValue,
      displayValue: `${
        (warmthValue - sliderRanges.warmthRange.default) /
        sliderRanges.warmthRange.divider
      }`,
      fieldName: "warmth",
    });
  }

  function renderRedSlider() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.red" }),
      className: "red-slider",
      range: sliderRanges.colourRange,
      value: redValue,
      setValue: setRedValue,
      displayValue: `${
        (redValue - sliderRanges.colourRange.default) /
        sliderRanges.colourRange.divider
      }%`,
      fieldName: "red",
    });
  }

  function renderGreenSlider() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.green" }),
      className: "green-slider",
      range: sliderRanges.colourRange,
      value: greenValue,
      setValue: setGreenValue,
      displayValue: `${
        (greenValue - sliderRanges.colourRange.default) /
        sliderRanges.colourRange.divider
      }%`,
      fieldName: "green",
    });
  }

  function renderBlueSlider() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.blue" }),
      className: "blue-slider",
      range: sliderRanges.colourRange,
      value: blueValue,
      setValue: setBlueValue,
      displayValue: `${
        (blueValue - sliderRanges.colourRange.default) /
        sliderRanges.colourRange.divider
      }%`,
      fieldName: "blue",
    });
  }

  function renderRotate() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.rotate" }),
      range: sliderRanges.rotateRange,
      value: rotateValue,
      setValue: setRotateValue,
      displayValue: `${
        (rotateValue - sliderRanges.rotateRange.default) /
        sliderRanges.rotateRange.divider
      }\xB0`,
      fieldName: "rotate",
    });
  }

  function renderScale() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.scale" }),
      range: sliderRanges.scaleRange,
      value: scaleValue,
      setValue: setScaleValue,
      displayValue: `${scaleValue / sliderRanges.scaleRange.divider}%`,
      fieldName: "scale",
    });
  }

  function renderAspectRatio() {
    return renderSlider({
      title: intl.formatMessage({ id: "effect_filters.aspect" }),
      range: sliderRanges.aspectRatioRange,
      value: aspectRatioValue,
      setValue: setAspectRatioValue,
      displayValue: `${
        (aspectRatioValue - sliderRanges.aspectRatioRange.default) /
        sliderRanges.aspectRatioRange.divider
      }`,
      fieldName: "aspect_ratio",
    });
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
    setContrastValue(sliderRanges.contrastRange.default);
    setBrightnessValue(sliderRanges.brightnessRange.default);
    setGammaValue(sliderRanges.gammaRange.default);
    setSaturateValue(sliderRanges.saturateRange.default);
    setHueRotateValue(sliderRanges.hueRotateRange.default);
    setWarmthValue(sliderRanges.warmthRange.default);
    setRedValue(sliderRanges.colourRange.default);
    setGreenValue(sliderRanges.colourRange.default);
    setBlueValue(sliderRanges.colourRange.default);
    setBlurValue(sliderRanges.blurRange.default);
  }

  function onResetTransforms() {
    setScaleValue(sliderRanges.scaleRange.default);
    setRotateValue(sliderRanges.rotateRange.default);
    setAspectRatioValue(sliderRanges.aspectRatioRange.default);
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

  function renderSaveButtons() {
    return (
      <>
        <div className="form-container edit-buttons-container row px-3 pt-3">
          <div className="edit-buttons mb-3 pl-0">
            <Button
              className="edit-button"
              variant="primary"
              disabled={!formik.dirty || !isEqual(formik.errors, {})}
              onClick={() => onSceneUpdate(formik.values)}
            >
              <FormattedMessage id="actions.save" />
            </Button>
            <Button
              className="edit-button"
              variant="danger"
              onClick={onDeleteAction}
              disabled={hasFiltersRef.current}
            >
              <FormattedMessage id="actions.delete" />
            </Button>
          </div>
        </div>
      </>
    );
  }

  function renderFilterContainer() {
    return <div id="video-filter-container" />;
  }

  if (isLoading) return <LoadingIndicator />;

  // On render update video style.
  updateSceneFilters(gammaValue, redValue, greenValue, blueValue, warmthValue);
  updateSceneFiltersStyle(
    aspectRatioValue,
    blurValue,
    brightnessValue,
    contrastValue,
    gammaValue,
    hueRotateValue,
    redValue,
    greenValue,
    blueValue,
    rotateValue,
    saturateValue,
    scaleValue,
    warmthValue
  );

  return (
    <div className="scene-filters-panel">
      <div className="scene-video-filter form-container">
        <Form noValidate onSubmit={formik.handleSubmit}>
          {renderSaveButtons()}
          <div className="row form-group">
            <span className="col-12">
              <h5>
                <FormattedMessage id="effect_filters.name" />
              </h5>
            </span>
          </div>
          {renderBrightness()}
          {renderContrast()}
          {renderGammaSlider()}
          {renderSaturate()}
          {renderHueRotateSlider()}
          {renderWarmth()}
          {renderRedSlider()}
          {renderGreenSlider()}
          {renderBlueSlider()}
          {renderBlur()}
          <div className="row form-group">
            <span className="col-12">
              <h5>
                <FormattedMessage id="effect_filters.name_transforms" />
              </h5>
            </span>
          </div>
          {renderRotate()}
          {renderScale()}
          {renderAspectRatio()}
        </Form>
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
    </div>
  );
};

export default SceneVideoFilterPanel;
