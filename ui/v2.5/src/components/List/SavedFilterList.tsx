import React, { useEffect, useRef, useState } from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  FormControl,
  InputGroup,
  Modal,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import {
  useFindSavedFilters,
  useSavedFilterDestroy,
  useSaveFilter,
  useSetDefaultFilter,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SavedFilterDataFragment } from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared";
import { PersistanceLevel } from "src/hooks/ListHook";
import { Icon } from "../Shared";

interface ISavedFilterListProps {
  filter: ListFilterModel;
  onSetFilter: (f: ListFilterModel) => void;
  persistState?: PersistanceLevel;
}

export const SavedFilterList: React.FC<ISavedFilterListProps> = ({
  filter,
  onSetFilter,
  persistState,
}) => {
  const Toast = useToast();
  const { data, error, loading, refetch } = useFindSavedFilters(filter.mode);
  const oldError = useRef(error);

  const [filterName, setFilterName] = useState("");
  const [saving, setSaving] = useState(false);
  const [deletingFilter, setDeletingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();
  const [overwritingFilter, setOverwritingFilter] = useState<
    SavedFilterDataFragment | undefined
  >();

  const [saveFilter] = useSaveFilter();
  const [destroyFilter] = useSavedFilterDestroy();
  const [setDefaultFilter] = useSetDefaultFilter();

  const savedFilters = data?.findSavedFilters ?? [];

  useEffect(() => {
    if (error && error !== oldError.current) {
      Toast.error(error);
    }

    oldError.current = error;
  }, [error, Toast, oldError]);

  async function onSaveFilter(name: string, id?: string) {
    const filterCopy = filter.clone();
    filterCopy.currentPage = 1;

    try {
      setSaving(true);
      await saveFilter({
        variables: {
          input: {
            id,
            mode: filter.mode,
            name,
            filter: JSON.stringify(filterCopy.getSavedQueryParameters()),
          },
        },
      });

      Toast.success({ content: "Filter saved" });
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

      Toast.success({ content: "Filter deleted" });
      refetch();
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
      setDeletingFilter(undefined);
    }
  }

  async function onSetDefaultFilter() {
    const filterCopy = filter.clone();
    filterCopy.currentPage = 1;

    try {
      setSaving(true);

      await setDefaultFilter({
        variables: {
          input: {
            mode: filter.mode,
            filter: JSON.stringify(filterCopy.getSavedQueryParameters()),
          },
        },
      });

      Toast.success({ content: "Default filter set" });
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
    }
  }

  function filterClicked(f: SavedFilterDataFragment) {
    const newFilter = filter.clone();
    newFilter.currentPage = 1;
    newFilter.configureFromQueryParameters(JSON.parse(f.filter));

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
            title="Overwrite"
            onClick={(e) => {
              setOverwritingFilter(item);
              e.stopPropagation();
            }}
          >
            <Icon icon="save" />
          </Button>
          <Button
            className="delete-button"
            variant="secondary"
            size="sm"
            title="Delete"
            onClick={(e) => {
              setDeletingFilter(item);
              e.stopPropagation();
            }}
          >
            <Icon icon="times" />
          </Button>
        </ButtonGroup>
      </div>
    );
  };

  function maybeRenderDeleteAlert() {
    if (!deletingFilter) {
      return;
    }

    return (
      <Modal show>
        <Modal.Body>
          Are you sure you want to delete saved query &quot;
          {deletingFilter.name}&quot;?
        </Modal.Body>
        <Modal.Footer>
          <Button
            variant="danger"
            onClick={() => onDeleteFilter(deletingFilter)}
          >
            Delete
          </Button>
          <Button
            variant="secondary"
            onClick={() => setDeletingFilter(undefined)}
          >
            Cancel
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }

  function maybeRenderOverwriteAlert() {
    if (!overwritingFilter) {
      return;
    }

    return (
      <Modal show>
        <Modal.Body>
          Are you sure you want to overwrite existing saved query &quot;
          {overwritingFilter.name}&quot;?
        </Modal.Body>
        <Modal.Footer>
          <Button
            variant="primary"
            onClick={() =>
              onSaveFilter(overwritingFilter.name, overwritingFilter.id)
            }
          >
            Overwrite
          </Button>
          <Button
            variant="secondary"
            onClick={() => setOverwritingFilter(undefined)}
          >
            Cancel
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }

  function renderSavedFilters() {
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
            (f) => !filterName || f.name.toLowerCase().includes(filterName)
          )
          .map((f) => (
            <SavedFilterItem key={f.name} item={f} />
          ))}
      </ul>
    );
  }

  function maybeRenderSetDefaultButton() {
    if (persistState === PersistanceLevel.ALL) {
      return (
        <Button
          className="set-as-default-button"
          variant="secondary"
          size="sm"
          onClick={() => onSetDefaultFilter()}
        >
          Set as default
        </Button>
      );
    }
  }

  return (
    <div>
      {maybeRenderDeleteAlert()}
      {maybeRenderOverwriteAlert()}
      <InputGroup>
        <FormControl
          className="bg-secondary text-white border-secondary"
          placeholder="Filter name..."
          value={filterName}
          onChange={(e) => setFilterName(e.target.value)}
        />
        <InputGroup.Append>
          <OverlayTrigger
            placement="top"
            overlay={<Tooltip id="filter-tooltip">Save filter</Tooltip>}
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
              <Icon icon="save" />
            </Button>
          </OverlayTrigger>
        </InputGroup.Append>
      </InputGroup>
      {renderSavedFilters()}
      {maybeRenderSetDefaultButton()}
    </div>
  );
};
