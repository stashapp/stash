import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieCard } from "./MovieCard";

interface IMovieCardGrid {
  movies: GQL.MovieDataFragment[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const MovieCardGrid: React.FC<IMovieCardGrid> = ({
  movies,
  selectedIds,
  onSelectChange,
}) => {
  return (
    <div className="row justify-content-center">
      {movies.map((movie) => (
        <MovieCard
          key={movie.id}
          movie={movie}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(movie.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(movie.id, selected, shiftKey)
          }
        />
      ))}
    </div>
  );
};
