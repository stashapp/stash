import { debounce } from "lodash";
import React, { SyntheticEvent, useCallback, useState } from "react";
import { Badge, Button, ButtonGroup, Dropdown, Form, OverlayTrigger, Tooltip } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

import { Criterion } from "../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode } from "../../models/list-filter/types";
import { AddFilter } from "./AddFilter";

interface IListFilterOperation {
  text: string;
  onClick: () => void;
}

interface IListFilterProps {
  onChangePageSize: (pageSize: number) => void;
  onChangeQuery: (query: string) => void;
  onChangeSortDirection: (sortDirection: "asc" | "desc") => void;
  onChangeSortBy: (sortBy: string) => void;
  onChangeDisplayMode: (displayMode: DisplayMode) => void;
  onAddCriterion: (criterion: Criterion, oldId?: string) => void;
  onRemoveCriterion: (criterion: Criterion) => void;
  zoomIndex?: number;
  onChangeZoom?: (zoomIndex: number) => void;
  onSelectAll?: () => void;
  onSelectNone?: () => void;
  otherOperations?: IListFilterOperation[];
  filter: ListFilterModel;
}

const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120"];

export const ListFilter: React.FC<IListFilterProps> = (props: IListFilterProps) => {
  const searchCallback = useCallback(
    debounce((event: any) => {
      props.onChangeQuery(event.target.value);
    }, 500), [props.onChangeQuery]
  );

  const [editingCriterion, setEditingCriterion] = useState<Criterion | undefined>(undefined);

  function onChangePageSize(event: SyntheticEvent<HTMLSelectElement>) {
    const val = event!.currentTarget!.value;
    props.onChangePageSize(parseInt(val, 10));
  }

  function onChangeQuery(event: SyntheticEvent<HTMLInputElement>) {
    event.persist();
    searchCallback(event);
  }

  function onChangeSortDirection(_: any) {
    if (props.filter.sortDirection === "asc") {
      props.onChangeSortDirection("desc");
    } else {
      props.onChangeSortDirection("asc");
    }
  }

  function onChangeSortBy(event: React.MouseEvent<any>) {
    props.onChangeSortBy(event.currentTarget.text);
  }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    props.onChangeDisplayMode(displayMode);
  }

  function onAddCriterion(criterion: Criterion, oldId?: string) {
    props.onAddCriterion(criterion, oldId);
  }

  function onCancelAddCriterion() {
    setEditingCriterion(undefined);
  }

  let removedCriterionId = "";
  function onRemoveCriterionTag(criterion?: Criterion) {
    if (!criterion) { return; }
    setEditingCriterion(undefined);
    removedCriterionId = criterion.getId();
    props.onRemoveCriterion(criterion);
  }
  function onClickCriterionTag(criterion?: Criterion) {
    if (!criterion || removedCriterionId !== "") { return; }
    setEditingCriterion(criterion);
  }

  function renderSortByOptions() {
    return props.filter.sortByOptions.map((option) => (
      <Dropdown.Item onClick={onChangeSortBy} key={option}>{option}</Dropdown.Item>
    ));
  }

  function renderDisplayModeOptions() {
    function getIcon(option: DisplayMode) {
      switch (option) {
        case DisplayMode.Grid: return "th-large";
        case DisplayMode.List: return "list";
        case DisplayMode.Wall: return "square";
      }
    }
    function getLabel(option: DisplayMode) {
      switch (option) {
        case DisplayMode.Grid: return "Grid";
        case DisplayMode.List: return "List";
        case DisplayMode.Wall: return "Wall";
      }
    }
    return props.filter.displayModeOptions.map((option) => (
      <OverlayTrigger overlay={<Tooltip id="display-mode-tooltip">{getLabel(option)}</Tooltip>}>
        <Button
          key={option}
          active={props.filter.displayMode === option}
          onClick={() => onChangeDisplayMode(option)}
        >
          <FontAwesomeIcon icon={getIcon(option)} />
        </Button>
      </OverlayTrigger>
    ));
  }

  function renderFilterTags() {
    return props.filter.criteria.map((criterion) => (
      <Badge
        className="tag-item"
        variant="secondary"
        onClick={() => onClickCriterionTag(criterion)}
      >
        {criterion.getLabel()}
        <Button onClick={() => onRemoveCriterionTag(criterion)}>
          <FontAwesomeIcon icon="times" />
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

  function renderSelectAll() {
    if (props.onSelectAll) {
      return <Dropdown.Item onClick={() => onSelectAll()}>Select All</Dropdown.Item>;
    }
  }

  function renderSelectNone() {
    if (props.onSelectNone) {
      return <Dropdown.Item onClick={() => onSelectNone()}>Select None</Dropdown.Item>;
    }
  }

  function renderMore() {
    let options = [
      renderSelectAll(),
      renderSelectNone()
    ];

    if (props.otherOperations) {
      props.otherOperations.forEach((o) => {
        options.push(<Dropdown.Item onClick={o.onClick}>{o.text}</Dropdown.Item>);
      });
    }

    if (options.length > 0) {
      return (
        <Dropdown>
          <Dropdown.Toggle variant="secondary" id="more-menu">
            <Button>
              <FontAwesomeIcon icon="ellipsis-h" />
            </Button>
          </Dropdown.Toggle>
          <Dropdown.Menu>
            {options}
          </Dropdown.Menu>
        </Dropdown>
      );
    }
  }

  function onChangeZoom(v : number) {
    if (props.onChangeZoom) {
      props.onChangeZoom(v);
    }
  } 

  function maybeRenderZoom() {
    if (props.onChangeZoom) {
      return (
        <span className="zoom-slider">
          <Form.Control
            type="range"
            min={0}
            max={3}
            onChange={(event: any) => onChangeZoom(Number.parseInt(event.target.value))}
          />
      </span>
      );
    }
  }

  function render() {
    return (
      <>
        <div className="filter-container">
          <Form.Control
            placeholder="Search..."
            value={props.filter.searchTerm}
            onChange={onChangeQuery}
            className="filter-item"
          />
          <Form.Control
            as="select"
            onChange={onChangePageSize}
            value={props.filter.itemsPerPage.toString()}
            className="filter-item"
          >
            { PAGE_SIZE_OPTIONS.map(s => <option value={s}>{s}</option>) }
          </Form.Control>
          <ButtonGroup className="filter-item">
            <Dropdown>
              <Dropdown.Toggle variant="secondary" id="more-menu">
                <Button>{props.filter.sortBy}</Button>
              </Dropdown.Toggle>
              <Dropdown.Menu>
                {renderSortByOptions()}
              </Dropdown.Menu>
            </Dropdown>

            <OverlayTrigger overlay={
              <Tooltip id="sort-direction-tooltip">{props.filter.sortDirection === "asc" ? "Ascending" : "Descending"}</Tooltip>
            }>
              <Button onClick={onChangeSortDirection}>
                <FontAwesomeIcon icon={props.filter.sortDirection === "asc" ? "caret-up" : "caret-down"} />
              </Button>
            </OverlayTrigger>
          </ButtonGroup>

          <AddFilter
            filter={props.filter}
            onAddCriterion={onAddCriterion}
            onCancel={onCancelAddCriterion}
            editingCriterion={editingCriterion}
          />

          <ButtonGroup className="filter-item">
            {renderDisplayModeOptions()}
          </ButtonGroup>

          {maybeRenderZoom()}

          <ButtonGroup className="filter-item">
            {renderMore()}
          </ButtonGroup>
        </div>
        <div style={{display: "flex", justifyContent: "center", margin: "10px auto"}}>
          {renderFilterTags()}
        </div>
      </>
    );
  }

  return render();
};
