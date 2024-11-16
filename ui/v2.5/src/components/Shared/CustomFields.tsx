import React from "react";
import { CollapseButton } from "./CollapseButton";
import { DetailItem } from "./DetailItem";

export type CustomFieldMap = {
  [key: string]: unknown;
};

interface ICustomFields {
  values: CustomFieldMap;
}

export const CustomFields: React.FC<ICustomFields> = ({ values }) => {
  if (Object.keys(values).length === 0) {
    return null;
  }

  return (
    // according to linter rule CSS classes shouldn't use underscores
    <div className="custom-fields">
      <CollapseButton text="Custom Fields">
        {Object.entries(values).map(([key, value]) => (
          <DetailItem
            key={key}
            id={`custom-field-${key}`}
            label={key}
            value={value}
            fullWidth={true}
          />
        ))}
      </CollapseButton>
    </div>
  );
};
