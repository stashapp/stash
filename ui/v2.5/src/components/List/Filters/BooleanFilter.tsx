import cloneDeep from "lodash-es/cloneDeep";
import React from "react";
import { Form } from "react-bootstrap";
import { BooleanCriterion } from "src/models/list-filter/criteria/criterion";
import { FormattedMessage } from "react-intl";

interface IBooleanFilter {
  criterion: BooleanCriterion;
  setCriterion: (c: BooleanCriterion) => void;
}

export const BooleanFilter: React.FC<IBooleanFilter> = ({
  criterion,
  setCriterion,
}) => {
  function onSelect(v: boolean) {
    const c = cloneDeep(criterion);
    if ((v && c.value === "true") || (!v && c.value === "false")) {
      c.value = "";
    } else {
      c.value = v ? "true" : "false";
    }

    setCriterion(c);
  }

  return (
    <div className="boolean-filter">
      <Form.Check
        id={`${criterion.getId()}-true`}
        onChange={() => onSelect(true)}
        checked={criterion.value === "true"}
        type="checkbox"
        label={<FormattedMessage id="true" />}
      />
      <Form.Check
        id={`${criterion.getId()}-false`}
        onChange={() => onSelect(false)}
        checked={criterion.value === "false"}
        type="checkbox"
        label={<FormattedMessage id="false" />}
      />
    </div>
  );
};
