import { Spinner } from "@blueprintjs/core";
import _ from "lodash";
import queryString from "query-string";
import React, { useEffect, useState, useRef } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { ListFilter } from "../components/list/ListFilter";
import { Pagination } from "../components/list/Pagination";
import { StashService } from "../core/StashService";
import { IBaseProps } from "../models";
import { Criterion } from "../models/list-filter/criteria/criterion";
import { ListFilterModel } from "../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../models/list-filter/types";
import { useInterfaceLocalForage } from "./LocalForage";

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

export interface IFilterListImpl {
  getData: (filter : ListFilterModel) => QueryHookResult<any, any>;
  getItems: (data: any) => any[];
  getCount: (data: any) => number;
}

const SceneFilterListImpl: IFilterListImpl = {
  getData: (filter : ListFilterModel) => { return StashService.useFindScenes(filter); },
  getItems: (data: any) => { return !!data && !!data.findScenes ? data.findScenes.scenes : []; },
  getCount: (data: any) => { return !!data && !!data.findScenes ? data.findScenes.count : 0; }
}

const SceneMarkerFilterListImpl: IFilterListImpl = {
  getData: (filter : ListFilterModel) => { return StashService.useFindSceneMarkers(filter); },
  getItems: (data: any) => { return !!data && !!data.findSceneMarkers ? data.findSceneMarkers.scene_markers : []; },
  getCount: (data: any) => { return !!data && !!data.findSceneMarkers ? data.findSceneMarkers.count : 0; }
}

const GalleryFilterListImpl: IFilterListImpl = {
  getData: (filter : ListFilterModel) => { return StashService.useFindGalleries(filter); },
  getItems: (data: any) => { return !!data && !!data.findGalleries ? data.findGalleries.galleries : []; },
  getCount: (data: any) => { return !!data && !!data.findGalleries ? data.findGalleries.count : 0; }
}

const StudioFilterListImpl: IFilterListImpl = {
  getData: (filter : ListFilterModel) => { return StashService.useFindStudios(filter); },
  getItems: (data: any) => { return !!data && !!data.findStudios ? data.findStudios.studios : []; },
  getCount: (data: any) => { return !!data && !!data.findStudios ? data.findStudios.count : 0; }
}

const PerformerFilterListImpl: IFilterListImpl = {
  getData: (filter : ListFilterModel) => { return StashService.useFindPerformers(filter); },
  getItems: (data: any) => { return !!data && !!data.findPerformers ? data.findPerformers.performers : []; },
  getCount: (data: any) => { return !!data && !!data.findPerformers ? data.findPerformers.count : 0; }
}

const MoviesFilterListImpl: IFilterListImpl = {
  getData: (filter : ListFilterModel) => { return StashService.useFindMovies(filter); },
  getItems: (data: any) => { return !!data && !!data.findMovies ? data.findMovies.movies : []; },
  getCount: (data: any) => { return !!data && !!data.findMovies ? data.findMovies.count : 0; }
}


