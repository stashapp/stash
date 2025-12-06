import React, { useCallback, useEffect, useMemo, useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindStudios,
  useFindStudios,
  useStudiosDestroy,
} from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { StudioTagger } from "src/components/Tagger/studios/StudioTagger";
import { StudioCardGrid } from "src/components/Studios/StudioCardGrid";
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
import cx from "classnames";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "src/extensions/ui";
import { Pagination } from "src/components/List/Pagination";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import {
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
  createBooleanCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { useFocus } from "src/extensions/hooks";
import {
  faFilter,
  faPencil,
  faPlus,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { DeleteEntityDialog } from "src/components/Shared/DeleteEntityDialog";
import { ExportDialog } from "src/components/Shared/ExportDialog";
import {
  SidebarTagsFilter,
  SidebarRatingFilter,
  SidebarStringFilter,
  SidebarNumberFilter,
  SidebarBooleanFilter,
  SidebarStashIDFilter,
  SidebarParentStudiosFilter,
  SidebarFilterSelector,
  FilterWrapper,
  SidebarDateFilter,
  SidebarIsMissingFilter,
} from "src/extensions/filters";
import { TagsCriterionOption } from "src/models/list-filter/criteria/tags";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { PatchContainerComponent } from "src/patch";
import { ListResultsHeader } from "src/extensions/ui";
import { FavoriteStudioCriterionOption } from "src/models/list-filter/criteria/favorite";
import { StashIDCriterionOption } from "src/models/list-filter/criteria/stash-ids";
import { ParentStudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import { TaggerContext } from "src/components/Tagger/context";
import { SidebarFilterDefinition } from "src/hooks/useSidebarFilters";
import { createMandatoryTimestampCriterionOption } from "src/models/list-filter/criteria/criterion";
import { StudioIsMissingCriterionOption } from "src/models/list-filter/criteria/is-missing";
import {
  useStudioFacetCounts,
  FacetCountsContext,
} from "src/hooks/useFacetCounts";

function useViewRandom(
  result: GQL.FindStudiosQueryResult,
  filter: ListFilterModel
) {
  const history = useHistory();

  const viewRandom = useCallback(async () => {
    // query for a random studio
    if (result.data?.findStudios) {
      const { count } = result.data.findStudios;

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
    }
  }, [result, filter, history]);

  return viewRandom;
}

function useAddKeybinds(
  result: GQL.FindStudiosQueryResult,
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

const StudioListContent: React.FC<{
  studios: GQL.SlimStudioDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromParent?: boolean;
}> = ({ studios, filter, selectedIds, onSelectChange, fromParent }) => {
  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <StudioCardGrid
        studios={studios as GQL.StudioDataFragment[]}
        zoomIndex={filter.zoomIndex}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
        fromParent={fromParent}
      />
    );
  }
  if (filter.displayMode === DisplayMode.Tagger) {
    return <StudioTagger studios={studios as GQL.StudioDataFragment[]} />;
  }

  return null;
};

export const MyStudiosFilterSidebarSections = PatchContainerComponent(
  "MyFilteredStudioList.SidebarSections"
);

// Define available filters for studios sidebar
const studioFilterDefinitions: SidebarFilterDefinition[] = [
  // Tier 1: Primary filters (visible by default)
  { id: "rating", messageId: "rating", defaultVisible: true },
  { id: "favorite", messageId: "favourite", defaultVisible: true },
  { id: "tags", messageId: "tags", defaultVisible: true },
  { id: "parent_studios", messageId: "parent_studios", defaultVisible: true },

  // Tier 2: Library stats
  { id: "scene_count", messageId: "scene_count", defaultVisible: false },
  { id: "image_count", messageId: "image_count", defaultVisible: false },
  { id: "gallery_count", messageId: "gallery_count", defaultVisible: false },
  { id: "tag_count", messageId: "tag_count", defaultVisible: false },
  { id: "child_count", messageId: "subsidiary_studio_count", defaultVisible: false },

  // Tier 3: Metadata
  { id: "name", messageId: "name", defaultVisible: false },
  { id: "aliases", messageId: "aliases", defaultVisible: false },
  { id: "details", messageId: "details", defaultVisible: false },
  { id: "url", messageId: "url", defaultVisible: false },

  // Tier 4: System
  { id: "is_missing", messageId: "isMissing", defaultVisible: false },
  { id: "ignore_auto_tag", messageId: "ignore_auto_tag", defaultVisible: false },
  { id: "created_at", messageId: "created_at", defaultVisible: false },
  { id: "updated_at", messageId: "updated_at", defaultVisible: false },
  { id: "stash_id", messageId: "stash_id", defaultVisible: false },
];

const SidebarContent: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  view?: View;
  showEditFilter: (editingCriterion?: string) => void;
  sidebarOpen: boolean;
  onClose?: () => void;
  count?: number;
  focus?: ReturnType<typeof useFocus>;
  clearAllCriteria: () => void;
  onFilterEditModeChange?: (isEditMode: boolean) => void;
}> = ({
  filter,
  setFilter,
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
  const UrlCriterionOption = createStringCriterionOption("url");
  const DetailsCriterionOption = createStringCriterionOption("details");
  const AliasesCriterionOption = createStringCriterionOption("aliases");
  const NameCriterionOption = createStringCriterionOption("name");
  const SceneCountCriterionOption = createMandatoryNumberCriterionOption("scene_count");
  const ImageCountCriterionOption = createMandatoryNumberCriterionOption("image_count");
  const GalleryCountCriterionOption = createMandatoryNumberCriterionOption("gallery_count");
  const TagCountCriterionOption = createMandatoryNumberCriterionOption("tag_count");
  const ChildCountCriterionOption = createMandatoryNumberCriterionOption("child_count", "subsidiary_studio_count");
  const IgnoreAutoTagCriterionOption = createBooleanCriterionOption("ignore_auto_tag");
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

      <MyStudiosFilterSidebarSections>
        <div className="sidebar-filters">
          <SidebarFilterSelector
            viewName="studios"
            filterDefinitions={studioFilterDefinitions}
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
            <FilterWrapper filterId="favorite">
              <SidebarBooleanFilter
                title={<FormattedMessage id="favourite" />}
                option={FavoriteStudioCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="favourite"
              />
            </FilterWrapper>
            <FilterWrapper filterId="tags">
              <SidebarTagsFilter
                title={<FormattedMessage id="tags" />}
                option={TagsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="tags"
              />
            </FilterWrapper>
            <FilterWrapper filterId="parent_studios">
              <SidebarParentStudiosFilter
                title={<FormattedMessage id="parent_studios" />}
                option={ParentStudiosCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="parent_studios"
              />
            </FilterWrapper>

            {/* Tier 2: Library Stats */}
            <FilterWrapper filterId="scene_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="scene_count" />}
                option={SceneCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="scene_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="image_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="image_count" />}
                option={ImageCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="image_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="gallery_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="gallery_count" />}
                option={GalleryCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="gallery_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="tag_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="tag_count" />}
                option={TagCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="tag_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="child_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="subsidiary_studio_count" />}
                option={ChildCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="child_count"
              />
            </FilterWrapper>

            {/* Tier 3: Metadata */}
            <FilterWrapper filterId="name">
              <SidebarStringFilter
                title={<FormattedMessage id="name" />}
                option={NameCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="name"
              />
            </FilterWrapper>
            <FilterWrapper filterId="aliases">
              <SidebarStringFilter
                title={<FormattedMessage id="aliases" />}
                option={AliasesCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="aliases"
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
            <FilterWrapper filterId="url">
              <SidebarStringFilter
                title={<FormattedMessage id="url" />}
                option={UrlCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="url"
              />
            </FilterWrapper>

            {/* Tier 4: System */}
            <FilterWrapper filterId="is_missing">
              <SidebarIsMissingFilter
                title={<FormattedMessage id="isMissing" />}
                option={StudioIsMissingCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="is_missing"
              />
            </FilterWrapper>
            <FilterWrapper filterId="ignore_auto_tag">
              <SidebarBooleanFilter
                title={<FormattedMessage id="ignore_auto_tag" />}
                option={IgnoreAutoTagCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="ignore_auto_tag"
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
            <FilterWrapper filterId="stash_id">
              <SidebarStashIDFilter
                title={<FormattedMessage id="stash_id" />}
                option={StashIDCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="stash_id"
              />
            </FilterWrapper>
          </SidebarFilterSelector>
        </div>
      </MyStudiosFilterSidebarSections>
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

const StudioListOperations: React.FC<{
  items: number;
  hasSelection: boolean;
  operations: IOperations[];
  onDelete: () => void;
  onCreateNew: () => void;
}> = ({ items, hasSelection, operations, onDelete, onCreateNew }) => {
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
              { entityType: intl.formatMessage({ id: "studio" }) }
            )}
          >
            <Icon icon={faPlus} />
          </Button>
        )}

        {hasSelection && (
          <>
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

interface IFilteredStudios {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
  fromParent?: boolean;
}

export const MyFilteredStudioList: React.FC<IFilteredStudios> = (props) => {
  const intl = useIntl();
  const history = useHistory();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, defaultSort, view, alterQuery, fromParent } = props;

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
        filterMode: GQL.FilterMode.Studios,
        defaultSort,
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
  const { counts: facetCounts, loading: facetLoading } = useStudioFacetCounts(filter, {
    isOpen: showSidebar ?? false,
    debounceMs: 300,
  });

  useEffect(() => {
    Mousetrap.bind("d d", () => {
      if (hasSelection) {
        onDelete();
      }
    });

    return () => {
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
    history.push("/studios/new");
  }

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

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.view_random" }),
      onClick: viewRandom,
      isDisplayed: () => totalCount > 1,
    },
    {
      text: intl.formatMessage(
        { id: "actions.create_entity" },
        { entityType: intl.formatMessage({ id: "studio" }) }
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
    <StudioListOperations
      items={items.length}
      hasSelection={hasSelection}
      operations={otherOperations}
      onDelete={onDelete}
      onCreateNew={onCreateNew}
    />
  );

  return (
    <TaggerContext>
      <div
        className={cx("item-list-container studio-list", {
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
                showEditFilter={showEditFilter}
                view={view}
                sidebarOpen={showSidebar}
                onClose={() => setShowSidebar(false)}
                count={cachedResult.loading ? undefined : totalCount}
                onFilterEditModeChange={setIsFilterEditMode}
                focus={searchFocus}
                clearAllCriteria={() => clearAllCriteria(true)}
              />
            </Sidebar>
            <SidebarPaneContent>
              <FilteredListToolbar2
                className="studio-list-toolbar"
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
                <StudioListContent
                  filter={effectiveFilter}
                  studios={items}
                  selectedIds={selectedIds}
                  onSelectChange={onSelectChange}
                  fromParent={fromParent}
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

// Backward compatibility wrapper
export const StudioList: React.FC<IFilteredStudios> = (props) => {
  return <MyFilteredStudioList {...props} />;
};
