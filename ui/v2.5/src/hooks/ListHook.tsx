import _ from "lodash";
import queryString from "query-string";
import React, {
  useCallback,
  useRef,
  useState,
  useEffect,
  useMemo,
} from "react";
import { ApolloError } from "@apollo/client";
import { useHistory, useLocation } from "react-router-dom";
import Mousetrap from "mousetrap";
import { IconProp } from "@fortawesome/fontawesome-svg-core";
import {
  SlimSceneDataFragment,
  SceneMarkerDataFragment,
  SlimGalleryDataFragment,
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
  FindImagesQueryResult,
  SlimImageDataFragment,
  FilterMode,
} from "src/core/generated-graphql";
import { useInterfaceLocalForage } from "src/hooks/LocalForage";
import { LoadingIndicator } from "src/components/Shared";
import { ListFilter } from "src/components/List/ListFilter";
import { FilterTags } from "src/components/List/FilterTags";
import { Pagination, PaginationIndex } from "src/components/List/Pagination";
import {
  useFindDefaultFilter,
  useFindScenes,
  useFindSceneMarkers,
  useFindImages,
  useFindMovies,
  useFindStudios,
  useFindGalleries,
  useFindPerformers,
  useFindTags,
} from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { ListFilterOptions } from "src/models/list-filter/filter-options";
import { getFilterOptions } from "src/models/list-filter/factory";
import { ButtonToolbar } from "react-bootstrap";
import { ListViewOptions } from "src/components/List/ListViewOptions";
import { ListOperationButtons } from "src/components/List/ListOperationButtons";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { AddFilterDialog } from "src/components/List/AddFilterDialog";
import { TextUtils } from "src/utils";
import { FormattedNumber } from "react-intl";

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
  onChangePage: (page: number) => void;
}

export interface IListHookOperation<T> {
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
  postRefetch?: boolean;
  icon?: IconProp;
  buttonVariant?: string;
}

export enum PersistanceLevel {
  // do not load default query or persist display mode
  NONE,
  // load default query, don't load or persist display mode
  ALL,
  // load and persist display mode only
  VIEW,
}

interface IListHookOptions<T, E> {
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
  defaultZoomIndex?: number;
  otherOperations?: IListHookOperation<T>[];
  renderContent: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>,
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
  getMetadataByline: (data: T) => React.ReactNode;
}

interface IRenderListProps {
  filter: ListFilterModel;
  filterOptions: ListFilterOptions;
  onChangePage: (page: number) => void;
  updateQueryParams: (filter: ListFilterModel) => void;
}

const RenderList = <
  QueryResult extends IQueryResult,
  QueryData extends IDataItem
