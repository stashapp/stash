import { faBan } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button, Form, FormControlProps, InputGroup } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "./Icon";

interface IBulkUpdateTextInputProps extends FormControlProps {
  valueChanged: (value: string | undefined) => void;
  unsetDisabled?: boolean;
  as?: React.ElementType;
}

export const BulkUpdateTextInput: React.FC<IBulkUpdateTextInputProps> = ({
  valueChanged,
  unsetDisabled,
  ...props
}) => {
  const intl = useIntl();

  const unset = props.value === undefined;

  const placeholderValue =
    props.value === undefined
      ? `<${intl.formatMessage({ id: "existing_value" })}>`
      : props.value === ""
      ? `<${intl.formatMessage({ id: "empty_value" })}>`
      : undefined;

  return (
    <InputGroup className="bulk-update-text-input">
      <Form.Control
        {...props}
        className="text-input"
        type="text"
        as={props.as}
        value={props.value ?? ""}
        placeholder={placeholderValue}
        onChange={(event) => valueChanged(event.currentTarget.value)}
      />
      <InputGroup.Append>
        {!unsetDisabled ? (
          <Button
            variant="secondary"
            onClick={() => valueChanged(undefined)}
            title={intl.formatMessage({ id: "actions.unset" })}
            disabled={unset}
          >
            <Icon icon={faBan} />
          </Button>
        ) : undefined}
      </InputGroup.Append>
    </InputGroup>
  );
};
