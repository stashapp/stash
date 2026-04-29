import React from "react";
import { useFindGroups } from "src/core/StashService";
import { GroupCard } from "./GroupCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const GroupRecommendationRow: React.FC<IProps> = PatchComponent(
  "GroupRecommendationRow",
  (props: IProps) => {
    const result = useFindGroups(props.filter);
    const count = result.data?.findGroups.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="group-recommendations"
        heading={props.header}
        url={`/groups?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div key={`_${i}`} className="group-skeleton skeleton-card"></div>
            ))
          : result.data?.findGroups.groups.map((g) => (
              <GroupCard key={g.id} group={g} />
            ))}
      </FilteredRecommendationRow>
    );
  }
);
