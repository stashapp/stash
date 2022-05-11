import React, { FunctionComponent } from "react";
import { FindMoviesQueryResult } from "src/core/generated-graphql";
import Slider from "react-slick";
import { MovieCard } from "./MovieCard";

interface IProps {
  isTouch: boolean;
  result: FindMoviesQueryResult;
  header: String;
}

export const MovieRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  function determineSlidesToScroll(prefered: number, cardCount: number) {
    if (props.isTouch) {
      return 1;
    } else if (cardCount! > prefered) {
      return prefered;
    } else {
      return cardCount;
    }
  }

  const cardCount = props.result.data?.findMovies.count;
  var settings = {
    dots: !props.isTouch,
    arrows: !props.isTouch,
    infinite: !props.isTouch,
    speed: 300,
    variableWidth: true,
    swipeToSlide: true,
    slidesToShow: cardCount! > 5 ? 5 : cardCount,
    slidesToScroll: determineSlidesToScroll(5, cardCount!),
    responsive: [
      {
        breakpoint: 1909,
        settings: {
          slidesToShow: cardCount! > 4 ? 4 : cardCount,
          slidesToScroll: determineSlidesToScroll(4, cardCount!),
        },
      },
      {
        breakpoint: 1542,
        settings: {
          slidesToShow: cardCount! > 3 ? 3 : cardCount,
          slidesToScroll: determineSlidesToScroll(3, cardCount!),
        },
      },
      {
        breakpoint: 1170,
        settings: {
          slidesToShow: cardCount! > 2 ? 2 : cardCount,
          slidesToScroll: determineSlidesToScroll(2, cardCount!),
        },
      },
      {
        breakpoint: 801,
        settings: {
          slidesToShow: 1,
          slidesToScroll: 1,
          dots: false,
        },
      },
    ],
  };

  return (
    <div className="recommendation-row movie-recommendations">
      <div className="recommendation-row-head">
        <div>
          <h2>{props.header}</h2>
        </div>
        <a href="/scenes?sortby=date&sortdir=desc">View all</a>
      </div>
      <Slider {...settings}>
        {props.result.data?.findMovies.movies.map((p) => (
          <MovieCard key={p.id} movie={p} />
        ))}
      </Slider>
    </div>
  );
};
