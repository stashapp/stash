import React, { useCallback, useEffect, useRef, useState } from "react";
import {
  Button,
  Col,
  InputGroup,
  Overlay,
  Popover,
  Form,
  Row,
  Dropdown,
} from "react-bootstrap";
import cx from "classnames";
import Mousetrap from "mousetrap";

import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import useInterval from "../Interval";
import usePageVisibility from "../PageVisibility";
import { useToast } from "../Toast";
import { FormattedMessage, useIntl } from "react-intl";
import { LightboxImage } from "./LightboxImage";
import { useConfigurationContext } from "../Config";
import { Link } from "react-router-dom";
import { OCounterButton } from "src/components/Scenes/SceneDetails/OCounterButton";
import {
  mutateImageIncrementO,
  mutateImageDecrementO,
  mutateImageResetO,
  useImageUpdate,
} from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { useInterfaceLocalForage } from "../LocalForage";
import { imageLightboxDisplayModeIntlMap } from "src/core/enums";
import { ILightboxImage, IChapter } from "./types";
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
  faBars,
  faImages,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { DeleteImagesDialog } from "src/components/Images/DeleteImagesDialog";
import { useDebounce } from "../debounce";
import { isVideo } from "src/utils/visualFile";
import { imageTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import type { LightboxHideReason } from "./context";

const CLASSNAME = "Lightbox";
const CLASSNAME_HEADER = `${CLASSNAME}-header`;
const CLASSNAME_LEFT_SPACER = `${CLASSNAME_HEADER}-left-spacer`;
const CLASSNAME_CHAPTERS = `${CLASSNAME_HEADER}-chapters`;
const CLASSNAME_CHAPTER_BUTTON = `${CLASSNAME_HEADER}-chapter-button`;
const CLASSNAME_INDICATOR = `${CLASSNAME_HEADER}-indicator`;
const CLASSNAME_OPTIONS = `${CLASSNAME_HEADER}-options`;
const CLASSNAME_OPTIONS_ICON = `${CLASSNAME_OPTIONS}-icon`;
const CLASSNAME_OPTIONS_INLINE = `${CLASSNAME_OPTIONS}-inline`;
const CLASSNAME_RIGHT = `${CLASSNAME_HEADER}-right`;
const CLASSNAME_FOOTER = `${CLASSNAME}-footer`;
const CLASSNAME_FOOTER_LEFT = `${CLASSNAME_FOOTER}-left`;
const CLASSNAME_FOOTER_CENTER = `${CLASSNAME_FOOTER}-center`;
const CLASSNAME_FOOTER_RIGHT = `${CLASSNAME_FOOTER}-right`;
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
const MIN_ZOOM = 0.1;
const SCROLL_ZOOM_TIMEOUT = 250;
const ZOOM_NONE_EPSILON = 0.015;

interface IProps {
  images: ILightboxImage[];
  isVisible: boolean;
  isLoading: boolean;
  initialIndex?: number;
  showNavigation: boolean;
  slideshowEnabled?: boolean;
  slideshowAutostart?: boolean;
  page?: number;
  pages?: number;
  pageSize?: number;
  totalCount?: number;
  pageCallback?: (props: { direction?: number; page?: number }) => void;
  chapters?: IChapter[];
  hide: (reason?: LightboxHideReason) => void;
  onDeleteImage?: (id: string) => void;
}

export const LightboxComponent: React.FC<IProps> = ({
  images,
  isVisible,
  isLoading,
  initialIndex = 0,
  showNavigation,
  slideshowEnabled = false,
  slideshowAutostart = false,
  page,
  pageSize = 40,
  totalCount,
  pageCallback,
  chapters = [],
  hide,
  onDeleteImage,
}) => {
  const [updateImage] = useImageUpdate();

  // zero-based
  const [index, setIndex] = useState<number | null>(null);
  const [movingLeft, setMovingLeft] = useState(false);
  const oldIndex = useRef<number | null>(null);
  const [instantTransition, setInstantTransition] = useState(false);
  const [isSwitchingPage, setIsSwitchingPage] = useState(true);
  // Synchronous mirror of isSwitchingPage. The nav handlers (handleLeft/
  // handleRight, reached by the arrow keys, the on-screen chevrons and the
  // image-edge clicks) fire on raw keydown/click events that can arrive faster
  // than React re-renders, so the isSwitchingPage *state* they close over is
  // stale. The ref reflects an in-flight page switch immediately — set true
  // before pageCallback — so rapid inputs are dropped and, crucially, the
  // index-range effect won't clamp a stale index while a (possibly
  // synchronously-cached) page is still swapping in. The landing index stays
  // under the handlers' control; this ref only gates, it never picks the index.
  const isSwitchingPageRef = useRef(true);
  const [isFullscreen, setFullscreen] = useState(false);
  const [showOptions, setShowOptions] = useState(false);
  const [showChapters, setShowChapters] = useState(false);
  const [imagesLoaded, setImagesLoaded] = useState(0);
  // The image the delete dialog acts on, captured when the dialog is opened.
  // Reading images[currentIndex] at confirm time instead would delete whatever
  // image the index has since moved to (e.g. a page-switch settle landing
  // while the dialog is open).
  const [deleteTarget, setDeleteTarget] = useState<ILightboxImage | null>(null);
  const lastDKeyTime = useRef<number>(0);
  const [navOffset, setNavOffset] = useState<React.CSSProperties | undefined>();

  // An in-flight page switch's intended landing, set synchronously by
  // startPageSwitch *before* pageCallback. The landing is an explicit target —
  // "last" or a concrete index (0 for a forward cross) — NOT derived from the
  // page-number direction (a wrap moves the page number opposite to nav
  // direction, so deriving from direction breaks first/last wraparound). It is
  // resolved by the settle effect once the new page has loaded.
  const pendingTarget = useRef<"last" | number | null>(null);
  // The page number and images array at switch start, used by the settle
  // effect to detect the new page's arrival.
  const switchFromPage = useRef<number | undefined>(page);
  const switchFromImages = useRef<ILightboxImage[]>(images);

  // The ref and the state must never disagree, so they are only written
  // through this setter.
  const setSwitching = useCallback((v: boolean) => {
    isSwitchingPageRef.current = v;
    setIsSwitchingPage(v);
  }, []);

  // Begin a page switch: raise the guard and record the landing intent
  // *before* pageCallback. An already-cached page can swap its images in
  // synchronously, so the guard must be up first or a rapid input (and the
  // index-range effect) races the swap. The index is parked at 0 — valid on
  // any non-empty page — while the new page loads, so the index-range clamp
  // can't fire on a shorter incoming page; the settle effect moves it to the
  // target once the page arrives.
  const startPageSwitch = useCallback(
    (target: "last" | number, cbArg: { direction?: number; page?: number }) => {
      setSwitching(true);
      switchFromPage.current = page;
      switchFromImages.current = images;
      pendingTarget.current = target;
      setIndex(0);
      pageCallback?.(cbArg);
    },
    [images, page, pageCallback, setSwitching]
  );

  const [zoom, setZoom] = useState(1);

  function updateZoom(v: number) {
    if (v < MIN_ZOOM) {
      setZoom(MIN_ZOOM);
    } else if (Math.abs(v - 1) < ZOOM_NONE_EPSILON) {
      // "snap to 1" effect: if new zoom is close to 1, set to 1
      setZoom(1);
    } else {
      setZoom(v);
    }
  }

  const [resetPosition, setResetPosition] = useState(false);

  const containerRef = useRef<HTMLDivElement | null>(null);
  const overlayTarget = useRef<HTMLButtonElement | null>(null);
  const carouselRef = useRef<HTMLDivElement | null>(null);
  const indicatorRef = useRef<HTMLDivElement | null>(null);
  const navRef = useRef<HTMLDivElement | null>(null);
  const clearIntervalCallback = useRef<() => void>();
  const resetIntervalCallback = useRef<() => void>();

  const allowNavigation = images.length > 1 || pageCallback;

  const Toast = useToast();
  const intl = useIntl();
  const { configuration: config } = useConfigurationContext();
  const [interfaceLocalForage, setInterfaceLocalForage] =
    useInterfaceLocalForage();

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

  const disableAnimation = config?.interface.imageLightbox.disableAnimation;

  function setSlideshowDelay(v: number) {
    setLightboxSettings({ slideshowDelay: v });
  }

  const scaleUp =
    lightboxSettings?.scaleUp ??
    config?.interface.imageLightbox.scaleUp ??
    false;

  const resetZoomOnNav =
    lightboxSettings?.resetZoomOnNav ??
    config?.interface.imageLightbox.resetZoomOnNav ??
    false;

  const scrollMode =
    lightboxSettings?.scrollMode ??
    config?.interface.imageLightbox.scrollMode ??
    GQL.ImageLightboxScrollMode.Zoom;

  const displayMode =
    lightboxSettings?.displayMode ??
    config?.interface.imageLightbox.displayMode ??
    GQL.ImageLightboxDisplayMode.FitXy;
  const oldDisplayMode = useRef(displayMode);

  function setDisplayMode(v: GQL.ImageLightboxDisplayMode) {
    setLightboxSettings({ displayMode: v });
  }

  // `slideshowInterval` controls the slideshow logic, while
  // `displayedSlideshowInterval` is for display purposes only. Keeping the two
  // separate lets us handle the logic however we want while still showing the
  // user something that makes sense.
  //
  // Autostart the slideshow on open when requested (e.g. opening a gallery's
  // lightbox from the galleries page). The component mounts fresh on each open,
  // so initialising the interval here starts playback immediately.
  const [slideshowInterval, setSlideshowInterval] = useState<number | null>(
    slideshowEnabled && slideshowAutostart ? slideshowDelay : null
  );

  const [displayedSlideshowInterval, setDisplayedSlideshowInterval] =
    useState<string>((slideshowDelay / SECONDS_TO_MS).toString());

  useEffect(() => {
    const target = pendingTarget.current;
    if (target !== null) {
      // A page switch is in flight. It completes when the page number AND the
      // images array have both changed: the parent delivers them together, and
      // requiring both means a parent that updates the page a commit ahead of
      // the images can't settle the switch against the outgoing page's array.
      // Callers that provide pageCallback without a page prop can never change
      // the page number, so for them fall back to the images identity alone.
      const pageChanged = page !== switchFromPage.current;
      const imagesChanged = images !== switchFromImages.current;
      const settled =
        page === undefined ? imagesChanged : pageChanged && imagesChanged;
      if (settled) {
        // Resolve the landing, clamped to the arrived page: a chapter's
        // numeric target can be stale (images deleted since the chapters
        // loaded), and "last" is -1 on an empty page (which the close effect
        // then handles).
        const resolved =
          target === "last"
            ? images.length - 1
            : Math.min(target, images.length - 1);
        setIndex(Math.max(0, resolved));
        pendingTarget.current = null;
        switchFromPage.current = page;
        switchFromImages.current = images;
        setSwitching(false);
      }
      return;
    }
    // No switch pending: clear the switching flag once the first page's images
    // have loaded (the flag starts true on open so the loader shows until then).
    if (isSwitchingPage && !isLoading && images.length > 0) {
      setSwitching(false);
    }
  }, [isSwitchingPage, isLoading, images, page, setSwitching]);

  const disableInstantTransition = useDebounce(
    () => setInstantTransition(false),
    400
  );

  const setInstant = useCallback(() => {
    setInstantTransition(true);
    disableInstantTransition();
  }, [disableInstantTransition]);

  useEffect(() => {
    if (images.length < 2) return;
    if (index === oldIndex.current) return;
    if (index === null) return;

    // reset zoom status
    // setResetZoom((r) => !r);
    // setZoomed(false);
    if (resetZoomOnNav) {
      setZoom(1);
    }
    setResetPosition((r) => !r);

    oldIndex.current = index;
  }, [index, images.length, resetZoomOnNav]);

  const getNavOffset = useCallback(() => {
    if (images.length < 2) return;
    if (index === undefined || index === null) return;

    if (navRef.current) {
      const currentThumb = navRef.current.children[index + 1];
      if (currentThumb instanceof HTMLImageElement) {
        const offset =
          -1 *
          (currentThumb.offsetLeft - document.documentElement.clientWidth / 2);

        return { left: `${offset}px` };
      }
    }
  }, [index, images.length]);

  // biome-ignore lint/correctness/useExhaustiveDependencies: intentional reset on images change
  useEffect(() => {
    // reset images loaded counter for new images
    setImagesLoaded(0);
  }, [images]);

  useEffect(() => {
    setNavOffset(getNavOffset() ?? undefined);
  }, [getNavOffset]);

  useEffect(() => {
    if (displayMode !== oldDisplayMode.current) {
      // reset zoom status
      // setResetZoom((r) => !r);
      // setZoomed(false);
      if (resetZoomOnNav) {
        setZoom(1);
      }
      setResetPosition((r) => !r);
    }
    oldDisplayMode.current = displayMode;
  }, [displayMode, resetZoomOnNav]);

  const selectIndex = (e: React.MouseEvent, i: number) => {
    setIndex(i);
    e.stopPropagation();
  };

  useEffect(() => {
    if (isVisible && index === null) setIndex(initialIndex);
  }, [initialIndex, isVisible, index]);

  useEffect(() => {
    if (!isVisible) return;

    document.body.style.overflow = "hidden";
    Mousetrap.pause();
    return () => {
      const fullscreenElement = document.fullscreenElement;
      if (
        fullscreenElement &&
        containerRef.current?.contains(fullscreenElement)
      ) {
        document.exitFullscreen();
      }
      document.body.style.overflow = "auto";
      Mousetrap.unpause();
    };
  }, [isVisible]);

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

  const close = useCallback(
    (reason: LightboxHideReason = "dismiss") => {
      hide(reason);
    },
    [hide]
  );

  const handleClose = (e: React.MouseEvent<HTMLDivElement>) => {
    const { className } = e.target as Element;
    if (className?.includes?.(CLASSNAME_IMAGE)) close();
  };

  const handleLeft = useCallback(
    (isUserAction = true) => {
      if (isSwitchingPageRef.current) return;

      if (disableAnimation) {
        setInstant();
      }

      setShowChapters(false);
      setMovingLeft(true);

      if (index === 0) {
        // go to previous page (landing on its last image), or loop back if no
        // callback is set
        if (pageCallback) {
          startPageSwitch("last", { direction: -1 });
        } else setIndex(images.length - 1);
      } else setIndex((index ?? 0) - 1);

      if (isUserAction && resetIntervalCallback.current) {
        resetIntervalCallback.current();
      }
    },
    [images, pageCallback, index, disableAnimation, setInstant, startPageSwitch]
  );

  const handleRight = useCallback(
    (isUserAction = true) => {
      if (isSwitchingPageRef.current) return;

      if (disableAnimation) {
        setInstant();
      }

      setMovingLeft(false);
      setShowChapters(false);

      if (index === images.length - 1) {
        // go to next page (landing on its first image), or loop back if no
        // callback is set
        if (pageCallback) {
          startPageSwitch(0, { direction: 1 });
        } else setIndex(0);
      } else setIndex((index ?? 0) + 1);

      if (isUserAction && resetIntervalCallback.current) {
        resetIntervalCallback.current();
      }
    },
    [images, pageCallback, index, disableAnimation, setInstant, startPageSwitch]
  );

  const firstScroll = useRef<number | null>(null);
  const inScrollGroup = useRef(false);

  const debouncedScrollReset = useDebounce(() => {
    firstScroll.current = null;
    inScrollGroup.current = false;
  }, SCROLL_ZOOM_TIMEOUT);

  useEffect(() => {
    return () => {
      disableInstantTransition.cancel();
      debouncedScrollReset.cancel();
    };
  }, [disableInstantTransition, debouncedScrollReset]);

  const handleKey = useCallback(
    (e: KeyboardEvent) => {
      if (e.repeat && (e.key === "ArrowRight" || e.key === "ArrowLeft"))
        setInstant();
      if (e.key === "ArrowLeft") handleLeft();
      else if (e.key === "ArrowRight") handleRight();
      else if (e.key === "Escape") close();
      else if (e.key === "d") {
        // Not while a page switch is in flight: the index is parked at 0 then,
        // so the shortcut would target an image the user isn't viewing.
        const image = images[index ?? initialIndex];
        if (!isSwitchingPageRef.current && image?.id !== undefined) {
          const now = Date.now();
          if (now - lastDKeyTime.current < 1000) {
            setDeleteTarget(image);
          }
          lastDKeyTime.current = now;
        }
      }
    },
    [setInstant, handleLeft, handleRight, close, images, index, initialIndex]
  );

  const [clearCallback, resetCallback] = useInterval(
    () => {
      handleRight(false);
    },
    slideshowEnabled ? slideshowInterval : null
  );

  resetIntervalCallback.current = resetCallback;
  clearIntervalCallback.current = clearCallback;

  useEffect(() => {
    const handleFullScreenChange = () => {
      if (clearIntervalCallback.current) {
        clearIntervalCallback.current();
      }
      setFullscreen(document.fullscreenElement !== null);
    };

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

  function imageLoaded() {
    setImagesLoaded((loaded) => loaded + 1);

    if (imagesLoaded === images.length - 1) {
      // all images are loaded - update the nav offset
      setNavOffset(getNavOffset() ?? undefined);
    }
  }

  const navItems = images.map((image, i) =>
    React.createElement(image.paths.preview !== "" ? "video" : "img", {
      loop: image.paths.preview !== "",
      autoPlay: image.paths.preview !== "",
      playsInline: image.paths.preview !== "",
      src:
        image.paths.preview !== ""
          ? (image.paths.preview ?? "")
          : (image.paths.thumbnail ?? ""),
      alt: "",
      className: cx(CLASSNAME_NAVIMAGE, {
        [CLASSNAME_NAVSELECTED]: i === index,
      }),
      onClick: (e: React.MouseEvent) => selectIndex(e, i),
      role: "presentation",
      loading: "lazy",
      key: image.paths.thumbnail,
      onLoad: imageLoaded,
    })
  );

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

  const currentIndex = index === null ? initialIndex : index;

  useEffect(() => {
    // Don't auto-close while images are still loading. Some entry points open
    // the lightbox with an empty image list and isLoading=true, then populate
    // it asynchronously. Only an empty list *after* loading means the last
    // image was deleted.
    if (isLoading) return;
    if (images.length === 0) {
      close();
    } else if (
      !isSwitchingPageRef.current &&
      index !== null &&
      index >= images.length
    ) {
      // Same-page shrink only — e.g. the last image was deleted and findImages
      // refetched a shorter list. While a page switch is in flight the index is
      // transiently stale against a just-swapped, shorter page; clamping it
      // here is what made a forward cross land on the page's *last* image. The
      // settle effect reconciles the switch instead, so skip the clamp then.
      setIndex(images.length - 1);
    }
  }, [images.length, index, close, isLoading]);

  function gotoPage(imageIndex: number) {
    // indexInPage is the chapter's explicit target; the settle effect lands on
    // it once the new page loads, so a chapter jump lands on the right image of
    // the new page rather than its first/last.
    const indexInPage = (imageIndex - 1) % pageSize;
    if (pageCallback) {
      const jumppage = Math.floor((imageIndex - 1) / pageSize) + 1;
      if (page !== jumppage) {
        startPageSwitch(indexInPage, { page: jumppage });
        setShowChapters(false);
        return;
      }
    }

    setIndex(indexInPage);
    setShowChapters(false);
  }

  function chapterHeader() {
    const imageNumber = (index ?? 0) + 1;
    const globalIndex = page
      ? (page - 1) * pageSize + imageNumber
      : imageNumber;

    let chapterTitle = "";
    chapters.forEach((chapter) => {
      if (chapter.image_index > globalIndex) {
        return;
      }
      chapterTitle = chapter.title;
    });

    return chapterTitle ?? "";
  }

  const renderChapterMenu = () => {
    if (chapters.length <= 0) return;

    const popoverContent = chapters.map(({ id, title, image_index }) => (
      <Dropdown.Item key={id} onClick={() => gotoPage(image_index)}>
        {" "}
        {title}
        {title.length > 0 ? " - #" : "#"}
        {image_index}
      </Dropdown.Item>
    ));

    return (
      <Dropdown
        show={showChapters}
        onToggle={() => setShowChapters(!showChapters)}
      >
        <Dropdown.Toggle className={`minimal ${CLASSNAME_CHAPTER_BUTTON}`}>
          <Icon icon={showChapters ? faTimes : faBars} />
        </Dropdown.Toggle>
        <Dropdown.Menu className={`${CLASSNAME_CHAPTERS}`}>
          {popoverContent}
        </Dropdown.Menu>
      </Dropdown>
    );
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
                checked={scaleUp}
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
                checked={resetZoomOnNav}
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
                value={scrollMode}
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

  function renderBody() {
    if (images.length === 0 || isLoading || isSwitchingPage) {
      return <LoadingIndicator />;
    }

    const currentImage: ILightboxImage | undefined = images[currentIndex];
    const title = currentImage ? imageTitle(currentImage) : undefined;

    function setRating(v: number | null) {
      if (currentImage?.id) {
        updateImage({
          variables: {
            input: {
              id: currentImage.id,
              rating100: v,
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

    const currentNumber =
      page !== undefined
        ? (page - 1) * pageSize + currentIndex + 1
        : currentIndex + 1;
    const displayTotal = totalCount ?? images.length;

    return (
      <>
        <div className={CLASSNAME_HEADER}>
          <div className={CLASSNAME_LEFT_SPACER}>{renderChapterMenu()}</div>
          <div className={CLASSNAME_INDICATOR}>
            <span>{chapterHeader()}</span>
            {displayTotal > 1 ? (
              <b ref={indicatorRef}>{`${currentNumber} / ${displayTotal}`}</b>
            ) : undefined}
          </div>
          <div className={CLASSNAME_RIGHT}>
            <div className={CLASSNAME_OPTIONS}>
              <div className={CLASSNAME_OPTIONS_ICON}>
                <Button
                  ref={overlayTarget}
                  variant="link"
                  title={intl.formatMessage({
                    id: "dialogs.lightbox.options",
                  })}
                  onClick={() => setShowOptions(!showOptions)}
                >
                  <Icon icon={faCog} />
                </Button>
                <Overlay
                  target={overlayTarget.current}
                  show={showOptions}
                  placement="bottom"
                  container={containerRef}
                  rootClose
                  onHide={() => setShowOptions(false)}
                >
                  {({ placement, arrowProps, show: _show, ...props }) => (
                    <div
                      className="popover"
                      {...props}
                      style={{ ...props.style }}
                    >
                      <Popover.Title>
                        {intl.formatMessage({
                          id: "dialogs.lightbox.options",
                        })}
                      </Popover.Title>
                      <Popover.Content>{renderOptionsForm()}</Popover.Content>
                    </div>
                  )}
                </Overlay>
              </div>
              <InputGroup className={CLASSNAME_OPTIONS_INLINE}>
                {renderOptionsForm()}
              </InputGroup>
            </div>
            {slideshowEnabled && (
              <Button
                variant="link"
                onClick={toggleSlideshow}
                title="Toggle Slideshow"
              >
                <Icon icon={slideshowInterval !== null ? faPause : faPlay} />
              </Button>
            )}
            {zoom !== 1 && (
              <Button
                variant="link"
                onClick={() => {
                  setResetPosition(!resetPosition);
                  setZoom(1);
                }}
                title="Reset zoom"
              >
                <Icon icon={faSearchMinus} />
              </Button>
            )}
            {document.fullscreenEnabled && (
              <Button
                variant="link"
                onClick={toggleFullscreen}
                title="Toggle Fullscreen"
              >
                <Icon icon={faExpand} />
              </Button>
            )}
            <Button
              variant="link"
              onClick={() => close()}
              title="Close Lightbox"
            >
              <Icon icon={faTimes} />
            </Button>
          </div>
        </div>
        <div className={CLASSNAME_DISPLAY}>
          {allowNavigation && (
            <Button
              variant="link"
              onClick={handleLeft}
              className={`${CLASSNAME_NAVBUTTON} d-none d-lg-block`}
            >
              <Icon icon={faChevronLeft} />
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
                    width={image.visual_files?.[0]?.width ?? 0}
                    height={image.visual_files?.[0]?.height ?? 0}
                    displayMode={displayMode}
                    scaleUp={scaleUp}
                    scrollMode={scrollMode}
                    resetPosition={resetPosition}
                    zoom={i === currentIndex ? zoom : 1}
                    scrollAttemptsBeforeChange={scrollAttemptsBeforeChange}
                    firstScroll={firstScroll}
                    inScrollGroup={inScrollGroup}
                    current={i === currentIndex}
                    alignBottom={movingLeft}
                    setZoom={updateZoom}
                    debouncedScrollReset={debouncedScrollReset}
                    onLeft={handleLeft}
                    onRight={handleRight}
                    isVideo={isVideo(image.visual_files?.[0] ?? {})}
                  />
                ) : undefined}
              </div>
            ))}
          </div>

          {allowNavigation && (
            <Button
              variant="link"
              onClick={handleRight}
              className={`${CLASSNAME_NAVBUTTON} d-none d-lg-block`}
            >
              <Icon icon={faChevronRight} />
            </Button>
          )}
        </div>
        {showNavigation && !isFullscreen && images.length > 1 && (
          <div className={CLASSNAME_NAV} style={navOffset} ref={navRef}>
            <Button
              variant="link"
              onClick={() => setIndex(images.length - 1)}
              className={CLASSNAME_NAVBUTTON}
            >
              <Icon icon={faArrowLeft} className="mr-4" />
            </Button>
            {navItems}
            <Button
              variant="link"
              onClick={() => setIndex(0)}
              className={CLASSNAME_NAVBUTTON}
            >
              <Icon icon={faArrowRight} className="ml-4" />
            </Button>
          </div>
        )}
        <div className={CLASSNAME_FOOTER}>
          <div className={CLASSNAME_FOOTER_LEFT}>
            {currentImage?.id !== undefined && (
              <>
                <div>
                  <OCounterButton
                    onDecrement={onDecrementClick}
                    onIncrement={onIncrementClick}
                    onReset={onResetClick}
                    value={currentImage?.o_counter ?? 0}
                  />
                </div>
                <RatingSystem
                  value={currentImage?.rating100}
                  onSetRating={(v) => setRating(v)}
                  clickToRate
                  withoutContext
                />
              </>
            )}
          </div>
          <div className={CLASSNAME_FOOTER_CENTER}>
            {currentImage && (
              <>
                <div
                  className="d-flex align-items-center justify-content-center"
                  style={{ gap: "0.5rem" }}
                >
                  <Link
                    className="image-link"
                    to={`/images/${currentImage.id}`}
                    replace
                    onClick={() => close("navigate")}
                  >
                    {title ?? ""}
                  </Link>
                  {currentImage.id !== undefined && (
                    <Button
                      className="minimal delete-button"
                      onClick={() => setDeleteTarget(currentImage)}
                      title={intl.formatMessage({ id: "actions.delete" })}
                    >
                      <Icon icon={faTrash} />
                    </Button>
                  )}
                </div>
                {currentImage.galleries?.length ? (
                  <Link
                    className="image-gallery-link"
                    to={`/galleries/${currentImage.galleries[0].id}`}
                    replace
                    onClick={() => close("navigate")}
                  >
                    <Icon icon={faImages} />
                    {galleryTitle(currentImage.galleries[0])}
                  </Link>
                ) : null}
              </>
            )}
          </div>
          <div className={CLASSNAME_FOOTER_RIGHT} />
        </div>
      </>
    );
  }

  if (!isVisible) {
    return null;
  }

  return (
    <div
      className={CLASSNAME}
      role="presentation"
      ref={containerRef}
      onClick={handleClose}
    >
      {renderBody()}
      {deleteTarget?.id !== undefined && (
        <DeleteImagesDialog
          selected={[
            {
              id: deleteTarget.id,
              visual_files: deleteTarget.visual_files ?? [],
            },
          ]}
          onClose={(confirmed) => {
            setDeleteTarget(null);
            if (confirmed) onDeleteImage?.(deleteTarget.id!);
          }}
        />
      )}
    </div>
  );
};

export default LightboxComponent;
