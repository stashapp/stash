import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupCard } from "src/components/Movies/MovieCard";

interface ISceneMoviePanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneGroupPanel: React.FC<ISceneMoviePanelProps> = (
  props: ISceneMoviePanelProps
) => {
  const cards = props.scene.movies.map((sceneMovie) => (
    <GroupCard
      key={sceneMovie.movie.id}
      group={sceneMovie.movie}
      sceneIndex={sceneMovie.scene_index ?? undefined}
    />
  ));

  return (
    <>
      <div className="row justify-content-center">{cards}</div>
    </>
  );
};

export default SceneGroupPanel;
