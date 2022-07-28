import React from "react";
import { Badge } from "react-bootstrap";
import { FormattedNumber, useIntl } from "react-intl";
import { TextUtils } from "src/utils";

interface IProps {
  formatCounter?: boolean;
  count: number;
}

export const Counter: React.FC<IProps> = ({ formatCounter = false, count }) => {
  const intl = useIntl();

  if (formatCounter) {
    const formated = TextUtils.formatCounter(count);
    return (
      <Badge className="left-spacing" pill variant="secondary">
        <FormattedNumber
          value={formated.size}
          maximumFractionDigits={formated.digits}
        />
        {formated.unit}
      </Badge>
    );
  } else {
    return (
      <Badge className="left-spacing" pill variant="secondary">
        {intl.formatNumber(count)}
      </Badge>
    );
  }
};

export default Counter;
