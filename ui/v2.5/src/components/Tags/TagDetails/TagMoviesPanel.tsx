import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { MovieList } from "src/components/Movies/MovieList";

export const TagMoviesPanel: React.FC<{
  active: boolean;
  tag: GQL.TagDataFragment;
}> = ({ active, tag }) => {
  const filterHook = useTagFilterHook(tag);
  return <MovieList filterHook={filterHook} alterQuery={active} />;
};
