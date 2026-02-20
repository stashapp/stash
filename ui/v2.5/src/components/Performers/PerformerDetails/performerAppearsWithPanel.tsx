import React from "react";
import * as GQL from "src/core/generated-graphql";
import { FilteredPerformerList } from "src/components/Performers/PerformerList";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";
import { PatchComponent } from "src/patch";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerAppearsWithPanel: React.FC<IPerformerDetailsProps> =
  PatchComponent("PerformerAppearsWithPanel", ({ active, performer }) => {
    const performerValue = {
      id: performer.id,
      label: performer.name ?? `Performer ${performer.id}`,
    };

    const extraCriteria = {
      performer: performerValue,
    };

    const filterHook = usePerformerFilterHook(performer);

    return (
      <FilteredPerformerList
        filterHook={filterHook}
        extraCriteria={extraCriteria}
        alterQuery={active}
        view={View.PerformerAppearsWith}
      />
    );
  });
