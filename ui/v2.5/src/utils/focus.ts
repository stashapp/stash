import { useRef } from "react";

const useFocus = () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const htmlElRef = useRef<any>();
  const setFocus = () => {
    const currentEl = htmlElRef.current;
    if (currentEl) {
      currentEl.focus();
    }
  };

  return [htmlElRef, setFocus] as const;
};

export default useFocus;
