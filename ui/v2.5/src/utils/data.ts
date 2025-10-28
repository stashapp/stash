export const filterData = <T>(data?: (T | null | undefined)[] | null) =>
  data ? (data.filter((item) => item) as T[]) : [];

export interface IHasID {
  id: string;
}

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

export interface IHasStoredID {
  stored_id?: string | null;
}

export function sortStoredIdObjects<T extends IHasStoredID>(
  scrapedObjects?: T[]
): T[] | undefined {
  if (!scrapedObjects) {
    return undefined;
  }
  const ret = scrapedObjects.filter((p) => !!p.stored_id);

  if (ret.length === 0) {
    return undefined;
  }

  // sort by id numerically
  ret.sort((a, b) => {
    return parseInt(a.stored_id!, 10) - parseInt(b.stored_id!, 10);
  });

  return ret;
}
