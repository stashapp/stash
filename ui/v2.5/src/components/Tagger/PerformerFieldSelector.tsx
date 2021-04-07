import React, { useState } from "react";
import { Button } from "react-bootstrap";

import { Modal, Icon } from "src/components/Shared";
import { TextUtils } from "src/utils";

interface IProps {
  fields: string[];
  show: boolean;
  excludedFields: string[];
  onSelect: (fields: string[]) => void;
}

const PerformerFieldSelect: React.FC<IProps> = ({ fields, show, excludedFields, onSelect }) => {
  const [excluded, setExcluded] = useState<Record<string, boolean>>(excludedFields.reduce(
    (dict, field) => (
      { ...dict, [field]: true }
    ), {}));

  const toggleField = (name: string) => (
    setExcluded({
      ...excluded,
      [name]: !excluded[name],
    })
  );

  const renderField = (name: string) => (
    <div className="mb-1" key={name}>
      <Button onClick={() => toggleField(name)} variant="secondary" className={excluded[name] ? 'text-muted' : 'text-success'}>
        <Icon icon={excluded[name] ? 'times' : 'check'} />
      </Button>
      <span className="ml-3">
        {TextUtils.capitalize(name)}
      </span>
    </div>
  );

  return (
    <Modal
      show={show}
      icon="list"
      dialogClassName="FieldSelect"
      accept={{
        text: "Save",
        onClick: () => onSelect(Object.keys(excluded).filter(f => excluded[f]))
      }}>
      <h4>Select tagged fields</h4>
      <div className="mb-2">These fields will be tagged by default. Click the button to toggle.</div>
      { fields.map(f => renderField(f)) }
    </Modal>
  );
}

export default PerformerFieldSelect;
