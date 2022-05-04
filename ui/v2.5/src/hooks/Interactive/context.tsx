import React, { useCallback, useEffect, useState } from "react";
import { ConfigurationContext } from "../Config";
import { useLocalForage } from "../LocalForage";
import { Interactive as InteractiveAPI } from "./interactive";

export enum ConnectionState {
  Missing,
  Disconnected,
  Error,
  Connecting,
  Syncing,
  Uploading,
  Ready,
}

export function connectionStateLabel(s: ConnectionState) {
  const prefix = "handy_connection_status";
  switch (s) {
    case ConnectionState.Missing:
      return `${prefix}.missing`;
    case ConnectionState.Connecting:
      return `${prefix}.connecting`;
    case ConnectionState.Disconnected:
      return `${prefix}.disconnected`;
    case ConnectionState.Error:
      return `${prefix}.error`;
    case ConnectionState.Syncing:
      return `${prefix}.syncing`;
    case ConnectionState.Uploading:
      return `${prefix}.uploading`;
    case ConnectionState.Ready:
      return `${prefix}.ready`;
  }
}

export interface IState {
  interactive: InteractiveAPI;
  state: ConnectionState;
  serverOffset: number;
  uploadScript: (funscriptPath: string) => Promise<void>;
  sync: () => Promise<void>;
}

export const InteractiveContext = React.createContext<IState>({
  interactive: new InteractiveAPI("", 0),
  state: ConnectionState.Missing,
  serverOffset: 0,
  uploadScript: () => {
    return Promise.resolve();
  },
  sync: () => {
    return Promise.resolve();
  },
});

const LOCAL_FORAGE_KEY = "interactive";

interface IInteractiveState {
  serverOffset: number;
}

export const InteractiveProvider: React.FC = ({ children }) => {
  const [{ data: config }, setConfig] = useLocalForage<IInteractiveState>(
    LOCAL_FORAGE_KEY,
    { serverOffset: 0 }
  );

  const { configuration: stashConfig } = React.useContext(ConfigurationContext);

  const [state, setState] = useState<ConnectionState>(ConnectionState.Missing);
  const [handyKey, setHandyKey] = useState<string | undefined>(undefined);
  const [currentScript, setCurrentScript] = useState<string | undefined>(
    undefined
  );
  const [scriptOffset, setScriptOffset] = useState<number>(0);
  const [interactive] = useState<InteractiveAPI>(new InteractiveAPI("", 0));

  useEffect(() => {
    if (!stashConfig) {
      return;
    }

    setHandyKey(stashConfig.interface.handyKey ?? undefined);
    setScriptOffset(stashConfig.interface.funscriptOffset ?? 0);
  }, [stashConfig]);

  useEffect(() => {
    if (!config) {
      return;
    }

    const oldKey = interactive.handyKey;

    interactive.handyKey = handyKey ?? "";
    interactive.scriptOffset = scriptOffset;

    if (oldKey !== interactive.handyKey && interactive.handyKey) {
      if (!config?.serverOffset) {
        setState(ConnectionState.Syncing);
        interactive.sync().then((offset) => {
          setConfig({ serverOffset: offset });
          setState(ConnectionState.Ready);
        });
      } else {
        interactive.setServerTimeOffset(config.serverOffset);
        setState(ConnectionState.Connecting);
        interactive.connect().then(() => {
          setState(ConnectionState.Ready);
        });
      }
    }
  }, [handyKey, scriptOffset, config, interactive, setConfig]);

  const sync = useCallback(async () => {
    if (!interactive.handyKey || state === ConnectionState.Syncing) {
      return;
    }

    setState(ConnectionState.Syncing);
    const offset = await interactive.sync();
    setConfig({ serverOffset: offset });
    setState(ConnectionState.Ready);
  }, [interactive, state, setConfig]);

  const uploadScript = useCallback(
    async (funscriptPath: string) => {
      if (
        !interactive.handyKey ||
        !funscriptPath ||
        funscriptPath === currentScript
      ) {
        return Promise.resolve();
      }

      setState(ConnectionState.Uploading);
      try {
        await interactive.uploadScript(funscriptPath);
      } catch (e) {
        setState(ConnectionState.Error);
        return;
      }
      setCurrentScript(funscriptPath);
      setState(ConnectionState.Ready);
    },
    [interactive, currentScript]
  );

  return (
    <InteractiveContext.Provider
      value={{
        interactive,
        state,
        serverOffset: config?.serverOffset ?? 0,
        uploadScript,
        sync,
      }}
    >
      {children}
    </InteractiveContext.Provider>
  );
};

export default InteractiveProvider;
