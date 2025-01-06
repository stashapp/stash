import { useRef, useEffect } from "react";

const useFocus = () => {
  const htmlElRef = useRef<HTMLInputElement | null>(null);
  const setFocus = () => {
    const currentEl = htmlElRef.current;
    if (currentEl) {
      currentEl.focus();
    }
  };

  // eslint-disable-next-line no-undef
  return [htmlElRef, setFocus] as const;
};

// focuses on the element only once on mount
export const useFocusOnce = (active?: boolean, override?: boolean) => {
  const [htmlElRef, setFocus] = useFocus();
  const focused = useRef(false);

  useEffect(() => {
    if ((!focused.current || override) && active) {
      setFocus();
      focused.current = true;
    }
  }, [setFocus, active, override]);

  return [htmlElRef, setFocus] as const;
};

export default useFocus;
