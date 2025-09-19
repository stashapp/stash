import React, { useCallback, useContext, useEffect, useMemo } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindPerformers,
  useFindPerformers,
  usePerformersDestroy,
} from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { PerformerTagger } from "../Tagger/performers/PerformerTagger";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { IPerformerCardExtraCriteria } from "./PerformerCard";
import { PerformerListTable } from "./PerformerListTable";
import { EditPerformersDialog } from "./EditPerformersDialog";
import { cmToImperial, cmToInches, kgToLbs } from "src/utils/units";
import TextUtils from "src/utils/text";
import { PerformerCardGrid } from "./PerformerCardGrid";
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
import cx from "classnames";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
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
import { Criterion } from "src/models/list-filter/criteria/criterion";
import useFocus from "src/utils/myFocus";
import {
  faPencil,
  faPlus,
  faSliders,
  faTimes,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { TaggerContext } from "../Tagger/context";
import { GenderCriterionOption } from "src/models/list-filter/criteria/gender";
import { SidebarGenderFilter } from "../List/Filters/GenderFilter";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import { TagsCriterionOption } from "src/models/list-filter/criteria/tags";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { FavoritePerformerCriterionOption } from "src/models/list-filter/criteria/favorite";
import { StashIDCriterionOption } from "src/models/list-filter/criteria/stash-ids";
import { SidebarStashIDFilter } from "../List/Filters/StashIDFilter";
import { CountryCriterionOption } from "src/models/list-filter/criteria/country";
import { SidebarCountryFilter } from "../List/Filters/CountryFilter";
import { TattoosCriterionOption } from "src/models/list-filter/criteria/tattoos";
import { SidebarStringFilter } from "../List/Filters/StringFilter";
import { PiercingsCriterionOption } from "src/models/list-filter/criteria/piercings";

function getItems(result: GQL.FindPerformersQueryResult) {
  return result?.data?.findPerformers?.performers ?? [];
}

function getCount(result: GQL.FindPerformersQueryResult) {
  return result?.data?.findPerformers?.count ?? 0;
}

function useOpenRandom(filter: ListFilterModel, count: number) {
  const history = useHistory();

  const openRandom = useCallback(async () => {
    // query for a random performer
    if (count === 0) {
      return;
    }

    const pages = Math.ceil(count / filter.itemsPerPage);
    const page = Math.floor(Math.random() * pages) + 1;

    const indexMax = Math.min(filter.itemsPerPage, count);
    const index = Math.floor(Math.random() * indexMax);
    const filterCopy = cloneDeep(filter);
    filterCopy.currentPage = page;
    filterCopy.sortBy = "random";
    const queryResults = await queryFindPerformers(filterCopy);
    const performer = queryResults.data.findPerformers.performers[index];
    if (performer) {
      history.push(`/performers/${performer.id}`);
    }
  }, [filter, count, history]);

  return openRandom;
}

function useAddKeybinds(filter: ListFilterModel, count: number) {
  const openRandom = useOpenRandom(filter, count);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      openRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [openRandom]);
}

export const FormatHeight = (height?: number | null) => {
  const intl = useIntl();
  if (!height) {
    return "";
  }

  const [feet, inches] = cmToImperial(height);

  return (
    <span className="performer-height">
      <span className="height-metric">
        {intl.formatNumber(height, {
          style: "unit",
          unit: "centimeter",
          unitDisplay: "short",
        })}
      </span>
      <span className="height-imperial">
        {intl.formatNumber(feet, {
          style: "unit",
          unit: "foot",
          unitDisplay: "narrow",
        })}
        {intl.formatNumber(inches, {
          style: "unit",
          unit: "inch",
          unitDisplay: "narrow",
        })}
      </span>
    </span>
  );
};

export const FormatAge = (
  birthdate?: string | null,
  deathdate?: string | null
) => {
  if (!birthdate) {
    return "";
  }
  const age = TextUtils.age(birthdate, deathdate);

  return (
    <span className="performer-age">
      <span className="age">{age}</span>
      <span className="birthdate"> ({birthdate})</span>
    </span>
  );
};

export const FormatWeight = (weight?: number | null) => {
  const intl = useIntl();
  if (!weight) {
    return "";
  }

  const lbs = kgToLbs(weight);

  return (
    <span className="performer-weight">
      <span className="weight-metric">
        {intl.formatNumber(weight, {
          style: "unit",
          unit: "kilogram",
          unitDisplay: "short",
        })}
      </span>
      <span className="weight-imperial">
        {intl.formatNumber(lbs, {
          style: "unit",
          unit: "pound",
          unitDisplay: "short",
        })}
      </span>
    </span>
  );
};

export const FormatCircumcised = (circumcised?: GQL.CircumisedEnum | null) => {
  const intl = useIntl();
  if (!circumcised) {
    return "";
  }

  return (
    <span className="penis-circumcised">
      {intl.formatMessage({
        id: "circumcised_types." + circumcised,
      })}
    </span>
  );
};

