import React from "react";
import { Route } from "react-router-dom";

export enum PluginComponentLocation {
  Main = "main",
  Navbar = "navbar",
}

interface IPluginComponentsMap {
  [key: string]: React.FC[];
}

interface IPluginPagesMap {
  [key: string]: React.FC;
}

export const pluginComponentsMap: IPluginComponentsMap = {};
let pluginPagesMap: IPluginPagesMap = {};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
(window as any).addPluginComponent = (
  location: PluginComponentLocation,
  component: React.FC
) => {
  if (!pluginComponentsMap[location]) {
    pluginComponentsMap[location] = [];
  }

  pluginComponentsMap[location].push(component);
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
(window as any).registerPluginPage = (
  path: string,
  component: React.FC
) => {
  pluginPagesMap = {
    ...pluginPagesMap,
    [path]: component,
  }
};

export const PluginRoutes: React.FC = () => {
  const routes = Object.entries(pluginPagesMap).map((e) => {
    const [path, component] = e;
    const prefixedPath = `/plugin/${path}`;
    return (
      <Route key={prefixedPath} path={prefixedPath} component={component} />
    );
  });

  return <>{routes}</>;
}

export function renderPluginComponents(location: PluginComponentLocation) {
  if (!pluginComponentsMap[location]) return null;

  return pluginComponentsMap[location].map((Component, index) => (
    <Component key={index} />
  ));
}

interface IPluginComponents {
  location: PluginComponentLocation;
}

export const PluginComponents: React.FC<IPluginComponents> = ({ location }) => {
  return <>{renderPluginComponents(location)}</>;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
(window as any).React = React;
