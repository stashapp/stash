import React from "react";
import { BooleanSetting } from "../Inputs";
import { PatchComponent } from "src/patch";

interface IItem {
  id: string;
  headingID: string;
}

interface ICheckboxGroupProps {
  groupId: string;
  items: IItem[];
  checkedIds?: string[];
  onChange?: (ids: string[]) => void;
}

export const CheckboxGroup: React.FC<ICheckboxGroupProps> = PatchComponent(
  "CheckboxGroup",
  ({ groupId, items, checkedIds = [], onChange }) => {
    function generateId(itemId: string) {
      return `${groupId}-${itemId}`;
    }

    return (
      <>
        {items.map(({ id, headingID }) => (
          <BooleanSetting
            key={id}
            id={generateId(id)}
            headingID={headingID}
            checked={checkedIds.includes(id)}
            onChange={(v) => {
              if (v) {
                onChange?.(
                  items
                    .map((item) => item.id)
                    .filter(
                      (itemId) =>
                        generateId(itemId) === generateId(id) ||
                        checkedIds.includes(itemId)
                    )
                );
              } else {
                onChange?.(
                  items
                    .map((item) => item.id)
                    .filter(
                      (itemId) =>
                        generateId(itemId) !== generateId(id) &&
                        checkedIds.includes(itemId)
                    )
                );
              }
            }}
          />
        ))}
      </>
    );
  }
);
