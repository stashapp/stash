import React, { FunctionComponent } from "react";
import { useFindStudios } from "src/core/StashService";
import Slider from "react-slick";
import { StudioCard } from "./StudioCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: String;
  linkText: String;
  index: number;
}

export const StudioRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const result = useFindStudios(props.filter);
  const cardCount = result.data?.findStudios.count;
  if (result.loading) {
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
        <Slider
          {...getSlickSliderSettings(props.filter.itemsPerPage!, props.isTouch)}
        >
          {[...Array(props.filter.itemsPerPage)].map((i) => (
            <div key={i} className="studio-skeleton skeleton-card"></div>
          ))}
        </Slider>
      </div>
    );
  }

  if (cardCount === 0) {
    return null;
  }

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
        {result.data?.findStudios.studios.map((s) => (
          <StudioCard key={s.id} studio={s} hideParent={true} />
        ))}
      </Slider>
    </div>
  );
};
