import cloneDeep from "lodash-es/cloneDeep";
import React from "react";
import { useHistory } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared";
import Mousetrap from "mousetrap";
import { FindSceneMarkersQueryResult } from "src/core/generated-graphql";
import { queryFindSceneMarkers } from "src/core/StashService";
import { NavUtils } from "src/utils";
import { useSceneMarkersList } from "src/hooks";
import { PersistanceLevel } from "src/hooks/ListHook";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";

interface ISceneMarkerList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const SceneMarkerList: React.FC<ISceneMarkerList> = ({ filterHook }) => {
  const intl = useIntl();
  const history = useHistory();
  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.play_random" }),
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
    filterHook,
    addKeybinds,
    persistState: PersistanceLevel.ALL,
  });

  async function playRandom(
    result: FindSceneMarkersQueryResult,
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
  const title_template = `${intl.formatMessage({
    id: "markers",
  })} ${TITLE_SUFFIX}`;

  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />

      {listData.template}
    </>
  );
};

export default SceneMarkerList;
