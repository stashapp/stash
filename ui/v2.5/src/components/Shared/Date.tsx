import React from "react";
import { FormattedDate as IntlDate } from "react-intl";
import { PatchComponent } from "src/patch";

// wraps FormattedDate to handle year or year/month dates
export const FormattedDate: React.FC<{
  value: string | number | Date | undefined;
}> = PatchComponent("Date", ({ value }) => {
  if (typeof value === "string") {
    // try parsing as year or year/month
    const yearMatch = value.match(/^(\d{4})$/);
    if (yearMatch) {
      const year = parseInt(yearMatch[1], 10);
      return (
        <IntlDate value={Date.UTC(year, 0)} year="numeric" timeZone="utc" />
      );
    }

    const yearMonthMatch = value.match(/^(\d{4})-(\d{2})$/);
    if (yearMonthMatch) {
      const year = parseInt(yearMonthMatch[1], 10);
      const month = parseInt(yearMonthMatch[2], 10) - 1;

      return (
        <IntlDate
          value={Date.UTC(year, month)}
          year="numeric"
          month="long"
          timeZone="utc"
        />
      );
    }
  }

  return <IntlDate value={value} format="long" timeZone="utc" />;
});
