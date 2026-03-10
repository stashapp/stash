import { faCalendar } from "@fortawesome/free-regular-svg-icons";
import React, { forwardRef, useMemo } from "react";
import { Button, InputGroup, Form } from "react-bootstrap";
import ReactDatePicker from "react-datepicker";
import TextUtils from "src/utils/text";
import { Icon } from "./Icon";

import "react-datepicker/dist/react-datepicker.css";
import { useIntl } from "react-intl";
import { PatchComponent } from "src/patch";
import { faBan } from "@fortawesome/free-solid-svg-icons";

interface IProps {
  groupClassName?: string;
  className?: string;
  disabled?: boolean;
  value: string;
  isTime?: boolean;
  onValueChange(value: string): void;
  placeholder?: string;
  placeholderOverride?: string;
  error?: string;
  append?: React.ReactNode;
}

const ShowPickerButton = forwardRef<
  HTMLButtonElement,
  {
    onClick: (event: React.MouseEvent) => void;
  }
>(({ onClick }, ref) => (
  <Button variant="secondary" onClick={onClick} ref={ref}>
    <Icon icon={faCalendar} />
  </Button>
));

const _DateInput: React.FC<IProps> = (props: IProps) => {
  const intl = useIntl();

  const {
    groupClassName = "date-input-group",
    className = "date-input text-input",
  } = props;

  const date = useMemo(() => {
    const toDate = props.isTime
      ? TextUtils.stringToFuzzyDateTime
      : TextUtils.stringToFuzzyDate;
    if (props.value) {
      const ret = toDate(props.value);
      if (ret && !Number.isNaN(ret.getTime())) {
        return ret;
      }
    }
  }, [props.value, props.isTime]);

  function maybeRenderButton() {
    if (!props.disabled) {
      const dateToString = props.isTime
        ? TextUtils.dateTimeToString
        : TextUtils.dateToString;

      return (
        <ReactDatePicker
          selected={date}
          onChange={(v) => {
            props.onValueChange(v ? dateToString(v) : "");
          }}
          customInput={<ShowPickerButton onClick={() => {}} />}
          showMonthDropdown
          showYearDropdown
          scrollableMonthYearDropdown
          scrollableYearDropdown
          maxDate={new Date()}
          yearDropdownItemNumber={100}
          portalId="date-picker-portal"
          showTimeSelect={props.isTime}
        />
      );
    }
  }

  const formatHint = intl.formatMessage({
    id: props.isTime ? "datetime_format" : "date_format",
  });

  const placeholderText = props.placeholder
    ? `${props.placeholder} (${formatHint})`
    : formatHint;

  return (
    <InputGroup hasValidation className={groupClassName}>
      <Form.Control
        className={className}
        disabled={props.disabled}
        value={props.value}
        onChange={(e) => props.onValueChange(e.currentTarget.value)}
        placeholder={
          !props.disabled
            ? props.placeholderOverride ?? placeholderText
            : undefined
        }
        isInvalid={!!props.error}
      />
      <InputGroup.Append>
        {props.append}
        {maybeRenderButton()}
      </InputGroup.Append>
      <Form.Control.Feedback type="invalid">
        {props.error}
      </Form.Control.Feedback>
    </InputGroup>
  );
};

export const DateInput = PatchComponent("DateInput", _DateInput);

interface IBulkUpdateDateInputProps
  extends Omit<IProps, "onValueChange" | "value"> {
  value: string | undefined;
  valueChanged: (value: string | undefined) => void;
  unsetDisabled?: boolean;
  as?: React.ElementType;
}

export const BulkUpdateDateInput: React.FC<IBulkUpdateDateInputProps> = ({
  valueChanged,
  unsetDisabled,
  ...props
}) => {
  const intl = useIntl();

  const unsetClassName = props.value === undefined ? "unset" : "";

  const unsetButton = !unsetDisabled ? (
    <Button
      variant="secondary"
      onClick={() => valueChanged(undefined)}
      title={intl.formatMessage({ id: "actions.unset" })}
    >
      <Icon icon={faBan} />
    </Button>
  ) : undefined;

  return (
    <DateInput
      {...props}
      value={props.value ?? ""}
      placeholderOverride={
        props.value === undefined
          ? `<${intl.formatMessage({ id: "existing_value" })}>`
          : undefined
      }
      onValueChange={(v) => valueChanged(v)}
      groupClassName="bulk-update-date-input"
      className={`${unsetClassName} input-control date-input`}
      append={unsetButton}
    />
  );
};
