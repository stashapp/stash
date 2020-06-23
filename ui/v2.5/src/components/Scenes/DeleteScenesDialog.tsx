import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { useScenesDestroy } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormattedMessage } from "react-intl";

interface IDeleteSceneDialogProps {
  selected: GQL.SlimSceneDataFragment[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteScenesDialog: React.FC<IDeleteSceneDialogProps> = (
  props: IDeleteSceneDialogProps
) => {
  const plural = props.selected.length > 1;

  const singleMessageId = "deleteSceneText";
  const pluralMessageId = "deleteScenesText";

  const singleMessage =
    "Are you sure you want to delete this scene? Unless the file is also deleted, this scene will be re-added when scan is performed.";
  const pluralMessage =
    "Are you sure you want to delete these scenes? Unless the files are also deleted, these scenes will be re-added when scan is performed.";

  const header = plural ? "Delete Scenes" : "Delete Scene";
  const toastMessage = plural ? "Deleted scenes" : "Deleted scene";
  const messageId = plural ? pluralMessageId : singleMessageId;
  const message = plural ? pluralMessage : singleMessage;

  const [deleteFile, setDeleteFile] = useState<boolean>(false);
  const [deleteGenerated, setDeleteGenerated] = useState<boolean>(true);

  const Toast = useToast();
  const [deleteScene] = useScenesDestroy(getScenesDeleteInput());

  // Network state
  const [isDeleting, setIsDeleting] = useState(false);

  function getScenesDeleteInput(): GQL.ScenesDestroyInput {
    return {
      ids: props.selected.map((scene) => scene.id),
      delete_file: deleteFile,
      delete_generated: deleteGenerated,
    };
  }

  async function onDelete() {
    setIsDeleting(true);
    try {
      await deleteScene();
      Toast.success({ content: toastMessage });
    } catch (e) {
      Toast.error(e);
    }
    setIsDeleting(false);
    props.onClose(true);
  }

  return (
    <Modal
      show
      icon="trash-alt"
      header={header}
      accept={{ variant: "danger", onClick: onDelete, text: "Delete" }}
      cancel={{
        onClick: () => props.onClose(false),
        text: "Cancel",
        variant: "secondary",
      }}
      isRunning={isDeleting}
    >
      <p>
        <FormattedMessage id={messageId} defaultMessage={message} />
      </p>
      <Form>
        <Form.Check
          checked={deleteFile}
          label="Delete file"
          onChange={() => setDeleteFile(!deleteFile)}
        />
        <Form.Check
          checked={deleteGenerated}
          label="Delete generated supporting files"
          onChange={() => setDeleteGenerated(!deleteGenerated)}
        />
      </Form>
    </Modal>
  );
};
