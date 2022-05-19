import React, { lazy, Suspense, useEffect } from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import { IntlProvider, CustomFormats } from "react-intl";
import { Helmet } from "react-helmet";
import { mergeWith } from "lodash";
import { ToastProvider } from "src/hooks/Toast";
import LightboxProvider from "src/hooks/Lightbox/context";
import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { initPolyfills } from "src/polyfills";

import locales from "src/locales";
import { useConfiguration, useSystemStatus } from "src/core/StashService";
import { flattenMessages } from "src/utils";
import Mousetrap from "mousetrap";
import MousetrapPause from "mousetrap-pause";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { MainNavbar } from "./components/MainNavbar";
import { PageNotFound } from "./components/PageNotFound";
import * as GQL from "./core/generated-graphql";
import { LoadingIndicator, TITLE_SUFFIX } from "./components/Shared";
import { ConfigurationProvider } from "./hooks/Config";
import { ManualProvider } from "./components/Help/context";
import { InteractiveProvider } from "./hooks/Interactive/context";

const Performers = lazy(() => import("./components/Performers/Performers"));
const Recommendations = lazy(
  () => import("./components/Recommendations/Recommendations")
);
const Scenes = lazy(() => import("./components/Scenes/Scenes"));
const Settings = lazy(() => import("./components/Settings/Settings"));
const Stats = lazy(() => import("./components/Stats"));
const Studios = lazy(() => import("./components/Studios/Studios"));
const Galleries = lazy(() => import("./components/Galleries/Galleries"));

const Movies = lazy(() => import("./components/Movies/Movies"));
const Tags = lazy(() => import("./components/Tags/Tags"));
const Images = lazy(() => import("./components/Images/Images"));
const Setup = lazy(() => import("./components/Setup/Setup"));
const Migrate = lazy(() => import("./components/Setup/Migrate"));

const SceneFilenameParser = lazy(
  () => import("./components/SceneFilenameParser/SceneFilenameParser")
);
const SceneDuplicateChecker = lazy(
  () => import("./components/SceneDuplicateChecker/SceneDuplicateChecker")
);

initPolyfills();

MousetrapPause(Mousetrap);

// Set fontawesome/free-solid-svg as default fontawesome icons
library.add(fas);

const intlFormats: CustomFormats = {
  date: {
    long: { year: "numeric", month: "long", day: "numeric" },
  },
};

function languageMessageString(language: string) {
  return language.replace(/-/, "");
}

export const App: React.FC = () => {
  const config = useConfiguration();
  const { data: systemStatusData } = useSystemStatus();
  const defaultLocale = "en-GB";
  const language =
    config.data?.configuration?.interface?.language ?? defaultLocale;
  const defaultMessageLanguage = languageMessageString(defaultLocale);
  const messageLanguage = languageMessageString(language);

  // use en-GB as default messages if any messages aren't found in the chosen language
  const mergedMessages = mergeWith(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (locales as any)[defaultMessageLanguage],
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (locales as any)[messageLanguage],
    (objVal, srcVal) => {
      if (srcVal === "") {
        return objVal;
      }
    }
  );
  const messages = flattenMessages(mergedMessages);

  const setupMatch = useRouteMatch(["/setup", "/migrate"]);

  // redirect to setup or migrate as needed
  useEffect(() => {
    if (!systemStatusData) {
      return;
    }

    if (
      window.location.pathname !== "/setup" &&
      systemStatusData.systemStatus.status === GQL.SystemStatusEnum.Setup
    ) {
      // redirect to setup page
      const newURL = new URL("/setup", window.location.toString());
      window.location.href = newURL.toString();
    }

    if (
      window.location.pathname !== "/migrate" &&
      systemStatusData.systemStatus.status ===
        GQL.SystemStatusEnum.NeedsMigration
    ) {
      // redirect to setup page
      const newURL = new URL("/migrate", window.location.toString());
      window.location.href = newURL.toString();
    }
  }, [systemStatusData]);

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
      <Suspense fallback={<LoadingIndicator />}>
        <Switch>
          <Route exact path="/" component={Recommendations} />
          <Route path="/scenes" component={Scenes} />
          <Route path="/images" component={Images} />
          <Route path="/galleries" component={Galleries} />
          <Route path="/performers" component={Performers} />
          <Route path="/tags" component={Tags} />
          <Route path="/studios" component={Studios} />
          <Route path="/movies" component={Movies} />
          <Route path="/stats" component={Stats} />
          <Route path="/settings" component={Settings} />
          <Route path="/sceneFilenameParser" component={SceneFilenameParser} />
          <Route
            path="/sceneDuplicateChecker"
            component={SceneDuplicateChecker}
          />
          <Route path="/setup" component={Setup} />
          <Route path="/migrate" component={Migrate} />
          <Route component={PageNotFound} />
        </Switch>
      </Suspense>
    );
  }

  return (
    <ErrorBoundary>
      <IntlProvider locale={language} messages={messages} formats={intlFormats}>
        <ConfigurationProvider
          configuration={config.data?.configuration}
          loading={config.loading}
        >
          <ToastProvider>
            <Suspense fallback={<LoadingIndicator />}>
              <LightboxProvider>
                <ManualProvider>
                  <InteractiveProvider>
                    <Helmet
                      titleTemplate={`%s ${TITLE_SUFFIX}`}
                      defaultTitle="Stash"
                    />
                    {maybeRenderNavbar()}
                    <div className="main container-fluid">
                      {renderContent()}
                    </div>
                  </InteractiveProvider>
                </ManualProvider>
              </LightboxProvider>
            </Suspense>
          </ToastProvider>
        </ConfigurationProvider>
      </IntlProvider>
    </ErrorBoundary>
  );
};
