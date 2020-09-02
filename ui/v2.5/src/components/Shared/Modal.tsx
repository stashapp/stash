import React from "react";
import { Button, Modal, Spinner, ModalProps } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { IconName } from "@fortawesome/fontawesome-svg-core";

interface IButton {
  text?: string;
  variant?: "danger" | "primary" | "secondary";
  onClick?: () => void;
}

interface IModal {
  show: boolean;
  onHide?: () => void;
  header?: string;
  icon?: IconName;
  cancel?: IButton;
  accept?: IButton;
  isRunning?: boolean;
  modalProps?: ModalProps;
}

const ModalComponent: React.FC<IModal> = ({
  children,
  show,
  icon,
  header,
  cancel,
  accept,
  onHide,
  isRunning,
  modalProps,
}) => (
  <Modal keyboard={false} onHide={onHide} show={show} {...modalProps}>
    <Modal.Header>
      {icon ? <Icon icon={icon} /> : ""}
      <span>{header ?? ""}</span>
    </Modal.Header>
    <Modal.Body>{children}</Modal.Body>
    <Modal.Footer>
      <div>
        {cancel ? (
          <Button
            disabled={isRunning}
            variant={cancel.variant ?? "primary"}
            onClick={cancel.onClick}
          >
            {cancel.text ?? "Cancel"}
          </Button>
        ) : (
          ""
        )}
        <Button
          disabled={isRunning}
          variant={accept?.variant ?? "primary"}
          onClick={accept?.onClick}
        >
          {isRunning ? (
            <Spinner animation="border" role="status" size="sm" />
          ) : (
            accept?.text ?? "Close"
          )}
        </Button>
      </div>
    </Modal.Footer>
  </Modal>
);

export default ModalComponent;
