import React from "react";
import { FormattedMessage } from "react-intl";

interface IDetailItem {
  id?: string | null;
  label?: React.ReactNode;
  value?: React.ReactNode;
  title?: string;
  fullWidth?: boolean;
}

export const DetailItem: React.FC<IDetailItem> = ({
  id,
  label,
  value,
  title,
  fullWidth,
}) => {
  if (!id || !value || value === "Na") {
    return <></>;
  }

  const message = label ?? <FormattedMessage id={id} />;

  return (
    // according to linter rule CSS classes shouldn't use underscores
    <div className={`detail-item ${id}`}>
      <span className={`detail-item-title ${id.replace("_", "-")}`}>
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
