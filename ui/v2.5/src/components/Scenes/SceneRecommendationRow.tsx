import React, { FunctionComponent } from "react";
import { useFindScenes } from "src/core/StashService";
import Slider from "react-slick";
import { SceneCard } from "./SceneCard";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  queue: SceneQueue;
  header: String;
  linkText: String;
  index: number;
}

export const SceneRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  function buildRow(slider: JSX.Element) {
    return (
      <div className="recommendation-row scene-recommendations">
        <div className="recommendation-row-head">
          <div>
            <h2>{props.header}</h2>
          </div>
          <a href={`/scenes?${props.filter.makeQueryParameters()}`}>
            {props.linkText}
          </a>
        </div>
        {slider}
      </div>
    );
  }

  const result = useFindScenes(props.filter);
  const cardCount = result.data?.findScenes.count;
  if (result.loading) {
    const slider = (
      <Slider
        {...getSlickSliderSettings(props.filter.itemsPerPage!, props.isTouch)}
      >
        {[...Array(props.filter.itemsPerPage)].map((i) => (
          <div key={i} className="scene-skeleton skeleton-card"></div>
        ))}
      </Slider>
    );
    return buildRow(slider);
  }

  if (cardCount === 0) {
    return null;
  }

  const slider = (
    <Slider {...getSlickSliderSettings(cardCount!, props.isTouch)}>
      {result.data?.findScenes.scenes.map((scene, index) => (
        <SceneCard
          key={scene.id}
          scene={scene}
          queue={props.queue}
          index={index}
          zoomIndex={1}
        />
      ))}
    </Slider>
  );
  return buildRow(slider);
};
