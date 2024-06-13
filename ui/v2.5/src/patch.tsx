import React from "react";
import { HoverPopover } from "./components/Shared/HoverPopover";
import { TagLink } from "./components/Shared/TagLink";
import { LoadingIndicator } from "./components/Shared/LoadingIndicator";

export const components: Record<string, Function> = {
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
export function before(component: string, fn: Function) {
  if (!beforeFns[component]) {
    beforeFns[component] = [];
  }
  beforeFns[component].push(fn);
}

export function instead(component: string, fn: Function) {
  if (insteadFns[component]) {
    throw new Error("instead has already been called for " + component);
  }
  insteadFns[component] = fn;
}

export function after(component: string, fn: Function) {
  if (!afterFns[component]) {
    afterFns[component] = [];
  }
  afterFns[component].push(fn);
}

export function RegisterComponent<T extends Function>(
  component: string,
  fn: T
) {
  // register with the plugin api
  if (components[component]) {
    throw new Error("Component " + component + " has already been registered");
  }

  components[component] = fn;

  return fn;
}

// patches a function to implement the before/instead/after functionality
export function PatchFunction<T extends Function>(name: string, fn: T) {
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
