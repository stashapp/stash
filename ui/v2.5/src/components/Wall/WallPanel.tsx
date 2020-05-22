import React from "react";
import * as GQL from "src/core/generated-graphql";
import { WallItem } from "./WallItem";

interface IWallPanelProps {
  scenes?: GQL.SlimSceneDataFragment[];
  sceneMarkers?: GQL.SceneMarkerDataFragment[];
  clickHandler?: (
    item: GQL.SlimSceneDataFragment | GQL.SceneMarkerDataFragment
  ) => void;
}

const calculateClass = (index: number, count: number) => {
  if (index === 0 && count > 5) return "transform-origin-top-left";
  if (index === 4 && count > 5) return "transform-origin-top-right";
  if (index < 5) return "transform-origin-top";
  if (count > 9 && count % 5 === 0 && index + 1 === count)
    return "transform-origin-bottom-right";
  if (count > 9 && count % 5 === 0 && index + 5 === count)
    return "transform-origin-bottom-left";
  if (count - (count % 5 || 5) <= index + 1) return "transform-origin-bottom";
  if (index % 5 === 4) return "transform-origin-right";
  if (index % 5 === 0) return "transform-origin-left";
  return "transform-origin-center";
};

export const WallPanel: React.FC<IWallPanelProps> = (
  props: IWallPanelProps
) => {
  const scenes = (props.scenes ?? []).map((scene, index, sceneArray) => (
    <WallItem
      key={scene.id}
      scene={scene}
      clickHandler={props.clickHandler}
      className={calculateClass(index, sceneArray.length)}
    />
  ));

  const sceneMarkers = (
    props.sceneMarkers ?? []
  ).map((marker, index, markerArray) => (
    <WallItem
      key={marker.id}
      sceneMarker={marker}
      clickHandler={props.clickHandler}
      className={calculateClass(index, markerArray.length)}
    />
  ));

  return (
    <div className="row">
      <div className="wall w-100 row justify-content-center">
        {scenes}
        {sceneMarkers}
      </div>
    </div>
  );
};
