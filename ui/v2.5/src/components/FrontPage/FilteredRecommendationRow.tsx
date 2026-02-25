import React from "react";
import { Link } from "react-router-dom";
import Slider from "@ant-design/react-slick";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";
import { PatchComponent } from "src/patch";
import { UnsupportedCriterion } from "src/models/list-filter/criteria/criterion";
import { PopoverCard, WarningHoverPopover } from "../Shared/HoverPopover";

interface IProps {
  className?: string;
  isTouch: boolean;
  filter: ListFilterModel;
  heading: string;
  count: number;
  loading: boolean;
  url: string;
}

export const FilteredRecommendationRow: React.FC<IProps> = PatchComponent(
  "FilteredRecommendationRow",
  (props) => {
    const cardCount = props.count;

    const unsupportedCriteria = props.filter.criteria.filter(
      (criterion) => criterion instanceof UnsupportedCriterion
    );

    const header = unsupportedCriteria.length ? (
      <div>
        <span>{props.heading}</span>
        <WarningHoverPopover
          placement="top"
          content={
            <PopoverCard>
              <FormattedMessage
                id="unsupported_criteria"
                values={{
                  criteria: unsupportedCriteria
                    .map((c) => c.criterionOption.type)
                    .join(", "),
                }}
              />
            </PopoverCard>
          }
        />
      </div>
    ) : (
      props.heading
    );

    if (!props.loading && !cardCount) {
      return null;
    }

    return (
      <RecommendationRow
        className={props.className}
        header={header}
        link={
          <Link to={props.url}>
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
          {props.children}
        </Slider>
      </RecommendationRow>
    );
  }
);
