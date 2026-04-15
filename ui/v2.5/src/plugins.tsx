import React, { useEffect } from "react";
import { PatchFunction } from "./patch";
import { usePlugins } from "./core/StashService";
import { useMemoOnce } from "./hooks/state";
import { uniq } from "lodash-es";
import useScript, { useCSS } from "./hooks/useScript";
import { PluginsQuery } from "./core/generated-graphql";
import { LoadingIndicator } from "./components/Shared/LoadingIndicator";
import { FormattedMessage } from "react-intl";
import { useToast } from "./hooks/Toast";

type PluginList = NonNullable<Required<PluginsQuery["plugins"]>>;

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

// load all plugins and their dependencies
// returns true when all plugins are loaded, regardess of success or failure
// if disableCustomizations is true, skip loading plugins entirely
function useLoadPlugins(disableCustomizations?: boolean) {
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
    // Skip loading plugin JS if customizations are disabled.
    // Note: We check inside useMemoOnce rather than early-returning from useLoadPlugins
    // to comply with React's rules of hooks - hooks must be called unconditionally.
    if (disableCustomizations) {
      return [[], true];
    }
    return [
      uniq(
        sortedPlugins
          ?.filter((plugin) => plugin.enabled && plugin.paths.javascript)
          .map((plugin) => plugin.paths.javascript!)
          .flat() ?? []
      ),
      !!sortedPlugins && !pluginsLoading && !pluginsError,
    ];
  }, [sortedPlugins, pluginsLoading, pluginsError, disableCustomizations]);

  const pluginCSS = useMemoOnce(() => {
    if (disableCustomizations) {
      return [[], true];
    }
    return [
      uniq(
        sortedPlugins
          ?.filter((plugin) => plugin.enabled && plugin.paths.css)
          .map((plugin) => plugin.paths.css!)
          .flat() ?? []
      ),
      !!sortedPlugins && !pluginsLoading && !pluginsError,
    ];
  }, [sortedPlugins, pluginsLoading, pluginsError, disableCustomizations]);

  const pluginJavascriptLoaded = useScript(
    pluginJavascripts ?? [],
    !!pluginJavascripts && !pluginsLoading && !pluginsError
  );
  useCSS(pluginCSS ?? [], !pluginsLoading && !pluginsError);

  return {
    loading: !pluginsLoading && !!pluginJavascripts && pluginJavascriptLoaded,
    error: pluginsError,
  };
}

interface IPluginsLoaderProps {
  disableCustomizations?: boolean;
}

export const PluginsLoader: React.FC<
  React.PropsWithChildren<IPluginsLoaderProps>
> = ({ disableCustomizations, children }) => {
  const Toast = useToast();
  const { loading: loaded, error } = useLoadPlugins(disableCustomizations);

  useEffect(() => {
    if (error) {
      Toast.error(`Error loading plugins: ${error.message}`);
    }
  }, [Toast, error]);

  if (!loaded && !error)
    return (
      <LoadingIndicator message={<FormattedMessage id="loading.plugins" />} />
    );

  return <>{children}</>;
};

export const PluginRoutes: React.FC<React.PropsWithChildren<{}>> =
  PatchFunction("PluginRoutes", (props: React.PropsWithChildren<{}>) => {
    return <>{props.children}</>;
  }) as React.FC;
