import React, { FunctionComponent, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { WallItem } from "./WallItem";
import "./Wall.scss";

interface IWallPanelProps {
  scenes?: GQL.SlimSceneDataFragment[];
  sceneMarkers?: GQL.SceneMarkerDataFragment[];
  clickHandler?: (
    item: GQL.SlimSceneDataFragment | GQL.SceneMarkerDataFragment
  ) => void;
}

export const WallPanel: FunctionComponent<IWallPanelProps> = (
  props: IWallPanelProps
) => {
  const [showOverlay, setShowOverlay] = useState<boolean>(false);

  function onOverlay(show: boolean) {
    setShowOverlay(show);
  }

  function getOrigin(index: number, rowSize: number, total: number): string {
    const isAtStart = index % rowSize === 0;
    const isAtEnd = index % rowSize === rowSize - 1;
    const endRemaining = total % rowSize;

    // First row
    if (total === 1) {
      return "top";
    }
    if (index === 0) {
      return "top left";
    }
    if (index === rowSize - 1 || (total < rowSize && index === total - 1)) {
      return "top right";
    }
    if (index < rowSize) {
      return "top";
    }

    // Bottom row
    if (isAtEnd && index === total - 1) {
      return "bottom right";
    }
    if (isAtStart && index === total - rowSize) {
      return "bottom left";
    }
    if (endRemaining !== 0 && index >= total - endRemaining) {
      return "bottom";
    }
    if (endRemaining === 0 && index >= total - rowSize) {
      return "bottom";
    }

    // Everything else
    if (isAtStart) {
      return "center left";
    }
    if (isAtEnd) {
      return "center right";
    }
    return "center";
  }

  function maybeRenderScenes() {
    if (props.scenes === undefined) {
      return;
    }
    return props.scenes.map((scene, index) => {
      const origin = getOrigin(index, 5, props.scenes!.length);
      return (
        <WallItem
          key={scene.id}
          scene={scene}
          onOverlay={onOverlay}
          clickHandler={props.clickHandler}
          origin={origin}
        />
      );
    });
  }

  function maybeRenderSceneMarkers() {
    if (props.sceneMarkers === undefined) {
      return;
    }
    return props.sceneMarkers.map((marker, index) => {
      const origin = getOrigin(index, 5, props.sceneMarkers!.length);
      return (
        <WallItem
          key={marker.id}
          sceneMarker={marker}
          onOverlay={onOverlay}
          clickHandler={props.clickHandler}
          origin={origin}
        />
      );
    });
  }

  function render() {
    const overlayClassName = showOverlay ? "visible" : "hidden";
    return (
      <>
        <div className={`wall-overlay ${overlayClassName}`} />
        <div className="wall grid">
          {maybeRenderScenes()}
          {maybeRenderSceneMarkers()}
        </div>
      </>
    );
  }

  return render();
};
