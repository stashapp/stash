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
}

export interface IListHookOptions {
  filterMode: FilterMode;
  props: IBaseProps;
  renderContent: (result: QueryHookResult<any, any>, filter: ListFilterModel) => JSX.Element | undefined;
}

export class ListHook {
  public static useList(options: IListHookOptions): IListHookData {
    const [filter, setFilter] = useState<ListFilterModel>(new ListFilterModel(options.filterMode));

    // Update the filter when the query parameters change
    useEffect(() => {
      const queryParams = queryString.parse(options.props.location.search);
      filter.configureFromQueryParameters(queryParams);
      setFilter(filter);
    }, [options.props.location.search]);

    let result: QueryHookResult<any, any>;
    let totalCount: number;

    switch (options.filterMode) {
      case FilterMode.Scenes: {
        result = StashService.useFindScenes(filter);
        totalCount = !!result.data && !!result.data.findScenes ? result.data.findScenes.count : 0;
        break;
      }
      case FilterMode.SceneMarkers: {
        result = StashService.useFindSceneMarkers(filter);
        totalCount = !!result.data && !!result.data.findSceneMarkers ? result.data.findSceneMarkers.count : 0;
        break;
      }
      case FilterMode.Galleries: {
        result = StashService.useFindGalleries(filter);
        totalCount = !!result.data && !!result.data.findGalleries ? result.data.findGalleries.count : 0;
        break;
      }
      case FilterMode.Studios: {
        result = StashService.useFindStudios(filter);
        totalCount = !!result.data && !!result.data.findStudios ? result.data.findStudios.count : 0;
        break;
      }
      case FilterMode.Performers: {
        result = StashService.useFindPerformers(filter);
        totalCount = !!result.data && !!result.data.findPerformers ? result.data.findPerformers.count : 0;
        break;
      }
      default: {
        console.error("REMOVE DEFAULT IN LIST HOOK");
        result = StashService.useFindScenes(filter);
        totalCount = !!result.data && !!result.data.findScenes ? result.data.findScenes.count : 0;
        break;
      }
    }

    // Update the query parameters when the data changes
    useEffect(() => {
      const location = Object.assign({}, options.props.history.location);
      location.search = filter.makeQueryParameters();
      options.props.history.replace(location);
    }, [result.data, filter.displayMode]);

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
          filter={filter}
        />
        {result.loading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
        {result.error ? <h1>{result.error.message}</h1> : undefined}
        {options.renderContent(result, filter)}
        <Pagination
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
          onChangePage={onChangePage}
        />
      </div>
    );

    return { filter, template, options };
  }
}
