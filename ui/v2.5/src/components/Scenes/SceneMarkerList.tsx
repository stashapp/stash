import cloneDeep from "lodash-es/cloneDeep";
import React from "react";
import { useHistory } from "react-router-dom";
import { useIntl } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindSceneMarkers,
  queryFindSceneMarkersByID,
  useFindSceneMarkers,
} from "src/core/StashService";
import { ItemList, ItemListContext, showWhenSelected } from "../List/ItemList";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { MarkerWallPanel } from "../Wall/WallPanel";
import { View } from "../List/views";
import { SceneMarkerCardsGrid } from "./SceneMarkerCardsGrid";
import { DeleteSceneMarkersDialog } from "./DeleteSceneMarkersDialog";
import { SceneQueue } from "src/models/sceneQueue";
import { ConfigurationContext } from "src/hooks/Config";
import { faPlay } from "@fortawesome/free-solid-svg-icons";

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
  defaultSort?: string;
}

export const SceneMarkerList: React.FC<ISceneMarkerList> = ({
  filterHook,
  view,
  alterQuery,
}) => {
  const intl = useIntl();
  const history = useHistory();
  const config = React.useContext(ConfigurationContext);
  const filterMode = GQL.FilterMode.SceneMarkers;

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

  async function playSelected(
    result: GQL.FindSceneMarkersQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    // populate queue and go to first scene
    const sceneMarkerIDs = Array.from(selectedIds.values());
    const queue = SceneQueue.fromSceneMarkerIDList(sceneMarkerIDs);
    const autoPlay =
      config.configuration?.interface.autostartVideoOnPlaySelected ?? false;

    // Fetch the first marker.
    const query = await queryFindSceneMarkersByID([Number(sceneMarkerIDs[0])]);
    const { scene_markers } = query.data.findSceneMarkers;
    const url = queue.makeLink(scene_markers[0].scene.id, {
      autoPlay: autoPlay,
      start: scene_markers[0].seconds,
      end: scene_markers[0].end_seconds,
    });

    history.push(url);
  }

  async function playRandom(
    result: GQL.FindSceneMarkersQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random scene
    if (result.data?.findSceneMarkers) {
      const { count } = result.data.findSceneMarkers;
      const pages = Math.ceil(count / filter.itemsPerPage);
      const page = Math.floor(Math.random() * pages) + 1;
      const indexMax = Math.min(filter.itemsPerPage, count);
      const index = Math.floor(Math.random() * indexMax);
      const filterCopy = cloneDeep(filter);
      filterCopy.currentPage = page;
      filterCopy.sortBy = "random";
      const queryResults = await queryFindSceneMarkers(filterCopy);
      const marker = queryResults.data.findSceneMarkers.scene_markers[index];
      if (marker) {
        const queue = SceneQueue.fromListFilterModel(filterCopy);
        const autoPlay =
          config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
        const url = queue.makeLink(marker.scene.id, {
          sceneIndex: index,
          continue: autoPlay,
          start: marker.seconds,
          end: marker.end_seconds,
          mode: "scene_marker",
        });
        history.push(url);
      }
    }
  }

  function renderContent(
    result: GQL.FindSceneMarkersQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void
  ) {
    if (!result.data?.findSceneMarkers) return;

    const queue = SceneQueue.fromListFilterModel(filter);

    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <MarkerWallPanel
          markers={result.data.findSceneMarkers.scene_markers}
          sceneQueue={queue}
        />
      );
    }

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <SceneMarkerCardsGrid
          markers={result.data.findSceneMarkers.scene_markers}
          queue={queue}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
  }

  function renderDeleteDialog(
    selectedSceneMarkers: GQL.SceneMarkerDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return (
      <DeleteSceneMarkersDialog
        selected={selectedSceneMarkers}
        onClose={onClose}
      />
    );
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
      selectable
    >
      <ItemList
        zoomable
        view={view}
        otherOperations={otherOperations}
        addKeybinds={addKeybinds}
        renderContent={renderContent}
        renderDeleteDialog={renderDeleteDialog}
      />
    </ItemListContext>
  );
};

export default SceneMarkerList;
