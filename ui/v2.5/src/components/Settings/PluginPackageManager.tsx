import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  evictQueries,
  getClient,
  queryAvailablePluginPackages,
  useInstalledPluginPackages,
  mutateInstallPluginPackages,
  mutateUninstallPluginPackages,
  mutateUpdatePluginPackages,
  pluginMutationImpactedQueries,
  isLoading,
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

  const { data, previousData, refetch, networkStatus, error } =
    useInstalledPluginPackages(loadUpgrades);

  const loading = isLoading(networkStatus);

  async function onUpdatePackages(packages: GQL.PackageSpecInput[]) {
    const r = await mutateUpdatePluginPackages(packages);

    setJobID(r.data?.updatePackages);
  }

  async function onUninstallPackages(packages: GQL.PackageSpecInput[]) {
    const r = await mutateUninstallPluginPackages(packages);

    setJobID(r.data?.uninstallPackages);
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
      refetch();
    }
  }

  // when loadUpgrades changes from false to true, data is set to undefined while the request is loading
  // so use previousData as a fallback, which will be the result when loadUpgrades was false,
  // to prevent displaying a "No packages found" message
  const installedPackages =
    data?.installedPackages ?? previousData?.installedPackages ?? [];

  return (
    <SettingSection headingID="config.plugins.installed_plugins">
      <div className="package-manager">
        <InstalledPackages
          loading={!!job || loading}
          error={error?.message}
          packages={installedPackages}
          onCheckForUpdates={onCheckForUpdates}
          onUpdatePackages={(packages) =>
            onUpdatePackages(
              packages.map((p) => ({
                id: p.package_id,
                sourceURL: p.sourceURL,
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
          updatesLoaded={loadUpgrades && !loading}
        />
      </div>
    </SettingSection>
  );
};

export const AvailablePluginPackages: React.FC = () => {
  const { general, loading: configLoading, error, saveGeneral } = useSettings();

  const [jobID, setJobID] = useState<string>();
  const { job } = useMonitorJob(jobID, () => onPackageChanges());

  // Get installed packages to filter them out from available list
  const { data: installedData } = useInstalledPluginPackages(false);
  const installedPackageIds = new Set(
    installedData?.installedPackages?.map((p) => p.package_id) ?? []
  );

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
    // Filter out already installed packages
    return data.availablePackages.filter(
      (pkg) => !installedPackageIds.has(pkg.package_id)
    );
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
