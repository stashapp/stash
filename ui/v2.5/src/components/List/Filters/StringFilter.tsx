import React, { useCallback, useMemo } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { SelectedItem, SelectedList } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import { ModifierSelectorButtons } from "../ModifierSelect";

interface IStringFilterProps {
  criterion: ModifierCriterion<string>;
  onValueChanged: (value: string) => void;
  placeholder?: string;
}

export function useStringCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ModifierCriterion<string>;

    const newCriterion = filter.makeCriterion(
      option.type
    ) as ModifierCriterion<string>;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ModifierCriterion<string>) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  return { criterion, setCriterion };
}

// Shared hook for modifier-based sidebar filters
export function useModifierCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const { criterion, setCriterion } = useStringCriterion(option, filter, setFilter);
  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  const onValueChange = useCallback((value: string) => {
    if (!value.trim()) {
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    const newCriterion = cloneDeep(criterion);
    newCriterion.modifier = criterion?.modifier ? criterion.modifier : defaultModifier;
    newCriterion.value = value;
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }, [criterion, setCriterion, filter, setFilter, option.type, defaultModifier]);

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
    onValueChange,
    onChangedModifierSelect
  };
}

// Shared component for selected items display
interface ISelectedItemsProps {
  criterion: ModifierCriterion<string> | null;
  defaultModifier: CriterionModifier;
  onChangedModifierSelect: (m: CriterionModifier) => void;
  onValueChange: (value: string) => void;
}

export const SelectedItems: React.FC<ISelectedItemsProps> = ({
  criterion,
  defaultModifier,
  onChangedModifierSelect,
  onValueChange
}) => {
  const intl = useIntl();

  return (
    <ul className="selected-list">
      {criterion?.modifier != defaultModifier && criterion?.modifier ? (
        <SelectedItem
          className="modifier-object"
          label={ModifierCriterion.getModifierLabel(intl, criterion.modifier)}
          onClick={() => onChangedModifierSelect(defaultModifier)}
        />
      ) : null}
      {criterion?.value ? (
        <SelectedItem
          label={criterion.value}
          onClick={() => onValueChange("")}
        />
      ) : null}
    </ul>
  );
};

// Shared component for modifier controls
interface IModifierControlsProps {
  modifierOptions: CriterionModifier[] | undefined;
  currentModifier: CriterionModifier;
  onChangedModifierSelect: (m: CriterionModifier) => void;
}

export const ModifierControls: React.FC<IModifierControlsProps> = ({
  modifierOptions,
  currentModifier,
  onChangedModifierSelect
}) => (
  <ModifierSelectorButtons
    options={modifierOptions}
    value={currentModifier}
    onChanged={onChangedModifierSelect}
  />
);

export const StringFilter: React.FC<IStringFilterProps> = ({
  criterion,
  onValueChanged,
  placeholder,
}) => {
  const intl = useIntl();

  function onValueChange(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged(event.target.value);
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={onValueChange}
          value={criterion.value}
          placeholder={placeholder || intl.formatMessage({ id: "search" })}
        />
      </Form.Group>
    </div>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  placeholder?: string;
  modifier?: CriterionModifier;
}

export const SidebarStringFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  placeholder,
}) => {
  const intl = useIntl();

  const {
    criterion,
    defaultModifier,
    modifierOptions,
    onValueChange,
    onChangedModifierSelect
  } = useModifierCriterion(option, filter, setFilter);

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      outsideCollapse={
        <SelectedItems
          criterion={criterion}
          defaultModifier={defaultModifier}
          onChangedModifierSelect={onChangedModifierSelect}
          onValueChange={onValueChange}
        />
      }
    >
      <div className="string-filter">
        <div className="filter-group">
          <ModifierControls
            modifierOptions={modifierOptions}
            currentModifier={criterion?.modifier || defaultModifier}
            onChangedModifierSelect={onChangedModifierSelect}
          />
          <Form.Control
            className="btn-secondary"
            onChange={(v) => onValueChange(v.target.value)}
            value={criterion?.value || ""}
            placeholder={placeholder || intl.formatMessage({ id: "search" })}
          />
        </div>
      </div>
    </SidebarSection>
  );
};

// Convenience exports for specific string filters
export const SidebarTattoosFilter: React.FC<ISidebarFilter> = (props) => (
  <SidebarStringFilter
    {...props}
    placeholder={props.placeholder || "tattoos"}
  />
);
