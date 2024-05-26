import React, {
  useCallback,
  useContext,
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import cloneDeep from "lodash-es/cloneDeep";
import isEqual from "lodash-es/isEqual";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { QueryResult } from "@apollo/client";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { useInterfaceLocalForage } from "src/hooks/LocalForage";
import { useHistory, useLocation } from "react-router-dom";
import { ConfigurationContext } from "src/hooks/Config";
import { getFilterOptions } from "src/models/list-filter/factory";
import { useFindDefaultFilter } from "src/core/StashService";
import { Pagination, PaginationIndex } from "./Pagination";
import { EditFilterDialog } from "src/components/List/EditFilterDialog";
import { ListFilter } from "./ListFilter";
import { FilterTags } from "./FilterTags";
import { DisplayModeSelect, ZoomSelect } from "./ListViewOptions";
import { ListOperationButtons } from "./ListOperationButtons";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { DisplayMode } from "src/models/list-filter/types";
import { ButtonToolbar } from "react-bootstrap";
import { useListSelect } from "src/hooks/listSelect";
import { useFilterConfig } from "./util";

export enum PersistanceLevel {
  // do not load default query or persist display mode
  NONE,
  // load default query, don't load or persist display mode
  ALL,
  // load and persist display mode only
  VIEW,
}

interface IDataItem {
  id: string;
}

export interface IItemListOperation<T extends QueryResult> {
  text: string;
  onClick: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => Promise<void>;
  isDisplayed?: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => boolean;
  postRefetch?: boolean;
  icon?: IconDefinition;
  buttonVariant?: string;
}

interface IItemListOptions<T extends QueryResult, E extends IDataItem> {
  filterMode: GQL.FilterMode;
  useResult: (filter: ListFilterModel) => T;
  getCount: (data: T) => number;
  renderMetadataByline?: (data: T) => React.ReactNode;
  getItems: (data: T) => E[];
}

interface IRenderListProps {
  filter: ListFilterModel;
  onChangePage: (page: number) => void;
  updateFilter: (filter: ListFilterModel) => void;
}

interface IItemListProps<T extends QueryResult, E extends IDataItem> {
  persistState?: PersistanceLevel;
  persistanceKey?: string;
  defaultSort?: string;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  filterDialog?: (
    criteria: Criterion<CriterionValue>[],
    setCriteria: (v: Criterion<CriterionValue>[]) => void
  ) => React.ReactNode;
  zoomable?: boolean;
  selectable?: boolean;
  alterQuery?: boolean;
  defaultZoomIndex?: number;
  otherOperations?: IItemListOperation<T>[];
  renderContent: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void,
    onChangePage: (page: number) => void,
    pageCount: number
  ) => React.ReactNode;
  renderEditDialog?: (
    selected: E[],
    onClose: (applied: boolean) => void
  ) => React.ReactNode;
  renderDeleteDialog?: (
    selected: E[],
    onClose: (confirmed: boolean) => void
  ) => React.ReactNode;
  addKeybinds?: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => () => void;
}

const getSelectedData = <I extends IDataItem>(
  data: I[],
  selectedIds: Set<string>
) => data.filter((value) => selectedIds.has(value.id));

/**
 * A factory function for ItemList components.
 * IMPORTANT: as the component manipulates the URL query string, if there are
 * ever multiple ItemLists rendered at once, all but one of them need to have
 * `alterQuery` set to false to prevent conflicts.
 */
