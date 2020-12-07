import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { FormattedPlural } from "react-intl";
import { NavUtils } from "src/utils";
import { BasicCard, TruncatedText } from "src/components/Shared";

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
          <span>
            {studio.scene_count}&nbsp;
            <FormattedPlural
              value={studio.scene_count ?? 0}
              one="scene"
              other="scenes"
            />
            .
          </span>
          {maybeRenderParent(studio, hideParent)}
          {maybeRenderChildren(studio)}
        </>
      }
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
};
