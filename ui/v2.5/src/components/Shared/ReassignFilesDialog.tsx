import React, { useState } from "react";
import {
  Modal,
  SceneSelect,
  GallerySelect,
  ImageSelect,
} from "src/components/Shared";
// import { useToast } from "src/hooks";
import { useIntl } from "react-intl";
import { faSignOutAlt } from "@fortawesome/free-solid-svg-icons";
import { Col, Form, Row } from "react-bootstrap";
import { FormUtils } from "src/utils";

interface IFile {
  id: string;
  path: string;
}

interface IReassignFilesDialogProps {
  type: "scenes" | "images" | "galleries";
  selected: IFile[];
  reassign: (ids?: string[]) => void;
  onClose: () => void;
}

export const ReassignFilesDialog: React.FC<IReassignFilesDialogProps> = (
  props: IReassignFilesDialogProps
) => {
  const [scenes, setScenes] = useState<{ id: string; title: string }[]>([]);
  const [createNew, setCreateNew] = useState<boolean>(false);

  const intl = useIntl();
  const singularEntity = intl.formatMessage({ id: "file" });
  const pluralEntity = intl.formatMessage({ id: "files" });

  const header = intl.formatMessage(
    { id: "dialogs.reassign_entity_title" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  // const toastMessage = intl.formatMessage(
  //   { id: "toast.reassign_past_tense" },
  //   { count: props.selected.length, singularEntity, pluralEntity }
  // );

  // const Toast = useToast();

  // Network state
  const [isDeleting] /* , setIsDeleting] */ = useState(false);

  async function onAccept() {
    // setIsDeleting(true);
    // try {
    //   await mutateDeleteFiles(props.selected.map((f) => f.id));
    //   Toast.success({ content: toastMessage });
    //   props.onClose(true);
    // } catch (e) {
    //   Toast.error(e);
    //   props.onClose(false);
    // }
    // setIsDeleting(false);
  }

  function renderSelect() {
    switch (props.type) {
      case "scenes":
        return (
          <SceneSelect
            selected={createNew ? [] : scenes}
            onSelect={(items) => setScenes(items)}
            disabled={createNew}
          />
        );
      case "images":
        return (
          <ImageSelect
            selected={createNew ? [] : scenes}
            onSelect={(items) => setScenes(items)}
            disabled={createNew}
          />
        );
      case "galleries":
        return (
          <GallerySelect
            selected={createNew ? [] : scenes}
            onSelect={(items) => setScenes(items)}
            disabled={createNew}
          />
        );
    }
  }

  function getCreateEntityID() {
    switch (props.type) {
      case "scenes":
        return "scene";
      case "images":
        return "image";
      case "galleries":
        return "gallery";
    }
  }

  return (
    <Modal
      show
      icon={faSignOutAlt}
      header={header}
      accept={{
        onClick: onAccept,
        text: intl.formatMessage({ id: "actions.reassign" }),
      }}
      cancel={{
        onClick: () => props.onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      isRunning={isDeleting}
    >
      <Form>
        <Form.Check
          id="create-new"
          checked={createNew}
          label={intl.formatMessage(
            {
              id: "dialogs.create_new_entity",
            },
            {
              entity: intl.formatMessage({ id: getCreateEntityID() }),
            }
          )}
          onChange={() => setCreateNew(!createNew)}
        />
        <Form.Group controlId="dest" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({
              id: "dialogs.reassign_files.destination",
            }),
            labelProps: {
              column: true,
              sm: 3,
              xl: 12,
            },
          })}
          <Col sm={9} xl={12}>
            {renderSelect()}
          </Col>
        </Form.Group>
      </Form>
    </Modal>
  );
};

export default ReassignFilesDialog;
