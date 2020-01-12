import React, { FunctionComponent, useState } from "react";
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
import Dvds from "./components/Dvds/Dvds";
import Tags from "./components/Tags/Tags";
import { SceneFilenameParser } from "./components/scenes/SceneFilenameParser";
import { Sidebar } from "./components/Sidebar";
import { IconName } from "@blueprintjs/core";

export interface IMenuItem {
  icon: IconName
  text: string
  href: string
}

interface IProps {}

export const App: FunctionComponent<IProps> = (props: IProps) => {
  const [menuOpen, setMenuOpen] = useState<boolean>(false);

  function getSidebarClosedClass() {
    if (!menuOpen) {
      return " sidebar-closed";
    }

    return "";
  }

  const menuItems: IMenuItem[] = [
    {
      icon: "video",
      text: "Scenes",
      href: "/scenes"
    },
    {
      href: "/scenes/markers",
      icon: "map-marker",
      text: "Markers"
    },
    {
      href: "/galleries",
      icon: "media",
      text: "Galleries"
    },
    {
      href: "/performers",
      icon: "person",
      text: "Performers"
    },
    {
      href: "/studios",
      icon: "mobile-video",
      text: "Studios"
    },
    {
      href: "/dvds",
      icon: "film",
      text: "Dvds"
    },
    {
      href: "/tags",
      icon: "tag",
      text: "Tags"
    }
  ];

  return (
    <div className="bp3-dark">
      <ErrorBoundary>
        <MainNavbar onMenuToggle={() => setMenuOpen(!menuOpen)} menuItems={menuItems}/>
        <Sidebar className={getSidebarClosedClass()} menuItems={menuItems}/>
        <div className={"main" + getSidebarClosedClass()}>
          <Switch>
            <Route exact={true} path="/" component={Stats} />
            <Route path="/scenes" component={Scenes} />
            {/* <Route path="/scenes/:id" component={Scene} /> */}
            <Route path="/galleries" component={Galleries} />
            <Route path="/performers" component={Performers} />
            <Route path="/tags" component={Tags} />
            <Route path="/studios" component={Studios} />
            <Route path="/dvds" component={Dvds} />
            <Route path="/settings" component={Settings} />
            <Route path="/sceneFilenameParser" component={SceneFilenameParser} />
            <Route component={PageNotFound} />
          </Switch>
        </div>
      </ErrorBoundary>
    </div>
  );
};
