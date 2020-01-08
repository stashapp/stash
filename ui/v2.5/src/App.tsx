import React from "react";
import { Route, Switch } from "react-router-dom";
import { ErrorBoundary } from "./components/ErrorBoundary";
import Galleries from "./components/Galleries/Galleries";
import { MainNavbar } from "./components/MainNavbar";
import { PageNotFound } from "./components/PageNotFound";
import Performers from "./components/performers/performers";
import Scenes from "./components/scenes/scenes";
import { Settings } from "./components/Settings/Settings";
import { Stats } from "./components/Stats";
import Studios from "./components/Studios/Studios";
import Tags from "./components/Tags/Tags";
import { SceneFilenameParser } from "./components/scenes/SceneFilenameParser";
import { ToastProvider } from './components/Shared/Toast';

import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'

import 'bootstrap/dist/css/bootstrap.min.css';

library.add(fas);

export const App: React.FC = () => (
  <div className="bp3-dark">
    <ErrorBoundary>
      <MainNavbar />
      <ToastProvider>
        <div className="main">
          <Switch>
            <Route exact={true} path="/" component={Stats} />
            <Route path="/scenes" component={Scenes} />
            {/* <Route path="/scenes/:id" component={Scene} /> */}
            <Route path="/galleries" component={Galleries} />
            <Route path="/performers" component={Performers} />
            <Route path="/tags" component={Tags} />
            <Route path="/studios" component={Studios} />
            <Route path="/settings" component={Settings} />
            <Route path="/sceneFilenameParser" component={SceneFilenameParser} />
            <Route component={PageNotFound} />
          </Switch>
        </div>
      </ToastProvider>
    </ErrorBoundary>
  </div>
);
