import React from "react";
import { Form } from "react-bootstrap";
import {
  Criterion,
  CriterionValue,
} from "../../../models/list-filter/criteria/criterion";
import { useDebouncedState } from "src/hooks/debounce";

interface IInputFilterProps {
  criterion: Criterion<CriterionValue>;
  onValueChanged: (value: string) => void;
}

export const InputFilter: React.FC<IInputFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const [value, setValue] = useDebouncedState<string>(
    criterion.value ? criterion.value.toString() : "",
    onValueChanged
  );

  return (
    <>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          type={criterion.criterionOption.inputType}
          onChange={(e) => setValue(e.target.value)}
          value={value}
        />
      </Form.Group>
    </>
  );
};
