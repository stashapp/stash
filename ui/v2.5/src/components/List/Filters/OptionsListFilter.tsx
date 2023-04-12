import cloneDeep from "lodash-es/cloneDeep";
import React from "react";
import { Form } from "react-bootstrap";
import {
  CriterionValue,
  Criterion,
} from "src/models/list-filter/criteria/criterion";

interface IOptionsListFilter {
  criterion: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
}

export const OptionsListFilter: React.FC<IOptionsListFilter> = ({
  criterion,
  setCriterion,
}) => {
  function onSelect(v: string) {
    const c = cloneDeep(criterion);
    if (c.value === v) {
      c.value = "";
    } else {
      c.value = v;
    }

    setCriterion(c);
  }

  const { options } = criterion.criterionOption;

  return (
    <div className="option-list-filter">
      {options?.map((o) => (
        <Form.Check
          id={`${criterion.getId()}-${o.toString()}`}
          key={o.toString()}
          onChange={() => onSelect(o.toString())}
          checked={criterion.value === o.toString()}
          type="checkbox"
          label={o.toString()}
        />
      ))}
    </div>
  );
};
