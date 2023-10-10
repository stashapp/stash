import React from "react";
import { Form } from "react-bootstrap";
import { FilterSelect, SelectObject } from "src/components/Shared/Select";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { ILabeledId } from "src/models/list-filter/types";

interface ILabeledIdFilterProps {
  criterion: Criterion<ILabeledId[]>;
  onValueChanged: (value: ILabeledId[]) => void;
}

export const LabeledIdFilter: React.FC<ILabeledIdFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const { criterionOption } = criterion;
  const { inputType } = criterionOption;

  if (
    inputType !== "performers" &&
    inputType !== "studios" &&
    inputType !== "scene_tags" &&
    inputType !== "performer_tags" &&
    inputType !== "tags" &&
    inputType !== "movies"
  ) {
    return null;
  }

  function onSelectionChanged(items: SelectObject[]) {
    onValueChanged(
      items.map((i) => ({
        id: i.id,
        label: i.name ?? i.title ?? "",
      }))
    );
  }

  return (
    <Form.Group>
      <FilterSelect
        type={inputType}
        isMulti
        onSelect={onSelectionChanged}
        ids={criterion.value.map((labeled) => labeled.id)}
        menuPortalTarget={document.body}
      />
    </Form.Group>
  );
};
