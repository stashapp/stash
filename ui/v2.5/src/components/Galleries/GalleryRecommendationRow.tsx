import React, { FunctionComponent } from "react";
import { FindGalleriesQueryResult } from "src/core/generated-graphql";
import Slider from "react-slick";
import { GalleryCard } from "./GalleryCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  result: FindGalleriesQueryResult;
  header: String;
  linkText: String;
}

export const GalleryRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const cardCount = props.result.data?.findGalleries.count;
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
        {props.result.data?.findGalleries.galleries.map((gallery) => (
          <GalleryCard key={gallery.id} gallery={gallery} zoomIndex={1} />
        ))}
      </Slider>
    </div>
  );
};
