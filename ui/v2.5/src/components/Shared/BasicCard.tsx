import React from "react";
import { Card, Form } from "react-bootstrap";
import { Link } from "react-router-dom";

interface IBasicCardProps {
  className?: string;
  linkClassName?: string;
  url: string;
  image: JSX.Element;
  details: JSX.Element;
  overlays?: JSX.Element;
  popovers?: JSX.Element;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const BasicCard: React.FC<IBasicCardProps> = (
  props: IBasicCardProps
) => {
  function handleImageClick(
    event: React.MouseEvent<HTMLAnchorElement, MouseEvent>
  ) {
    const { shiftKey } = event;

    if (!props.onSelectedChanged) {
      return;
    }

    if (props.selecting) {
      props.onSelectedChanged(!props.selected, shiftKey);
      event.preventDefault();
    }
  }

  function handleDrag(event: React.DragEvent<HTMLAnchorElement>) {
    if (props.selecting) {
      event.dataTransfer.setData("text/plain", "");
      event.dataTransfer.setDragImage(new Image(), 0, 0);
    }
  }

  function handleDragOver(event: React.DragEvent<HTMLAnchorElement>) {
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

  return (
    <Card className={props.className}>
      {maybeRenderCheckbox()}

      <div className="image-section">
        <Link
          to={props.url}
          className={props.linkClassName}
          onClick={handleImageClick}
          onDragStart={handleDrag}
          onDragOver={handleDragOver}
          draggable={props.onSelectedChanged && props.selecting}
        >
          {props.image}
        </Link>
        {props.overlays}
      </div>
      <div className="card-section">
        {props.details}
      </div>

      {props.popovers}
    </Card>
  );
};
