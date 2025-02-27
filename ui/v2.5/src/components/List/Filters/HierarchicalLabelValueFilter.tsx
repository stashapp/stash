import React from "react";
import { Form } from "react-bootstrap";
import { defineMessages, MessageDescriptor, useIntl } from "react-intl";
import { FilterSelect, SelectObject } from "src/components/Shared/Select";
import { ModifierCriterion } from "src/models/list-filter/criteria/criterion";
import { IHierarchicalLabelValue } from "src/models/list-filter/types";
import { NumberField } from "src/utils/form";

interface IHierarchicalLabelValueFilterProps {
  criterion: ModifierCriterion<IHierarchicalLabelValue>;
  onValueChanged: (value: IHierarchicalLabelValue) => void;
}

export const HierarchicalLabelValueFilter: React.FC<
  IHierarchicalLabelValueFilterProps
> = ({ criterion, onValueChanged }) => {
  const criterionOption = criterion.modifierCriterionOption();
  const { type, inputType } = criterionOption;

  const intl = useIntl();

  if (
    inputType !== "studios" &&
    inputType !== "tags" &&
    inputType !== "scene_tags" &&
    inputType !== "performer_tags" &&
    inputType !== "groups"
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
    if (inputType === "groups") {
      return "include-sub-groups";
    }
    if (type === "children") {
      return "include-parent-tags";
    }
    console.log(inputType);
    return "include-sub-tags";
  }

  function criterionOptionTypeToIncludeUIString(): MessageDescriptor {
    let id: string;
    if (inputType === "studios") {
      id = "include_sub_studios";
    } else if (inputType === "groups") {
      id = "include_sub_groups";
    } else if (type === "children") {
      id = "include_parent_tags";
    } else {
      id = "include_sub_tags";
    }

    return {
      id,
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
          <NumberField
            className="btn-secondary"
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
