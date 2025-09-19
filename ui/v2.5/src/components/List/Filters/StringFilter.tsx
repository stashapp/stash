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

  // const criteria = filter.criteriaFor(
  //   option.type
  // ) as ModifierCriterion<string>[];
  // const criterion = criteria.length > 0 ? criteria[0] : null;
  const {criterion, setCriterion} = useStringCriterion(option, filter, setFilter);
  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption.defaultModifier;
  const modifierOptions = modifierCriterionOption.modifierOptions;

  function onValueChange(value: string) {
    if (!value.trim()) {
      // Remove criterion if empty
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    // const newCriterion = criterion ? criterion.clone() : option.makeCriterion();
    const newCriterion = cloneDeep(criterion);
    newCriterion.modifier = criterion?.modifier ? criterion.modifier : defaultModifier;
    newCriterion.value = value;
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      console.log("onChangedModifierSelect", m);
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      outsideCollapse={
        <ul className="selected-list">
          {criterion?.modifier != defaultModifier ? (
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
      }
    >
      <div className="string-filter">
        <div className="filter-group">
        <ModifierSelectorButtons
            options={modifierOptions}
            value={criterion?.modifier ? criterion.modifier : defaultModifier}
            onChanged={onChangedModifierSelect}
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
