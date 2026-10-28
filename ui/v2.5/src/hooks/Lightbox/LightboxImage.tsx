import React, {
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
  useCallback,
} from "react";
import { flushSync } from "react-dom";
import useResizeObserver from "@react-hook/resize-observer";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useDebounce } from "src/hooks/debounce";

const ZOOM_STEP = 1.1;
const ZOOM_FACTOR = 700;
const SCROLL_GROUP_THRESHOLD = 8;
const SCROLL_GROUP_EXIT_THRESHOLD = 4;
const SCROLL_INFINITE_THRESHOLD = 10;
const SCROLL_PAN_STEP = 75;
const SCROLL_PAN_FACTOR = 2;
// Zoom clamp. The upper cap is image-size-aware (see maxZoomForPixels): only a
// *large* image's composited layer can exhaust GPU memory at high zoom and crash
// the tab on phones, so large images get the smaller device-tied cap while small
// images (cheap to raster at any scale) reach the full ceiling on every device.
// (PINCH_MIN must match Lightbox.tsx's updateZoom MIN_ZOOM, and both clamp to the
// same maxZoomForPixels, so the pinch focal-pan uses the same effective scale
// change the zoom is actually clamped to, or the pan diverges off-screen.)
const PINCH_MIN_ZOOM = 0.1;

function computeDeviceMaxZoom(): number {
  // Chrome/Android expose approximate RAM; scale the ceiling with it.
  const mem = (navigator as Navigator & { deviceMemory?: number }).deviceMemory;
  if (typeof mem === "number" && mem > 0) {
    if (mem >= 8) return 10;
    if (mem >= 4) return 6;
    return 4;
  }
  // Safari/iOS has no deviceMemory. Phones have far less GPU headroom than
  // tablets, so use the smaller screen side (CSS px) as a phone-vs-tablet proxy:
  // phones are < ~430, the smallest iPads ~744.
  const minSide =
    typeof window !== "undefined" && window.screen
      ? Math.min(window.screen.width, window.screen.height)
      : 0;
  return minSide > 0 && minSide < 600 ? 5 : 10;
}

export const DEVICE_MAX_ZOOM = computeDeviceMaxZoom();

// Touch gesture tuning (#2538). A single-finger swipe commits once it travels
// far enough OR fast enough; the axis (horizontal/vertical/pan) is locked once
// movement passes the lock threshold so a gesture can't change its mind midway.
const SWIPE_COMMIT_DISTANCE = 60; // px of travel that commits a nav/delete swipe
const SWIPE_VELOCITY = 0.3; // px/ms; a quick flick commits below the distance
const AXIS_LOCK_THRESHOLD = 10; // px of travel before the axis is decided
const TAP_MAX_MOVE = 10; // px; movement below this is a tap, not a swipe
const DOUBLE_TAP_MS = 300; // max gap between taps to count as a double-tap
const DOUBLE_TAP_SLOP = 30; // px; max distance between the two taps
const PANNABLE_EPSILON = 1; // px tolerance when testing if the image overflows
// Above this natural resolution, the first zoom-in forces an expensive full-res
// raster (a multi-second stall on huge images). On a double-tap (which jumps
// straight to the target zoom, with no in-between feedback) we cover that stall
// with a busy spinner; a pinch needs no spinner because the live finger-tracking
// is its own feedback. Cheaper images never trigger it either way. This is also
// the threshold for the conservative zoom cap below - only images this big risk
// an OOM.
const LARGE_IMAGE_PIXELS = 30_000_000; // ~30 MP
const ZOOM_SPINNER_MIN_MS = 400; // min spinner time when we decode-to-ready

// Hard ceiling on zoom for any image on any device.
export const MAX_ZOOM = 10;

// Per-image upper zoom cap. A large image (> LARGE_IMAGE_PIXELS) gets the
// smaller device-tied cap because its composited layer can exhaust GPU memory at
// high zoom; a small image is cheap to raster at any scale, so it reaches the
// full ceiling everywhere. `pixels` is the image's natural width * height.
export function maxZoomForPixels(pixels: number): number {
  return pixels > LARGE_IMAGE_PIXELS ? DEVICE_MAX_ZOOM : MAX_ZOOM;
}
const CLASSNAME = "Lightbox";
const CLASSNAME_CAROUSEL = `${CLASSNAME}-carousel`;
const CLASSNAME_IMAGE = `${CLASSNAME_CAROUSEL}-image`;

function calculateDefaultZoom(
  width: number,
  height: number,
  boundWidth: number,
  boundHeight: number,
  displayMode: GQL.ImageLightboxDisplayMode,
  scaleUp: boolean
) {
  // set initial zoom level based on options
  let xZoom: number;
  let yZoom: number;
  let newZoom = 1;
  switch (displayMode) {
    case GQL.ImageLightboxDisplayMode.FitXy:
      xZoom = boundWidth / width;
      yZoom = boundHeight / height;

      if (!scaleUp) {
        xZoom = Math.min(xZoom, 1);
        yZoom = Math.min(yZoom, 1);
      }
      newZoom = Math.min(xZoom, yZoom);
      break;
    case GQL.ImageLightboxDisplayMode.FitX:
      newZoom = boundWidth / width;

      if (!scaleUp) {
        newZoom = Math.min(newZoom, 1);
      }
      break;
    case GQL.ImageLightboxDisplayMode.Original:
      newZoom = 1;
      break;
  }

  return newZoom;
}

