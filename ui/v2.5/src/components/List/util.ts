import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import Mousetrap from "mousetrap";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useHistory, useLocation } from "react-router-dom";
import { isEqual, isFunction } from "lodash-es";
import { QueryResult } from "@apollo/client";
import { IHasID } from "src/utils/data";
import { useConfigurationContext } from "src/hooks/Config";
import { View } from "./views";
import { usePrevious } from "src/hooks/state";
import * as GQL from "src/core/generated-graphql";
import { DisplayMode } from "src/models/list-filter/types";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { useInterfaceLocalForage } from "src/hooks/LocalForage";

function locationEquals(
  loc1: ReturnType<typeof useLocation> | undefined,
  loc2: ReturnType<typeof useLocation>
) {
  return loc1 && loc1.pathname === loc2.pathname && loc1.search === loc2.search;
}

/**
 * Sync a filter to/from the URL query string.
 * When the filter changes, the URL is updated via `history.replace`.
 * When the location changes, the filter is re-parsed from the query string.
 * @param filter - Current filter state.
 * @param setFilter - Raw React state setter (not the URL-synced wrapper).
 * @param options.defaultFilter - Fallback filter used when the URL has no query string.
 * @param options.active - Whether URL sync is active (e.g. false for hidden detail panels).
 */
export function useFilterURL(
  filter: ListFilterModel,
  setFilter: React.Dispatch<React.SetStateAction<ListFilterModel>>,
  options?: {
    defaultFilter?: ListFilterModel;
    active?: boolean;
  }
) {
  const { defaultFilter, active = true } = options ?? {};

  const history = useHistory();
  const location = useLocation();
  const prevLocation = usePrevious(location);

  // when the filter changes, update the URL
  const updateFilter = useCallback(
    (
      value: ListFilterModel | ((prevState: ListFilterModel) => ListFilterModel)
    ) => {
      const newFilter = isFunction(value) ? value(filter) : value;

      if (active) {
        const newParams = newFilter.makeQueryParameters();
        history.replace({ ...history.location, search: newParams });
      } else {
        // set the filter without updating the URL
        setFilter(newFilter);
      }
    },
    [history, active, setFilter, filter]
  );

  // This hook runs on every page location change (ie navigation),
  // and updates the filter accordingly.
  useEffect(() => {
    // don't apply if active is false
    // also don't apply if location is unchanged
    if (!active || locationEquals(prevLocation, location)) return;

    // re-init to load default filter on empty new query params
    if (!location.search) {
      if (defaultFilter) updateFilter(defaultFilter.clone());
      return;
    }

    // the query has changed, update filter if necessary
    setFilter((prevFilter) => {
      let newFilter = prevFilter.empty();
      newFilter.configureFromQueryString(location.search);
      if (!isEqual(newFilter, prevFilter)) {
        // filter may have changed if random seed was set, update the URL
        const newParams = newFilter.makeQueryParameters();
        if (newParams !== location.search) {
          history.replace({ ...history.location, search: newParams });
        }

        return newFilter;
      } else {
        return prevFilter;
      }
    });
  }, [
    active,
    prevLocation,
    location,
    defaultFilter,
    setFilter,
    updateFilter,
    history,
  ]);

  return { setFilter: updateFilter };
}

/**
 * Load a user's saved default filter for the given view from server config.
 * Falls back to the empty filter if no saved default exists.
 * @param emptyFilter - A base filter to clone and apply saved options onto.
 * @param view - The view to look up a saved default for.
 */
export function useDefaultFilter(emptyFilter: ListFilterModel, view?: View) {
  const { configuration: config } = useConfigurationContext();

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

  const retFilter = defaultFilter ?? emptyFilter;

  return { defaultFilter: retFilter };
}

function useEmptyFilter(props: {
  filterMode: GQL.FilterMode;
  defaultSort?: string;
  config?: GQL.ConfigDataFragment;
}) {
  const { filterMode, defaultSort, config } = props;

  const emptyFilter = useMemo(
    () =>
      new ListFilterModel(filterMode, config, {
        defaultSortBy: defaultSort,
      }),
    [config, filterMode, defaultSort]
  );

  return emptyFilter;
}

