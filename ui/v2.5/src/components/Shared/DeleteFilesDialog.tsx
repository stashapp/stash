import React, { useState } from "react";
import { mutateDeleteFiles } from "src/core/StashService";
import { ModalComponent } from "./Modal";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import { FormattedMessage, useIntl } from "react-intl";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";

interface IFile {
  id: string;
  path: string;
}

interface IDeleteSceneDialogProps {
  selected: IFile[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteFilesDialog: React.FC<IDeleteSceneDialogProps> = (
  props: IDeleteSceneDialogProps
) => {
  const intl = useIntl();
  const singularEntity = intl.formatMessage({ id: "file" });
  const pluralEntity = intl.formatMessage({ id: "files" });

  const header = intl.formatMessage(
    { id: "dialogs.delete_entity_title" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const toastMessage = intl.formatMessage(
    { id: "toast.delete_past_tense" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const message = intl.formatMessage(
    { id: "dialogs.delete_entity_simple_desc" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );

  const Toast = useToast();

  // Network state
  const [isDeleting, setIsDeleting] = useState(false);

  const context = React.useContext(ConfigurationContext);
  const config = context?.configuration;

  async function onDelete() {
    setIsDeleting(true);
    try {
      await mutateDeleteFiles(props.selected.map((f) => f.id));
      Toast.success(toastMessage);
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
      props.onClose(false);
    }
    setIsDeleting(false);
  }

  function renderDeleteFileAlert() {
    const deletedFiles = props.selected.map((f) => f.path);

    const deleteTrashPath = config?.general.deleteTrashPath;
    const deleteAlertId = deleteTrashPath
      ? "dialogs.delete_alert_to_trash"
      : "dialogs.delete_alert";

    return (
      <div className="delete-dialog alert alert-danger text-break">
        <p className="font-weight-bold">
          <FormattedMessage
            values={{
              count: props.selected.length,
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
      {renderDeleteFileAlert()}
    </ModalComponent>
  );
};
