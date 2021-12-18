import { ApolloError } from "@apollo/client/errors";
import { debounce } from "lodash";
import React, {
  useState,
  useEffect,
  useMemo,
  useCallback,
  useRef,
} from "react";
import { Spinner } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import {
  useConfiguration,
  useConfigureDefaults,
  useConfigureDLNA,
  useConfigureGeneral,
  useConfigureInterface,
  useConfigureScraping,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { withoutTypename } from "src/utils";
import { Icon } from "../Shared";

export interface ISettingsContextState {
  loading: boolean;
  error: ApolloError | undefined;
  general: GQL.ConfigGeneralInput;
  interface: GQL.ConfigInterfaceInput;
  defaults: GQL.ConfigDefaultSettingsInput;
  scraping: GQL.ConfigScrapingInput;
  dlna: GQL.ConfigDlnaInput;

  // apikey isn't directly settable, so expose it here
  apiKey: string;

  saveGeneral: (input: Partial<GQL.ConfigGeneralInput>) => void;
  saveInterface: (input: Partial<GQL.ConfigInterfaceInput>) => void;
  saveDefaults: (input: Partial<GQL.ConfigDefaultSettingsInput>) => void;
  saveScraping: (input: Partial<GQL.ConfigScrapingInput>) => void;
  saveDLNA: (input: Partial<GQL.ConfigDlnaInput>) => void;
}

export const SettingStateContext = React.createContext<ISettingsContextState>({
  loading: false,
  error: undefined,
  general: {},
  interface: {},
  defaults: {},
  scraping: {},
  dlna: {},
  apiKey: "",
  saveGeneral: () => {},
  saveInterface: () => {},
  saveDefaults: () => {},
  saveScraping: () => {},
  saveDLNA: () => {},
});

export const SettingsContext: React.FC = ({ children }) => {
  const Toast = useToast();

  const { data, error, loading } = useConfiguration();
  const initialRef = useRef(false);

  const [general, setGeneral] = useState<GQL.ConfigGeneralInput>({});
  const [pendingGeneral, setPendingGeneral] = useState<
    GQL.ConfigGeneralInput | undefined
  >();
  const [updateGeneralConfig] = useConfigureGeneral();

  const [iface, setIface] = useState<GQL.ConfigInterfaceInput>({});
  const [pendingInterface, setPendingInterface] = useState<
    GQL.ConfigInterfaceInput | undefined
  >();
  const [updateInterfaceConfig] = useConfigureInterface();

  const [defaults, setDefaults] = useState<GQL.ConfigDefaultSettingsInput>({});
  const [pendingDefaults, setPendingDefaults] = useState<
    GQL.ConfigDefaultSettingsInput | undefined
  >();
  const [updateDefaultsConfig] = useConfigureDefaults();

  const [scraping, setScraping] = useState<GQL.ConfigScrapingInput>({});
  const [pendingScraping, setPendingScraping] = useState<
    GQL.ConfigScrapingInput | undefined
  >();
  const [updateScrapingConfig] = useConfigureScraping();

  const [dlna, setDLNA] = useState<GQL.ConfigDlnaInput>({});
  const [pendingDLNA, setPendingDLNA] = useState<
    GQL.ConfigDlnaInput | undefined
  >();
  const [updateDLNAConfig] = useConfigureDLNA();

  const [updateSuccess, setUpdateSuccess] = useState(false);

  const [apiKey, setApiKey] = useState("");

  useEffect(() => {
    // only initialise once - assume we have control over these settings and
    // they aren't modified elsewhere
    if (!data?.configuration || error || initialRef.current) return;
    initialRef.current = true;

    setGeneral({ ...withoutTypename(data.configuration.general) });
    setIface({ ...withoutTypename(data.configuration.interface) });
    setDefaults({ ...withoutTypename(data.configuration.defaults) });
    setScraping({ ...withoutTypename(data.configuration.scraping) });
    setDLNA({ ...withoutTypename(data.configuration.dlna) });
    setApiKey(data.configuration.general.apiKey);
  }, [data, error]);

  const resetSuccess = useMemo(
    () =>
      debounce(() => {
        setUpdateSuccess(false);
      }, 4000),
    []
  );

  const onSuccess = useCallback(() => {
    setUpdateSuccess(true);
    resetSuccess();
  }, [resetSuccess]);

  // saves the configuration if no further changes are made after a half second
  const saveGeneralConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigGeneralInput) => {
        try {
          await updateGeneralConfig({
            variables: {
              input,
            },
          });

          setPendingGeneral(undefined);
          onSuccess();
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, updateGeneralConfig, onSuccess]
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
  const saveInterfaceConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigInterfaceInput) => {
        try {
          await updateInterfaceConfig({
            variables: {
              input,
            },
          });

          setPendingInterface(undefined);
          onSuccess();
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, updateInterfaceConfig, onSuccess]
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
  const saveDefaultsConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigDefaultSettingsInput) => {
        try {
          await updateDefaultsConfig({
            variables: {
              input,
            },
          });

          setPendingDefaults(undefined);
          onSuccess();
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, updateDefaultsConfig, onSuccess]
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
  const saveScrapingConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigScrapingInput) => {
        try {
          await updateScrapingConfig({
            variables: {
              input,
            },
          });

          setPendingScraping(undefined);
          onSuccess();
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, updateScrapingConfig, onSuccess]
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
  const saveDLNAConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigDlnaInput) => {
        try {
          await updateDLNAConfig({
            variables: {
              input,
            },
          });

          setPendingDLNA(undefined);
          onSuccess();
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, updateDLNAConfig, onSuccess]
  );

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

  function maybeRenderLoadingIndicator() {
    if (
      pendingGeneral ||
      pendingInterface ||
      pendingDefaults ||
      pendingScraping ||
      pendingDLNA
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
        <div className="loading-indicator">
          <Icon icon="check-circle" className="fa-fw" />
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
        saveGeneral,
        saveInterface,
        saveDefaults,
        saveScraping,
        saveDLNA,
      }}
    >
      {maybeRenderLoadingIndicator()}
      {children}
    </SettingStateContext.Provider>
  );
};
