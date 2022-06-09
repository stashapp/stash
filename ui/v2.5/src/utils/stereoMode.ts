import * as GQL from "../core/generated-graphql";

export const stringStereoModeMap = new Map<string, GQL.StereoModeEnum>([
  ["MONO", GQL.StereoModeEnum.Mono],
  ["TOP_BOTTOM", GQL.StereoModeEnum.TopBottom],
  ["LEFT_RIGHT", GQL.StereoModeEnum.LeftRight],
]);

export const stereoModeToString = (value?: GQL.StereoModeEnum | string) => {
  if (!value) {
    return undefined;
  }

  const foundEntry = Array.from(stringStereoModeMap.entries()).find((e) => {
    return e[1] === value;
  });

  if (foundEntry) {
    return foundEntry[0];
  }
};

export const stringToStereoMode = (
  value?: string | null,
  caseInsensitive?: boolean
): GQL.StereoModeEnum | undefined => {
  if (!value) {
    return undefined;
  }

  const existing = Object.entries(GQL.StereoModeEnum).find(
    (e) => e[1] === value
  );
  if (existing) return existing[1];

  const ret = stringStereoModeMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }

  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringStereoModeMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const stereoModeStrings = Array.from(stringStereoModeMap.keys());
