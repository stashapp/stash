import { faCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";
import { FormattedMessage } from "react-intl";
import {
  ConnectionState,
  connectionStateLabel,
  InteractiveContext,
} from "./context";

export const SceneInteractiveStatus: React.FC = ({}) => {
  const { state, error } = React.useContext(InteractiveContext);

  function getStateClass() {
    switch (state) {
      case ConnectionState.Connecting:
        return "interactive-status-connecting";
      case ConnectionState.Disconnected:
        return "interactive-status-disconnected";
      case ConnectionState.Error:
        return "interactive-status-error";
      case ConnectionState.Syncing:
        return "interactive-status-uploading";
      case ConnectionState.Uploading:
        return "interactive-status-syncing";
      case ConnectionState.Ready:
        return "interactive-status-ready";
    }

    return "";
  }

  if (state === ConnectionState.Missing) {
    return <></>;
  }

  return (
    <div className={`scene-interactive-status ${getStateClass()}`}>
      <FontAwesomeIcon pulse icon={faCircle} size="xs" />
      <span className="status-text">
        <FormattedMessage id={connectionStateLabel(state)} />
        {error && <span>: {error}</span>}
      </span>
    </div>
  );
};

export default SceneInteractiveStatus;