/** Props shared by all filtered-list hooks. */
export interface IFilterStateHook {
  filterMode: GQL.FilterMode;
  defaultFilter?: ListFilterModel;
  defaultSort?: string;
  view?: View;
  useURL?: boolean;
}

/**
 * Top-level hook that wires together filter state, URL sync, default filters,
 * and persisted display-mode for a filtered list page.
 * @param props.filterMode - The GraphQL FilterMode enum value.
 * @param props.defaultSort - Default sort field name.
 * @param props.config - Server configuration (injected by the caller).
 * @param props.view - The current view, used for saved defaults and display-mode persistence.
 * @param props.useURL - Whether to sync filter state to/from the URL query string.
 * @param props.defaultFilter - Optional override default filter.
 */
export function useFilterState(
  props: IFilterStateHook & {
    config?: GQL.ConfigDataFragment;
  }
) {
  const {
    filterMode,
    defaultSort,
    config,
    view,
    useURL,
    defaultFilter: propDefaultFilter,
  } = props;

  const [filter, setFilterState] = useState<ListFilterModel>(
    () =>
      new ListFilterModel(filterMode, config, { defaultSortBy: defaultSort })
  );

  const emptyFilter = useEmptyFilter({ filterMode, defaultSort, config });

  const { defaultFilter: defaultFilterFromConfig } = useDefaultFilter(
    emptyFilter,
    view
  );

  const { setFilter } = useFilterURL(filter, setFilterState, {
    defaultFilter: propDefaultFilter ?? defaultFilterFromConfig,
    active: useURL,
  });

  const location = useLocation();
  const [
    { data: interfaceData, loading: interfaceLoading },
    setInterfaceLocalForage,
  ] = useInterfaceLocalForage();

  // on mount, restore persisted displayMode if URL doesn't specify one
  useEffect(() => {
    if (interfaceLoading || !interfaceData || !view) return;
    const persisted = interfaceData.viewConfig?.[view]?.displayMode;
    if (persisted === undefined) return;
    if (location.search.includes("disp=")) return;

    setFilter((cv) => {
      if (cv.displayMode === persisted) return cv;
      return cv.setDisplayMode(persisted);
    });
  }, [view, interfaceData, interfaceLoading, location.search, setFilter]);

  // persist displayMode on change
  const prevDisplayMode = usePrevious(filter.displayMode);
  useEffect(() => {
    if (!view || interfaceLoading) return;
    if (prevDisplayMode === undefined) return;
    if (filter.displayMode === prevDisplayMode) return;

    setInterfaceLocalForage((prev) => ({
      ...prev,
      viewConfig: {
        ...prev.viewConfig,
        [view]: {
          ...prev.viewConfig?.[view],
          displayMode: filter.displayMode,
        },
      },
    }));
  }, [
    view,
    filter.displayMode,
    prevDisplayMode,
    interfaceLoading,
    setInterfaceLocalForage,
  ]);

  return { filter, setFilter };
}

/**
 * Convenience wrapper around common filter mutations (page, display mode, zoom, criteria).
 * All returned callbacks are memoised and stable.
 */
export function useFilterOperations(props: {
  filter: ListFilterModel;
  setFilter: (
    value: ListFilterModel | ((prevState: ListFilterModel) => ListFilterModel)
  ) => void;
}) {
  const { setFilter } = props;

  const setPage = useCallback(
    (p: number) => {
      setFilter((cv) => cv.changePage(p));
    },
    [setFilter]
  );

  const setDisplayMode = useCallback(
    (displayMode: DisplayMode) => {
      setFilter((cv) => cv.setDisplayMode(displayMode));
    },
    [setFilter]
  );

  const setZoom = useCallback(
    (newZoomIndex: number) => {
      setFilter((cv) => cv.setZoom(newZoomIndex));
    },
    [setFilter]
  );

  const removeCriterion = useCallback(
    (removedCriterion: Criterion) => {
      setFilter((cv) =>
        cv.removeCriterion(removedCriterion.criterionOption.type)
      );
    },
    [setFilter]
  );

  const clearAllCriteria = useCallback(
    (includeSearchTerm = false) => {
      setFilter((cv) => cv.clearCriteria(includeSearchTerm));
    },
    [setFilter]
  );

  return {
    setPage,
    setDisplayMode,
    setZoom,
    removeCriterion,
    clearAllCriteria,
  };
}

