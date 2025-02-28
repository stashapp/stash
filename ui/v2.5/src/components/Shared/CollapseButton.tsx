import {
  faChevronDown,
  faChevronRight,
  faChevronUp,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Collapse } from "react-bootstrap";
import { Icon } from "./Icon";

interface IProps {
  className?: string;
  text: React.ReactNode;
}

export const CollapseButton: React.FC<React.PropsWithChildren<IProps>> = (
  props: React.PropsWithChildren<IProps>
) => {
  const [open, setOpen] = useState(false);

  return (
    <div className={props.className}>
      <Button
        onClick={() => setOpen(!open)}
        className="minimal collapse-button"
      >
        <Icon icon={open ? faChevronDown : faChevronRight} fixedWidth />
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
        <Icon icon={buttonIcon} fixedWidth />
      </Button>
    </span>
  );
};
