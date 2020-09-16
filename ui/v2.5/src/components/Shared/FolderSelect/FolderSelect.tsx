import React, { useEffect, useState } from "react";
import { FormattedMessage } from "react-intl";
import { Button, InputGroup, Form, Modal } from "react-bootstrap";
import { LoadingIndicator } from "src/components/Shared";
import { useDirectory } from "src/core/StashService";

interface IProps {
  onClose: (directory?: string) => void;
}

export const FolderSelect: React.FC<IProps> = (props: IProps) => {
  const [currentDirectory, setCurrentDirectory] = useState<string>("");
  const { data, error, loading } = useDirectory(currentDirectory);

  useEffect(() => {
    if (currentDirectory === "" && data?.directory.path)
      setCurrentDirectory(data.directory.path);
  }, [currentDirectory, data]);

  const selectableDirectories: string[] = data?.directory.directories ?? [];

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

  return (
    <Modal
      show
      onHide={() => props.onClose()}
      title=""
    >
      <Modal.Header>Select Directory</Modal.Header>
      <Modal.Body>
        <div className="dialog-content">
          {error ? <h1>{error.message}</h1> : ""}
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
        <Button variant="success" onClick={() => props.onClose(currentDirectory)}>
          Add
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