/**
 * Bind keyboard shortcuts for list navigation (left/right for pagination,
 * "f" for filter, "s a/n/i" for select-all/none/invert).
 * Shortcuts are scoped to the component lifecycle.
 */
export function useListKeyboardShortcuts(props: {
  currentPage?: number;
  onChangePage?: (page: number) => void;
  showEditFilter?: () => void;
  pages?: number;
  onSelectAll?: () => void;
  onSelectNone?: () => void;
  onInvertSelection?: () => void;
}) {
  const {
    currentPage,
    onChangePage,
    showEditFilter,
    pages = 0,
    onSelectAll,
    onSelectNone,
    onInvertSelection,
  } = props;

  // set up hotkeys
  useEffect(() => {
    if (showEditFilter) {
      Mousetrap.bind("f", (e) => {
        showEditFilter();
        // prevent default behavior of typing f in a text field
        // otherwise the filter dialog closes, the query field is focused and
        // f is typed.
        e.preventDefault();
      });

      return () => {
        Mousetrap.unbind("f");
      };
    }
  }, [showEditFilter]);

  useEffect(() => {
    if (!currentPage || !changePage || !pages) return;

    function changePage(page: number) {
      if (!currentPage || !onChangePage || !pages) return;
      if (page >= 1 && page <= pages) {
        onChangePage(page);
      }
    }

    Mousetrap.bind("right", () => {
      changePage(currentPage + 1);
    });
    Mousetrap.bind("left", () => {
      changePage(currentPage - 1);
    });
    Mousetrap.bind("shift+right", () => {
      changePage(Math.min(pages, currentPage + 10));
    });
    Mousetrap.bind("shift+left", () => {
      changePage(Math.max(1, currentPage - 10));
    });
    Mousetrap.bind("ctrl+end", () => {
      changePage(pages);
    });
    Mousetrap.bind("ctrl+home", () => {
      changePage(1);
    });

    return () => {
      Mousetrap.unbind("right");
      Mousetrap.unbind("left");
      Mousetrap.unbind("shift+right");
      Mousetrap.unbind("shift+left");
      Mousetrap.unbind("ctrl+end");
      Mousetrap.unbind("ctrl+home");
    };
  }, [currentPage, onChangePage, pages]);

  useEffect(() => {
    Mousetrap.bind("s a", () => onSelectAll?.());
    Mousetrap.bind("s n", () => onSelectNone?.());
    Mousetrap.bind("s i", () => onInvertSelection?.());

    return () => {
      Mousetrap.unbind("s a");
      Mousetrap.unbind("s n");
      Mousetrap.unbind("s i");
    };
  }, [onSelectAll, onSelectNone, onInvertSelection]);
}

/**
 * Track multi-select state for a list of items with shift-click support.
 * @param items - The current page of items.
 * @returns Selected items, selected IDs set, and selection helper callbacks.
 */
