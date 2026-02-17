import React, { useCallback, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import Mousetrap from "mousetrap";
import { useHistory } from "react-router-dom";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindGroups,
  useFindGroups,
  useGroupsDestroy,
} from "src/core/StashService";
import { useFilteredItemList } from "../List/ItemList";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { GroupCardGrid } from "./GroupCardGrid";
import { EditGroupsDialog } from "./EditGroupsDialog";
import { View } from "../List/views";
import {
  FilteredListToolbar,
  IItemListOperation,
} from "../List/FilteredListToolbar";
import { PatchComponent, PatchContainerComponent } from "src/patch";
import useFocus from "src/utils/focus";
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "../Shared/Sidebar";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/FilterSidebar";
import {
  IListFilterOperation,
  ListOperations,
} from "../List/ListOperationButtons";
import cx from "classnames";
import { FilterTags } from "../List/FilterTags";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { LoadedContent } from "../List/PagedList";
import { SidebarStudiosFilter } from "../List/Filters/StudiosFilter";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { Button } from "react-bootstrap";

const GroupList: React.FC<{
  groups: GQL.ListGroupDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
}> = PatchComponent(
  "GroupList",
  ({ groups, filter, selectedIds, onSelectChange, fromGroupId, onMove }) => {
    if (groups.length === 0) {
      return null;
    }

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <GroupCardGrid
          groups={groups ?? []}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
          fromGroupId={fromGroupId}
          onMove={onMove}
        />
      );
    }

    return null;
  }
);

const GroupFilterSidebarSections = PatchContainerComponent(
  "FilteredGroupList.SidebarSections"
);

const SidebarContent: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  sidebarOpen: boolean;
  onClose?: () => void;
  showEditFilter: (editingCriterion?: string) => void;
  count?: number;
  focus?: ReturnType<typeof useFocus>;
}> = ({
  filter,
  setFilter,
  filterHook,
  view,
  showEditFilter,
  sidebarOpen,
  onClose,
  count,
  focus,
}) => {
  const showResultsId =
    count !== undefined ? "actions.show_count_results" : "actions.show_results";

  const hideStudios = view === View.StudioScenes;

  return (
    <>
      <FilteredSidebarHeader
        sidebarOpen={sidebarOpen}
        showEditFilter={showEditFilter}
        filter={filter}
        setFilter={setFilter}
        view={view}
        focus={focus}
      />

      <GroupFilterSidebarSections>
        {!hideStudios && (
          <SidebarStudiosFilter
            filter={filter}
            setFilter={setFilter}
            filterHook={filterHook}
          />
        )}
        <SidebarTagsFilter
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarRatingFilter filter={filter} setFilter={setFilter} />
      </GroupFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
    </>
  );
};

interface IGroupListContext {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultFilter?: ListFilterModel;
  view?: View;
  alterQuery?: boolean;
}

interface IGroupList extends IGroupListContext {
  fromGroupId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
  otherOperations?: IItemListOperation<GQL.FindGroupsQueryResult>[];
}

function useViewRandom(filter: ListFilterModel, count: number) {
  const history = useHistory();

  const viewRandom = useCallback(async () => {
    // query for a random scene
    if (count === 0) {
      return;
    }

    const index = Math.floor(Math.random() * count);
    const filterCopy = cloneDeep(filter);
    filterCopy.itemsPerPage = 1;
    filterCopy.currentPage = index + 1;
    const singleResult = await queryFindGroups(filterCopy);
    if (singleResult.data.findGroups.groups.length === 1) {
      const { id } = singleResult.data.findGroups.groups[0];
      // navigate to the image player page
      history.push(`/groups/${id}`);
    }
  }, [history, filter, count]);

  return viewRandom;
}

function useAddKeybinds(filter: ListFilterModel, count: number) {
  const viewRandom = useViewRandom(filter, count);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      viewRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [viewRandom]);
}

