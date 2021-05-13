import React from "react";
import { Form } from "react-bootstrap";

interface IItem {
  id: string;
  label: string;
}

interface ICheckboxGroupProps {
  groupId: string;
  items: IItem[];
  checkedIds?: string[];
  onChange?: (ids: string[]) => void;
}

export const CheckboxGroup: React.FC<ICheckboxGroupProps> = ({
  groupId,
  items,
  checkedIds = [],
  onChange,
}) => {
  function generateId(itemId: string) {
    return `${groupId}-${itemId}`;
  }

  return (
    <>
      {items.map(({ id, label }) => (
        <Form.Check
          key={id}
          type="checkbox"
          id={generateId(id)}
          label={label}
          checked={checkedIds.includes(id)}
          onChange={(event) => {
            const target = event.currentTarget;
            if (target.checked) {
              onChange?.(
                items
                  .map((item) => item.id)
                  .filter(
                    (itemId) =>
                      generateId(itemId) === target.id ||
                      checkedIds.includes(itemId)
                  )
              );
            } else {
              onChange?.(
                items
                  .map((item) => item.id)
                  .filter(
                    (itemId) =>
                      generateId(itemId) !== target.id &&
                      checkedIds.includes(itemId)
                  )
              );
            }
          }}
        />
      ))}
    </>
  );
};
