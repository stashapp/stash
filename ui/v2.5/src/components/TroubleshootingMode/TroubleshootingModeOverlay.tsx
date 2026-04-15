import React from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { faBug } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { useTroubleshootingMode } from "./useTroubleshootingMode";

export const TroubleshootingModeOverlay: React.FC = () => {
  const { isActive, isLoading, disable } = useTroubleshootingMode();

  if (!isActive) {
    return null;
  }

  return (
    <div className="troubleshooting-mode-overlay">
      <div className="troubleshooting-mode-alert">
        <span>
          <Icon icon={faBug} className="mr-2" />
          <FormattedMessage id="config.ui.troubleshooting_mode.overlay_message" />
        </span>
        <Button variant="link" onClick={disable} disabled={isLoading}>
          <FormattedMessage id="config.ui.troubleshooting_mode.exit" />
        </Button>
      </div>
    </div>
  );
};
