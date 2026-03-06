import { QueryResult } from "@apollo/client";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useShowEditFilter } from "src/components/List/EditFilterDialog";
import { IHasID } from "src/utils/data";
import { useModal } from "src/hooks/modal";
import {
  IFilterStateHook,
  IQueryResultHook,
  useEnsureValidPage,
  useFilterOperations,
  useFilterState,
  useListKeyboardShortcuts,
  useListSelect,
  useQueryResult,
  useScrollToTopOnPageChange,
} from "./util";
import { useConfigurationContext } from "src/hooks/Config";

interface IFilteredItemList<
  T extends QueryResult,
  E extends IHasID = IHasID,
  M = unknown
> {
  filterStateProps: IFilterStateHook;
  queryResultProps: IQueryResultHook<T, E, M>;
}

// Provides the common state and behaviour for filtered item list components
export function useFilteredItemList<
  T extends QueryResult,
  E extends IHasID = IHasID,
  M = unknown
>(props: IFilteredItemList<T, E, M>) {
  const { configuration: config } = useConfigurationContext();

  // States
  const filterState = useFilterState({
    config,
    ...props.filterStateProps,
  });

  const { filter, setFilter } = filterState;

  const queryResult = useQueryResult({
    filter,
    ...props.queryResultProps,
  });
  const { result, items, totalCount, pages, metadataInfo } = queryResult;

  const listSelect = useListSelect(items);
  const { onSelectAll, onSelectNone, onInvertSelection } = listSelect;

  const modalState = useModal();
  const { showModal, closeModal } = modalState;

  // Utility hooks
  const { setPage } = useFilterOperations({ filter, setFilter });

  // scroll to the top of the page when the page changes
  useScrollToTopOnPageChange(filter.currentPage, result.loading);

  // ensure that the current page is valid
  useEnsureValidPage(filter, totalCount, setFilter);

  const showEditFilter = useShowEditFilter({
    showModal,
    closeModal,
    filter,
    setFilter,
  });

  useListKeyboardShortcuts({
    currentPage: filter.currentPage,
    onChangePage: setPage,
    onSelectAll,
    onSelectNone,
    onInvertSelection,
    pages,
    showEditFilter,
  });

  return {
    filterState,
    queryResult,
    metadataInfo,
    listSelect,
    modalState,
    showEditFilter,
  };
}

export const showWhenSelected = <T extends QueryResult>(
  result: T,
  filter: ListFilterModel,
  selectedIds: Set<string>
) => {
  return selectedIds.size > 0;
};

export const showWhenSingleSelection = <T extends QueryResult>(
  result: T,
  filter: ListFilterModel,
  selectedIds: Set<string>
) => {
  return selectedIds.size == 1;
};

export const showWhenNoneSelected = <T extends QueryResult>(
  result: T,
  filter: ListFilterModel,
  selectedIds: Set<string>
) => {
  return selectedIds.size === 0;
};
