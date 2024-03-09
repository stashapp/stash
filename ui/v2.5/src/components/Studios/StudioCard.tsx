import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import {
  GridCard,
  calculateCardWidth,
} from "src/components/Shared/GridCard/GridCard";
import { ButtonGroup } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import { RatingBanner } from "../Shared/RatingBanner";
import ScreenUtils from "src/utils/screen";

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
                {studio.child_studios.length} studios
              </Link>
            ),
          }}
        />
      </div>
    );
  }
}

export const StudioCard: React.FC<IProps> = (props: IProps) => {
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

  function maybeRenderMoviesPopoverButton() {
    if (!props.studio.movie_count) return;

    return (
      <PopoverCountButton
        className="movie-count"
        type="movie"
        count={props.studio.movie_count}
        url={NavUtils.makeStudioMoviesUrl(props.studio)}
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

  function maybeRenderPopoverButtonGroup() {
    if (
      props.studio.scene_count ||
      props.studio.image_count ||
      props.studio.gallery_count ||
      props.studio.movie_count ||
      props.studio.performer_count
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderMoviesPopoverButton()}
            {maybeRenderImagesPopoverButton()}
            {maybeRenderGalleriesPopoverButton()}
            {maybeRenderPerformersPopoverButton()}
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
      popovers={maybeRenderPopoverButtonGroup()}
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
