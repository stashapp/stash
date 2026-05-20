import React, { useState } from "react";
import { FormattedMessage } from "react-intl";
import { Button, Modal } from "react-bootstrap";
import TextUtils from "src/utils/text";
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
    TextUtils.stripQuotes(currentValue ?? "")
  );
  const handleChangeDirectory = (v: string) =>
    setCurrentDirectory(TextUtils.stripQuotes(v));

  return (
    <Modal show onHide={() => onClose()} title="">
      <Modal.Header>
        <FormattedMessage id="actions.select_directory" />
      </Modal.Header>
      <Modal.Body>
        <div className="dialog-content">
          <FolderSelect
            currentDirectory={currentDirectory}
            onChangeDirectory={handleChangeDirectory}
          />
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={() => onClose()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
        <Button
          variant="success"
          onClick={() => onClose(TextUtils.stripQuotes(currentDirectory))}
        >
          <FormattedMessage id="actions.confirm" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
