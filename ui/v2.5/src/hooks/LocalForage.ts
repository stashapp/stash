import localForage from "localforage";
import isEqual from "lodash-es/isEqual";
import React, { Dispatch, SetStateAction, useEffect } from "react";
import { View } from "src/components/List/views";
import { ConfigImageLightboxInput } from "src/core/generated-graphql";
import { DisplayMode } from "src/models/list-filter/types";

interface IInterfaceQueryConfig {
  filter: string;
  itemsPerPage: number;
  currentPage: number;
}

/** Per-view UI configuration persisted to IndexedDB. */
export interface IViewConfig {
  showSidebar?: boolean;
  displayMode?: DisplayMode;
}

type IQueryConfig = Record<string, IInterfaceQueryConfig>;

interface IInterfaceConfig {
  queryConfig: IQueryConfig;
  imageLightbox: ConfigImageLightboxInput;
  // Partial is required because using View makes the key mandatory
  viewConfig: Partial<Record<View, IViewConfig>>;
}

/** Persisted changelog acknowledgement state. */
export interface IChangelogConfig {
  versions: Record<string, boolean>;
}

interface ILocalForage<T> {
  data?: T;
  error: Error | null;
  loading: boolean;
}

const Loading: Record<string, boolean> = {};
const Cache: Record<string, object> = {};

/**
 * Read/write a JSON-serialisable value to IndexedDB via localforage.
 * The value is cached in memory so subsequent reads are synchronous.
 * @param key - Storage key (scoped under the stash app).
 * @param defaultValue - Fallback when no data exists in storage.
 * @returns A tuple of `{ data, error, loading }` and a setter.
 */
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

/** Convenience hook for reading/writing the interface config (view config, lightbox settings, etc.). */
export const useInterfaceLocalForage = () =>
  useLocalForage<IInterfaceConfig>("interface");

/** Convenience hook for reading/writing changelog acknowledgement state. */
export const useChangelogStorage = () =>
  useLocalForage<IChangelogConfig>("changelog");
