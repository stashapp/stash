import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupList } from "src/components/Groups/GroupList";
import { useStudioFilterHook } from "src/core/studios";
import { View } from "src/components/List/views";

interface IStudioGroupsPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
  showChildStudioContent?: boolean;
}

export const StudioGroupsPanel: React.FC<IStudioGroupsPanel> = ({
  active,
  studio,
  showChildStudioContent,
}) => {
  const filterHook = useStudioFilterHook(studio, showChildStudioContent);
  return (
    <GroupList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioGroups}
    />
  );
};
