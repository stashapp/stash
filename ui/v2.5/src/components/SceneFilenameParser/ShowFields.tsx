import React, { useState } from 'react';
import {
  Button,
  Collapse
} from 'react-bootstrap';
import { Icon } from 'src/components/Shared';

interface IShowFieldsProps {
  fields: Map<string, boolean>;
  onShowFieldsChanged: (fields: Map<string, boolean>) => void;
}

export const ShowFields = (props: IShowFieldsProps) => {
  const [open, setOpen] = useState(false);

  function handleClick(label: string) {
    const copy = new Map<string, boolean>(props.fields);
    copy.set(label, !props.fields.get(label));
    props.onShowFieldsChanged(copy);
  }

  const fieldRows = [...props.fields.entries()].map(([label, enabled]) => (
    <Button
      className="minimal d-block"
      key={label}
      onClick={() => {
        handleClick(label);
      }}
    >
      <Icon icon={enabled ? "check" : "times"} />
      <span>{label}</span>
    </Button>
  ));

  return (
    <div>
      <Button onClick={() => setOpen(!open)} className="minimal">
        <Icon icon={open ? "chevron-down" : "chevron-right"} />
        <span>Display fields</span>
      </Button>
      <Collapse in={open}>
        <div>{fieldRows}</div>
      </Collapse>
    </div>
  );
}
