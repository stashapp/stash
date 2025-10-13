import React, { useCallback, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { useFilteredItemList } from "../List/ItemList";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { queryFindGalleries, useFindGalleries } from "src/core/StashService";
import GalleryWallCard from "./GalleryWallCard";
import { EditGalleriesDialog } from "./EditGalleriesDialog";
import { DeleteGalleriesDialog } from "./DeleteGalleriesDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { GalleryListTable } from "./GalleryListTable";
import { GalleryCardGrid } from "./GalleryGridCard";
import { View } from "../List/views";
import useFocus from "src/utils/focus";
import { Sidebar, SidebarPane, useSidebarState } from "../Shared/Sidebar";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/FilterSidebar";
import cx from "classnames";
import {
  FilteredListToolbar2,
  ToolbarFilterSection,
  ToolbarSelectionSection,
} from "../List/ListToolbar";
import { ListResultsHeader } from "../List/ListResultsHeader";
import { LoadedContent } from "../List/PagedList";
import { Pagination } from "../List/Pagination";
import { PatchContainerComponent } from "src/patch";
import { SidebarStudiosFilter } from "../List/Filters/StudiosFilter";
import { SidebarPerformersFilter } from "../List/Filters/PerformersFilter";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { OrganizedCriterionOption } from "src/models/list-filter/criteria/organized";
import { Button } from "react-bootstrap";
import { ListOperations } from "../List/ListOperationButtons";

const GalleryList: React.FC<{
  galleries: GQL.SlimGalleryDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}> = ({ galleries, filter, selectedIds, onSelectChange }) => {
  if (galleries.length === 0) {
    return null;
  }

  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <GalleryCardGrid
        galleries={galleries}
        selectedIds={selectedIds}
        zoomIndex={filter.zoomIndex}
        onSelectChange={onSelectChange}
      />
    );
  }
  if (filter.displayMode === DisplayMode.List) {
    return (
      <GalleryListTable
        galleries={galleries}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
      />
    );
  }
  if (filter.displayMode === DisplayMode.Wall) {
    return (
      <div className="row">
        <div className={`GalleryWall zoom-${filter.zoomIndex}`}>
          {galleries.map((gallery) => (
            <GalleryWallCard key={gallery.id} gallery={gallery} />
          ))}
        </div>
      </div>
    );
  }

  return null;
};

const GalleryFilterSidebarSections = PatchContainerComponent(
  "FilteredGalleryList.SidebarSections"
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

      <GalleryFilterSidebarSections>
        {!hideStudios && (
          <SidebarStudiosFilter
            filter={filter}
            setFilter={setFilter}
            filterHook={filterHook}
          />
        )}
        <SidebarPerformersFilter
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarTagsFilter
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarRatingFilter filter={filter} setFilter={setFilter} />
        <SidebarBooleanFilter
          title={<FormattedMessage id="organized" />}
          data-type={OrganizedCriterionOption.type}
          option={OrganizedCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
      </GalleryFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
    </>
  );
};

interface IGalleryList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
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
    const singleResult = await queryFindGalleries(filterCopy);
    if (singleResult.data.findGalleries.galleries.length === 1) {
      const { id } = singleResult.data.findGalleries.galleries[0];
      // navigate to the image player page
      history.push(`/galleries/${id}`);
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

export const FilteredGalleryList = (props: IGalleryList) => {
  const intl = useIntl();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, view, alterQuery } = props;

  // States
  const {
    showSidebar,
    setShowSidebar,
    loading: sidebarStateLoading,
  } = useSidebarState(view);

  const { filterState, queryResult, modalState, listSelect, showEditFilter } =
    useFilteredItemList({
      filterStateProps: {
        filterMode: GQL.FilterMode.Galleries,
        view,
        useURL: alterQuery,
      },
      queryResultProps: {
        useResult: useFindGalleries,
        getCount: (r) => r.data?.findGalleries.count ?? 0,
        getItems: (r) => r.data?.findGalleries.galleries ?? [],
        filterHook,
      },
    });

  const { filter, setFilter, loading: filterLoading } = filterState;

  const { effectiveFilter, result, cachedResult, items, totalCount } =
    queryResult;

  const {
    selectedIds,
    selectedItems,
    onSelectChange,
    onSelectAll,
    onSelectNone,
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
          galleries: {
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
      <EditGalleriesDialog
        selected={selectedItems}
        onClose={onCloseEditDelete}
      />
    );
  }

  function onDelete() {
    showModal(
      <DeleteGalleriesDialog
        selected={selectedItems}
        onClose={onCloseEditDelete}
      />
    );
  }

  const otherOperations = [
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
  if (filterLoading || sidebarStateLoading) return null;

  return (
    <div
      className={cx("item-list-container gallery-list", {
        "hide-sidebar": !showSidebar,
      })}
    >
      {modal}

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
        <div>
          <FilteredListToolbar2
            className="gallery-list-toolbar"
            hasSelection={hasSelection}
            filterSection={
              <ToolbarFilterSection
                filter={filter}
                onSetFilter={setFilter}
                onToggleSidebar={() => setShowSidebar(!showSidebar)}
                onEditCriterion={(c) => showEditFilter(c?.criterionOption.type)}
                onRemoveCriterion={removeCriterion}
                onRemoveAllCriterion={() => clearAllCriteria(true)}
                onEditSearchTerm={() => {
                  setShowSidebar(true);
                  setSearchFocus(true);
                }}
                onRemoveSearchTerm={() => setFilter(filter.clearSearchTerm())}
              />
            }
            selectionSection={
              <ToolbarSelectionSection
                selected={selectedIds.size}
                onToggleSidebar={() => setShowSidebar(!showSidebar)}
                onSelectAll={() => onSelectAll()}
                onSelectNone={() => onSelectNone()}
              />
            }
            operationSection={
              <ListOperations
                items={items.length}
                hasSelection={hasSelection}
                operations={otherOperations}
                onEdit={onEdit}
                onDelete={onDelete}
                entityType="gallery"
              />
            }
          />

          <ListResultsHeader
            loading={cachedResult.loading}
            filter={filter}
            totalCount={totalCount}
            onChangeFilter={(newFilter) => setFilter(newFilter)}
          />

          <LoadedContent loading={result.loading} error={result.error}>
            <GalleryList
              filter={effectiveFilter}
              galleries={items}
              selectedIds={selectedIds}
              onSelectChange={onSelectChange}
            />
          </LoadedContent>

          {totalCount > filter.itemsPerPage && (
            <div className="pagination-footer">
              <Pagination
                itemsPerPage={filter.itemsPerPage}
                currentPage={filter.currentPage}
                totalItems={totalCount}
                onChangePage={setPage}
                pagePopupPlacement="top"
              />
            </div>
          )}
        </div>
      </SidebarPane>
    </div>
  );
};
