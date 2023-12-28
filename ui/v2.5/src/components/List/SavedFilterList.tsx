import React, { useContext, useEffect, useRef, useState } from "react";
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
import { faSave, faTimes } from "@fortawesome/free-solid-svg-icons";
import { DefaultFilters, IUIConfig } from "src/core/config";
import { ConfigurationContext } from "src/hooks/Config";
import { IHierarchicalLabelValue } from "src/models/list-filter/types";

interface ISavedFilterListProps {
  filter: ListFilterModel;
  onSetFilter: (f: ListFilterModel) => void;
  view?: View;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const SavedFilterList: React.FC<ISavedFilterListProps> = ({
  filter,
  onSetFilter,
  view,
  filterHook,
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
  const { configuration } = useContext(ConfigurationContext);
  const ui = (configuration?.ui ?? {}) as IUIConfig;
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
            object_filter: filterCopy.makeSavedFindFilter(),
            ui_options: filterCopy.makeUIOptions(),
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

      // TODO - this is a horrible temporary hack to work around stupid viper
      let existingDefaultFilters: DefaultFilters;
      try {
        existingDefaultFilters = JSON.parse(ui.defaultFilters ?? "{}");
      } catch (e) {
        // ignore
        existingDefaultFilters = {};
      }

      let objectFilter = filterCopy.makeSavedFindFilter();
      // Remove Subview Filter from default Filter
      if (filterHook) {
        let subViewFilter = filterHook(
          new ListFilterModel(filterCopy.mode, undefined)
        );
        subViewFilter.criteria.forEach((criterion) => {
          let subViewCriterionValue =
            criterion.value as IHierarchicalLabelValue;
          let subViewType = criterion.criterionOption.type;
          if (
            Object.keys(objectFilter).indexOf(subViewType) > -1 &&
            objectFilter[subViewType].modifier === criterion.modifier
          ) {
            let value = objectFilter[subViewType]
              .value as IHierarchicalLabelValue;
            value.items = value.items.filter((item) => {
              let ret = true;
              subViewCriterionValue.items.forEach((subViewItem) => {
                if (ret) {
                  ret = item.id != subViewItem.id;
                }
              });
              return ret;
            });
            value.excluded = value.excluded.filter((excluded) => {
              let ret = true;
              subViewCriterionValue.excluded.forEach((subViewExcluded) => {
                if (ret) {
                  ret = excluded.id != subViewExcluded.id;
                }
              });
              return ret;
            });
            objectFilter[subViewType].value = value;
            if (value.items.length === 0 && value.excluded.length === 0) {
              delete objectFilter[subViewType];
            }
          }
        });
      }

      const newDefaultFilters = JSON.stringify({
        ...existingDefaultFilters,
        [view.toString()]: {
          mode: filter.mode,
          find_filter: filter.makeFindFilter(),
          object_filter: objectFilter,
          ui_options: filterCopy.makeUIOptions(),
        },
      });

      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            defaultFilters: newDefaultFilters,
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
          <Button
            className="set-as-default-button"
            variant="secondary"
            size="sm"
            onClick={() => onSetDefaultFilter()}
          >
            {intl.formatMessage({ id: "actions.set_as_default" })}
          </Button>
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
