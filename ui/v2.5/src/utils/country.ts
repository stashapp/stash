import Countries from "i18n-iso-countries";
import EN from "i18n-iso-countries/langs/en.json";
import DE from "i18n-iso-countries/langs/de.json";
import ES from "i18n-iso-countries/langs/es.json";
import IT from "i18n-iso-countries/langs/it.json";
import PT from "i18n-iso-countries/langs/pt.json";
import SV from "i18n-iso-countries/langs/sv.json";
import ZH from "i18n-iso-countries/langs/zh.json";
import TW from "src/locales/countryNames/zh-TW.json";

Countries.registerLocale(EN);
Countries.registerLocale(DE);
Countries.registerLocale(ES);
Countries.registerLocale(IT);
Countries.registerLocale(PT);
Countries.registerLocale(SV);
Countries.registerLocale(ZH);
Countries.registerLocale(TW);

const getLocaleCode = (code: string) => {
  if (code === "zh-CN") return "zh";
  if (code === "zh-TW") return "tw";
  return code.slice(0, 2);
};

export const getCountryByISO = (
  iso: string | null | undefined,
  locale: string = "en"
): string | undefined => {
  if (!iso) return;

  return Countries.getName(iso, getLocaleCode(locale));
};

export const getCountries = (locale: string = "en") =>
  Object.entries(Countries.getNames(getLocaleCode(locale))).map(
    ([code, name]) => ({
      label: name,
      value: code,
    })
  );
