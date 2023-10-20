import React from "react";
import { Route } from "react-router-dom";

interface IPluginComponentsMap {
  [key: string]: React.FC[];
}

interface IPluginPagesMap {
  [key: string]: React.FC;
}

export const pluginComponents: IPluginComponentsMap = {};
export let pluginPages: IPluginPagesMap = {};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function registerPluginPage(path: string, component: React.FC) {
  pluginPages = {
    ...pluginPages,
    [path]: component,
  };
}

export function addPluginComponent(
  location: PluginComponentLocation,
  component: React.FC
) {
  if (!pluginComponents[location]) {
    pluginComponents[location] = [];
  }

  pluginComponents[location].push(component);
}

export enum PluginComponentLocation {
  Main = "main",
  Navbar = "navbar",
}

export const PluginRoutes: React.FC = () => {
  const routes = Object.entries(pluginPages).map((e) => {
    const [path, component] = e;
    const prefixedPath = `/plugin/${path}`;
    return (
      <Route key={prefixedPath} path={prefixedPath} component={component} />
    );
  });

  return <>{routes}</>;
};

export function renderPluginComponents(location: PluginComponentLocation) {
  if (!pluginComponents[location]) return null;

  return pluginComponents[location].map((Component, index) => (
    <Component key={index} />
  ));
}

interface IPluginComponents {
  location: PluginComponentLocation;
}

export const PluginComponents: React.FC<IPluginComponents> = ({ location }) => {
  return <>{renderPluginComponents(location)}</>;
};

interface ICardHooks {
  Image?: React.FC;
  Overlays?: React.FC;
  Details?: React.FC;
  Popovers?: React.FC;
}

interface IPluginComponentHooks {
  SceneCard?: ICardHooks;
}

type PluginHookType = keyof IPluginComponentHooks;

export let pluginComponentHooks: IPluginComponentHooks = {};

export function registerCardComponentHooks(
  location: PluginHookType,
  component: ICardHooks
) {
  pluginComponentHooks = {
    ...pluginComponentHooks,
    [location]: component,
  };
}
