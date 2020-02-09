import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconName } from "@fortawesome/fontawesome-svg-core";

interface IIcon {
  icon: IconName;
  className?: string;
  color?: string;
}

const Icon: React.FC<IIcon> = ({ icon, className, color }) => (
  <FontAwesomeIcon
    icon={icon}
    className={`fa-icon ${className}`}
    color={color}
  />
);

export default Icon;
