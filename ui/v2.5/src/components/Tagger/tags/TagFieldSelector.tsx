import { faCheck, faList, faTimes } from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Row, Col } from "react-bootstrap";
import { useIntl } from "react-intl";

import { ModalComponent } from "../../Shared/Modal";
import { Icon } from "../../Shared/Icon";
import { TAG_FIELDS } from "../constants";

interface IProps {
  show: boolean;
  excludedFields: string[];
  onSelect: (fields: string[]) => void;
}

const TagFieldSelect: React.FC<IProps> = ({
  show,
  excludedFields,
  onSelect,
}) => {
  const intl = useIntl();
  const [excluded, setExcluded] = useState<Record<string, boolean>>(
    excludedFields
      .filter((field) => TAG_FIELDS.includes(field))
      .reduce((dict, field) => ({ ...dict, [field]: true }), {})
  );

  const toggleField = (field: string) =>
    setExcluded({
      ...excluded,
      [field]: !excluded[field],
    });

  const renderField = (field: string) => (
    <Col xs={6} className="mb-1" key={field}>
      <Button
        onClick={() => toggleField(field)}
        variant="secondary"
        className={excluded[field] ? "text-muted" : "text-success"}
      >
        <Icon icon={excluded[field] ? faTimes : faCheck} />
      </Button>
      <span className="ml-3">{intl.formatMessage({ id: field })}</span>
    </Col>
  );

  return (
    <ModalComponent
      show={show}
      icon={faList}
      dialogClassName="FieldSelect"
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: () =>
          onSelect(Object.keys(excluded).filter((f) => excluded[f])),
      }}
    >
      <h4>Select tagged fields</h4>
      <div className="mb-2">
        These fields will be tagged by default. Click the button to toggle.
      </div>
      <Row>{TAG_FIELDS.map((f) => renderField(f))}</Row>
    </ModalComponent>
  );
};

export default TagFieldSelect;
