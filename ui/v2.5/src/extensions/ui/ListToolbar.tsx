/**
 * ListToolbar Component
 *
 * Custom toolbar for filtered lists with filter tags and selection support.
 *
 * Components:
 * - ToolbarFilterSection: Sidebar toggle + filter tags + filter button
 * - ToolbarSelectionSection: Selection info + actions
 * - FilteredListToolbar2: Main toolbar combining the above
 */

import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { ListFilterModel } from "src/models/list-filter/filter";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { FilterTags } from "src/extensions/ui/FilterTags";
import cx from "classnames";
import { Button, ButtonToolbar } from "react-bootstrap";
import { FilterButton } from "src/components/List/Filters/FilterButton";
import { Icon } from "src/components/Shared/Icon";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { SidebarToggleButton } from "src/components/Shared/Sidebar";
import { PatchComponent } from "src/patch";

export const ToolbarFilterSection: React.FC<{
  filter: ListFilterModel;
  onToggleSidebar: () => void;
  onSetFilter: (filter: ListFilterModel) => void;
  onEditCriterion: (c?: Criterion) => void;
  onRemoveCriterion: (criterion: Criterion, valueIndex?: number) => void;
  onRemoveAllCriterion: () => void;
  onEditSearchTerm: () => void;
  onRemoveSearchTerm: () => void;
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
  }) => {
    const { criteria, searchTerm } = filter;

    return (
      <>
        <div className="my filter-section">
          <SidebarToggleButton onClick={onToggleSidebar} />
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
          <FilterButton
            onClick={() => onEditCriterion()}
            count={criteria.length}
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
          <span>{selected} selected</span>
        </div>
        {operations}
        <div className="empty-space" />
        <Button variant="link" onClick={() => onSelectAll()}>
          <FormattedMessage id="actions.select_all" />
        </Button>
        <Button
          variant="secondary"
          className="minimal"
          onClick={() => onSelectNone()}
          title={intl.formatMessage({ id: "actions.select_none" })}
        >
          <Icon icon={faTimes} />
        </Button>
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
      <div className="filtered-list-toolbar-operations">{operationSection}</div>
    </ButtonToolbar>
  );
};
