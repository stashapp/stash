import React from "react";
import { useIntl } from "react-intl";
import { getCountryByISO } from "src/utils";

interface ICountryFlag {
  country?: string | null;
  className?: string;
}

const CountryFlag: React.FC<ICountryFlag> = ({
  className,
  country: isoCountry,
}) => {
  const { locale } = useIntl();

  const country = getCountryByISO(isoCountry, locale);

  if (!isoCountry || !country) return <></>;

  return (
    <span
      className={`${
        className ?? ""
      } flag-icon flag-icon-${isoCountry.toLowerCase()}`}
      title={country}
    />
  );
};

export default CountryFlag;
