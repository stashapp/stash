import { Button, Modal } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { PatchComponent } from "src/patch";

export interface IAlertModalProps {
  text: JSX.Element | string;
  show?: boolean;
  confirmButtonText?: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export const AlertModal = PatchComponent(
  "AlertModal",
  (props: IAlertModalProps) => {
    return (
      <Modal show={props.show}>
        <Modal.Body>{props.text}</Modal.Body>
        <Modal.Footer>
          <Button variant="danger" onClick={() => props.onConfirm()}>
            {props.confirmButtonText ?? (
              <FormattedMessage id="actions.confirm" />
            )}
          </Button>
          <Button variant="secondary" onClick={() => props.onCancel()}>
            <FormattedMessage id="actions.cancel" />
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }
);