interface IDimension {
  width: number;
  height: number;
}

// Track the container's size live via ResizeObserver. The box dimensions feed
// fit-zoom, centering, the pannable/axis decisions and the pan-edge bounds, so
// they must follow window resizes and phone orientation changes - measuring once
// at mount leaves a zoomed image mis-centered after a rotate. Returns a stable
// ref to attach to the element (also fixes the prior createRef()-per-render).
function useContainerDimensions<T extends HTMLElement = HTMLDivElement>(): [
  React.MutableRefObject<T | null>,
  IDimension,
] {
  const ref = useRef<T | null>(null);
  const [dimension, setDimension] = useState<IDimension>({
    width: 0,
    height: 0,
  });

  // Measure synchronously on mount (before paint) so the fit-zoom is known for
  // the first render. Otherwise the box stays 0 until the debounced observer
  // fires, and a large image renders at full natural size for ~120ms first
  // (an expensive raster) before fit-scaling down.
  useLayoutEffect(() => {
    const el = ref.current;
    if (!el) return;
    const r = el.getBoundingClientRect();
    if (r.width && r.height) setDimension({ width: r.width, height: r.height });
  }, []);

  // Debounced so an *animated* rotate - which fires a burst of intermediate
  // sizes - re-lays-out only once it settles. Used for subsequent resizes /
  // rotation only; the initial size comes from the layout effect above.
  const onResize = useDebounce((entry: ResizeObserverEntry) => {
    const { inlineSize: width, blockSize: height } = entry.contentBoxSize[0];
    setDimension({ width, height });
  }, 120);

  useResizeObserver(ref, onResize);

  return [ref, dimension];
}

interface IProps {
  src: string;
  width: number;
  height: number;
  displayMode: GQL.ImageLightboxDisplayMode;
  scaleUp: boolean;
  scrollMode: GQL.ImageLightboxScrollMode;
  resetPosition?: boolean;
  zoom: number;
  scrollAttemptsBeforeChange: number;
  // these refs must be outside of LightboxImage,
  // since they need to be shared between all LightboxImages
  firstScroll: React.MutableRefObject<number | null>;
  inScrollGroup: React.MutableRefObject<boolean>;
  current: boolean;
  // set to true to align image with bottom instead of top
  alignBottom?: boolean;
  setZoom: (v: number) => void;
  debouncedScrollReset: () => void;
  onLeft: () => void;
  onRight: () => void;
  // swipe-up-to-delete; opens the confirmation dialog (never bypassed)
  onSwipeDelete?: () => void;
  // swipe-down-to-close; dismisses the lightbox
  onSwipeClose?: () => void;
  // ask the parent to re-center the image (toggles its resetPosition)
  onResetPosition?: () => void;
  // whether horizontal swipes should navigate (false for single-image galleries)
  navigationEnabled: boolean;
  isVideo: boolean;
}

