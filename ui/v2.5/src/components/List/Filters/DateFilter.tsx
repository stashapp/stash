import React, { useEffect } from "react";
import { Button, Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { IDateValue } from "../../../models/list-filter/types";
import { Criterion } from "../../../models/list-filter/criteria/criterion";
import { Icon } from "src/components/Shared/Icon";
import { faCheck } from "@fortawesome/free-solid-svg-icons";

interface IDateFilterProps {
  criterion: Criterion<IDateValue>;
  onValueChanged: (value: IDateValue) => void;
}

export const DateFilter: React.FC<IDateFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  const [value, setValue] = React.useState({ ...criterion.value });

  useEffect(() => {
    setValue({ ...criterion.value });
  }, [criterion.value]);

  function onChanged(
    event: React.ChangeEvent<HTMLInputElement>,
    property: "value" | "value2"
  ) {
    const newValue = event.target.value;
    const valueCopy = { ...value };

    valueCopy[property] = newValue;
    setValue(valueCopy);
  }

  function isValid() {
    if (
      criterion.modifier === CriterionModifier.Between ||
      criterion.modifier === CriterionModifier.NotBetween
    ) {
      return value.value !== undefined && value.value2 !== undefined;
    }

    return true;
  }

  function confirm() {
    onValueChanged(value);
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
          type="text"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
          placeholder={
            intl.formatMessage({ id: "criterion.value" }) + " (YYYY-MM-DD)"
          }
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
          type="text"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
          placeholder={
            intl.formatMessage({ id: "criterion.greater_than" }) +
            " (YYYY-MM-DD)"
          }
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
            intl.formatMessage({ id: "criterion.less_than" }) + " (YYYY-MM-DD)"
          }
        />
      </Form.Group>
    );
  }

  return (
    <>
      {equalsControl}
      {lowerControl}
      {upperControl}
      <Button disabled={!isValid()} onClick={() => confirm()}>
        <Icon icon={faCheck} />
      </Button>
    </>
  );
};
