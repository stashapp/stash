import React, { useCallback, useContext, useEffect, useMemo, useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindScenes, useFindScenes } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { Tagger } from "../Tagger/scenes/SceneTagger";
import { IPlaySceneOptions, SceneQueue } from "src/models/sceneQueue";
import { SceneWallPanel } from "../Wall/WallPanel";
import { SceneListTable } from "./SceneListTable";
import { EditScenesDialog } from "./EditScenesDialog";
import { DeleteScenesDialog } from "./DeleteScenesDialog";
import { GenerateDialog } from "../Dialogs/GenerateDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { SceneCardsGrid } from "./SceneCardsGrid";
import { TaggerContext } from "../Tagger/context";
import { IdentifyDialog } from "../Dialogs/IdentifyDialog/IdentifyDialog";
import { ConfigurationContext } from "src/hooks/Config";
import { faPlay } from "@fortawesome/free-solid-svg-icons";
import { SceneMergeModal } from "./SceneMergeDialog";
import { objectTitle } from "src/core/files";
import TextUtils from "src/utils/text";
import { View } from "../List/views";
import { FileSize } from "../Shared/FileSize";
import { PagedList } from "../List/PagedList";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import { IListFilterOperation } from "../List/ListOperationButtons";
import { FilteredListToolbar } from "../List/FilteredListToolbar";
import { useFilteredItemList } from "../List/ItemList";
import { FilterTags } from "../List/FilterTags";
import { Sidebar, SidebarPane } from "../Shared/Sidebar";

function renderMetadataByline(result: GQL.FindScenesQueryResult) {
  const duration = result?.data?.findScenes?.duration;
  const size = result?.data?.findScenes?.filesize;

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
      {size ? (
        <span className="scenes-size">
          <FileSize size={size} />
        </span>
      ) : undefined}
      )
    </span>
  );
}

function usePlayScene() {
  const history = useHistory();

  const playScene = useCallback(
    (queue: SceneQueue, sceneID: string, options: IPlaySceneOptions) => {
      history.push(queue.makeLink(sceneID, options));
    },
    [history]
  );

  return playScene;
}

function usePlaySelected(selectedIds: Set<string>) {
  const { configuration: config } = useContext(ConfigurationContext);
  const playScene = usePlayScene();

  const playSelected = useCallback(() => {
    // populate queue and go to first scene
    const sceneIDs = Array.from(selectedIds.values());
    const queue = SceneQueue.fromSceneIDList(sceneIDs);
    const autoPlay = config?.interface.autostartVideoOnPlaySelected ?? false;
    playScene(queue, sceneIDs[0], { autoPlay });
  }, [selectedIds, config?.interface.autostartVideoOnPlaySelected, playScene]);

  return playSelected;
}

function usePlayRandom(filter: ListFilterModel, count: number) {
  const { configuration: config } = useContext(ConfigurationContext);
  const playScene = usePlayScene();

  const playRandom = useCallback(async () => {
    // query for a random scene
    if (count === 0) {
      return;
    }

    const pages = Math.ceil(count / filter.itemsPerPage);
    const page = Math.floor(Math.random() * pages) + 1;

    const indexMax = Math.min(filter.itemsPerPage, count);
    const index = Math.floor(Math.random() * indexMax);
    const filterCopy = cloneDeep(filter);
    filterCopy.currentPage = page;
    filterCopy.sortBy = "random";
    const queryResults = await queryFindScenes(filterCopy);
    const scene = queryResults.data.findScenes.scenes[index];
    if (scene) {
      // navigate to the image player page
      const queue = SceneQueue.fromListFilterModel(filterCopy);
      const autoPlay = config?.interface.autostartVideoOnPlaySelected ?? false;
      playScene(queue, scene.id, { sceneIndex: index, autoPlay });
    }
  }, [
    filter,
    count,
    config?.interface.autostartVideoOnPlaySelected,
    playScene,
  ]);

  return playRandom;
}

function useAddKeybinds(filter: ListFilterModel, count: number) {
  const playRandom = usePlayRandom(filter, count);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      playRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [playRandom]);
}

const SceneList: React.FC<{
  scenes: GQL.SlimSceneDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
}> = ({ scenes, filter, selectedIds, onSelectChange, fromGroupId }) => {
  const queue = useMemo(() => SceneQueue.fromListFilterModel(filter), [filter]);

  if (scenes.length === 0) {
    return null;
  }

  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <SceneCardsGrid
        scenes={scenes}
        queue={queue}
        zoomIndex={filter.zoomIndex}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
        fromGroupId={fromGroupId}
      />
    );
  }
  if (filter.displayMode === DisplayMode.List) {
    return (
      <SceneListTable
        scenes={scenes}
        queue={queue}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
      />
    );
  }
  if (filter.displayMode === DisplayMode.Wall) {
    return <SceneWallPanel scenes={scenes} sceneQueue={queue} />;
  }
  if (filter.displayMode === DisplayMode.Tagger) {
    return <Tagger scenes={scenes} queue={queue} />;
  }

  return null;
};

