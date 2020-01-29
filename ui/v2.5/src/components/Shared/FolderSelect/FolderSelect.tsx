import React, { useEffect, useState } from "react";
import { Button, InputGroup, Form, Modal } from "react-bootstrap";
import { LoadingIndicator } from 'src/components/Shared';
import { StashService } from "src/core/StashService";

interface IProps {
  directories: string[];
  onDirectoriesChanged: (directories: string[]) => void;
}

export const FolderSelect: React.FC<IProps> = (props: IProps) => {
  const [currentDirectory, setCurrentDirectory] = useState<string>("");
  const [isDisplayingDialog, setIsDisplayingDialog] = useState<boolean>(false);
  const [selectedDirectories, setSelectedDirectories] = useState<string[]>([]);
  const { data, error, loading } = StashService.useDirectories(
    currentDirectory
  );

  useEffect(() => {
    setSelectedDirectories(props.directories);
  }, [props.directories]);

  const selectableDirectories: string[] = data?.directories ?? [];

  function onSelectDirectory() {
    selectedDirectories.push(currentDirectory);
    setSelectedDirectories(selectedDirectories);
    setCurrentDirectory("");
    setIsDisplayingDialog(false);
    props.onDirectoriesChanged(selectedDirectories);
  }

  function onRemoveDirectory(directory: string) {
    const newSelectedDirectories = selectedDirectories.filter(
      dir => dir !== directory
    );
    setSelectedDirectories(newSelectedDirectories);
    props.onDirectoriesChanged(newSelectedDirectories);
  }

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
                onChange={(e: any) => setCurrentDirectory(e.target.value)}
                defaultValue={currentDirectory}
              />
              <InputGroup.Append>
                {!data || !data.directories || loading ? (
                  <LoadingIndicator inline />
                ) : (
                  ""
                )}
              </InputGroup.Append>
            </InputGroup>
            <ul className="folder-list">
              {selectableDirectories.map(path => {
                return (
                  <li className="folder-item">
                    <Button
                      variant="link"
                      key={path}
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
          <Button variant="success" onClick={() => onSelectDirectory()}>Add</Button>
        </Modal.Footer>
      </Modal>
    );
  }

  return (
    <>
      {error ? <h1>{error.message}</h1> : ""}
      {renderDialog()}
      <Form.Group>
        {selectedDirectories.map(path => {
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

      <Button variant="secondary" onClick={() => setIsDisplayingDialog(true)}>Add Directory</Button>
    </>
  );
};
