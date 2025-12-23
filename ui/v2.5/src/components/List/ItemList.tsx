import React, {
  PropsWithChildren,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import * as GQL from "src/core/generated-graphql";
import { QueryResult } from "@apollo/client";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  EditFilterDialog,
  useShowEditFilter,
} from "src/components/List/EditFilterDialog";
import { FilterTags } from "./FilterTags";
import { View } from "./views";
import { IHasID } from "src/utils/data";
import {
  ListContext,
  QueryResultContext,
  useListContext,
  useQueryResultContext,
} from "./ListProvider";
import { FilterContext, SetFilterURL, useFilter } from "./FilterProvider";
import { useModal } from "src/hooks/modal";
import {
  IFilterStateHook,
  IQueryResultHook,
  useDefaultFilter,
  useEnsureValidPage,
  useFilterOperations,
  useFilterState,
  useListKeyboardShortcuts,
  useListSelect,
  useQueryResult,
  useScrollToTopOnPageChange,
} from "./util";
import {
  FilteredListToolbar,
  IFilteredListToolbar,
  IItemListOperation,
} from "./FilteredListToolbar";
import { PagedList } from "./PagedList";
import { useConfigurationContext } from "src/hooks/Config";
import { useZoomKeybinds } from "./ZoomSlider";
import { DisplayMode } from "src/models/list-filter/types";

interface IFilteredItemList<T extends QueryResult, E extends IHasID = IHasID> {
  filterStateProps: IFilterStateHook;
  queryResultProps: IQueryResultHook<T, E>;
}

// Provides the common state and behaviour for filtered item list components
export function useFilteredItemList<
  T extends QueryResult,
  E extends IHasID = IHasID
>(props: IFilteredItemList<T, E>) {
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
  const { result, items, totalCount, pages } = queryResult;

  const listSelect = useListSelect(items);
  const { onSelectAll, onSelectNone } = listSelect;

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
    pages,
    showEditFilter,
  });

  return {
    filterState,
    queryResult,
    listSelect,
    modalState,
    showEditFilter,
  };
}

interface IItemListProps<T extends QueryResult, E extends IHasID, M = unknown> {
  view?: View;
  otherOperations?: IItemListOperation<T>[];
  renderContent: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void,
    onChangePage: (page: number) => void,
    pageCount: number
  ) => React.ReactNode;
  renderMetadataByline?: (data: T, metadataInfo?: M) => React.ReactNode;
  renderEditDialog?: (
    selected: E[],
    onClose: (applied: boolean) => void
  ) => React.ReactNode;
  renderDeleteDialog?: (
    selected: E[],
    onClose: (confirmed: boolean) => void
  ) => React.ReactNode;
  addKeybinds?: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => () => void;
  renderToolbar?: (props: IFilteredListToolbar) => React.ReactNode;
}

