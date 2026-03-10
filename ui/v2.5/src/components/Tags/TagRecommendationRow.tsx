import React from "react";
import { useFindTags } from "src/core/StashService";
import { TagCard } from "./TagCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const TagRecommendationRow: React.FC<IProps> = PatchComponent(
  "TagRecommendationRow",
  (props) => {
    const result = useFindTags(props.filter);
    const count = result.data?.findTags.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="tag-recommendations"
        heading={props.header}
        url={`/tags?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div key={`_${i}`} className="tag-skeleton skeleton-card"></div>
            ))
          : result.data?.findTags.tags.map((p) => (
              <TagCard key={p.id} tag={p} zoomIndex={0} />
            ))}
      </FilteredRecommendationRow>
    );
  }
);
