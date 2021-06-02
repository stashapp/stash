import React, { useState } from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  FormControl,
  InputGroup,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import { FilterMode } from "src/models/list-filter/types";
import { Icon } from "../Shared";

interface ISavedFilterListProps {
  filterMode: FilterMode;
  onSetFilter: () => void;
}

interface ISavedFilter {
  name: string;
}

const SavedFilter: React.FC<ISavedFilter> = ({ name }) => {
  return (
    <Dropdown.Item>
      <span>{name}</span>
      <ButtonGroup>
        <Button
          className="save-button"
          variant="secondary"
          size="sm"
          title="Overwrite"
        >
          <Icon icon="save" />
        </Button>
        <Button
          className="delete-button"
          variant="secondary"
          size="sm"
          title="Delete"
        >
          <Icon icon="times" />
        </Button>
      </ButtonGroup>
    </Dropdown.Item>
  );
};

export const SavedFilterList: React.FC<ISavedFilterListProps> = () =>
  //   {
  //   // filterMode,
  //   // onSetFilter,
  // }
  {
    const [filterName, setFilterName] = useState("");
    const savedFilters = [
      {
        name: "bar",
      },
      {
        name: "foo",
      },
    ];

    return (
      <div>
        <InputGroup>
          <FormControl
            className="bg-secondary text-white border-secondary"
            placeholder="Filter name..."
            value={filterName}
            onChange={(e) => setFilterName(e.target.value)}
          />
          <InputGroup.Append>
            <OverlayTrigger
              placement="top"
              overlay={<Tooltip id="filter-tooltip">Save filter</Tooltip>}
            >
              <Button
                disabled={!filterName}
                variant="secondary"
                onClick={() => {}}
              >
                <Icon icon="save" />
              </Button>
            </OverlayTrigger>
          </InputGroup.Append>
        </InputGroup>
        <ul className="saved-filter-list">
          {savedFilters
            .filter(
              (f) => !filterName || f.name.toLowerCase().includes(filterName)
            )
            .map((f) => (
              <SavedFilter key={f.name} name={f.name} />
            ))}
        </ul>
        <Button className="set-as-default-button" variant="secondary" size="sm">
          Set as default
        </Button>
      </div>
    );
  };
