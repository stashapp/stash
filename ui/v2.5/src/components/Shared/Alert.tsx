import React from "react";
import { Button, Modal } from "react-bootstrap";
import { FormattedMessage } from "react-intl";

export interface IAlertModalProps {
  text: JSX.Element | string;
  confirmVariant?: string;
  show?: boolean;
  confirmButtonText?: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export const AlertModal: React.FC<IAlertModalProps> = ({
  text,
  show,
  confirmVariant = "danger",
  confirmButtonText,
  onConfirm,
  onCancel,
}) => {
  return (
    <Modal show={show}>
      <Modal.Body>{text}</Modal.Body>
      <Modal.Footer>
        <Button variant={confirmVariant} onClick={() => onConfirm()}>
          {confirmButtonText ?? <FormattedMessage id="actions.confirm" />}
        </Button>
        <Button variant="secondary" onClick={() => onCancel()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
