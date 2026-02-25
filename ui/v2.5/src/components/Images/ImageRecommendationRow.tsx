import React from "react";
import { useFindImages } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { ImageCard } from "./ImageCard";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const ImageRecommendationRow: React.FC<IProps> = PatchComponent(
  "ImageRecommendationRow",
  (props: IProps) => {
    const result = useFindImages(props.filter);
    const count = result.data?.findImages.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="images-recommendations"
        heading={props.header}
        url={`/images?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div key={`_${i}`} className="image-skeleton skeleton-card"></div>
            ))
          : result.data?.findImages.images.map((i) => (
              <ImageCard key={i.id} image={i} zoomIndex={1} />
            ))}
      </FilteredRecommendationRow>
    );
  }
);
