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
  onLeft: () => void;
  onRight: () => void;
}

export const LightboxImage: React.FC<IProps> = ({
  src,
  onLeft,
  onRight,
  mode,
}) => {
  const [zoom, setZoom] = useState(0);
  const [moving, setMoving] = useState(false);
  const [positionX, setPositionX] = useState(0);
  const [positionY, setPositionY] = useState(0);
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);
  const [boxWidth, setBoxWidth] = useState(0);
  const [boxHeight, setBoxHeight] = useState(0);

  const container = React.createRef<HTMLDivElement>();
  const startPoints = useRef<number[]>([0, 0]);

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
      setZoom(1);
      setPositionX(0);
      setPositionY(0);
      return;
    }

    // set initial zoom level based on options
    let xZoom: number;
    let yZoom: number;
    let newZoom = 1;
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
    const newPositionY = Math.min((boxHeight - height) / 2, 0);

    setZoom(newZoom);
    setPositionX(newPositionX);
    setPositionY(newPositionY);
  }, [width, height, boxWidth, boxHeight, mode]);

  function onImageScroll(ev: React.WheelEvent<HTMLDivElement>) {
    const percent = ev.deltaY > 0 ? ZOOM_STEP : 1 / ZOOM_STEP;
    setZoom(zoom * percent);
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

  return (
    /* eslint-disable-next-line jsx-a11y/no-static-element-interactions */
    <div
      ref={container}
      className={`${CLASSNAME_IMAGE}`}
      onWheel={(e) => onImageScroll(e)}
      onMouseDown={(e) => onImageMouseDown(e)}
      onMouseMove={(e) => onImageMouseOver(e)}
    >
      {zoom ? (
        <picture
          style={{
            transform: `translate(${positionX}px, ${positionY}px) scale(${zoom})`,
          }}
        >
          <source srcSet={src} media="(min-width: 800px)" />
          <img src={src} alt="" draggable={false} />
        </picture>
      ) : undefined}
    </div>
  );
};
