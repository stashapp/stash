import {
  faChevronDown,
  faChevronRight,
  faChevronUp,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Collapse } from "react-bootstrap";
import { Icon } from "./Icon";

interface IProps {
  text: React.ReactNode;
}

export const CollapseButton: React.FC<React.PropsWithChildren<IProps>> = (
  props: React.PropsWithChildren<IProps>
) => {
  const [open, setOpen] = useState(false);

  return (
    <div>
      <Button
        onClick={() => setOpen(!open)}
        className="minimal collapse-button"
      >
        <Icon icon={open ? faChevronDown : faChevronRight} />
        <span>{props.text}</span>
      </Button>
      <Collapse in={open}>
        <div>{props.children}</div>
      </Collapse>
    </div>
  );
};

export const ExpandCollapseButton: React.FC<{
  collapsed: boolean;
  setCollapsed: (collapsed: boolean) => void;
}> = ({ collapsed, setCollapsed }) => {
  const buttonIcon = collapsed ? faChevronDown : faChevronUp;

  return (
    <span className="detail-expand-collapse">
      <Button
        className="minimal expand-collapse"
        onClick={() => setCollapsed(!collapsed)}
      >
        <Icon className="fa-fw" icon={buttonIcon} />
      </Button>
    </span>
  );
};
