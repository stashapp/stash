import React, { useMemo } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IPhashDistanceValue } from "../../../models/list-filter/types";
import { CriterionOption, ModifierCriterion } from "../../../models/list-filter/criteria/criterion";
import { CriterionModifier } from "src/core/generated-graphql";
import { NumberField } from "src/utils/form";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { PhashCriterion } from "src/models/list-filter/criteria/phash";

interface IPhashFilterProps {
  criterion: ModifierCriterion<IPhashDistanceValue>;
  onValueChanged: (value: IPhashDistanceValue) => void;
}

export const PhashFilter: React.FC<IPhashFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const { value } = criterion;

  function valueChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      value: event.target.value,
      distance: criterion.value.distance,
    });
  }

  function distanceChanged(event: React.ChangeEvent<HTMLInputElement>) {
    let distance = parseInt(event.target.value);
    if (distance < 0 || isNaN(distance)) {
      distance = 0;
    }

    onValueChanged({
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
            <NumberField
              className="btn-secondary"
              onChange={distanceChanged}
              value={value ? value.distance : ""}
              placeholder={intl.formatMessage({ id: "distance" })}
            />
          </Form.Group>
        )}
    </div>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

const any = "any";
const none = "none";

export const SidebarPhashFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {
  const intl = useIntl();

  const anyLabel = `(${intl.formatMessage({
    id: "criterion_modifier_values.any",
  })})`;
  const noneLabel = `(${intl.formatMessage({
    id: "criterion_modifier_values.none",
  })})`;

  const anyOption = useMemo(
    () => ({
      id: "any",
      label: anyLabel,
      className: "modifier-object",
    }),
    [anyLabel]
  );

  const noneOption = useMemo(
    () => ({
      id: "none",
      label: noneLabel,
      className: "modifier-object",
    }),
    [noneLabel]
  );

  const options: Option[] = useMemo(() => {
    return [anyOption, noneOption];
  }, [anyOption, noneOption])

  const criteria = filter.criteriaFor(option.type) as PhashCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (criterion.modifier === CriterionModifier.NotNull) {
      return [anyOption];
    } else if (criterion.modifier === CriterionModifier.IsNull) {
      return [noneOption];
    }

    return [];
  }, [anyOption, noneOption, criterion]);


  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (item.id === any) {
      newCriterion.modifier = CriterionModifier.NotNull;
      // newCriterion.value
    } else if (item.id === none) {
      newCriterion.modifier = CriterionModifier.IsNull;
    }

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
      />
    </>
  );
};