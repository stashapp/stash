import { debounce } from "lodash";
import React, { useState } from "react";
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
  Col,
  Row,
} from "react-bootstrap";

import { Icon } from "src/components/Shared";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { AddFilter } from "./AddFilter";

interface IListFilterOperation {
  text: string;
  onClick: () => void;
}

interface IListFilterProps {
  onChangePageSize: (pageSize: number) => void;
  onChangeQuery: (query: string) => void;
  onChangeSortDirection: (sortDirection: SortDirectionEnum) => void;
  onChangeSortBy: (sortBy: string) => void;
  onSortReshuffle: () => void;
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

export const ListFilter: React.FC<IListFilterProps> = (
  props: IListFilterProps
) => {
  const searchCallback = debounce((value: string) => {
    props.onChangeQuery(value);
  }, 500);

  const [editingCriterion, setEditingCriterion] = useState<
    Criterion | undefined
  >(undefined);

  function onChangePageSize(event: React.ChangeEvent<HTMLSelectElement>) {
    const val = event.currentTarget.value;
    props.onChangePageSize(parseInt(val, 10));
  }

  function onChangeQuery(event: React.FormEvent<HTMLInputElement>) {
    searchCallback(event.currentTarget.value);
  }

  function onEdit() {}

  function onChangeSortDirection() {
    if (props.filter.sortDirection === SortDirectionEnum.Asc) {
      props.onChangeSortDirection(SortDirectionEnum.Desc);
    } else {
      props.onChangeSortDirection(SortDirectionEnum.Asc);
    }
  }

  function onChangeSortBy(event: React.MouseEvent<SafeAnchorProps>) {
    const target = event.currentTarget as HTMLAnchorElement;
    props.onChangeSortBy(target.text);
  }

  function onReshuffleRandomSort() {
    props.onSortReshuffle();
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
    if (!criterion) {
      return;
    }
    setEditingCriterion(undefined);
    removedCriterionId = criterion.getId();
    props.onRemoveCriterion(criterion);
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

    const option = DisplayMode.Grid;
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
        <Form.Control
          className="zoom-slider d-none d-sm-inline-flex"
          type="range"
          min={0}
          max={3}
          defaultValue={1}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChangeZoom(Number.parseInt(e.currentTarget.value, 10))
          }
        />
      );
    }
  }

  function maybeRenderMultiOps() {
    if (true) {
      return (
        <>
          <ButtonGroup className="mr-1">
            <OverlayTrigger
              overlay={
                <Tooltip id="edit">Edit</Tooltip>
              }
            >
              <Button variant="secondary" onClick={onEdit}>
                <Icon icon="pencil-alt" />
              </Button>
            </OverlayTrigger>
          </ButtonGroup>

          <ButtonGroup className="mr-1">
            <OverlayTrigger
              overlay={
                <Tooltip id="delete">Delete</Tooltip>
              }
            >
              <Button variant="danger" onClick={onEdit}>
                <Icon icon="trash" />
              </Button>
            </OverlayTrigger>
          </ButtonGroup>
        </>
      )
    }
  }

  function render() {
    return (
      <>
        <div className="form-row align-items-center justify-content-center">
          <Col sm={12} md={6} xl={4} lg={5} className="my-1">
            <Row className="justify-content-center">
              <Col xs={6} className="px-1">
                <InputGroup>
                  <FormControl 
                    placeholder="Search..."
                    defaultValue={props.filter.searchTerm}
                    onInput={onChangeQuery}
                    className="bg-secondary text-white border-secondary" />

                  <InputGroup.Append>
                    <AddFilter
                      filter={props.filter}
                      onAddCriterion={onAddCriterion}
                      onCancel={onCancelAddCriterion}
                      editingCriterion={editingCriterion}
                    />
                  </InputGroup.Append>
                </InputGroup>
              </Col>
              
              <Col xs="auto" className="px-1">
                <ButtonGroup>
                  <Dropdown as={ButtonGroup}>
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
                </ButtonGroup>
              </Col>

              <Col xs="auto" className="px-1">
                <Form.Control
                  as="select"
                  onChange={onChangePageSize}
                  value={props.filter.itemsPerPage.toString()}
                  className="btn-secondary"
                >
                  {PAGE_SIZE_OPTIONS.map((s) => (
                    <option value={s} key={s}>
                      {s}
                    </option>
                  ))}
                </Form.Control>
              </Col>
            </Row>
          </Col>      
          
          <Col sm={12} md="auto" className="my-1">
            <Row className="align-items-center justify-content-center">

            {maybeRenderMultiOps()}

            <ButtonGroup className="mr-3">
              {renderMore()}
            </ButtonGroup>

            <ButtonGroup className="mr-3">
              {renderDisplayModeOptions()}
            </ButtonGroup>

            <ButtonGroup>
              {maybeRenderZoom()}
            </ButtonGroup>
            </Row>
          </Col>      

          <Col xs="auto">
            
          </Col>
        </div>
        <div className="d-flex justify-content-center">
          {renderFilterTags()}
        </div>
      </>
    );
  }

  return render();
};
