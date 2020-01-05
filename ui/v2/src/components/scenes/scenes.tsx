import React from "react";
import { Route, Switch } from "react-router-dom";
import { Scene } from "./SceneDetails/Scene";
import { SceneMarkerList } from "./SceneMarkerList";
import { SceneListPage } from "./SceneListPage";

const Scenes = () => (
  <Switch>
    <Route exact={true} path="/scenes" component={SceneListPage} />
    <Route exact={true} path="/scenes/markers" component={SceneMarkerList} />
    <Route path="/scenes/:id" component={Scene} />
  </Switch>
);

export default Scenes;
