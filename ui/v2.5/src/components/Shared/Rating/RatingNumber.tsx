import React, { useRef } from "react";

export interface IRatingNumberProps {
  value: number | null;
  onSetRating?: (value: number | null) => void;
  disabled?: boolean;
}

export const RatingNumber: React.FC<IRatingNumberProps> = (
  props: IRatingNumberProps
) => {
  const text = ((props.value ?? 0) / 10).toFixed(1);
  const useValidation = useRef(true);

  function stepChange() {
    useValidation.current = false;
  }

  function nonStepChange() {
    useValidation.current = true;
  }

  function setCursorPosition(
    target: HTMLInputElement,
    pos: number,
    endPos?: number
  ) {
    // This is a workaround to a missing feature where you can't set cursor position in input numbers.
    // See https://stackoverflow.com/questions/33406169/failed-to-execute-setselectionrange-on-htmlinputelement-the-input-elements
    target.type = "text";

    target.setSelectionRange(pos, endPos ?? pos);
    target.type = "number";
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (!props.onSetRating) {
      return;
    }

    let val = e.target.value;
    if (!useValidation.current) {
      e.target.value = Number(val).toFixed(1);
      const tempVal = Number(val) * 10;
      props.onSetRating(tempVal || null);
      useValidation.current = true;
      return;
    }

    const match = /(\d?)(\d?)(.?)((\d)?)/g.exec(val);
    const matchOld = /(\d?)(\d?)(.?)((\d{0,2})?)/g.exec(text ?? "");

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
      let tempVal = Number(value) * 10;
      props.onSetRating(tempVal || null);

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

      setCursorPosition(e.target, cursorPosition);
    }
  }

  if (props.disabled) {
    return (
      <div className="rating-number disabled">
        <span>{Number((props.value ?? 0) / 10).toFixed(1)}</span>
      </div>
    );
  } else {
    return (
      <div className="rating-number">
        <input
          className="text-input form-control"
          name="ratingnumber"
          type="number"
          onMouseDown={stepChange}
          onKeyDown={nonStepChange}
          onChange={handleChange}
          value={text}
          min="0.0"
          step="0.1"
          max="10"
          placeholder="0.0"
        />
      </div>
    );
  }
};
