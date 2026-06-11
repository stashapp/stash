import { useCallback } from "react";
import { IUIConfig } from "src/core/config";
import { useConfigurationContext } from "./Config";

export function useAutoTagTrigger(
  onRun: () => void,
  onOpenConfirm: () => void,
  override?: boolean | null
) {
  const { configuration } = useConfigurationContext();
  const ui = configuration?.ui as IUIConfig | undefined;
  const disabled = override ?? ui?.disableAutoTagWarning ?? false;
  return useCallback(() => {
    if (disabled) {
      onRun();
      return;
    }
    onOpenConfirm();
  }, [disabled, onRun, onOpenConfirm]);
}
