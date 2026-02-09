import React, { Dispatch, useState } from "react";
import { Badge, Button, Card, Collapse, Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { useConfigurationContext } from "src/hooks/Config";

import { ITaggerConfig } from "./constants";
import FieldSelector from "./FieldSelector";

interface ITaggerConfigProps {
  show: boolean;
  config: ITaggerConfig;
  setConfig: Dispatch<ITaggerConfig>;
  excludedFields: string[];
  onFieldsChange: (fields: string[]) => void;
  fields: string[];
  entityName: string;
  extraConfig?: React.ReactNode;
}

const TaggerConfig: React.FC<ITaggerConfigProps> = ({
  show,
  config,
  setConfig,
  excludedFields,
  onFieldsChange,
  fields,
  entityName,
  extraConfig,
}) => {
  const { configuration: stashConfig } = useConfigurationContext();
  const [showExclusionModal, setShowExclusionModal] = useState(false);

  const handleInstanceSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedEndpoint = e.currentTarget.value;
    setConfig({
      ...config,
      selectedEndpoint,
    });
  };

  const stashBoxes = stashConfig?.general.stashBoxes ?? [];

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
              <Form.Group
                controlId="stash-box-endpoint"
                className="align-items-center row no-gutters mt-4"
              >
                <Form.Label className="mr-4">
                  <FormattedMessage id="tagger.config.active_stash-box_instance" />
                </Form.Label>
                <Form.Control
                  as="select"
                  value={config.selectedEndpoint}
                  className="col-md-4 col-6 input-control"
                  disabled={!stashBoxes.length}
                  onChange={handleInstanceSelect}
                >
                  {!stashBoxes.length && (
                    <option>
                      <FormattedMessage id="tagger.config.no_instances_found" />
                    </option>
                  )}
                  {stashConfig?.general.stashBoxes.map((i) => (
                    <option value={i.endpoint} key={i.endpoint}>
                      {i.endpoint}
                    </option>
                  ))}
                </Form.Control>
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
