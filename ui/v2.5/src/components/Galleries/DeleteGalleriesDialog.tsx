import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { useGalleryDestroy } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormattedMessage } from "react-intl";

interface IDeleteGalleryDialogProps {
  selected: Partial<GQL.GalleryDataFragment>[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteGalleriesDialog: React.FC<IDeleteGalleryDialogProps> = (
  props: IDeleteGalleryDialogProps
) => {
  const plural = props.selected.length > 1;

  const singleMessageId = "deleteGalleryText";
  const pluralMessageId = "deleteGallerysText";

  const singleMessage =
    "Are you sure you want to delete this gallery? Galleries for zip files will be re-added during the next scan unless the zip file is also deleted.";
  const pluralMessage =
    "Are you sure you want to delete these galleries? Galleries for zip files will be re-added during the next scan unless the zip files are also deleted.";

  const header = plural ? "Delete Galleries" : "Delete Gallery";
  const toastMessage = plural ? "Deleted galleries" : "Deleted gallery";
  const messageId = plural ? pluralMessageId : singleMessageId;
  const message = plural ? pluralMessage : singleMessage;

  const [deleteFile, setDeleteFile] = useState<boolean>(false);
  const [deleteGenerated, setDeleteGenerated] = useState<boolean>(true);

  const Toast = useToast();
  const [deleteGallery] = useGalleryDestroy(getGalleriesDeleteInput());

  // Network state
  const [isDeleting, setIsDeleting] = useState(false);

  function getGalleriesDeleteInput(): GQL.GalleryDestroyInput {
    return {
      ids: props.selected.map((gallery) => gallery.id!),
      delete_file: deleteFile,
      delete_generated: deleteGenerated,
    };
  }

  async function onDelete() {
    setIsDeleting(true);
    try {
      await deleteGallery();
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
          id="delete-file"
          checked={deleteFile}
          label="Delete zip file and any images not attached to any other gallery."
          onChange={() => setDeleteFile(!deleteFile)}
        />
        <Form.Check
          id="delete-generated"
          checked={deleteGenerated}
          label="Delete generated supporting files"
          onChange={() => setDeleteGenerated(!deleteGenerated)}
        />
      </Form>
    </Modal>
  );
};
