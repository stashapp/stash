import React, {
  MutableRefObject,
  useEffect,
  useRef,
  useState,
  useCallback,
} from "react";
import useResizeObserver from "@react-hook/resize-observer";
import * as GQL from "src/core/generated-graphql";
import cx from "classnames";

const ZOOM_STEP = 1.1;
const ZOOM_FACTOR = 700;
const SCROLL_GROUP_THRESHOLD = 8;
const SCROLL_GROUP_EXIT_THRESHOLD = 4;
const SCROLL_INFINITE_THRESHOLD = 10;
const SCROLL_PAN_STEP = 75;
const SCROLL_PAN_FACTOR = 2;
const CLASSNAME = "Lightbox";
const CLASSNAME_CAROUSEL = `${CLASSNAME}-carousel`;
const CLASSNAME_IMAGE = `${CLASSNAME_CAROUSEL}-image`;
const CLASSNAME_IMAGE_PAN = `${CLASSNAME_IMAGE}-pan`;

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

export const useContainerDimensions = <
  T extends HTMLElement = HTMLDivElement
>(): [MutableRefObject<T | null>, IDimension] => {
  const target = useRef<T | null>(null);
  const [dimension, setDimension] = useState<IDimension>({
    width: 0,
    height: 0,
  });

  useResizeObserver(target, (entry) => {
    const { inlineSize: width, blockSize: height } = entry.contentBoxSize[0];
    setDimension({ width, height });
  });

  return [target, dimension];
};

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
  moveCarousel: (v: number) => void;
  releaseCarousel: (
    ev: React.PointerEvent,
    swipeDuration: number,
    cancelled: boolean
  ) => void;
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
  moveCarousel,
  releaseCarousel,
  isVideo,
}) => {
  const [defaultZoom, setDefaultZoom] = useState<number | null>(null);
  const [positionX, setPositionX] = useState(0);
  const [positionY, setPositionY] = useState(0);
  const [imageWidth, setImageWidth] = useState(width);
  const [imageHeight, setImageHeight] = useState(height);
  const dimensionsProvided = width > 0 && height > 0;
  const [containerRef, { width: boxWidth, height: boxHeight }] =
    useContainerDimensions();

  const resetPositionRef = useRef(resetPosition);

  // Panning and swipe navigation are tracked in startPoint. Pinch zoom is
  // tracked in prevDiff. They are undefined if no action of that type is in
  // progress.
  const startPoint = useRef<number[] | undefined>();
  const startTime = useRef<number>(0);
  const pointerCache = useRef<React.PointerEvent[]>([]);
  const prevDiff = useRef<number | undefined>();

  const scrollAttempts = useRef(0);

  useEffect(() => {
    function toggleVideoPlay() {
      if (containerRef.current) {
        let openVideo = containerRef.current.getElementsByTagName("video");
        if (openVideo.length > 0) {
          let rect = openVideo[0].getBoundingClientRect();
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
  }, [containerRef]);

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

  const calcPanBounds = useCallback(
    (appliedZoom: number) => {
      const xRange = Math.max(appliedZoom * imageWidth - boxWidth, 0);
      const yRange = Math.max(appliedZoom * imageHeight - boxHeight, 0);
      const nonZero = xRange != 0 || yRange != 0;
      return {
        minX: -xRange / 2,
        maxX: xRange / 2,
        minY: -yRange / 2,
        maxY: yRange / 2,
        nonZero,
      };
    },
    [imageWidth, boxWidth, imageHeight, boxHeight]
  );
  const panBounds =
    defaultZoom !== null
      ? calcPanBounds(defaultZoom * zoom)
      : { minX: 0, maxX: 0, minY: 0, maxY: 0, nonZero: false };

  const minMaxY = useCallback(
    (appliedZoom: number) => {
      const minY = Math.min((boxHeight - appliedZoom * imageHeight) / 2, 0);
      const maxY = Math.max((appliedZoom * imageHeight - boxHeight) / 2, 0);

      return [minY, maxY];
    },
    [imageHeight, boxHeight]
  );

  const calculateInitialPosition = useCallback(
    (appliedZoom: number) => {
      // If image is smaller than container, place in center. Otherwise, align
      // the left side of the image with the left side of the container, and
      // align either the top or bottom of the image with the corresponding
      // edge of container, depending on whether navigation is forwards or
      // backwards.
      const [minY, maxY] = minMaxY(appliedZoom);
      const newPositionX = Math.max(
        (appliedZoom * imageWidth - boxWidth) / 2,
        0
      );
      const newPositionY = alignBottom ? minY : maxY;

      return [newPositionX, newPositionY];
    },
    [boxWidth, imageWidth, alignBottom, minMaxY]
  );

  useEffect(() => {
    // don't set anything until we have the dimensions
    if (!imageWidth || !imageHeight || !boxWidth || !boxHeight) {
      return;
    }

    if (!scaleUp && imageWidth < boxWidth && imageHeight < boxHeight) {
      setDefaultZoom(1);
      setPositionX(0);
      setPositionY(0);
      return;
    }

    // set initial zoom level based on options
    const newZoom = calculateDefaultZoom(
      imageWidth,
      imageHeight,
      boxWidth,
      boxHeight,
      displayMode,
      scaleUp
    );

    setDefaultZoom(newZoom);

    const [newPositionX, newPositionY] = calculateInitialPosition(newZoom * 1);

    setPositionX(newPositionX);
    setPositionY(newPositionY);

    if (alignBottom) {
      scrollAttempts.current = scrollAttemptsBeforeChange;
    } else {
      scrollAttempts.current = -scrollAttemptsBeforeChange;
    }
  }, [
    imageWidth,
    imageHeight,
    boxWidth,
    boxHeight,
    displayMode,
    scaleUp,
    alignBottom,
    calculateInitialPosition,
    scrollAttemptsBeforeChange,
  ]);

  useEffect(() => {
    if (defaultZoom === null) {
      return;
    }
    if (resetPosition !== resetPositionRef.current) {
      resetPositionRef.current = resetPosition;

      const [x, y] = calculateInitialPosition(zoom * defaultZoom);
      setPositionX(x);
      setPositionY(y);
    }
  }, [
    zoom,
    defaultZoom,
    resetPosition,
    resetPositionRef,
    calculateInitialPosition,
  ]);

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
    if (!current || defaultZoom === null) return;

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
    if (defaultZoom === null) return;

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
      case GQL.ImageLightboxScrollMode.Zoom:
        let percent: number;
        if (infinite) {
          percent = 1 - ev.deltaY / ZOOM_FACTOR;
        } else {
          percent = ev.deltaY < 0 ? ZOOM_STEP : 1 / ZOOM_STEP;
        }
        setZoom(zoom * percent);
        const bounds = calcPanBounds(defaultZoom * zoom * percent);
        setPositionX(Math.max(bounds.minX, Math.min(bounds.maxX, positionX)));
        setPositionY(Math.max(bounds.minY, Math.min(bounds.maxY, positionY)));
        break;
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

  function onPointerDown(ev: React.PointerEvent) {
    // replace pointer event with the same id, if applicable
    pointerCache.current = pointerCache.current.filter(
      (e) => e.pointerId !== ev.pointerId
    );

    pointerCache.current.push(ev);
    prevDiff.current = undefined;

    startTime.current = ev.timeStamp;
    if (pointerCache.current.length === 1) {
      startPoint.current = [ev.clientX, ev.clientY];
    } else if (
      pointerCache.current.length === 2 &&
      startPoint.current !== undefined
    ) {
      const centerX = Math.abs(ev.clientX + startPoint.current[0]) / 2;
      const centerY = Math.abs(ev.clientY + startPoint.current[1]) / 2;
      startPoint.current = [centerX, centerY];
    }
  }

  function onPointerUp(ev: React.PointerEvent) {
    let found = false;

    for (let i = 0; i < pointerCache.current.length; i++) {
      if (pointerCache.current[i].pointerId === ev.pointerId) {
        pointerCache.current.splice(i, 1);
        found = true;
        break;
      }
    }

    if (!found || pointerCache.current.length !== 0) {
      if (pointerCache.current.length === 1) {
        // If we are transitioning from pinch zoom to pan, reset this
        // so we don't pan relative to the old center point.
        startPoint.current = [
          pointerCache.current[0].clientX,
          pointerCache.current[0].clientY,
        ];
      }
      return;
    }

    if (ev.pointerType === "touch" && startPoint.current !== null) {
      // Swipe navigation
      releaseCarousel(ev, ev.timeStamp - startTime.current, false);
    }

    if (
      ev.button === 0 &&
      ev.timeStamp - startTime.current <= 200 &&
      startPoint.current !== undefined &&
      ev.clientX === startPoint.current[0] &&
      ev.clientY === startPoint.current[1]
    ) {
      // Click or tap navigation

      if (ev.clientX >= window.innerWidth / 2) {
        onRight();
      } else {
        onLeft();
      }
    }
  }

  function onPointerCancel(ev: React.PointerEvent) {
    for (let i = 0; i < pointerCache.current.length; i++) {
      if (pointerCache.current[i].pointerId === ev.pointerId) {
        pointerCache.current.splice(i, 1);
        if (ev.pointerType === "touch" && pointerCache.current.length === 0) {
          releaseCarousel(ev, ev.timeStamp - startTime.current, true);
        }
        return;
      }
    }
  }

  function onPointerMove(ev: React.PointerEvent) {
    // find the event in the cache
    const cachedIndex = pointerCache.current.findIndex(
      (c) => c.pointerId === ev.pointerId
    );

    if (cachedIndex === -1 || defaultZoom === null) return;

    pointerCache.current[cachedIndex] = ev;

    if (pointerCache.current.length === 2 && startPoint.current !== undefined) {
      // Pinch zoom

      // compare the difference between the two pointers
      const ev1 = pointerCache.current[0];
      const ev2 = pointerCache.current[1];
      const diffX = Math.abs(ev1.clientX - ev2.clientX);
      const diffY = Math.abs(ev1.clientY - ev2.clientY);
      const diff = Math.sqrt(diffX ** 2 + diffY ** 2);
      const centerX = Math.abs(ev1.clientX + ev2.clientX) / 2;
      const centerY = Math.abs(ev1.clientY + ev2.clientY) / 2;
      const deltaX = centerX - startPoint.current[0];
      const deltaY = centerY - startPoint.current[1];
      startPoint.current = [centerX, centerY];

      if (prevDiff.current !== undefined) {
        const diffDiff = diff - prevDiff.current;
        const factor = (Math.abs(diffDiff) / 20) * 0.1 + 1;

        let newZoom = diffDiff > 0 ? zoom * factor : zoom / factor;
        setZoom(newZoom);
        const bounds = calcPanBounds(defaultZoom * newZoom);

        const newPositionX = Math.max(
          bounds.minX,
          Math.min(
            bounds.maxX,
            (diffDiff > 0 ? positionX * factor : positionX / factor) + deltaX
          )
        );
        const newPositionY = Math.max(
          bounds.minY,
          Math.min(
            bounds.maxY,
            (diffDiff > 0 ? positionY * factor : positionY / factor) + deltaY
          )
        );

        setPositionX(newPositionX);
        setPositionY(newPositionY);
      }

      prevDiff.current = diff;
    } else if (
      pointerCache.current.length === 1 &&
      startPoint.current !== undefined &&
      pointerCache.current[0].pointerType === "touch" &&
      panBounds.minX === panBounds.maxX
    ) {
      // Swipe navigation (touch only, and only when panning is not possible)
      const deltaX = ev.clientX - startPoint.current[0];
      startPoint.current = [ev.clientX, ev.clientY];
      moveCarousel(deltaX);
    } else if (
      pointerCache.current.length === 1 &&
      startPoint.current !== undefined
    ) {
      // Panning

      if (!ev.buttons) {
        return;
      }

      const deltaX = ev.clientX - startPoint.current[0];
      const deltaY = ev.clientY - startPoint.current[1];
      startPoint.current = [ev.clientX, ev.clientY];

      const newPositionX = Math.max(
        panBounds.minX,
        Math.min(panBounds.maxX, positionX + deltaX)
      );
      const newPositionY = Math.max(
        panBounds.minY,
        Math.min(panBounds.maxY, positionY + deltaY)
      );

      setPositionX(newPositionX);
      setPositionY(newPositionY);
    }
  }

  const ImageView = isVideo ? "video" : "img";

  return (
    <div
      ref={containerRef}
      className={cx(CLASSNAME_IMAGE, {
        [CLASSNAME_IMAGE_PAN]: panBounds.nonZero,
      })}
      style={{ touchAction: "none" }}
      onWheel={(e) => onContainerScroll(e)}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={onPointerUp}
      onPointerCancel={onPointerCancel}
    >
      {defaultZoom ? (
        <picture>
          {/* eslint-disable-next-line jsx-a11y/no-noninteractive-element-interactions */}
          <ImageView
            style={{
              touchAction: "none",
              position: "relative",
              left: "50%",
              top: "50%",
              transform: `translate(-50%, -50%) translate(${positionX}px, ${positionY}px) scale(${
                defaultZoom * zoom
              })`,
            }}
            loop={isVideo}
            src={src}
            alt=""
            draggable={false}
            onWheel={current ? (e) => onImageScroll(e) : undefined}
          />
        </picture>
      ) : undefined}
    </div>
  );
};
