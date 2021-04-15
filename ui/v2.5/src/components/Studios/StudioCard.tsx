import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { NavUtils } from "src/utils";
import { BasicCard, TruncatedText } from "src/components/Shared";
import { ButtonGroup } from "react-bootstrap";
import { PopoverCountButton } from "../Shared/PopoverCountButton";

interface IProps {
  studio: GQL.StudioDataFragment;
  hideParent?: boolean;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

function maybeRenderParent(
  studio: GQL.StudioDataFragment,
  hideParent?: boolean
) {
  if (!hideParent && studio.parent_studio) {
    return (
      <div>
        Part of&nbsp;
        <Link to={`/studios/${studio.parent_studio.id}`}>
          {studio.parent_studio.name}
        </Link>
        .
      </div>
    );
  }
}

function maybeRenderChildren(studio: GQL.StudioDataFragment) {
  if (studio.child_studios.length > 0) {
    return (
      <div>
        Parent of&nbsp;
        <Link to={NavUtils.makeChildStudiosUrl(studio)}>
          {studio.child_studios.length} studios
        </Link>
        .
      </div>
    );
  }
}

export const StudioCard: React.FC<IProps> = ({
  studio,
  hideParent,
  selecting,
  selected,
  onSelectedChanged,
}) => {
  function maybeRenderScenesPopoverButton() {
    if (!studio.scene_count) return;

    return (
      <PopoverCountButton
        type="scene"
        count={studio.scene_count}
        url={NavUtils.makeStudioScenesUrl(studio)}
      />
    );
  }

  function maybeRenderImagesPopoverButton() {
    if (!studio.image_count) return;

    return (
      <PopoverCountButton
        type="image"
        count={studio.image_count}
        url={NavUtils.makeStudioImagesUrl(studio)}
      />
    );
  }

  function maybeRenderGalleriesPopoverButton() {
    if (!studio.gallery_count) return;

    return (
      <PopoverCountButton
        type="gallery"
        count={studio.gallery_count}
        url={NavUtils.makeStudioGalleriesUrl(studio)}
      />
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (studio.scene_count || studio.image_count || studio.gallery_count) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderImagesPopoverButton()}
            {maybeRenderGalleriesPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <BasicCard
      className="studio-card"
      url={`/studios/${studio.id}`}
      linkClassName="studio-card-header"
      image={
        <img
          className="studio-card-image"
          alt={studio.name}
          src={studio.image_path ?? ""}
        />
      }
      details={
        <>
          <h5>
            <TruncatedText text={studio.name} />
          </h5>
          {maybeRenderParent(studio, hideParent)}
          {maybeRenderChildren(studio)}
          {maybeRenderPopoverButtonGroup()}
        </>
      }
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
};