>({
  filter,
  filterOptions,
  onChangePage,
  addKeybinds,
  useData,
  getCount,
  getData,
  getMetadataByline,
  otherOperations,
  renderContent,
  zoomable,
  selectable,
  renderEditDialog,
  renderDeleteDialog,
  updateQueryParams,
  filterDialog,
  persistState,
}: IListHookOptions<QueryResult, QueryData> &
  IQuery<QueryResult, QueryData> &
  IRenderListProps) => {
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [lastClickedId, setLastClickedId] = useState<string | undefined>();

  const [editingCriterion, setEditingCriterion] = useState<
    Criterion<CriterionValue> | undefined
  >(undefined);
  const [newCriterion, setNewCriterion] = useState(false);

  const result = useData(filter);
  const totalCount = getCount(result);
  const metadataByline = getMetadataByline(result);
  const items = getData(result);
  const pages = Math.ceil(totalCount / filter.itemsPerPage);

  // handle case where page is more than there are pages
  useEffect(() => {
    if (pages > 0 && filter.currentPage > pages) {
      onChangePage(pages);
    }
  }, [pages, filter.currentPage, onChangePage]);

  useEffect(() => {
    Mousetrap.bind("f", () => setNewCriterion(true));
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
    const newFilter = _.cloneDeep(filter);
    newFilter.zoomIndex = newZoomIndex;
    updateQueryParams(newFilter);
  }

  async function onOperationClicked(o: IListHookOperation<QueryResult>) {
    await o.onClick(result, filter, selectedIds);
    if (o.postRefetch) {
      result.refetch();
    }
  }

  const operations =
    otherOperations &&
    otherOperations.map((o) => ({
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

  const renderPagination = () => (
    <Pagination
      itemsPerPage={filter.itemsPerPage}
      currentPage={filter.currentPage}
      totalItems={totalCount}
      metadataByline={metadataByline}
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
        <PaginationIndex
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
          metadataByline={metadataByline}
        />
        {renderContent(result, filter, selectedIds, onChangePage, pages)}
        <PaginationIndex
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
          metadataByline={metadataByline}
        />
        {renderPagination()}
      </>
    );
  }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = _.cloneDeep(filter);
    newFilter.displayMode = displayMode;
    updateQueryParams(newFilter);
  }

  function onAddCriterion(
    criterion: Criterion<CriterionValue>,
    oldId?: string
  ) {
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
    setEditingCriterion(undefined);
    setNewCriterion(false);
  }

  function onRemoveCriterion(removedCriterion: Criterion<CriterionValue>) {
    const newFilter = _.cloneDeep(filter);
    newFilter.criteria = newFilter.criteria.filter(
      (criterion) => criterion.getId() !== removedCriterion.getId()
    );
    newFilter.currentPage = 1;
    updateQueryParams(newFilter);
  }

  function updateCriteria(c: Criterion<CriterionValue>[]) {
    const newFilter = _.cloneDeep(filter);
    newFilter.criteria = c.slice();
    setNewCriterion(false);
  }

  function onCancelAddCriterion() {
    setEditingCriterion(undefined);
    setNewCriterion(false);
  }

  const content = (
    <div>
      <ButtonToolbar className="align-items-center justify-content-center mb-2">
        <ListFilter
          onFilterUpdate={updateQueryParams}
          filter={filter}
          filterOptions={filterOptions}
          openFilterDialog={() => setNewCriterion(true)}
          filterDialogOpen={newCriterion ?? editingCriterion}
          persistState={persistState}
        />
        <ListOperationButtons
          onSelectAll={selectable ? onSelectAll : undefined}
          onSelectNone={selectable ? onSelectNone : undefined}
          otherOperations={operations}
          itemsSelected={selectedIds.size > 0}
          onEdit={renderEditDialog ? onEdit : undefined}
          onDelete={renderDeleteDialog ? onDelete : undefined}
        />
        <ListViewOptions
          displayMode={filter.displayMode}
          displayModeOptions={filterOptions.displayModeOptions}
          onSetDisplayMode={onChangeDisplayMode}
          zoomIndex={zoomable ? filter.zoomIndex : undefined}
          onSetZoom={zoomable ? onChangeZoom : undefined}
        />
      </ButtonToolbar>
      <FilterTags
        criteria={filter.criteria}
        onEditCriterion={(c) => setEditingCriterion(c)}
        onRemoveCriterion={onRemoveCriterion}
      />
      {(newCriterion || editingCriterion) && !filterDialog && (
        <AddFilterDialog
          filterOptions={filterOptions}
          onAddCriterion={onAddCriterion}
          onCancel={onCancelAddCriterion}
          editingCriterion={editingCriterion}
          existingCriterions={filter.criteria}
        />
      )}
      {newCriterion &&
        filterDialog &&
        filterDialog(filter.criteria, (c) => updateCriteria(c))}
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
  const filterOptions = getFilterOptions(options.filterMode);

  const history = useHistory();
  const location = useLocation();
  const [interfaceState, setInterfaceState] = useInterfaceLocalForage();
  // If persistState is false we don't care about forage and consider it initialised
  const [forageInitialised, setForageInitialised] = useState(
    !options.persistState
  );
  // Store initial pathname to prevent hooks from operating outside this page
  const originalPathName = useRef(location.pathname);
  const persistanceKey = options.persistanceKey ?? options.filterMode;

  const defaultSort = options.defaultSort ?? filterOptions.defaultSortBy;
  const defaultDisplayMode = filterOptions.displayModeOptions[0];
  const [filter, setFilter] = useState<ListFilterModel>(
    new ListFilterModel(
      options.filterMode,
      queryString.parse(location.search),
      defaultSort,
      defaultDisplayMode,
      options.defaultZoomIndex
    )
  );

  const updateInterfaceConfig = useCallback(
    (updatedFilter: ListFilterModel, level: PersistanceLevel) => {
      if (level === PersistanceLevel.VIEW) {
        setInterfaceState((prevState) => {
          if (!prevState.queryConfig) {
            prevState.queryConfig = {};
          }
          return {
            ...prevState,
            queryConfig: {
              ...prevState.queryConfig,
              [persistanceKey]: {
                ...prevState.queryConfig[persistanceKey],
                filter: queryString.stringify({
                  ...queryString.parse(
                    prevState.queryConfig[persistanceKey]?.filter ?? ""
                  ),
                  disp: updatedFilter.displayMode,
                }),
              },
            },
          };
        });
      }
    },
    [persistanceKey, setInterfaceState]
  );

  const {
    data: defaultFilter,
    loading: defaultFilterLoading,
  } = useFindDefaultFilter(options.filterMode);

  const updateQueryParams = useCallback(
    (listFilter: ListFilterModel) => {
      setFilter(listFilter);
      const newLocation = { ...location };
      newLocation.search = listFilter.makeQueryParameters();
      history.replace(newLocation);
      if (options.persistState) {
        updateInterfaceConfig(listFilter, options.persistState);
      }
    },
    [setFilter, history, location, options.persistState, updateInterfaceConfig]
  );

  useEffect(() => {
    if (
      // defer processing this until forage is initialised and
      // default filter is loaded
      interfaceState.loading ||
      defaultFilterLoading ||
      // Only update query params on page the hook was mounted on
      history.location.pathname !== originalPathName.current
    )
      return;

    if (!forageInitialised) setForageInitialised(true);

    const newFilter = filter.clone();
    let update = false;

    // Compare constructed filter with current filter.
    // If different it's the result of navigation, and we update the filter.
    if (
      history.location.search &&
      history.location.search !== `?${filter.makeQueryParameters()}`
    ) {
      newFilter.configureFromQueryParameters(
        queryString.parse(history.location.search)
      );
      update = true;
    }

    // if default query is set and no search params are set, then
    // load the default query
    // #1512 - use default query only if persistState is ALL
    if (
      options.persistState === PersistanceLevel.ALL &&
      !location.search &&
      defaultFilter?.findDefaultFilter
    ) {
      newFilter.currentPage = 1;
      try {
        newFilter.configureFromQueryParameters(
          JSON.parse(defaultFilter.findDefaultFilter.filter)
        );
      } catch (err) {
        console.log(err);
        // ignore
      }
      // #1507 - reset random seed when loaded
      newFilter.randomSeed = -1;
      update = true;
    }

    // set the display type if persisted
    const storedQuery = interfaceState.data?.queryConfig?.[persistanceKey];

    if (options.persistState === PersistanceLevel.VIEW && storedQuery) {
      const storedFilter = queryString.parse(storedQuery.filter);

      if (storedFilter.disp !== undefined) {
        const displayMode = Number.parseInt(storedFilter.disp as string, 10);
        if (displayMode !== newFilter.displayMode) {
          newFilter.displayMode = displayMode;
          update = true;
        }
      }
    }

    if (update) {
      updateQueryParams(newFilter);
    }
  }, [
    defaultSort,
    defaultDisplayMode,
    filter,
    interfaceState,
    history,
    location.search,
    updateQueryParams,
    defaultFilter,
    defaultFilterLoading,
    persistanceKey,
    forageInitialised,
    options.persistState,
  ]);

  const onChangePage = useCallback(
    (page: number) => {
      const newFilter = _.cloneDeep(filter);
      newFilter.currentPage = page;
      updateQueryParams(newFilter);
      window.scrollTo(0, 0);
    },
    [filter, updateQueryParams]
  );

  const renderFilter = useMemo(() => {
    return !options.filterHook
      ? filter
      : options.filterHook(_.cloneDeep(filter));
  }, [filter, options]);

  const { contentTemplate, onSelectChange } = RenderList({
    ...options,
    filter: renderFilter,
    filterOptions,
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
    onChangePage,
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
    getMetadataByline: (result: FindScenesQueryResult) => {
      const duration = result?.data?.findScenes?.duration;
      const size = result?.data?.findScenes?.filesize;
      const filesize = size ? TextUtils.fileSize(size) : undefined;

      if (!duration && !size) {
        return;
      }

      const separator = duration && size ? " - " : "";

      return (
        <span className="scenes-stats">
          &nbsp;(
          {duration ? (
            <span className="scenes-duration">
              {TextUtils.secondsAsTimeString(duration, 3)}
            </span>
          ) : undefined}
          {separator}
          {size && filesize ? (
            <span className="scenes-size">
              <FormattedNumber
                value={filesize.size}
                maximumFractionDigits={TextUtils.fileSizeFractionalDigits(
                  filesize.unit
                )}
              />
              {` ${TextUtils.formatFileSizeUnit(filesize.unit)}`}
            </span>
          ) : undefined}
          )
        </span>
      );
    },
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
    getMetadataByline: () => [],
  });

export const useImagesList = (
  props: IListHookOptions<FindImagesQueryResult, SlimImageDataFragment>
) =>
  useList<FindImagesQueryResult, SlimImageDataFragment>({
    ...props,
    filterMode: FilterMode.Images,
    useData: useFindImages,
    getData: (result: FindImagesQueryResult) =>
      result?.data?.findImages?.images ?? [],
    getCount: (result: FindImagesQueryResult) =>
      result?.data?.findImages?.count ?? 0,
    getMetadataByline: (result: FindImagesQueryResult) => {
      const megapixels = result?.data?.findImages?.megapixels;
      const size = result?.data?.findImages?.filesize;
      const filesize = size ? TextUtils.fileSize(size) : undefined;

      if (!megapixels && !size) {
        return;
      }

      const separator = megapixels && size ? " - " : "";

      return (
        <span className="images-stats">
          &nbsp;(
          {megapixels ? (
            <span className="images-megapixels">
              <FormattedNumber value={megapixels} /> Megapixels
            </span>
          ) : undefined}
          {separator}
          {size && filesize ? (
            <span className="images-size">
              <FormattedNumber
                value={filesize.size}
                maximumFractionDigits={TextUtils.fileSizeFractionalDigits(
                  filesize.unit
                )}
              />
              {` ${TextUtils.formatFileSizeUnit(filesize.unit)}`}
            </span>
          ) : undefined}
          )
        </span>
      );
    },
  });

export const useGalleriesList = (
  props: IListHookOptions<FindGalleriesQueryResult, SlimGalleryDataFragment>
) =>
  useList<FindGalleriesQueryResult, SlimGalleryDataFragment>({
    ...props,
    filterMode: FilterMode.Galleries,
    useData: useFindGalleries,
    getData: (result: FindGalleriesQueryResult) =>
      result?.data?.findGalleries?.galleries ?? [],
    getCount: (result: FindGalleriesQueryResult) =>
      result?.data?.findGalleries?.count ?? 0,
    getMetadataByline: () => [],
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
    getMetadataByline: () => [],
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
    getMetadataByline: () => [],
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
    getMetadataByline: () => [],
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
    getMetadataByline: () => [],
  });

export const showWhenSelected = <T extends IQueryResult>(
  _result: T,
  _filter: ListFilterModel,
  selectedIds: Set<string>
) => {
  return selectedIds.size > 0;
};
