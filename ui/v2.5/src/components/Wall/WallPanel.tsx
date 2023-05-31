import React, { MouseEvent } from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
import { WallItem, WallItemData, WallItemType } from "./WallItem";

interface IWallPanelProps<T extends WallItemType> {
  type: T;
  data: WallItemData[T][];
  sceneQueue?: SceneQueue;
  zoomIndex?: number;
  clickHandler?: (e: MouseEvent, item: WallItemData[T]) => void;
}

const calculateClass = (index: number, count: number, columns: number) => {
  // No Classname if only 1 column
  if (columns === 1) {
    return "no-transform";
  }
  const lastIndex = count - 1;
  const lastRowIndex = Math.floor(lastIndex / columns) * columns;
  
  // First position and more than one row
  if (index === 0 && count > columns) return "transform-origin-top-left";
  // MaxColumn+1 position and more than one row
  if (index === columns - 1 && count > columns) return "transform-origin-top-right";
  // Top row
  if (index < columns) return "transform-origin-top";
  // Two or more rows, with full last row and index is last
  if (count > lastRowIndex + columns && (index + 1) === lastIndex) return "transform-origin-bottom-right";
  // Two or more rows, with full last row and index is fifth to last
  if (count > lastRowIndex + columns && (index + columns) === lastIndex) return "transform-origin-bottom-left";
  // Multiple of columns minus one
  if ((index + 1) % columns === 0) return "transform-origin-right";
  // Multiple of columns
  if (index % columns === 0) {
    // Position is equal or larger than first position in last row
    if (lastIndex - (lastIndex % columns || columns) <= index) return "transform-origin-bottom";
    return "transform-origin-left";
  }
  // Position is equal or larger than first position in last row
  if (lastIndex - (lastIndex % columns || columns) <= index + 1) return "transform-origin-bottom";
  // Default
  return "transform-origin-center";
};

const WallPanel = <T extends WallItemType>({
  type,
  data,
  sceneQueue,
  clickHandler,
  zoomIndex,
}: IWallPanelProps<T>) => {
  function renderItems() {
    let column = 5 - (zoomIndex ?? 0);
    return data.map((item, index, arr) => (
      <WallItem
        type={type}
        key={item.id}
        index={index}
        data={item}
        zoomIndex={zoomIndex}
        sceneQueue={sceneQueue}
        clickHandler={clickHandler}
        className={calculateClass(index, arr.length, column)}
      />
    ));
  }

  return (
    <div className="row">
      <div className="wall w-100 row justify-content-center">
        {renderItems()}
      </div>
    </div>
  );
};

interface IImageWallPanelProps {
  images: GQL.SlimImageDataFragment[];
  clickHandler?: (e: MouseEvent, item: GQL.SlimImageDataFragment) => void;
}

export const ImageWallPanel: React.FC<IImageWallPanelProps> = ({
  images,
  clickHandler,
}) => {
  return <WallPanel type="image" data={images} clickHandler={clickHandler} />;
};

interface IMarkerWallPanelProps {
  markers: GQL.SceneMarkerDataFragment[];
  clickHandler?: (e: MouseEvent, item: GQL.SceneMarkerDataFragment) => void;
}

export const MarkerWallPanel: React.FC<IMarkerWallPanelProps> = ({
  markers,
  clickHandler,
}) => {
  return (
    <WallPanel type="sceneMarker" data={markers} clickHandler={clickHandler} />
  );
};

interface ISceneWallPanelProps {
  scenes: GQL.SlimSceneDataFragment[];
  sceneQueue?: SceneQueue;
  zoomIndex?: number;
  clickHandler?: (e: MouseEvent, item: GQL.SlimSceneDataFragment) => void;
}

export const SceneWallPanel: React.FC<ISceneWallPanelProps> = ({
  scenes,
  sceneQueue,
  clickHandler,
  zoomIndex,
}) => {
  return (
    <WallPanel
      type="scene"
      data={scenes}
      sceneQueue={sceneQueue}
      zoomIndex={zoomIndex}
      clickHandler={clickHandler}
    />
  );
};
