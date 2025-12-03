import React, { useCallback, useEffect } from "react";
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
import { GroupCardGrid } from "./GroupCardGrid";
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
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "../Shared/Sidebar";
import cx from "classnames";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/MyFilterSidebar";
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
import { EditGroupsDialog } from "./EditGroupsDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { SidebarPerformersFilter } from "../List/Filters/PerformersFilter";
import { SidebarStudiosFilter } from "../List/Filters/StudiosFilter";
import { PerformersCriterionOption } from "src/models/list-filter/criteria/performers";
import { StudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import { TagsCriterionOption } from "src/models/list-filter/criteria/tags";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { SidebarStringFilter } from "../List/Filters/StringFilter";
import { SidebarNumberFilter } from "../List/Filters/NumberFilter";
import { PatchContainerComponent } from "src/patch";
import { SidebarDateFilter } from "../List/Filters/DateFilter";
import { ListResultsHeader } from "../List/MyListResultsHeader";

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
}) => {
  const showResultsId =
    count !== undefined ? "actions.show_count_results" : "actions.show_results";

  const containingGroupCountCriterionOption =
    createMandatoryNumberCriterionOption("containing_group_count");
  const subGroupCountCriterionOption =
    createMandatoryNumberCriterionOption("sub_group_count");
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

      <MyGroupsFilterSidebarSections>
        <div className="sidebar-filters">
          <div className="sidebar-section-header">
            <Icon icon={faFilter} />
            <FormattedMessage id="filters" />
          </div>
          <SidebarPerformersFilter
            title={<FormattedMessage id="performers" />}
            option={PerformersCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="performers"
          />
          <SidebarStudiosFilter
            title={<FormattedMessage id="studios" />}
            option={StudiosCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="studios"
          />
          <SidebarTagsFilter
            title={<FormattedMessage id="tags" />}
            option={TagsCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="tags"
          />
          <SidebarDateFilter
            title={<FormattedMessage id="date" />}
            data-type={DateCriterionOption.type}
            option={DateCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="date"
          />
          <SidebarRatingFilter
            title={<FormattedMessage id="rating" />}
            option={RatingCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="rating"
          />
          <SidebarStringFilter
            title={<FormattedMessage id="url" />}
            data-type={UrlCriterionOption.type}
            option={UrlCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="url"
          />
          <SidebarNumberFilter
            title={<FormattedMessage id="containing_group_count" />}
            option={containingGroupCountCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="containing_group_count"
          />
          <SidebarNumberFilter
            title={<FormattedMessage id="sub_group_count" />}
            option={subGroupCountCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="sub_group_count"
          />
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
    sectionOpen,
    setSectionOpen,
  } = useSidebarState(view);

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

      <SidebarStateContext.Provider value={{ sectionOpen, setSectionOpen }}>
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
      </SidebarStateContext.Provider>
    </div>
  );
};

// Backward compatibility wrapper
export const GroupList: React.FC<IFilteredGroups> = (props) => {
  return <MyFilteredGroupList {...props} />;
};
