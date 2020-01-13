import React, { FunctionComponent } from "react";
import { IBaseProps } from "../../models/base-props";
import { SceneList } from "./SceneList";

interface ISceneListPageProps extends IBaseProps {}

export const SceneListPage: FunctionComponent<ISceneListPageProps> = (props: ISceneListPageProps) => {
  return <SceneList base={props}/>;
};