export function useListSelect<T extends IHasID = IHasID>(items: T[]) {
  const [itemsSelected, setItemsSelected] = useState<T[]>([]);
  const [lastClickedId, setLastClickedId] = useState<string>();

  // TODO - this doesn't get updated when items changes
  const selectedIds = useMemo(() => {
    const newSelectedIds = new Set<string>();
    itemsSelected.forEach((item) => {
      newSelectedIds.add(item.id);
    });

    return newSelectedIds;
  }, [itemsSelected]);

  // const prevItems = usePrevious(items);

  // #5341 - HACK/TODO: this is a regression of previous behaviour. I don't like the idea
  // of keeping selected items that are no longer in the list, since its not
  // clear to the user that the item is still selected, but there is now an expectation of
  // this behaviour.
  // useEffect(() => {
  //   if (prevItems === items) {
  //     return;
  //   }

  //   // filter out any selectedIds that are no longer in the list
  //   const newSelectedIds = new Set<string>();

  //   selectedIds.forEach((id) => {
  //     if (items.some((item) => item.id === id)) {
  //       newSelectedIds.add(id);
  //     }
  //   });

  //   setSelectedIds(newSelectedIds);
  // }, [prevItems, items, selectedIds]);

  function singleSelect(id: string, selected: boolean) {
    setLastClickedId(id);

    setItemsSelected((prevSelected) => {
      if (selected) {
        // prevent duplicates
        if (prevSelected.some((v) => v.id === id)) {
          return prevSelected;
        }

        const item = items.find((i) => i.id === id);
        if (item) {
          return [...prevSelected, item];
        }
        return prevSelected;
      } else {
        return prevSelected.filter((item) => item.id !== id);
      }
    });
  }

  function selectRange(startIndex: number, endIndex: number) {
    let start = startIndex;
    let end = endIndex;
    if (start > end) {
      const tmp = start;
      start = end;
      end = tmp;
    }

    const subset = items.slice(start, end + 1);

    // prevent duplicates
    const toAdd = subset.filter((item) => !selectedIds.has(item.id));

    const newSelected = itemsSelected.concat(toAdd);
    setItemsSelected(newSelected);
  }

  function multiSelect(id: string) {
    let startIndex = 0;
    let thisIndex = -1;

    if (lastClickedId) {
      startIndex = items.findIndex((item) => {
        return item.id === lastClickedId;
      });
    }

    thisIndex = items.findIndex((item) => {
      return item.id === id;
    });

    selectRange(startIndex, thisIndex);
  }

  function onSelectChange(id: string, selected: boolean, shiftKey: boolean) {
    if (shiftKey) {
      multiSelect(id);
    } else {
      singleSelect(id, selected);
    }
  }

  function onSelectAll() {
    // #5341 - HACK/TODO: maintaining legacy behaviour of replacing selected items with
    // all items on the current page. To be consistent with the existing behaviour, it
    // should probably _add_ all items on the current page to the selected items.
    setItemsSelected([...items]);
    setLastClickedId(undefined);
  }

  function onSelectNone() {
    setItemsSelected([]);
    setLastClickedId(undefined);
  }

  function onInvertSelection() {
    setItemsSelected((prevSelected) => {
      const selectedSet = new Set(prevSelected.map((item) => item.id));
      return items.filter((item) => !selectedSet.has(item.id));
    });
    setLastClickedId(undefined);
  }

  // TODO - this is for backwards compatibility
  const getSelected = useCallback(() => itemsSelected, [itemsSelected]);

  // convenience state
  const hasSelection = itemsSelected.length > 0;

  return {
    selectedItems: itemsSelected,
    selectedIds,
    getSelected,
    onSelectChange,
    onSelectAll,
    onSelectNone,
    onInvertSelection,
    hasSelection,
  };
}

/** Inferred return type of {@link useListSelect}. */
export type IListSelect<T extends IHasID = IHasID> = ReturnType<
  typeof useListSelect<T>
>;

// returns true if the filter has changed in a way that impacts the total count
function totalCountImpacted(
  oldFilter: ListFilterModel,
  newFilter: ListFilterModel
) {
  return (
    oldFilter.criteria.length !== newFilter.criteria.length ||
    oldFilter.criteria.some((c) => {
      const newCriterion = newFilter.criteria.find(
        (nc) => nc.getId() === c.getId()
      );
      return !newCriterion || !isEqual(c, newCriterion);
    })
  );
}

/**
 * Cache a query result and only update when the filter changes in a way that
 * affects the total count. Prevents pagination flicker during page/sort changes.
 * @param filter - The current filter.
 * @param result - The latest Apollo query result.
 */
export function useCachedQueryResult<T extends QueryResult>(
  filter: ListFilterModel,
  result: T
) {
  const [cachedResult, setCachedResult] = useState(result);
  const lastFilterRef = useRef(filter);

  // if we are only changing the page or sort, don't update the result count
  useEffect(() => {
    if (!result.loading) {
      setCachedResult(result);
    } else {
      if (totalCountImpacted(lastFilterRef.current, filter)) {
        setCachedResult(result);
      }
    }

    lastFilterRef.current = filter;
  }, [filter, result]);

  return cachedResult;
}

/** Configuration for a list query result hook. */
export interface IQueryResultHook<
  T extends QueryResult,
  E extends IHasID = IHasID,
  M = unknown
