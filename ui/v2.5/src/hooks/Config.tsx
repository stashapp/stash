import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";

const BLUR_IMG_KEY = "blurImage";
interface IContextProps {
  configuration?: GQL.ConfigDataFragment;
  loading?: boolean;
  imageBlurred: boolean;
}

export interface IContext extends IContextProps {
  enableImageBlur: () => void;
  disableImageBlur: () => void;
}

export const ConfigurationContext = React.createContext<IContext>({
  enableImageBlur: () => {},
  disableImageBlur: () => {},
  imageBlurred: false,
});

export const getInitialImageBlur = (): boolean => {
  const value = window.localStorage.getItem(BLUR_IMG_KEY);

  if (value === null) {
    return false;
  }

  return JSON.parse(value);
};

export const ConfigurationProvider: React.FC<IContextProps> = ({
  loading,
  configuration,
  imageBlurred,
  children,
}) => {
  const [blurImages, setBlurImages] = useState<boolean>(imageBlurred);

  const enableImageBlur = () => {
    window.localStorage.setItem(BLUR_IMG_KEY, JSON.stringify(true));
    setBlurImages(true);
  };

  const disableImageBlur = () => {
    window.localStorage.setItem(BLUR_IMG_KEY, JSON.stringify(false));
    setBlurImages(false);
  };

  return (
    <ConfigurationContext.Provider
      value={{
        configuration,
        imageBlurred: blurImages,
        loading,
        enableImageBlur,
        disableImageBlur,
      }}
    >
      {children}
    </ConfigurationContext.Provider>
  );
};
