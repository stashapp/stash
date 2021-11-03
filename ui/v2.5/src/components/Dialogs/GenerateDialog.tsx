import React, { useState, useEffect, useMemo } from "react";
import { Form, Button, Collapse } from "react-bootstrap";
import { mutateMetadataGenerate } from "src/core/StashService";
import { Modal, Icon } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { ConfigurationContext } from "src/hooks/Config";
import { DirectorySelectionDialog } from "../Settings/SettingsTasksPanel/DirectorySelectionDialog";
import { Manual } from "../Help/Manual";

interface IGenerateOptions {
  options: GQL.GenerateMetadataInput;
  setOptions: (s: GQL.GenerateMetadataInput) => void;
}

const GenerateOptions: React.FC<IGenerateOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const intl = useIntl();

  const [previewOptionsOpen, setPreviewOptionsOpen] = useState(false);

  const previewOptions: GQL.GeneratePreviewOptionsInput =
    options.previewOptions ?? {};

  function setOptions(input: Partial<GQL.GenerateMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  function setPreviewOptions(input: Partial<GQL.GeneratePreviewOptionsInput>) {
    setOptions({
      previewOptions: {
        ...previewOptions,
        ...input,
      },
    });
  }

  return (
    <Form.Group>
      <Form.Group>
        <Form.Check
          id="preview-task"
          checked={options.previews ?? false}
          label={intl.formatMessage({
            id: "dialogs.scene_gen.video_previews",
          })}
          onChange={() => setOptions({ previews: !options.previews })}
        />
        <div className="d-flex flex-row">
          <div>↳</div>
          <Form.Check
            id="image-preview-task"
            checked={options.imagePreviews ?? false}
            disabled={!options.previews}
            label={intl.formatMessage({
              id: "dialogs.scene_gen.image_previews",
            })}
            onChange={() =>
              setOptions({ imagePreviews: !options.imagePreviews })
            }
            className="ml-2 flex-grow"
          />
        </div>
      </Form.Group>

      <Form.Group>
        <Button
          onClick={() => setPreviewOptionsOpen(!previewOptionsOpen)}
          className="minimal pl-0 no-focus"
        >
          <Icon icon={previewOptionsOpen ? "chevron-down" : "chevron-right"} />
          <span>
            {intl.formatMessage({
              id: "dialogs.scene_gen.preview_options",
            })}
          </span>
        </Button>
        <Form.Group>
          <Collapse in={previewOptionsOpen}>
            <Form.Group className="mt-2">
              <Form.Group id="preview-preset">
                <h6>
                  {intl.formatMessage({
                    id: "dialogs.scene_gen.preview_preset_head",
                  })}
                </h6>
                <Form.Control
                  className="w-auto input-control"
                  as="select"
                  value={previewOptions.previewPreset ?? GQL.PreviewPreset.Slow}
                  onChange={(e) =>
                    setPreviewOptions({
                      previewPreset: e.currentTarget.value as GQL.PreviewPreset,
                    })
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
                  value={previewOptions.previewSegments?.toString() ?? ""}
                  onChange={(e) =>
                    setPreviewOptions({
                      previewSegments: Number.parseInt(
                        e.currentTarget.value,
                        10
                      ),
                    })
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
                  value={
                    previewOptions.previewSegmentDuration?.toString() ?? ""
                  }
                  onChange={(e) =>
                    setPreviewOptions({
                      previewSegmentDuration: Number.parseFloat(
                        e.currentTarget.value
                      ),
                    })
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
                  value={previewOptions.previewExcludeStart ?? ""}
                  onChange={(e) =>
                    setPreviewOptions({
                      previewExcludeStart: e.currentTarget.value,
                    })
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
                  value={previewOptions.previewExcludeEnd ?? ""}
                  onChange={(e) =>
                    setPreviewOptions({
                      previewExcludeEnd: e.currentTarget.value,
                    })
                  }
                />
                <Form.Text className="text-muted">
                  {intl.formatMessage({
                    id: "dialogs.scene_gen.preview_exclude_end_time_desc",
                  })}
                </Form.Text>
              </Form.Group>
            </Form.Group>
          </Collapse>
        </Form.Group>
      </Form.Group>

      <Form.Group>
        <Form.Check
          id="sprite-task"
          checked={options.sprites ?? false}
          label={intl.formatMessage({ id: "dialogs.scene_gen.sprites" })}
          onChange={() => setOptions({ sprites: !options.sprites })}
        />
        <Form.Group>
          <Form.Check
            id="marker-task"
            checked={options.markers ?? false}
            label={intl.formatMessage({ id: "dialogs.scene_gen.markers" })}
            onChange={() => setOptions({ markers: !options.markers })}
          />
          <div className="d-flex flex-row">
            <div>↳</div>
            <Form.Group>
              <Form.Check
                id="marker-image-preview-task"
                checked={options.markerImagePreviews ?? false}
                disabled={!options.markers}
                label={intl.formatMessage({
                  id: "dialogs.scene_gen.marker_image_previews",
                })}
                onChange={() =>
                  setOptions({
                    markerImagePreviews: !options.markerImagePreviews,
                  })
                }
                className="ml-2 flex-grow"
              />
              <Form.Check
                id="marker-screenshot-task"
                checked={options.markerScreenshots ?? false}
                disabled={!options.markers}
                label={intl.formatMessage({
                  id: "dialogs.scene_gen.marker_screenshots",
                })}
                onChange={() =>
                  setOptions({ markerScreenshots: !options.markerScreenshots })
                }
                className="ml-2 flex-grow"
              />
            </Form.Group>
          </div>
        </Form.Group>

        <Form.Group>
          <Form.Check
            id="transcode-task"
            checked={options.transcodes ?? false}
            label={intl.formatMessage({ id: "dialogs.scene_gen.transcodes" })}
            onChange={() => setOptions({ transcodes: !options.transcodes })}
          />
          <Form.Check
            id="phash-task"
            checked={options.phashes ?? false}
            label={intl.formatMessage({ id: "dialogs.scene_gen.phash" })}
            onChange={() => setOptions({ phashes: !options.phashes })}
          />
        </Form.Group>

        <hr />
        <Form.Group>
          <Form.Check
            id="overwrite"
            checked={options.overwrite ?? false}
            label={intl.formatMessage({ id: "dialogs.scene_gen.overwrite" })}
            onChange={() => setOptions({ overwrite: !options.overwrite })}
          />
        </Form.Group>
      </Form.Group>
    </Form.Group>
  );
};

interface ISceneGenerateDialog {
  selectedIds?: string[];
  onClose: () => void;
}

export const GenerateDialog: React.FC<ISceneGenerateDialog> = ({
  selectedIds,
  onClose,
}) => {
  const { configuration } = React.useContext(ConfigurationContext);

  function getDefaultOptions(): GQL.GenerateMetadataInput {
    return {
      sprites: true,
      phashes: true,
      previews: true,
      markers: true,
      previewOptions: {
        previewSegments: 0,
        previewSegmentDuration: 0,
        previewPreset: GQL.PreviewPreset.Slow,
      },
    };
  }

  const [options, setOptions] = useState<GQL.GenerateMetadataInput>(
    getDefaultOptions()
  );
  const [paths, setPaths] = useState<string[]>([]);
  const [showManual, setShowManual] = useState(false);
  const [settingPaths, setSettingPaths] = useState(false);
  const [animation, setAnimation] = useState(true);

  const intl = useIntl();
  const Toast = useToast();

  useEffect(() => {
    if (!configuration) return;

    if (configuration.general) {
      const { general } = configuration;
      setOptions((existing) => ({
        ...existing,
        previewOptions: {
          ...existing.previewOptions,
          previewSegments:
            general.previewSegments ?? existing.previewOptions?.previewSegments,
          previewSegmentDuration:
            general.previewSegmentDuration ??
            existing.previewOptions?.previewSegmentDuration,
          previewExcludeStart:
            general.previewExcludeStart ??
            existing.previewOptions?.previewExcludeStart,
          previewExcludeEnd:
            general.previewExcludeEnd ??
            existing.previewOptions?.previewExcludeEnd,
          previewPreset:
            general.previewPreset ?? existing.previewOptions?.previewPreset,
        },
      }));
    }
  }, [configuration]);

  const selectionStatus = useMemo(() => {
    if (selectedIds) {
      return (
        <Form.Group id="selected-generate-ids">
          <FormattedMessage
            id="config.tasks.generate.generating_scenes"
            values={{
              num: selectedIds.length,
              scene: intl.formatMessage(
                {
                  id: "countables.scenes",
                },
                {
                  count: selectedIds.length,
                }
              ),
            }}
          />
          .
        </Form.Group>
      );
    }
    const message = paths.length ? (
      <div>
        <FormattedMessage id="config.tasks.generate.generating_from_paths" />:
        <ul>
          {paths.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      </div>
    ) : (
      <span>
        <FormattedMessage
          id="config.tasks.generate.generating_scenes"
          values={{
            num: intl.formatMessage({ id: "all" }),
            scene: intl.formatMessage(
              {
                id: "countables.scenes",
              },
              {
                count: 0,
              }
            ),
          }}
        />
        .
      </span>
    );

    function onClick() {
      setAnimation(false);
      setSettingPaths(true);
    }

    return (
      <Form.Group className="dialog-selected-folders">
        <div>
          {message}
          <div>
            <Button
              title={intl.formatMessage({ id: "actions.select_folders" })}
              onClick={() => onClick()}
            >
              <Icon icon="folder-open" />
            </Button>
          </div>
        </div>
      </Form.Group>
    );
  }, [selectedIds, intl, paths]);

  async function onGenerate() {
    try {
      await mutateMetadataGenerate(options);
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.generate" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  if (settingPaths) {
    return (
      <DirectorySelectionDialog
        animation={false}
        allowEmpty
        initialPaths={paths}
        onClose={(p) => {
          if (p) {
            setPaths(p);
          }
          setSettingPaths(false);
        }}
      />
    );
  }

  if (showManual) {
    return (
      <Manual
        animation={false}
        show
        onClose={() => setShowManual(false)}
        defaultActiveTab="Tasks.md"
      />
    );
  }

  return (
    <Modal
      show
      modalProps={{ animation, size: "lg" }}
      icon="cogs"
      header={intl.formatMessage({ id: "actions.generate" })}
      accept={{
        onClick: onGenerate,
        text: intl.formatMessage({ id: "actions.generate" }),
      }}
      cancel={{
        onClick: () => onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
    >
      <Form>
        {selectionStatus}
        <GenerateOptions options={options} setOptions={setOptions} />
      </Form>
    </Modal>
  );
};
