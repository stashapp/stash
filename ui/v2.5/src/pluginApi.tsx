import React from "react";
import ReactDOM from "react-dom";
import * as ReactRouterDOM from "react-router-dom";
import Mousetrap from "mousetrap";
import MousetrapPause from "mousetrap-pause";
import NavUtils from "./utils/navigation";
import * as GQL from "src/core/generated-graphql";
import * as StashService from "src/core/StashService";
import * as Apollo from "@apollo/client";
import * as Bootstrap from "react-bootstrap";
import * as Intl from "react-intl";
import * as FontAwesomeSolid from "@fortawesome/free-solid-svg-icons";
import * as FontAwesomeRegular from "@fortawesome/free-regular-svg-icons";
import * as FontAwesomeBrands from "@fortawesome/free-brands-svg-icons";
import * as ReactSelect from "react-select";
import { useSpriteInfo } from "./hooks/sprite";
import { useToast } from "./hooks/Toast";
import Event from "./hooks/event";
import { after, before, components, instead, RegisterComponent } from "./patch";
import { useSettings } from "./components/Settings/context";
import { useInteractive } from "./hooks/Interactive/context";
import InteractiveUtils from "./hooks/Interactive/utils";
import { useLightbox, useGalleryLightbox } from "./hooks/Lightbox/hooks";

// due to code splitting, some components may not have been loaded when a plugin
// page is loaded. This function will load all components passed to it.
// The components need to be imported here. Any required imports will be added
// to the loadableComponents object in the plugin api.
async function loadComponents(c: (() => Promise<unknown>)[]) {
  await Promise.all(c.map((fn) => fn()));
}

// useLoadComponents is a hook that loads all components passed to it.
// It returns a boolean indicating whether the components are still loading.
function useLoadComponents(toLoad: (() => Promise<unknown>)[]) {
  const [loading, setLoading] = React.useState(true);
  const [componentList] = React.useState(toLoad);

  async function load(c: (() => Promise<unknown>)[]) {
    await loadComponents(c);
    setLoading(false);
  }

  React.useEffect(() => {
    setLoading(true);
    load(componentList);
  }, [componentList]);

  return loading;
}

function registerRoute(path: string, component: React.FC) {
  before("PluginRoutes", function (props: React.PropsWithChildren<{}>) {
    return [
      {
        children: (
          <>
            {props.children}
            <ReactRouterDOM.Route path={path} component={component} />
          </>
        ),
      },
    ];
  });
}

export const PluginApi = {
  React,
  ReactDOM,
  GQL,
  libraries: {
    ReactRouterDOM,
    Bootstrap,
    Apollo,
    Intl,
    FontAwesomeRegular,
    FontAwesomeSolid,
    FontAwesomeBrands,
    Mousetrap,
    MousetrapPause,
    ReactSelect,
  },
  register: {
    // register a route to be added to the main router
    route: registerRoute,
    // register a component to be added to the components object
    component: RegisterComponent,
  },
  loadableComponents: {
    // add any lazy loaded imports here - this is coarse-grained and will load all components
    // in the import
    Performers: () => import("./components/Performers/Performers"),
    FrontPage: () => import("./components/FrontPage/FrontPage"),
    Scenes: () => import("./components/Scenes/Scenes"),
    Settings: () => import("./components/Settings/Settings"),
    Stats: () => import("./components/Stats"),
    Studios: () => import("./components/Studios/Studios"),
    Galleries: () => import("./components/Galleries/Galleries"),
    Groups: () => import("./components/Groups/Groups"),
    Tags: () => import("./components/Tags/Tags"),
    Images: () => import("./components/Images/Images"),

    SubmitStashBoxDraft: () => import("src/components/Dialogs/SubmitDraft"),
    GenerateDialog: () => import("./components/Dialogs/GenerateDialog"),

    ScenePlayer: () => import("src/components/ScenePlayer/ScenePlayer"),

    GalleryViewer: () => import("src/components/Galleries/GalleryViewer"),

    DeleteScenesDialog: () => import("./components/Scenes/DeleteScenesDialog"),
    SceneList: () => import("./components/Scenes/SceneList"),
    SceneMarkerList: () => import("./components/Scenes/SceneMarkerList"),
    Scene: () => import("./components/Scenes/SceneDetails/Scene"),
    SceneCreate: () => import("./components/Scenes/SceneDetails/SceneCreate"),

    ExternalPlayerButton: () =>
      import("./components/Scenes/SceneDetails/ExternalPlayerButton"),
    QueueViewer: () => import("./components/Scenes/SceneDetails/QueueViewer"),
    SceneMarkersPanel: () =>
      import("./components/Scenes/SceneDetails/SceneMarkersPanel"),
    SceneFileInfoPanel: () =>
      import("./components/Scenes/SceneDetails/SceneFileInfoPanel"),
    SceneDetailPanel: () =>
      import("./components/Scenes/SceneDetails/SceneDetailPanel"),
    SceneHistoryPanel: () =>
      import("./components/Scenes/SceneDetails/SceneHistoryPanel"),
    SceneGroupPanel: () =>
      import("./components/Scenes/SceneDetails/SceneGroupPanel"),
    SceneGalleriesPanel: () =>
      import("./components/Scenes/SceneDetails/SceneGalleriesPanel"),
    SceneVideoFilterPanel: () =>
      import("./components/Scenes/SceneDetails/SceneVideoFilterPanel"),
    SceneScrapeDialog: () =>
      import("./components/Scenes/SceneDetails/SceneScrapeDialog"),
    SceneQueryModal: () =>
      import("./components/Scenes/SceneDetails/SceneQueryModal"),

    LightboxComponent: () => import("src/hooks/Lightbox/Lightbox"),

    // intentionally omitting these for now
    // Setup: () => import("./components/Setup/Setup"),
    // Migrate: () => import("./components/Setup/Migrate"),
    // SceneFilenameParser: () => import("./components/SceneFilenameParser/SceneFilenameParser"),
    // SceneDuplicateChecker: () => import("./components/SceneDuplicateChecker/SceneDuplicateChecker"),
    // Manual: () => import("./Manual"),

    // individual components here
    // add components as needed for plugins that provide pages
    SceneCard: () => import("./components/Scenes/SceneCard"),
    PerformerSelect: () => import("./components/Performers/PerformerSelect"),
    TagLink: () => import("./components/Shared/TagLink"),
    PerformerCard: () => import("./components/Performers/PerformerCard"),
  },
  components,
  utils: {
    InteractiveUtils,
    NavUtils,
    StashService,
    loadComponents,
  },
  hooks: {
    useLoadComponents,
    useSpriteInfo,
    useToast,
    useSettings,
    useInteractive,
    useLightbox,
    useGalleryLightbox,
  },
  patch: {
    // intercept the arguments of supported functions
    // the provided function should accept the arguments and return the new arguments
    before,
    // replace a function with a new one implementation
    // the provided function will be called with the arguments passed to the original function
    // and the original function as the last argument
    // only one instead function can be registered per component
    instead,
    // intercept the result of supported functions
    // the provided function will be called with the arguments passed to the original function
    // and the result of the original function
    after,
  },
  Event: Event,
};

export default PluginApi;

interface IWindow {
  PluginApi: typeof PluginApi;
}

const localWindow = window as unknown as IWindow;

// export the plugin api to the window object
localWindow.PluginApi = PluginApi;
