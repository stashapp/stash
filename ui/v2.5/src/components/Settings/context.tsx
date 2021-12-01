import { ApolloError } from "@apollo/client/errors";
import { debounce } from "lodash";
import React, { useState, useEffect, useMemo } from "react";
import { useIntl } from "react-intl";
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
  const intl = useIntl();
  const Toast = useToast();

  const { data, error, loading } = useConfiguration();

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

  const [apiKey, setApiKey] = useState("");

  useEffect(() => {
    if (!data?.configuration || error) return;

    setGeneral({ ...withoutTypename(data.configuration.general) });
    setIface({ ...withoutTypename(data.configuration.interface) });
    setDefaults({ ...withoutTypename(data.configuration.defaults) });
    setScraping({ ...withoutTypename(data.configuration.scraping) });
    setDLNA({ ...withoutTypename(data.configuration.dlna) });
    setApiKey(data.configuration.general.apiKey);
  }, [data, error]);

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

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingGeneral(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateGeneralConfig]
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

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingInterface(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateInterfaceConfig]
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

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingDefaults(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateDefaultsConfig]
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

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingScraping(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateScrapingConfig]
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

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingDLNA(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateDLNAConfig]
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
      {children}
    </SettingStateContext.Provider>
  );
};
