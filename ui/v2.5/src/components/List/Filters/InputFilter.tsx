import React from "react";
import { Form } from "react-bootstrap";
import {
  ModifierCriterion,
  CriterionValue,
} from "../../../models/list-filter/criteria/criterion";

interface IInputFilterProps {
  criterion: ModifierCriterion<CriterionValue>;
  onValueChanged: (value: string) => void;
}

export const InputFilter: React.FC<IInputFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  function onChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged(event.target.value);
  }

  return (
    <>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          type={criterion.modifierCriterionOption().inputType}
          onChange={onChanged}
          value={criterion.value ? criterion.value.toString() : ""}
        />
      </Form.Group>
    </>
  );
};
