import React from "react";
import { FormattedMessage } from "react-intl";

interface IDetailItem {
  id?: string | null;
  className?: string;
  label?: React.ReactNode;
  value?: React.ReactNode;
  labelTitle?: string;
  title?: string;
  fullWidth?: boolean;
  showEmpty?: boolean;
}

export const DetailItem: React.FC<IDetailItem> = ({
  id,
  className = "",
  label,
  value,
  labelTitle,
  title,
  fullWidth,
  showEmpty = false,
}) => {
  if (!id || (!showEmpty && (!value || value === "Na"))) {
    return <></>;
  }

  const message = label ?? <FormattedMessage id={id} />;

  // according to linter rule CSS classes shouldn't use underscores
  const sanitisedID = id.replace(/_/g, "-");

  return (
    <div className={`detail-item ${id} ${className}`}>
      <span className={`detail-item-title ${sanitisedID}`} title={labelTitle}>
        {message}
        {fullWidth ? ":" : ""}
      </span>
      <span className={`detail-item-value ${sanitisedID}`} title={title}>
        {value}
      </span>
    </div>
  );
};
