import {
  faCheck,
  faChevronDown,
  faChevronRight,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Collapse } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";

interface IShowFieldsProps {
  fields: Map<string, boolean>;
  onShowFieldsChanged: (fields: Map<string, boolean>) => void;
}

export const ShowFields: React.FC<IShowFieldsProps> = (props) => {
  const intl = useIntl();
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
      <Icon icon={enabled ? faCheck : faTimes} />
      <span>{label}</span>
    </Button>
  ));

  return (
    <div>
      <Button onClick={() => setOpen(!open)} className="minimal">
        <Icon icon={open ? faChevronDown : faChevronRight} />
        <span>
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.display_fields",
          })}
        </span>
      </Button>
      <Collapse in={open}>
        <div>{fieldRows}</div>
      </Collapse>
    </div>
  );
};
