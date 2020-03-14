import { Button, Modal } from "react-bootstrap";
import React, { useState } from "react";
import { ImageInput } from "src/components/Shared";

interface IProps {
  objectName?: string;
  isNew: boolean;
  isEditing: boolean;
  onToggleEdit: () => void;
  onSave: () => void;
  onDelete: () => void;
  onAutoTag?: () => void;
  onImageChange: (event: React.FormEvent<HTMLInputElement>) => void;
  onBackImageChange?: (event: React.FormEvent<HTMLInputElement>) => void;
}

export const DetailsEditNavbar: React.FC<IProps> = (props: IProps) => {
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  function renderEditButton() {
    if (props.isNew) return;
    return (
      <Button
        variant="primary"
        className="edit"
        onClick={() => props.onToggleEdit()}
      >
        {props.isEditing ? "Cancel" : "Edit"}
      </Button>
    );
  }

  function renderSaveButton() {
    if (!props.isEditing) return;

    return (
      <Button variant="success" className="save" onClick={() => props.onSave()}>
        Save
      </Button>
    );
  }

  function renderDeleteButton() {
    if (props.isNew || props.isEditing) return;
    return (
      <Button
        variant="danger"
        className="delete d-none d-sm-block"
        onClick={() => setIsDeleteAlertOpen(true)}
      >
        Delete
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
        text="Back image..."
        onImageChange={props.onBackImageChange}
      />
    );
  }

  function renderAutoTagButton() {
    if (props.isNew || props.isEditing) return;

    if (props.onAutoTag) {
      return (
        <Button
          variant="secondary"
          onClick={() => {
            if (props.onAutoTag) {
              props.onAutoTag();
            }
          }}
        >
          Auto Tag
        </Button>
      );
    }
  }

  function renderDeleteAlert() {
    return (
      <Modal show={isDeleteAlertOpen}>
        <Modal.Body>
          Are you sure you want to delete {props.objectName}?
        </Modal.Body>
        <Modal.Footer>
          <Button variant="danger" onClick={props.onDelete}>
            Delete
          </Button>
          <Button
            variant="secondary"
            onClick={() => setIsDeleteAlertOpen(false)}
          >
            Cancel
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }

  return (
    <div className="details-edit">
      {renderEditButton()}
      <ImageInput
        isEditing={props.isEditing}
        text={props.onBackImageChange ? "Front image..." : undefined}
        onImageChange={props.onImageChange}
      />
      {renderBackImageInput()}
      {renderAutoTagButton()}
      {renderSaveButton()}
      {renderDeleteButton()}
      {renderDeleteAlert()}
    </div>
  );
};
