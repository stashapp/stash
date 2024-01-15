import React, { useCallback } from "react";
import { Pagination } from "../List/Pagination";
import { DisplayModeSelect } from "../List/ListViewOptions";
import { DisplayMode } from "src/models/list-filter/types";
import { PageSizeSelect, SortBySelect } from "../List/ListFilter";
import { SortDirectionEnum } from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { FormattedMessage } from "react-intl";
import { faTimes } from "@fortawesome/free-solid-svg-icons";

interface IDefaultListHeaderProps {
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  totalItems: number;
  actionButtons?: React.ReactNode;
}

const DefaultListHeader: React.FC<IDefaultListHeaderProps> = ({
  filter,
  setFilter,
  totalItems,
  actionButtons,
}) => {
  const filterOptions = getFilterOptions(filter.mode);

  // function onChangeZoom(newZoomIndex: number) {
  //   const newFilter = filter.clone();
  //   newFilter.zoomIndex = newZoomIndex;
  //   setFilter(newFilter);
  // }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = filter.clone();
    newFilter.displayMode = displayMode;
    setFilter(newFilter);
  }

  function onChangePageSize(val: number) {
    const newFilter = filter.clone();
    newFilter.itemsPerPage = val;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  function onChangeSortDirection(dir: SortDirectionEnum) {
    const newFilter = filter.clone();
    newFilter.sortDirection = dir;
    setFilter(newFilter);
  }

  function onChangeSortBy(eventKey: string | null) {
    const newFilter = filter.clone();
    newFilter.sortBy = eventKey ?? undefined;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  function onReshuffleRandomSort() {
    const newFilter = filter.clone();
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
    setFilter(newFilter);
  }

  const onChangePage = useCallback(
    (page: number) => {
      const newFilter = filter.clone();
      newFilter.currentPage = page;
      setFilter(newFilter);

      // if the current page has a detail-header, then
      // scroll up relative to that rather than 0, 0
      const detailHeader = document.querySelector(".detail-header");
      if (detailHeader) {
        window.scrollTo(0, detailHeader.scrollHeight - 50);
      } else {
        window.scrollTo(0, 0);
      }
    },
    [filter, setFilter]
  );

  return (
    <div className="list-header">
      <div>{actionButtons}</div>
      <div>
        <Pagination
          currentPage={filter.currentPage}
          itemsPerPage={filter.itemsPerPage}
          totalItems={totalItems}
          onChangePage={onChangePage}
        />
      </div>
      <div>
        <SortBySelect
          sortBy={filter.sortBy}
          direction={filter.sortDirection}
          options={filterOptions.sortByOptions}
          setSortBy={onChangeSortBy}
          setDirection={onChangeSortDirection}
          onReshuffleRandomSort={onReshuffleRandomSort}
        />
        <PageSizeSelect
          pageSize={filter.itemsPerPage}
          setPageSize={onChangePageSize}
        />
        {/* <div>
          <ZoomSelect
            minZoom={0}
            maxZoom={3}
            zoomIndex={filter.zoomIndex}
            onChangeZoom={onChangeZoom}
          />
        </div> */}
        <DisplayModeSelect
          displayMode={filter.displayMode}
          displayModeOptions={filterOptions.displayModeOptions}
          onSetDisplayMode={onChangeDisplayMode}
        />
      </div>
    </div>
  );
};

interface ISelectedListHeader {
  selectedIds: Set<string>;
  onSelectAll: () => void;
  onSelectNone: () => void;
  selectedButtons?: (selectedIds: Set<string>) => React.ReactNode;
}

export const SelectedListHeader: React.FC<ISelectedListHeader> = ({
  selectedIds,
  onSelectAll,
  onSelectNone,
  selectedButtons = () => null,
}) => {
  return (
    <div className="list-header">
      <div>
        <span>{selectedIds.size} items selected</span>
        <Button variant="link" onClick={() => onSelectAll()}>
          <FormattedMessage id="actions.select_all" />
        </Button>
      </div>
      {selectedButtons(selectedIds)}
      <div>
        <Button className="minimal select-none" onClick={() => onSelectNone()}>
          <Icon icon={faTimes} />
        </Button>
      </div>
    </div>
  );
};

export interface IListHeaderProps
  extends IDefaultListHeaderProps,
    ISelectedListHeader {}

export const ListHeader: React.FC<IListHeaderProps> = (props) => {
  if (props.selectedIds.size === 0) {
    return <DefaultListHeader {...props} />;
  } else {
    return <SelectedListHeader {...props} />;
  }
};
