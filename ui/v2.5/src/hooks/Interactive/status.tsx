import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";
import {
  ConnectionState,
  connectionStateLabel,
  InteractiveContext,
} from "./context";

export const SceneInteractiveStatus: React.FC = ({}) => {
  const { state } = React.useContext(InteractiveContext);

  function getStateClass() {
    switch (state) {
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
      <FontAwesomeIcon pulse icon="circle" size="xs" />
      <span className="status-text">{connectionStateLabel(state)}</span>
    </div>
  );
};

export default SceneInteractiveStatus;
