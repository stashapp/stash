import React, { MouseEvent } from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
import { WallItem, WallItemData, WallItemType } from "./WallItem";

interface IWallPanelProps<T extends WallItemType> {
  type: T;
  data: WallItemData[T][];
  sceneQueue?: SceneQueue;
  clickHandler?: (e: MouseEvent, item: WallItemData[T]) => void;
}

const calculateClass = (index: number, count: number) => {
  // First position and more than one row
  if (index === 0 && count > 5) return "transform-origin-top-left";
  // Fifth position and more than one row
  if (index === 4 && count > 5) return "transform-origin-top-right";
  // Top row
  if (index < 5) return "transform-origin-top";
  // Two or more rows, with full last row and index is last
  if (count > 9 && count % 5 === 0 && index + 1 === count)
    return "transform-origin-bottom-right";
  // Two or more rows, with full last row and index is fifth to last
  if (count > 9 && count % 5 === 0 && index + 5 === count)
    return "transform-origin-bottom-left";
  // Multiple of five minus one
  if (index % 5 === 4) return "transform-origin-right";
  // Multiple of five
  if (index % 5 === 0) return "transform-origin-left";
  // Position is equal or larger than first position in last row
  if (count - (count % 5 || 5) <= index + 1) return "transform-origin-bottom";
  // Default
  return "transform-origin-center";
};

const WallPanel = <T extends WallItemType>({
  type,
  data,
  sceneQueue,
  clickHandler,
}: IWallPanelProps<T>) => {
  function renderItems() {
    return data.map((item, index, arr) => (
      <WallItem
        type={type}
        key={item.id}
        index={index}
        data={item}
        sceneQueue={sceneQueue}
        clickHandler={clickHandler}
        className={calculateClass(index, arr.length)}
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
