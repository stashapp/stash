import React, { PropsWithChildren } from "react";
import { Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";

interface ITask {
  headingID?: string;
  description?: React.ReactNode;
}

export const Task: React.FC<PropsWithChildren<ITask>> = ({
  children,
  headingID,
  description,
}) => (
  <div className="task">
    {headingID ? (
      <h6>
        <FormattedMessage id={headingID} />
      </h6>
    ) : undefined}
    {children}
    {description ? (
      <Form.Text className="text-muted">{description}</Form.Text>
    ) : undefined}
  </div>
);
