import { VideoCodecEnum } from "src/core/generated-graphql";

const stringVideoCodecMap = new Map<string, VideoCodecEnum>([
  ["h264", VideoCodecEnum.H264],
  ["hevc", VideoCodecEnum.Hevc],
  ["av1", VideoCodecEnum.Av1],
  ["wmv3", VideoCodecEnum.Wmv3],
  ["vc1", VideoCodecEnum.Vc1],
]);

export const stringToVideoCodec = (
  value?: string | null,
  caseInsensitive?: boolean
) => {
  if (!value) {
    return undefined;
  }

  const ret = stringVideoCodecMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }

  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringVideoCodecMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const videoCodecStrings = Array.from(stringVideoCodecMap.keys());
