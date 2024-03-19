/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react-hooks/exhaustive-deps */
import { debounce, DebouncedFunc, DebounceSettings } from "lodash-es";
import { useCallback, useEffect, useRef, useState } from "react";

export function useDebounce<T extends (...args: any) => any>(
  fn: T,
  wait?: number,
  options?: DebounceSettings
): DebouncedFunc<T> {
  const func = useRef<T>(fn);
  func.current = fn;
  return useCallback(
    debounce(
      function (this: any) {
        return func.current.apply(this, arguments as any);
      },
      wait,
      options
    ),
    [wait, options?.leading, options?.trailing, options?.maxWait]
  );
}

export function useDebouncedState<T>(
  value: T,
  setValue: (value: T) => void,
  wait: number = 250,
  options?: DebounceSettings
) {
  const [stageValue, setStageValue] = useState<T>(value);

  useEffect(() => {
    setStageValue(value);
  }, [value]);

  const setValueCallback = useDebounce(
    () => {
      setValue(stageValue);
    },
    wait,
    options
  );

  function onSetStageValue(v: T) {
    setStageValue(v);
    setValueCallback();
  }

  return [stageValue, onSetStageValue] as [T, (value: T) => void];
}
