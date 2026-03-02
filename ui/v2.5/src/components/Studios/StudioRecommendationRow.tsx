import React from "react";
import { useFindStudios } from "src/core/StashService";
import { StudioCard } from "./StudioCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const StudioRecommendationRow: React.FC<IProps> = PatchComponent(
  "StudioRecommendationRow",
  (props) => {
    const result = useFindStudios(props.filter);
    const count = result.data?.findStudios.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="studio-recommendations"
        heading={props.header}
        url={`/studios?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
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
      </FilteredRecommendationRow>
    );
  }
);
