import React, { ReactNode } from "react";
import { FormattedMessage } from "react-intl";

interface IProps {
  message?: React.ReactNode;
  error: string | ReactNode;
}

export const ErrorMessage: React.FC<IProps> = (props) => {
  const { error, message = <FormattedMessage id="errors.header" /> } = props;

  return (
    <div className="row ErrorMessage">
      <h2 className="ErrorMessage-content">
        {message}: {error}
      </h2>
    </div>
  );
};
