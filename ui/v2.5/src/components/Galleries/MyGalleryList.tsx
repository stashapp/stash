import React, { useCallback, useEffect } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindGalleries, useFindGalleries } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import GalleryWallCard from "./GalleryWallCard";
import { EditGalleriesDialog } from "./EditGalleriesDialog";
import { DeleteGalleriesDialog } from "./DeleteGalleriesDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { GalleryListTable } from "./GalleryListTable";
import { GalleryCardGrid } from "./GalleryGridCard";
import { View } from "../List/views";
import { LoadedContent } from "../List/PagedList";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import {
  OperationDropdown,
  OperationDropdownItem,
} from "../List/ListOperationButtons";
import { useFilteredItemList } from "../List/ItemList";
import { FilterTags } from "../List/MyFilterTags";
import { Sidebar, SidebarPane, useSidebarState } from "../Shared/Sidebar";
import { SidebarStudiosFilter } from "../List/Filters/StudiosFilter";
import { StudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "src/models/list-filter/criteria/tags";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import cx from "classnames";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { OrganizedCriterionOption } from "src/models/list-filter/criteria/organized";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/MyFilterSidebar";
import { PatchContainerComponent } from "src/patch";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { Button, ButtonGroup, ButtonToolbar } from "react-bootstrap";
import { FilterButton } from "../List/Filters/FilterButton";
import { Icon } from "../Shared/Icon";
import { ListViewOptions } from "../List/ListViewOptions";
import { PageSizeSelector, SortBySelect } from "../List/ListFilter";
import { createMandatoryNumberCriterionOption, Criterion } from "src/models/list-filter/criteria/criterion";
import useFocus from "src/utils/myFocus";
import {
  faPencil,
  faPlus,
  faSliders,
  faTimes,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { TaggerContext } from "../Tagger/context";
import { SidebarPerformersFilter } from "../List/Filters/PerformersFilter";
import { PerformersCriterionOption } from "src/models/list-filter/criteria/performers";
import { SidebarPathFilter } from "../List/Filters/PathFilter";
import { PathCriterionOption } from "src/models/list-filter/criteria/path";
import { SidebarNumberFilter } from "../List/Filters/NumberFilter";

function getItems(result: GQL.FindGalleriesQueryResult) {
  return result?.data?.findGalleries?.galleries ?? [];
}

function getCount(result: GQL.FindGalleriesQueryResult) {
  return result?.data?.findGalleries?.count ?? 0;
}

function useViewRandom(
  result: GQL.FindGalleriesQueryResult,
  filter: ListFilterModel
) {
  const history = useHistory();

  const viewRandom = useCallback(async () => {
    // query for a random image
    if (result.data?.findGalleries) {
      const { count } = result.data.findGalleries;

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
    }
  }, [result, filter, history]);

  return viewRandom;
}

function useAddKeybinds(
  result: GQL.FindGalleriesQueryResult,
  filter: ListFilterModel
) {
  const viewRandom = useViewRandom(result, filter);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      viewRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [viewRandom]);
}

const GalleryListContent: React.FC<{
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
        <div className="GalleryWall">
          {galleries.map((gallery) => (
            <GalleryWallCard key={gallery.id} gallery={gallery} />
          ))}
        </div>
      </div>
    );
  }

  return null;
};

