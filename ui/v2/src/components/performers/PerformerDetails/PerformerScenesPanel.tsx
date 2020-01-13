import React, { FunctionComponent } from "react";
import * as GQL from "../../../core/generated-graphql";
import { IBaseProps } from "../../../models";
import { SceneList } from "../../scenes/SceneList";
import { PerformersCriterion } from "../../../models/list-filter/criteria/performers";
import { ListFilterModel } from "../../../models/list-filter/filter";

interface IPerformerDetailsProps {
  performer: Partial<GQL.PerformerDataFragment>
  base: IBaseProps
}

export const PerformerScenesPanel: FunctionComponent<IPerformerDetailsProps> = (props: IPerformerDetailsProps) => {

  function filterHook(filter: ListFilterModel) {
    let performerValue = {id: props.performer.id!, label: props.performer.name!};
    // if performers is already present, then we modify it, otherwise add
    let performerCriterion = filter.criteria.find((c) => {
      return c.type === "performers";
    });

    if (performerCriterion && 
        (performerCriterion.modifier === GQL.CriterionModifier.IncludesAll || 
         performerCriterion.modifier === GQL.CriterionModifier.Includes)) {
      // add the performer if not present
      if (!performerCriterion.value.find((p : any) => {
        return p.id === props.performer.id;
      })) {
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

  return (
    <SceneList 
      base={props.base} 
      subComponent={true} 
      filterHook={filterHook}
    />
  );
}