import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieList } from "src/components/Movies/MovieList";
import { performerFilterHook } from "src/core/performers";

interface IPerformerDetailsProps {
  performer: Partial<GQL.PerformerDataFragment>;
}

export const PerformerMoviesPanel: React.FC<IPerformerDetailsProps> = ({
  performer,
}) => {
  return <MovieList filterHook={performerFilterHook(performer)} />;
};
