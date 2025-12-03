import React, { HTMLAttributes, useEffect, useMemo, useState } from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  Form,
  FormControl,
  InputGroup,
  Modal,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import {
  useConfigureUISetting,
  useFindSavedFilters,
  useSavedFilterDestroy,
  useSaveFilter,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  FilterMode,
  SavedFilterDataFragment,
} from "src/core/generated-graphql";
import { View } from "./views";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { faBookmark, faSave, faTimes } from "@fortawesome/free-solid-svg-icons";
import { AlertModal } from "../Shared/Alert";
import cx from "classnames";
import { TruncatedInlineText } from "../Shared/TruncatedText";
import { OperationButton } from "../Shared/OperationButton";
import { createPortal } from "react-dom";

const ExistingSavedFilterList: React.FC<{
  name: string;
  onSelect: (value: SavedFilterDataFragment) => void;
  savedFilters: SavedFilterDataFragment[];
  disabled?: boolean;
}> = ({ name, onSelect, savedFilters: existing, disabled = false }) => {
  const filtered = useMemo(() => {
    if (!name) return existing;

    return existing.filter((f) =>
      f.name.toLowerCase().includes(name.toLowerCase())
    );
  }, [existing, name]);

  return (
    <ul className="existing-filter-list">
      {filtered.map((f) => (
        <li key={f.id}>
          <Button
            className="minimal"
            variant="link"
            onClick={() => onSelect(f)}
            disabled={disabled}
          >
            {f.name}
          </Button>
        </li>
      ))}
    </ul>
  );
};

export const SaveFilterDialog: React.FC<{
  mode: FilterMode;
  onClose: (name?: string, id?: string) => void;
  isSaving?: boolean;
}> = ({ mode, onClose, isSaving = false }) => {
  const intl = useIntl();
  const [filterName, setFilterName] = useState("");

  const { data } = useFindSavedFilters(mode);

  const overwritingFilter = useMemo(() => {
    const savedFilters = data?.findSavedFilters ?? [];
    return savedFilters.find(
      (f) => f.name.toLowerCase() === filterName.toLowerCase()
    );
  }, [data?.findSavedFilters, filterName]);

  return (
    <Modal show className="save-filter-dialog">
      <Modal.Header>
        <FormattedMessage id="actions.save_filter" />
      </Modal.Header>
      <Modal.Body>
        <Form.Group>
          <Form.Label>
            <FormattedMessage id="filter_name" />
          </Form.Label>
          <FormControl
            className="bg-secondary text-white border-secondary"
            placeholder={`${intl.formatMessage({ id: "filter_name" })}…`}
            value={filterName}
            onChange={(e) => setFilterName(e.target.value)}
            disabled={isSaving}
          />
        </Form.Group>

        <ExistingSavedFilterList
          name={filterName}
          onSelect={(f) => setFilterName(f.name)}
          savedFilters={data?.findSavedFilters ?? []}
        />

        {!!overwritingFilter && (
          <span className="saved-filter-overwrite-warning">
            <FormattedMessage
              id="dialogs.overwrite_filter_warning"
              values={{
                entityName: overwritingFilter.name,
              }}
            />
          </span>
        )}
      </Modal.Body>
      <Modal.Footer>
        <Button
          variant="secondary"
          onClick={() => onClose()}
          disabled={isSaving}
        >
          {intl.formatMessage({ id: "actions.cancel" })}
        </Button>
        <OperationButton
          loading={isSaving}
          variant="primary"
          onClick={() => onClose(filterName, overwritingFilter?.id)}
        >
          {intl.formatMessage({ id: "actions.save" })}
        </OperationButton>
      </Modal.Footer>
    </Modal>
  );
};

