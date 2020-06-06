import React, { useState } from "react";
import { Button, Form } from "react-bootstrap";
import { mutateMetadataGenerate } from "src/core/StashService";
import { PreviewPreset } from "src/core/generated-graphql";
import { useToast } from "src/hooks";

export const GenerateButton: React.FC = () => {
  const Toast = useToast();
  const [sprites, setSprites] = useState(true);
  const [previews, setPreviews] = useState(true);
  const [markers, setMarkers] = useState(true);
  const [transcodes, setTranscodes] = useState(false);
  const [thumbnails, setThumbnails] = useState(false);
  const [imagePreviews, setImagePreviews] = useState(false);
  const [previewPreset, setPreviewPreset] = useState<string>(
    PreviewPreset.Slow
  );

  async function onGenerate() {
    try {
      await mutateMetadataGenerate({
        sprites,
        previews,
        imagePreviews: previews && imagePreviews,
        markers,
        transcodes,
        thumbnails,
        previewPreset: (previewPreset as PreviewPreset) ?? undefined,
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
        <Form.Group controlId="preview-preset" className="mt-2">
          <Form.Label>
            <h6>Preview encoding preset</h6>
          </Form.Label>
          <Form.Control
            as="select"
            value={previewPreset}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
              setPreviewPreset(e.currentTarget.value)
            }
            disabled={!previews}
            className="col-1"
          >
            {Object.keys(PreviewPreset).map((p) => (
              <option value={p.toLowerCase()} key={p}>
                {p}
              </option>
            ))}
          </Form.Control>
          <Form.Text className="text-muted">
            The preset regulates size, quality and encoding time of preview
            generation. Presets beyond “slow” have diminishing returns and are
            not recommended.
          </Form.Text>
        </Form.Group>
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
