import { faCheck } from "@fortawesome/free-solid-svg-icons";
import React, { useEffect } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import {
  Criterion,
  CriterionValue,
} from "../../../models/list-filter/criteria/criterion";

interface IInputFilterProps {
  criterion: Criterion<CriterionValue>;
  onValueChanged: (value: string) => void;
}

export const InputFilter: React.FC<IInputFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const [value, setValue] = React.useState<string>(
    criterion.value ? criterion.value.toString() : ""
  );

  useEffect(() => {
    setValue(criterion.value ? criterion.value.toString() : "");
  }, [criterion.value]);

  function onChanged(event: React.ChangeEvent<HTMLInputElement>) {
    setValue(event.target.value);
  }

  function onConfirm() {
    onValueChanged(value);
  }

  return (
    <Form.Group>
      <InputGroup>
        <Form.Control
          className="btn-secondary"
          type={criterion.criterionOption.inputType}
          onChange={onChanged}
          value={value}
        />
        <InputGroup.Append>
          <Button
            disabled={!value}
            variant="primary"
            onClick={() => {
              onConfirm();
            }}
          >
            <Icon icon={faCheck} />
          </Button>
        </InputGroup.Append>
      </InputGroup>
    </Form.Group>
  );
};
