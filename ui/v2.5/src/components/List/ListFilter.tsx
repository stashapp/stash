import cloneDeep from "lodash-es/cloneDeep";
import React, { useCallback, useEffect, useRef, useState } from "react";
import cx from "classnames";
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
  Popover,
  Overlay,
} from "react-bootstrap";

import { Icon } from "../Shared/Icon";
import { ListFilterModel } from "src/models/list-filter/filter";
import useFocus from "src/utils/focus";
import { FormattedMessage, useIntl } from "react-intl";
import { SavedFilterDropdown } from "./SavedFilterList";
import {
  faCaretDown,
  faCaretUp,
  faCheck,
  faRandom,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { FilterButton } from "./Filters/FilterButton";
import { useDebounce } from "src/hooks/debounce";
import { View } from "./views";

interface IListFilterProps {
  onFilterUpdate: (newFilter: ListFilterModel) => void;
  filter: ListFilterModel;
  view?: View;
  openFilterDialog: () => void;
}

const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120", "250", "500", "1000"];

export const ListFilter: React.FC<IListFilterProps> = ({
  onFilterUpdate,
  filter,
  openFilterDialog,
  view,
}) => {
  const [customPageSizeShowing, setCustomPageSizeShowing] = useState(false);
  const [queryRef, setQueryFocus] = useFocus();
  const [queryClearShowing, setQueryClearShowing] = useState(
    !!filter.searchTerm
  );
  const perPageSelect = useRef(null);
  const [perPageInput, perPageFocus] = useFocus();

  const filterOptions = filter.options;

  const searchQueryUpdated = useCallback(
    (value: string) => {
      const newFilter = cloneDeep(filter);
      newFilter.searchTerm = value;
      newFilter.currentPage = 1;
      onFilterUpdate(newFilter);
    },
    [filter, onFilterUpdate]
  );

  const searchCallback = useDebounce((value: string) => {
    const newFilter = cloneDeep(filter);
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

  useEffect(() => {
    if (customPageSizeShowing) {
      perPageFocus();
    }
  }, [customPageSizeShowing, perPageFocus]);

  // clear search input when filter is cleared
  useEffect(() => {
    if (!filter.searchTerm) {
      if (queryRef.current) queryRef.current.value = "";
      setQueryClearShowing(false);
    }
  }, [filter.searchTerm, queryRef]);

  function onChangePageSize(val: string) {
    if (val === "custom") {
      // added timeout since Firefox seems to trigger the rootClose immediately
      // without it
      setTimeout(() => setCustomPageSizeShowing(true), 0);
      return;
    }

    setCustomPageSizeShowing(false);

    let pp = parseInt(val, 10);
    if (Number.isNaN(pp) || pp <= 0) {
      return;
    }

    const newFilter = cloneDeep(filter);
    newFilter.itemsPerPage = pp;
    newFilter.currentPage = 1;
    onFilterUpdate(newFilter);
  }

  function onChangeQuery(event: React.FormEvent<HTMLInputElement>) {
    searchCallback(event.currentTarget.value);
    setQueryClearShowing(!!event.currentTarget.value);
  }

  function onClearQuery() {
    if (queryRef.current) queryRef.current.value = "";
    searchQueryUpdated("");
    setQueryFocus();
    setQueryClearShowing(false);
  }

  function onChangeSortDirection() {
    const newFilter = cloneDeep(filter);
    if (filter.sortDirection === SortDirectionEnum.Asc) {
      newFilter.sortDirection = SortDirectionEnum.Desc;
    } else {
      newFilter.sortDirection = SortDirectionEnum.Asc;
    }

    onFilterUpdate(newFilter);
  }

  function onChangeSortBy(eventKey: string | null) {
    const newFilter = cloneDeep(filter);
    newFilter.sortBy = eventKey ?? undefined;
    newFilter.currentPage = 1;
    onFilterUpdate(newFilter);
  }

  function onReshuffleRandomSort() {
    const newFilter = cloneDeep(filter);
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

  function render() {
    const currentSortBy = filterOptions.sortByOptions.find(
      (o) => o.value === filter.sortBy
    );

    const pageSizeOptions = PAGE_SIZE_OPTIONS.map((o) => {
      return {
        label: o,
        value: o,
      };
    });
    const currentPerPage = filter.itemsPerPage.toString();
    if (!pageSizeOptions.find((o) => o.value === currentPerPage)) {
      pageSizeOptions.push({ label: currentPerPage, value: currentPerPage });
      pageSizeOptions.sort(
        (a, b) => parseInt(a.value, 10) - parseInt(b.value, 10)
      );
    }

    pageSizeOptions.push({
      label: `${intl.formatMessage({ id: "custom" })}...`,
      value: "custom",
    });

    return (
      <>
        <div className="mb-2 mr-2 d-flex">
          <div className="flex-grow-1 query-text-field-group">
            <FormControl
              ref={queryRef}
              placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
              defaultValue={filter.searchTerm}
              onInput={onChangeQuery}
              className="query-text-field bg-secondary text-white border-secondary"
            />
            <Button
              variant="secondary"
              onClick={onClearQuery}
              title={intl.formatMessage({ id: "actions.clear" })}
              className={cx(
                "query-text-field-clear",
                queryClearShowing ? "" : "d-none"
              )}
            >
              <Icon icon={faTimes} />
            </Button>
          </div>
        </div>

        <ButtonGroup className="mr-2 mb-2">
          <SavedFilterDropdown
            filter={filter}
            onSetFilter={(f) => {
              onFilterUpdate(f);
            }}
            view={view}
          />
          <OverlayTrigger
            placement="top"
            overlay={
              <Tooltip id="filter-tooltip">
                <FormattedMessage id="search_filter.name" />
              </Tooltip>
            }
          >
            <FilterButton onClick={() => openFilterDialog()} filter={filter} />
          </OverlayTrigger>
        </ButtonGroup>

        <Dropdown as={ButtonGroup} className="mr-2 mb-2">
          <InputGroup.Prepend>
            <Dropdown.Toggle variant="secondary">
              {currentSortBy
                ? intl.formatMessage({ id: currentSortBy.messageID })
                : ""}
            </Dropdown.Toggle>
          </InputGroup.Prepend>
          <Dropdown.Menu className="bg-secondary text-white">
            {renderSortByOptions()}
          </Dropdown.Menu>
          <OverlayTrigger
            overlay={
              <Tooltip id="sort-direction-tooltip">
                {filter.sortDirection === SortDirectionEnum.Asc
                  ? intl.formatMessage({ id: "ascending" })
                  : intl.formatMessage({ id: "descending" })}
              </Tooltip>
            }
          >
            <Button variant="secondary" onClick={onChangeSortDirection}>
              <Icon
                icon={
                  filter.sortDirection === SortDirectionEnum.Asc
                    ? faCaretUp
                    : faCaretDown
                }
              />
            </Button>
          </OverlayTrigger>
          {filter.sortBy === "random" && (
            <OverlayTrigger
              overlay={
                <Tooltip id="sort-reshuffle-tooltip">
                  {intl.formatMessage({ id: "actions.reshuffle" })}
                </Tooltip>
              }
            >
              <Button variant="secondary" onClick={onReshuffleRandomSort}>
                <Icon icon={faRandom} />
              </Button>
            </OverlayTrigger>
          )}
        </Dropdown>

        <div className="mb-2">
          <Form.Control
            as="select"
            ref={perPageSelect}
            onChange={(e) => onChangePageSize(e.target.value)}
            value={filter.itemsPerPage.toString()}
            className="btn-secondary"
          >
            {pageSizeOptions.map((s) => (
              <option value={s.value} key={s.value}>
                {s.label}
              </option>
            ))}
          </Form.Control>
          <Overlay
            target={perPageSelect.current}
            show={customPageSizeShowing}
            placement="bottom"
            rootClose
            onHide={() => setCustomPageSizeShowing(false)}
          >
            <Popover id="custom_pagesize_popover">
              <Form inline>
                <InputGroup>
                  <Form.Control
                    type="number"
                    min={1}
                    className="text-input"
                    ref={perPageInput}
                    onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) => {
                      if (e.key === "Enter") {
                        onChangePageSize(
                          (perPageInput.current as HTMLInputElement)?.value ??
                            ""
                        );
                        e.preventDefault();
                      }
                    }}
                  />
                  <InputGroup.Append>
                    <Button
                      variant="primary"
                      onClick={() =>
                        onChangePageSize(
                          (perPageInput.current as HTMLInputElement)?.value ??
                            ""
                        )
                      }
                    >
                      <Icon icon={faCheck} />
                    </Button>
                  </InputGroup.Append>
                </InputGroup>
              </Form>
            </Popover>
          </Overlay>
        </div>
      </>
    );
  }

  return render();
};
