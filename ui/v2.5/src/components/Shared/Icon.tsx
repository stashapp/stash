import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconProp, library } from "@fortawesome/fontawesome-svg-core";
import { faStar as fasStar } from "@fortawesome/free-solid-svg-icons";
import { faStar as farStar } from "@fortawesome/free-regular-svg-icons";

// need these to use far and fas styles of stars
library.add(fasStar, farStar);

interface IIcon {
  icon: IconProp;
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
