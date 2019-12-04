import _ from "lodash";
import React, { FunctionComponent } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindScenesQuery, FindScenesVariables, SlimSceneDataFragment } from "../../core/generated-graphql";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";
import { SceneCard } from "./SceneCard";
import { SceneListTable } from "./SceneListTable";
import { SceneSelectedOptions } from "./SceneSelectedOptions";

interface ISceneListProps extends IBaseProps {}

export const SceneList: FunctionComponent<ISceneListProps> = (props: ISceneListProps) => {
  const listData = ListHook.useList({
    filterMode: FilterMode.Scenes,
    props,
    zoomable: true,
    selectable: true,
    renderContent,
    renderSelectedOptions
  });

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
