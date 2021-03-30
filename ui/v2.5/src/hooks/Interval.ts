import { useEffect, useRef, useState } from "react";
import noop from "lodash/noop";

const useInterval = (
  callback: () => void,
  delay: number | null = 5000
): (() => void)[] => {
  const savedCallback = useRef<() => void>();
  const savedIntervalId = useRef<NodeJS.Timeout>();
  const [savedDelay, setSavedDelay] = useState<number | null>(delay);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    setSavedDelay(delay);
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
    }

    if (savedDelay !== null) {
      savedIntervalId.current = setInterval(tick, savedDelay);
    }
  }

  useEffect(() => {
    cancel();

    const tick = () => {
      if (savedCallback.current) savedCallback.current();
    };

    if (delay !== null) {
      savedIntervalId.current = setInterval(tick, delay);
      return cancel;
    }
  }, [callback, delay]);

  return delay ? [cancel, reset] : [noop, noop];
};

export default useInterval;
