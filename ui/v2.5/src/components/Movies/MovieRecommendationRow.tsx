import React from "react";
import { Link } from "react-router-dom";
import { useFindMovies } from "src/core/StashService";
import Slider from "@ant-design/react-slick";
import { MovieCard } from "./MovieCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const MovieRecommendationRow: React.FC<IProps> = (props: IProps) => {
  const result = useFindMovies(props.filter);
  const cardCount = result.data?.findMovies.count;

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="movie-recommendations"
      header={props.header}
      link={
        <Link to={`/movies?${props.filter.makeQueryParameters()}`}>
          <FormattedMessage id="view_all" />
        </Link>
      }
    >
      <Slider
        {...getSlickSliderSettings(
          cardCount ? cardCount : props.filter.itemsPerPage,
          props.isTouch
        )}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div key={`_${i}`} className="movie-skeleton skeleton-card"></div>
            ))
          : result.data?.findMovies.movies.map((m) => (
              <MovieCard key={m.id} movie={m} />
            ))}
      </Slider>
    </RecommendationRow>
  );
};
