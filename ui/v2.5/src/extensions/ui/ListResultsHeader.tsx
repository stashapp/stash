/**
 * ListResultsHeader Component
 *
 * Results header showing pagination, sort options, and view controls.
 *
 * Features:
 * - Pagination index display
 * - Sort by selector
 * - Page size selector
 * - View mode options (grid/list/tagger)
 */

import React from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import { PaginationIndex } from "src/components/List/Pagination";
import { ButtonToolbar } from "react-bootstrap";
import { ListViewOptions } from "src/components/List/ListViewOptions";
import { PageSizeSelector, SortBySelect } from "src/components/List/ListFilter";
import cx from "classnames";

export const ListResultsHeader: React.FC<{
  className?: string;
  loading: boolean;
  filter: ListFilterModel;
  totalCount: number;
  metadataByline?: React.ReactNode;
  onChangeFilter: (filter: ListFilterModel) => void;
}> = ({
  className,
  loading,
  filter,
  totalCount,
  metadataByline,
  onChangeFilter,
}) => {
  return (
    <ButtonToolbar className={cx(className, "list-results-header")}>
      <div>
        <PaginationIndex
          loading={loading}
          itemsPerPage={filter.itemsPerPage}
          currentPage={filter.currentPage}
          totalItems={totalCount}
          metadataByline={metadataByline}
        />
      </div>
      <div>
        <SortBySelect
          options={filter.options.sortByOptions}
          sortBy={filter.sortBy}
          sortDirection={filter.sortDirection}
          onChangeSortBy={(s) =>
            onChangeFilter(filter.setSortBy(s ?? undefined))
          }
          onChangeSortDirection={() =>
            onChangeFilter(filter.toggleSortDirection())
          }
          onReshuffleRandomSort={() =>
            onChangeFilter(filter.reshuffleRandomSort())
          }
        />
        <PageSizeSelector
          pageSize={filter.itemsPerPage}
          setPageSize={(s) => onChangeFilter(filter.setPageSize(s))}
        />
        <ListViewOptions
          displayMode={filter.displayMode}
          zoomIndex={filter.zoomIndex}
          displayModeOptions={filter.options.displayModeOptions}
          onSetDisplayMode={(mode) =>
            onChangeFilter(filter.setDisplayMode(mode))
          }
          onSetZoom={(zoom) => onChangeFilter(filter.setZoom(zoom))}
        />
      </div>
    </ButtonToolbar>
  );
};
