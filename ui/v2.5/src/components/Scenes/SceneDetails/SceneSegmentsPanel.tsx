import React, { useState, useEffect } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { SceneSegmentForm } from "./SceneSegmentForm";

interface ISceneSegmentsPanelProps {
  sceneId: string;
  isVisible: boolean;
}

export const SceneSegmentsPanel: React.FC<ISceneSegmentsPanelProps> = ({
  sceneId,
  isVisible,
}) => {
  const { data, loading, refetch } = GQL.useFindSceneSegmentsQuery({
    variables: { scene_id: sceneId },
  });
  const [isEditorOpen, setIsEditorOpen] = useState<boolean>(false);
  const [editingSegment, setEditingSegment] =
    useState<GQL.SceneSegmentDataFragment>();

  // set up hotkeys
  useEffect(() => {
    if (!isVisible) return;

    Mousetrap.bind("s", () => onOpenEditor());

    return () => {
      Mousetrap.unbind("s");
    };
  });

  if (loading) return null;

  function onOpenEditor(segment?: GQL.SceneSegmentDataFragment) {
    setIsEditorOpen(true);
    setEditingSegment(segment ?? undefined);
  }

  const closeEditor = () => {
    setEditingSegment(undefined);
    setIsEditorOpen(false);
    refetch();
  };

  if (isEditorOpen)
    return (
      <SceneSegmentForm
        sceneID={sceneId}
        segment={editingSegment}
        onClose={closeEditor}
      />
    );

  const segments = data?.findSceneSegments ?? [];

  return (
    <div className="scene-segments-panel">
      <h4>
        <FormattedMessage id="scene.segments" defaultMessage="Segments" />
      </h4>
      <Button onClick={() => onOpenEditor()} className="mb-3">
        <FormattedMessage
          id="actions.create_segment"
          defaultMessage="Create Segment"
        />
      </Button>

      {segments.length === 0 ? (
        <div className="text-muted">
          <FormattedMessage
            id="scene.no_segments"
            defaultMessage="No segments defined for this scene"
          />
        </div>
      ) : (
        <div className="segments-list">
          <table className="table table-sm">
            <thead>
              <tr>
                <th>Title</th>
                <th>Start</th>
                <th>End</th>
                <th>Duration</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {segments.map((segment) => {
                const duration = segment.end_seconds - segment.start_seconds;
                return (
                  <tr key={segment.id}>
                    <td>{segment.title}</td>
                    <td>{formatTime(segment.start_seconds)}</td>
                    <td>{formatTime(segment.end_seconds)}</td>
                    <td>{formatTime(duration)}</td>
                    <td>
                      <Button
                        size="sm"
                        variant="secondary"
                        onClick={() => onOpenEditor(segment)}
                        className="mr-2"
                      >
                        Edit
                      </Button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, "0")}`;
}

export default SceneSegmentsPanel;
