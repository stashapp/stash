import React, { useState, useEffect } from "react";
import { Button, ButtonGroup, InputGroup, Form } from 'react-bootstrap';
import { Icon } from 'src/components/Shared'
import { TextUtils } from "src/utils";

interface IProps {
  disabled?: boolean
  numericValue: number
  onValueChange(valueAsNumber: number): void
  onReset?(): void
}

export const DurationInput: React.FC<IProps> = (props: IProps) => {
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
      if (thisSplit === undefined) {
        return 0;
      }

      let thisInt = parseInt(thisSplit, 10);
      if (isNaN(thisInt)) {
        return 0;
      }

      seconds += factor * thisInt;
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
      >
        <Button
          disabled={props.disabled}
          onClick={() => increment()}
        >
          <Icon icon="chevron-up" />
        </Button>
        <Button
          disabled={props.disabled}
          onClick={() => decrement()}
        >
          <Icon icon="chevron-down" />
        </Button>
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
          onClick={() => onReset()}
        >
          <Icon icon="clock" />
        </Button>
      )
    }
  }

  return (
    <Form.Group>
      <InputGroup>
        <Form.Control
          disabled={props.disabled}
          value={value}
          onChange={(e : any) => setValue(e.target.value)}
          onBlur={() => props.onValueChange(stringToSeconds(value))}
          placeholder="hh:mm:ss"
        />
        <InputGroup.Append>
          { maybeRenderReset() }
          { renderButtons() }
        </InputGroup.Append>
      </InputGroup>
    </Form.Group>
  )
};
