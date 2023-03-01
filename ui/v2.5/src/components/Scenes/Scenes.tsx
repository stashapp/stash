import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared/constants";
import { PersistanceLevel } from "src/hooks/ListHook";
import { lazy_component } from "src/utils/lazy_component";

const SceneList = lazy_component(() => import("./SceneList"));
const SceneMarkerList = lazy_component(() => import("./SceneMarkerList"));
const Scene = lazy_component(() => import("./SceneDetails/Scene"));
const SceneCreate = lazy_component(() => import("./SceneDetails/SceneCreate"));

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
        <Route exact path="/scenes/new" component={SceneCreate} />
        <Route path="/scenes/:id" component={Scene} />
      </Switch>
    </>
  );
};
export default Scenes;
