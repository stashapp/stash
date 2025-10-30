import localForage from "localforage";
import isEqual from "lodash-es/isEqual";
import React, { Dispatch, SetStateAction, useEffect } from "react";
import { View } from "src/components/List/views";
import { ConfigImageLightboxInput } from "src/core/generated-graphql";

interface IInterfaceQueryConfig {
  filter: string;
  itemsPerPage: number;
  currentPage: number;
}

export interface IViewConfig {
  showSidebar?: boolean;
}

type IQueryConfig = Record<string, IInterfaceQueryConfig>;

interface IInterfaceConfig {
  queryConfig: IQueryConfig;
  imageLightbox: ConfigImageLightboxInput;
  // Partial is required because using View makes the key mandatory
  viewConfig: Partial<Record<View, IViewConfig>>;
}

export interface IChangelogConfig {
  versions: Record<string, boolean>;
}

interface ILocalForage<T> {
  data?: T;
  error: Error | null;
  loading: boolean;
}

const Loading: Record<string, boolean> = {};
const Cache: Record<string, {}> = {};

export function useLocalForage<T extends {}>(
  key: string,
  defaultValue: T = {} as T
): [ILocalForage<T>, Dispatch<SetStateAction<T>>] {
  const [error, setError] = React.useState<Error | null>(null);
  const [data, setData] = React.useState<T>(Cache[key] as T);
  const [loading, setLoading] = React.useState(Loading[key]);

  useEffect(() => {
    async function runAsync() {
      try {
        let parsed = await localForage.getItem<T>(key);
        if (typeof parsed === "string") {
          parsed = JSON.parse(parsed ?? "null");
        }
        if (parsed !== null) {
          setData(parsed);
          Cache[key] = parsed;
        } else {
          setData(defaultValue);
          Cache[key] = defaultValue;
        }
        setError(null);
      } catch (err) {
        if (err instanceof Error) setError(err);
        Cache[key] = defaultValue;
      } finally {
        Loading[key] = false;
        setLoading(false);
      }
    }

    if (!loading && !Cache[key]) {
      Loading[key] = true;
      setLoading(true);
      runAsync();
    }
  }, [loading, key, defaultValue]);

  useEffect(() => {
    if (!isEqual(Cache[key], data)) {
      Cache[key] = {
        ...Cache[key],
        ...data,
      };
      localForage.setItem(key, Cache[key]);
    }
  });

  const isLoading = loading || loading === undefined;

  return [{ data, error, loading: isLoading }, setData];
}

export const useInterfaceLocalForage = () =>
  useLocalForage<IInterfaceConfig>("interface");

export const useChangelogStorage = () =>
  useLocalForage<IChangelogConfig>("changelog");
