import React, { useState, useEffect } from "react";
import { Form, Button, Collapse } from "react-bootstrap";
import {
  mutateMetadataGenerate,
  useConfiguration,
} from "src/core/StashService";
import { Modal, Icon } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";

interface ISceneGenerateDialogProps {
  selectedIds: string[];
  onClose: () => void;
}

export const SceneGenerateDialog: React.FC<ISceneGenerateDialogProps> = (
  props: ISceneGenerateDialogProps
) => {
  const { data, error, loading } = useConfiguration();

  const [sprites, setSprites] = useState(true);
  const [previews, setPreviews] = useState(true);
  const [markers, setMarkers] = useState(true);
  const [transcodes, setTranscodes] = useState(false);
  const [overwrite, setOverwrite] = useState(true);
  const [imagePreviews, setImagePreviews] = useState(false);

  const [previewSegments, setPreviewSegments] = useState<number>(0);
  const [previewSegmentDuration, setPreviewSegmentDuration] = useState<number>(
    0
  );
  const [previewExcludeStart, setPreviewExcludeStart] = useState<
    string | undefined
  >(undefined);
  const [previewExcludeEnd, setPreviewExcludeEnd] = useState<
    string | undefined
  >(undefined);
  const [previewPreset, setPreviewPreset] = useState<string>(
    GQL.PreviewPreset.Slow
  );

  const [previewOptionsOpen, setPreviewOptionsOpen] = useState(false);

  const Toast = useToast();

  useEffect(() => {
    if (!data?.configuration) return;

    const conf = data.configuration;
    if (conf.general) {
      setPreviewSegments(conf.general.previewSegments);
      setPreviewSegmentDuration(conf.general.previewSegmentDuration);
      setPreviewExcludeStart(conf.general.previewExcludeStart);
      setPreviewExcludeEnd(conf.general.previewExcludeEnd);
      setPreviewPreset(conf.general.previewPreset);
    }
  }, [data]);

  async function onGenerate() {
    try {
      await mutateMetadataGenerate({
        sprites,
        previews,
        imagePreviews: previews && imagePreviews,
        markers,
        transcodes,
        thumbnails: false,
        overwrite,
        sceneIDs: props.selectedIds,
        previewOptions: {
          previewPreset: (previewPreset as GQL.PreviewPreset) ?? undefined,
          previewSegments,
          previewSegmentDuration,
          previewExcludeStart,
          previewExcludeEnd,
        },
      });
      Toast.success({ content: "Started generating" });
    } catch (e) {
      Toast.error(e);
    } finally {
      props.onClose();
    }
  }

  if (error) {
    Toast.error(error);
    props.onClose();
  }

  if (loading) {
    return <></>;
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
          <div className="my-2">
            <Button
              onClick={() => setPreviewOptionsOpen(!previewOptionsOpen)}
              className="minimal pl-0 no-focus"
            >
              <Icon
                icon={previewOptionsOpen ? "chevron-down" : "chevron-right"}
              />
              <span>Preview Options</span>
            </Button>
            <Collapse in={previewOptionsOpen}>
              <div>
                <Form.Group id="transcode-size">
                  <h6>Preview encoding preset</h6>
                  <Form.Control
                    className="w-auto input-control"
                    as="select"
                    value={previewPreset}
                    onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
                      setPreviewPreset(e.currentTarget.value)
                    }
                  >
                    {Object.keys(GQL.PreviewPreset).map((p) => (
                      <option value={p.toLowerCase()} key={p}>
                        {p}
                      </option>
                    ))}
                  </Form.Control>
                  <Form.Text className="text-muted">
                    The preset regulates size, quality and encoding time of
                    preview generation. Presets beyond “slow” have diminishing
                    returns and are not recommended.
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-segments">
                  <h6>Number of segments in preview</h6>
                  <Form.Control
                    className="col col-sm-6 text-input"
                    type="number"
                    value={previewSegments.toString()}
                    onInput={(e: React.FormEvent<HTMLInputElement>) =>
                      setPreviewSegments(
                        Number.parseInt(e.currentTarget.value, 10)
                      )
                    }
                  />
                  <Form.Text className="text-muted">
                    Number of segments in preview files.
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-segment-duration">
                  <h6>Preview segment duration</h6>
                  <Form.Control
                    className="col col-sm-6 text-input"
                    type="number"
                    value={previewSegmentDuration.toString()}
                    onInput={(e: React.FormEvent<HTMLInputElement>) =>
                      setPreviewSegmentDuration(
                        Number.parseFloat(e.currentTarget.value)
                      )
                    }
                  />
                  <Form.Text className="text-muted">
                    Duration of each preview segment, in seconds.
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-exclude-start">
                  <h6>Exclude start time</h6>
                  <Form.Control
                    className="col col-sm-6 text-input"
                    defaultValue={previewExcludeStart}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setPreviewExcludeStart(e.currentTarget.value)
                    }
                  />
                  <Form.Text className="text-muted">
                    Exclude the first x seconds from scene previews. This can be
                    a value in seconds, or a percentage (eg 2%) of the total
                    scene duration.
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-exclude-start">
                  <h6>Exclude end time</h6>
                  <Form.Control
                    className="col col-sm-6 text-input"
                    defaultValue={previewExcludeEnd}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setPreviewExcludeEnd(e.currentTarget.value)
                    }
                  />
                  <Form.Text className="text-muted">
                    Exclude the last x seconds from scene previews. This can be
                    a value in seconds, or a percentage (eg 2%) of the total
                    scene duration.
                  </Form.Text>
                </Form.Group>
              </div>
            </Collapse>
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
