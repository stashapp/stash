export const filterData = <T>(data?: (T | null | undefined)[] | null) =>
  data ? (data.filter((item) => item) as T[]) : [];

interface ITypename {
  __typename?: string;
}

export function withoutTypename<T extends ITypename>(o: T) {
  const { __typename, ...ret } = o;
  return ret;
}
