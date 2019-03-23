import {
  Button,
  Classes,
  Dialog,
  InputGroup,
  Spinner,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useState } from "react";
import { StashService } from "../../../core/StashService";

interface IProps {
  directories: string[];
  onDirectoriesChanged: (directories: string[]) => void;
}

export const FolderSelect: FunctionComponent<IProps> = (props: IProps) => {
  const [currentDirectory, setCurrentDirectory] = useState<string>("");
  const [isDisplayingDialog, setIsDisplayingDialog] = useState<boolean>(false);
  const [selectableDirectories, setSelectableDirectories] = useState<string[]>([]);
  const [selectedDirectories, setSelectedDirectories] = useState<string[]>([]);
  const { data, error, loading } = StashService.useDirectories(currentDirectory);

  useEffect(() => {
    setSelectedDirectories(props.directories);
  }, [props.directories]);

  useEffect(() => {
    if (!data || !data.directories || !!error) { return; }
    setSelectableDirectories(StashService.nullToUndefined(data.directories));
  }, [data]);

  function onSelectDirectory() {
    selectedDirectories.push(currentDirectory);
    setSelectedDirectories(selectedDirectories);
    setCurrentDirectory("");
    setIsDisplayingDialog(false);
    props.onDirectoriesChanged(selectedDirectories);
  }

  function onRemoveDirectory(directory: string) {
    const newSelectedDirectories = selectedDirectories.filter((dir) => dir !== directory);
    setSelectedDirectories(newSelectedDirectories);
    props.onDirectoriesChanged(newSelectedDirectories);
  }

  function renderDialog() {
    return (
      <Dialog
        isOpen={isDisplayingDialog}
        onClose={() => setIsDisplayingDialog(false)}
        title="Select Directory"
      >
        <div className="dialog-content">
          <InputGroup
            large={true}
            placeholder="File path"
            onChange={(e: any) => setCurrentDirectory(e.target.value)}
            value={currentDirectory}
            rightElement={(!data || !data.directories || loading) ? <Spinner size={Spinner.SIZE_SMALL} /> : undefined}
          />
          {selectableDirectories.map((path) => {
            return <div key={path} onClick={() => setCurrentDirectory(path)}>{path}</div>;
          })}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={() => onSelectDirectory()}>Add</Button>
          </div>
        </div>
      </Dialog>
    );
  }

  return (
    <>
      {!!error ? error : undefined}
      {renderDialog()}
      {selectedDirectories.map((path) => {
        return <div key={path}>{path} <a onClick={() => onRemoveDirectory(path)}>Remove</a></div>;
      })}
      <Button small={true} onClick={() => setIsDisplayingDialog(true)}>Add Directory</Button>
    </>
  );
};
