import React from "react";
import {
  FontAwesomeIcon,
  FontAwesomeIconProps,
} from "@fortawesome/react-fontawesome";
import { PatchComponent } from "src/patch";

export const Icon: React.FC<FontAwesomeIconProps> = PatchComponent(
  "Icon",
  (props) => (
    <FontAwesomeIcon
      {...props}
      className={`fa-icon ${props.className ?? ""}`}
    />
  )
);
