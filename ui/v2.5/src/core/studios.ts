import * as GQL from "src/core/generated-graphql";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { ListFilterModel } from "src/models/list-filter/filter";
import React from "react";
import { ConfigurationContext } from "src/hooks/Config";
import { IUIConfig } from "./config";

export const useStudioFilterHook = (studio: GQL.StudioDataFragment) => {
  const config = React.useContext(ConfigurationContext);
  return (filter: ListFilterModel) => {
    const studioValue = { id: studio.id, label: studio.name };
    // if studio is already present, then we modify it, otherwise add
    let studioCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "studios";
    }) as StudiosCriterion | undefined;

    if (studioCriterion) {
      // we should be showing studio only. Remove other values
      studioCriterion.value.items = [studioValue];
      studioCriterion.modifier = GQL.CriterionModifier.Includes;
    } else {
      studioCriterion = new StudiosCriterion();
      studioCriterion.value = {
        items: [studioValue],
        excluded: [],
        depth: (config?.configuration?.ui as IUIConfig)?.showChildStudioContent
          ? -1
          : 0,
      };
      studioCriterion.modifier = GQL.CriterionModifier.Includes;
      filter.criteria.push(studioCriterion);
    }

    return filter;
  };
};
