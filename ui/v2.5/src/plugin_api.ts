import React from "react";
import {
  addPluginComponent,
  registerCardComponentHooks,
  registerPluginPage,
} from "./plugins";
import { Link } from "react-router-dom";
import NavUtils from "./utils/navigation";
import { HoverPopover } from "./components/Shared/HoverPopover";
import { TagLink } from "./components/Shared/TagLink";

export const PluginApi = {
  React,
  ReactReactorDOM: {
    Link,
  },
  register: {
    page: registerPluginPage,
    component: addPluginComponent,
    cardComponentHook: registerCardComponentHooks,
  },
  components: {
    HoverPopover,
    TagLink,
  },
  utils: {
    NavUtils: NavUtils,
  }
};

export default PluginApi;

interface IWindow {
  PluginApi: typeof PluginApi;
}

const localWindow = window as unknown as IWindow;

// export the plugin api to the window object
localWindow.PluginApi = PluginApi;
