import React, { Dispatch, useState } from "react";
import { Badge, Button, Card, Collapse, Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { useConfigurationContext } from "src/hooks/Config";

import { ITaggerConfig } from "../constants";
import StudioFieldSelector from "./StudioFieldSelector";

interface IConfigProps {
  show: boolean;
  config: ITaggerConfig;
  setConfig: Dispatch<ITaggerConfig>;
}

const Config: React.FC<IConfigProps> = ({ show, config, setConfig }) => {
  const { configuration: stashConfig } = useConfigurationContext();
  const [showExclusionModal, setShowExclusionModal] = useState(false);

  const excludedFields = config.excludedStudioFields ?? [];

  const handleInstanceSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedEndpoint = e.currentTarget.value;
    setConfig({
      ...config,
      selectedEndpoint,
    });
  };

  const stashBoxes = stashConfig?.general.stashBoxes ?? [];

  const handleFieldSelect = (fields: string[]) => {
    setConfig({ ...config, excludedStudioFields: fields });
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
              <Form.Group
                controlId="create-parent"
                className="align-items-center"
              >
                <Form.Check
                  label={
                    <FormattedMessage id="studio_tagger.config.create_parent_label" />
                  }
                  checked={config.createParentStudios}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                    setConfig({
                      ...config,
                      createParentStudios: e.currentTarget.checked,
                    })
                  }
                />
                <Form.Text>
                  <FormattedMessage id="studio_tagger.config.create_parent_desc" />
                </Form.Text>
              </Form.Group>
              <Form.Group controlId="excluded-studio-fields">
                <h6>
                  <FormattedMessage id="studio_tagger.config.excluded_fields" />
                </h6>
                <span>
                  {excludedFields.length > 0 ? (
                    excludedFields.map((f) => (
                      <Badge variant="secondary" className="tag-item" key={f}>
                        <FormattedMessage id={f} />
                      </Badge>
                    ))
                  ) : (
                    <FormattedMessage id="studio_tagger.config.no_fields_are_excluded" />
                  )}
                </span>
                <Form.Text>
                  <FormattedMessage id="studio_tagger.config.these_fields_will_not_be_changed_when_updating_studios" />
                </Form.Text>
                <Button
                  onClick={() => setShowExclusionModal(true)}
                  className="mt-2"
                >
                  <FormattedMessage id="studio_tagger.config.edit_excluded_fields" />
                </Button>
              </Form.Group>
              <Form.Group
                controlId="stash-box-endpoint"
                className="align-items-center row no-gutters mt-4"
              >
                <Form.Label className="mr-4">
                  <FormattedMessage id="studio_tagger.config.active_stash-box_instance" />
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
                      <FormattedMessage id="studio_tagger.config.no_instances_found" />
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
      <StudioFieldSelector
        show={showExclusionModal}
        onSelect={handleFieldSelect}
        excludedFields={excludedFields}
      />
    </>
  );
};

export default Config;
