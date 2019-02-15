import {
  Card,
  Elevation,
  H4,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { TextUtils } from "../../utils/text";

interface IPerformerCardProps {
  performer: GQL.PerformerDataFragment;
  ageFromDate?: string;
}

export const PerformerCard: FunctionComponent<IPerformerCardProps> = (props: IPerformerCardProps) => {
  const age = TextUtils.age(props.performer.birthdate, props.ageFromDate);
  const ageString = `${age} years old${!!props.ageFromDate ? " in this scene." : "."}`;

  function maybeRenderFavoriteBanner() {
    if (props.performer.favorite === false) { return; }
    return (
      <div className={`rating-banner rating-5`}>
        FAVORITE
      </div>
    );
  }

  return (
    <Card
      className="grid-item"
      elevation={Elevation.ONE}
    >
      <Link
        to={`/performers/${props.performer.id}`}
        className="performer previewable image"
        style={{backgroundImage: `url(${props.performer.image_path})`}}
      >
        {maybeRenderFavoriteBanner()}
      </Link>
      <div className="card-section">
        <H4 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
          {props.performer.name}
        </H4>
        {age !== 0 ? <span className="bp3-text-muted block">{ageString}</span> : undefined}
        <span className="bp3-text-muted block">Stars in {props.performer.scene_count} scenes.</span>
      </div>
    </Card>
  );
};
