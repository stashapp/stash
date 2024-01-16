import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { PerformerList } from "src/components/Performers/PerformerList";
import { View } from "src/components/List/views";

interface ITagPerformersPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
}

export const TagPerformersPanel: React.FC<ITagPerformersPanel> = ({
  active,
  tag,
}) => {
  const filterHook = useTagFilterHook(tag);
  return (
    <PerformerList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagPerformers}
    />
  );
};
