import React from "react";
import Icon from "src/components/Shared/Icon";

interface ISuccessIconProps {
  className?: string;
}

const SuccessIcon: React.FC<ISuccessIconProps> = ({ className }) => (
  <Icon icon={["far", "check-circle"]} className={className} color="#0f9960" />
);

export default SuccessIcon;
