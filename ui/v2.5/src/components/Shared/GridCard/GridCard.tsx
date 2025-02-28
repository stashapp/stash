import React, {
  MutableRefObject,
  PropsWithChildren,
  useRef,
  useState,
} from "react";
import { Card, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import cx from "classnames";
import { TruncatedText } from "../TruncatedText";
import ScreenUtils from "src/utils/screen";
import useResizeObserver from "@react-hook/resize-observer";
import { Icon } from "../Icon";
import { faGripLines } from "@fortawesome/free-solid-svg-icons";
import { DragSide, useDragMoveSelect } from "./dragMoveSelect";

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

  // move logic - both of the following are required to enable move dragging
  objectId?: string; // required for move dragging
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
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

const Checkbox: React.FC<{
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}> = ({ selected = false, onSelectedChanged }) => {
  let shiftKey = false;

  return (
    <Form.Control
      type="checkbox"
      // #2750 - add mousetrap class to ensure keyboard shortcuts work
      className="card-check mousetrap"
      checked={selected}
      onChange={() => onSelectedChanged!(!selected, shiftKey)}
      onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => {
        shiftKey = event.shiftKey;
        event.stopPropagation();
      }}
    />
  );
};

const DragHandle: React.FC<{
  setInHandle: (inHandle: boolean) => void;
}> = ({ setInHandle }) => {
  function onMouseEnter() {
    setInHandle(true);
  }

  function onMouseLeave() {
    setInHandle(false);
  }

  return (
    <span onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
      <Icon className="card-drag-handle" icon={faGripLines} />
    </span>
  );
};

const Controls: React.FC<PropsWithChildren<{}>> = ({ children }) => {
  return <div className="card-controls">{children}</div>;
};

const MoveTarget: React.FC<{ dragSide: DragSide }> = ({ dragSide }) => {
  if (dragSide === undefined) {
    return null;
  }

  return (
    <div
      className={`move-target move-target-${
        dragSide === DragSide.BEFORE ? "before" : "after"
      }`}
    ></div>
  );
};

export const GridCard: React.FC<ICardProps> = (props: ICardProps) => {
  const { setInHandle, moveTarget, dragProps } = useDragMoveSelect({
    selecting: props.selecting || false,
    selected: props.selected || false,
    onSelectedChanged: props.onSelectedChanged,
    objectId: props.objectId,
    onMove: props.onMove,
  });

  function handleImageClick(event: React.MouseEvent<HTMLElement, MouseEvent>) {
    const { shiftKey } = event;

    if (!props.onSelectedChanged) {
      return;
    }

    if (props.selecting) {
      props.onSelectedChanged(!props.selected, shiftKey);
      event.preventDefault();
      event.stopPropagation();
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
      {...dragProps}
      style={
        props.width && !ScreenUtils.isMobile()
          ? { width: `${props.width}px` }
          : {}
      }
    >
      {moveTarget !== undefined && <MoveTarget dragSide={moveTarget} />}
      <Controls>
        {props.onSelectedChanged && (
          <Checkbox
            selected={props.selected}
            onSelectedChanged={props.onSelectedChanged}
          />
        )}

        {!!props.objectId && props.onMove && (
          <DragHandle setInHandle={setInHandle} />
        )}
      </Controls>

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
