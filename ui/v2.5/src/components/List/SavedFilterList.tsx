import React, {
  HTMLAttributes,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
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
  savedFilterFields,
  useConfigureUISetting,
  useFindSavedFilters,
  useRenameSavedFilter,
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
import {
  faBookmark,
  faEllipsisV,
  faSave,
} from "@fortawesome/free-solid-svg-icons";
import { AlertModal } from "../Shared/Alert";
import cx from "classnames";
import { TruncatedInlineText } from "../Shared/TruncatedText";
import { OperationButton } from "../Shared/OperationButton";
import { createPortal } from "react-dom";
import { PatchFunction } from "src/patch";
import {
  OperationDropdown,
  OperationDropdownItem,
} from "./ListOperationButtons";

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

// Renames a saved filter in place. A modal would be portalled to the document
// body, which the toolbar's dropdown reads as a click outside itself, closing
// it mid-rename.
const SavedFilterNameInput: React.FC<{
  name: string;
  onCommit: (name: string) => void;
  onCancel: () => void;
}> = ({ name: currentName, onCommit, onCancel }) => {
  const intl = useIntl();
  const [name, setName] = useState(currentName);

  function commit() {
    // the backend trims the name and rejects an empty one
    const newName = name.trim();

    if (!newName || newName === currentName) {
      onCancel();
    } else {
      onCommit(newName);
    }
  }

  return (
    <FormControl
      className="bg-secondary text-white border-secondary saved-filter-name-input"
      size="sm"
      autoFocus
      value={name}
      placeholder={`${intl.formatMessage({ id: "filter_name" })}…`}
      onChange={(e) => setName(e.target.value)}
      onBlur={onCancel}
      onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
          commit();
          e.preventDefault();
        } else if (e.key === "Escape") {
          // an enclosing dropdown closes on escape from two separate events -
          // react-overlays' Dropdown on keydown, useRootClose on keyup - so
          // stop both, and abandon on keyup while this input still exists
          e.stopPropagation();
          e.preventDefault();
        }
      }}
      onKeyUp={(e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Escape") {
          e.stopPropagation();
          onCancel();
        }
      }}
    />
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

// the operations available for a single saved filter. The menu is portalled to
// the document body so that it is not clipped by the scrolling filter list.
const SavedFilterOperations: React.FC<{
  // unique across every saved filter list on the page, not just within one -
  // the toolbar dropdown stays mounted once opened, so both lists render a
  // row per filter at the same time
  id: string;
  onOperationChosen?: () => void;
  onOverwrite: () => void;
  onRename: () => void;
  onSetDefault?: () => void;
  onDelete: () => void;
}> = ({
  id,
  onOperationChosen,
  onOverwrite,
  onRename,
  onSetDefault,
  onDelete,
}) => {
  const intl = useIntl();

  function operation(action: () => void) {
    return () => {
      onOperationChosen?.();
      action();
    };
  }

  return (
    <OperationDropdown
      className="saved-filter-operations"
      menuPortalTarget={document.body}
      icon={faEllipsisV}
      id={`saved-filter-operations-${id}`}
      title={intl.formatMessage({ id: "operations" })}
      toggleClassName="minimal"
      size="sm"
      // don't load the filter when interacting with its operations
      onClick={(e) => e.stopPropagation()}
    >
      <OperationDropdownItem
        text={intl.formatMessage({ id: "actions.overwrite" })}
        onClick={operation(onOverwrite)}
      />
      <OperationDropdownItem
        text={intl.formatMessage({ id: "actions.rename" })}
        onClick={operation(onRename)}
      />
      {onSetDefault && (
        <OperationDropdownItem
          text={intl.formatMessage({ id: "actions.set_as_default" })}
          onClick={operation(onSetDefault)}
        />
      )}
      <OperationDropdownItem
        // text-danger rather than a rule of our own - the menu items carry
        // bootstrap's text-white, which is !important and so cannot be
        // overridden by specificity
        className="text-danger"
        text={intl.formatMessage({ id: "actions.delete" })}
        onClick={operation(onDelete)}
      />
    </OperationDropdown>
  );
};

