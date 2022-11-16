import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneList } from "src/components/Scenes/SceneList";
import { usePerformerFilterHook } from "src/core/performers";

interface IPerformerDetailsProps {
  performer: GQL.PerformerDataFragment;
}

export const PerformerScenesPanel: React.FC<IPerformerDetailsProps> = ({
  performer,
}) => {
  const filterHook = usePerformerFilterHook(performer);
  return <SceneList filterHook={filterHook} />;
};
