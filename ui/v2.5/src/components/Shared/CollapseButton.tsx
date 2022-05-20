import React, { useState } from "react";
import { Button, Collapse } from "react-bootstrap";
import Icon from "src/components/Shared/Icon";

interface IProps {
  text: string;
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
        <Icon icon={open ? "chevron-down" : "chevron-right"} />
        <span>{props.text}</span>
      </Button>
      <Collapse in={open}>
        <div>{props.children}</div>
      </Collapse>
    </div>
  );
};
