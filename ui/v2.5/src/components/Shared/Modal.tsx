import React from "react";
import { Button, Modal, Spinner, ModalProps } from "react-bootstrap";
import { ButtonVariant } from "react-bootstrap/types";
import { Icon } from "./Icon";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { FormattedMessage } from "react-intl";

interface IButton {
  text?: string;
  variant?: ButtonVariant;
  onClick?: () => void;
}

interface IModal {
  show: boolean;
  onHide?: () => void;
  header?: JSX.Element | string;
  icon?: IconDefinition;
  cancel?: IButton;
  accept?: IButton;
  isRunning?: boolean;
  disabled?: boolean;
  modalProps?: ModalProps;
  dialogClassName?: string;
  footerButtons?: React.ReactNode;
  leftFooterButtons?: React.ReactNode;
}

const defaultOnHide = () => {};

export const ModalComponent: React.FC<IModal> = ({
  children,
  show,
  icon,
  header,
  cancel,
  accept,
  onHide,
  isRunning,
  disabled,
  modalProps,
  dialogClassName,
  footerButtons,
  leftFooterButtons,
}) => (
  <Modal
    className="ModalComponent"
    keyboard={false}
    onHide={onHide ?? defaultOnHide}
    show={show}
    dialogClassName={dialogClassName}
    {...modalProps}
  >
    <Modal.Header>
      {icon ? <Icon icon={icon} /> : ""}
      <span>{header ?? ""}</span>
    </Modal.Header>
    <Modal.Body>{children}</Modal.Body>
    <Modal.Footer className="ModalFooter">
      <div>{leftFooterButtons}</div>
      <div>
        {footerButtons}
        {cancel ? (
          <Button
            disabled={isRunning}
            variant={cancel.variant ?? "primary"}
            onClick={cancel.onClick}
            className="ml-2"
          >
            {cancel.text ?? (
              <FormattedMessage
                id="actions.cancel"
                defaultMessage="Cancel"
                description="Cancels the current action and dismisses the modal."
              />
            )}
          </Button>
        ) : (
          ""
        )}
        <Button
          disabled={isRunning || disabled}
          variant={accept?.variant ?? "primary"}
          onClick={accept?.onClick}
          className="ml-2"
        >
          {isRunning ? (
            <Spinner animation="border" role="status" size="sm" />
          ) : (
            accept?.text ?? (
              <FormattedMessage
                id="actions.close"
                defaultMessage="Close"
                description="Closes the current modal."
              />
            )
          )}
        </Button>
      </div>
    </Modal.Footer>
  </Modal>
);
