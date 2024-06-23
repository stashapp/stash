import { useContext, useMemo } from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import * as GQL from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";
import { View } from "./views";

export function useDefaultFilter(mode: GQL.FilterMode, view?: View) {
  const emptyFilter = useMemo(() => new ListFilterModel(mode), [mode]);
  const { configuration: config, loading } = useContext(ConfigurationContext);

  const defaultFilter = useMemo(() => {
    if (view && config?.ui.defaultFilters?.[view]) {
      const savedFilter = config.ui.defaultFilters[view]!;
      const newFilter = emptyFilter.clone();

      newFilter.currentPage = 1;
      try {
        newFilter.configureFromSavedFilter(savedFilter);
      } catch (err) {
        console.log(err);
        // ignore
      }
      // #1507 - reset random seed when loaded
      newFilter.randomSeed = -1;
      return newFilter;
    }
  }, [view, config?.ui.defaultFilters, emptyFilter]);

  const retFilter = loading ? undefined : defaultFilter ?? emptyFilter;

  return { defaultFilter: retFilter, loading };
}
