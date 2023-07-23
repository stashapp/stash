import React from "react";
import { useIntl } from "react-intl";
import { getCountryByISO } from "src/utils/country";

interface IDetailItem {
  header?: string | null;
  value?: any;
}

export const DetailItem: React.FC<IDetailItem> = ({ header, value }) => {
  if (!value || value === "Na") return <></>;

  return (
    <div className="quick-detail-item">
      <span className="quick-detail-header">{header}</span>
      <span className="quick-detail-value">{value} </span>
    </div>
  );
};
