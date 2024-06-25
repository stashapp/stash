import React, { Suspense, useEffect, useState } from "react";
import {
  Route,
  Switch,
  useHistory,
  useLocation,
  useRouteMatch,
} from "react-router-dom";
import { IntlProvider, CustomFormats } from "react-intl";
import { Helmet } from "react-helmet";
import cloneDeep from "lodash-es/cloneDeep";
import mergeWith from "lodash-es/mergeWith";
import { ToastProvider } from "src/hooks/Toast";
import { LightboxProvider } from "src/hooks/Lightbox/context";
import { initPolyfills } from "src/polyfills";

import locales, { registerCountry } from "src/locales";
import {
  useConfiguration,
  useConfigureUI,
  usePlugins,
  useSystemStatus,
} from "src/core/StashService";
import flattenMessages from "./utils/flattenMessages";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import MousetrapPause from "mousetrap-pause";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { MainNavbar } from "./components/MainNavbar";
import { PageNotFound } from "./components/PageNotFound";
import * as GQL from "./core/generated-graphql";
import { makeTitleProps } from "./hooks/title";
import { LoadingIndicator } from "./components/Shared/LoadingIndicator";

import { ConfigurationProvider } from "./hooks/Config";
import { ManualProvider } from "./components/Help/context";
import { InteractiveProvider } from "./hooks/Interactive/context";
import { ReleaseNotesDialog } from "./components/Dialogs/ReleaseNotesDialog";
import { releaseNotes } from "./docs/en/ReleaseNotes";
import { getPlatformURL } from "./core/createClient";
import { lazyComponent } from "./utils/lazyComponent";
import { isPlatformUniquelyRenderedByApple } from "./utils/apple";
import useScript, { useCSS } from "./hooks/useScript";
import { useMemoOnce } from "./hooks/state";
import Event from "./hooks/event";
import { uniq } from "lodash-es";

import { PluginRoutes } from "./plugins";

// import plugin_api to run code
import "./pluginApi";
import { ConnectionMonitor } from "./ConnectionMonitor";
import { PatchFunction } from "./patch";

const Performers = lazyComponent(
  () => import("./components/Performers/Performers")
);
const FrontPage = lazyComponent(
  () => import("./components/FrontPage/FrontPage")
);
const Scenes = lazyComponent(() => import("./components/Scenes/Scenes"));
const Settings = lazyComponent(() => import("./components/Settings/Settings"));
const Stats = lazyComponent(() => import("./components/Stats"));
const Studios = lazyComponent(() => import("./components/Studios/Studios"));
const Galleries = lazyComponent(
  () => import("./components/Galleries/Galleries")
);

const Groups = lazyComponent(() => import("./components/Movies/Movies"));
const Tags = lazyComponent(() => import("./components/Tags/Tags"));
const Images = lazyComponent(() => import("./components/Images/Images"));
const Setup = lazyComponent(() => import("./components/Setup/Setup"));
const Migrate = lazyComponent(() => import("./components/Setup/Migrate"));

const SceneFilenameParser = lazyComponent(
  () => import("./components/SceneFilenameParser/SceneFilenameParser")
);
const SceneDuplicateChecker = lazyComponent(
  () => import("./components/SceneDuplicateChecker/SceneDuplicateChecker")
);

const appleRendering = isPlatformUniquelyRenderedByApple();

initPolyfills();

MousetrapPause(Mousetrap);

const intlFormats: CustomFormats = {
  date: {
    long: { year: "numeric", month: "long", day: "numeric" },
  },
};

const defaultLocale = "en-GB";

function languageMessageString(language: string) {
  return language.replace(/-/, "");
}

type PluginList = NonNullable<Required<GQL.PluginsQuery["plugins"]>>;

