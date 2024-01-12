import React, { useMemo, useState } from "react";
import { PaginationIndex } from "../List/Pagination";
import { DisplayMode } from "src/models/list-filter/types";
import { FilterMode, FindScenesQueryResult } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useFindScenes } from "src/core/StashService";
import { SceneCardsGrid } from "./SceneCardsGrid";
import SceneQueue from "src/models/sceneQueue";
import { SceneListTable } from "./SceneListTable";
import { SceneWallPanel } from "../Wall/WallPanel";
import { Tagger } from "../Tagger/scenes/SceneTagger";
import { TaggerContext } from "../Tagger/context";
import { FormattedNumber, useIntl } from "react-intl";
import cx from "classnames";
import TextUtils from "src/utils/text";
import { useListSelect } from "src/hooks/listSelect";
import { IItemListOperation } from "../List/ItemList";
import { FilterSidebar } from "../List/FilterSidebar";
import { ListHeader } from "../List/ListHeader";

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
        <ListHeader
          filter={filter}
          setFilter={setFilter}
          totalItems={totalCount}
          filterHidden={!showFilter}
          onShowFilter={() => setShowFilter(true)}
          selectedIds={selectedIds}
          onSelectAll={onSelectAll}
          onSelectNone={onSelectNone}
          otherOperations={operations}
        />
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
