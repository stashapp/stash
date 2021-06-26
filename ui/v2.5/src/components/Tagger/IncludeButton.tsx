import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared";

interface IIncludeExcludeButton {
  exclude: boolean;
  setExclude: (v: boolean) => void;
}

export const IncludeExcludeButton: React.FC<IIncludeExcludeButton> = ({
  exclude,
  setExclude,
}) => (
  <Button
    onClick={() => setExclude(!exclude)}
    variant="minimal"
    className={`${
      exclude ? "text-danger" : "text-success"
    } include-exclude-button`}
  >
    <Icon className="fa-fw" icon={exclude ? "times" : "check"} />
  </Button>
);

interface IOptionalField {
  exclude: boolean;
  setExclude: (v: boolean) => void;
}

export const OptionalField: React.FC<IOptionalField> = ({
  exclude,
  setExclude,
  children,
}) => (
  <div className={`optional-field ${!exclude ? "included" : "excluded"}`}>
    <IncludeExcludeButton exclude={exclude} setExclude={setExclude} />
    {children}
  </div>
);
