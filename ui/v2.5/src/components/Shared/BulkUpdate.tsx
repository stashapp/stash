import { faBan } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import {
  Button,
  Col,
  Form,
  FormControlProps,
  InputGroup,
  Row,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "./Icon";
import * as FormUtils from "src/utils/form";

interface IBulkUpdateTextInputProps extends Omit<FormControlProps, "value"> {
  valueChanged: (value: string | null | undefined) => void;
  value: string | null | undefined;
  unsetDisabled?: boolean;
  as?: React.ElementType;
}

export const BulkUpdateTextInput: React.FC<IBulkUpdateTextInputProps> = ({
  valueChanged,
  unsetDisabled,
  ...props
}) => {
  const intl = useIntl();

  const value = props.value === null ? "" : props.value ?? undefined;
  const unset = value === undefined;

  const placeholderValue = unset
    ? `<${intl.formatMessage({ id: "existing_value" })}>`
    : value === ""
    ? `<${intl.formatMessage({ id: "empty_value" })}>`
    : undefined;

  return (
    <InputGroup className="bulk-update-text-input">
      <Form.Control
        {...props}
        className="text-input"
        type="text"
        as={props.as}
        value={value ?? ""}
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

export const BulkUpdateFormGroup: React.FC<{
  name: string;
  messageId?: string;
  inline?: boolean;
}> = ({ name, messageId = name, inline = true, children }) => {
  if (inline) {
    return (
      <Form.Group controlId={name} data-field={name} as={Row}>
        {FormUtils.renderLabel({
          title: <FormattedMessage id={messageId} />,
        })}
        <Col xs={9}>{children}</Col>
      </Form.Group>
    );
  }

  return (
    <Form.Group controlId={name} data-field={name}>
      <Form.Label>
        <FormattedMessage id={messageId} />
      </Form.Label>
      {children}
    </Form.Group>
  );
};