export const LoadFilterDialog: React.FC<{
  mode: FilterMode;
  onClose: (filter?: SavedFilterDataFragment) => void;
}> = ({ mode, onClose }) => {
  const intl = useIntl();
  const [filterName, setFilterName] = useState("");

  const { data } = useFindSavedFilters(mode);

  return (
    <Modal show className="load-filter-dialog">
      <Modal.Header>
        <FormattedMessage id="actions.load_filter" />
      </Modal.Header>
      <Modal.Body>
        <Form.Group>
          <Form.Label>
            <FormattedMessage id="filter_name" />
          </Form.Label>
          <FormControl
            className="bg-secondary text-white border-secondary"
            placeholder={`${intl.formatMessage({ id: "filter_name" })}…`}
            value={filterName}
            onChange={(e) => setFilterName(e.target.value)}
          />
        </Form.Group>

        <ExistingSavedFilterList
          name={filterName}
          onSelect={(f) => onClose(f)}
          savedFilters={data?.findSavedFilters ?? []}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={() => onClose()}>
          {intl.formatMessage({ id: "actions.cancel" })}
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

const DeleteAlert: React.FC<{
  deletingFilter: SavedFilterDataFragment | undefined;
  onClose: (confirm?: boolean) => void;
}> = ({ deletingFilter, onClose }) => {
  if (!deletingFilter) {
    return null;
  }

  return (
    <Modal show>
      <Modal.Body>
        <FormattedMessage
          id="dialogs.delete_confirm"
          values={{
            entityName: deletingFilter.name,
          }}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="danger" onClick={() => onClose(true)}>
          <FormattedMessage id="actions.delete" />
        </Button>
        <Button variant="secondary" onClick={() => onClose()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

const OverwriteAlert: React.FC<{
  overwritingFilter: SavedFilterDataFragment | undefined;
  onClose: (confirm?: boolean) => void;
}> = ({ overwritingFilter, onClose }) => {
  if (!overwritingFilter) {
    return null;
  }

  return (
    <Modal show>
      <Modal.Body>
        <FormattedMessage
          id="dialogs.overwrite_filter_warning"
          values={{
            entityName: overwritingFilter.name,
          }}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="primary" onClick={() => onClose(true)}>
          <FormattedMessage id="actions.overwrite" />
        </Button>
        <Button variant="secondary" onClick={() => onClose()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

interface ISavedFilterListProps {
  filter: ListFilterModel;
  onSetFilter: (f: ListFilterModel) => void;
  view?: View;
  menuPortalTarget?: Element | DocumentFragment;
}

export const SavedFilterList: React.FC<ISavedFilterListProps> = ({
  filter,
  onSetFilter,
  view,
}) => {
  const Toast = useToast();
  const intl = useIntl();

  const { data, error, loading, refetch } = useFindSavedFilters(filter.mode);

  const [filterName, setFilterName] = useState("");
  const [saving, setSaving] = useState(false);
  const [deletingFilter, setDeletingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();
  const [overwritingFilter, setOverwritingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();

  const saveFilter = useSaveFilter();
  const [destroyFilter] = useSavedFilterDestroy();
  const [saveUISetting] = useConfigureUISetting();

  const savedFilters = data?.findSavedFilters ?? [];

  async function onSaveFilter(name: string, id?: string) {
    const filterCopy = filter.clone();

    try {
      setSaving(true);
      await saveFilter(filterCopy, name, id);

      Toast.success(
        intl.formatMessage(
          {
            id: "toast.saved_entity",
          },
          {
            entity: intl.formatMessage({ id: "filter" }).toLocaleLowerCase(),
          }
        )
      );
      setFilterName("");
      setOverwritingFilter(undefined);
      refetch();
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
    }
  }

  async function onDeleteFilter(f: SavedFilterDataFragment) {
    try {
      setSaving(true);

      await destroyFilter({
        variables: {
          input: {
            id: f.id,
          },
        },
      });

      Toast.success(
        intl.formatMessage(
          {
            id: "toast.delete_past_tense",
          },
          {
            count: 1,
            singularEntity: intl.formatMessage({ id: "filter" }),
            pluralEntity: intl.formatMessage({ id: "filters" }),
          }
        )
      );
      refetch();
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
      setDeletingFilter(undefined);
    }
  }

  async function onSetDefaultFilter() {
    if (!view) {
      return;
    }

    const filterCopy = filter.clone();

    try {
      setSaving(true);

      await saveUISetting({
        variables: {
          key: `defaultFilters.${view.toString()}`,
          value: {
            mode: filter.mode,
            find_filter: filterCopy.makeFindFilter(),
            object_filter: filterCopy.makeSavedFilter(),
            ui_options: filterCopy.makeSavedUIOptions(),
          },
        },
      });

      Toast.success(
        intl.formatMessage({
          id: "toast.default_filter_set",
        })
      );
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
    }
  }

  function filterClicked(f: SavedFilterDataFragment) {
    const newFilter = filter.clone();

    newFilter.currentPage = 1;
    // #1795 - reset search term if not present in saved filter
    newFilter.searchTerm = "";
    newFilter.configureFromSavedFilter(f);
    // #1507 - reset random seed when loaded
    newFilter.randomSeed = -1;

    onSetFilter(newFilter);
  }

  interface ISavedFilterItem {
    item: SavedFilterDataFragment;
  }
  const SavedFilterItem: React.FC<ISavedFilterItem> = ({ item }) => {
    return (
      <div className="dropdown-item-container">
        <Dropdown.Item onClick={() => filterClicked(item)} title={item.name}>
          <span>{item.name}</span>
        </Dropdown.Item>
        <ButtonGroup>
          <Button
            className="save-button"
            variant="secondary"
            size="sm"
            title={intl.formatMessage({ id: "actions.overwrite" })}
            onClick={(e) => {
              setOverwritingFilter(item);
              e.stopPropagation();
            }}
          >
            <Icon icon={faSave} />
          </Button>
          <Button
            className="delete-button"
            variant="secondary"
            size="sm"
            title={intl.formatMessage({ id: "actions.delete" })}
            onClick={(e) => {
              setDeletingFilter(item);
              e.stopPropagation();
            }}
          >
            <Icon icon={faTimes} />
          </Button>
        </ButtonGroup>
      </div>
    );
  };

  function renderSavedFilters() {
    if (error) return <h6 className="text-center">{error.message}</h6>;

    if (loading || saving) {
      return (
        <div className="loading">
          <LoadingIndicator message="" />
        </div>
      );
    }

    return (
      <ul className="saved-filter-list">
        {savedFilters
          .filter(
            (f) =>
              !filterName ||
              f.name.toLowerCase().includes(filterName.toLowerCase())
          )
          .map((f) => (
            <SavedFilterItem key={f.name} item={f} />
          ))}
      </ul>
    );
  }

  function maybeRenderSetDefaultButton() {
    if (view) {
      return (
        <div className="mt-1">
          <Dropdown.Item
            as={Button}
            title={intl.formatMessage({ id: "actions.set_as_default" })}
            className="set-as-default-button"
            variant="secondary"
            size="sm"
            onClick={() => onSetDefaultFilter()}
          >
            {intl.formatMessage({ id: "actions.set_as_default" })}
          </Dropdown.Item>
        </div>
      );
    }
  }

  return (
    <>
      <DeleteAlert
        deletingFilter={deletingFilter}
        onClose={(confirm) => {
          if (confirm) {
            onDeleteFilter(deletingFilter!);
          }
          setDeletingFilter(undefined);
        }}
      />
      <OverwriteAlert
        overwritingFilter={overwritingFilter}
        onClose={(confirm) => {
          if (confirm) {
            onSaveFilter(overwritingFilter!.name, overwritingFilter!.id);
          }
          setOverwritingFilter(undefined);
        }}
      />
      <InputGroup>
        <FormControl
          className="bg-secondary text-white border-secondary"
          placeholder={`${intl.formatMessage({ id: "filter_name" })}…`}
          value={filterName}
          onChange={(e) => setFilterName(e.target.value)}
        />
        <InputGroup.Append>
          <OverlayTrigger
            placement="top"
            overlay={
              <Tooltip id="filter-tooltip">
                <FormattedMessage id="actions.save_filter" />
              </Tooltip>
            }
          >
            <Button
              disabled={
                !filterName || !!savedFilters.find((f) => f.name === filterName)
              }
              variant="secondary"
              onClick={() => {
                onSaveFilter(filterName);
              }}
            >
              <Icon icon={faSave} />
            </Button>
          </OverlayTrigger>
        </InputGroup.Append>
      </InputGroup>
      {renderSavedFilters()}
      {maybeRenderSetDefaultButton()}
    </>
  );
};

interface ISavedFilterItem {
  item: SavedFilterDataFragment;
  onClick: () => void;
  onDelete: () => void;
  selected?: boolean;
}

const SavedFilterItem: React.FC<ISavedFilterItem> = ({
  item,
  onClick,
  onDelete,
  selected = false,
}) => {
  const intl = useIntl();

  return (
    <li className="saved-filter-item">
      <a onClick={onClick}>
        <div className="label-group">
          <TruncatedInlineText
            className={cx("no-icon-margin", { selected })}
            text={item.name}
          />
        </div>
        <div>
          <Button
            className="delete-button"
            variant="minimal"
            size="sm"
            title={intl.formatMessage({ id: "actions.delete" })}
            onClick={(e) => {
              onDelete();
              e.stopPropagation();
            }}
          >
            <Icon fixedWidth icon={faTimes} />
          </Button>
        </div>
      </a>
    </li>
  );
};

const SavedFilters: React.FC<{
  error?: string;
  loading?: boolean;
  saving?: boolean;
  savedFilters: SavedFilterDataFragment[];
  onFilterClicked: (f: SavedFilterDataFragment) => void;
  onDeleteClicked: (f: SavedFilterDataFragment) => void;
  currentFilterID?: string;
}> = ({
  error,
  loading,
  saving,
  savedFilters,
  onFilterClicked,
  onDeleteClicked,
  currentFilterID,
}) => {
  if (error) return <h6 className="text-center">{error}</h6>;

  if (loading || saving) {
    return (
      <div className="loading">
        <LoadingIndicator message="" />
      </div>
    );
  }

  return (
    <ul className="saved-filter-list">
      {savedFilters.map((f) => (
        <SavedFilterItem
          key={f.name}
          item={f}
          onClick={() => onFilterClicked(f)}
          onDelete={() => onDeleteClicked(f)}
          selected={currentFilterID === f.id}
        />
      ))}
    </ul>
  );
};

export const SidebarSavedFilterList: React.FC<ISavedFilterListProps> = ({
  filter,
  onSetFilter,
  view,
}) => {
  const Toast = useToast();
  const intl = useIntl();

  const [currentSavedFilter, setCurrentSavedFilter] = useState<{
    id: string;
    set: boolean;
  }>();

  const { data, error, loading, refetch } = useFindSavedFilters(filter.mode);

  const [filterName, setFilterName] = useState("");
  const [saving, setSaving] = useState(false);
  const [deletingFilter, setDeletingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();
  const [showSaveDialog, setShowSaveDialog] = useState(false);
  const [settingDefault, setSettingDefault] = useState(false);

  const saveFilter = useSaveFilter();
  const [destroyFilter] = useSavedFilterDestroy();
  const [saveUISetting] = useConfigureUISetting();

  const filteredFilters = useMemo(() => {
    const savedFilters = data?.findSavedFilters ?? [];
    if (!filterName) return savedFilters;

    return savedFilters.filter(
      (f) =>
        !filterName || f.name.toLowerCase().includes(filterName.toLowerCase())
    );
  }, [data?.findSavedFilters, filterName]);

  // handle when filter is changed to de-select the current filter
  useEffect(() => {
    // HACK - first change will be from setting the filter
    // second change is likely from somewhere else
    setCurrentSavedFilter((v) => {
      if (!v) return v;

      if (v.set) {
        setCurrentSavedFilter({ id: v.id, set: false });
      } else {
        setCurrentSavedFilter(undefined);
      }
    });
  }, [filter]);

  async function onSaveFilter(name: string, id?: string) {
    try {
      setSaving(true);
      await saveFilter(filter, name, id);

      Toast.success(
        intl.formatMessage(
          {
            id: "toast.saved_entity",
          },
          {
            entity: intl.formatMessage({ id: "filter" }).toLocaleLowerCase(),
          }
        )
      );
      setFilterName("");
      setShowSaveDialog(false);
      refetch();
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
    }
  }

  async function onDeleteFilter(f: SavedFilterDataFragment) {
    try {
      setSaving(true);

      await destroyFilter({
        variables: {
          input: {
            id: f.id,
          },
        },
      });

      Toast.success(
        intl.formatMessage(
          {
            id: "toast.delete_past_tense",
          },
          {
            count: 1,
            singularEntity: intl.formatMessage({ id: "filter" }),
            pluralEntity: intl.formatMessage({ id: "filters" }),
          }
        )
      );
      refetch();
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
      setDeletingFilter(undefined);
    }
  }

  async function onSetDefaultFilter() {
    if (!view) {
      return;
    }

    const filterCopy = filter.clone();

    try {
      setSaving(true);

      await saveUISetting({
        variables: {
          key: `defaultFilters.${view.toString()}`,
          value: {
            mode: filter.mode,
            find_filter: filterCopy.makeFindFilter(),
            object_filter: filterCopy.makeSavedFilter(),
            ui_options: filterCopy.makeSavedUIOptions(),
          },
        },
      });

      Toast.success(
        intl.formatMessage({
          id: "toast.default_filter_set",
        })
      );
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
      setSettingDefault(false);
    }
  }

  function filterClicked(f: SavedFilterDataFragment) {
    const newFilter = filter.clone();

    newFilter.currentPage = 1;
    // #1795 - reset search term if not present in saved filter
    newFilter.searchTerm = "";
    newFilter.configureFromSavedFilter(f);
    // #1507 - reset random seed when loaded
    newFilter.randomSeed = -1;

    setCurrentSavedFilter({ id: f.id, set: true });
    onSetFilter(newFilter);
  }

  return (
    <div className="sidebar-saved-filter-list-container">
      <DeleteAlert
        deletingFilter={deletingFilter}
        onClose={(confirm) => {
          if (confirm) {
            onDeleteFilter(deletingFilter!);
          }
          setDeletingFilter(undefined);
        }}
      />
      {showSaveDialog && (
        <SaveFilterDialog
          mode={filter.mode}
          onClose={(name, id) => {
            setShowSaveDialog(false);
            if (name) {
              onSaveFilter(name, id);
            }
          }}
        />
      )}
      <AlertModal
        show={!!settingDefault}
        text={<FormattedMessage id="dialogs.set_default_filter_confirm" />}
        confirmVariant="primary"
        onConfirm={() => onSetDefaultFilter()}
        onCancel={() => setSettingDefault(false)}
      />

      <div className="toolbar">
        <Button
          className="minimal save-filter-button"
          size="sm"
          onClick={() => setShowSaveDialog(true)}
        >
          <span>
            <FormattedMessage id="actions.save_filter" />
          </span>
        </Button>
        <Button
          className="minimal set-as-default-button"
          variant="secondary"
          size="sm"
          onClick={() => setSettingDefault(true)}
        >
          <FormattedMessage id="actions.set_as_default" />
        </Button>
      </div>

      <FormControl
        className="bg-secondary text-white border-secondary saved-filter-search-input"
        placeholder={`${intl.formatMessage({ id: "filter_name" })}…`}
        value={filterName}
        onChange={(e) => setFilterName(e.target.value)}
      />
      <SavedFilters
        error={error?.message}
        loading={loading}
        saving={saving}
        savedFilters={filteredFilters}
        onFilterClicked={filterClicked}
        onDeleteClicked={setDeletingFilter}
        currentFilterID={currentSavedFilter?.id}
      />
    </div>
  );
};

export const SavedFilterDropdown: React.FC<ISavedFilterListProps> = (props) => {
  const SavedFilterDropdownRef = React.forwardRef<
    HTMLDivElement,
    HTMLAttributes<HTMLDivElement>
  >(({ style, className }: HTMLAttributes<HTMLDivElement>, ref) => (
    <div ref={ref} style={style} className={className}>
      <SavedFilterList {...props} />
    </div>
  ));
  SavedFilterDropdownRef.displayName = "SavedFilterDropdown";

  const menu = (
    <Dropdown.Menu
      as={SavedFilterDropdownRef}
      className="saved-filter-list-menu"
    />
  );

  return (
    <Dropdown as={ButtonGroup} className="saved-filter-dropdown">
      <OverlayTrigger
        placement="top"
        overlay={
          <Tooltip id="filter-tooltip">
            <FormattedMessage id="search_filter.saved_filters" />
          </Tooltip>
        }
      >
        <Dropdown.Toggle variant="secondary">
          <Icon icon={faBookmark} />
        </Dropdown.Toggle>
      </OverlayTrigger>
      {props.menuPortalTarget
        ? createPortal(menu, props.menuPortalTarget)
        : menu}
    </Dropdown>
  );
};
