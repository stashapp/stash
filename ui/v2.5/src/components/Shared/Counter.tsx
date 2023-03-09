import React from "react";
import { Badge } from "react-bootstrap";
import { FormattedNumber, useIntl } from "react-intl";
import TextUtils from "src/utils/text";

interface IProps {
  abbreviateCounter?: boolean;
  count: number;
}

export const Counter: React.FC<IProps> = ({
  abbreviateCounter = false,
  count,
}) => {
  const intl = useIntl();

  if (abbreviateCounter) {
    const formated = TextUtils.abbreviateCounter(count);
    return (
      <Badge
        className="left-spacing"
        pill
        variant="secondary"
        data-value={intl.formatNumber(count)}
      >
        <FormattedNumber
          value={formated.size}
          maximumFractionDigits={formated.digits}
        />
        {formated.unit}
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
