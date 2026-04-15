import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { FilteredPerformerList } from "src/components/Performers/PerformerList";
import { View } from "src/components/List/views";

interface ITagPerformersPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}

export const TagPerformersPanel: React.FC<ITagPerformersPanel> = ({
  active,
  tag,
  showSubTagContent,
}) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return (
    <FilteredPerformerList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagPerformers}
    />
  );
};
