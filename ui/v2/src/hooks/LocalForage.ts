import localForage from "localforage";
import _ from "lodash";
import React from "react";

interface IInterfaceWallConfig {
}
export interface IInterfaceConfig {
  wall: IInterfaceWallConfig;
  queries: any;
}

type ValidTypes = IInterfaceConfig | undefined;

interface ILocalForage<T> {
  data: T;
  setData: React.Dispatch<React.SetStateAction<T>>;
  error: Error | null;
}

export function useInterfaceLocalForage(): ILocalForage<IInterfaceConfig | undefined> {
  const result = useLocalForage("interface");
  // Set defaults
  React.useEffect(() => {
    if (result.data === undefined) {
      result.setData({
        wall: {
          // nothing here currently
        },
        queries: {}
      });
    }
  });
  return result;
}

function useLocalForage(item: string): ILocalForage<ValidTypes> {
  const [json, setJson] = React.useState<ValidTypes>(undefined);

  const prevJson = React.useRef<ValidTypes>(undefined);
  React.useEffect(() => {
    async function runAsync() {
      if (typeof json !== "undefined" && !_.isEqual(json, prevJson.current)) {
        await localForage.setItem(item, JSON.stringify(json));
      }
      prevJson.current = json;
    }
    runAsync();
  });

  const [err, setErr] = React.useState(null);
  React.useEffect(() => {
    async function runAsync() {
      try {
        const serialized = await localForage.getItem<any>(item);
        const parsed = JSON.parse(serialized);
        if (typeof json === "undefined" && !Object.is(parsed, null)) {
          setErr(null);
          setJson(parsed);
        }
      } catch (error) {
        setErr(error);
      }
    }
    runAsync();
  });

  return {data: json, setData: setJson, error: err};
}
