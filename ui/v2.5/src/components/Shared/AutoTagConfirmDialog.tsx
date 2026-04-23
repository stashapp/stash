import React, { useContext, useEffect, useState } from "react";
import { Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { faExclamationTriangle } from "@fortawesome/free-solid-svg-icons";
import { useConfigureInterface } from "src/core/StashService";
import { SettingStateContext } from "src/components/Settings/context";
import { useToast } from "src/hooks/Toast";
import { ModalComponent } from "./Modal";
import { Icon } from "./Icon";

interface IAutoTagConfirmDialog {
  show: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export const AutoTagWarning = () => (
  <>
    <p>
      <FormattedMessage id="config.tasks.auto_tag_based_on_filenames" />
    </p>
    <p>
      <FormattedMessage id="config.tasks.auto_tag_confirm" />
    </p>
    <p className="lead">
      <Icon icon={faExclamationTriangle} className="text-warning" />
      <FormattedMessage id="config.tasks.auto_tag_warning" />
    </p>
  </>
);

export const AutoTagConfirmDialog: React.FC<IAutoTagConfirmDialog> = ({
  show,
  onConfirm,
  onCancel,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [dontShowAgain, setDontShowAgain] = useState(false);
  const [configureInterface] = useConfigureInterface();
  const settingsContext = useContext(SettingStateContext);

  useEffect(() => {
    if (show) {
      setDontShowAgain(false);
    }
  }, [show]);

  function handleConfirm() {
    if (dontShowAgain) {
      if (settingsContext) {
        settingsContext.saveInterface({ disableAutoTagWarning: true });
      } else {
        configureInterface({
          variables: { input: { disableAutoTagWarning: true } },
        }).catch((e) => Toast.error(e));
      }
    }
    onConfirm();
  }

  return (
    <ModalComponent
      show={show}
      icon={faExclamationTriangle}
      header={intl.formatMessage({ id: "actions.auto_tag" })}
      accept={{
        text: intl.formatMessage({ id: "actions.confirm" }),
        variant: "danger",
        onClick: handleConfirm,
      }}
      cancel={{
        onClick: onCancel,
      }}
    >
      <AutoTagWarning />
      <Form.Check
        id="auto-tag-dont-show-again"
        checked={dontShowAgain}
        onChange={(e) => setDontShowAgain(e.currentTarget.checked)}
        label={intl.formatMessage({
          id: "dialogs.dont_show_again",
        })}
      />
    </ModalComponent>
  );
};
