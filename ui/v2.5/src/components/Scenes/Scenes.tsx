import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared/constants";
import { PersistanceLevel } from "../List/ItemList";
import { lazyComponent } from "src/utils/lazyComponent";

const SceneList = lazyComponent(() => import("./SceneList"));
const SceneMarkerList = lazyComponent(() => import("./SceneMarkerList"));
const Scene = lazyComponent(() => import("./SceneDetails/Scene"));
const SceneCreate = lazyComponent(() => import("./SceneDetails/SceneCreate"));

const Scenes: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "scenes",
  })} ${TITLE_SUFFIX}`;
  const marker_title_template = `${intl.formatMessage({
    id: "markers",
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
        <Route
          exact
          path="/scenes/markers"
          render={() => (
            <>
              <Helmet
                defaultTitle={marker_title_template}
                titleTemplate={`%s | ${marker_title_template}`}
              />
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
