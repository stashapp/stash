import React from "react";
import * as GQL from "src/core/generated-graphql";
import { FilteredAudioList } from "src/components/Audios/AudioList";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";
import { PatchComponent } from "src/patch";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerAudiosPanel: React.FC<IPerformerDetailsProps> =
  PatchComponent("PerformerAudiosPanel", ({ active, performer }) => {
    const filterHook = usePerformerFilterHook(performer);
    return (
      <FilteredAudioList
        filterHook={filterHook}
        alterQuery={active}
        view={View.PerformerAudios}
      />
    );
  });
