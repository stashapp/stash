import React from "react";
import * as GQL from "src/core/generated-graphql";
import { PerformerList } from "src/components/Performers/PerformerList";
import { usePerformerFilterHook } from "src/core/performers";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerAppearsWithPanel: React.FC<IPerformerDetailsProps> = ({
  active,
  performer,
}) => {
  const filterHook = usePerformerFilterHook(performer);

  return <PerformerList filterHook={filterHook} alterQuery={active} />;
};
