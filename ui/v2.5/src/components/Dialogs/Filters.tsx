import React, { useState, useEffect } from "react";
import { Form, Button, ListGroup } from "react-bootstrap";
import { Modal, Icon } from "src/components/Shared";
import { FormattedMessage, useIntl } from "react-intl";
import { SavedFilter } from "src/core/generated-graphql";

interface IFilterEditor {
  isNew: boolean;
  availableFilters: SavedFilter[];
  filter: SavedFilter;
  saveFilter: (s?: SavedFilter) => void;
}

export const FiltersEditor: React.FC<IFilterEditor> = ({
  isNew,
  availableFilters,
  filter: initialFilter,
  saveFilter,
}) => {
  const [filter, setFilter] = useState<SavedFilter>(initialFilter);
  const intl = useIntl();

  // if id is empty, then we are adding a new filter
  const headerMsgId = isNew ? "actions.add" : "dialogs.edit_entity_title";
  const acceptMsgId = isNew ? "actions.add" : "actions.confirm";

  function handleFilterSelect(e: React.ChangeEvent<HTMLSelectElement>) {
    const selectedFilter = availableFilters.find(
      (s) => s.id === e.currentTarget.value
    );
    if (!selectedFilter) return;

    setFilter({
      ...filter,
      id: selectedFilter.id,
      name: selectedFilter.name,
      filter: selectedFilter.filter,
      mode: selectedFilter.mode,
    });
  }

  return (
    <Modal
      dialogClassName="identify-filter-editor"
      modalProps={{ animation: false, size: "lg" }}
      show
      icon={isNew ? "plus" : "pencil-alt"}
      header={intl.formatMessage(
        { id: headerMsgId },
        {
          count: 1,
          singularEntity: filter?.name,
          pluralEntity: filter?.name,
        }
      )}
      accept={{
        onClick: () => saveFilter(filter),
        text: intl.formatMessage({ id: acceptMsgId }),
      }}
      cancel={{
        onClick: () => saveFilter(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
    >
      <Form>
        {isNew && (
          <Form.Group>
            <h5>
              <FormattedMessage id="filter" />
            </h5>
            <Form.Control
              as="select"
              value={filter.id}
              className="input-control"
              onChange={handleFilterSelect}
            >
              {availableFilters.map((i) => (
                <option value={i.id} key={i.id}>
                  {i.name}
                </option>
              ))}
            </Form.Control>
          </Form.Group>
        )}
      </Form>
    </Modal>
  );
};

interface IFiltersList {
  filters: SavedFilter[];
  setFilters: (f: SavedFilter[]) => void;
  editFilter: (f?: SavedFilter) => void;
  canAdd: boolean;
}

export const FiltersList: React.FC<IFiltersList> = ({
  filters,
  setFilters,
  editFilter,
  canAdd,
}) => {
  const [tempFilters, setTempFilters] = useState(filters);
  const [dragIndex, setDragIndex] = useState<number | undefined>();
  const [mouseOverIndex, setMouseOverIndex] = useState<number | undefined>();

  useEffect(() => {
    setTempFilters([...filters]);
  }, [filters]);

  function removeFilter(index: number) {
    const newFilters = [...filters];
    newFilters.splice(index, 1);
    setFilters(newFilters);
  }

  function onDragStart(event: React.DragEvent<HTMLElement>, index: number) {
    event.dataTransfer.effectAllowed = "move";
    setDragIndex(index);
  }

  function onDragOver(event: React.DragEvent<HTMLElement>, index?: number) {
    if (dragIndex !== undefined && index !== undefined && index !== dragIndex) {
      const newFilters = [...tempFilters];
      const moved = newFilters.splice(dragIndex, 1);
      newFilters.splice(index, 0, moved[0]);
      setTempFilters(newFilters);
      setDragIndex(index);
    }

    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDragOverDefault(event: React.DragEvent<HTMLDivElement>) {
    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDrop() {
    // assume we've already set the temp filter list
    // feed it up
    setFilters(tempFilters!);
    setDragIndex(undefined);
    setMouseOverIndex(undefined);
  }

  return (
    <Form.Group className="scraper-filters" onDragOver={onDragOverDefault}>
      <h5>Filters</h5>
      <ListGroup as="ul" className="scraper-filter-list">
        {tempFilters.map((s, index) => (
          <ListGroup.Item
            as="li"
            key={s.id}
            className="d-flex justify-content-between align-items-center"
            draggable={mouseOverIndex === index}
            onDragStart={(e) => onDragStart(e, index)}
            onDragEnter={(e) => onDragOver(e, index)}
            onDrop={() => onDrop()}
          >
            <div>
              <div
                className="minimal text-muted drag-handle"
                onMouseEnter={() => setMouseOverIndex(index)}
                onMouseLeave={() => setMouseOverIndex(undefined)}
              >
                <Icon icon="grip-vertical" />
              </div>
              {s.name}
            </div>
            <div>
              <Button
                className="minimal text-danger"
                onClick={() => removeFilter(index)}
              >
                <Icon icon="minus" />
              </Button>
            </div>
          </ListGroup.Item>
        ))}
      </ListGroup>
      {canAdd && (
        <div className="text-right">
          <Button
            className="minimal add-scraper-filter-button"
            onClick={() => editFilter()}
          >
            <Icon icon="plus" />
          </Button>
        </div>
      )}
    </Form.Group>
  );
};
