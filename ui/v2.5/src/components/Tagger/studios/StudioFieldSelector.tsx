import { faCheck, faList, faTimes } from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Row, Col } from "react-bootstrap";
import { useIntl } from "react-intl";

import { ModalComponent } from "../../Shared/Modal";
import { Icon } from "../../Shared/Icon";
import TextUtils from "src/utils/text";

interface IProps {
  fields: string[];
  show: boolean;
  excludedFields: string[];
  onSelect: (fields: string[]) => void;
}

const StudioFieldSelect: React.FC<IProps> = ({
  fields,
  show,
  excludedFields,
  onSelect,
}) => {
  const intl = useIntl();
  const [excluded, setExcluded] = useState<Record<string, boolean>>(
    excludedFields.reduce((dict, field) => ({ ...dict, [field]: true }), {})
  );

  const toggleField = (name: string) =>
    setExcluded({
      ...excluded,
      [name]: !excluded[name],
    });

  const renderField = (name: string) => (
    <Col xs={6} className="mb-1" key={name}>
      <Button
        onClick={() => toggleField(name)}
        variant="secondary"
        className={excluded[name] ? "text-muted" : "text-success"}
      >
        <Icon icon={excluded[name] ? faTimes : faCheck} />
      </Button>
      <span className="ml-3">{TextUtils.capitalize(name)}</span>
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
      <Row>{fields.map((f) => renderField(f))}</Row>
    </ModalComponent>
  );
};

export default StudioFieldSelect;
