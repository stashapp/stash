import { faEllipsis } from "@fortawesome/free-solid-svg-icons";
import React, { useMemo } from "react";
import { Button, InputGroup, Form } from "react-bootstrap";
import ReactDatePicker from "react-datepicker";
import TextUtils from "src/utils/text";
import { Icon } from "./Icon";

import "react-datepicker/dist/react-datepicker.css";

interface IProps {
  disabled?: boolean;
  value: string | undefined;
  onValueChange(value: string): void;
  placeholder?: string;
  error?: string;
}

export const DateInput: React.FC<IProps> = (props: IProps) => {
  const date = useMemo(() => {
    if (props.value) {
      const ret = TextUtils.stringToFuzzyDate(props.value);
      if (!ret || isNaN(ret.getTime())) {
        return undefined;
      }

      return ret;
    }
  }, [props.value]);

  function maybeRenderButton() {
    if (!props.disabled) {
      const ShowPickerButton = ({
        onClick,
      }: {
        onClick: (
          event: React.MouseEvent<HTMLButtonElement, MouseEvent>
        ) => void;
      }) => (
        <Button variant="secondary" onClick={onClick}>
          <Icon icon={faEllipsis} />
        </Button>
      );

      return (
        <ReactDatePicker
          selected={date}
          onChange={(v) => {
            props.onValueChange(v ? TextUtils.dateToString(v) : "");
          }}
          customInput={React.createElement(ShowPickerButton)}
          showMonthDropdown
          showYearDropdown
          scrollableMonthYearDropdown
          scrollableYearDropdown
          maxDate={new Date()}
          yearDropdownItemNumber={100}
        />
      );
    }
  }

  return (
    <div>
      <InputGroup hasValidation>
        <Form.Control
          className="date-input text-input"
          disabled={props.disabled}
          value={props.value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            props.onValueChange(e.currentTarget.value)
          }
          placeholder={
            !props.disabled
              ? props.placeholder
                ? `${props.placeholder} ("YYYY-MM-DD")`
                : "YYYY-MM-DD"
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
