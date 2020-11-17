import React, { useState } from "react";
import { Button, Modal } from "react-bootstrap";
import { FolderSelect } from "./FolderSelect";

interface IProps {
  onClose: (directory?: string) => void;
}

export const FolderSelectDialog: React.FC<IProps> = (props: IProps) => {
  const [currentDirectory, setCurrentDirectory] = useState<string>("");

  return (
    <Modal show onHide={() => props.onClose()} title="">
      <Modal.Header>Select Directory</Modal.Header>
      <Modal.Body>
        <div className="dialog-content">
          <FolderSelect
            currentDirectory={currentDirectory}
            setCurrentDirectory={(v) => setCurrentDirectory(v)}
          />
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button
          variant="success"
          onClick={() => props.onClose(currentDirectory)}
        >
          Add
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
