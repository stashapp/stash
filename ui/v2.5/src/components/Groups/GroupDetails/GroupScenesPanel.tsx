import React from "react";
import * as GQL from "src/core/generated-graphql";
import {
  GroupsCriterion,
  GroupsCriterionOption,
} from "src/models/list-filter/criteria/groups";
import { ListFilterModel } from "src/models/list-filter/filter";
import { FilteredSceneList } from "src/components/Scenes/SceneList";
import { View } from "src/components/List/views";

interface IGroupScenesPanel {
  active: boolean;
  group: GQL.GroupDataFragment;
  showSubGroupContent?: boolean;
}

function useFilterHook(
  group: Pick<GQL.GroupDataFragment, "id" | "name">,
  showSubGroupContent?: boolean
) {
  return (filter: ListFilterModel) => {
    const groupValue = { id: group.id, label: group.name };
    // if group is already present, then we modify it, otherwise add
    let groupCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "groups";
    }) as GroupsCriterion | undefined;

    if (
      groupCriterion &&
      (groupCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        groupCriterion.modifier === GQL.CriterionModifier.Includes)
    ) {
      // add the group if not present
      if (
        !groupCriterion.value.items.find((p) => {
          return p.id === group.id;
        })
      ) {
        groupCriterion.value.items.push(groupValue);
      }

      groupCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      groupCriterion = new GroupsCriterion(GroupsCriterionOption);
      groupCriterion.value = {
        items: [groupValue],
        depth: showSubGroupContent ? -1 : 0,
        excluded: [],
      };
      filter.criteria.push(groupCriterion);
    }

    return filter;
  };
}

export const GroupScenesPanel: React.FC<IGroupScenesPanel> = ({
  active,
  group,
  showSubGroupContent,
}) => {
  const filterHook = useFilterHook(group, showSubGroupContent);

  if (group && group.id) {
    return (
      <FilteredSceneList
        filterHook={filterHook}
        defaultSort="group_scene_number"
        alterQuery={active}
        view={View.GroupScenes}
        fromGroupId={group.id}
      />
    );
  }
  return <></>;
};
