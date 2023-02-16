import React from "react";
import { Spinner } from "react-bootstrap";
import cx from "classnames";

interface ILoadingProps {
  message?: string;
  inline?: boolean;
  small?: boolean;
  card?: boolean;
}

const CLASSNAME = "LoadingIndicator";
const CLASSNAME_MESSAGE = `${CLASSNAME}-message`;

export const LoadingIndicator: React.FC<ILoadingProps> = ({
  message,
  inline = false,
  small = false,
  card = false,
}) => (
  <div className={cx(CLASSNAME, { inline, small, "card-based": card })}>
    <Spinner animation="border" role="status" size={small ? "sm" : undefined}>
      <span className="sr-only">Loading...</span>
    </Spinner>
    {message !== "" && (
      <h4 className={CLASSNAME_MESSAGE}>{message ?? "Loading..."}</h4>
    )}
  </div>
);