export function makeItemList<T extends QueryResult, E extends IDataItem>({
  filterMode,
  useResult,
  getCount,
  renderMetadataByline,
  getItems,
}: IItemListOptions<T, E>) {
  const filterOptions = getFilterOptions(filterMode);

  const RenderList: React.FC<IItemListProps<T, E> & IRenderListProps> = ({
    filter,
    filterHook,
    onChangePage: _onChangePage,
    updateFilter,
    persistState,
    zoomable,
    selectable,
    otherOperations,
    renderContent,
    renderEditDialog,
    renderDeleteDialog,
    addKeybinds,
  }) => {
    const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

    const [editingCriterion, setEditingCriterion] = useState<string>();
    const [showEditFilter, setShowEditFilter] = useState(false);

    const { criterionOptions, setCriterionOptions } = useFilterConfig(
      filter.mode
    );

    const effectiveFilter = useMemo(() => {
      if (filterHook) {
        return filterHook(cloneDeep(filter));
      }
      return filter;
    }, [filter, filterHook]);

    const result = useResult(effectiveFilter);
    const [totalCount, setTotalCount] = useState(0);
    const [metadataByline, setMetadataByline] = useState<React.ReactNode>();
    const items = useMemo(() => getItems(result), [result]);

    const { selectedIds, onSelectChange, onSelectAll, onSelectNone } =
      useListSelect(items);

    const [arePaging, setArePaging] = useState(false);
    const hidePagination = !arePaging && result.loading;

    // useLayoutEffect to set total count before paint, avoiding a 0 being displayed
    useLayoutEffect(() => {
      if (result.loading) return;
      setArePaging(false);

      setTotalCount(getCount(result));
      setMetadataByline(renderMetadataByline?.(result));
    }, [result]);

    const onChangePage = useCallback(
      (page: number) => {
        setArePaging(true);
        _onChangePage(page);
      },
      [_onChangePage]
    );

    // handle case where page is more than there are pages
    useEffect(() => {
      const pages = Math.ceil(totalCount / filter.itemsPerPage);
      if (pages > 0 && filter.currentPage > pages) {
        onChangePage(pages);
      }
    }, [filter, onChangePage, totalCount]);

    // set up hotkeys
    useEffect(() => {
      Mousetrap.bind("f", (e) => {
        setShowEditFilter(true);
        // prevent default behavior of typing f in a text field
        // otherwise the filter dialog closes, the query field is focused and
        // f is typed.
        e.preventDefault();
      });

      return () => {
        Mousetrap.unbind("f");
      };
    }, []);
    useEffect(() => {
      const pages = Math.ceil(totalCount / filter.itemsPerPage);
      Mousetrap.bind("right", () => {
        if (filter.currentPage < pages) {
          onChangePage(filter.currentPage + 1);
        }
      });
      Mousetrap.bind("left", () => {
        if (filter.currentPage > 1) {
          onChangePage(filter.currentPage - 1);
        }
      });
      Mousetrap.bind("shift+right", () => {
        onChangePage(Math.min(pages, filter.currentPage + 10));
      });
      Mousetrap.bind("shift+left", () => {
        onChangePage(Math.max(1, filter.currentPage - 10));
      });
      Mousetrap.bind("ctrl+end", () => {
        onChangePage(pages);
      });
      Mousetrap.bind("ctrl+home", () => {
        onChangePage(1);
      });

      return () => {
        Mousetrap.unbind("right");
        Mousetrap.unbind("left");
        Mousetrap.unbind("shift+right");
        Mousetrap.unbind("shift+left");
        Mousetrap.unbind("ctrl+end");
        Mousetrap.unbind("ctrl+home");
      };
    }, [filter, onChangePage, totalCount]);
    useEffect(() => {
      if (addKeybinds) {
        const unbindExtras = addKeybinds(result, filter, selectedIds);
        return () => {
          unbindExtras();
        };
      }
    }, [addKeybinds, result, filter, selectedIds]);

    function onChangeZoom(newZoomIndex: number) {
      const newFilter = cloneDeep(filter);
      newFilter.zoomIndex = newZoomIndex;
      updateFilter(newFilter);
    }

    async function onOperationClicked(o: IItemListOperation<T>) {
      await o.onClick(result, filter, selectedIds);
      if (o.postRefetch) {
        result.refetch();
      }
    }

    const operations = otherOperations?.map((o) => ({
      text: o.text,
      onClick: () => {
        onOperationClicked(o);
      },
      isDisplayed: () => {
        if (o.isDisplayed) {
          return o.isDisplayed(result, filter, selectedIds);
        }

        return true;
      },
      icon: o.icon,
      buttonVariant: o.buttonVariant,
    }));

    function onEdit() {
      setIsEditDialogOpen(true);
    }

    function onEditDialogClosed(applied: boolean) {
      if (applied) {
        onSelectNone();
      }
      setIsEditDialogOpen(false);

      // refetch
      result.refetch();
    }

    function onDelete() {
      setIsDeleteDialogOpen(true);
    }

    function onDeleteDialogClosed(deleted: boolean) {
      if (deleted) {
        onSelectNone();
      }
      setIsDeleteDialogOpen(false);

      // refetch
      result.refetch();
    }

    function renderPagination() {
      if (hidePagination) return;
      return (
        <div className="filter-container">
          <Pagination
            itemsPerPage={filter.itemsPerPage}
            currentPage={filter.currentPage}
            totalItems={totalCount}
            metadataByline={metadataByline}
            onChangePage={onChangePage}
          />
        </div>
      );
    }

    function renderPaginationIndex() {
      if (hidePagination) return;
      return (
        <PaginationIndex
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
          metadataByline={metadataByline}
        />
      );
    }

    function maybeRenderContent() {
      if (result.loading) {
        return <LoadingIndicator />;
      }
      if (result.error) {
        return <h1>{result.error.message}</h1>;
      }

      const pages = Math.ceil(totalCount / filter.itemsPerPage);
      return (
        <>
          {renderContent(
            result,
            // #4780 - use effectiveFilter to ensure filterHook is applied
            effectiveFilter,
            selectedIds,
            onSelectChange,
            onChangePage,
            pages
          )}
          {!!pages && (
            <>
              {renderPaginationIndex()}
              {renderPagination()}
            </>
          )}
        </>
      );
    }

    function onChangeDisplayMode(displayMode: DisplayMode) {
      const newFilter = cloneDeep(filter);
      newFilter.displayMode = displayMode;
      updateFilter(newFilter);
    }

    function onRemoveCriterion(removedCriterion: Criterion<CriterionValue>) {
      const newFilter = cloneDeep(filter);
      newFilter.criteria = newFilter.criteria.filter(
        (criterion) => criterion.getId() !== removedCriterion.getId()
      );
      newFilter.currentPage = 1;
      updateFilter(newFilter);
    }

    function onClearAllCriteria() {
      const newFilter = cloneDeep(filter);
      newFilter.criteria = [];
      newFilter.currentPage = 1;
      updateFilter(newFilter);
    }

    function onApplyEditFilter(f?: ListFilterModel) {
      setShowEditFilter(false);
      setEditingCriterion(undefined);

      if (!f) return;
      updateFilter(f);
    }

    const minZoom = 0;
    const maxZoom = 3;

    return (
      <div className="item-list-container">
        <ButtonToolbar className="justify-content-center">
          <ListFilter
            onFilterUpdate={updateFilter}
            filter={filter}
            filterOptions={filterOptions}
            openFilterDialog={() => setShowEditFilter(true)}
            persistState={persistState}
          />
          <div className="ml-2 mb-2">
            <ListOperationButtons
              onSelectAll={selectable ? onSelectAll : undefined}
              onSelectNone={selectable ? onSelectNone : undefined}
              otherOperations={operations}
              itemsSelected={selectedIds.size > 0}
              onEdit={renderEditDialog ? onEdit : undefined}
              onDelete={renderDeleteDialog ? onDelete : undefined}
            />
          </div>
          <div className="mb-2 d-inline-flex">
            <DisplayModeSelect
              displayMode={filter.displayMode}
              displayModeOptions={filterOptions.displayModeOptions}
              onSetDisplayMode={onChangeDisplayMode}
            />
            {!!zoomable && filter.displayMode === DisplayMode.Grid && (
              <div className="ml-2 d-none d-sm-inline-flex">
                <ZoomSelect
                  minZoom={minZoom}
                  maxZoom={maxZoom}
                  zoomIndex={filter.zoomIndex ?? minZoom}
                  onChangeZoom={onChangeZoom}
                />
              </div>
            )}
          </div>
        </ButtonToolbar>
        <FilterTags
          criteria={filter.criteria}
          onEditCriterion={(c) => setEditingCriterion(c.criterionOption.type)}
          onRemoveCriterion={onRemoveCriterion}
          onRemoveAll={() => onClearAllCriteria()}
        />
        {(showEditFilter || editingCriterion) && (
          <EditFilterDialog
            filter={filter}
            criterionOptions={criterionOptions}
            setCriterionOptions={setCriterionOptions}
            onClose={onApplyEditFilter}
            editingCriterion={editingCriterion}
          />
        )}
        {isEditDialogOpen &&
          renderEditDialog &&
          renderEditDialog(getSelectedData(items, selectedIds), (applied) =>
            onEditDialogClosed(applied)
          )}
        {isDeleteDialogOpen &&
          renderDeleteDialog &&
          renderDeleteDialog(getSelectedData(items, selectedIds), (deleted) =>
            onDeleteDialogClosed(deleted)
          )}
        {renderPagination()}
        {renderPaginationIndex()}
        {maybeRenderContent()}
      </div>
    );
  };

  const ItemList: React.FC<IItemListProps<T, E>> = (props) => {
    const {
      persistState,
      persistanceKey = filterMode,
      defaultSort = filterOptions.defaultSortBy,
      defaultZoomIndex,
      alterQuery = true,
    } = props;

    const history = useHistory();
    const location = useLocation();
    const [interfaceState, setInterfaceState] = useInterfaceLocalForage();
    const [filterInitialised, setFilterInitialised] = useState(false);
    const { configuration: config } = useContext(ConfigurationContext);

    const lastPathname = useRef(location.pathname);
    const defaultDisplayMode = filterOptions.displayModeOptions[0];
    const [filter, setFilter] = useState<ListFilterModel>(
      () => new ListFilterModel(filterMode)
    );

    const updateSavedFilter = useCallback(
      (updatedFilter: ListFilterModel) => {
        setInterfaceState((prevState) => {
          if (!prevState.queryConfig) {
            prevState.queryConfig = {};
          }

          const oldFilter = prevState.queryConfig[persistanceKey]?.filter ?? "";
          const newFilter = new URLSearchParams(oldFilter);
          newFilter.set("disp", String(updatedFilter.displayMode));

          return {
            ...prevState,
            queryConfig: {
              ...prevState.queryConfig,
              [persistanceKey]: {
                ...prevState.queryConfig[persistanceKey],
                filter: newFilter.toString(),
              },
            },
          };
        });
      },
      [persistanceKey, setInterfaceState]
    );

    const { data: defaultFilter, loading: defaultFilterLoading } =
      useFindDefaultFilter(filterMode);

    const updateQueryParams = useCallback(
      (newFilter: ListFilterModel) => {
        if (!alterQuery) return;

        const newParams = newFilter.makeQueryParameters();
        history.replace({ ...history.location, search: newParams });
      },
      [alterQuery, history]
    );

    const updateFilter = useCallback(
      (newFilter: ListFilterModel) => {
        setFilter(newFilter);
        updateQueryParams(newFilter);
        if (persistState === PersistanceLevel.VIEW) {
          updateSavedFilter(newFilter);
        }
      },
      [persistState, updateSavedFilter, updateQueryParams]
    );

    // 'Startup' hook, initialises the filters
    useEffect(() => {
      // Only run once
      if (filterInitialised) return;

      let newFilter = new ListFilterModel(filterMode, config, defaultZoomIndex);
      let loadDefault = true;
      if (alterQuery && location.search) {
        loadDefault = false;
        newFilter.configureFromQueryString(location.search);
      }

      if (persistState === PersistanceLevel.ALL) {
        // only set default filter if uninitialised
        if (loadDefault) {
          // wait until default filter is loaded
          if (defaultFilterLoading) return;

          if (defaultFilter?.findDefaultFilter) {
            newFilter.currentPage = 1;
            try {
              newFilter.configureFromSavedFilter(
                defaultFilter.findDefaultFilter
              );
            } catch (err) {
              console.log(err);
              // ignore
            }
            // #1507 - reset random seed when loaded
            newFilter.randomSeed = -1;
          }
        }
      } else if (persistState === PersistanceLevel.VIEW) {
        // wait until forage is initialised
        if (interfaceState.loading) return;

        const storedQuery = interfaceState.data?.queryConfig?.[persistanceKey];
        if (persistState === PersistanceLevel.VIEW && storedQuery) {
          const displayMode = new URLSearchParams(storedQuery.filter).get(
            "disp"
          );
          if (displayMode) {
            newFilter.displayMode = Number.parseInt(displayMode, 10);
          }
        }
      }
      setFilter(newFilter);
      updateQueryParams(newFilter);

      setFilterInitialised(true);
    }, [
      filterInitialised,
      location,
      config,
      defaultSort,
      defaultDisplayMode,
      defaultZoomIndex,
      alterQuery,
      persistState,
      updateQueryParams,
      defaultFilter,
      defaultFilterLoading,
      interfaceState,
      persistanceKey,
    ]);

    // This hook runs on every page location change (ie navigation),
    // and updates the filter accordingly.
    useEffect(() => {
      if (!filterInitialised || !alterQuery) return;

      // re-init if the pathname has changed
      if (location.pathname !== lastPathname.current) {
        lastPathname.current = location.pathname;
        setFilterInitialised(false);
        return;
      }

      // re-init to load default filter on empty new query params
      if (!location.search) {
        setFilterInitialised(false);
        return;
      }

      // the query has changed, update filter if necessary
      setFilter((prevFilter) => {
        let newFilter = prevFilter.clone();
        newFilter.configureFromQueryString(location.search);
        if (!isEqual(newFilter, prevFilter)) {
          return newFilter;
        } else {
          return prevFilter;
        }
      });
    }, [filterInitialised, alterQuery, location]);

    const onChangePage = useCallback(
      (page: number) => {
        const newFilter = cloneDeep(filter);
        newFilter.currentPage = page;
        updateFilter(newFilter);

        // if the current page has a detail-header, then
        // scroll up relative to that rather than 0, 0
        const detailHeader = document.querySelector(".detail-header");
        if (detailHeader) {
          window.scrollTo(0, detailHeader.scrollHeight - 50);
        } else {
          window.scrollTo(0, 0);
        }
      },
      [filter, updateFilter]
    );

    if (!filterInitialised) return null;

    return (
      <RenderList
        filter={filter}
        onChangePage={onChangePage}
        updateFilter={updateFilter}
        {...props}
      />
    );
  };

  return ItemList;
}

export const showWhenSelected = <T extends QueryResult>(
  result: T,
  filter: ListFilterModel,
  selectedIds: Set<string>
) => {
  return selectedIds.size > 0;
};
