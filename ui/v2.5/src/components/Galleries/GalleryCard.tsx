import { Button, ButtonGroup } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  GridCard,
  HoverPopover,
  Icon,
  TagLink,
  TruncatedText,
} from "src/components/Shared";
import { PopoverCountButton } from "src/components/Shared/PopoverCountButton";
import { NavUtils, TextUtils } from "src/utils";
import { ConfigurationContext } from "src/hooks/Config";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import { RatingBanner } from "../Shared/RatingBanner";

interface IProps {
  gallery: GQL.SlimGalleryDataFragment;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const GalleryCard: React.FC<IProps> = (props) => {
  const { configuration } = React.useContext(ConfigurationContext);
  const showStudioAsText = configuration?.interface.showStudioAsText ?? false;

  function maybeRenderScenePopoverButton() {
    if (props.gallery.scenes.length === 0) return;

    const popoverContent = props.gallery.scenes.map((scene) => (
      <TagLink key={scene.id} scene={scene} />
    ));

    return (
      <HoverPopover
        className="scene-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon="play-circle" />
          <span>{props.gallery.scenes.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.gallery.tags.length <= 0) return;

    const popoverContent = props.gallery.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} tagType="gallery" />
    ));

    return (
      <HoverPopover
        className="tag-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon="tag" />
          <span>{props.gallery.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.gallery.performers.length <= 0) return;

    return <PerformerPopoverButton performers={props.gallery.performers} />;
  }

  function maybeRenderImagesPopoverButton() {
    if (!props.gallery.image_count) return;

    return (
      <PopoverCountButton
        className="image-count"
        type="image"
        count={props.gallery.image_count}
        url={NavUtils.makeGalleryImagesUrl(props.gallery)}
      />
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

  function maybeRenderOrganized() {
    if (props.gallery.organized) {
      return (
        <div className="organized">
          <Button className="minimal">
            <Icon icon="box" />
          </Button>
        </div>
      );
    }
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      props.gallery.scenes.length > 0 ||
      props.gallery.performers.length > 0 ||
      props.gallery.tags.length > 0 ||
      props.gallery.organized ||
      props.gallery.image_count > 0
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderImagesPopoverButton()}
            {maybeRenderTagPopoverButton()}
            {maybeRenderPerformerPopoverButton()}
            {maybeRenderScenePopoverButton()}
            {maybeRenderOrganized()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <GridCard
      className={`gallery-card zoom-${props.zoomIndex}`}
      url={`/galleries/${props.gallery.id}`}
      title={
        props.gallery.title
          ? props.gallery.title
          : TextUtils.fileNameFromPath(props.gallery.path ?? "")
      }
      linkClassName="gallery-card-header"
      image={
        <>
          {props.gallery.cover ? (
            <img
              className="gallery-card-image"
              alt={props.gallery.title ?? ""}
              src={`${props.gallery.cover.paths.thumbnail}`}
            />
          ) : undefined}
          <RatingBanner rating={props.gallery.rating} />
        </>
      }
      overlays={maybeRenderSceneStudioOverlay()}
      details={
        <>
          <span>{props.gallery.date}</span>
          <p>
            <TruncatedText text={props.gallery.details} lineCount={3} />
          </p>
        </>
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
