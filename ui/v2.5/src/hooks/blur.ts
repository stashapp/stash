import cx from "classnames";
import { useContext } from "react";
import { BLURRED_CLASSNAME } from "src/components/Shared/StashImage";
import { ConfigurationContext } from "./Config";

export const useImageBlur = () => {
  const { imageBlurred } = useContext(ConfigurationContext);
  const blurClassName = (name: string): string => {
    return cx(name, imageBlurred ? BLURRED_CLASSNAME : "");
  };

  return { imageBlurred, blurClassName };
};
