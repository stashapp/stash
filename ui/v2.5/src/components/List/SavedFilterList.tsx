import React, { HTMLAttributes, useState } from "react";
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
  useConfigureUI,
  useFindSavedFilters,
  useSavedFilterDestroy,
  useSaveFilter,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SavedFilterDataFragment } from "src/core/generated-graphql";
import { View } from "./views";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { faBookmark, faSave, faTimes } from "@fortawesome/free-solid-svg-icons";

interface ISavedFilterListProps {
  filter: ListFilterModel;
  onSetFilter: (f: ListFilterModel) => void;
  view?: View;
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

  const [saveFilter] = useSaveFilter();
  const [destroyFilter] = useSavedFilterDestroy();
  const [saveUI] = useConfigureUI();

  const savedFilters = data?.findSavedFilters ?? [];

  async function onSaveFilter(name: string, id?: string) {
    const filterCopy = filter.clone();

    try {
      setSaving(true);
      await saveFilter({
        variables: {
          input: {
            id,
            mode: filter.mode,
            name,
            find_filter: filterCopy.makeFindFilter(),
            object_filter: filterCopy.makeSavedFilter(),
            ui_options: filterCopy.makeSavedUIOptions(),
          },
        },
      });

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

      await saveUI({
        variables: {
          partial: {
            defaultFilters: {
              [view.toString()]: {
                mode: filter.mode,
                find_filter: filterCopy.makeFindFilter(),
                object_filter: filterCopy.makeSavedFilter(),
                ui_options: filterCopy.makeSavedUIOptions(),
              },
            },
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

  function maybeRenderDeleteAlert() {
    if (!deletingFilter) {
      return;
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
          <Button
            variant="danger"
            onClick={() => onDeleteFilter(deletingFilter)}
          >
            {intl.formatMessage({ id: "actions.delete" })}
          </Button>
          <Button
            variant="secondary"
            onClick={() => setDeletingFilter(undefined)}
          >
            {intl.formatMessage({ id: "actions.cancel" })}
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
          <FormattedMessage
            id="dialogs.overwrite_filter_confirm"
            values={{
              entityName: overwritingFilter.name,
            }}
          />
        </Modal.Body>
        <Modal.Footer>
          <Button
            variant="primary"
            onClick={() =>
              onSaveFilter(overwritingFilter.name, overwritingFilter.id)
            }
          >
            {intl.formatMessage({ id: "actions.overwrite" })}
          </Button>
          <Button
            variant="secondary"
            onClick={() => setOverwritingFilter(undefined)}
          >
            {intl.formatMessage({ id: "actions.cancel" })}
          </Button>
        </Modal.Footer>
      </Modal>
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
      {maybeRenderDeleteAlert()}
      {maybeRenderOverwriteAlert()}
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

  return (
    <Dropdown>
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
      <Dropdown.Menu
        as={SavedFilterDropdownRef}
        className="saved-filter-list-menu"
      />
    </Dropdown>
  );
};
