import { Card } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { FormattedPlural } from "react-intl";

interface IProps {
  studio: GQL.StudioDataFragment;
}

export const StudioCard: React.FC<IProps> = ({ studio }) => {
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
      </div>
    </Card>
  );
};
