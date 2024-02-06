import React from "react";
import ReactDOM from "react-dom";
import * as ReactRouterDOM from "react-router-dom";
import Mousetrap from "mousetrap";
import MousetrapPause from "mousetrap-pause";
import NavUtils from "./utils/navigation";
import { HoverPopover } from "./components/Shared/HoverPopover";
import { TagLink } from "./components/Shared/TagLink";
import { LoadingIndicator } from "./components/Shared/LoadingIndicator";
import * as GQL from "src/core/generated-graphql";
import * as StashService from "src/core/StashService";
import * as Apollo from "@apollo/client";
import * as Bootstrap from "react-bootstrap";
import * as Intl from "react-intl";
import * as FontAwesomeSolid from "@fortawesome/free-solid-svg-icons";
import * as FontAwesomeRegular from "@fortawesome/free-regular-svg-icons";
import { useSpriteInfo } from "./hooks/sprite";

// due to code splitting, some components may not have been loaded when a plugin
// page is loaded. This function will load all components passed to it.
// The components need to be imported here. Any required imports will be added
// to the loadableComponents object in the plugin api.
async function loadComponents(c: (() => Promise<unknown>)[]) {
  await Promise.all(c.map((fn) => fn()));
}

// useLoadComponents is a hook that loads all components passed to it.
// It returns a boolean indicating whether the components are still loading.
function useLoadComponents(components: (() => Promise<unknown>)[]) {
  const [loading, setLoading] = React.useState(true);
  const [componentList] = React.useState(components);

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

const components: Record<string, Function> = {
  HoverPopover,
  TagLink,
  LoadingIndicator,
};

const beforeFns: Record<string, Function[]> = {};
const insteadFns: Record<string, Function> = {};
const afterFns: Record<string, Function[]> = {};

// patch functions
// registers a patch to a function. Before functions are expected to return the
// new arguments to be passed to the function.
function before(component: string, fn: Function) {
  if (!beforeFns[component]) {
    beforeFns[component] = [];
  }
  beforeFns[component].push(fn);
}

function instead(component: string, fn: Function) {
  if (insteadFns[component]) {
    throw new Error("instead has already been called for " + component);
  }
  insteadFns[component] = fn;
}

function after(component: string, fn: Function) {
  if (!afterFns[component]) {
    afterFns[component] = [];
  }
  afterFns[component].push(fn);
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

export function RegisterComponent(component: string, fn: Function) {
  // register with the plugin api
  components[component] = fn;

  return fn;
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
};

// patches a function to implement the before/instead/after functionality
export function PatchFunction(name: string, fn: Function) {
  return new Proxy(fn, {
    apply(target, ctx, args) {
      let result;

      for (const beforeFn of beforeFns[name] || []) {
        args = beforeFn.apply(ctx, args);
      }
      if (insteadFns[name]) {
        result = insteadFns[name].apply(ctx, args.concat(target));
      } else {
        result = target.apply(ctx, args);
      }
      for (const afterFn of afterFns[name] || []) {
        result = afterFn.apply(ctx, args.concat(result));
      }
      return result;
    },
  });
}

// patches a component and registers it in the pluginapi components object
export function PatchComponent<T>(
  component: string,
  fn: React.FC<T>
): React.FC<T> {
  const ret = PatchFunction(component, fn);

  // register with the plugin api
  RegisterComponent(component, ret);
  return ret as React.FC<T>;
}

// patches a component and registers it in the pluginapi components object
export function PatchContainerComponent(
  component: string
): React.FC<React.PropsWithChildren<{}>> {
  const fn = (props: React.PropsWithChildren<{}>) => {
    return <>{props.children}</>;
  };

  return PatchComponent(component, fn);
}

export default PluginApi;

interface IWindow {
  PluginApi: typeof PluginApi;
}

const localWindow = window as unknown as IWindow;

// export the plugin api to the window object
localWindow.PluginApi = PluginApi;
