import { faCheck, faTimes } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared/Icon";

interface IIncludeExcludeButton {
  exclude: boolean;
  disabled?: boolean;
  setExclude: (v: boolean) => void;
}

export const IncludeExcludeButton: React.FC<IIncludeExcludeButton> = ({
  exclude,
  disabled,
  setExclude,
}) => (
  <Button
    onClick={() => setExclude(!exclude)}
    disabled={disabled}
    variant="minimal"
    className={`${
      exclude ? "text-danger" : "text-success"
    } include-exclude-button`}
  >
    <Icon className="fa-fw" icon={exclude ? faTimes : faCheck} />
  </Button>
);

interface IOptionalField {
  exclude: boolean;
  title?: string;
  disabled?: boolean;
  setExclude: (v: boolean) => void;
}

export const OptionalField: React.FC<IOptionalField> = ({
  exclude,
  setExclude,
  children,
  title,
}) => {
  return (
    <div className={`optional-field ${!exclude ? "included" : "excluded"}`}>
      <IncludeExcludeButton exclude={exclude} setExclude={setExclude} />
      {title && <span className="optional-field-title">{title}</span>}
      <div className="optional-field-content">{children}</div>
    </div>
  );
};