export const FormatPenisLength = (penis_length?: number | null) => {
  const intl = useIntl();
  if (!penis_length) {
    return "";
  }

  const inches = cmToInches(penis_length);

  return (
    <span className="performer-penis-length">
      <span className="penis-length-metric">
        {intl.formatNumber(penis_length, {
          style: "unit",
          unit: "centimeter",
          unitDisplay: "short",
          maximumFractionDigits: 2,
        })}
      </span>
      <span className="penis-length-imperial">
        {intl.formatNumber(inches, {
          style: "unit",
          unit: "inch",
          unitDisplay: "narrow",
          maximumFractionDigits: 2,
        })}
      </span>
    </span>
  );
};

const PerformerListContent: React.FC<{
  performers: GQL.SlimPerformerDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}> = ({ performers, filter, selectedIds, onSelectChange, extraCriteria }) => {
  if (performers.length === 0) {
    return null;
  }

  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <PerformerCardGrid
        performers={performers as any}
        zoomIndex={filter.zoomIndex}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
        extraCriteria={extraCriteria}
      />
    );
  }
  if (filter.displayMode === DisplayMode.List) {
    return (
      <PerformerListTable
        performers={performers as any}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
      />
    );
  }
  if (filter.displayMode === DisplayMode.Tagger) {
    return <PerformerTagger performers={performers as any} />;
  }

  return null;
};

export const MyPerformersFilterSidebarSections = PatchContainerComponent(
  "MyFilteredPerformerList.SidebarSections"
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

      <MyPerformersFilterSidebarSections>
        <SidebarTagsFilter
          title={<FormattedMessage id="tags" />}
          data-type={TagsCriterionOption.type}
          option={TagsCriterionOption}
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarGenderFilter
          title={<FormattedMessage id="gender" />}
          data-type={GenderCriterionOption.type}
          option={GenderCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarCountryFilter
          title={<FormattedMessage id="country" />}
          data-type={CountryCriterionOption.type}
          option={CountryCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarStringFilter
          title={<FormattedMessage id="piercings" />}
          data-type={PiercingsCriterionOption.type}
          option={PiercingsCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarStringFilter
          title={<FormattedMessage id="tattoos" />}
          data-type={TattoosCriterionOption.type}
          option={TattoosCriterionOption}
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
          title={<FormattedMessage id="favourite" />}
          data-type={FavoritePerformerCriterionOption.type}
          option={FavoritePerformerCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarStashIDFilter
          title={<FormattedMessage id="stash_id" />}
          data-type={StashIDCriterionOption.type}
          option={StashIDCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
      </MyPerformersFilterSidebarSections>

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
  items: GQL.SlimPerformerDataFragment[];
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
                { entityType: intl.formatMessage({ id: "performer" }) }
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

          <OperationDropdown className="performer-list-operations">
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
    <ButtonToolbar className="performer-list-header">
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

interface IFilteredPerformers {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
  extraCriteria?: IPerformerCardExtraCriteria;
}

export const MyFilteredPerformerList = (props: IFilteredPerformers) => {
  const intl = useIntl();
  const history = useHistory();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, defaultSort, view, alterQuery, extraCriteria } = props;

  // States
  const {
    showSidebar,
    setShowSidebar,
    loading: sidebarStateLoading,
  } = useSidebarState(view);

  const { filterState, queryResult, modalState, listSelect, showEditFilter } =
    useFilteredItemList({
      filterStateProps: {
        filterMode: GQL.FilterMode.Performers,
        defaultSort,
        view,
        useURL: alterQuery,
      },
      queryResultProps: {
        useResult: useFindPerformers,
        getCount: (r) => r.data?.findPerformers.count ?? 0,
        getItems: (r) => r.data?.findPerformers.performers ?? [],
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

  const openRandom = useOpenRandom(filter, totalCount);

  function onCreateNew() {
    history.push("/performers/new");
  }

  function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          performers: {
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
      <EditPerformersDialog
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
        singularEntity={intl.formatMessage({ id: "performer" })}
        pluralEntity={intl.formatMessage({ id: "performers" })}
        destroyMutation={usePerformersDestroy}
      />
    );
  }

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.open_random" }),
      onClick: openRandom,
      isDisplayed: () => totalCount > 1,
    },
    {
      text: intl.formatMessage(
        { id: "actions.create_entity" },
        { entityType: intl.formatMessage({ id: "performer" }) }
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
        className={cx("item-list-container performer-list", {
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
              className={cx("performer-list-toolbar", {
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
              <PerformerListContent
                filter={effectiveFilter}
                performers={items}
                selectedIds={selectedIds}
                onSelectChange={onSelectChange}
                extraCriteria={extraCriteria}
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
export const PerformerList: React.FC<IFilteredPerformers> = (props) => {
  return <MyFilteredPerformerList {...props} />;
};

export default MyFilteredPerformerList;
