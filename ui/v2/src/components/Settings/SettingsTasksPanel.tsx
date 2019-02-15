import {
  Alert,
  Button,
  Divider,
  FormGroup,
  H1,
  H4,
  H6,
  InputGroup,
  Tag,
} from "@blueprintjs/core";
import React, { FunctionComponent, useState } from "react";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { TextUtils } from "../../utils/text";

interface IProps {}

export const SettingsTasksPanel: FunctionComponent<IProps> = (props: IProps) => {
  const [isImportAlertOpen, setIsImportAlertOpen] = useState<boolean>(false);

  function onImport() {
    setIsImportAlertOpen(false);
    StashService.queryMetadataImport();
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

  return (
    <>
      {renderImportAlert()}

      <H4>Library</H4>
      <FormGroup
        helperText="Scan for new content and add it to the database."
        labelFor="scan"
        inline={true}
      >
        <Button id="scan" text="Scan" onClick={() => StashService.queryMetadataScan()} />
      </FormGroup>
      <Divider />

      <H4>Generated Content</H4>
      <FormGroup
        helperText="Generate supporting image, sprite, video, vtt and other files."
        labelFor="generate"
        inline={true}
      >
        <Button id="generate" text="Generate" onClick={() => StashService.queryMetadataGenerate()} />
      </FormGroup>
      <FormGroup
        helperText="TODO"
        labelFor="clean"
        inline={true}
      >
        <Button id="clean" text="Clean" onClick={() => StashService.queryMetadataClean()} />
      </FormGroup>
      <Divider />

      <H4>Metadata</H4>
      <FormGroup
        helperText="Export the database content into JSON format"
        labelFor="export"
        inline={true}
      >
        <Button id="export" text="Export" onClick={() => StashService.queryMetadataExport()} />
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
