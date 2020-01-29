import { Spinner } from "@blueprintjs/core";
import _ from "lodash";
import queryString from "query-string";
import React, { useEffect, useState } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { ListFilter } from "../components/list/ListFilter";
import { Pagination } from "../components/list/Pagination";
import { StashService } from "../core/StashService";
import { IBaseProps } from "../models";
import { Criterion } from "../models/list-filter/criteria/criterion";
import { ListFilterModel } from "../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../models/list-filter/types";

export interface IListHookData {
  filter: ListFilterModel;
  template: JSX.Element;
  options: IListHookOptions;
  onSelectChange: (id: string, selected : boolean, shiftKey: boolean) => void;
}

interface IListHookOperation {
  text: string;
  onClick: (result: QueryHookResult<any, any>, filter: ListFilterModel, selectedIds: Set<string>) => void;
}

export interface IListHookOptions {
  filterMode: FilterMode;
  subComponent?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  props: IBaseProps;
  zoomable?: boolean;
  otherOperations?: IListHookOperation[];
  renderContent: (result: QueryHookResult<any, any>, filter: ListFilterModel, selectedIds: Set<string>, zoomIndex: number) => JSX.Element | undefined;
  renderSelectedOptions?: (result: QueryHookResult<any, any>, selectedIds: Set<string>) => JSX.Element | undefined;
}

