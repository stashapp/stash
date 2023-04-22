import React, { useMemo } from "react";
import { Form } from "react-bootstrap";
import {
  Criterion,
  CriterionValue,
} from "../../../models/list-filter/criteria/criterion";

interface IOptionsFilterProps {
  criterion: Criterion<CriterionValue>;
  onValueChanged: (value: CriterionValue) => void;
}

export const OptionsFilter: React.FC<IOptionsFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  function onChanged(event: React.ChangeEvent<HTMLSelectElement>) {
    onValueChanged(event.target.value);
  }

  const options = useMemo(() => {
    const ret = criterion.criterionOption.options?.slice() ?? [];

    ret.unshift("");

    return ret;
  }, [criterion.criterionOption.options]);

  return (
    <Form.Group>
      <Form.Control
        as="select"
        onChange={onChanged}
        value={criterion.value.toString()}
        className="btn-secondary"
      >
        {options.map((c) => (
          <option key={c.toString()} value={c.toString()}>
            {c ? c : "---"}
          </option>
        ))}
      </Form.Control>
    </Form.Group>
  );
};
