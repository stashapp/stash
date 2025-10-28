import React from "react";
import { useIntl } from "react-intl";
import { getCountryByISO } from "src/utils/country";
import { OverlayTrigger, Tooltip } from "react-bootstrap";

interface ICountryFlag {
  country?: string | null;
  className?: string;
  includeName?: boolean;
  includeOverlay?: boolean;
}

export const CountryFlag: React.FC<ICountryFlag> = ({
  className,
  country: isoCountry,
  includeName,
  includeOverlay,
}) => {
  const { locale } = useIntl();

  const country = getCountryByISO(isoCountry, locale);

  if (!isoCountry || !country) return <></>;

  return (
    <>
      {includeName ? country : ""}
      {includeOverlay ? (
        <OverlayTrigger
          overlay={<Tooltip id="{country}-tooltip">{country}</Tooltip>}
        >
          <span
            className={`${className ?? ""} fi fi-${isoCountry.toLowerCase()}`}
          />
        </OverlayTrigger>
      ) : (
        <span
          className={`${className ?? ""} fi fi-${isoCountry.toLowerCase()}`}
        />
      )}
    </>
  );
};
