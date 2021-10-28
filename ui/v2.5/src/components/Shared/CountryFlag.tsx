import React from "react";
import { useIntl } from "react-intl";
import { getCountryByISO } from "src/utils";

interface ICountryFlag {
  country?: string | null;
  className?: string;
}

const CountryFlag: React.FC<ICountryFlag> = ({ className, country }) => {
  const { locale } = useIntl();
  if (!country) return <></>;

  return (
    <span
      className={`${
        className ?? ""
      } flag-icon flag-icon-${country.toLowerCase()}`}
      title={getCountryByISO(country, locale)}
    />
  );
};

export default CountryFlag;
