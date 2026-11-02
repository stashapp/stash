import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IStashIDValue } from "../../../models/list-filter/types";
import { ModifierCriterion } from "../../../models/list-filter/criteria/criterion";
import { CriterionModifier } from "src/core/generated-graphql";
import { useConfigurationContext } from "src/hooks/Config";
import { stashboxDisplayName } from "src/utils/stashbox";

interface IStashIDFilterProps {
  criterion: ModifierCriterion<IStashIDValue>;
  onValueChanged: (value: IStashIDValue) => void;
}

export const StashIDFilter: React.FC<IStashIDFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();
  const { value } = criterion;
  const stashBoxes = configuration.general.stashBoxes;
  const selectedEndpoint = value?.endpoint ?? "";

  function onEndpointChanged(event: React.ChangeEvent<HTMLSelectElement>) {
    onValueChanged({
      endpoint: event.target.value,
      stashID: criterion.value.stashID,
    });
  }

  function onStashIDChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      stashID: event.target.value,
      endpoint: criterion.value.endpoint,
    });
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          as="select"
          className="btn-secondary"
          onChange={onEndpointChanged}
          value={selectedEndpoint}
          aria-label={intl.formatMessage({ id: "stash_id_endpoint" })}
        >
          <option value="">
            {intl.formatMessage({ id: "stash_id_endpoint_any" })}
          </option>
          {stashBoxes.map((stashBox, index) => (
            <option value={stashBox.endpoint} key={stashBox.endpoint}>
              {stashboxDisplayName(stashBox.name, index)}
            </option>
          ))}
          {selectedEndpoint &&
            !stashBoxes.some(
              (stashBox) => stashBox.endpoint === selectedEndpoint
            ) && <option value={selectedEndpoint}>{selectedEndpoint}</option>}
        </Form.Control>
      </Form.Group>
      {criterion.modifier !== CriterionModifier.IsNull &&
        criterion.modifier !== CriterionModifier.NotNull && (
          <Form.Group>
            <Form.Control
              className="btn-secondary"
              onChange={onStashIDChanged}
              value={value ? value.stashID : ""}
              placeholder={intl.formatMessage({ id: "stash_id" })}
            />
          </Form.Group>
        )}
    </div>
  );
};
