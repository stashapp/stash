import { AudioCodecEnum } from "src/core/generated-graphql";

const stringAudioCodecMap = new Map<string, AudioCodecEnum>([
  ["aac", AudioCodecEnum.Aac],
]);

export const stringToAudioCodec = (
  value?: string | null,
  caseInsensitive?: boolean
) => {
  if (!value) {
    return undefined;
  }

  const ret = stringAudioCodecMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }

  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringAudioCodecMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const audioCodecStrings = Array.from(stringAudioCodecMap.keys());
