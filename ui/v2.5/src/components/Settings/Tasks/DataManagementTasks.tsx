import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Col, Form, Row } from "react-bootstrap";
import {
  mutateMigrateHashNaming,
  mutateMetadataExport,
  mutateBackupDatabase,
  mutateMetadataImport,
  mutateVerifyPaths,
  mutateAnonymiseDatabase,
  mutateMigrateSceneScreenshots,
  mutateMigrateBlobs,
  mutateOptimiseDatabase,
  mutateCleanGenerated,
  mutatePurgeMissing,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import downloadFile from "src/utils/download";
import { ModalComponent } from "src/components/Shared/Modal";
import { ImportDialog } from "./ImportDialog";
import * as GQL from "src/core/generated-graphql";
import { SettingSection } from "../SettingSection";
import { BooleanSetting, Setting } from "../Inputs";
import { ManualLink } from "src/components/Help/context";
import { Icon } from "src/components/Shared/Icon";
import { useConfigurationContext } from "src/hooks/Config";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";
import {
  faBoxArchive,
  faExclamationTriangle,
  faMinus,
  faPlus,
  faQuestionCircle,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { CleanGeneratedDialog } from "./CleanGeneratedDialog";
import { useSettings } from "../context";

interface IVerifyDialog {
  pathSelection?: boolean;
  dryRun: boolean;
  purgeMissing: boolean;
  onClose: (paths?: string[]) => void;
}

const VerifyDialog: React.FC<IVerifyDialog> = ({
  pathSelection = false,
  dryRun,
  purgeMissing,
  onClose,
}) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();

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

  let msg: React.ReactNode;
  if (dryRun) {
    msg = (
      <p>{intl.formatMessage({ id: "actions.tasks.dry_mode_selected" })}</p>
    );
  } else if (purgeMissing) {
    msg = (
      <p>
        {intl.formatMessage({
          id: "actions.tasks.verify_purge_confirm_message",
        })}
      </p>
    );
  }

  return (
    <ModalComponent
      show
      header={<FormattedMessage id="actions.verify_files" />}
      icon={faTrashAlt}
      disabled={pathSelection && paths.length === 0}
      accept={{
        text: intl.formatMessage({ id: "actions.verify_files" }),
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
              onChangeDirectory={setCurrentDirectory}
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
    </ModalComponent>
  );
};

interface IVerifyOptions {
  options: GQL.VerifyPathsInput;
  setOptions: (s: GQL.VerifyPathsInput) => void;
}

const VerifyOptions: React.FC<IVerifyOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  function setOptions(input: Partial<GQL.VerifyPathsInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <>
      <BooleanSetting
        id="verify-ignore-zip-contents"
        checked={options.checkZipFileContents ?? false}
        headingID="config.tasks.verify_check_zip_contents"
        subHeadingID="config.tasks.verify_check_zip_contents_desc"
        onChange={(v) => setOptions({ checkZipFileContents: v })}
      />
      <BooleanSetting
        id="verify-dryrun"
        checked={options.dryRun ?? false}
        headingID="config.tasks.dry_run"
        onChange={(v) => setOptions({ dryRun: v })}
      />
      <BooleanSetting
        id="verify-purge-missing"
        checked={(!options.dryRun && options.purgeMissing) ?? false}
        disabled={options.dryRun ?? false}
        headingID="config.tasks.verify_purge_missing"
        onChange={(v) => setOptions({ purgeMissing: v })}
      />
    </>
  );
};

interface IPurgeMissingDialog {
  pathSelection?: boolean;
  dryRun: boolean;
  onClose: (paths?: string[]) => void;
}

const PurgeMissingDialog: React.FC<IPurgeMissingDialog> = ({
  pathSelection = false,
  dryRun,
  onClose,
}) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();

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

  let msg: React.ReactNode;
  if (dryRun) {
    msg = (
      <p>{intl.formatMessage({ id: "actions.tasks.dry_mode_selected" })}</p>
    );
  } else {
    msg = (
      <p>
        {intl.formatMessage({
          id: "actions.tasks.purge_missing_confirm_message",
        })}
      </p>
    );
  }

  return (
    <ModalComponent
      show
      icon={faTrashAlt}
      header={<FormattedMessage id="config.tasks.purge_missing" />}
      disabled={pathSelection && paths.length === 0}
      accept={{
        text: intl.formatMessage({ id: "config.tasks.purge_missing" }),
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
              onChangeDirectory={setCurrentDirectory}
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
    </ModalComponent>
  );
};

interface IPurgeMissingOptions {
  options: GQL.PurgeMissingInput;
  setOptions: (s: GQL.PurgeMissingInput) => void;
}

const PurgeMissingOptions: React.FC<IPurgeMissingOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  function setOptions(input: Partial<GQL.PurgeMissingInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <BooleanSetting
      id="purge-missing-dryrun"
      checked={options.dryRun ?? false}
      headingID="config.tasks.dry_run"
      onChange={(v) => setOptions({ dryRun: v })}
    />
  );
};

const BackupDialog: React.FC<{
  onClose: (
    confirmed?: boolean,
    download?: boolean,
    includeBlobs?: boolean
  ) => void;
}> = ({ onClose }) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();

  const includeBlobsDefault =
    configuration?.general.blobsStorage === GQL.BlobsStorageType.Filesystem;
  const backupDir =
    configuration.general.backupDirectoryPath ||
    `<${intl.formatMessage({
      id: "config.general.backup_directory_path.heading",
    })}>`;

  const [download, setDownload] = useState(false);
  const [includeBlobs, setIncludeBlobs] = useState(includeBlobsDefault);

  let msg: React.ReactNode;
  if (!includeBlobs) {
    msg = intl.formatMessage(
      { id: "config.tasks.backup_database.sqlite" },
      {
        filename_format: (
          <code>[origFilename].sqlite.[schemaVersion].[YYYYMMDD_HHMMSS]</code>
        ),
      }
    );
  } else {
    msg = intl.formatMessage(
      { id: "config.tasks.backup_database.zip" },
      {
        filename_format: (
          <code>
            [origFilename].sqlite.[schemaVersion].[YYYYMMDD_HHMMSS].zip
          </code>
        ),
      }
    );
  }

  const warning =
    includeBlobs !== includeBlobsDefault ? (
      <p className="lead">
        <Icon icon={faExclamationTriangle} className="text-warning" />
        <FormattedMessage id="config.tasks.backup_database.warning_blobs" />
      </p>
    ) : null;

  const acceptID = download
    ? "config.tasks.backup_database.download"
    : "actions.backup";

  return (
    <ModalComponent
      show
      icon={faBoxArchive}
      accept={{
        text: intl.formatMessage({ id: acceptID }),
        onClick: () => onClose(true, download, includeBlobs),
      }}
      cancel={{
        onClick: () => onClose(),
        variant: "secondary",
      }}
    >
      <div className="dialog-container">
        <Form.Group>
          <h5>
            <FormattedMessage id="config.tasks.backup_database.destination" />
          </h5>
          <Form.Check
            type="radio"
            id="backup-directory"
            checked={!download}
            onChange={() => setDownload(false)}
            label={intl.formatMessage(
              {
                id: "config.tasks.backup_database.to_directory",
              },
              {
                directory: <code>{backupDir}</code>,
              }
            )}
          />

          <Form.Check
            type="radio"
            id="backup-download"
            checked={download}
            onChange={() => setDownload(true)}
            label={intl.formatMessage({
              id: "config.tasks.backup_database.download",
            })}
          />
        </Form.Group>

        <SettingSection>
          <BooleanSetting
            id="backup-include-blobs"
            // if includeBlobsDefault is false, then blobs are in the database, so we check the box and disable it
            checked={includeBlobs || !includeBlobsDefault}
            headingID="config.tasks.backup_database.include_blobs"
            onChange={(v) => setIncludeBlobs(v)}
            // if includeBlobsDefault is false, then blobs are in the database
            disabled={!includeBlobsDefault}
          />
        </SettingSection>

        <p>{msg}</p>
        {warning}
      </div>
    </ModalComponent>
  );
};

interface IDataManagementTasks {
  setIsBackupRunning: (v: boolean) => void;
  setIsAnonymiseRunning: (v: boolean) => void;
}

export const DataManagementTasks: React.FC<IDataManagementTasks> = ({
  setIsBackupRunning,
  setIsAnonymiseRunning,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [dialogOpen, setDialogOpenState] = useState({
    importAlert: false,
    import: false,
    backup: false,
    verify: false,
    verifyAlert: false,
    purgeMissing: false,
    purgeMissingAlert: false,
    cleanGenerated: false,
  });

  const [verifyOptions, setVerifyOptions] = useState<GQL.VerifyPathsInput>({
    dryRun: false,
    checkZipFileContents: false,
    purgeMissing: false,
  });

  const [purgeMissingOptions, setPurgeMissingOptions] =
    useState<GQL.PurgeMissingInput>({
      dryRun: false,
    });

  const [migrateBlobsOptions, setMigrateBlobsOptions] =
    useState<GQL.MigrateBlobsInput>({
      deleteOld: true,
    });

  const [migrateSceneScreenshotsOptions, setMigrateSceneScreenshotsOptions] =
    useState<GQL.MigrateSceneScreenshotsInput>({
      deleteFiles: false,
      overwriteExisting: false,
    });

  const { ui, saveUI, loading } = useSettings();

  const { taskDefaults } = ui;

  useEffect(() => {
    if (loading) {
      return;
    }

    if (taskDefaults?.verify) {
      setVerifyOptions(taskDefaults.verify);
    }
    if (taskDefaults?.purgeMissing) {
      setPurgeMissingOptions(taskDefaults.purgeMissing);
    }
  }, [taskDefaults, loading]);

  function configureDefaults(partial: Record<string, object>) {
    saveUI({ taskDefaults: { ...partial } });
  }

  function onSetVerifyOptions(s: GQL.VerifyPathsInput) {
    configureDefaults({ verify: s });
    setVerifyOptions(s);
  }

  function onSetPurgeMissingOptions(s: GQL.PurgeMissingInput) {
    configureDefaults({ purgeMissing: s });
    setPurgeMissingOptions(s);
  }

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
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.import" }) }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderImportAlert() {
    return (
      <ModalComponent
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
      </ModalComponent>
    );
  }

  function renderImportDialog() {
    if (!dialogOpen.import) {
      return;
    }

    return <ImportDialog onClose={() => setDialogOpen({ import: false })} />;
  }

  async function onVerify(paths?: string[]) {
    try {
      await mutateVerifyPaths({
        ...verifyOptions,
        paths,
      });

      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.verify_files" }) }
        )
      );
    } catch (e) {
      Toast.error(e);
    } finally {
      setDialogOpen({ verify: false });
    }
  }

  function onVerifyClicked() {
    if (verifyOptions.dryRun || !verifyOptions.purgeMissing) {
      onVerify();
      return;
    }

    setDialogOpen({ verifyAlert: true });
  }

  async function onPurgeMissing(paths?: string[]) {
    try {
      await mutatePurgeMissing({
        ...purgeMissingOptions,
        paths,
      });

      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "config.tasks.purge_missing",
            }),
          }
        )
      );
    } catch (e) {
      Toast.error(e);
    } finally {
      setDialogOpen({ purgeMissing: false });
    }
  }

  function onPurgeMissingClicked() {
    if (purgeMissingOptions.dryRun) {
      onPurgeMissing();
      return;
    }

    setDialogOpen({ purgeMissingAlert: true });
  }

  async function onCleanGenerated(options: GQL.CleanGeneratedInput) {
    try {
      await mutateCleanGenerated({
        ...options,
      });

      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "actions.clean_generated",
            }),
          }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onMigrateHashNaming() {
    try {
      await mutateMigrateHashNaming();
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "actions.hash_migration",
            }),
          }
        )
      );
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onMigrateSceneScreenshots() {
    try {
      await mutateMigrateSceneScreenshots(migrateSceneScreenshotsOptions);
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "actions.migrate_scene_screenshots",
            }),
          }
        )
      );
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onMigrateBlobs() {
    try {
      await mutateMigrateBlobs(migrateBlobsOptions);
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "actions.migrate_blobs",
            }),
          }
        )
      );
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onExport() {
    try {
      await mutateMetadataExport();
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.export" }) }
        )
      );
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onBackup(download?: boolean, includeBlobs?: boolean) {
    try {
      setIsBackupRunning(true);
      const ret = await mutateBackupDatabase({
        download,
        includeBlobs,
      });

      // download the result
      if (download && ret.data?.backupDatabase) {
        const link = ret.data.backupDatabase;
        downloadFile(link);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsBackupRunning(false);
    }
  }

  async function onOptimiseDatabase() {
    try {
      await mutateOptimiseDatabase();
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "actions.optimise_database",
            }),
          }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onAnonymise(download?: boolean) {
    try {
      setIsAnonymiseRunning(true);
      const ret = await mutateAnonymiseDatabase({
        download,
      });

      // download the result
      if (download && ret.data?.anonymiseDatabase) {
        const link = ret.data.anonymiseDatabase;
        downloadFile(link);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsAnonymiseRunning(false);
    }
  }

  return (
    <Form.Group>
      {renderImportAlert()}
      {renderImportDialog()}
      {dialogOpen.verifyAlert || dialogOpen.verify ? (
        <VerifyDialog
          dryRun={verifyOptions.dryRun ?? false}
          purgeMissing={verifyOptions.purgeMissing ?? false}
          pathSelection={dialogOpen.verify}
          onClose={(p) => {
            // undefined means cancelled
            if (p !== undefined) {
              if (dialogOpen.verifyAlert) {
                // don't provide paths
                onVerify();
              } else {
                onVerify(p);
              }
            }

            setDialogOpen({
              verify: false,
              verifyAlert: false,
            });
          }}
        />
      ) : null}
      {dialogOpen.purgeMissingAlert || dialogOpen.purgeMissing ? (
        <PurgeMissingDialog
          dryRun={purgeMissingOptions.dryRun ?? false}
          pathSelection={dialogOpen.purgeMissing}
          onClose={(p) => {
            // undefined means cancelled
            if (p !== undefined) {
              if (dialogOpen.purgeMissingAlert) {
                // don't provide paths
                onPurgeMissing();
              } else {
                onPurgeMissing(p);
              }
            }

            setDialogOpen({
              purgeMissing: false,
              purgeMissingAlert: false,
            });
          }}
        />
      ) : null}
      {dialogOpen.cleanGenerated && (
        <CleanGeneratedDialog
          onClose={(options) => {
            if (options) {
              onCleanGenerated(options);
            }

            setDialogOpen({ cleanGenerated: false });
          }}
        />
      )}
      {dialogOpen.backup && (
        <BackupDialog
          onClose={(confirmed, download, includeBlobs) => {
            if (confirmed) {
              onBackup(download, includeBlobs);
            }

            setDialogOpen({ backup: false });
          }}
        />
      )}

      <SettingSection headingID="config.tasks.maintenance">
        <div className="setting-group">
          <Setting
            heading={
              <>
                <FormattedMessage id="actions.verify_files" />
                <ManualLink tab="Tasks">
                  <Icon icon={faQuestionCircle} />
                </ManualLink>
              </>
            }
            subHeadingID="config.tasks.verify_files_desc"
          >
            <Button
              variant="danger"
              type="submit"
              onClick={() => onVerifyClicked()}
            >
              <FormattedMessage id="actions.verify_files" />
            </Button>
            <Button
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ verify: true })}
            >
              <FormattedMessage id="actions.selective_verify" />…
            </Button>
          </Setting>
          <VerifyOptions
            options={verifyOptions}
            setOptions={(o) => onSetVerifyOptions(o)}
          />
        </div>

        <div className="setting-group">
          <Setting
            heading={
              <>
                <FormattedMessage id="config.tasks.purge_missing" />
                <ManualLink tab="Tasks">
                  <Icon icon={faQuestionCircle} />
                </ManualLink>
              </>
            }
            subHeadingID="config.tasks.purge_missing_desc"
          >
            <Button
              variant="danger"
              type="submit"
              onClick={() => onPurgeMissingClicked()}
            >
              <FormattedMessage id="config.tasks.purge_missing" />
            </Button>
            <Button
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ purgeMissing: true })}
            >
              <FormattedMessage id="config.tasks.selective_purge_missing" />…
            </Button>
          </Setting>
          <PurgeMissingOptions
            options={purgeMissingOptions}
            setOptions={(o) => onSetPurgeMissingOptions(o)}
          />
        </div>

        <div className="setting-group">
          <Setting
            heading={<FormattedMessage id="actions.clean_generated" />}
            subHeadingID="config.tasks.clean_generated.description"
          >
            <Button
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ cleanGenerated: true })}
            >
              <FormattedMessage id="actions.clean_generated" />…
            </Button>
          </Setting>
        </div>

        <Setting
          headingID="actions.optimise_database"
          subHeading={
            <>
              <FormattedMessage id="config.tasks.optimise_database" />
              <br />
              <FormattedMessage id="config.tasks.optimise_database_warning" />
            </>
          }
        >
          <Button
            id="optimiseDatabase"
            variant="danger"
            onClick={() => onOptimiseDatabase()}
          >
            <FormattedMessage id="actions.optimise_database" />
          </Button>
        </Setting>
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
          heading={
            <>
              <FormattedMessage id="actions.backup" />
              <ManualLink tab="Tasks">
                <Icon icon={faQuestionCircle} />
              </ManualLink>
            </>
          }
          subHeading={intl.formatMessage({
            id: "config.tasks.backup_database.description",
          })}
        >
          <Button
            id="backup"
            variant="secondary"
            type="submit"
            onClick={() => setDialogOpen({ backup: true })}
          >
            <FormattedMessage id="actions.backup" />…
          </Button>
        </Setting>
      </SettingSection>

      <SettingSection headingID="actions.anonymise">
        <Setting
          headingID="actions.anonymise"
          subHeading={intl.formatMessage(
            { id: "config.tasks.anonymise_database" },
            {
              filename_format: (
                <code>
                  [origFilename].anonymous.sqlite.[schemaVersion].[YYYYMMDD_HHMMSS]
                </code>
              ),
            }
          )}
        >
          <Button
            id="anonymise"
            variant="secondary"
            type="submit"
            onClick={() => onAnonymise()}
          >
            <FormattedMessage id="actions.anonymise" />
          </Button>
        </Setting>

        <Setting
          headingID="actions.download_anonymised"
          subHeadingID="config.tasks.anonymise_and_download"
        >
          <Button
            id="anonymiseDownload"
            variant="secondary"
            type="submit"
            onClick={() => onAnonymise(true)}
          >
            <FormattedMessage id="actions.download_anonymised" />
          </Button>
        </Setting>
      </SettingSection>

      <SettingSection headingID="config.tasks.migrations">
        <Setting
          advanced
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

        <div className="setting-group">
          <Setting
            headingID="actions.migrate_blobs"
            subHeadingID="config.tasks.migrate_blobs.description"
          >
            <Button
              id="migrateBlobs"
              variant="danger"
              onClick={() => onMigrateBlobs()}
            >
              <FormattedMessage id="actions.migrate_blobs" />
            </Button>
          </Setting>

          <BooleanSetting
            id="migrate-blobs-delete-old"
            checked={migrateBlobsOptions.deleteOld ?? false}
            headingID="config.tasks.migrate_blobs.delete_old"
            onChange={(v) =>
              setMigrateBlobsOptions({ ...migrateBlobsOptions, deleteOld: v })
            }
          />
        </div>

        <div className="setting-group">
          <Setting
            headingID="actions.migrate_scene_screenshots"
            subHeadingID="config.tasks.migrate_scene_screenshots.description"
          >
            <Button
              id="migrateSceneScreenshots"
              variant="danger"
              onClick={() => onMigrateSceneScreenshots()}
            >
              <FormattedMessage id="actions.migrate_scene_screenshots" />
            </Button>
          </Setting>

          <BooleanSetting
            id="migrate-scene-screenshots-overwrite-existing"
            checked={migrateSceneScreenshotsOptions.overwriteExisting ?? false}
            headingID="config.tasks.migrate_scene_screenshots.overwrite_existing"
            onChange={(v) =>
              setMigrateSceneScreenshotsOptions({
                ...migrateSceneScreenshotsOptions,
                overwriteExisting: v,
              })
            }
          />

          <BooleanSetting
            id="migrate-scene-screenshots-delete-files"
            checked={migrateSceneScreenshotsOptions.deleteFiles ?? false}
            headingID="config.tasks.migrate_scene_screenshots.delete_files"
            onChange={(v) =>
              setMigrateSceneScreenshotsOptions({
                ...migrateSceneScreenshotsOptions,
                deleteFiles: v,
              })
            }
          />
        </div>
      </SettingSection>
    </Form.Group>
  );
};
