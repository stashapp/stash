import { faCalendar } from "@fortawesome/free-regular-svg-icons";
import React, { forwardRef, useMemo } from "react";
import { Button, InputGroup, Form } from "react-bootstrap";
import ReactDatePicker from "react-datepicker";
import TextUtils from "src/utils/text";
import { Icon } from "./Icon";

import "react-datepicker/dist/react-datepicker.css";
import { useIntl } from "react-intl";
import { PatchComponent } from "src/patch";

interface IProps {
  disabled?: boolean;
  value: string;
  isTime?: boolean;
  onValueChange(value: string): void;
  placeholder?: string;
  error?: string;
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

  const placeholderText = intl.formatMessage({
    id: props.isTime ? "datetime_format" : "date_format",
  });

  return (
    <div>
      <InputGroup hasValidation>
        <Form.Control
          className="date-input text-input"
          disabled={props.disabled}
          value={props.value}
          onChange={(e) => props.onValueChange(e.currentTarget.value)}
          placeholder={
            !props.disabled
              ? props.placeholder
                ? `${props.placeholder} (${placeholderText})`
                : placeholderText
              : undefined
          }
          isInvalid={!!props.error}
        />
        <InputGroup.Append>{maybeRenderButton()}</InputGroup.Append>
        <Form.Control.Feedback type="invalid">
          {props.error}
        </Form.Control.Feedback>
      </InputGroup>
    </div>
  );
};

export const DateInput = PatchComponent("DateInput", _DateInput);
