import React, { useCallback, useEffect } from "react";
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
import { StudioTagger } from "../Tagger/studios/StudioTagger";
import { StudioCardGrid } from "./StudioCardGrid";
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
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { TagsCriterionOption } from "src/models/list-filter/criteria/tags";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { SidebarStringFilter } from "../List/Filters/StringFilter";
import { SidebarNumberFilter } from "../List/Filters/NumberFilter";
import { PatchContainerComponent } from "src/patch";
import { ListResultsHeader } from "../List/MyListResultsHeader";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { FavoriteStudioCriterionOption } from "src/models/list-filter/criteria/favorite";
import { StashIDCriterionOption } from "src/models/list-filter/criteria/stash-ids";
import { SidebarStashIDFilter } from "../List/Filters/StashIDFilter";
import { ParentStudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import { SidebarParentStudiosFilter } from "../List/Filters/StudiosFilter";
import { TaggerContext } from "../Tagger/context";

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

  const UrlCriterionOption = createStringCriterionOption("url");
  const DetailsCriterionOption = createStringCriterionOption("details");
  const AliasesCriterionOption = createStringCriterionOption("aliases");
  const SceneCountCriterionOption =
    createMandatoryNumberCriterionOption("scene_count");
  const ImageCountCriterionOption =
    createMandatoryNumberCriterionOption("image_count");
  const GalleryCountCriterionOption =
    createMandatoryNumberCriterionOption("gallery_count");
  const TagCountCriterionOption =
    createMandatoryNumberCriterionOption("tag_count");
  const ChildCountCriterionOption = createMandatoryNumberCriterionOption(
    "child_count",
    "subsidiary_studio_count"
  );

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
          <div className="sidebar-section-header">
            <Icon icon={faFilter} />
            <FormattedMessage id="filters" />
          </div>
          <SidebarParentStudiosFilter
            title={<FormattedMessage id="parent_studios" />}
            option={ParentStudiosCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="parent_studios"
          />
          <SidebarTagsFilter
            title={<FormattedMessage id="tags" />}
            option={TagsCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="tags"
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
          <SidebarStringFilter
            title={<FormattedMessage id="details" />}
            data-type={DetailsCriterionOption.type}
            option={DetailsCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="details"
          />
          <SidebarStringFilter
            title={<FormattedMessage id="aliases" />}
            data-type={AliasesCriterionOption.type}
            option={AliasesCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="aliases"
          />
          <SidebarNumberFilter
            title={<FormattedMessage id="scene_count" />}
            data-type={SceneCountCriterionOption.type}
            option={SceneCountCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="scene_count"
          />
          <SidebarNumberFilter
            title={<FormattedMessage id="image_count" />}
            data-type={ImageCountCriterionOption.type}
            option={ImageCountCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="image_count"
          />
          <SidebarNumberFilter
            title={<FormattedMessage id="gallery_count" />}
            data-type={GalleryCountCriterionOption.type}
            option={GalleryCountCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="gallery_count"
          />
          <SidebarNumberFilter
            title={<FormattedMessage id="tag_count" />}
            data-type={TagCountCriterionOption.type}
            option={TagCountCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="tag_count"
          />
          <SidebarNumberFilter
            title={<FormattedMessage id="subsidiary_studio_count" />}
            data-type={ChildCountCriterionOption.type}
            option={ChildCountCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="child_count"
          />
          <SidebarBooleanFilter
            title={<FormattedMessage id="favourite" />}
            data-type={FavoriteStudioCriterionOption.type}
            option={FavoriteStudioCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="favourite"
          />
          <SidebarStashIDFilter
            title={<FormattedMessage id="stash_id" />}
            data-type={StashIDCriterionOption.type}
            option={StashIDCriterionOption}
            filter={filter}
            setFilter={setFilter}
            sectionID="stash_id"
          />
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
    sectionOpen,
    setSectionOpen,
  } = useSidebarState(view);

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
        </SidebarStateContext.Provider>
      </div>
    </TaggerContext>
  );
};

// Backward compatibility wrapper
export const StudioList: React.FC<IFilteredStudios> = (props) => {
  return <MyFilteredStudioList {...props} />;
};
