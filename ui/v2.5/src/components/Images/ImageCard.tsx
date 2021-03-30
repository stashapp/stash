import React from "react";
import { Button, ButtonGroup, Card, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import {
  Icon,
  TagLink,
  HoverPopover,
  SweatDrops,
  TruncatedText,
} from "src/components/Shared";
import { TextUtils } from "src/utils";

interface IImageCardProps {
  image: GQL.SlimImageDataFragment;
  selecting?: boolean;
  selected: boolean | undefined;
  zoomIndex: number;
  onSelectedChanged: (selected: boolean, shiftKey: boolean) => void;
}

export const ImageCard: React.FC<IImageCardProps> = (
  props: IImageCardProps
) => {
  function maybeRenderRatingBanner() {
    if (!props.image.rating) {
      return;
    }
    return (
      <div
        className={`rating-banner ${
          props.image.rating ? `rating-${props.image.rating}` : ""
        }`}
      >
        RATING: {props.image.rating}
      </div>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.image.tags.length <= 0) return;

    const popoverContent = props.image.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} tagType="image" />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="tag" />
          <span>{props.image.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.image.performers.length <= 0) return;

    const popoverContent = props.image.performers.map((performer) => (
      <div className="performer-tag-container row" key={performer.id}>
        <Link
          to={`/performers/${performer.id}`}
          className="performer-tag col m-auto zoom-2"
        >
          <img
            className="image-thumbnail"
            alt={performer.name ?? ""}
            src={performer.image_path ?? ""}
          />
        </Link>
        <TagLink key={performer.id} performer={performer} className="d-block" />
      </div>
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="user" />
          <span>{props.image.performers.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderOCounter() {
    if (props.image.o_counter) {
      return (
        <div>
          <Button className="minimal">
            <span className="fa-icon">
              <SweatDrops />
            </span>
            <span>{props.image.o_counter}</span>
          </Button>
        </div>
      );
    }
  }

  function maybeRenderOrganized() {
    if (props.image.organized) {
      return (
        <div>
          <Button className="minimal">
            <Icon icon="box" />
          </Button>
        </div>
      );
    }
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      props.image.tags.length > 0 ||
      props.image.performers.length > 0 ||
      props.image.o_counter ||
      props.image.organized
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderTagPopoverButton()}
            {maybeRenderPerformerPopoverButton()}
            {maybeRenderOCounter()}
            {maybeRenderOrganized()}
          </ButtonGroup>
        </>
      );
    }
  }

  function handleImageClick(
    event: React.MouseEvent<HTMLAnchorElement, MouseEvent>
  ) {
    const { shiftKey } = event;

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

    if (props.selecting && !props.selected) {
      props.onSelectedChanged(true, shiftKey);
    }

    ev.dataTransfer.dropEffect = "move";
    ev.preventDefault();
  }

  function isPortrait() {
    const { file } = props.image;
    const width = file.width ? file.width : 0;
    const height = file.height ? file.height : 0;
    return height > width;
  }

  let shiftKey = false;

  return (
    <Card className={`image-card zoom-${props.zoomIndex}`}>
      <Form.Control
        type="checkbox"
        className="image-card-check"
        checked={props.selected}
        onChange={() => props.onSelectedChanged(!props.selected, shiftKey)}
        onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => {
          // eslint-disable-next-line prefer-destructuring
          shiftKey = event.shiftKey;
          event.stopPropagation();
        }}
      />

      <div className="image-section">
        <Link
          to={`/images/${props.image.id}`}
          className="image-card-link"
          onClick={handleImageClick}
          onDragStart={handleDrag}
          onDragOver={handleDragOver}
          draggable={props.selecting}
        >
          <div className={cx("image-card-preview", { portrait: isPortrait() })}>
            <img
              className="image-card-preview-image"
              alt={props.image.title ?? ""}
              src={props.image.paths.thumbnail ?? ""}
            />
          </div>
          {maybeRenderRatingBanner()}
        </Link>
      </div>
      <div className="card-section">
        <h5 className="card-section-title">
          <TruncatedText
            text={
              props.image.title
                ? props.image.title
                : TextUtils.fileNameFromPath(props.image.path)
            }
            lineCount={2}
          />
        </h5>
      </div>

      {maybeRenderPopoverButtonGroup()}
    </Card>
  );
};
