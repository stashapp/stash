import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneList } from "src/components/Scenes/SceneList";
import { useStudioFilterHook } from "src/core/studios";
import { View } from "src/components/List/views";

interface IStudioScenesPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
  showChildStudioContent?: boolean;
}

export const StudioScenesPanel: React.FC<IStudioScenesPanel> = ({
  active,
  studio,
  showChildStudioContent,
}) => {
  const filterHook = useStudioFilterHook(studio, showChildStudioContent);
  return (
    <SceneList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioScenes}
    />
  );
};
