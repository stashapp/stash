import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import {
  faChevronLeft,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";

export const CollapseDivider: React.FC<{
  collapsed: boolean;
  setCollapsed: (v: boolean) => void;
}> = ({ collapsed, setCollapsed }) => {
  const icon = collapsed ? faChevronRight : faChevronLeft;

  return (
    <div className={cx("collapse-divider", { collapsed })}>
      <Button onClick={() => setCollapsed(!collapsed)}>
        <div>
          <Icon className="fa-fw" icon={icon} />
        </div>
        <div></div>
      </Button>
    </div>
  );
};
