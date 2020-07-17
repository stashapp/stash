import React, { useState } from "react";
import { Form, Button } from "react-bootstrap";
import { mutateMetadataGenerate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";

interface ISceneGenerateDialogProps {
  selectedIds: string[];
  onClose: () => void;
}

export const SceneGenerateDialog: React.FC<ISceneGenerateDialogProps> = (
  props: ISceneGenerateDialogProps
) => {
  const [sprites, setSprites] = useState(true);
  const [previews, setPreviews] = useState(true);
  const [markers, setMarkers] = useState(true);
  const [transcodes, setTranscodes] = useState(false);
  const [overwrite, setOverwrite] = useState(true);
  const [imagePreviews, setImagePreviews] = useState(false);
  const [previewPreset, setPreviewPreset] = useState<string>(
    GQL.PreviewPreset.Slow
  );

  const Toast = useToast();

  async function onGenerate() {
    try {
      await mutateMetadataGenerate({
        sprites,
        previews,
        imagePreviews: previews && imagePreviews,
        markers,
        transcodes,
        thumbnails: false,
        previewPreset: (previewPreset as GQL.PreviewPreset) ?? undefined,
        overwrite,
        sceneIDs: props.selectedIds,
      });
      Toast.success({ content: "Started generating" });
    } catch (e) {
      Toast.error(e);
    } finally {
      props.onClose();
    }
  }
  
  return (
    <Modal
      show
      icon="cogs"
      header="Generate"
      accept={{ onClick: onGenerate, text: "Generate" }}
      cancel={{
        onClick: () => props.onClose(),
        text: "Cancel",
        variant: "secondary",
      }}
    >
      <Form>
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
                className="w-auto"
                as="select"
                value={previewPreset}
                onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
                setPreviewPreset(e.currentTarget.value)
                }
                disabled={!previews}
            >
                {Object.keys(GQL.PreviewPreset).map((p) => (
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

            <hr />
            <Form.Check
            id="overwrite"
            checked={overwrite}
            label="Overwrite existing generated files"
            onChange={() => setOverwrite(!overwrite)}
            />
        </Form.Group>
      </Form>
    </Modal>
  );
};