export const FilteredGroupList = PatchComponent(
  "FilteredGroupList",
  (props: IGroupList) => {
    const intl = useIntl();

    const searchFocus = useFocus();

    const {
      filterHook,
      view,
      alterQuery,
      onMove,
      fromGroupId,
      otherOperations: providedOperations = [],
      defaultFilter,
    } = props;

    const withSidebar = view !== View.GroupSubGroups;
    const filterable = view !== View.GroupSubGroups;
    const sortable = view !== View.GroupSubGroups;

    // States
    const {
      showSidebar,
      setShowSidebar,
      sectionOpen,
      setSectionOpen,
      loading: sidebarStateLoading,
    } = useSidebarState(view);

    const { filterState, queryResult, modalState, listSelect, showEditFilter } =
      useFilteredItemList({
        filterStateProps: {
          filterMode: GQL.FilterMode.Groups,
          defaultFilter,
          view,
          useURL: alterQuery,
        },
        queryResultProps: {
          useResult: useFindGroups,
          getCount: (r) => r.data?.findGroups.count ?? 0,
          getItems: (r) => r.data?.findGroups.groups ?? [],
          filterHook,
        },
      });

    const { filter, setFilter } = filterState;

    const { effectiveFilter, result, cachedResult, items, totalCount } =
      queryResult;

    const {
      selectedIds,
      selectedItems,
      onSelectChange,
      onSelectAll,
      onSelectNone,
      onInvertSelection,
      hasSelection,
    } = listSelect;

    const { modal, showModal, closeModal } = modalState;

    // Utility hooks
    const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
      filter,
      setFilter,
    });

    useAddKeybinds(filter, totalCount);
    useFilteredSidebarKeybinds({
      showSidebar,
      setShowSidebar,
    });

    useEffect(() => {
      Mousetrap.bind("e", () => {
        if (hasSelection) {
          onEdit?.();
        }
      });

      Mousetrap.bind("d d", () => {
        if (hasSelection) {
          onDelete?.();
        }
      });

      return () => {
        Mousetrap.unbind("e");
        Mousetrap.unbind("d d");
      };
    });

    const onCloseEditDelete = useCloseEditDelete({
      closeModal,
      onSelectNone,
      result,
    });

    const viewRandom = useViewRandom(filter, totalCount);

    function onExport(all: boolean) {
      showModal(
        <ExportDialog
          exportInput={{
            groups: {
              ids: Array.from(selectedIds.values()),
              all: all,
            },
          }}
          onClose={() => closeModal()}
        />
      );
    }

    function onEdit() {
      showModal(
        <EditGroupsDialog
          selected={selectedItems}
          onClose={onCloseEditDelete}
        />
      );
    }

    function onDelete() {
      showModal(
        <DeleteEntityDialog
          selected={selectedItems}
          onClose={onCloseEditDelete}
          singularEntity={intl.formatMessage({ id: "group" })}
          pluralEntity={intl.formatMessage({ id: "groups" })}
          destroyMutation={useGroupsDestroy}
        />
      );
    }

    const convertedExtraOperations: IListFilterOperation[] =
      providedOperations.map((o) => ({
        ...o,
        isDisplayed: o.isDisplayed
          ? () => o.isDisplayed!(result, filter, selectedIds)
          : undefined,
        onClick: () => {
          o.onClick(result, filter, selectedIds);
        },
      }));

    const otherOperations = [
      ...convertedExtraOperations,
      {
        text: intl.formatMessage({ id: "actions.select_all" }),
        onClick: () => onSelectAll(),
        isDisplayed: () => totalCount > 0,
      },
      {
        text: intl.formatMessage({ id: "actions.select_none" }),
        onClick: () => onSelectNone(),
        isDisplayed: () => hasSelection,
      },
      {
        text: intl.formatMessage({ id: "actions.invert_selection" }),
        onClick: () => onInvertSelection(),
        isDisplayed: () => totalCount > 0,
      },
      {
        text: intl.formatMessage({ id: "actions.view_random" }),
        onClick: viewRandom,
      },
      {
        text: intl.formatMessage({ id: "actions.export" }),
        onClick: () => onExport(false),
        isDisplayed: () => hasSelection,
      },
      {
        text: intl.formatMessage({ id: "actions.export_all" }),
        onClick: () => onExport(true),
      },
    ];

    // render
    if (sidebarStateLoading) return null;

    const operations = (
      <ListOperations
        items={items.length}
        hasSelection={hasSelection}
        operations={otherOperations}
        onEdit={onEdit}
        onDelete={onDelete}
        operationsMenuClassName="group-list-operations-dropdown"
      />
    );

    const content = (
      <>
        <FilteredListToolbar
          filter={filter}
          listSelect={listSelect}
          setFilter={setFilter}
          showEditFilter={showEditFilter}
          onDelete={onDelete}
          onEdit={onEdit}
          operationComponent={operations}
          view={view}
          zoomable
          filterable={filterable}
          sortable={sortable}
        />

        <FilterTags
          criteria={filter.criteria}
          onEditCriterion={(c) => showEditFilter(c.criterionOption.type)}
          onRemoveCriterion={removeCriterion}
          onRemoveAll={clearAllCriteria}
        />

        <div className="pagination-index-container">
          <Pagination
            currentPage={filter.currentPage}
            itemsPerPage={filter.itemsPerPage}
            totalItems={totalCount}
            onChangePage={(page) => setFilter(filter.changePage(page))}
          />
          <PaginationIndex
            loading={cachedResult.loading}
            itemsPerPage={filter.itemsPerPage}
            currentPage={filter.currentPage}
            totalItems={totalCount}
          />
        </div>

        <LoadedContent loading={result.loading} error={result.error}>
          <GroupList
            filter={effectiveFilter}
            groups={items}
            selectedIds={selectedIds}
            onSelectChange={onSelectChange}
            fromGroupId={fromGroupId}
            onMove={onMove}
          />
        </LoadedContent>

        {totalCount > filter.itemsPerPage && (
          <div className="pagination-footer-container">
            <div className="pagination-footer">
              <Pagination
                itemsPerPage={filter.itemsPerPage}
                currentPage={filter.currentPage}
                totalItems={totalCount}
                onChangePage={setPage}
                pagePopupPlacement="top"
              />
            </div>
          </div>
        )}
      </>
    );

    if (!withSidebar) {
      return content;
    }

    return (
      <div
        className={cx("item-list-container group-list", {
          "hide-sidebar": !showSidebar,
        })}
      >
        {modal}

        <SidebarStateContext.Provider value={{ sectionOpen, setSectionOpen }}>
          <SidebarPane hideSidebar={!showSidebar}>
            <Sidebar hide={!showSidebar} onHide={() => setShowSidebar(false)}>
              <SidebarContent
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                showEditFilter={showEditFilter}
                view={view}
                sidebarOpen={showSidebar}
                onClose={() => setShowSidebar(false)}
                count={cachedResult.loading ? undefined : totalCount}
                focus={searchFocus}
              />
            </Sidebar>
            <SidebarPaneContent
              onSidebarToggle={() => setShowSidebar(!showSidebar)}
            >
              {content}
            </SidebarPaneContent>
          </SidebarPane>
        </SidebarStateContext.Provider>
      </div>
    );
  }
);
