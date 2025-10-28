import React, { useState } from "react";
import { FormattedMessage } from "react-intl";
import { Button, Modal } from "react-bootstrap";
import { FolderSelect } from "./FolderSelect";

interface IProps {
  defaultValue?: string;
  onClose: (directory?: string) => void;
}

export const FolderSelectDialog: React.FC<IProps> = ({
  defaultValue: currentValue,
  onClose,
}) => {
  const [currentDirectory, setCurrentDirectory] = useState<string>(
    currentValue ?? ""
  );

  return (
    <Modal show onHide={() => onClose()} title="">
      <Modal.Header>Select Directory</Modal.Header>
      <Modal.Body>
        <div className="dialog-content">
          <FolderSelect
            currentDirectory={currentDirectory}
            onChangeDirectory={setCurrentDirectory}
          />
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={() => onClose()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
        <Button variant="success" onClick={() => onClose(currentDirectory)}>
          <FormattedMessage id="actions.confirm" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
