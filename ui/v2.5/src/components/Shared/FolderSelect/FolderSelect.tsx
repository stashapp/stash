import React, { useEffect, useState } from "react";
import { FormattedMessage } from "react-intl";
import { Button, InputGroup, Form, Modal } from "react-bootstrap";
import { LoadingIndicator } from "src/components/Shared";
import { useDirectory } from "src/core/StashService";

interface IProps {
  directories: string[];
  onDirectoriesChanged: (directories: string[]) => void;
}

export const FolderSelect: React.FC<IProps> = (props: IProps) => {
  const [currentDirectory, setCurrentDirectory] = useState<string>("");
  const [isDisplayingDialog, setIsDisplayingDialog] = useState<boolean>(false);
  const [selectedDirectories, setSelectedDirectories] = useState<string[]>([]);
  const { data, error, loading } = useDirectory(currentDirectory);

  useEffect(() => {
    setSelectedDirectories(props.directories);
  }, [props.directories]);

  useEffect(() => {
    if (currentDirectory === "" && data?.directory.path)
      setCurrentDirectory(data.directory.path);
  }, [currentDirectory, data]);

  const selectableDirectories: string[] = data?.directory.directories ?? [];

  function onSelectDirectory() {
    selectedDirectories.push(currentDirectory);
    setSelectedDirectories(selectedDirectories);
    setCurrentDirectory("");
    setIsDisplayingDialog(false);
    props.onDirectoriesChanged(selectedDirectories);
  }

  function onRemoveDirectory(directory: string) {
    const newSelectedDirectories = selectedDirectories.filter(
      (dir) => dir !== directory
    );
    setSelectedDirectories(newSelectedDirectories);
    props.onDirectoriesChanged(newSelectedDirectories);
  }

  const topDirectory = data?.directory?.parent ? (
    <li className="folder-list-parent folder-list-item">
      <Button
        variant="link"
        onClick={() =>
          data.directory.parent && setCurrentDirectory(data.directory.parent)
        }
      >
        <FormattedMessage defaultMessage="Up a directory" id="up-dir" />
      </Button>
    </li>
  ) : null;

  function renderDialog() {
    return (
      <Modal
        show={isDisplayingDialog}
        onHide={() => setIsDisplayingDialog(false)}
        title=""
      >
        <Modal.Header>Select Directory</Modal.Header>
        <Modal.Body>
          <div className="dialog-content">
            <InputGroup>
              <Form.Control
                placeholder="File path"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setCurrentDirectory(e.currentTarget.value)
                }
                value={currentDirectory}
                spellCheck={false}
              />
              <InputGroup.Append>
                {!data || !data.directory || loading ? (
                  <LoadingIndicator inline />
                ) : (
                  ""
                )}
              </InputGroup.Append>
            </InputGroup>
            <ul className="folder-list">
              {topDirectory}
              {selectableDirectories.map((path) => {
                return (
                  <li key={path} className="folder-list-item">
                    <Button
                      variant="link"
                      onClick={() => setCurrentDirectory(path)}
                    >
                      {path}
                    </Button>
                  </li>
                );
              })}
            </ul>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="success" onClick={() => onSelectDirectory()}>
            Add
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }

  return (
    <>
      {error ? <h1>{error.message}</h1> : ""}
      {renderDialog()}
      <Form.Group>
        {selectedDirectories.map((path) => {
          return (
            <div key={path}>
              {path}{" "}
              <Button variant="link" onClick={() => onRemoveDirectory(path)}>
                Remove
              </Button>
            </div>
          );
        })}
      </Form.Group>

      <Button variant="secondary" onClick={() => setIsDisplayingDialog(true)}>
        Add Directory
      </Button>
    </>
  );
};
