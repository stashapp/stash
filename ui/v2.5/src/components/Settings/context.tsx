import { ApolloError } from "@apollo/client/errors";
import { debounce } from "lodash";
import React, { useState, useEffect, useMemo } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { useConfiguration, useConfigureGeneral } from "src/core/StashService";
import { useToast } from "src/hooks";
import { withoutTypename } from "src/utils";

export interface ISettingsContextState {
  loading: boolean;
  error: ApolloError | undefined;
  general: GQL.ConfigGeneralInput;

  // apikey isn't directly settable, so expose it here
  apiKey: string;

  saveGeneral: (input: Partial<GQL.ConfigGeneralInput>) => void;
}

export const SettingStateContext = React.createContext<ISettingsContextState>({
  loading: false,
  error: undefined,
  general: {},
  apiKey: "",
  saveGeneral: () => {},
});

export const SettingsContext: React.FC = ({ children }) => {
  const intl = useIntl();
  const Toast = useToast();

  const { data, error, loading } = useConfiguration();

  const [general, setGeneral] = useState<GQL.ConfigGeneralInput>({});
  const [apiKey, setApiKey] = useState("");
  const [pendingGeneral, setPendingGeneral] = useState<
    GQL.ConfigGeneralInput | undefined
  >();

  const [updateGeneralConfig] = useConfigureGeneral();

  useEffect(() => {
    if (!data?.configuration || error) return;

    setGeneral({ ...withoutTypename(data.configuration.general) });
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

  return (
    <SettingStateContext.Provider
      value={{
        loading,
        error,
        apiKey,
        general,
        saveGeneral,
      }}
    >
      {children}
    </SettingStateContext.Provider>
  );
};
