import React, { useRef } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { INumberValue } from "../../../models/list-filter/types";
import { Criterion } from "../../../models/list-filter/criteria/criterion";
import { ConfigurationContext } from "../../../hooks/Config";
import {
  convertFromRatingFormat
} from "../../../components/Scenes/SceneDetails/RatingSystem";
import * as GQL from "src/core/generated-graphql";

interface IDurationFilterProps {
  criterion: Criterion<INumberValue>;
  onValueChanged: (value: INumberValue) => void;
  configuration: GQL.ConfigDataFragment | undefined;
}

export const RatingFilter: React.FC<IDurationFilterProps> = ({
  criterion,
  onValueChanged,
  configuration
}) => {
  const intl = useIntl();

  const valueStage = useRef<INumberValue>(criterion.value);

  function onChanged(
    event: React.ChangeEvent<HTMLInputElement>,
    property: "value" | "value2"
  ) {
    const value = parseInt(event.target.value, 10);
    valueStage.current[property] = !Number.isNaN(value) ? convertFromRatingFormat(value, configuration?.interface.ratingSystem) : 0;
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
        <Form.Control
          className="btn-secondary"
          type="number"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          onBlur={onBlurInput}
          defaultValue={criterion.value?.value ?? ""}
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
        <Form.Control
          className="btn-secondary"
          type="number"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(
              e,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          onBlur={onBlurInput}
          defaultValue={
            (criterion.modifier === CriterionModifier.LessThan
              ? criterion.value?.value
              : criterion.value?.value2) ?? ""
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
