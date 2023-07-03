import { faMinus } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import { Icon } from "./Icon";

interface IStringListInputProps {
  value: string[];
  setValue: (value: string[]) => void;
  placeholder?: string;
  className?: string;
  errors?: string;
  errorIdx?: number[];
}

export const StringListInput: React.FC<IStringListInputProps> = (props) => {
  const values = props.value.concat("");

  function valueChanged(idx: number, value: string) {
    const newValues = values
      .map((v, i) => {
        const ret = idx !== i ? v : value;
        return ret;
      })
      .filter((v, i) => i < values.length - 2 || v);
    props.setValue(newValues);
  }

  function removeValue(idx: number) {
    const newValues = props.value.filter((_v, i) => i !== idx);

    props.setValue(newValues);
  }

  return (
    <>
      <div className={`string-list-input ${props.errors ? "is-invalid" : ""}`}>
        <Form.Group>
          {values.map((v, i) => (
            <InputGroup className={props.className} key={i}>
              <Form.Control
                className={`text-input ${
                  props.errorIdx?.includes(i) ? "is-invalid" : ""
                }`}
                value={v}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  valueChanged(i, e.currentTarget.value)
                }
                placeholder={props.placeholder}
              />
              <InputGroup.Append>
                <Button
                  variant="danger"
                  onClick={() => removeValue(i)}
                  disabled={i === values.length - 1}
                >
                  <Icon icon={faMinus} />
                </Button>
              </InputGroup.Append>
            </InputGroup>
          ))}
        </Form.Group>
      </div>
      <div className="invalid-feedback mt-n2">{props.errors}</div>
    </>
  );
};
