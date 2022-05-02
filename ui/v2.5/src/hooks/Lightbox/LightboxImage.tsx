import React, { useEffect, useRef, useState, useCallback } from "react";
import * as GQL from "src/core/generated-graphql";

const ZOOM_STEP = 1.1;
const SCROLL_PAN_STEP = 75;
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

interface IProps {
  src: string;
  displayMode: GQL.ImageLightboxDisplayMode;
  scaleUp: boolean;
  scrollMode: GQL.ImageLightboxScrollMode;
  resetPosition?: boolean;
  zoom: number;
  scrollAttemptsBeforeChange: number;
  current: boolean;
  // set to true to align image with bottom instead of top
  alignBottom?: boolean;
  setZoom: (v: number) => void;
  onLeft: () => void;
  onRight: () => void;
}

export const LightboxImage: React.FC<IProps> = ({
  src,
  onLeft,
  onRight,
  displayMode,
  scaleUp,
  scrollMode,
  alignBottom,
  zoom,
  scrollAttemptsBeforeChange,
  current,
  setZoom,
  resetPosition,
}) => {
  const [defaultZoom, setDefaultZoom] = useState(1);
  const [moving, setMoving] = useState(false);
  const [positionX, setPositionX] = useState(0);
  const [positionY, setPositionY] = useState(0);
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);
  const [boxWidth, setBoxWidth] = useState(0);
  const [boxHeight, setBoxHeight] = useState(0);

  const mouseDownEvent = useRef<MouseEvent>();
  const resetPositionRef = useRef(resetPosition);

  const container = React.createRef<HTMLDivElement>();
  const startPoints = useRef<number[]>([0, 0]);
  const pointerCache = useRef<React.PointerEvent<HTMLDivElement>[]>([]);
  const prevDiff = useRef<number | undefined>();

  const scrollAttempts = useRef(0);

  useEffect(() => {
    const box = container.current;
    if (box) {
      setBoxWidth(box.offsetWidth);
      setBoxHeight(box.offsetHeight);
    }
  }, [container]);

  useEffect(() => {
    let mounted = true;
    const img = new Image();
    function onLoad() {
      if (mounted) {
        setWidth(img.width);
        setHeight(img.height);
      }
    }

    img.onload = onLoad;
    img.src = src;

    return () => {
      mounted = false;
    };
  }, [src]);

  const minMaxY = useCallback(
    (appliedZoom: number) => {
      let minY, maxY: number;
      const inBounds = appliedZoom * height <= boxHeight;

      // NOTE: I don't even know how these work, but they do
      if (!inBounds) {
        if (height > boxHeight) {
          minY =
            (appliedZoom * height - height) / 2 -
            appliedZoom * height +
            boxHeight;
          maxY = (appliedZoom * height - height) / 2;
        } else {
          minY = (boxHeight - appliedZoom * height) / 2;
          maxY = (appliedZoom * height - boxHeight) / 2;
        }
      } else {
        minY = Math.min((boxHeight - height) / 2, 0);
        maxY = minY;
      }

      return [minY, maxY];
    },
    [height, boxHeight]
  );

  const calculateInitialPosition = useCallback(
    (appliedZoom: number) => {
      // Center image from container's center
      const newPositionX = Math.min((boxWidth - width) / 2, 0);
      let newPositionY: number;

      if (displayMode === GQL.ImageLightboxDisplayMode.FitXy) {
        newPositionY = Math.min((boxHeight - height) / 2, 0);
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
    [displayMode, boxWidth, width, boxHeight, height, alignBottom, minMaxY]
  );

  useEffect(() => {
    // don't set anything until we have the dimensions
    if (!width || !height || !boxWidth || !boxHeight) {
      return;
    }

    if (!scaleUp && width < boxWidth && height < boxHeight) {
      setDefaultZoom(1);
      setPositionX(0);
      setPositionY(0);
      return;
    }

    // set initial zoom level based on options
    const newZoom = calculateDefaultZoom(
      width,
      height,
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
    width,
    height,
    boxWidth,
    boxHeight,
    displayMode,
    scaleUp,
    alignBottom,
    calculateInitialPosition,
    scrollAttemptsBeforeChange,
  ]);

  useEffect(() => {
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

  function getScrollMode(ev: React.WheelEvent<HTMLDivElement>) {
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

  function onContainerScroll(ev: React.WheelEvent<HTMLDivElement>) {
    // don't zoom if mouse isn't over image
    if (getScrollMode(ev) === GQL.ImageLightboxScrollMode.PanY) {
      onImageScroll(ev);
    }
  }

  function onImageScrollPanY(ev: React.WheelEvent<HTMLDivElement>) {
    if (current) {
      const [minY, maxY] = minMaxY(zoom * defaultZoom);

      const scrollable = positionY !== maxY || positionY !== minY;

      let newPositionY =
        positionY + (ev.deltaY < 0 ? SCROLL_PAN_STEP : -SCROLL_PAN_STEP);

      // #2389 - if scroll up and at top, then go to previous image
      // if scroll down and at bottom, then go to next image
      if (newPositionY > maxY && positionY === maxY) {
        // #2535 - require additional scrolls before changing page
        if (
          !scrollable ||
          scrollAttempts.current <= -scrollAttemptsBeforeChange
        ) {
          onLeft();
        } else {
          scrollAttempts.current--;
        }
      } else if (newPositionY < minY && positionY === minY) {
        // #2535 - require additional scrolls before changing page
        if (
          !scrollable ||
          scrollAttempts.current >= scrollAttemptsBeforeChange
        ) {
          onRight();
        } else {
          scrollAttempts.current++;
        }
      } else {
        scrollAttempts.current = 0;

        // ensure image doesn't go offscreen
        newPositionY = Math.max(newPositionY, minY);
        newPositionY = Math.min(newPositionY, maxY);

        setPositionY(newPositionY);
      }

      ev.stopPropagation();
    }
  }

  function onImageScroll(ev: React.WheelEvent<HTMLDivElement>) {
    const percent = ev.deltaY < 0 ? ZOOM_STEP : 1 / ZOOM_STEP;

    switch (getScrollMode(ev)) {
      case GQL.ImageLightboxScrollMode.Zoom:
        setZoom(zoom * percent);
        break;
      case GQL.ImageLightboxScrollMode.PanY:
        onImageScrollPanY(ev);
        break;
    }
  }

  function onImageMouseOver(ev: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    if (!moving) return;

    if (!ev.buttons) {
      setMoving(false);
      return;
    }

    const posX = ev.pageX - startPoints.current[0];
    const posY = ev.pageY - startPoints.current[1];
    startPoints.current = [ev.pageX, ev.pageY];

    setPositionX(positionX + posX);
    setPositionY(positionY + posY);
  }

  function onImageMouseDown(ev: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    startPoints.current = [ev.pageX, ev.pageY];
    setMoving(true);

    mouseDownEvent.current = ev.nativeEvent;
  }

  function onImageMouseUp(ev: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    if (ev.button !== 0) return;

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

  function onTouchStart(ev: React.TouchEvent<HTMLDivElement>) {
    ev.preventDefault();
    if (ev.touches.length === 1) {
      startPoints.current = [ev.touches[0].pageX, ev.touches[0].pageY];
      setMoving(true);
    }
  }

  function onTouchMove(ev: React.TouchEvent<HTMLDivElement>) {
    if (!moving) return;

    if (ev.touches.length === 1) {
      const posX = ev.touches[0].pageX - startPoints.current[0];
      const posY = ev.touches[0].pageY - startPoints.current[1];
      startPoints.current = [ev.touches[0].pageX, ev.touches[0].pageY];

      setPositionX(positionX + posX);
      setPositionY(positionY + posY);
    }
  }

  function onPointerDown(ev: React.PointerEvent<HTMLDivElement>) {
    // replace pointer event with the same id, if applicable
    pointerCache.current = pointerCache.current.filter(
      (e) => e.pointerId !== ev.pointerId
    );

    pointerCache.current.push(ev);
    prevDiff.current = undefined;
  }

  function onPointerUp(ev: React.PointerEvent<HTMLDivElement>) {
    for (let i = 0; i < pointerCache.current.length; i++) {
      if (pointerCache.current[i].pointerId === ev.pointerId) {
        pointerCache.current.splice(i, 1);
        break;
      }
    }
  }

  function onPointerMove(ev: React.PointerEvent<HTMLDivElement>) {
    // find the event in the cache
    const cachedIndex = pointerCache.current.findIndex(
      (c) => c.pointerId === ev.pointerId
    );
    if (cachedIndex !== -1) {
      pointerCache.current[cachedIndex] = ev;
    }

    // compare the difference between the two pointers
    if (pointerCache.current.length === 2) {
      const ev1 = pointerCache.current[0];
      const ev2 = pointerCache.current[1];
      const diffX = Math.abs(ev1.clientX - ev2.clientX);
      const diffY = Math.abs(ev1.clientY - ev2.clientY);
      const diff = Math.sqrt(diffX ** 2 + diffY ** 2);

      if (prevDiff.current !== undefined) {
        const diffDiff = diff - prevDiff.current;
        const factor = (Math.abs(diffDiff) / 20) * 0.1 + 1;

        if (diffDiff > 0) {
          setZoom(zoom * factor);
        } else if (diffDiff < 0) {
          setZoom((zoom * 1) / factor);
        }
      }

      prevDiff.current = diff;
    }
  }

  return (
    <div
      ref={container}
      className={`${CLASSNAME_IMAGE}`}
      onWheel={(e) => onContainerScroll(e)}
    >
      {defaultZoom ? (
        <picture
          style={{
            transform: `translate(${positionX}px, ${positionY}px) scale(${
              defaultZoom * zoom
            })`,
          }}
        >
          <source srcSet={src} media="(min-width: 800px)" />
          {/* eslint-disable-next-line jsx-a11y/no-noninteractive-element-interactions */}
          <img
            src={src}
            alt=""
            draggable={false}
            style={{ touchAction: "none" }}
            onWheel={current ? (e) => onImageScroll(e) : undefined}
            onMouseDown={(e) => onImageMouseDown(e)}
            onMouseUp={(e) => onImageMouseUp(e)}
            onMouseMove={(e) => onImageMouseOver(e)}
            onTouchStart={(e) => onTouchStart(e)}
            onTouchMove={(e) => onTouchMove(e)}
            onPointerDown={(e) => onPointerDown(e)}
            onPointerUp={(e) => onPointerUp(e)}
            onPointerMove={(e) => onPointerMove(e)}
          />
        </picture>
      ) : undefined}
    </div>
  );
};
