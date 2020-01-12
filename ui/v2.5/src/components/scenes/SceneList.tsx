import React from "react";
import _ from "lodash";
import { QueryHookResult } from "react-apollo-hooks";
import { useHistory } from 'react-router-dom';
import { FindScenesQuery, FindScenesVariables, SlimSceneDataFragment } from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { useScenesList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";
import { SceneCard } from "./SceneCard";
import { SceneListTable } from "./SceneListTable";
import { SceneSelectedOptions } from "./SceneSelectedOptions";

export const SceneList: React.FC = () => {
  const history = useHistory();
  const otherOperations = [
    {
      text: "Play Random",
      onClick: playRandom,
    }
  ];

  const listData = useScenesList({
    zoomable: true,
    otherOperations: otherOperations,
    renderContent,
    renderSelectedOptions
  });

  async function playRandom(result: QueryHookResult<FindScenesQuery, FindScenesVariables>, filter: ListFilterModel) {
    // query for a random scene
    if (result.data && result.data.findScenes) {
      let count = result.data.findScenes.count;

      let index = Math.floor(Math.random() * count);
      let filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await StashService.queryFindScenes(filterCopy);
      if (singleResult && singleResult.data && singleResult.data.findScenes && singleResult.data.findScenes.scenes.length === 1) {
        let id = singleResult!.data!.findScenes!.scenes[0].id;
        // navigate to the scene player page
        history.push("/scenes/" + id + "?autoplay=true");
      }
    }
  }

  function renderSelectedOptions(result: QueryHookResult<FindScenesQuery, FindScenesVariables>, selectedIds: Set<string>) {
    // find the selected items from the ids
    if (!result.data || !result.data.findScenes) { return undefined; }

    var scenes = result.data.findScenes.scenes;

    var selectedScenes : SlimSceneDataFragment[] = [];
    selectedIds.forEach((id) => {
      var scene = scenes.find((scene) => {
        return scene.id === id;
      });

      if (scene) {
        selectedScenes.push(scene);
      }
    });

    return (
      <>
      <SceneSelectedOptions selected={selectedScenes} onScenesUpdated={() => { return; }}/>
      </>
    );
  }

  function renderSceneCard(scene : SlimSceneDataFragment, selectedIds: Set<string>, zoomIndex: number) {
    return (
      <SceneCard
        key={scene.id}
        scene={scene}
        zoomIndex={zoomIndex}
        selected={selectedIds.has(scene.id)}
        onSelectedChanged={(selected: boolean, shiftKey: boolean) => listData.onSelectChange(scene.id, selected, shiftKey)}
      />
    )
  }

  function renderContent(result: QueryHookResult<FindScenesQuery, FindScenesVariables>, filter: ListFilterModel, selectedIds: Set<string>, zoomIndex: number) {
    if (!result.data || !result.data.findScenes) { return; }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="grid">
          {result.data.findScenes.scenes.map((scene) => renderSceneCard(scene, selectedIds, zoomIndex))}
        </div>
      );
    } else if (filter.displayMode === DisplayMode.List) {
      return <SceneListTable scenes={result.data.findScenes.scenes}/>;
    } else if (filter.displayMode === DisplayMode.Wall) {
      return <WallPanel scenes={result.data.findScenes.scenes} />;
    }
  }

  return listData.template;
};
