import * as GQL from "../core/generated-graphql";

export const stringProjectionMap = new Map<string, GQL.ProjectionEnum>([
  ["AUTO", GQL.ProjectionEnum.Auto],
  ["FLAT", GQL.ProjectionEnum.Flat],
  ["DOME", GQL.ProjectionEnum.Dome],
  ["SPHERE", GQL.ProjectionEnum.Sphere],
  ["FISHEYE", GQL.ProjectionEnum.Fisheye],
  ["MKX200", GQL.ProjectionEnum.Mkx200],
  ["RF52", GQL.ProjectionEnum.Rf52],
  ["CUBE", GQL.ProjectionEnum.Cube],
  ["EAC", GQL.ProjectionEnum.Eac],
]);

export const projectionToString = (value?: GQL.ProjectionEnum | string) => {
  if (!value) {
    return undefined;
  }

  const foundEntry = Array.from(stringProjectionMap.entries()).find((e) => {
    return e[1] === value;
  });

  if (foundEntry) {
    return foundEntry[0];
  }
};

export const stringToProjection = (
  value?: string | null,
  caseInsensitive?: boolean
): GQL.ProjectionEnum | undefined => {
  if (!value) {
    return undefined;
  }

  const existing = Object.entries(GQL.ProjectionEnum).find(
    (e) => e[1] === value
  );
  if (existing) return existing[1];

  const ret = stringProjectionMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }

  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringProjectionMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const projectionStrings = Array.from(stringProjectionMap.keys());