// sort plugins by their dependencies
function sortPlugins(plugins: PluginList) {
  type Node = { id: string; afters: string[] };

  let nodes: Record<string, Node> = {};
  let sorted: PluginList = [];
  let visited: Record<string, boolean> = {};

  plugins.forEach((v) => {
    let from = v.id;

    if (!nodes[from]) nodes[from] = { id: from, afters: [] };

    v.requires?.forEach((to) => {
      if (!nodes[to]) nodes[to] = { id: to, afters: [] };
      if (!nodes[to].afters.includes(from)) nodes[to].afters.push(from);
    });
  });

  function visit(idstr: string, ancestors: string[] = []) {
    let node = nodes[idstr];
    const { id } = node;

    if (visited[idstr]) return;

    ancestors.push(id);
    visited[idstr] = true;
    node.afters.forEach(function (afterID) {
      if (ancestors.indexOf(afterID) >= 0)
        throw new Error("closed chain : " + afterID + " is in " + id);
      visit(afterID.toString(), ancestors.slice());
    });

    const plugin = plugins.find((v) => v.id === id);
    if (plugin) {
      sorted.unshift(plugin);
    }
  }

  Object.keys(nodes).forEach((n) => {
    visit(n);
  });

  return sorted;
}

const AppContainer: React.FC<React.PropsWithChildren<{}>> = PatchFunction(
  "App",
  (props: React.PropsWithChildren<{}>) => {
    return <>{props.children}</>;
  }
) as React.FC;

