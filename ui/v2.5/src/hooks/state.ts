import React, { useCallback, Dispatch, SetStateAction } from "react";

// useInitialState is an extension of the useState hook.
// It maintains a state, but additionally exposes a setInitialState function.
// When setInitialState is called, the current state is only updated if the current
// state is unchanged from the initial state. This means that the current state will
// only be updated if explicitly called, or if the initial state is changed and the current
// state is not dirty.
export function useInitialState<T>(
  initialValue: T
): [T, Dispatch<SetStateAction<T>>, Dispatch<T>] {
  const [, setInitialValueInternal] = React.useState<T>(initialValue);
  const [value, setValue] = React.useState<T>(initialValue);

  const setInitialValue = useCallback((v: T) => {
    setInitialValueInternal((currentInitial) => {
      if (v === currentInitial) {
        return currentInitial;
      }

      setValue((currentValue) => {
        if (currentInitial === currentValue) {
          return v;
        }

        return currentValue;
      });

      return v;
    });
  }, []);

  return [value, setValue, setInitialValue];
}

// useCompare is a hook that returns true if the value has changed since the last render.
export function useCompare<T>(val: T) {
  const prevVal = usePrevious(val);
  return prevVal !== val;
}

// usePrevious is a hook that returns the previous value of a variable.
export function usePrevious<T>(value: T) {
  const ref = React.useRef<T>();
  React.useEffect(() => {
    ref.current = value;
  }, [value]);
  return ref.current;
}
