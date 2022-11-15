import React from "react";
import { useIntl } from "react-intl";
import { CountryFlag } from "src/components/Shared";
import { getCountryByISO } from "src/utils";

interface IProps {
  country: string | undefined;
  showFlag?: boolean;
}

const CountryLabel: React.FC<IProps> = ({ country, showFlag = true }) => {
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

export default CountryLabel;
