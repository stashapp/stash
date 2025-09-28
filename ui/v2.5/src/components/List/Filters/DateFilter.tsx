import React, { useCallback, useMemo } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { IDateValue } from "../../../models/list-filter/types";
import {
  ModifierCriterion,
  CriterionOption,
} from "../../../models/list-filter/criteria/criterion";
import { DateInput } from "src/components/Shared/DateInput";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { SelectedItem } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import { ModifierControls } from "./StringFilter";
import { ModifierSelectorButtons } from "../ModifierSelect";

interface IDateFilterProps {
  criterion: ModifierCriterion<IDateValue>;
  onValueChanged: (value: IDateValue) => void;
}

// Hook for date-based sidebar filters
export function useDateCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ModifierCriterion<IDateValue>;

    const newCriterion = filter.makeCriterion(
      option.type
    ) as ModifierCriterion<IDateValue>;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ModifierCriterion<IDateValue>) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      // if (c.isValid()) newCriteria.push(c);
      newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  const onValueChanged = useCallback(
    (value: IDateValue) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.value = value;
      setCriterion(newCriterion);
    },
    [criterion, filter, setFilter, option.type, defaultModifier]
  );

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  return {
    criterion,
    setCriterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  };
}

// Component for selected date items display
interface IDateSelectedItemsProps {
  criterion: ModifierCriterion<IDateValue> | null;
  defaultModifier: CriterionModifier;
  onChangedModifierSelect: (m: CriterionModifier) => void;
  onClear: () => void;
}

export const DateSelectedItems: React.FC<IDateSelectedItemsProps> = ({
  criterion,
  defaultModifier,
  onChangedModifierSelect,
  onClear,
}) => {
  if (criterion?.value.value === "") {
    return null;
  }

  const intl = useIntl();

  const getValueLabel = () => {
    if (!criterion?.value) return null;

    const { value, value2 } = criterion.value;

    switch (criterion.modifier) {
      case CriterionModifier.Equals:
        return value;
      case CriterionModifier.NotEquals:
        return `≠ ${value}`;
      case CriterionModifier.GreaterThan:
        return `> ${value}`;
      case CriterionModifier.LessThan:
        return `< ${value}`;
      case CriterionModifier.Between:
        return `${value} - ${value2}`;
      case CriterionModifier.NotBetween:
        return `not ${value} - ${value2}`;
      default:
        return value;
    }
  };

  const valueLabel = getValueLabel();

  return (
    <ul className="selected-list">
      {criterion?.modifier != defaultModifier && criterion?.modifier ? (
        <SelectedItem
          className="modifier-object"
          label={ModifierCriterion.getModifierLabel(intl, criterion.modifier)}
          onClick={() => onChangedModifierSelect(defaultModifier)}
        />
      ) : null}
      {valueLabel ? (
        <SelectedItem label={valueLabel} onClick={onClear} />
      ) : null}
    </ul>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

export const DateFilter: React.FC<IDateFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  const { value } = criterion;

  function onChanged(newValue: string, property: "value" | "value2") {
    const valueCopy = { ...value };

    valueCopy[property] = newValue;
    onValueChanged(valueCopy);
  }

  let equalsControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals
  ) {
    equalsControl = (
      <Form.Group>
        <DateInput
          value={value?.value ?? ""}
          onValueChange={(v) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.value" })}
        />
      </Form.Group>
    );
  }

  let lowerControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.GreaterThan ||
    criterion.modifier === CriterionModifier.Between ||
    criterion.modifier === CriterionModifier.NotBetween
  ) {
    lowerControl = (
      <Form.Group>
        <DateInput
          value={value?.value ?? ""}
          onValueChange={(v) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.greater_than" })}
        />
      </Form.Group>
    );
  }

  let upperControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.LessThan ||
    criterion.modifier === CriterionModifier.Between ||
    criterion.modifier === CriterionModifier.NotBetween
  ) {
    upperControl = (
      <Form.Group>
        <DateInput
          value={
            (criterion.modifier === CriterionModifier.LessThan
              ? value?.value
              : value?.value2) ?? ""
          }
          onValueChange={(v) =>
            onChanged(
              v,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          placeholder={intl.formatMessage({ id: "criterion.less_than" })}
        />
      </Form.Group>
    );
  }

  return (
    <>
      {equalsControl}
      {lowerControl}
      {upperControl}
    </>
  );
};

export const SidebarDateFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {
  const intl = useIntl();

  const {
    criterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  } = useDateCriterion(option, filter, setFilter);

  const modifierSelector = useMemo(() => {
    return (
      <ModifierSelectorButtons
        options={modifierOptions}
        value={criterion.modifier}
        onChanged={onChangedModifierSelect}
      />
    );
  }, [
    modifierOptions,
    onChangedModifierSelect,
    criterion.modifier,
  ]);

  const valueControl = useMemo(() => {
    return (
      <DateFilter criterion={criterion} onValueChanged={onValueChanged} />
    );
  }, [criterion]);


  const onClear = useCallback(() => {
    setFilter(filter.removeCriterion(option.type));
  }, [filter, setFilter, option.type]);

  const onChanged = useCallback(
    (newValue: string, property: "value" | "value2") => {
      const currentValue = criterion?.value || { value: "", value2: "" };
      const valueCopy = { ...currentValue };

      valueCopy[property] = newValue;
      onValueChanged(valueCopy);
    },
    [criterion?.value, onValueChanged]
  );

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      outsideCollapse={
        <DateSelectedItems
          criterion={criterion}
          defaultModifier={defaultModifier}
          onChangedModifierSelect={onChangedModifierSelect}
          onClear={onClear}
        />
      }
    >
      <div className="date-filter">
        <div className="filter-group">
          {modifierSelector}
          {valueControl}
        </div>
      </div>
    </SidebarSection>
  );
};
