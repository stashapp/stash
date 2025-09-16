import React, { useMemo } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IStashIDValue } from "../../../models/list-filter/types";
import { ModifierCriterion, CriterionOption } from "../../../models/list-filter/criteria/criterion";
import { CriterionModifier } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { StashIDCriterion } from "src/models/list-filter/criteria/stash-ids";
import { Option, SidebarListFilter } from "./SidebarListFilter";

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

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

const any = "any";
const none = "none";

export const SidebarStashIDFilter: React.FC<ISidebarFilter> = ({
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
  }, [anyOption, noneOption]);

  const criteria = filter.criteriaFor(option.type) as StashIDCriterion[];
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
      newCriterion.value = { endpoint: "", stashID: "" };
    } else if (item.id === none) {
      newCriterion.modifier = CriterionModifier.IsNull;
      newCriterion.value = { endpoint: "", stashID: "" };
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
        singleValue={true}
      />
    </>
  );
};
