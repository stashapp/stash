import React, { useRef, useState } from "react";
import { CollapseButton } from "./CollapseButton";
import { DetailItem } from "./DetailItem";
import { Col, Form, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { cloneDeep } from "@apollo/client/utilities";

export type CustomFieldMap = {
  [key: string]: unknown;
};

interface ICustomFields {
  values: CustomFieldMap;
}

export const CustomFields: React.FC<ICustomFields> = ({ values }) => {
  if (Object.keys(values).length === 0) {
    return null;
  }

  return (
    // according to linter rule CSS classes shouldn't use underscores
    <div className="custom-fields">
      <CollapseButton text="Custom Fields">
        {Object.entries(values).map(([key, value]) => (
          <DetailItem
            key={key}
            id={`custom-field-${key}`}
            label={key}
            value={value}
            fullWidth={true}
          />
        ))}
      </CollapseButton>
    </div>
  );
};

const CustomFieldInput: React.FC<{
  field: string;
  value: unknown;
  onChange: (field: string, value: unknown) => void;
}> = ({ field, value, onChange }) => {
  const intl = useIntl();
  const [currentField, setCurrentField] = useState(field);
  const [currentValue, setCurrentValue] = useState(value);

  const fieldRef = useRef<HTMLInputElement>(null);
  const valueRef = useRef<HTMLInputElement>(null);

  function onBlur(event: React.FocusEvent<HTMLInputElement>) {
    // don't fire an update if the user is tabbing between fields
    // this prevents focus being stolen from the field
    if (
      currentField &&
      (event.relatedTarget === valueRef.current ||
        event.relatedTarget === fieldRef.current)
    ) {
      return;
    }

    onChange(currentField, currentValue);
  }

  return (
    <Row className="custom-fields-row">
      <Col xs={6}>
        <Form.Control
          ref={fieldRef}
          className="input-control"
          type="text"
          value={(currentField as string) ?? ""}
          placeholder={intl.formatMessage({ id: "field" })}
          onChange={(event) => setCurrentField(event.currentTarget.value)}
          onBlur={onBlur}
        />
      </Col>
      <Col xs={6}>
        <Form.Control
          ref={valueRef}
          className="input-control"
          type="text"
          value={(currentValue as string) ?? ""}
          placeholder={currentField}
          onChange={(event) => setCurrentValue(event.currentTarget.value)}
          onBlur={onBlur}
        />
      </Col>
    </Row>
  );
};

interface ICustomFieldsInput {
  values: CustomFieldMap;
  onChange: (values: CustomFieldMap) => void;
}

export const CustomFieldsInput: React.FC<ICustomFieldsInput> = ({
  values,
  onChange,
}) => {
  function fieldChanged(
    currentField: string,
    newField: string,
    value: unknown
  ) {
    let newValues = cloneDeep(values);
    delete newValues[currentField];
    if (newField !== "") {
      newValues[newField] = value;
    }
    onChange(newValues);
  }

  return (
    <Row className="custom-fields-input">
      <Col xl={9}>
        <Row className="custom-fields-input-header">
          <Form.Label column xs={6}>
            <FormattedMessage id="field" />
          </Form.Label>
          <Form.Label column xs={6}>
            <FormattedMessage id="value" />
          </Form.Label>
        </Row>
        {Object.entries(values).map(([field, value]) => (
          <CustomFieldInput
            key={field}
            field={field}
            value={value}
            onChange={(newField, newValue) =>
              fieldChanged(field, newField, newValue)
            }
          />
        ))}
        <CustomFieldInput
          key={Object.keys(values).length}
          field=""
          value=""
          onChange={(field, value) => fieldChanged("", field, value)}
        />
      </Col>
    </Row>
  );
};
