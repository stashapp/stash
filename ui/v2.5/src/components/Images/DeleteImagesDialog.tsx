import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { useImagesDestroy } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormattedMessage } from "react-intl";

interface IDeleteImageDialogProps {
  selected: GQL.SlimImageDataFragment[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteImagesDialog: React.FC<IDeleteImageDialogProps> = (
  props: IDeleteImageDialogProps
) => {
  const plural = props.selected.length > 1;

  const singleMessageId = "deleteImageText";
  const pluralMessageId = "deleteImagesText";

  const singleMessage =
    "Are you sure you want to delete this image? Unless the file is also deleted, this image will be re-added when scan is performed.";
  const pluralMessage =
    "Are you sure you want to delete these images? Unless the files are also deleted, these images will be re-added when scan is performed.";

  const header = plural ? "Delete Images" : "Delete Image";
  const toastMessage = plural ? "Deleted images" : "Deleted image";
  const messageId = plural ? pluralMessageId : singleMessageId;
  const message = plural ? pluralMessage : singleMessage;

  const [deleteFile, setDeleteFile] = useState<boolean>(false);
  const [deleteGenerated, setDeleteGenerated] = useState<boolean>(true);

  const Toast = useToast();
  const [deleteImage] = useImagesDestroy(getImagesDeleteInput());

  // Network state
  const [isDeleting, setIsDeleting] = useState(false);

  function getImagesDeleteInput(): GQL.ImagesDestroyInput {
    return {
      ids: props.selected.map((image) => image.id),
      delete_file: deleteFile,
      delete_generated: deleteGenerated,
    };
  }

  async function onDelete() {
    setIsDeleting(true);
    try {
      await deleteImage();
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
