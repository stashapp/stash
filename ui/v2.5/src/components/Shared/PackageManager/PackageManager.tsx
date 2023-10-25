import { Button, Form, Table } from "react-bootstrap";
import React, { useState, useMemo, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Icon";
import {
  faChevronDown,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";

type PackageSpec = GQL.PackageSpecInput & { name: string };

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

interface IPackage {
  id: string;
  name: string;
}

function filterPackages<T extends IPackage>(packages: T[], filter: string) {
  if (!filter) return packages;

  return packages.filter((pkg) => {
    return (
      pkg.name.toLowerCase().includes(filter.toLowerCase()) ||
      pkg.id.toLowerCase().includes(filter.toLowerCase())
    );
  });
}

const InstalledPackagesList: React.FC<{
  filter: string;
  loading?: boolean;
  updatesLoaded: boolean;
  packages: GQL.Package[];
  checkedPackages: GQL.Package[];
  setCheckedPackages: React.Dispatch<React.SetStateAction<GQL.Package[]>>;
}> = ({
  filter,
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

  const filteredPackages = useMemo(() => {
    return filterPackages(packages, filter);
  }, [filter, packages]);

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
        {filteredPackages.map((pkg) => (
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
  filter: string;
  setFilter: (s: string) => void;
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
  filter,
  setFilter,
}) => {
  const intl = useIntl();
  // TODO - alert for uninstall

  return (
    <div className="package-manager-toolbar">
      <div>
        <Form.Control
          placeholder={`${intl.formatMessage({ id: "filter" })}...`}
          className="text-input"
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
        />
      </div>
      <div>
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
  const [filter, setFilter] = useState("");

  const filteredPackages = useMemo(() => {
    return filterPackages(checkedPackages, filter);
  }, [checkedPackages, filter]);

  useEffect(() => {
    setCheckedPackages((prev) => {
      const newVal = prev.filter((pkg) =>
        packages.find((p) => p.id === pkg.id)
      );
      if (newVal.length !== prev.length) {
        return newVal;
      }

      return prev;
    });
  }, [checkedPackages, packages]);

  return (
    <div className="installed-packages">
      <InstalledPackagesToolbar
        filter={filter}
        setFilter={(f) => setFilter(f)}
        loading={loading}
        checkedPackages={filteredPackages}
        onCheckForUpdates={onCheckForUpdates}
        onUpdatePackages={() => onUpdatePackages(filteredPackages)}
        onUninstallPackages={() => onUninstallPackages(filteredPackages)}
      />
      <InstalledPackagesList
        filter={filter}
        loading={loading}
        packages={packages}
        // use original checked packages so that check boxes are not affected by filter
        checkedPackages={checkedPackages}
        setCheckedPackages={setCheckedPackages}
        updatesLoaded={updatesLoaded}
      />
    </div>
  );
};

const AvailablePackagesToolbar: React.FC<{
  filter: string;
  setFilter: (s: string) => void;
  loading?: boolean;
  checkedPackages: GQL.PackageSpecInput[];
  onInstallPackages: () => void;
}> = ({ checkedPackages, onInstallPackages, loading, filter, setFilter }) => {
  const intl = useIntl();

  return (
    <div className="package-manager-toolbar">
      <div>
        <Form.Control
          placeholder={`${intl.formatMessage({ id: "filter" })}...`}
          className="text-input"
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
        />
      </div>
      <div>
        <Button
          variant="primary"
          disabled={!checkedPackages.length || loading}
          onClick={() => onInstallPackages()}
        >
          <FormattedMessage id="package_manager.install" />
        </Button>
      </div>
    </div>
  );
};

const AvailablePackagesList: React.FC<{
  filter: string;
  loading?: boolean;
  sources: GQL.PackageSource[];
  packages: Record<string, GQL.Package[]>;
  loadSource: (source: string) => void;
  checkedPackages: PackageSpec[];
  setCheckedPackages: React.Dispatch<React.SetStateAction<PackageSpec[]>>;
}> = ({
  sources,
  packages,
  loadSource,
  checkedPackages,
  setCheckedPackages,
  loading,
  filter,
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

  const filteredPackages = useMemo(() => {
    const map: Record<string, GQL.Package[]> = {};

    Object.keys(packages).forEach((source) => {
      map[source] = filterPackages(packages[source], filter);
    });

    return map;
  }, [filter, packages]);

  function togglePackage(sourceURL: string, pkg: GQL.Package) {
    if (loading) return;
    function isPackage(s: GQL.PackageSpecInput) {
      return s.id === pkg.id && s.sourceURL === sourceURL;
    }

    setCheckedPackages((prev) => {
      if (prev.find((s) => s.id === pkg.id && s.sourceURL === sourceURL)) {
        return prev.filter((n) => !isPackage(n));
      } else {
        return prev.concat({ id: pkg.id, sourceURL, name: pkg.name });
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
            packages[sourceURL]?.map((pkg) => ({
              id: pkg.id,
              name: pkg.name,
              sourceURL,
            })) ?? []
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
      ? filteredPackages[source.url]?.map((pkg) => (
          <tr key={pkg.id}>
            <td colSpan={2}>
              <Form.Check
                checked={checkedMap[source.url]?.[pkg.id] ?? false}
                onChange={() => togglePackage(source.url, pkg)}
                disabled={loading}
              />
            </td>
            <td
              className="package-cell"
              onClick={() => togglePackage(source.url, pkg)}
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
  const [checkedPackages, setCheckedPackages] = useState<PackageSpec[]>([]);
  const [filter, setFilter] = useState("");

  const filteredPackages = useMemo(() => {
    return filterPackages(checkedPackages, filter);
  }, [checkedPackages, filter]);

  function toPackageSpecInput(i: PackageSpec[]): GQL.PackageSpecInput[] {
    return i.map((pkg) => ({
      id: pkg.id,
      sourceURL: pkg.sourceURL,
    }));
  }

  return (
    <div className="installed-packages">
      <AvailablePackagesToolbar
        filter={filter}
        setFilter={(f) => setFilter(f)}
        loading={loading}
        checkedPackages={filteredPackages}
        onInstallPackages={() =>
          onInstallPackages(toPackageSpecInput(filteredPackages))
        }
      />
      <AvailablePackagesList
        filter={filter}
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
