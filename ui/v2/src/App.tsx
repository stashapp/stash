import React, { FunctionComponent, useEffect, useState } from "react";
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
import { Sidebar } from "./components/Sidebar";

interface IProps {}

export const App: FunctionComponent<IProps> = (props: IProps) => {
  const [menuOpen, setMenuOpen] = useState<boolean>(getInitialMenuState());

  function getInitialMenuState() {
    return window.innerWidth > 768;
  }

  function getSidebarClosedClass() {
    if (!menuOpen) {
      return " sidebar-closed";
    }

    return "";
  }

  return (
    <div className="bp3-dark">
      <ErrorBoundary>
        <MainNavbar onMenuToggle={() => setMenuOpen(!menuOpen)}/>
        <Sidebar className={getSidebarClosedClass()}/>
        <div className={"main" + getSidebarClosedClass()}>
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
      </ErrorBoundary>
    </div>
  );
};
