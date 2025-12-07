import cloneDeep from "lodash-es/cloneDeep";
import React, { useMemo } from "react";
import { Form } from "react-bootstrap";
import {
  BooleanCriterion,
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { FormattedMessage, useIntl } from "react-intl";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";

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
        type="radio"
        label={<FormattedMessage id="true" />}
      />
      <Form.Check
        id={`${criterion.getId()}-false`}
        onChange={() => onSelect(false)}
        checked={criterion.value === "false"}
        type="radio"
        label={<FormattedMessage id="false" />}
      />
    </div>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarBooleanFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();

  const trueLabel = intl.formatMessage({
    id: "true",
  });
  const falseLabel = intl.formatMessage({
    id: "false",
  });

  const trueOption = useMemo(
    () => ({
      id: "true",
      label: trueLabel,
    }),
    [trueLabel]
  );

  const falseOption = useMemo(
    () => ({
      id: "false",
      label: falseLabel,
    }),
    [falseLabel]
  );

  const criteria = filter.criteriaFor(option.type) as BooleanCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (criterion.value === "true") {
      return [trueOption];
    } else if (criterion.value === "false") {
      return [falseOption];
    }

    return [];
  }, [trueOption, falseOption, criterion]);

  const options: Option[] = useMemo(() => {
    return [trueOption, falseOption].filter((o) => !selected.includes(o));
  }, [selected, trueOption, falseOption]);

  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    newCriterion.value = item.id;

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect() {
    setFilter(filter.removeCriterion(option.type));
  }

  return (
    <>
      <SidebarListFilter
        title={title}
        candidates={options}
        onSelect={onSelect}
        onUnselect={onUnselect}
        selected={selected}
        singleValue
        sectionID={sectionID}
      />
    </>
  );
};
