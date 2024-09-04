import React from "react";
import ReactDOM from "react-dom";
import * as ReactRouterDOM from "react-router-dom";
import Mousetrap from "mousetrap";
import MousetrapPause from "mousetrap-pause";
import NavUtils from "./utils/navigation";
import * as GQL from "src/core/generated-graphql";
import * as StashService from "src/core/StashService";
import * as Apollo from "@apollo/client";
import * as Bootstrap from "react-bootstrap";
import * as Intl from "react-intl";
import * as FontAwesomeSolid from "@fortawesome/free-solid-svg-icons";
import * as FontAwesomeRegular from "@fortawesome/free-regular-svg-icons";
import { useSpriteInfo } from "./hooks/sprite";
import { useToast } from "./hooks/Toast";
import Event from "./hooks/event";
import { before, instead, after, components, RegisterComponent } from "./patch";
import { useSettings } from "./components/Settings/context";
import { useInteractive } from "./hooks/Interactive/context";

// due to code splitting, some components may not have been loaded when a plugin
// page is loaded. This function will load all components passed to it.
// The components need to be imported here. Any required imports will be added
// to the loadableComponents object in the plugin api.
async function loadComponents(c: (() => Promise<unknown>)[]) {
  await Promise.all(c.map((fn) => fn()));
}

// useLoadComponents is a hook that loads all components passed to it.
// It returns a boolean indicating whether the components are still loading.
function useLoadComponents(toLoad: (() => Promise<unknown>)[]) {
  const [loading, setLoading] = React.useState(true);
  const [componentList] = React.useState(toLoad);

  async function load(c: (() => Promise<unknown>)[]) {
    await loadComponents(c);
    setLoading(false);
  }

  React.useEffect(() => {
    setLoading(true);
    load(componentList);
  }, [componentList]);

  return loading;
}

function registerRoute(path: string, component: React.FC) {
  before("PluginRoutes", function (props: React.PropsWithChildren<{}>) {
    return [
      {
        children: (
          <>
            {props.children}
            <ReactRouterDOM.Route path={path} component={component} />
          </>
        ),
      },
    ];
  });
}

export const PluginApi = {
  React,
  ReactDOM,
  GQL,
  libraries: {
    ReactRouterDOM,
    Bootstrap,
    Apollo,
    Intl,
    FontAwesomeRegular,
    FontAwesomeSolid,
    Mousetrap,
    MousetrapPause,
  },
  register: {
    // register a route to be added to the main router
    route: registerRoute,
    // register a component to be added to the components object
    component: RegisterComponent,
  },
  loadableComponents: {
    // add components as needed for plugins that provide pages
    SceneCard: () => import("./components/Scenes/SceneCard"),
  },
  components,
  utils: {
    NavUtils,
    StashService,
    loadComponents,
  },
  hooks: {
    useLoadComponents,
    useSpriteInfo,
    useToast,
    useSettings,
    useInteractive,
  },
  patch: {
    // intercept the arguments of supported functions
    // the provided function should accept the arguments and return the new arguments
    before,
    // replace a function with a new one implementation
    // the provided function will be called with the arguments passed to the original function
    // and the original function as the last argument
    // only one instead function can be registered per component
    instead,
    // intercept the result of supported functions
    // the provided function will be called with the arguments passed to the original function
    // and the result of the original function
    after,
  },
  Event: Event,
};

export default PluginApi;

interface IWindow {
  PluginApi: typeof PluginApi;
}

const localWindow = window as unknown as IWindow;

// export the plugin api to the window object
localWindow.PluginApi = PluginApi;
