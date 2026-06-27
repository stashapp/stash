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
} from "@fortawesome/free-solid-svg-icons";
import { ISortByOption } from "src/models/list-filter/filter-options";
import { useConfigurationContext } from "src/hooks/Config";
import ClearableInput from "../Shared/ClearableInput";
import useFocus from "src/utils/focus";
import ScreenUtils from "src/utils/screen";

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

  const sortedOptions = useMemo(() => {
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

  const dropdownItems = useMemo(() => {
    return filteredOptions.map((option) => (
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
  }, [filteredOptions, onChangeSortBy]);

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
