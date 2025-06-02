import { useEffect, useState } from "react";

export interface ILoadResults<T> {
  results: T;
  loading: boolean;
}

export function useCacheResults<T>(data: ILoadResults<T>) {
  const [results, setResults] = useState<T | undefined>(
    !data.loading ? data.results : undefined
  );

  useEffect(() => {
    if (!data.loading) {
      setResults(data.results);
    }
  }, [data.loading, data.results]);

  return { loading: data.loading, results };
}
