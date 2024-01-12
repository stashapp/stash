import React, { useCallback, useMemo, useState } from "react";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { DisplayModeSelect, ZoomSelect } from "../List/ListViewOptions";
import { DisplayMode } from "src/models/list-filter/types";
import { PageSizeSelect, SortBySelect } from "../List/ListFilter";
import {
  FilterMode,
  FindScenesQueryResult,
  SortDirectionEnum,
} from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useFindScenes } from "src/core/StashService";
import { SceneCardsGrid } from "./SceneCardsGrid";
import SceneQueue from "src/models/sceneQueue";
import { SceneListTable } from "./SceneListTable";
import { SceneWallPanel } from "../Wall/WallPanel";
import { Tagger } from "../Tagger/scenes/SceneTagger";
import { TaggerContext } from "../Tagger/context";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import {
  faChevronRight,
  faPlay,
  faShuffle,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";
import TextUtils from "src/utils/text";
import { FilterButton } from "../List/Filters/FilterButton";
import { useListSelect } from "src/hooks/listSelect";
import {
  IListFilterOperation,
  ListOperationButtons,
} from "../List/ListOperationButtons";
import { IItemListOperation } from "../List/ItemList";
import { FilterSidebar } from "../List/FilterSidebar";

export const DefaultListHeader: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  totalItems: number;
  filterHidden: boolean;
  onShowFilter: () => void;
}> = ({ filter, setFilter, totalItems, filterHidden, onShowFilter }) => {
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

  function onChangeSortDirection(dir: SortDirectionEnum) {
    const newFilter = filter.clone();
    newFilter.sortDirection = dir;
    setFilter(newFilter);
  }

  function onChangeSortBy(eventKey: string | null) {
    const newFilter = filter.clone();
    newFilter.sortBy = eventKey ?? undefined;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  function onReshuffleRandomSort() {
    const newFilter = filter.clone();
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
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
        {filterHidden && (
          <FilterButton
            filter={filter}
            icon={faChevronRight}
            onClick={() => onShowFilter()}
          />
        )}
        <PageSizeSelect
          pageSize={filter.itemsPerPage}
          setPageSize={onChangePageSize}
        />
        <Pagination
          currentPage={filter.currentPage}
          itemsPerPage={filter.itemsPerPage}
          totalItems={totalItems}
          onChangePage={onChangePage}
          pagesToShow={1}
        />
      </div>
      <div>
        <div>
          <Button className="play-scenes-button" variant="secondary">
            <Icon icon={faPlay} />
          </Button>
          <Button className="shuffle-scenes-button" variant="secondary">
            <Icon icon={faShuffle} />
          </Button>
        </div>
        <SortBySelect
          sortBy={filter.sortBy}
          direction={filter.sortDirection}
          options={filterOptions.sortByOptions}
          setSortBy={onChangeSortBy}
          setDirection={onChangeSortDirection}
          onReshuffleRandomSort={onReshuffleRandomSort}
        />
        <div>
          <ZoomSelect
            minZoom={0}
            maxZoom={3}
            zoomIndex={filter.zoomIndex}
            onChangeZoom={onChangeZoom}
          />
        </div>
        <DisplayModeSelect
          displayMode={filter.displayMode}
          displayModeOptions={filterOptions.displayModeOptions}
          onSetDisplayMode={onChangeDisplayMode}
        />
      </div>
    </div>
  );
};

export const SelectedListHeader: React.FC<{
  selectedIds: Set<string>;
  onSelectAll: () => void;
  onSelectNone: () => void;
  otherOperations: IListFilterOperation[];
}> = ({ selectedIds, onSelectAll, onSelectNone, otherOperations }) => {
  return (
    <div className="list-header">
      <div>
        <span>{selectedIds.size} items selected</span>
        <Button variant="link" onClick={() => onSelectAll()}>
          <FormattedMessage id="actions.select_all" />
        </Button>
      </div>
      <div>
        <Button className="play-scenes-button" variant="secondary">
          <Icon icon={faPlay} />
        </Button>

        <ListOperationButtons
          itemsSelected
          onEdit={() => {}}
          onDelete={() => {}}
          otherOperations={otherOperations}
        />
      </div>
      <div>
        <Button className="minimal select-none" onClick={() => onSelectNone()}>
          <Icon icon={faTimes} />
        </Button>
      </div>
    </div>
  );
};

export const ScenesPage: React.FC = ({}) => {
  const intl = useIntl();

  const [filter, setFilter] = useState<ListFilterModel>(
    () => new ListFilterModel(FilterMode.Scenes)
  );
  const [showFilter, setShowFilter] = useState(true);

  const result = useFindScenes(filter);
  const items = result.data?.findScenes.scenes ?? [];
  const { selectedIds, onSelectChange, onSelectAll, onSelectNone } =
    useListSelect(items);

  const totalCount = useMemo(
    () => result.data?.findScenes.count ?? 0,
    [result.data?.findScenes.count]
  );

  const metadataByline = useMemo(() => {
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
  }, [result]);

  function renderScenes() {
    if (!result.data?.findScenes) return;

    const queue = SceneQueue.fromListFilterModel(filter);

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <SceneCardsGrid
          scenes={items}
          queue={queue}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <SceneListTable
          scenes={items}
          queue={queue}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return <SceneWallPanel scenes={items} sceneQueue={queue} />;
    }
    if (filter.displayMode === DisplayMode.Tagger) {
      return (
        <TaggerContext>
          <Tagger scenes={items} queue={queue} />
        </TaggerContext>
      );
    }
  }

  // async function playSelected(
  //   result: FindScenesQueryResult,
  //   filter: ListFilterModel,
  //   selectedIds: Set<string>
  // ) {
  //   // populate queue and go to first scene
  //   // const sceneIDs = Array.from(selectedIds.values());
  //   // const queue = SceneQueue.fromSceneIDList(sceneIDs);
  //   // const autoPlay =
  //   //   config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
  //   // playScene(queue, sceneIDs[0], { autoPlay });
  // }

  // async function playRandom(r: FindScenesQueryResult, filter: ListFilterModel) {
  // query for a random scene
  // if (result.data?.findScenes) {
  //   const { count } = result.data.findScenes;
  //   const pages = Math.ceil(count / filter.itemsPerPage);
  //   const page = Math.floor(Math.random() * pages) + 1;
  //   const indexMax = Math.min(filter.itemsPerPage, count);
  //   const index = Math.floor(Math.random() * indexMax);
  //   const filterCopy = filter.clone();
  //   filterCopy.currentPage = page;
  //   filterCopy.sortBy = "random";
  //   const queryResults = await queryFindScenes(filterCopy);
  //   const scene = queryResults.data.findScenes.scenes[index];
  //   if (scene) {
  //     // navigate to the image player page
  //     const queue = SceneQueue.fromListFilterModel(filterCopy);
  //     const autoPlay =
  //       config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
  //     playScene(queue, scene.id, { sceneIndex: index, autoPlay });
  //   }
  // }
  // }

  // async function onMerge(
  //   r: FindScenesQueryResult,
  //   f: ListFilterModel,
  //   selectedIds: Set<string>
  // ) {
  //   // const selected =
  //   //   result.data?.findScenes.scenes
  //   //     .filter((s) => selectedIds.has(s.id))
  //   //     .map((s) => {
  //   //       return {
  //   //         id: s.id,
  //   //         title: objectTitle(s),
  //   //       };
  //   //     }) ?? [];
  //   // setMergeScenes(selected);
  // }

  // async function onExport() {
  //   // setIsExportAll(false);
  //   // setIsExportDialogOpen(true);
  // }

  // async function onExportAll() {
  //   // setIsExportAll(true);
  //   // setIsExportDialogOpen(true);
  // }

  const otherOperations: IItemListOperation<FindScenesQueryResult>[] = [
    {
      text: intl.formatMessage({ id: "actions.play_random" }),
      onClick: async () => {}, // playRandom,
    },
    {
      text: `${intl.formatMessage({ id: "actions.generate" })}…`,
      onClick: async () => {}, // setIsGenerateDialogOpen(true),
    },
    {
      text: `${intl.formatMessage({ id: "actions.identify" })}…`,
      onClick: async () => {}, // setIsIdentifyDialogOpen(true),
    },
    {
      text: `${intl.formatMessage({ id: "actions.merge" })}…`,
      onClick: async () => {}, // onMerge,
    },
    {
      text: intl.formatMessage({ id: "actions.export" }),
      onClick: async () => {}, // onExport,
    },
    {
      text: intl.formatMessage({ id: "actions.export_all" }),
      onClick: async () => {}, // onExportAll,
    },
  ];

  async function onOperationClicked(
    o: IItemListOperation<FindScenesQueryResult>
  ) {
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

  return (
    <div id="scenes-page">
      {showFilter && (
        <FilterSidebar
          onHide={() => setShowFilter(false)}
          filter={filter}
          setFilter={(f) => setFilter(f)}
        />
      )}
      <div className={cx("scenes-page-results", { expanded: !showFilter })}>
        {selectedIds.size === 0 ? (
          <DefaultListHeader
            filter={filter}
            setFilter={setFilter}
            totalItems={totalCount}
            filterHidden={!showFilter}
            onShowFilter={() => setShowFilter(true)}
          />
        ) : (
          <SelectedListHeader
            selectedIds={selectedIds}
            onSelectAll={onSelectAll}
            onSelectNone={onSelectNone}
            otherOperations={operations}
          />
        )}
        <div className="scenes-page-items">
          <PaginationIndex
            itemsPerPage={filter.itemsPerPage}
            currentPage={filter.currentPage}
            totalItems={totalCount}
            metadataByline={metadataByline}
          />
          {renderScenes()}
        </div>
      </div>
    </div>
  );
};

export default ScenesPage;
