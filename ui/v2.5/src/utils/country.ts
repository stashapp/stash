import Countries from "i18n-iso-countries";
import { getLocaleCode } from "src/locales";

export const getCountryByISO = (
  iso: string | null | undefined,
  locale: string = "en"
): string | undefined => {
  if (!iso) return;

  const ret = Countries.getName(iso, getLocaleCode(locale));
  if (ret) {
    return ret;
  }

  // fallback to english if locale is not en
  if (locale !== "en") {
    return Countries.getName(iso, "en");
  }
};

export const getCountries = (locale: string = "en") => {
  let countries = Countries.getNames(getLocaleCode(locale));

  if (!countries.length) {
    countries = Countries.getNames("en");
  }

  return Object.entries(countries).map(([code, name]) => ({
    label: name,
    value: code,
  }));
};
