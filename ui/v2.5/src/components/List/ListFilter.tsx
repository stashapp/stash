import _, { debounce } from "lodash";
import React, { useState, useEffect } from "react";
import Mousetrap from "mousetrap";
import { SortDirectionEnum } from "src/core/generated-graphql";
import {
  Badge,
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
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { AddFilter } from "./AddFilter";

interface IListFilterOperation {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
}

interface IListFilterProps {
  onFilterUpdate: (newFilter: ListFilterModel) => void;
  onSelectAll?: () => void;
  onSelectNone?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  otherOperations?: IListFilterOperation[];
  filter: ListFilterModel;
  filterOptions: ListFilterOptions;
  itemsSelected?: boolean;
}

const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120", "250", "500", "1000"];

export const ListFilter: React.FC<IListFilterProps> = (
  props: IListFilterProps
) => {
  const [queryRef, setQueryFocus] = useFocus();

  const searchCallback = debounce((value: string) => {
    const newFilter = _.cloneDeep(props.filter);
    newFilter.searchTerm = value;
    newFilter.currentPage = 1;
    props.onFilterUpdate(newFilter);
  }, 500);

  const [editingCriterion, setEditingCriterion] = useState<
    Criterion<CriterionValue> | undefined
  >(undefined);

  const intl = useIntl();

  useEffect(() => {
    Mousetrap.bind("/", (e) => {
      setQueryFocus();
      e.preventDefault();
    });

    Mousetrap.bind("r", () => onReshuffleRandomSort());
    Mousetrap.bind("s a", () => onSelectAll());
    Mousetrap.bind("s n", () => onSelectNone());

    if (props.itemsSelected) {
      Mousetrap.bind("e", () => {
        if (props.onEdit) {
          props.onEdit();
        }
      });

      Mousetrap.bind("d d", () => {
        if (props.onDelete) {
          props.onDelete();
        }
      });
    }

    return () => {
      Mousetrap.unbind("/");
      Mousetrap.unbind("r");
      Mousetrap.unbind("s a");
      Mousetrap.unbind("s n");

      if (props.itemsSelected) {
        Mousetrap.unbind("e");
        Mousetrap.unbind("d d");
      }
    };
  });

  function onChangePageSize(event: React.ChangeEvent<HTMLSelectElement>) {
    const val = event.currentTarget.value;

    const newFilter = _.cloneDeep(props.filter);
    newFilter.itemsPerPage = parseInt(val, 10);
    newFilter.currentPage = 1;
    props.onFilterUpdate(newFilter);
  }

  function onChangeQuery(event: React.FormEvent<HTMLInputElement>) {
    searchCallback(event.currentTarget.value);
  }

  function onChangeSortDirection() {
    const newFilter = _.cloneDeep(props.filter);
    if (props.filter.sortDirection === SortDirectionEnum.Asc) {
      newFilter.sortDirection = SortDirectionEnum.Desc;
    } else {
      newFilter.sortDirection = SortDirectionEnum.Asc;
    }

    props.onFilterUpdate(newFilter);
  }

  function onChangeSortBy(eventKey: string | null) {
    const newFilter = _.cloneDeep(props.filter);
    newFilter.sortBy = eventKey ?? undefined;
    newFilter.currentPage = 1;
    props.onFilterUpdate(newFilter);
  }

  function onReshuffleRandomSort() {
    const newFilter = _.cloneDeep(props.filter);
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
    props.onFilterUpdate(newFilter);
  }

  function onAddCriterion(
    criterion: Criterion<CriterionValue>,
    oldId?: string
  ) {
    const newFilter = _.cloneDeep(props.filter);

    // Find if we are editing an existing criteria, then modify that.  Or create a new one.
    const existingIndex = newFilter.criteria.findIndex((c) => {
      // If we modified an existing criterion, then look for the old id.
      const id = oldId || criterion.getId();
      return c.getId() === id;
    });
    if (existingIndex === -1) {
      newFilter.criteria.push(criterion);
    } else {
      newFilter.criteria[existingIndex] = criterion;
    }

    // Remove duplicate modifiers
    newFilter.criteria = newFilter.criteria.filter((obj, pos, arr) => {
      return arr.map((mapObj) => mapObj.getId()).indexOf(obj.getId()) === pos;
    });

    newFilter.currentPage = 1;
    props.onFilterUpdate(newFilter);
  }

  function onCancelAddCriterion() {
    setEditingCriterion(undefined);
  }

  function onRemoveCriterion(removedCriterion: Criterion<CriterionValue>) {
    const newFilter = _.cloneDeep(props.filter);
    newFilter.criteria = newFilter.criteria.filter(
      (criterion) => criterion.getId() !== removedCriterion.getId()
    );
    newFilter.currentPage = 1;
    props.onFilterUpdate(newFilter);
  }

  let removedCriterionId = "";
  function onRemoveCriterionTag(criterion?: Criterion<CriterionValue>) {
    if (!criterion) {
      return;
    }
    setEditingCriterion(undefined);
    removedCriterionId = criterion.getId();
    onRemoveCriterion(criterion);
  }

  function onClickCriterionTag(criterion?: Criterion<CriterionValue>) {
    if (!criterion || removedCriterionId !== "") {
      return;
    }
    setEditingCriterion(criterion);
  }

  function renderSortByOptions() {
    return props.filterOptions.sortByOptions
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

  function renderFilterTags() {
    return props.filter.criteria.map((criterion) => (
      <Badge
        className="tag-item"
        variant="secondary"
        key={criterion.getId()}
        onClick={() => onClickCriterionTag(criterion)}
      >
        {criterion.getLabel(intl)}
        <Button
          variant="secondary"
          onClick={() => onRemoveCriterionTag(criterion)}
        >
          <Icon icon="times" />
        </Button>
      </Badge>
    ));
  }

  function onSelectAll() {
    if (props.onSelectAll) {
      props.onSelectAll();
    }
  }

  function onSelectNone() {
    if (props.onSelectNone) {
      props.onSelectNone();
    }
  }

  function onEdit() {
    if (props.onEdit) {
      props.onEdit();
    }
  }

  function onDelete() {
    if (props.onDelete) {
      props.onDelete();
    }
  }

  function renderSelectAll() {
    if (props.onSelectAll) {
      return (
        <Dropdown.Item
          key="select-all"
          className="bg-secondary text-white"
          onClick={() => onSelectAll()}
        >
          Select All
        </Dropdown.Item>
      );
    }
  }

  function renderSelectNone() {
    if (props.onSelectNone) {
      return (
        <Dropdown.Item
          key="select-none"
          className="bg-secondary text-white"
          onClick={() => onSelectNone()}
        >
          Select None
        </Dropdown.Item>
      );
    }
  }

  function renderMore() {
    const options = [renderSelectAll(), renderSelectNone()].filter((o) => o);

    if (props.otherOperations) {
      props.otherOperations
        .filter((o) => {
          if (!o.isDisplayed) {
            return true;
          }

          return o.isDisplayed();
        })
        .forEach((o) => {
          options.push(
            <Dropdown.Item
              key={o.text}
              className="bg-secondary text-white"
              onClick={o.onClick}
            >
              {o.text}
            </Dropdown.Item>
          );
        });
    }

    if (options.length > 0) {
      return (
        <Dropdown>
          <Dropdown.Toggle variant="secondary" id="more-menu">
            <Icon icon="ellipsis-h" />
          </Dropdown.Toggle>
          <Dropdown.Menu className="bg-secondary text-white">
            {options}
          </Dropdown.Menu>
        </Dropdown>
      );
    }
  }

  function maybeRenderSelectedButtons() {
    if (props.itemsSelected && (props.onEdit || props.onDelete)) {
      return (
        <ButtonGroup className="ml-2">
          {props.onEdit && (
            <OverlayTrigger overlay={<Tooltip id="edit">Edit</Tooltip>}>
              <Button variant="secondary" onClick={onEdit}>
                <Icon icon="pencil-alt" />
              </Button>
            </OverlayTrigger>
          )}

          {props.onDelete && (
            <OverlayTrigger overlay={<Tooltip id="delete">Delete</Tooltip>}>
              <Button variant="danger" onClick={onDelete}>
                <Icon icon="trash" />
              </Button>
            </OverlayTrigger>
          )}
        </ButtonGroup>
      );
    }
  }

  function render() {
    const currentSortBy = props.filterOptions.sortByOptions.find(
      (o) => o.value === props.filter.sortBy
    );

    return (
      <>
        <div className="d-flex">
          <InputGroup className="mr-2 flex-grow-1">
            <FormControl
              ref={queryRef}
              placeholder="Search..."
              defaultValue={props.filter.searchTerm}
              onInput={onChangeQuery}
              className="bg-secondary text-white border-secondary w-50"
            />

            <InputGroup.Append>
              <AddFilter
                filterOptions={props.filterOptions}
                onAddCriterion={onAddCriterion}
                onCancel={onCancelAddCriterion}
                editingCriterion={editingCriterion}
              />
            </InputGroup.Append>
          </InputGroup>

          <Dropdown as={ButtonGroup} className="mr-2">
            <Dropdown.Toggle split variant="secondary" id="more-menu">
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
                  {props.filter.sortDirection === SortDirectionEnum.Asc
                    ? "Ascending"
                    : "Descending"}
                </Tooltip>
              }
            >
              <Button variant="secondary" onClick={onChangeSortDirection}>
                <Icon
                  icon={
                    props.filter.sortDirection === SortDirectionEnum.Asc
                      ? "caret-up"
                      : "caret-down"
                  }
                />
              </Button>
            </OverlayTrigger>
            {props.filter.sortBy === "random" && (
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
          value={props.filter.itemsPerPage.toString()}
          className="btn-secondary mx-1"
        >
          {PAGE_SIZE_OPTIONS.map((s) => (
            <option value={s} key={s}>
              {s}
            </option>
          ))}
        </Form.Control>

        {maybeRenderSelectedButtons()}

        <div className="mx-2">{renderMore()}</div>

        <div className="d-flex justify-content-center">
          {renderFilterTags()}
        </div>
      </>
    );
  }

  return render();
};