function getFilterListImpl(filterMode: FilterMode) {
  switch (filterMode) {
    case FilterMode.Scenes: {
      return SceneFilterListImpl;
    }
    case FilterMode.SceneMarkers: {
      return SceneMarkerFilterListImpl;
    }
    case FilterMode.Galleries: {
      return GalleryFilterListImpl;
    }
    case FilterMode.Studios: {
      return StudioFilterListImpl;
    }
    case FilterMode.Performers: {
      return PerformerFilterListImpl;
    }
    case FilterMode.Movies: {
      return MoviesFilterListImpl;
    }
    default: {
      console.error("REMOVE DEFAULT IN LIST HOOK");
      return SceneFilterListImpl;
    }
  }
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

function updateFromQueryString(queryStr: string, setFilter: (value: React.SetStateAction<ListFilterModel>) => void, forageData?: any) {
  const queryParams = queryString.parse(queryStr);
  setFilter((f) => {
    const newFilter = _.cloneDeep(f);
    newFilter.configureFromQueryParameters(queryParams);

    if (forageData) {
      const forageParams = queryString.parse(forageData.filter);
      newFilter.overridePrefs(queryParams, forageParams);      
    }

    return newFilter;
  });
}

export class ListHook {
  public static useList(options: IListHookOptions): IListHookData {
    const [filter, setFilter] = useState<ListFilterModel>(new ListFilterModel(options.filterMode));
    const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
    const [lastClickedId, setLastClickedId] = useState<string | undefined>(undefined);
    const [totalCount, setTotalCount] = useState<number>(0);
    const [zoomIndex, setZoomIndex] = useState<number>(1);

    const [interfaceForage, setInterfaceForage] = useInterfaceLocalForage();
    const forageInitialised = useRef<boolean>(false);

    const filterListImpl = getFilterListImpl(options.filterMode);

    // Initialise from interface forage when loaded
    useEffect(() => {
      function updateFromLocalForage(queryData: any) {
        const queryParams = queryString.parse(queryData.filter);
        
        setFilter((f) => {
          const newFilter = _.cloneDeep(f);
          newFilter.configureFromQueryParameters(queryParams);
          newFilter.currentPage = queryData.currentPage;
          newFilter.itemsPerPage = queryData.itemsPerPage;
          return newFilter;
        });
      }

      function initialise() {
        forageInitialised.current = true;

        let forageData: any;

        if (interfaceForage.data && interfaceForage.data.queries[options.filterMode]) {
          forageData = interfaceForage.data.queries[options.filterMode];
        }

        if (!options.props!.location.search && forageData) {
          // we have some data, try to load it
          updateFromLocalForage(forageData);
        } else {
          // use query string instead - include the forageData to include the following
          // preferences if not specified: displayMode, itemsPerPage, sortBy and sortDir
          updateFromQueryString(options.props!.location.search, setFilter, forageData);
        }
      }

      // don't use query parameters for sub-components
      if (!options.subComponent) {
        // initialise once when the forage is loaded
        if (!forageInitialised.current && !interfaceForage.loading) {
          initialise();
          return;
        }
      }
    }, [interfaceForage.data, interfaceForage.loading, options.props, options.filterMode, options.subComponent]);

    // Update the filter when the query parameters change
    useEffect(() => {
      // don't use query parameters for sub-components
      if (!options.subComponent) {
        // only update from the URL if the forage is initialised
        if (forageInitialised.current) {
          updateFromQueryString(options.props!.location.search, setFilter);
        }
      }
    }, [options.props, options.filterMode, options.subComponent]);

    function getFilter() {
      if (!options.filterHook) {
        return filter;
      }

      // make a copy of the filter and call the hook
      let newFilter = _.cloneDeep(filter);
      return options.filterHook(newFilter);
    }

    const result = filterListImpl.getData(getFilter());

    useEffect(() => {
      setTotalCount(filterListImpl.getCount(result.data));

      // select none when data changes
      onSelectNone();
      setLastClickedId(undefined);
    }, [result.data, filterListImpl])

    // Update the query parameters when the data changes
    
    useEffect(() => {
      // don't use query parameters for sub-components
      if (!options.subComponent) {
        // don't update this until local forage is loaded
        if (forageInitialised.current) {
          const location = Object.assign({}, options.props.history.location);
          const includePrefs = true;
          location.search = "?" + filter.makeQueryParameters(includePrefs);

          if (location.search !== options.props.history.location.search) {
            options.props.history.replace(location);
          }

          setInterfaceForage((d) => {
            const dataClone = _.cloneDeep(d);
            dataClone!.queries[options.filterMode] = {
              filter: location.search,
              itemsPerPage: filter.itemsPerPage,
              currentPage: filter.currentPage
            };
            return dataClone;
          });
        }
      }
    }, [result.data, filter, options.subComponent, options.filterMode, options.props.history, setInterfaceForage]);

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
        startIndex = filterListImpl.getItems(result.data).findIndex((item) => {
          return item.id === lastClickedId;
        });
      }

      thisIndex = filterListImpl.getItems(result.data).findIndex((item) => {
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
  
      const subset = filterListImpl.getItems(result.data).slice(startIndex, endIndex + 1);
      const newSelectedIds : Set<string> = new Set();

      subset.forEach((item) => {
        newSelectedIds.add(item.id);
      });

      setSelectedIds(newSelectedIds);
    }

    function onSelectAll() {
      const newSelectedIds : Set<string> = new Set();
      filterListImpl.getItems(result.data).forEach((item) => {
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

    function maybeRenderContent() {
      if (!result.loading && !result.error) {
        return options.renderContent(result, filter, selectedIds, zoomIndex);
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

    function getTemplate() {
      if (!options.subComponent && !forageInitialised.current) {
        return (
          <div>
            {!result.error ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
            {result.error ? <h1>{result.error.message}</h1> : undefined}
          </div>
        )
      } else {
        return (
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
            {result.loading || (!options.subComponent && !forageInitialised.current) ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
            {result.error ? <h1>{result.error.message}</h1> : undefined}
            {maybeRenderContent()}
            {maybeRenderPagination()}
          </div>
        )
      }
    }

    return { filter, template: getTemplate(), options, onSelectChange };
  }
}