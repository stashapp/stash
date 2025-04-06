import React, { ReactNode } from "react";
import { Alert } from "react-bootstrap";
import { FormattedMessage } from "react-intl";

interface IProps {
  message?: React.ReactNode;
  error: string | ReactNode;
}

export const ErrorMessage: React.FC<IProps> = (props) => {
  const { error, message = <FormattedMessage id="errors.header" /> } = props;

  return (
    <div className="ErrorMessage-container">
      <Alert variant="danger" className="ErrorMessage">
        <Alert.Heading className="ErrorMessage-header">{message}</Alert.Heading>
        <div className="ErrorMessage-content code">{error}</div>
      </Alert>
    </div>
  );
};
