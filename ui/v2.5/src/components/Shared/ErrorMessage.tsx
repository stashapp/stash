import React, { ReactNode } from "react";

interface IProps {
  error: string | ReactNode;
}

const ErrorMessage: React.FC<IProps> = ({ error }) => (
  <div className="row ErrorMessage">
    <h2 className="ErrorMessage-content">Error: {error}</h2>
  </div>
);

export default ErrorMessage;
