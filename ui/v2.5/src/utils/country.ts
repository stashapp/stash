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
  "Slovak Republic": "SK",
  Iran: "IR",
  Moldova: "MD",
  Laos: "LA",
};

const getISOCountry = (country: string | null | undefined) => {
  if (!country) return null;

  const code =
    fuzzyDict[country] ?? Countries.getAlpha2Code(country, "en") ?? country;
  // Check if code is valid alpha2 iso
  if (!Countries.alpha2ToAlpha3(code)) return null;

  return {
    code,
    name: Countries.getName(code, "en"),
  };
};

export const getCountryByISO = (iso: string | null | undefined) => {
  if (!iso) return null;

  return Countries.getName(iso, "en") ?? null;
};

export default getISOCountry;