interface IFilteredScenes {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
  fromGroupId?: string;
}

export const FilteredSceneList = (props: IFilteredScenes) => {
  const intl = useIntl();
  const history = useHistory();

  const { filterHook, defaultSort, view, alterQuery, fromGroupId } = props;

  // States
  const [showSidebar, setShowSidebar] = useState(true);

  const { filterState, queryResult, modalState, listSelect, showEditFilter } =
    useFilteredItemList({
      filterStateProps: {
        filterMode: GQL.FilterMode.Scenes,
        defaultSort,
        view,
        useURL: alterQuery,
      },
      queryResultProps: {
        useResult: useFindScenes,
        getCount: (r) => r.data?.findScenes.count ?? 0,
        getItems: (r) => r.data?.findScenes.scenes ?? [],
        filterHook,
      },
    });

  const { filter, setFilter, loading: filterLoading } = filterState;

  const { effectiveFilter, result, cachedResult, items, totalCount } =
    queryResult;

  const {
    selectedIds,
    selectedItems,
    onSelectChange,
    onSelectNone,
    hasSelection,
  } = listSelect;

  const { modal, showModal, closeModal } = modalState;

  // Utility hooks
  const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
    filter,
    setFilter,
  });

  useAddKeybinds(filter, totalCount);

  const onCloseEditDelete = useCloseEditDelete({
    closeModal,
    onSelectNone,
    result,
  });

  const metadataByline = useMemo(() => {
    if (cachedResult.loading) return "";

    return renderMetadataByline(cachedResult) ?? "";
  }, [cachedResult]);

  const playSelected = usePlaySelected(selectedIds);
  const playRandom = usePlayRandom(filter, totalCount);

  function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          scenes: {
            ids: Array.from(selectedIds.values()),
            all: all,
          },
        }}
        onClose={() => closeModal()}
      />
    );
  }

  function onMerge() {
    const selected =
      selectedItems.map((s) => {
        return {
          id: s.id,
          title: objectTitle(s),
        };
      }) ?? [];
    showModal(
      <SceneMergeModal
        scenes={selected}
        onClose={(mergedID?: string) => {
          closeModal();
          if (mergedID) {
            history.push(`/scenes/${mergedID}`);
          }
        }}
        show
      />
    );
  }

  const otherOperations: IListFilterOperation[] = [
    {
      text: intl.formatMessage({ id: "actions.play_selected" }),
      onClick: playSelected,
      isDisplayed: () => hasSelection,
      icon: faPlay,
    },
    {
      text: intl.formatMessage({ id: "actions.play_random" }),
      onClick: playRandom,
    },
    {
      text: `${intl.formatMessage({ id: "actions.generate" })}…`,
      onClick: () =>
        showModal(
          <GenerateDialog
            type="scene"
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => closeModal()}
          />
        ),
      isDisplayed: () => hasSelection,
    },
    {
      text: `${intl.formatMessage({ id: "actions.identify" })}…`,
      onClick: () =>
        showModal(
          <IdentifyDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => closeModal()}
          />
        ),
      isDisplayed: () => hasSelection,
    },
    {
      text: `${intl.formatMessage({ id: "actions.merge" })}…`,
      onClick: () => onMerge(),
      isDisplayed: () => hasSelection,
    },
    {
      text: intl.formatMessage({ id: "actions.export" }),
      onClick: () => onExport(false),
      isDisplayed: () => hasSelection,
    },
    {
      text: intl.formatMessage({ id: "actions.export_all" }),
      onClick: () => onExport(true),
    },
  ];

  // render
  if (filterLoading) return null;

  return (
    <TaggerContext>
      <div className="item-list-container">
        {modal}

        <FilteredListToolbar
          filter={filter}
          setFilter={setFilter}
          showEditFilter={showEditFilter}
          view={view}
          listSelect={listSelect}
          onEdit={() =>
            showModal(
              <EditScenesDialog
                selected={selectedItems}
                onClose={onCloseEditDelete}
              />
            )
          }
          onDelete={() => {
            showModal(
              <DeleteScenesDialog
                selected={selectedItems}
                onClose={onCloseEditDelete}
              />
            );
          }}
          operations={otherOperations}
          zoomable
        />

        <SidebarPane>
          <Sidebar hide={!showSidebar}>
          </Sidebar>
          <div>
            <FilterTags
              criteria={filter.criteria}
              onEditCriterion={(c) => showEditFilter(c.criterionOption.type)}
              onRemoveCriterion={removeCriterion}
              onRemoveAll={() => clearAllCriteria()}
            />

            <PagedList
              result={result}
              cachedResult={cachedResult}
              filter={filter}
              totalCount={totalCount}
              onChangePage={setPage}
              metadataByline={metadataByline}
            >
              <SceneList
                filter={effectiveFilter}
                scenes={items}
                selectedIds={selectedIds}
                onSelectChange={onSelectChange}
                fromGroupId={fromGroupId}
              />
            </PagedList>
          </div>
        </SidebarPane>
      </div>
    </TaggerContext>
  );
};

export default FilteredSceneList;
