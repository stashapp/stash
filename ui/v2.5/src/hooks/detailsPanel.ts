import { useEffect, useState } from "react";

function shouldLoadStickyHeader() {
  return document.documentElement.scrollTop > 50;
}

export function useLoadStickyHeader() {
  const [load, setLoad] = useState(shouldLoadStickyHeader());

  useEffect(() => {
    const onScroll = () => {
      setLoad(shouldLoadStickyHeader());
    };

    window.addEventListener("scroll", onScroll);
    return () => {
      window.removeEventListener("scroll", onScroll);
    };
  }, []);

  return load;
}
