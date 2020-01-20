import React, { useState } from "react";
import { Button, Form } from "react-bootstrap";
import { StashService } from "src/core/StashService";
import { useToast } from "src/hooks";

export const GenerateButton: React.FC = () => {
  const Toast = useToast();
  const [sprites, setSprites] = useState(true);
  const [previews, setPreviews] = useState(true);
  const [markers, setMarkers] = useState(true);
  const [transcodes, setTranscodes] = useState(true);

  async function onGenerate() {
    try {
      await StashService.queryMetadataGenerate({
        sprites,
        previews,
        markers,
        transcodes
      });
      Toast.success({ content: "Started generating" });
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <Form.Group>
      <Form.Check
        id="sprite-task"
        checked={sprites}
        label="Sprites (for the scene scrubber)"
        onChange={() => setSprites(!sprites)}
      />
      <Form.Check
        id="preview-task"
        checked={previews}
        label="Previews (video previews which play when hovering over a scene)"
        onChange={() => setPreviews(!previews)}
      />
      <Form.Check
        id="marker-task"
        checked={markers}
        label="Markers (20 second videos which begin at the given timecode)"
        onChange={() => setMarkers(!markers)}
      />
      <Form.Check
        id="transcode-task"
        checked={transcodes}
        label="Transcodes (MP4 conversions of unsupported video formats)"
        onChange={() => setTranscodes(!transcodes)}
      />
      <Button id="generate" type="submit" onClick={() => onGenerate()}>
        Generate
      </Button>
      <Form.Text className="text-muted">
        Generate supporting image, sprite, video, vtt and other files.
      </Form.Text>
    </Form.Group>
  );
};
