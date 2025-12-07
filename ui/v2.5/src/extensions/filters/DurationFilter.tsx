import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import { DurationInput } from "src/components/Shared/DurationInput";
import { INumberValue } from "src/models/list-filter/types";
import { ModifierCriterion } from "src/models/list-filter/criteria/criterion";

interface IDurationFilterProps {
  criterion: ModifierCriterion<INumberValue>;
  onValueChanged: (value: INumberValue) => void;
}

export const DurationFilter: React.FC<IDurationFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  function onChanged(v: number | null, property: "value" | "value2") {
    const { value } = criterion;
    value[property] = v ?? undefined;
    onValueChanged(value);
  }

  function renderTop() {
    let placeholder: string;
    if (
      criterion.modifier === CriterionModifier.GreaterThan ||
      criterion.modifier === CriterionModifier.Between ||
      criterion.modifier === CriterionModifier.NotBetween
    ) {
      placeholder = intl.formatMessage({ id: "criterion.greater_than" });
    } else if (criterion.modifier === CriterionModifier.LessThan) {
      placeholder = intl.formatMessage({ id: "criterion.less_than" });
    } else {
      placeholder = intl.formatMessage({ id: "criterion.value" });
    }

    return (
      <Form.Group>
        <DurationInput
          value={criterion.value?.value}
          setValue={(v) => onChanged(v, "value")}
          placeholder={placeholder}
        />
      </Form.Group>
    );
  }

  function renderBottom() {
    if (
      criterion.modifier !== CriterionModifier.Between &&
      criterion.modifier !== CriterionModifier.NotBetween
    ) {
      return;
    }

    return (
      <Form.Group>
        <DurationInput
          value={criterion.value?.value2}
          setValue={(v) => onChanged(v, "value2")}
          placeholder={intl.formatMessage({ id: "criterion.less_than" })}
        />
      </Form.Group>
    );
  }

  return (
    <>
      {renderTop()}
      {renderBottom()}
    </>
  );
};
