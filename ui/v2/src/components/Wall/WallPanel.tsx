import _ from "lodash";
import React, { FunctionComponent, useState, useEffect } from "react";
import * as GQL from "../../core/generated-graphql";
import "./Wall.scss";
import { WallItem, IWallItemPosition } from "./WallItem";
import justifiedLayout from "justified-layout";

interface IWallPanelProps {
  scenes?: GQL.SlimSceneDataFragment[];
  sceneMarkers?: GQL.SceneMarkerDataFragment[];
  clickHandler?: (item: GQL.SlimSceneDataFragment | GQL.SceneMarkerDataFragment) => void;
}

export const WallPanel: FunctionComponent<IWallPanelProps> = (props: IWallPanelProps) => {
  const [showOverlay, setShowOverlay] = useState<boolean>(false);
  const [wallItemPositions, setWallItemPositions] = useState<IWallItemPosition[]>([]);

  useEffect(() => {
    if (!props.scenes) {
      setWallItemPositions([]);
    } else {
      // need to get all of the aspect ratios
      var w = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
      w = w * 0.8;

      let newSceneAspectRatios = props.scenes.map((scene) => {
        const defaultAspectRatio = 4 / 3;
        if (!scene.file.width || !scene.file.height) {
          return defaultAspectRatio;
        }
        return scene.file.width / scene.file.height;
      });

      const rowHeight = 290;
      const heightTolerance = 0.1;
      let layoutGeo = justifiedLayout(newSceneAspectRatios, {
        containerWidth: w, 
        targetRowHeight: rowHeight, 
        targetRowHeightTolerance: heightTolerance,
        boxSpacing: { horizontal: 0, vertical: 0 }
      });

      setWallItemPositions(layoutGeo.boxes as IWallItemPosition[]);
    }
  }, [props.scenes]);

  function onOverlay(show: boolean) {
    setShowOverlay(show);
  }

  function getOrigin(index: number, rowSize: number, total: number): string {
    const isAtStart = index % rowSize === 0;
    const isAtEnd = index % rowSize === rowSize - 1;
    const endRemaining = total % rowSize;

    // First row
    if (total === 1) { return "top"; }
    if (index === 0) { return "top left"; }
    if (index === rowSize - 1 || (total < rowSize && index === total - 1)) { return "top right"; }
    if (index < rowSize) { return "top"; }

    // Bottom row
    if (isAtEnd && index === total - 1) { return "bottom right"; }
    if (isAtStart && index === total - rowSize) { return "bottom left"; }
    if (endRemaining !== 0 && index >= total - endRemaining) { return "bottom"; }
    if (endRemaining === 0 && index >= total - rowSize) { return "bottom"; }

    // Everything else
    if (isAtStart) { return "center left"; }
    if (isAtEnd) { return "center right"; }
    return "center";
  }

  function maybeRenderScenes() {
    if (props.scenes === undefined) { return; }
    return props.scenes.map((scene, index) => {
      const origin = getOrigin(index, 5, props.scenes!.length);
      return (
        <WallItem
          key={scene.id}
          scene={scene}
          onOverlay={onOverlay}
          clickHandler={props.clickHandler}
          origin={origin}
          position={wallItemPositions[index]}
        />
      );
    });
  }

  function maybeRenderSceneMarkers() {
    if (props.sceneMarkers === undefined) { return; }
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
