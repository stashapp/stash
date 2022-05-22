import { faMinus, faPlus } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import Icon from "src/components/Shared/Icon";

interface IStringListInputProps {
  value: string[];
  setValue: (value: string[]) => void;
  defaultNewValue?: string;
  className?: string;
  errors?: string;
}

export const StringListInput: React.FC<IStringListInputProps> = (props) => {
  function valueChanged(idx: number, value: string) {
    const newValues = props.value.map((v, i) => {
      const ret = idx !== i ? v : value;
      return ret;
    });
    props.setValue(newValues);
  }

  function removeValue(idx: number) {
    const newValues = props.value.filter((_v, i) => i !== idx);

    props.setValue(newValues);
  }

  function addValue() {
    const newValues = props.value.concat(props.defaultNewValue ?? "");

    props.setValue(newValues);
  }

  return (
    <>
      <div className={`string-list-input ${props.errors ? "is-invalid" : ""}`}>
        {props.value && props.value.length > 0 && (
          <Form.Group>
            {props.value &&
              props.value.map((v, i) => (
                // eslint-disable-next-line react/no-array-index-key
                <InputGroup className={props.className} key={i}>
                  <Form.Control
                    className="text-input"
                    value={v}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      valueChanged(i, e.currentTarget.value)
                    }
                  />
                  <InputGroup.Append>
                    <Button variant="danger" onClick={() => removeValue(i)}>
                      <Icon icon={faMinus} />
                    </Button>
                  </InputGroup.Append>
                </InputGroup>
              ))}
          </Form.Group>
        )}
        <Button className="minimal" size="sm" onClick={() => addValue()}>
          <Icon icon={faPlus} />
        </Button>
      </div>
      <div className="invalid-feedback">{props.errors}</div>
    </>
  );
};
