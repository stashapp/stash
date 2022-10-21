import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneCard } from "src/components/Scenes/SceneCard";
import { ConfigurationContext } from "src/hooks/Config";

interface IGalleryScenesPanelProps {
  scenes: GQL.SlimSceneDataFragment[];
}

export const GalleryScenesPanel: React.FC<IGalleryScenesPanelProps> = ({
  scenes,
}) => {
  const [touchPreviewActive, setTouchPreviewActive] = React.useState("");
  const { isTouch } = React.useContext(ConfigurationContext);

  return (
    <div className="container gallery-scenes">
      {scenes.map((scene) => (
        <SceneCard
          scene={scene}
          key={scene.id}
          isTouchPreviewActive={isTouch && touchPreviewActive === scene.id}
          onTouchPreview={() => {
            setTouchPreviewActive(scene.id);
          }}
        />
      ))}
    </div>
  );
};
