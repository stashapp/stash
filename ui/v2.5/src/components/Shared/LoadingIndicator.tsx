import React from "react";
import { Spinner } from "react-bootstrap";
import cx from "classnames";

interface ILoadingProps {
  message?: string;
  inline?: boolean;
}

const CLASSNAME = "LoadingIndicator";
const CLASSNAME_MESSAGE = `${CLASSNAME}-message`;

const LoadingIndicator: React.FC<ILoadingProps> = ({
  message,
  inline = false,
}) => (
  <div className={cx(CLASSNAME, { inline })}>
    <Spinner animation="border" role="status">
      <span className="sr-only">Loading...</span>
    </Spinner>
    <h4 className={CLASSNAME_MESSAGE}>{message ?? "Loading..."}</h4>
  </div>
);

export default LoadingIndicator;
