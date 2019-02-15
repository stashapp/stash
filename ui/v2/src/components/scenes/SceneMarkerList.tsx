import _ from "lodash";
import React, { FunctionComponent } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindSceneMarkersQuery, FindSceneMarkersVariables } from "../../core/generated-graphql";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";
import { WallPanel } from "../Wall/WallPanel";

interface IProps extends IBaseProps {}

export const SceneMarkerList: FunctionComponent<IProps> = (props: IProps) => {
  const listData = ListHook.useList({
    filterMode: FilterMode.SceneMarkers,
    props,
    renderContent,
  });

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
