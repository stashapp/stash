import React from "react";
import { Link } from "react-router-dom";
import { useFindGroups } from "src/core/StashService";
import Slider from "@ant-design/react-slick";
import { GroupCard } from "./GroupCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const GroupRecommendationRow: React.FC<IProps> = (props: IProps) => {
  const result = useFindGroups(props.filter);
  const cardCount = result.data?.findGroups.count;

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="group-recommendations"
      header={props.header}
      link={
        <Link to={`/groups?${props.filter.makeQueryParameters()}`}>
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
              <div key={`_${i}`} className="group-skeleton skeleton-card"></div>
            ))
          : result.data?.findGroups.groups.map((g) => (
              <GroupCard key={g.id} group={g} />
            ))}
      </Slider>
    </RecommendationRow>
  );
};
