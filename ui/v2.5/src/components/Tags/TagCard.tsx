import { PatchComponent } from "src/patch";
import { Button, ButtonGroup } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import { FormattedMessage } from "react-intl";
import { TruncatedText } from "../Shared/TruncatedText";
import { GridCard } from "../Shared/GridCard/GridCard";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import { Icon } from "../Shared/Icon";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";
import { useTagUpdate } from "src/core/StashService";

interface IProps {
  tag: GQL.TagDataFragment | GQL.TagListDataFragment;
  cardWidth?: number;
  zoomIndex: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

const TagCardPopovers: React.FC<IProps> = PatchComponent(
  "TagCard.Popovers",
  ({ tag }) => {
    return (
      <>
        <hr />
        <ButtonGroup className="card-popovers">
          <PopoverCountButton
            className="scene-count"
            type="scene"
            count={tag.scene_count}
            url={NavUtils.makeTagScenesUrl(tag)}
            showZero={false}
          />
          <PopoverCountButton
            className="image-count"
            type="image"
            count={tag.image_count}
            url={NavUtils.makeTagImagesUrl(tag)}
            showZero={false}
          />
          <PopoverCountButton
            className="gallery-count"
            type="gallery"
            count={tag.gallery_count}
            url={NavUtils.makeTagGalleriesUrl(tag)}
            showZero={false}
          />
          <PopoverCountButton
            className="group-count"
            type="group"
            count={tag.group_count}
            url={NavUtils.makeTagGroupsUrl(tag)}
            showZero={false}
          />
          <PopoverCountButton
            className="marker-count"
            type="marker"
            count={tag.scene_marker_count}
            url={NavUtils.makeTagSceneMarkersUrl(tag)}
            showZero={false}
          />
          <PopoverCountButton
            className="performer-count"
            type="performer"
            count={tag.performer_count}
            url={NavUtils.makeTagPerformersUrl(tag)}
            showZero={false}
          />
          <PopoverCountButton
            className="studio-count"
            type="studio"
            count={tag.studio_count}
            url={NavUtils.makeTagStudiosUrl(tag)}
            showZero={false}
          />
        </ButtonGroup>
      </>
    );
  }
);

const TagCardOverlays: React.FC<IProps> = PatchComponent(
  "TagCard.Overlays",
  ({ tag }) => {
    const [updateTag] = useTagUpdate();

    function renderFavoriteIcon() {
      return (
        <Link to="" onClick={(e) => e.preventDefault()}>
          <Button
            className={cx(
              "minimal",
              "mousetrap",
              "favorite-button",
              tag.favorite ? "favorite" : "not-favorite"
            )}
            onClick={() => onToggleFavorite!(!tag.favorite)}
          >
            <Icon icon={faHeart} size="2x" />
          </Button>
        </Link>
      );
    }

    function onToggleFavorite(v: boolean) {
      if (tag.id) {
        updateTag({
          variables: {
            input: {
              id: tag.id,
              favorite: v,
            },
          },
        });
      }
    }

    return <>{renderFavoriteIcon()}</>;
  }
);

const TagCardDetails: React.FC<IProps> = PatchComponent(
  "TagCard.Details",
  ({ tag }) => {
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

    return (
      <>
        {maybeRenderDescription()}
        {maybeRenderParents()}
        {maybeRenderChildren()}
      </>
    );
  }
);

const TagCardImage: React.FC<IProps> = PatchComponent(
  "TagCard.Image",
  ({ tag }) => {
    return (
      <>
        <img
          loading="lazy"
          className="tag-card-image"
          alt={tag.name}
          src={tag.image_path ?? ""}
        />
      </>
    );
  }
);

const TagCardTitle: React.FC<IProps> = PatchComponent(
  "TagCard.Title",
  ({ tag }) => {
    return <>{tag.name ?? ""}</>;
  }
);

export const TagCard: React.FC<IProps> = PatchComponent("TagCard", (props) => {
  const { tag, cardWidth, zoomIndex, selecting, selected, onSelectedChanged } =
    props;

  return (
    <GridCard
      className={`tag-card zoom-${zoomIndex}`}
      url={`/tags/${tag.id}`}
      width={cardWidth}
      title={<TagCardTitle {...props} />}
      linkClassName="tag-card-header"
      image={<TagCardImage {...props} />}
      details={<TagCardDetails {...props} />}
      overlays={<TagCardOverlays {...props} />}
      popovers={<TagCardPopovers {...props} />}
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
});
