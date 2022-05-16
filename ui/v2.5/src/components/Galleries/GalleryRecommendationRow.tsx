import React, { FunctionComponent } from "react";
import { useFindGalleries } from "src/core/StashService";
import Slider from "react-slick";
import { GalleryCard } from "./GalleryCard";
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

export const GalleryRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const result = useFindGalleries(props.filter);
  const cardCount = result.data?.findGalleries.count;
  props.loadingArr[props.index] = result.loading;

  return (
    <div className="recommendation-row gallery-recommendations">
      <div className="recommendation-row-head">
        <div>
          <h2>{props.header}</h2>
        </div>
        <a href={`/galleries?${props.filter.makeQueryParameters()}`}>
          {props.linkText}
        </a>
      </div>
      <Slider {...getSlickSliderSettings(cardCount!, props.isTouch)}>
        {result.data?.findGalleries.galleries.map((g) => (
          <GalleryCard key={g.id} gallery={g} zoomIndex={1} />
        ))}
      </Slider>
    </div>
  );
};
