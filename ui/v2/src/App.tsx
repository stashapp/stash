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

export class App extends React.Component {
  public render() {
    return (
      <div className="bp3-dark">
        <MainNavbar />
        <ErrorBoundary>
          <Switch>
            <Route exact={true} path="/" component={Stats} />
            <Route path="/scenes" component={Scenes} />
            {/* <Route path="/scenes/:id" component={Scene} /> */}
            <Route path="/galleries" component={Galleries} />
            <Route path="/performers" component={Performers} />
            <Route path="/tags" component={Tags} />
            <Route path="/studios" component={Studios} />
            <Route path="/settings" component={Settings} />
            <Route component={PageNotFound} />
          </Switch>
        </ErrorBoundary>
      </div>
    );
  }
}
