import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IPhashDistanceValue } from "../../../models/list-filter/types";
import { Criterion } from "../../../models/list-filter/criteria/criterion";
import { CriterionModifier } from "src/core/generated-graphql";
import { useDebouncedState } from "src/hooks/debounce";

interface IPhashFilterProps {
  criterion: Criterion<IPhashDistanceValue>;
  onValueChanged: (value: IPhashDistanceValue) => void;
}

export const PhashFilter: React.FC<IPhashFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const [value, setValue] = useDebouncedState(criterion.value, onValueChanged);

  function valueChanged(event: React.ChangeEvent<HTMLInputElement>) {
    setValue({
      value: event.target.value,
      distance: criterion.value.distance,
    });
  }

  function distanceChanged(event: React.ChangeEvent<HTMLInputElement>) {
    const distanceStr = event.target.value;

    if (distanceStr === "") {
      setValue({
        value: criterion.value.value,
        distance: undefined,
      });
      return;
    }

    let distance = parseInt(event.target.value);
    if (distance < 0 || isNaN(distance)) {
      distance = 0;
    }

    setValue({
      distance,
      value: criterion.value.value,
    });
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={valueChanged}
          value={value ? value.value : ""}
          placeholder={intl.formatMessage({ id: "media_info.phash" })}
        />
      </Form.Group>
      {criterion.modifier !== CriterionModifier.IsNull &&
        criterion.modifier !== CriterionModifier.NotNull && (
          <Form.Group>
            <Form.Control
              className="btn-secondary"
              onChange={distanceChanged}
              type="number"
              value={value ? value.distance : ""}
              placeholder={intl.formatMessage({ id: "distance" })}
            />
          </Form.Group>
        )}
    </div>
  );
};
