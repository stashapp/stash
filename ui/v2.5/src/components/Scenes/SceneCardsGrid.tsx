import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneCard } from "./SceneCard";

interface ISceneCardsGrid {
  scenes: GQL.SlimSceneDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  onSceneClick?: (id: string) => void;
}

export const SceneCardsGrid: React.FC<ISceneCardsGrid> = ({
  scenes,
  selectedIds,
  zoomIndex,
  onSelectChange,
  onSceneClick,
}) => {
  function sceneClicked(sceneID: string) {
    if (onSceneClick) {
      onSceneClick(sceneID);
    }
  }

  return (
    <div className="row justify-content-center">
      {scenes.map((scene) => (
        <SceneCard
          key={scene.id}
          scene={scene}
          zoomIndex={zoomIndex}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(scene.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(scene.id, selected, shiftKey)
          }
          onSceneClicked={() => sceneClicked(scene.id)}
        />
      ))}
    </div>
  );
};
