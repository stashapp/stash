import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import { DurationInput } from "src/components/Shared/DurationInput";
import { INumberValue } from "src/models/list-filter/types";
import { Criterion } from "src/models/list-filter/criteria/criterion";

interface IDurationFilterProps {
  criterion: Criterion<INumberValue>;
  onValueChanged: (value: INumberValue) => void;
}

export const DurationFilter: React.FC<IDurationFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  function onChanged(valueAsNumber: number, property: "value" | "value2") {
    const { value } = criterion;
    value[property] = valueAsNumber;
    onValueChanged(value);
  }

  let equalsControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals
  ) {
    equalsControl = (
      <Form.Group>
        <DurationInput
          numericValue={criterion.value?.value}
          onValueChange={(v: number) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.value" })}
        />
      </Form.Group>
    );
  }

  let lowerControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.GreaterThan ||
    criterion.modifier === CriterionModifier.Between ||
    criterion.modifier === CriterionModifier.NotBetween
  ) {
    lowerControl = (
      <Form.Group>
        <DurationInput
          numericValue={criterion.value?.value}
          onValueChange={(v: number) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.greater_than" })}
        />
      </Form.Group>
    );
  }

  let upperControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.LessThan ||
    criterion.modifier === CriterionModifier.Between ||
    criterion.modifier === CriterionModifier.NotBetween
  ) {
    upperControl = (
      <Form.Group>
        <DurationInput
          numericValue={
            criterion.modifier === CriterionModifier.LessThan
              ? criterion.value?.value
              : criterion.value?.value2
          }
          onValueChange={(v: number) =>
            onChanged(
              v,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          placeholder={intl.formatMessage({ id: "criterion.less_than" })}
        />
      </Form.Group>
    );
  }

  return (
    <>
      {equalsControl}
      {lowerControl}
      {upperControl}
    </>
  );
};
