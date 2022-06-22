import {
  faChevronDown,
  faChevronUp,
  faClock,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState, useEffect } from "react";
import { Button, ButtonGroup, InputGroup, Form } from "react-bootstrap";
import Icon from "src/components/Shared/Icon";
import { DurationUtils } from "src/utils";

interface IProps {
  disabled?: boolean;
  numericValue: number | undefined;
  mandatory?: boolean;
  onValueChange(
    valueAsNumber: number | undefined,
    valueAsString?: string
  ): void;
  onReset?(): void;
  className?: string;
  placeholder?: string;
}

export const DurationInput: React.FC<IProps> = (props: IProps) => {
  const [value, setValue] = useState<string | undefined>(
    props.numericValue !== undefined
      ? DurationUtils.secondsToString(props.numericValue)
      : undefined
  );

  useEffect(() => {
    if (props.numericValue !== undefined || props.mandatory) {
      setValue(DurationUtils.secondsToString(props.numericValue ?? 0));
    } else {
      setValue(undefined);
    }
  }, [props.numericValue, props.mandatory]);

  function increment() {
    if (value === undefined) {
      return;
    }

    let seconds = DurationUtils.stringToSeconds(value);
    seconds += 1;
    props.onValueChange(seconds, DurationUtils.secondsToString(seconds));
  }

  function decrement() {
    if (value === undefined) {
      return;
    }

    let seconds = DurationUtils.stringToSeconds(value);
    seconds -= 1;
    props.onValueChange(seconds, DurationUtils.secondsToString(seconds));
  }

  function renderButtons() {
    if (!props.disabled) {
      return (
        <ButtonGroup vertical>
          <Button
            variant="secondary"
            className="duration-button"
            disabled={props.disabled}
            onClick={() => increment()}
          >
            <Icon icon={faChevronUp} />
          </Button>
          <Button
            variant="secondary"
            className="duration-button"
            disabled={props.disabled}
            onClick={() => decrement()}
          >
            <Icon icon={faChevronDown} />
          </Button>
        </ButtonGroup>
      );
    }
  }

  function onReset() {
    if (props.onReset) {
      props.onReset();
    }
  }

  function maybeRenderReset() {
    if (props.onReset) {
      return (
        <Button variant="secondary" onClick={onReset}>
          <Icon icon={faClock} />
        </Button>
      );
    }
  }

  return (
    <div className={`duration-input ${props.className}`}>
      <InputGroup>
        <Form.Control
          className="duration-control text-input"
          disabled={props.disabled}
          value={value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setValue(e.currentTarget.value)
          }
          onBlur={() => {
            if (props.mandatory || (value !== undefined && value !== "")) {
              props.onValueChange(DurationUtils.stringToSeconds(value), value);
            } else {
              props.onValueChange(undefined);
            }
          }}
          placeholder={
            !props.disabled
              ? props.placeholder
                ? `${props.placeholder} (hh:mm:ss)`
                : "hh:mm:ss"
              : undefined
          }
        />
        <InputGroup.Append>
          {maybeRenderReset()}
          {renderButtons()}
        </InputGroup.Append>
      </InputGroup>
    </div>
  );
};
