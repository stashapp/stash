import { VideoCodecEnum } from "src/core/generated-graphql";

const stringVideoCodecMap = new Map<string, VideoCodecEnum>([
  ["av1", VideoCodecEnum.Av1],
  ["h264", VideoCodecEnum.H264],
  ["hevc", VideoCodecEnum.Hevc],
  ["mpeg2video", VideoCodecEnum.Mpeg2Video],
  ["mpeg4", VideoCodecEnum.Mpeg4],
  ["vc1", VideoCodecEnum.Vc1],
  ["vp6f", VideoCodecEnum.Vp6F],
  ["wmv1", VideoCodecEnum.Wmv1],
  ["wmv2", VideoCodecEnum.Wmv2],
  ["wmv3", VideoCodecEnum.Wmv3],
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
