import { Card, Button, ButtonGroup, Form } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { FormattedPlural } from "react-intl";
import { HoverPopover, Icon, TagLink } from "../Shared";
import { useConfiguration } from "src/core/StashService";

interface IProps {
  gallery: GQL.GalleryDataFragment;
  selecting?: boolean;
  selected: boolean | undefined;
  zoomIndex: number;
  onSelectedChanged: (selected: boolean, shiftKey: boolean) => void;
}

export const GalleryCard: React.FC<IProps> = (props) => {
  const config = useConfiguration();
  const showStudioAsText =
    config?.data?.configuration.interface.showStudioAsText ?? false;
    
  function maybeRenderScenePopoverButton() {
    if (!props.gallery.scene) return;

    const popoverContent = (
      <TagLink key={props.gallery.scene.id} scene={props.gallery.scene} />
    );

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Link to={`/scenes/${props.gallery.scene.id}`}>
          <Button className="minimal">
            <Icon icon="play-circle" />
          </Button>
        </Link>
      </HoverPopover>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.gallery.tags.length <= 0) return;

    const popoverContent = props.gallery.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="tag" />
          <span>{props.gallery.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.gallery.performers.length <= 0) return;

    const popoverContent = props.gallery.performers.map((performer) => (
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
          <span>{props.gallery.performers.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderSceneStudioOverlay() {
    if (!props.gallery.studio) return;

    return (
      <div className="scene-studio-overlay">
        <Link to={`/studios/${props.gallery.studio.id}`}>
          {showStudioAsText ? (
            props.gallery.studio.name
          ) : (
            <img
              className="image-thumbnail"
              alt={props.gallery.studio.name}
              src={props.gallery.studio.image_path ?? ""}
            />
          )}
        </Link>
      </div>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (props.gallery.scene || props.gallery.performers.length > 0 || props.gallery.tags.length > 0) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderTagPopoverButton()}
            {maybeRenderPerformerPopoverButton()}
            {maybeRenderScenePopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  function maybeRenderRatingBanner() {
    if (!props.gallery.rating) {
      return;
    }
    return (
      <div
        className={`rating-banner ${
          props.gallery.rating ? `rating-${props.gallery.rating}` : ""
        }`}
      >
        RATING: {props.gallery.rating}
      </div>
    );
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

  let shiftKey = false;

  return (
    <Card className={`gallery-card zoom-${props.zoomIndex}`}>
      <Form.Control
        type="checkbox"
        className="gallery-card-check"
        checked={props.selected}
        onChange={() => props.onSelectedChanged(!props.selected, shiftKey)}
        onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => {
          // eslint-disable-next-line prefer-destructuring
          shiftKey = event.shiftKey;
          event.stopPropagation();
        }}
      />

      <div className="gallery-section">
        <Link
          to={`/galleries/${props.gallery.id}`}
          className="gallery-card-header"
          onClick={handleImageClick}
          onDragStart={handleDrag}
          onDragOver={handleDragOver}
          draggable={props.selecting}
        >
          {props.gallery.cover ? (
            <img
              className="gallery-card-image"
              alt={props.gallery.title ?? ""}
              src={`${props.gallery.cover.paths.thumbnail}`}
            />
          ) : undefined}
          {maybeRenderRatingBanner()}
        </Link>
        {maybeRenderSceneStudioOverlay()}
      </div>
      <div className="card-section">
        <Link to={`/galleries/${props.gallery.id}`}>
          <h5 className="card-section-title">
            {props.gallery.title ?? props.gallery.path}
          </h5>
        </Link>
        <span>
          {props.gallery.images.length}&nbsp;
          <FormattedPlural
            value={props.gallery.images.length ?? 0}
            one="image"
            other="images"
          />
          .
        </span>
      </div>
      {maybeRenderPopoverButtonGroup()}
    </Card>
  );
};