> {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  useResult: (filter: ListFilterModel) => T;
  useMetadataInfo?: (filter: ListFilterModel) => M;
  getCount: (data: T) => number;
  getItems: (data: T) => E[];
}

/**
 * Execute a GraphQL query driven by a ListFilterModel and return paginated
 * results with caching metadata info (total count without page/sort offset).
 * @param props.filterHook - Optional transform applied to the filter before querying.
 * @param props.useResult - Apollo hook that maps a filter to a query result.
 * @param props.useMetadataInfo - Optional hook for fetching unfiltered metadata counts.
 * @param props.getCount - Extract total count from the query result.
 * @param props.getItems - Extract item list from the query result.
 * @param props.filter - The current filter.
 */
export function useQueryResult<
  T extends QueryResult,
  E extends IHasID = IHasID,
  M = unknown
>(
  props: IQueryResultHook<T, E, M> & {
    filter: ListFilterModel;
  }
) {
  const { filter, filterHook, useResult, useMetadataInfo, getItems, getCount } =
    props;

  const effectiveFilter = useMemo(() => {
    if (filterHook) {
      return filterHook(filter.clone());
    }
    return filter;
  }, [filter, filterHook]);

  // metadata filter is the effective filter with the sort, page size and page number removed
  const metadataFilter = useMemo(
    () => effectiveFilter.metadataInfo(),
    [effectiveFilter]
  );

  const result = useResult(effectiveFilter);
  const metadataInfo = useMetadataInfo?.(metadataFilter);

  // use cached query result for pagination and metadata rendering
  const cachedResult = useCachedQueryResult(effectiveFilter, result);

  const items = useMemo(() => getItems(result), [getItems, result]);
  const totalCount = useMemo(
    () => getCount(cachedResult),
    [getCount, cachedResult]
  );

  const pages = Math.ceil(totalCount / filter.itemsPerPage);

  return {
    effectiveFilter,
    metadataInfo,
    result,
    cachedResult,
    items,
    totalCount,
    pages,
  };
}

/**
 * Collect common logic for closing an edit/delete dialog: close the modal,
 * clear selection, and refetch the list if changes were applied.
 * @param props.onSelectNone - Clear current selection.
 * @param props.closeModal - Close the dialog.
 * @param props.result - Apollo query result to refetch.
 * @returns A callback to invoke with `true` if changes were applied.
 */
export function useCloseEditDelete(props: {
  onSelectNone: () => void;
  closeModal: () => void;
  result: QueryResult;
}) {
  const { onSelectNone, closeModal, result } = props;

  const onCloseEditDelete = useCallback(
    (applied?: boolean) => {
      closeModal();
      if (applied) {
        onSelectNone();

        // refetch
        result.refetch();
      }
    },
    [onSelectNone, closeModal, result]
  );

  return onCloseEditDelete;
}

/**
 * Scroll to the top of the page (or below the detail header) when the page changes.
 * Only scrolls when `loading` is false and the page has actually changed.
 */
export function useScrollToTopOnPageChange(
  currentPage: number,
  loading: boolean
) {
  const prevPage = usePrevious(currentPage);

  // only scroll to top if the page has changed and is not loading
  useEffect(() => {
    if (loading || currentPage === prevPage || prevPage === undefined) {
      return;
    }

    // if the current page has a detail-header, then
    // scroll up relative to that rather than 0, 0
    const detailHeader = document.querySelector(".detail-header");
    if (detailHeader) {
      window.scrollTo(0, detailHeader.scrollHeight - 50);
    } else {
      window.scrollTo(0, 0);
    }
  }, [prevPage, currentPage, loading]);
}

/**
 * Ensure the current page does not exceed the total number of pages.
 * If it does, clamp to the last page.
 */
export function useEnsureValidPage(
  filter: ListFilterModel,
  totalCount: number,
  setFilter: React.Dispatch<React.SetStateAction<ListFilterModel>>
) {
  useEffect(() => {
    const totalPages = Math.ceil(totalCount / filter.itemsPerPage);

    if (totalPages > 0 && filter.currentPage > totalPages) {
      setFilter((prevFilter) => prevFilter.changePage(totalPages));
    }
  }, [filter, totalCount, setFilter]);
}
