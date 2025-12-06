import React, { useCallback, useEffect, useMemo, useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindGalleries, useFindGalleries } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import GalleryWallCard from "src/components/Galleries/GalleryWallCard";
import { EditGalleriesDialog } from "src/components/Galleries/EditGalleriesDialog";
import { DeleteGalleriesDialog } from "src/components/Galleries/DeleteGalleriesDialog";
import { ExportDialog } from "src/components/Shared/ExportDialog";
import { GalleryListTable } from "src/components/Galleries/GalleryListTable";
import { GalleryCardGrid } from "src/components/Galleries/GalleryGridCard";
import { View } from "src/components/List/views";
import { LoadedContent } from "src/components/List/PagedList";
import { useCloseEditDelete, useFilterOperations } from "src/components/List/util";
import {
  OperationDropdown,
  OperationDropdownItem,
} from "src/components/List/ListOperationButtons";
import {
  FilteredListToolbar2,
  ToolbarFilterSection,
  ToolbarSelectionSection,
} from "src/extensions/ui";
import { useFilteredItemList } from "src/components/List/ItemList";
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "src/components/Shared/Sidebar";
import {
  SidebarStudiosFilter,
  SidebarTagsFilter,
  SidebarRatingFilter,
  SidebarBooleanFilter,
  SidebarPerformersFilter,
  SidebarPathFilter,
  SidebarNumberFilter,
  SidebarStringFilter,
  SidebarDateFilter,
  SidebarPerformerTagsFilter,
  SidebarFilterSelector,
  FilterWrapper,
  SidebarIsMissingFilter,
  SidebarResolutionFilter,
} from "src/extensions/filters";
import { StudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "src/models/list-filter/criteria/tags";
import cx from "classnames";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { OrganizedCriterionOption } from "src/models/list-filter/criteria/organized";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "src/extensions/ui";
import { PatchContainerComponent } from "src/patch";
import { Pagination } from "src/components/List/Pagination";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import {
  createDateCriterionOption,
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { useFocus } from "src/extensions/hooks";
import {
  faFilter,
  faPencil,
  faPlus,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { TaggerContext } from "src/components/Tagger/context";
import { PerformersCriterionOption } from "src/models/list-filter/criteria/performers";
import { PathCriterionOption } from "src/models/list-filter/criteria/path";
import { ListResultsHeader } from "src/extensions/ui";
import { SidebarFilterDefinition } from "src/hooks/useSidebarFilters";
import { createMandatoryTimestampCriterionOption } from "src/models/list-filter/criteria/criterion";
import { GalleryIsMissingCriterionOption } from "src/models/list-filter/criteria/is-missing";
import {
  useGalleryFacetCounts,
  FacetCountsContext,
} from "src/hooks/useFacetCounts";
import { PerformerFavoriteCriterionOption } from "src/models/list-filter/criteria/favorite";
import { AverageResolutionCriterionOption } from "src/models/list-filter/criteria/resolution";
import { HasChaptersCriterionOption } from "src/models/list-filter/criteria/has-chapters";

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

// Define available filters for galleries sidebar
const galleryFilterDefinitions: SidebarFilterDefinition[] = [
  // Tier 1: Primary filters (visible by default)
  { id: "rating", messageId: "rating", defaultVisible: true },
  { id: "date", messageId: "date", defaultVisible: true },
  { id: "tags", messageId: "tags", defaultVisible: true },
  { id: "performers", messageId: "performers", defaultVisible: true },
  { id: "studios", messageId: "studios", defaultVisible: true },
  { id: "organized", messageId: "organized", defaultVisible: true },

  // Tier 2: Discovery
  { id: "performer_tags", messageId: "performer_tags", defaultVisible: false },
  { id: "performer_favorite", messageId: "performer_favorite", defaultVisible: false },
  { id: "average_resolution", messageId: "resolution", defaultVisible: false },
  { id: "has_chapters", messageId: "hasChapters", defaultVisible: false },

  // Tier 3: Library stats
  { id: "image_count", messageId: "image_count", defaultVisible: false },
  { id: "tag_count", messageId: "tag_count", defaultVisible: false },
  { id: "performer_count", messageId: "performer_count", defaultVisible: false },
  { id: "performer_age", messageId: "performer_age", defaultVisible: false },

  // Tier 4: Metadata
  { id: "title", messageId: "title", defaultVisible: false },
  { id: "code", messageId: "scene_code", defaultVisible: false },
  { id: "details", messageId: "details", defaultVisible: false },
  { id: "photographer", messageId: "photographer", defaultVisible: false },
  { id: "url", messageId: "url", defaultVisible: false },
  { id: "path", messageId: "path", defaultVisible: false },
  { id: "file_count", messageId: "zip_file_count", defaultVisible: false },

  // Tier 5: Technical
  { id: "checksum", messageId: "media_info.checksum", defaultVisible: false },

  // Tier 6: System
  { id: "is_missing", messageId: "isMissing", defaultVisible: false },
  { id: "created_at", messageId: "created_at", defaultVisible: false },
  { id: "updated_at", messageId: "updated_at", defaultVisible: false },
];

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
  clearAllCriteria: () => void;
  onFilterEditModeChange?: (isEditMode: boolean) => void;
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
  clearAllCriteria,
  onFilterEditModeChange,
}) => {
  const showResultsId =
    count !== undefined ? "actions.show_count_results" : "actions.show_results";

  // Criterion options
  const fileCountCriterionOption = createMandatoryNumberCriterionOption("file_count", "zip_file_count");
  const imageCountCriterionOption = createMandatoryNumberCriterionOption("image_count");
  const tagCountCriterionOption = createMandatoryNumberCriterionOption("tag_count");
  const performerCountCriterionOption = createMandatoryNumberCriterionOption("performer_count");
  const performerAgeCriterionOption = createMandatoryNumberCriterionOption("performer_age");
  const UrlCriterionOption = createStringCriterionOption("url");
  const TitleCriterionOption = createStringCriterionOption("title");
  const CodeCriterionOption = createStringCriterionOption("code", "scene_code");
  const DetailsCriterionOption = createStringCriterionOption("details");
  const PhotographerCriterionOption = createStringCriterionOption("photographer");
  const ChecksumCriterionOption = createStringCriterionOption("checksum", "media_info.checksum");
  const DateCriterionOption = createDateCriterionOption("date");
  const CreatedAtCriterionOption = createMandatoryTimestampCriterionOption("created_at");
  const UpdatedAtCriterionOption = createMandatoryTimestampCriterionOption("updated_at");

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
          <SidebarFilterSelector
            viewName="galleries"
            filterDefinitions={galleryFilterDefinitions}
            headerContent={
              <>
                <Icon icon={faFilter} />
                <FormattedMessage id="filters" />
              </>
            }
            onEditModeChange={onFilterEditModeChange}
          >
            {/* Tier 1: Primary Filters */}
            <FilterWrapper filterId="rating">
              <SidebarRatingFilter
                title={<FormattedMessage id="rating" />}
                option={RatingCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="rating"
              />
            </FilterWrapper>
            <FilterWrapper filterId="date">
              <SidebarDateFilter
                title={<FormattedMessage id="date" />}
                option={DateCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="date"
              />
            </FilterWrapper>
            <FilterWrapper filterId="tags">
              <SidebarTagsFilter
                title={<FormattedMessage id="tags" />}
                option={TagsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="tags"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performers">
              <SidebarPerformersFilter
                title={<FormattedMessage id="performers" />}
                option={PerformersCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="performers"
              />
            </FilterWrapper>
            <FilterWrapper filterId="studios">
              <SidebarStudiosFilter
                title={<FormattedMessage id="studios" />}
                option={StudiosCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="studios"
              />
            </FilterWrapper>
            <FilterWrapper filterId="organized">
              <SidebarBooleanFilter
                title={<FormattedMessage id="organized" />}
                option={OrganizedCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="organized"
              />
            </FilterWrapper>

            {/* Tier 2: Discovery */}
            <FilterWrapper filterId="performer_tags">
              <SidebarPerformerTagsFilter
                title={<FormattedMessage id="performer_tags" />}
                option={PerformerTagsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="performer_tags"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performer_favorite">
              <SidebarBooleanFilter
                title={<FormattedMessage id="performer_favorite" />}
                option={PerformerFavoriteCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performer_favorite"
              />
            </FilterWrapper>
            <FilterWrapper filterId="average_resolution">
              <SidebarResolutionFilter
                title={<FormattedMessage id="resolution" />}
                option={AverageResolutionCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="average_resolution"
              />
            </FilterWrapper>
            <FilterWrapper filterId="has_chapters">
              <SidebarBooleanFilter
                title={<FormattedMessage id="hasChapters" />}
                option={HasChaptersCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="has_chapters"
              />
            </FilterWrapper>

            {/* Tier 3: Library Stats */}
            <FilterWrapper filterId="image_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="image_count" />}
                option={imageCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="image_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="tag_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="tag_count" />}
                option={tagCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="tag_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performer_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="performer_count" />}
                option={performerCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performer_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performer_age">
              <SidebarNumberFilter
                title={<FormattedMessage id="performer_age" />}
                option={performerAgeCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performer_age"
              />
            </FilterWrapper>

            {/* Tier 4: Metadata */}
            <FilterWrapper filterId="title">
              <SidebarStringFilter
                title={<FormattedMessage id="title" />}
                option={TitleCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="title"
              />
            </FilterWrapper>
            <FilterWrapper filterId="code">
              <SidebarStringFilter
                title={<FormattedMessage id="scene_code" />}
                option={CodeCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="code"
              />
            </FilterWrapper>
            <FilterWrapper filterId="details">
              <SidebarStringFilter
                title={<FormattedMessage id="details" />}
                option={DetailsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="details"
              />
            </FilterWrapper>
            <FilterWrapper filterId="photographer">
              <SidebarStringFilter
                title={<FormattedMessage id="photographer" />}
                option={PhotographerCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="photographer"
              />
            </FilterWrapper>
            <FilterWrapper filterId="url">
              <SidebarStringFilter
                title={<FormattedMessage id="url" />}
                option={UrlCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="url"
              />
            </FilterWrapper>
            <FilterWrapper filterId="path">
              <SidebarPathFilter
                title={<FormattedMessage id="path" />}
                option={PathCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="path"
              />
            </FilterWrapper>
            <FilterWrapper filterId="file_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="zip_file_count" />}
                option={fileCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="file_count"
              />
            </FilterWrapper>

            {/* Tier 5: Technical */}
            <FilterWrapper filterId="checksum">
              <SidebarStringFilter
                title={<FormattedMessage id="media_info.checksum" />}
                option={ChecksumCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="checksum"
              />
            </FilterWrapper>

            {/* Tier 6: System */}
            <FilterWrapper filterId="is_missing">
              <SidebarIsMissingFilter
                title={<FormattedMessage id="isMissing" />}
                option={GalleryIsMissingCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="is_missing"
              />
            </FilterWrapper>
            <FilterWrapper filterId="created_at">
              <SidebarDateFilter
                title={<FormattedMessage id="created_at" />}
                option={CreatedAtCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="created_at"
                isTime
              />
            </FilterWrapper>
            <FilterWrapper filterId="updated_at">
              <SidebarDateFilter
                title={<FormattedMessage id="updated_at" />}
                option={UpdatedAtCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="updated_at"
                isTime
              />
            </FilterWrapper>
          </SidebarFilterSelector>
        </div>
      </MyGalleriesFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
      <div className="clear-all-filters">
        <Button
          className="clear-all-filters-button"
          variant="secondary"
          onClick={() => clearAllCriteria()}
          // TODO: add message
          title="Clear All Filters"
        >
          <FormattedMessage id="Clear All Filters" />
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
  items: number;
  hasSelection: boolean;
  operations: IOperations[];
  onEdit: () => void;
  onDelete: () => void;
  onCreateNew: () => void;
}> = ({ items, hasSelection, operations, onEdit, onDelete, onCreateNew }) => {
  const intl = useIntl();

  return (
    <div className="list-operations">
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

        <OperationDropdown
          className="list-operations"
          menuPortalTarget={document.body}
        >
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
    sectionOpen: baseSectionOpen,
    setSectionOpen: baseSetSectionOpen,
  } = useSidebarState(view);

  // Track filter customization edit mode
  const [isFilterEditMode, setIsFilterEditMode] = useState(false);

  const sectionOpen = useMemo(() => {
    if (isFilterEditMode) {
      const closedSections: Record<string, boolean> = {};
      Object.keys(baseSectionOpen).forEach((key) => {
        closedSections[key] = false;
      });
      return closedSections;
    }
    return baseSectionOpen;
  }, [isFilterEditMode, baseSectionOpen]);

  const setSectionOpen = useCallback(
    (section: string, open: boolean) => {
      if (isFilterEditMode) return;
      baseSetSectionOpen(section, open);
    },
    [isFilterEditMode, baseSetSectionOpen]
  );

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

  // Fetch facet counts for sidebar filters
  const { counts: facetCounts, loading: facetLoading } = useGalleryFacetCounts(filter, {
    isOpen: showSidebar ?? false,
    debounceMs: 300,
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

  const operations = (
    <GalleryListOperations
      items={items.length}
      hasSelection={hasSelection}
      operations={otherOperations}
      onEdit={onEdit}
      onDelete={onDelete}
      onCreateNew={onCreateNew}
    />
  );

  return (
    <TaggerContext>
      <div
        className={cx("item-list-container gallery-list", {
          "hide-sidebar": !showSidebar,
        })}
      >
        {modal}

        <SidebarStateContext.Provider value={{ sectionOpen, setSectionOpen, disabled: isFilterEditMode }}>
          <FacetCountsContext.Provider value={{ counts: facetCounts, loading: facetLoading }}>
          <SidebarPane hideSidebar={!showSidebar}>
            <Sidebar hide={!showSidebar} onHide={() => setShowSidebar(false)}>
              <SidebarContent
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                showEditFilter={showEditFilter}
                view={view}
                onFilterEditModeChange={setIsFilterEditMode}
                sidebarOpen={showSidebar}
                onClose={() => setShowSidebar(false)}
                count={cachedResult.loading ? undefined : totalCount}
                focus={searchFocus}
                clearAllCriteria={() => clearAllCriteria(true)}
              />
            </Sidebar>
            <SidebarPaneContent>
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
                    onRemoveSearchTerm={() =>
                      setFilter(filter.clearSearchTerm())
                    }
                  />
                }
                selectionSection={
                  <ToolbarSelectionSection
                    selected={selectedIds.size}
                    onToggleSidebar={() => setShowSidebar(!showSidebar)}
                    onSelectAll={() => onSelectAll()}
                    onSelectNone={() => onSelectNone()}
                    operations={operations}
                  />
                }
                operationSection={operations}
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
            </SidebarPaneContent>
          </SidebarPane>
          </FacetCountsContext.Provider>
        </SidebarStateContext.Provider>
      </div>
    </TaggerContext>
  );
};

// Keep the old component for backward compatibility
export const GalleryList: React.FC<IFilteredGalleries> = (props) => {
  return <MyFilteredGalleryList {...props} />;
};

export default MyFilteredGalleryList;