export const MyGalleriesFilterSidebarSections = PatchContainerComponent(
  "MyFilteredGalleryList.SidebarSections"
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
  
  const fileCountCriterionOption = createMandatoryNumberCriterionOption("file_count");

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

      <MyGalleriesFilterSidebarSections>
        <SidebarStudiosFilter
          title={<FormattedMessage id="studios" />}
          data-type={StudiosCriterionOption.type}
          option={StudiosCriterionOption}
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarPerformersFilter
          title={<FormattedMessage id="performers" />}
          data-type={PerformersCriterionOption.type}
          option={PerformersCriterionOption}
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarTagsFilter
          title={<FormattedMessage id="performer_tags" />}
          data-type={PerformerTagsCriterionOption.type}
          option={PerformerTagsCriterionOption}
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarTagsFilter
          title={<FormattedMessage id="tags" />}
          data-type={TagsCriterionOption.type}
          option={TagsCriterionOption}
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarPathFilter
          title={<FormattedMessage id="path" />}
          data-type={PathCriterionOption.type}
          option={PathCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarNumberFilter
          title={<FormattedMessage id="file_count" />}
          data-type={fileCountCriterionOption.type}
          option={fileCountCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarRatingFilter
          title={<FormattedMessage id="rating" />}
          data-type={RatingCriterionOption.type}
          option={RatingCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarBooleanFilter
          title={<FormattedMessage id="organized" />}
          data-type={OrganizedCriterionOption.type}
          option={OrganizedCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
      </MyGalleriesFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
    </>
  );
};

interface IOperations {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
  className?: string;
}

const ListToolbarContent: React.FC<{
  searchTerm: string;
  criteria: Criterion[];
  items: GQL.SlimGalleryDataFragment[];
  selectedIds: Set<string>;
  operations: IOperations[];
  onToggleSidebar: () => void;
  onEditCriterion: (c?: Criterion) => void;
  onRemoveCriterion: (criterion: Criterion, valueIndex?: number) => void;
  onRemoveAllCriterion: () => void;
  onEditSearchTerm: () => void;
  onRemoveSearchTerm: () => void;
  onSelectAll: () => void;
  onSelectNone: () => void;
  onEdit: () => void;
  onDelete: () => void;
  onCreateNew: () => void;
}> = ({
  searchTerm,
  criteria,
  items,
  selectedIds,
  operations,
  onToggleSidebar,
  onEditCriterion,
  onRemoveCriterion,
  onRemoveAllCriterion,
  onEditSearchTerm,
  onRemoveSearchTerm,
  onSelectAll,
  onSelectNone,
  onEdit,
  onDelete,
  onCreateNew,
}) => {
  const intl = useIntl();

  const hasSelection = selectedIds.size > 0;

  const sidebarToggle = (
    <Button
      className="sidebar-toggle-button ignore-sidebar-outside-click"
      variant="secondary"
      onClick={() => onToggleSidebar()}
      title={intl.formatMessage({ id: "actions.sidebar.toggle" })}
    >
      <Icon icon={faSliders} />
    </Button>
  );

  return (
    <>
      {!hasSelection && (
        <div className="my-filter-toolbar">
          {sidebarToggle}
          <FilterTags
            searchTerm={searchTerm}
            criteria={criteria}
            onEditCriterion={onEditCriterion}
            onRemoveCriterion={onRemoveCriterion}
            onRemoveAll={onRemoveAllCriterion}
            onEditSearchTerm={onEditSearchTerm}
            onRemoveSearchTerm={onRemoveSearchTerm}
            truncateOnOverflow
          />
          <FilterButton
            onClick={() => onEditCriterion()}
            count={criteria.length}
            title={intl.formatMessage({ id: "actions.sidebar.toggle" })}
          />
        </div>
      )}
      {hasSelection && (
        <div className="selected-items-info">
          <Button
            variant="secondary"
            className="minimal"
            onClick={() => onSelectNone()}
            title={intl.formatMessage({ id: "actions.select_none" })}
          >
            <Icon icon={faTimes} />
          </Button>
          <span>{selectedIds.size} selected</span>
          <Button variant="link" onClick={() => onSelectAll()}>
            <FormattedMessage id="actions.select_all" />
          </Button>
        </div>
      )}
      <div>
        <ButtonGroup>
          {!hasSelection && (
            <Button
              className="create-new-button"
              variant="secondary"
              onClick={() => onCreateNew()}
              title={intl.formatMessage(
                { id: "actions.create_entity" },
                { entityType: intl.formatMessage({ id: "gallery" }) }
              )}
            >
              <Icon icon={faPlus} />
            </Button>
          )}

          {hasSelection && (
            <>
              <Button variant="secondary" onClick={() => onEdit()}>
                <Icon icon={faPencil} />
              </Button>
              <Button
                variant="danger"
                className="btn-danger-minimal"
                onClick={() => onDelete()}
              >
                <Icon icon={faTrash} />
              </Button>
            </>
          )}

          <OperationDropdown className="gallery-list-operations">
            {operations.map((o) => {
              if (o.isDisplayed && !o.isDisplayed()) {
                return null;
              }

              return (
                <OperationDropdownItem
                  key={o.text}
                  onClick={o.onClick}
                  text={o.text}
                  className={o.className}
                />
              );
            })}
          </OperationDropdown>
        </ButtonGroup>
      </div>
    </>
  );
};

const ListResultsHeader: React.FC<{
  loading: boolean;
  filter: ListFilterModel;
  totalCount: number;
  onChangeFilter: (filter: ListFilterModel) => void;
}> = ({ loading, filter, totalCount, onChangeFilter }) => {
  return (
    <ButtonToolbar className="gallery-list-header">
      <div>
        <PaginationIndex
          loading={loading}
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
        />
      </div>
      <div>
        <SortBySelect
          options={filter.options.sortByOptions}
          sortBy={filter.sortBy}
          sortDirection={filter.sortDirection}
          onChangeSortBy={(s) =>
            onChangeFilter(filter.setSortBy(s ?? undefined))
          }
          onChangeSortDirection={() =>
            onChangeFilter(filter.toggleSortDirection())
          }
          onReshuffleRandomSort={() =>
            onChangeFilter(filter.reshuffleRandomSort())
          }
        />
        <PageSizeSelector
          pageSize={filter.itemsPerPage}
          setPageSize={(s) => onChangeFilter(filter.setPageSize(s))}
        />
        <ListViewOptions
          displayMode={filter.displayMode}
          zoomIndex={filter.zoomIndex}
          displayModeOptions={filter.options.displayModeOptions}
          onSetDisplayMode={(mode) =>
            onChangeFilter(filter.setDisplayMode(mode))
          }
          onSetZoom={(zoom) => onChangeFilter(filter.setZoom(zoom))}
        />
      </div>
    </ButtonToolbar>
  );
};

interface IFilteredGalleries {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
}

export const MyFilteredGalleryList = (props: IFilteredGalleries) => {
  const intl = useIntl();
  const history = useHistory();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, defaultSort, view, alterQuery } = props;

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
        defaultSort,
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
  const { setPage, removeCriterion, myClearAllCriteria } = useFilterOperations({
    filter,
    setFilter,
  });

  useAddKeybinds(result, filter);
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

  const viewRandom = useViewRandom(result, filter);

  function onCreateNew() {
    history.push("/galleries/new");
  }

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
      isDisplayed: () => totalCount > 1,
    },
    {
      text: intl.formatMessage(
        { id: "actions.create_entity" },
        { entityType: intl.formatMessage({ id: "gallery" }) }
      ),
      onClick: () => onCreateNew(),
      isDisplayed: () => !hasSelection,
      className: "create-new-item",
    },
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
    <TaggerContext>
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
            <ButtonToolbar
              className={cx("gallery-list-toolbar", {
                "has-selection": hasSelection,
              })}
            >
              <ListToolbarContent
                searchTerm={filter.searchTerm}
                criteria={filter.criteria}
                items={items}
                selectedIds={selectedIds}
                operations={otherOperations}
                onToggleSidebar={() => setShowSidebar(!showSidebar)}
                onEditCriterion={(c) => showEditFilter(c?.criterionOption.type)}
                onRemoveCriterion={removeCriterion}
                onRemoveAllCriterion={() => myClearAllCriteria(true)}
                onEditSearchTerm={() => {
                  setShowSidebar(true);
                  setSearchFocus(true);
                }}
                onRemoveSearchTerm={() => setFilter(filter.clearSearchTerm())}
                onSelectAll={() => onSelectAll()}
                onSelectNone={() => onSelectNone()}
                onEdit={onEdit}
                onDelete={onDelete}
                onCreateNew={onCreateNew}
              />
            </ButtonToolbar>

            <ListResultsHeader
              loading={cachedResult.loading}
              filter={filter}
              totalCount={totalCount}
              onChangeFilter={(newFilter) => setFilter(newFilter)}
            />

            <LoadedContent loading={result.loading} error={result.error}>
              <GalleryListContent
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
    </TaggerContext>
  );
};

// Keep the old component for backward compatibility
export const GalleryList: React.FC<IFilteredGalleries> = (props) => {
  return <MyFilteredGalleryList {...props} />;
};

export default MyFilteredGalleryList;
