import React, { useState, useMemo } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  evictQueries,
  getClient,
  queryAvailablePluginPackages,
  useInstalledPluginPackages,
  useInstalledPluginPackagesStatus,
  mutateInstallPluginPackages,
  mutateUninstallPluginPackages,
  mutateUpdatePluginPackages,
  pluginMutationImpactedQueries,
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

export const InstalledPluginPackages: React.FC = () => {
  const [loadUpgrades, setLoadUpgrades] = useState(false);
  const [jobID, setJobID] = useState<string>();
  const { job } = useMonitorJob(jobID, () => onPackageChanges());

  const {
    data: installedPlugins,
    refetch: refetchPackages1,
    loading: loading1,
    error: error1,
  } = useInstalledPluginPackages({
    skip: loadUpgrades,
  });

  const {
    data: withStatus,
    refetch: refetchPackages2,
    loading: loading2,
    error: error2,
  } = useInstalledPluginPackagesStatus({
    skip: !loadUpgrades,
  });

  async function onUpdatePackages(packages: GQL.PackageSpecInput[]) {
    const r = await mutateUpdatePluginPackages(packages);

    setJobID(r.data?.updatePackages);
  }

  async function onUninstallPackages(packages: GQL.PackageSpecInput[]) {
    const r = await mutateUninstallPluginPackages(packages);

    setJobID(r.data?.uninstallPackages);
  }

  function refetchPackages() {
    refetchPackages1();
    refetchPackages2();
  }

  function onPackageChanges() {
    // job is complete, refresh all local data
    const ac = getClient();
    evictQueries(ac.cache, pluginMutationImpactedQueries);
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

  const loading = !!job || loading1 || loading2;
  const error = error1 || error2;

  return (
    <SettingSection headingID="config.plugins.installed_plugins">
      <div className="package-manager">
        <InstalledPackages
          loading={loading}
          error={error?.message}
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

  const [jobID, setJobID] = useState<string>();
  const { job } = useMonitorJob(jobID, () => onPackageChanges());

  async function onInstallPackages(packages: GQL.PackageSpecInput[]) {
    const r = await mutateInstallPluginPackages(packages);

    setJobID(r.data?.installPackages);
  }

  function onPackageChanges() {
    // job is complete, refresh all local data
    const ac = getClient();
    evictQueries(ac.cache, pluginMutationImpactedQueries);
  }

  async function loadSource(source: string): Promise<RemotePackage[]> {
    const { data } = await queryAvailablePluginPackages(source);
    return data.availablePackages;
  }

  function addSource(source: GQL.PackageSource) {
    saveGeneral({
      pluginPackageSources: [...(general.pluginPackageSources ?? []), source],
    });
  }

  function editSource(existing: GQL.PackageSource, changed: GQL.PackageSource) {
    saveGeneral({
      pluginPackageSources: general.pluginPackageSources?.map((s) =>
        s.url === existing.url ? changed : s
      ),
    });
  }

  function deleteSource(source: GQL.PackageSource) {
    saveGeneral({
      pluginPackageSources: general.pluginPackageSources?.filter(
        (s) => s.url !== source.url
      ),
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

  const sources = general?.pluginPackageSources ?? [];

  return (
    <SettingSection headingID="config.plugins.available_plugins">
      <div className="package-manager">
        <AvailablePackages
          loading={loading}
          onInstallPackages={onInstallPackages}
          renderDescription={renderDescription}
          loadSource={(source) => loadSource(source)}
          sources={sources}
          addSource={addSource}
          editSource={editSource}
          deleteSource={deleteSource}
        />
      </div>
    </SettingSection>
  );
};
