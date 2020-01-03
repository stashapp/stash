import { Card } from 'react-bootstrap';
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";

interface IProps {
  studio: GQL.StudioDataFragment;
}

export const StudioCard: React.FC<IProps> = (props: IProps) => {
  return (
    <Card
      className="col-4"
    >
      <Link
        to={`/studios/${props.studio.id}`}
        className="studio previewable image"
        style={{backgroundImage: `url(${props.studio.image_path})`}}
      />
      <div className="card-section">
        <h4 className="text-truncate">
          {props.studio.name}
        </h4>
        <span>{props.studio.scene_count} scenes.</span>
      </div>
    </Card>
  );
};
