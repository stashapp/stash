import React, { useCallback, useMemo } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { INumberValue } from "../../../models/list-filter/types";
import {
  NumberCriterion,
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { NumberField } from "src/utils/form";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { SelectedItem } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import { ModifierControls } from "./StringFilter";
import { ModifierSelectorButtons } from "../ModifierSelect";

interface IDurationFilterProps {
  criterion: NumberCriterion;
  onValueChanged: (value: INumberValue) => void;
}

// Hook for number-based sidebar filters
export function useNumberCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as NumberCriterion;

    const newCriterion = filter.makeCriterion(option.type) as NumberCriterion;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: NumberCriterion) => {
      const newFilter = cloneDeep(filter);

      if (false && !c.isValid()) {
        // TODO: need to remove false condition once we have a way to validate the criterion from sidebar
        return;
        // remove from the filter if present
        const newCriteria = filter.criteria.filter((cc) => {
          return cc.criterionOption.type !== c.criterionOption.type;
        });
  
        newFilter.criteria = newCriteria;
      } else {
        let found = false;
  
        const newCriteria = filter.criteria.map((cc) => {
          if (cc.criterionOption.type === c.criterionOption.type) {
            found = true;
            return c;
          }
  
          return cc;
        });
  
        if (!found) {
          newCriteria.push(c);
        }
  
        newFilter.criteria = newCriteria;
      }
      setFilter(newFilter);
    },
    [option.type, setFilter, filter]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  function onValueChanged(value: INumberValue) {
    const newCriterion = cloneDeep(criterion);
    newCriterion.value = value;
    setCriterion(newCriterion);
  }

  return {
    criterion,
    setCriterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  };
}

// Component for selected number items display
interface INumberSelectedItemsProps {
  criterion: NumberCriterion | null;
  defaultModifier: CriterionModifier;
  onChangedModifierSelect: (m: CriterionModifier) => void;
  onClear: () => void;
}

export const NumberSelectedItems: React.FC<INumberSelectedItemsProps> = ({
  criterion,
  defaultModifier,
  onChangedModifierSelect,
  onClear,
}) => {
  if (criterion?.value.value === undefined) {
    return null;
  }
  const intl = useIntl();

  const getValueLabel = () => {
    if (!criterion?.value) return null;

    const { value, value2 } = criterion.value;

    switch (criterion.modifier) {
      case CriterionModifier.Equals:
        return value?.toString();
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
        return value?.toString();
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

export const NumberFilter: React.FC<IDurationFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  const { value } = criterion;

  function onChanged(
    event: React.ChangeEvent<HTMLInputElement>,
    property: "value" | "value2"
  ) {
    const numericValue = parseInt(event.target.value, 10);
    const valueCopy = { ...value };

    valueCopy[property] = !Number.isNaN(numericValue) ? numericValue : 0;
    onValueChanged(valueCopy);
  }

  let equalsControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals
  ) {
    equalsControl = (
      <Form.Group>
        <NumberField
          className="btn-secondary"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
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
        <NumberField
          className="btn-secondary"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
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
        <NumberField
          className="btn-secondary"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(
              e,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          value={
            (criterion.modifier === CriterionModifier.LessThan
              ? value?.value
              : value?.value2) ?? ""
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

export const SidebarNumberFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {

  const {
    criterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  } = useNumberCriterion(option, filter, setFilter);

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
      <NumberFilter criterion={criterion} onValueChanged={onValueChanged} />
    );
  }, [criterion]);


  const onClear = useCallback(() => {
    setFilter(filter.removeCriterion(option.type));
  }, [filter, setFilter, option.type]);
    
  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      outsideCollapse={
        <NumberSelectedItems
          criterion={criterion}
          defaultModifier={defaultModifier}
          onChangedModifierSelect={onChangedModifierSelect}
          onClear={onClear}
        />
      }
    >
      <div className="number-filter">
        <div className="filter-group">
          {modifierSelector}
          {valueControl}
        </div>
      </div>
    </SidebarSection>
  );
};
