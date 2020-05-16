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

export const WallPanel: React.FC<IWallPanelProps> = (
  props: IWallPanelProps
) => {
  const scenes = (props.scenes ?? []).map(scene => (
    <WallItem
      key={scene.id}
      scene={scene}
      clickHandler={props.clickHandler}
    />
  ));

  const sceneMarkers = (props.sceneMarkers ?? []).map(marker => (
    <WallItem
      key={marker.id}
      sceneMarker={marker}
      clickHandler={props.clickHandler}
    />
  ));

  return (
    <div className="row">
      <div className="wall row justify-content-center">
        {scenes}
        {sceneMarkers}
      </div>
    </div>
  );
};
