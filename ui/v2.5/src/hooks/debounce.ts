/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react-hooks/exhaustive-deps */
import { debounce, DebouncedFunc, DebounceSettings } from "lodash-es";
import { useCallback, useRef, useState } from "react";

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
  initialValue: T,
  setValue: (v: T) => void,
  wait?: number
): [T, (v: T) => void, (v: T) => void] {
  const [displayedState, setDisplayedState] = useState(initialValue);

  const debouncedSetValue = useDebounce(setValue, wait);
  const onChange = useCallback(
    (input: T) => {
      setDisplayedState(input);
      debouncedSetValue(input);
    },
    [debouncedSetValue, setDisplayedState]
  );

  const setInstant = useCallback(
    (v: T) => {
      setDisplayedState(v);
      setValue(v);
    },
    [setValue]
  );

  return [displayedState, onChange, setInstant];
}
