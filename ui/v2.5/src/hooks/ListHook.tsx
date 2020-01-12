import _ from "lodash";
import queryString from "query-string";
import React, { useState } from "react";
import { Spinner } from 'react-bootstrap';
import { QueryHookResult } from "react-apollo-hooks";
import { ApolloError } from 'apollo-client';
import { useHistory } from 'react-router-dom';
import {
  FindScenesQuery,
  FindScenesVariables,
  SlimSceneDataFragment,
  FindSceneMarkersQuery,
  FindSceneMarkersVariables,
  FindSceneMarkersSceneMarkers,
  FindGalleriesQuery,
  FindGalleriesVariables,
  GalleryDataFragment,
  FindStudiosQuery,
  FindStudiosVariables,
  StudioDataFragment,
  FindPerformersQuery,
  FindPerformersVariables,
  PerformerDataFragment
} from 'src/core/generated-graphql';
import { ListFilter } from "../components/list/ListFilter";
import { Pagination } from "../components/list/Pagination";
import { StashService } from "../core/StashService";
import { Criterion } from "../models/list-filter/criteria/criterion";
import { ListFilterModel } from "../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../models/list-filter/types";

interface IListHookData {
  filter: ListFilterModel;
  template: JSX.Element;
  onSelectChange: (id: string, selected : boolean, shiftKey: boolean) => void;
}

interface IListHookOperation<T> {
  text: string;
  onClick: (result: T, filter: ListFilterModel, selectedIds: Set<string>) => void;
}

