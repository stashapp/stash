import React from "react";
import { Dropdown } from "react-bootstrap";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { Icon } from "../Shared/Icon";
import { faEllipsisH } from "@fortawesome/free-solid-svg-icons";

export interface IListOperation {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
  icon?: IconDefinition;
  buttonVariant?: string;
}

interface IListOperationDropdownProps {
  operations: IListOperation[];
}

export const ListOperationDropdown: React.FC<IListOperationDropdownProps> = ({
  operations,
}) => {
  const options = operations
    .filter((o) => {
      if (!o.isDisplayed) {
        return true;
      }

      return o.isDisplayed();
    })
    .map((o) => (
      <Dropdown.Item
        key={o.text}
        className="bg-secondary text-white"
        onClick={o.onClick}
      >
        {o.text}
      </Dropdown.Item>
    ));

  if (options.length > 0) {
    return (
      <Dropdown>
        <Dropdown.Toggle variant="secondary" id="more-menu">
          <Icon icon={faEllipsisH} />
        </Dropdown.Toggle>
        <Dropdown.Menu className="bg-secondary text-white">
          {options}
        </Dropdown.Menu>
      </Dropdown>
    );
  }

  return null;
};
