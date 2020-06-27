import React, { useState, useEffect } from "react";
import { Button } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { WallPanel } from "src/components/Wall/WallPanel";
import { PrimaryTags } from "./PrimaryTags";
import { SceneMarkerForm } from "./SceneMarkerForm";

interface ISceneMarkersPanelProps {
  scene: GQL.SceneDataFragment;
  isVisible: boolean;
  onClickMarker: (marker: GQL.SceneMarkerDataFragment) => void;
}

export const SceneMarkersPanel: React.FC<ISceneMarkersPanelProps> = (
  props: ISceneMarkersPanelProps
) => {
  const [isEditorOpen, setIsEditorOpen] = useState<boolean>(false);
  const [editingMarker, setEditingMarker] = useState<
    GQL.SceneMarkerDataFragment
  >();

  // set up hotkeys
  useEffect(() => {
    if (props.isVisible) {
      Mousetrap.bind("n", () => onOpenEditor());

      return () => {
        Mousetrap.unbind("n");
      };
    }
  });

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
          clickHandler={(marker) => {
            window.scrollTo(0, 0);
            onClickMarker(marker as GQL.SceneMarkerDataFragment);
          }}
        />
      </div>
    </>
  );
};
