import React, { FunctionComponent } from "react";
import { FindScenesQueryResult } from "src/core/generated-graphql";
import Slider from "react-slick";
import { SceneCard } from "./SceneCard";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  result: FindScenesQueryResult;
  queue: SceneQueue;
  header: String;
  linkText: String;
}

export const SceneRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  const cardCount = props.result.data?.findScenes.count;
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
        {props.result.data?.findScenes.scenes.map((scene, index) => (
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
