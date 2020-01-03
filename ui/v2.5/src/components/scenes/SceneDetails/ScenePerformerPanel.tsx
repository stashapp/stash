import React, { FunctionComponent } from "react";
import * as GQL from "../../../core/generated-graphql";
import { PerformerCard } from "../../performers/PerformerCard";

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
