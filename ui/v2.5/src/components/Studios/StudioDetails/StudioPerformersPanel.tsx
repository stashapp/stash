import React from "react";
import * as GQL from "src/core/generated-graphql";
import { studioFilterHook } from "src/core/studios";
import { PerformerList } from "src/components/Performers/PerformerList";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";

interface IStudioPerformersPanel {
  studio: Partial<GQL.StudioDataFragment>;
}

export const StudioPerformersPanel: React.FC<IStudioPerformersPanel> = ({
  studio,
}) => {
  const studioCriterion = new StudiosCriterion();
  studioCriterion.value = [
    { id: studio.id!, label: studio.name || `Studio ${studio.id}` },
  ];

  const extraCriteria = {
    scenes: [studioCriterion],
    images: [studioCriterion],
    galleries: [studioCriterion],
  };

  return (
    <PerformerList
      filterHook={studioFilterHook(studio)}
      extraCriteria={extraCriteria}
    />
  );
};
