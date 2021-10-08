import React, { useEffect } from "react";

const useIndeterminate = (
  ref: React.RefObject<HTMLInputElement>,
  value: boolean | undefined
) => {
  useEffect(() => {
    if (ref.current) {
      // eslint-disable-next-line no-param-reassign
      ref.current.indeterminate = value === undefined;
    }
  }, [ref, value]);
};

export default useIndeterminate;
