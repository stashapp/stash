import { ButtonGroup } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import { FormattedMessage } from "react-intl";
import { TruncatedText } from "../Shared/TruncatedText";
import { GridCard, calculateCardWidth } from "../Shared/GridCard";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import ScreenUtils from "src/utils/screen";

interface IProps {
  tag: GQL.TagDataFragment;
  containerWidth?: number;
  zoomIndex: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const TagCard: React.FC<IProps> = ({
  tag,
  containerWidth,
  zoomIndex,
  selecting,
  selected,
  onSelectedChanged,
}) => {
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (!containerWidth || zoomIndex === undefined || ScreenUtils.isMobile())
      return;

    let zoomValue = zoomIndex;
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
      containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [containerWidth, zoomIndex]);

  function maybeRenderDescription() {
    if (tag.description) {
      return (
        <TruncatedText
          className="tag-description"
          text={tag.description}
          lineCount={3}
        />
      );
    }
  }

  function maybeRenderParents() {
    if (tag.parents.length === 1) {
      const parent = tag.parents[0];
      return (
        <div className="tag-parent-tags">
          <FormattedMessage
            id="sub_tag_of"
            values={{
              parent: <Link to={`/tags/${parent.id}`}>{parent.name}</Link>,
            }}
          />
        </div>
      );
    }

    if (tag.parents.length > 1) {
      return (
        <div className="tag-parent-tags">
          <FormattedMessage
            id="sub_tag_of"
            values={{
              parent: (
                <Link to={NavUtils.makeParentTagsUrl(tag)}>
                  {tag.parents.length}&nbsp;
                  <FormattedMessage
                    id="countables.tags"
                    values={{ count: tag.parents.length }}
                  />
                </Link>
              ),
            }}
          />
        </div>
      );
    }
  }

  function maybeRenderChildren() {
    if (tag.children.length > 0) {
      return (
        <div className="tag-sub-tags">
          <FormattedMessage
            id="parent_of"
            values={{
              children: (
                <Link to={NavUtils.makeChildTagsUrl(tag)}>
                  {tag.children.length}&nbsp;
                  <FormattedMessage
                    id="countables.tags"
                    values={{ count: tag.children.length }}
                  />
                </Link>
              ),
            }}
          />
        </div>
      );
    }
  }

  function maybeRenderScenesPopoverButton() {
    if (!tag.scene_count) return;

    return (
      <PopoverCountButton
        className="scene-count"
        type="scene"
        count={tag.scene_count}
        url={NavUtils.makeTagScenesUrl(tag)}
      />
    );
  }

  function maybeRenderSceneMarkersPopoverButton() {
    if (!tag.scene_marker_count) return;

    return (
      <PopoverCountButton
        className="marker-count"
        type="marker"
        count={tag.scene_marker_count}
        url={NavUtils.makeTagSceneMarkersUrl(tag)}
      />
    );
  }

  function maybeRenderImagesPopoverButton() {
    if (!tag.image_count) return;

    return (
      <PopoverCountButton
        className="image-count"
        type="image"
        count={tag.image_count}
        url={NavUtils.makeTagImagesUrl(tag)}
      />
    );
  }

  function maybeRenderGalleriesPopoverButton() {
    if (!tag.gallery_count) return;

    return (
      <PopoverCountButton
        className="gallery-count"
        type="gallery"
        count={tag.gallery_count}
        url={NavUtils.makeTagGalleriesUrl(tag)}
      />
    );
  }

  function maybeRenderPerformersPopoverButton() {
    if (!tag.performer_count) return;

    return (
      <PopoverCountButton
        className="performer-count"
        type="performer"
        count={tag.performer_count}
        url={NavUtils.makeTagPerformersUrl(tag)}
      />
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (tag) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderImagesPopoverButton()}
            {maybeRenderGalleriesPopoverButton()}
            {maybeRenderSceneMarkersPopoverButton()}
            {maybeRenderPerformersPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <GridCard
      className={`tag-card zoom-${zoomIndex}`}
      url={`/tags/${tag.id}`}
      width={cardWidth}
      title={tag.name ?? ""}
      linkClassName="tag-card-header"
      image={
        <img
          loading="lazy"
          className="tag-card-image"
          alt={tag.name}
          src={tag.image_path ?? ""}
        />
      }
      details={
        <>
          {maybeRenderDescription()}
          {maybeRenderParents()}
          {maybeRenderChildren()}
        </>
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
};
