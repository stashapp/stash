import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IStashIDValue } from "../../../models/list-filter/types";
import { ModifierCriterion } from "../../../models/list-filter/criteria/criterion";
import { CriterionModifier } from "src/core/generated-graphql";

interface IStashIDFilterProps {
  criterion: ModifierCriterion<IStashIDValue>;
  onValueChanged: (value: IStashIDValue) => void;
}

export const StashIDFilter: React.FC<IStashIDFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const { value } = criterion;

  function onEndpointChanged(event: React.ChangeEvent<HTMLInputElement>) {
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
          className="btn-secondary"
          onChange={onEndpointChanged}
          value={value ? value.endpoint : ""}
          placeholder={intl.formatMessage({ id: "stash_id_endpoint" })}
        />
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
