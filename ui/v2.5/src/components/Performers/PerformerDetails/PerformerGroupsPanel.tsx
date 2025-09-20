import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupList } from "src/components/Groups/MyGroupList";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";
import { PatchComponent } from "src/patch";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerGroupsPanel: React.FC<IPerformerDetailsProps> =
  PatchComponent("PerformerGroupsPanel", ({ active, performer }) => {
    const filterHook = usePerformerFilterHook(performer);
    return (
      <GroupList
        filterHook={filterHook}
        alterQuery={active}
        view={View.PerformerGroups}
      />
    );
  });
