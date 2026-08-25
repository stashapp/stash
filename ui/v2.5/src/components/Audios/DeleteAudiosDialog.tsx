import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { useAudiosDestroy } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { useConfigurationContext } from "src/hooks/Config";
import { FormattedMessage, useIntl } from "react-intl";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";

interface IDeleteAudioDialogProps {
  selected: GQL.SlimAudioDataFragment[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteAudiosDialog: React.FC<IDeleteAudioDialogProps> = (
  props: IDeleteAudioDialogProps
) => {
  const intl = useIntl();
  const singularEntity = intl.formatMessage({ id: "audio" });
  const pluralEntity = intl.formatMessage({ id: "audios" });

  const header = intl.formatMessage(
    { id: "dialogs.delete_entity_title" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const toastMessage = intl.formatMessage(
    { id: "toast.delete_past_tense" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const message = intl.formatMessage(
    { id: "dialogs.delete_entity_desc" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );

  const { configuration: config } = useConfigurationContext();

  const [deleteFile, setDeleteFile] = useState<boolean>(
    config?.defaults.deleteFile ?? false
  );

  const Toast = useToast();
  const [deleteAudio] = useAudiosDestroy(getAudiosDeleteInput());

  // Network state
  const [isDeleting, setIsDeleting] = useState(false);

  function getAudiosDeleteInput(): GQL.AudiosDestroyInput {
    return {
      ids: props.selected.map((audio) => audio.id),
      delete_file: deleteFile,
    };
  }

  async function onDelete() {
    setIsDeleting(true);
    try {
      await deleteAudio();
      Toast.success(toastMessage);
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
      props.onClose(false);
    }
    setIsDeleting(false);
  }

  function maybeRenderDeleteFileAlert() {
    if (!deleteFile) {
      return;
    }

    const deletedFiles: string[] = [];

    props.selected.forEach((a) => {
      deletedFiles.push(...a.files.map((f) => f.path));
    });

    const deleteTrashPath = config?.general.deleteTrashPath;
    const deleteAlertId = deleteTrashPath
      ? "dialogs.delete_alert_to_trash"
      : "dialogs.delete_alert";

    return (
      <div className="delete-dialog alert alert-danger text-break">
        <p className="font-weight-bold">
          <FormattedMessage
            values={{
              count: deletedFiles.length,
              singularEntity: intl.formatMessage({ id: "file" }),
              pluralEntity: intl.formatMessage({ id: "files" }),
            }}
            id={deleteAlertId}
          />
        </p>
        <ul>
          {deletedFiles.slice(0, 5).map((s) => (
            <li key={s}>{s}</li>
          ))}
          {deletedFiles.length > 5 && (
            <FormattedMessage
              values={{
                count: deletedFiles.length - 5,
                singularEntity: intl.formatMessage({ id: "file" }),
                pluralEntity: intl.formatMessage({ id: "files" }),
              }}
              id="dialogs.delete_object_overflow"
            />
          )}
        </ul>
      </div>
    );
  }

  return (
    <ModalComponent
      show
      icon={faTrashAlt}
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
          label={intl.formatMessage({
            id: "actions.delete_file",
          })}
          onChange={() => setDeleteFile(!deleteFile)}
        />
      </Form>
    </ModalComponent>
  );
};

export default DeleteAudiosDialog;
