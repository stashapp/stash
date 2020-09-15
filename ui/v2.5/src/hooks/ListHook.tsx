import _ from "lodash";
import queryString from "query-string";
import React, { useCallback, useRef, useState, useEffect } from "react";
import { ApolloError } from "@apollo/client";
import { useHistory, useLocation } from "react-router-dom";
import {
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
  FindTagsQueryResult,
  TagDataFragment,
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
  useFindTags,
} from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { FilterMode } from "src/models/list-filter/types";

const getSelectedData = <I extends IDataItem>(
  result: I[],
  selectedIds: Set<string>
) => {
  // find the selected items from the ids
  const selectedResults: I[] = [];

  selectedIds.forEach((id) => {
    const item = result.find((s) => s.id === id);

    if (item) {
      selectedResults.push(item);
    }
  });

  return selectedResults;
};

interface IListHookData {
  filter: ListFilterModel;
  template: React.ReactElement;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

interface IListHookOperation<T> {
  text: string;
  onClick: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => void;
  isDisplayed?: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => boolean;
}

interface IListHookOptions<T, E> {
  persistState?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  zoomable?: boolean;
  selectable?: boolean;
  defaultZoomIndex?: number;
  otherOperations?: IListHookOperation<T>[];
  renderContent: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
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

interface IDataItem {
  id: string;
}
interface IQueryResult {
  error?: ApolloError;
  loading: boolean;
  refetch: () => void;
}

interface IQuery<T extends IQueryResult, T2 extends IDataItem> {
  filterMode: FilterMode;
  useData: (filter: ListFilterModel) => T;
  getData: (data: T) => T2[];
  getCount: (data: T) => number;
}

interface IRenderListProps {
  filter: ListFilterModel;
  onChangePage: (page: number) => void;
  updateQueryParams: (filter: ListFilterModel) => void;
}

const RenderList = <
  QueryResult extends IQueryResult,
  QueryData extends IDataItem
>({
  defaultZoomIndex,
  filter,
  onChangePage,
  addKeybinds,
  useData,
  getCount,
  getData,
  otherOperations,
  renderContent,
  zoomable,
  selectable,
  renderEditDialog,
  renderDeleteDialog,
  updateQueryParams,
}: IListHookOptions<QueryResult, QueryData> &
  IQuery<QueryResult, QueryData> &
  IRenderListProps) => {
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [lastClickedId, setLastClickedId] = useState<string | undefined>();
  const [zoomIndex, setZoomIndex] = useState<number>(defaultZoomIndex ?? 1);

  const result = useData(filter);
  const totalCount = getCount(result);
  const items = getData(result);

  useEffect(() => {
    Mousetrap.bind("right", () => {
      const maxPage = totalCount / filter.itemsPerPage;
      if (filter.currentPage < maxPage) {
        onChangePage(filter.currentPage + 1);
      }
    });
    Mousetrap.bind("left", () => {
      if (filter.currentPage > 1) {
        onChangePage(filter.currentPage - 1);
      }
    });
    Mousetrap.bind("shift+right", () => {
      const maxPage = totalCount / filter.itemsPerPage + 1;
      onChangePage(Math.min(maxPage, filter.currentPage + 10));
    });
    Mousetrap.bind("shift+left", () => {
      onChangePage(Math.max(1, filter.currentPage - 10));
    });
    Mousetrap.bind("ctrl+end", () => {
      const maxPage = totalCount / filter.itemsPerPage + 1;
      onChangePage(maxPage);
    });
    Mousetrap.bind("ctrl+home", () => {
      onChangePage(1);
    });

    let unbindExtras: () => void;
    if (addKeybinds) {
      unbindExtras = addKeybinds(result, filter, selectedIds);
    }

    return () => {
      Mousetrap.unbind("right");
      Mousetrap.unbind("left");
      Mousetrap.unbind("shift+right");
      Mousetrap.unbind("shift+left");
      Mousetrap.unbind("ctrl+end");
      Mousetrap.unbind("ctrl+home");

      if (unbindExtras) {
        unbindExtras();
      }
    };
  });

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

  const operations =
    otherOperations &&
    otherOperations.map((o) => ({
      text: o.text,
      onClick: () => {
        o.onClick(result, filter, selectedIds);
      },
      isDisplayed: () => {
        if (o.isDisplayed) {
          return o.isDisplayed(result, filter, selectedIds);
        }

        return true;
      },
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

  const renderPagination = () => (
    <Pagination
      itemsPerPage={filter.itemsPerPage}
      currentPage={filter.currentPage}
      totalItems={totalCount}
      onChangePage={onChangePage}
    />
  );

  function maybeRenderContent() {
    if (result.loading || result.error) {
      return;
    }

    return (
      <>
        {renderPagination()}
        {renderContent(result, filter, selectedIds, zoomIndex)}
        <PaginationIndex
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
        />
        {renderPagination()}
      </>
    );
  }

  const content = (
    <div>
      <ListFilter
        onFilterUpdate={updateQueryParams}
        onSelectAll={selectable ? onSelectAll : undefined}
        onSelectNone={selectable ? onSelectNone : undefined}
        zoomIndex={zoomable ? zoomIndex : undefined}
        onChangeZoom={zoomable ? onChangeZoom : undefined}
        otherOperations={operations}
        itemsSelected={selectedIds.size > 0}
        onEdit={renderEditDialog ? onEdit : undefined}
        onDelete={renderDeleteDialog ? onDelete : undefined}
        filter={filter}
      />
      {isEditDialogOpen &&
        renderEditDialog &&
        renderEditDialog(
          getSelectedData(getData(result), selectedIds),
          (applied) => onEditDialogClosed(applied)
        )}
      {isDeleteDialogOpen &&
        renderDeleteDialog &&
        renderDeleteDialog(
          getSelectedData(getData(result), selectedIds),
          (deleted) => onDeleteDialogClosed(deleted)
        )}
      {result.loading ? <LoadingIndicator /> : undefined}
      {result.error ? <h1>{result.error.message}</h1> : undefined}
      {maybeRenderContent()}
    </div>
  );

  return { contentTemplate: content, onSelectChange };
};

const useList = <QueryResult extends IQueryResult, QueryData extends IDataItem>(
  options: IListHookOptions<QueryResult, QueryData> &
    IQuery<QueryResult, QueryData>
): IListHookData => {
  const history = useHistory();
  const location = useLocation();
  const [interfaceState, setInterfaceState] = useInterfaceLocalForage();
  // If persistState is false we don't care about forage and consider it initialised
  const [forageInitialised, setForageInitialised] = useState(
    !options.persistState
  );
  // Store initial pathname to prevent hooks from operating outside this page
  const originalPathName = useRef(location.pathname);

  const [filter, setFilter] = useState<ListFilterModel>(
    new ListFilterModel(options.filterMode, queryString.parse(location.search))
  );

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
    if (
      interfaceState.loading ||
      // Only update query params on page the hook was mounted on
      history.location.pathname !== originalPathName.current
    )
      return;

    if (!forageInitialised) setForageInitialised(true);

    if (!options.persistState) return;

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
    options.filterMode,
    forageInitialised,
    updateInterfaceConfig,
    options.persistState,
  ]);

  function updateQueryParams(listFilter: ListFilterModel) {
    setFilter(listFilter);
    const newLocation = { ...location };
    newLocation.search = listFilter.makeQueryParameters();
    history.replace(newLocation);
    if (options.persistState) {
      updateInterfaceConfig(listFilter);
    }
  }

  const onChangePage = (page: number) => {
    const newFilter = _.cloneDeep(filter);
    newFilter.currentPage = page;
    updateQueryParams(newFilter);
  };

  const renderFilter = !options.filterHook
    ? filter
    : options.filterHook(_.cloneDeep(filter));

  const { contentTemplate, onSelectChange } = RenderList({
    ...options,
    filter: renderFilter,
    onChangePage,
    updateQueryParams,
  });

  const template = !forageInitialised ? (
    <LoadingIndicator />
  ) : (
    <>{contentTemplate}</>
  );

  return {
    filter,
    template,
    onSelectChange,
  };
};

export const useScenesList = (
  props: IListHookOptions<FindScenesQueryResult, SlimSceneDataFragment>
) =>
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
  props: IListHookOptions<FindSceneMarkersQueryResult, SceneMarkerDataFragment>
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
  props: IListHookOptions<FindGalleriesQueryResult, GalleryDataFragment>
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
  props: IListHookOptions<FindStudiosQueryResult, StudioDataFragment>
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
  props: IListHookOptions<FindPerformersQueryResult, PerformerDataFragment>
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

export const useMoviesList = (
  props: IListHookOptions<FindMoviesQueryResult, MovieDataFragment>
) =>
  useList<FindMoviesQueryResult, MovieDataFragment>({
    ...props,
    filterMode: FilterMode.Movies,
    useData: useFindMovies,
    getData: (result: FindMoviesQueryResult) =>
      result?.data?.findMovies?.movies ?? [],
    getCount: (result: FindMoviesQueryResult) =>
      result?.data?.findMovies?.count ?? 0,
  });

export const useTagsList = (
  props: IListHookOptions<FindTagsQueryResult, TagDataFragment>
) =>
  useList<FindTagsQueryResult, TagDataFragment>({
    ...props,
    filterMode: FilterMode.Tags,
    useData: useFindTags,
    getData: (result: FindTagsQueryResult) =>
      result?.data?.findTags?.tags ?? [],
    getCount: (result: FindTagsQueryResult) =>
      result?.data?.findTags?.count ?? 0,
  });

export const showWhenSelected = <T extends IQueryResult>(
  _result: T,
  _filter: ListFilterModel,
  selectedIds: Set<string>
) => {
  return selectedIds.size > 0;
};
