import React, { FunctionComponent } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindMoviesQuery, FindMoviesVariables } from "../../core/generated-graphql";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";
import { MovieCard } from "./MovieCard";

interface IProps extends IBaseProps {}

export const MovieList: FunctionComponent<IProps> = (props: IProps) => {
  const listData = ListHook.useList({
    filterMode: FilterMode.Movies,
    props,
    renderContent,
  });

  function renderContent(result: QueryHookResult<FindMoviesQuery, FindMoviesVariables>, filter: ListFilterModel) {
    if (!result.data || !result.data.findMovies) { return; }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="grid">
          {result.data.findMovies.movies.map((movie) => (<MovieCard key={movie.id} movie={movie}/>))}
        </div>
      );
    } else if (filter.displayMode === DisplayMode.List) {
      return <h1>TODO</h1>;
    } 
  }

  return listData.template;
};
