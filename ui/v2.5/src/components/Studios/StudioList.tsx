import React, { useCallback, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory, useLocation } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindStudios,
  useFindStudios,
  useStudiosDestroy,
} from "src/core/StashService";
import { useFilteredItemList } from "../List/ItemList";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { StudioTagger } from "../Tagger/studios/StudioTagger";
import { StudioCardGrid } from "./StudioCardGrid";
import { View } from "../List/views";
import { EditStudiosDialog } from "./EditStudiosDialog";
import {
  FilteredListToolbar,
  IItemListOperation,
} from "../List/FilteredListToolbar";
import { PatchComponent, PatchContainerComponent } from "src/patch";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import { ListOperations } from "../List/ListOperationButtons";
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "../Shared/Sidebar";
import useFocus from "src/utils/focus";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/FilterSidebar";
import { FilterTags } from "../List/FilterTags";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { LoadedContent } from "../List/PagedList";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { FavoriteStudioCriterionOption } from "src/models/list-filter/criteria/favorite";
import { Button } from "react-bootstrap";
import cx from "classnames";

const StudioList: React.FC<{
  studios: GQL.StudioDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromParent?: boolean;
}> = PatchComponent(
  "StudioList",
  ({ studios, filter, selectedIds, onSelectChange, fromParent }) => {
    if (studios.length === 0) {
      return null;
    }

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <StudioCardGrid
          studios={studios}
          zoomIndex={filter.zoomIndex}
          fromParent={fromParent}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return <h1>TODO</h1>;
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return <h1>TODO</h1>;
    }
    if (filter.displayMode === DisplayMode.Tagger) {
      return <StudioTagger studios={studios} />;
    }

    return null;
  }
);

const StudioFilterSidebarSections = PatchContainerComponent(
  "FilteredStudioList.SidebarSections"
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

      <StudioFilterSidebarSections>
        <SidebarTagsFilter
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarRatingFilter filter={filter} setFilter={setFilter} />
        <SidebarBooleanFilter
          title={<FormattedMessage id="favourite" />}
          filter={filter}
          setFilter={setFilter}
          option={FavoriteStudioCriterionOption}
          sectionID="favourite"
        />
      </StudioFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
    </>
  );
};

interface IStudioList {
  fromParent?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
  extraOperations?: IItemListOperation<GQL.FindStudiosQueryResult>[];
}

function useViewRandom(filter: ListFilterModel, count: number) {
  const history = useHistory();

  const viewRandom = useCallback(async () => {
    // query for a random studio
    if (count === 0) {
      return;
    }

    const index = Math.floor(Math.random() * count);
    const filterCopy = cloneDeep(filter);
    filterCopy.itemsPerPage = 1;
    filterCopy.currentPage = index + 1;
    const singleResult = await queryFindStudios(filterCopy);
    if (singleResult.data.findStudios.studios.length === 1) {
      const { id } = singleResult.data.findStudios.studios[0];
      // navigate to the studio page
      history.push(`/studios/${id}`);
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

export const FilteredStudioList = PatchComponent(
  "FilteredStudioList",
  (props: IStudioList) => {
    const intl = useIntl();
    const history = useHistory();
    const location = useLocation();

    const searchFocus = useFocus();

    const { filterHook, view, alterQuery, extraOperations = [] } = props;

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
          filterMode: GQL.FilterMode.Studios,
          view,
          useURL: alterQuery,
        },
        queryResultProps: {
          useResult: useFindStudios,
          getCount: (r) => r.data?.findStudios.count ?? 0,
          getItems: (r) => r.data?.findStudios.studios ?? [],
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

    function onCreateNew() {
      let queryParam = new URLSearchParams(location.search).get("q");
      let newPath = "/studios/new";
      if (queryParam) {
        newPath += "?q=" + encodeURIComponent(queryParam);
      }
      history.push(newPath);
    }

    const viewRandom = useViewRandom(filter, totalCount);

    function onExport(all: boolean) {
      showModal(
        <ExportDialog
          exportInput={{
            studios: {
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
        <EditStudiosDialog
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
          singularEntity={intl.formatMessage({ id: "studio" })}
          pluralEntity={intl.formatMessage({ id: "studios" })}
          destroyMutation={useStudiosDestroy}
        />
      );
    }

    const convertedExtraOperations = extraOperations.map((op) => ({
      text: op.text,
      onClick: () => op.onClick(result, filter, selectedIds),
      isDisplayed: () => op.isDisplayed?.(result, filter, selectedIds) ?? true,
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
        onCreateNew={onCreateNew}
        entityType={intl.formatMessage({ id: "studio" })}
        operationsMenuClassName="studio-list-operations-dropdown"
      />
    );

    return (
      <div
        className={cx("item-list-container studio-list", {
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
                <StudioList
                  filter={effectiveFilter}
                  studios={items}
                  selectedIds={selectedIds}
                  onSelectChange={onSelectChange}
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
            </SidebarPaneContent>
          </SidebarPane>
        </SidebarStateContext.Provider>
      </div>
    );
  }
);
