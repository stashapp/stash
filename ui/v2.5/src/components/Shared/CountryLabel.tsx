import React from "react";
import { useIntl } from "react-intl";
import { CountryFlag } from "./CountryFlag";
import { getCountryByISO } from "src/utils/country";

interface IProps {
  country: string | undefined;
  showFlag?: boolean;
}

export const CountryLabel: React.FC<IProps> = ({
  country,
  showFlag = true,
}) => {
  const { locale } = useIntl();

  // #3063 - use alpha2 values only
  const fromISO =
    country?.length === 2 ? getCountryByISO(country, locale) : undefined;

  return (
    <div>
      {showFlag && <CountryFlag className="mr-2" country={country} />}
      <span>{fromISO ?? country}</span>
    </div>
  );
};
