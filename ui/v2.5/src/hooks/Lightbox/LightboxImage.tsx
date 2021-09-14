import React, { useEffect, useRef, useState } from "react";

const ZOOM_STEP = 1.1;
const CLASSNAME = "Lightbox";
const CLASSNAME_CAROUSEL = `${CLASSNAME}-carousel`;
const CLASSNAME_IMAGE = `${CLASSNAME_CAROUSEL}-image`;

export enum DisplayMode {
  ORIGINAL = "ORIGINAL",
  FIT_XY = "FIT_XY",
  FIT_X = "FIT_X",
}

interface IProps {
  src: string;
  mode: DisplayMode;
  resetZoom?: boolean;
  onLeft: () => void;
  onRight: () => void;
  onZoomed?: () => void;
}

interface IInitialPosition {
  x: number;
  y: number;
  zoom: number;
}

export const LightboxImage: React.FC<IProps> = ({
  src,
  onLeft,
  onRight,
  mode,
  onZoomed,
  resetZoom,
}) => {
  const [zoom, setZoom] = useState(0);
  const [moving, setMoving] = useState(false);
  const [positionX, setPositionX] = useState(0);
  const [positionY, setPositionY] = useState(0);
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);
  const [boxWidth, setBoxWidth] = useState(0);
  const [boxHeight, setBoxHeight] = useState(0);

  const initialPosition = useRef<IInitialPosition | undefined>();
  const container = React.createRef<HTMLDivElement>();
  const startPoints = useRef<number[]>([0, 0]);
  const pointerCache = useRef<React.PointerEvent<HTMLDivElement>[]>([]);
  const prevDiff = useRef<number | undefined>();

  useEffect(() => {
    const box = container.current;
    if (box) {
      setBoxWidth(box.offsetWidth);
      setBoxHeight(box.offsetHeight);
    }
  }, [container]);

  useEffect(() => {
    const img = new Image();
    function onLoad() {
      setWidth(img.width);
      setHeight(img.height);
    }

    img.onload = onLoad;
    img.src = src;
  }, [src]);

  useEffect(() => {
    // don't set anything until we have the heights
    if (!width || !height || !boxWidth || !boxHeight) {
      return;
    }

    if (width < boxWidth && height < boxHeight) {
      initialPosition.current = {
        zoom: 1,
        x: 0,
        y: 0,
      };
      setZoom(1);
      setPositionX(0);
      setPositionY(0);
      return;
    }

    // set initial zoom level based on options
    let xZoom: number;
    let yZoom: number;
    let newZoom = 1;
    let newPositionY = 0;
    switch (mode) {
      case DisplayMode.FIT_XY:
        xZoom = Math.min(boxWidth / width, 1);
        yZoom = Math.min(boxHeight / height, 1);
        newZoom = Math.min(xZoom, yZoom);
        break;
      case DisplayMode.FIT_X:
        xZoom = Math.min(boxWidth / width, 1);
        newZoom = Math.min(xZoom, 1);
        break;
      case DisplayMode.ORIGINAL:
        newZoom = 1;
        break;
    }

    // Center image from container's center
    const newPositionX = Math.min((boxWidth - width) / 2, 0);

    // if fitting to screen, then centre, other
    if (mode === DisplayMode.FIT_XY) {
      newPositionY = Math.min((boxHeight - height) / 2, 0);
    } else {
      // otherwise, align top of image with container
      newPositionY = Math.min(((height * newZoom) - height) / 2, 0);
    }

    initialPosition.current = {
      zoom: newZoom,
      x: newPositionX,
      y: newPositionY,
    };
    setZoom(newZoom);
    setPositionX(newPositionX);
    setPositionY(newPositionY);
  }, [width, height, boxWidth, boxHeight, mode]);

  useEffect(() => {
    if (initialPosition.current !== undefined) {
      setZoom(initialPosition.current.zoom);
      setPositionX(initialPosition.current.x);
      setPositionY(initialPosition.current.y);
    }
  }, [resetZoom, initialPosition]);

  function zoomed() {
    if (onZoomed) {
      onZoomed();
    }
  }

  function onImageScroll(ev: React.WheelEvent<HTMLDivElement>) {
    const percent = ev.deltaY < 0 ? ZOOM_STEP : 1 / ZOOM_STEP;
    setZoom(zoom * percent);
    zoomed();
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

    const target = ev.currentTarget;
    target.addEventListener("mouseup", onImageMouseUp);
    setTimeout(() => {
      target.removeEventListener("mouseup", onImageMouseUp);
    }, 200);
  }

  function onImageMouseUp(ev: MouseEvent) {
    // must be a click
    if (
      ev.pageX !== startPoints.current[0] ||
      ev.pageY !== startPoints.current[1]
    ) {
      return;
    }

    if (ev.offsetX >= (ev.target as HTMLElement).offsetWidth / 2) {
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
          setZoom(zoom / factor);
        }
        zoomed();
      }

      prevDiff.current = diff;
    }
  }

  return (
    <div ref={container} className={`${CLASSNAME_IMAGE}`}>
      {zoom ? (
        <picture
          style={{
            transform: `translate(${positionX}px, ${positionY}px) scale(${zoom})`,
          }}
        >
          <source srcSet={src} media="(min-width: 800px)" />
          {/* eslint-disable-next-line jsx-a11y/no-noninteractive-element-interactions */}
          <img
            src={src}
            alt=""
            draggable={false}
            style={{ touchAction: "none" }}
            onWheel={(e) => onImageScroll(e)}
            onMouseDown={(e) => onImageMouseDown(e)}
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
