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
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
} as { [key: string]: any };

export default localeLoader;
