import { faMinus } from "@fortawesome/free-solid-svg-icons";
import React, { ComponentType } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import { Icon } from "./Icon";

interface IListInputComponentProps {
  value: string;
  setValue: (value: string) => void;
  placeholder?: string;
  className?: string;
  readOnly?: boolean;
}

interface IListInputAppendProps {
  value: string;
}

export interface IStringListInputProps {
  value: string[];
  setValue: (value: string[]) => void;
  inputComponent?: ComponentType<IListInputComponentProps>;
  appendComponent?: ComponentType<IListInputAppendProps>;
  placeholder?: string;
  className?: string;
  errors?: string;
  errorIdx?: number[];
  readOnly?: boolean;
}

export const StringInput: React.FC<IListInputComponentProps> = ({
  className,
  placeholder,
  value,
  setValue,
  readOnly = false,
}) => {
  return (
    <Form.Control
      className={`text-input ${className ?? ""}`}
      value={value}
      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
        setValue(e.currentTarget.value)
      }
      placeholder={placeholder}
      readOnly={readOnly}
    />
  );
};

export const StringListInput: React.FC<IStringListInputProps> = (props) => {
  const Input = props.inputComponent ?? StringInput;
  const AppendComponent = props.appendComponent;
  const values = props.value.concat("");

  function valueChanged(idx: number, value: string) {
    const newValues = props.value.slice();
    newValues[idx] = value;

    // if we cleared the last string, delete it from the array entirely
    if (!value && idx === newValues.length - 1) {
      newValues.splice(newValues.length - 1);
    }

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
              <Input
                value={v}
                setValue={(value) => valueChanged(i, value)}
                placeholder={props.placeholder}
                className={props.errorIdx?.includes(i) ? "is-invalid" : ""}
                readOnly={props.readOnly}
              />
              <InputGroup.Append>
                {AppendComponent && <AppendComponent value={v} />}
                {!props.readOnly && (
                  <Button
                    variant="danger"
                    onClick={() => removeValue(i)}
                    disabled={i === values.length - 1}
                  >
                    <Icon icon={faMinus} />
                  </Button>
                )}
              </InputGroup.Append>
            </InputGroup>
          ))}
        </Form.Group>
      </div>
      <div className="invalid-feedback mt-n2">{props.errors}</div>
    </>
  );
};
