import { ComponentType, lazy } from "react";

interface ILazyComponentError {
  __lazy_component_error?: true;
}

export const is_lazy_component_error = (e: unknown) => {
  return !!(e as ILazyComponentError).__lazy_component_error;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const lazy_component = <T extends ComponentType<any>>(
  factory: Parameters<typeof lazy<T>>[0]
) => {
  return lazy<T>(async () => {
    try {
      return await factory();
    } catch (e) {
      // set flag to identify lazy component loading errors
      (e as ILazyComponentError).__lazy_component_error = true;
      throw e;
    }
  });
};
