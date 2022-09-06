import React, { useCallback, useEffect, useRef, useState } from "react";
import {
  Button,
  Col,
  InputGroup,
  Overlay,
  Popover,
  Form,
  Row,
} from "react-bootstrap";


import { Icon } from "src/components/Shared";
import { useInterval, usePageVisibility, useToast } from "src/hooks";
import { FormattedMessage, useIntl } from "react-intl";
import { ConfigurationContext } from "../Config";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { OCounterButton } from "src/components/Scenes/SceneDetails/OCounterButton";
import {
  useImageUpdate,
  mutateImageIncrementO,
  mutateImageDecrementO,
  mutateImageResetO,
} from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { useInterfaceLocalForage } from "../LocalForage";
import { imageLightboxDisplayModeIntlMap } from "src/core/enums";
import { ILightboxImage } from "./types";
import {
  faArrowLeft,
  faArrowRight,
  faChevronLeft,
  faChevronRight,
  faCog,
  faExpand,
  faPause,
  faPlay,
  faSearchMinus,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";

import PhotoSwipeLightbox from 'photoswipe/lightbox';
import 'photoswipe/style.css';

const CLASSNAME = "Lightbox";

const DEFAULT_SLIDESHOW_DELAY = 5000;
const SECONDS_TO_MS = 1000;
const MIN_VALID_INTERVAL_SECONDS = 1;

interface IProps {
  images: ILightboxImage[];
  isVisible: boolean;
  initialIndex?: number;
  showNavigation: boolean;
  slideshowEnabled?: boolean;
  pageHeader?: string;
  pageCallback?: (direction: number) => void;
  hide: () => void;
}

export const LightboxComponent: React.FC<IProps> = ({
  images,
  initialIndex = 0,
  showNavigation,
  slideshowEnabled = false,
}) => {

  const Toast = useToast();
  const intl = useIntl();
  const { configuration: config } = React.useContext(ConfigurationContext);
  const [
    interfaceLocalForage,
    setInterfaceLocalForage,
  ] = useInterfaceLocalForage();
  const [updateImage] = useImageUpdate();
  const lightboxSettings = interfaceLocalForage.data?.imageLightbox;

  function setLightboxSettings(v: Partial<GQL.ConfigImageLightboxInput>) {
    setInterfaceLocalForage((prev) => {
      return {
        ...prev,
        imageLightbox: {
          ...prev.imageLightbox,
          ...v,
        },
      };
    });
  }

  function setScaleUp(value: boolean) {
    setLightboxSettings({ scaleUp: value });
  }

  function setResetZoomOnNav(v: boolean) {
    setLightboxSettings({ resetZoomOnNav: v });
  }

  function setScrollMode(v: GQL.ImageLightboxScrollMode) {
    setLightboxSettings({ scrollMode: v });
  }

  const configuredDelay = config?.interface.imageLightbox.slideshowDelay
    ? config.interface.imageLightbox.slideshowDelay * SECONDS_TO_MS
    : undefined;

  const savedDelay = lightboxSettings?.slideshowDelay
    ? lightboxSettings.slideshowDelay * SECONDS_TO_MS
    : undefined;

  const slideshowDelay =
    savedDelay ?? configuredDelay ?? DEFAULT_SLIDESHOW_DELAY;

  const scrollAttemptsBeforeChange = Math.max(
    0,
    config?.interface.imageLightbox.scrollAttemptsBeforeChange ?? 0
  );

  function setSlideshowDelay(v: number) {
    setLightboxSettings({ slideshowDelay: v });
  }

  const displayMode =
    lightboxSettings?.displayMode ?? GQL.ImageLightboxDisplayMode.FitXy;
  const oldDisplayMode = useRef(displayMode);

  function setDisplayMode(v: GQL.ImageLightboxDisplayMode) {
    setLightboxSettings({ displayMode: v });
  }

  // slideshowInterval is used for controlling the logic
  // displaySlideshowInterval is for display purposes only
  // keeping them separate and independant allows us to handle the logic however we want
  // while still displaying something that makes sense to the user
  const [slideshowInterval, setSlideshowInterval] = useState<number | null>(
    null
  );

  const [
    displayedSlideshowInterval,
    setDisplayedSlideshowInterval,
  ] = useState<string>((slideshowDelay / SECONDS_TO_MS).toString());

 

  const toggleSlideshow = useCallback(() => {
    if (slideshowInterval) {
      setSlideshowInterval(null);
    } else {
      setSlideshowInterval(slideshowDelay);
    }
  }, [slideshowInterval, slideshowDelay]);

  // stop slideshow when the page is hidden
  usePageVisibility((hidden: boolean) => {
    if (hidden) {
      setSlideshowInterval(null);
    }
  });

  const onDelayChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let numberValue = Number.parseInt(e.currentTarget.value, 10);
    setDisplayedSlideshowInterval(e.currentTarget.value);

    // Without this exception, the blocking of updates for invalid values is even weirder
    if (e.currentTarget.value === "-" || e.currentTarget.value === "") {
      return;
    }

    numberValue =
      numberValue >= MIN_VALID_INTERVAL_SECONDS
        ? numberValue
        : MIN_VALID_INTERVAL_SECONDS;

    setSlideshowDelay(numberValue);

    if (slideshowInterval !== null) {
      setSlideshowInterval(numberValue * SECONDS_TO_MS);
    }
  };

  // #2451: making OptionsForm an inline component means it
  // get re-rendered each time. This makes the text
  // field lose focus on input. Use function instead.
  function renderOptionsForm() {
    return (
      <>
        {slideshowEnabled ? (
          <Form.Group controlId="delay" as={Row} className="form-container">
            <Col xs={4}>
              <Form.Label className="col-form-label">
                <FormattedMessage id="dialogs.lightbox.delay" />
              </Form.Label>
            </Col>
            <Col xs={8}>
              <Form.Control
                type="number"
                className="text-input"
                min={1}
                value={displayedSlideshowInterval ?? 0}
                onChange={onDelayChange}
                size="sm"
              />
            </Col>
          </Form.Group>
        ) : undefined}

        <Form.Group controlId="displayMode" as={Row}>
          <Col xs={4}>
            <Form.Label className="col-form-label">
              <FormattedMessage id="dialogs.lightbox.display_mode.label" />
            </Form.Label>
          </Col>
          <Col xs={8}>
            <Form.Control
              as="select"
              onChange={(e) =>
                setDisplayMode(e.target.value as GQL.ImageLightboxDisplayMode)
              }
              value={displayMode}
              className="btn-secondary mx-1 mb-1"
            >
              {Array.from(imageLightboxDisplayModeIntlMap.entries()).map(
                (v) => (
                  <option key={v[0]} value={v[0]}>
                    {intl.formatMessage({
                      id: v[1],
                    })}
                  </option>
                )
              )}
            </Form.Control>
          </Col>
        </Form.Group>
        <Form.Group>
          <Form.Group controlId="scaleUp" as={Row} className="mb-1">
            <Col>
              <Form.Check
                type="checkbox"
                label={intl.formatMessage({
                  id: "dialogs.lightbox.scale_up.label",
                })}
                checked={lightboxSettings?.scaleUp ?? false}
                disabled={displayMode === GQL.ImageLightboxDisplayMode.Original}
                onChange={(v) => setScaleUp(v.currentTarget.checked)}
              />
            </Col>
          </Form.Group>
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "dialogs.lightbox.scale_up.description",
            })}
          </Form.Text>
        </Form.Group>
        <Form.Group>
          <Form.Group controlId="resetZoomOnNav" as={Row} className="mb-1">
            <Col>
              <Form.Check
                type="checkbox"
                label={intl.formatMessage({
                  id: "dialogs.lightbox.reset_zoom_on_nav",
                })}
                checked={lightboxSettings?.resetZoomOnNav ?? false}
                onChange={(v) => setResetZoomOnNav(v.currentTarget.checked)}
              />
            </Col>
          </Form.Group>
        </Form.Group>
        <Form.Group controlId="scrollMode">
          <Form.Group as={Row} className="mb-1">
            <Col xs={4}>
              <Form.Label className="col-form-label">
                <FormattedMessage id="dialogs.lightbox.scroll_mode.label" />
              </Form.Label>
            </Col>
            <Col xs={8}>
              <Form.Control
                as="select"
                onChange={(e) =>
                  setScrollMode(e.target.value as GQL.ImageLightboxScrollMode)
                }
                value={
                  lightboxSettings?.scrollMode ??
                  GQL.ImageLightboxScrollMode.Zoom
                }
                className="btn-secondary mx-1 mb-1"
              >
                <option
                  value={GQL.ImageLightboxScrollMode.Zoom}
                  key={GQL.ImageLightboxScrollMode.Zoom}
                >
                  {intl.formatMessage({
                    id: "dialogs.lightbox.scroll_mode.zoom",
                  })}
                </option>
                <option
                  value={GQL.ImageLightboxScrollMode.PanY}
                  key={GQL.ImageLightboxScrollMode.PanY}
                >
                  {intl.formatMessage({
                    id: "dialogs.lightbox.scroll_mode.pan_y",
                  })}
                </option>
              </Form.Control>
            </Col>
          </Form.Group>
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "dialogs.lightbox.scroll_mode.description",
            })}
          </Form.Text>
        </Form.Group>
      </>
    );
  }
  const currentIndex = 0;//TODO
  const currentImage: ILightboxImage | undefined = images[currentIndex];
  function setRating(v: number | null) {
    if (currentImage?.id) {
      updateImage({
        variables: {
          input: {
            id: currentImage.id,
            rating: v,
          },
        },
      });
    }
  }

  async function onIncrementClick() {
    if (currentImage?.id === undefined) return;
    try {
      await mutateImageIncrementO(currentImage.id);
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onDecrementClick() {
    if (currentImage?.id === undefined) return;
    try {
      await mutateImageDecrementO(currentImage.id);
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onResetClick() {
    if (currentImage?.id === undefined) return;
    try {
      await mutateImageResetO(currentImage?.id);
    } catch (e) {
      Toast.error(e);
    }
  }

  useEffect(() => {
    let lightbox = new PhotoSwipeLightbox({
      gallery: '#' + CLASSNAME,
      children: 'a',
      pswpModule: () => import('photoswipe'),
    });
    lightbox.init();

    return () => {
      lightbox.destroy();
    };
  }, []);

  return (
    <div className="pswp-gallery" id={CLASSNAME}>
      {images.map((image, i) => (
        <a
          href={image.paths.image as string}
          data-pswp-width={image.width}
          data-pswp-height={image.height}
          key={CLASSNAME + '-' + i}
          target="_blank"
          rel="noreferrer"
        >
          <img src={image.paths.thumbnail as string} alt="" />
        </a>
      ))}
    </div>
  );
};

export default LightboxComponent;
