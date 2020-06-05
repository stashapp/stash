import { Card } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { FormattedPlural } from "react-intl";
import { NavUtils } from "src/utils";

interface IProps {
  studio: GQL.StudioDataFragment;
  hideParent?: boolean;
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

export const StudioCard: React.FC<IProps> = ({ studio, hideParent }) => {
  return (
    <Card className="studio-card">
      <Link to={`/studios/${studio.id}`} className="studio-card-header">
        <img
          className="studio-card-image"
          alt={studio.name}
          src={studio.image_path ?? ""}
        />
      </Link>
      <div className="card-section">
        <h5 className="text-truncate">{studio.name}</h5>
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
      </div>
    </Card>
  );
};
