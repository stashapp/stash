import { useEffect } from "react";

const useScript = (url: string, condition?: boolean) => {
  useEffect(() => {
    const script = document.createElement("script");

    script.src = url;
    script.async = true;

    if (condition) {
      document.head.appendChild(script);
    }

    return () => {
      if (condition) {
        document.head.removeChild(script);
      }
    };
  }, [url, condition]);
};

export default useScript;
