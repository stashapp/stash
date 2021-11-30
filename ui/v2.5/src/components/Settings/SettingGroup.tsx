import React from "react";
import { Card } from "react-bootstrap";
import { useIntl } from "react-intl";
import { PropsWithChildren } from "react-router/node_modules/@types/react";

interface ISettingGroup {
  headingID: string;
  subHeadingID?: string;
}

export const SettingGroup: React.FC<PropsWithChildren<ISettingGroup>> = ({
  children,
  headingID,
  subHeadingID,
}) => {
  const intl = useIntl();

  return (
    <div className="setting-group">
      <h1>{intl.formatMessage({ id: headingID })}</h1>
      {subHeadingID ? (
        <h2>{intl.formatMessage({ id: subHeadingID })}</h2>
      ) : undefined}
      <Card>{children}</Card>
    </div>
  );
};
