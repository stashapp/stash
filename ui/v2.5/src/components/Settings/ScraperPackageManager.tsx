import React, { useEffect, useState, useMemo, useContext } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  evictQueries,
  getClient,
  queryAvailableScraperPackages,
  useInstallScraperPackages,
  useInstalledScraperPackages,
  useInstalledScraperPackagesStatus,
  useUninstallScraperPackages,
  useUpdateScraperPackages,
} from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { useMonitorJob } from "src/utils/job";
import { PackageManager } from "../Shared/PackageManager/PackageManager";

export const ScraperPackageManager: React.FC = () => {
  const { configuration } = useContext(ConfigurationContext);
  const [loadUpgrades, setLoadUpgrades] = useState(false);
  const [sourcePackages, setSourcePackages] = useState<
    Record<string, GQL.Package[]>
  >({});
  const [sources, setSources] = useState<GQL.PackageSource[]>();
  const [sourcesLoaded, setSourcesLoaded] = useState<Record<string, boolean>>(
    {}
  );
  const [jobID, setJobID] = useState<string>();
  const { job } = useMonitorJob(jobID, () => onPackageChanges());

  const { data: installedScrapers, refetch: refetchPackages1 } =
    useInstalledScraperPackages({
      skip: loadUpgrades,
    });

  const {
    data: withStatus,
    refetch: refetchPackages2,
    loading: statusLoading,
  } = useInstalledScraperPackagesStatus({
    skip: !loadUpgrades,
  });

  const [installPackages] = useInstallScraperPackages();
  const [updatePackages] = useUpdateScraperPackages();
  const [uninstallPackages] = useUninstallScraperPackages();

  async function onInstallPackages(packages: GQL.PackageSpecInput[]) {
    const r = await installPackages({
      variables: {
        packages,
      },
    });

    setJobID(r.data?.installPackages);
  }

  async function onUpdatePackages(packages: GQL.PackageSpecInput[]) {
    const r = await updatePackages({
      variables: {
        packages,
      },
    });

    setJobID(r.data?.updatePackages);
  }

  async function onUninstallPackages(packages: string[]) {
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
    const impactedQueries = [
      GQL.ListPerformerScrapersDocument,
      GQL.ListSceneScrapersDocument,
      GQL.ListMovieScrapersDocument,
    ];

    refetchPackages();

    // job is complete, refresh all local data
    const ac = getClient();
    evictQueries(ac.cache, impactedQueries);
  }

  function onCheckForUpdates() {
    if (!loadUpgrades) {
      setLoadUpgrades(true);
    } else {
      refetchPackages();
    }
  }

  useEffect(() => {
    if (!sources && configuration?.general.scraperPackageSources) {
      setSources(configuration.general.scraperPackageSources);
    }
  }, [sources, configuration?.general.scraperPackageSources]);

  const installedPackages = useMemo(() => {
    if (withStatus?.installedPackages) {
      return withStatus.installedPackages;
    }

    return installedScrapers?.installedPackages ?? [];
  }, [installedScrapers, withStatus]);

  async function loadSource(source: string) {
    if (sourcesLoaded[source]) {
      return;
    }

    const { data } = await queryAvailableScraperPackages(source);

    setSourcePackages((prev) => {
      return {
        ...prev,
        [source]: data.availablePackages,
      };
    });

    setSourcesLoaded((prev) => {
      return {
        ...prev,
        [source]: true,
      };
    });
  }

  return (
    <PackageManager
      loading={!!job || statusLoading}
      installedPackages={installedPackages}
      updatesLoaded={loadUpgrades}
      onCheckForUpdates={() => onCheckForUpdates()}
      sources={sources ?? []}
      sourcePackages={sourcePackages}
      onLoadSource={(source) => loadSource(source)}
      onInstallPackages={(packages) => onInstallPackages(packages)}
      onUpdatePackages={(packages) =>
        onUpdatePackages(
          packages.map((p) => ({ id: p.id, sourceURL: p.upgrade!.sourceURL }))
        )
      }
      onUninstallPackages={(packages) =>
        onUninstallPackages(packages.map((p) => p.id))
      }
    />
  );
};
