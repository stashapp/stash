import React from "react";
import { QueryResult } from "@apollo/client";
import { ListFilterModel } from "src/models/list-filter/filter";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { PageSizeSelector, SearchTermInput, SortBySelect } from "./ListFilter";
import { ListViewButtonGroup } from "./ListViewOptions";
import {
  IListFilterOperation,
  ListOperationButtons,
} from "./ListOperationButtons";
import { Button, ButtonGroup, ButtonToolbar } from "react-bootstrap";
import { View } from "./views";
import { IListSelect, useFilterOperations } from "./util";
import { SavedFilterDropdown } from "./SavedFilterList";
import { FilterButton } from "./Filters/FilterButton";
import { Icon } from "../Shared/Icon";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { faSquareCheck } from "@fortawesome/free-regular-svg-icons";
import { useIntl } from "react-intl";
import cx from "classnames";

const SelectionSection: React.FC<{
  filter: ListFilterModel;
  selected: number;
  onSelectAll: () => void;
  onSelectNone: () => void;
}> = ({ selected, onSelectAll, onSelectNone }) => {
  const intl = useIntl();

  return (
    <div className="selected-items-info">
      <Button
        variant="secondary"
        className="minimal"
        onClick={() => onSelectNone()}
        title={intl.formatMessage({ id: "actions.select_none" })}
      >
        <Icon icon={faTimes} />
      </Button>
      <span className="selected-count">{selected}</span>
      <Button
        variant="secondary"
        className="minimal"
        onClick={() => onSelectAll()}
        title={intl.formatMessage({ id: "actions.select_all" })}
      >
        <Icon icon={faSquareCheck} />
      </Button>
    </div>
  );
};

export interface IItemListOperation<T extends QueryResult> {
  text: string;
  onClick: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => Promise<void>;
  isDisplayed?: (
    result: T,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) => boolean;
  postRefetch?: boolean;
  icon?: IconDefinition;
  buttonVariant?: string;
}

export interface IFilteredListToolbar {
  filter: ListFilterModel;
  setFilter: (
    value: ListFilterModel | ((prevState: ListFilterModel) => ListFilterModel)
  ) => void;
  showEditFilter: () => void;
  view?: View;
  listSelect: IListSelect;
  onEdit?: () => void;
  onDelete?: () => void;
  operations?: IListFilterOperation[];
  operationComponent?: React.ReactNode;
  zoomable?: boolean;
}

export const FilteredListToolbar: React.FC<IFilteredListToolbar> = ({
  filter,
  setFilter,
  showEditFilter,
  view,
  listSelect,
  onEdit,
  onDelete,
  operations,
  operationComponent,
  zoomable = false,
}) => {
  const filterOptions = filter.options;
  const { setDisplayMode, setZoom } = useFilterOperations({
    filter,
    setFilter,
  });
  const { selectedIds, onSelectAll, onSelectNone } = listSelect;
  const hasSelection = selectedIds.size > 0;

  const renderOperations = operationComponent ?? (
    <ListOperationButtons
      onSelectAll={onSelectAll}
      onSelectNone={onSelectNone}
      otherOperations={operations}
      itemsSelected={selectedIds.size > 0}
      onEdit={onEdit}
      onDelete={onDelete}
    />
  );

  return (
    <ButtonToolbar
      className={cx("filtered-list-toolbar", { "has-selection": hasSelection })}
    >
      {hasSelection ? (
        <SelectionSection
          filter={filter}
          selected={selectedIds.size}
          onSelectAll={onSelectAll}
          onSelectNone={onSelectNone}
        />
      ) : (
        <>
          <SearchTermInput filter={filter} onFilterUpdate={setFilter} />

          <ButtonGroup>
            <SavedFilterDropdown
              filter={filter}
              onSetFilter={setFilter}
              view={view}
            />
            <FilterButton
              onClick={() => showEditFilter()}
              count={filter.count()}
            />
          </ButtonGroup>

          <SortBySelect
            sortBy={filter.sortBy}
            sortDirection={filter.sortDirection}
            options={filterOptions.sortByOptions}
            onChangeSortBy={(e) => setFilter(filter.setSortBy(e ?? undefined))}
            onChangeSortDirection={() =>
              setFilter(filter.toggleSortDirection())
            }
            onReshuffleRandomSort={() =>
              setFilter(filter.reshuffleRandomSort())
            }
          />

          <PageSizeSelector
            pageSize={filter.itemsPerPage}
            setPageSize={(size) => setFilter(filter.setPageSize(size))}
          />
        </>
      )}

      {renderOperations}

      <ListViewButtonGroup
        displayMode={filter.displayMode}
        displayModeOptions={filterOptions.displayModeOptions}
        onSetDisplayMode={setDisplayMode}
        zoomIndex={zoomable ? filter.zoomIndex : undefined}
        onSetZoom={zoomable ? setZoom : undefined}
      />
    </ButtonToolbar>
  );
};
