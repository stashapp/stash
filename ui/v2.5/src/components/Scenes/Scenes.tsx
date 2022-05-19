import React, { lazy } from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared";
import { PersistanceLevel } from "src/hooks/ListHook";

const SceneList = lazy(() => import("./SceneList"));
const SceneMarkerList = lazy(() => import("./SceneMarkerList"));
const Scene = lazy(() => import("./SceneDetails/Scene"));

const Scenes: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "scenes",
  })} ${TITLE_SUFFIX}`;
  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />
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
    </>
  );
};
export default Scenes;
