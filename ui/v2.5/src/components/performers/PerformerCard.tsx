import React from "react";
import { Card } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { NavUtils, TextUtils } from "src/utils";

interface IPerformerCardProps {
  performer: GQL.PerformerDataFragment;
  ageFromDate?: string;
}

export const PerformerCard: React.FC<IPerformerCardProps> = (
  props: IPerformerCardProps
) => {
  const age = TextUtils.age(props.performer.birthdate, props.ageFromDate);
  const ageString = `${age} years old${
    props.ageFromDate ? " in this scene." : "."
  }`;

  function maybeRenderFavoriteBanner() {
    if (props.performer.favorite === false) {
      return;
    }
    return <div className="rating-banner rating-5">FAVORITE</div>;
  }

  return (
    <Card className="grid-item">
      <Link
        to={`/performers/${props.performer.id}`}
        className="performer previewable image"
        style={{ backgroundImage: `url(${props.performer.image_path})` }}
      >
        {maybeRenderFavoriteBanner()}
      </Link>
      <div className="card-section">
        <h5 className="text-truncate">{props.performer.name}</h5>
        {age !== 0 ? <div className="text-muted">{ageString}</div> : ""}
        <div className="text-muted">
          Stars in {props.performer.scene_count}{" "}
          <Link to={NavUtils.makePerformerScenesUrl(props.performer)}>
            scenes
          </Link>
          .
        </div>
      </div>
    </Card>
  );
};
