import React, { FunctionComponent } from "react";
import * as GQL from "src/core/generated-graphql";
import { PerformerCard } from "src/components/performers/PerformerCard";

interface IScenePerformerPanelProps {
  scene: GQL.SceneDataFragment;
}

export const ScenePerformerPanel: FunctionComponent<IScenePerformerPanelProps> = (props: IScenePerformerPanelProps) => {
  const cards = props.scene.performers.map((performer) => (
    <PerformerCard key={performer.id} performer={performer} ageFromDate={props.scene.date} />
  ));

  return (
    <>
      <div className="grid">
        {cards}
      </div>
    </>
  );
};
