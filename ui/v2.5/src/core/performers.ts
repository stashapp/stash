import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import * as GQL from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";

export const performerFilterHook = (
  performer: Partial<GQL.PerformerDataFragment>
) => {
  return (filter: ListFilterModel) => {
    const performerValue = { id: performer.id!, label: performer.name! };
    // if performers is already present, then we modify it, otherwise add
    let performerCriterion = filter.criteria.find((c) => {
      return c.type === "performers";
    }) as PerformersCriterion;

    if (
      performerCriterion &&
      (performerCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        performerCriterion.modifier === GQL.CriterionModifier.Includes)
    ) {
      // add the performer if not present
      if (
        !performerCriterion.value.find((p) => {
          return p.id === performer.id;
        })
      ) {
        performerCriterion.value.push(performerValue);
      }

      performerCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      performerCriterion = new PerformersCriterion();
      performerCriterion.value = [performerValue];
      filter.criteria.push(performerCriterion);
    }

    return filter;
  };
};
