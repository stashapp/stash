import React, { MouseEvent, useEffect, useMemo, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "src/components/Shared/Icon";
import { GalleryLink, TagLink } from "src/components/Shared/TagLink";
import { HoverPopover } from "src/components/Shared/HoverPopover";
import { SweatDrops } from "src/components/Shared/SweatDrops";
import { PerformerPopoverButton } from "src/components/Shared/PerformerPopoverButton";
import {
  GridCard,
  calculateCardWidth,
} from "src/components/Shared/GridCard/GridCard";
import { RatingBanner } from "src/components/Shared/RatingBanner";
import {
  faBox,
  faImages,
  faSearch,
  faTag,
} from "@fortawesome/free-solid-svg-icons";
import { objectTitle } from "src/core/files";
import { TruncatedText } from "../Shared/TruncatedText";
import ScreenUtils from "src/utils/screen";
import { StudioOverlay } from "../Shared/GridCard/StudioOverlay";

interface IImageCardProps {
  image: GQL.SlimImageDataFragment;
  containerWidth?: number;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  onPreview?: (ev: MouseEvent) => void;
}

export const ImageCard: React.FC<IImageCardProps> = (
  props: IImageCardProps
) => {
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

  const file = useMemo(
    () =>
      props.image.visual_files.length > 0
        ? props.image.visual_files[0]
        : undefined,
    [props.image]
  );

  function maybeRenderTagPopoverButton() {
    if (props.image.tags.length <= 0) return;

    const popoverContent = props.image.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} linkType="image" />
    ));

    return (
      <HoverPopover
        className="tag-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon={faTag} />
          <span>{props.image.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.image.performers.length <= 0) return;

    return (
      <PerformerPopoverButton
        performers={props.image.performers}
        linkType="image"
      />
    );
  }

  function maybeRenderOCounter() {
    if (props.image.o_counter) {
      return (
        <div className="o-count">
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

  function maybeRenderGallery() {
    if (props.image.galleries.length <= 0) return;

    const popoverContent = props.image.galleries.map((gallery) => (
      <GalleryLink key={gallery.id} gallery={gallery} />
    ));

    return (
      <HoverPopover
        className="gallery-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon={faImages} />
          <span>{props.image.galleries.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderOrganized() {
    if (props.image.organized) {
      return (
        <div className="organized">
          <Button className="minimal">
            <Icon icon={faBox} />
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
      props.image.galleries.length > 0 ||
      props.image.organized
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderTagPopoverButton()}
            {maybeRenderPerformerPopoverButton()}
            {maybeRenderOCounter()}
            {maybeRenderGallery()}
            {maybeRenderOrganized()}
          </ButtonGroup>
        </>
      );
    }
  }

  function isPortrait() {
    const width = file?.width ? file.width : 0;
    const height = file?.height ? file.height : 0;
    return height > width;
  }

  const source =
    props.image.paths.preview != ""
      ? props.image.paths.preview ?? ""
      : props.image.paths.thumbnail ?? "";
  const video = source.includes("preview");
  const ImagePreview = video ? "video" : "img";

  return (
    <GridCard
      className={`image-card zoom-${props.zoomIndex}`}
      url={`/images/${props.image.id}`}
      width={cardWidth}
      title={objectTitle(props.image)}
      linkClassName="image-card-link"
      image={
        <>
          <div className={cx("image-card-preview", { portrait: isPortrait() })}>
            <ImagePreview
              loop={video}
              autoPlay={video}
              className="image-card-preview-image"
              alt={props.image.title ?? ""}
              src={source}
            />
            {props.onPreview ? (
              <div className="preview-button">
                <Button onClick={props.onPreview}>
                  <Icon icon={faSearch} />
                </Button>
              </div>
            ) : undefined}
          </div>
          <RatingBanner rating={props.image.rating100} />
        </>
      }
      details={
        <div className="image-card__details">
          <span className="image-card__date">{props.image.date}</span>
          <TruncatedText
            className="image-card__description"
            text={props.image.details}
            lineCount={3}
          />
        </div>
      }
      overlays={<StudioOverlay studio={props.image.studio} />}
      popovers={maybeRenderPopoverButtonGroup()}
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
