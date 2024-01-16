import React from "react";
import * as GQL from "src/core/generated-graphql";
import { PerformerList } from "src/components/Performers/PerformerList";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerAppearsWithPanel: React.FC<IPerformerDetailsProps> = ({
  active,
  performer,
}) => {
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
      view={View.PerformerAppearsWith}
    />
  );
};
