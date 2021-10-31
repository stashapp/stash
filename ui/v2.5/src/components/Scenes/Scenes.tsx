import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared";
import { PersistanceLevel } from "src/hooks/ListHook";
import Scene from "./SceneDetails/Scene";
import { SceneList } from "./SceneList";
import { SceneMarkerList } from "./SceneMarkerList";

const Scenes: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage(
    {
      id: "countables.scenes",
    },
    { count: 100 }
  )} ${TITLE_SUFFIX}`;
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
