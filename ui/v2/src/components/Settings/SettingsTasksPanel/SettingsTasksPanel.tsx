import {
  Alert,
  Button,
  Checkbox,
  Divider,
  FormGroup,
  H4,
  AnchorButton,
  ProgressBar,
  H5,
} from "@blueprintjs/core";
import React, { FunctionComponent, useState, useEffect } from "react";
import { StashService } from "../../../core/StashService";
import { ErrorUtils } from "../../../utils/errors";
import { ToastUtils } from "../../../utils/toasts";
import { GenerateButton } from "./GenerateButton";
import { Link } from "react-router-dom";

interface IProps {}

export const SettingsTasksPanel: FunctionComponent<IProps> = (props: IProps) => {
  const [isImportAlertOpen, setIsImportAlertOpen] = useState<boolean>(false);
  const [isCleanAlertOpen, setIsCleanAlertOpen] = useState<boolean>(false);
  const [nameFromMetadata, setNameFromMetadata] = useState<boolean>(true);
  const [status, setStatus] = useState<string>("");
  const [progress, setProgress] = useState<number | undefined>(undefined);

  const [autoTagPerformers, setAutoTagPerformers] = useState<boolean>(true);
  const [autoTagStudios, setAutoTagStudios] = useState<boolean>(true);
  const [autoTagTags, setAutoTagTags] = useState<boolean>(true);

  const jobStatus = StashService.useJobStatus();
  const metadataUpdate = StashService.useMetadataUpdate();

  function statusToText(status : string) {
    switch(status) {
      case "Idle":
        return "Idle";
      case "Scan":
        return "Scanning for new content";
      case "Generate":
        return "Generating supporting files";
      case "Clean":
        return "Cleaning the database";
      case "Export":
        return "Exporting to JSON";
      case "Import":
        return "Importing from JSON";
      case "Auto Tag":
        return "Auto tagging scenes";
    }

    return "Idle";
  }

  useEffect(() => {
    if (!!jobStatus.data && !!jobStatus.data.jobStatus) {
      setStatus(statusToText(jobStatus.data.jobStatus.status));
      var newProgress = jobStatus.data.jobStatus.progress;
      if (newProgress < 0) {
        setProgress(undefined);
      } else {
        setProgress(newProgress);
      }
    }
  }, [jobStatus.data]);

  useEffect(() => {
    if (!!metadataUpdate.data && !!metadataUpdate.data.metadataUpdate) {
      setStatus(statusToText(metadataUpdate.data.metadataUpdate.status));
      var newProgress = metadataUpdate.data.metadataUpdate.progress;
      if (newProgress < 0) {
        setProgress(undefined);
      } else {
        setProgress(newProgress);
      }
    }
  }, [metadataUpdate.data]);

  function onImport() {
    setIsImportAlertOpen(false);
    StashService.queryMetadataImport().then(() => { jobStatus.refetch()});
  }

  function renderImportAlert() {
    return (
      <Alert
        cancelButtonText="Cancel"
        confirmButtonText="Import"
        icon="trash"
        intent="danger"
        isOpen={isImportAlertOpen}
        onCancel={() => setIsImportAlertOpen(false)}
        onConfirm={() => onImport()}
      >
        <p>
          Are you sure you want to import?  This will delete the database and re-import from
          your exported metadata.
        </p>
      </Alert>
    );
  }

  function onClean() {
    setIsCleanAlertOpen(false);
    StashService.queryMetadataClean().then(() => { jobStatus.refetch()});
  }

  function renderCleanAlert() {
    return (
      <Alert
        cancelButtonText="Cancel"
        confirmButtonText="Clean"
        icon="trash"
        intent="danger"
        isOpen={isCleanAlertOpen}
        onCancel={() => setIsCleanAlertOpen(false)}
        onConfirm={() => onClean()}
      >
        <p>
          Are you sure you want to Clean?
          This will delete db information and generated content
          for all scenes that are no longer found in the filesystem.
        </p>
      </Alert>
    );
  }

  async function onScan() {
    try {
      await StashService.queryMetadataScan({nameFromMetadata});
      ToastUtils.success("Started scan");
      jobStatus.refetch();
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  function getAutoTagInput() {
    var wildcard = ["*"];
    return {
      performers: autoTagPerformers ? wildcard : [],
      studios: autoTagStudios ? wildcard : [],
      tags: autoTagTags ? wildcard : []
    }
  }

  async function onAutoTag() {
    try {
      await StashService.queryMetadataAutoTag(getAutoTagInput());
      ToastUtils.success("Started auto tagging");
      jobStatus.refetch();
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  function maybeRenderStop() {
    if (!status || status === "Idle") {
      return undefined;
    }

    return (
      <>
      <FormGroup>
        <Button id="stop" text="Stop" intent="danger" onClick={() => StashService.queryStopJob().then(() => jobStatus.refetch())} />
      </FormGroup>
      </>
    );
  }

  function renderJobStatus() {
    return (
      <>
      <FormGroup>
        <H5>Status: {status}</H5>
        {!!status && status !== "Idle" ? <ProgressBar value={progress}/> : undefined}
      </FormGroup>
      {maybeRenderStop()}
      </>
    );
  }

  return (
    <>
      {renderImportAlert()}
      {renderCleanAlert()}

      <H4>Running Jobs</H4>

      {renderJobStatus()}

      <Divider/>

      <H4>Library</H4>
      <FormGroup
        helperText="Scan for new content and add it to the database."
        labelFor="scan"
        inline={true}
      >
        <Checkbox
          checked={nameFromMetadata}
          label="Set name from metadata (if present)"
          onChange={() => setNameFromMetadata(!nameFromMetadata)}
        />
        <Button id="scan" text="Scan" onClick={() => onScan()} />
      </FormGroup>

      <Divider />

      <H4>Auto Tagging</H4>

      <FormGroup
        helperText="Auto-tag content based on filenames."
        labelFor="autoTag"
        inline={true}
      >
        <Checkbox
          checked={autoTagPerformers}
          label="Performers"
          onChange={() => setAutoTagPerformers(!autoTagPerformers)}
        />
        <Checkbox
          checked={autoTagStudios}
          label="Studios"
          onChange={() => setAutoTagStudios(!autoTagStudios)}
        />
        <Checkbox
          checked={autoTagTags}
          label="Tags"
          onChange={() => setAutoTagTags(!autoTagTags)}
        />
        <Button id="autoTag" text="Auto Tag" onClick={() => onAutoTag()} />
      </FormGroup>

      <FormGroup>
        <Link className="bp3-button" to={"/sceneFilenameParser"}>
          Scene Filename Parser
        </Link>
      </FormGroup>
      <Divider />

      <H4>Generated Content</H4>
      <GenerateButton />
      <FormGroup
        helperText="Check for missing files and remove them from the database. This is a destructive action."
        labelFor="clean"
        inline={true}
      >
        <Button id="clean" text="Clean" intent="danger" onClick={() => setIsCleanAlertOpen(true)} />
      </FormGroup>
      <Divider />

      <H4>Metadata</H4>
      <FormGroup
        helperText="Export the database content into JSON format"
        labelFor="export"
        inline={true}
      >
        <Button id="export" text="Export" onClick={() => StashService.queryMetadataExport().then(() => { jobStatus.refetch()})} />
      </FormGroup>

      <FormGroup
        helperText="Import from exported JSON.  This is a destructive action."
        labelFor="import"
        inline={true}
      >
        <Button id="import" text="Import" intent="danger" onClick={() => setIsImportAlertOpen(true)} />
      </FormGroup>
    </>
  );
};
