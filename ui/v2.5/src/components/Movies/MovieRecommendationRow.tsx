import React, { FunctionComponent } from "react";
import { FindMoviesQueryResult } from "src/core/generated-graphql";
import Slider from "react-slick";
import { MovieCard } from "./MovieCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  result: FindMoviesQueryResult;
  header: String;
  linkText: String;
}

export const MovieRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const cardCount = props.result.data?.findMovies.count;
  return (
    <div className="recommendation-row movie-recommendations">
      <div className="recommendation-row-head">
        <div>
          <h2>{props.header}</h2>
        </div>
        <a href={`/movies?${props.filter.makeQueryParameters()}`}>
          {props.linkText}
        </a>
      </div>
      <Slider {...getSlickSliderSettings(cardCount!, props.isTouch)}>
        {props.result.data?.findMovies.movies.map((p) => (
          <MovieCard key={p.id} movie={p} />
        ))}
      </Slider>
    </div>
  );
};
