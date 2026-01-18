import React from "react";
import { Alert, Button } from "react-bootstrap";
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
      <Alert variant="danger" className="mb-0">
        <div className="d-flex justify-content-between align-items-center">
          <span>
            <Icon icon={faBug} className="mr-2" />
            <FormattedMessage id="config.ui.troubleshooting_mode.overlay_message" />
          </span>
          <Button
            variant="danger"
            size="sm"
            onClick={disable}
            disabled={isLoading}
          >
            <FormattedMessage id="config.ui.troubleshooting_mode.exit" />
          </Button>
        </div>
      </Alert>
    </div>
  );
};
