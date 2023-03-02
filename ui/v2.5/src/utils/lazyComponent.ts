import { ComponentType, lazy } from "react";

interface ILazyComponentError {
  __lazyComponentError?: true;
}

export const isLazyComponentError = (e: unknown) => {
  return !!(e as ILazyComponentError).__lazyComponentError;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const lazyComponent = <T extends ComponentType<any>>(
  factory: Parameters<typeof lazy<T>>[0]
) => {
  return lazy<T>(async () => {
    try {
      return await factory();
    } catch (e) {
      // set flag to identify lazy component loading errors
      (e as ILazyComponentError).__lazyComponentError = true;
      throw e;
    }
  });
};
