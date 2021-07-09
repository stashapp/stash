import React, { useRef } from "react";
import { Form } from "react-bootstrap";
import { CriterionModifier } from "../../../core/generated-graphql";
import { INumberValue } from "../../../models/list-filter/types";
import { Criterion } from "../../../models/list-filter/criteria/criterion";

interface IDurationFilterProps {
  criterion: Criterion<INumberValue>;
  onValueChanged: (value: INumberValue) => void;
}

export const NumberFilter: React.FC<IDurationFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const valueStage = useRef<INumberValue>(criterion.value);

  function onChanged(
    event: React.ChangeEvent<HTMLInputElement>,
    property: "value" | "value2"
  ) {
    valueStage.current[property] = parseInt(event.target.value, 10);
  }

  function onBlurInput() {
    onValueChanged(valueStage.current);
  }

  let equalsControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals
  ) {
    equalsControl = (
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          type="number"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          onBlur={onBlurInput}
          defaultValue={criterion.value?.value ?? ""}
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
        <Form.Control
          className="btn-secondary"
          type="number"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          onBlur={onBlurInput}
          defaultValue={criterion.value?.value ?? ""}
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
        <Form.Control
          className="btn-secondary"
          type="number"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, criterion.modifier === CriterionModifier.LessThan? "value" : "value2")
          }
          onBlur={onBlurInput}
          defaultValue={(criterion.modifier === CriterionModifier.LessThan ? criterion.value?.value : criterion.value?.value2) ?? ""}
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
