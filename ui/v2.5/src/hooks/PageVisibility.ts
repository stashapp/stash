import { useEffect, useRef } from "react";

const usePageVisibility = (
  visibilityChangeCallback: () => void
) : void => {
  const savedVisibilityChangedCallback = useRef<() => void>();

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

    if (typeof document.addEventListener === "undefined" || hidden === undefined) {
      // this browser doesn't have support for modern event listeners or the Page Visibility API
      return;
    }

    savedVisibilityChangedCallback.current = visibilityChangeCallback;

    document.addEventListener(visibilityChange, savedVisibilityChangedCallback.current);

    return () => {
      if (savedVisibilityChangedCallback.current) {
        document.removeEventListener(visibilityChange, savedVisibilityChangedCallback.current);
      }
    }
  }, [visibilityChangeCallback]);
}

export default usePageVisibility;
