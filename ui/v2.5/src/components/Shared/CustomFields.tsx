import React, { useEffect, useMemo, useRef, useState } from "react";
import { CollapseButton } from "./CollapseButton";
import { DetailItem } from "./DetailItem";
import { Button, Col, Form, InputGroup, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { cloneDeep } from "@apollo/client/utilities";
import { Icon } from "./Icon";
import { faMinus, faPlus } from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";

export type CustomFieldMap = {
  [key: string]: unknown;
};

interface ICustomFields {
  values: CustomFieldMap;
}

export const CustomFields: React.FC<ICustomFields> = ({ values }) => {
  const intl = useIntl();
  if (Object.keys(values).length === 0) {
    return null;
  }

  return (
    // according to linter rule CSS classes shouldn't use underscores
    <div className="custom-fields">
      <CollapseButton text={intl.formatMessage({ id: "custom_fields.title" })}>
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
  isNew?: boolean;
}> = ({ field, value, onChange, isNew = false }) => {
  const intl = useIntl();
  const [currentField, setCurrentField] = useState(field);
  const [currentValue, setCurrentValue] = useState(value);

  const fieldRef = useRef<HTMLInputElement>(null);
  const valueRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setCurrentField(field);
    setCurrentValue(value);
  }, [field, value]);

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

    // only update on existing fields
    if (!isNew) {
      onChange(currentField, currentValue);
    }
  }

  function onAdd() {
    onChange(currentField, currentValue);
    setCurrentField("");
    setCurrentValue("");
  }

  function onDelete() {
    onChange("", "");
  }

  return (
    <Row className={cx("custom-fields-row", { "custom-fields-new": isNew })}>
      <Col sm={3} xl={2}>
        {isNew ? (
          <Form.Control
            ref={fieldRef}
            className="input-control"
            type="text"
            value={(currentField as string) ?? ""}
            placeholder={intl.formatMessage({ id: "custom_fields.field" })}
            onChange={(event) => setCurrentField(event.currentTarget.value)}
            onBlur={onBlur}
          />
        ) : (
          <Form.Label>{currentField}</Form.Label>
        )}
      </Col>
      <Col sm={9} xl={7}>
        <InputGroup>
          <Form.Control
            ref={valueRef}
            className="input-control"
            type="text"
            value={(currentValue as string) ?? ""}
            placeholder={currentField}
            onChange={(event) => setCurrentValue(event.currentTarget.value)}
            onBlur={onBlur}
          />
          <InputGroup.Append>
            {isNew ? (
              <Button
                className="custom-fields-add"
                variant="success"
                onClick={() => onAdd()}
                disabled={!currentField}
              >
                <Icon icon={faPlus} />
              </Button>
            ) : (
              <Button
                className="custom-fields-remove"
                variant="danger"
                onClick={() => onDelete()}
              >
                <Icon icon={faMinus} />
              </Button>
            )}
          </InputGroup.Append>
        </InputGroup>
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
  const intl = useIntl();

  const fields = useMemo(() => {
    const ret = Object.keys(values);
    ret.sort();
    return ret;
  }, [values]);

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
    <CollapseButton
      className="custom-fields-input"
      text={intl.formatMessage({ id: "custom_fields.title" })}
    >
      <Row>
        <Col xl={12}>
          <Row className="custom-fields-input-header">
            <Form.Label column sm={3} xl={2}>
              <FormattedMessage id="custom_fields.field" />
            </Form.Label>
            <Form.Label column sm={9} xl={7}>
              <FormattedMessage id="custom_fields.value" />
            </Form.Label>
          </Row>
          {fields.map((field) => (
            <CustomFieldInput
              key={field}
              field={field}
              value={values[field]}
              onChange={(newField, newValue) =>
                fieldChanged(field, newField, newValue)
              }
            />
          ))}
          <CustomFieldInput
            field=""
            value=""
            onChange={(field, value) => fieldChanged("", field, value)}
            isNew
          />
        </Col>
      </Row>
    </CollapseButton>
  );
};
