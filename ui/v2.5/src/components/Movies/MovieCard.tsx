import { Card } from "react-bootstrap";
import React, { FunctionComponent } from "react";
import { FormattedPlural } from "react-intl";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";

interface IProps {
  movie: GQL.MovieDataFragment;
  sceneIndex?: number;
}

export const MovieCard: FunctionComponent<IProps> = (props: IProps) => {
  function maybeRenderRatingBanner() {
    if (!props.movie.rating) {
      return;
    }
    return (
      <div
        className={`rating-banner ${
          props.movie.rating ? `rating-${props.movie.rating}` : ""
        }`}
      >
        RATING: {props.movie.rating}
      </div>
    );
  }

  function maybeRenderSceneNumber() {
    if (!props.sceneIndex) {
      return (
        <span>
          {props.movie.scene_count}&nbsp;
          <FormattedPlural
            value={props.movie.scene_count ?? 0}
            one="scene"
            other="scenes"
          />
        </span>
      );
    }

    return <span>Scene number: {props.sceneIndex}</span>;
  }

  return (
    <Card className="movie-card">
      <Link to={`/movies/${props.movie.id}`} className="movie-card-header">
        <img
          className="movie-card-image"
          alt={props.movie.name ?? ""}
          src={props.movie.front_image_path ?? ""}
        />
        {maybeRenderRatingBanner()}
      </Link>
      <div className="card-section">
        <h5 className="text-truncate">{props.movie.name}</h5>
        {maybeRenderSceneNumber()}
      </div>
    </Card>
  );
};
