import React, { useRef } from "react";
import * as GQL from "src/core/generated-graphql";
import { MovieCard } from "./MovieCard";
import { useContainerDimensions } from "../Shared/GridCard/GridCard";

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
  const componentRef = useRef<HTMLDivElement>(null);
  const { width } = useContainerDimensions(componentRef);
  return (
    <div className="row justify-content-center" ref={componentRef}>
      {movies.map((p) => (
        <MovieCard
          key={p.id}
          containerWidth={width}
          movie={p}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(p.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(p.id, selected, shiftKey)
          }
        />
      ))}
    </div>
  );
};
