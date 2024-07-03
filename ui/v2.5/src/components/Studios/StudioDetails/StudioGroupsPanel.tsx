import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupList } from "src/components/Groups/GroupList";
import { useStudioFilterHook } from "src/core/studios";
import { View } from "src/components/List/views";

interface IStudioGroupsPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioGroupsPanel: React.FC<IStudioGroupsPanel> = ({
  active,
  studio,
}) => {
  const filterHook = useStudioFilterHook(studio);
  return (
    <GroupList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioGroups}
    />
  );
};
