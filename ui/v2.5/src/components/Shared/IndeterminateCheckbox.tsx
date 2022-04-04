import React, { useEffect } from "react";
import { Form, FormCheckProps } from "react-bootstrap";

const useIndeterminate = (
  ref: React.RefObject<HTMLInputElement>,
  value: boolean | undefined
) => {
  useEffect(() => {
    if (ref.current) {
      // eslint-disable-next-line no-param-reassign
      ref.current.indeterminate = value === undefined;
    }
  }, [ref, value]);
};

interface IIndeterminateCheckbox extends FormCheckProps {
  setChecked: (v: boolean | undefined) => void;
  allowIndeterminate?: boolean;
  indeterminateClassname?: string;
}

export const IndeterminateCheckbox: React.FC<IIndeterminateCheckbox> = ({
  checked,
  setChecked,
  allowIndeterminate,
  indeterminateClassname,
  ...props
}) => {
  const ref = React.createRef<HTMLInputElement>();

  useIndeterminate(ref, checked);

  function cycleState() {
    const undefAllowed = allowIndeterminate ?? true;
    if (undefAllowed && checked) {
      return undefined;
    }
    if ((!undefAllowed && checked) || checked === undefined) {
      return false;
    }
    return true;
  }

  return (
    <Form.Check
      {...props}
      className={`${props.className ?? ""} ${
        checked === undefined ? indeterminateClassname : ""
      }`}
      ref={ref}
      checked={checked ?? false}
      onChange={() => setChecked(cycleState())}
    />
  );
};
