import { useEffect, useState } from "react";
import { getWSClient, useWSState } from "./core/StashService";
import { useToast } from "./hooks/Toast";
import { useIntl } from "react-intl";

export const ConnectionMonitor: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();

  const { state } = useWSState(getWSClient());
  const [cachedState, setCacheState] = useState<typeof state>(state);

  useEffect(() => {
    if (cachedState === "connecting" && state === "error") {
      Toast.error(
        intl.formatMessage({
          id: "connection_monitor.websocket_connection_failed",
        })
      );
    }

    if (state === "connected" && cachedState === "error") {
      Toast.success(
        intl.formatMessage({
          id: "connection_monitor.websocket_connection_reestablished",
        })
      );
    }

    setCacheState(state);
  }, [state, cachedState, Toast, intl]);

  return null;
};
