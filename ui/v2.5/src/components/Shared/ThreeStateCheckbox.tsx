import { faCheck, faMinus, faTimes } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "./Icon";

interface IThreeStateCheckbox {
  value: boolean | undefined;
  setValue: (v: boolean | undefined) => void;
  allowUndefined?: boolean;
  label?: React.ReactNode;
  disabled?: boolean;
}

export const ThreeStateCheckbox: React.FC<IThreeStateCheckbox> = ({
  value,
  setValue,
  allowUndefined,
  label,
  disabled = false,
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

  const icon = value === undefined ? faMinus : value ? faCheck : faTimes;
  const labelClassName =
    value === undefined ? "unset" : value ? "checked" : "not-checked";

  return (
    <span className={`three-state-checkbox ${labelClassName}`}>
      <Button
        onClick={() => setValue(cycleState())}
        className="minimal"
        disabled={disabled}
      >
        <Icon icon={icon} className="fa-fw" />
      </Button>
      <span className="label">{label}</span>
    </span>
  );
};
