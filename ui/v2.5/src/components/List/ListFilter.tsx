import _, { debounce } from "lodash";
import React, { useState, useEffect } from "react";
import { SortDirectionEnum } from "src/core/generated-graphql";
import {
  Badge,
  Button,
  ButtonGroup,
  Dropdown,
  Form,
  OverlayTrigger,
  Tooltip,
  SafeAnchorProps,
  InputGroup,
  FormControl,
  ButtonToolbar,
} from "react-bootstrap";

import { Icon } from "src/components/Shared";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { useFocus } from "src/utils";
import { AddFilter } from "./AddFilter";

interface IListFilterOperation {
  text: string;
  onClick: () => void;
}

interface IListFilterProps {
  subComponent?: boolean;
  onFilterUpdate: (newFilter: ListFilterModel) => void;
  zoomIndex?: number;
  onChangeZoom?: (zoomIndex: number) => void;
  onSelectAll?: () => void;
  onSelectNone?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  otherOperations?: IListFilterOperation[];
  filter: ListFilterModel;
  itemsSelected?: boolean;
}

const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120"];
const minZoom = 0;
const maxZoom = 3;

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
    Criterion | undefined
  >(undefined);

  useEffect(() => {
    Mousetrap.bind("/", (e) => {
      setQueryFocus();
      e.preventDefault();
    });

    Mousetrap.bind("r", () => onReshuffleRandomSort());
    Mousetrap.bind("v g", () => {
      if (props.filter.displayModeOptions.includes(DisplayMode.Grid)) {
        onChangeDisplayMode(DisplayMode.Grid);
      }
    });
    Mousetrap.bind("v l", () => {
      if (props.filter.displayModeOptions.includes(DisplayMode.List)) {
        onChangeDisplayMode(DisplayMode.List);
      }
    });
    Mousetrap.bind("v w", () => {
      if (props.filter.displayModeOptions.includes(DisplayMode.Wall)) {
        onChangeDisplayMode(DisplayMode.Wall);
      }
    });
    Mousetrap.bind("+", () => {
      if (
        props.onChangeZoom &&
        props.zoomIndex !== undefined &&
        props.zoomIndex < maxZoom
      ) {
        props.onChangeZoom(props.zoomIndex + 1);
      }
    });
    Mousetrap.bind("-", () => {
      if (
        props.onChangeZoom &&
        props.zoomIndex !== undefined &&
        props.zoomIndex > minZoom
      ) {
        props.onChangeZoom(props.zoomIndex - 1);
      }
    });
    Mousetrap.bind("s a", () => onSelectAll());
    Mousetrap.bind("s n", () => onSelectNone());

    if (!props.subComponent && props.itemsSelected) {
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
      Mousetrap.unbind("v g");
      Mousetrap.unbind("v l");
      Mousetrap.unbind("v w");
      Mousetrap.unbind("+");
      Mousetrap.unbind("-");
      Mousetrap.unbind("s a");
      Mousetrap.unbind("s n");

      if (!props.subComponent && props.itemsSelected) {
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

  function onChangeSortBy(event: React.MouseEvent<SafeAnchorProps>) {
    const target = event.currentTarget as HTMLAnchorElement;

    const newFilter = _.cloneDeep(props.filter);
    newFilter.sortBy = target.text;
    newFilter.currentPage = 1;
    props.onFilterUpdate(newFilter);
  }

  function onReshuffleRandomSort() {
    const newFilter = _.cloneDeep(props.filter);
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
    props.onFilterUpdate(newFilter);
  }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = _.cloneDeep(props.filter);
    newFilter.displayMode = displayMode;
    props.onFilterUpdate(newFilter);
  }

  function onAddCriterion(criterion: Criterion, oldId?: string) {
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

  function onRemoveCriterion(removedCriterion: Criterion) {
    const newFilter = _.cloneDeep(props.filter);
    newFilter.criteria = newFilter.criteria.filter(
      (criterion) => criterion.getId() !== removedCriterion.getId()
    );
    newFilter.currentPage = 1;
    props.onFilterUpdate(newFilter);
  }

  let removedCriterionId = "";
  function onRemoveCriterionTag(criterion?: Criterion) {
    if (!criterion) {
      return;
    }
    setEditingCriterion(undefined);
    removedCriterionId = criterion.getId();
    onRemoveCriterion(criterion);
  }

  function onClickCriterionTag(criterion?: Criterion) {
    if (!criterion || removedCriterionId !== "") {
      return;
    }
    setEditingCriterion(criterion);
  }

  function renderSortByOptions() {
    return props.filter.sortByOptions.map((option) => (
      <Dropdown.Item
        onClick={onChangeSortBy}
        key={option}
        className="bg-secondary text-white"
      >
        {option}
      </Dropdown.Item>
    ));
  }

  function renderDisplayModeOptions() {
    function getIcon(option: DisplayMode) {
      switch (option) {
        case DisplayMode.Grid:
          return "th-large";
        case DisplayMode.List:
          return "list";
        case DisplayMode.Wall:
          return "square";
      }
    }
    function getLabel(option: DisplayMode) {
      switch (option) {
        case DisplayMode.Grid:
          return "Grid";
        case DisplayMode.List:
          return "List";
        case DisplayMode.Wall:
          return "Wall";
      }
    }

    return props.filter.displayModeOptions.map((option) => (
      <OverlayTrigger
        key={option}
        overlay={
          <Tooltip id="display-mode-tooltip">{getLabel(option)}</Tooltip>
        }
      >
        <Button
          variant="secondary"
          active={props.filter.displayMode === option}
          onClick={() => onChangeDisplayMode(option)}
        >
          <Icon icon={getIcon(option)} />
        </Button>
      </OverlayTrigger>
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
        {criterion.getLabel()}
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
    const options = [renderSelectAll(), renderSelectNone()];

    if (props.otherOperations) {
      props.otherOperations.forEach((o) => {
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

  function onChangeZoom(v: number) {
    if (props.onChangeZoom) {
      props.onChangeZoom(v);
    }
  }

  function maybeRenderZoom() {
    if (props.onChangeZoom) {
      return (
        <div className="align-middle">
          <Form.Control
            className="zoom-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={minZoom}
            max={maxZoom}
            value={props.zoomIndex}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              onChangeZoom(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </div>
      );
    }
  }

  function maybeRenderSelectedButtons() {
    if (props.itemsSelected) {
      return (
        <>
          {props.onEdit ? (
            <ButtonGroup className="mr-1">
              <OverlayTrigger overlay={<Tooltip id="edit">Edit</Tooltip>}>
                <Button variant="secondary" onClick={onEdit}>
                  <Icon icon="pencil-alt" />
                </Button>
              </OverlayTrigger>
            </ButtonGroup>
          ) : undefined}

          {props.onDelete ? (
            <ButtonGroup className="mr-1">
              <OverlayTrigger overlay={<Tooltip id="delete">Delete</Tooltip>}>
                <Button variant="danger" onClick={onDelete}>
                  <Icon icon="trash" />
                </Button>
              </OverlayTrigger>
            </ButtonGroup>
          ) : undefined}
        </>
      );
    }
  }

  function render() {
    return (
      <>
        <ButtonToolbar className="align-items-center justify-content-center">
          <ButtonGroup className="mr-3 my-1">
            <InputGroup className="mr-2">
              <FormControl
                ref={queryRef}
                placeholder="Search..."
                defaultValue={props.filter.searchTerm}
                onInput={onChangeQuery}
                className="bg-secondary text-white border-secondary"
              />

              <InputGroup.Append>
                <AddFilter
                  filter={props.filter}
                  onAddCriterion={onAddCriterion}
                  onCancel={onCancelAddCriterion}
                  editingCriterion={editingCriterion}
                />
              </InputGroup.Append>
            </InputGroup>

            <Dropdown as={ButtonGroup} className="mr-2">
              <Dropdown.Toggle split variant="secondary" id="more-menu">
                {props.filter.sortBy}
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

            <Form.Control
              as="select"
              onChange={onChangePageSize}
              value={props.filter.itemsPerPage.toString()}
              className="btn-secondary mr-1"
            >
              {PAGE_SIZE_OPTIONS.map((s) => (
                <option value={s} key={s}>
                  {s}
                </option>
              ))}
            </Form.Control>
          </ButtonGroup>

          <ButtonGroup className="mr-3 my-1">
            {maybeRenderSelectedButtons()}
            {renderMore()}
          </ButtonGroup>

          <ButtonGroup className="my-1">
            {renderDisplayModeOptions()}
            {maybeRenderZoom()}
          </ButtonGroup>
        </ButtonToolbar>

        <div className="d-flex justify-content-center">
          {renderFilterTags()}
        </div>
      </>
    );
  }

  return render();
};
