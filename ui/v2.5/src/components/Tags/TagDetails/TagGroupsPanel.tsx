import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { EnhancedGroupList as GroupList } from "src/extensions/facets/enhanced";

export const TagGroupsPanel: React.FC<{
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}> = ({ active, tag, showSubTagContent }) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return <GroupList filterHook={filterHook} alterQuery={active} />;
};
