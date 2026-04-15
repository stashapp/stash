import React, { useEffect } from "react";

export const useOnOutsideClick = (
  ref: React.RefObject<HTMLElement>,
  callback?: () => void,
  excludeClassName?: string
) => {
  useEffect(() => {
    if (!callback) return;

    /**
     * Alert if clicked on outside of element
     */
    function handleClickOutside(event: MouseEvent) {
      if (
        ref.current &&
        event.target instanceof Node &&
        !ref.current.contains(event.target) &&
        !(
          excludeClassName &&
          (event.target as HTMLElement).closest(`.${excludeClassName}`)
        )
      ) {
        callback?.();
      }
    }
    // Bind the event listener
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      // Unbind the event listener on clean up
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [ref, callback, excludeClassName]);
};
