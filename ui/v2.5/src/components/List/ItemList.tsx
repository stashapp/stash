import React, { useCallback, useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { QueryResult } from "@apollo/client";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { EditFilterDialog } from "src/components/List/EditFilterDialog";
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
  useDefaultFilter,
  useEnsureValidPage,
  useListKeyboardShortcuts,
  useScrollToTopOnPageChange,
} from "./util";
import { FilteredListToolbar, IItemListOperation } from "./FilteredListToolbar";
import { PagedList } from "./PagedList";

interface IItemListOptions<T extends QueryResult, E extends IHasID> {
  filterMode: GQL.FilterMode;
  useResult: (filter: ListFilterModel) => T;
  getCount: (data: T) => number;
  renderMetadataByline?: (data: T) => React.ReactNode;
  getItems: (data: T) => E[];
}

interface IItemListProps<T extends QueryResult, E extends IHasID> {
  view?: View;
  defaultSort?: string;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  filterDialog?: (
    criteria: Criterion<CriterionValue>[],
    setCriteria: (v: Criterion<CriterionValue>[]) => void
  ) => React.ReactNode;
  zoomable?: boolean;
  selectable?: boolean;
  alterQuery?: boolean;
  defaultZoomIndex?: number;
  otherOperations?: IItemListOperation<T>[];
  renderContent: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void,
    onChangePage: (page: number) => void,
    pageCount: number
  ) => React.ReactNode;
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
}

/**
 * A factory function for ItemList components.
 * IMPORTANT: as the component manipulates the URL query string, if there are
 * ever multiple ItemLists rendered at once, all but one of them need to have
 * `alterQuery` set to false to prevent conflicts.
 */
export function makeItemList<T extends QueryResult, E extends IHasID>({
  filterMode,
  useResult,
  getCount,
  renderMetadataByline,
  getItems,
}: IItemListOptions<T, E>) {
  const RenderList: React.FC<IItemListProps<T, E>> = ({
    view,
    zoomable,
    otherOperations,
    renderContent,
    renderEditDialog,
    renderDeleteDialog,
    addKeybinds,
  }) => {
    const { filter, setFilter: updateFilter } = useFilter();
    const { effectiveFilter, result, cachedResult } = useQueryResultContext<
      T,
      E
    >();
    const {
      selectedIds,
      getSelected,
      onSelectChange,
      onSelectAll,
      onSelectNone,
    } = useListContext<E>();

    const { modal, showModal, closeModal } = useModal();

    const totalCount = useMemo(() => getCount(cachedResult), [cachedResult]);
    const metadataByline = useMemo(() => {
      if (cachedResult.loading) return "";

      return renderMetadataByline?.(cachedResult) ?? "";
    }, [cachedResult]);

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

    useEffect(() => {
      if (addKeybinds) {
        const unbindExtras = addKeybinds(result, effectiveFilter, selectedIds);
        return () => {
          unbindExtras();
        };
      }
    }, [addKeybinds, result, effectiveFilter, selectedIds]);

    async function onOperationClicked(o: IItemListOperation<T>) {
      await o.onClick(result, effectiveFilter, selectedIds);
      if (o.postRefetch) {
        result.refetch();
      }
    }

    const operations = otherOperations?.map((o) => ({
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

    function onEdit() {
      if (!renderEditDialog) {
        return;
      }

      showModal(
        renderEditDialog(getSelected(), (applied) =>
          onEditDialogClosed(applied)
        )
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

    function onRemoveCriterion(removedCriterion: Criterion<CriterionValue>) {
      updateFilter(
        filter.removeCriterion(removedCriterion.criterionOption.type)
      );
    }

    function onClearAllCriteria() {
      updateFilter(filter.clearCriteria());
    }

    return (
      <div className="item-list-container">
        <FilteredListToolbar
          showEditFilter={showEditFilter}
          view={view}
          operations={operations}
          zoomable={zoomable}
          onEdit={renderEditDialog ? onEdit : undefined}
          onDelete={renderDeleteDialog ? onDelete : undefined}
        />
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

  const ItemList: React.FC<IItemListProps<T, E>> = (props) => {
    const { view, filterHook, selectable, alterQuery = true } = props;

    const [filter, setFilterState] = useState<ListFilterModel>(
      () => new ListFilterModel(filterMode)
    );

    const { defaultFilter, loading: defaultFilterLoading } = useDefaultFilter(
      filterMode,
      view
    );

    // scroll to the top of the page when the page changes
    useScrollToTopOnPageChange(filter.currentPage);

    if (defaultFilterLoading) return null;

    return (
      <FilterContext filter={filter} setFilter={setFilterState}>
        <SetFilterURL defaultFilter={defaultFilter} setURL={alterQuery}>
          <QueryResultContext
            filterHook={filterHook}
            useResult={useResult}
            getCount={getCount}
            getItems={getItems}
          >
            {({ items }) => (
              <ListContext selectable={selectable} items={items}>
                <RenderList {...props} />
              </ListContext>
            )}
          </QueryResultContext>
        </SetFilterURL>
      </FilterContext>
    );
  };

  return ItemList;
}

export const showWhenSelected = <T extends QueryResult>(
  result: T,
  filter: ListFilterModel,
  selectedIds: Set<string>
) => {
  return selectedIds.size > 0;
};
