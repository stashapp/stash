import { faMinus } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import { Icon } from "./Icon";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage } from "react-intl";

export interface IPerformerAliasListInputProps {
  value: GQL.PerformerAliasInput[];
  setValue: (value: GQL.PerformerAliasInput[]) => void;
  placeholder?: string;
  className?: string;
  errors?: string;
  errorIdx?: number[];
  readOnly?: boolean;
}

export const PerformerAliasListInput: React.FC<
  IPerformerAliasListInputProps
> = (props) => {
  const values = props.value.concat({ alias: "", ignore_auto_tag: true });

  function valueChanged(
    idx: number,
    field: "alias" | "ignore_auto_tag",
    value: string | boolean
  ) {
    const newValues = props.value.slice();
    if (!newValues[idx]) {
      newValues[idx] = { alias: "", ignore_auto_tag: true };
    }
    newValues[idx] = { ...newValues[idx], [field]: value };

    // if we cleared the last string, delete it from the array entirely
    if (!newValues[idx].alias && idx === newValues.length - 1) {
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
              <Form.Control
                className={`text-input ${
                  props.errorIdx?.includes(i) ? "is-invalid" : ""
                }`}
                value={v.alias}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  valueChanged(i, "alias", e.currentTarget.value)
                }
                placeholder={props.placeholder}
                readOnly={props.readOnly}
              />
              <div className="d-flex align-items-center ml-3 mr-3">
                <Form.Check
                  id={`ignore_auto_tag_${i}`}
                  className="mb-0"
                  custom
                  type="switch"
                  label={
                    <FormattedMessage id="auto_tag" defaultMessage="Auto tag" />
                  }
                  checked={!v.ignore_auto_tag}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                    valueChanged(i, "ignore_auto_tag", !e.currentTarget.checked)
                  }
                  disabled={props.readOnly || i === values.length - 1}
                />
              </div>
              <InputGroup.Append>
                {!props.readOnly && (
                  <Button
                    variant="danger"
                    onClick={() => removeValue(i)}
                    disabled={i === values.length - 1}
                    size="sm"
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
