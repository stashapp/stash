import React from "react";
import { Link } from "react-router-dom";
import { useFindStudios } from "src/core/StashService";
import Slider from "@ant-design/react-slick";
import { StudioCard } from "./StudioCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const StudioRecommendationRow: React.FC<IProps> = (props) => {
  const result = useFindStudios(props.filter);
  const cardCount = result.data?.findStudios.count;

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="studio-recommendations"
      header={props.header}
      link={
        <Link to={`/studios?${props.filter.makeQueryParameters()}`}>
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
              <div
                key={`_${i}`}
                className="studio-skeleton skeleton-card"
              ></div>
            ))
          : result.data?.findStudios.studios.map((s) => (
              <StudioCard key={s.id} studio={s} hideParent={true} />
            ))}
      </Slider>
    </RecommendationRow>
  );
};
