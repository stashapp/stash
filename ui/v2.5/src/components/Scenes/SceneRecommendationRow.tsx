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
  loadingArr: boolean[];
  index: number;
}

export const SceneRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const result = useFindScenes(props.filter);
  const cardCount = result.data?.findScenes.count;
  props.loadingArr[props.index] = result.loading;
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
    </div>
  );
};
