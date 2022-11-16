import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieList } from "src/components/Movies/MovieList";
import { usePerformerFilterHook } from "src/core/performers";

interface IPerformerDetailsProps {
  performer: GQL.PerformerDataFragment;
}

export const PerformerMoviesPanel: React.FC<IPerformerDetailsProps> = ({
  performer,
}) => {
  const filterHook = usePerformerFilterHook(performer);
  return <MovieList filterHook={filterHook} />;
};
