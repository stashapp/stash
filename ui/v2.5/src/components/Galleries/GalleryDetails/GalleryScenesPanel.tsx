import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneCard } from "src/components/Scenes/SceneCard";

interface IGalleryScenesPanelProps {
  scenes: GQL.SlimSceneDataFragment[];
}

export const GalleryScenesPanel: React.FC<IGalleryScenesPanelProps> = ({
  scenes,
}) => (
  <div className="container gallery-scenes">
    {scenes.map((scene) => (
      <SceneCard scene={scene} key={scene.id} />
    ))}
  </div>
);
