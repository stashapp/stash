import { clone } from "lodash-es";

export const filterData = <T>(data?: (T | null | undefined)[] | null) =>
  data ? (data.filter((item) => item) as T[]) : [];

export interface ITypename {
  __typename?: string;
}

const hasTypename = (value: unknown): value is ITypename =>
  !!(value as ITypename)?.__typename;

const processNoneObjValue = (value: unknown): unknown =>
  Array.isArray(value)
    ? value.map((v) =>
        hasTypename(v) ? withoutTypename(v) : processNoneObjValue(v)
      )
    : value;

export function withoutTypename<T extends ITypename>(
  o: T
): Omit<T, "__typename"> {
  const { __typename, ...data } = o;

  return Object.entries(data).reduce(
    (ret, [key, value]) => ({
      ...ret,
      [key]: hasTypename(value)
        ? withoutTypename(value)
        : processNoneObjValue(value),
    }),
    {} as Omit<T, "__typename">
  );
}

// excludeFields removes fields from data that are in the excluded object
export function excludeFields(
  data: { [index: string]: unknown },
  excluded: Record<string, boolean>
) {
  Object.keys(data).forEach((k) => {
    if (excluded[k] || !data[k]) {
      data[k] = undefined;
    }
  });
}

export interface IHasID {
  id: string;
}

export function sortIdObjectList<T extends IHasID>(list?: T[] | null) {
  if (!list) {
    return;
  }

  const ret = clone(list);
  // sort by id numerically
  ret.sort((a, b) => {
    return parseInt(a.id, 10) - parseInt(b.id, 10);
  });

  return ret;
}
