import _ from "lodash";
import React, { FunctionComponent } from "react";
import * as GQL from "../../../core/generated-graphql";
import { IBaseProps } from "../../../models";
import { SceneList } from "../../scenes/SceneList";
import { PerformersCriterion } from "../../../models/list-filter/criteria/performers";

interface IPerformerDetailsProps {
  performer: Partial<GQL.PerformerDataFragment>
  base: IBaseProps
}

export const PerformerScenesPanel: FunctionComponent<IPerformerDetailsProps> = (props: IPerformerDetailsProps) => {

  function makeCriteria() {
    let performerCriterion = new PerformersCriterion();
    performerCriterion.value = [{id: props.performer.id!, label: props.performer.name!}];
    return [performerCriterion];
  }

  return (
    <SceneList 
      base={props.base} 
      subComponent={true} 
      criteria={makeCriteria()}
    />
  );
}