import React, { useCallback, useContext, useEffect, useState } from "react";
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
  initialised: boolean;
  currentScript?: string;
  error?: string;
  initialise: () => Promise<void>;
  uploadScript: (funscriptPath: string) => Promise<void>;
  sync: () => Promise<void>;
}

export const InteractiveContext = React.createContext<IState>({
  interactive: new InteractiveAPI("", 0),
  state: ConnectionState.Missing,
  serverOffset: 0,
  initialised: false,
  initialise: () => {
    return Promise.resolve();
  },
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
  const [useStashHostedFunscript, setUseStashHostedFunscript] =
    useState<boolean>(false);
  const [interactive] = useState<InteractiveAPI>(new InteractiveAPI("", 0));

  const [initialised, setInitialised] = useState(false);
  const [error, setError] = useState<string | undefined>();

  const initialise = useCallback(async () => {
    setError(undefined);

    if (!config?.serverOffset) {
      setState(ConnectionState.Syncing);
      const offset = await interactive.sync();
      setConfig({ serverOffset: offset });
      setState(ConnectionState.Ready);
      setInitialised(true);
    } else {
      interactive.setServerTimeOffset(config.serverOffset);
      setState(ConnectionState.Connecting);
      try {
        await interactive.connect();
        setState(ConnectionState.Ready);
        setInitialised(true);
      } catch (e) {
        if (e instanceof Error) {
          setError(e.message ?? e.toString());
          setState(ConnectionState.Error);
        }
      }
    }
  }, [config, interactive, setConfig]);

  useEffect(() => {
    if (!stashConfig) {
      return;
    }

    setHandyKey(stashConfig.interface.handyKey ?? undefined);
    setScriptOffset(stashConfig.interface.funscriptOffset ?? 0);
    setUseStashHostedFunscript(
      stashConfig.interface.useStashHostedFunscript ?? false
    );
  }, [stashConfig]);

  useEffect(() => {
    if (!config) {
      return;
    }

    const oldKey = interactive.handyKey;

    interactive.handyKey = handyKey ?? "";
    interactive.scriptOffset = scriptOffset;
    interactive.useStashHostedFunscript = useStashHostedFunscript;

    if (oldKey !== interactive.handyKey && interactive.handyKey) {
      initialise();
    }
  }, [
    handyKey,
    scriptOffset,
    useStashHostedFunscript,
    config,
    interactive,
    initialise,
  ]);

  const sync = useCallback(async () => {
    if (
      !interactive.handyKey ||
      state === ConnectionState.Syncing ||
      !initialised
    ) {
      return;
    }

    setState(ConnectionState.Syncing);
    const offset = await interactive.sync();
    setConfig({ serverOffset: offset });
    setState(ConnectionState.Ready);
  }, [interactive, state, setConfig, initialised]);

  const uploadScript = useCallback(
    async (funscriptPath: string) => {
      interactive.pause();
      if (
        !interactive.handyKey ||
        !funscriptPath ||
        funscriptPath === currentScript
      ) {
        return Promise.resolve();
      }

      setState(ConnectionState.Uploading);
      try {
        await interactive.uploadScript(
          funscriptPath,
          stashConfig?.general?.apiKey
        );
        setCurrentScript(funscriptPath);
        setState(ConnectionState.Ready);
      } catch (e) {
        setState(ConnectionState.Error);
      }
    },
    [interactive, currentScript, stashConfig]
  );

  return (
    <InteractiveContext.Provider
      value={{
        interactive,
        state,
        error,
        currentScript,
        serverOffset: config?.serverOffset ?? 0,
        initialised,
        initialise,
        uploadScript,
        sync,
      }}
    >
      {children}
    </InteractiveContext.Provider>
  );
};

export const useInteractive = () => {
  return useContext(InteractiveContext);
};
export default InteractiveProvider;
