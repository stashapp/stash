import React from "react";
import { Button, Form } from "react-bootstrap";
import { CriterionModifier } from "src/core/generated-graphql";
import { ModifierCriterion } from "src/models/list-filter/criteria/criterion";
import cx from "classnames";
import { useIntl } from "react-intl";

const defaultOptions = [
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
  CriterionModifier.Equals,
  CriterionModifier.NotEquals,
  CriterionModifier.Includes,
  CriterionModifier.Excludes,
  CriterionModifier.GreaterThan,
  CriterionModifier.LessThan,
  CriterionModifier.Between,
  CriterionModifier.NotBetween,
];

interface IModifierSelect {
  options?: CriterionModifier[];
  value: CriterionModifier;
  onChanged: (m: CriterionModifier) => void;
}

export const ModifierSelectorButtons: React.FC<IModifierSelect> = ({
  options = defaultOptions,
  value,
  onChanged,
}) => {
  const intl = useIntl();

  return (
    <Form.Group className="modifier-options">
      {options.map((m) => (
        <Button
          className={cx("modifier-option", {
            selected: value === m,
          })}
          key={m}
          onClick={() => onChanged(m)}
        >
          {ModifierCriterion.getModifierLabel(intl, m)}
        </Button>
      ))}
    </Form.Group>
  );
};

export const ModifierSelect: React.FC<IModifierSelect> = ({
  options = defaultOptions,
  value,
  onChanged,
}) => {
  const intl = useIntl();

  return (
    <Form.Control
      as="select"
      onChange={(e) => onChanged(e.target.value as CriterionModifier)}
      value={value}
      className="btn-secondary modifier-selector"
    >
      {options.map((m) => (
        <option key={m} value={m}>
          {ModifierCriterion.getModifierLabel(intl, m)}
        </option>
      ))}
    </Form.Control>
  );
};
