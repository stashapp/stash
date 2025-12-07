import React, { useCallback, useEffect, useMemo, useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindGroups,
  useFindGroups,
  useGroupsDestroy,
} from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { GroupCardGrid } from "src/components/Groups/GroupCardGrid";
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
import { EditGroupsDialog } from "src/components/Groups/EditGroupsDialog";
import { DeleteEntityDialog } from "src/components/Shared/DeleteEntityDialog";
import { ExportDialog } from "src/components/Shared/ExportDialog";
import {
  SidebarPerformersFilter,
  SidebarStudiosFilter,
  SidebarGroupsFilter,
  SidebarTagsFilter,
  SidebarRatingFilter,
  SidebarStringFilter,
  SidebarNumberFilter,
  SidebarDateFilter,
  SidebarFilterSelector,
  FilterWrapper,
  SidebarBooleanFilter,
  SidebarDurationFilter,
  SidebarIsMissingFilter,
} from "src/extensions/filters";
import { PerformersCriterionOption } from "src/models/list-filter/criteria/performers";
import { StudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import {
  ContainingGroupsCriterionOption,
  SubGroupsCriterionOption,
} from "src/models/list-filter/criteria/groups";
import { TagsCriterionOption } from "src/models/list-filter/criteria/tags";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { PatchContainerComponent } from "src/patch";
import { ListResultsHeader } from "src/extensions/ui";
import { SidebarFilterDefinition } from "src/extensions/hooks/useSidebarFilters";
import {
  createMandatoryTimestampCriterionOption,
  createDurationCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { GroupIsMissingCriterionOption } from "src/models/list-filter/criteria/is-missing";
import {
  useGroupFacetCounts,
  FacetCountsContext,
} from "src/extensions/hooks/useFacetCounts";

function useViewRandom(
  result: GQL.FindGroupsQueryResult,
  filter: ListFilterModel
) {
  const history = useHistory();

  const viewRandom = useCallback(async () => {
    // query for a random group
    if (result.data?.findGroups) {
      const { count } = result.data.findGroups;

      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindGroups(filterCopy);
      if (singleResult.data.findGroups.groups.length === 1) {
        const { id } = singleResult.data.findGroups.groups[0];
        // navigate to the group page
        history.push(`/groups/${id}`);
      }
    }
  }, [result, filter, history]);

  return viewRandom;
}

function useAddKeybinds(
  result: GQL.FindGroupsQueryResult,
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

const GroupListContent: React.FC<{
  groups: GQL.SlimGroupDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
}> = ({ groups, filter, selectedIds, onSelectChange, fromGroupId, onMove }) => {
  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <GroupCardGrid
        groups={groups as GQL.GroupDataFragment[]}
        zoomIndex={filter.zoomIndex}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
        fromGroupId={fromGroupId}
        onMove={onMove}
      />
    );
  }

  return null;
};

export const MyGroupsFilterSidebarSections = PatchContainerComponent(
  "MyFilteredGroupList.SidebarSections"
);

// Define available filters for groups sidebar
const groupFilterDefinitions: SidebarFilterDefinition[] = [
  // Tier 1: Primary filters (visible by default)
  { id: "rating", messageId: "rating", defaultVisible: true },
  { id: "date", messageId: "date", defaultVisible: true },
  { id: "tags", messageId: "tags", defaultVisible: true },
  { id: "performers", messageId: "performers", defaultVisible: true },
  { id: "studios", messageId: "studios", defaultVisible: true },
  { id: "duration", messageId: "duration", defaultVisible: true },
  { id: "containing_groups", messageId: "containing_groups", defaultVisible: false },
  { id: "sub_groups", messageId: "sub_groups", defaultVisible: false },

  // Tier 2: Library stats
  { id: "containing_group_count", messageId: "containing_group_count", defaultVisible: false },
  { id: "sub_group_count", messageId: "sub_group_count", defaultVisible: false },
  { id: "tag_count", messageId: "tag_count", defaultVisible: false },

  // Tier 3: Metadata
  { id: "name", messageId: "name", defaultVisible: false },
  { id: "director", messageId: "director", defaultVisible: false },
  { id: "synopsis", messageId: "synopsis", defaultVisible: false },
  { id: "url", messageId: "url", defaultVisible: false },

  // Tier 4: System
  { id: "is_missing", messageId: "isMissing", defaultVisible: false },
  { id: "created_at", messageId: "created_at", defaultVisible: false },
  { id: "updated_at", messageId: "updated_at", defaultVisible: false },
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
  const containingGroupCountCriterionOption = createMandatoryNumberCriterionOption("containing_group_count");
  const subGroupCountCriterionOption = createMandatoryNumberCriterionOption("sub_group_count");
  const tagCountCriterionOption = createMandatoryNumberCriterionOption("tag_count");
  const UrlCriterionOption = createStringCriterionOption("url");
  const NameCriterionOption = createStringCriterionOption("name");
  const DirectorCriterionOption = createStringCriterionOption("director");
  const SynopsisCriterionOption = createStringCriterionOption("synopsis");
  const DateCriterionOption = createDateCriterionOption("date");
  const DurationCriterionOption = createDurationCriterionOption("duration");
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

      <MyGroupsFilterSidebarSections>
        <div className="sidebar-filters">
          <SidebarFilterSelector
            viewName="groups"
            filterDefinitions={groupFilterDefinitions}
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
                sectionID="tags"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performers">
              <SidebarPerformersFilter
                title={<FormattedMessage id="performers" />}
                option={PerformersCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performers"
              />
            </FilterWrapper>
            <FilterWrapper filterId="studios">
              <SidebarStudiosFilter
                title={<FormattedMessage id="studios" />}
                option={StudiosCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="studios"
              />
            </FilterWrapper>
            <FilterWrapper filterId="duration">
              <SidebarDurationFilter
                title={<FormattedMessage id="duration" />}
                option={DurationCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="duration"
              />
            </FilterWrapper>
            <FilterWrapper filterId="containing_groups">
              <SidebarGroupsFilter
                title={<FormattedMessage id="containing_groups" />}
                option={ContainingGroupsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="containing_groups"
              />
            </FilterWrapper>
            <FilterWrapper filterId="sub_groups">
              <SidebarGroupsFilter
                title={<FormattedMessage id="sub_groups" />}
                option={SubGroupsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="sub_groups"
              />
            </FilterWrapper>

            {/* Tier 2: Library Stats */}
            <FilterWrapper filterId="containing_group_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="containing_group_count" />}
                option={containingGroupCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="containing_group_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="sub_group_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="sub_group_count" />}
                option={subGroupCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="sub_group_count"
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
            <FilterWrapper filterId="director">
              <SidebarStringFilter
                title={<FormattedMessage id="director" />}
                option={DirectorCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="director"
              />
            </FilterWrapper>
            <FilterWrapper filterId="synopsis">
              <SidebarStringFilter
                title={<FormattedMessage id="synopsis" />}
                option={SynopsisCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="synopsis"
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
                option={GroupIsMissingCriterionOption}
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
      </MyGroupsFilterSidebarSections>
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

const GroupListOperations: React.FC<{
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
              { entityType: intl.formatMessage({ id: "group" }) }
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

interface IFilteredGroups {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
  fromGroupId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
}

export const MyFilteredGroupList: React.FC<IFilteredGroups> = (props) => {
  const intl = useIntl();
  const history = useHistory();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, defaultSort, view, alterQuery, fromGroupId, onMove } =
    props;

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
        filterMode: GQL.FilterMode.Groups,
        defaultSort,
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
  const { counts: facetCounts, loading: facetLoading } = useGroupFacetCounts(filter, {
    isOpen: showSidebar ?? false,
    debounceMs: 300,
  });

  useEffect(() => {
    Mousetrap.bind("e", () => {
      if (hasSelection) {
        onEdit();
      }
    });

    Mousetrap.bind("d d", () => {
      if (hasSelection) {
        onDelete();
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
    history.push("/groups/new");
  }

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
      <EditGroupsDialog selected={selectedItems} onClose={onCloseEditDelete} />
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

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.view_random" }),
      onClick: viewRandom,
      isDisplayed: () => totalCount > 1,
    },
    {
      text: intl.formatMessage(
        { id: "actions.create_entity" },
        { entityType: intl.formatMessage({ id: "group" }) }
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
    <GroupListOperations
      items={items.length}
      hasSelection={hasSelection}
      operations={otherOperations}
      onEdit={onEdit}
      onDelete={onDelete}
      onCreateNew={onCreateNew}
    />
  );

  return (
    <div
      className={cx("item-list-container group-list", {
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
              onFilterEditModeChange={setIsFilterEditMode}
              onClose={() => setShowSidebar(false)}
              count={cachedResult.loading ? undefined : totalCount}
              focus={searchFocus}
              clearAllCriteria={() => clearAllCriteria(true)}
            />
          </Sidebar>
          <SidebarPaneContent>
            <FilteredListToolbar2
              className="group-list-toolbar"
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
              <GroupListContent
                filter={effectiveFilter}
                groups={items}
                selectedIds={selectedIds}
                onSelectChange={onSelectChange}
                fromGroupId={fromGroupId}
                onMove={onMove}
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
  );
};

// Backward compatibility wrapper
export const GroupList: React.FC<IFilteredGroups> = (props) => {
  return <MyFilteredGroupList {...props} />;
};
