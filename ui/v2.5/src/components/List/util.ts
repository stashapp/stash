import { useCallback, useContext, useEffect, useMemo, useState } from "react";
import Mousetrap from "mousetrap";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useHistory, useLocation } from "react-router-dom";
import { isEqual, isFunction } from "lodash-es";
import { QueryResult } from "@apollo/client";
import { IHasID } from "src/utils/data";
import { ConfigurationContext } from "src/hooks/Config";
import { View } from "./views";
import { usePrevious } from "src/hooks/state";

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
    if (!active) return;

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
    location.search,
    defaultFilter,
    setFilter,
    updateFilter,
    history,
  ]);

  return { setFilter: updateFilter };
}

export function useDefaultFilter(emptyFilter: ListFilterModel, view?: View) {
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

export function useListSelect<T extends { id: string }>(items: T[]) {
  const [itemsSelected, setItemsSelected] = useState<T[]>([]);
  const [lastClickedId, setLastClickedId] = useState<string>();

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

  const getSelected = useCallback(() => itemsSelected, [itemsSelected]);

  return {
    selectedIds,
    getSelected,
    onSelectChange,
    onSelectAll,
    onSelectNone,
  };
}

export type IListSelect<T extends IHasID> = ReturnType<typeof useListSelect<T>>;

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
