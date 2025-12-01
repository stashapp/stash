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
import {
  FilteredListToolbar2,
  ToolbarFilterSection,
  ToolbarSelectionSection,
} from "../List/MyListToolbar";
import { useFilteredItemList } from "../List/ItemList";
import { Sidebar, SidebarPane, useSidebarState } from "../Shared/MySidebar";
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
import { Pagination } from "../List/Pagination";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import {
  createDateCriterionOption,
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import useFocus from "src/utils/myFocus";
import {
  faFilter,
  faPencil,
  faPlus,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { TaggerContext } from "../Tagger/context";
import { SidebarPerformersFilter } from "../List/Filters/PerformersFilter";
import { PerformersCriterionOption } from "src/models/list-filter/criteria/performers";
import { SidebarPathFilter } from "../List/Filters/PathFilter";
import { PathCriterionOption } from "src/models/list-filter/criteria/path";
import { SidebarNumberFilter } from "../List/Filters/NumberFilter";
import { SidebarStringFilter } from "../List/Filters/StringFilter";
import { SidebarDateFilter } from "../List/Filters/DateFilter";
import { SidebarPerformerTagsFilter } from "../List/Filters/PerformerTagsFilter";
import { ListResultsHeader } from "../List/MyListResultsHeader";

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

  const fileCountCriterionOption =
    createMandatoryNumberCriterionOption("file_count");
  const UrlCriterionOption = createStringCriterionOption("url");
  const DateCriterionOption = createDateCriterionOption("date");

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
        <div className="sidebar-filters">
          <div className="sidebar-section-header">
            <Icon icon={faFilter} />
            <FormattedMessage id="filters" />
          </div>
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
          <SidebarPerformerTagsFilter
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
          <SidebarDateFilter
            title={<FormattedMessage id="date" />}
            data-type={DateCriterionOption.type}
            option={DateCriterionOption}
            filter={filter}
            setFilter={setFilter}
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
          <SidebarStringFilter
            title={<FormattedMessage id="url" />}
            data-type={UrlCriterionOption.type}
            option={UrlCriterionOption}
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
        </div>
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

const GalleryListOperations: React.FC<{
  hasSelection: boolean;
  operations: IOperations[];
  onEdit: () => void;
  onDelete: () => void;
  onCreateNew: () => void;
}> = ({ hasSelection, operations, onEdit, onDelete, onCreateNew }) => {
  const intl = useIntl();

  return (
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
  const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
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
            <FilteredListToolbar2
              className="gallery-list-toolbar"
              hasSelection={hasSelection}
              filterSection={
                <ToolbarFilterSection
                  filter={filter}
                  onSetFilter={setFilter}
                  onToggleSidebar={() => setShowSidebar(!showSidebar)}
                  onEditCriterion={(c) =>
                    showEditFilter(c?.criterionOption.type)
                  }
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
                <GalleryListOperations
                  hasSelection={hasSelection}
                  operations={otherOperations}
                  onEdit={onEdit}
                  onDelete={onDelete}
                  onCreateNew={onCreateNew}
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
