import React from "react";
import { Button, Modal } from "react-bootstrap";
import { FormattedMessage } from "react-intl";

export interface IAlertModalProps {
  text: JSX.Element | string;
  show?: boolean;
  confirmButtonText?: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export const AlertModal: React.FC<IAlertModalProps> = ({
  text,
  show,
  confirmButtonText,
  onConfirm,
  onCancel,
}) => {
  return (
    <Modal show={show}>
      <Modal.Body>{text}</Modal.Body>
      <Modal.Footer>
        <Button variant="danger" onClick={() => onConfirm()}>
          {confirmButtonText ?? <FormattedMessage id="actions.confirm" />}
        </Button>
        <Button variant="secondary" onClick={() => onCancel()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
