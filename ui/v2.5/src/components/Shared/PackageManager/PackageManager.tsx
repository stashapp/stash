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

export type InstalledPackage = Omit<GQL.Package, "requires">;

const InstalledPackagesList: React.FC<{
  filter: string;
  loading?: boolean;
  updatesLoaded: boolean;
  packages: InstalledPackage[];
  checkedPackages: InstalledPackage[];
  setCheckedPackages: React.Dispatch<React.SetStateAction<InstalledPackage[]>>;
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

  function togglePackage(pkg: InstalledPackage) {
    if (loading) return;

    setCheckedPackages((prev) => {
      if (prev.includes(pkg)) {
        return prev.filter((n) => n.id !== pkg.id);
      } else {
        return prev.concat(pkg);
      }
    });
  }

  function rowClassname(pkg: InstalledPackage) {
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
  checkedPackages: InstalledPackage[];
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
  packages: InstalledPackage[];
  updatesLoaded: boolean;
  onCheckForUpdates: () => void;
  onUpdatePackages: (packages: InstalledPackage[]) => void;
  onUninstallPackages: (packages: InstalledPackage[]) => void;
}> = ({
  packages,
  onCheckForUpdates,
  updatesLoaded,
  onUpdatePackages,
  onUninstallPackages,
  loading,
}) => {
  const [checkedPackages, setCheckedPackages] = useState<InstalledPackage[]>(
    []
  );
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
  installEnabled: boolean;
  onInstallPackages: () => void;
  selectedOnly: boolean;
  setSelectedOnly: (v: boolean) => void;
}> = ({
  installEnabled,
  onInstallPackages,
  loading,
  filter,
  setFilter,
  selectedOnly,
  setSelectedOnly,
}) => {
  const intl = useIntl();

  const selectedOnlyId = !selectedOnly
    ? "package_manager.selected_only"
    : "package_manager.show_all";

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
          variant="secondary"
          onClick={() => setSelectedOnly(!selectedOnly)}
        >
          <FormattedMessage id={selectedOnlyId} />
        </Button>
        <Button
          variant="primary"
          disabled={!installEnabled || loading}
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

export type RemotePackage = Omit<GQL.Package, "requires"> & {
  requires: { id: string }[];
};

const SourcePackagesList: React.FC<{
  filter: string;
  selectedOnly: boolean;
  loading?: boolean;
  source: GQL.PackageSource;
  loadSource: () => Promise<RemotePackage[]>;
  selectedPackages: RemotePackage[];
  setSelectedPackages: React.Dispatch<React.SetStateAction<RemotePackage[]>>;
  editSource: () => void;
  deleteSource: () => void;
}> = ({
  source,
  loadSource,
  selectedPackages,
  setSelectedPackages,
  loading,
  filter,
  selectedOnly,
  editSource,
  deleteSource,
}) => {
  const intl = useIntl();
  const [packages, setPackages] = useState<RemotePackage[]>();
  const [sourceOpen, setSourceOpen] = useState(false);

  const checkedMap = useMemo(() => {
    const map: Record<string, boolean> = {};

    selectedPackages.forEach((pkg) => {
      map[pkg.id] = true;
    });
    return map;
  }, [selectedPackages]);

  const sourceChecked = useMemo(() => {
    return packages && Object.keys(checkedMap).length === packages.length;
  }, [checkedMap, packages]);

  const filteredPackages = useMemo(() => {
    if (!packages) return [];

    let ret = filterPackages(packages, filter);

    if (selectedOnly) {
      ret = ret.filter((pkg) => checkedMap[pkg.id]);
    }

    return ret;
  }, [filter, packages, selectedOnly, checkedMap]);

  function togglePackage(pkg: RemotePackage) {
    if (loading || !packages) return;

    setSelectedPackages((prev) => {
      const selected = prev.find((p) => p.id === pkg.id);

      if (selected) {
        return prev.filter((n) => n.id !== pkg.id);
      } else {
        // also include required packages
        const toAdd = [pkg];
        pkg.requires.forEach((r) => {
          // find the required package
          const requiredSelected = prev.find((p) => p.id === r.id);
          const required = packages.find((p) => p.id === r.id);

          if (!requiredSelected && required) {
            toAdd.push(required);
          }
        });

        return prev.concat(...toAdd);
      }
    });
  }

  function toggleSource() {
    if (loading || packages === undefined) return;

    if (sourceChecked) {
      setSelectedPackages([]);
    } else {
      setSelectedPackages(packages.slice());
    }
  }

  async function toggleSourceOpen() {
    if (packages === undefined) {
      // need to load
      try {
        const loaded = await loadSource();
        setPackages(loaded);
      } catch (e) {
        // TODO - handle
        console.error(e);
      }
    }

    setSourceOpen((prev) => !prev);
  }

  function renderCollapseButton() {
    return (
      <Button
        variant="minimal"
        size="sm"
        className="package-collapse-button"
        onClick={() => toggleSourceOpen()}
      >
        <Icon icon={sourceOpen ? faChevronDown : faChevronRight} />
      </Button>
    );
  }

  function renderRequiredBy(pkg: RemotePackage) {
    const requiredBy = selectedPackages.filter((p) => {
      return p.requires.find((r) => r.id === pkg.id);
    });

    if (!requiredBy.length) return;

    return (
      <div className="package-required-by">
        <FormattedMessage
          id="package_manager.required_by"
          values={{ packages: requiredBy.map((p) => p.name).join(", ") }}
        />
      </div>
    );
  }

  const children = sourceOpen
    ? filteredPackages.map((pkg) => (
        <tr key={pkg.id}>
          <td colSpan={2}>
            <Form.Check
              checked={checkedMap[pkg.id] ?? false}
              onChange={() => togglePackage(pkg)}
              disabled={loading}
            />
          </td>
          <td className="package-cell" onClick={() => togglePackage(pkg)}>
            <span className="package-name">{pkg.name}</span>
            <span className="package-id">{pkg.id}</span>
          </td>
          <td>
            <span className="package-version">{pkg.version}</span>
            <span className="package-date">{formatDate(pkg.date)}</span>
          </td>
          <td>
            {renderRequiredBy(pkg)}
            <div>{pkg.description}</div>
          </td>
        </tr>
      )) ?? []
    : [];

  return (
    <>
      <tr key={source.url} className="package-source">
        <td>
          {packages !== undefined ? (
            <Form.Check
              checked={sourceChecked ?? false}
              onChange={() => toggleSource()}
              disabled={loading}
            />
          ) : undefined}
        </td>
        <td>{renderCollapseButton()}</td>
        <td colSpan={2} onClick={() => toggleSourceOpen()}>
          <span>{source.name ?? source.url}</span>
        </td>
        <td className="source-controls">
          <Button
            size="sm"
            variant="primary"
            title={intl.formatMessage({ id: "actions.edit" })}
            onClick={() => editSource()}
          >
            <FormattedMessage id="actions.edit" />
          </Button>
          <Button
            size="sm"
            variant="danger"
            title={intl.formatMessage({ id: "actions.delete" })}
            onClick={() => deleteSource()}
          >
            <FormattedMessage id="actions.delete" />
          </Button>
        </td>
      </tr>
      {...children}
    </>
  );
};

const AvailablePackagesList: React.FC<{
  filter: string;
  selectedOnly: boolean;
  loading?: boolean;
  sources: GQL.PackageSource[];
  loadSource: (source: string) => Promise<RemotePackage[]>;
  selectedPackages: Record<string, RemotePackage[]>; // map of source url to selected packages
  setSelectedPackages: React.Dispatch<
    React.SetStateAction<Record<string, RemotePackage[]>>
  >;
  addSource: (src: GQL.PackageSource) => void;
  editSource: (existing: GQL.PackageSource, changed: GQL.PackageSource) => void;
  deleteSource: (source: GQL.PackageSource) => void;
}> = ({
  sources,
  loadSource,
  selectedPackages,
  setSelectedPackages,
  loading,
  filter,
  selectedOnly,
  addSource,
  editSource,
  deleteSource,
}) => {
  const [deletingSource, setDeletingSource] = useState<GQL.PackageSource>();
  const [editingSource, setEditingSource] = useState<GQL.PackageSource>();
  const [addingSource, setAddingSource] = useState(false);

  function onDeleteSource() {
    if (!deletingSource) return;

    deleteSource(deletingSource);
    setDeletingSource(undefined);
  }

  function setSelectedSourcePackages(
    src: GQL.PackageSource,
    v: RemotePackage[] | ((prevState: RemotePackage[]) => RemotePackage[])
  ) {
    setSelectedPackages((prev) => {
      const existing = prev[src.url] ?? [];
      const next = typeof v === "function" ? v(existing) : v;

      return {
        ...prev,
        [src.url]: next,
      };
    });
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
            {sources.map((src) => (
              <SourcePackagesList
                key={src.url}
                filter={filter}
                selectedOnly={selectedOnly}
                loading={loading}
                source={src}
                loadSource={() => loadSource(src.url)}
                selectedPackages={selectedPackages[src.url] ?? []}
                setSelectedPackages={(v) => setSelectedSourcePackages(src, v)}
                editSource={() => setEditingSource(src)}
                deleteSource={() => setDeletingSource(src)}
              />
            ))}
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
  loadSource: (source: string) => Promise<RemotePackage[]>;
  onInstallPackages: (packages: GQL.PackageSpecInput[]) => void;
  addSource: (src: GQL.PackageSource) => void;
  editSource: (existing: GQL.PackageSource, changed: GQL.PackageSource) => void;
  deleteSource: (source: GQL.PackageSource) => void;
}> = ({
  sources,
  loadSource,
  onInstallPackages,
  loading,
  addSource,
  editSource,
  deleteSource,
}) => {
  const [checkedPackages, setCheckedPackages] = useState<
    Record<string, RemotePackage[]>
  >({});
  const [filter, setFilter] = useState("");
  const [selectedOnly, setSelectedOnly] = useState(false);

  const installEnabled = useMemo(() => {
    return Object.values(checkedPackages).some((s) => s.length > 0);
  }, [checkedPackages]);

  function toPackageSpecInput(): GQL.PackageSpecInput[] {
    const ret: GQL.PackageSpecInput[] = [];
    Object.keys(checkedPackages).forEach((sourceURL) => {
      checkedPackages[sourceURL].forEach((pkg) => {
        ret.push({ id: pkg.id, sourceURL });
      });
    });
    return ret;
  }

  return (
    <div className="available-packages">
      <AvailablePackagesToolbar
        filter={filter}
        setFilter={(f) => setFilter(f)}
        loading={loading}
        installEnabled={installEnabled}
        onInstallPackages={() => onInstallPackages(toPackageSpecInput())}
        selectedOnly={selectedOnly}
        setSelectedOnly={setSelectedOnly}
      />
      <AvailablePackagesList
        filter={filter}
        selectedOnly={selectedOnly}
        loading={loading}
        sources={sources}
        loadSource={loadSource}
        selectedPackages={checkedPackages}
        setSelectedPackages={setCheckedPackages}
        addSource={addSource}
        editSource={editSource}
        deleteSource={deleteSource}
      />
    </div>
  );
};
