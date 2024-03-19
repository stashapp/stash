import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IStashIDValue } from "../../../models/list-filter/types";
import { Criterion } from "../../../models/list-filter/criteria/criterion";
import { CriterionModifier } from "src/core/generated-graphql";
import { useDebouncedState } from "src/hooks/debounce";

interface IStashIDFilterProps {
  criterion: Criterion<IStashIDValue>;
  onValueChanged: (value: IStashIDValue) => void;
}

export const StashIDFilter: React.FC<IStashIDFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const [value, setValue] = useDebouncedState(criterion.value, onValueChanged);

  function onEndpointChanged(event: React.ChangeEvent<HTMLInputElement>) {
    setValue({
      endpoint: event.target.value,
      stashID: criterion.value.stashID,
    });
  }

  function onStashIDChanged(event: React.ChangeEvent<HTMLInputElement>) {
    setValue({
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
