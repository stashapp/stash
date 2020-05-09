import React from "react";
import { getISOCountry } from "src/utils";

interface ICountryFlag {
  country?: string | null;
  className?: string;
}

const CountryFlag: React.FC<ICountryFlag> = ({ className, country }) => {
  const ISOCountry = getISOCountry(country);
  if (!ISOCountry?.code) return <></>;

  return (
    <span
      className={`${
        className ?? ""
      } flag-icon flag-icon-${ISOCountry.code.toLowerCase()}`}
      title={ISOCountry.name}
    />
  );
};

export default CountryFlag;
