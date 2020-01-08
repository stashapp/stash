import React, { FunctionComponent, useState, useEffect } from "react";
import { InputGroup, ButtonGroup, Button, IInputGroupProps, HTMLInputProps, ControlGroup } from "@blueprintjs/core";
import { DurationUtils } from "../../utils/duration";
import { FIXED, NUMERIC_INPUT } from "@blueprintjs/core/lib/esm/common/classes";

interface IProps {
  disabled?: boolean
  numericValue: number
  onValueChange(valueAsNumber: number): void
  onReset?(): void
}

export const DurationInput: FunctionComponent<HTMLInputProps & IProps> = (props: IProps) => {
  const [value, setValue] = useState<string>(DurationUtils.secondsToString(props.numericValue));

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
      <ButtonGroup
       vertical={true}
       className={FIXED}
      >
        <Button
          icon="chevron-up"
          disabled={props.disabled}
          onClick={() => increment()}
        />
        <Button
          icon="chevron-down"
          disabled={props.disabled}
          onClick={() => decrement()}
        />
      </ButtonGroup>
    )
  }

  function onReset() {
    if (props.onReset) {
      props.onReset();
    }
  }

  function maybeRenderReset() {
    if (props.onReset) {
      return (
        <Button
          icon="time"
          onClick={() => onReset()}
        />
      )
    }
  }

  return (
    <ControlGroup className={NUMERIC_INPUT}>
      <InputGroup
        disabled={props.disabled}
        value={value}
        onChange={(e : any) => setValue(e.target.value)}
        onBlur={() => props.onValueChange(DurationUtils.stringToSeconds(value))}
        placeholder="hh:mm:ss"
        rightElement={maybeRenderReset()}
      />
      {renderButtons()}
    </ControlGroup>
  )
};