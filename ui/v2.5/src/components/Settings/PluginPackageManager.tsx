import React, { useEffect, useState, useMemo } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  evictQueries,
  getClient,
  queryAvailablePluginPackages,
  useInstallPluginPackages,
  useInstalledPluginPackages,
  useInstalledPluginPackagesStatus,
  useUninstallPluginPackages,
  useUpdatePluginPackages,
} from "src/core/StashService";
import { useMonitorJob } from "src/utils/job";
import {
  AvailablePackages,
  InstalledPackages,
  RemotePackage,
} from "../Shared/PackageManager/PackageManager";
import { useSettings } from "./context";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { SettingSection } from "./SettingSection";

const impactedPackageChangeQueries = [
  GQL.PluginsDocument,
  GQL.PluginTasksDocument,
  GQL.InstalledPluginPackagesDocument,
  GQL.InstalledPluginPackagesStatusDocument,
];

export const InstalledPluginPackages: React.FC = () => {
  const [loadUpgrades, setLoadUpgrades] = useState(false);
  const [jobID, setJobID] = useState<string>();
  const { job } = useMonitorJob(jobID, () => onPackageChanges());

  const { data: installedPlugins, refetch: refetchPackages1 } =
    useInstalledPluginPackages({
      skip: loadUpgrades,
    });

  const {
    data: withStatus,
    refetch: refetchPackages2,
    loading: statusLoading,
  } = useInstalledPluginPackagesStatus({
    skip: !loadUpgrades,
  });

  const [updatePackages] = useUpdatePluginPackages();
  const [uninstallPackages] = useUninstallPluginPackages();

  async function onUpdatePackages(packages: GQL.PackageSpecInput[]) {
    const r = await updatePackages({
      variables: {
        packages,
      },
    });

    setJobID(r.data?.updatePackages);
  }

  async function onUninstallPackages(packages: GQL.PackageSpecInput[]) {
    const r = await uninstallPackages({
      variables: {
        packages,
      },
    });

    setJobID(r.data?.uninstallPackages);
  }

  function refetchPackages() {
    refetchPackages1();
    refetchPackages2();
  }

  function onPackageChanges() {
    // job is complete, refresh all local data
    const ac = getClient();
    evictQueries(ac.cache, impactedPackageChangeQueries);
  }

  function onCheckForUpdates() {
    if (!loadUpgrades) {
      setLoadUpgrades(true);
    } else {
      refetchPackages();
    }
  }

  const installedPackages = useMemo(() => {
    if (withStatus?.installedPackages) {
      return withStatus.installedPackages;
    }

    return installedPlugins?.installedPackages ?? [];
  }, [installedPlugins, withStatus]);

  const loading = !!job || statusLoading;

  return (
    <SettingSection headingID="config.plugins.installed_plugins">
      <div className="package-manager">
        <InstalledPackages
          loading={loading}
          packages={installedPackages}
          onCheckForUpdates={onCheckForUpdates}
          onUpdatePackages={(packages) =>
            onUpdatePackages(
              packages.map((p) => ({
                id: p.package_id,
                sourceURL: p.upgrade!.sourceURL,
              }))
            )
          }
          onUninstallPackages={(packages) =>
            onUninstallPackages(
              packages.map((p) => ({
                id: p.package_id,
                sourceURL: p.sourceURL,
              }))
            )
          }
          updatesLoaded={loadUpgrades}
        />
      </div>
    </SettingSection>
  );
};

export const AvailablePluginPackages: React.FC = () => {
  const { general, loading: configLoading, error, saveGeneral } = useSettings();

  const [sources, setSources] = useState<GQL.PackageSource[]>();
  const [jobID, setJobID] = useState<string>();
  const { job } = useMonitorJob(jobID, () => onPackageChanges());

  const [installPackages] = useInstallPluginPackages();

  async function onInstallPackages(packages: GQL.PackageSpecInput[]) {
    const r = await installPackages({
      variables: {
        packages,
      },
    });

    setJobID(r.data?.installPackages);
  }

  function onPackageChanges() {
    // job is complete, refresh all local data
    const ac = getClient();
    evictQueries(ac.cache, impactedPackageChangeQueries);
  }

  useEffect(() => {
    if (!sources && !configLoading && general.pluginPackageSources) {
      setSources(general.pluginPackageSources);
    }
  }, [sources, configLoading, general.pluginPackageSources]);

  async function loadSource(source: string): Promise<RemotePackage[]> {
    const { data } = await queryAvailablePluginPackages(source);
    return data.availablePackages;
  }

  function addSource(source: GQL.PackageSource) {
    saveGeneral({
      pluginPackageSources: [...(general.pluginPackageSources ?? []), source],
    });

    setSources((prev) => {
      return [...(prev ?? []), source];
    });
  }

  function editSource(existing: GQL.PackageSource, changed: GQL.PackageSource) {
    saveGeneral({
      pluginPackageSources: general.pluginPackageSources?.map((s) =>
        s.url === existing.url ? changed : s
      ),
    });

    setSources((prev) => {
      return prev?.map((s) => (s.url === existing.url ? changed : s));
    });
  }

  function deleteSource(source: GQL.PackageSource) {
    saveGeneral({
      pluginPackageSources: general.pluginPackageSources?.filter(
        (s) => s.url !== source.url
      ),
    });

    setSources((prev) => {
      return prev?.filter((s) => s.url !== source.url);
    });
  }

  function renderDescription(pkg: RemotePackage) {
    if (pkg.metadata.description) {
      return pkg.metadata.description;
    }
  }

  if (error) return <h1>{error.message}</h1>;
  if (configLoading) return <LoadingIndicator />;

  const loading = !!job;

  return (
    <SettingSection headingID="config.plugins.available_plugins">
      <div className="package-manager">
        <AvailablePackages
          loading={loading}
          onInstallPackages={onInstallPackages}
          renderDescription={renderDescription}
          loadSource={(source) => loadSource(source)}
          sources={sources ?? []}
          addSource={addSource}
          editSource={editSource}
          deleteSource={deleteSource}
        />
      </div>
    </SettingSection>
  );
};
