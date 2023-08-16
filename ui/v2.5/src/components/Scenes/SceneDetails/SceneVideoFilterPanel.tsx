import React, {
  useEffect,
  useState,
  useMemo,
  useRef,
  useCallback,
} from "react";
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
  updateVideoFilters,
  updateVideoStyle,
  sliderRanges,
} from "src/utils/videoFilter";
import * as yup from "yup";
import {
  useSceneFilterCreate,
  useSceneFilterUpdate,
  useSceneFilterDestroy,
} from "src/core/StashService";
import isEqual from "lodash-es/isEqual";

interface ISceneVideoFilterPanelProps {
  scene: GQL.SceneDataFragment;
  filter?: GQL.SceneFilterDataFragment;
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

export const SceneVideoFilterPanel: React.FC<ISceneVideoFilterPanelProps> = ({
  filter,
  isVisible,
  ...props
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const isNew = props.scene.scene_filters[0] === undefined;
  const [sceneFilterCreate] = useSceneFilterCreate();
  const [sceneFilterUpdate] = useSceneFilterUpdate();
  const [sceneFilterDestroy] = useSceneFilterDestroy();

  const [contrastValue, setContrastValue] = useState(
    props.scene.scene_filters[0]?.contrast ?? sliderRanges.contrastRange.default
  );
  const [brightnessValue, setBrightnessValue] = useState(
    props.scene.scene_filters[0]?.brightness ??
      sliderRanges.brightnessRange.default
  );
  const [gammaValue, setGammaValue] = useState(
    props.scene.scene_filters[0]?.gamma ?? sliderRanges.gammaRange.default
  );
  const [saturateValue, setSaturateValue] = useState(
    props.scene.scene_filters[0]?.saturate ?? sliderRanges.saturateRange.default
  );
  const [hueRotateValue, setHueRotateValue] = useState(
    props.scene.scene_filters[0]?.hue_rotate ??
      sliderRanges.hueRotateRange.default
  );
  const [warmthValue, setWarmthValue] = useState(
    props.scene.scene_filters[0]?.warmth ?? sliderRanges.warmthRange.default
  );
  const [redValue, setRedValue] = useState(
    props.scene.scene_filters[0]?.red ?? sliderRanges.colourRange.default
  );
  const [greenValue, setGreenValue] = useState(
    props.scene.scene_filters[0]?.green ?? sliderRanges.colourRange.default
  );
  const [blueValue, setBlueValue] = useState(
    props.scene.scene_filters[0]?.blue ?? sliderRanges.colourRange.default
  );
  const [blurValue, setBlurValue] = useState(
    props.scene.scene_filters[0]?.blur ?? sliderRanges.blurRange.default
  );
  const [rotateValue, setRotateValue] = useState(
    props.scene.scene_filters[0]?.rotate ?? sliderRanges.rotateRange.default
  );
  const [scaleValue, setScaleValue] = useState(
    props.scene.scene_filters[0]?.scale ?? sliderRanges.scaleRange.default
  );
  const [aspectRatioValue, setAspectRatioValue] = useState(
    props.scene.scene_filters[0]?.aspect_ratio ??
      sliderRanges.aspectRatioRange.default
  );

  console.log("isNew: " + isNew);

  const setDefaultFilterValues = useCallback(() => {
    setContrastValue(sliderRanges.contrastRange.default);
    setBrightnessValue(sliderRanges.brightnessRange.default);
    setGammaValue(sliderRanges.contrastRange.default);
    setSaturateValue(sliderRanges.saturateRange.default);
    setHueRotateValue(sliderRanges.hueRotateRange.default);
    setWarmthValue(sliderRanges.warmthRange.default);
    setRedValue(sliderRanges.contrastRange.default);
    setGreenValue(sliderRanges.colourRange.default);
    setBlueValue(sliderRanges.colourRange.default);
    setBlurValue(sliderRanges.blurRange.default);
    setScaleValue(sliderRanges.scaleRange.default);
    setRotateValue(sliderRanges.rotateRange.default);
    setAspectRatioValue(sliderRanges.aspectRatioRange.default);
  }, []);

  const setFilterFilterValues = useCallback(() => {
    setContrastValue(props.scene.scene_filters[0].contrast);
    setBrightnessValue(props.scene.scene_filters[0].brightness);
    setGammaValue(props.scene.scene_filters[0].gamma);
    setSaturateValue(props.scene.scene_filters[0].saturate);
    setHueRotateValue(props.scene.scene_filters[0].hue_rotate);
    setWarmthValue(props.scene.scene_filters[0].warmth);
    setRedValue(props.scene.scene_filters[0].red);
    setGreenValue(props.scene.scene_filters[0].green);
    setBlueValue(props.scene.scene_filters[0].blue);
    setBlurValue(props.scene.scene_filters[0].blur);
    setScaleValue(props.scene.scene_filters[0].scale);
    setRotateValue(props.scene.scene_filters[0].rotate);
    setAspectRatioValue(props.scene.scene_filters[0].aspect_ratio);
  }, [props.scene.scene_filters]);

  const prevVideoFiltersRef = useRef<boolean>(false);
  const prevDefaultFiltersRef = useRef<boolean>(false);

  const setCurrentOrPreviousVideoState = (
    isFilter: boolean,
    isDefaultFilter: boolean
  ) => {
    prevVideoFiltersRef.current = isFilter;
    prevDefaultFiltersRef.current = isDefaultFilter;
  };

  // Video filter changing between scenes
  // Initialize values once when component mounts
  useEffect(() => {
    if (prevVideoFiltersRef.current && props.scene.scene_filters.length === 0) {
      setCurrentOrPreviousVideoState(false, true);
      setDefaultFilterValues();
    } else if (
      prevVideoFiltersRef.current &&
      props.scene.scene_filters.length > 0
    ) {
      setCurrentOrPreviousVideoState(true, false);
      setFilterFilterValues();
    } else if (
      prevDefaultFiltersRef.current &&
      props.scene.scene_filters.length === 0
    ) {
      setCurrentOrPreviousVideoState(false, true);
    } else if (
      prevDefaultFiltersRef.current &&
      props.scene.scene_filters.length > 0
    ) {
      setCurrentOrPreviousVideoState(true, false);
      setFilterFilterValues();
    } else if (
      !prevDefaultFiltersRef.current &&
      props.scene.scene_filters.length > 0
    ) {
      setCurrentOrPreviousVideoState(true, false);
      setFilterFilterValues();
    }
  }, [props.scene, setFilterFilterValues, setDefaultFilterValues]);

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
      contrast: filter?.contrast ?? 100,
      brightness: filter?.brightness ?? 100,
      gamma: filter?.gamma ?? 100,
      saturate: filter?.saturate ?? 100,
      hue_rotate: filter?.hue_rotate ?? 0,
      warmth: filter?.warmth ?? 100,
      red: filter?.red ?? 100,
      green: filter?.green ?? 100,
      blue: filter?.blue ?? 100,
      blur: filter?.blur ?? 0,
      rotate: filter?.rotate ?? 2,
      scale: filter?.scale ?? 100,
      aspect_ratio: filter?.aspect_ratio ?? 150,
    }),
    [filter]
  );

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSave(values),
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
          onDelete();
        }
      });

      return () => {
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");
      };
    }
  });

  async function onDelete() {
    setIsLoading(true);
    try {
      const result = await sceneFilterDestroy({
        variables: { id: props.scene.scene_filters[0].id },
      });
      if (result) {
        Toast.success({
          content: "Scene filter successfully deleted",
        });
        onResetTransforms();
        onResetFilters();
      }
      return;
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onSave(input: InputValues) {
    try {
      if (isNew == true) {
        const result = await sceneFilterCreate({
          variables: {
            scene_id: props.scene.id,
            ...input,
          },
        });
        if (result) {
          Toast.success({
            content: "Scene filter successfully created",
          });
          onResetTransforms();
          onResetFilters();
        }
      } else {
        const result = await sceneFilterUpdate({
          variables: {
            scene_id: props.scene.id,
            id: props.scene.scene_filters[0].id,
            ...input,
          },
        });
        if (result) {
          Toast.success({
            content: "Scene filter successfully updated",
          });
        }
      }
    } catch (e) {
      Toast.error(e);
    }
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
              disabled={
                (!isNew && !formik.dirty) || !isEqual(formik.errors, {})
              }
              onClick={() => formik.submitForm()}
            >
              <FormattedMessage id="actions.save" />
            </Button>
            <Button
              className="edit-button"
              variant="danger"
              onClick={() => onDelete()}
              disabled={isNew}
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
  updateVideoFilters(gammaValue, redValue, greenValue, blueValue, warmthValue);
  updateVideoStyle(
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
  );
};
export default SceneVideoFilterPanel;
