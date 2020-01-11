import _ from "lodash";
import React from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindSceneMarkersQuery, FindSceneMarkersVariables } from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { NavUtils } from "src/utils";
import { ListHook } from "src/hooks";
import { IBaseProps } from "src/models/base-props";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode, FilterMode } from "src/models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";

interface IProps extends IBaseProps {}

export const SceneMarkerList: React.FC<IProps> = (props: IProps) => {
  const otherOperations = [{
    text: "Play Random",
    onClick: playRandom
  }];

  const listData = ListHook.useList({
    filterMode: FilterMode.SceneMarkers,
    otherOperations: otherOperations,
    props,
    renderContent,
  });

  async function playRandom(result: QueryHookResult<FindSceneMarkersQuery, FindSceneMarkersVariables>, filter: ListFilterModel) {
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
        let url = NavUtils.makeSceneMarkerUrl(singleResult.data.findSceneMarkers.scene_markers[0])
        props.history.push(url);
      }
    }
  }

  function renderContent(
    result: QueryHookResult<FindSceneMarkersQuery, FindSceneMarkersVariables>,
    filter: ListFilterModel,
  ) {
    if (!result?.data?.findSceneMarkers)
      return;
    if (filter.displayMode === DisplayMode.Wall) {
      return <WallPanel sceneMarkers={result.data.findSceneMarkers.scene_markers} />;
    }
  }

  return listData.template;
};
