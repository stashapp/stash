import {
  Card,
  Elevation,
  H4,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";


interface IProps {
  dvd: GQL.DvdDataFragment;
 // scene: GQL.SceneDataFragment;
}


export const DvdCard: FunctionComponent<IProps> = (props: IProps) => {
   return (
    <Card
      className="grid-item"
      elevation={Elevation.ONE}
    >
      <Link
        to={`/dvds/${props.dvd.id}`}
        className="dvd previewable image"
        style={{backgroundImage: `url(${props.dvd.frontimage_path})`}}
      />
      <div className="card-section">
        <H4 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
          {props.dvd.name}
        </H4>
        
        <span className="bp3-text-muted block">{props.dvd.scene_count} scenes.</span>
          
        
      </div>
    </Card>
  );
};
