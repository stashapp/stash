import React from "react";
import { FormattedMessage } from "react-intl";

interface IDetailItem {
  id?: string | null;
  label?: React.ReactNode;
  value?: React.ReactNode;
  labelTitle?: string;
  title?: string;
  fullWidth?: boolean;
  showEmpty?: boolean;
}

export const DetailItem: React.FC<IDetailItem> = ({
  id,
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

  return (
    // according to linter rule CSS classes shouldn't use underscores
    <div className={`detail-item ${id}`}>
      <span
        className={`detail-item-title ${id.replace("_", "-")}`}
        title={labelTitle}
      >
        {message}
        {fullWidth ? ":" : ""}
      </span>
      <span
        className={`detail-item-value ${id.replace("_", "-")}`}
        title={title}
      >
        {value}
      </span>
    </div>
  );
};
