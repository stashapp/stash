import React from "react";
import { Badge } from "react-bootstrap";
import { FormattedNumber, useIntl } from "react-intl";
import TextUtils from "src/utils/text";

interface IProps {
  abbreviateCounter?: boolean;
  count: number;
  hideZero?: boolean;
  hideOne?: boolean;
}

export const Counter: React.FC<IProps> = ({
  abbreviateCounter = false,
  count,
  hideZero = false,
  hideOne = false,
}) => {
  const intl = useIntl();

  if (hideZero && count === 0) return null;
  if (hideOne && count === 1) return null;

  if (abbreviateCounter) {
    const formatted = TextUtils.abbreviateCounter(count);
    return (
      <Badge
        className="left-spacing"
        pill
        variant="secondary"
        data-value={intl.formatNumber(count)}
      >
        <FormattedNumber
          value={formatted.size}
          maximumFractionDigits={formatted.digits}
        />
        {formatted.unit}
      </Badge>
    );
  } else {
    return (
      <Badge
        className="left-spacing"
        pill
        variant="secondary"
        data-value={intl.formatNumber(count)}
      >
        {intl.formatNumber(count)}
      </Badge>
    );
  }
};
