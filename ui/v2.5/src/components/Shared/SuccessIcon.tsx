import React from "react";
import { Icon } from "src/components/Shared";

interface ISuccessIconProps {
  className?: string;
}

const SuccessIcon: React.FC<ISuccessIconProps> = ({ className }) => (
  <Icon icon="check" className={className} color="#0f9960" />
);

export default SuccessIcon;
