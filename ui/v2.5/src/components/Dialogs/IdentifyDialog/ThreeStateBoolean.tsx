import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";

interface IThreeStateBoolean {
  id: string;
  value: boolean | undefined;
  setValue: (v: boolean | undefined) => void;
  allowUndefined?: boolean;
  label?: React.ReactNode;
  disabled?: boolean;
  defaultValue?: boolean;
}

export const ThreeStateBoolean: React.FC<IThreeStateBoolean> = ({
  id,
  value,
  setValue,
  allowUndefined = true,
  label,
  disabled,
  defaultValue,
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

  function getBooleanText(v: boolean) {
    if (v) {
      return intl.formatMessage({ id: "true" });
    }
    return intl.formatMessage({ id: "false" });
  }

  function getButtonText(v: boolean | undefined) {
    if (v === undefined) {
      const defaultVal =
        defaultValue !== undefined ? (
          <span className="default-value">
            {" "}
            ({getBooleanText(defaultValue)})
          </span>
        ) : (
          ""
        );
      return (
        <span>
          {intl.formatMessage({ id: "use_default" })}
          {defaultVal}
        </span>
      );
    }

    return getBooleanText(v);
  }

  function renderModeButton(v: boolean | undefined) {
    return (
      <Form.Check
        type="radio"
        id={`${id}-value-${v ?? "undefined"}`}
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
