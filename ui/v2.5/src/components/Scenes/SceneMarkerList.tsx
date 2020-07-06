import _ from "lodash";
import React from "react";
import { useHistory } from "react-router-dom";
import { FindSceneMarkersQueryResult } from "src/core/generated-graphql";
import { queryFindSceneMarkers } from "src/core/StashService";
import { NavUtils } from "src/utils";
import { useSceneMarkersList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";

interface ISceneMarkerList {
  subComponent?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const SceneMarkerList: React.FC<ISceneMarkerList> = ({
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
    result: FindSceneMarkersQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      playRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  const listData = useSceneMarkersList({
    otherOperations,
    renderContent,
    subComponent,
    filterHook,
    addKeybinds,
  });

  async function playRandom(
    result: FindSceneMarkersQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random scene
    if (result.data?.findSceneMarkers) {
      const { count } = result.data.findSceneMarkers;

      const index = Math.floor(Math.random() * count);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindSceneMarkers(filterCopy);
      if (singleResult?.data?.findSceneMarkers?.scene_markers?.length === 1) {
        // navigate to the scene player page
        const url = NavUtils.makeSceneMarkerUrl(
          singleResult.data.findSceneMarkers.scene_markers[0]
        );
        history.push(url);
      }
    }
  }

  function renderContent(
    result: FindSceneMarkersQueryResult,
    filter: ListFilterModel
  ) {
    if (!result?.data?.findSceneMarkers) return;
    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <WallPanel sceneMarkers={result.data.findSceneMarkers.scene_markers} />
      );
    }
  }

  return listData.template;
};
