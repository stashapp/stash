import React, { useState } from "react";
import { Button, Form } from "react-bootstrap";
import { mutateMetadataGenerate } from "src/core/StashService";
import { useToast } from "src/hooks";

export const GenerateButton: React.FC = () => {
  const Toast = useToast();
  const [sprites, setSprites] = useState(true);
  const [previews, setPreviews] = useState(true);
  const [markers, setMarkers] = useState(true);
  const [transcodes, setTranscodes] = useState(false);
  const [thumbnails, setThumbnails] = useState(false);
  const [imagePreviews, setImagePreviews] = useState(false);

  async function onGenerate() {
    try {
      await mutateMetadataGenerate({
        sprites,
        previews,
        imagePreviews: previews && imagePreviews,
        markers,
        transcodes,
        thumbnails,
      });
      Toast.success({ content: "Started generating" });
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <>
      <Form.Group>
        <Form.Check
          id="preview-task"
          checked={previews}
          label="Previews (video previews which play when hovering over a scene)"
          onChange={() => setPreviews(!previews)}
        />
        <div className="d-flex flex-row">
          <div>↳</div>
          <Form.Check
            id="image-preview-task"
            checked={imagePreviews}
            disabled={!previews}
            label="Image Previews (animated WebP previews, only required if Preview Type is set to Animated Image)"
            onChange={() => setImagePreviews(!imagePreviews)}
            className="ml-2 flex-grow"
          />
        </div>
        <Form.Check
          id="sprite-task"
          checked={sprites}
          label="Sprites (for the scene scrubber)"
          onChange={() => setSprites(!sprites)}
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
        <Form.Check
          id="thumbnail-task"
          checked={thumbnails}
          label="Gallery thumbnails (thumbnails for all the gallery images)"
          onChange={() => setThumbnails(!thumbnails)}
        />
      </Form.Group>
      <Form.Group>
        <Button
          id="generate"
          variant="secondary"
          type="submit"
          onClick={() => onGenerate()}
        >
          Generate
        </Button>
        <Form.Text className="text-muted">
          Generate supporting image, sprite, video, vtt and other files.
        </Form.Text>
      </Form.Group>
    </>
  );
};
