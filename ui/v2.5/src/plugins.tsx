import React from "react";
import { PatchFunction } from "./patch";

export const PluginRoutes: React.FC<React.PropsWithChildren<{}>> =
  PatchFunction("PluginRoutes", (props: React.PropsWithChildren<{}>) => {
    return <>{props.children}</>;
  }) as React.FC;
