import React, { useMemo } from "react";
import { IListSelect, useCachedQueryResult, useListSelect } from "./util";
import { isFunction } from "lodash-es";
import { IHasID } from "src/utils/data";
import { useFilter } from "./FilterProvider";
import { ListFilterModel } from "src/models/list-filter/filter";
import { QueryResult } from "@apollo/client";

interface IListContextOptions<T extends IHasID> {
  selectable?: boolean;
  items: T[];
}

export type IListContextState<T extends IHasID = IHasID> = IListSelect<T> & {
  selectable: boolean;
  items: T[];
};

export const ListStateContext = React.createContext<IListContextState | null>(
  null
);

export const ListContext = <T extends IHasID = IHasID>(
  props: IListContextOptions<T> & {
    children?:
      | ((props: IListContextState) => React.ReactNode)
      | React.ReactNode;
  }
) => {
  const { selectable = false, items, children } = props;

  const listSelect = useListSelect(items);

  const state: IListContextState<T> = {
    selectable,
    items,
    ...listSelect,
  };

  return (
    <ListStateContext.Provider value={state}>
      {isFunction(children)
        ? (children as (props: IListContextState) => React.ReactNode)(state)
        : children}
    </ListStateContext.Provider>
  );
};

export function useListContext<T extends IHasID = IHasID>() {
  const context = React.useContext(ListStateContext);

  if (context === null) {
    throw new Error("useListContext must be used within a ListStateContext");
  }

  return context as IListContextState<T>;
}

const emptyState: IListContextState = {
  selectable: false,
  selectedIds: new Set(),
  getSelected: () => [],
  onSelectChange: () => {},
  onSelectAll: () => {},
  onSelectNone: () => {},
  items: [],
  hasSelection: false,
  selectedItems: [],
};

export function useListContextOptional<T extends IHasID = IHasID>() {
  const context = React.useContext(ListStateContext);

  if (context === null) {
    return emptyState as IListContextState<T>;
  }

  return context as IListContextState<T>;
}

interface IQueryResultContextOptions<
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

export interface IQueryResultContextState<
  T extends QueryResult = QueryResult,
  E extends IHasID = IHasID,
  M = unknown
> {
  effectiveFilter: ListFilterModel;
  result: T;
  cachedResult: T;
  metadataInfo?: M;
  items: E[];
  totalCount: number;
}

export const QueryResultStateContext =
  React.createContext<IQueryResultContextState | null>(null);

export const QueryResultContext = <
  T extends QueryResult,
  E extends IHasID = IHasID,
  M = unknown
>(
  props: IQueryResultContextOptions<T, E, M> & {
    children?:
      | ((props: IQueryResultContextState<T, E, M>) => React.ReactNode)
      | React.ReactNode;
  }
) => {
  const {
    filterHook,
    useResult,
    useMetadataInfo,
    getItems,
    getCount,
    children,
  } = props;

  const { filter } = useFilter();
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

  // use cached query result for pagination
  const cachedResult = useCachedQueryResult(effectiveFilter, result);

  const items = useMemo(() => getItems(result), [getItems, result]);
  const totalCount = useMemo(
    () => getCount(cachedResult),
    [getCount, cachedResult]
  );

  const state: IQueryResultContextState<T, E, M> = {
    effectiveFilter,
    result,
    cachedResult,
    items,
    totalCount,
    metadataInfo,
  };

  return (
    <QueryResultStateContext.Provider value={state}>
      {isFunction(children)
        ? (children as (props: IQueryResultContextState) => React.ReactNode)(
            state
          )
        : children}
    </QueryResultStateContext.Provider>
  );
};

export function useQueryResultContext<
  T extends QueryResult,
  E extends IHasID = IHasID,
  M = unknown
>() {
  const context = React.useContext(QueryResultStateContext);

  if (context === null) {
    throw new Error(
      "useQueryResultContext must be used within a ListStateContext"
    );
  }

  return context as IQueryResultContextState<T, E, M>;
}
