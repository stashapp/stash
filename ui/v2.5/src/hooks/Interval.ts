import { useEffect, useRef } from "react";
import noop from "lodash/noop";

const useInterval = (
  callback: () => void,
  delay: number | null = 5000
): (() => void) => {
  const savedCallback = useRef<() => void>();
  const savedIntervalId = useRef<NodeJS.Timeout>();

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  const cancel = () => {
    const intervalId = savedIntervalId.current;
    if (intervalId) {
      savedIntervalId.current = undefined;
      clearInterval(intervalId);
    }
  };

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

  return delay ? cancel : noop;
};

export default useInterval;
