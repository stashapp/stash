import React, { useMemo } from "react";
import { SortDirectionEnum } from "src/core/generated-graphql";
import {
  Button,
  ButtonGroup,
  Dropdown,
  OverlayTrigger,
  Tooltip,
  InputGroup,
} from "react-bootstrap";

import { Icon } from "../Shared/Icon";
import { useIntl } from "react-intl";
import {
  faCaretDown,
  faCaretUp,
  faRandom,
  faThumbtack,
} from "@fortawesome/free-solid-svg-icons";
import { ISortByOption } from "src/models/list-filter/filter-options";
import { useConfigurationContext } from "src/hooks/Config";
import ClearableInput from "../Shared/ClearableInput";
import useFocus from "src/utils/focus";
import ScreenUtils from "src/utils/screen";

interface IMessageValue {
  value: string;
  message: string;
}

export const SortBySelect: React.FC<{
  className?: string;
  sortBy: string | undefined;
  sortDirection: SortDirectionEnum;
  options: ISortByOption[];
  pinnedOptions?: string[];
  togglePinSortBy?: (sortBy: string) => void;
  onChangeSortBy: (eventKey: string | null) => void;
  onChangeSortDirection: () => void;
  onReshuffleRandomSort: () => void;
}> = ({
  className,
  sortBy,
  sortDirection,
  options,
  pinnedOptions = [],
  togglePinSortBy,
  onChangeSortBy,
  onChangeSortDirection,
  onReshuffleRandomSort,
}) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();
  const { sfwContentMode } = configuration.interface;

  const [filter, setFilter] = React.useState("");

  const focusOnOpen = !ScreenUtils.isTouch();
  const focusRef = useFocus();
  const [, setFocus] = focusRef;

  const currentSortBy = options.find((o) => o.value === sortBy);
  const currentSortByMessageID = currentSortBy
    ? !sfwContentMode
      ? currentSortBy.messageID
      : (currentSortBy.sfwMessageID ?? currentSortBy.messageID)
    : "";

  const sortedOptions: IMessageValue[] = useMemo(() => {
    return options
      .map((o) => {
        const messageID = !sfwContentMode
          ? o.messageID
          : (o.sfwMessageID ?? o.messageID);
        return {
          message: intl.formatMessage({ id: messageID }),
          value: o.value,
        };
      })
      .sort((a, b) => a.message.localeCompare(b.message));
  }, [intl, options, sfwContentMode]);

  const filteredOptions = useMemo(() => {
    if (!filter) return sortedOptions;

    return sortedOptions.filter((o) =>
      o.message.toLowerCase().includes(filter.toLowerCase())
    );
  }, [sortedOptions, filter]);

  const filteredPinnedOptions = useMemo(() => {
    return filteredOptions.filter((o) => pinnedOptions.includes(o.value));
  }, [filteredOptions, pinnedOptions]);

  const unpinnedOptions = useMemo(() => {
    return filteredOptions.filter((o) => !pinnedOptions.includes(o.value));
  }, [filteredOptions, pinnedOptions]);

  const dropdownItems = useMemo(() => {
    function togglePin(
      e: React.MouseEvent<HTMLElement>,
      option: IMessageValue
    ) {
      e.stopPropagation();
      e.preventDefault();
      togglePinSortBy?.(option.value);
    }

    function renderOptionMessage(option: IMessageValue, isPin: boolean) {
      return (
        <Dropdown.Item
          onSelect={onChangeSortBy}
          key={option.value}
          className="bg-secondary text-white sort-by-option"
          eventKey={option.value}
          data-value={option.value}
        >
          {option.message}
          {togglePinSortBy && (
            <Button
              className="pin-sort-by-button"
              variant="minimal"
              onClick={(e) => togglePin(e, option)}
            >
              <Icon icon={faThumbtack} className={isPin ? "" : "tilted"} />
            </Button>
          )}
        </Dropdown.Item>
      );
    }

    return (
      <>
        {filteredPinnedOptions.map((option) =>
          renderOptionMessage(option, true)
        )}
        {filteredPinnedOptions.length > 0 && <Dropdown.Divider />}
        {unpinnedOptions.map((option) => renderOptionMessage(option, false))}
      </>
    );
  }, [filteredPinnedOptions, unpinnedOptions, onChangeSortBy, togglePinSortBy]);

  return (
    <Dropdown
      as={ButtonGroup}
      className={`${className ?? ""} sort-by-select`}
      onToggle={(v) => {
        if (focusOnOpen && v) setTimeout(() => setFocus(true), 0);
      }}
    >
      <InputGroup.Prepend>
        <Dropdown.Toggle variant="secondary">
          {currentSortBy
            ? intl.formatMessage({ id: currentSortByMessageID })
            : ""}
        </Dropdown.Toggle>
      </InputGroup.Prepend>
      <Dropdown.Menu className="bg-secondary text-white">
        <div className="sort-by-filter-container">
          <ClearableInput
            placeholder={`${intl.formatMessage({ id: "filter" })}...`}
            value={filter}
            setValue={setFilter}
            focus={focusRef}
          />
        </div>

        {dropdownItems}
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
