import React from "react";
import { FindMoviesQueryResult } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { useMoviesList } from "src/hooks/ListHook";
import { MovieCard } from "./MovieCard";

export const MovieList: React.FC = () => {
  const listData = useMoviesList({
    renderContent,
    persistState: true,
  });

  function renderContent(
    result: FindMoviesQueryResult,
    filter: ListFilterModel
  ) {
    if (!result.data?.findMovies) {
      return;
    }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row justify-content-center">
          {result.data.findMovies.movies.map((p) => (
            <MovieCard key={p.id} movie={p} />
          ))}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return <h1>TODO</h1>;
    }
  }

  return listData.template;
};
