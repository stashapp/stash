import React from "react";
import { useFindGalleries } from "src/core/StashService";
import { GalleryCard } from "./GalleryCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const GalleryRecommendationRow: React.FC<IProps> = PatchComponent(
  "GalleryRecommendationRow",
  (props) => {
    const result = useFindGalleries(props.filter);
    const count = result.data?.findGalleries.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="gallery-recommendations"
        heading={props.header}
        url={`/galleries?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
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
      </FilteredRecommendationRow>
    );
  }
);
