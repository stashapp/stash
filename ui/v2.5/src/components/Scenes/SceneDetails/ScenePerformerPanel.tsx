import React, { FunctionComponent } from "react";
import * as GQL from "src/core/generated-graphql";
import { PerformerCard } from "src/components/Performers/PerformerCard";

interface IScenePerformerPanelProps {
  scene: GQL.SceneDataFragment;
}

export const ScenePerformerPanel: FunctionComponent<IScenePerformerPanelProps> = (
  props: IScenePerformerPanelProps
) => {
  const cards = props.scene.performers.map((performer) => (
    <PerformerCard
      key={performer.id}
      performer={performer}
      ageFromDate={props.scene.date ?? undefined}
    />
  ));

  return (
    <>
      <div className="row justify-content-center">{cards}</div>
    </>
  );
};