export const ItemList = <T extends QueryResult, E extends IHasID, M = unknown>(
  props: IItemListProps<T, E, M>
) => {
  const {
    view,
    otherOperations,
    renderContent,
    renderEditDialog,
    renderDeleteDialog,
    renderMetadataByline,
    addKeybinds,
    renderToolbar: providedToolbar,
  } = props;

  const { filter, setFilter: updateFilter } = useFilter();
  const { effectiveFilter, result, metadataInfo, cachedResult, totalCount } =
    useQueryResultContext<T, E, M>();
  const listSelect = useListContext<E>();
  const {
    selectedIds,
    getSelected,
    onSelectChange,
    onSelectAll,
    onSelectNone,
  } = listSelect;

  // scroll to the top of the page when the page changes
  useScrollToTopOnPageChange(filter.currentPage, result.loading);

  const { modal, showModal, closeModal } = useModal();

  const metadataByline = useMemo(() => {
    if (cachedResult.loading) return "";

    return renderMetadataByline?.(cachedResult, metadataInfo) ?? "";
  }, [renderMetadataByline, cachedResult, metadataInfo]);

  const pages = Math.ceil(totalCount / filter.itemsPerPage);

  const onChangePage = useCallback(
    (p: number) => {
      updateFilter(filter.changePage(p));
    },
    [filter, updateFilter]
  );

  useEnsureValidPage(filter, totalCount, updateFilter);

  const showEditFilter = useCallback(
    (editingCriterion?: string) => {
      function onApplyEditFilter(f: ListFilterModel) {
        closeModal();
        updateFilter(f);
      }

      showModal(
        <EditFilterDialog
          filter={filter}
          onApply={onApplyEditFilter}
          onCancel={() => closeModal()}
          editingCriterion={editingCriterion}
        />
      );
    },
    [filter, updateFilter, showModal, closeModal]
  );

  useListKeyboardShortcuts({
    currentPage: filter.currentPage,
    onChangePage,
    onSelectAll,
    onSelectNone,
    pages,
    showEditFilter,
  });

  const zoomable =
    filter.displayMode === DisplayMode.Grid ||
    filter.displayMode === DisplayMode.Wall;

  useZoomKeybinds({
    zoomIndex: zoomable ? filter.zoomIndex : undefined,
    onChangeZoom: (zoom) => updateFilter(filter.setZoom(zoom)),
  });

  useEffect(() => {
    if (addKeybinds) {
      const unbindExtras = addKeybinds(result, effectiveFilter, selectedIds);
      return () => {
        unbindExtras();
      };
    }
  }, [addKeybinds, result, effectiveFilter, selectedIds]);

  const operations = useMemo(() => {
    async function onOperationClicked(o: IItemListOperation<T>) {
      await o.onClick(result, effectiveFilter, selectedIds);
      if (o.postRefetch) {
        result.refetch();
      }
    }

    return otherOperations?.map((o) => ({
      text: o.text,
      onClick: () => {
        onOperationClicked(o);
      },
      isDisplayed: () => {
        if (o.isDisplayed) {
          return o.isDisplayed(result, effectiveFilter, selectedIds);
        }

        return true;
      },
      icon: o.icon,
      buttonVariant: o.buttonVariant,
    }));
  }, [result, effectiveFilter, selectedIds, otherOperations]);

  function onEdit() {
    if (!renderEditDialog) {
      return;
    }

    showModal(
      renderEditDialog(getSelected(), (applied) => onEditDialogClosed(applied))
    );
  }

  function onEditDialogClosed(applied: boolean) {
    if (applied) {
      onSelectNone();
    }
    closeModal();

    // refetch
    result.refetch();
  }

  function onDelete() {
    if (!renderDeleteDialog) {
      return;
    }

    showModal(
      renderDeleteDialog(getSelected(), (deleted) =>
        onDeleteDialogClosed(deleted)
      )
    );
  }

  function onDeleteDialogClosed(deleted: boolean) {
    if (deleted) {
      onSelectNone();
    }
    closeModal();

    // refetch
    result.refetch();
  }

  function onRemoveCriterion(removedCriterion: Criterion, valueIndex?: number) {
    if (valueIndex === undefined) {
      updateFilter(
        filter.removeCriterion(removedCriterion.criterionOption.type)
      );
    } else {
      updateFilter(
        filter.removeCustomFieldCriterion(
          removedCriterion.criterionOption.type,
          valueIndex
        )
      );
    }
  }

  function onClearAllCriteria() {
    updateFilter(filter.clearCriteria());
  }

  const filterListToolbarProps: IFilteredListToolbar = {
    filter,
    setFilter: updateFilter,
    listSelect,
    showEditFilter,
    view: view,
    operations: operations,
    zoomable: zoomable,
    onEdit: renderEditDialog ? onEdit : undefined,
    onDelete: renderDeleteDialog ? onDelete : undefined,
  };

  return (
    <div className="item-list-container">
      {providedToolbar ? (
        providedToolbar(filterListToolbarProps)
      ) : (
        <FilteredListToolbar {...filterListToolbarProps} />
      )}
      <FilterTags
        criteria={filter.criteria}
        onEditCriterion={(c) => showEditFilter(c.criterionOption.type)}
        onRemoveCriterion={onRemoveCriterion}
        onRemoveAll={() => onClearAllCriteria()}
      />
      {modal}

      <PagedList
        result={result}
        cachedResult={cachedResult}
        filter={filter}
        totalCount={totalCount}
        onChangePage={onChangePage}
        metadataByline={metadataByline}
      >
        {renderContent(
          result,
          // #4780 - use effectiveFilter to ensure filterHook is applied
          effectiveFilter,
          selectedIds,
          onSelectChange,
          onChangePage,
          pages
        )}
      </PagedList>
    </div>
  );
};

interface IItemListContextProps<
  T extends QueryResult,
  E extends IHasID,
  M = unknown
> {
  filterMode: GQL.FilterMode;
  defaultSort?: string;
  defaultFilter?: ListFilterModel;
  useResult: (filter: ListFilterModel) => T;
  useMetadataInfo?: (filter: ListFilterModel) => M;
  getCount: (data: T) => number;
  getItems: (data: T) => E[];
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
  selectable?: boolean;
}

// Provides the contexts for the ItemList component. Includes functionality to scroll
// to top on page change.
export const ItemListContext = <
  T extends QueryResult,
  E extends IHasID,
  M = unknown
>(
  props: PropsWithChildren<IItemListContextProps<T, E, M>>
) => {
  const {
    filterMode,
    defaultSort,
    defaultFilter: providedDefaultFilter,
    useResult,
    useMetadataInfo,
    getCount,
    getItems,
    view,
    filterHook,
    alterQuery = true,
    selectable,
    children,
  } = props;

  const { configuration: config } = useConfigurationContext();

  const emptyFilter = useMemo(
    () =>
      providedDefaultFilter?.clone() ??
      new ListFilterModel(filterMode, config, {
        defaultSortBy: defaultSort,
      }),
    [config, filterMode, defaultSort, providedDefaultFilter]
  );

  const [filter, setFilterState] = useState<ListFilterModel>(
    () =>
      new ListFilterModel(filterMode, config, { defaultSortBy: defaultSort })
  );

  const { defaultFilter } = useDefaultFilter(emptyFilter, view);

  return (
    <FilterContext filter={filter} setFilter={setFilterState}>
      <SetFilterURL defaultFilter={defaultFilter} setURL={alterQuery}>
        <QueryResultContext
          filterHook={filterHook}
          useResult={useResult}
          useMetadataInfo={useMetadataInfo}
          getCount={getCount}
          getItems={getItems}
        >
          {({ items }) => (
            <ListContext selectable={selectable} items={items}>
              {children}
            </ListContext>
          )}
        </QueryResultContext>
      </SetFilterURL>
    </FilterContext>
  );
};

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
