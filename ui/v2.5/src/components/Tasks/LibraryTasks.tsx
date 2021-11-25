import React, { useState, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, ButtonGroup, Card, Form } from "react-bootstrap";
import {
  mutateMetadataScan,
  mutateMetadataAutoTag,
  mutateMetadataGenerate,
  useConfigureDefaults,
  mutateMetadataClean,
} from "src/core/StashService";
import { withoutTypename } from "src/utils";
import { ConfigurationContext } from "src/hooks/Config";
import { PropsWithChildren } from "react-router/node_modules/@types/react";
import { IdentifyDialog } from "../Dialogs/IdentifyDialog/IdentifyDialog";
import { GenerateDialog } from "../Dialogs/GenerateDialog";
import * as GQL from "src/core/generated-graphql";
import { DirectorySelectionDialog } from "./DirectorySelectionDialog";
import { ScanOptions } from "./ScanOptions";
import { useToast } from "src/hooks";
import { Modal } from "../Shared";

interface ITask {
  description?: React.ReactNode;
}

const Task: React.FC<PropsWithChildren<ITask>> = ({
  children,
  description,
}) => (
  <div className="task">
    {children}
    {description ? (
      <Form.Text className="text-muted">{description}</Form.Text>
    ) : undefined}
  </div>
);

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

interface ICleanOptions {
  options: GQL.CleanMetadataInput;
  setOptions: (s: GQL.CleanMetadataInput) => void;
}

const CleanOptions: React.FC<ICleanOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const intl = useIntl();

  function setOptions(input: Partial<GQL.CleanMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <Form.Group>
      <Form.Check
        id="clean-dryrun"
        checked={options.dryRun}
        label={intl.formatMessage({ id: "config.tasks.only_dry_run" })}
        onChange={() => setOptions({ dryRun: !options.dryRun })}
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
    generate: false,
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
  const [cleanOptions, setCleanOptions] = useState<GQL.CleanMetadataInput>({
    dryRun: false,
  });

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
  }, [configuration]);

  function setDialogOpen(s: Partial<DialogOpenState>) {
    setDialogOpenState((v) => {
      return { ...v, ...s };
    });
  }

  function renderCleanDialog() {
    let msg;
    if (cleanOptions.dryRun) {
      msg = (
        <p>{intl.formatMessage({ id: "actions.tasks.dry_mode_selected" })}</p>
      );
    } else {
      msg = (
        <p>
          {intl.formatMessage({ id: "actions.tasks.clean_confirm_message" })}
        </p>
      );
    }

    return (
      <Modal
        show={dialogOpen.clean}
        icon="trash-alt"
        accept={{
          text: intl.formatMessage({ id: "actions.clean" }),
          variant: "danger",
          onClick: onClean,
        }}
        cancel={{ onClick: () => setDialogOpen({ clean: false }) }}
      >
        {msg}
      </Modal>
    );
  }

  async function onClean() {
    try {
      await mutateMetadataClean(cleanOptions);

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.clean" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setDialogOpen({ clean: false });
    }
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

  function maybeRenderGenerateDialog() {
    if (!dialogOpen.generate) return;

    return (
      <GenerateDialog onClose={() => setDialogOpen({ generate: false })} />
    );
  }

  async function onGenerateClicked() {
    // check if defaults are set for generate
    // if not, then open the dialog
    if (!configuration) {
      return;
    }

    const { generate } = configuration?.defaults;
    if (!generate) {
      setDialogOpen({ generate: true });
    } else {
      mutateMetadataGenerate(withoutTypename(generate));
    }
  }

  return (
    <Form.Group>
      {renderCleanDialog()}
      {renderScanDialog()}
      {renderAutoTagDialog()}
      {maybeRenderIdentifyDialog()}
      {maybeRenderGenerateDialog()}

      <Form.Group>
        <h5>{intl.formatMessage({ id: "library" })}</h5>

        <Card className="task-group">
          <Task
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

          <Task
            description={intl.formatMessage({
              id: "config.tasks.cleanup_desc",
            })}
          >
            <CleanOptions
              options={cleanOptions}
              setOptions={(o) => setCleanOptions(o)}
            />
            <Button
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ clean: true })}
            >
              <FormattedMessage id="actions.clean" />…
            </Button>
          </Task>
        </Card>
      </Form.Group>

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.tasks.generated_content" })}</h5>

        <Card className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.generate_desc",
            })}
          >
            <ButtonGroup className="ellipsis-button">
              <Button
                variant="secondary"
                type="submit"
                onClick={() => onGenerateClicked()}
              >
                <FormattedMessage id="actions.generate" />
              </Button>
              <Button
                variant="secondary"
                onClick={() => setDialogOpen({ generate: true })}
              >
                …
              </Button>
            </ButtonGroup>
          </Task>
        </Card>
      </Form.Group>
    </Form.Group>
  );
};
