import React from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import { isFunction } from "lodash-es";
import { useFilterURL } from "./util";

interface IFilterContextOptions {
  filter: ListFilterModel;
  setFilter: React.Dispatch<React.SetStateAction<ListFilterModel>>;
}

export interface IFilterContextState {
  filter: ListFilterModel;
  setFilter: React.Dispatch<React.SetStateAction<ListFilterModel>>;
}

export const FilterStateContext =
  React.createContext<IFilterContextState | null>(null);

export const FilterContext = (
  props: IFilterContextOptions & {
    children?:
      | ((props: IFilterContextState) => React.ReactNode)
      | React.ReactNode;
  }
) => {
  const { filter, setFilter, children } = props;

  const state = {
    filter,
    setFilter,
  };

  return (
    <FilterStateContext.Provider value={state}>
      {isFunction(children)
        ? (children as (props: IFilterContextState) => React.ReactNode)(state)
        : children}
    </FilterStateContext.Provider>
  );
};

export function useFilter() {
  const context = React.useContext(FilterStateContext);

  if (context === null) {
    throw new Error("useFilter must be used within a FilterStateContext");
  }

  return context;
}

// This component is used to set the filter from the URL.
// It replaces the setFilter function to set the URL instead.
// It also loads the default filter if the URL is empty.
export const SetFilterURL = (props: {
  defaultFilter?: ListFilterModel;
  setURL?: boolean;
  children?:
    | ((props: IFilterContextState) => React.ReactNode)
    | React.ReactNode;
}) => {
  const { defaultFilter, setURL = true, children } = props;

  const { filter, setFilter: setFilterOrig } = useFilter();

  const { setFilter } = useFilterURL(filter, setFilterOrig, {
    defaultFilter,
    active: setURL,
  });

  return (
    <FilterContext filter={filter} setFilter={setFilter}>
      {children}
    </FilterContext>
  );
};
