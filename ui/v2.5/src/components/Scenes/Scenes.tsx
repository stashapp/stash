import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { PersistanceLevel } from "../List/ItemList";
import { lazyComponent } from "src/utils/lazyComponent";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";

const SceneList = lazyComponent(() => import("./SceneList"));
const SceneMarkerList = lazyComponent(() => import("./SceneMarkerList"));
const Scene = lazyComponent(() => import("./SceneDetails/Scene"));
const SceneCreate = lazyComponent(() => import("./SceneDetails/SceneCreate"));

const Scenes: React.FC = () => {
  useScrollToTopOnMount();

  return <SceneList persistState={PersistanceLevel.ALL} />;
};

const SceneMarkers: React.FC = () => {
  useScrollToTopOnMount();

  const titleProps = useTitleProps({ id: "markers" });
  return (
    <>
      <Helmet {...titleProps} />
      <SceneMarkerList />
    </>
  );
};

const SceneRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "scenes" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/scenes" component={Scenes} />
        <Route exact path="/scenes/markers" component={SceneMarkers} />
        <Route exact path="/scenes/new" component={SceneCreate} />
        <Route path="/scenes/:id" component={Scene} />
      </Switch>
    </>
  );
};

export default SceneRoutes;
