import {
  faChevronDown,
  faChevronRight,
  faChevronUp,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Collapse, CollapseProps } from "react-bootstrap";
import { Icon } from "./Icon";

interface IProps {
  className?: string;
  text: React.ReactNode;
  collapseProps?: Partial<CollapseProps>;
}

export const CollapseButton: React.FC<React.PropsWithChildren<IProps>> = (
  props: React.PropsWithChildren<IProps>
) => {
  const [open, setOpen] = useState(false);

  return (
    <div className={props.className}>
      <div className="collapse-header">
        <Button
          onClick={() => setOpen(!open)}
          className="minimal collapse-button"
        >
          <Icon icon={open ? faChevronDown : faChevronRight} fixedWidth />
          <span>{props.text}</span>
        </Button>
      </div>
      <Collapse in={open} {...props.collapseProps}>
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
