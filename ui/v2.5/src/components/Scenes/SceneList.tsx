import React from "react";
import _ from "lodash";
import { useHistory } from "react-router-dom";
import {
  FindScenesQueryResult,
  SlimSceneDataFragment,
} from "src/core/generated-graphql";
import { queryFindScenes } from "src/core/StashService";
import { useScenesList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";
import { SceneCard } from "./SceneCard";
import { SceneListTable } from "./SceneListTable";
import { EditScenesDialog } from "./EditScenesDialog";
import { DeleteScenesDialog } from "./DeleteScenesDialog";

interface ISceneList {
  subComponent?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const SceneList: React.FC<ISceneList> = ({
  subComponent,
  filterHook,
}) => {
  const history = useHistory();
  const otherOperations = [
    {
      text: "Play Random",
      onClick: playRandom,
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

  const listData = useScenesList({
    zoomable: true,
    otherOperations,
    renderContent,
    renderEditDialog: renderEditScenesDialog,
    renderDeleteDialog: renderDeleteScenesDialog,
    subComponent,
    filterHook,
    addKeybinds,
  });

  async function playRandom(
    result: FindScenesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random scene
    if (result.data && result.data.findScenes) {
      const { count } = result.data.findScenes;

      const index = Math.floor(Math.random() * count);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindScenes(filterCopy);
      if (
        singleResult &&
        singleResult.data &&
        singleResult.data.findScenes &&
        singleResult.data.findScenes.scenes.length === 1
      ) {
        const { id } = singleResult!.data!.findScenes!.scenes[0];
        // navigate to the scene player page
        history.push(`/scenes/${id}?autoplay=true`);
      }
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

  function renderDeleteScenesDialog(
    selectedScenes: SlimSceneDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return (
      <>
        <DeleteScenesDialog selected={selectedScenes} onClose={onClose} />
      </>
    );
  }

  function renderSceneCard(
    scene: SlimSceneDataFragment,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    return (
      <SceneCard
        key={scene.id}
        scene={scene}
        zoomIndex={zoomIndex}
        selecting={selectedIds.size > 0}
        selected={selectedIds.has(scene.id)}
        onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
          listData.onSelectChange(scene.id, selected, shiftKey)
        }
      />
    );
  }

  function renderContent(
    result: FindScenesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    if (!result.data || !result.data.findScenes) {
      return;
    }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row justify-content-center">
          {result.data.findScenes.scenes.map((scene) =>
            renderSceneCard(scene, selectedIds, zoomIndex)
          )}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return <SceneListTable scenes={result.data.findScenes.scenes} />;
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return <WallPanel scenes={result.data.findScenes.scenes} />;
    }
  }

  return listData.template;
};
