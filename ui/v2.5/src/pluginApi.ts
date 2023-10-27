import React from "react";
import { addPluginComponent, registerPluginPage } from "./plugins";
import { Link } from "react-router-dom";
import NavUtils from "./utils/navigation";
import { HoverPopover } from "./components/Shared/HoverPopover";
import { TagLink } from "./components/Shared/TagLink";

const components: Record<string, Function> = {
  HoverPopover,
  TagLink,
};

const beforeFns: Record<string, Function[]> = {};
const insteadFns: Record<string, Function> = {};
const afterFns: Record<string, Function[]> = {};

// patch functions
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

export const PluginApi = {
  React,
  ReactReactorDOM: {
    Link,
  },
  register: {
    page: registerPluginPage,
    component: addPluginComponent,
  },
  components,
  utils: {
    NavUtils: NavUtils,
  },
  patch: {
    before,
    instead,
    after,
  },
};

export function PatchFunction(name: string, fn: Function) {
  return new Proxy(fn, {
    apply(target, ctx, args) {
      let result;

      for (const beforeFn of beforeFns[name] || []) {
        beforeFn.apply(ctx, args);
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

export function PatchComponent(component: string, fn: Function) {
  // register with the plugin api
  PluginApi.components[component] = fn;

  return PatchFunction(component, fn);
}

export function RegisterComponent(component: string, fn: Function) {
  // register with the plugin api
  PluginApi.components[component] = fn;

  return fn;
}

export default PluginApi;

interface IWindow {
  PluginApi: typeof PluginApi;
}

const localWindow = window as unknown as IWindow;

// export the plugin api to the window object
localWindow.PluginApi = PluginApi;
