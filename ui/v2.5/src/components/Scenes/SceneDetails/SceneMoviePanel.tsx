import React, { FunctionComponent } from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieCard } from "src/components/Movies/MovieCard";

interface ISceneMoviePanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneMoviePanel: FunctionComponent<ISceneMoviePanelProps> = (
  props: ISceneMoviePanelProps
) => {
  const cards = props.scene.movies.map((sceneMovie) => (
    <MovieCard
      key={sceneMovie.movie.id}
      movie={sceneMovie.movie}
      sceneIndex={sceneMovie.scene_index ?? undefined}
    />
  ));

  return (
    <>
      <div className="row justify-content-center">{cards}</div>
    </>
  );
};

export default SceneMoviePanel;
