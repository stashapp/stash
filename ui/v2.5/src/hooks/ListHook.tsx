import _ from "lodash";
import queryString from "query-string";
import React, { useCallback, useState, useEffect } from "react";
import { ApolloError } from "apollo-client";
import { useHistory, useLocation } from "react-router-dom";
import {
  SortDirectionEnum,
  SlimSceneDataFragment,
  SceneMarkerDataFragment,
  GalleryDataFragment,
  StudioDataFragment,
  PerformerDataFragment,
  FindScenesQueryResult,
  FindSceneMarkersQueryResult,
  FindGalleriesQueryResult,
  FindStudiosQueryResult,
  FindPerformersQueryResult,
  FindMoviesQueryResult,
  MovieDataFragment,
} from "src/core/generated-graphql";
import {
  useInterfaceLocalForage,
  IInterfaceConfig,
} from "src/hooks/LocalForage";
import { LoadingIndicator } from "src/components/Shared";
import { ListFilter } from "src/components/List/ListFilter";
import { Pagination, PaginationIndex } from "src/components/List/Pagination";
import {
  useFindScenes,
  useFindSceneMarkers,
  useFindMovies,
  useFindStudios,
  useFindGalleries,
  useFindPerformers,
} from "src/core/StashService";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode, FilterMode } from "src/models/list-filter/types";

interface IListHookData {
  filter: ListFilterModel;
  template: JSX.Element;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

interface IListHookOperation<T> {
  text: string;
  onClick: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => void;
}

interface IListHookOptions<T> {
  subComponent?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  zoomable?: boolean;
  otherOperations?: IListHookOperation<T>[];
  renderContent: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) => JSX.Element | undefined;
  renderSelectedOptions?: (
    result: T,
    selectedIds: Set<string>
  ) => JSX.Element | undefined;
}

interface IDataItem {
  id: string;
}
interface IQueryResult {
  error?: ApolloError;
  loading: boolean;
}

interface IQuery<T extends IQueryResult, T2 extends IDataItem> {
  filterMode: FilterMode;
  useData: (filter: ListFilterModel) => T;
  getData: (data: T) => T2[];
  getCount: (data: T) => number;
}

