import React, { useState, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, ButtonGroup, Card, Form } from "react-bootstrap";
import {
  mutateMetadataScan,
  mutateMetadataIdentify,
  mutateMetadataAutoTag,
  mutateMetadataGenerate,
  useConfigureDefaults,
} from "src/core/StashService";
import { withoutTypename } from "src/utils";
import { ConfigurationContext } from "src/hooks/Config";
import { PropsWithChildren } from "react-router/node_modules/@types/react";
import { CleanDialog } from "../Dialogs/CleanDialog";
import { AutoTagDialog } from "../Dialogs/AutoTagDialog";
import { IdentifyDialog } from "../Dialogs/IdentifyDialog/IdentifyDialog";
import { GenerateDialog } from "../Dialogs/GenerateDialog";
import * as GQL from "src/core/generated-graphql";
import { DirectorySelectionDialog } from "./DirectorySelectionDialog";
import { ScanOptions } from "./ScanOptions";
import { useToast } from "src/hooks";

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

  type DialogOpenState = typeof dialogOpen;

  const { configuration } = React.useContext(ConfigurationContext);

  useEffect(() => {
    if (!configuration?.defaults) {
      return;
    }

    const { scan } = configuration.defaults;

    if (scan) {
      setScanOptions(withoutTypename(scan));
    }
  }, [configuration]);

  function setDialogOpen(s: Partial<DialogOpenState>) {
    setDialogOpenState((v) => {
      return { ...v, ...s };
    });
  }

  function renderCleanDialog() {
    if (!dialogOpen.clean) {
      return;
    }

    return <CleanDialog onClose={() => setDialogOpen({ clean: false })} />;
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
      await configureDefaults({
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

    return <AutoTagDialog onClose={() => setDialogOpen({ autoTag: false })} />;
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

  async function onIdentifyClicked() {
    // check if defaults are set for identify
    // if not, then open the dialog
    if (!configuration) {
      return;
    }

    const { identify } = configuration?.defaults;
    if (!identify) {
      setDialogOpen({ identify: true });
    } else {
      mutateMetadataIdentify(withoutTypename(identify));
    }
  }

  async function onAutoTagClicked() {
    // check if defaults are set for auto tag
    // if not, then open the dialog
    if (!configuration) {
      return;
    }

    const { autoTag } = configuration?.defaults;
    if (!autoTag) {
      setDialogOpen({ autoTag: true });
    } else {
      mutateMetadataAutoTag(withoutTypename(autoTag));
    }
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
            <ButtonGroup className="ellipsis-button">
              <Button
                variant="secondary"
                type="submit"
                onClick={() => onIdentifyClicked()}
              >
                <FormattedMessage id="actions.identify" />
              </Button>
              <Button
                variant="secondary"
                onClick={() => setDialogOpen({ identify: true })}
              >
                …
              </Button>
            </ButtonGroup>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.auto_tag_based_on_filenames",
            })}
          >
            <ButtonGroup className="ellipsis-button">
              <Button
                variant="secondary"
                type="submit"
                onClick={() => onAutoTagClicked()}
              >
                <FormattedMessage id="actions.auto_tag" />
              </Button>
              <Button
                variant="secondary"
                onClick={() => setDialogOpen({ autoTag: true })}
              >
                …
              </Button>
            </ButtonGroup>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.cleanup_desc",
            })}
          >
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
