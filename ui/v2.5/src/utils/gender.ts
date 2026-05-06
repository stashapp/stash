import * as GQL from "../core/generated-graphql";

export const stringGenderMap = new Map<string, GQL.GenderEnum>([
  ["Female", GQL.GenderEnum.Female],
  ["Male", GQL.GenderEnum.Male],
  ["Transgender Male", GQL.GenderEnum.TransgenderMale],
  ["Transgender Female", GQL.GenderEnum.TransgenderFemale],
  ["Intersex", GQL.GenderEnum.Intersex],
  ["Non-Binary", GQL.GenderEnum.NonBinary],
]);

export const genderList = [
  GQL.GenderEnum.Female,
  GQL.GenderEnum.Male,
  GQL.GenderEnum.TransgenderFemale,
  GQL.GenderEnum.TransgenderMale,
  GQL.GenderEnum.Intersex,
  GQL.GenderEnum.NonBinary,
];

export const genderToString = (value?: GQL.GenderEnum | string | null) => {
  if (!value) {
    return undefined;
  }

  const foundEntry = Array.from(stringGenderMap.entries()).find((e) => {
    return e[1] === value;
  });

  if (foundEntry) {
    return foundEntry[0];
  }
};

export const stringToGender = (
  value?: string | null,
  caseInsensitive?: boolean
): GQL.GenderEnum | undefined => {
  if (!value) {
    return undefined;
  }

  const existing = Object.entries(GQL.GenderEnum).find((e) => e[1] === value);
  if (existing) return existing[1];

  const ret = stringGenderMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }

  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringGenderMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const genderStrings = Array.from(stringGenderMap.keys());