const useList = <QueryResult extends IQueryResult, QueryData extends IDataItem>(
  options: IListHookOptions<QueryResult> & IQuery<QueryResult, QueryData>
): IListHookData => {
  const [interfaceState, setInterfaceState] = useInterfaceLocalForage();
  const [forageInitialised, setForageInitialised] = useState(false);
  const history = useHistory();
  const location = useLocation();
  const [filter, setFilter] = useState<ListFilterModel>(
    new ListFilterModel(
      options.filterMode,
      options.subComponent ? undefined : queryString.parse(location.search)
    )
  );
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [lastClickedId, setLastClickedId] = useState<string | undefined>();
  const [zoomIndex, setZoomIndex] = useState<number>(1);

  const result = options.useData(getFilter());
  const totalCount = options.getCount(result);
  const items = options.getData(result);

  const updateInterfaceConfig = useCallback(
    (updatedFilter: ListFilterModel) => {
      setInterfaceState((config) => {
        const data = { ...config } as IInterfaceConfig;
        data.queries = {
          [options.filterMode]: {
            filter: updatedFilter.makeQueryParameters(),
            itemsPerPage: updatedFilter.itemsPerPage,
            currentPage: updatedFilter.currentPage,
          },
        };
        return data;
      });
    },
    [options.filterMode, setInterfaceState]
  );

  useEffect(() => {
    if (interfaceState.loading) return;
    if (!forageInitialised) setForageInitialised(true);

    // Don't use query parameters for sub-components
    if (options.subComponent) return;

    const storedQuery = interfaceState.data?.queries?.[options.filterMode];
    if (!storedQuery) return;

    const queryFilter = queryString.parse(history.location.search);
    const storedFilter = queryString.parse(storedQuery.filter);
    const query = history.location.search
      ? {
          sortby: storedFilter.sortby,
          sortdir: storedFilter.sortdir,
          disp: storedFilter.disp,
          perPage: storedFilter.perPage,
          ...queryFilter,
        }
      : storedFilter;

    const newFilter = new ListFilterModel(options.filterMode, query);

    // Compare constructed filter with current filter.
    // If different it's the result of navigation, and we update the filter.
    const newLocation = { ...history.location };
    newLocation.search = newFilter.makeQueryParameters();
    if (newLocation.search !== filter.makeQueryParameters()) {
      setFilter(newFilter);
      updateInterfaceConfig(newFilter);
    }
    // If constructed search is different from current, update it as well
    if (newLocation.search !== location.search) {
      newLocation.search = newFilter.makeQueryParameters();
      history.replace(newLocation);
    }
  }, [
    filter,
    interfaceState.data,
    interfaceState.loading,
    history,
    location.search,
    options.subComponent,
    options.filterMode,
    forageInitialised,
    updateInterfaceConfig,
  ]);

  function getFilter() {
    if (!options.filterHook) {
      return filter;
    }

    // make a copy of the filter and call the hook
    const newFilter = _.cloneDeep(filter);
    return options.filterHook(newFilter);
  }

  function updateQueryParams(listFilter: ListFilterModel) {
    setFilter(listFilter);
    if (!options.subComponent) {
      const newLocation = { ...location };
      newLocation.search = listFilter.makeQueryParameters();
      history.replace(newLocation);
      updateInterfaceConfig(listFilter);
    }
  }

  function onChangePageSize(pageSize: number) {
    const newFilter = _.cloneDeep(filter);
    newFilter.itemsPerPage = pageSize;
    newFilter.currentPage = 1;
    updateQueryParams(newFilter);
  }

  function onChangeQuery(query: string) {
    const newFilter = _.cloneDeep(filter);
    newFilter.searchTerm = query;
    newFilter.currentPage = 1;
    updateQueryParams(newFilter);
  }

  function onChangeSortDirection(sortDirection: SortDirectionEnum) {
    const newFilter = _.cloneDeep(filter);
    newFilter.sortDirection = sortDirection;
    updateQueryParams(newFilter);
  }

  function onChangeSortBy(sortBy: string) {
    const newFilter = _.cloneDeep(filter);
    newFilter.sortBy = sortBy;
    newFilter.currentPage = 1;
    updateQueryParams(newFilter);
  }

  function onSortReshuffle() {
    const newFilter = _.cloneDeep(filter);
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
    updateQueryParams(newFilter);
  }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = _.cloneDeep(filter);
    newFilter.displayMode = displayMode;
    updateQueryParams(newFilter);
  }

  function onAddCriterion(criterion: Criterion, oldId?: string) {
    const newFilter = _.cloneDeep(filter);

    // Find if we are editing an existing criteria, then modify that.  Or create a new one.
    const existingIndex = newFilter.criteria.findIndex((c) => {
      // If we modified an existing criterion, then look for the old id.
      const id = oldId || criterion.getId();
      return c.getId() === id;
    });
    if (existingIndex === -1) {
      newFilter.criteria.push(criterion);
    } else {
      newFilter.criteria[existingIndex] = criterion;
    }

    // Remove duplicate modifiers
    newFilter.criteria = newFilter.criteria.filter((obj, pos, arr) => {
      return arr.map((mapObj) => mapObj.getId()).indexOf(obj.getId()) === pos;
    });

    newFilter.currentPage = 1;
    updateQueryParams(newFilter);
  }

  function onRemoveCriterion(removedCriterion: Criterion) {
    const newFilter = _.cloneDeep(filter);
    newFilter.criteria = newFilter.criteria.filter(
      (criterion) => criterion.getId() !== removedCriterion.getId()
    );
    newFilter.currentPage = 1;
    updateQueryParams(newFilter);
  }

  function onChangePage(page: number) {
    const newFilter = _.cloneDeep(filter);
    newFilter.currentPage = page;
    updateQueryParams(newFilter);
  }

  function singleSelect(id: string, selected: boolean) {
    setLastClickedId(id);

    const newSelectedIds = _.clone(selectedIds);
    if (selected) {
      newSelectedIds.add(id);
    } else {
      newSelectedIds.delete(id);
    }

    setSelectedIds(newSelectedIds);
  }

  function selectRange(startIndex: number, endIndex: number) {
    let start = startIndex;
    let end = endIndex;
    if (start > end) {
      const tmp = start;
      start = end;
      end = tmp;
    }

    const subset = items.slice(start, end + 1);
    const newSelectedIds: Set<string> = new Set();

    subset.forEach((item) => {
      newSelectedIds.add(item.id);
    });

    setSelectedIds(newSelectedIds);
  }

  function multiSelect(id: string) {
    let startIndex = 0;
    let thisIndex = -1;

    if (lastClickedId) {
      startIndex = items.findIndex((item) => {
        return item.id === lastClickedId;
      });
    }

    thisIndex = items.findIndex((item) => {
      return item.id === id;
    });

    selectRange(startIndex, thisIndex);
  }

  function onSelectChange(id: string, selected: boolean, shiftKey: boolean) {
    if (shiftKey) {
      multiSelect(id);
    } else {
      singleSelect(id, selected);
    }
  }

  function onSelectAll() {
    const newSelectedIds: Set<string> = new Set();
    items.forEach((item) => {
      newSelectedIds.add(item.id);
    });

    setSelectedIds(newSelectedIds);
    setLastClickedId(undefined);
  }

  function onSelectNone() {
    const newSelectedIds: Set<string> = new Set();
    setSelectedIds(newSelectedIds);
    setLastClickedId(undefined);
  }

  function onChangeZoom(newZoomIndex: number) {
    setZoomIndex(newZoomIndex);
  }

  const otherOperations = options.otherOperations
    ? options.otherOperations.map((o) => {
        return {
          text: o.text,
          onClick: () => {
            o.onClick(result, filter, selectedIds);
          },
        };
      })
    : undefined;

  function maybeRenderContent() {
    if (!result.loading && !result.error) {
      return options.renderContent(result, filter, selectedIds, zoomIndex);
    }
  }

  function maybeRenderPaginationIndex() {
    if (!result.loading && !result.error) {
      return (
        <PaginationIndex
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
        />
      );
    }
  }

  function maybeRenderPagination() {
    if (!result.loading && !result.error) {
      return (
        <Pagination
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
          onChangePage={onChangePage}
        />
      );
    }
  }

  const template = (
    <div>
      <ListFilter
        onChangePageSize={onChangePageSize}
        onChangeQuery={onChangeQuery}
        onChangeSortDirection={onChangeSortDirection}
        onChangeSortBy={onChangeSortBy}
        onSortReshuffle={onSortReshuffle}
        onChangeDisplayMode={onChangeDisplayMode}
        onAddCriterion={onAddCriterion}
        onRemoveCriterion={onRemoveCriterion}
        onSelectAll={onSelectAll}
        onSelectNone={onSelectNone}
        zoomIndex={options.zoomable ? zoomIndex : undefined}
        onChangeZoom={options.zoomable ? onChangeZoom : undefined}
        otherOperations={otherOperations}
        filter={filter}
      />
      {options.renderSelectedOptions && selectedIds.size > 0
        ? options.renderSelectedOptions(result, selectedIds)
        : undefined}
      {(result.loading || !forageInitialised) && <LoadingIndicator />}
      {result.error && <h1>{result.error.message}</h1>}
      {maybeRenderPagination()}
      {maybeRenderContent()}
      {maybeRenderPaginationIndex()}
      {maybeRenderPagination()}
    </div>
  );

  return { filter, template, onSelectChange };
};

