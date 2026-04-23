import { useCallback } from "react";
import { useConfigurationContext } from "./Config";

export function useAutoTagTrigger(
  onRun: () => void,
  onOpenConfirm: () => void,
  override?: boolean | null
) {
  const { configuration } = useConfigurationContext();
  const disabled =
    override ?? configuration?.interface.disableAutoTagWarning ?? false;
  return useCallback(() => {
    if (disabled) {
      onRun();
      return;
    }
    onOpenConfirm();
  }, [disabled, onRun, onOpenConfirm]);
}
