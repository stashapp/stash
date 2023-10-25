import { Button, Form, Table } from "react-bootstrap";
import React, { useEffect, useState, useMemo, useContext } from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  queryAvailableScraperPackages,
  useInstallScraperPackages,
  useInstalledScraperPackages,
  useInstalledScraperPackagesStatus,
  useUninstallScraperPackages,
  useUpdateScraperPackages,
} from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { Icon } from "../Icon";
import {
  faChevronDown,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";
import { useMonitorJob } from "src/utils/job";

function formatVersion(
  version: string | undefined | null,
  date: string | undefined | null
) {
  let ret = version ?? "";
  if (date) {
    const parsedDate = new Date(date);
    if (version) {
      ret += "-";
    }
    ret += `${parsedDate.toISOString()}`;
  }

  return ret;
}

function formatDate(date: string | undefined | null) {
  if (!date) return;

  return new Date(date).toISOString();
}

const InstalledPackagesList: React.FC<{
  loading?: boolean;
  updatesLoaded: boolean;
  packages: GQL.Package[];
  checkedPackages: GQL.Package[];
  setCheckedPackages: React.Dispatch<React.SetStateAction<GQL.Package[]>>;
}> = ({
  packages,
  checkedPackages,
  setCheckedPackages,
  updatesLoaded,
  loading,
}) => {
  const checkedMap = useMemo(() => {
    const map: Record<string, boolean> = {};
    checkedPackages.forEach((pkg) => {
      map[pkg.id] = true;
    });
    return map;
  }, [checkedPackages]);

  const allChecked = useMemo(() => {
    return checkedPackages.length === packages.length;
  }, [checkedPackages, packages]);

  function toggleAllChecked() {
    setCheckedPackages(allChecked ? [] : packages.slice());
  }

  function togglePackage(pkg: GQL.Package) {
    if (loading) return;

    setCheckedPackages((prev) => {
      if (prev.includes(pkg)) {
        return prev.filter((n) => n.id !== pkg.id);
      } else {
        return prev.concat(pkg);
      }
    });
  }

  function rowClassname(pkg: GQL.Package) {
    if (pkg.upgrade?.package.version) {
      return "package-update-available";
    }
  }

  return (
    <Table>
      <thead>
        <tr>
          <th>
            <Form.Check
              checked={allChecked ?? false}
              onChange={toggleAllChecked}
              disabled={loading}
            />
          </th>
          <th>
            <FormattedMessage id="package_manager.package" />
          </th>
          <th>
            <FormattedMessage id="package_manager.installed_version" />
          </th>
          {updatesLoaded ? (
            <th>
              <FormattedMessage id="package_manager.latest_version" />
            </th>
          ) : undefined}
        </tr>
      </thead>
      <tbody>
        {packages.map((pkg) => (
          <tr key={pkg.id} className={rowClassname(pkg)}>
            <td>
              <Form.Check
                checked={checkedMap[pkg.id] ?? false}
                disabled={loading}
                onChange={() => togglePackage(pkg)}
              />
            </td>
            <td>
              <span className="package-name">{pkg.name}</span>
              <span className="package-id">{pkg.id}</span>
            </td>
            <td>
              <span className="package-version">{pkg.version}</span>
              <span className="package-date">{formatDate(pkg.date)}</span>
            </td>
            {updatesLoaded ? (
              <td>
                {formatVersion(
                  pkg.upgrade?.package.version,
                  pkg.upgrade?.package.date
                )}
              </td>
            ) : undefined}
          </tr>
        ))}
      </tbody>
    </Table>
  );
};

const InstalledPackagesToolbar: React.FC<{
  loading?: boolean;
  checkedPackages: GQL.Package[];
  onCheckForUpdates: () => void;
  onUpdatePackages: () => void;
  onUninstallPackages: () => void;
}> = ({
  loading,
  checkedPackages,
  onCheckForUpdates,
  onUpdatePackages,
  onUninstallPackages,
}) => {
  // TODO - alert for uninstall

  return (
    <div className="package-manager-toolbar">
      <Button
        variant="primary"
        onClick={() => onCheckForUpdates()}
        disabled={loading}
      >
        <FormattedMessage id="package_manager.check_for_updates" />
      </Button>
      <Button
        variant="primary"
        disabled={!checkedPackages.length || loading}
        onClick={() => onUpdatePackages()}
      >
        <FormattedMessage id="package_manager.update" />
      </Button>
      <Button
        variant="danger"
        disabled={!checkedPackages.length || loading}
        onClick={() => onUninstallPackages()}
      >
        <FormattedMessage id="package_manager.uninstall" />
      </Button>
    </div>
  );
};

const InstalledPackages: React.FC<{
  loading?: boolean;
  packages: GQL.Package[];
  updatesLoaded: boolean;
  onCheckForUpdates: () => void;
  onUpdatePackages: (packages: GQL.Package[]) => void;
  onUninstallPackages: (packages: GQL.Package[]) => void;
}> = ({
  packages,
  onCheckForUpdates,
  updatesLoaded,
  onUpdatePackages,
  onUninstallPackages,
  loading,
}) => {
  const [checkedPackages, setCheckedPackages] = useState<GQL.Package[]>([]);

  return (
    <div className="installed-packages">
      <InstalledPackagesToolbar
        loading={loading}
        checkedPackages={checkedPackages}
        onCheckForUpdates={onCheckForUpdates}
        onUpdatePackages={() => onUpdatePackages(checkedPackages)}
        onUninstallPackages={() => onUninstallPackages(checkedPackages)}
      />
      <InstalledPackagesList
        loading={loading}
        packages={packages}
        checkedPackages={checkedPackages}
        setCheckedPackages={setCheckedPackages}
        updatesLoaded={updatesLoaded}
      />
    </div>
  );
};

const AvailablePackagesToolbar: React.FC<{
  loading?: boolean;
  checkedPackages: GQL.PackageSpecInput[];
  onInstallPackages: () => void;
}> = ({ checkedPackages, onInstallPackages }) => {
  return (
    <div className="package-manager-toolbar">
      <Button
        variant="primary"
        disabled={!checkedPackages.length}
        onClick={() => onInstallPackages()}
      >
        <FormattedMessage id="package_manager.install" />
      </Button>
    </div>
  );
};

const AvailablePackagesList: React.FC<{
  loading?: boolean;
  sources: GQL.PackageSource[];
  packages: Record<string, GQL.Package[]>;
  loadSource: (source: string) => void;
  checkedPackages: GQL.PackageSpecInput[];
  setCheckedPackages: React.Dispatch<
    React.SetStateAction<GQL.PackageSpecInput[]>
  >;
}> = ({
  sources,
  packages,
  loadSource,
  checkedPackages,
  setCheckedPackages,
  loading,
}) => {
  const [sourceOpen, setSourceOpen] = useState<Record<string, boolean>>({});

  const checkedMap = useMemo(() => {
    const map: Record<string, Record<string, boolean>> = {};
    checkedPackages.forEach((pkg) => {
      if (!map[pkg.sourceURL]) {
        map[pkg.sourceURL] = {};
      }
      map[pkg.sourceURL][pkg.id] = true;
    });
    return map;
  }, [checkedPackages]);

  const sourceChecked = useMemo(() => {
    const map: Record<string, boolean> = {};

    Object.keys(checkedMap).forEach((source) => {
      map[source] =
        Object.keys(checkedMap[source]).length === packages[source]?.length;
    });

    return map;
  }, [checkedMap, packages]);

  function togglePackage(sourceURL: string, id: string) {
    if (loading) return;
    function isPackage(s: GQL.PackageSpecInput) {
      return s.id === id && s.sourceURL === sourceURL;
    }

    setCheckedPackages((prev) => {
      if (prev.find((s) => s.id === id && s.sourceURL === sourceURL)) {
        return prev.filter((n) => !isPackage(n));
      } else {
        return prev.concat({ id, sourceURL });
      }
    });
  }

  function toggleSource(sourceURL: string) {
    if (loading) return;
    if (sourceOpen[sourceURL] === undefined) {
      return;
    }

    if (sourceChecked[sourceURL]) {
      setCheckedPackages((prev) => {
        return prev.filter((n) => n.sourceURL !== sourceURL);
      });
    } else {
      setCheckedPackages((prev) => {
        return prev
          .filter((n) => n.sourceURL !== sourceURL)
          .concat(
            packages[sourceURL]?.map((pkg) => ({ id: pkg.id, sourceURL })) ?? []
          );
      });
    }
  }

  function toggleSourceOpen(source: string) {
    if (sourceOpen[source] === undefined) {
      // need to load
      loadSource(source);
    }

    setSourceOpen((prev) => {
      return {
        ...prev,
        [source]: !prev[source],
      };
    });
  }

  function renderCollapseButton(source: string) {
    return (
      <Button
        variant="minimal"
        className="package-collapse-button"
        onClick={() => toggleSourceOpen(source)}
      >
        <Icon icon={sourceOpen[source] ? faChevronDown : faChevronRight} />
      </Button>
    );
  }

  function renderSource(source: GQL.PackageSource) {
    const children = sourceOpen[source.url]
      ? packages[source.url]?.map((pkg) => (
          <tr key={pkg.id}>
            <td colSpan={2}>
              <Form.Check
                checked={checkedMap[source.url]?.[pkg.id] ?? false}
                onChange={() => togglePackage(source.url, pkg.id)}
                disabled={loading}
              />
            </td>
            <td
              className="package-cell"
              onClick={() => togglePackage(source.url, pkg.id)}
            >
              <span className="package-name">{pkg.name}</span>
              <span className="package-id">{pkg.id}</span>
            </td>
            <td>
              <span className="package-version">{pkg.version}</span>
              <span className="package-date">{formatDate(pkg.date)}</span>
            </td>
            <td>
              {formatVersion(
                pkg.upgrade?.package.version,
                pkg.upgrade?.package.date
              )}
            </td>
            <td>{pkg.description}</td>
          </tr>
        )) ?? []
      : [];

    return [
      <tr key={source.url} className="package-source">
        <td>
          {sourceOpen[source.url] !== undefined ? (
            <Form.Check
              checked={sourceChecked[source.url] ?? false}
              onChange={() => toggleSource(source.url)}
              disabled={loading}
            />
          ) : undefined}
        </td>
        <td>{renderCollapseButton(source.url)}</td>
        <td colSpan={4} onClick={() => toggleSourceOpen(source.url)}>
          <FormattedMessage id="package_manager.source" />
          {": "}
          <span>{source.name ?? source.url}</span>
        </td>
      </tr>,
      ...children,
    ];
  }

  return (
    <Table>
      <thead>
        <tr>
          <th></th>
          <th></th>
          <th>
            <FormattedMessage id="package_manager.package" />
          </th>
          <th>
            <FormattedMessage id="package_manager.version" />
          </th>
          <th>
            <FormattedMessage id="package_manager.description" />
          </th>
        </tr>
      </thead>
      <tbody>{sources.map((pkg) => renderSource(pkg))}</tbody>
    </Table>
  );
};

const AvailablePackages: React.FC<{
  loading?: boolean;
  sources: GQL.PackageSource[];
  packages: Record<string, GQL.Package[]>;
  loadSource: (source: string) => void;
  onInstallPackages: (packages: GQL.PackageSpecInput[]) => void;
}> = ({ sources, packages, loadSource, onInstallPackages, loading }) => {
  const [checkedPackages, setCheckedPackages] = useState<
    GQL.PackageSpecInput[]
  >([]);

  return (
    <div className="installed-packages">
      <AvailablePackagesToolbar
        loading={loading}
        checkedPackages={checkedPackages}
        onInstallPackages={() => onInstallPackages(checkedPackages)}
      />
      <AvailablePackagesList
        loading={loading}
        sources={sources}
        loadSource={loadSource}
        packages={packages}
        checkedPackages={checkedPackages}
        setCheckedPackages={setCheckedPackages}
      />
    </div>
  );
};

export interface IPackageManagerProps {
  loading?: boolean;
  installedPackages: GQL.Package[];
  onCheckForUpdates: () => void;
  updatesLoaded: boolean;

  sources: GQL.PackageSource[];
  sourcePackages: Record<string, GQL.Package[]>;
  onLoadSource: (source: string) => void;

  onInstallPackages: (packages: GQL.PackageSpecInput[]) => void;
  onUpdatePackages: (packages: GQL.Package[]) => void;
  onUninstallPackages: (packages: GQL.Package[]) => void;
}

export const PackageManager: React.FC<IPackageManagerProps> = ({
  loading,
  installedPackages,
  updatesLoaded,
  onCheckForUpdates,
  sources,
  sourcePackages,
  onLoadSource,
  onInstallPackages,
  onUpdatePackages,
  onUninstallPackages,
}) => {
  return (
    <div className="package-manager">
      <div>
        <h3>
          <FormattedMessage id={"config.scraping.installed_scrapers"} />
        </h3>
        <InstalledPackages
          loading={loading}
          packages={installedPackages}
          onCheckForUpdates={onCheckForUpdates}
          onUpdatePackages={onUpdatePackages}
          onUninstallPackages={onUninstallPackages}
          updatesLoaded={updatesLoaded}
        />
      </div>

      <div>
        <h3>
          <FormattedMessage id={"config.scraping.available_scrapers"} />
        </h3>
        <AvailablePackages
          loading={loading}
          onInstallPackages={onInstallPackages}
          loadSource={onLoadSource}
          sources={sources}
          packages={sourcePackages}
        />
      </div>
    </div>
  );
};

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
  const { job } = useMonitorJob(jobID, () => refetchPackages());

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
