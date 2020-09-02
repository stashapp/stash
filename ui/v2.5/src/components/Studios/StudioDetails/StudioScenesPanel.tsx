import React from "react";
import * as GQL from "src/core/generated-graphql";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SceneList } from "src/components/Scenes/SceneList";

interface IStudioScenesPanel {
  studio: Partial<GQL.StudioDataFragment>;
}

export const StudioScenesPanel: React.FC<IStudioScenesPanel> = ({ studio }) => {
  function filterHook(filter: ListFilterModel) {
    const studioValue = { id: studio.id!, label: studio.name! };
    // if studio is already present, then we modify it, otherwise add
    let studioCriterion = filter.criteria.find((c) => {
      return c.type === "studios";
    }) as StudiosCriterion;

    if (
      studioCriterion &&
      (studioCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        studioCriterion.modifier === GQL.CriterionModifier.Includes)
    ) {
      // add the studio if not present
      if (
        !studioCriterion.value.find((p) => {
          return p.id === studio.id;
        })
      ) {
        studioCriterion.value.push(studioValue);
      }

      studioCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      studioCriterion = new StudiosCriterion();
      studioCriterion.value = [studioValue];
      filter.criteria.push(studioCriterion);
    }

    return filter;
  }

  return <SceneList filterHook={filterHook} />;
};
