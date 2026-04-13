import React, { useState } from "react";
import { Badge, Button, Card, Collapse, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import FieldSelector from "./FieldSelector";
import { Icon } from "../Shared/Icon";
import { faCog } from "@fortawesome/free-solid-svg-icons";

interface ITaggerConfigProps {
  show: boolean;
  excludedFields: string[];
  onFieldsChange: (fields: string[]) => void;
  fields: string[];
  entityName: string;
  extraConfig?: React.ReactNode;
}

const TaggerConfig: React.FC<ITaggerConfigProps> = ({
  show,
  excludedFields,
  onFieldsChange,
  fields,
  entityName,
  extraConfig,
}) => {
  const [showExclusionModal, setShowExclusionModal] = useState(false);

  const handleFieldSelect = (selectedFields: string[]) => {
    onFieldsChange(selectedFields);
    setShowExclusionModal(false);
  };

  return (
    <>
      <Collapse in={show}>
        <Card>
          <div className="row">
            <h4 className="col-12">
              <FormattedMessage id="configuration" />
            </h4>
            <hr className="w-100" />
            <div className="col-md-6">
              {extraConfig}
              <Form.Group controlId="excluded-fields">
                <h6>
                  <FormattedMessage id="tagger.config.excluded_fields" />
                </h6>
                <span>
                  {excludedFields.length > 0 ? (
                    excludedFields.map((f) => (
                      <Badge variant="secondary" className="tag-item" key={f}>
                        <FormattedMessage id={f} />
                      </Badge>
                    ))
                  ) : (
                    <FormattedMessage id="tagger.config.no_fields_are_excluded" />
                  )}
                </span>
                <Form.Text>
                  <FormattedMessage
                    id="tagger.config.fields_will_not_be_changed"
                    values={{ entity: entityName }}
                  />
                </Form.Text>
                <Button
                  onClick={() => setShowExclusionModal(true)}
                  className="mt-2"
                >
                  <FormattedMessage id="tagger.config.edit_excluded_fields" />
                </Button>
              </Form.Group>
            </div>
          </div>
        </Card>
      </Collapse>
      <FieldSelector
        show={showExclusionModal}
        fields={fields}
        onSelect={handleFieldSelect}
        excludedFields={excludedFields}
      />
    </>
  );
};

export default TaggerConfig;

export const ConfigButton: React.FC<{
  onClick: () => void;
  showConfig: boolean;
}> = ({ onClick, showConfig }) => {
  const intl = useIntl();

  const showHideConfigId = showConfig
    ? "actions.hide_configuration"
    : "actions.show_configuration";

  return (
    <Button
      onClick={onClick}
      title={intl.formatMessage({ id: showHideConfigId })}
    >
      <Icon className="fa-fw" icon={faCog} />
    </Button>
  );
};
