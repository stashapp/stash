import { useCallback, useEffect, useMemo, useState } from "react";
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

function locationEquals(
  loc1: ReturnType<typeof useLocation> | undefined,
  loc2: ReturnType<typeof useLocation>
) {
  return loc1 && loc1.pathname === loc2.pathname && loc1.search === loc2.search;
}

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

export interface IFilterStateHook {
  filterMode: GQL.FilterMode;
  defaultSort?: string;
  view?: View;
  useURL?: boolean;
}

export function useFilterState(
  props: IFilterStateHook & {
    config?: GQL.ConfigDataFragment;
  }
) {
  const { filterMode, defaultSort, config, view, useURL } = props;

  const [filter, setFilterState] = useState<ListFilterModel>(
    () =>
      new ListFilterModel(filterMode, config, { defaultSortBy: defaultSort })
  );

  const emptyFilter = useEmptyFilter({ filterMode, defaultSort, config });

  const { defaultFilter } = useDefaultFilter(emptyFilter, view);

  const { setFilter } = useFilterURL(filter, setFilterState, {
    defaultFilter,
    active: useURL,
  });

  return { filter, setFilter };
}

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

export function useListKeyboardShortcuts(props: {
  currentPage?: number;
  onChangePage?: (page: number) => void;
  showEditFilter?: () => void;
  pages?: number;
  onSelectAll?: () => void;
  onSelectNone?: () => void;
}) {
  const {
    currentPage,
    onChangePage,
    showEditFilter,
    pages = 0,
    onSelectAll,
    onSelectNone,
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

    return () => {
      Mousetrap.unbind("s a");
      Mousetrap.unbind("s n");
    };
  }, [onSelectAll, onSelectNone]);
}

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
    hasSelection,
  };
}

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

// this hook caches a query result and count, and only updates it when the filter changes
// in a way that would impact the result count
// it is used to prevent the result count/pagination from flickering when changing pages or sorting
export function useCachedQueryResult<T extends QueryResult>(
  filter: ListFilterModel,
  result: T
) {
  const [cachedResult, setCachedResult] = useState(result);
  const [lastFilter, setLastFilter] = useState(filter);

  // if we are only changing the page or sort, don't update the result count
  useEffect(() => {
    if (!result.loading) {
      setCachedResult(result);
    } else {
      if (totalCountImpacted(lastFilter, filter)) {
        setCachedResult(result);
      }
    }

    setLastFilter(filter);
  }, [filter, result, lastFilter]);

  return cachedResult;
}

export interface IQueryResultHook<
  T extends QueryResult,
  E extends IHasID = IHasID
> {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  useResult: (filter: ListFilterModel) => T;
  getCount: (data: T) => number;
  getItems: (data: T) => E[];
}

export function useQueryResult<
  T extends QueryResult,
  E extends IHasID = IHasID
>(
  props: IQueryResultHook<T, E> & {
    filter: ListFilterModel;
  }
) {
  const { filter, filterHook, useResult, getItems, getCount } = props;

  const effectiveFilter = useMemo(() => {
    if (filterHook) {
      return filterHook(filter.clone());
    }
    return filter;
  }, [filter, filterHook]);

  const result = useResult(effectiveFilter);

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
    result,
    cachedResult,
    items,
    totalCount,
    pages,
  };
}

// this hook collects the common logic when closing the edit/delete dialog
// if applied is true, then the list should be refetched and selection cleared
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

export function useScrollToTopOnPageChange(
  currentPage: number,
  loading: boolean
) {
  const prevPage = usePrevious(currentPage);

  // scroll to the top of the page when the page changes
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

// handle case where page is more than there are pages
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
