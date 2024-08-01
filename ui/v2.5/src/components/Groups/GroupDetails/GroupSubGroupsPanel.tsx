import React from "react";
import * as GQL from "src/core/generated-graphql";
import { View } from "src/components/List/views";
import { GroupList } from "../GroupList";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  ContainingGroupsCriterionOption,
  GroupsCriterion,
} from "src/models/list-filter/criteria/groups";

export const useContainingGroupFilterHook = (
  group: Pick<GQL.StudioDataFragment, "id" | "name">,
  showSubGroupContent?: boolean
) => {
  return (filter: ListFilterModel) => {
    const groupValue = { id: group.id, label: group.name };
    // if studio is already present, then we modify it, otherwise add
    let groupCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "containing_groups";
    }) as GroupsCriterion | undefined;

    if (groupCriterion) {
      // add the group if not present
      if (
        !groupCriterion.value.items.find((p) => {
          return p.id === group.id;
        })
      ) {
        groupCriterion.value.items.push(groupValue);
      }
    } else {
      groupCriterion = new GroupsCriterion(ContainingGroupsCriterionOption);
      groupCriterion.value = {
        items: [groupValue],
        excluded: [],
        depth: showSubGroupContent ? -1 : 0,
      };
      groupCriterion.modifier = GQL.CriterionModifier.Includes;
      filter.criteria.push(groupCriterion);
    }

    return filter;
  };
};

interface IGroupSubGroupsPanel {
  active: boolean;
  group: GQL.GroupDataFragment;
  showSubGroupContent?: boolean;
}

export const GroupSubGroupsPanel: React.FC<IGroupSubGroupsPanel> = ({
  active,
  group,
  showSubGroupContent,
}) => {
  const filterHook = useContainingGroupFilterHook(group, showSubGroupContent);

  return (
    <GroupList
      filterHook={filterHook}
      alterQuery={active}
      view={View.GroupSubGroups}
      fromGroupId={group.id}
    />
  );
};
