import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconProp, SizeProp, library } from "@fortawesome/fontawesome-svg-core";
import { faStar as fasStar } from "@fortawesome/free-solid-svg-icons";
import {
  faCheckCircle as farCheckCircle,
  faStar as farStar,
} from "@fortawesome/free-regular-svg-icons";

// need these to use far and fas styles of stars
library.add(fasStar, farStar, farCheckCircle);

interface IIcon {
  icon: IconProp;
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
