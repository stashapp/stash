import React, { useState } from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { mutateMetadataGenerate } from "src/core/StashService";
import { useToast } from "src/hooks";

export const GenerateButton: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();
  const [sprites, setSprites] = useState(true);
  const [phashes, setPhashes] = useState(true);
  const [previews, setPreviews] = useState(true);
  const [markers, setMarkers] = useState(true);
  const [transcodes, setTranscodes] = useState(false);
  const [imagePreviews, setImagePreviews] = useState(false);
  const [markerImagePreviews, setMarkerImagePreviews] = useState(false);
  const [markerScreenshots, setMarkerScreenshots] = useState(false);

  async function onGenerate() {
    try {
      await mutateMetadataGenerate({
        sprites,
        phashes,
        previews,
        imagePreviews: previews && imagePreviews,
        markers,
        markerImagePreviews: markers && markerImagePreviews,
        markerScreenshots: markers && markerScreenshots,
        transcodes,
      });
      Toast.success({
        content: intl.formatMessage({
          id: "toast.added_generation_job_to_queue",
        }),
      });
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
          label={intl.formatMessage({ id: "dialogs.scene_gen.video_previews" })}
          onChange={() => setPreviews(!previews)}
        />
        <div className="d-flex flex-row">
          <div>↳</div>
          <Form.Check
            id="image-preview-task"
            checked={imagePreviews}
            disabled={!previews}
            label={intl.formatMessage({
              id: "dialogs.scene_gen.image_previews",
            })}
            onChange={() => setImagePreviews(!imagePreviews)}
            className="ml-2 flex-grow"
          />
        </div>
        <Form.Check
          id="sprite-task"
          checked={sprites}
          label={intl.formatMessage({ id: "dialogs.scene_gen.sprites" })}
          onChange={() => setSprites(!sprites)}
        />
        <Form.Check
          id="marker-task"
          checked={markers}
          label={intl.formatMessage({ id: "dialogs.scene_gen.markers" })}
          onChange={() => setMarkers(!markers)}
        />
        <div className="d-flex flex-row">
          <div>↳</div>
          <Form.Group>
            <Form.Check
              id="marker-image-preview-task"
              checked={markerImagePreviews}
              disabled={!markers}
              label={intl.formatMessage({
                id: "dialogs.scene_gen.marker_image_previews",
              })}
              onChange={() => setMarkerImagePreviews(!markerImagePreviews)}
              className="ml-2 flex-grow"
            />
            <Form.Check
              id="marker-screenshot-task"
              checked={markerScreenshots}
              disabled={!markers}
              label={intl.formatMessage({
                id: "dialogs.scene_gen.marker_screenshots",
              })}
              onChange={() => setMarkerScreenshots(!markerScreenshots)}
              className="ml-2 flex-grow"
            />
          </Form.Group>
        </div>
        <Form.Check
          id="transcode-task"
          checked={transcodes}
          label={intl.formatMessage({ id: "dialogs.scene_gen.transcodes" })}
          onChange={() => setTranscodes(!transcodes)}
        />
        <Form.Check
          id="phash-task"
          checked={phashes}
          label={intl.formatMessage({ id: "dialogs.scene_gen.phash" })}
          onChange={() => setPhashes(!phashes)}
        />
      </Form.Group>
      <Form.Group>
        <Button
          id="generate"
          variant="secondary"
          type="submit"
          onClick={() => onGenerate()}
        >
          <FormattedMessage id="actions.generate" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.generate_desc" })}
        </Form.Text>
      </Form.Group>
    </>
  );
};
