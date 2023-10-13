import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieList } from "src/components/Movies/MovieList";
import { useStudioFilterHook } from "src/core/studios";
import { View } from "src/components/List/views";

interface IStudioMoviesPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioMoviesPanel: React.FC<IStudioMoviesPanel> = ({
  active,
  studio,
}) => {
  const filterHook = useStudioFilterHook(studio);
  return (
    <MovieList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioMovies}
    />
  );
};
