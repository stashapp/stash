import React, { useCallback, useEffect, useRef, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  Button,
  Col,
  InputGroup,
  OverlayTrigger,
  Popover,
  Form,
  Row,
} from "react-bootstrap";
import cx from "classnames";
import Mousetrap from "mousetrap";
import debounce from "lodash/debounce";

import { Icon, LoadingIndicator } from "src/components/Shared";
import { useInterval, usePageVisibility } from "src/hooks";
import { useConfiguration } from "src/core/StashService";
import { FormattedMessage, useIntl } from "react-intl";
import { DisplayMode, LightboxImage } from "./LightboxImage";

const CLASSNAME = "Lightbox";
const CLASSNAME_HEADER = `${CLASSNAME}-header`;
const CLASSNAME_LEFT_SPACER = `${CLASSNAME_HEADER}-left-spacer`;
const CLASSNAME_INDICATOR = `${CLASSNAME_HEADER}-indicator`;
const CLASSNAME_OPTIONS = `${CLASSNAME_HEADER}-options`;
const CLASSNAME_OPTIONS_ICON = `${CLASSNAME_OPTIONS}-icon`;
const CLASSNAME_OPTIONS_INLINE = `${CLASSNAME_OPTIONS}-inline`;
const CLASSNAME_RIGHT = `${CLASSNAME_HEADER}-right`;
const CLASSNAME_DISPLAY = `${CLASSNAME}-display`;
const CLASSNAME_CAROUSEL = `${CLASSNAME}-carousel`;
const CLASSNAME_INSTANT = `${CLASSNAME_CAROUSEL}-instant`;
const CLASSNAME_IMAGE = `${CLASSNAME_CAROUSEL}-image`;
const CLASSNAME_NAVBUTTON = `${CLASSNAME}-navbutton`;
const CLASSNAME_NAV = `${CLASSNAME}-nav`;
const CLASSNAME_NAVIMAGE = `${CLASSNAME_NAV}-image`;
const CLASSNAME_NAVSELECTED = `${CLASSNAME_NAV}-selected`;

const DEFAULT_SLIDESHOW_DELAY = 5000;
const SECONDS_TO_MS = 1000;
const MIN_VALID_INTERVAL_SECONDS = 1;

type Image = Pick<GQL.Image, "paths">;
interface IProps {
  images: Image[];
  isVisible: boolean;
  isLoading: boolean;
  initialIndex?: number;
  showNavigation: boolean;
  slideshowEnabled?: boolean;
  pageHeader?: string;
  pageCallback?: (direction: number) => boolean;
  hide: () => void;
}

