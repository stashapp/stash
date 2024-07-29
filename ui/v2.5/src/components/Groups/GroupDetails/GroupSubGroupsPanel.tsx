import React from "react";
import * as GQL from "src/core/generated-graphql";
import { View } from "src/components/List/views";
import { GroupList } from "../GroupList";
import { useContainingGroupFilterHook } from "src/core/groups";

interface IGroupSubGroupsPanel {
  active: boolean;
  group: GQL.GroupDataFragment;
}

export const GroupSubGroupsPanel: React.FC<IGroupSubGroupsPanel> = ({
  active,
  group,
}) => {
  const filterHook = useContainingGroupFilterHook(group);

  return (
    <GroupList
      filterHook={filterHook}
      alterQuery={active}
      view={View.GroupSubGroups}
    />
  );
};
