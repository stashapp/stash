import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconDefinition, SizeProp } from "@fortawesome/fontawesome-svg-core";

interface IIcon {
  icon: IconDefinition;
  className?: string;
  color?: string;
  size?: SizeProp;
}

const Icon: React.FC<IIcon> = ({ icon, className, color, size }) => (
  <FontAwesomeIcon
    icon={icon}
    className={`fa-icon ${className ?? ""}`}
    color={color}
    size={size}
  />
);

export default Icon;
