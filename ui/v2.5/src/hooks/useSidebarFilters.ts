import { useContext } from "react";
import { useConfigureUI } from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { useToast } from "./Toast";

export interface SidebarFilterDefinition {
  id: string;
  messageId: string;
  defaultVisible?: boolean;
}

export const useSidebarFilters = (
  viewName: string,
  filterDefinitions: SidebarFilterDefinition[]
) => {
  const Toast = useToast();

  const { configuration } = useContext(ConfigurationContext);
  const [saveUI] = useConfigureUI();

  const ui = configuration?.ui;

  // Get default visible filters
  const defaultVisibleFilters = filterDefinitions
    .filter((f) => f.defaultVisible !== false)
    .map((f) => f.id);

  // Get saved visible filters, or use defaults
  const visibleFilters: string[] =
    ui?.sidebarFilters?.[viewName] ?? defaultVisibleFilters;

  // Check if a filter is visible
  const isFilterVisible = (filterId: string): boolean => {
    return visibleFilters.includes(filterId);
  };

  // Save visible filters
  async function saveVisibleFilters(updatedFilters: string[]) {
    try {
      await saveUI({
        variables: {
          input: {
            ...ui,
            sidebarFilters: {
              ...ui?.sidebarFilters,
              [viewName]: updatedFilters,
            },
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  // Toggle a filter's visibility
  async function toggleFilter(filterId: string) {
    const newFilters = visibleFilters.includes(filterId)
      ? visibleFilters.filter((f) => f !== filterId)
      : [...visibleFilters, filterId];
    await saveVisibleFilters(newFilters);
  }

  // Show all filters
  async function showAllFilters() {
    const allFilterIds = filterDefinitions.map((f) => f.id);
    await saveVisibleFilters(allFilterIds);
  }

  // Reset to defaults
  async function resetToDefaults() {
    await saveVisibleFilters(defaultVisibleFilters);
  }

  return {
    visibleFilters,
    isFilterVisible,
    toggleFilter,
    saveVisibleFilters,
    showAllFilters,
    resetToDefaults,
    filterDefinitions,
  };
};

