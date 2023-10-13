import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieList } from "src/components/Movies/MovieList";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerMoviesPanel: React.FC<IPerformerDetailsProps> = ({
  active,
  performer,
}) => {
  const filterHook = usePerformerFilterHook(performer);
  return (
    <MovieList
      filterHook={filterHook}
      alterQuery={active}
      view={View.PerformerMovies}
    />
  );
};
