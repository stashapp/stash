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

  const unsetClassName = props.value === undefined ? "unset" : "";

  return (
    <InputGroup className={`bulk-update-text-input ${unsetClassName}`}>
      <Form.Control
        {...props}
        className="input-control"
        type="text"
        as={props.as}
        value={props.value ?? ""}
        placeholder={
          props.value === undefined
            ? `<${intl.formatMessage({ id: "existing_value" })}>`
            : undefined
        }
        onChange={(event) => valueChanged(event.currentTarget.value)}
      />
      {!unsetDisabled ? (
        <Button
          variant="secondary"
          onClick={() => valueChanged(undefined)}
          title={intl.formatMessage({ id: "actions.unset" })}
        >
          <Icon icon={faBan} />
        </Button>
      ) : undefined}
    </InputGroup>
  );
};
