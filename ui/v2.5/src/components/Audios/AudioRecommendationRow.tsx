import React from "react";
import { useFindAudios } from "src/core/StashService";
import { AudioCard } from "./AudioCard";
import { ListFilterModel } from "src/models/list-filter/filter";
import { PatchComponent } from "src/patch";
import { FilteredRecommendationRow } from "../FrontPage/FilteredRecommendationRow";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const AudioRecommendationRow: React.FC<IProps> = PatchComponent(
  "AudioRecommendationRow",
  (props) => {
    const result = useFindAudios(props.filter);
    const count = result.data?.findAudios.count ?? 0;

    return (
      <FilteredRecommendationRow
        className="audio-recommendations"
        heading={props.header}
        url={`/audios?${props.filter.makeQueryParameters()}`}
        count={count}
        loading={result.loading}
        isTouch={props.isTouch}
        filter={props.filter}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div key={`_${i}`} className="audio-skeleton skeleton-card"></div>
            ))
          : result.data?.findAudios.audios.map((audio, index) => (
              <AudioCard
                key={audio.id}
                audio={audio}
                index={index}
                zoomIndex={1}
              />
            ))}
      </FilteredRecommendationRow>
    );
  }
);
