import React from "react";
import { Link } from "react-router-dom";
import { useFindImages } from "src/core/StashService";
import Slider from "@ant-design/react-slick";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";
import { ImageCard } from "./ImageCard";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const ImageRecommendationRow: React.FC<IProps> = (props: IProps) => {
  const result = useFindImages(props.filter);
  const cardCount = result.data?.findImages.count;

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="images-recommendations"
      header={props.header}
      link={
        <Link to={`/images?${props.filter.makeQueryParameters()}`}>
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
              <div key={`_${i}`} className="image-skeleton skeleton-card"></div>
            ))
          : result.data?.findImages.images.map((i) => (
              <ImageCard key={i.id} image={i} zoomIndex={1} />
            ))}
      </Slider>
    </RecommendationRow>
  );
};
