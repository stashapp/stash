import React, { MutableRefObject, useRef, useState } from "react";
import { Card, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import cx from "classnames";
import { TruncatedText } from "../TruncatedText";
import ScreenUtils from "src/utils/screen";
import useResizeObserver from "@react-hook/resize-observer";

interface ICardProps {
  className?: string;
  linkClassName?: string;
  thumbnailSectionClassName?: string;
  width?: number;
  url: string;
  pretitleIcon?: JSX.Element;
  title: JSX.Element | string;
  image: JSX.Element;
  details?: JSX.Element;
  overlays?: JSX.Element;
  popovers?: JSX.Element;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  resumeTime?: number;
  duration?: number;
  interactiveHeatmap?: string;
}

export const calculateCardWidth = (
  containerWidth: number,
  preferredWidth: number
) => {
  const containerPadding = 30;
  const cardMargin = 10;
  let maxUsableWidth = containerWidth - containerPadding;
  let maxElementsOnRow = Math.ceil(maxUsableWidth / preferredWidth);
  return maxUsableWidth / maxElementsOnRow - cardMargin;
};

interface IDimension {
  width: number;
  height: number;
}

export const useContainerDimensions = <T extends HTMLElement = HTMLDivElement>(
  sensitivityThreshold = 20
): [MutableRefObject<T | null>, IDimension] => {
  const target = useRef<T | null>(null);
  const [dimension, setDimension] = useState<IDimension>({
    width: 0,
    height: 0,
  });

  useResizeObserver(target, (entry) => {
    const { inlineSize: width, blockSize: height } = entry.contentBoxSize[0];
    let difference = Math.abs(dimension.width - width);
    // Only adjust when width changed by a significant margin. This addresses the cornercase that sees
    // the dimensions toggle back and forward when the window is adjusted perfectly such that overflow
    // is trigger then immediable disabled because of a resize event then continues this loop endlessly.
    // the scrollbar size varies between platforms. Windows is apparently around 17 pixels.
    if (difference > sensitivityThreshold) {
      setDimension({ width, height });
    }
  });

  return [target, dimension];
};

export const GridCard: React.FC<ICardProps> = (props: ICardProps) => {
  function handleImageClick(event: React.MouseEvent<HTMLElement, MouseEvent>) {
    const { shiftKey } = event;

    if (!props.onSelectedChanged) {
      return;
    }

    if (props.selecting) {
      props.onSelectedChanged(!props.selected, shiftKey);
      event.preventDefault();
    }
  }

  function handleDrag(event: React.DragEvent<HTMLElement>) {
    if (props.selecting) {
      event.dataTransfer.setData("text/plain", "");
      event.dataTransfer.setDragImage(new Image(), 0, 0);
    }
  }

  function handleDragOver(event: React.DragEvent<HTMLElement>) {
    const ev = event;
    const shiftKey = false;

    if (!props.onSelectedChanged) {
      return;
    }

    if (props.selecting && !props.selected) {
      props.onSelectedChanged(true, shiftKey);
    }

    ev.dataTransfer.dropEffect = "move";
    ev.preventDefault();
  }

  let shiftKey = false;

  function maybeRenderCheckbox() {
    if (props.onSelectedChanged) {
      return (
        <Form.Control
          type="checkbox"
          // #2750 - add mousetrap class to ensure keyboard shortcuts work
          className="card-check mousetrap"
          checked={props.selected}
          onChange={() => props.onSelectedChanged!(!props.selected, shiftKey)}
          onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => {
            shiftKey = event.shiftKey;
            event.stopPropagation();
          }}
        />
      );
    }
  }

  function maybeRenderInteractiveHeatmap() {
    if (props.interactiveHeatmap) {
      return (
        <img
          loading="lazy"
          src={props.interactiveHeatmap}
          alt="interactive heatmap"
          className="interactive-heatmap"
        />
      );
    }
  }

  function maybeRenderProgressBar() {
    if (
      props.resumeTime &&
      props.duration &&
      props.duration > props.resumeTime
    ) {
      const percentValue = (100 / props.duration) * props.resumeTime;
      const percentStr = percentValue + "%";
      return (
        <div title={Math.round(percentValue) + "%"} className="progress-bar">
          <div style={{ width: percentStr }} className="progress-indicator" />
        </div>
      );
    }
  }

  return (
    <Card
      className={cx(props.className, "grid-card")}
      onClick={handleImageClick}
      onDragStart={handleDrag}
      onDragOver={handleDragOver}
      draggable={props.onSelectedChanged && props.selecting}
      style={
        props.width && !ScreenUtils.isMobile()
          ? { width: `${props.width}px` }
          : {}
      }
    >
      {maybeRenderCheckbox()}

      <div className={cx(props.thumbnailSectionClassName, "thumbnail-section")}>
        <Link
          to={props.url}
          className={props.linkClassName}
          onClick={handleImageClick}
        >
          {props.image}
        </Link>
        {props.overlays}
        {maybeRenderProgressBar()}
      </div>
      {maybeRenderInteractiveHeatmap()}
      <div className="card-section">
        <Link to={props.url} onClick={handleImageClick}>
          <h5 className="card-section-title flex-aligned">
            {props.pretitleIcon}
            <TruncatedText text={props.title} lineCount={2} />
          </h5>
        </Link>
        {props.details}
      </div>

      {props.popovers}
    </Card>
  );
};
