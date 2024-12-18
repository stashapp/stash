import Countries from "i18n-iso-countries";

export const localeCountries = {
  af: () => import("i18n-iso-countries/langs/af.json"),
  bn: () => import("i18n-iso-countries/langs/bn.json"),
  ca: () => import("i18n-iso-countries/langs/ca.json"),
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
  id: () => import("i18n-iso-countries/langs/id.json"),
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
  afZA: () => import("./af-ZA.json"),
  bnBD: () => import("./bn-BD.json"),
  caES: () => import("./ca-ES.json"),
  csCZ: () => import("./cs-CZ.json"),
  daDK: () => import("./da-DK.json"),
  deDE: () => import("./de-DE.json"),
  enGB: () => import("./en-GB.json"),
  enUS: () => import("./en-US.json"),
  esES: () => import("./es-ES.json"),
  etEE: () => import("./et-EE.json"),
  faIR: () => import("./fa-IR.json"),
  fiFI: () => import("./fi-FI.json"),
  frFR: () => import("./fr-FR.json"),
  hrHR: () => import("./hr-HR.json"),
  huHU: () => import("./hu-HU.json"),
  idID: () => import("./id-ID.json"),
  itIT: () => import("./it-IT.json"),
  jaJP: () => import("./ja-JP.json"),
  koKR: () => import("./ko-KR.json"),
  // nbNO: () => import("./nb-NO.json"),
  // neNP: () => import("./ne-NP.json"),
  nlNL: () => import("./nl-NL.json"),
  plPL: () => import("./pl-PL.json"),
  ptBR: () => import("./pt-BR.json"),
  roRO: () => import("./ro-RO.json"),
  ruRU: () => import("./ru-RU.json"),
  svSE: () => import("./sv-SE.json"),
  thTH: () => import("./th-TH.json"),
  trTR: () => import("./tr-TR.json"),
  ukUA: () => import("./uk-UA.json"),
  zhCN: () => import("./zh-CN.json"),
  zhTW: () => import("./zh-TW.json"),
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
} as { [key: string]: any };

export default localeLoader;