interface IListHookOptions<T> {
  zoomable?: boolean;
  otherOperations?: IListHookOperation<T>[];
  renderContent: (result: T, filter: ListFilterModel, selectedIds: Set<string>, zoomIndex: number) => JSX.Element | undefined;
  renderSelectedOptions?: (result: T, selectedIds: Set<string>) => JSX.Element | undefined;
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

type ScenesQuery = QueryHookResult<FindScenesQuery, FindScenesVariables>;
export const useScenesList = (props:IListHookOptions<ScenesQuery>) => (
  useList<ScenesQuery, SlimSceneDataFragment>({
    ...props,
    filterMode: FilterMode.Scenes,
    useData: StashService.useFindScenes,
    getData: (result:ScenesQuery) => (result?.data?.findScenes?.scenes ?? []),
    getCount: (result:ScenesQuery) => (result?.data?.findScenes?.count ?? 0)
  })
)

type SceneMarkersQuery = QueryHookResult<FindSceneMarkersQuery, FindSceneMarkersVariables>;
export const useSceneMarkersList = (props:IListHookOptions<SceneMarkersQuery>) => (
  useList<SceneMarkersQuery, FindSceneMarkersSceneMarkers>({
    ...props,
    filterMode: FilterMode.SceneMarkers,
    useData: StashService.useFindSceneMarkers,
    getData: (result:SceneMarkersQuery) => (result?.data?.findSceneMarkers?.scene_markers?? []),
    getCount: (result:SceneMarkersQuery) => (result?.data?.findSceneMarkers?.count ?? 0)
  })
)

type GalleriesQuery = QueryHookResult<FindGalleriesQuery, FindGalleriesVariables>;
export const useGalleriesList = (props:IListHookOptions<GalleriesQuery>) => (
  useList<GalleriesQuery, GalleryDataFragment>({
    ...props,
    filterMode: FilterMode.Galleries,
    useData: StashService.useFindGalleries,
    getData: (result:GalleriesQuery) => (result?.data?.findGalleries?.galleries ?? []),
    getCount: (result:GalleriesQuery) => (result?.data?.findGalleries?.count ?? 0)
  })
)

type StudiosQuery = QueryHookResult<FindStudiosQuery, FindStudiosVariables>;
export const useStudiosList = (props:IListHookOptions<StudiosQuery>) => (
  useList<StudiosQuery, StudioDataFragment>({
    ...props,
    filterMode: FilterMode.Studios,
    useData: StashService.useFindStudios,
    getData: (result:StudiosQuery) => (result?.data?.findStudios?.studios ?? []),
    getCount: (result:StudiosQuery) => (result?.data?.findStudios?.count ?? 0)
  })
)

type PerformersQuery = QueryHookResult<FindPerformersQuery, FindPerformersVariables>;
export const usePerformersList = (props:IListHookOptions<PerformersQuery>) => (
  useList<PerformersQuery, PerformerDataFragment>({
    ...props,
    filterMode: FilterMode.Performers,
    useData: StashService.useFindPerformers,
    getData: (result:PerformersQuery) => (result?.data?.findPerformers?.performers ?? []),
    getCount: (result:PerformersQuery) => (result?.data?.findPerformers?.count ?? 0)
  })
)

const useList = <QueryResult extends IQueryResult, QueryData extends IDataItem>(
  options: (IListHookOptions<QueryResult> & IQuery<QueryResult, QueryData>)
): IListHookData => {
  const history = useHistory();
  const [filter, setFilter] = useState<ListFilterModel>(new ListFilterModel(options.filterMode, queryString.parse(history.location.search)));
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [lastClickedId, setLastClickedId] = useState<string | undefined>();
  const [zoomIndex, setZoomIndex] = useState<number>(1);

  const result = options.useData(filter);
  const totalCount = options.getCount(result);
  const items = options.getData(result);

  function updateQueryParams(filter:ListFilterModel) {
    const newLocation = Object.assign({}, history.location);
    newLocation.search = filter.makeQueryParameters();
    history.replace(newLocation);
  }

  function onChangePageSize(pageSize: number) {
    const newFilter = _.cloneDeep(filter);
    newFilter.itemsPerPage = pageSize;
    newFilter.currentPage = 1;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onChangeQuery(query: string) {
    const newFilter = _.cloneDeep(filter);
    newFilter.searchTerm = query;
    newFilter.currentPage = 1;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onChangeSortDirection(sortDirection: "asc" | "desc") {
    const newFilter = _.cloneDeep(filter);
    newFilter.sortDirection = sortDirection;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onChangeSortBy(sortBy: string) {
    const newFilter = _.cloneDeep(filter);
    newFilter.sortBy = sortBy;
    newFilter.currentPage = 1;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = _.cloneDeep(filter);
    newFilter.displayMode = displayMode;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onAddCriterion(criterion: Criterion, oldId?: string) {
    const newFilter = _.cloneDeep(filter);

    // Find if we are editing an existing criteria, then modify that.  Or create a new one.
    const existingIndex = newFilter.criteria.findIndex((c) => {
      // If we modified an existing criterion, then look for the old id.
      const id = !!oldId ? oldId : criterion.getId();
      return c.getId() === id;
    });
    if (existingIndex === -1) {
      newFilter.criteria.push(criterion);
    } else {
      newFilter.criteria[existingIndex] = criterion;
    }

    // Remove duplicate modifiers
    newFilter.criteria = newFilter.criteria.filter((obj, pos, arr) => {
      return arr.map((mapObj: any) => mapObj.getId()).indexOf(obj.getId()) === pos;
    });

    newFilter.currentPage = 1;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onRemoveCriterion(removedCriterion: Criterion) {
    const newFilter = _.cloneDeep(filter);
    newFilter.criteria = newFilter.criteria.filter((criterion) => criterion.getId() !== removedCriterion.getId());
    newFilter.currentPage = 1;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onChangePage(page: number) {
    const newFilter = _.cloneDeep(filter);
    newFilter.currentPage = page;
    setFilter(newFilter);
    updateQueryParams(newFilter);
  }

  function onSelectChange(id: string, selected : boolean, shiftKey: boolean) {
    if (shiftKey) {
      multiSelect(id);
    } else {
      singleSelect(id, selected);
    }
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

  function selectRange(startIndex : number, endIndex : number) {
    if (startIndex > endIndex) {
      let tmp = startIndex;
      startIndex = endIndex;
      endIndex = tmp;
    }

    const subset = items.slice(startIndex, endIndex + 1);
    const newSelectedIds : Set<string> = new Set();

    subset.forEach((item) => {
      newSelectedIds.add(item.id);
    });

    setSelectedIds(newSelectedIds);
  }

  function onSelectAll() {
    const newSelectedIds : Set<string> = new Set();
    items.forEach((item) => {
      newSelectedIds.add(item.id);
    });

    setSelectedIds(newSelectedIds);
    setLastClickedId(undefined);
  }

  function onSelectNone() {
    const newSelectedIds : Set<string> = new Set();
    setSelectedIds(newSelectedIds);
    setLastClickedId(undefined);
  }

  function onChangeZoom(newZoomIndex : number) {
    setZoomIndex(newZoomIndex);
  }

  const otherOperations = options.otherOperations ? options.otherOperations.map((o) => {
    return {
      text: o.text,
      onClick: () => {
        o.onClick(result, filter, selectedIds);
      }
    }
  }) : undefined;

  const template = (
    <div>
      <ListFilter
        onChangePageSize={onChangePageSize}
        onChangeQuery={onChangeQuery}
        onChangeSortDirection={onChangeSortDirection}
        onChangeSortBy={onChangeSortBy}
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
      {options.renderSelectedOptions && selectedIds.size > 0 ? options.renderSelectedOptions(result, selectedIds) : undefined}
      {result.loading ? <Spinner animation="border" variant="light" /> : undefined}
      {result.error ? <h1>{result.error.message}</h1> : undefined}
      {options.renderContent(result, filter, selectedIds, zoomIndex)}
      <Pagination
        itemsPerPage={filter.itemsPerPage}
        currentPage={filter.currentPage}
        totalItems={totalCount}
        onChangePage={onChangePage}
        loading={result.loading}
      />
    </div>
  );

  return { filter, template, onSelectChange };
}
