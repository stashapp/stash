import React, { useState } from "react";
import { useSceneMarkersDestroy } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { useIntl } from "react-intl";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";

interface IDeleteSceneMarkersDialogProps {
  selected: GQL.SceneMarkerDataFragment[];
  onClose: (confirmed: boolean) => void;
}

export const DeleteSceneMarkersDialog: React.FC<
  IDeleteSceneMarkersDialogProps
> = (props: IDeleteSceneMarkersDialogProps) => {
  const intl = useIntl();
  const singularEntity = intl.formatMessage({ id: "marker" });
  const pluralEntity = intl.formatMessage({ id: "markers" });

  const header = intl.formatMessage(
    { id: "dialogs.delete_object_title" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const toastMessage = intl.formatMessage(
    { id: "toast.delete_past_tense" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );
  const message = intl.formatMessage(
    { id: "dialogs.delete_object_desc" },
    { count: props.selected.length, singularEntity, pluralEntity }
  );

  const Toast = useToast();
  const [deleteSceneMarkers] = useSceneMarkersDestroy(
    getSceneMarkersDeleteInput()
  );

  // Network state
  const [isDeleting, setIsDeleting] = useState(false);

  function getSceneMarkersDeleteInput(): GQL.SceneMarkersDestroyMutationVariables {
    return {
      ids: props.selected.map((marker) => marker.id),
    };
  }

  async function onDelete() {
    setIsDeleting(true);
    try {
      await deleteSceneMarkers();
      Toast.success(toastMessage);
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
      props.onClose(false);
    }
    setIsDeleting(false);
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
    </ModalComponent>
  );
};

export default DeleteSceneMarkersDialog;
