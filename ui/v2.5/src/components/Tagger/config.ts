import { useCallback } from "react";
import { useConfigurationContext } from "src/hooks/Config";
import { initialConfig, ITaggerConfig } from "./constants";
import { useConfigureUISetting } from "src/core/StashService";

export function useTaggerConfig() {
  const { configuration: stashConfig } = useConfigurationContext();
  const [saveUISetting] = useConfigureUISetting();

  const config = stashConfig?.ui.taggerConfig ?? initialConfig;

  const setConfig = useCallback(
    (c: ITaggerConfig) => {
      saveUISetting({ variables: { key: "taggerConfig", value: c } });
    },
    [saveUISetting]
  );

  return { config, setConfig };
}
