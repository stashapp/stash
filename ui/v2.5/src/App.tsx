import React from "react";
import { Route, Switch } from "react-router-dom";
import { ToastProvider } from "src/hooks/Toast";
import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { ErrorBoundary } from "./components/ErrorBoundary";
import Galleries from "./components/Galleries/Galleries";
import { MainNavbar } from "./components/MainNavbar";
import { PageNotFound } from "./components/PageNotFound";
import Performers from "./components/Performers/Performers";
import Scenes from "./components/Scenes/Scenes";
import { Settings } from "./components/Settings/Settings";
import { Stats } from "./components/Stats";
import Studios from "./components/Studios/Studios";
import { TagList } from "./components/Tags/TagList";
import { SceneFilenameParser } from "./components/SceneFilenameParser/SceneFilenameParser";

library.add(fas);

export const App: React.FC = () => (
  <div className="bp3-dark">
    <ErrorBoundary>
      <ToastProvider>
        <MainNavbar />
        <div className="main container-fluid">
          <Switch>
            <Route exact path="/" component={Stats} />
            <Route path="/scenes" component={Scenes} />
            {/* <Route path="/scenes/:id" component={Scene} /> */}
            <Route path="/galleries" component={Galleries} />
            <Route path="/performers" component={Performers} />
            <Route path="/tags" component={TagList} />
            <Route path="/studios" component={Studios} />
            <Route path="/settings" component={Settings} />
            <Route
              path="/sceneFilenameParser"
              component={SceneFilenameParser}
            />
            <Route component={PageNotFound} />
          </Switch>
        </div>
      </ToastProvider>
    </ErrorBoundary>
  </div>
);
