import React, { FunctionComponent } from "react";
import { useFindGalleries } from "src/core/StashService";
import Slider from "react-slick";
import { GalleryCard } from "./GalleryCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const GalleryRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const result = useFindGalleries(props.filter);
  const cardCount = result.data?.findGalleries.count;

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="gallery-recommendations"
      header={props.header}
      link={
        <a href={`/galleries?${props.filter.makeQueryParameters()}`}>
          <FormattedMessage id="view_all" />
        </a>
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
              <div
                key={`_${i}`}
                className="gallery-skeleton skeleton-card"
              ></div>
            ))
          : result.data?.findGalleries.galleries.map((g) => (
              <GalleryCard key={g.id} gallery={g} zoomIndex={1} />
            ))}
      </Slider>
    </RecommendationRow>
  );
};
