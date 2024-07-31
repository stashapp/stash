import React from "react";
import { QueryResult } from "@apollo/client";
import { ListFilterModel } from "src/models/list-filter/filter";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { ListFilter } from "./ListFilter";
import { ListViewOptions } from "./ListViewOptions";
import {
  IListFilterOperation,
  ListOperationButtons,
} from "./ListOperationButtons";
import { DisplayMode } from "src/models/list-filter/types";
import { ButtonToolbar } from "react-bootstrap";
import { View } from "./views";
import { useListContext } from "./ListProvider";
import { useFilter } from "./FilterProvider";

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

export const FilteredListToolbar: React.FC<{
  showEditFilter: (editingCriterion?: string) => void;
  view?: View;
  onEdit?: () => void;
  onDelete?: () => void;
  operations?: IListFilterOperation[];
  zoomable?: boolean;
}> = ({
  showEditFilter,
  view,
  onEdit,
  onDelete,
  operations,
  zoomable = false,
}) => {
  const { getSelected, onSelectAll, onSelectNone } = useListContext();
  const { filter, setFilter } = useFilter();

  const filterOptions = filter.options;

  function onChangeDisplayMode(displayMode: DisplayMode) {
    setFilter(filter.setDisplayMode(displayMode));
  }

  function onChangeZoom(newZoomIndex: number) {
    setFilter(filter.setZoom(newZoomIndex));
  }

  return (
    <ButtonToolbar className="justify-content-center">
      <ListFilter
        onFilterUpdate={setFilter}
        filter={filter}
        openFilterDialog={() => showEditFilter()}
        view={view}
      />
      <ListOperationButtons
        onSelectAll={onSelectAll}
        onSelectNone={onSelectNone}
        otherOperations={operations}
        itemsSelected={getSelected().length > 0}
        onEdit={onEdit}
        onDelete={onDelete}
      />
      <ListViewOptions
        displayMode={filter.displayMode}
        displayModeOptions={filterOptions.displayModeOptions}
        onSetDisplayMode={onChangeDisplayMode}
        zoomIndex={zoomable ? filter.zoomIndex : undefined}
        onSetZoom={zoomable ? onChangeZoom : undefined}
      />
    </ButtonToolbar>
  );
};
