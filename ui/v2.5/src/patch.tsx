import React from "react";

export let components: Record<string, Function> = {};

const beforeFns: Record<string, Function[]> = {};
const insteadFns: Record<string, Function[]> = {};
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

// registers a patch to a function. Instead functions receive the original arguments,
// plus the next function to call. In order for all instead functions to be called,
// it is expected that the provided next() function will be called.
export function instead(component: string, fn: Function) {
  if (!insteadFns[component]) {
    insteadFns[component] = [];
  }
  insteadFns[component].push(fn);
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
    // only throw an error in production, in development we allow
    // multiple registrations to allow for hot reloading of components
    if (!import.meta.env.DEV) {
      throw new Error(
        "Component " + component + " has already been registered"
      );
    }
  }

  components[component] = fn;

  return fn;
}

/* eslint-disable @typescript-eslint/no-explicit-any */
function runInstead(
  fns: Function[],
  targetFn: Function,
  thisArg: any,
  argArray: any[]
) {
  if (!fns.length) {
    return targetFn.apply(thisArg, argArray);
  }

  let i = 1;
  function next(): any {
    if (i >= fns.length) {
      return targetFn;
    }

    const thisTarget = fns[i++];
    return new Proxy(thisTarget, {
      apply: function (target, ctx, args) {
        return target.apply(ctx, args.concat(next()));
      },
    });
  }

  return fns[0].apply(thisArg, argArray.concat(next()));
}
/* eslint-enable @typescript-eslint/no-explicit-any */

// patches a function to implement the before/instead/after functionality
export function PatchFunction<T extends Function>(name: string, fn: T) {
  return new Proxy(fn, {
    apply(target, ctx, args) {
      let result;

      for (const beforeFn of beforeFns[name] || []) {
        args = beforeFn.apply(ctx, args);
      }
      if (insteadFns[name]) {
        result = runInstead(insteadFns[name], target, ctx, args);
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
export function PatchContainerComponent<T = {}>(
  component: string
): React.FC<React.PropsWithChildren<T>> {
  const fn: React.FC<React.PropsWithChildren<T>> = (
    props: React.PropsWithChildren<T>
  ) => {
    return <>{props.children}</>;
  };

  return PatchComponent(component, fn);
}
