import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import {
  GridCard,
  calculateCardWidth,
} from "src/components/Shared/GridCard/GridCard";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { TagLink } from "../Shared/TagLink";
import { Button, ButtonGroup } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import { RatingBanner } from "../Shared/RatingBanner";
import ScreenUtils from "src/utils/screen";
import { FavoriteIcon } from "../Shared/FavoriteIcon";
import { useStudioUpdate } from "src/core/StashService";
import { faTag } from "@fortawesome/free-solid-svg-icons";

interface IProps {
  studio: GQL.StudioDataFragment;
  containerWidth?: number;
  hideParent?: boolean;
  selecting?: boolean;
  selected?: boolean;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

function maybeRenderParent(
  studio: GQL.StudioDataFragment,
  hideParent?: boolean
) {
  if (!hideParent && studio.parent_studio) {
    return (
      <div className="studio-parent-studios">
        <FormattedMessage
          id="part_of"
          values={{
            parent: (
              <Link to={`/studios/${studio.parent_studio.id}`}>
                {studio.parent_studio.name}
              </Link>
            ),
          }}
        />
      </div>
    );
  }
}

function maybeRenderChildren(studio: GQL.StudioDataFragment) {
  if (studio.child_studios.length > 0) {
    return (
      <div className="studio-child-studios">
        <FormattedMessage
          id="parent_of"
          values={{
            children: (
              <Link to={NavUtils.makeChildStudiosUrl(studio)}>
                {studio.child_studios.length}&nbsp;
                <FormattedMessage
                  id="countables.studios"
                  values={{ count: studio.child_studios.length }}
                />
              </Link>
            ),
          }}
        />
      </div>
    );
  }
}

export const StudioCard: React.FC<IProps> = (props: IProps) => {
  const [updateStudio] = useStudioUpdate();
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (
      !props.containerWidth ||
      props.zoomIndex === undefined ||
      ScreenUtils.isMobile()
    )
      return;

    let zoomValue = props.zoomIndex;
    console.log(zoomValue);
    let preferredCardWidth: number;
    switch (zoomValue) {
      case 0:
        preferredCardWidth = 280;
        break;
      case 1:
        preferredCardWidth = 340;
        break;
      case 2:
        preferredCardWidth = 420;
        break;
      case 3:
        preferredCardWidth = 560;
    }
    let fittedCardWidth = calculateCardWidth(
      props.containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [props.containerWidth, props.zoomIndex]);

  function onToggleFavorite(v: boolean) {
    if (props.studio.id) {
      updateStudio({
        variables: {
          input: {
            id: props.studio.id,
            favorite: v,
          },
        },
      });
    }
  }

  function maybeRenderScenesPopoverButton() {
    if (!props.studio.scene_count) return;

    return (
      <PopoverCountButton
        className="scene-count"
        type="scene"
        count={props.studio.scene_count}
        url={NavUtils.makeStudioScenesUrl(props.studio)}
      />
    );
  }

  function maybeRenderImagesPopoverButton() {
    if (!props.studio.image_count) return;

    return (
      <PopoverCountButton
        className="image-count"
        type="image"
        count={props.studio.image_count}
        url={NavUtils.makeStudioImagesUrl(props.studio)}
      />
    );
  }

  function maybeRenderGalleriesPopoverButton() {
    if (!props.studio.gallery_count) return;

    return (
      <PopoverCountButton
        className="gallery-count"
        type="gallery"
        count={props.studio.gallery_count}
        url={NavUtils.makeStudioGalleriesUrl(props.studio)}
      />
    );
  }

  function maybeRenderGroupsPopoverButton() {
    if (!props.studio.group_count) return;

    return (
      <PopoverCountButton
        className="group-count"
        type="group"
        count={props.studio.group_count}
        url={NavUtils.makeStudioGroupsUrl(props.studio)}
      />
    );
  }

  function maybeRenderPerformersPopoverButton() {
    if (!props.studio.performer_count) return;

    return (
      <PopoverCountButton
        className="performer-count"
        type="performer"
        count={props.studio.performer_count}
        url={NavUtils.makeStudioPerformersUrl(props.studio)}
      />
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.studio.tags.length <= 0) return;

    const popoverContent = props.studio.tags.map((tag) => (
      <TagLink key={tag.id} linkType="studio" tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal tag-count">
          <Icon icon={faTag} />
          <span>{props.studio.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      props.studio.scene_count ||
      props.studio.image_count ||
      props.studio.gallery_count ||
      props.studio.group_count ||
      props.studio.performer_count ||
      props.studio.tags.length > 0
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderGroupsPopoverButton()}
            {maybeRenderImagesPopoverButton()}
            {maybeRenderGalleriesPopoverButton()}
            {maybeRenderPerformersPopoverButton()}
            {maybeRenderTagPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <GridCard
      className={`studio-card zoom-${props.zoomIndex}`}
      url={`/studios/${props.studio.id}`}
      width={cardWidth}
      title={props.studio.name}
      linkClassName="studio-card-header"
      image={
        <img
          loading="lazy"
          className="studio-card-image"
          alt={props.studio.name}
          src={props.studio.image_path ?? ""}
        />
      }
      details={
        <div className="studio-card__details">
          {maybeRenderParent(props.studio, props.hideParent)}
          {maybeRenderChildren(props.studio)}
          <RatingBanner rating={props.studio.rating100} />
        </div>
      }
      overlays={
        <FavoriteIcon
          favorite={props.studio.favorite}
          onToggleFavorite={(v) => onToggleFavorite(v)}
          size="2x"
          className="hide-not-favorite"
        />
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
