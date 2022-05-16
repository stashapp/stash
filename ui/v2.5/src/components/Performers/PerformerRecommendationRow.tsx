import React, { FunctionComponent } from "react";
import { useFindPerformers } from "src/core/StashService";
import Slider from "react-slick";
import { PerformerCard } from "./PerformerCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: String;
  linkText: String;
  loadingArr: boolean[];
  index: number;
}

export const PerformerRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const result = useFindPerformers(props.filter);
  const cardCount = result.data?.findPerformers.count;
  props.loadingArr[props.index] = result.loading;
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
        {result.data?.findPerformers.performers.map((p) => (
          <PerformerCard key={p.id} performer={p} />
        ))}
      </Slider>
    </div>
  );
};
