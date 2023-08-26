import { useEffect } from "react";

export function useScrollToTopOnMount() {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);
}
