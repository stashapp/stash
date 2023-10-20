import React from "react";
import { ScenePreview } from "./components/Scenes/SceneCard";
import {
  addPluginComponent,
  registerCardComponentHooks,
  registerPluginPage,
} from "./plugins";

export const PluginApi = {
  React,
  register: {
    page: registerPluginPage,
    component: addPluginComponent,
    cardComponentHook: registerCardComponentHooks,
  },
  components: {
    ScenePreview,
  },
};

export default PluginApi;

interface IWindow {
  PluginApi: typeof PluginApi;
}

const localWindow = window as unknown as IWindow;

// export the plugin api to the window object
localWindow.PluginApi = PluginApi;
