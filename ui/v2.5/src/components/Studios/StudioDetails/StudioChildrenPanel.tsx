import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ParentStudiosCriterion } from "src/models/list-filter/criteria/studios";
import { ListFilterModel } from "src/models/list-filter/filter";
import { StudioList } from "../StudioList";
import { View } from "src/components/List/views";

function useFilterHook(studio: GQL.StudioDataFragment) {
  return (filter: ListFilterModel) => {
    const studioValue = { id: studio.id!, label: studio.name! };
    // if studio is already present, then we modify it, otherwise add
    let parentStudioCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "parents";
    }) as ParentStudiosCriterion | undefined;

    if (
      parentStudioCriterion &&
      (parentStudioCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        parentStudioCriterion.modifier === GQL.CriterionModifier.Includes)
    ) {
      // add the studio if not present
      if (
        !parentStudioCriterion.value.find((p) => {
          return p.id === studio.id;
        })
      ) {
        parentStudioCriterion.value.push(studioValue);
      }

      parentStudioCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      parentStudioCriterion = new ParentStudiosCriterion();
      parentStudioCriterion.value = [studioValue];
      filter.criteria.push(parentStudioCriterion);
    }

    return filter;
  };
}

interface IStudioChildrenPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioChildrenPanel: React.FC<IStudioChildrenPanel> = ({
  active,
  studio,
}) => {
  const filterHook = useFilterHook(studio);

  return (
    <StudioList
      fromParent
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioChildren}
    />
  );
};
