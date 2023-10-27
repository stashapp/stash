import React from "react";
import { PatchFunction } from "./pluginApi";

export const PluginRoutes: React.FC<React.PropsWithChildren<{}>> =
  PatchFunction("PluginRoutes", (props: React.PropsWithChildren<{}>) => {
    return <>{props.children}</>;
  }) as React.FC;
