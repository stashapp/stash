import Countries from "i18n-iso-countries";
import english from "i18n-iso-countries/langs/en.json";

Countries.registerLocale(english);

export const getCountryByISO = (
  iso: string | null | undefined
): string | undefined => {
  if (!iso) return;

  return Countries.getName(iso, "en");
};

export const getCountries = (language: string = "en") =>
  Object.entries(Countries.getNames(language)).map(([code, name]) => ({
    label: name,
    value: code,
  }));
