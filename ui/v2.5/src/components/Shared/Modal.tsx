import React from "react";
import { Button, Modal } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { IconName } from "@fortawesome/fontawesome-svg-core";

interface IButton {
  text?: string;
  variant?: "danger" | "primary";
  onClick?: () => void;
}

interface IModal {
  show: boolean;
  onHide?: () => void;
  header?: string;
  icon?: IconName;
  cancel?: IButton;
  accept?: IButton;
}

const ModalComponent: React.FC<IModal> = ({
  children,
  show,
  icon,
  header,
  cancel,
  accept,
  onHide
}) => (
  <Modal keyboard={false} onHide={onHide} show={show}>
    <Modal.Header>
      {icon ? <Icon icon={icon} /> : ""}
      <span>{header ?? ""}</span>
    </Modal.Header>
    <Modal.Body>{children}</Modal.Body>
    <Modal.Footer>
      <div>
        {cancel ? (
          <Button
            variant={cancel.variant ?? "primary"}
            onClick={cancel.onClick}
          >
            {cancel.text ?? "Cancel"}
          </Button>
        ) : (
          ""
        )}
        <Button
          variant={accept?.variant ?? "primary"}
          onClick={accept?.onClick}
        >
          {accept?.text ?? "Close"}
        </Button>
      </div>
    </Modal.Footer>
  </Modal>
);

export default ModalComponent;
