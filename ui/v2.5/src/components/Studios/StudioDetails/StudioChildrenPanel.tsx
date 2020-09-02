import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ParentStudiosCriterion } from "src/models/list-filter/criteria/studios";
import { ListFilterModel } from "src/models/list-filter/filter";
import { StudioList } from "../StudioList";

interface IStudioChildrenPanel {
  studio: Partial<GQL.StudioDataFragment>;
}

export const StudioChildrenPanel: React.FC<IStudioChildrenPanel> = ({
  studio,
}) => {
  function filterHook(filter: ListFilterModel) {
    const studioValue = { id: studio.id!, label: studio.name! };
    // if studio is already present, then we modify it, otherwise add
    let parentStudioCriterion = filter.criteria.find((c) => {
      return c.type === "parent_studios";
    }) as ParentStudiosCriterion;

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
  }

  return <StudioList fromParent filterHook={filterHook} />;
};
