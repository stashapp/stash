import React from "react";
import { Card } from "react-bootstrap";
import { useIntl } from "react-intl";
import { PropsWithChildren } from "react-router/node_modules/@types/react";

interface ISettingGroup {
  id?: string;
  headingID?: string;
  subHeadingID?: string;
}

export const SettingSection: React.FC<PropsWithChildren<ISettingGroup>> = ({
  id,
  children,
  headingID,
  subHeadingID,
}) => {
  const intl = useIntl();

  return (
    <div className="setting-section" id={id}>
      <h1>{headingID ? intl.formatMessage({ id: headingID }) : undefined}</h1>
      {subHeadingID ? (
        <div className="sub-heading">
          {intl.formatMessage({ id: subHeadingID })}
        </div>
      ) : undefined}
      <Card>{children}</Card>
    </div>
  );
};
