import React from "react";
import { getCountryByISO } from "src/utils";

interface ICountryFlag {
  country?: string | null;
  className?: string;
}

const CountryFlag: React.FC<ICountryFlag> = ({ className, country }) => {
  if (!country) return <></>;

  return (
    <span
      className={`${
        className ?? ""
      } flag-icon flag-icon-${country.toLowerCase()}`}
      title={getCountryByISO(country)}
    />
  );
};

export default CountryFlag;
