import React from "react";
import { useIntl } from "react-intl";
import { getCountryByISO } from "src/utils/country";

interface ICountryFlag {
  country?: string | null;
  className?: string;
}

export const CountryFlag: React.FC<ICountryFlag> = ({
  className,
  country: isoCountry,
}) => {
  const { locale } = useIntl();

  const country = getCountryByISO(isoCountry, locale);

  if (!isoCountry || !country) return <></>;

  return (
    <span
      className={`${className ?? ""} fi fi-${isoCountry.toLowerCase()}`}
      title={country}
    />
  );
};
