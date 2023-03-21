import React from "react";
import * as GQL from "src/core/generated-graphql";
import { usePerformerFilterHook } from "src/core/performers";
import { PerformerList } from "src/components/Performers/PerformerList";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";

interface IPerformerAppearsWithPanel {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerAppearsWithPanel: React.FC<
  IPerformerAppearsWithPanel
> = ({ active, performer }) => {
  const criterion = new PerformersCriterion();
  criterion.value = [
    { id: performer.id, label: performer.name || `Performer ${performer.id}` },
  ];

  const performerValue = {
    id: performer.id,
    label: performer.name ?? `Performer ${performer.id}`,
  };

  const extraCriteria = {
    performer: performerValue,
  };

  const filterHook = usePerformerFilterHook(performer);

  return (
    <PerformerList
      filterHook={filterHook}
      extraCriteria={extraCriteria}
      alterQuery={active}
    />
  );
};