export const LightboxImage: React.FC<IProps> = ({
  src,
  width,
  height,
  displayMode,
  scaleUp,
  scrollMode,
  resetPosition,
  zoom,
  scrollAttemptsBeforeChange,
  firstScroll,
  inScrollGroup,
  current,
  alignBottom,
  setZoom,
  debouncedScrollReset,
  onLeft,
  onRight,
  onSwipeDelete,
  onSwipeClose,
  onResetPosition,
  navigationEnabled,
  isVideo,
}) => {
  const [defaultZoom, setDefaultZoom] = useState(1);
  const [moving, setMoving] = useState(false);
  const [positionX, setPositionX] = useState(0);
  const [positionY, setPositionY] = useState(0);
  const [imageWidth, setImageWidth] = useState(width);
  const [imageHeight, setImageHeight] = useState(height);
  // live container size (follows window resize / orientation change)
  const [container, { width: boxWidth, height: boxHeight }] =
    useContainerDimensions<HTMLDivElement>();
  const dimensionsProvided = width > 0 && height > 0;

  const mouseDownEvent = useRef<MouseEvent>();
  const resetPositionRef = useRef(resetPosition);

  const startPoints = useRef<number[]>([0, 0]);
  const pointerCache = useRef<React.PointerEvent[]>([]);
  // Snapshot taken once at the start of a pinch. Pinch zoom/pan is computed
  // absolutely from this each frame (no per-frame accumulation or DOM reads),
  // which is what keeps it stable - see onPointerMove.
  const pinch = useRef<{
    dist: number;
    zoom: number;
    posX: number;
    posY: number;
    c0x: number;
    c0y: number;
    fx: number;
    fy: number;
  } | null>(null);

  // single-finger touch gesture state (#2538)
  const touchStart = useRef<{ x: number; y: number; t: number } | null>(null);
  const gestureAxis = useRef<"none" | "horizontal" | "vertical" | "pan">(
    "none"
  );
  const lastTap = useRef<{ x: number; y: number; t: number } | null>(null);
  // timestamp of the last touchend, used to suppress the "ghost" mouse click
  // iOS synthesizes after a tap (which would otherwise trigger half-screen nav)
  const lastTouchEnd = useRef(0);

  // Busy spinner for the first (expensive) zoom-in of a large image. `zoomedOnce`
  // gates it to once per image (the raster is cached afterwards).
  const [zoomBusy, setZoomBusy] = useState(false);
  // True while a 2-finger pinch is active. The busy spinner is unmounted during
  // a pinch: its always-running rotate animation (needed so it can appear while
  // the main thread is blocked - see runZoomIn) otherwise competes with the
  // pinch's own heavy compositing on a huge image and makes it janky. We never
  // show a spinner for pinch anyway, so dropping it here is free.
  const [pinching, setPinching] = useState(false);
  const zoomedOnce = useRef(false);
  const zoomSpinnerShownAt = useRef(0);
  const isLargeImage = imageWidth * imageHeight > LARGE_IMAGE_PIXELS;
  // Per-image zoom ceiling (large images get the conservative device cap).
  const maxZoom = maxZoomForPixels(imageWidth * imageHeight);

  // Show the busy spinner now (flushSync so it paints before the blocking work).
  function showZoomBusy() {
    zoomSpinnerShownAt.current = performance.now();
    flushSync(() => setZoomBusy(true));
  }

  // Clear it, but not before it's been visible at least `minMs` (so it always
  // registers even when the decode finishes quickly).
  function clearZoomBusy(minMs = ZOOM_SPINNER_MIN_MS) {
    const wait = minMs - (performance.now() - zoomSpinnerShownAt.current);
    if (wait > 0) window.setTimeout(() => setZoomBusy(false), wait);
    else setZoomBusy(false);
  }

  // biome-ignore lint/correctness/useExhaustiveDependencies: explicitly want to reset the zoom state whenever the image src changes
  useEffect(() => {
    zoomedOnce.current = false;
    setZoomBusy(false);
  }, [src]);

  const scrollAttempts = useRef(0);

  useEffect(() => {
    function toggleVideoPlay() {
      if (container.current) {
        const openVideo = container.current.getElementsByTagName("video");
        if (openVideo.length > 0) {
          const rect = openVideo[0].getBoundingClientRect();
          if (Math.abs(rect.x) < document.body.clientWidth / 2) {
            openVideo[0].play();
          } else {
            openVideo[0].pause();
          }
        }
      }
    }

    setTimeout(() => {
      toggleVideoPlay();
    }, 250);
  }, [container]);

  useEffect(() => {
    if (dimensionsProvided) {
      return;
    }
    let mounted = true;
    const img = new Image();
    function onLoad() {
      if (mounted) {
        setImageWidth(img.width);
        setImageHeight(img.height);
      }
    }

    img.onload = onLoad;
    img.src = src;

    return () => {
      mounted = false;
    };
  }, [src, dimensionsProvided]);

  const minMaxY = useCallback(
    (appliedZoom: number) => {
      let minY: number, maxY: number;
      const inBounds = appliedZoom * imageHeight <= boxHeight;

      // NOTE: I don't even know how these work, but they do
      if (!inBounds) {
        if (imageHeight > boxHeight) {
          minY =
            (appliedZoom * imageHeight - imageHeight) / 2 -
            appliedZoom * imageHeight +
            boxHeight;
          maxY = (appliedZoom * imageHeight - imageHeight) / 2;
        } else {
          minY = (boxHeight - appliedZoom * imageHeight) / 2;
          maxY = (appliedZoom * imageHeight - boxHeight) / 2;
        }
      } else {
        minY = Math.min((boxHeight - imageHeight) / 2, 0);
        maxY = minY;
      }

      return [minY, maxY];
    },
    [imageHeight, boxHeight]
  );

  // Horizontal pan bounds - mirror of minMaxY (the transform scales about the
  // wrapper centre and the image lays out flex-start when it overflows, so X is
  // symmetric to Y). [minX, maxX] is the range of positionX that keeps the
  // scaled image covering the viewport, i.e. you can't pan past an edge.
  const minMaxX = useCallback(
    (appliedZoom: number) => {
      let minX: number, maxX: number;
      const inBounds = appliedZoom * imageWidth <= boxWidth;

      if (!inBounds) {
        if (imageWidth > boxWidth) {
          minX =
            (appliedZoom * imageWidth - imageWidth) / 2 -
            appliedZoom * imageWidth +
            boxWidth;
          maxX = (appliedZoom * imageWidth - imageWidth) / 2;
        } else {
          minX = (boxWidth - appliedZoom * imageWidth) / 2;
          maxX = (appliedZoom * imageWidth - boxWidth) / 2;
        }
      } else {
        minX = Math.min((boxWidth - imageWidth) / 2, 0);
        maxX = minX;
      }

      return [minX, maxX];
    },
    [imageWidth, boxWidth]
  );

  const calculateInitialPosition = useCallback(
    (appliedZoom: number) => {
      // Center image from container's center
      const newPositionX = Math.min((boxWidth - imageWidth) / 2, 0);
      let newPositionY: number;

      if (displayMode === GQL.ImageLightboxDisplayMode.FitXy) {
        newPositionY = Math.min((boxHeight - imageHeight) / 2, 0);
      } else {
        // otherwise, align image with container
        const [minY, maxY] = minMaxY(appliedZoom);
        if (!alignBottom) {
          newPositionY = maxY;
        } else {
          newPositionY = minY;
        }
      }

      return [newPositionX, newPositionY];
    },
    [
      displayMode,
      boxWidth,
      imageWidth,
      boxHeight,
      imageHeight,
      alignBottom,
      minMaxY,
    ]
  );

  // identity of the currently-positioned view; a change means a genuinely new
  // image/display option (reset to initial), as opposed to a pure box resize.
  const viewKey = useRef<string>("");
  // last observed box size, to tell a resize from an unrelated re-run.
  const prevBox = useRef<{ w: number; h: number }>({ w: 0, h: 0 });

  useEffect(() => {
    // don't set anything until we have the dimensions
    if (!imageWidth || !imageHeight || !boxWidth || !boxHeight) {
      return;
    }

    const smallUnscaled =
      !scaleUp && imageWidth < boxWidth && imageHeight < boxHeight;
    // the fit baseline for the *current* box
    const newZoom = smallUnscaled
      ? 1
      : calculateDefaultZoom(
          imageWidth,
          imageHeight,
          boxWidth,
          boxHeight,
          displayMode,
          scaleUp
        );

    const key = `${src}|${displayMode}|${scaleUp}|${alignBottom}`;
    const oldBoxW = prevBox.current.w;
    const oldBoxH = prevBox.current.h;
    const boxResized = boxWidth !== oldBoxW || boxHeight !== oldBoxH;
    prevBox.current = { w: boxWidth, h: boxHeight };

    // refit to the box and reset to the initial centered position
    const resetToInitial = () => {
      setDefaultZoom(newZoom);
      if (smallUnscaled) {
        setPositionX(0);
        setPositionY(0);
      } else {
        const [nx, ny] = calculateInitialPosition(newZoom * 1);
        setPositionX(nx);
        setPositionY(ny);
      }
      scrollAttempts.current = alignBottom
        ? scrollAttemptsBeforeChange
        : -scrollAttemptsBeforeChange;
    };

    if (key !== viewKey.current) {
      // genuinely new image / display option
      viewKey.current = key;
      resetToInitial();
    } else if (boxResized) {
      if (Math.abs(zoom - 1) < 0.05) {
        // viewing at fit: refit + recentre for the new orientation
        resetToInitial();
      } else {
        // zoomed in: preserve the absolute scale (keep defaultZoom, so the image
        // stays the same size and still overflows = pannable) and the focal
        // point. `position` is a screen-space translation, so shifting it by half
        // the box-size change keeps the content under the viewport centre fixed
        // (scale/layout-independent); then clamp into the new bounds.
        const applied = defaultZoom * zoom;
        const dx = (boxWidth - oldBoxW) / 2;
        const dy = (boxHeight - oldBoxH) / 2;
        const [minX, maxX] = minMaxX(applied);
        const [minY, maxY] = minMaxY(applied);
        setPositionX((px) => Math.min(Math.max(px + dx, minX), maxX));
        setPositionY((py) => Math.min(Math.max(py + dy, minY), maxY));
      }
    } else {
      // unrelated re-run (e.g. a zoom change): keep the fit baseline in sync
      // without touching position, so the pinch/zoom path is untouched
      setDefaultZoom(newZoom);
    }
  }, [
    src,
    imageWidth,
    imageHeight,
    boxWidth,
    boxHeight,
    displayMode,
    scaleUp,
    alignBottom,
    zoom,
    defaultZoom,
    calculateInitialPosition,
    minMaxX,
    minMaxY,
    scrollAttemptsBeforeChange,
  ]);

  useEffect(() => {
    if (resetPosition !== resetPositionRef.current) {
      resetPositionRef.current = resetPosition;

      const [x, y] = calculateInitialPosition(zoom * defaultZoom);
      setPositionX(x);
      setPositionY(y);
    }
  }, [zoom, defaultZoom, resetPosition, calculateInitialPosition]);

  function getScrollMode(ev: React.WheelEvent) {
    if (ev.shiftKey) {
      switch (scrollMode) {
        case GQL.ImageLightboxScrollMode.Zoom:
          return GQL.ImageLightboxScrollMode.PanY;
        case GQL.ImageLightboxScrollMode.PanY:
          return GQL.ImageLightboxScrollMode.Zoom;
      }
    }

    return scrollMode;
  }

  function onContainerScroll(ev: React.WheelEvent) {
    // don't zoom if mouse isn't over image
    if (getScrollMode(ev) === GQL.ImageLightboxScrollMode.PanY) {
      onImageScroll(ev);
    }
  }

  function onLeftScroll(
    ev: React.WheelEvent,
    scrollable: boolean,
    infinite: boolean
  ) {
    if (infinite) {
      // for infinite scrolls, only change once per scroll "group"
      if (ev.deltaY <= -SCROLL_GROUP_THRESHOLD) {
        if (!inScrollGroup.current) {
          onLeft();
        }
      }
    } else {
      // #2535 - require additional scrolls before changing page
      if (
        !scrollable ||
        scrollAttempts.current <= -scrollAttemptsBeforeChange
      ) {
        scrollAttempts.current = 0;
        onLeft();
      } else {
        scrollAttempts.current--;
      }
    }
  }

  function onRightScroll(
    ev: React.WheelEvent,
    scrollable: boolean,
    infinite: boolean
  ) {
    if (infinite) {
      // for infinite scrolls, only change once per scroll "group"
      if (ev.deltaY >= SCROLL_GROUP_THRESHOLD) {
        if (!inScrollGroup.current) {
          onRight();
        }
      }
    } else {
      // #2535 - require additional scrolls before changing page
      if (!scrollable || scrollAttempts.current >= scrollAttemptsBeforeChange) {
        scrollAttempts.current = 0;
        onRight();
      } else {
        scrollAttempts.current++;
      }
    }
  }

  function onImageScrollPanY(ev: React.WheelEvent, infinite: boolean) {
    if (!current) return;

    const [minY, maxY] = minMaxY(zoom * defaultZoom);

    const scrollable = positionY !== maxY || positionY !== minY;

    let newPositionY: number;
    if (infinite) {
      newPositionY = positionY - ev.deltaY / SCROLL_PAN_FACTOR;
    } else {
      newPositionY =
        positionY + (ev.deltaY < 0 ? SCROLL_PAN_STEP : -SCROLL_PAN_STEP);
    }

    // #2389 - if scroll up and at top, then go to previous image
    // if scroll down and at bottom, then go to next image
    if (newPositionY > maxY && positionY === maxY) {
      onLeftScroll(ev, scrollable, infinite);
    } else if (newPositionY < minY && positionY === minY) {
      onRightScroll(ev, scrollable, infinite);
    } else {
      scrollAttempts.current = 0;

      // ensure image doesn't go offscreen
      newPositionY = Math.max(newPositionY, minY);
      newPositionY = Math.min(newPositionY, maxY);

      setPositionY(newPositionY);
    }

    ev.stopPropagation();
  }

  function onImageScroll(ev: React.WheelEvent) {
    const absDeltaY = Math.abs(ev.deltaY);
    const firstDeltaY = firstScroll.current;
    // detect infinite scrolling (mousepad, mouse with infinite scrollwheel)
    const infinite =
      // scrolling is infinite if deltaY is small
      absDeltaY < SCROLL_INFINITE_THRESHOLD ||
      // or if scroll events come quickly and the first one was small
      (firstDeltaY !== null &&
        Math.abs(firstDeltaY) < SCROLL_INFINITE_THRESHOLD);

    switch (getScrollMode(ev)) {
      case GQL.ImageLightboxScrollMode.Zoom: {
        let percent: number;
        if (infinite) {
          percent = 1 - ev.deltaY / ZOOM_FACTOR;
        } else {
          percent = ev.deltaY < 0 ? ZOOM_STEP : 1 / ZOOM_STEP;
        }
        setZoom(zoom * percent);
        break;
      }
      case GQL.ImageLightboxScrollMode.PanY:
        onImageScrollPanY(ev, infinite);
        break;
    }
    if (firstDeltaY === null) {
      firstScroll.current = ev.deltaY;
    }
    if (absDeltaY >= SCROLL_GROUP_THRESHOLD) {
      inScrollGroup.current = true;
    } else if (absDeltaY <= SCROLL_GROUP_EXIT_THRESHOLD) {
      // only "exit" the scroll group if speed has slowed considerably
      inScrollGroup.current = false;
    }
    debouncedScrollReset();
  }

  function onImageMouseOver(ev: React.MouseEvent) {
    if (!moving) return;

    if (!ev.buttons) {
      setMoving(false);
      return;
    }

    const posX = ev.pageX - startPoints.current[0];
    const posY = ev.pageY - startPoints.current[1];
    startPoints.current = [ev.pageX, ev.pageY];

    // clamp to the image bounds so a drag stops at the edge (matches touch pan)
    const appliedZoom = defaultZoom * zoom;
    const [minX, maxX] = minMaxX(appliedZoom);
    const [minY, maxY] = minMaxY(appliedZoom);
    setPositionX(Math.min(Math.max(positionX + posX, minX), maxX));
    setPositionY(Math.min(Math.max(positionY + posY, minY), maxY));
  }

  function onImageMouseDown(ev: React.MouseEvent) {
    startPoints.current = [ev.pageX, ev.pageY];
    setMoving(true);

    mouseDownEvent.current = ev.nativeEvent;
  }

  function onImageMouseUp(ev: React.MouseEvent) {
    if (ev.button !== 0) return;

    // ignore the synthesized click that follows a touch tap (iOS doesn't always
    // honour touchstart-preventDefault here) so touch taps don't navigate -
    // on touch, zoom is via double-tap and navigation via swipe. Real mouse
    // clicks have no preceding touchend and pass through.
    if (ev.timeStamp - lastTouchEnd.current < 700) return;

    if (
      !mouseDownEvent.current ||
      ev.timeStamp - mouseDownEvent.current.timeStamp > 200
    ) {
      // not a click - ignore
      return;
    }

    // must be a click
    if (
      ev.pageX !== startPoints.current[0] ||
      ev.pageY !== startPoints.current[1]
    ) {
      return;
    }

    if (ev.nativeEvent.offsetX >= (ev.target as HTMLElement).offsetWidth / 2) {
      onRight();
    } else {
      onLeft();
    }
  }

  function onTouchStart(ev: React.TouchEvent) {
    ev.preventDefault();

    // a second finger means a pinch, which is driven by the pointer handlers;
    // abandon any single-finger gesture so the two don't fight
    if (ev.touches.length !== 1 || pointerCache.current.length >= 2) {
      touchStart.current = null;
      gestureAxis.current = "none";
      setMoving(false);
      return;
    }

    const t = ev.touches[0];
    touchStart.current = { x: t.pageX, y: t.pageY, t: ev.timeStamp };
    startPoints.current = [t.pageX, t.pageY];
    gestureAxis.current = "none";
    setMoving(true);
  }

  function onTouchMove(ev: React.TouchEvent) {
    if (!moving || !touchStart.current) return;
    // pinch in progress - let the pointer handlers own it
    if (ev.touches.length !== 1 || pointerCache.current.length >= 2) return;

    const t = ev.touches[0];

    // Decide the gesture axis once, when movement first passes the lock
    // threshold. A drag pans only when the image actually overflows that axis;
    // otherwise it's a swipe candidate (navigate / delete), committed on release.
    if (gestureAxis.current === "none") {
      const dx = t.pageX - touchStart.current.x;
      const dy = t.pageY - touchStart.current.y;
      if (Math.hypot(dx, dy) < AXIS_LOCK_THRESHOLD) return;

      const appliedZoom = defaultZoom * zoom;
      if (Math.abs(dx) > Math.abs(dy)) {
        const horizPannable =
          imageWidth * appliedZoom > boxWidth + PANNABLE_EPSILON;
        gestureAxis.current = horizPannable ? "pan" : "horizontal";
      } else {
        const vertPannable =
          imageHeight * appliedZoom > boxHeight + PANNABLE_EPSILON;
        gestureAxis.current = vertPannable ? "pan" : "vertical";
      }
    }

    if (gestureAxis.current === "pan") {
      const posX = t.pageX - startPoints.current[0];
      const posY = t.pageY - startPoints.current[1];
      startPoints.current = [t.pageX, t.pageY];

      // clamp to the image bounds so a drag stops at the edge instead of
      // flinging the image off the viewport
      const appliedZoom = defaultZoom * zoom;
      const [minX, maxX] = minMaxX(appliedZoom);
      const [minY, maxY] = minMaxY(appliedZoom);
      setPositionX(Math.min(Math.max(positionX + posX, minX), maxX));
      setPositionY(Math.min(Math.max(positionY + posY, minY), maxY));
    }
  }

  function onTouchEnd(ev: React.TouchEvent) {
    setMoving(false);
    lastTouchEnd.current = ev.timeStamp;

    const start = touchStart.current;
    const axis = gestureAxis.current;
    touchStart.current = null;
    gestureAxis.current = "none";

    // pinch end, or a gesture we never started tracking
    if (!start || pointerCache.current.length >= 2) return;
    // a pan already applied its movement during touchmove
    if (axis === "pan") return;

    const t = ev.changedTouches[0];
    if (!t) return;

    const dx = t.pageX - start.x;
    const dy = t.pageY - start.y;
    const dt = Math.max(1, ev.timeStamp - start.t);

    // little movement: it's a tap (single-tap is a no-op, double-tap zooms).
    // Pass client coords - the double-tap focal zoom anchors on them and matches
    // them against the container's getBoundingClientRect (also client-space).
    if (Math.hypot(dx, dy) < TAP_MAX_MOVE) {
      handleTap(t.clientX, t.clientY, ev.timeStamp);
      return;
    }

    // committed once it travels far enough or is flicked fast enough
    const committed = (d: number) =>
      Math.abs(d) > SWIPE_COMMIT_DISTANCE || Math.abs(d) / dt > SWIPE_VELOCITY;

    if (axis === "horizontal" && navigationEnabled) {
      if (committed(dx)) {
        if (dx > 0) onLeft();
        else onRight();
      }
    } else if (axis === "vertical" && committed(dy)) {
      // swipe up opens the delete confirmation dialog; swipe down closes the
      // lightbox. Both only fire when the
      // image isn't vertically pannable - i.e. at fit zoom - since a zoomed-in
      // vertical drag locks to "pan" instead.
      if (dy < 0) {
        onSwipeDelete?.();
      } else {
        onSwipeClose?.();
      }
    }
  }

  function handleTap(x: number, y: number, time: number) {
    const prev = lastTap.current;
    if (
      prev &&
      time - prev.t < DOUBLE_TAP_MS &&
      Math.hypot(x - prev.x, y - prev.y) < DOUBLE_TAP_SLOP
    ) {
      lastTap.current = null;
      toggleDoubleTapZoom(x, y);
      return;
    }
    lastTap.current = { x, y, t: time };
    // A single tap on touch is intentionally a no-op: navigation is via swipe
    // and zoom via double-tap. Acting on the first tap here would make a
    // double-tap impossible (or require a 300ms tap-delay hack). Desktop
    // left/right-half click navigation (onImageMouseUp) is unaffected.
  }

  function toggleDoubleTapZoom(tapX: number, tapY: number) {
    if (Math.abs(zoom - 1) >= 0.05) {
      // already zoomed: reset to fit and re-centre
      setZoom(1);
      onResetPosition?.();
      return;
    }

    // At fit: zoom toward the tapped point to 1:1 native (à la PhotoSwipe).
    // defaultZoom is the fit scale, so 1/defaultZoom is exactly 100% native;
    // clamp to the pinch ceiling, and to a sensible minimum so the gesture
    // always does something on images that already fit.
    const target = Math.min(Math.max(1 / (defaultZoom || 1), 2), maxZoom);
    const box = container.current?.getBoundingClientRect();
    runZoomIn(() => {
      if (box) {
        // anchor the zoom on the tap point (same focal math as pinch)
        const r = target / zoom;
        const c0x =
          box.left + (imageWidth > box.width ? imageWidth : box.width) / 2;
        const c0y =
          box.top + (imageHeight > box.height ? imageHeight : box.height) / 2;
        setPositionX((tapX - c0x) * (1 - r) + r * positionX);
        setPositionY((tapY - c0y) * (1 - r) + r * positionY);
      }
      setZoom(target);
    });
  }

  // Apply an expensive zoom-in. The first zoom of a large image triggers a
  // multi-second full-res re-raster on the main thread (on iOS Safari it jumps
  // straight to sharp afterwards, with no feedback meanwhile). We cover that
  // stall with a busy spinner that lives on its own compositor layer (see the
  // overlay JSX): its opacity + rotate animation are composited independently of
  // the main thread, so it keeps spinning while the main thread is blocked.
  //
  // We deliberately do NOT pre-decode the bitmap: img.decode() forces the whole
  // full-res image (~width*height*4 bytes, hundreds of MB) into RAM regardless
  // of how far you zoom, which is redundant with the raster the transform
  // triggers anyway and just raises peak memory. Letting the browser raster
  // on-demand keeps the footprint lower on every device.
  function runZoomIn(apply: () => void) {
    if (!isLargeImage || zoomedOnce.current) {
      apply();
      return;
    }
    zoomedOnce.current = true;
    // Reveal the spinner, yield a real frame (rAF then setTimeout(0) - that
    // hands the opacity change to the compositor before apply() blocks;
    // double-rAF is unreliable on iOS), then apply the zoom and clear once the
    // heavy paint lands.
    showZoomBusy();
    requestAnimationFrame(() =>
      window.setTimeout(() => {
        apply();
        requestAnimationFrame(() => clearZoomBusy());
      }, 0)
    );
  }

  function onPointerDown(ev: React.PointerEvent) {
    // replace pointer event with the same id, if applicable
    pointerCache.current = pointerCache.current.filter(
      (e) => e.pointerId !== ev.pointerId
    );

    pointerCache.current.push(ev);
    // (re)start the pinch snapshot; it's captured on the first 2-pointer move
    pinch.current = null;
    // unmount the busy spinner while pinching (see `pinching`)
    if (pointerCache.current.length >= 2) setPinching(true);
  }

  function onPointerUp(ev: React.PointerEvent) {
    for (let i = 0; i < pointerCache.current.length; i++) {
      if (pointerCache.current[i].pointerId === ev.pointerId) {
        pointerCache.current.splice(i, 1);
        break;
      }
    }
    if (pointerCache.current.length < 2) {
      pinch.current = null;
      // pinch over: remount the spinner so its animation is warm for a later tap
      setPinching(false);
    }
  }

  // Pinch-to-zoom, computed ABSOLUTELY from a snapshot taken at the start of the
  // gesture (à la timmywil/panzoom) rather than accumulated per frame. Each move
  // derives the target zoom and pan purely from the snapshot + the live finger
  // distance/midpoint, so there's no frame-to-frame drift, no exponential
  // divergence, and no per-frame getBoundingClientRect (all of which made the
  // incremental version fling the image off-screen / crash the tab).
  function onPointerMove(ev: React.PointerEvent) {
    // find the event in the cache
    const cachedIndex = pointerCache.current.findIndex(
      (c) => c.pointerId === ev.pointerId
    );
    if (cachedIndex !== -1) {
      pointerCache.current[cachedIndex] = ev;
    }

    if (pointerCache.current.length !== 2) return;

    const ev1 = pointerCache.current[0];
    const ev2 = pointerCache.current[1];
    const dist = Math.hypot(
      ev1.clientX - ev2.clientX,
      ev1.clientY - ev2.clientY
    );
    const midX = (ev1.clientX + ev2.clientX) / 2;
    const midY = (ev1.clientY + ev2.clientY) / 2;

    if (!pinch.current) {
      // Snapshot once: the finger distance, the current zoom/pan, the focal
      // (the initial midpoint we anchor on), and the wrapper's untransformed
      // centre. The image lays out at natural size, so it aligns flex-start
      // when it overflows the box (the usual case) and is margin:auto-centred
      // only when smaller.
      const box = container.current?.getBoundingClientRect();
      pinch.current = {
        dist: dist || 1,
        zoom,
        posX: positionX,
        posY: positionY,
        fx: midX,
        fy: midY,
        c0x: box
          ? box.left + (imageWidth > box.width ? imageWidth : box.width) / 2
          : midX,
        c0y: box
          ? box.top + (imageHeight > box.height ? imageHeight : box.height) / 2
          : midY,
      };
      return;
    }

    const p = pinch.current;
    // zoom scales with the finger-spread ratio (natural pinch feel), clamped
    const toZoom = Math.min(
      Math.max((p.zoom * dist) / p.dist, PINCH_MIN_ZOOM),
      maxZoom
    );
    // keep the focal content point fixed: with scale about the wrapper centre,
    // pan' = (focal - centre)(1 - r) + r * panAtStart, where r = scale ratio
    // from the start of the gesture. Bounded because toZoom is clamped.
    const r = toZoom / p.zoom;
    setPositionX((p.fx - p.c0x) * (1 - r) + r * p.posX);
    setPositionY((p.fy - p.c0y) * (1 - r) + r * p.posY);
    setZoom(toZoom);

    // mark the first zoom done so a later double-tap on this image skips the
    // spinner (the pinch has already decoded it). No spinner for pinch itself -
    // the live gesture is its own feedback.
    if (isLargeImage && !zoomedOnce.current) zoomedOnce.current = true;
  }

  const ImageView = isVideo ? "video" : "img";

  return (
    <div
      ref={container}
      className={`${CLASSNAME_IMAGE}`}
      onWheel={(e) => onContainerScroll(e)}
    >
      {defaultZoom ? (
        /* The transform is applied to this wrapper rather than the <img>
           to work around a Safari rendering bug: `transform: scale` on
           an <img> with very large intrinsic dimensions distorts the
           image's aspect ratio. See #5087. */
        <div
          className={`${CLASSNAME_IMAGE}-wrapper`}
          style={{
            transform: `translate(${positionX}px, ${positionY}px) scale(${
              defaultZoom * zoom
            })`,
          }}
        >
          <source srcSet={src} media="(min-width: 800px)" />
          {/* XXbiome-ignore jsx-a11y/no-noninteractive-element-interactions */}
          <ImageView
            loop={isVideo}
            src={src}
            alt=""
            draggable={false}
            style={{ touchAction: "none" }}
            onWheel={current ? (e) => onImageScroll(e) : undefined}
            onMouseDown={onImageMouseDown}
            onMouseUp={onImageMouseUp}
            onMouseMove={onImageMouseOver}
            onTouchStart={onTouchStart}
            onTouchMove={onTouchMove}
            onTouchEnd={onTouchEnd}
            onTouchCancel={onTouchEnd}
            onPointerDown={onPointerDown}
            onPointerUp={onPointerUp}
            onPointerMove={onPointerMove}
          />
        </div>
      ) : undefined}
      {current && !pinching ? (
        // Busy indicator for the first zoom-in of a large image (fixed-centred
        // over the viewport; non-interactive so it never blocks the gesture).
        // Mounted (with its rotate animation already running) and toggled via
        // `opacity` so it sits on its own compositor layer (translateZ(0)) - the
        // animation must be live on the compositor *before* the zoom blocks the
        // main thread, or it never paints during the stall (a freshly-mounted
        // one doesn't commit in time). It's unmounted during a pinch only, where
        // that perpetual animation would otherwise fight the pinch's compositing.
        <div
          aria-hidden={!zoomBusy}
          style={{
            position: "fixed",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%) translateZ(0)",
            zIndex: 10,
            pointerEvents: "none",
            opacity: zoomBusy ? 1 : 0,
            transition: "opacity 80ms linear",
          }}
        >
          <LoadingIndicator inline message="" />
        </div>
      ) : undefined}
    </div>
  );
};