export const LightboxComponent: React.FC<IProps> = ({
  images,
  isVisible,
  isLoading,
  initialIndex = 0,
  showNavigation,
  slideshowEnabled = false,
  pageHeader,
  pageCallback,
  hide,
}) => {
  const [index, setIndex] = useState<number | null>(null);
  const oldIndex = useRef<number | null>(null);
  const [instantTransition, setInstantTransition] = useState(false);
  const [isSwitchingPage, setIsSwitchingPage] = useState(false);
  const [isFullscreen, setFullscreen] = useState(false);
  const [displayMode, setDisplayMode] = useState(DisplayMode.FIT_XY);
  const [zoomed, setZoomed] = useState(false);
  const [resetZoom, setResetZoom] = useState(false);

  const containerRef = useRef<HTMLDivElement | null>(null);
  const carouselRef = useRef<HTMLDivElement | null>(null);
  const indicatorRef = useRef<HTMLDivElement | null>(null);
  const navRef = useRef<HTMLDivElement | null>(null);
  const clearIntervalCallback = useRef<() => void>();
  const resetIntervalCallback = useRef<() => void>();

  const intl = useIntl();
  const config = useConfiguration();

  const userSelectedSlideshowDelayOrDefault =
    config?.data?.configuration.interface.slideshowDelay ??
    DEFAULT_SLIDESHOW_DELAY;

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
  ] = useState<string>(
    (userSelectedSlideshowDelayOrDefault / SECONDS_TO_MS).toString()
  );

  useEffect(() => {
    if (!isSwitchingPage) return;
    setIsSwitchingPage(false);
    if (index === -1) setIndex(images.length - 1);
  }, [isSwitchingPage, images, index]);

  const disableInstantTransition = debounce(
    () => setInstantTransition(false),
    400
  );

  const setInstant = useCallback(() => {
    setInstantTransition(true);
    disableInstantTransition();
  }, [disableInstantTransition]);

  useEffect(() => {
    // reset zoom status
    setZoomed(false);

    if (images.length < 2) return;
    if (index === oldIndex.current) return;
    if (index === null) return;

    if (carouselRef.current)
      carouselRef.current.style.left = `${index * -100}vw`;
    if (indicatorRef.current)
      indicatorRef.current.innerHTML = `${index + 1} / ${images.length}`;
    if (navRef.current) {
      const currentThumb = navRef.current.children[index + 1];
      if (currentThumb instanceof HTMLImageElement) {
        const offset =
          -1 *
          (currentThumb.offsetLeft - document.documentElement.clientWidth / 2);
        navRef.current.style.left = `${offset}px`;

        const previouslySelected = navRef.current.getElementsByClassName(
          CLASSNAME_NAVSELECTED
        )?.[0];
        if (previouslySelected)
          previouslySelected.className = CLASSNAME_NAVIMAGE;

        currentThumb.className = `${CLASSNAME_NAVIMAGE} ${CLASSNAME_NAVSELECTED}`;
      }
    }

    oldIndex.current = index;
  }, [index, images.length]);

  const selectIndex = (e: React.MouseEvent, i: number) => {
    setIndex(i);
    e.stopPropagation();
  };

  useEffect(() => {
    if (isVisible) {
      if (index === null) setIndex(initialIndex);
      document.body.style.overflow = "hidden";
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (Mousetrap as any).pause();
    }
  }, [initialIndex, isVisible, setIndex, index]);

  const toggleSlideshow = useCallback(() => {
    if (slideshowInterval) {
      setSlideshowInterval(null);
    } else if (
      displayedSlideshowInterval !== null &&
      typeof displayedSlideshowInterval !== "undefined"
    ) {
      const intervalNumber = Number.parseInt(displayedSlideshowInterval, 10);
      setSlideshowInterval(intervalNumber * SECONDS_TO_MS);
    } else {
      setSlideshowInterval(userSelectedSlideshowDelayOrDefault);
    }
  }, [
    slideshowInterval,
    userSelectedSlideshowDelayOrDefault,
    displayedSlideshowInterval,
  ]);

  usePageVisibility(() => {
    toggleSlideshow();
  });

  const close = useCallback(() => {
    if (!isFullscreen) {
      hide();
      document.body.style.overflow = "auto";
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (Mousetrap as any).unpause();
    } else document.exitFullscreen();
  }, [isFullscreen, hide]);

  const handleClose = (e: React.MouseEvent<HTMLDivElement>) => {
    const { className } = e.target as Element;
    if (className && className.includes && className.includes(CLASSNAME_IMAGE))
      close();
  };

  const handleLeft = useCallback(
    (isUserAction = true) => {
      if (isSwitchingPage || index === -1) return;

      if (index === 0) {
        if (pageCallback) {
          setIsSwitchingPage(true);
          setIndex(-1);
          // Check if calling page wants to swap page
          const repage = pageCallback(-1);
          if (!repage) {
            setIsSwitchingPage(false);
            setIndex(0);
          }
        } else setIndex(images.length - 1);
      } else setIndex((index ?? 0) - 1);

      if (isUserAction && resetIntervalCallback.current) {
        resetIntervalCallback.current();
      }
    },
    [
      images,
      setIndex,
      pageCallback,
      isSwitchingPage,
      resetIntervalCallback,
      index,
    ]
  );

  const handleRight = useCallback(
    (isUserAction = true) => {
      if (isSwitchingPage) return;

      if (index === images.length - 1) {
        if (pageCallback) {
          setIsSwitchingPage(true);
          setIndex(0);
          const repage = pageCallback?.(1);
          if (!repage) {
            setIsSwitchingPage(false);
            setIndex(images.length - 1);
          }
        } else setIndex(0);
      } else setIndex((index ?? 0) + 1);

      if (isUserAction && resetIntervalCallback.current) {
        resetIntervalCallback.current();
      }
    },
    [
      images,
      setIndex,
      pageCallback,
      isSwitchingPage,
      resetIntervalCallback,
      index,
    ]
  );

  const handleKey = useCallback(
    (e: KeyboardEvent) => {
      if (e.repeat && (e.key === "ArrowRight" || e.key === "ArrowLeft"))
        setInstant();
      if (e.key === "ArrowLeft") handleLeft();
      else if (e.key === "ArrowRight") handleRight();
      else if (e.key === "Escape") close();
    },
    [setInstant, handleLeft, handleRight, close]
  );
  const handleFullScreenChange = () => {
    if (clearIntervalCallback.current) {
      clearIntervalCallback.current();
    }
    setFullscreen(document.fullscreenElement !== null);
  };

  const [clearCallback, resetCallback] = useInterval(
    () => {
      handleRight(false);
    },
    slideshowEnabled ? slideshowInterval : null
  );

  resetIntervalCallback.current = resetCallback;
  clearIntervalCallback.current = clearCallback;

  useEffect(() => {
    if (isVisible) {
      document.addEventListener("keydown", handleKey);
      document.addEventListener("fullscreenchange", handleFullScreenChange);
    }
    return () => {
      document.removeEventListener("keydown", handleKey);
      document.removeEventListener("fullscreenchange", handleFullScreenChange);
    };
  }, [isVisible, handleKey]);

  const toggleFullscreen = useCallback(() => {
    if (!isFullscreen) containerRef.current?.requestFullscreen();
    else document.exitFullscreen();
  }, [isFullscreen]);

  const handleSlideshowIntervalChange = (newSlideshowInterval: number) => {
    setSlideshowInterval(newSlideshowInterval);
  };

  const navItems = images.map((image, i) => (
    <img
      src={image.paths.thumbnail ?? ""}
      alt=""
      className={cx(CLASSNAME_NAVIMAGE, {
        [CLASSNAME_NAVSELECTED]: i === index,
      })}
      onClick={(e: React.MouseEvent) => selectIndex(e, i)}
      role="presentation"
      loading="lazy"
      key={image.paths.thumbnail}
    />
  ));

  const onDelayChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let numberValue = Number.parseInt(e.currentTarget.value, 10);
    // Without this exception, the blocking of updates for invalid values is even weirder
    if (e.currentTarget.value === "-" || e.currentTarget.value === "") {
      setDisplayedSlideshowInterval(e.currentTarget.value);
      return;
    }

    setDisplayedSlideshowInterval(e.currentTarget.value);
    if (slideshowInterval !== null) {
      numberValue =
        numberValue >= MIN_VALID_INTERVAL_SECONDS
          ? numberValue
          : MIN_VALID_INTERVAL_SECONDS;
      handleSlideshowIntervalChange(numberValue * SECONDS_TO_MS);
    }
  };

  const currentIndex = index === null ? initialIndex : index;

  const OptionsForm: React.FC<{}> = () => (
    <>
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
            id="delay-input"
          />
        </Col>
      </Form.Group>
      <Form.Group controlId="displayMode" as={Row}>
        <Col xs={4}>
          <Form.Label className="col-form-label">
            <FormattedMessage id="dialogs.lightbox.display_mode.label" />
          </Form.Label>
        </Col>
        <Col xs={8}>
          <Form.Control
            as="select"
            onChange={(e) => setDisplayMode(e.target.value as DisplayMode)}
            value={displayMode}
            className="btn-secondary mx-1 mb-1"
          >
            <option value={DisplayMode.ORIGINAL} key={DisplayMode.ORIGINAL}>
              {intl.formatMessage({
                id: "dialogs.lightbox.display_mode.original",
              })}
            </option>
            <option value={DisplayMode.FIT_XY} key={DisplayMode.FIT_XY}>
              {intl.formatMessage({
                id: "dialogs.lightbox.display_mode.fit_to_screen",
              })}
            </option>
            <option value={DisplayMode.FIT_X} key={DisplayMode.FIT_X}>
              {intl.formatMessage({
                id: "dialogs.lightbox.display_mode.fit_horizontally",
              })}
            </option>
          </Form.Control>
        </Col>
      </Form.Group>
    </>
  );

  const optionsPopover = (
    <Popover id="slideshow-options">
      <Popover.Title>Options</Popover.Title>
      <Popover.Content>
        <OptionsForm />
      </Popover.Content>
    </Popover>
  );

  const element = isVisible ? (
    <div
      className={CLASSNAME}
      role="presentation"
      ref={containerRef}
      onClick={handleClose}
    >
      {images.length > 0 && !isLoading && !isSwitchingPage ? (
        <>
          <div className={CLASSNAME_HEADER}>
            <div className={CLASSNAME_LEFT_SPACER} />
            <div className={CLASSNAME_INDICATOR}>
              <span>{pageHeader}</span>
              <b ref={indicatorRef}>
                {`${currentIndex + 1} / ${images.length}`}
              </b>
            </div>
            <div className={CLASSNAME_RIGHT}>
              {slideshowEnabled && (
                <>
                  <div className={CLASSNAME_OPTIONS}>
                    <div className={CLASSNAME_OPTIONS_ICON}>
                      <OverlayTrigger
                        trigger="click"
                        placement="bottom"
                        overlay={optionsPopover}
                      >
                        <Button variant="link" title="Options">
                          <Icon icon="cog" />
                        </Button>
                      </OverlayTrigger>
                    </div>
                    <InputGroup className={CLASSNAME_OPTIONS_INLINE}>
                      <OptionsForm />
                    </InputGroup>
                  </div>
                  <Button
                    variant="link"
                    onClick={toggleSlideshow}
                    title="Toggle Slideshow"
                  >
                    <Icon
                      icon={slideshowInterval !== null ? "pause" : "play"}
                    />
                  </Button>
                </>
              )}
              {zoomed && (
                <Button
                  variant="link"
                  onClick={() => {
                    setResetZoom(!resetZoom);
                    setZoomed(false);
                  }}
                  title="Reset zoom"
                >
                  <Icon icon="search-minus" />
                </Button>
              )}
              {document.fullscreenEnabled && (
                <Button
                  variant="link"
                  onClick={toggleFullscreen}
                  title="Toggle Fullscreen"
                >
                  <Icon icon="expand" />
                </Button>
              )}
              <Button
                variant="link"
                onClick={() => close()}
                title="Close Lightbox"
              >
                <Icon icon="times" />
              </Button>
            </div>
          </div>
          <div className={CLASSNAME_DISPLAY}>
            {images.length > 1 && (
              <Button
                variant="link"
                onClick={handleLeft}
                className={`${CLASSNAME_NAVBUTTON} d-none d-lg-block`}
              >
                <Icon icon="chevron-left" />
              </Button>
            )}

            <div
              className={cx(CLASSNAME_CAROUSEL, {
                [CLASSNAME_INSTANT]: instantTransition,
              })}
              style={{ left: `${currentIndex * -100}vw` }}
              ref={carouselRef}
            >
              {images.map((image, i) => (
                <div className={`${CLASSNAME_IMAGE}`} key={image.paths.image}>
                  {i >= currentIndex - 1 && i <= currentIndex + 1 ? (
                    <LightboxImage
                      src={image.paths.image ?? ""}
                      mode={displayMode}
                      onLeft={() => handleLeft(true)}
                      onRight={handleRight}
                      onZoomed={() => setZoomed(true)}
                      resetZoom={resetZoom}
                    />
                  ) : undefined}
                </div>
              ))}
            </div>

            {images.length > 1 && (
              <Button
                variant="link"
                onClick={handleRight}
                className={`${CLASSNAME_NAVBUTTON} d-none d-lg-block`}
              >
                <Icon icon="chevron-right" />
              </Button>
            )}
          </div>
          {showNavigation && !isFullscreen && images.length > 1 && (
            <div className={CLASSNAME_NAV} ref={navRef}>
              <Button
                variant="link"
                onClick={() => setIndex(images.length - 1)}
                className={CLASSNAME_NAVBUTTON}
              >
                <Icon icon="arrow-left" className="mr-4" />
              </Button>
              {navItems}
              <Button
                variant="link"
                onClick={() => setIndex(0)}
                className={CLASSNAME_NAVBUTTON}
              >
                <Icon icon="arrow-right" className="ml-4" />
              </Button>
            </div>
          )}
        </>
      ) : (
        <LoadingIndicator />
      )}
    </div>
  ) : (
    <></>
  );

  return element;
};
