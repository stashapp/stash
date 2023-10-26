import { Button, Form, Table } from "react-bootstrap";
import React, { useState, useMemo, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Icon";
import {
  faChevronDown,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";
import { SettingModal } from "src/components/Settings/Inputs";
import * as yup from "yup";
import { FormikErrors, yupToFormErrors } from "formik";
import { AlertModal } from "../Alert";

type PackageSpec = GQL.PackageSpecInput & { name: string };

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
    <div className="package-manager-table-container">
      <Table>
        <thead>
          <tr>
            <th className="button-cell">
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
                  <span className="package-version">
                    {pkg.upgrade?.package.version}
                  </span>
                  <span className="package-date">
                    {formatDate(pkg.upgrade?.package.date)}
                  </span>
                </td>
              ) : undefined}
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
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

export const InstalledPackages: React.FC<{
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

const EditSourceModal: React.FC<{
  sources: GQL.PackageSource[];
  existing?: GQL.PackageSource;
  onClose: (source?: GQL.PackageSource) => void;
}> = ({ existing, sources, onClose }) => {
  const intl = useIntl();

  const schema = yup.object({
    name: yup
      .string()
      .required()
      .test({
        name: "name",
        test: (value) => {
          if (!value) return true;
          const found = sources.find((s) => s.name === value);
          return !found || found === existing;
        },
        message: intl.formatMessage({ id: "validation.unique" }),
      }),
    url: yup
      .string()
      .required()
      .test({
        name: "url",
        test: (value) => {
          if (!value) return true;
          const found = sources.find((s) => s.url === value);
          return !found || found === existing;
        },
        message: intl.formatMessage({ id: "validation.unique" }),
      }),
  });

  type InputValues = yup.InferType<typeof schema>;
  function validate(
    v: GQL.PackageSource | undefined
  ): FormikErrors<InputValues> | undefined {
    try {
      schema.validateSync(v, { abortEarly: false });
    } catch (e) {
      return yupToFormErrors(e);
    }
  }

  const headerID = !!existing
    ? "package_manager.edit_source"
    : "package_manager.add_source";

  function renderField(
    v: GQL.PackageSource | undefined,
    setValue: (v: GQL.PackageSource | undefined) => void
  ) {
    const errors = validate(v);

    return (
      <>
        <Form.Group id="package-source-name">
          <h6>
            <FormattedMessage id="package_manager.source.name" />
          </h6>
          <Form.Control
            placeholder={intl.formatMessage({
              id: "package_manager.source.name",
            })}
            className="text-input"
            value={v?.name ?? ""}
            isInvalid={!!errors?.name}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setValue({ ...v!, name: e.currentTarget.value })
            }
          />
          <Form.Control.Feedback type="invalid">
            {errors?.name}
          </Form.Control.Feedback>
        </Form.Group>

        <Form.Group id="package-source-url">
          <h6>
            <FormattedMessage id="package_manager.source.url" />
          </h6>
          <Form.Control
            placeholder={intl.formatMessage({
              id: "package_manager.source.url",
            })}
            className="text-input"
            value={v?.url}
            isInvalid={!!errors?.url}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setValue({ ...v!, url: e.currentTarget.value.trim() })
            }
          />
          <Form.Control.Feedback type="invalid">
            {errors?.url}
          </Form.Control.Feedback>
        </Form.Group>
      </>
    );
  }

  return (
    <SettingModal<GQL.PackageSource>
      headingID={headerID}
      value={existing ?? { url: "", name: "" }}
      validate={(v) => validate(v) === undefined}
      renderField={renderField}
      close={onClose}
    />
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
  addSource: (src: GQL.PackageSource) => void;
  editSource: (existing: GQL.PackageSource, changed: GQL.PackageSource) => void;
  deleteSource: (source: GQL.PackageSource) => void;
}> = ({
  sources,
  packages,
  loadSource,
  checkedPackages,
  setCheckedPackages,
  loading,
  filter,
  addSource,
  editSource,
  deleteSource,
}) => {
  const intl = useIntl();
  const [sourceOpen, setSourceOpen] = useState<Record<string, boolean>>({});
  const [deletingSource, setDeletingSource] = useState<GQL.PackageSource>();
  const [editingSource, setEditingSource] = useState<GQL.PackageSource>();
  const [addingSource, setAddingSource] = useState(false);

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
        size="sm"
        className="package-collapse-button"
        onClick={() => toggleSourceOpen(source)}
      >
        <Icon icon={sourceOpen[source] ? faChevronDown : faChevronRight} />
      </Button>
    );
  }

  function onDeleteSource() {
    if (!deletingSource) return;

    deleteSource(deletingSource);
    setDeletingSource(undefined);
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
        <td colSpan={2} onClick={() => toggleSourceOpen(source.url)}>
          <span>{source.name ?? source.url}</span>
        </td>
        <td className="source-controls">
          <Button
            size="sm"
            variant="primary"
            title={intl.formatMessage({ id: "actions.edit" })}
            onClick={() => setEditingSource(source)}
          >
            <FormattedMessage id="actions.edit" />
          </Button>
          <Button
            size="sm"
            variant="danger"
            title={intl.formatMessage({ id: "actions.delete" })}
            onClick={() => setDeletingSource(source)}
          >
            <FormattedMessage id="actions.delete" />
          </Button>
        </td>
      </tr>,
      ...children,
    ];
  }

  return (
    <>
      <AlertModal
        show={!!deletingSource}
        text={
          <FormattedMessage
            id="package_manager.confirm_delete_source"
            values={{ name: deletingSource?.name, url: deletingSource?.url }}
          />
        }
        onConfirm={() => onDeleteSource()}
        onCancel={() => setDeletingSource(undefined)}
      />

      {editingSource || addingSource ? (
        <EditSourceModal
          sources={sources}
          existing={editingSource}
          onClose={(v) => {
            if (v) {
              if (addingSource) addSource(v);
              else if (editingSource) editSource(editingSource, v);
            }
            setEditingSource(undefined);
            setAddingSource(false);
          }}
        />
      ) : undefined}

      <div className="package-manager-table-container">
        <Table>
          <thead>
            <tr>
              <th className="button-cell"></th>
              <th className="button-cell"></th>
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
          <tbody>
            {sources.map((pkg) => renderSource(pkg))}
            <tr className="package-source">
              <td colSpan={2}></td>
              <td colSpan={3} onClick={() => setAddingSource(true)}>
                <Button size="sm" variant="success">
                  <FormattedMessage id="package_manager.add_source" />
                </Button>
              </td>
            </tr>
          </tbody>
        </Table>
      </div>
    </>
  );
};

export const AvailablePackages: React.FC<{
  loading?: boolean;
  sources: GQL.PackageSource[];
  packages: Record<string, GQL.Package[]>;
  loadSource: (source: string) => void;
  onInstallPackages: (packages: GQL.PackageSpecInput[]) => void;
  addSource: (src: GQL.PackageSource) => void;
  editSource: (existing: GQL.PackageSource, changed: GQL.PackageSource) => void;
  deleteSource: (source: GQL.PackageSource) => void;
}> = ({
  sources,
  packages,
  loadSource,
  onInstallPackages,
  loading,
  addSource,
  editSource,
  deleteSource,
}) => {
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
    <div className="available-packages">
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
        addSource={addSource}
        editSource={editSource}
        deleteSource={deleteSource}
      />
    </div>
  );
};
