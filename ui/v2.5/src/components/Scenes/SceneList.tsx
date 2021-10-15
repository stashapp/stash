import React, { useState } from "react";
import _ from "lodash";
import { useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import { IconProp } from "@fortawesome/fontawesome-svg-core";
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
import { SceneGenerateDialog } from "./SceneGenerateDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { SceneCardsGrid } from "./SceneCardsGrid";
import { TaggerContext } from "../Tagger/context";

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
  const [isGenerateDialogOpen, setIsGenerateDialogOpen] = useState(false);
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.play_selected" }),
      onClick: playSelected,
      isDisplayed: showWhenSelected,
      icon: "play" as IconProp,
    },
    {
      text: intl.formatMessage({ id: "actions.play_random" }),
      onClick: playRandom,
    },
    {
      text: intl.formatMessage({ id: "actions.generate" }),
      onClick: generate,
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
    queue.playScene(history, sceneIDs[0], { autoPlay: true });
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
      const index = Math.floor(Math.random() * filter.itemsPerPage);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.currentPage = page;
      filterCopy.sortBy = "random";
      const queryResults = await queryFindScenes(filterCopy);
      if (queryResults.data.findScenes.scenes.length > index) {
        const { id } = queryResults!.data!.findScenes!.scenes[index];
        // navigate to the image player page
        const queue = SceneQueue.fromListFilterModel(filterCopy);
        queue.playScene(history, id, { sceneIndex: index, autoPlay: true });
      }
    }
  }

  async function generate() {
    setIsGenerateDialogOpen(true);
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
          <SceneGenerateDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => {
              setIsGenerateDialogOpen(false);
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
        {maybeRenderSceneExportDialog(selectedIds)}
        {renderScenes(result, filter, selectedIds)}
      </>
    );
  }

  return <TaggerContext>{listData.template}</TaggerContext>;
};