export const App: React.FC = () => {
  const config = useConfiguration();
  const [saveUI] = useConfigureUI();

  const { data: systemStatusData } = useSystemStatus();

  const language =
    config.data?.configuration?.interface?.language ?? defaultLocale;

  // use en-GB as default messages if any messages aren't found in the chosen language
  const [messages, setMessages] = useState<{}>();
  const [customMessages, setCustomMessages] = useState<{}>();

  useEffect(() => {
    (async () => {
      try {
        const res = await fetch(getPlatformURL("customlocales"));
        if (res.ok) {
          setCustomMessages(await res.json());
        }
      } catch (err) {
        console.log(err);
      }
    })();
  }, []);

  useEffect(() => {
    const setLocale = async () => {
      const defaultMessageLanguage = languageMessageString(defaultLocale);
      const messageLanguage = languageMessageString(language);

      // register countries for the chosen language
      await registerCountry(language);

      const defaultMessages = (await locales[defaultMessageLanguage]()).default;
      const mergedMessages = cloneDeep(Object.assign({}, defaultMessages));
      const chosenMessages = (await locales[messageLanguage]()).default;

      mergeWith(
        mergedMessages,
        chosenMessages,
        customMessages,
        (objVal, srcVal) => {
          if (srcVal === "") {
            return objVal;
          }
        }
      );

      const newMessages = flattenMessages(mergedMessages);

      yup.setLocale({
        mixed: {
          required: newMessages["validation.required"],
        },
      });

      setMessages(newMessages);
    };

    setLocale();
  }, [customMessages, language]);

  const {
    data: plugins,
    loading: pluginsLoading,
    error: pluginsError,
  } = usePlugins();

  const sortedPlugins = useMemoOnce(() => {
    return [
      sortPlugins(plugins?.plugins ?? []),
      !pluginsLoading && !pluginsError,
    ];
  }, [plugins?.plugins, pluginsLoading, pluginsError]);

  const pluginJavascripts = useMemoOnce(() => {
    return [
      uniq(
        sortedPlugins
          ?.filter((plugin) => plugin.enabled && plugin.paths.javascript)
          .map((plugin) => plugin.paths.javascript!)
          .flat() ?? []
      ),
      !!sortedPlugins && !pluginsLoading && !pluginsError,
    ];
  }, [sortedPlugins, pluginsLoading, pluginsError]);

  const pluginCSS = useMemoOnce(() => {
    return [
      uniq(
        sortedPlugins
          ?.filter((plugin) => plugin.enabled && plugin.paths.css)
          .map((plugin) => plugin.paths.css!)
          .flat() ?? []
      ),
      !!sortedPlugins && !pluginsLoading && !pluginsError,
    ];
  }, [sortedPlugins, pluginsLoading, pluginsError]);

  useScript(pluginJavascripts ?? [], !pluginsLoading && !pluginsError);
  useCSS(pluginCSS ?? [], !pluginsLoading && !pluginsError);

  const location = useLocation();
  const history = useHistory();
  const setupMatch = useRouteMatch(["/setup", "/migrate"]);

  // dispatch event when location changes
  useEffect(() => {
    Event.dispatch("location", "", { location });
  }, [location]);

  // redirect to setup or migrate as needed
  useEffect(() => {
    if (!systemStatusData) {
      return;
    }

    const { status } = systemStatusData.systemStatus;

    if (
      location.pathname !== "/setup" &&
      status === GQL.SystemStatusEnum.Setup
    ) {
      // redirect to setup page
      history.push("/setup");
    }

    if (
      location.pathname !== "/migrate" &&
      status === GQL.SystemStatusEnum.NeedsMigration
    ) {
      // redirect to migrate page
      history.replace("/migrate");
    }
  }, [systemStatusData, setupMatch, history, location]);

  function maybeRenderNavbar() {
    // don't render navbar for setup views
    if (!setupMatch) {
      return <MainNavbar />;
    }
  }

  function renderContent() {
    if (!systemStatusData) {
      return <LoadingIndicator />;
    }

    return (
      <ErrorBoundary>
        <Suspense fallback={<LoadingIndicator />}>
          <Switch>
            <Route exact path="/" component={FrontPage} />
            <Route path="/scenes" component={Scenes} />
            <Route path="/images" component={Images} />
            <Route path="/galleries" component={Galleries} />
            <Route path="/performers" component={Performers} />
            <Route path="/tags" component={Tags} />
            <Route path="/studios" component={Studios} />
            <Route path="/groups" component={Groups} />
            <Route path="/stats" component={Stats} />
            <Route path="/settings" component={Settings} />
            <Route
              path="/sceneFilenameParser"
              component={SceneFilenameParser}
            />
            <Route
              path="/sceneDuplicateChecker"
              component={SceneDuplicateChecker}
            />
            <Route path="/setup" component={Setup} />
            <Route path="/migrate" component={Migrate} />
            <PluginRoutes />
            <Route component={PageNotFound} />
          </Switch>
        </Suspense>
      </ErrorBoundary>
    );
  }

  function maybeRenderReleaseNotes() {
    if (setupMatch || !systemStatusData || config.loading || config.error) {
      return;
    }

    const lastNoteSeen = config.data?.configuration.ui.lastNoteSeen;
    const notes = releaseNotes.filter((n) => {
      return !lastNoteSeen || n.date > lastNoteSeen;
    });

    if (notes.length === 0) return;

    return (
      <ReleaseNotesDialog
        notes={notes}
        onClose={() => {
          saveUI({
            variables: {
              input: {
                ...config.data?.configuration.ui,
                lastNoteSeen: notes[0].date,
              },
            },
          });
        }}
      />
    );
  }

  const titleProps = makeTitleProps();

  return (
    <AppContainer>
      <ErrorBoundary>
        {messages ? (
          <IntlProvider
            locale={language}
            messages={messages}
            formats={intlFormats}
          >
            <ConfigurationProvider
              configuration={config.data?.configuration}
              loading={config.loading}
            >
              {maybeRenderReleaseNotes()}
              <ToastProvider>
                <ConnectionMonitor />
                <Suspense fallback={<LoadingIndicator />}>
                  <LightboxProvider>
                    <ManualProvider>
                      <InteractiveProvider>
                        <Helmet {...titleProps} />
                        {maybeRenderNavbar()}
                        <div
                          className={`main container-fluid ${
                            appleRendering ? "apple" : ""
                          }`}
                        >
                          {renderContent()}
                        </div>
                      </InteractiveProvider>
                    </ManualProvider>
                  </LightboxProvider>
                </Suspense>
              </ToastProvider>
            </ConfigurationProvider>
          </IntlProvider>
        ) : null}
      </ErrorBoundary>
    </AppContainer>
  );
};
