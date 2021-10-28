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

  return (
    <div>
      {showFlag && <CountryFlag country={country} />}
      <span className="ml-2">{getCountryByISO(country, locale)}</span>
    </div>
  );
};

export default CountryLabel;
