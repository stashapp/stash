import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupCard } from "src/components/Movies/MovieCard";

interface ISceneGroupPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneGroupPanel: React.FC<ISceneGroupPanelProps> = (
  props: ISceneGroupPanelProps
) => {
  const cards = props.scene.movies.map((sceneGroup) => (
    <GroupCard
      key={sceneGroup.movie.id}
      group={sceneGroup.movie}
      sceneIndex={sceneGroup.scene_index ?? undefined}
    />
  ));

  return (
    <>
      <div className="row justify-content-center">{cards}</div>
    </>
  );
};

export default SceneGroupPanel;
