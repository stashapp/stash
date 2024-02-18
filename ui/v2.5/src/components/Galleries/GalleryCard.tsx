import { Button, ButtonGroup, OverlayTrigger, Tooltip } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { GridCard, calculateCardWidth } from "../Shared/GridCard/GridCard";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { SceneLink, TagLink } from "../Shared/TagLink";
import { TruncatedText } from "../Shared/TruncatedText";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import NavUtils from "src/utils/navigation";
import { RatingBanner } from "../Shared/RatingBanner";
import { faBox, faPlayCircle, faTag } from "@fortawesome/free-solid-svg-icons";
import { galleryTitle } from "src/core/galleries";
import ScreenUtils from "src/utils/screen";
import { StudioOverlay } from "../Shared/GridCard/StudioOverlay";

interface IProps {
  gallery: GQL.SlimGalleryDataFragment;
  containerWidth?: number;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const GalleryCard: React.FC<IProps> = (props) => {
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (
      !props.containerWidth ||
      props.zoomIndex === undefined ||
      ScreenUtils.isMobile()
    )
      return;

    let zoomValue = props.zoomIndex;
    let preferredCardWidth: number;
    switch (zoomValue) {
      case 0:
        preferredCardWidth = 240;
        break;
      case 1:
        preferredCardWidth = 340;
        break;
      case 2:
        preferredCardWidth = 480;
        break;
      case 3:
        preferredCardWidth = 640;
    }
    let fittedCardWidth = calculateCardWidth(
      props.containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [props, props.containerWidth, props.zoomIndex]);

  function maybeRenderScenePopoverButton() {
    if (props.gallery.scenes.length === 0) return;

    const popoverContent = props.gallery.scenes.map((scene) => (
      <SceneLink key={scene.id} scene={scene} />
    ));

    return (
      <HoverPopover
        className="scene-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon={faPlayCircle} />
          <span>{props.gallery.scenes.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.gallery.tags.length <= 0) return;

    const popoverContent = props.gallery.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} linkType="gallery" />
    ));

    return (
      <HoverPopover
        className="tag-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon={faTag} />
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

  function maybeRenderOrganized() {
    if (props.gallery.organized) {
      return (
        <OverlayTrigger
          overlay={<Tooltip id="organised-tooltip">{"Organized"}</Tooltip>}
          placement="bottom"
        >
          <div className="organized">
            <Button className="minimal">
              <Icon icon={faBox} />
            </Button>
          </div>
        </OverlayTrigger>
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
      width={cardWidth}
      title={galleryTitle(props.gallery)}
      linkClassName="gallery-card-header"
      image={
        <>
          {props.gallery.cover ? (
            <img
              loading="lazy"
              className="gallery-card-image"
              alt={props.gallery.title ?? ""}
              src={`${props.gallery.cover.paths.thumbnail}`}
            />
          ) : undefined}
          <RatingBanner rating={props.gallery.rating100} />
        </>
      }
      overlays={<StudioOverlay studio={props.gallery.studio} />}
      details={
        <div className="gallery-card__details">
          <span className="gallery-card__date">{props.gallery.date}</span>
          <TruncatedText
            className="gallery-card__description"
            text={props.gallery.details}
            lineCount={3}
          />
        </div>
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
