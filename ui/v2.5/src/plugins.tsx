import React from "react";

export enum PluginComponentLocation {
  Main = "main",
  Navbar = "navbar",
}

interface IPluginComponentsMap {
  [key: string]: React.FC[];
}

export const pluginComponentsMap: IPluginComponentsMap = {};

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
