import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useStudioFilterHook } from "src/core/studios";
import { PerformerList } from "src/components/Performers/PerformerList";

interface IStudioPerformersPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioPerformersPanel: React.FC<IStudioPerformersPanel> = ({
  active,
  studio,
}) => {
  const filterHook = useStudioFilterHook(studio);

  return <PerformerList filterHook={filterHook} alterQuery={active} />;
};
