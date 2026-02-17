import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { FilteredGroupList } from "src/components/Groups/GroupList";
import { View } from "src/components/List/views";

export const TagGroupsPanel: React.FC<{
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}> = ({ active, tag, showSubTagContent }) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return (
    <FilteredGroupList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagGroups}
    />
  );
};
