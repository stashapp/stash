import { useEffect } from "react";

const usePageVisibility = (
  visibilityChangeCallback: (hidden: boolean) => void
): void => {
  useEffect(() => {
    const callback = () => visibilityChangeCallback(document.hidden);
    document.addEventListener("visibilitychange", callback);

    return () => {
      document.removeEventListener("visibilitychange", callback);
    };
  }, [visibilityChangeCallback]);
};

export default usePageVisibility;
