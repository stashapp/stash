import React, { useState, useEffect } from "react";
import { Button, ButtonGroup, InputGroup, Form } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { DurationUtils } from "src/utils";

interface IProps {
  disabled?: boolean;
  numericValue: number;
  onValueChange(valueAsNumber: number): void;
  onReset?(): void;
}

export const DurationInput: React.FC<IProps> = (props: IProps) => {
  const [value, setValue] = useState<string>(
    DurationUtils.secondsToString(props.numericValue)
  );

  useEffect(() => {
    setValue(DurationUtils.secondsToString(props.numericValue));
  }, [props.numericValue]);

  function increment() {
    let seconds = DurationUtils.stringToSeconds(value);
    seconds += 1;
    props.onValueChange(seconds);
  }

  function decrement() {
    let seconds = DurationUtils.stringToSeconds(value);
    seconds -= 1;
    props.onValueChange(seconds);
  }

  function renderButtons() {
    return (
      <ButtonGroup vertical>
        <Button disabled={props.disabled} onClick={() => increment()}>
          <Icon icon="chevron-up" />
        </Button>
        <Button disabled={props.disabled} onClick={() => decrement()}>
          <Icon icon="chevron-down" />
        </Button>
      </ButtonGroup>
    );
  }

  function onReset() {
    if (props.onReset) {
      props.onReset();
    }
  }

  function maybeRenderReset() {
    if (props.onReset) {
      return (
        <Button onClick={() => onReset()}>
          <Icon icon="clock" />
        </Button>
      );
    }
  }

  return (
    <Form.Group>
      <InputGroup>
        <Form.Control
          disabled={props.disabled}
          value={value}
          onChange={(e: any) => setValue(e.target.value)}
          onBlur={() =>
            props.onValueChange(DurationUtils.stringToSeconds(value))
          }
          placeholder="hh:mm:ss"
        />
        <InputGroup.Append>
          {maybeRenderReset()}
          {renderButtons()}
        </InputGroup.Append>
      </InputGroup>
    </Form.Group>
  );
};
