import {
  faChevronDown,
  faChevronUp,
  faClock,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, ButtonGroup, InputGroup, Form } from "react-bootstrap";
import { Icon } from "./Icon";
import DurationUtils from "src/utils/duration";

interface IProps {
  disabled?: boolean;
  value: number | undefined;
  setValue(value: number | undefined): void;
  onReset?(): void;
  className?: string;
  placeholder?: string;
}

export const DurationInput: React.FC<IProps> = ({
  disabled,
  value,
  setValue,
  onReset,
  className,
  placeholder,
}) => {
  const [tmpValue, setTmpValue] = useState<string>();

  function onChange(e: React.ChangeEvent<HTMLInputElement>) {
    setTmpValue(e.currentTarget.value);
  }

  function onBlur() {
    if (tmpValue !== undefined) {
      setValue(DurationUtils.stringToSeconds(tmpValue));
      setTmpValue(undefined);
    }
  }

  function increment() {
    setTmpValue(undefined);
    setValue((value ?? 0) + 1);
  }

  function decrement() {
    setTmpValue(undefined);
    setValue((value ?? 0) - 1);
  }

  function renderButtons() {
    if (!disabled) {
      return (
        <ButtonGroup vertical>
          <Button
            variant="secondary"
            className="duration-button"
            onClick={() => increment()}
          >
            <Icon icon={faChevronUp} />
          </Button>
          <Button
            variant="secondary"
            className="duration-button"
            onClick={() => decrement()}
          >
            <Icon icon={faChevronDown} />
          </Button>
        </ButtonGroup>
      );
    }
  }

  function maybeRenderReset() {
    if (onReset) {
      return (
        <Button variant="secondary" onClick={() => onReset()}>
          <Icon icon={faClock} />
        </Button>
      );
    }
  }

  let inputValue = "";
  if (tmpValue !== undefined) {
    inputValue = tmpValue;
  } else if (value !== undefined) {
    inputValue = DurationUtils.secondsToString(value);
  }

  return (
    <div className={`duration-input ${className}`}>
      <InputGroup>
        <Form.Control
          className="duration-control text-input"
          disabled={disabled}
          value={inputValue}
          onChange={onChange}
          onBlur={onBlur}
          placeholder={placeholder ? `${placeholder} (hh:mm:ss)` : "hh:mm:ss"}
        />
        <InputGroup.Append>
          {maybeRenderReset()}
          {renderButtons()}
        </InputGroup.Append>
      </InputGroup>
    </div>
  );
};
