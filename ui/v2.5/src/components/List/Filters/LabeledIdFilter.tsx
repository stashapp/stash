import React from "react";
import { Form } from "react-bootstrap";
import { FilterSelect, SelectObject } from "src/components/Shared/Select";
import { galleryTitle } from "src/core/galleries";
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
    inputType !== "scenes" &&
    inputType !== "groups" &&
    inputType !== "galleries"
  ) {
    return null;
  }

  function getLabel(i: SelectObject) {
    if (inputType === "galleries") {
      return galleryTitle(i);
    }

    return i.name ?? i.title ?? "";
  }

  function onSelectionChanged(items: SelectObject[]) {
    onValueChanged(
      items.map((i) => ({
        id: i.id,
        label: getLabel(i),
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