export class ListHook {
  public static useList(options: IListHookOptions): IListHookData {
    const [filter, setFilter] = useState<ListFilterModel>(new ListFilterModel(options.filterMode));
    const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
    const [lastClickedId, setLastClickedId] = useState<string | undefined>(undefined);
    const [totalCount, setTotalCount] = useState<number>(0);
    const [zoomIndex, setZoomIndex] = useState<number>(1);

    // Update the filter when the query parameters change
    // don't use query parameters for sub-components
    if (!options.subComponent) {
      useEffect(() => {
        const queryParams = queryString.parse(options.props!.location.search);
        const newFilter = _.cloneDeep(filter);
        newFilter.configureFromQueryParameters(queryParams);
        setFilter(newFilter);

        // TODO: Need this side effect to update the query params properly
        filter.configureFromQueryParameters(queryParams);
      }, [options.props.location.search]);
    }

    function getFilter() {
      if (!options.filterHook) {
        return filter;
      }

      // make a copy of the filter and call the hook
      let newFilter = _.cloneDeep(filter);
      return options.filterHook(newFilter);
    }

    let result: QueryHookResult<any, any>;

    let getData: (filter : ListFilterModel) => QueryHookResult<any, any>;
    let getItems: () => any[];
    let getCount: () => number;

    switch (options.filterMode) {
      case FilterMode.Scenes: {
        getData = (filter : ListFilterModel) => { return StashService.useFindScenes(filter); }
        getItems = () => { return !!result.data && !!result.data.findScenes ? result.data.findScenes.scenes : []; }
        getCount = () => { return !!result.data && !!result.data.findScenes ? result.data.findScenes.count : 0; }
        break;
      }
      case FilterMode.SceneMarkers: {
        getData = (filter : ListFilterModel) => { return StashService.useFindSceneMarkers(filter); }
        getItems = () => { return !!result.data && !!result.data.findSceneMarkers ? result.data.findSceneMarkers.scene_markers : []; }
        getCount = () => { return !!result.data && !!result.data.findSceneMarkers ? result.data.findSceneMarkers.count : 0; }
        break;
      }
      case FilterMode.Galleries: {
        getData = (filter : ListFilterModel) => { return StashService.useFindGalleries(filter); }
        getItems = () => { return !!result.data && !!result.data.findGalleries ? result.data.findGalleries.galleries : []; }
        getCount = () => { return !!result.data && !!result.data.findGalleries ? result.data.findGalleries.count : 0; }
        break;
      }
      case FilterMode.Studios: {
        getData = (filter : ListFilterModel) => { return StashService.useFindStudios(filter); }
        getItems = () => { return !!result.data && !!result.data.findStudios ? result.data.findStudios.studios : []; }
        getCount = () => { return !!result.data && !!result.data.findStudios ? result.data.findStudios.count : 0; }
        break;
      }
      case FilterMode.Movies: {
        getData = (filter : ListFilterModel) => { return StashService.useFindMovies(filter); }
        getItems = () => { return !!result.data && !!result.data.findMovies ? result.data.findMovies.movies : []; }
        getCount = () => { return !!result.data && !!result.data.findMovies ? result.data.findMovies.count : 0; }
        break;
      }

      case FilterMode.Performers: {
        getData = (filter : ListFilterModel) => { return StashService.useFindPerformers(filter); }
        getItems = () => { return !!result.data && !!result.data.findPerformers ? result.data.findPerformers.performers : []; }
        getCount = () => { return !!result.data && !!result.data.findPerformers ? result.data.findPerformers.count : 0; }
        break;
      }
      default: {
        console.error("REMOVE DEFAULT IN LIST HOOK");
        getData = (filter : ListFilterModel) => { return StashService.useFindScenes(filter); }
        getItems = () => { return !!result.data && !!result.data.findScenes ? result.data.findScenes.scenes : []; }
        getCount = () => { return !!result.data && !!result.data.findScenes ? result.data.findScenes.count : 0; }
        break;
      }
    }

    result = getData(getFilter());

    useEffect(() => {
      setTotalCount(getCount());

      // select none when data changes
      onSelectNone();
      setLastClickedId(undefined);
    }, [result.data])

    // Update the query parameters when the data changes
    // don't use query parameters for sub-components
    if (!options.subComponent) {
      useEffect(() => {
        const location = Object.assign({}, options.props.history.location);
        location.search = filter.makeQueryParameters();
        options.props.history.replace(location);
      }, [result.data, filter.displayMode]);
    }

    // Update the total count
    useEffect(() => {
      const newFilter = _.cloneDeep(filter);
      newFilter.totalCount = totalCount;
      setFilter(newFilter);
    }, [totalCount]);

    function onChangePageSize(pageSize: number) {
      const newFilter = _.cloneDeep(filter);
      newFilter.itemsPerPage = pageSize;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    }

    function onChangeQuery(query: string) {
      const newFilter = _.cloneDeep(filter);
      newFilter.searchTerm = query;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    }

    function onChangeSortDirection(sortDirection: "asc" | "desc") {
      const newFilter = _.cloneDeep(filter);
      newFilter.sortDirection = sortDirection;
      setFilter(newFilter);
    }

    function onChangeSortBy(sortBy: string) {
      const newFilter = _.cloneDeep(filter);
      newFilter.sortBy = sortBy;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    }

    function onChangeDisplayMode(displayMode: DisplayMode) {
      const newFilter = _.cloneDeep(filter);
      newFilter.displayMode = displayMode;
      setFilter(newFilter);
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
    }

    function onRemoveCriterion(removedCriterion: Criterion) {
      const newFilter = _.cloneDeep(filter);
      newFilter.criteria = newFilter.criteria.filter((criterion) => criterion.getId() !== removedCriterion.getId());
      newFilter.currentPage = 1;
      setFilter(newFilter);
    }

    function onChangePage(page: number) {
      const newFilter = _.cloneDeep(filter);
      newFilter.currentPage = page;
      setFilter(newFilter);
    }

    function onSelectChange(id: string, selected : boolean, shiftKey: boolean) {
      if (shiftKey) {
        multiSelect(id, selected);
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

    function multiSelect(id: string, selected : boolean) {
      let startIndex = 0;
      let thisIndex = -1;
  
      if (!!lastClickedId) {
        startIndex = getItems().findIndex((item) => {
          return item.id === lastClickedId;
        });
      }

      thisIndex = getItems().findIndex((item) => {
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
  
      const subset = getItems().slice(startIndex, endIndex + 1);
      const newSelectedIds : Set<string> = new Set();

      subset.forEach((item) => {
        newSelectedIds.add(item.id);
      });

      setSelectedIds(newSelectedIds);
    }

    function onSelectAll() {
      const newSelectedIds : Set<string> = new Set();
      getItems().forEach((item) => {
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
        {result.loading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
        {result.error ? <h1>{result.error.message}</h1> : undefined}
        {options.renderContent(result, filter, selectedIds, zoomIndex)}
        <Pagination
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
          onChangePage={onChangePage}
        />
      </div>
    );

    return { filter, template, options, onSelectChange };
  }
}
