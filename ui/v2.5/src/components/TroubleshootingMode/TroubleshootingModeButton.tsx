import React, { useState } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { faBug } from "@fortawesome/free-solid-svg-icons";
import { ModalComponent } from "src/components/Shared/Modal";
import { useTroubleshootingMode } from "./useTroubleshootingMode";

const DIALOG_ITEMS = [
  "config.ui.troubleshooting_mode.dialog_item_plugins",
  "config.ui.troubleshooting_mode.dialog_item_css",
  "config.ui.troubleshooting_mode.dialog_item_js",
  "config.ui.troubleshooting_mode.dialog_item_locales",
] as const;

export const TroubleshootingModeButton: React.FC = () => {
  const intl = useIntl();
  const [showDialog, setShowDialog] = useState(false);
  const { enable, isLoading } = useTroubleshootingMode();

  return (
    <>
      <div className="troubleshooting-mode-button">
        <Button variant="primary" size="sm" onClick={() => setShowDialog(true)}>
          <FormattedMessage id="config.ui.troubleshooting_mode.button" />
        </Button>
      </div>

      <ModalComponent
        show={showDialog}
        onHide={() => setShowDialog(false)}
        header={intl.formatMessage({
          id: "config.ui.troubleshooting_mode.dialog_title",
        })}
        icon={faBug}
        accept={{
          text: intl.formatMessage({
            id: "config.ui.troubleshooting_mode.enable",
          }),
          variant: "primary",
          onClick: enable,
        }}
        cancel={{
          onClick: () => setShowDialog(false),
          variant: "secondary",
        }}
        isRunning={isLoading}
      >
        <p>
          <FormattedMessage id="config.ui.troubleshooting_mode.dialog_description" />
        </p>
        <ul>
          {DIALOG_ITEMS.map((id) => (
            <li key={id}>
              <FormattedMessage id={id} />
            </li>
          ))}
        </ul>
        <p>
          <FormattedMessage id="config.ui.troubleshooting_mode.dialog_log_level" />
        </p>
        <p className="text-muted">
          <FormattedMessage id="config.ui.troubleshooting_mode.dialog_reload_note" />
        </p>
      </ModalComponent>
    </>
  );
};
