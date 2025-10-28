/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react-hooks/exhaustive-deps */
import { debounce, DebouncedFunc, DebounceSettings } from "lodash-es";
import { useCallback, useRef } from "react";

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
