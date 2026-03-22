import React from "react";
import { useIntl } from "react-intl";
import { faExclamationTriangle } from "@fortawesome/free-solid-svg-icons";
import { ModalComponent } from "./Modal";

interface IAutoTagConfirmDialog {
  show: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export const AutoTagConfirmDialog: React.FC<IAutoTagConfirmDialog> = ({
  show,
  onConfirm,
  onCancel,
}) => {
  const intl = useIntl();

  return (
    <ModalComponent
      show={show}
      icon={faExclamationTriangle}
      header={intl.formatMessage({ id: "actions.auto_tag" })}
      accept={{
        text: intl.formatMessage({ id: "actions.confirm" }),
        variant: "danger",
        onClick: onConfirm,
      }}
      cancel={{
        onClick: onCancel,
      }}
    >
      <p>
        {intl.formatMessage({
          id: "config.tasks.auto_tag_based_on_filenames",
        })}
      </p>
      <p>{intl.formatMessage({ id: "config.tasks.auto_tag_confirm" })}</p>
    </ModalComponent>
  );
};
