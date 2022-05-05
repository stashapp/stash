import { useEffect, useRef } from "react";

const usePageVisibility = (
  visibilityChangeCallback: (hidden: boolean) => void
): void => {
  const savedVisibilityChangedCallback = useRef<(hidden: boolean) => void>();

  useEffect(() => {
    // resolve event names for different browsers
    let hidden = "";
    let visibilityChange = "";

    if (typeof document.hidden !== "undefined") {
      hidden = "hidden";
      visibilityChange = "visibilitychange";
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } else if (typeof (document as any).msHidden !== "undefined") {
      hidden = "msHidden";
      visibilityChange = "msvisibilitychange";
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } else if (typeof (document as any).webkitHidden !== "undefined") {
      hidden = "webkitHidden";
      visibilityChange = "webkitvisibilitychange";
    }

    if (
      typeof document.addEventListener === "undefined" ||
      hidden === undefined
    ) {
      // this browser doesn't have support for modern event listeners or the Page Visibility API
      return;
    }

    savedVisibilityChangedCallback.current = visibilityChangeCallback;

    function fireCallback() {
      const callback = savedVisibilityChangedCallback.current;
      if (callback) {
        const isHidden = document.visibilityState !== "visible";
        callback(isHidden);
      }
    }

    document.addEventListener(visibilityChange, fireCallback);

    return () => {
      document.removeEventListener(visibilityChange, fireCallback);
    };
  }, [visibilityChangeCallback]);
};

export default usePageVisibility;
