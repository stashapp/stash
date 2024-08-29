import cloneDeep from "lodash-es/cloneDeep";
import React from "react";
import { useHistory } from "react-router-dom";
import { useIntl } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindSceneMarkers,
  useFindSceneMarkers,
} from "src/core/StashService";
import NavUtils from "src/utils/navigation";
import { ItemList, ItemListContext } from "../List/ItemList";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { MarkerWallPanel } from "../Wall/WallPanel";
import { View } from "../List/views";

function getItems(result: GQL.FindSceneMarkersQueryResult) {
  return result?.data?.findSceneMarkers?.scene_markers ?? [];
}

function getCount(result: GQL.FindSceneMarkersQueryResult) {
  return result?.data?.findSceneMarkers?.count ?? 0;
}

interface ISceneMarkerList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
}

export const SceneMarkerList: React.FC<ISceneMarkerList> = ({
  filterHook,
  view,
  alterQuery,
}) => {
  const intl = useIntl();
  const history = useHistory();

  const filterMode = GQL.FilterMode.SceneMarkers;

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.play_random" }),
      onClick: playRandom,
    },
  ];

  function addKeybinds(
    result: GQL.FindSceneMarkersQueryResult,
    filter: ListFilterModel
  ) {
    Mousetrap.bind("p r", () => {
      playRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }

  async function playRandom(
    result: GQL.FindSceneMarkersQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random scene
    if (result.data?.findSceneMarkers) {
      const { count } = result.data.findSceneMarkers;

      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindSceneMarkers(filterCopy);
      if (singleResult.data.findSceneMarkers.scene_markers.length === 1) {
        // navigate to the scene player page
        const url = NavUtils.makeSceneMarkerUrl(
          singleResult.data.findSceneMarkers.scene_markers[0]
        );
        history.push(url);
      }
    }
  }

  function renderContent(
    result: GQL.FindSceneMarkersQueryResult,
    filter: ListFilterModel
  ) {
    if (!result.data?.findSceneMarkers) return;

    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <MarkerWallPanel markers={result.data.findSceneMarkers.scene_markers} />
      );
    }
  }

  return (
    <ItemListContext
      filterMode={filterMode}
      useResult={useFindSceneMarkers}
      getItems={getItems}
      getCount={getCount}
      alterQuery={alterQuery}
      filterHook={filterHook}
      view={view}
    >
      <ItemList
        view={view}
        otherOperations={otherOperations}
        addKeybinds={addKeybinds}
        renderContent={renderContent}
      />
    </ItemListContext>
  );
};

export default SceneMarkerList;
