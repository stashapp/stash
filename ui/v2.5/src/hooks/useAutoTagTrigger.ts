import { useConfigurationContext } from "./Config";

// Centralises the "skip warning" check so every entry point to auto tag
// consults the same interface flag.
export function useAutoTagTrigger(
  onRun: () => void,
  onOpenConfirm: () => void
) {
  const { configuration } = useConfigurationContext();
  return () => {
    if (configuration?.interface.disableAutoTagWarning) {
      onRun();
      return;
    }
    onOpenConfirm();
  };
}