import React from "react";
import { Route, Switch } from "react-router-dom";
import { PersistanceLevel } from "src/hooks/ListHook";
import Scene from "./SceneDetails/Scene";
import { SceneList } from "./SceneList";
import { SceneMarkerList } from "./SceneMarkerList";

const Scenes = () => (
  <Switch>
    <Route
      exact
      path="/scenes"
      render={(props) => (
        <SceneList persistState={PersistanceLevel.ALL} {...props} />
      )}
    />
    <Route exact path="/scenes/markers" component={SceneMarkerList} />
    <Route path="/scenes/:id" component={Scene} />
  </Switch>
);

export default Scenes;
