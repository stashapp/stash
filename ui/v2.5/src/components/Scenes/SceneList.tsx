import React, { useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedNumber, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindScenes, useFindScenes } from "src/core/StashService";
import { ItemList, ItemListContext, showWhenSelected } from "../List/ItemList";
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

function getItems(result: GQL.FindScenesQueryResult) {
  return result?.data?.findScenes?.scenes ?? [];
}

function getCount(result: GQL.FindScenesQueryResult) {
  return result?.data?.findScenes?.count ?? 0;
}

function renderMetadataByline(result: GQL.FindScenesQueryResult) {
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
}

interface ISceneList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
}

export const SceneList: React.FC<ISceneList> = ({
  filterHook,
  defaultSort,
  view,
  alterQuery,
}) => {
  const intl = useIntl();
  const history = useHistory();
  const config = React.useContext(ConfigurationContext);
  const [isGenerateDialogOpen, setIsGenerateDialogOpen] = useState(false);
  const [mergeScenes, setMergeScenes] =
    useState<{ id: string; title: string }[]>();
  const [isIdentifyDialogOpen, setIsIdentifyDialogOpen] = useState(false);
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const filterMode = GQL.FilterMode.Scenes;

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.play_selected" }),
      onClick: playSelected,
      isDisplayed: showWhenSelected,
      icon: faPlay,
    },
    {
      text: intl.formatMessage({ id: "actions.play_random" }),
      onClick: playRandom,
    },
    {
      text: `${intl.formatMessage({ id: "actions.generate" })}…`,
      onClick: async () => setIsGenerateDialogOpen(true),
      isDisplayed: showWhenSelected,
    },
    {
      text: `${intl.formatMessage({ id: "actions.identify" })}…`,
      onClick: async () => setIsIdentifyDialogOpen(true),
      isDisplayed: showWhenSelected,
    },
    {
      text: `${intl.formatMessage({ id: "actions.merge" })}…`,
      onClick: onMerge,
      isDisplayed: showWhenSelected,
    },
    {
      text: intl.formatMessage({ id: "actions.export" }),
      onClick: onExport,
      isDisplayed: showWhenSelected,
    },
    {
      text: intl.formatMessage({ id: "actions.export_all" }),
      onClick: onExportAll,
    },
  ];

  function addKeybinds(
    result: GQL.FindScenesQueryResult,
    filter: ListFilterModel
  ) {
    Mousetrap.bind("p r", () => {
      playRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }

  function playScene(
    queue: SceneQueue,
    sceneID: string,
    options: IPlaySceneOptions
  ) {
    history.push(queue.makeLink(sceneID, options));
  }

  async function playSelected(
    result: GQL.FindScenesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    // populate queue and go to first scene
    const sceneIDs = Array.from(selectedIds.values());
    const queue = SceneQueue.fromSceneIDList(sceneIDs);
    const autoPlay =
      config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
    playScene(queue, sceneIDs[0], { autoPlay });
  }

  async function playRandom(
    result: GQL.FindScenesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random scene
    if (result.data?.findScenes) {
      const { count } = result.data.findScenes;

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
        const autoPlay =
          config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
        playScene(queue, scene.id, { sceneIndex: index, autoPlay });
      }
    }
  }

  async function onMerge(
    result: GQL.FindScenesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    const selected =
      result.data?.findScenes.scenes
        .filter((s) => selectedIds.has(s.id))
        .map((s) => {
          return {
            id: s.id,
            title: objectTitle(s),
          };
        }) ?? [];

    setMergeScenes(selected);
  }

  async function onExport() {
    setIsExportAll(false);
    setIsExportDialogOpen(true);
  }

  async function onExportAll() {
    setIsExportAll(true);
    setIsExportDialogOpen(true);
  }

  function renderContent(
    result: GQL.FindScenesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void
  ) {
    function maybeRenderSceneGenerateDialog() {
      if (isGenerateDialogOpen) {
        return (
          <GenerateDialog
            type="scene"
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => setIsGenerateDialogOpen(false)}
          />
        );
      }
    }

    function maybeRenderSceneIdentifyDialog() {
      if (isIdentifyDialogOpen) {
        return (
          <IdentifyDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => setIsIdentifyDialogOpen(false)}
          />
        );
      }
    }

    function maybeRenderSceneExportDialog() {
      if (isExportDialogOpen) {
        return (
          <ExportDialog
            exportInput={{
              scenes: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => setIsExportDialogOpen(false)}
          />
        );
      }
    }

    function renderMergeDialog() {
      if (mergeScenes) {
        return (
          <SceneMergeModal
            scenes={mergeScenes}
            onClose={(mergedID?: string) => {
              setMergeScenes(undefined);
              if (mergedID) {
                history.push(`/scenes/${mergedID}`);
              }
            }}
            show
          />
        );
      }
    }

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
            onSelectChange={onSelectChange}
          />
        );
      }
      if (filter.displayMode === DisplayMode.List) {
        return (
          <SceneListTable
            scenes={result.data.findScenes.scenes}
            queue={queue}
            selectedIds={selectedIds}
            onSelectChange={onSelectChange}
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
        return <Tagger scenes={result.data.findScenes.scenes} queue={queue} />;
      }
    }

    return (
      <>
        {maybeRenderSceneGenerateDialog()}
        {maybeRenderSceneIdentifyDialog()}
        {maybeRenderSceneExportDialog()}
        {renderMergeDialog()}
        {renderScenes()}
      </>
    );
  }

  function renderEditDialog(
    selectedScenes: GQL.SlimSceneDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return <EditScenesDialog selected={selectedScenes} onClose={onClose} />;
  }

  function renderDeleteDialog(
    selectedScenes: GQL.SlimSceneDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return <DeleteScenesDialog selected={selectedScenes} onClose={onClose} />;
  }

  return (
    <TaggerContext>
      <ItemListContext
        filterMode={filterMode}
        defaultSort={defaultSort}
        useResult={useFindScenes}
        getItems={getItems}
        getCount={getCount}
        alterQuery={alterQuery}
        filterHook={filterHook}
        view={view}
        selectable
      >
        <ItemList
          zoomable
          view={view}
          otherOperations={otherOperations}
          addKeybinds={addKeybinds}
          renderContent={renderContent}
          renderEditDialog={renderEditDialog}
          renderDeleteDialog={renderDeleteDialog}
          renderMetadataByline={renderMetadataByline}
        />
      </ItemListContext>
    </TaggerContext>
  );
};

export default SceneList;
