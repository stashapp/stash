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

  const fromISO = getCountryByISO(country, locale);

  return (
    <div>
      {showFlag && <CountryFlag className="mr-2" country={country} />}
      <span>{fromISO ?? country}</span>
    </div>
  );
};

export default CountryLabel;