export const useScenesList = (props: IListHookOptions<FindScenesQueryResult>) =>
  useList<FindScenesQueryResult, SlimSceneDataFragment>({
    ...props,
    filterMode: FilterMode.Scenes,
    useData: useFindScenes,
    getData: (result: FindScenesQueryResult) =>
      result?.data?.findScenes?.scenes ?? [],
    getCount: (result: FindScenesQueryResult) =>
      result?.data?.findScenes?.count ?? 0,
  });

export const useSceneMarkersList = (
  props: IListHookOptions<FindSceneMarkersQueryResult>
) =>
  useList<FindSceneMarkersQueryResult, SceneMarkerDataFragment>({
    ...props,
    filterMode: FilterMode.SceneMarkers,
    useData: useFindSceneMarkers,
    getData: (result: FindSceneMarkersQueryResult) =>
      result?.data?.findSceneMarkers?.scene_markers ?? [],
    getCount: (result: FindSceneMarkersQueryResult) =>
      result?.data?.findSceneMarkers?.count ?? 0,
  });

export const useGalleriesList = (
  props: IListHookOptions<FindGalleriesQueryResult>
) =>
  useList<FindGalleriesQueryResult, GalleryDataFragment>({
    ...props,
    filterMode: FilterMode.Galleries,
    useData: useFindGalleries,
    getData: (result: FindGalleriesQueryResult) =>
      result?.data?.findGalleries?.galleries ?? [],
    getCount: (result: FindGalleriesQueryResult) =>
      result?.data?.findGalleries?.count ?? 0,
  });

export const useStudiosList = (
  props: IListHookOptions<FindStudiosQueryResult>
) =>
  useList<FindStudiosQueryResult, StudioDataFragment>({
    ...props,
    filterMode: FilterMode.Studios,
    useData: useFindStudios,
    getData: (result: FindStudiosQueryResult) =>
      result?.data?.findStudios?.studios ?? [],
    getCount: (result: FindStudiosQueryResult) =>
      result?.data?.findStudios?.count ?? 0,
  });

export const usePerformersList = (
  props: IListHookOptions<FindPerformersQueryResult>
) =>
  useList<FindPerformersQueryResult, PerformerDataFragment>({
    ...props,
    filterMode: FilterMode.Performers,
    useData: useFindPerformers,
    getData: (result: FindPerformersQueryResult) =>
      result?.data?.findPerformers?.performers ?? [],
    getCount: (result: FindPerformersQueryResult) =>
      result?.data?.findPerformers?.count ?? 0,
  });

export const useMoviesList = (props: IListHookOptions<FindMoviesQueryResult>) =>
  useList<FindMoviesQueryResult, MovieDataFragment>({
    ...props,
    filterMode: FilterMode.Movies,
    useData: useFindMovies,
    getData: (result: FindMoviesQueryResult) =>
      result?.data?.findMovies?.movies ?? [],
    getCount: (result: FindMoviesQueryResult) =>
      result?.data?.findMovies?.count ?? 0,
  });
