import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Col, Form, Row } from "react-bootstrap";
import {
  mutateMigrateHashNaming,
  mutateMetadataExport,
  mutateBackupDatabase,
  mutateMetadataImport,
  mutateMetadataClean,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { downloadFile } from "src/utils";
import { Modal } from "../../Shared";
import { ImportDialog } from "./ImportDialog";
import * as GQL from "src/core/generated-graphql";
import { SettingSection } from "../SettingSection";
import { BooleanSetting, Setting } from "../Inputs";
import { ManualLink } from "src/components/Help/context";
import { Icon } from "src/components/Shared";
import { ConfigurationContext } from "src/hooks/Config";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";
import {
  faMinus,
  faPlus,
  faQuestionCircle,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";

interface ICleanDialog {
  pathSelection?: boolean;
  dryRun: boolean;
  onClose: (paths?: string[]) => void;
}

const CleanDialog: React.FC<ICleanDialog> = ({
  pathSelection = false,
  dryRun,
  onClose,
}) => {
  const intl = useIntl();
  const { configuration } = React.useContext(ConfigurationContext);

  const libraryPaths = configuration?.general.stashes.map((s) => s.path);

  const [paths, setPaths] = useState<string[]>([]);
  const [currentDirectory, setCurrentDirectory] = useState<string>("");

  function removePath(p: string) {
    setPaths(paths.filter((path) => path !== p));
  }

  function addPath(p: string) {
    if (p && !paths.includes(p)) {
      setPaths(paths.concat(p));
    }
  }

  let msg;
  if (dryRun) {
    msg = (
      <p>{intl.formatMessage({ id: "actions.tasks.dry_mode_selected" })}</p>
    );
  } else {
    msg = (
      <p>{intl.formatMessage({ id: "actions.tasks.clean_confirm_message" })}</p>
    );
  }

  return (
    <Modal
      show
      icon={faTrashAlt}
      disabled={pathSelection && paths.length === 0}
      accept={{
        text: intl.formatMessage({ id: "actions.clean" }),
        variant: "danger",
        onClick: () => onClose(paths),
      }}
      cancel={{ onClick: () => onClose() }}
    >
      <div className="dialog-container">
        <div className="mb-3">
          {paths.map((p) => (
            <Row className="align-items-center mb-1" key={p}>
              <Form.Label column xs={10}>
                {p}
              </Form.Label>
              <Col xs={2} className="d-flex justify-content-end">
                <Button
                  className="ml-auto"
                  size="sm"
                  variant="danger"
                  title={intl.formatMessage({ id: "actions.delete" })}
                  onClick={() => removePath(p)}
                >
                  <Icon icon={faMinus} />
                </Button>
              </Col>
            </Row>
          ))}

          {pathSelection ? (
            <FolderSelect
              currentDirectory={currentDirectory}
              setCurrentDirectory={(v) => setCurrentDirectory(v)}
              defaultDirectories={libraryPaths}
              appendButton={
                <Button
                  variant="secondary"
                  onClick={() => addPath(currentDirectory)}
                >
                  <Icon icon={faPlus} />
                </Button>
              }
            />
          ) : undefined}
        </div>

        {msg}
      </div>
    </Modal>
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
  function setOptions(input: Partial<GQL.CleanMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <>
      <BooleanSetting
        id="clean-dryrun"
        checked={options.dryRun}
        headingID="config.tasks.only_dry_run"
        onChange={(v) => setOptions({ dryRun: v })}
      />
    </>
  );
};

interface IDataManagementTasks {
  setIsBackupRunning: (v: boolean) => void;
}

export const DataManagementTasks: React.FC<IDataManagementTasks> = ({
  setIsBackupRunning,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [dialogOpen, setDialogOpenState] = useState({
    importAlert: false,
    import: false,
    clean: false,
    cleanAlert: false,
  });

  const [cleanOptions, setCleanOptions] = useState<GQL.CleanMetadataInput>({
    dryRun: false,
  });

  type DialogOpenState = typeof dialogOpen;

  function setDialogOpen(s: Partial<DialogOpenState>) {
    setDialogOpenState((v) => {
      return { ...v, ...s };
    });
  }

  async function onImport() {
    setDialogOpen({ importAlert: false });
    try {
      await mutateMetadataImport();
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.import" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderImportAlert() {
    return (
      <Modal
        show={dialogOpen.importAlert}
        icon={faTrashAlt}
        accept={{
          text: intl.formatMessage({ id: "actions.import" }),
          variant: "danger",
          onClick: onImport,
        }}
        cancel={{ onClick: () => setDialogOpen({ importAlert: false }) }}
      >
        <p>{intl.formatMessage({ id: "actions.tasks.import_warning" })}</p>
      </Modal>
    );
  }

  function renderImportDialog() {
    if (!dialogOpen.import) {
      return;
    }

    return <ImportDialog onClose={() => setDialogOpen({ import: false })} />;
  }

  async function onClean(paths?: string[]) {
    try {
      await mutateMetadataClean({
        ...cleanOptions,
        paths,
      });

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

  async function onMigrateHashNaming() {
    try {
      await mutateMigrateHashNaming();
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "actions.hash_migration",
            }),
          }
        ),
      });
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onExport() {
    try {
      await mutateMetadataExport();
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.backup" }) }
        ),
      });
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onBackup(download?: boolean) {
    try {
      setIsBackupRunning(true);
      const ret = await mutateBackupDatabase({
        download,
      });

      // download the result
      if (download && ret.data && ret.data.backupDatabase) {
        const link = ret.data.backupDatabase;
        downloadFile(link);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsBackupRunning(false);
    }
  }

  return (
    <Form.Group>
      {renderImportAlert()}
      {renderImportDialog()}
      {dialogOpen.cleanAlert || dialogOpen.clean ? (
        <CleanDialog
          dryRun={cleanOptions.dryRun}
          pathSelection={dialogOpen.clean}
          onClose={(p) => {
            // undefined means cancelled
            if (p !== undefined) {
              if (dialogOpen.cleanAlert) {
                // don't provide paths
                onClean();
              } else {
                onClean(p);
              }
            }

            setDialogOpen({
              clean: false,
              cleanAlert: false,
            });
          }}
        />
      ) : (
        dialogOpen.clean
      )}

      <SettingSection headingID="config.tasks.maintenance">
        <div className="setting-group">
          <Setting
            heading={
              <>
                <FormattedMessage id="actions.clean" />
                <ManualLink tab="Tasks">
                  <Icon icon={faQuestionCircle} />
                </ManualLink>
              </>
            }
            subHeadingID="config.tasks.cleanup_desc"
          >
            <Button
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ cleanAlert: true })}
            >
              <FormattedMessage id="actions.clean" />…
            </Button>
            <Button
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ clean: true })}
            >
              <FormattedMessage id="actions.selective_clean" />…
            </Button>
          </Setting>
          <CleanOptions
            options={cleanOptions}
            setOptions={(o) => setCleanOptions(o)}
          />
        </div>
      </SettingSection>

      <SettingSection headingID="metadata">
        <Setting
          headingID="actions.full_export"
          subHeadingID="config.tasks.export_to_json"
        >
          <Button
            id="export"
            variant="secondary"
            type="submit"
            onClick={() => onExport()}
          >
            <FormattedMessage id="actions.full_export" />
          </Button>
        </Setting>

        <Setting
          headingID="actions.full_import"
          subHeadingID="config.tasks.import_from_exported_json"
        >
          <Button
            id="import"
            variant="danger"
            type="submit"
            onClick={() => setDialogOpen({ importAlert: true })}
          >
            <FormattedMessage id="actions.full_import" />
          </Button>
        </Setting>

        <Setting
          headingID="actions.import_from_file"
          subHeadingID="config.tasks.incremental_import"
        >
          <Button
            id="partial-import"
            variant="danger"
            type="submit"
            onClick={() => setDialogOpen({ import: true })}
          >
            <FormattedMessage id="actions.import_from_file" />
          </Button>
        </Setting>
      </SettingSection>

      <SettingSection headingID="actions.backup">
        <Setting
          headingID="actions.backup"
          subHeading={intl.formatMessage(
            { id: "config.tasks.backup_database" },
            {
              filename_format: (
                <code>
                  [origFilename].sqlite.[schemaVersion].[YYYYMMDD_HHMMSS]
                </code>
              ),
            }
          )}
        >
          <Button
            id="backup"
            variant="secondary"
            type="submit"
            onClick={() => onBackup()}
          >
            <FormattedMessage id="actions.backup" />
          </Button>
        </Setting>

        <Setting
          headingID="actions.download_backup"
          subHeadingID="config.tasks.backup_and_download"
        >
          <Button
            id="backupDownload"
            variant="secondary"
            type="submit"
            onClick={() => onBackup(true)}
          >
            <FormattedMessage id="actions.download_backup" />
          </Button>
        </Setting>
      </SettingSection>

      <SettingSection headingID="config.tasks.migrations">
        <Setting
          headingID="actions.rename_gen_files"
          subHeadingID="config.tasks.migrate_hash_files"
        >
          <Button
            id="migrateHashNaming"
            variant="danger"
            onClick={() => onMigrateHashNaming()}
          >
            <FormattedMessage id="actions.rename_gen_files" />
          </Button>
        </Setting>
      </SettingSection>
    </Form.Group>
  );
};
