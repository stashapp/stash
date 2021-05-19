export const filterData = <T>(data?: (T | null | undefined)[] | null) =>
  data ? (data.filter((item) => item) as T[]) : [];
