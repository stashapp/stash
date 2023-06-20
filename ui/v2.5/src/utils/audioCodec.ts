import { AudioCodecEnum } from "src/core/generated-graphql";

const stringAudioCodecMap = new Map<string, AudioCodecEnum>([
  ["aac", AudioCodecEnum.Aac],
  ["ac3", AudioCodecEnum.Ac3],
  ["adpcm_ima_wav", AudioCodecEnum.AdpcmImaWav],
  ["mp3", AudioCodecEnum.Mp3],
  ["pcm_s16le", AudioCodecEnum.PcmS16Le],
  ["wmapro", AudioCodecEnum.Wmapro],
  ["wmav2", AudioCodecEnum.Wmav2],
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
