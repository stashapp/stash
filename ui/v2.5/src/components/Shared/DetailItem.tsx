import React from "react";
import { FormattedMessage } from "react-intl";

interface IDetailItem {
  id?: string | null;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value?: any;
  title?: string;
  fullWidth?: boolean;
}

export const DetailItem: React.FC<IDetailItem> = ({
  id,
  value,
  title,
  fullWidth,
}) => {
  if (!id || !value || value === "Na") {
    return <></>;
  }

  const message = <FormattedMessage id={id} />;

  return (
    <div className={`detail-item ${id}`}>
      <span className={`detail-item-title ${id}`}>
        {message}
        {fullWidth ? ":" : ""}
      </span>
      <span className={`detail-item-value ${id}`} title={title}>
        {value}
      </span>
    </div>
  );
};
