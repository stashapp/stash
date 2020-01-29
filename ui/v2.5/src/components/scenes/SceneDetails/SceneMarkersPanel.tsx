import React, { useState } from "react";
import { Button } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { WallPanel } from "src/components/Wall/WallPanel";
import { JWUtils } from "src/utils";
import { PrimaryTags } from "./PrimaryTags";
import { SceneMarkerForm } from "./SceneMarkerForm";

interface ISceneMarkersPanelProps {
  scene: GQL.SceneDataFragment;
  onClickMarker: (marker: GQL.SceneMarkerDataFragment) => void;
}

export const SceneMarkersPanel: React.FC<ISceneMarkersPanelProps> = (
  props: ISceneMarkersPanelProps
) => {
  const [isEditorOpen, setIsEditorOpen] = useState<boolean>(false);
  const [editingMarker, setEditingMarker] = useState<
    GQL.SceneMarkerDataFragment
  >();

  const jwplayer = JWUtils.getPlayer();

  function onOpenEditor(marker?: GQL.SceneMarkerDataFragment) {
    setIsEditorOpen(true);
    setEditingMarker(marker ?? undefined);
  }

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    props.onClickMarker(marker);
  }

  const closeEditor = () => {
    setEditingMarker(undefined);
    setIsEditorOpen(false);
  };

  if (isEditorOpen)
    return (
      <SceneMarkerForm
        sceneID={props.scene.id}
        editingMarker={editingMarker}
        playerPosition={jwplayer.getPlayer?.().playerPosition}
        onClose={closeEditor}
      />
    );

  return (
    <>
      <Button onClick={() => onOpenEditor()}>Create Marker</Button>
      <div className="container">
        <PrimaryTags
          sceneMarkers={props.scene.scene_markers ?? []}
          onClickMarker={onClickMarker}
          onEdit={onOpenEditor}
        />
      </div>
      <div className="row">
        <WallPanel
          sceneMarkers={props.scene.scene_markers}
          clickHandler={marker => {
            window.scrollTo(0, 0);
            onClickMarker(marker as any);
          }}
        />
      </div>
    </>
  );
};
