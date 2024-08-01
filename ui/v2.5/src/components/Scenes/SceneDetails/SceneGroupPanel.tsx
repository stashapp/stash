import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupCard } from "src/components/Groups/GroupCard";

interface ISceneGroupPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneGroupPanel: React.FC<ISceneGroupPanelProps> = (
  props: ISceneGroupPanelProps
) => {
  const cards = props.scene.groups.map((sceneGroup) => (
    <GroupCard
      key={sceneGroup.group.id}
      group={sceneGroup.group}
      sceneNumber={sceneGroup.scene_index ?? undefined}
    />
  ));

  return (
    <>
      <div className="row justify-content-center">{cards}</div>
    </>
  );
};

export default SceneGroupPanel;
