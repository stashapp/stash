import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieList } from "src/components/Movies/MovieList";
import { studioFilterHook } from "src/core/studios";

interface IStudioMoviesPanel {
  studio: Partial<GQL.StudioDataFragment>;
}

export const StudioMoviesPanel: React.FC<IStudioMoviesPanel> = ({ studio }) => {
  return <MovieList filterHook={studioFilterHook(studio)} />;
};
