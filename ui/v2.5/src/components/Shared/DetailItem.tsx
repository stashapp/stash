import React from "react";
import { FormattedMessage } from "react-intl";

interface IDetailItem {
  id?: string | null;
  value?: any;
  title?: string;
}

export const DetailItem: React.FC<IDetailItem> = ({ id, value, title }) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  if (!id || !value || value === "Na") return <></>;

  const message = <FormattedMessage id={id} />;

  return (
    <div className={`detail-item ${id}`}>
      <span className={`detail-item-title ${id}`}>{message}</span>
      <span className={`detail-item-value ${id}`} title={title}>
        {value}
      </span>
    </div>
  );
};
