import React from "react";
import { useFindSceneMarkers } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SceneMarkerCard } from "./SceneMarkerCard";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const SceneMarkerRecommendationRow: React.FC<IProps> = PatchComponent(
  "SceneMarkerRecommendationRow",
  (props) => {
    const result = useFindSceneMarkers(props.filter);
    const count = result.data?.findSceneMarkers.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="scene-marker-recommendations"
        heading={props.header}
        url={`/scenes/markers?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div
                key={`_${i}`}
                className="scene-marker-skeleton skeleton-card"
              ></div>
            ))
          : result.data?.findSceneMarkers.scene_markers.map((marker, index) => (
              <SceneMarkerCard
                key={marker.id}
                marker={marker}
                index={index}
                zoomIndex={1}
              />
            ))}
      </FilteredRecommendationRow>
    );
  }
);
