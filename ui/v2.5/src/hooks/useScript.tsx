import { useEffect } from "react";

const useScript = (url: string, condition?: boolean) => {
  useEffect(() => {
    const script = document.createElement("script");

    script.src = url;
    script.async = true;

    if (condition) {
      document.body.appendChild(script);
    }

    return () => {
      if (condition) {
        document.body.removeChild(script);
      }
    };
  }, [url, condition]);
};

export default useScript;
