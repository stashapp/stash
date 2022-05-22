import React, { useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import {
  FindScenesQueryResult,
  SlimSceneDataFragment,
} from "src/core/generated-graphql";
import { queryFindScenes } from "src/core/StashService";
import { useScenesList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { showWhenSelected, PersistanceLevel } from "src/hooks/ListHook";
import Tagger from "src/components/Tagger";
import { SceneQueue } from "src/models/sceneQueue";
import { WallPanel } from "../Wall/WallPanel";
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

interface ISceneList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  persistState?: PersistanceLevel.ALL;
}

export const SceneList: React.FC<ISceneList> = ({
  filterHook,
  defaultSort,
  persistState,
}) => {
  const intl = useIntl();
  const history = useHistory();
  const config = React.useContext(ConfigurationContext);
  const [isGenerateDialogOpen, setIsGenerateDialogOpen] = useState(false);
  const [isIdentifyDialogOpen, setIsIdentifyDialogOpen] = useState(false);
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

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
      onClick: generate,
      isDisplayed: showWhenSelected,
    },
    {
      text: `${intl.formatMessage({ id: "actions.identify" })}…`,
      onClick: identify,
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

  const addKeybinds = (
    result: FindScenesQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      playRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  const renderDeleteDialog = (
    selectedScenes: SlimSceneDataFragment[],
    onClose: (confirmed: boolean) => void
  ) => <DeleteScenesDialog selected={selectedScenes} onClose={onClose} />;

  const listData = useScenesList({
    zoomable: true,
    selectable: true,
    otherOperations,
    defaultSort,
    renderContent,
    renderEditDialog: renderEditScenesDialog,
    renderDeleteDialog,
    filterHook,
    addKeybinds,
    persistState,
  });

  async function playSelected(
    result: FindScenesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    // populate queue and go to first scene
    const sceneIDs = Array.from(selectedIds.values());
    const queue = SceneQueue.fromSceneIDList(sceneIDs);
    const autoPlay =
      config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
    const cont =
      config.configuration?.interface.continuePlaylistDefault ?? false;
    queue.playScene(history, sceneIDs[0], { autoPlay, continue: cont });
  }

  async function playRandom(
    result: FindScenesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random scene
    if (result.data && result.data.findScenes) {
      const { count } = result.data.findScenes;

      const pages = Math.ceil(count / filter.itemsPerPage);
      const page = Math.floor(Math.random() * pages) + 1;

      const indexMax =
        filter.itemsPerPage < count ? filter.itemsPerPage : count;
      const index = Math.floor(Math.random() * indexMax);
      const filterCopy = cloneDeep(filter);
      filterCopy.currentPage = page;
      filterCopy.sortBy = "random";
      const queryResults = await queryFindScenes(filterCopy);
      if (queryResults.data.findScenes.scenes.length > index) {
        const { id } = queryResults!.data!.findScenes!.scenes[index];
        // navigate to the image player page
        const queue = SceneQueue.fromListFilterModel(filterCopy);
        const autoPlay =
          config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
        const cont =
          config.configuration?.interface.continuePlaylistDefault ?? false;
        queue.playScene(history, id, {
          sceneIndex: index,
          autoPlay,
          continue: cont,
        });
      }
    }
  }

  async function generate() {
    setIsGenerateDialogOpen(true);
  }

  async function identify() {
    setIsIdentifyDialogOpen(true);
  }

  async function onExport() {
    setIsExportAll(false);
    setIsExportDialogOpen(true);
  }

  async function onExportAll() {
    setIsExportAll(true);
    setIsExportDialogOpen(true);
  }

  function maybeRenderSceneGenerateDialog(selectedIds: Set<string>) {
    if (isGenerateDialogOpen) {
      return (
        <>
          <GenerateDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => {
              setIsGenerateDialogOpen(false);
            }}
          />
        </>
      );
    }
  }

  function maybeRenderSceneIdentifyDialog(selectedIds: Set<string>) {
    if (isIdentifyDialogOpen) {
      return (
        <>
          <IdentifyDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => {
              setIsIdentifyDialogOpen(false);
            }}
          />
        </>
      );
    }
  }

  function maybeRenderSceneExportDialog(selectedIds: Set<string>) {
    if (isExportDialogOpen) {
      return (
        <>
          <ExportDialog
            exportInput={{
              scenes: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => {
              setIsExportDialogOpen(false);
            }}
          />
        </>
      );
    }
  }

  function renderEditScenesDialog(
    selectedScenes: SlimSceneDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return (
      <>
        <EditScenesDialog selected={selectedScenes} onClose={onClose} />
      </>
    );
  }

  function renderScenes(
    result: FindScenesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    if (!result.data || !result.data.findScenes) {
      return;
    }

    const queue = SceneQueue.fromListFilterModel(filter);

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <SceneCardsGrid
          scenes={result.data.findScenes.scenes}
          queue={queue}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={(id, selected, shiftKey) =>
            listData.onSelectChange(id, selected, shiftKey)
          }
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <SceneListTable
          scenes={result.data.findScenes.scenes}
          queue={queue}
          selectedIds={selectedIds}
          onSelectChange={(id, selected, shiftKey) =>
            listData.onSelectChange(id, selected, shiftKey)
          }
        />
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <WallPanel scenes={result.data.findScenes.scenes} sceneQueue={queue} />
      );
    }
    if (filter.displayMode === DisplayMode.Tagger) {
      return <Tagger scenes={result.data.findScenes.scenes} queue={queue} />;
    }
  }

  function renderContent(
    result: FindScenesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    return (
      <>
        {maybeRenderSceneGenerateDialog(selectedIds)}
        {maybeRenderSceneIdentifyDialog(selectedIds)}
        {maybeRenderSceneExportDialog(selectedIds)}
        {renderScenes(result, filter, selectedIds)}
      </>
    );
  }

  return <TaggerContext>{listData.template}</TaggerContext>;
};

export default SceneList;
