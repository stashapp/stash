import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneList } from "src/components/Scenes/SceneList";
import { useStudioFilterHook } from "src/core/studios";

interface IStudioScenesPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioScenesPanel: React.FC<IStudioScenesPanel> = ({
  active,
  studio,
}) => {
  const filterHook = useStudioFilterHook(studio);
  return <SceneList filterHook={filterHook} alterQuery={active} />;
};
