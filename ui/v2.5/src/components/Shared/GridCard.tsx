import React from "react";
import { Card, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import cx from "classnames";
import TruncatedText from "./TruncatedText";

interface ICardProps {
  className?: string;
  linkClassName?: string;
  thumbnailSectionClassName?: string;
  url: string;
  pretitleIcon?: JSX.Element;
  title: string;
  image: JSX.Element;
  details?: JSX.Element;
  overlays?: JSX.Element;
  popovers?: JSX.Element;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  interactiveHeatmap?: string;
}

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
          className="card-check"
          checked={props.selected}
          onChange={() => props.onSelectedChanged!(!props.selected, shiftKey)}
          onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => {
            // eslint-disable-next-line prefer-destructuring
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
          src={props.interactiveHeatmap}
          alt="interactive heatmap"
          className="interactive-heatmap"
        />
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
