import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { FilteredStudioList } from "src/components/Studios/StudioList";
import { View } from "src/components/List/views";

interface ITagStudiosPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}

export const TagStudiosPanel: React.FC<ITagStudiosPanel> = ({
  active,
  tag,
  showSubTagContent,
}) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return (
    <FilteredStudioList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagStudios}
    />
  );
};
