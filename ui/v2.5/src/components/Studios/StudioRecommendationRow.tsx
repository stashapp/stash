import React, { FunctionComponent } from "react";
import { FindStudiosQueryResult } from "src/core/generated-graphql";
import Slider from "react-slick";
import { StudioCard } from "./StudioCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  result: FindStudiosQueryResult;
  header: String;
  linkText: String;
}

export const StudioRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const cardCount = props.result.data?.findStudios.count;
  return (
    <div className="recommendation-row studio-recommendations">
      <div className="recommendation-row-head">
        <div>
          <h2>{props.header}</h2>
        </div>
        <a href={`/studios?${props.filter.makeQueryParameters()}`}>
          {props.linkText}
        </a>
      </div>
      <Slider {...getSlickSliderSettings(cardCount!, props.isTouch)}>
        {props.result.data?.findStudios.studios.map((studio) => (
          <StudioCard key={studio.id} studio={studio} hideParent={true} />
        ))}
      </Slider>
    </div>
  );
};
