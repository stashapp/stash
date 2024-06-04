import cx from "classnames";
import React, { useContext } from "react";
import { ConfigurationContext } from "src/hooks/Config";

type IImageProps = JSX.IntrinsicElements["img"];

export const BLURRED_CLASSNAME = "blurred";

const StashImage: React.FC<IImageProps> = ({ className, ...rest }) => {
  const configContext = useContext(ConfigurationContext);

  return (
    <StashImage
      className={cx(
        className,
        configContext.imageBlurred ? BLURRED_CLASSNAME : ""
      )}
      {...rest}
    />
  );
};

export default StashImage;
