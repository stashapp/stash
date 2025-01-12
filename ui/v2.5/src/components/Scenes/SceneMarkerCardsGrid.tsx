import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneMarkerCard } from "./SceneMarkerCard";
import { useContainerDimensions } from "../Shared/GridCard/GridCard";
import { SceneQueue } from "src/models/sceneQueue";

interface ISceneMarkerCardsGrid {
  markers: GQL.SceneMarkerDataFragment[];
  queue?: SceneQueue;
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const SceneMarkerCardsGrid: React.FC<ISceneMarkerCardsGrid> = ({
  markers,
  queue,
  selectedIds,
  zoomIndex,
  onSelectChange,
}) => {
  const [componentRef, { width }] = useContainerDimensions();
  return (
    <div className="row justify-content-center" ref={componentRef}>
      {markers.map((marker, index) => (
        <SceneMarkerCard
          key={marker.id}
          containerWidth={width}
          marker={marker}
          queue={queue}
          index={index}
          zoomIndex={zoomIndex}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(marker.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(marker.id, selected, shiftKey)
          }
        />
      ))}
    </div>
  );
};
