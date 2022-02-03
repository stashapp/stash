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
