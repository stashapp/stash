import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { PersistanceLevel } from "../List/ItemList";
import { lazyComponent } from "src/utils/lazyComponent";

const SceneList = lazyComponent(() => import("./SceneList"));
const SceneMarkerList = lazyComponent(() => import("./SceneMarkerList"));
const Scene = lazyComponent(() => import("./SceneDetails/Scene"));
const SceneCreate = lazyComponent(() => import("./SceneDetails/SceneCreate"));

const Scenes: React.FC = () => {
  const titleProps = useTitleProps({ id: "scenes" });
  const markerTitleProps = useTitleProps({ id: "markers" });

  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route
          exact
          path="/scenes"
          render={(props) => (
            <SceneList persistState={PersistanceLevel.ALL} {...props} />
          )}
        />
        <Route
          exact
          path="/scenes/markers"
          render={() => (
            <>
              <Helmet {...markerTitleProps} />
              <SceneMarkerList />
            </>
          )}
        />
        <Route exact path="/scenes/new" component={SceneCreate} />
        <Route path="/scenes/:id" component={Scene} />
      </Switch>
    </>
  );
};
export default Scenes;
