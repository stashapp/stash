import { useEffect, useMemo, useState } from "react";

const useScript = (urls: string | string[], condition: boolean = true) => {
  // array of booleans to track the loading state of each script
  const [loadStates, setLoadStates] = useState<boolean[]>();

  const urlArray = useMemo(() => {
    if (!Array.isArray(urls)) {
      return [urls];
    }

    return urls;
  }, [urls]);

  useEffect(() => {
    if (condition) {
      setLoadStates(urlArray.map(() => false));
    }

    const scripts = urlArray.map((url) => {
      const script = document.createElement("script");

      script.src = url;
      script.async = false;
      script.defer = true;

      function onLoad() {
        setLoadStates((prev) =>
          prev!.map((state, i) => (i === urlArray.indexOf(url) ? true : state))
        );
      }
      script.addEventListener("load", onLoad);
      script.addEventListener("error", onLoad); // handle error as well

      return script;
    });

    if (condition) {
      scripts.forEach((script) => {
        document.head.appendChild(script);
      });
    }

    return () => {
      if (condition) {
        scripts.forEach((script) => {
          document.head.removeChild(script);
        });
      }
    };
  }, [urlArray, condition]);

  return (
    condition &&
    loadStates &&
    (loadStates.length === 0 || loadStates.every((state) => state))
  );
};

export const useCSS = (urls: string | string[], condition?: boolean) => {
  const urlArray = useMemo(() => {
    if (!Array.isArray(urls)) {
      return [urls];
    }

    return urls;
  }, [urls]);

  useEffect(() => {
    const links = urlArray.map((url) => {
      const link = document.createElement("link");

      link.href = url;
      link.rel = "stylesheet";
      link.type = "text/css";
      return link;
    });

    if (condition) {
      links.forEach((link) => {
        document.head.appendChild(link);
      });
    }

    return () => {
      if (condition) {
        links.forEach((link) => {
          document.head.removeChild(link);
        });
      }
    };
  }, [urlArray, condition]);
};

export default useScript;
