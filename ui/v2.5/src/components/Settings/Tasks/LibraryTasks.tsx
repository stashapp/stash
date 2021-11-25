import React, { useState, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataScan,
  mutateMetadataAutoTag,
  mutateMetadataGenerate,
  useConfigureDefaults,
} from "src/core/StashService";
import { withoutTypename } from "src/utils";
import { ConfigurationContext } from "src/hooks/Config";
import { IdentifyDialog } from "../../Dialogs/IdentifyDialog/IdentifyDialog";
import * as GQL from "src/core/generated-graphql";
import { DirectorySelectionDialog } from "./DirectorySelectionDialog";
import { ScanOptions } from "./ScanOptions";
import { useToast } from "src/hooks";
import { GenerateOptions } from "./GenerateOptions";
import { Task } from "./Task";

interface IAutoTagOptions {
  options: GQL.AutoTagMetadataInput;
  setOptions: (s: GQL.AutoTagMetadataInput) => void;
}

const AutoTagOptions: React.FC<IAutoTagOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const intl = useIntl();

  const { performers, studios, tags } = options;
  const wildcard = ["*"];

  function toggle(v?: GQL.Maybe<string[]>) {
    if (!v?.length) {
      return wildcard;
    }
    return [];
  }

  function setOptions(input: Partial<GQL.AutoTagMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <Form.Group>
      <Form.Check
        id="autotag-performers"
        checked={!!performers?.length}
        label={intl.formatMessage({ id: "performers" })}
        onChange={() => setOptions({ performers: toggle(performers) })}
      />
      <Form.Check
        id="autotag-studios"
        checked={!!studios?.length}
        label={intl.formatMessage({ id: "studios" })}
        onChange={() => setOptions({ studios: toggle(studios) })}
      />
      <Form.Check
        id="autotag-tags"
        checked={!!tags?.length}
        label={intl.formatMessage({ id: "tags" })}
        onChange={() => setOptions({ tags: toggle(tags) })}
      />
    </Form.Group>
  );
};

export const LibraryTasks: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  const [configureDefaults] = useConfigureDefaults();

  const [dialogOpen, setDialogOpenState] = useState({
    clean: false,
    scan: false,
    autoTag: false,
    identify: false,
  });

  const [scanOptions, setScanOptions] = useState<GQL.ScanMetadataInput>({});
  const [
    autoTagOptions,
    setAutoTagOptions,
  ] = useState<GQL.AutoTagMetadataInput>({
    performers: ["*"],
    studios: ["*"],
    tags: ["*"],
  });

  function getDefaultGenerateOptions(): GQL.GenerateMetadataInput {
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

  const [
    generateOptions,
    setGenerateOptions,
  ] = useState<GQL.GenerateMetadataInput>(getDefaultGenerateOptions());

  type DialogOpenState = typeof dialogOpen;

  const { configuration } = React.useContext(ConfigurationContext);

  useEffect(() => {
    if (!configuration?.defaults) {
      return;
    }

    const { scan, autoTag } = configuration.defaults;

    if (scan) {
      setScanOptions(withoutTypename(scan));
    }
    if (autoTag) {
      setAutoTagOptions(withoutTypename(autoTag));
    }

    if (configuration?.defaults.generate) {
      const { generate } = configuration.defaults;
      setGenerateOptions(withoutTypename(generate));
    } else if (configuration?.general) {
      // backwards compatibility
      const { general } = configuration;
      setGenerateOptions((existing) => ({
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

  function setDialogOpen(s: Partial<DialogOpenState>) {
    setDialogOpenState((v) => {
      return { ...v, ...s };
    });
  }

  function renderScanDialog() {
    if (!dialogOpen.scan) {
      return;
    }

    return <DirectorySelectionDialog onClose={onScanDialogClosed} />;
  }

  function onScanDialogClosed(paths?: string[]) {
    if (paths) {
      runScan(paths);
    }

    setDialogOpen({ scan: false });
  }

  async function runScan(paths?: string[]) {
    try {
      configureDefaults({
        variables: {
          input: {
            scan: scanOptions,
          },
        },
      });

      await mutateMetadataScan({
        ...scanOptions,
        paths,
      });

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.scan" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderAutoTagDialog() {
    if (!dialogOpen.autoTag) {
      return;
    }

    return <DirectorySelectionDialog onClose={onAutoTagDialogClosed} />;
  }

  function onAutoTagDialogClosed(paths?: string[]) {
    if (paths) {
      runAutoTag(paths);
    }

    setDialogOpen({ autoTag: false });
  }

  async function runAutoTag(paths?: string[]) {
    try {
      configureDefaults({
        variables: {
          input: {
            autoTag: autoTagOptions,
          },
        },
      });

      await mutateMetadataAutoTag({
        ...autoTagOptions,
        paths,
      });

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.auto_tag" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function maybeRenderIdentifyDialog() {
    if (!dialogOpen.identify) return;

    return (
      <IdentifyDialog onClose={() => setDialogOpen({ identify: false })} />
    );
  }

  async function onGenerateClicked() {
    try {
      configureDefaults({
        variables: {
          input: {
            generate: generateOptions,
          },
        },
      });

      await mutateMetadataGenerate(generateOptions);
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.generate" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <Form.Group>
      {renderScanDialog()}
      {renderAutoTagDialog()}
      {maybeRenderIdentifyDialog()}

      <Form.Group>
        <h5>{intl.formatMessage({ id: "library" })}</h5>

        <div className="task-group">
          <Task
            headingID="actions.scan"
            description={intl.formatMessage({
              id: "config.tasks.scan_for_content_desc",
            })}
          >
            <ScanOptions options={scanOptions} setOptions={setScanOptions} />
            <Button
              variant="secondary"
              type="submit"
              className="mr-2"
              onClick={() => runScan()}
            >
              <FormattedMessage id="actions.scan" />
            </Button>

            <Button
              variant="secondary"
              type="submit"
              className="mr-2"
              onClick={() => setDialogOpen({ scan: true })}
            >
              <FormattedMessage id="actions.selective_scan" />…
            </Button>
          </Task>

          <Task
            headingID="config.tasks.identify.heading"
            description={intl.formatMessage({
              id: "config.tasks.identify.description",
            })}
          >
            <Button
              variant="secondary"
              type="submit"
              onClick={() => setDialogOpen({ identify: true })}
            >
              <FormattedMessage id="actions.identify" />…
            </Button>
          </Task>

          <Task
            headingID="config.tasks.auto_tagging"
            description={intl.formatMessage({
              id: "config.tasks.auto_tag_based_on_filenames",
            })}
          >
            <AutoTagOptions
              options={autoTagOptions}
              setOptions={(o) => setAutoTagOptions(o)}
            />

            <Button
              variant="secondary"
              type="submit"
              className="mr-2"
              onClick={() => runAutoTag()}
            >
              <FormattedMessage id="actions.auto_tag" />
            </Button>
            <Button
              variant="secondary"
              type="submit"
              onClick={() => setDialogOpen({ autoTag: true })}
            >
              <FormattedMessage id="actions.selective_auto_tag" />…
            </Button>
          </Task>
        </div>
      </Form.Group>

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.tasks.generated_content" })}</h5>

        <div className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.generate_desc",
            })}
          >
            <GenerateOptions
              options={generateOptions}
              setOptions={setGenerateOptions}
            />
            <Button
              variant="secondary"
              type="submit"
              onClick={() => onGenerateClicked()}
            >
              <FormattedMessage id="actions.generate" />
            </Button>
          </Task>
        </div>
      </Form.Group>
    </Form.Group>
  );
};
