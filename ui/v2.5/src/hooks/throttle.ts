/* biome-ignore-all lint/suspicious/noExplicitAny: don't know how to change this to be properly typed */
/* biome-ignore-all lint/correctness/useExhaustiveDependencies: seems like a false positive for function arguments */
import { DebouncedFunc, DebounceSettings, throttle } from "lodash-es";
import { useCallback, useRef } from "react";

export function useThrottle<T extends (...args: any) => any>(
  fn: T,
  wait?: number,
  options?: DebounceSettings
): DebouncedFunc<T> {
  const func = useRef<T>(fn);
  func.current = fn;
  return useCallback(
    throttle(
      function (this: any) {
        return func.current.apply(this, arguments as any);
      },
      wait,
      options
    ),
    [wait, options?.leading, options?.trailing, options?.maxWait]
  );
}
