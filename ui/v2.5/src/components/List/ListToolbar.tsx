import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { ListFilterModel } from "src/models/list-filter/filter";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { FilterTags } from "src/extensions/ui";
import cx from "classnames";
import { Button, ButtonGroup, ButtonToolbar } from "react-bootstrap";
import { FilterButton } from "../List/Filters/FilterButton";
import { Icon } from "../Shared/Icon";
import { SearchTermInput } from "../List/ListFilter";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { SidebarToggleButton } from "../Shared/Sidebar";
import { PatchComponent } from "src/patch";
import { SavedFilterDropdown } from "./SavedFilterList";
import { View } from "./views";

export const ToolbarFilterSection: React.FC<{
  filter: ListFilterModel;
  onToggleSidebar: () => void;
  onSetFilter: (filter: ListFilterModel) => void;
  onEditCriterion: (c?: Criterion) => void;
  onRemoveCriterion: (criterion: Criterion, valueIndex?: number) => void;
  onRemoveAllCriterion: () => void;
  onEditSearchTerm: () => void;
  onRemoveSearchTerm: () => void;
  view?: View;
}> = PatchComponent(
  "ToolbarFilterSection",
  ({
    filter,
    onToggleSidebar,
    onSetFilter,
    onEditCriterion,
    onRemoveCriterion,
    onRemoveAllCriterion,
    onEditSearchTerm,
    onRemoveSearchTerm,
    view,
  }) => {
    const { criteria, searchTerm } = filter;

    return (
      <>
        <div className="search-container">
          <SearchTermInput filter={filter} onFilterUpdate={onSetFilter} />
        </div>
        <div className="filter-section">
          <ButtonGroup>
            <SidebarToggleButton onClick={onToggleSidebar} />
            <SavedFilterDropdown
              filter={filter}
              onSetFilter={onSetFilter}
              view={view}
              menuPortalTarget={document.body}
            />
            <FilterButton
              onClick={() => onEditCriterion()}
              count={criteria.length}
            />
          </ButtonGroup>
          <FilterTags
            searchTerm={searchTerm}
            criteria={criteria}
            onEditCriterion={onEditCriterion}
            onRemoveCriterion={onRemoveCriterion}
            onRemoveAll={onRemoveAllCriterion}
            onEditSearchTerm={onEditSearchTerm}
            onRemoveSearchTerm={onRemoveSearchTerm}
            truncateOnOverflow
          />
        </div>
      </>
    );
  }
);

export const ToolbarSelectionSection: React.FC<{
  selected: number;
  onToggleSidebar: () => void;
  operations?: React.ReactNode;
  onSelectAll: () => void;
  onSelectNone: () => void;
}> = PatchComponent(
  "ToolbarSelectionSection",
  ({ selected, onToggleSidebar, operations, onSelectAll, onSelectNone }) => {
    const intl = useIntl();

    return (
      <div className="toolbar-selection-section">
        <div className="selected-items-info">
          <SidebarToggleButton onClick={onToggleSidebar} />
          <Button
            variant="secondary"
            className="minimal"
            onClick={() => onSelectNone()}
            title={intl.formatMessage({ id: "actions.select_none" })}
          >
            <Icon icon={faTimes} />
          </Button>
          <span>{selected} selected</span>
          <Button variant="link" onClick={() => onSelectAll()}>
            <FormattedMessage id="actions.select_all" />
          </Button>
        </div>
        {operations}
        <div className="empty-space" />
      </div>
    );
  }
);

// TODO - rename to FilteredListToolbar once all list components have been updated
// TODO - and expose to plugins
export const FilteredListToolbar2: React.FC<{
  className?: string;
  hasSelection: boolean;
  filterSection: React.ReactNode;
  selectionSection: React.ReactNode;
  operationSection: React.ReactNode;
}> = ({
  className,
  hasSelection,
  filterSection,
  selectionSection,
  operationSection,
}) => {
  return (
    <ButtonToolbar
      className={cx(className, "filtered-list-toolbar", {
        "has-selection": hasSelection,
      })}
    >
      {!hasSelection ? filterSection : selectionSection}
      {!hasSelection ? (
        <div className="filtered-list-toolbar-operations">
          {operationSection}
        </div>
      ) : null}
    </ButtonToolbar>
  );
};
