import React from "react";
import { Link } from "react-router-dom";
import { useFindSceneMarkers } from "src/core/StashService";
import Slider from "@ant-design/react-slick";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { FormattedMessage } from "react-intl";
import { SceneMarkerCard } from "./SceneMarkerCard";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  header: string;
}

export const SceneMarkerRecommendationRow: React.FC<IProps> = (props) => {
  const result = useFindSceneMarkers(props.filter);
  const cardCount = result.data?.findSceneMarkers.count;

  if (!result.loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="scene-marker-recommendations"
      header={props.header}
      link={
        <Link to={`/scenes/markers?${props.filter.makeQueryParameters()}`}>
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
      </Slider>
    </RecommendationRow>
  );
};
