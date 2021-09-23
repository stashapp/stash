import React, { useState, useEffect } from "react";
import { Form, Button, Collapse } from "react-bootstrap";
import {
  mutateMetadataGenerate,
  useConfiguration,
} from "src/core/StashService";
import { Modal, Icon } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { useIntl } from "react-intl";

interface ISceneGenerateDialogProps {
  selectedIds: string[];
  onClose: () => void;
}

export const SceneGenerateDialog: React.FC<ISceneGenerateDialogProps> = (
  props: ISceneGenerateDialogProps
) => {
  const { data, error, loading } = useConfiguration();

  const [sprites, setSprites] = useState(true);
  const [phashes, setPhashes] = useState(true);
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
  const [markerImagePreviews, setMarkerImagePreviews] = useState(false);
  const [markerScreenshots, setMarkerScreenshots] = useState(false);

  const [previewOptionsOpen, setPreviewOptionsOpen] = useState(false);

  const intl = useIntl();
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
        phashes,
        previews,
        imagePreviews: previews && imagePreviews,
        markers,
        markerImagePreviews: markers && markerImagePreviews,
        markerScreenshots: markers && markerScreenshots,
        transcodes,
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
      header={intl.formatMessage({ id: "actions.generate" })}
      accept={{
        onClick: onGenerate,
        text: intl.formatMessage({ id: "actions.generate" }),
      }}
      cancel={{
        onClick: () => props.onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
    >
      <Form>
        <Form.Group>
          <Form.Check
            id="preview-task"
            checked={previews}
            label={intl.formatMessage({
              id: "dialogs.scene_gen.video_previews",
            })}
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
          <div className="my-2">
            <Button
              onClick={() => setPreviewOptionsOpen(!previewOptionsOpen)}
              className="minimal pl-0 no-focus"
            >
              <Icon
                icon={previewOptionsOpen ? "chevron-down" : "chevron-right"}
              />
              <span>
                {intl.formatMessage({
                  id: "dialogs.scene_gen.preview_options",
                })}
              </span>
            </Button>
            <Collapse in={previewOptionsOpen}>
              <div>
                <Form.Group id="transcode-size">
                  <h6>
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_preset_head",
                    })}
                  </h6>
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
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_preset_desc",
                    })}
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-segments">
                  <h6>
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_seg_count_head",
                    })}
                  </h6>
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
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_seg_count_desc",
                    })}
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-segment-duration">
                  <h6>
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_seg_duration_head",
                    })}
                  </h6>
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
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_seg_duration_desc",
                    })}
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-exclude-start">
                  <h6>
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_exclude_start_time_head",
                    })}
                  </h6>
                  <Form.Control
                    className="col col-sm-6 text-input"
                    defaultValue={previewExcludeStart}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setPreviewExcludeStart(e.currentTarget.value)
                    }
                  />
                  <Form.Text className="text-muted">
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_exclude_start_time_desc",
                    })}
                  </Form.Text>
                </Form.Group>

                <Form.Group id="preview-exclude-start">
                  <h6>
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_exclude_end_time_head",
                    })}
                  </h6>
                  <Form.Control
                    className="col col-sm-6 text-input"
                    defaultValue={previewExcludeEnd}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setPreviewExcludeEnd(e.currentTarget.value)
                    }
                  />
                  <Form.Text className="text-muted">
                    {intl.formatMessage({
                      id: "dialogs.scene_gen.preview_exclude_end_time_desc",
                    })}
                  </Form.Text>
                </Form.Group>
              </div>
            </Collapse>
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

          <hr />
          <Form.Check
            id="overwrite"
            checked={overwrite}
            label={intl.formatMessage({ id: "dialogs.scene_gen.overwrite" })}
            onChange={() => setOverwrite(!overwrite)}
          />
        </Form.Group>
      </Form>
    </Modal>
  );
};
