import { useEffect, useMemo } from "react";

const useScript = (urls: string | string[], condition?: boolean) => {
  const urlArray = useMemo(() => {
    if (!Array.isArray(urls)) {
      return [urls];
    }

    return urls;
  }, [urls]);

  useEffect(() => {
    const scripts = urlArray.map((url) => {
      const script = document.createElement("script");

      script.src = url;
      script.defer = true;
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
