import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { useFindScenes } from "src/core/StashService";
import Slider from "@ant-design/react-slick";
import { SceneCard } from "./SceneCard";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const SceneRecommendationRow: React.FC<IProps> = (props) => {
  const result = useFindScenes(props.filter);
  const cardCount = result.data?.findScenes.count;

  const queue = useMemo(() => {
    return SceneQueue.fromListFilterModel(props.filter);
  }, [props.filter]);

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="scene-recommendations"
      header={props.header}
      link={
        <Link to={`/scenes?${props.filter.makeQueryParameters()}`}>
          <FormattedMessage id="view_all" />
        </Link>
      }
    >
      <Slider
        {...getSlickSliderSettings(
          cardCount ? cardCount : props.filter.itemsPerPage,
          props.isTouch
        )}
      >
        {result.loading
          ? [...Array(props.filter.itemsPerPage)].map((i) => (
              <div key={`_${i}`} className="scene-skeleton skeleton-card"></div>
            ))
          : result.data?.findScenes.scenes.map((scene, index) => (
              <SceneCard
                key={scene.id}
                scene={scene}
                queue={queue}
                index={index}
                zoomIndex={1}
              />
            ))}
      </Slider>
    </RecommendationRow>
  );
};
