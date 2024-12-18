import { ApolloError } from "@apollo/client/errors";
import {
  faCheckCircle,
  faTimesCircle,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState, useEffect, useCallback, useRef } from "react";
import { Spinner } from "react-bootstrap";
import { IUIConfig } from "src/core/config";
import * as GQL from "src/core/generated-graphql";
import {
  useConfiguration,
  useConfigureDefaults,
  useConfigureDLNA,
  useConfigureGeneral,
  useConfigureInterface,
  useConfigurePlugin,
  useConfigureScraping,
  useConfigureUI,
} from "src/core/StashService";
import { useDebounce } from "src/hooks/debounce";
import { useToast } from "src/hooks/Toast";
import { withoutTypename } from "src/utils/data";
import { Icon } from "../Shared/Icon";

type PluginConfigs = Record<string, Record<string, unknown>>;

export interface ISettingsContextState {
  loading: boolean;
  error: ApolloError | undefined;
  general: GQL.ConfigGeneralInput;
  interface: GQL.ConfigInterfaceInput;
  defaults: GQL.ConfigDefaultSettingsInput;
  scraping: GQL.ConfigScrapingInput;
  dlna: GQL.ConfigDlnaInput;
  ui: IUIConfig;
  plugins: PluginConfigs;

  advancedMode: boolean;

  // apikey isn't directly settable, so expose it here
  apiKey: string;

  saveGeneral: (input: Partial<GQL.ConfigGeneralInput>) => void;
  saveInterface: (input: Partial<GQL.ConfigInterfaceInput>) => void;
  saveDefaults: (input: Partial<GQL.ConfigDefaultSettingsInput>) => void;
  saveScraping: (input: Partial<GQL.ConfigScrapingInput>) => void;
  saveDLNA: (input: Partial<GQL.ConfigDlnaInput>) => void;
  saveUI: (input: Partial<IUIConfig>) => void;
  savePluginSettings: (pluginID: string, input: {}) => void;
  setAdvancedMode: (value: boolean) => void;

  refetch: () => void;
}

function noop() {}

const emptyState: ISettingsContextState = {
  loading: false,
  error: undefined,
  general: {},
  interface: {},
  defaults: {},
  scraping: {},
  dlna: {},
  ui: {},
  plugins: {},

  advancedMode: false,

  apiKey: "",

  saveGeneral: noop,
  saveInterface: noop,
  saveDefaults: noop,
  saveScraping: noop,
  saveDLNA: noop,
  saveUI: noop,
  savePluginSettings: noop,
  setAdvancedMode: noop,

  refetch: noop,
};

export const SettingStateContext =
  React.createContext<ISettingsContextState | null>(null);

export const useSettings = () => {
  const context = React.useContext(SettingStateContext);

  if (context === null) {
    throw new Error("useSettings must be used within a SettingsContext");
  }

  return context;
};

export function useSettingsOptional(): ISettingsContextState {
  const context = React.useContext(SettingStateContext);

  if (context === null) {
    return emptyState;
  }

  return context;
}

