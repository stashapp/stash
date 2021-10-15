import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { useScenesDestroy } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { ConfigurationContext } from "src/hooks/Config";
import { useIntl } from "react-intl";

interface IDeleteSceneDialogProps {
  selected: GQL.SlimSceneDataFragment[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteScenesDialog: React.FC<IDeleteSceneDialogProps> = (
  props: IDeleteSceneDialogProps
) => {
  const intl = useIntl();
  const singularEntity = intl.formatMessage({ id: "scene" });
  const pluralEntity = intl.formatMessage({ id: "scenes" });

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
      <Form>
        <Form.Check
          id="delete-file"
          checked={deleteFile}
          label={intl.formatMessage({ id: "actions.delete_file" })}
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
