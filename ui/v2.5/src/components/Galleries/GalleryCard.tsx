import { Button, ButtonGroup, OverlayTrigger, Tooltip } from "react-bootstrap";
import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { GridCard } from "../Shared/GridCard/GridCard";
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
import { StudioOverlay } from "../Shared/GridCard/StudioOverlay";
import { GalleryPreviewScrubber } from "./GalleryPreviewScrubber";
import cx from "classnames";
import { useHistory } from "react-router-dom";
import { PatchComponent } from "src/patch";

interface IGalleryPreviewProps {
  gallery: GQL.SlimGalleryDataFragment;
  onScrubberClick?: (index: number) => void;
}

export const GalleryPreview: React.FC<IGalleryPreviewProps> = ({
  gallery,
  onScrubberClick,
}) => {
  const [imgSrc, setImgSrc] = useState<string | undefined>(
    gallery.paths.cover ?? undefined
  );

  return (
    <div className={cx("gallery-card-cover")}>
      {!!imgSrc && (
        <img
          loading="lazy"
          className="gallery-card-image"
          alt={gallery.title ?? ""}
          src={imgSrc}
        />
      )}
      {gallery.image_count > 0 && (
        <GalleryPreviewScrubber
          previewPath={gallery.paths.preview}
          defaultPath={gallery.paths.cover ?? ""}
          imageCount={gallery.image_count}
          onClick={onScrubberClick}
          onPathChanged={setImgSrc}
        />
      )}
    </div>
  );
};

interface IGalleryCardProps {
  gallery: GQL.SlimGalleryDataFragment;
  cardWidth?: number;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

const GalleryCardPopovers = PatchComponent(
  "GalleryCard.Popovers",
  (props: IGalleryCardProps) => {
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

      return (
        <PerformerPopoverButton
          performers={props.gallery.performers}
          linkType="gallery"
        />
      );
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

    return <>{maybeRenderPopoverButtonGroup()}</>;
  }
);

const GalleryCardDetails = PatchComponent(
  "GalleryCard.Details",
  (props: IGalleryCardProps) => {
    return (
      <div className="gallery-card__details">
        <span className="gallery-card__date">{props.gallery.date}</span>
        <TruncatedText
          className="gallery-card__description"
          text={props.gallery.details}
          lineCount={3}
        />
      </div>
    );
  }
);

const GalleryCardOverlays = PatchComponent(
  "GalleryCard.Overlays",
  (props: IGalleryCardProps) => {
    return <StudioOverlay studio={props.gallery.studio} />;
  }
);

const GalleryCardImage = PatchComponent(
  "GalleryCard.Image",
  (props: IGalleryCardProps) => {
    const history = useHistory();

    return (
      <>
        <GalleryPreview
          gallery={props.gallery}
          onScrubberClick={(i) => {
            history.push(`/galleries/${props.gallery.id}/images/${i}`);
          }}
        />
        <RatingBanner rating={props.gallery.rating100} />
      </>
    );
  }
);

export const GalleryCard = PatchComponent(
  "GalleryCard",
  (props: IGalleryCardProps) => {
    return (
      <GridCard
        className={`gallery-card zoom-${props.zoomIndex}`}
        url={`/galleries/${props.gallery.id}`}
        width={props.cardWidth}
        title={galleryTitle(props.gallery)}
        linkClassName="gallery-card-header"
        image={<GalleryCardImage {...props} />}
        overlays={<GalleryCardOverlays {...props} />}
        details={<GalleryCardDetails {...props} />}
        popovers={<GalleryCardPopovers {...props} />}
        selected={props.selected}
        selecting={props.selecting}
        onSelectedChanged={props.onSelectedChanged}
      />
    );
  }
);
