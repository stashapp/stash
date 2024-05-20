import { Button, Modal } from "react-bootstrap";
import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { ImageInput } from "./ImageInput";
import cx from "classnames";

interface IProps {
  objectName?: string;
  isNew: boolean;
  isEditing: boolean;
  onToggleEdit: () => void;
  onSave: () => void;
  saveDisabled?: boolean;
  onDelete: () => void;
  onAutoTag?: () => void;
  autoTagDisabled?: boolean;
  onImageChange: (event: React.FormEvent<HTMLInputElement>) => void;
  onBackImageChange?: (event: React.FormEvent<HTMLInputElement>) => void;
  onImageChangeURL?: (url: string) => void;
  onBackImageChangeURL?: (url: string) => void;
  onClearImage?: () => void;
  onClearBackImage?: () => void;
  acceptSVG?: boolean;
  customButtons?: JSX.Element;
  classNames?: string;
  children?: JSX.Element | JSX.Element[];
}

export const DetailsEditNavbar: React.FC<IProps> = (props: IProps) => {
  const intl = useIntl();
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  function renderEditButton() {
    if (props.isNew) return;
    return (
      <Button
        variant="primary"
        className="edit"
        onClick={() => props.onToggleEdit()}
      >
        {props.isEditing
          ? intl.formatMessage({ id: "actions.cancel" })
          : intl.formatMessage({ id: "actions.edit" })}
      </Button>
    );
  }

  function renderSaveButton() {
    if (!props.isEditing) return;

    return (
      <Button
        variant="success"
        className="save"
        disabled={props.saveDisabled}
        onClick={() => props.onSave()}
      >
        <FormattedMessage id="actions.save" />
      </Button>
    );
  }

  function renderDeleteButton() {
    if (props.isNew || props.isEditing) return;
    return (
      <Button
        variant="danger"
        className="delete"
        onClick={() => setIsDeleteAlertOpen(true)}
      >
        <FormattedMessage id="actions.delete" />
      </Button>
    );
  }

  function renderBackImageInput() {
    if (!props.isEditing || !props.onBackImageChange) {
      return;
    }
    return (
      <ImageInput
        isEditing={props.isEditing}
        text={intl.formatMessage({ id: "actions.set_back_image" })}
        onImageChange={props.onBackImageChange}
        onImageURL={props.onBackImageChangeURL}
      />
    );
  }

  function renderAutoTagButton() {
    if (props.isNew || props.isEditing) return;

    if (props.onAutoTag) {
      return (
        <div>
          <Button
            variant="secondary"
            disabled={props.autoTagDisabled}
            onClick={() => {
              if (props.onAutoTag) {
                props.onAutoTag();
              }
            }}
          >
            <FormattedMessage id="actions.auto_tag" />
          </Button>
        </div>
      );
    }
  }

  function renderDeleteAlert() {
    return (
      <Modal show={isDeleteAlertOpen}>
        <Modal.Body>
          <FormattedMessage
            id="dialogs.delete_confirm"
            values={{ entityName: props.objectName }}
          />
        </Modal.Body>
        <Modal.Footer>
          <Button variant="danger" onClick={props.onDelete}>
            <FormattedMessage id="actions.delete" />
          </Button>
          <Button
            variant="secondary"
            onClick={() => setIsDeleteAlertOpen(false)}
          >
            <FormattedMessage id="actions.cancel" />
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }

  return (
    <div className={cx("details-edit", props.classNames)}>
      {renderEditButton()}
      <ImageInput
        isEditing={props.isEditing}
        text={
          props.onBackImageChange
            ? intl.formatMessage({ id: "actions.set_front_image" })
            : undefined
        }
        onImageChange={props.onImageChange}
        onImageURL={props.onImageChangeURL}
        acceptSVG={props.acceptSVG ?? false}
      />
      {props.isEditing && props.onClearImage ? (
        <div>
          <Button
            className="mr-2"
            variant="danger"
            onClick={() => props.onClearImage!()}
          >
            {props.onClearBackImage
              ? intl.formatMessage({ id: "actions.clear_front_image" })
              : intl.formatMessage({ id: "actions.clear_image" })}
          </Button>
        </div>
      ) : null}
      {renderBackImageInput()}
      {props.isEditing && props.onClearBackImage ? (
        <div>
          <Button
            className="mr-2"
            variant="danger"
            onClick={() => props.onClearBackImage!()}
          >
            {intl.formatMessage({ id: "actions.clear_back_image" })}
          </Button>
        </div>
      ) : null}
      {renderAutoTagButton()}
      {props.customButtons}
      {renderSaveButton()}
      {renderDeleteButton()}
      {renderDeleteAlert()}
    </div>
  );
};
