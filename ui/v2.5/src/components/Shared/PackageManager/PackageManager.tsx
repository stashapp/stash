import { Button, Form, Table } from "react-bootstrap";
import React, { useEffect, useState, useMemo, useContext } from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  queryAvailableScraperPackages,
  useInstalledScraperPackages,
  useInstalledScraperPackagesStatus,
} from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { Icon } from "../Icon";
import {
  faChevronDown,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";

type packageSpec = {
  source: string;
  id: string;
};

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
  updatesLoaded: boolean;
  packages: GQL.Package[];
  checkedPackages: string[];
  setCheckedPackages: React.Dispatch<React.SetStateAction<string[]>>;
}> = ({ packages, checkedPackages, setCheckedPackages, updatesLoaded }) => {
  const checkedMap = useMemo(() => {
    const map: Record<string, boolean> = {};
    checkedPackages.forEach((id) => {
      map[id] = true;
    });
    return map;
  }, [checkedPackages]);

  const allChecked = useMemo(() => {
    return checkedPackages.length === packages.length;
  }, [checkedPackages, packages]);

  function toggleAllChecked() {
    setCheckedPackages(allChecked ? [] : packages.map((pkg) => pkg.id));
  }

  function togglePackage(id: string) {
    setCheckedPackages((prev) => {
      if (prev.includes(id)) {
        return prev.filter((n) => n !== id);
      } else {
        return prev.concat(id);
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
            <Form.Check checked={allChecked} onChange={toggleAllChecked} />
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
                checked={checkedMap[pkg.id]}
                onChange={() => togglePackage(pkg.id)}
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
  checkedPackages: string[];
  onCheckForUpdates: () => void;
  onUpdatePackages: () => void;
  onUninstallPackages: () => void;
}> = ({
  checkedPackages,
  onCheckForUpdates,
  onUpdatePackages,
  onUninstallPackages,
}) => {
  // TODO - alert for uninstall

  return (
    <div className="package-manager-toolbar">
      <Button variant="primary" onClick={() => onCheckForUpdates()}>
        <FormattedMessage id="package_manager.check_for_updates" />
      </Button>
      <Button
        variant="primary"
        disabled={!checkedPackages.length}
        onClick={() => onUpdatePackages()}
      >
        <FormattedMessage id="package_manager.update" />
      </Button>
      <Button
        variant="danger"
        disabled={!checkedPackages.length}
        onClick={() => onUninstallPackages()}
      >
        <FormattedMessage id="package_manager.uninstall" />
      </Button>
    </div>
  );
};

const InstalledPackages: React.FC<{
  packages: GQL.Package[];
  updatesLoaded: boolean;
  onCheckForUpdates: () => void;
}> = ({ packages, onCheckForUpdates, updatesLoaded }) => {
  const [checkedPackages, setCheckedPackages] = useState<string[]>([]);

  return (
    <div className="installed-packages">
      <InstalledPackagesToolbar
        checkedPackages={checkedPackages}
        onCheckForUpdates={onCheckForUpdates}
        onUpdatePackages={() => {}}
        onUninstallPackages={() => {}}
      />
      <InstalledPackagesList
        packages={packages}
        checkedPackages={checkedPackages}
        setCheckedPackages={setCheckedPackages}
        updatesLoaded={updatesLoaded}
      />
    </div>
  );
};

const AvailablePackagesToolbar: React.FC<{
  checkedPackages: packageSpec[];
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
  sources: GQL.PackageSource[];
  packages: Record<string, GQL.Package[]>;
  loadSource: (source: string) => void;
  checkedPackages: packageSpec[];
  setCheckedPackages: React.Dispatch<React.SetStateAction<packageSpec[]>>;
}> = ({
  sources,
  packages,
  loadSource,
  checkedPackages,
  setCheckedPackages,
}) => {
  const [sourceOpen, setSourceOpen] = useState<Record<string, boolean>>({});

  const checkedMap = useMemo(() => {
    const map: Record<string, Record<string, boolean>> = {};
    checkedPackages.forEach((pkg) => {
      if (!map[pkg.source]) {
        map[pkg.source] = {};
      }
      map[pkg.source][pkg.id] = true;
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

  function togglePackage(source: string, id: string) {
    function isPackage(s: packageSpec) {
      return s.id === id && s.source === source;
    }

    setCheckedPackages((prev) => {
      if (prev.find((s) => s.id === id && s.source === source)) {
        return prev.filter((n) => !isPackage(n));
      } else {
        return prev.concat({ id, source });
      }
    });
  }

  function toggleSource(source: string) {
    if (sourceOpen[source] === undefined) {
      return;
    }

    if (sourceChecked[source]) {
      setCheckedPackages((prev) => {
        return prev.filter((n) => n.source !== source);
      });
    } else {
      setCheckedPackages((prev) => {
        return prev
          .filter((n) => n.source !== source)
          .concat(
            packages[source]?.map((pkg) => ({ id: pkg.id, source })) ?? []
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
                checked={checkedMap[source.url]?.[pkg.id]}
                onChange={() => togglePackage(source.url, pkg.id)}
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
              checked={sourceChecked[source.url]}
              onChange={() => toggleSource(source.url)}
            />
          ) : undefined}
        </td>
        <td>{renderCollapseButton(source.url)}</td>
        <td colSpan={4} onClick={() => toggleSource(source.url)}>
          {source.name ?? source.url}
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
  sources: GQL.PackageSource[];
  packages: Record<string, GQL.Package[]>;
  loadSource: (source: string) => void;
  onInstallPackages: () => void;
}> = ({ sources, packages, loadSource }) => {
  const [checkedPackages, setCheckedPackages] = useState<packageSpec[]>([]);

  return (
    <div className="installed-packages">
      <AvailablePackagesToolbar
        checkedPackages={checkedPackages}
        onInstallPackages={() => {}}
      />
      <AvailablePackagesList
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
  installedPackages: GQL.Package[];
  onCheckForUpdates: () => void;
  updatesLoaded: boolean;

  sources: GQL.PackageSource[];
  sourcePackages: Record<string, GQL.Package[]>;
  onLoadSource: (source: string) => void;
}

export const PackageManager: React.FC<IPackageManagerProps> = ({
  installedPackages,
  updatesLoaded,
  onCheckForUpdates,
  sources,
  sourcePackages,
  onLoadSource,
}) => {
  return (
    <div className="package-manager">
      <div>
        <h3>
          <FormattedMessage id={"config.scraping.installed_scrapers"} />
        </h3>
        <InstalledPackages
          packages={installedPackages}
          onCheckForUpdates={onCheckForUpdates}
          updatesLoaded={updatesLoaded}
        />
      </div>

      <div>
        <h3>
          <FormattedMessage id={"config.scraping.available_scrapers"} />
        </h3>
        <AvailablePackages
          onInstallPackages={() => {}}
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

  const { data: installedScrapers } = useInstalledScraperPackages();

  const { data: withStatus, refetch } = useInstalledScraperPackagesStatus({
    skip: !loadUpgrades,
  });

  function onCheckForUpdates() {
    if (!loadUpgrades) {
      setLoadUpgrades(true);
    } else {
      refetch();
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
      installedPackages={installedPackages}
      updatesLoaded={loadUpgrades}
      onCheckForUpdates={() => onCheckForUpdates()}
      sources={sources ?? []}
      sourcePackages={sourcePackages}
      onLoadSource={(source) => loadSource(source)}
    />
  );
};
