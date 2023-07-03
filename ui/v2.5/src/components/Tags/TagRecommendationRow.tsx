import React from "react";
import { Link } from "react-router-dom";
import { useFindTags } from "src/core/StashService";
import Slider from "@ant-design/react-slick";
import { TagCard } from "./TagCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const TagRecommendationRow: React.FC<IProps> = (props) => {
  const result = useFindTags(props.filter);
  const cardCount = result.data?.findTags.count;

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="tag-recommendations"
      header={props.header}
      link={
        <Link to={`/tags?${props.filter.makeQueryParameters()}`}>
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
              <div key={`_${i}`} className="tag-skeleton skeleton-card"></div>
            ))
          : result.data?.findTags.tags.map((p) => (
              <TagCard key={p.id} tag={p} zoomIndex={0} />
            ))}
      </Slider>
    </RecommendationRow>
  );
};
