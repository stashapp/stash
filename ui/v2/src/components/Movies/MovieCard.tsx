import {
  Card,
  Elevation,
  H4,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { ColorUtils } from "../../utils/color";

interface IProps {
  movie: GQL.MovieDataFragment;
  fromscene: boolean;
 // scene: GQL.SceneDataFragment;
}


export const MovieCard: FunctionComponent<IProps> = (props: IProps) => {
  function maybeRenderRatingBanner() {
    if (!props.movie.rating) { return; }
    return (
      <div className={`rating-banner ${ColorUtils.classForRating(parseInt(props.movie.rating,10))}`}>
        RATING: {props.movie.rating}
      </div>
    );
  }
  
  function maybeRenderSceneNumber() {
    if (!props.fromscene) {
    return (
      <div className="card-section">
         <H4 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
          {props.movie.name}
         </H4>
         <span className="bp3-text-muted block">{props.movie.scene_count} scenes.</span>
      </div>
    );  
    } else {
    return (
      <div className="card-section">
        <H4 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
        {props.movie.name}
        </H4>
        <span className="bp3-text-muted block">Scene number: {props.movie.scene_index}</span>   
      </div>
    );  
      }
  }
  
  return (
    <Card
      className="grid-item"
      elevation={Elevation.ONE}
    >
      
      <Link
        to={`/movies/${props.movie.id}`}
        className="movie previewable image"
        style={{backgroundImage: `url(${props.movie.front_image_path})`}}
      >
      {maybeRenderRatingBanner()}
      </Link>
      {maybeRenderSceneNumber()}

    </Card>
  );
};
