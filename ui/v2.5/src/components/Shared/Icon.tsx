import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconDefinition, SizeProp } from "@fortawesome/fontawesome-svg-core";
import { PatchComponent } from "src/patch";

interface IIcon {
  icon: IconDefinition;
  className?: string;
  color?: string;
  size?: SizeProp;
}

export const Icon: React.FC<IIcon> = PatchComponent(
  "Icon",
  ({ icon, className, color, size }) => (
    <FontAwesomeIcon
      icon={icon}
      className={`fa-icon ${className ?? ""}`}
      color={color}
      size={size}
    />
  )
);
