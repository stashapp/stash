import _, { debounce } from "lodash";
import React, { HTMLAttributes, useEffect } from "react";
import Mousetrap from "mousetrap";
import { SortDirectionEnum } from "src/core/generated-graphql";
import {
  Button,
  ButtonGroup,
  Dropdown,
  Form,
  OverlayTrigger,
  Tooltip,
  InputGroup,
  FormControl,
} from "react-bootstrap";

import { Icon } from "src/components/Shared";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useFocus } from "src/utils";
import { ListFilterOptions } from "src/models/list-filter/filter-options";
import { useIntl } from "react-intl";
import { SavedFilterList } from "./SavedFilterList";

interface IListFilterProps {
  onFilterUpdate: (newFilter: ListFilterModel) => void;
  filter: ListFilterModel;
  filterOptions: ListFilterOptions;
  filterDialogOpen?: boolean;
  openFilterDialog: () => void;
}

const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120", "250", "500", "1000"];

export const ListFilter: React.FC<IListFilterProps> = ({
  onFilterUpdate,
  filter,
  filterOptions,
  filterDialogOpen,
  openFilterDialog,
}) => {
  const [queryRef, setQueryFocus] = useFocus();

  const searchCallback = debounce((value: string) => {
    const newFilter = _.cloneDeep(filter);
    newFilter.searchTerm = value;
    newFilter.currentPage = 1;
    onFilterUpdate(newFilter);
  }, 500);

  const intl = useIntl();

  useEffect(() => {
    Mousetrap.bind("/", (e) => {
      setQueryFocus();
      e.preventDefault();
    });

    Mousetrap.bind("r", () => onReshuffleRandomSort());

    return () => {
      Mousetrap.unbind("/");
      Mousetrap.unbind("r");
    };
  });

  function onChangePageSize(event: React.ChangeEvent<HTMLSelectElement>) {
    const val = event.currentTarget.value;

    const newFilter = _.cloneDeep(filter);
    newFilter.itemsPerPage = parseInt(val, 10);
    newFilter.currentPage = 1;
    onFilterUpdate(newFilter);
  }

  function onChangeQuery(event: React.FormEvent<HTMLInputElement>) {
    searchCallback(event.currentTarget.value);
  }

  function onChangeSortDirection() {
    const newFilter = _.cloneDeep(filter);
    if (filter.sortDirection === SortDirectionEnum.Asc) {
      newFilter.sortDirection = SortDirectionEnum.Desc;
    } else {
      newFilter.sortDirection = SortDirectionEnum.Asc;
    }

    onFilterUpdate(newFilter);
  }

  function onChangeSortBy(eventKey: string | null) {
    const newFilter = _.cloneDeep(filter);
    newFilter.sortBy = eventKey ?? undefined;
    newFilter.currentPage = 1;
    onFilterUpdate(newFilter);
  }

  function onReshuffleRandomSort() {
    const newFilter = _.cloneDeep(filter);
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
    onFilterUpdate(newFilter);
  }

  function renderSortByOptions() {
    return filterOptions.sortByOptions
      .map((o) => {
        return {
          message: intl.formatMessage({ id: o.messageID }),
          value: o.value,
        };
      })
      .sort((a, b) => a.message.localeCompare(b.message))
      .map((option) => (
        <Dropdown.Item
          onSelect={onChangeSortBy}
          key={option.value}
          className="bg-secondary text-white"
          eventKey={option.value}
        >
          {option.message}
        </Dropdown.Item>
      ));
  }

  const SavedFilterDropdown = React.forwardRef<
    HTMLDivElement,
    HTMLAttributes<HTMLDivElement>
  >(({ style, className }, ref) => (
    <div ref={ref} style={style} className={className}>
      <SavedFilterList
        filter={filter}
        onSetFilter={(f) => {
          onFilterUpdate(f);
        }}
      />
    </div>
  ));

  function render() {
    const currentSortBy = filterOptions.sortByOptions.find(
      (o) => o.value === filter.sortBy
    );

    return (
      <>
        <div className="d-flex">
          <InputGroup className="mr-2 flex-grow-1">
            <InputGroup.Prepend>
              <Dropdown>
                <OverlayTrigger
                  placement="top"
                  overlay={<Tooltip id="filter-tooltip">Saved filters</Tooltip>}
                >
                  <Dropdown.Toggle variant="secondary">
                    <Icon icon="bookmark" />
                  </Dropdown.Toggle>
                </OverlayTrigger>
                <Dropdown.Menu
                  as={SavedFilterDropdown}
                  className="saved-filter-list-menu"
                />
              </Dropdown>
            </InputGroup.Prepend>
            <FormControl
              ref={queryRef}
              placeholder="Search..."
              defaultValue={filter.searchTerm}
              onInput={onChangeQuery}
              className="query-text-field bg-secondary text-white border-secondary"
            />

            <InputGroup.Append>
              <OverlayTrigger
                placement="top"
                overlay={<Tooltip id="filter-tooltip">Filter</Tooltip>}
              >
                <Button
                  variant="secondary"
                  onClick={() => openFilterDialog()}
                  active={filterDialogOpen}
                >
                  <Icon icon="filter" />
                </Button>
              </OverlayTrigger>
            </InputGroup.Append>
          </InputGroup>

          <Dropdown as={ButtonGroup} className="mr-2">
            <Dropdown.Toggle variant="secondary">
              {currentSortBy
                ? intl.formatMessage({ id: currentSortBy.messageID })
                : ""}
            </Dropdown.Toggle>
            <Dropdown.Menu className="bg-secondary text-white">
              {renderSortByOptions()}
            </Dropdown.Menu>
            <OverlayTrigger
              overlay={
                <Tooltip id="sort-direction-tooltip">
                  {filter.sortDirection === SortDirectionEnum.Asc
                    ? "Ascending"
                    : "Descending"}
                </Tooltip>
              }
            >
              <Button variant="secondary" onClick={onChangeSortDirection}>
                <Icon
                  icon={
                    filter.sortDirection === SortDirectionEnum.Asc
                      ? "caret-up"
                      : "caret-down"
                  }
                />
              </Button>
            </OverlayTrigger>
            {filter.sortBy === "random" && (
              <OverlayTrigger
                overlay={
                  <Tooltip id="sort-reshuffle-tooltip">Reshuffle</Tooltip>
                }
              >
                <Button variant="secondary" onClick={onReshuffleRandomSort}>
                  <Icon icon="random" />
                </Button>
              </OverlayTrigger>
            )}
          </Dropdown>
        </div>

        <Form.Control
          as="select"
          onChange={onChangePageSize}
          value={filter.itemsPerPage.toString()}
          className="btn-secondary mx-1"
        >
          {PAGE_SIZE_OPTIONS.map((s) => (
            <option value={s} key={s}>
              {s}
            </option>
          ))}
        </Form.Control>
      </>
    );
  }

  return render();
};
