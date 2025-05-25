import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useStudioFilterHook } from "src/core/studios";
import { PerformerList } from "src/components/Performers/PerformerList";
import { View } from "src/components/List/views";

interface IStudioPerformersPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
  showChildStudioContent?: boolean;
}

export const StudioPerformersPanel: React.FC<IStudioPerformersPanel> = ({
  active,
  studio,
  showChildStudioContent,
}) => {
  const filterHook = useStudioFilterHook(studio, showChildStudioContent);

  return (
    <PerformerList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioPerformers}
    />
  );
};
