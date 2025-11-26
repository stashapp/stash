import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
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
  Popover,
  Overlay,
} from "react-bootstrap";

import { Icon } from "../Shared/Icon";
import { ListFilterModel } from "src/models/list-filter/filter";
import useFocus from "src/utils/focus";
import { useIntl } from "react-intl";
import {
  faCaretDown,
  faCaretUp,
  faCheck,
  faRandom,
} from "@fortawesome/free-solid-svg-icons";
import { useDebounce } from "src/hooks/debounce";
import { ClearableInput } from "../Shared/ClearableInput";
import { useStopWheelScroll } from "src/utils/form";
import { ISortByOption } from "src/models/list-filter/filter-options";
import { useConfigurationContext } from "src/hooks/Config";

export function useDebouncedSearchInput(
  filter: ListFilterModel,
  setFilter: (filter: ListFilterModel) => void
) {
  const callback = useCallback(
    (value: string) => {
      const newFilter = filter.clone();
      newFilter.searchTerm = value;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    },
    [filter, setFilter]
  );

  const onClear = useCallback(() => callback(""), [callback]);

  const searchCallback = useDebounce(callback, 500);

  return { searchCallback, onClear };
}

export const SearchTermInput: React.FC<{
  filter: ListFilterModel;
  onFilterUpdate: (newFilter: ListFilterModel) => void;
  focus?: ReturnType<typeof useFocus>;
}> = ({ filter, onFilterUpdate, focus: providedFocus }) => {
  const intl = useIntl();
  const [localInput, setLocalInput] = useState(filter.searchTerm);

  const localFocus = useFocus();
  const focus = providedFocus ?? localFocus;
  const [, setQueryFocus] = focus;

  useEffect(() => {
    setLocalInput(filter.searchTerm);
  }, [filter.searchTerm]);

  const { searchCallback, onClear } = useDebouncedSearchInput(
    filter,
    onFilterUpdate
  );

  useEffect(() => {
    Mousetrap.bind("/", (e) => {
      setQueryFocus();
      e.preventDefault();
    });

    return () => {
      Mousetrap.unbind("/");
    };
  });

  function onSetQuery(value: string) {
    setLocalInput(value);

    if (!value) {
      onClear();
    }

    searchCallback(value);
  }

  return (
    <ClearableInput
      className="search-term-input"
      focus={focus}
      value={localInput}
      setValue={onSetQuery}
      placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
    />
  );
};

const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120", "250", "500", "1000"];

export const PageSizeSelector: React.FC<{
  pageSize: number;
  setPageSize: (pageSize: number) => void;
}> = ({ pageSize, setPageSize }) => {
  const intl = useIntl();

  const perPageSelect = useRef(null);
  const [perPageInput, perPageFocus] = useFocus();
  const [customPageSizeShowing, setCustomPageSizeShowing] = useState(false);

  useEffect(() => {
    if (customPageSizeShowing) {
      perPageFocus();
    }
  }, [customPageSizeShowing, perPageFocus]);

  useStopWheelScroll(perPageInput);

  const pageSizeOptions = useMemo(() => {
    const ret = PAGE_SIZE_OPTIONS.map((o) => {
      return {
        label: o,
        value: o,
      };
    });
    const currentPerPage = pageSize.toString();
    if (!ret.find((o) => o.value === currentPerPage)) {
      ret.push({ label: currentPerPage, value: currentPerPage });
      ret.sort((a, b) => parseInt(a.value, 10) - parseInt(b.value, 10));
    }

    ret.push({
      label: `${intl.formatMessage({ id: "custom" })}...`,
      value: "custom",
    });

    return ret;
  }, [intl, pageSize]);

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

    setPageSize(pp);
  }

  return (
    <div className="page-size-selector">
      <Form.Control
        as="select"
        ref={perPageSelect}
        onChange={(e) => onChangePageSize(e.target.value)}
        value={pageSize.toString()}
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
              {/* can't use NumberField because of the ref */}
              <Form.Control
                type="number"
                min={1}
                className="text-input"
                ref={perPageInput}
                onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) => {
                  if (e.key === "Enter") {
                    onChangePageSize(
                      (perPageInput.current as HTMLInputElement)?.value ?? ""
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
                      (perPageInput.current as HTMLInputElement)?.value ?? ""
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
  );
};

export const SortBySelect: React.FC<{
  className?: string;
  sortBy: string | undefined;
  sortDirection: SortDirectionEnum;
  options: ISortByOption[];
  onChangeSortBy: (eventKey: string | null) => void;
  onChangeSortDirection: () => void;
  onReshuffleRandomSort: () => void;
}> = ({
  className,
  sortBy,
  sortDirection,
  options,
  onChangeSortBy,
  onChangeSortDirection,
  onReshuffleRandomSort,
}) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();
  const { sfwContentMode } = configuration.interface;

  const currentSortBy = options.find((o) => o.value === sortBy);
  const currentSortByMessageID = currentSortBy
    ? !sfwContentMode
      ? currentSortBy.messageID
      : currentSortBy.sfwMessageID ?? currentSortBy.messageID
    : "";

  function renderSortByOptions() {
    return options
      .map((o) => {
        const messageID = !sfwContentMode
          ? o.messageID
          : o.sfwMessageID ?? o.messageID;
        return {
          message: intl.formatMessage({ id: messageID }),
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
          data-value={option.value}
        >
          {option.message}
        </Dropdown.Item>
      ));
  }

  return (
    <Dropdown as={ButtonGroup} className={`${className ?? ""} sort-by-select`}>
      <InputGroup.Prepend>
        <Dropdown.Toggle variant="secondary">
          {currentSortBy
            ? intl.formatMessage({ id: currentSortByMessageID })
            : ""}
        </Dropdown.Toggle>
      </InputGroup.Prepend>
      <Dropdown.Menu className="bg-secondary text-white">
        {renderSortByOptions()}
      </Dropdown.Menu>
      <OverlayTrigger
        overlay={
          <Tooltip id="sort-direction-tooltip">
            {sortDirection === SortDirectionEnum.Asc
              ? intl.formatMessage({ id: "ascending" })
              : intl.formatMessage({ id: "descending" })}
          </Tooltip>
        }
      >
        <Button variant="secondary" onClick={onChangeSortDirection}>
          <Icon
            icon={
              sortDirection === SortDirectionEnum.Asc ? faCaretUp : faCaretDown
            }
          />
        </Button>
      </OverlayTrigger>
      {sortBy === "random" && (
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
  );
};
