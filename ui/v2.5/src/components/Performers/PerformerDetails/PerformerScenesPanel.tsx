import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MyFilteredSceneList } from "src/components/Scenes/MySceneList";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";
import { PatchComponent } from "src/patch";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerScenesPanel: React.FC<IPerformerDetailsProps> =
  PatchComponent("PerformerScenesPanel", ({ active, performer }) => {
    const filterHook = usePerformerFilterHook(performer);
    return (
      <MyFilteredSceneList
        filterHook={filterHook}
        alterQuery={active}
        view={View.PerformerScenes}
      />
    );
  });
