import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useGroupFilterHook } from "src/core/groups";
import { PerformerList } from "src/components/Performers/PerformerList";
import { GroupsCriterion, GroupsCriterionOption } from "src/models/list-filter/criteria/groups";
import { View } from "src/components/List/views";

interface IGroupPerformersPanel {
  active: boolean;
  group: GQL.GroupDataFragment;
  showChildGroupContent?: boolean;
}

export const GroupPerformersPanel: React.FC<IGroupPerformersPanel> = ({
  active,
  group,
  showChildGroupContent,
}) => {
  const groupCriterion = new GroupsCriterion(GroupsCriterionOption);
  groupCriterion.value = {
    items: [{ id: group.id!, label: group.name || `Group ${group.id}` }],
    excluded: [],
    depth: 0,
  };

  const extraCriteria = {
    scenes: [groupCriterion],
  };

  const filterHook = useGroupFilterHook(group, showChildGroupContent);

  return (
    <PerformerList
      filterHook={filterHook}
      extraCriteria={extraCriteria}
      alterQuery={active}
      view={View.GroupPerformers}
    />
  );
}; 