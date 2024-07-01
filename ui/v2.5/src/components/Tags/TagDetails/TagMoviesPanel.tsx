import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { GroupList } from "src/components/Movies/MovieList";

export const TagGroupsPanel: React.FC<{
  active: boolean;
  tag: GQL.TagDataFragment;
}> = ({ active, tag }) => {
  const filterHook = useTagFilterHook(tag);
  return <GroupList filterHook={filterHook} alterQuery={active} />;
};
