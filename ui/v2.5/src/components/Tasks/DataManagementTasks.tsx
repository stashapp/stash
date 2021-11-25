import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Card, Form } from "react-bootstrap";
import {
  mutateMigrateHashNaming,
  mutateMetadataExport,
  mutateBackupDatabase,
  mutateMetadataImport,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { downloadFile } from "src/utils";
import { PropsWithChildren } from "react-router/node_modules/@types/react";
import { Modal } from "../Shared";
import { ImportDialog } from "./ImportDialog";

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
    scan: false,
    autoTag: false,
    identify: false,
    generate: false,
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
        icon="trash-alt"
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

      <Form.Group>
        <h5>{intl.formatMessage({ id: "metadata" })}</h5>
        <Card className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.export_to_json",
            })}
          >
            <Button
              id="export"
              variant="secondary"
              type="submit"
              onClick={() => onExport()}
            >
              <FormattedMessage id="actions.full_export" />
            </Button>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.import_from_exported_json",
            })}
          >
            <Button
              id="import"
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ importAlert: true })}
            >
              <FormattedMessage id="actions.full_import" />
            </Button>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.incremental_import",
            })}
          >
            <Button
              id="partial-import"
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ import: true })}
            >
              <FormattedMessage id="actions.import_from_file" />
            </Button>
          </Task>
        </Card>
      </Form.Group>

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "actions.backup" })}</h5>
        <Card className="task-group">
          <Task
            description={intl.formatMessage(
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
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.backup_and_download",
            })}
          >
            <Button
              id="backupDownload"
              variant="secondary"
              type="submit"
              onClick={() => onBackup(true)}
            >
              <FormattedMessage id="actions.download_backup" />
            </Button>
          </Task>
        </Card>
      </Form.Group>

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.tasks.migrations" })}</h5>

        <Card className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.migrate_hash_files",
            })}
          >
            <Button
              id="migrateHashNaming"
              variant="danger"
              onClick={() => onMigrateHashNaming()}
            >
              <FormattedMessage id="actions.rename_gen_files" />
            </Button>
          </Task>
        </Card>
      </Form.Group>
    </Form.Group>
  );
};
