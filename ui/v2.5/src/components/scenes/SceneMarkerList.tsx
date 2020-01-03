import _ from "lodash";
import React, { FunctionComponent } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindSceneMarkersQuery, FindSceneMarkersVariables } from "../../core/generated-graphql";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";
import { StashService } from "../../core/StashService";
import { NavigationUtils } from "../../utils/navigation";

interface IProps extends IBaseProps {}

export const SceneMarkerList: FunctionComponent<IProps> = (props: IProps) => {
  const otherOperations = [
    {
      text: "Play Random",
      onClick: playRandom,
    }
  ];

  const listData = ListHook.useList({
    filterMode: FilterMode.SceneMarkers,
    otherOperations: otherOperations,
    props,
    renderContent,
  });

  async function playRandom(result: QueryHookResult<FindSceneMarkersQuery, FindSceneMarkersVariables>, filter: ListFilterModel, selectedIds: Set<string>) {
    // query for a random scene
    if (result.data && result.data.findSceneMarkers) {
      let count = result.data.findSceneMarkers.count;

      let index = Math.floor(Math.random() * count);
      let filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await StashService.queryFindSceneMarkers(filterCopy);
      if (singleResult && singleResult.data && singleResult.data.findSceneMarkers && singleResult.data.findSceneMarkers.scene_markers.length === 1) {
        // navigate to the scene player page
        let url = NavigationUtils.makeSceneMarkerUrl(singleResult!.data!.findSceneMarkers!.scene_markers[0])
        props.history.push(url);
      }
    }
  }

  function renderContent(
    result: QueryHookResult<FindSceneMarkersQuery, FindSceneMarkersVariables>,
    filter: ListFilterModel,
  ) {
    if (!result.data || !result.data.findSceneMarkers) { return; }
    if (filter.displayMode === DisplayMode.Wall) {
      return <WallPanel sceneMarkers={result.data.findSceneMarkers.scene_markers} />;
    }
  }

  return listData.template;
};
