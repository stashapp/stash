import { faCheckCircle } from "@fortawesome/free-regular-svg-icons";
import React from "react";
import { Icon } from "./Icon";

interface ISuccessIconProps {
  className?: string;
}

export const SuccessIcon: React.FC<ISuccessIconProps> = ({ className }) => (
  <Icon icon={faCheckCircle} className={className} color="#0f9960" />
);
