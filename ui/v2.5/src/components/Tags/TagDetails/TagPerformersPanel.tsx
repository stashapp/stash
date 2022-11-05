import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { PerformerList } from "src/components/Performers/PerformerList";

interface ITagPerformersPanel {
  tag: GQL.TagDataFragment;
}

export const TagPerformersPanel: React.FC<ITagPerformersPanel> = ({ tag }) => {
  const filterHook = useTagFilterHook(tag);
  return <PerformerList filterHook={filterHook} />;
};
