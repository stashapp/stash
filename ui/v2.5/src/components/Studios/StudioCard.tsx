import { Card } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";

interface IProps {
  studio: GQL.StudioDataFragment;
}

export const StudioCard: React.FC<IProps> = ({ studio }) => {
  return (
    <Card className="studio-card">
      <Link to={`/studios/${studio.id}`} className="studio-image">
        <img alt={studio.name} src={studio.image_path ?? ""} />
      </Link>
      <div className="card-section">
        <h5 className="text-truncate">{studio.name}</h5>
        <span>{studio.scene_count} scenes.</span>
      </div>
    </Card>
  );
};
