import { useCallback, useMemo } from "react";
import { ListFilterModel } from "src/models/list-filter/filter";

type ListFilterHook = (filter: ListFilterModel) => ListFilterModel;

export const useFavoritesFirstFilterHook = (
  showFavoritesFirst: boolean,
  filterHook?: ListFilterHook
) => {
  const favoritesFirstHook = useCallback((f: ListFilterModel) => {
    const sortBy = f.sortBy ?? "name";

    if (sortBy.startsWith("favorites_first_")) {
      return f;
    }

    if (sortBy === "random") {
      if (f.randomSeed === -1) {
        // Match ListFilterModel random seed generation for stable paging.
        f.randomSeed = Math.floor(Math.random() * 10 ** 8);
      }
      f.sortBy = `favorites_first_random_${f.randomSeed.toString()}`;
      return f;
    }

    f.sortBy = `favorites_first_${sortBy}`;
    return f;
  }, []);

  return useMemo(() => {
    if (!showFavoritesFirst) {
      return filterHook;
    }
    if (!filterHook) {
      return favoritesFirstHook;
    }

    return (f: ListFilterModel) => favoritesFirstHook(filterHook(f));
  }, [showFavoritesFirst, filterHook, favoritesFirstHook]);
};