export const SettingsContext: React.FC = ({ children }) => {
  const Toast = useToast();

  const { data, error, loading, refetch } = useConfiguration();
  const initialRef = useRef(false);

  const [general, setGeneral] = useState<GQL.ConfigGeneralInput>({});
  const [pendingGeneral, setPendingGeneral] =
    useState<GQL.ConfigGeneralInput>();
  const [updateGeneralConfig] = useConfigureGeneral();

  const [iface, setIface] = useState<GQL.ConfigInterfaceInput>({});
  const [pendingInterface, setPendingInterface] =
    useState<GQL.ConfigInterfaceInput>();
  const [updateInterfaceConfig] = useConfigureInterface();

  const [defaults, setDefaults] = useState<GQL.ConfigDefaultSettingsInput>({});
  const [pendingDefaults, setPendingDefaults] =
    useState<GQL.ConfigDefaultSettingsInput>();
  const [updateDefaultsConfig] = useConfigureDefaults();

  const [scraping, setScraping] = useState<GQL.ConfigScrapingInput>({});
  const [pendingScraping, setPendingScraping] =
    useState<GQL.ConfigScrapingInput>();
  const [updateScrapingConfig] = useConfigureScraping();

  const [dlna, setDLNA] = useState<GQL.ConfigDlnaInput>({});
  const [pendingDLNA, setPendingDLNA] = useState<GQL.ConfigDlnaInput>();
  const [updateDLNAConfig] = useConfigureDLNA();

  const [ui, setUI] = useState<IUIConfig>({});
  const [pendingUI, setPendingUI] = useState<{}>();
  const [updateUIConfig] = useConfigureUI();

  const [plugins, setPlugins] = useState<PluginConfigs>({});
  const [pendingPlugins, setPendingPlugins] = useState<PluginConfigs>();
  const [updatePluginConfig] = useConfigurePlugin();

  const [updateSuccess, setUpdateSuccess] = useState<boolean>();

  const [apiKey, setApiKey] = useState("");

  useEffect(() => {
    if (!data?.configuration || error) return;

    // always set api key
    setApiKey(data.configuration.general.apiKey);

    // only initialise once - assume we have control over these settings and
    // they aren't modified elsewhere
    if (initialRef.current) return;
    initialRef.current = true;

    setGeneral({ ...withoutTypename(data.configuration.general) });
    setIface({ ...withoutTypename(data.configuration.interface) });
    setDefaults({ ...withoutTypename(data.configuration.defaults) });
    setScraping({ ...withoutTypename(data.configuration.scraping) });
    setDLNA({ ...withoutTypename(data.configuration.dlna) });
    setUI(data.configuration.ui);
    setPlugins(data.configuration.plugins);
  }, [data, error]);

  const resetSuccess = useDebounce(() => setUpdateSuccess(undefined), 4000);

  const onSuccess = useCallback(() => {
    setUpdateSuccess(true);
    resetSuccess();
  }, [resetSuccess]);

  const onError = useCallback(
    (err) => {
      Toast.error(err);
      setUpdateSuccess(false);
    },
    [Toast]
  );

  // saves the configuration if no further changes are made after a half second
  const saveGeneralConfig = useDebounce(
    async (input: GQL.ConfigGeneralInput) => {
      try {
        setUpdateSuccess(undefined);
        await updateGeneralConfig({
          variables: {
            input,
          },
        });

        setPendingGeneral(undefined);
        onSuccess();
      } catch (e) {
        onError(e);
      }
    },
    500
  );

  useEffect(() => {
    if (!pendingGeneral) {
      return;
    }

    saveGeneralConfig(pendingGeneral);
  }, [pendingGeneral, saveGeneralConfig]);

  function saveGeneral(input: Partial<GQL.ConfigGeneralInput>) {
    if (!general) {
      return;
    }

    setGeneral({
      ...general,
      ...input,
    });

    setPendingGeneral((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  // saves the configuration if no further changes are made after a half second
  const saveInterfaceConfig = useDebounce(
    async (input: GQL.ConfigInterfaceInput) => {
      try {
        setUpdateSuccess(undefined);
        await updateInterfaceConfig({
          variables: {
            input,
          },
        });

        setPendingInterface(undefined);
        onSuccess();
      } catch (e) {
        onError(e);
      }
    },
    500
  );

  useEffect(() => {
    if (!pendingInterface) {
      return;
    }

    saveInterfaceConfig(pendingInterface);
  }, [pendingInterface, saveInterfaceConfig]);

  function saveInterface(input: Partial<GQL.ConfigInterfaceInput>) {
    if (!iface) {
      return;
    }

    setIface({
      ...iface,
      ...input,
    });

    setPendingInterface((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  // saves the configuration if no further changes are made after a half second
  const saveDefaultsConfig = useDebounce(
    async (input: GQL.ConfigDefaultSettingsInput) => {
      try {
        setUpdateSuccess(undefined);
        await updateDefaultsConfig({
          variables: {
            input,
          },
        });

        setPendingDefaults(undefined);
        onSuccess();
      } catch (e) {
        onError(e);
      }
    },
    500
  );

  useEffect(() => {
    if (!pendingDefaults) {
      return;
    }

    saveDefaultsConfig(pendingDefaults);
  }, [pendingDefaults, saveDefaultsConfig]);

  function saveDefaults(input: Partial<GQL.ConfigDefaultSettingsInput>) {
    if (!defaults) {
      return;
    }

    setDefaults({
      ...defaults,
      ...input,
    });

    setPendingDefaults((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  // saves the configuration if no further changes are made after a half second
  const saveScrapingConfig = useDebounce(
    async (input: GQL.ConfigScrapingInput) => {
      try {
        setUpdateSuccess(undefined);
        await updateScrapingConfig({
          variables: {
            input,
          },
        });

        setPendingScraping(undefined);
        onSuccess();
      } catch (e) {
        onError(e);
      }
    },
    500
  );

  useEffect(() => {
    if (!pendingScraping) {
      return;
    }

    saveScrapingConfig(pendingScraping);
  }, [pendingScraping, saveScrapingConfig]);

  function saveScraping(input: Partial<GQL.ConfigScrapingInput>) {
    if (!scraping) {
      return;
    }

    setScraping({
      ...scraping,
      ...input,
    });

    setPendingScraping((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  // saves the configuration if no further changes are made after a half second
  const saveDLNAConfig = useDebounce(async (input: GQL.ConfigDlnaInput) => {
    try {
      setUpdateSuccess(undefined);
      await updateDLNAConfig({
        variables: {
          input,
        },
      });

      setPendingDLNA(undefined);
      onSuccess();
    } catch (e) {
      onError(e);
    }
  }, 500);

  useEffect(() => {
    if (!pendingDLNA) {
      return;
    }

    saveDLNAConfig(pendingDLNA);
  }, [pendingDLNA, saveDLNAConfig]);

  function saveDLNA(input: Partial<GQL.ConfigDlnaInput>) {
    if (!dlna) {
      return;
    }

    setDLNA({
      ...dlna,
      ...input,
    });

    setPendingDLNA((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  type UIConfigInput = GQL.Scalars["Map"]["input"];

  // saves the configuration if no further changes are made after a half second
  const saveUIConfig = useDebounce(async (input: Partial<IUIConfig>) => {
    try {
      setUpdateSuccess(undefined);
      await updateUIConfig({
        variables: {
          partial: input as UIConfigInput,
        },
      });

      setPendingUI(undefined);
      onSuccess();
    } catch (e) {
      onError(e);
    }
  }, 500);

  useEffect(() => {
    if (!pendingUI) {
      return;
    }

    saveUIConfig(pendingUI);
  }, [pendingUI, saveUIConfig]);

  function saveUI(input: IUIConfig) {
    if (!ui) {
      return;
    }

    setUI({
      ...ui,
      ...input,
    });

    setPendingUI((current) => {
      return {
        ...current,
        ...input,
      };
    });
  }

  function setAdvancedMode(value: boolean) {
    saveUI({
      advancedMode: value,
    });
  }

  // saves the configuration if no further changes are made after a half second
  const savePluginConfig = useDebounce(async (input: PluginConfigs) => {
    try {
      setUpdateSuccess(undefined);

      for (const pluginID in input) {
        await updatePluginConfig({
          variables: {
            plugin_id: pluginID,
            input: input[pluginID],
          },
        });
      }

      setPendingPlugins(undefined);
      onSuccess();
    } catch (e) {
      onError(e);
    }
  }, 500);

  useEffect(() => {
    if (!pendingPlugins) {
      return;
    }

    savePluginConfig(pendingPlugins);
  }, [pendingPlugins, savePluginConfig]);

  function savePluginSettings(
    pluginID: string,
    input: Record<string, unknown>
  ) {
    if (!plugins) {
      return;
    }

    setPlugins({
      ...plugins,
      [pluginID]: input,
    });

    setPendingPlugins((current) => {
      if (!current) {
        // use full UI object to ensure nothing is wiped
        return {
          ...plugins,
          [pluginID]: input,
        };
      }
      return {
        ...current,
        [pluginID]: input,
      };
    });
  }

  function maybeRenderLoadingIndicator() {
    if (updateSuccess === false) {
      return (
        <div className="loading-indicator failed">
          <Icon icon={faTimesCircle} className="fa-fw" />
        </div>
      );
    }

    if (
      pendingGeneral ||
      pendingInterface ||
      pendingDefaults ||
      pendingScraping ||
      pendingDLNA ||
      pendingUI ||
      pendingPlugins
    ) {
      return (
        <div className="loading-indicator">
          <Spinner animation="border" role="status">
            <span className="sr-only">Loading...</span>
          </Spinner>
        </div>
      );
    }

    if (updateSuccess) {
      return (
        <div className="loading-indicator success">
          <Icon icon={faCheckCircle} className="fa-fw" />
        </div>
      );
    }
  }

  return (
    <SettingStateContext.Provider
      value={{
        loading,
        error,
        apiKey,
        general,
        interface: iface,
        defaults,
        scraping,
        dlna,
        ui,
        plugins,
        advancedMode: ui.advancedMode ?? false,
        saveGeneral,
        saveInterface,
        saveDefaults,
        saveScraping,
        saveDLNA,
        saveUI,
        refetch,
        savePluginSettings,
        setAdvancedMode,
      }}
    >
      {maybeRenderLoadingIndicator()}
      {children}
    </SettingStateContext.Provider>
  );
};
