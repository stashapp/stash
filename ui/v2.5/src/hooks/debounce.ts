/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react-hooks/exhaustive-deps */
import { DebounceSettings } from "lodash-es";
import debounce, { DebouncedFunc } from "lodash-es/debounce";
import React, { useCallback } from "react";

export function useDebounce<T extends (...args: any) => any>(
  fn: T,
  deps: React.DependencyList,
  wait?: number,
  options?: DebounceSettings
): DebouncedFunc<T> {
  return useCallback(debounce(fn, wait, options), [...deps, wait, options]);
}

// Convenience hook for use with state setters
export function useDebouncedSetState<S>(
  fn: React.Dispatch<React.SetStateAction<S>>,
  wait?: number,
  options?: DebounceSettings
): DebouncedFunc<React.Dispatch<React.SetStateAction<S>>> {
  return useDebounce(fn, [], wait, options);
}
