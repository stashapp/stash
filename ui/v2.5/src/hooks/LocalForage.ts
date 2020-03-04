import localForage from "localforage";
import _ from "lodash";
import React, { Dispatch, SetStateAction, useEffect } from "react";

interface IInterfaceWallConfig {}
interface IInterfaceQueryConfig {
  filter: string;
  itemsPerPage: number;
  currentPage: number;
}

export interface IInterfaceConfig {
  wall?: IInterfaceWallConfig;
  queries?: Record<string, IInterfaceQueryConfig>;
}

type ValidTypes = IInterfaceConfig;
type Key = "interface";

interface ILocalForage<T> {
  data?: T;
  error: Error | null;
  loading: boolean;
}

const Loading: Record<string, boolean> = {};
const Cache: Record<string, ValidTypes> = {};

function useLocalForage(
  key: Key
): [ILocalForage<ValidTypes>, Dispatch<SetStateAction<ValidTypes>>] {
  const [error, setError] = React.useState(null);
  const [data, setData] = React.useState(Cache[key]);
  const [loading, setLoading] = React.useState(Loading[key]);

  useEffect(() => {
    async function runAsync() {
      try {
        const serialized = await localForage.getItem<string>(key);
        const parsed = JSON.parse(serialized);
        if (!Object.is(parsed, null)) {
          setData(parsed);
          Cache[key] = parsed;
        }
        else {
          setData({});
          Cache[key] = {};
        }
        setError(null);
      } catch (err) {
        setError(err);
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
  }, [loading, data, key]);

  useEffect(() => {
    if (!_.isEqual(Cache[key], data)) {
      Cache[key] = _.merge(Cache[key], data);
      localForage.setItem(key, JSON.stringify(Cache[key]));
    }
  });

  const isLoading = loading || loading === undefined;

  return [{ data, error, loading: isLoading }, setData];
}

export function useInterfaceLocalForage(): [
  ILocalForage<IInterfaceConfig>,
  Dispatch<SetStateAction<IInterfaceConfig>>
] {
  return useLocalForage("interface");
}
