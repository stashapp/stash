import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { useGalleryDestroy } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { ConfigurationContext } from "src/hooks/Config";
import { FormattedMessage, useIntl } from "react-intl";

interface IDeleteGalleryDialogProps {
  selected: GQL.SlimGalleryDataFragment[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteGalleriesDialog: React.FC<IDeleteGalleryDialogProps> = (
  props: IDeleteGalleryDialogProps
) => {
  const intl = useIntl();
  const singularEntity = intl.formatMessage({ id: "gallery" });
  const pluralEntity = intl.formatMessage({ id: "galleries" });

  const header = intl.formatMessage(
    { id: "dialogs.delete_entity_title" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const toastMessage = intl.formatMessage(
    { id: "toast.delete_entity" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const message = intl.formatMessage(
    { id: "dialogs.delete_entity_desc" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );

  const { configuration: config } = React.useContext(ConfigurationContext);

  const [deleteFile, setDeleteFile] = useState<boolean>(
    config?.interface.deleteFileDefault ?? false
  );
  const [deleteGenerated, setDeleteGenerated] = useState<boolean>(
    config?.interface.deleteGeneratedDefault ?? true
  );

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

  function maybeRenderDeleteFileAlert() {
    if (!deleteFile) {
      return;
    }

    return (
      <div className="delete-dialog alert alert-danger text-break">
        <p className="font-weight-bold">
          <FormattedMessage
            values={{
              count: props.selected.length,
              singularEntity,
              pluralEntity,
            }}
            id="dialogs.delete_alert"
          />
        </p>
        <ul>
          {props.selected.slice(0, 5).map((s) => (
            <li>{s.path}</li>
          ))}
          {props.selected.length > 5 && (
            <FormattedMessage
              values={{
                count: props.selected.length - 5,
                singularEntity,
                pluralEntity,
              }}
              id="dialogs.delete_object_overflow"
            />
          )}
        </ul>
      </div>
    );
  }

  return (
    <Modal
      show
      icon="trash-alt"
      header={header}
      accept={{
        variant: "danger",
        onClick: onDelete,
        text: intl.formatMessage({ id: "actions.delete" }),
      }}
      cancel={{
        onClick: () => props.onClose(false),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      isRunning={isDeleting}
    >
      <p>{message}</p>
      {maybeRenderDeleteFileAlert()}
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
          label={intl.formatMessage({
            id: "actions.delete_generated_supporting_files",
          })}
          onChange={() => setDeleteGenerated(!deleteGenerated)}
        />
      </Form>
    </Modal>
  );
};
