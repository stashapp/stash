import Countries from "i18n-iso-countries";

export const localeCountries = {
  bn: () => import("i18n-iso-countries/langs/bn.json"),
  cs: () => import("i18n-iso-countries/langs/cs.json"),
  da: () => import("i18n-iso-countries/langs/da.json"),
  de: () => import("i18n-iso-countries/langs/de.json"),
  en: () => import("i18n-iso-countries/langs/en.json"),
  es: () => import("i18n-iso-countries/langs/es.json"),
  et: () => import("i18n-iso-countries/langs/et.json"),
  fa: () => import("i18n-iso-countries/langs/fa.json"),
  fi: () => import("i18n-iso-countries/langs/fi.json"),
  fr: () => import("i18n-iso-countries/langs/fr.json"),
  hu: () => import("i18n-iso-countries/langs/hu.json"),
  hr: () => import("i18n-iso-countries/langs/hr.json"),
  it: () => import("i18n-iso-countries/langs/it.json"),
  ja: () => import("i18n-iso-countries/langs/ja.json"),
  ko: () => import("i18n-iso-countries/langs/ko.json"),
  nl: () => import("i18n-iso-countries/langs/nl.json"),
  pl: () => import("i18n-iso-countries/langs/pl.json"),
  pt: () => import("i18n-iso-countries/langs/pt.json"),
  ro: () => import("i18n-iso-countries/langs/ro.json"),
  ru: () => import("i18n-iso-countries/langs/ru.json"),
  sv: () => import("i18n-iso-countries/langs/sv.json"),
  th: () => import("i18n-iso-countries/langs/th.json"),
  tr: () => import("i18n-iso-countries/langs/tr.json"),
  uk: () => import("i18n-iso-countries/langs/uk.json"),
  zh: () => import("i18n-iso-countries/langs/zh.json"),
  tw: () => import("src/locales/countryNames/zh-TW.json"),
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
} as { [key: string]: any };

export const getLocaleCode = (code: string) => {
  if (code === "zh-CN") return "zh";
  if (code === "zh-TW") return "tw";
  return code.slice(0, 2);
};

export async function registerCountry(locale: string) {
  const localeCode = getLocaleCode(locale);
  const countries = await localeCountries[localeCode]();
  Countries.registerLocale(countries);
}

export const localeLoader = {
  deDE: () => import("./de-DE.json"),
  enGB: () => import("./en-GB.json"),
  enUS: () => import("./en-US.json"),
  esES: () => import("./es-ES.json"),
  ptBR: () => import("./pt-BR.json"),
  frFR: () => import("./fr-FR.json"),
  itIT: () => import("./it-IT.json"),
  fiFI: () => import("./fi-FI.json"),
  svSE: () => import("./sv-SE.json"),
  zhTW: () => import("./zh-TW.json"),
  zhCN: () => import("./zh-CN.json"),
  hrHR: () => import("./hr-HR.json"),
  nlNL: () => import("./nl-NL.json"),
  ruRU: () => import("./ru-RU.json"),
  trTR: () => import("./tr-TR.json"),
  jaJP: () => import("./ja-JP.json"),
  plPL: () => import("./pl-PL.json"),
  daDK: () => import("./da-DK.json"),
  koKR: () => import("./ko-KR.json"),
  ukUA: () => import("./uk-UA.json"),
  bnBD: () => import("./bn-BD.json"),
  csCZ: () => import("./cs-CZ.json"),
  etEE: () => import("./et-EE.json"),
  faIR: () => import("./fa-IR.json"),
  huHU: () => import("./hu-HU.json"),
  roRO: () => import("./ro-RO.json"),
  thTH: () => import("./th-TH.json"),
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
} as { [key: string]: any };

export default localeLoader;
