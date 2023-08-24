import React from "react";
import { Form } from "react-bootstrap";
import { defineMessages, MessageDescriptor, useIntl } from "react-intl";
import { FilterSelect, SelectObject } from "src/components/Shared/Select";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { IHierarchicalLabelValue } from "src/models/list-filter/types";

interface IHierarchicalLabelValueFilterProps {
  criterion: Criterion<IHierarchicalLabelValue>;
  onValueChanged: (value: IHierarchicalLabelValue) => void;
}

export const HierarchicalLabelValueFilter: React.FC<
  IHierarchicalLabelValueFilterProps
> = ({ criterion, onValueChanged }) => {
  const { criterionOption } = criterion;
  const { type, inputType } = criterionOption;

  const intl = useIntl();

  if (
    inputType !== "studios" &&
    inputType !== "tags" &&
    inputType !== "scene_tags" &&
    inputType !== "performer_tags"
  ) {
    return null;
  }

  const messages = defineMessages({
    studio_depth: {
      id: "studio_depth",
      defaultMessage: "Levels (empty for all)",
    },
  });

  function onSelectionChanged(items: SelectObject[]) {
    const { value } = criterion;
    value.items = items.map((i) => ({
      id: i.id,
      label: i.name ?? i.title ?? "",
    }));
    onValueChanged(value);
  }

  function onDepthChanged(depth: number) {
    const { value } = criterion;
    value.depth = depth;
    onValueChanged(value);
  }

  function criterionOptionTypeToIncludeID(): string {
    if (inputType === "studios") {
      return "include-sub-studios";
    }
    if (type === "children") {
      return "include-parent-tags";
    }
    return "include-sub-tags";
  }

  function criterionOptionTypeToIncludeUIString(): MessageDescriptor {
    const optionType =
      inputType === "studios"
        ? "include_sub_studios"
        : type === "children"
        ? "include_parent_tags"
        : "include_sub_tags";
    return {
      id: optionType,
    };
  }

  return (
    <>
      <Form.Group>
        <FilterSelect
          type={inputType}
          isMulti
          onSelect={onSelectionChanged}
          ids={criterion.value.items.map((labeled) => labeled.id)}
          menuPortalTarget={document.body}
        />
      </Form.Group>

      <Form.Group>
        <Form.Check
          id={criterionOptionTypeToIncludeID()}
          checked={criterion.value.depth !== 0}
          label={intl.formatMessage(criterionOptionTypeToIncludeUIString())}
          onChange={() => onDepthChanged(criterion.value.depth !== 0 ? 0 : -1)}
        />
      </Form.Group>

      {criterion.value.depth !== 0 && (
        <Form.Group>
          <Form.Control
            className="btn-secondary"
            type="number"
            placeholder={intl.formatMessage(messages.studio_depth)}
            onChange={(e) =>
              onDepthChanged(e.target.value ? parseInt(e.target.value, 10) : -1)
            }
            defaultValue={
              criterion.value && criterion.value.depth !== -1
                ? criterion.value.depth
                : ""
            }
            min="1"
          />
        </Form.Group>
      )}
    </>
  );
};
