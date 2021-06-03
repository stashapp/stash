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
      Toast.success({ content: "Added import task to queue" });
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
        <p>
          Are you sure you want to import? This will delete the database and
          re-import from your exported metadata.
        </p>
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
        <p>
          Dry Mode selected. No actual deleting will take place, only logging.
        </p>
      );
    } else {
      msg = (
        <p>
          Are you sure you want to Clean? This will delete database information
          and generated content for all scenes and galleries that are no longer
          found in the filesystem.
        </p>
      );
    }

    return (
      <Modal
        show={isCleanAlertOpen}
        icon="trash-alt"
        accept={{ text: "Clean", variant: "danger", onClick: onClean }}
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
      });
      Toast.success({ content: "Added scan to job queue" });
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
      Toast.success({ content: "Added Auto tagging job to queue" });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onPluginTaskClicked(plugin: Plugin, operation: PluginTask) {
    await mutateRunPluginTask(plugin.id, operation.name);
    Toast.success({ content: `Added ${operation.name} job to queue` });
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

    return (
      <>
        <hr />
        <h5>Plugin Tasks</h5>
        {plugins.data.plugins.map((o) => {
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
      Toast.success({ content: "Added hash migration task to queue" });
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onExport() {
    try {
      await mutateMetadataExport();
      Toast.success({ content: "Added export task to queue" });
    } catch (err) {
      Toast.error(err);
    }
  }

  if (isBackupRunning) {
    return <LoadingIndicator message="Backup up database" />;
  }

  return (
    <>
      {renderImportAlert()}
      {renderCleanAlert()}
      {renderImportDialog()}
      {renderScanDialog()}
      {renderAutoTagDialog()}

      <h4>Job Queue</h4>

      <JobTable />

      <hr />

      <h5>Library</h5>
      <Form.Group>
        <Form.Check
          id="use-file-metadata"
          checked={useFileMetadata}
          label="Set name, date, details from metadata (if present)"
          onChange={() => setUseFileMetadata(!useFileMetadata)}
        />
        <Form.Check
          id="strip-file-extension"
          checked={stripFileExtension}
          label="Don't include file extension as part of the title"
          onChange={() => setStripFileExtension(!stripFileExtension)}
        />
        <Form.Check
          id="scan-generate-previews"
          checked={scanGeneratePreviews}
          label="Generate previews during scan (video previews which play when hovering over a scene)"
          onChange={() => setScanGeneratePreviews(!scanGeneratePreviews)}
        />
        <div className="d-flex flex-row">
          <div>↳</div>
          <Form.Check
            id="scan-generate-image-previews"
            checked={scanGenerateImagePreviews}
            disabled={!scanGeneratePreviews}
            label="Generate image previews during scan (animated WebP previews, only required if Preview Type is set to Animated Image)"
            onChange={() =>
              setScanGenerateImagePreviews(!scanGenerateImagePreviews)
            }
            className="ml-2 flex-grow"
          />
        </div>
        <Form.Check
          id="scan-generate-sprites"
          checked={scanGenerateSprites}
          label="Generate sprites during scan (for the scene scrubber)"
          onChange={() => setScanGenerateSprites(!scanGenerateSprites)}
        />
        <Form.Check
          id="scan-generate-phashes"
          checked={scanGeneratePhashes}
          label="Generate phashes during scan (for deduplication and scene identification)"
          onChange={() => setScanGeneratePhashes(!scanGeneratePhashes)}
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
          Scan for new content and add it to the database.
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>Auto Tagging</h5>

      <Form.Group>
        <Form.Check
          id="autotag-performers"
          checked={autoTagPerformers}
          label="Performers"
          onChange={() => setAutoTagPerformers(!autoTagPerformers)}
        />
        <Form.Check
          id="autotag-studios"
          checked={autoTagStudios}
          label="Studios"
          onChange={() => setAutoTagStudios(!autoTagStudios)}
        />
        <Form.Check
          id="autotag-tags"
          checked={autoTagTags}
          label="Tags"
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
          Auto-tag content based on filenames.
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>Generated Content</h5>
      <GenerateButton />

      <hr />
      <h5>Maintenance</h5>
      <Form.Group>
        <Form.Check
          id="clean-dryrun"
          checked={cleanDryRun}
          label="Only perform a dry run. Don't remove anything"
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
          Check for missing files and remove them from the database. This is a
          destructive action.
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>Metadata</h5>
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
          Exports the database content into JSON format in the metadata
          directory.
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
          Import from exported JSON in the metadata directory. Wipes the
          existing database.
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
          Incremental import from a supplied export zip file.
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>Backup</h5>
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
          Performs a backup of the database to the same directory as the
          database, with the filename format{" "}
          <code>[origFilename].sqlite.[schemaVersion].[YYYYMMDD_HHMMSS]</code>
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
          Performs a backup of the database and downloads the resulting file.
        </Form.Text>
      </Form.Group>

      {renderPlugins()}

      <hr />

      <h5>Migrations</h5>

      <Form.Group>
        <Button
          id="migrateHashNaming"
          variant="danger"
          onClick={() => onMigrateHashNaming()}
        >
          <FormattedMessage id="actions.rename_gen_files" />
        </Button>
        <Form.Text className="text-muted">
          Used after changing the Generated file naming hash to rename existing
          generated files to the new hash format.
        </Form.Text>
      </Form.Group>
    </>
  );
};
