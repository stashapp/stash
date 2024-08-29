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

  const {
    selectedIds,
    getSelected,
    onSelectChange,
    onSelectAll,
    onSelectNone,
  } = useListSelect(items);

  const state: IListContextState<T> = {
    selectable,
    selectedIds,
    getSelected,
    onSelectChange,
    onSelectAll,
    onSelectNone,
    items,
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

interface IQueryResultContextOptions<
  T extends QueryResult,
  E extends IHasID = IHasID
> {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  useResult: (filter: ListFilterModel) => T;
  getCount: (data: T) => number;
  getItems: (data: T) => E[];
}

export interface IQueryResultContextState<
  T extends QueryResult = QueryResult,
  E extends IHasID = IHasID
> {
  effectiveFilter: ListFilterModel;
  result: T;
  cachedResult: T;
  items: E[];
  totalCount: number;
}

export const QueryResultStateContext =
  React.createContext<IQueryResultContextState | null>(null);

export const QueryResultContext = <
  T extends QueryResult,
  E extends IHasID = IHasID
>(
  props: IQueryResultContextOptions<T, E> & {
    children?:
      | ((props: IQueryResultContextState<T, E>) => React.ReactNode)
      | React.ReactNode;
  }
) => {
  const { filterHook, useResult, getItems, getCount, children } = props;

  const { filter } = useFilter();
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

  const state: IQueryResultContextState<T, E> = {
    effectiveFilter,
    result,
    cachedResult,
    items,
    totalCount,
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
  E extends IHasID = IHasID
>() {
  const context = React.useContext(QueryResultStateContext);

  if (context === null) {
    throw new Error(
      "useQueryResultContext must be used within a ListStateContext"
    );
  }

  return context as IQueryResultContextState<T, E>;
}
