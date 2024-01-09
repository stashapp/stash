import React, { useCallback, useMemo, useState } from "react";
import { Pagination } from "../List/Pagination";
import { ListViewOptions, ZoomSelect } from "../List/ListViewOptions";
import { DisplayMode } from "src/models/list-filter/types";
import { PageSizeSelect, SearchField, SortBySelect } from "../List/ListFilter";
import { FilterMode, SortDirectionEnum } from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useFindScenes } from "src/core/StashService";
import { SceneCardsGrid } from "./SceneCardsGrid";
import SceneQueue from "src/models/sceneQueue";
import { SceneListTable } from "./SceneListTable";
import { SceneWallPanel } from "../Wall/WallPanel";
import { Tagger } from "../Tagger/scenes/SceneTagger";
import { TaggerContext } from "../Tagger/context";
import useFocus from "src/utils/focus";

const SceneFilter: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
}> = ({ filter, setFilter }) => {
  const [queryRef, setQueryFocus] = useFocus();

  const searchQueryUpdated = useCallback(
    (value: string) => {
      const newFilter = filter.clone();
      newFilter.searchTerm = value;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    },
    [filter, setFilter]
  );

  return (
    <div className="scene-filter">
      <SearchField
        searchTerm={filter.searchTerm}
        setSearchTerm={searchQueryUpdated}
        queryRef={queryRef}
        setQueryFocus={setQueryFocus}
      />
    </div>
  );
};

export const ListHeader: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  totalItems: number;
}> = ({ filter, setFilter, totalItems }) => {
  const filterOptions = getFilterOptions(filter.mode);

  function onChangeZoom(newZoomIndex: number) {
    const newFilter = filter.clone();
    newFilter.zoomIndex = newZoomIndex;
    setFilter(newFilter);
  }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = filter.clone();
    newFilter.displayMode = displayMode;
    setFilter(newFilter);
  }

  function onChangePageSize(val: number) {
    const newFilter = filter.clone();
    newFilter.itemsPerPage = val;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  const onChangePage = useCallback(
    (page: number) => {
      const newFilter = filter.clone();
      newFilter.currentPage = page;
      setFilter(newFilter);

      // if the current page has a detail-header, then
      // scroll up relative to that rather than 0, 0
      const detailHeader = document.querySelector(".detail-header");
      if (detailHeader) {
        window.scrollTo(0, detailHeader.scrollHeight - 50);
      } else {
        window.scrollTo(0, 0);
      }
    },
    [filter, setFilter]
  );

  return (
    <div className="list-header">
      <div>
        <Pagination
          currentPage={filter.currentPage}
          itemsPerPage={filter.itemsPerPage}
          totalItems={totalItems}
          onChangePage={onChangePage}
        />
        <PageSizeSelect
          pageSize={filter.itemsPerPage}
          setPageSize={onChangePageSize}
        />
      </div>
      <div>
        <SortBySelect
          sortBy="title"
          direction={SortDirectionEnum.Asc}
          options={filterOptions.sortByOptions}
          setSortBy={() => {}}
          setDirection={() => {}}
          onReshuffleRandomSort={() => {}}
        />
        <div>
          <ZoomSelect
            minZoom={0}
            maxZoom={3}
            zoomIndex={filter.zoomIndex}
            onChangeZoom={onChangeZoom}
          />
        </div>
        <ListViewOptions
          displayMode={filter.displayMode}
          displayModeOptions={filterOptions.displayModeOptions}
          onSetDisplayMode={onChangeDisplayMode}
        />
      </div>
    </div>
  );
};

export const ScenesPage: React.FC = ({}) => {
  const [filter, setFilter] = useState<ListFilterModel>(
    () => new ListFilterModel(FilterMode.Scenes)
  );

  const result = useFindScenes(filter);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const totalCount = useMemo(
    () => result.data?.findScenes.count ?? 0,
    [result.data?.findScenes.count]
  );

  function renderScenes() {
    if (!result.data?.findScenes) return;

    const queue = SceneQueue.fromListFilterModel(filter);

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <SceneCardsGrid
          scenes={result.data.findScenes.scenes}
          queue={queue}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={() => {}}
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <SceneListTable
          scenes={result.data.findScenes.scenes}
          queue={queue}
          selectedIds={selectedIds}
          onSelectChange={() => {}}
        />
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <SceneWallPanel
          scenes={result.data.findScenes.scenes}
          sceneQueue={queue}
        />
      );
    }
    if (filter.displayMode === DisplayMode.Tagger) {
      return (
        <TaggerContext>
          <Tagger scenes={result.data.findScenes.scenes} queue={queue} />
        </TaggerContext>
      );
    }
  }

  return (
    <div id="scenes-page">
      <SceneFilter filter={filter} setFilter={(f) => setFilter(f)} />
      <div className="scenes-page-results">
        <ListHeader
          filter={filter}
          setFilter={(f) => setFilter(f)}
          totalItems={totalCount}
        />
        <div className="scenes-page-items">{renderScenes()}</div>
      </div>
    </div>
  );
};

export default ScenesPage;
