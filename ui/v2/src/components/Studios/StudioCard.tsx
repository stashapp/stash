import {
  Card,
  Elevation,
  H4,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";

interface IProps {
  studio: GQL.StudioDataFragment;
}

export const StudioCard: FunctionComponent<IProps> = (props: IProps) => {
  return (
    <Card
      className="grid-item"
      elevation={Elevation.ONE}
    >
      <Link
        to={`/studios/${props.studio.id}`}
        className="studio previewable image"
        style={{backgroundImage: `url(${props.studio.image_path})`}}
      />
      <div className="card-section">
        <H4 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
          {props.studio.name}
        </H4>
        <span className="bp3-text-muted block">{props.studio.scene_count} scenes.</span>
      </div>
    </Card>
  );
};
