import React, { useState } from "react";
import { defineMessages, FormattedMessage, useIntl } from "react-intl";
import { FetchResult } from "@apollo/client";

import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";

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
}

const messages = defineMessages({
  deleteHeader: {
    id: "delete-header",
    defaultMessage:
      "Delete {count, plural, =1 {{singularEntity}} other {{pluralEntity}}}",
  },
  deleteToast: {
    id: "delete-toast",
    defaultMessage:
      "Deleted {count, plural, =1 {{singularEntity}} other {{pluralEntity}}}",
  },
  deleteMessage: {
    id: "delete-message",
    defaultMessage:
      "Are you sure you want to delete {count, plural, =1 {this {singularEntity}} other {these {pluralEntity}}}?",
  },
  overflowMessage: {
    id: "overflow-message",
    defaultMessage:
      "...and {count} other {count, plural, =1 {{ singularEntity}} other {{ pluralEntity }}}.",
  },
});

const DeleteEntityDialog: React.FC<IDeleteEntityDialogProps> = ({
  selected,
  onClose,
  singularEntity,
  pluralEntity,
  destroyMutation,
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
      Toast.success({
        content: intl.formatMessage(messages.deleteToast, {
          count,
          singularEntity,
          pluralEntity,
        }),
      });
    } catch (e) {
      Toast.error(e);
    }
    setIsDeleting(false);
    onClose(true);
  }

  return (
    <Modal
      show
      icon="trash-alt"
      header={intl.formatMessage(messages.deleteHeader, {
        count,
        singularEntity,
        pluralEntity,
      })}
      accept={{ variant: "danger", onClick: onDelete, text: intl.formatMessage({ id: 'actions.delete' }) }}
      cancel={{
        onClick: () => onClose(false),
        text: "Cancel",
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
          <li>{s.name}</li>
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
    </Modal>
  );
};

export default DeleteEntityDialog;
