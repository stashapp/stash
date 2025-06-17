import React, { useCallback, useContext, useEffect, useMemo } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindScenes, useFindScenes } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { Tagger } from "../Tagger/scenes/SceneTagger";
import { IPlaySceneOptions, SceneQueue } from "src/models/sceneQueue";
import { SceneWallPanel } from "./SceneWallPanel";
import { SceneListTable } from "./SceneListTable";
import { EditScenesDialog } from "./EditScenesDialog";
import { DeleteScenesDialog } from "./DeleteScenesDialog";
import { GenerateDialog } from "../Dialogs/GenerateDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { SceneCardsGrid } from "./SceneCardsGrid";
import { TaggerContext } from "../Tagger/context";
import { IdentifyDialog } from "../Dialogs/IdentifyDialog/IdentifyDialog";
import { ConfigurationContext } from "src/hooks/Config";
import { faEllipsisH, faPlay, faPlus } from "@fortawesome/free-solid-svg-icons";
import { SceneMergeModal } from "./SceneMergeDialog";
import { objectTitle } from "src/core/files";
import TextUtils from "src/utils/text";
import { View } from "../List/views";
import { FileSize } from "../Shared/FileSize";
import { LoadedContent } from "../List/PagedList";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import { IListFilterOperation } from "../List/ListOperationButtons";
import { useFilteredItemList } from "../List/ItemList";
import { FilterTags } from "../List/FilterTags";
import {
  Sidebar,
  SidebarIcon,
  SidebarPane,
  useSidebarState,
} from "../Shared/Sidebar";
import { SidebarPerformersFilter } from "../List/Filters/PerformersFilter";
import { SidebarStudiosFilter } from "../List/Filters/StudiosFilter";
import { PerformersCriterionOption } from "src/models/list-filter/criteria/performers";
import { StudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import { TagsCriterionOption } from "src/models/list-filter/criteria/tags";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import cx from "classnames";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { OrganizedCriterionOption } from "src/models/list-filter/criteria/organized";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/FilterSidebar";
import { PatchContainerComponent } from "src/patch";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { Button, ButtonGroup, ButtonToolbar } from "react-bootstrap";
import { FilterButton } from "../List/Filters/FilterButton";
import { Icon } from "../Shared/Icon";

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

  const { configuration: config } = useContext(ConfigurationContext);
  const cont = config?.interface.continuePlaylistDefault ?? false;
  const autoPlay = config?.interface.autostartVideoOnPlaySelected ?? false;

  const playScene = useCallback(
    (queue: SceneQueue, sceneID: string, options?: IPlaySceneOptions) => {
      history.push(
        queue.makeLink(sceneID, { autoPlay, continue: cont, ...options })
      );
    },
    [history, cont, autoPlay]
  );

  return playScene;
}

function usePlaySelected(selectedIds: Set<string>) {
  const playScene = usePlayScene();

  const playSelected = useCallback(() => {
    // populate queue and go to first scene
    const sceneIDs = Array.from(selectedIds.values());
    const queue = SceneQueue.fromSceneIDList(sceneIDs);

    playScene(queue, sceneIDs[0]);
  }, [selectedIds, playScene]);

  return playSelected;
}

function usePlayFirst() {
  const playScene = usePlayScene();

  const playFirst = useCallback(
    (queue: SceneQueue, sceneID: string, index: number) => {
      // populate queue and go to first scene
      playScene(queue, sceneID, { sceneIndex: index });
    },
    [playScene]
  );

  return playFirst;
}

function usePlayRandom(filter: ListFilterModel, count: number) {
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
      playScene(queue, scene.id, { sceneIndex: index });
    }
  }, [filter, count, playScene]);

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

const ScenesFilterSidebarSections = PatchContainerComponent(
  "FilteredSceneList.SidebarSections"
);

