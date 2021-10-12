import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";

interface IThreeStateBoolean {
  value: boolean | undefined;
  setValue: (v: boolean | undefined) => void;
  allowUndefined?: boolean;
  label?: React.ReactNode;
  disabled?: boolean;
}

export const ThreeStateBoolean: React.FC<IThreeStateBoolean> = ({
  value,
  setValue,
  allowUndefined = true,
  label,
  disabled,
}) => {
  const intl = useIntl();

  if (!allowUndefined) {
    return (
      <Form.Check
        disabled={disabled}
        checked={value}
        label={label}
        onChange={() => setValue(!value)}
      />
    );
  }

  function getButtonText(v: boolean | undefined) {
    if (v === undefined) {
      return intl.formatMessage({ id: "use_default" });
    }
    if (v) {
      return intl.formatMessage({ id: "true" });
    }
    return intl.formatMessage({ id: "false" });
  }

  function renderModeButton(v: boolean | undefined) {
    return (
      <Form.Check
        type="radio"
        id={`value-${v ?? "undefined"}`}
        checked={value === v}
        onChange={() => setValue(v)}
        disabled={disabled}
        label={getButtonText(v)}
      />
    );
  }

  return (
    <Form.Group>
      <h6>{label}</h6>
      <Form.Group>
        {renderModeButton(undefined)}
        {renderModeButton(false)}
        {renderModeButton(true)}
      </Form.Group>
    </Form.Group>
  );
};
