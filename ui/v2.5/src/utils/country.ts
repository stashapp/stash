import Countries from "i18n-iso-countries";
import english from "i18n-iso-countries/langs/en.json";

Countries.registerLocale(english);

const fuzzyDict: Record<string, string> = {
  USA: "US",
  "United States": "US",
  America: "US",
  American: "US",
  Czechia: "CZ",
  England: "GB",
  "United Kingdom": "GB",
  Russia: "RU",
};

const getISOCode = (country: string | null | undefined) => {
  if (!country) return null;

  if (fuzzyDict[country]) return fuzzyDict[country];

  return Countries.getAlpha2Code(country, "en");
};

export default getISOCode;
