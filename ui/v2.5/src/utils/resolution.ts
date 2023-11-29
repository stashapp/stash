import { ResolutionEnum } from "src/core/generated-graphql";

const stringResolutionMap = new Map<string, ResolutionEnum>([
  ["144p", ResolutionEnum.VeryLow],
  ["240p", ResolutionEnum.Low],
  ["360p", ResolutionEnum.R360P],
  ["480p", ResolutionEnum.Standard],
  ["540p", ResolutionEnum.WebHd],
  ["720p", ResolutionEnum.StandardHd],
  ["1080p", ResolutionEnum.FullHd],
  ["1440p", ResolutionEnum.QuadHd],
  // ["1920p", ResolutionEnum.VrHd],
  ["4k", ResolutionEnum.FourK],
  ["5k", ResolutionEnum.FiveK],
  ["6k", ResolutionEnum.SixK],
  ["7k", ResolutionEnum.SevenK],
  ["8k", ResolutionEnum.EightK],
  ["Huge", ResolutionEnum.Huge],
]);

export const stringToResolution = (
  value?: string | null,
  caseInsensitive?: boolean
) => {
  if (!value) {
    return undefined;
  }

  const ret = stringResolutionMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }

  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringResolutionMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const resolutionStrings = Array.from(stringResolutionMap.keys());
