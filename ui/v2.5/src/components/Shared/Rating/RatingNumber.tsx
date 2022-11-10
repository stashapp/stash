import React, { useState } from "react";

export interface IRatingNumberProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
}

export const RatingNumber: React.FC<IRatingNumberProps> = (
  props: IRatingNumberProps
) => {
  const [input, setInput] = useState<string | "0.0">();
  const [previous, setPrevious] = useState<string | "0.0">();
  const [useValidation, setValidation] = useState<boolean | true>();
  function stepChange() {
    setValidation(false);
  }

  function nonStepChange() {
    setValidation(true);
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    let val = e.target.value;
    if (!useValidation && props.onSetRating != null) {
      e.target.value = Number(val).toFixed(1);
      setInput(Number(val).toFixed(1));
      setPrevious(Number(val).toFixed(1));
      let tempVal = Number(val) * 10;
      props.onSetRating(tempVal != 0 ? tempVal : undefined);
      setValidation(true);
      return;
    }
    const match = /(\d?)(\d?)(.?)((\d)?)/g.exec(val);
    const matchOld = /(\d?)(\d?)(.?)((\d{0,2})?)/g.exec(previous ?? "");

    if (match == null || props.onSetRating == null) {
      return;
    }

    if (match[2] && !(match[2] == "0" && match[1] == "1")) {
      match[2] = "";
    }
    if (match[4] == null || match[4] == "") {
      match[4] = "0";
    }
    let value = match[1] + match[2] + "." + match[4];
    e.target.value = value;
    if (val.length > 0) {
      if (Number(value) > 10) {
        value = "10.0";
      }
      e.target.value = Number(value).toFixed(1);
      setInput(Number(value).toFixed(1));
      setPrevious(Number(value).toFixed(1));
      let tempVal = Number(value) * 10;
      props.onSetRating(tempVal != 0 ? tempVal : undefined);

      // This is a workaround to a missing feature where you can't set cursor position in input numbers.
      // See https://stackoverflow.com/questions/33406169/failed-to-execute-setselectionrange-on-htmlinputelement-the-input-elements
      e.target.type = "text";
      let cursorPosition = 0;
      if (match[2] && !match[4]) {
        cursorPosition = 3;
      } else if (matchOld != null && match[1] !== matchOld[1]) {
        cursorPosition = 2;
      } else if (
        matchOld != null &&
        match[1] === matchOld[1] &&
        match[2] === matchOld[2] &&
        match[4] === matchOld[4]
      ) {
        cursorPosition = 2;
      }
      e.target.setSelectionRange(cursorPosition, cursorPosition);
      e.target.type = "number";
    }
  }

  if (props.disabled) {
    return <text>{Number((props.value ?? 0) / 10).toFixed(1)}</text>;
  } else {
    return (
      <div>
        <input
          className="text-input"
          type="number"
          onMouseDown={stepChange}
          onKeyDown={nonStepChange}
          onChange={handleChange}
          value={input}
          defaultValue={
            props.value == null || props.value == undefined
              ? "0.0"
              : Number(props.value / 10).toFixed(1)
          }
          min="0.0"
          step="0.1"
          max="10"
          style={{ fontSize: "22px", padding: "4px" }}
          placeholder="0.0"
        />
      </div>
    );
  }
};
