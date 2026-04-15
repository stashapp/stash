import React from "react";
import { useFindPerformers } from "src/core/StashService";
import { PerformerCard } from "./PerformerCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const PerformerRecommendationRow: React.FC<IProps> = PatchComponent(
  "PerformerRecommendationRow",
  (props) => {
    const result = useFindPerformers(props.filter);
    const count = result.data?.findPerformers.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="performer-recommendations"
        heading={props.header}
        url={`/performers?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div
                key={`_${i}`}
                className="performer-skeleton skeleton-card"
              ></div>
            ))
          : result.data?.findPerformers.performers.map((p) => (
              <PerformerCard key={p.id} performer={p} />
            ))}
      </FilteredRecommendationRow>
    );
  }
);
