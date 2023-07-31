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