interface ISavedFilterListProps {
  filter: ListFilterModel;
  onSetFilter: (f: ListFilterModel) => void;
  view?: View;
  menuPortalTarget?: Element | DocumentFragment;
  // Called when an operation is chosen, so that an enclosing dropdown can
  // stay open - the operation happens inside it. A callback rather than an
  // event check, since useRootClose tests containment during capture, and by
  // the time the close arrives the clicked item may be gone.
  onOperationChosen?: () => void;
}

export interface ISavedFilterLoaded {
  savedFilter: SavedFilterDataFragment;
  filter: ListFilterModel;
  source: "dialog" | "sidebar" | "toolbar";
  view?: View;
}

export const notifySavedFilterLoaded = PatchFunction(
  "SavedFilter.Loaded",
  (event: ISavedFilterLoaded) => event
);

// renaming and setting a default operate on an existing saved filter rather
// than on the currently applied filter, so these are shared between the
// toolbar and sidebar saved filter lists
function useSavedFilterOperations(props: {
  savedFilters: SavedFilterDataFragment[];
  view?: View;
  setSaving: (saving: boolean) => void;
  refetch: () => void;
}) {
  const { savedFilters, view, setSaving, refetch } = props;

  const Toast = useToast();
  const intl = useIntl();

  const renameFilter = useRenameSavedFilter();
  const [saveUISetting] = useConfigureUISetting();

  // resolves false if the rename was rejected, so that the caller can leave
  // the name being edited in place rather than discarding it
  async function renameSavedFilter(f: SavedFilterDataFragment, name: string) {
    // saved filter names are unique per mode. The unique index is
    // case-sensitive, but SaveFilterDialog matches names case-insensitively,
    // so do the same here rather than allowing a confusing "foo"/"Foo" pair.
    const existing = savedFilters.find(
      (o) => o.id !== f.id && o.name.toLowerCase() === name.toLowerCase()
    );

    if (existing) {
      Toast.error(
        intl.formatMessage(
          {
            id: "toast.filter_name_exists",
          },
          {
            entityName: existing.name,
          }
        )
      );
      return false;
    }

    try {
      setSaving(true);

      await renameFilter(f, name);

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
      refetch();
      return true;
    } catch (err) {
      Toast.error(err);
      return false;
    } finally {
      setSaving(false);
    }
  }

  async function setSavedFilterAsDefault(f: SavedFilterDataFragment) {
    if (!view) {
      return;
    }

    try {
      setSaving(true);

      await saveUISetting({
        variables: {
          key: `defaultFilters.${view.toString()}`,
          value: savedFilterFields(f),
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

  return { renameSavedFilter, setSavedFilterAsDefault };
}

export const SavedFilterList: React.FC<ISavedFilterListProps> = ({
  filter,
  onSetFilter,
  view,
  onOperationChosen,
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
  const [renamingFilter, setRenamingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();
  const [defaultingFilter, setDefaultingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();

  const saveFilter = useSaveFilter();
  const [destroyFilter] = useSavedFilterDestroy();

  const { renameSavedFilter, setSavedFilterAsDefault } =
    useSavedFilterOperations({
      savedFilters: data?.findSavedFilters ?? [],
      view,
      setSaving,
      refetch,
    });

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

  function filterClicked(f: SavedFilterDataFragment) {
    const newFilter = filter.clone();

    newFilter.currentPage = 1;
    // #1795 - reset search term if not present in saved filter
    newFilter.searchTerm = "";
    newFilter.configureFromSavedFilter(f);
    // #1507 - reset random seed when loaded
    newFilter.randomSeed = -1;

    onSetFilter(newFilter);
    notifySavedFilterLoaded({
      filter: newFilter,
      savedFilter: f,
      source: "toolbar",
      view,
    });
  }

  // a plain render function rather than a component: a component declared in
  // here would be a new type on every render, so React would remount the rows
  // and discard the in-place rename the user is part way through
  function renderSavedFilterItem(item: SavedFilterDataFragment) {
    if (renamingFilter?.id === item.id) {
      return (
        <div className="dropdown-item-container" key={item.name}>
          <SavedFilterNameInput
            name={item.name}
            onCommit={async (name) => {
              if (await renameSavedFilter(item, name)) {
                setRenamingFilter(undefined);
              }
            }}
            onCancel={() => setRenamingFilter(undefined)}
          />
        </div>
      );
    }

    return (
      <div className="dropdown-item-container" key={item.name}>
        <Dropdown.Item onClick={() => filterClicked(item)} title={item.name}>
          <span>{item.name}</span>
        </Dropdown.Item>
        <SavedFilterOperations
          id={`toolbar-${item.id}`}
          onOperationChosen={onOperationChosen}
          onOverwrite={() => setOverwritingFilter(item)}
          onRename={() => setRenamingFilter(item)}
          onSetDefault={view ? () => setDefaultingFilter(item) : undefined}
          onDelete={() => setDeletingFilter(item)}
        />
      </div>
    );
  }

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
          .map((f) => renderSavedFilterItem(f))}
      </ul>
    );
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
      <AlertModal
        show={!!defaultingFilter}
        text={<FormattedMessage id="dialogs.set_default_filter_confirm" />}
        confirmVariant="primary"
        onConfirm={() => {
          setSavedFilterAsDefault(defaultingFilter!);
          setDefaultingFilter(undefined);
        }}
        onCancel={() => setDefaultingFilter(undefined)}
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
    </>
  );
};

interface ISavedFilterItem {
  item: SavedFilterDataFragment;
  onClick: () => void;
  onOverwrite: () => void;
  onRename: () => void;
  onRenameCommit: (name: string) => void;
  onRenameCancel: () => void;
  onSetDefault?: () => void;
  onDelete: () => void;
  renaming?: boolean;
  selected?: boolean;
}

const SavedFilterItem: React.FC<ISavedFilterItem> = ({
  item,
  onClick,
  onOverwrite,
  onRename,
  onRenameCommit,
  onRenameCancel,
  onSetDefault,
  onDelete,
  renaming = false,
  selected = false,
}) => {
  if (renaming) {
    return (
      <li className="saved-filter-item">
        <SavedFilterNameInput
          name={item.name}
          onCommit={onRenameCommit}
          onCancel={onRenameCancel}
        />
      </li>
    );
  }

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
          <SavedFilterOperations
            id={`sidebar-${item.id}`}
            onOverwrite={onOverwrite}
            onRename={onRename}
            onSetDefault={onSetDefault}
            onDelete={onDelete}
          />
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
  onOverwriteClicked: (f: SavedFilterDataFragment) => void;
  onRenameClicked: (f: SavedFilterDataFragment) => void;
  onSetDefaultClicked?: (f: SavedFilterDataFragment) => void;
  onDeleteClicked: (f: SavedFilterDataFragment) => void;
  onRenameCommitted: (f: SavedFilterDataFragment, name: string) => void;
  onRenameCancelled: () => void;
  renamingFilterID?: string;
  currentFilterID?: string;
}> = ({
  error,
  loading,
  saving,
  savedFilters,
  onFilterClicked,
  onOverwriteClicked,
  onRenameClicked,
  onSetDefaultClicked,
  onDeleteClicked,
  onRenameCommitted,
  onRenameCancelled,
  renamingFilterID,
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
          onOverwrite={() => onOverwriteClicked(f)}
          onRename={() => onRenameClicked(f)}
          onSetDefault={
            onSetDefaultClicked ? () => onSetDefaultClicked(f) : undefined
          }
          onDelete={() => onDeleteClicked(f)}
          onRenameCommit={(name) => onRenameCommitted(f, name)}
          onRenameCancel={onRenameCancelled}
          renaming={renamingFilterID === f.id}
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
  const [overwritingFilter, setOverwritingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();
  const [renamingFilter, setRenamingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();
  const [defaultingFilter, setDefaultingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();
  const [showSaveDialog, setShowSaveDialog] = useState(false);
  const [settingDefault, setSettingDefault] = useState(false);

  const saveFilter = useSaveFilter();
  const [destroyFilter] = useSavedFilterDestroy();
  const [saveUISetting] = useConfigureUISetting();

  const { renameSavedFilter, setSavedFilterAsDefault } =
    useSavedFilterOperations({
      savedFilters: data?.findSavedFilters ?? [],
      view,
      setSaving,
      refetch,
    });

  const filteredFilters = useMemo(() => {
    const savedFilters = data?.findSavedFilters ?? [];
    if (!filterName) return savedFilters;

    return savedFilters.filter(
      (f) =>
        !filterName || f.name.toLowerCase().includes(filterName.toLowerCase())
    );
  }, [data?.findSavedFilters, filterName]);

  // handle when filter is changed to de-select the current filter
  // biome-ignore lint/correctness/useExhaustiveDependencies: intentionally only want to trigger when filter changes
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
    notifySavedFilterLoaded({
      filter: newFilter,
      savedFilter: f,
      source: "sidebar",
      view,
    });
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
      <OverwriteAlert
        overwritingFilter={overwritingFilter}
        onClose={(confirm) => {
          if (confirm) {
            onSaveFilter(overwritingFilter!.name, overwritingFilter!.id);
          }
          setOverwritingFilter(undefined);
        }}
      />
      <AlertModal
        show={!!defaultingFilter}
        text={<FormattedMessage id="dialogs.set_default_filter_confirm" />}
        confirmVariant="primary"
        onConfirm={() => {
          setSavedFilterAsDefault(defaultingFilter!);
          setDefaultingFilter(undefined);
        }}
        onCancel={() => setDefaultingFilter(undefined)}
      />
      {showSaveDialog && (
        <SaveFilterDialog
          mode={filter.mode}
          onClose={(name, id) => {
            setShowSaveDialog(false);
            if (name) {
              // filterName is the list search box here, not a name field
              setFilterName("");
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
        onOverwriteClicked={setOverwritingFilter}
        onRenameClicked={setRenamingFilter}
        onSetDefaultClicked={view ? setDefaultingFilter : undefined}
        onDeleteClicked={setDeletingFilter}
        onRenameCommitted={async (f, name) => {
          if (await renameSavedFilter(f, name)) {
            setRenamingFilter(undefined);
          }
        }}
        onRenameCancelled={() => setRenamingFilter(undefined)}
        renamingFilterID={renamingFilter?.id}
        currentFilterID={currentSavedFilter?.id}
      />
    </div>
  );
};

// Declared out here so that its type is stable. Built during render, every
// re-render would be a new component type, remounting the list and discarding
// an in-progress rename.
const SavedFilterDropdownMenu = React.forwardRef<
  HTMLDivElement,
  HTMLAttributes<HTMLDivElement> & { listProps: ISavedFilterListProps }
>(({ style, className, listProps }, ref) => (
  <div ref={ref} style={style} className={className}>
    <SavedFilterList {...listProps} />
  </div>
));
SavedFilterDropdownMenu.displayName = "SavedFilterDropdownMenu";

export const SavedFilterDropdown: React.FC<ISavedFilterListProps> = (props) => {
  const [show, setShow] = useState(false);

  const keepOpenRef = useRef(false);

  // see ISavedFilterListProps.onOperationChosen
  const keepOpen = useCallback(() => {
    keepOpenRef.current = true;

    // released once the click has finished propagating - it asks to close more
    // than once, and every one of those should be ignored
    window.setTimeout(() => {
      keepOpenRef.current = false;
    }, 0);
  }, []);

  const menu = (
    <Dropdown.Menu
      as={SavedFilterDropdownMenu}
      className="saved-filter-list-menu"
      listProps={{ ...props, onOperationChosen: keepOpen }}
    />
  );

  return (
    <Dropdown
      as={ButtonGroup}
      className="saved-filter-dropdown"
      show={show}
      onToggle={(nextShow) => {
        if (!nextShow && keepOpenRef.current) {
          return;
        }

        setShow(nextShow);
      }}
    >
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
