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
  const result = {} as Omit<T, "__typename">;

  for (const [key, value] of Object.entries(o)) {
    if (key === "__typename") {
      continue;
    }

    (result as Record<string, unknown>)[key] = hasTypename(value)
      ? withoutTypename(value)
      : processNoneObjValue(value);
  }

  return result;
}

export function listToMap(list: string[]) {
  const map: Record<string, boolean> = {};
  list.forEach((item) => {
    map[item] = true;
  });
  return map;
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

export function uniqIDStoredIDs<T extends IHasStoredID>(objs: T[]) {
  return objs.filter((o, i) => {
    return objs.findIndex((oo) => oo.stored_id === o.stored_id) === i;
  });
}

export function idToStoredID(o: { id: string; name: string }) {
  return {
    stored_id: o.id,
    name: o.name,
  };
}

export function getActiveSortColumn(
  sortMap: Record<string, string>,
  sortBy?: string | null
): string | undefined {
  const reverseMap = Object.fromEntries(
    Object.entries(sortMap).map(([k, v]) => [v, k])
  );
  return reverseMap[sortBy ?? ""] ?? sortBy;
}
