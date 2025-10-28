import React, { useState } from "react";
import { defineMessages, FormattedMessage, useIntl } from "react-intl";
import { FetchResult } from "@apollo/client";

import { ModalComponent } from "./Modal";
import { useToast } from "src/hooks/Toast";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";

interface IDeletionEntity {
  id: string;
  name?: string | null;
}

type DestroyMutation = (input: {
  ids: string[];
}) => [() => Promise<FetchResult>, {}];

interface IDeleteEntityDialogProps {
  selected: IDeletionEntity[];
  onClose: (confirmed: boolean) => void;
  singularEntity: string;
  pluralEntity: string;
  destroyMutation: DestroyMutation;
  onDeleted?: () => void;
}

const messages = defineMessages({
  deleteHeader: {
    id: "dialogs.delete_object_title",
  },
  deleteToast: {
    id: "toast.delete_past_tense",
  },
  deleteMessage: {
    id: "dialogs.delete_object_desc",
  },
  overflowMessage: {
    id: "dialogs.delete_object_overflow",
  },
});

export const DeleteEntityDialog: React.FC<IDeleteEntityDialogProps> = ({
  selected,
  onClose,
  singularEntity,
  pluralEntity,
  destroyMutation,
  onDeleted,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [deleteEntities] = destroyMutation({ ids: selected.map((p) => p.id) });
  const count = selected.length;

  // Network state
  const [isDeleting, setIsDeleting] = useState(false);

  async function onDelete() {
    setIsDeleting(true);
    try {
      await deleteEntities();
      if (onDeleted) {
        onDeleted();
      }
      Toast.success(
        intl.formatMessage(messages.deleteToast, {
          count,
          singularEntity,
          pluralEntity,
        })
      );
    } catch (e) {
      Toast.error(e);
    }
    setIsDeleting(false);
    onClose(true);
  }

  return (
    <ModalComponent
      show
      icon={faTrashAlt}
      header={intl.formatMessage(messages.deleteHeader, {
        count,
        singularEntity,
        pluralEntity,
      })}
      accept={{
        variant: "danger",
        onClick: onDelete,
        text: intl.formatMessage({ id: "actions.delete" }),
      }}
      cancel={{
        onClick: () => onClose(false),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      isRunning={isDeleting}
    >
      <p>
        <FormattedMessage
          values={{ count, singularEntity, pluralEntity }}
          {...messages.deleteMessage}
        />
      </p>
      <ul>
        {selected.slice(0, 10).map((s) => (
          <li key={s.name}>{s.name}</li>
        ))}
        {selected.length > 10 && (
          <FormattedMessage
            values={{
              count: selected.length - 10,
              singularEntity,
              pluralEntity,
            }}
            {...messages.overflowMessage}
          />
        )}
      </ul>
    </ModalComponent>
  );
};
