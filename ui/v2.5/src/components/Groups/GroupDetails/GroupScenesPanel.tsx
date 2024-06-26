import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupsCriterion } from "src/models/list-filter/criteria/groups";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SceneList } from "src/components/Scenes/SceneList";
import { View } from "src/components/List/views";

interface IGroupScenesPanel {
  active: boolean;
  group: GQL.GroupDataFragment;
}

export const GroupScenesPanel: React.FC<IGroupScenesPanel> = ({
  active,
  group,
}) => {
  function filterHook(filter: ListFilterModel) {
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
        !groupCriterion.value.find((p) => {
          return p.id === group.id;
        })
      ) {
        groupCriterion.value.push(groupValue);
      }

      groupCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      groupCriterion = new GroupsCriterion();
      groupCriterion.value = [groupValue];
      filter.criteria.push(groupCriterion);
    }

    return filter;
  }

  if (group && group.id) {
    return (
      <SceneList
        filterHook={filterHook}
        defaultSort="group_scene_number"
        alterQuery={active}
        view={View.GroupScenes}
      />
    );
  }
  return <></>;
};
