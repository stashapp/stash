import React, { FunctionComponent, useState, useEffect } from "react";
import { InputGroup, ButtonGroup, Button, IInputGroupProps, HTMLInputProps, ControlGroup } from "@blueprintjs/core";
import { TextUtils } from "../../utils/text";
import { FIXED, NUMERIC_INPUT } from "@blueprintjs/core/lib/esm/common/classes";

interface IProps {
  disabled?: boolean
  numericValue: number
  onValueChange(valueAsNumber: number): void
  onReset?(): void
}

export const DurationInput: FunctionComponent<HTMLInputProps & IProps> = (props: IProps) => {
  const [value, setValue] = useState<string>(secondsToString(props.numericValue));

  useEffect(() => {
    setValue(secondsToString(props.numericValue));
  }, [props.numericValue]);

  function secondsToString(seconds : number) {
    let ret = TextUtils.secondsToTimestamp(seconds);

    if (ret.startsWith("00:")) {
      ret = ret.substr(3);

      if (ret.startsWith("0")) {
        ret = ret.substr(1);
      }
    }

    return ret;
  }

  function stringToSeconds(v : string) {
    if (!v) {
      return 0;
    }
    
    let splits = v.split(":");

    if (splits.length > 3) {
      return 0;
    }

    let seconds = 0;
    let factor = 1;
    while(splits.length > 0) {
      let thisSplit = splits.pop();
      if (thisSplit == undefined) {
        return 0;
      }

      seconds += factor * parseInt(thisSplit, 10);
      factor *= 60;
    }

    return seconds;
  }

  function increment() {
    let seconds = stringToSeconds(value);
    seconds += 1;
    props.onValueChange(seconds);
  }

  function decrement() {
    let seconds = stringToSeconds(value);
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
        onBlur={() => props.onValueChange(stringToSeconds(value))}
        placeholder="hh:mm:ss"
        rightElement={maybeRenderReset()}
      />
      {renderButtons()}
    </ControlGroup>
  )
};