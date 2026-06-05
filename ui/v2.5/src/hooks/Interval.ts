import { useEffect, useRef, useState } from "react";

const MIN_VALID_INTERVAL = 1000;

function noop() {}

// it's possible that the warning is correct, but we'd need to test it explcitly
// disabling warning for now
// biome-ignore-start lint/correctness/useExhaustiveDependencies: intentional
const useInterval = (
  callback: () => void,
  delay: number | null = 5000
): (() => void)[] => {
  const savedCallback = useRef<() => void>();
  const savedIntervalId = useRef<number>();
  const [savedDelay, setSavedDelay] = useState<number | null>(delay);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    let validDelay: number | null;
    if (delay !== null) {
      validDelay = delay >= MIN_VALID_INTERVAL ? delay : MIN_VALID_INTERVAL;
    } else {
      validDelay = delay;
    }

    setSavedDelay(validDelay);
  }, [delay]);

  const cancel = () => {
    const intervalId = savedIntervalId.current;
    if (intervalId) {
      savedIntervalId.current = undefined;
      clearInterval(intervalId);
    }
  };

  const reset = () => {
    cancel();

    const tick = () => {
      if (savedCallback.current) savedCallback.current();
    };

    if (savedDelay !== null) {
      savedIntervalId.current = setInterval(tick, savedDelay);
    }
  };

  useEffect(() => {
    cancel();

    const tick = () => {
      if (savedCallback.current) savedCallback.current();
    };

    if (savedDelay !== null) {
      savedIntervalId.current = setInterval(tick, savedDelay);
      return cancel;
    }
  }, [callback, savedDelay]);

  return delay ? [cancel, reset] : [noop, noop];
};
// biome-ignore-end lint/correctness/useExhaustiveDependencies: intentional

export default useInterval;
