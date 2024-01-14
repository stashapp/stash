import {
  faChevronDown,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, ButtonGroup, Collapse } from "react-bootstrap";
import { Icon } from "./Icon";

interface IProps {
  text: string;
  rightControls?: React.ReactNode;
}

export const CollapseButton: React.FC<React.PropsWithChildren<IProps>> = (
  props: React.PropsWithChildren<IProps>
) => {
  const [open, setOpen] = useState(false);

  return (
    <div className="collapse-button-container">
      <ButtonGroup className="collapse-button-group">
        <Button
          onClick={() => setOpen(!open)}
          className="minimal collapse-button"
        >
          <span>
            <Icon icon={open ? faChevronDown : faChevronRight} />
            <span>{props.text}</span>
          </span>
        </Button>
        {props.rightControls}
      </ButtonGroup>
      <Collapse in={open}>
        <div>{props.children}</div>
      </Collapse>
    </div>
  );
};