const SidebarContent: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  view?: View;
  sidebarOpen: boolean;
  onClose?: () => void;
  showEditFilter: (editingCriterion?: string) => void;
}> = ({ filter, setFilter, view, showEditFilter, sidebarOpen, onClose }) => {
  return (
    <>
      <FilteredSidebarHeader
        sidebarOpen={sidebarOpen}
        onClose={onClose}
        showEditFilter={showEditFilter}
        filter={filter}
        setFilter={setFilter}
        view={view}
      />

      <ScenesFilterSidebarSections>
        <SidebarStudiosFilter
          title={<FormattedMessage id="studios" />}
          data-type={StudiosCriterionOption.type}
          option={StudiosCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarPerformersFilter
          title={<FormattedMessage id="performers" />}
          data-type={PerformersCriterionOption.type}
          option={PerformersCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarTagsFilter
          title={<FormattedMessage id="tags" />}
          data-type={TagsCriterionOption.type}
          option={TagsCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarRatingFilter
          title={<FormattedMessage id="rating" />}
          data-type={RatingCriterionOption.type}
          option={RatingCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
        <SidebarBooleanFilter
          title={<FormattedMessage id="organized" />}
          data-type={OrganizedCriterionOption.type}
          option={OrganizedCriterionOption}
          filter={filter}
          setFilter={setFilter}
        />
      </ScenesFilterSidebarSections>
    </>
  );
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
  const {
    showSidebar,
    setShowSidebar,
    loading: sidebarStateLoading,
  } = useSidebarState(view);

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
  useFilteredSidebarKeybinds({
    showSidebar,
    setShowSidebar,
  });

  const onCloseEditDelete = useCloseEditDelete({
    closeModal,
    onSelectNone,
    result,
  });

  const metadataByline = useMemo(() => {
    if (cachedResult.loading) return "";

    return renderMetadataByline(cachedResult) ?? "";
  }, [cachedResult]);

  const queue = useMemo(() => SceneQueue.fromListFilterModel(filter), [filter]);

  const playSelected = usePlaySelected(selectedIds);
  const playRandom = usePlayRandom(filter, totalCount);
  const playFirst = usePlayFirst();

  function onPlay() {
    if (items.length === 0) {
      return;
    }

    // if there are selected items, play those
    if (selectedIds.size > 0) {
      playSelected();
      return;
    }

    // otherwise, play the first item in the list
    const sceneID = items[0].id;
    playFirst(queue, sceneID, 0);
  }

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
  if (filterLoading || sidebarStateLoading) return null;

  return (
    <TaggerContext>
      <div
        className={cx("item-list-container scene-list", {
          "hide-sidebar": !showSidebar,
        })}
      >
        {modal}

        <SidebarPane hideSidebar={!showSidebar}>
          <Sidebar hide={!showSidebar} onHide={() => setShowSidebar(false)}>
            <SidebarContent
              filter={filter}
              setFilter={setFilter}
              showEditFilter={showEditFilter}
              view={view}
              sidebarOpen={showSidebar}
              onClose={() => setShowSidebar(false)}
            />
          </Sidebar>
          <div>
            <ButtonToolbar className="scene-list-toolbar">
              <div>
                <FilterButton
                  onClick={() => showEditFilter()}
                  count={filter.count()}
                />
              </div>
              <div>
                <ButtonGroup>
                  {!!items.length && (
                    <Button 
                      variant="secondary" 
                      onClick={() => onPlay()} 
                      title={intl.formatMessage({ id: "actions.play" })}
                    >
                      <Icon icon={faPlay} />
                    </Button>
                  )}
                  <Button variant="secondary">
                    <Icon icon={faPlus} />
                  </Button>
                  <Button variant="secondary">
                    <Icon icon={faEllipsisH} />
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => setShowSidebar(!showSidebar)}
                  >
                    <SidebarIcon />
                  </Button>
                </ButtonGroup>
              </div>
            </ButtonToolbar>

            <FilterTags
              criteria={filter.criteria}
              onEditCriterion={(c) => showEditFilter(c.criterionOption.type)}
              onRemoveCriterion={removeCriterion}
              onRemoveAll={() => clearAllCriteria()}
            />

            <PaginationIndex
              loading={cachedResult.loading}
              itemsPerPage={filter.itemsPerPage}
              currentPage={filter.currentPage}
              totalItems={totalCount}
              metadataByline={metadataByline}
            />

            <LoadedContent loading={result.loading} error={result.error}>
              <SceneList
                filter={effectiveFilter}
                scenes={items}
                selectedIds={selectedIds}
                onSelectChange={onSelectChange}
                fromGroupId={fromGroupId}
              />
            </LoadedContent>

            {totalCount > filter.itemsPerPage && (
              <div className="pagination-footer">
                <Pagination
                  itemsPerPage={filter.itemsPerPage}
                  currentPage={filter.currentPage}
                  totalItems={totalCount}
                  metadataByline={metadataByline}
                  onChangePage={setPage}
                  pagePopupPlacement="top"
                />
              </div>
            )}
          </div>
        </SidebarPane>
      </div>
    </TaggerContext>
  );
};

export default FilteredSceneList;
