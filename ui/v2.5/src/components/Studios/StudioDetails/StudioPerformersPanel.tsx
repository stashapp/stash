import React from "react";
import * as GQL from "src/core/generated-graphql";
import { PerformerList } from "src/components/Performers/PerformerList";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { ConfigurationContext } from "src/hooks/Config";
import { IUIConfig } from "src/core/config";

interface IStudioPerformersPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioPerformersPanel: React.FC<IStudioPerformersPanel> = ({
  active,
  studio,
}) => {
  const studioCriterion = new StudiosCriterion();
  const config = React.useContext(ConfigurationContext);
  const depth = (config?.configuration?.ui as IUIConfig)?.showChildStudioContent
    ? -1
    : 0;

  studioCriterion.value = {
    items: [{ id: studio.id!, label: studio.name || `Studio ${studio.id}` }],
    excluded: [],
    depth: depth,
  };

  const extraCriteria = {
    scenes: [studioCriterion],
    images: [studioCriterion],
    galleries: [studioCriterion],
    movies: [studioCriterion],
  };

  const queryArgs = {
    id: studio.id,
    depth: depth,
    type: "STUDIO",
  };

  return (
    <PerformerList
      queryArgs={queryArgs}
      extraCriteria={extraCriteria}
      alterQuery={active}
    />
  );
};
