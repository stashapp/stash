import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from ".";

interface IThreeStateCheckbox {
  value: boolean | undefined;
  setValue: (v: boolean | undefined) => void;
  allowUndefined?: boolean;
  label?: React.ReactNode;
}

export const ThreeStateCheckbox: React.FC<IThreeStateCheckbox> = ({
  value,
  setValue,
  allowUndefined,
  label,
}) => {
  function cycleState() {
    const undefAllowed = allowUndefined ?? true;
    if (undefAllowed && value) {
      return undefined;
    }
    if ((!undefAllowed && value) || value === undefined) {
      return false;
    }
    return true;
  }

  const icon = value === undefined ? "square" : value ? "check" : "times";
  const labelClassName =
    value === undefined ? "unset" : value ? "checked" : "not-checked";

  return (
    <span className={`three-state-checkbox ${labelClassName}`}>
      <Button onClick={() => setValue(cycleState())} className="minimal">
        <Icon icon={icon} className="fa-fw" />
      </Button>
      <span className="label">{label}</span>
    </span>
  );
};
