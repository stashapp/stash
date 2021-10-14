import React, { Dispatch, useState } from "react";
import { Badge, Button, Card, Collapse, Form } from "react-bootstrap";
import { ConfigurationContext } from "src/hooks/Config";

import { TextUtils } from "src/utils";
import { ITaggerConfig, PERFORMER_FIELDS } from "../constants";
import PerformerFieldSelector from "../PerformerFieldSelector";

interface IConfigProps {
  show: boolean;
  config: ITaggerConfig;
  setConfig: Dispatch<ITaggerConfig>;
}

const Config: React.FC<IConfigProps> = ({ show, config, setConfig }) => {
  const { configuration: stashConfig } = React.useContext(ConfigurationContext);
  const [showExclusionModal, setShowExclusionModal] = useState(false);

  const excludedFields = config.excludedPerformerFields ?? [];

  const handleInstanceSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedEndpoint = e.currentTarget.value;
    setConfig({
      ...config,
      selectedEndpoint,
    });
  };

  const stashBoxes = stashConfig?.general.stashBoxes ?? [];

  const handleFieldSelect = (fields: string[]) => {
    setConfig({ ...config, excludedPerformerFields: fields });
    setShowExclusionModal(false);
  };

  return (
    <>
      <Collapse in={show}>
        <Card>
          <div className="row">
            <h4 className="col-12">Configuration</h4>
            <hr className="w-100" />
            <div className="col-md-6">
              <Form.Group controlId="excluded-performer-fields">
                <h6>Excluded fields:</h6>
                <span>
                  {excludedFields.length > 0
                    ? excludedFields.map((f) => (
                        <Badge variant="secondary" className="tag-item">
                          {TextUtils.capitalize(f)}
                        </Badge>
                      ))
                    : "No fields are excluded"}
                </span>
                <Form.Text>
                  These fields will not be changed when updating performers.
                </Form.Text>
                <Button
                  onClick={() => setShowExclusionModal(true)}
                  className="mt-2"
                >
                  Edit Excluded Fields
                </Button>
              </Form.Group>
              <Form.Group
                controlId="stash-box-endpoint"
                className="align-items-center row no-gutters mt-4"
              >
                <Form.Label className="mr-4">
                  Active stash-box instance:
                </Form.Label>
                <Form.Control
                  as="select"
                  value={config.selectedEndpoint}
                  className="col-md-4 col-6 input-control"
                  disabled={!stashBoxes.length}
                  onChange={handleInstanceSelect}
                >
                  {!stashBoxes.length && <option>No instances found</option>}
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
      <PerformerFieldSelector
        fields={PERFORMER_FIELDS}
        show={showExclusionModal}
        onSelect={handleFieldSelect}
        excludedFields={excludedFields}
      />
    </>
  );
};

export default Config;
