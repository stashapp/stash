import React from "react";
import * as GQL from "src/core/generated-graphql";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SceneList } from "src/components/Scenes/SceneList";

interface IPerformerDetailsProps {
  performer: Partial<GQL.PerformerDataFragment>;
}

export const PerformerScenesPanel: React.FC<IPerformerDetailsProps> = ({
  performer,
}) => {
  function filterHook(filter: ListFilterModel) {
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
  }

  return <SceneList filterHook={filterHook} />;
};
