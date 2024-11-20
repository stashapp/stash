import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { ITimestampValue } from "../../../models/list-filter/types";
import { ModifierCriterion } from "../../../models/list-filter/criteria/criterion";
import { DateInput } from "src/components/Shared/DateInput";

interface ITimestampFilterProps {
  criterion: ModifierCriterion<ITimestampValue>;
  onValueChanged: (value: ITimestampValue) => void;
}

export const TimestampFilter: React.FC<ITimestampFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  const { value } = criterion;

  function onChanged(newValue: string, property: "value" | "value2") {
    const valueCopy = { ...value };

    valueCopy[property] = newValue;
    onValueChanged(valueCopy);
  }

  let equalsControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals
  ) {
    equalsControl = (
      <Form.Group>
        <DateInput
          value={value?.value ?? ""}
          onValueChange={(v) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.value" })}
          isTime
        />
        {/* <Form.Control
          className="btn-secondary"
          type="text"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
          placeholder={
            intl.formatMessage({ id: "criterion.value" }) +
            " (YYYY-MM-DD HH:MM)"
          }
        /> */}
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
        <DateInput
          value={value?.value ?? ""}
          onValueChange={(v) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.greater_than" })}
          isTime
        />
        {/* <Form.Control
          className="btn-secondary"
          type="text"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
          placeholder={
            intl.formatMessage({ id: "criterion.greater_than" }) +
            " (YYYY-MM-DD HH:MM)"
          }
        /> */}
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
        <DateInput
          value={
            (criterion.modifier === CriterionModifier.LessThan
              ? value?.value
              : value?.value2) ?? ""
          }
          onValueChange={(v) =>
            onChanged(
              v,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          placeholder={intl.formatMessage({ id: "criterion.less_than" })}
          isTime
        />
        {/* <Form.Control
          className="btn-secondary"
          type="text"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(
              e,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          value={
            (criterion.modifier === CriterionModifier.LessThan
              ? value?.value
              : value?.value2) ?? ""
          }
          placeholder={
            intl.formatMessage({ id: "criterion.less_than" }) +
            " (YYYY-MM-DD HH:MM)"
          }
        /> */}
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
