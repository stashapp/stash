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
  if (
    criterion.criterionOption.type !== "performers" &&
    criterion.criterionOption.type !== "studios" &&
    criterion.criterionOption.type !== "parent_studios" &&
    criterion.criterionOption.type !== "tags" &&
    criterion.criterionOption.type !== "sceneTags" &&
    criterion.criterionOption.type !== "performerTags" &&
    criterion.criterionOption.type !== "parentTags" &&
    criterion.criterionOption.type !== "childTags" &&
    criterion.criterionOption.type !== "movies"
  )
    return null;

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
        type={criterion.criterionOption.type}
        isMulti
        onSelect={onSelectionChanged}
        ids={criterion.value.map((labeled) => labeled.id)}
        menuPortalTarget={document.body}
      />
    </Form.Group>
  );
};
