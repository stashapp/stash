import React, { FunctionComponent } from "react";
import { FindPerformersQueryResult } from "src/core/generated-graphql";
import Slider from "react-slick";
import { PerformerCard } from "./PerformerCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  result: FindPerformersQueryResult;
  header: String;
  linkText: String;
}

export const PerformerRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const cardCount = props.result.data?.findPerformers.count;
  return (
    <div className="recommendation-row performer-recommendations">
      <div className="recommendation-row-head">
        <div>
          <h2>{props.header}</h2>
        </div>
        <a href={`/performers?${props.filter.makeQueryParameters()}`}>
          {props.linkText}
        </a>
      </div>
      <Slider {...getSlickSliderSettings(cardCount!, props.isTouch)}>
        {props.result.data?.findPerformers.performers.map((p) => (
          <PerformerCard key={p.id} performer={p} />
        ))}
      </Slider>
    </div>
  );
};
