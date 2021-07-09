import React from "react";
import { Form } from "react-bootstrap";
import { CriterionModifier } from "../../../core/generated-graphql";
import { DurationInput } from "../../Shared";
import { INumberValue } from "../../../models/list-filter/types";
import { Criterion } from "../../../models/list-filter/criteria/criterion";

interface IDurationFilterProps {
  criterion: Criterion<INumberValue>;
  onValueChanged: (value: INumberValue) => void;
}

export const DurationFilter: React.FC<IDurationFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  function onChanged(
    valueAsNumber: number,
    property: "value" | "value2"
  ) {
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
          numericValue={criterion.modifier === CriterionModifier.LessThan ? criterion.value?.value : criterion.value?.value2}
          onValueChange={(v: number) => onChanged(v, criterion.modifier === CriterionModifier.LessThan ? "value" : "value2")}
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
