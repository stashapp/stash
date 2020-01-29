import React, { FunctionComponent } from "react";
import * as GQL from "../../../core/generated-graphql";
import { MovieCard } from "../../Movies/MovieCard";

interface ISceneMoviePanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneMoviePanel: FunctionComponent<ISceneMoviePanelProps> = (props: ISceneMoviePanelProps) => {
  const cards = props.scene.movies.map((movie) => (
    <MovieCard key={movie.id} movie={movie} fromscene={true} />
    
  ));

  return (
    <>
      <div className="grid">
        {cards}
      </div>
    </>
  );
};
