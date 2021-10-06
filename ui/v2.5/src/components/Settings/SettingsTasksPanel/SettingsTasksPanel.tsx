import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataImport,
  mutateMetadataClean,
  mutateMetadataScan,
  mutateMetadataAutoTag,
  mutateMetadataExport,
  mutateMigrateHashNaming,
  usePlugins,
  mutateRunPluginTask,
  mutateBackupDatabase,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator, Modal } from "src/components/Shared";
import { downloadFile } from "src/utils";
import { GenerateButton } from "./GenerateButton";
import { ImportDialog } from "./ImportDialog";
import { DirectorySelectionDialog } from "./DirectorySelectionDialog";
import { JobTable } from "./JobTable";

type Plugin = Pick<GQL.Plugin, "id">;
type PluginTask = Pick<GQL.PluginTask, "name" | "description">;

export const SettingsTasksPanel: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  const [isImportAlertOpen, setIsImportAlertOpen] = useState<boolean>(false);
  const [isCleanAlertOpen, setIsCleanAlertOpen] = useState<boolean>(false);
  const [isImportDialogOpen, setIsImportDialogOpen] = useState<boolean>(false);
  const [isScanDialogOpen, setIsScanDialogOpen] = useState<boolean>(false);
  const [isAutoTagDialogOpen, setIsAutoTagDialogOpen] = useState<boolean>(
    false
  );
  const [isBackupRunning, setIsBackupRunning] = useState<boolean>(false);
  const [useFileMetadata, setUseFileMetadata] = useState<boolean>(false);
  const [stripFileExtension, setStripFileExtension] = useState<boolean>(false);
  const [scanGeneratePreviews, setScanGeneratePreviews] = useState<boolean>(
    false
  );
  const [scanGenerateSprites, setScanGenerateSprites] = useState<boolean>(
    false
  );
  const [scanGeneratePhashes, setScanGeneratePhashes] = useState<boolean>(
    false
  );
  const [scanGenerateThumbnails, setScanGenerateThumbnails] = useState<boolean>(
    false
  );
  const [cleanDryRun, setCleanDryRun] = useState<boolean>(false);
  const [
    scanGenerateImagePreviews,
    setScanGenerateImagePreviews,
  ] = useState<boolean>(false);

  const [autoTagPerformers, setAutoTagPerformers] = useState<boolean>(true);
  const [autoTagStudios, setAutoTagStudios] = useState<boolean>(true);
  const [autoTagTags, setAutoTagTags] = useState<boolean>(true);

  const plugins = usePlugins();

  async function onImport() {
    setIsImportAlertOpen(false);
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
        show={isImportAlertOpen}
        icon="trash-alt"
        accept={{
          text: intl.formatMessage({ id: "actions.import" }),
          variant: "danger",
          onClick: onImport,
        }}
        cancel={{ onClick: () => setIsImportAlertOpen(false) }}
      >
        <p>{intl.formatMessage({ id: "actions.tasks.import_warning" })}</p>
      </Modal>
    );
  }

  function onClean() {
    setIsCleanAlertOpen(false);
    mutateMetadataClean({
      dryRun: cleanDryRun,
    });
  }

  function renderCleanAlert() {
    let msg;
    if (cleanDryRun) {
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
        show={isCleanAlertOpen}
        icon="trash-alt"
        accept={{
          text: intl.formatMessage({ id: "actions.clean" }),
          variant: "danger",
          onClick: onClean,
        }}
        cancel={{ onClick: () => setIsCleanAlertOpen(false) }}
      >
        {msg}
      </Modal>
    );
  }

  function renderImportDialog() {
    if (!isImportDialogOpen) {
      return;
    }

    return <ImportDialog onClose={() => setIsImportDialogOpen(false)} />;
  }

  function renderScanDialog() {
    if (!isScanDialogOpen) {
      return;
    }

    return <DirectorySelectionDialog onClose={onScanDialogClosed} />;
  }

  function onScanDialogClosed(paths?: string[]) {
    if (paths) {
      onScan(paths);
    }

    setIsScanDialogOpen(false);
  }

  async function onScan(paths?: string[]) {
    try {
      await mutateMetadataScan({
        paths,
        useFileMetadata,
        stripFileExtension,
        scanGeneratePreviews,
        scanGenerateImagePreviews,
        scanGenerateSprites,
        scanGeneratePhashes,
        scanGenerateThumbnails,
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
    if (!isAutoTagDialogOpen) {
      return;
    }

    return <DirectorySelectionDialog onClose={onAutoTagDialogClosed} />;
  }

  function onAutoTagDialogClosed(paths?: string[]) {
    if (paths) {
      onAutoTag(paths);
    }

    setIsAutoTagDialogOpen(false);
  }

  function getAutoTagInput(paths?: string[]) {
    const wildcard = ["*"];
    return {
      paths,
      performers: autoTagPerformers ? wildcard : [],
      studios: autoTagStudios ? wildcard : [],
      tags: autoTagTags ? wildcard : [],
    };
  }

  async function onAutoTag(paths?: string[]) {
    try {
      await mutateMetadataAutoTag(getAutoTagInput(paths));
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

  async function onPluginTaskClicked(plugin: Plugin, operation: PluginTask) {
    await mutateRunPluginTask(plugin.id, operation.name);
    Toast.success({
      content: intl.formatMessage(
        { id: "config.tasks.added_job_to_queue" },
        { operation_name: operation.name }
      ),
    });
  }

  function renderPluginTasks(plugin: Plugin, pluginTasks: PluginTask[]) {
    if (!pluginTasks) {
      return;
    }

    return pluginTasks.map((o) => {
      return (
        <div key={o.name}>
          <Button
            onClick={() => onPluginTaskClicked(plugin, o)}
            className="mt-3"
            variant="secondary"
            size="sm"
          >
            {o.name}
          </Button>
          {o.description ? (
            <Form.Text className="text-muted">{o.description}</Form.Text>
          ) : undefined}
        </div>
      );
    });
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

  function renderPlugins() {
    if (!plugins.data || !plugins.data.plugins) {
      return;
    }

    const taskPlugins = plugins.data.plugins.filter(
      (p) => p.tasks && p.tasks.length > 0
    );

    return (
      <>
        <hr />
        <h5>{intl.formatMessage({ id: "config.tasks.plugin_tasks" })}</h5>
        {taskPlugins.map((o) => {
          return (
            <div key={`${o.id}`} className="mb-3">
              <h6>{o.name}</h6>
              {renderPluginTasks(o, o.tasks ?? [])}
              <hr />
            </div>
          );
        })}
      </>
    );
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

  if (isBackupRunning) {
    return (
      <LoadingIndicator
        message={intl.formatMessage({ id: "config.tasks.backing_up_database" })}
      />
    );
  }

  return (
    <>
      {renderImportAlert()}
      {renderCleanAlert()}
      {renderImportDialog()}
      {renderScanDialog()}
      {renderAutoTagDialog()}

      <h4>{intl.formatMessage({ id: "config.tasks.job_queue" })}</h4>

      <JobTable />

      <hr />

      <h5>{intl.formatMessage({ id: "library" })}</h5>
      <Form.Group>
        <Form.Check
          id="use-file-metadata"
          checked={useFileMetadata}
          label={intl.formatMessage({
            id: "config.tasks.set_name_date_details_from_metadata_if_present",
          })}
          onChange={() => setUseFileMetadata(!useFileMetadata)}
        />
        <Form.Check
          id="strip-file-extension"
          checked={stripFileExtension}
          label={intl.formatMessage({
            id: "config.tasks.dont_include_file_extension_as_part_of_the_title",
          })}
          onChange={() => setStripFileExtension(!stripFileExtension)}
        />
        <Form.Check
          id="scan-generate-previews"
          checked={scanGeneratePreviews}
          label={intl.formatMessage({
            id: "config.tasks.generate_video_previews_during_scan",
          })}
          onChange={() => setScanGeneratePreviews(!scanGeneratePreviews)}
        />
        <div className="d-flex flex-row">
          <div>↳</div>
          <Form.Check
            id="scan-generate-image-previews"
            checked={scanGenerateImagePreviews}
            disabled={!scanGeneratePreviews}
            label={intl.formatMessage({
              id: "config.tasks.generate_previews_during_scan",
            })}
            onChange={() =>
              setScanGenerateImagePreviews(!scanGenerateImagePreviews)
            }
            className="ml-2 flex-grow"
          />
        </div>
        <Form.Check
          id="scan-generate-sprites"
          checked={scanGenerateSprites}
          label={intl.formatMessage({
            id: "config.tasks.generate_sprites_during_scan",
          })}
          onChange={() => setScanGenerateSprites(!scanGenerateSprites)}
        />
        <Form.Check
          id="scan-generate-phashes"
          checked={scanGeneratePhashes}
          label={intl.formatMessage({
            id: "config.tasks.generate_phashes_during_scan",
          })}
          onChange={() => setScanGeneratePhashes(!scanGeneratePhashes)}
        />
        <Form.Check
          id="scan-generate-thumbnails"
          checked={scanGenerateThumbnails}
          label={intl.formatMessage({
            id: "config.tasks.generate_thumbnails_during_scan",
          })}
          onChange={() => setScanGenerateThumbnails(!scanGenerateThumbnails)}
        />
      </Form.Group>
      <Form.Group>
        <Button
          className="mr-2"
          variant="secondary"
          type="submit"
          onClick={() => onScan()}
        >
          <FormattedMessage id="actions.scan" />
        </Button>
        <Button
          variant="secondary"
          type="submit"
          onClick={() => setIsScanDialogOpen(true)}
        >
          <FormattedMessage id="actions.selective_scan" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.scan_for_content_desc" })}
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>{intl.formatMessage({ id: "config.tasks.auto_tagging" })}</h5>

      <Form.Group>
        <Form.Check
          id="autotag-performers"
          checked={autoTagPerformers}
          label={intl.formatMessage({ id: "performers" })}
          onChange={() => setAutoTagPerformers(!autoTagPerformers)}
        />
        <Form.Check
          id="autotag-studios"
          checked={autoTagStudios}
          label={intl.formatMessage({ id: "studios" })}
          onChange={() => setAutoTagStudios(!autoTagStudios)}
        />
        <Form.Check
          id="autotag-tags"
          checked={autoTagTags}
          label={intl.formatMessage({ id: "tags" })}
          onChange={() => setAutoTagTags(!autoTagTags)}
        />
      </Form.Group>
      <Form.Group>
        <Button
          variant="secondary"
          type="submit"
          className="mr-2"
          onClick={() => onAutoTag()}
        >
          <FormattedMessage id="actions.auto_tag" />
        </Button>
        <Button
          variant="secondary"
          type="submit"
          onClick={() => setIsAutoTagDialogOpen(true)}
        >
          <FormattedMessage id="actions.selective_auto_tag" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "config.tasks.auto_tag_based_on_filenames",
          })}
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>{intl.formatMessage({ id: "config.tasks.generated_content" })}</h5>
      <GenerateButton />

      <hr />
      <h5>{intl.formatMessage({ id: "config.tasks.maintenance" })}</h5>
      <Form.Group>
        <Form.Check
          id="clean-dryrun"
          checked={cleanDryRun}
          label={intl.formatMessage({ id: "config.tasks.only_dry_run" })}
          onChange={() => setCleanDryRun(!cleanDryRun)}
        />
      </Form.Group>
      <Form.Group>
        <Button
          id="clean"
          variant="danger"
          onClick={() => setIsCleanAlertOpen(true)}
        >
          <FormattedMessage id="actions.clean" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.cleanup_desc" })}
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>{intl.formatMessage({ id: "metadata" })}</h5>
      <Form.Group>
        <Button
          id="export"
          variant="secondary"
          type="submit"
          onClick={() => onExport()}
        >
          <FormattedMessage id="actions.full_export" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.export_to_json" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Button
          id="import"
          variant="danger"
          onClick={() => setIsImportAlertOpen(true)}
        >
          <FormattedMessage id="actions.full_import" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.import_from_exported_json" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Button
          id="partial-import"
          variant="danger"
          onClick={() => setIsImportDialogOpen(true)}
        >
          <FormattedMessage id="actions.import_from_file" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.incremental_import" })}
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>{intl.formatMessage({ id: "actions.backup" })}</h5>
      <Form.Group>
        <Button
          id="backup"
          variant="secondary"
          type="submit"
          onClick={() => onBackup()}
        >
          <FormattedMessage id="actions.backup" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage(
            { id: "config.tasks.backup_database" },
            {
              filename_format: (
                <code>
                  [origFilename].sqlite.[schemaVersion].[YYYYMMDD_HHMMSS]
                </code>
              ),
            }
          )}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Button
          id="backupDownload"
          variant="secondary"
          type="submit"
          onClick={() => onBackup(true)}
        >
          <FormattedMessage id="actions.download_backup" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.backup_and_download" })}
        </Form.Text>
      </Form.Group>

      {renderPlugins()}

      <hr />

      <h5>{intl.formatMessage({ id: "config.tasks.migrations" })}</h5>

      <Form.Group>
        <Button
          id="migrateHashNaming"
          variant="danger"
          onClick={() => onMigrateHashNaming()}
        >
          <FormattedMessage id="actions.rename_gen_files" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.migrate_hash_files" })}
        </Form.Text>
      </Form.Group>
    </>
  );
};
